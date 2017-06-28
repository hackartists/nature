package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func YAMLConfig(file string, i interface{}) error {
	c, e := ioutil.ReadFile(file)

	if e != nil {
		return e
	}

	e = yaml.Unmarshal(c, i)

	if e != nil {
		return e
	}

	return nil
}
