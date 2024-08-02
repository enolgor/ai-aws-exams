package conf

import (
	"encoding/json"
	"os"
	"path"
)

type Config struct {
	file          string
	LastUsedEmail string `json:"last_used_email"`
}

func Load(appConfigDir string) (*Config, error) {
	config := &Config{
		file: path.Join(appConfigDir, "config.json"),
	}
	if _, err := os.Stat(config.file); err != nil {
		return config, config.Save()
	}
	f, err := os.Open(config.file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	err = dec.Decode(config)
	return config, err
}

func (config *Config) Save() error {
	f, err := os.OpenFile(config.file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	return enc.Encode(config)
}
