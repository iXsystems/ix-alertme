package main

import (
	"os"
	"flag"
	"fmt"
	"encoding/json"
)

var Config Configuration

func CheckForUpdates(){
  updates := pluginUpdates()
  tmp, _ := json.Marshal(updates)
  fmt.Println( string(tmp) )
}

func ListPlugins(){

}

func showHelp(){
  fmt.Println("Provide the -h or --help flags to see usage information")
  os.Exit(1)
}

func main() {
  if len(os.Args) < 2 {
    showHelp()
  }
  //Define CLI flags
  configfile := flag.String("c","/usr/local/etc/ix-alertme.json", "Use custom configuration file")
  flag.Parse()
  Config = loadConfiguration(*configfile)
  //Now parse the input arguments and do the things
  cmd := "unknown"
  if len(flag.Args()) < 1 {
    showHelp()
  }
  cmd = flag.Args()[0]
  fmt.Println("Got Command:", cmd, "ConfigFile:", *configfile)
  switch cmd {
    case "check-updates" : CheckForUpdates()
    default: 
    fmt.Println("Unknown Option:",cmd); showHelp()
  }
}
