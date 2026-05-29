package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("Success_load", func(t *testing.T) {
		cfg, err := Load([]string{"../../config"})
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
	})
	t.Run("Fail_path_not_exist", func(t *testing.T) {
		cfg, err := Load([]string{"/nonexistent/path"})
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
	t.Run("Fail_empty_path", func(t *testing.T) {
		cfg, err := Load([]string{})
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
}
