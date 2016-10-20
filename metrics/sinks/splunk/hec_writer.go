package splunk

import (
    "fmt"
    "crypto/tls"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "os"
    "strings"
    
    "github.com/golang/glog"
)

type HecWriter struct {
    config Config
    client *http.Client
}

func NewHecWriter(config *Config) *HecWriter {
    tr := &http.Transport {
        TLSClientConfig: &tls.Config{ InsecureSkipVerify: true },
    }
    client := &http.Client{ Transport: tr }
    return &HecWriter {
        config: *config,
        client: client,
    }
}

func toJson(event *event, config *Config) *[]byte {
    hostname, _ := os.Hostname()
    eventMap := make(map[string]string)
    for _, attribute := range (*event).attributes {
        eventMap[attribute.key] = fmt.Sprintf("%v", attribute.value)
    }
    hecEvent := map[string]interface{} {
        "index": (*config).Index,
        "sourcetype": (*config).SourceType,
        "source": (*config).Source,
        "time": (*event).timestamp.UnixNano() / 1000000,
        "host": hostname,
        "event": eventMap,
    }
    json, _ := json.Marshal(hecEvent)
    return &json
}

func (this *HecWriter) write(events *[]event) {
    hecEndpointUrl := fmt.Sprintf("%v://%v:%v%v", 
        this.config.Scheme, 
        this.config.Host, 
        this.config.Port, 
        this.config.Path)
    for _, event := range *events {
        eventBytes := toJson(&event, &this.config)
        req, requestError := http.NewRequest("POST", hecEndpointUrl, 
            strings.NewReader(string(*eventBytes)))
        glog.V(2).Info("hec_token=", this.config.HecToken)
        hecAuth := fmt.Sprintf("Splunk %v", this.config.HecToken)
        if requestError != nil {
            glog.Errorf("Failed to create HEC request. url=%s error=%s", 
                hecEndpointUrl, requestError)  
        } else {
            req.Header.Set("Authorization", hecAuth)
            req.Header.Set("Content-Type", "application/json")

            resp, err := this.client.Do(req)
            if err != nil {
                if !this.config.TestRun {
                    glog.Errorf("Splunk HEC endpoint reports an error. error='%s'", err)
                }
            } else {
                defer resp.Body.Close()
                responseBody, _ := ioutil.ReadAll(resp.Body)    
                glog.V(2).Infof("Succeeded to write one event via Splunk HEC, response=%s", responseBody)    
            }    
        }
    }
}
