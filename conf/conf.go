package conf

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Conf is ...
type Conf struct {
	Secret string `toml:"secret"`
	Token  string `toml:"token"`
}

// LoadConf is ...
func LoadConf(path string) (*Conf, error) {
	confFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf *Conf
	_, err = toml.Decode(string(confFile), &conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
