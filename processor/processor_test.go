// +build small

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
	"testing"
	"time"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProcessor(t *testing.T) {
	processor := New()
	Convey("Create processor", t, func() {
		Convey("So processor should not be nil", func() {
			So(processor, ShouldNotBeNil)
		})
		Convey("So processor should be of type statisticsProcessor", func() {
			So(processor, ShouldHaveSameTypeAs, &Plugin{})
		})
		Convey("processor.GetConfigPolicy should return a config policy", func() {
			configPolicy, err := processor.GetConfigPolicy()
			Convey("So config policy should be a plugin.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, plugin.ConfigPolicy{})
			})
			Convey("So err should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestProcess(t *testing.T) {
	Convey("Test processing of metrics with correct configuration", t, func() {
		newPlugin := New()
		config := plugin.Config{}
		config[configLogRegexp] = defaultLogRegexp
		config[configMessageRegexp] = defaultMessageRegexp
		config[configTimeRegexp] = defaultTimeRegexp

		Convey("Testing with exemplary logs from keystone-apache-public-access.log", func() {
			logs := []string{
				`127.0.0.1 - - [07/Dec/2016:06:00:12 -0500] "GET /v3/users/fa2b2986c200431b8119035d4a47d420/projects HTTP/1.1" 200 446 21747 "-" "python-keystoneclient"`,
				`127.0.0.1 - - [07/Dec/2016:09:57:29 -0500] "GET / HTTP/1.0" 300 587 2824 "-" "-"`,
				`127.0.0.1 - - [07/Dec/2016:10:59:19 -0500] "GET / HTTP/1.1" 300 627 2026 "-" "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; FSL 7.5.11.01005)"`,
			}
			mts := make([]plugin.Metric, 0)
			for i := range logs {
				mt := plugin.Metric{
					Namespace: plugin.NewNamespace("intel", "logs", "metric", "log", "message"),
					Timestamp: time.Now(),
					Tags:      make(map[string]string),
					Data:      logs[i]}
				mts = append(mts, mt)
			}

			metrics, err := newPlugin.Process(mts, config)

			So(err, ShouldBeNil)
			So(len(metrics), ShouldEqual, 3)

			So(metrics[0].Timestamp.Format(timeFormat), ShouldEqual, `07 Dec 2016 06:00:12 -0500`)
			So(metrics[0].Data.(string), ShouldEqual, `"GET /v3/users/fa2b2986c200431b8119035d4a47d420/projects HTTP/1.1" 200 446 21747 "-" "python-keystoneclient"`)
			So(metrics[0].Tags["client_ip"], ShouldEqual, `127.0.0.1`)
			So(metrics[0].Tags["http_method"], ShouldEqual, `GET`)
			So(metrics[0].Tags["http_url"], ShouldEqual, `/v3/users/fa2b2986c200431b8119035d4a47d420/projects`)
			So(metrics[0].Tags["http_version"], ShouldEqual, `1.1`)
			So(metrics[0].Tags["http_status"], ShouldEqual, `200`)
			So(metrics[0].Tags["http_response_size"], ShouldEqual, `446`)
			So(metrics[0].Tags["http_response_time"], ShouldEqual, `21747`)

			So(metrics[1].Timestamp.Format(timeFormat), ShouldEqual, `07 Dec 2016 09:57:29 -0500`)
			So(metrics[1].Data.(string), ShouldEqual, `"GET / HTTP/1.0" 300 587 2824 "-" "-"`)
			So(metrics[1].Tags["client_ip"], ShouldEqual, `127.0.0.1`)
			So(metrics[1].Tags["http_method"], ShouldEqual, `GET`)
			So(metrics[1].Tags["http_url"], ShouldEqual, `/`)
			So(metrics[1].Tags["http_version"], ShouldEqual, `1.0`)
			So(metrics[1].Tags["http_status"], ShouldEqual, `300`)
			So(metrics[1].Tags["http_response_size"], ShouldEqual, `587`)
			So(metrics[1].Tags["http_response_time"], ShouldEqual, `2824`)

			So(metrics[2].Timestamp.Format(timeFormat), ShouldEqual, `07 Dec 2016 10:59:19 -0500`)
			So(metrics[2].Data.(string), ShouldEqual, `"GET / HTTP/1.1" 300 627 2026 "-" "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; FSL 7.5.11.01005)"`)
			So(metrics[2].Tags["client_ip"], ShouldEqual, `127.0.0.1`)
			So(metrics[2].Tags["http_method"], ShouldEqual, `GET`)
			So(metrics[2].Tags["http_url"], ShouldEqual, `/`)
			So(metrics[2].Tags["http_version"], ShouldEqual, `1.1`)
			So(metrics[2].Tags["http_status"], ShouldEqual, `300`)
			So(metrics[2].Tags["http_response_size"], ShouldEqual, `627`)
			So(metrics[2].Tags["http_response_time"], ShouldEqual, `2026`)

		})

		Convey("Testing with exemplary logs from keystone-apache-admin-access.log", func() {
			logs := []string{
				`127.0.0.1 - - [06/Dec/2016:09:20:06 -0500] "POST /v3/projects HTTP/1.1" 201 271 11516 "-" "python-keystoneclient"`,
				`127.0.0.1 - - [06/Dec/2016:09:20:06 -0500] "GET /v3/projects/61e0ccdef20f409f8d54e80150fbed6d HTTP/1.1" 200 271 4748 "-" "python-keystoneclient"`,
				`127.0.0.1 - - [06/Dec/2016:09:20:06 -0500] "PUT /v3/projects/61e0ccdef20f409f8d54e80150fbed6d/users/a83065ce7d6c4a19a5eee41b6ac169e8/roles/5d74ba2ac7474733b31e13fcdb812458 HTTP/1.1" 204 - 9968 "-" "python-keystoneclient"`,
			}
			mts := make([]plugin.Metric, 0)
			for i := range logs {
				mt := plugin.Metric{
					Namespace: plugin.NewNamespace("intel", "logs", "metric", "log", "message"),
					Timestamp: time.Now(),
					Tags:      make(map[string]string),
					Data:      logs[i]}
				mts = append(mts, mt)
			}

			metrics, err := newPlugin.Process(mts, config)

			So(err, ShouldBeNil)
			So(len(metrics), ShouldEqual, 3)

			So(metrics[0].Timestamp.Format(timeFormat), ShouldEqual, `06 Dec 2016 09:20:06 -0500`)
			So(metrics[0].Data.(string), ShouldEqual, `"POST /v3/projects HTTP/1.1" 201 271 11516 "-" "python-keystoneclient"`)
			So(metrics[0].Tags["client_ip"], ShouldEqual, `127.0.0.1`)
			So(metrics[0].Tags["http_method"], ShouldEqual, `POST`)
			So(metrics[0].Tags["http_url"], ShouldEqual, `/v3/projects`)
			So(metrics[0].Tags["http_version"], ShouldEqual, `1.1`)
			So(metrics[0].Tags["http_status"], ShouldEqual, `201`)
			So(metrics[0].Tags["http_response_size"], ShouldEqual, `271`)
			So(metrics[0].Tags["http_response_time"], ShouldEqual, `11516`)

			So(metrics[1].Timestamp.Format(timeFormat), ShouldEqual, `06 Dec 2016 09:20:06 -0500`)
			So(metrics[1].Data.(string), ShouldEqual, `"GET /v3/projects/61e0ccdef20f409f8d54e80150fbed6d HTTP/1.1" 200 271 4748 "-" "python-keystoneclient"`)
			So(metrics[1].Tags["client_ip"], ShouldEqual, `127.0.0.1`)
			So(metrics[1].Tags["http_method"], ShouldEqual, `GET`)
			So(metrics[1].Tags["http_url"], ShouldEqual, `/v3/projects/61e0ccdef20f409f8d54e80150fbed6d`)
			So(metrics[1].Tags["http_version"], ShouldEqual, `1.1`)
			So(metrics[1].Tags["http_status"], ShouldEqual, `200`)
			So(metrics[1].Tags["http_response_size"], ShouldEqual, `271`)
			So(metrics[1].Tags["http_response_time"], ShouldEqual, `4748`)

			So(metrics[2].Timestamp.Format(timeFormat), ShouldEqual, `06 Dec 2016 09:20:06 -0500`)
			So(metrics[2].Data.(string), ShouldEqual, `"PUT /v3/projects/61e0ccdef20f409f8d54e80150fbed6d/users/a83065ce7d6c4a19a5eee41b6ac169e8/roles/5d74ba2ac7474733b31e13fcdb812458 HTTP/1.1" 204 - 9968 "-" "python-keystoneclient"`)
			So(metrics[2].Tags["client_ip"], ShouldEqual, `127.0.0.1`)
			So(metrics[2].Tags["http_method"], ShouldEqual, `PUT`)
			So(metrics[2].Tags["http_url"], ShouldEqual, `/v3/projects/61e0ccdef20f409f8d54e80150fbed6d/users/a83065ce7d6c4a19a5eee41b6ac169e8/roles/5d74ba2ac7474733b31e13fcdb812458`)
			So(metrics[2].Tags["http_version"], ShouldEqual, `1.1`)
			So(metrics[2].Tags["http_status"], ShouldEqual, `204`)
			So(metrics[2].Tags["http_response_size"], ShouldEqual, `-`)
			So(metrics[2].Tags["http_response_time"], ShouldEqual, `9968`)
		})

		Convey("Testing with exemplary logs from horizon-access.log", func() {
			logs := []string{
				`127.0.0.1 - - [07/Dec/2016:10:59:26 -0500] "GET / HTTP/1.1" 302 296 "-" "-"`,
				`127.0.0.1 - - [07/Dec/2016:10:59:28 -0500] "GET / HTTP/1.0" 302 297 "-" "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; FSL 7.5.11.01005)"`,
				`127.0.0.1 - - [07/Dec/2016:10:59:31 -0500] "GET /wp/wp-content/plugins/allwebmenus-wordpress-menu-plugin/readme.txt HTTP/1.1" 404 6469 "-" "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; FSL 7.5.11.01005)"`,
			}
			mts := make([]plugin.Metric, 0)
			for i := range logs {
				mt := plugin.Metric{
					Namespace: plugin.NewNamespace("intel", "logs", "metric", "log", "message"),
					Timestamp: time.Now(),
					Tags:      make(map[string]string),
					Data:      logs[i]}
				mts = append(mts, mt)
			}

			metrics, err := newPlugin.Process(mts, config)

			So(err, ShouldBeNil)
			So(len(metrics), ShouldEqual, 3)

			So(metrics[0].Timestamp.Format(timeFormat), ShouldEqual, `07 Dec 2016 10:59:26 -0500`)
			So(metrics[0].Data.(string), ShouldEqual, `"GET / HTTP/1.1" 302 296 "-" "-"`)
			So(metrics[0].Tags["client_ip"], ShouldEqual, `127.0.0.1`)
			So(metrics[0].Tags["http_method"], ShouldEqual, `GET`)
			So(metrics[0].Tags["http_url"], ShouldEqual, `/`)
			So(metrics[0].Tags["http_version"], ShouldEqual, `1.1`)
			So(metrics[0].Tags["http_status"], ShouldEqual, `302`)
			So(metrics[0].Tags["http_response_size"], ShouldEqual, `296`)
			So(metrics[0].Tags["http_response_time"], ShouldEqual, `"-"`)

			So(metrics[1].Timestamp.Format(timeFormat), ShouldEqual, `07 Dec 2016 10:59:28 -0500`)
			So(metrics[1].Data.(string), ShouldEqual, `"GET / HTTP/1.0" 302 297 "-" "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; FSL 7.5.11.01005)"`)
			So(metrics[1].Tags["client_ip"], ShouldEqual, `127.0.0.1`)
			So(metrics[1].Tags["http_method"], ShouldEqual, `GET`)
			So(metrics[1].Tags["http_url"], ShouldEqual, `/`)
			So(metrics[1].Tags["http_version"], ShouldEqual, `1.0`)
			So(metrics[1].Tags["http_status"], ShouldEqual, `302`)
			So(metrics[1].Tags["http_response_size"], ShouldEqual, `297`)
			So(metrics[1].Tags["http_response_time"], ShouldEqual, `"-"`)

			So(metrics[2].Timestamp.Format(timeFormat), ShouldEqual, `07 Dec 2016 10:59:31 -0500`)
			So(metrics[2].Data.(string), ShouldEqual, `"GET /wp/wp-content/plugins/allwebmenus-wordpress-menu-plugin/readme.txt HTTP/1.1" 404 6469 "-" "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; FSL 7.5.11.01005)"`)
			So(metrics[2].Tags["client_ip"], ShouldEqual, `127.0.0.1`)
			So(metrics[2].Tags["http_method"], ShouldEqual, `GET`)
			So(metrics[2].Tags["http_url"], ShouldEqual, `/wp/wp-content/plugins/allwebmenus-wordpress-menu-plugin/readme.txt`)
			So(metrics[2].Tags["http_version"], ShouldEqual, `1.1`)
			So(metrics[2].Tags["http_status"], ShouldEqual, `404`)
			So(metrics[2].Tags["http_response_size"], ShouldEqual, `6469`)
			So(metrics[2].Tags["http_response_time"], ShouldEqual, `"-"`)
		})

		Convey("Testing without fields in regexp", func() {
			newPlugin := New()

			config[configMessageRegexp] = defaultMessageRegexp
			config[configTimeRegexp] = defaultTimeRegexp

			mts := []plugin.Metric{
				plugin.Metric{
					Namespace: plugin.NewNamespace("intel", "logs", "metric", "log", "message"),
					Timestamp: time.Now(),
					Data:      `127.0.0.1 - - [07/Dec/2016:06:00:12 -0500] "GET /v3/users/fa2b2986c200431b8119035d4a47d420/projects HTTP/1.1" 200 446 21747 "-" "python-keystoneclient"`},
			}

			Convey("Missing message field in log regexp", func() {
				config[configLogRegexp] = `(?P<client_ip>\S+) (\S{1,}) (\S{1,}) [[](?P<timestamp>\d{2}[/]\S+[/]\d{4}[:]\d{2}[:]\d{2}[:]\d{2} \S\d+)[]] (?P<info>.*)`
				metrics, err := newPlugin.Process(mts, config)
				So(err, ShouldBeNil)
				So(len(metrics), ShouldEqual, 1)
				So(metrics[0].Data, ShouldEqual, `127.0.0.1 - - [07/Dec/2016:06:00:12 -0500] "GET /v3/users/fa2b2986c200431b8119035d4a47d420/projects HTTP/1.1" 200 446 21747 "-" "python-keystoneclient"`)
			})

			Convey("Missing timestamp field in log regexp", func() {
				config[configLogRegexp] = `(?P<client_ip>\S+) (\S{1,}) (\S{1,}) [[](?P<time>\d{2}[/]\S+[/]\d{4}[:]\d{2}[:]\d{2}[:]\d{2} \S\d+)[]] (?P<message>.*)`
				metrics, err := newPlugin.Process(mts, config)
				So(err, ShouldBeNil)
				So(len(metrics), ShouldEqual, 1)
				So(metrics[0].Data, ShouldEqual, `"GET /v3/users/fa2b2986c200431b8119035d4a47d420/projects HTTP/1.1" 200 446 21747 "-" "python-keystoneclient"`)
			})
		})
	})

	Convey("Test processing of metrics with errors", t, func() {
		newPlugin := New()

		mts := []plugin.Metric{
			plugin.Metric{
				Namespace: plugin.NewNamespace("intel", "logs", "metric", "log", "message"),
				Timestamp: time.Now(),
				Tags:      make(map[string]string),
				Data:      `127.0.0.1 - - [07/Dec/2016:06:00:12 -0500] "GET /v3/users/fa2b2986c200431b8119035d4a47d420/projects HTTP/1.1" 200 446 21747 "-" "python-keystoneclient"`},
		}

		Convey("Testing with missing configurable parameter", func() {
			Convey("Missing log regexp", func() {
				config := plugin.Config{}
				config[configMessageRegexp] = defaultMessageRegexp
				config[configTimeRegexp] = defaultTimeRegexp
				metrics, err := newPlugin.Process(mts, config)
				So(err, ShouldNotBeNil)
				So(metrics, ShouldBeNil)
			})

			Convey("Missing message regexp", func() {
				config := plugin.Config{}
				config[configLogRegexp] = defaultLogRegexp
				config[configTimeRegexp] = defaultTimeRegexp
				metrics, err := newPlugin.Process(mts, config)
				So(err, ShouldNotBeNil)
				So(metrics, ShouldBeNil)
			})

			Convey("Missing time regexp", func() {
				config := plugin.Config{}
				config[configLogRegexp] = defaultLogRegexp
				config[configMessageRegexp] = defaultMessageRegexp
				metrics, err := newPlugin.Process(mts, config)
				So(err, ShouldNotBeNil)
				So(metrics, ShouldBeNil)
			})
		})

		Convey("Testing with wrong configurable parameter", func() {
			config := plugin.Config{}
			config[configLogRegexp] = defaultLogRegexp
			config[configMessageRegexp] = defaultMessageRegexp
			config[configTimeRegexp] = defaultTimeRegexp

			Convey("Missing day label in time regexp", func() {
				config[configTimeRegexp] = `(?P<day123>\d{2})/(?P<month>[a-zA-Z]+)/(?P<year>\d{4}):(?P<hour>\d{2}):(?P<minutes>\d{2}):(?P<seconds>\d{2}) (?P<timezone>.\d+)`
				metrics, err := newPlugin.Process(mts, config)
				So(err, ShouldBeNil)
				So(len(metrics), ShouldEqual, 1)
				So(metrics[0].Data, ShouldEqual, `127.0.0.1 - - [07/Dec/2016:06:00:12 -0500] "GET /v3/users/fa2b2986c200431b8119035d4a47d420/projects HTTP/1.1" 200 446 21747 "-" "python-keystoneclient"`)
			})

			Convey("Wrong time regexp", func() {
				config[configTimeRegexp] = `\d(+`
				metrics, err := newPlugin.Process(mts, config)
				So(err, ShouldNotBeNil)
				So(metrics, ShouldBeNil)
			})

			Convey("Wrong message regexp", func() {
				config[configMessageRegexp] = `\d(+`
				metrics, err := newPlugin.Process(mts, config)
				So(err, ShouldNotBeNil)
				So(metrics, ShouldBeNil)
			})

			Convey("Wrong log regexp", func() {
				config[configLogRegexp] = `\d(+`
				metrics, err := newPlugin.Process(mts, config)
				So(err, ShouldNotBeNil)
				So(metrics, ShouldBeNil)
			})
		})
	})

	Convey("Test processing of metrics with unexpected log message", t, func() {
		newPlugin := New()
		config := plugin.Config{}
		config[configLogRegexp] = defaultLogRegexp
		config[configMessageRegexp] = defaultMessageRegexp
		config[configTimeRegexp] = defaultTimeRegexp

		Convey("Unexpected time format", func() {
			mts := []plugin.Metric{
				plugin.Metric{
					Namespace: plugin.NewNamespace("intel", "logs", "metric", "log", "message"),
					Timestamp: time.Now(),
					Tags:      make(map[string]string),
					Data:      `127.0.0.1 - - [07/WrongMonth/2016:06:00:12 -0500] "GET /v3/users/fa2b2986c200431b8119035d4a47d420/projects HTTP/1.1" 200 446 21747 "-" "python-keystoneclient"`},
			}
			metrics, err := newPlugin.Process(mts, config)
			So(err, ShouldBeNil)
			So(len(metrics), ShouldEqual, 1)
			So(metrics[0].Data, ShouldEqual, `127.0.0.1 - - [07/WrongMonth/2016:06:00:12 -0500] "GET /v3/users/fa2b2986c200431b8119035d4a47d420/projects HTTP/1.1" 200 446 21747 "-" "python-keystoneclient"`)
		})

		Convey("Unexpected metric data type", func() {
			mts := []plugin.Metric{
				plugin.Metric{
					Namespace: plugin.NewNamespace("intel", "logs", "metric", "log", "message"),
					Timestamp: time.Now(),
					Tags:      make(map[string]string),
					Data:      123},
			}
			metrics, err := newPlugin.Process(mts, config)
			So(err, ShouldBeNil)
			So(len(metrics), ShouldEqual, 1)
			So(metrics[0].Data, ShouldEqual, 123)
		})
	})
}
