package main

import(
	"encoding/json"
	"os"
	"io/ioutil"
	"os/user"
)

type Repo struct {
	Name string 	`json:"name"`
	Url string		`json:"url"`
}

type Configuration struct {
	RepoList []Repo	`json:"repos"`
	InstallDir string	`json:"install_dir"`
}

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func loadConfiguration(location string) Configuration {
  var config Configuration
  srchpaths := []string{ "/usr/local/", "/usr/", "/" }
  if location == "" { 
    // Default config file location - look for it in the search paths
    for _, dir := range(srchpaths) {
      if fileExists( dir+"etc/ix-alertme.json" ){
        // Found config file
        location = dir+"etc/ix-alertme.json"
        break
      }
    }
  }
  //Load the local config file
  file, err := os.Open(location)
  defer file.Close()
  if err == nil {
    tmp, err := ioutil.ReadAll(file)
    if err == nil { 
      json.Unmarshal(tmp, &config)
    }
  }
  //Load the default values if not specified
  if config.InstallDir == "" { 
    if os.Getuid() == 0 {
      config.InstallDir = "/usr/local/ix-alertme/plugins"
    } else {
      cuser, _ := user.Current()
      config.InstallDir = cuser.HomeDir;
      config.InstallDir = config.InstallDir + "/.local/ix-alertme/plugins"
    }
  }

  //Append the default repository to the list if it is empty
  //fmt.Println("Repos:", config.RepoList, len(config.RepoList))
  if len(config.RepoList) <1 {
    //fmt.Println("Repo List Empty")
    var defaultrepo Repo
    defaultrepo.Name = "ix-alertme"
    defaultrepo.Url = "https://raw.githubusercontent.com/iXsystems/ix-alertme/master/provider-plugins"
    config.RepoList = append(config.RepoList, defaultrepo)
  }
  //Return the configuration
  return config
}
