package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"encoding/json"
)

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {
    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        PrintError("File Unavailable: "+url)
        return err
    }
    defer resp.Body.Close()

    // Create the file
    out, err := os.Create(filepath)
    if err != nil {
        PrintError("File cannot be created: "+filepath)
        return err
    }
    defer out.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    return err
}

func FetchPluginIndex(repo Repo) ([]PluginIndexManifest, error) {
    var info []PluginIndexManifest
    // Get the data
    //fmt.Println("Fetch Plugin index:", repo.Url+"/index.json")
    resp, err := http.Get(repo.Url+"/index.json")
    if err != nil {
        PrintError("Repo index unavailable: "+repo.Name)
        return info, err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil { 
      PrintError("Repo index unreadable: "+repo.Name)
      return info, err 
    }
    //Now convert the response into the data structure and return it
    err = json.Unmarshal(body, &info)
    if err != nil { 
      PrintError("Repo index malformed: "+repo.Name)
    }
    return info, err
}

func FetchPluginManifest(repo Repo, name string) (PluginFullManifest, error) {
    var info PluginFullManifest
    // Get the data
    resp, err := http.Get(repo.Url+"/"+name+"/manifest.json")
    if err != nil {
        PrintError( "Unable to fetch plugin manifest: "+name+", From Repo: "+repo.Name )
        return info, err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    //Now convert the response into the data structure and return it
    err = json.Unmarshal(body, &info)
    if err != nil {
      PrintError( "Plugin manifest malformed: "+name+", From Repo: "+repo.Name )
    }
    return info, err
}
