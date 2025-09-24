package config

import (
	"gopkg.in/yaml.v3"
	"io"
)

func UnmarshalYAML(r io.Reader, cfg any) error {
	return yaml.NewDecoder(r).Decode(&cfg)
}
