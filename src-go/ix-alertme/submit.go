package main

import(
	"encoding/json"
	"os"
	"io/ioutil"
	"os/exec"
	"fmt"
)

type AlertText struct {
	Html string 		`json:"html"`
	PlainText string	`json:"plain"`
}

type AlertAPI struct {
	Text AlertText		`json:"text"`
	Settings map[string]interface{}	`json:"settings"`
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
        submitAlert(manifest, alert) //This will run in a parallel thead, so many of these can run at a time
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
