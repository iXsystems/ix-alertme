package main

import (
	"time"
	"net/url"
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
	RemoteUrl url.URL		`json:"url"`
	ExtractWith string		`json:"extract_with"`
	Sha256 string			`json:"sha256_checksum"`
}

type PluginDependencies struct {
	Pkg []PkgDependency	`json:"freebsd"`
	File []FileDependency	`json:"file"`
	Archive []FileDependency	`json:"archive"`
}

type PluginIndexManifest struct {
	Name string				`json:"name"`
	Summary string			`json:"summary"`
	Description string			`json:"description"`
	IconUrl url.URL			`json:"icon_url"`
	Version string				`json:"version"`
	VersionReleased time.Time	`json:"date_released"`
	RepoName string
}

type Person struct {
	Name string	`json:"name"`
	Email string	`json:"email"`
	Url	url.URL	`json:"site_url"`
}

type PluginFullManifest struct {
	Name string				`json:"name"`
	Summary string			`json:"summary"`
	Description string			`json:"description"`
	IconUrl url.URL			`json:"icon_url"`
	Version string				`json:"version"`
	VersionReleased time.Time	`json:"date_released"`
	Maintainers []Person		`json:"maintainer"`
	RepoName string
	Depends	PluginDependencies	`json:"depends"`
	Exec FileDependency		`json:"exec"`
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

func availablePlugins() map[string]PluginIndexManifest {
  out := make(map[string]PluginIndexManifest)
  for index := range(Config.RepoList) {
    repo := Config.RepoList[index]
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

func pluginUpdates() map[string]PluginIndexManifest {
  out := make(map[string]PluginIndexManifest)
  installed := installedPlugins()
  if len(installed) <1 { return out } //nothing installed
  available := availablePlugins()
  for name, plugin := range installed {
    if aplugin, ok := available[name] ; ok {
      //Plugin available remotely. Compare release dates to see if remote is newer
      if aplugin.VersionReleased.After(plugin.VersionReleased) {
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
