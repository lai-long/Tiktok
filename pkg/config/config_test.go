package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	cfg, err := Load([]string{"./"})
	if err != nil {
		t.Fatal(err)
		return
	}
	if cfg == nil {
		t.Fatal("cfg is nil")
		return
	}
}
