package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConf(t *testing.T) {
	_, err := LoadConf("../etc/conf.toml")
	assert.NoError(t, err)
}
