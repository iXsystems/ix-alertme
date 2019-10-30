package main

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/marcw/pagerduty"
	"fmt"
)

type API struct {
	Authtoken string	`json:"authtoken"`
	Type string		`json:"type"`
	Title string		`json:"title"`
	From string		`json:"from"`
        Service string		`json:"service"`
}

type AlertText struct {
	Html string 		`json:"html"`
	PlainText string	`json:"plain"`
}

type AlertAPI struct {
	Text AlertText		`json:"text"`
	Settings API		`json:"settings"`
}

func readAPI(path string) AlertAPI {
  var api AlertAPI
  tmp, err := ioutil.ReadFile(os.Args[1])
  if err != nil { fmt.Println("Cannot read API file: ", path) ; os.Exit(1) } //cannot read input JSON
  err = json.Unmarshal(tmp, &api)
  if err != nil { fmt.Println("Cannot read API JSON: ", path) ; os.Exit(1) } //cannot read input JSON
  return api
}

func main() {
  // Parse the input API
  api := readAPI(os.Args[1])
  event := pagerduty.NewTriggerEvent(api.Settings.Authtoken, api.Settings.Title)
  event.Description = api.Text.PlainText
  event.Details["text"] = api.Text.PlainText
  response, status, err := pagerduty.Submit(event)
  fmt.Println("Got response:", response);
  fmt.Println("Got status:", status);
  if err != nil {
    fmt.Println("Error sending pagerduty incident:", err)
    os.Exit(1)
  }
  os.Exit(0)
}
