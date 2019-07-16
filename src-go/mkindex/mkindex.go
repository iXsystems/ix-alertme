package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
)

type PluginIndexManifest struct {
	Name string				`json:"name"`
	Summary string			`json:"summary"`
	Description string			`json:"description"`
	IconUrl string				`json:"icon_url"`
	Version string				`json:"version"`
	VersionReleased string		`json:"date_released"`
	Tags	 []string				`json:"tags"`
}

var TotalIndex []PluginIndexManifest;

func importManifest(path string){
  //fmt.Println("Check Manifest: " +path)
  tmp, err := ioutil.ReadFile(path);
  if err != nil { return }
  var plugin PluginIndexManifest
  err = json.Unmarshal(tmp, &plugin)
  if err == nil {
    fmt.Println("Valid Plugin: "+plugin.Name)
    TotalIndex = append(TotalIndex, plugin)
  } else {
    fmt.Println("Error In Manifest: "+path, err)
  }
}

func main() {
  // Get the directory input
  cdir := os.Args[1]
  // Scan for subdirectories and import manifests
  filelist, _ := ioutil.ReadDir(cdir)
  for i := range(filelist) {
    if filelist[i].IsDir() {
      if _, err := os.Stat(cdir+"/"+filelist[i].Name()+"/manifest.json") ; err == nil {
        importManifest(cdir+"/"+filelist[i].Name()+"/manifest.json")
      }
    }
  }
  // Now write out the new index file
  tmp, _ := json.Marshal(TotalIndex)
  ioutil.WriteFile(cdir+"/index.json", tmp, 0644)
  fmt.Println("Index Updated")
}
