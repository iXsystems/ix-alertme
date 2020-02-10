package main

import (
	"os"
	"io/ioutil"
	"encoding/json"
	jira "gopkg.in/andygrunwald/go-jira.v1"
	"fmt"
)

type API struct {
	Username string		`json:"username"`
	Password string		`json:"password"`
	Server string		  `json:"server"`
	Project string		`json:"project"`
	Assignee string 	`json:"assignee"`
}

type AlertText struct {
	Html string 		  `json:"html"`
	PlainText string	`json:"plain"`
}

type AlertAPI struct {
	Text AlertText		`json:"text"`
	Settings API		  `json:"settings"`
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
  tp := jira.BasicAuthTransport{
    Username: api.Settings.Username,
    Password: api.Settings.Password,
  }

  jiraClient, err := jira.NewClient(tp.Client(), api.Settings.Server)
  if err != nil {
    panic(err)
  }

  i := jira.Issue{
    Fields: &jira.IssueFields{
      Assignee: &jira.User{
        Name: api.Settings.Assignee,
      },
      Reporter: &jira.User{
        Name: api.Settings.Username,
      },
      Description: "Test Issue",
      Type: jira.IssueType{
        Name: "Bug",
      },
      Project: jira.Project{
        Key: api.Settings.Project,
      },
      Summary: api.Text.PlainText,
    },
  }
  issue, response, err := jiraClient.Issue.Create(&i)
  if err != nil {
    fmt.Printf(err.Error())
    body, _ := ioutil.ReadAll(response.Body)
    fmt.Println(string(body))
    panic(err)
  }

  fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)
}
