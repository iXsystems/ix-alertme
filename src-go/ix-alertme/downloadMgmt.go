package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"encoding/json"
	"path"
	"crypto/sha256"
	"errors"
	"github.com/mholt/archiver"
)

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string, checksum string) error {
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
    // Compare the checksum if provided
    if checksum != "" {
      f, err := os.Open("file.txt")
      if err == nil {
        defer f.Close()
        h := sha256.New()
        if _, err := io.Copy(h, f); err == nil {
          if string(h.Sum(nil)) != checksum {
            err = errors.New("Checksum Mismatch: "+filepath)
            os.Remove(filepath)
          }
        }
      }
    } //end checksum empty check
    return err
}

func FetchPluginIndex(repo Repo) ([]PluginIndexManifest, error) {
    var info []PluginIndexManifest
    // Get the data
    //PrintDebug("Fetch Plugin index: "+repo.Url+"/index.json")
    resp, err := http.Get(repo.Url+"/index.json")
    if err != nil {
        PrintError("Repo index unavailable: "+repo.Name)
        return info, err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil { 
      PrintError("Repo index unreadable: "+repo.Name)
      PrintError( string(body) )
      return info, err 
    }
    //Now convert the response into the data structure and return it
    err = json.Unmarshal(body, &info)
    if err != nil { 
      PrintError("Repo index malformed: "+repo.Name)
      PrintError( string(body) )
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

func InstallFileDependency( file FileDependency, installdir string) error {
  var err error = nil
  if file.IsArchive || file.IsCompressed{
    // Archive. Need to download to a temporary dir and then extract into install dir
    tmpfile := os.TempDir()+"/"+path.Base( path.Clean(file.RemoteUrl) )
    err = DownloadFile(tmpfile, file.RemoteUrl, file.Sha256)
    if err == nil {
      if file.IsArchive {
        err = archiver.Unarchive(tmpfile, installdir)
      } else if file.IsCompressed {
        err = archiver.DecompressFile(tmpfile, installdir)
      }
    }
    os.Remove(tmpfile) //finished with the temporary file
  } else {
    // File. Just download directly to the install dir
    filename := file.Filename
    if filename == "" {
      //Pull the filename off of the URL
      filename = path.Base( path.Clean(file.RemoteUrl) )
    }
    err = DownloadFile(installdir+"/"+filename, file.RemoteUrl, file.Sha256)
  }
  return err
}
