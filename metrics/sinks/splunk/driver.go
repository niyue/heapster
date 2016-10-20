// Copyright 2016 Google Inc. All Rights Reserved.
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
    "net/url"
    "time"
    "sync"
    
    "github.com/golang/glog"
    "k8s.io/heapster/metrics/core"
)

type splunkSink struct {
    config Config
    sync.RWMutex
}

type eventAttribute struct {
    key string
    value interface {}
}

type event struct {
    attributes []eventAttribute
    timestamp time.Time
}

func (this *splunkSink) Name() string {
    return "Splunk Sink"
}

func (this *splunkSink) Stop() {
    // Do nothing.
}

func (this *splunkSink) ExportData(batch *core.DataBatch) {
    this.Lock()
    defer this.Unlock()
    events := processBatch(batch)
    hecWriter := NewHecWriter(&this.config)
    hecWriter.write(events)
}

func NewSplunkSink(uri *url.URL) (core.DataSink, error) {
    config, err := NewConfig(uri)
    if (err != nil) {
        return nil, err
    }
    sink := splunkSink{
        config: config,
    }

    glog.Info("Splunk sink is created")
    return &sink, nil
} 

func processBatch(batch *core.DataBatch) *[]event {
    glog.Info("Exporting data into Splunk sink metric_sets=", 
        len(batch.MetricSets))
    
    timestamp := batch.Timestamp.UTC()
    events := make([]event, 0) 
    for metricSetName, metricSet := range batch.MetricSets {
        eventsInMetricSetPointer := processMetricSet(metricSetName, metricSet, timestamp)
        events = append(events, *eventsInMetricSetPointer...)
    }
    
    glog.Info("Batch is exported into Splunk sink. events=", len(events))
    return &events
}

func processMetricSet(
    metricSetName string,
    metricSet *core.MetricSet, 
    timestamp time.Time) *[]event {
    glog.V(2).Info("Add new metric set. name=", metricSetName)
    events := make([]event, 0, len(metricSet.MetricValues))
    for metricName, metricValue := range metricSet.MetricValues {
        eventPointer := processMetric(metricName, &metricValue, timestamp, metricSet.Labels)
        events = append(events, *eventPointer) 
    }
    
    for _, labeldMetric := range metricSet.LabeledMetrics {
        eventPointer := processLabeledMetric(&labeldMetric, timestamp, metricSet.Labels) 
        events = append(events, *eventPointer) 
    }
    return &events
}

func processMetric(
    metricName string, 
    metricValue *core.MetricValue,
    timestamp time.Time,
    labelSets ...map[string]string) *event {
    attributes := make([]eventAttribute, 0)
    for _, labels := range labelSets {
        for k, v := range labels {
            attribute := eventAttribute {
                key: k,
                value: v,
            }
            attributes = append(attributes, attribute)
        }    
    }
    
    attributes = append(attributes, eventAttribute { 
        key: metricName,
        value: metricValue.GetValue(),
    })
    glog.V(2).Info("Add a new metric event. attributes=", attributes)
    return &event {
        attributes: attributes,
        timestamp: timestamp,
    }
}

func processLabeledMetric(
    labeledMetric *core.LabeledMetric,
    timestamp time.Time,
    metricSetLabels map[string]string) *event {
    glog.V(2).Info("Add a new labeled metric.")
    return processMetric(
        labeledMetric.Name,
        &labeledMetric.MetricValue,
        timestamp, 
        labeledMetric.Labels,
        metricSetLabels)
}


