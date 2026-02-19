package config

import (
	"os"
	"path/filepath"

	"github.com/awnumar/memguard"
	"github.com/spf13/viper"
)

type JWTTokenKeys struct {
	Public  *memguard.Enclave
	Private *memguard.Enclave
}

func NewJWTConfig() (*JWTTokenKeys, error) {
	pubPath := viper.GetString("jwt.pubPath")
	privPath := viper.GetString("jwt.privPath")

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	pubPath = filepath.Join(wd, pubPath)
	privPath = filepath.Join(wd, privPath)

	privRaw, err := os.ReadFile(privPath)
	if err != nil {
		return nil, err
	}

	pubRaw, err := os.ReadFile(pubPath)
	if err != nil {
		return nil, err
	}

	keys := &JWTTokenKeys{
		Public:  memguard.NewEnclave(pubRaw),
		Private: memguard.NewEnclave(privRaw),
	}

	for i := range privRaw {
		privRaw[i] = 0
	}
	for i := range pubRaw {
		pubRaw[i] = 0
	}

	return keys, nil
}
