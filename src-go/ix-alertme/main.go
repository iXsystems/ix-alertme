package main

import (
	"os"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"encoding/json"
)

var Config Configuration

// Internal simplification functions
func PrintDebug(error string) {
  fmt.Fprintln(os.Stderr, "[Debug] "+error)
}
func PrintError(error string) {
  fmt.Fprintln(os.Stderr, "[Error] "+error)
}

func CheckForUpdates() error {
  //fmt.Println("Check For Updates")
  updates, err := pluginUpdates()
  tmp, _ := json.Marshal(updates)
  fmt.Println( string(tmp) )
  return err
}

func ListRemotePlugins() error {
  if *pluginsSearchName == "" {
    updates, err := availablePlugins(*pluginsSearchRepo)
    tmp, _ := json.Marshal(updates)
    fmt.Println( string(tmp) )
    return err
  } else {
    updates, err := findPlugin(*pluginsSearchRepo, *pluginsSearchName)
    tmp, _ := json.Marshal(updates)
    fmt.Println( string(tmp) )
    return err
  }

}

func ListLocalPlugins() error {
  //fmt.Println("List Local Plugins")
  updates, err := installedPlugins()
  tmp, _ := json.Marshal(updates)
  fmt.Println( string(tmp) )
  return err
}

// Define all the CLI input flags and subcommands
var (
  app = kingpin.New("ix-alertme", "Alert Notification Plugin System")
  Configfile = app.Flag("config", "Use alternate configuration file").Short('c').Default("/usr/local/etc/ix-alertme.json").String()

  plugins 		= app.Command("plugins", "Plugin Management Functionality")
    pluginsSearch 	= plugins.Command("search", "Search for available plugins")
      pluginsSearchName = pluginsSearch.Arg("name", "Show full details for a specific plugin").Default("").String()
      pluginsSearchRepo = pluginsSearch.Flag("repo", "Restrict to a specific repository").Short('r').String()
    pluginsList 		= plugins.Command("list", "List all installed plugins")
    pluginsScan		= plugins.Command("scan", "Scan for updates to installed plugins")
    pluginsInstall		= plugins.Command("install", "Download and install a plugin")
      pluginsInstallName = pluginsInstall.Arg("name", "Name of the plugin to install").Required().String()
      pluginsInstallRepo = pluginsInstall.Flag("repo", "Restrict to a specific repository").Short('r').String()
    pluginsRemove	= plugins.Command("remove", "Delete an installed plugin")
      pluginsRemoveName = pluginsRemove.Arg("name", "Name of the plugin to remove").Required().String()
)

func main() {
  app.Version("0.1")
  app.UsageTemplate(kingpin.CompactUsageTemplate)
  app.HelpFlag.Short('h')
  var err error
  switch kingpin.MustParse(app.Parse(os.Args[1:])) {
    case pluginsSearch.FullCommand():
	Config = loadConfiguration(*Configfile)
        err = ListRemotePlugins()
    case pluginsList.FullCommand():
	Config = loadConfiguration(*Configfile)
        err = ListLocalPlugins()
    case pluginsScan.FullCommand():
	Config = loadConfiguration(*Configfile)
        err = CheckForUpdates()
    case pluginsInstall.FullCommand():
	Config = loadConfiguration(*Configfile)
        err = installPlugin(*pluginsInstallName, *pluginsInstallRepo)
    case pluginsRemove.FullCommand():
	Config = loadConfiguration(*Configfile)
        err = uninstallPlugin(*pluginsRemoveName)
    default:
      app.Fatalf("%s","Unknown subcommand")
  }
  if err != nil {
    fmt.Println(err) //standard output
    os.Exit(1)
  } else {
    os.Exit(0)
  }  
}
