package main

type PluginSettings struct {
	Value int
}

type PluginApi struct {
	AlertText string
	Settings PluginSettings
}
