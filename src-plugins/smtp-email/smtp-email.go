package main

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"net/smtp"
	"strings"
	"strconv"
	"fmt"
)

type API struct {
	Mailserver string	`json:"mailserver"`
	MailserverPort int	`json:"mailserver_port"`
	AuthType string	`json:"auth_type"`
	AuthUser string	`json:"auth_user"`
	AuthPass string	`json:"auth_pass"`
	FromAddr string	`json:"from"`
	ToAddr []string	`json:"to"`
	BccAddr []string	`json:"bcc"`
	CcAddr []string	`json:"cc"`
	Subject string		`json:"subject"`
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

func assembleBody(api AlertAPI) []byte {
  // This SMTP integration does not seem to support sending HTML content.
  //  Prefer plaintext if it is available

  var lines []string
  lines = append(lines, "From: "+api.Settings.FromAddr)
  lines = append(lines, "To: "+strings.Join(api.Settings.ToAddr,",") )
  if len(api.Settings.CcAddr) > 0 {
    lines = append(lines, "Cc: "+strings.Join(api.Settings.CcAddr,",") )
  }
  lines = append(lines, "Subject: "+api.Settings.Subject )
  body := api.Text.Html
  if body == "" {
    body = api.Text.PlainText
  } else {
    //Set the header lines to tag it as an html message
    lines = append(lines, "Mime-Version: 1.0;")
    lines = append(lines, "Content-Type: text/html; charset=\"UTF-8\";")
  }
  lines = append(lines, "\r\n"+body)

  msg := []byte( strings.Join(lines, "\r\n")+"\r\n")
  return msg
}

func setupAuth(api AlertAPI) smtp.Auth {
  var auth smtp.Auth
  if api.Settings.AuthType == "plain" || api.Settings.AuthPass == "" {
    auth = smtp.PlainAuth("", api.Settings.AuthUser, api.Settings.AuthPass, api.Settings.Mailserver)
  } else {
    //Unknown / none
  }
  return auth
}

func main() {
  //fmt.Println("Sending smtp-email...")
  // Parse the input API
  api := readAPI(os.Args[1])
  //Setup the authentication
  auth := setupAuth(api)
  //Send the email(s)
  toall := append(api.Settings.ToAddr, api.Settings.BccAddr...)
  toall = append(toall, api.Settings.CcAddr...)
  err := smtp.SendMail( api.Settings.Mailserver+":"+strconv.Itoa(api.Settings.MailserverPort), 
		auth, api.Settings.FromAddr, 
		toall, assembleBody(api) )
  if err != nil {
    fmt.Println("Error sending smtp-email:", err)
    os.Exit(1)
  }
  os.Exit(0)
}
