package main

import (
	"time"
	"os"
	"io/ioutil"
	"encoding/json"
)
type PkgDependency struct {
	PkgName string	`json:"pkg"`
	PortOrigin string	`json:"port"`
}

type FileDependency struct {
	Filename string		`json:"filename"`
	RemoteUrl string		`json:"url"`
	ExtractWith string		`json:"extract_with"`
	Sha256 string			`json:"sha256_checksum"`
}

type PluginDependencies struct {
	Pkg []PkgDependency		`json:"freebsd"`
	File []FileDependency		`json:"file"`
	Archive []FileDependency	`json:"archive"`
}

type PluginIndexManifest struct {
	Name string				`json:"name"`
	Summary string			`json:"summary"`
	Description string			`json:"description"`
	IconUrl string				`json:"icon_url"`
	Version string				`json:"version"`
	VersionReleased string		`json:"date_released"`
	RepoName string			`json:"repository"`
}

type Person struct {
	Name string	`json:"name"`
	Email string	`json:"email"`
	Url	string	`json:"site_url"`
}

type SetOpts struct {
	Field string	`json:"fieldname"`
	Description string	`json:"description"`
	Default interface{}	`json:"default"`
	Type interface{}	`json:"type"`
	Required bool		`json:"is_required"`
}
type PluginFullManifest struct {
	Name string				`json:"name"`
	Summary string			`json:"summary"`
	Description string			`json:"description"`
	IconUrl string				`json:"icon_url"`
	Version string				`json:"version"`
	VersionReleased string		`json:"date_released"`
	Tags []string				`json:"tags"`
	Maintainers []Person		`json:"maintainer"`
	RepoName string			`json:"repository"`
	Depends	PluginDependencies	`json:"depends"`
	Exec FileDependency			`json:"exec"`
	API []SetOpts				`json:"api"`
}

func Timestamp(t string) time.Time {
  //Run through all the supported timestamp formats and exit with the first successful match
  stamp, err := time.Parse("2006-01-02T15:04:05",t);
  if err == nil { return stamp }
  stamp, err = time.Parse("2006-01-02",t);
  if err == nil { return stamp }
  //Always exit at the end with a valid time structure
  return time.Now().AddDate(-10,0,0) //Now minus 10 years
}

func installedPlugins() map[string]PluginFullManifest {
  out := make(map[string]PluginFullManifest)
  filelist, _ := ioutil.ReadDir(Config.InstallDir)
  for i := range(filelist) {
    if filelist[i].IsDir() {
      if _, err := os.Stat(Config.InstallDir+"/"+filelist[i].Name()+"/manifest.json") ; os.IsExist(err) {
	tmp, err := ioutil.ReadFile(Config.InstallDir+"/"+filelist[i].Name()+"/manifest.json")
        if err != nil { continue }
        var tman PluginFullManifest
         json.Unmarshal(tmp, tman)
         if(tman.Name != ""){ out[tman.Name] = tman }
      }
    }
  }
  return out
}

func availablePlugins(repolimit string) map[string]PluginIndexManifest {
  out := make(map[string]PluginIndexManifest)
  for index := range(Config.RepoList) {
    repo := Config.RepoList[index]
    if repolimit != "" && repo.Name != repolimit { continue } //wrong repository
    //Show the general info about all available plugins
    list, err := FetchPluginIndex(repo)
    if err != nil { continue }
    for pindex := range(list) {
      plugin := list[pindex]
      if _, ok := out[plugin.Name] ; !ok {
        plugin.RepoName = repo.Name
        out[plugin.Name] = plugin
      }
    }
  }
  return out
}

func findPlugin(repolimit string, name string) PluginFullManifest {
  var out PluginFullManifest
  for index := range(Config.RepoList) {
    repo := Config.RepoList[index]
    if repolimit != "" && repo.Name != repolimit { continue } //wrong repository
    //Got a specific plugin to find and show details for
    plugin, err := FetchPluginManifest(repo, name)
    if err != nil { continue } //not available in this repo
    plugin.RepoName = repo.Name
    return plugin
  }
  return out
}

func pluginUpdates() map[string]PluginIndexManifest {
  out := make(map[string]PluginIndexManifest)
  installed := installedPlugins()
  if len(installed) <1 { return out } //nothing installed
  available := availablePlugins("")
  for name, plugin := range installed {
    if aplugin, ok := available[name] ; ok {
      //Plugin available remotely. Compare release dates to see if remote is newer
      if Timestamp(aplugin.VersionReleased).After( Timestamp(plugin.VersionReleased)) {
        //Remote plugin newer - flag it as an update
        out[name] = aplugin
      }
    } else {
      //Plugin not available remotely (at all - removed from listings?)
      out[name] = aplugin //empty structure
    }
  }
  return out
}
