# Snap plugin processor - logs-regexp

Snap plugin intended to process logs using regular expressions.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

### Installation
#### Download logs-regexp plugin binary:
You can get the pre-built binaries for your OS and architecture from the plugin's [GitHub Releases](https://github.com/intelsdi-x/snap-plugin-processor-logs-regexp/releases) page.
Download the plugin from the latest release and load it into `snapteld` (`/opt/snap/plugins` is the default location for Snap packages).

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-processor-logs-regexp

Clone repo into `$GOPATH/src/github/intelsdi-x/`:
```
$ git clone https://github.com/<yourGithubID>/snap-plugin-processor-logs-regexp
```
Build the plugin by running make in repo:
```
$ make
```
This builds the plugin in `./build`

### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap#getting-started)

## Documentation

The intention of this plugin is to parse logs according to the form defined by the user in task manifest.

The plugin processes log line using regular expressions with names. These names indicate tags which are added to metric.
The regular expression can contain special names (`timestamp` and `message`) to replace metric timestamp and value.

The plugin can be configured by following parameters (all parameters are optional):
- `regexp_log` - regular expression with method to parse log line, special names:
    - `message` - the main message in log which will be set as a metric value
    - `timestamp` - the log timestamp which will be set as a metric timestamp.
- `regexp_message` - regular expression with method to parse the main message in log (indicated by `message`).
- `regexp_time` - regular expression with method to parse the log timestamp (indicated by `timestamp`), this expression needs to contain following names:
   `day`, `month`, `year`, `hour`, `seconds`, `minutes`, `seconds`, `timezone`.

Notice: Special characters in regular expressions needs to be escaped.

On default plugin processes following log line (metric value which is returned by collector plugin):

```
127.0.0.1 - - [07/Dec/2016:06:00:12 -0500] "GET /v3/users/fa2b2986c200431b8119035d4a47d420/projects HTTP/1.1" 200 446 21747 "-" "python-keystoneclient"
```
to Snap's metric with:
- Data: `"GET /v3/users/fa2b2986c200431b8119035d4a47d420/projects HTTP/1.1" 200 446 21747 "-" "python-keystoneclient"`
- Timestamp: `07/Dec/2016:06:00:12 -0500`
- Tags:
    - `client_ip`: `127.0.0.1`
    - `http_method`:`GET`
    - `http_url`: `/v3/users/fa2b2986c200431b8119035d4a47d420/projects`
    - `http_version`: `1.1`
    - `http_status`: `200`
    - `http_response_size`: `446`
    - `http_response_time`: `21747`

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-processor-logs-regexp/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-processor-logs-regexp/pulls).

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[Snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements

* Author: [Katarzyna Kujawa](https://github.com/katarzyna-z)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
