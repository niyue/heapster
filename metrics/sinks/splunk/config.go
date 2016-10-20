// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package splunk
import (
    "net"
    "net/url"
    "strconv" 
    "strings"
    
    "github.com/golang/glog"
)

type Config struct {
    Scheme string
    Host string
    Port int
    Path string
    HecToken string
    Index string
    SourceType string
    Source string
    TestRun bool
}

func NewConfig(uri *url.URL) (Config, error) {
    config := Config{}
    config.Scheme = uri.Scheme
    host, port, _ := net.SplitHostPort(uri.Host)
    config.Host = host
    portValue, _ := strconv.ParseInt(port, 10, 0)
    config.Port = int(portValue)
    config.Path = uri.Path

    opts, _ := url.ParseQuery(uri.RawQuery)
    
    glog.V(2).Info("Splunk sink config is created. uri=", uri)
    
    if len(opts["hec-token"]) > 0 {
        config.HecToken = opts["hec-token"][0]
    }
    
    if len(opts["index"]) > 0 {
        config.Index = opts["index"][0]
    } else {
        config.Index = "heapster_metrics"
    }
    if len(opts["sourcetype"]) > 0 {
        config.SourceType = opts["sourcetype"][0]
    }
    if len(opts["source"]) > 0 {
        config.Source = opts["source"][0]
    }
    if len(opts["testrun"]) > 0 {
        config.TestRun = strings.Compare(opts["testrun"][0], "true") == 0
    }
    glog.V(2).Info("Splunk sink config is parsed. config=", config)
    return config, nil
}