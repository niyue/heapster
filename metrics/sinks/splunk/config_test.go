package splunk

import (
    "testing"
    "net/url"

    "github.com/stretchr/testify/assert"
)

func TestConfigWithToken(t *testing.T) {
    splunkURL, _ := url.Parse("http://localhost:8088/hec?hec-token=foo-token")
    config, _ := NewConfig(splunkURL)
    assert.Equal(t, "foo-token", config.HecToken)
}

func TestConfigSchemeAndHostAndPort(t *testing.T) {
    splunkURL, _ := url.Parse("http://localhost:8088/hec?hec-token=foo-token")
    config, _ := NewConfig(splunkURL)
    assert.Equal(t, "http", config.Scheme)
    assert.Equal(t, "localhost", config.Host)
    assert.Equal(t, 8088, config.Port)
}

func TestConfigWithIndex(t *testing.T) {
    splunkURL, _ := url.Parse("http://localhost:8088/hec?hec-token=foo-token&index=main")
    config, _ := NewConfig(splunkURL)
    assert.Equal(t, "main", config.Index)
}

func TestConfigWithDefaultIndex(t *testing.T) {
    splunkURL, _ := url.Parse("http://localhost:8088/hec?hec-token=foo-token")
    config, _ := NewConfig(splunkURL)
    assert.Equal(t, "heapster_metrics", config.Index)
}

