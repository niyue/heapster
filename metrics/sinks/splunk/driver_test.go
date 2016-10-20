package splunk

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/heapster/metrics/core"
)

func NewTestSplunkSink() core.DataSink {
	hecEndpointUrl, err := url.Parse(
        "https://localhost:8088/services/collector?index=main&hec-token=00000000-0000-0000-0000-000000000000&testrun=true")
	if err != nil {
		panic(err)
	}
	testSplunkSink, _ := NewSplunkSink(hecEndpointUrl)
	return testSplunkSink
}

func TestGetSinkName(t *testing.T) {
	splunkSink := NewTestSplunkSink()
	assert.Equal(t, "Splunk Sink", splunkSink.Name())
}

func TestStopSink(t *testing.T) {
	splunkSink := NewTestSplunkSink()
	splunkSink.Stop()
}

func TestExportData(t *testing.T) {
	splunkSink := NewTestSplunkSink()
	batch := createDataBatch()
	splunkSink.ExportData(batch)
}

func TestProcessBatch(t *testing.T) {
    batch := createDataBatch()
    eventsPointer := processBatch(batch)
    assert.Equal(t, 4, len(*eventsPointer))
    event := (*eventsPointer)[0]
    assert.Equal(t, 4, len(event.attributes))
    keys := make(map[string]interface {})
    for _, a := range event.attributes {
        keys[a.key] = a.value
    }
    assert.Equal(t, "123", keys["namespace_id"])
    assert.Equal(t, "my_container", keys["container_name"])
    assert.Equal(t, int64(123456), keys["cpu/limit"].(int64))
}

func createLabel() map[string]string {
	l := make(map[string]string)
	l["namespace_id"] = "123"
	l["container_name"] = "my_container"
	l[core.LabelPodId.Key] = "aaaa-bbbb-cccc-dddd"
	return l
}

func createMetricSet() core.MetricSet {
	label := createLabel()
	metricSet := core.MetricSet{
		Labels: label,
		MetricValues: map[string]core.MetricValue{
			"cpu/limit": {
				ValueType:  core.ValueInt64,
				MetricType: core.MetricCumulative,
				IntValue:   123456,
			},
		},
        LabeledMetrics: []core.LabeledMetric {
            core.LabeledMetric {
                Name: "network-tx",
                Labels: map[string]string{
                    "pod": "example",
                    "container_id": "007",
                },
                MetricValue: core.MetricValue {
                    IntValue: 10,
                    ValueType: core.ValueInt64,
                },
            },
        },
	}
	return metricSet
}

func createDataBatch() *core.DataBatch {
	metricSet1 := createMetricSet()

	metricSet2 := createMetricSet()

	timestamp := time.Now()
	data := core.DataBatch{
		Timestamp: timestamp,
		MetricSets: map[string]*core.MetricSet{
			"pod1": &metricSet1,
			"pod2": &metricSet2,
		},
	}
	return &data
}
