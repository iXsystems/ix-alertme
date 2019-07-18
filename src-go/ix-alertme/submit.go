package main

import(
	"encoding/json"
	"os"
	"io/ioutil"
	"os/exec"
	"fmt"
	"errors"
)

type AlertText struct {
	Html string 		`json:"html"`
	PlainText string	`json:"plain"`
}

type AlertAPI struct {
	Text AlertText		`json:"text"`
	Settings map[string]interface{}	`json:"settings"`
}

func validateSettings(plugin PluginFullManifest, alert AlertAPI) error {
  for index := range(plugin.API) {
    opt := plugin.API[index]
    if val, ok := alert.Settings[opt.Field] ; ok {
      // Verify the value/type of this field
      switch(opt.Type){
        case "":

        case "list":

	default:
	  fmt.Println("Got Value: ", val)
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

func sendAlerts(settingsFile string, textFile string) error {
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
          submitAlert(manifest, alert)
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
