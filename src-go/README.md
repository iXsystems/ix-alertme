# ix-alertme Utilities
## Package Dependencies
* go
   * FreeBSD/TrueOS (lang/go): `pkg install go`

## Build Instructions
* `go get`
* `go build`

## ix-alertme
Primary utility that manages alert backend plugins and interacts with them.

### Usage
Run `ix-alertme help` to see complete usage information.

#### Config File Format
The default configuration file is located at "/usr/local/etc/ix-alertme.json", but can be pointed to any other location via the "-c <path>" global command-line flag.

Default configuration file:
```
{
  "install_dir" : "",
  "repos" : [
    { "name" : "ix-alertme", "url" : "https://raw.githubusercontent.com/iXsystems/ix-alertme/master/provider-plugins" }
  ]
}
```
If the "install_dir" is unspecified, then plugins will get installed into "/usr/local/ix-alertme/plugins" when run as the root user, and "~/.local/ix-alertme/plugins" if run as a non-root user.

The repository URL should point to the directory containing the "index.json" provided by the remote repository. All of the individual manifests for the plugins are relative to that directory.

#### Plugin Settings File Format
ix-alertme does not save or manage the plugin settings. Instead, when an alert is submitted the user must provide the path to the settings file to use for that particular alert submission. This allows for the user to maintain as many different plugin settings structures as desired (such as for a multi-user system where each user has their own settings). Or even to store/maintain plugin settings in an external storage format (such as a database), and simply use a temporary in-memory file for the current alert submission.

**Example Settings File**
```
{
  "plugin_1" : {
    "api_field_1" : "example"
  },
  "plugin_2" : {
    "p2_api_field_1" : 10
  }
}
```
When sending an alert, ix-alertme will attempt to send the alert using *every* plugin that is configured. In this example, it will verify that the API settings for "plugin_1" and "plugin_2" are valid before submitting the alert to both plugins.


## mkindex
Small script which automates the generation/update of the index.json file for a plugin repository.

### Usage
* `mkindex <index directory>`
