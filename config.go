package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/mcuadros/go-defaults"
)

type Wallet struct {
	Address string `toml:"address"`
	Name    string `toml:"name"`
	Group   string `toml:"group"`
}

type Chain struct {
	Name        string   `toml:"name"`
	LCDEndpoint string   `toml:"lcd-endpoint"`
	Wallets     []Wallet `toml:"wallets"`
}

func (w Wallet) Validate() error {
	if w.Address == "" {
		return fmt.Errorf("address for wallet is not specified")
	}

	return nil
}

func (c *Chain) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("empty chain name")
	}

	if c.LCDEndpoint == "" {
		return fmt.Errorf("no LCD endpoint provided")
	}

	if len(c.Wallets) == 0 {
		return fmt.Errorf("no wallets provided")
	}

	for _, wallet := range c.Wallets {
		if err := wallet.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type Config struct {
	LogConfig LogConfig `toml:"log"`
	Chains    []Chain   `toml:"chains"`
}

type LogConfig struct {
	LogLevel   string `toml:"level" default:"info"`
	JSONOutput bool   `toml:"json" default:"false"`
}

func (c *Config) Validate() error {
	if len(c.Chains) == 0 {
		return fmt.Errorf("no chains provided")
	}

	for index, chain := range c.Chains {
		if err := chain.Validate(); err != nil {
			return fmt.Errorf("error in chain %d: %s", index, err)
		}
	}

	return nil
}

func GetConfig(path string) (*Config, error) {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	configString := string(configBytes)

	configStruct := Config{}
	if _, err = toml.Decode(configString, &configStruct); err != nil {
		return nil, err
	}

	defaults.SetDefaults(&configStruct)
	return &configStruct, nil
}
