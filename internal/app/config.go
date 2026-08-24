package app

import (
	"os"
	"path/filepath"
)

type Config struct {
	DataDir         string
	Address         string
	ShutdownSeconds int
}

func DefaultConfig() Config {
	data := os.Getenv("FIRELINE_DATA_DIR")
	if data == "" {
		data = filepath.Join(".", "data")
	}
	address := os.Getenv("FIRELINE_ADDR")
	if address == "" {
		address = "127.0.0.1:19700"
	}
	return Config{DataDir: data, Address: address, ShutdownSeconds: 5}
}
func (c Config) Valid() bool { return c.DataDir != "" && c.Address != "" && c.ShutdownSeconds > 0 }
