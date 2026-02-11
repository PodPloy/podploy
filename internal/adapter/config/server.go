package config

import (
	koanf "github.com/knadh/koanf/v2"
	// file "github.com/knadh/koanf/v2/providers/file"
	// toml "github.com/knadh/koanf/v2/parsers/toml"
)

func LoadConfig(path string) (*koanf.Koanf, error) {
	k := koanf.New(".")

	// if err := k.Load(file.Provider(path), toml.Parser()); err != nil {
	// 	return nil, err
	// }

	return k, nil
}