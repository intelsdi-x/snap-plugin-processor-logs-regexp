/*
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
*/

package processor

import (
	"fmt"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

const (
	//Name of the plugin
	Name = "logs-regexp"
	//Version of the plugin
	Version = 1

	timeFormat           = "02 Jan 2006 15:04:05 -0700"
	configLogRegexp      = "regexp_log"
	configMessageRegexp  = "regexp_message"
	configTimeRegexp     = "regexp_time"
	defaultLogRegexp     = `(?P<client_ip>\S+) (\S{1,}) (\S{1,}) [[](?P<timestamp>\d{2}[/]\S+[/]\d{4}[:]\d{2}[:]\d{2}[:]\d{2} \S\d+)[]] (?P<message>.*)`
	defaultMessageRegexp = `(?P<http_method>[A-Z]{3,}) (?P<http_url>/\S*) HTTP/(?P<http_version>\d+.\d+)" (?P<http_status>\d*) (?P<http_response_size>\S*) (?P<http_response_time>\S*)`
	defaultTimeRegexp    = `(?P<day>\d{2})/(?P<month>[a-zA-Z]+)/(?P<year>\d{4}):(?P<hour>\d{2}):(?P<minutes>\d{2}):(?P<seconds>\d{2}) (?P<timezone>.\d+)`
)

var timeLabels = []string{"day", "month", "year", "hour", "seconds", "minutes", "seconds", "timezone"}

type Plugin struct {
}

// New() returns a new instance of the plugin
func New() *Plugin {
	p := &Plugin{}
	return p
}

// GetConfigPolicy returns the config policy
func (p *Plugin) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()
	policy.AddNewStringRule([]string{""}, configLogRegexp, false, plugin.SetDefaultString(defaultLogRegexp))
	policy.AddNewStringRule([]string{""}, configMessageRegexp, false, plugin.SetDefaultString(defaultMessageRegexp))
	policy.AddNewStringRule([]string{""}, configTimeRegexp, false, plugin.SetDefaultString(defaultTimeRegexp))
	return *policy, nil
}

// Process processes the data
func (p *Plugin) Process(metrics []plugin.Metric, cfg plugin.Config) ([]plugin.Metric, error) {
	logRgx, msgRgx, timeRgx, err := getCheckConfig(cfg)
	if err != nil {
		return nil, err
	}

	for i, m := range metrics {
		logBlock, ok := m.Data.(string)
		if !ok {
			warnFields := map[string]interface{}{
				"namespace": m.Namespace.Strings(),
				"data":      m.Data,
			}
			log.WithFields(warnFields).Warn("unexpected data type, plugin processes only strings")
			continue
		}
		fields := make(map[string]string, 0)
		err := parse(logBlock, logRgx, fields)
		if err != nil {
			warnFields := map[string]interface{}{
				"namespace":     m.Namespace.Strings(),
				"data":          m.Data,
				configLogRegexp: logRgx,
			}
			log.WithFields(warnFields).Warn(err)
			continue
		}

		timeStr, ok := fields["timestamp"]
		if ok {
			timestamp, err := parseTime(timeStr, timeRgx)
			if err != nil {
				warnFields := map[string]interface{}{
					"namespace":      m.Namespace.Strings(),
					"data":           m.Data,
					"time":           timeStr,
					configTimeRegexp: timeRgx,
				}
				log.WithFields(warnFields).Warn(err)
				continue
			}
			//replace metric timestamp with timestamp from log
			metrics[i].Timestamp = timestamp
			delete(fields, "timestamp")
		}

		msg, ok := fields["message"]
		if ok {
			err = parse(msg, msgRgx, fields)
			if err != nil {
				warnFields := map[string]interface{}{
					"namespace":         m.Namespace.Strings(),
					"data":              m.Data,
					"message":           msg,
					configMessageRegexp: msgRgx,
				}
				log.WithFields(warnFields).Warn(err)
				continue
			}

			//replace metric data with main message from log
			metrics[i].Data = msg
			delete(fields, "message")
		}

		if len(fields) != 0 && metrics[i].Tags == nil {
			metrics[i].Tags = make(map[string]string)
		}

		//tags with an empty string as a key are not added
		delete(fields, "")

		//add tags from log
		for k, v := range fields {
			metrics[i].Tags[k] = v
		}
	}
	return metrics, nil
}

func getCheckConfigVar(cfg plugin.Config, cfgVarName string) (*regexp.Regexp, error) {
	expr, err := cfg.GetString(cfgVarName)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", cfgVarName, err)
	}
	rgx, err := regexp.Compile(expr)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", cfgVarName, err)
	}
	return rgx, nil

}
func getCheckConfig(cfg plugin.Config) (*regexp.Regexp, *regexp.Regexp, *regexp.Regexp, error) {
	logRgx, err := getCheckConfigVar(cfg, configLogRegexp)
	if err != nil {
		return nil, nil, nil, err
	}

	msgRgx, err := getCheckConfigVar(cfg, configMessageRegexp)
	if err != nil {
		return nil, nil, nil, err
	}

	timeRgx, err := getCheckConfigVar(cfg, configTimeRegexp)
	if err != nil {
		return nil, nil, nil, err
	}
	return logRgx, msgRgx, timeRgx, nil
}

func parse(message string, rgx *regexp.Regexp, fields map[string]string) error {
	match := rgx.FindStringSubmatch(message)
	for i, name := range rgx.SubexpNames() {
		if i > 0 && i <= len(match) {
			fields[name] = match[i]
		}
	}
	return nil
}

func parseTime(timeStr string, timeRgx *regexp.Regexp) (time.Time, error) {
	var timeStamp time.Time
	timeFields := make(map[string]string)

	//parse time string
	err := parse(timeStr, timeRgx, timeFields)
	if err != nil {
		return timeStamp, err
	}

	//validate required time fields
	for _, label := range timeLabels {
		_, ok := timeFields[label]
		if !ok {
			return timeStamp, fmt.Errorf("cannot parse log timestamp, missing required time label: %s", label)
		}
	}

	//format a new time string
	s := fmt.Sprintf("%s %s %s %s:%s:%s %s",
		timeFields["day"], timeFields["month"], timeFields["year"],
		timeFields["hour"], timeFields["minutes"], timeFields["seconds"],
		timeFields["timezone"])

	//parse to time.Time
	timeStamp, err = time.Parse(timeFormat, s)
	if err != nil {
		return timeStamp, err
	}
	return timeStamp, nil
}
