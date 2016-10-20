package splunk

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestToJson(t *testing.T) {
    eventPointer := &event {
        attributes: []eventAttribute {
            eventAttribute{ key: "name", value: "foo"},
            eventAttribute{ key: "pod", value: "bar"},
        },
    }
    configPointer := &Config {
        Index: "main",
        SourceType: "heapster-metrics",
        Source: "kubernetes",
    }
    eventJson := toJson(eventPointer, configPointer)
    assert.NotNil(t, eventJson)
}

func TestToJsonWithEmptyAttributesAndSourceSourceType(t *testing.T) {
    eventPointer := &event {
    }
    configPointer := &Config {
    }
    eventJson := toJson(eventPointer, configPointer)
    assert.NotNil(t, eventJson)
}