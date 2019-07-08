package main

import (
	"time"
)

type PluginIndexManifest struct {
	Name string				`json:"name"`
	Description string			`json:"description"`
	IconUrl string				`json:"icon_url"`
	Version string				`json:"version"`
	VersionReleased time.Time	`json:"date_released"`
	RepoName string
}

type PluginFullManifest struct {
	Name string
	Description string
	IconUrl string
	Version string
	VersionReleased time.Time	
}

func installedPlugins() map[string]PluginIndexManifest {
  out := make(map[string]PluginIndexManifest)

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

