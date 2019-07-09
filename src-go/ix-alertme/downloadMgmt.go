package main

import (
	//"fmt"
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
        return err
    }
    defer resp.Body.Close()

    // Create the file
    out, err := os.Create(filepath)
    if err != nil {
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
        //fmt.Println(" - Error Fetching index", err)
        return info, err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil { return info, err }
    //fmt.Println("Got Body:", string(body))
    //Now convert the response into the data structure and return it
    err = json.Unmarshal(body, &info)
    return info, err
}

func FetchPluginManifest(repo Repo, name string) (PluginFullManifest, error) {
    var info PluginFullManifest
    // Get the data
    resp, err := http.Get(repo.Url+"/"+name+"/manifest.json")
    if err != nil {
        return info, err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    //Now convert the response into the data structure and return it
    err = json.Unmarshal(body, info)
    return info, err
}
