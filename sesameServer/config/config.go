package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	MasterPort         int
	SlavePorts         string
	DurationMinutes    int
	Token              string
	ARR_IP_HeaderField string
	URLPostfix         string
}

func PrintConfigExample() {
	fmt.Println(defaultConfig())
}

func defaultConfig() string {
	var conf Config
	conf.MasterPort = 1337
	conf.SlavePorts = `22`
	conf.DurationMinutes = 40
	conf.Token = "secretPhrase"
	conf.ARR_IP_HeaderField = "SomeCustomHeaderField"
	conf.URLPostfix = ""
	confBytes, _ := json.Marshal(conf)
	return string(confBytes)
}

//Try to load config - if not found - create and load again.
//If cannot create or load - print error, help text and exit(1)
func LoadConfig() (*Config, error) {
	var conf Config

	// create default config.json if not exist
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		fmt.Println("config.json is not found. Default config will be created")
		if makeDefaultConfig() != nil {
			fmt.Println("Cannot create default config.json")
			fmt.Println(err.Error())
			return &conf, err
		}
	}

	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Cannot read config.json")
		fmt.Println(err.Error())
		return &conf, err
	}

	if json.Unmarshal(bytes, &conf) != nil {
		fmt.Println("Cannot unmarshal config.json. Check format, or delete config.json. Default config will be created at next run")
		fmt.Println(err.Error())
		return &conf, err
	}

	return &conf, conf.Verify()
}

func makeDefaultConfig() error {
	return ioutil.WriteFile("config.json", []byte(defaultConfig()), 0777)
}

func (t *Config) Verify() error {
	if t.DurationMinutes < 2 || t.DurationMinutes > 43200 {
		return errors.New("Config error: Duration in minutes should be in range 2-43200")
	}

	if t.MasterPort < 0 || t.MasterPort > 65535 {
		return errors.New("Config error: Master port should be in range 0-65535")
	}

	if t.SlavePorts == "" {
		return errors.New("Config error: Slave ports should not be empty. Example (without quoutes) '22,33-51'")
	}

	if t.Token == "" {
		return errors.New("Config error: Token should not be empty.")
	}

	if t.URLPostfix != "" {
		if len(t.URLPostfix) < 16 {
			return errors.New("Config error: URL postfix should has at least 16 characters ([a-z,A-Z,0-9]) or should be empty")
		}
	}

	return nil
}
