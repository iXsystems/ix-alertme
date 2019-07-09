package main

import(
	"encoding/json"
	"os"
	"io/ioutil"
)

type Repo struct {
	Name string 	`json:"name"`
	Url string		`json:"url"`
}

type Configuration struct {
	RepoList []Repo	`json:"repos"`
	InstallDir string	`json:"install_dir"`
	
}

func loadConfiguration(location string) Configuration {
  var config Configuration
  //Load the local config file
  file, err := os.Open(location)
  defer file.Close()
  if err != nil {
    tmp, err := ioutil.ReadAll(file)
    if err != nil { json.Unmarshal(tmp, config) }
  }
  //Load the default values if not specified
  if config.InstallDir == "" { config.InstallDir = "/usr/local/ix-alertme/plugins" }

  //Append the default repository to the list if it is empty
  if len(config.RepoList) <1 {
    var defaultrepo Repo
    defaultrepo.Name = "ix-alertme"
    defaultrepo.Url = "https://raw.githubusercontent.com/iXsystems/ix-alertme/master/provider-plugins"
    config.RepoList = append(config.RepoList, defaultrepo)
  }
  //Return the configuration
  return config
}
