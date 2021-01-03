package jenkins

import (
	"bytes"
	"testing"
)

func TestNewJenkinsClient(t *testing.T) {
	c := newJenkinsClient(&Config{})
	if c == nil {
		t.Errorf("Expected populated client")
	}

	c = newJenkinsClient(&Config{
		CACert: bytes.NewBufferString("certificate"),
	})
	if string(c.Requester.CACert) != "certificate" {
		t.Errorf("Initialization did not extract certificate data")
	}
}
