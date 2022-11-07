DISCONTINUATION OF PROJECT. 

This project will no longer be maintained by Intel.

This project has been identified as having known security escapes.

Intel has ceased development and contributions including, but not limited to, maintenance, bug fixes, new releases, or updates, to this project.  

Intel no longer accepts patches to this project.
<!--
http://www.apache.org/licenses/LICENSE-2.0.txt


    Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# DISCONTINUATION OF PROJECT 

**This project will no longer be maintained by Intel.  Intel will not provide or guarantee development of or support for this project, including but not limited to, maintenance, bug fixes, new releases or updates.  Patches to this project are no longer accepted by Intel. If you have an ongoing need to use this project, are interested in independently developing it, or would like to maintain patches for the community, please create your own fork of the project.**



# Snap plugin processor - logs-regexp

Snap plugin intended to process logs using regular expressions.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Examples](#examples)
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

### Examples

This is an example running [snap-plugin-collector-logs](https://github.com/intelsdi-x/snap-plugin-collector-logs),
processing collected regexp-logs and writing post-processed data to a file.
It is assumed that you are using the latest Snap binary and plugins.

In one terminal window, open the Snap daemon (in this case with logging set to 1 and trust disabled) with appropriate configuration needed by logs collector.
To do that properly, please follow the instruction on [snap-plugin-collector-logs](https://github.com/intelsdi-x/snap-plugin-collector-logs).
```
$ snapteld -l 1 -t 0 --config config.json
```
In another terminal window:

Download and load plugins:
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-logs/latest/linux/x86_64/snap-plugin-collector-log
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-processor-logs-regexp/latest/linux/x86_64/snap-plugin-processor-logs-regexp
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
$ snaptel plugin load snap-plugin-collector-logs
$ snaptel plugin load snap-plugin-processor-logs-regexp
$ snaptel plugin load snap-plugin-publisher-file
```
Create a task manifest - see examplary task manifests in [examples/tasks](examples/tasks/):
```json
{
  "version": 1,
  "schedule": {
    "type": "simple",
    "interval": "15s"
  },
  "workflow": {
    "collect": {
      "metrics": {
        "/intel/logs/*": {}
      },
      "process": [
        {
          "plugin_name": "logs-regexp",
          "config": {
            "regexp_log": "(?P<client_ip>\\S+) (\\S{1,}) (\\S{1,}) [[](?P<timestamp>\\d{2}[/]\\S+[/]\\d{4}[:]\\d{2}[:]\\d{2}[:]\\d{2} \\S\\d+)[]] (?P<message>.*)",
            "regexp_message": "(?P<http_method>[A-Z]{3,}) (?P<http_url>/\\S*) HTTP/(?P<http_version>\\d+.\\d+)\" (?P<http_status>\\d*) (?P<http_response_size>\\S*) (?P<http_response_time>\\S*)",
            "regexp_time" : "(?P<day>\\d{2})/(?P<month>[a-zA-Z]+)/(?P<year>\\d{4}):(?P<hour>\\d{2}):(?P<minutes>\\d{2}):(?P<seconds>\\d{2}) (?P<timezone>.\\d+)"
          },
          "process": null,
          "publish": [
            {
              "plugin_name": "file",
              "config": {
                "file": "/tmp/published_logs_with_config.log"
              }
            }
          ]
        }
      ]
    }
  }
}
```

Create a task:
```
$ snaptel task create -t task-config.json
```

To stop task:
```
$ snaptel task stop <task_id>
```

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
