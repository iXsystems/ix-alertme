package main

import(
	"encoding/json"
	"os"
	"io/ioutil"
	"os/exec"
	"fmt"
	"errors"
	"reflect"
)

type AlertText struct {
	Html string 		`json:"html"`
	PlainText string	`json:"plain"`
}

type AlertAPI struct {
	Text AlertText		`json:"text"`
	Settings map[string]interface{}	`json:"settings"`
}

func validateValue(opt SetOpts, val interface{}, inslice bool) bool {
  vtyp := reflect.ValueOf(val)
  ok := false
  if !inslice && opt.IsArray && vtyp.Kind()==reflect.Slice {
    //Need to check each value in this slice
    for elem := range( val.([]interface{}) ) {
      ok := validateValue(opt, val.([]interface{})[elem], true)
      if !ok { break }
    }

  } else {
    // Verify the value/type of this field
    needtyp := ""
    if reflect.ValueOf(opt.Type).Kind() == reflect.Slice {
      needtyp, ok = opt.Type.([]interface{})[0].(string)
    } else {
      needtyp, ok = opt.Type.(string)
    }
    if !ok { return false } //could not read type of field
    switch(needtyp){
      case "bool":
        _, ok = val.(bool)
      case "integer":
        _, ok = val.(int)
      case "float":
        _, ok = val.(float32)
      case "":
        _, ok = val.(string)
      case "regex":
        _, ok = val.(string)
      case "select":
        list, _ := opt.Type.([]interface{})
        for index := range(list) {
          if index <1 {continue } //skip the initial type field (already used)
          if val == list[index] {
            //got an exact match
            ok = true 
            break
          } else if reflect.ValueOf(list[index]).Kind() == reflect.Slice {
            // Got the expanded type where there is a description of the option in the second field
            if val == list[index].([]interface{})[0] {
              ok = true 
              break
            }
          }
        }

      default:
        fmt.Println("Unknown type of option: ", needtyp, "Provided Value: ", val)
    }
  }
  return ok
}

func validateSettings(plugin PluginFullManifest, alert AlertAPI) error {
  for index := range(plugin.API) {
    opt := plugin.API[index]
    if val, ok := alert.Settings[opt.Field] ; ok {
      if ! validateValue(opt, val, false){
        fmt.Println("["+plugin.Name+"] Invalid API Setting: ", opt.Field)
        return errors.New("Invalid API")
      }
    } else {
      //Field missing - see if there is a default value and add that in
      if opt.Default != nil {
        alert.Settings[opt.Field] = opt.Default //insert the default value

      } else if opt.Required {
        fmt.Println("["+plugin.Name+"] Missing API Setting: ", opt.Field)
        return errors.New("Missing API")
      }
    }
  }
  return nil
}

func sendAlerts(settingsFile string, textFile string, checkonly bool) error {
  // Load the local settings file
  settings := make(map[string]interface{})
  file, err := os.Open(settingsFile)
  defer file.Close()
  if err == nil {
    tmp, err := ioutil.ReadAll(file)
    if err == nil { json.Unmarshal(tmp, &settings) }
  }
  if err != nil { return err }
  // Load the local text file
  var text AlertText
  file2, err := os.Open(textFile)
  defer file2.Close()
  if err == nil {
    tmp, err := ioutil.ReadAll(file2)
    if err == nil { json.Unmarshal(tmp, &text) }
  }
  if err != nil { return err }
  // Now load the list of installed plugins
  installed, err := installedPlugins()
  if err != nil { return err }
  // Now go through the list and send out alerts for any plugins that are setup
  for plugin, manifest := range(installed) {
    if set, ok := settings[plugin] ; ok {
      // Got an installed plugin that has current settings
      var alert AlertAPI;
        alert.Text = text
        alert.Settings, ok = set.(map[string]interface{})
      if ok {
        if err = validateSettings(manifest, alert) ; err == nil {
          if !checkonly {
            fmt.Println("Sending alert via plugin: ", plugin)
            submitAlert(manifest, alert)
          } else {
            fmt.Println("Would send alert via plugin: ", plugin)
          }
        } else {
          fmt.Println("[Skipped] Invalid Plugin Settings: ", plugin)
        }
      }
    }
  }
  return nil
}

func submitAlert(plugin PluginFullManifest, alert AlertAPI){
  // First determine the path to the executable within the plugin dir
  execpath := Config.InstallDir+"/"+plugin.Name+"/"+plugin.Exec
  // Now save the settings file for input to the executable
  tmp, _ := json.Marshal(alert)
  tmpFile, err := ioutil.TempFile(os.TempDir(), ".*.json")
  if err != nil { return }
  _, err = tmpFile.Write(tmp)
  tmpFile.Close()
  if err != nil { os.Remove(tmpFile.Name()) ; return }
  // Now call the command with the input file path
  //fmt.Println("Using tmp file: ", tmpFile.Name())
  cmd := exec.Command(execpath, tmpFile.Name())
  info, err := cmd.Output()
  if err != nil { fmt.Println( string(info)) }
  // Now remove the temporary file 
  os.Remove(tmpFile.Name())
}
