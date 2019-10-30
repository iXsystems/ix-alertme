package main

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/PagerDuty/go-pagerduty"
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

func mkIncidentFromAPI(api AlertAPI) pagerduty.CreateIncidentOptions {
  var incident pagerduty.CreateIncidentOptions
  incident.Type = api.Settings.Type
  incident.Title = api.Settings.Title
  incident.Service = &pagerduty.APIReference{ api.Settings.Service ,"" }
  incident.Priority = &pagerduty.APIReference{ "warning" ,"" }
  incident.Body = &pagerduty.APIDetails{ "incident_body", api.Text.PlainText }
  return incident
}

func main() {
  // Parse the input API
  api := readAPI(os.Args[1])

  client := pagerduty.NewClient(api.Settings.Authtoken);
  ciopts := mkIncidentFromAPI(api) //Create Incident Options
  tmp, _ := json.MarshalIndent(ciopts, "", "  ")
  fmt.Println("Create Incident:", string(tmp) )
  _, err := client.CreateIncident( api.Settings.From, &ciopts);

  if err != nil {
    fmt.Println("Error sending pagerduty incident:", err)
    os.Exit(1)
  }
  os.Exit(0)
}
