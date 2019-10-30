package main

import(
	"encoding/json"
	"os"
	"io/ioutil"
	"os/exec"
	"fmt"
	"errors"
	"reflect"
	"regexp"
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
      ok = validateValue(opt, val.([]interface{})[elem], true)
      if !ok { break }
    }

  } else if val == nil && !opt.Required {
    //Optional value not provided
    ok = true

  } else {
    // Verify the value/type of this field
    switch(opt.Value.Type){
      case "bool":
        _, ok = val.(bool)

      case "integer":
        numf, ok2 := val.(float64)
        num := int(numf)
        if(ok2){ ok2 = ( numf == float64(num) ) } //Verify a whole number was provided
        if(ok2 && opt.Value.Min != nil){ ok2 = (num >= int(*opt.Value.Min)) }
        if(ok2 && opt.Value.Max != nil){ ok2 = (num <= int(*opt.Value.Max)) }
        ok = ok2;

      case "float":
        num, ok2 := val.(float64)
        if(ok && opt.Value.Min != nil){ ok2 = (num >= *opt.Value.Min) }
        if(ok && opt.Value.Max != nil){ ok2 = (num <= *opt.Value.Max) }
        ok = ok2;

      case "string":
        text, ok2 := val.(string)
        if(ok2 && opt.Value.Regex != ""){ ok2, _ = regexp.MatchString(opt.Value.Regex, text) }
        ok = ok2;

      case "select":
        for _, valid := range(opt.Value.Select) {
          if val == valid {
            //got an exact match
            ok = true 
            break
          }
        }

      default:
        fmt.Println("Unknown type of option: ", opt.Value.Type, "Provided Value: ", val)
    }
  }
  if !ok && !opt.Required {
    // Option check failed, but option also not required.
    fmt.Println("Optional argument not provided:", opt.Field)
    ok = true
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
      if opt.Value.Default != nil {
        alert.Settings[opt.Field] = opt.Value.Default //insert the default value

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
  // Ensure the plugin is flagged as executable (pre-1.0 bug)
  finfo, _ := os.Stat(execpath)
  fmode := finfo.Mode();
  fmt.Println("Got fmode:", fmod, execpath)
  if (fmode & 0111) != 0 {
    os.Chmod(execpath, 0755);
  }
  // Now save the settings file for input to the executable
  tmp, _ := json.Marshal(alert)
  tmpFile, err := ioutil.TempFile(os.TempDir(), ".*.json")
  if err != nil { fmt.Println("Error creating tmp file") ; return }
  _, err = tmpFile.Write(tmp)
  tmpFile.Close()
  if err != nil { os.Remove(tmpFile.Name()) ; fmt.Println("Error writing temporary file:", tmpFile.Name()) ; return }
  // Now call the command with the input file path
  //fmt.Println("Using tmp file: ", tmpFile.Name())
  cmd := exec.Command(execpath, tmpFile.Name())
  info, err := cmd.Output()
  if err != nil { fmt.Println( string(info)) }
  // Now remove the temporary file 
  os.Remove(tmpFile.Name())
}
