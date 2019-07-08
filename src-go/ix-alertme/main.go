package main

import (
	"os"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"encoding/json"
)

var Config Configuration

func CheckForUpdates(){
  fmt.Println("Check For Updates")
  updates := pluginUpdates()
  tmp, _ := json.Marshal(updates)
  fmt.Println( string(tmp) )
}

func ListRemotePlugins(){
  fmt.Println("List Remote Plugins")
  updates := availablePlugins()
  tmp, _ := json.Marshal(updates)
  fmt.Println( string(tmp) )
}

func ListLocalPlugins(){
  fmt.Println("List Local Plugins")
  updates := installedPlugins()
  tmp, _ := json.Marshal(updates)
  fmt.Println( string(tmp) )
}

// Define all the CLI input flags and subcommands
var (
  app = kingpin.New("ix-alertme", "Alert Notification Plugin System")
  Configfile = *app.Flag("config", "Custom config file").Short('c').Default("/usr/local/etc/ix-alertme.json").String()

  plugins 		= app.Command("plugins", "Manage Plugins")
  pluginsSearch 	= plugins.Command("search", "Search for available plugins")
  pluginsList 		= plugins.Command("list", "List all installed plugins")
  pluginsScan		= plugins.Command("scan", "Scan for updates to installed plugins")
)

func main() {
  app.Version("0.1")
  switch kingpin.MustParse(app.Parse(os.Args[1:])) {
    case pluginsSearch.FullCommand():
        ListRemotePlugins()
    case pluginsList.FullCommand():    
        ListLocalPlugins()
    case pluginsScan.FullCommand():
        CheckForUpdates()
    default:
      app.Fatalf("%s","Unknown subcommand")
  }
}
