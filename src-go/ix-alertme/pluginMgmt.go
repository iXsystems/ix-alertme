package main

import (
	"fmt"
	"time"
	"os"
	"io/ioutil"
	"encoding/json"
	"errors"
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

func installedPlugins() (map[string]PluginFullManifest, error) {
  out := make(map[string]PluginFullManifest)
  filelist, _ := ioutil.ReadDir(Config.InstallDir)
  for i := range(filelist) {
    if filelist[i].IsDir() {
      //PrintDebug("Check Dir: "+Config.InstallDir+"/"+filelist[i].Name()+"/manifest.json")
      if _, err := os.Stat(Config.InstallDir+"/"+filelist[i].Name()+"/manifest.json") ; err == nil {
        tmp, err := ioutil.ReadFile(Config.InstallDir+"/"+filelist[i].Name()+"/manifest.json")
        //PrintDebug("Read File:"+string(tmp))
        if err != nil { continue }
        var tman PluginFullManifest
         json.Unmarshal(tmp, &tman)
         if(tman.Name != ""){ out[tman.Name] = tman }
      }
    }
  }
  return out, nil
}

func availablePlugins(repolimit string) (map[string]PluginIndexManifest, error) {
  out := make(map[string]PluginIndexManifest)
  var err error
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
  return out, err
}

func findPlugin(repolimit string, name string) (PluginFullManifest, error) {
  var out PluginFullManifest
  var err error
  for index := range(Config.RepoList) {
    repo := Config.RepoList[index]
    if repolimit != "" && repo.Name != repolimit { continue } //wrong repository
    //Got a specific plugin to find and show details for
    plugin, err := FetchPluginManifest(repo, name)
    if err != nil { continue } //not available in this repo
    plugin.RepoName = repo.Name
    return plugin, nil
  }
  return out, err
}

func pluginUpdates() (map[string]PluginIndexManifest, error) {
  out := make(map[string]PluginIndexManifest)
  var err error
  installed, err := installedPlugins()
  if len(installed) <1 { return out, nil } //nothing installed
  available, err := availablePlugins("")
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
  return out, err
}

func uninstallPlugin(name string) error {
  installdir := Config.InstallDir+"/"+name
  err := os.RemoveAll(installdir)
  //Optional Later - uninstall package dependencies
  // Needs special handling to prevent dependency breakage for other plugins
  if err == nil {
    fmt.Println("Plugin removed: ", name)
  }
  return err
}

func installPlugin(name string, repolimit string) error {
  plugin, _ := findPlugin(repolimit, name)
  if plugin.Name == "" {
    //Could not find plugin to install
    return errors.New("Plugin Unavailable:"+name)
  }
  installdir := Config.InstallDir+"/"+name
  var err error
  err = nil
  // Verify that the plugin is not already installed

  // Perform installation
    //Create the directory
    err = os.MkdirAll(installdir, 0744)
    // Install FreeBSD packages
    for i := range(plugin.Depends.Pkg) {
      if err != nil { break }
      pkg := plugin.Depends.Pkg[i]
      // See if the pkg is already installed
      PrintDebug("TODO - Check for package: "+pkg.PkgName)
      // Install the pkg if necessary

    }
    // Download any archives into the install dir
    for i := range(plugin.Depends.File) {
      if err != nil { break }
      err = InstallFileDependency(plugin.Depends.File[i], installdir+"/files")
    }
    // Download any files into the install dir
    for i := range(plugin.Depends.Archive) {
      if err != nil { break }
      err = InstallFileDependency(plugin.Depends.Archive[i], installdir+"/files")
    }
    // Download the exec file into the install dir
    if err == nil {
      err = InstallFileDependency(plugin.Exec, installdir)
    }
    // Save the manifest into the install dir
    if err == nil {
      tmp, _ := json.Marshal(plugin)
      err = ioutil.WriteFile(installdir+"/manifest.json", tmp, 0644)
    } else {
      // Had Error: Cleanup the partially-installed plugin
      os.RemoveAll(installdir)
    }
  if err == nil {
    fmt.Println("Plugin installed: ", name)
  }
  return err
}


