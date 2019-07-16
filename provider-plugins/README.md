# Provider Plugins
These are the manifest files and/or scripts needed to utilize backend providers. Every plugin needs to be a completely self-contained tool which can be installed into any arbitrary directory on a target system. These plugins must also take a single JSON file as input, with a format specified below.

## Plugin Input API
```
{
  "text" : {
    "html" : "<html encoding=Utf-8><p>Html Form of the alert</p></html>",
    "plaintext" : "Plaintext form of the alert"
  },
  "settings" : {
    "<custom_API_field>" : <value>
  }
}
```

## Plugin Manifest
Each plugin provides a single JSON manifest file with all the necessary information about the plugin itself (name, version, etc).

### Manifest Field Details
* ***name*** (string) : Name of the plugin
* ***summary*** (string) : Short summary of what the plugin does
* ***description*** (string) : Longer description of what the plugin does
* ***icon_url*** (string) : [optional] URL to fetch the icon that represents the plugin (*.jpg or *.png only)
* ***version*** (string) : Version string for the plugin (informational only)
* ***date_released*** (string) : Timestamp in the format "yyyy-MM-dd" or "yyyy-MM-ddThh:mm:ss" (all times in UTC)
   * Example: "2019-01-02" or "2019-01-02T05:30:00" for January 2, 2019 at 5:30 AM
   * This date/time stamp is used to determine if the plugin has any updates available. This value should **always** increase, and never go backwards between releases of a plugin.
* ***tags*** (Json Array of strings) : [optional] Additional tags to aid in plugin searches
* ***maintainer*** (Json Array of Objects) : List of maintainers for the plugin. Object format for each maintainer is listed below:
   * ***name*** (string) : Name of the maintainer (Example: "John Doe")
   * ***email*** (string) : Email address for the maintainer (Example: "john.doe@example.net")
   * ***site_url*** (string) : URL to a website for the maintainer (Example: "https://github.com/john_doe_124978")
* ***depends*** (Json Array of Objects) : This object lists the pieces of the plugin itself. Everything listed here will get extracted into the same directory. The location of the install directory **is not consistent** between systems, nor is it fixed for a single system, so that a plugin could be installed multiple times on the same system as needed. Object format for each dependency is listed below
   * ***url*** (string) : URL for where to fetch the file from (always use HTTPS if possible)
   * ***sha256_checksum*** (string) : Checksum of the file for post-download verification.
   * ***extract*** (boolean) : [optional] Flag whether the downloaded file needs to be extracted (such as when downloading an archive of files). False by default.
      * Supported file formats: *.zip, *.tar (and compressed variants like *.tar.gz), and *.rar
      * See [https://github.com/mholt/archiver](https://github.com/mholt/archiver) for the full list of supported formats.
   * ***decompress*** (boolean) : [optional] Flag whether the downloaded file needs to be decompressed. For compressed archives, use the "extract" flag instead. False by default.
      * See [https://github.com/mholt/archiver](https://github.com/mholt/archiver) for the full list of supported formats.
   * ***filename*** (string) : [optional] If extraction and decompression are not needed, this can be provided to change the name of the resulting file in the plugin directory.
* ***exec*** (string) : Name of the binary from the plugin directory to execute. Must be installed via a "depends" entry.
* ***api*** (Json Array of Objects) : List of API fields which the plugin supports or needs in order to function. Object format for a single api entry is listed below:
   * ***fieldname*** (string) : Name of the JSON field for this API input.
   * ***is_required*** (boolean) : [optional] Indicate whether this field is required or not. False by default.
   * ***summary*** (string) : Short summary of how this field is used.
   * ***is_array*** (boolean) : [optional] Indicate whether this field should be an array of values. False by default.
      * Note that this flag may not be used with the special "Lists" type of inputs.
   * ***type*** (see below) : This defines any rules/checks for validating the input(s)
      * *Numbers*
         * "integer" for an integer value, or ["integer", min, max] for a specific range of values
         * "float" for a decimel value, or ["float", min, max] for a specific range of values
      * *Text*
         * "" for any string value, or ["regex", "<regular_expression>"] for a string that exactly matches the regular expression.
      * *Lists*
         * ["select", A, B, C] to prompt the user to select **one** of the options (A, B, or C in this case)
            * The options can have short summaries with the format [A, "Option A"]
            * Example: ["select", ["A", "Option A details"], [ "SomeB", "Some Option B"] ]
      * *True/False*
         * "bool" indicates that a true/false value is required
   * ***default*** (anything - see examples) : [optional] Default value for this field if nothing is provided
      * Any valid JSON can be placed here. "strings", numbers (5.5), booleans (true/false), or even arrays of values.
      * It is recommended to avoid using Json Objects as values, as these are not enforcable via the API check mechanisms.


### Manifest Example
```
{
  "name" : "example",
  "summary" : "Example manifest",
  "description" : "Example manifest for learning purposes. This can be copied as a template for future plugins as well.",
  "icon_url" : "https://my.example.net/icon.png",
  "version" : "1.0",
  "date_released" : "2019-07-16",
  "tags" : ["example","plugin","manifest"],
  "maintainer" : [
    {"name" : "John Doe", "email" : "john.doe@example.net", "site_url" : "http://my.example.net" }
  ],
  "depends" : [
    { "url" : "https://my.example.net/alert-plugin-example", "filename" : "example_binary", "sha256_checksum" : "ABCDEFGHIJ123456789" }
  ],
  "exec" : "example_binary",
  "api" : [
    {"fieldname" : "booltest", "summary" : "Example of a true/false input", "type" : "bool", "default" : false },
    {"fieldname" : "stringtest", "summary" : "Example of a string input", "type" : "", "default" : "default text" },
    {"fieldname" : "integertest", "summary" : "Example of an integer input from 0-100", "type" : ["integer", 0, 100], "default" : 50 },
    {"fieldname" : "selecttest", "summary" : "Example of a list selection input", "type" : ["select","A","B",["C","Option C Summary"]], "default" : "A" }
  ]
}
```
