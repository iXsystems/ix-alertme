package main

import (
	"os"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"encoding/json"
)

var Config Configuration

// Internal simplification functions
func PrintDebug(error string){
  fmt.Fprintln(os.Stderr, "[Debug] "+error)
}
func PrintError(error string){
  fmt.Fprintln(os.Stderr, "[Error] "+error)
}

func CheckForUpdates(){
  //fmt.Println("Check For Updates")
  updates := pluginUpdates()
  tmp, _ := json.Marshal(updates)
  fmt.Println( string(tmp) )
}

func ListRemotePlugins(){
  if *pluginsSearchName == "" {
    updates := availablePlugins(*pluginsSearchRepo)
    tmp, _ := json.Marshal(updates)
    fmt.Println( string(tmp) )
  } else {
    updates := findPlugin(*pluginsSearchRepo, *pluginsSearchName)
    tmp, _ := json.Marshal(updates)
    fmt.Println( string(tmp) )
  }
}

func ListLocalPlugins(){
  //fmt.Println("List Local Plugins")
  updates := installedPlugins()
  tmp, _ := json.Marshal(updates)
  fmt.Println( string(tmp) )
}

// Define all the CLI input flags and subcommands
var (
  app = kingpin.New("ix-alertme", "Alert Notification Plugin System")
  Configfile = app.Flag("config", "Use alternate configuration file").Short('c').Default("/usr/local/etc/ix-alertme.json").String()

  plugins 		= app.Command("plugins", "Plugin Management Functionality")
  pluginsSearch 	= plugins.Command("search", "Search for available plugins")
    pluginsSearchName = pluginsSearch.Arg("name", "Show full details for a specific plugin").Default("").String()
    pluginsSearchRepo = pluginsSearch.Flag("repo", "Restrict to this specific repository").Short('r').Default("").String()
  pluginsList 		= plugins.Command("list", "List all installed plugins")
  pluginsScan		= plugins.Command("scan", "Scan for updates to installed plugins")
)

func main() {
  app.Version("0.1")
  app.UsageTemplate(kingpin.CompactUsageTemplate)
  app.HelpFlag.Short('h')
  switch kingpin.MustParse(app.Parse(os.Args[1:])) {
    case pluginsSearch.FullCommand():
	Config = loadConfiguration(*Configfile)
        ListRemotePlugins()
    case pluginsList.FullCommand():
	Config = loadConfiguration(*Configfile)
        ListLocalPlugins()
    case pluginsScan.FullCommand():
	Config = loadConfiguration(*Configfile)
        CheckForUpdates()
    default:
      app.Fatalf("%s","Unknown subcommand")
  }
}
