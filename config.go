package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type Config struct {
	Test struct {
		Host     string `json:"Host"`
		Password string `json:"Password"`
	} `json:"test"`

	Logfile           string `json:"Logfile"`
	Loglevel          int    `json:"Loglevel"`
	BeeObjectPath     string `json:"BeeObjectPath"`
	PIDFile           string `json:"PIDFile"`
	HTTPport          string `json:"HTTPport"`
	HTTPSport         string `json:"HTTPSport"`
	HTTPScert         string `json:"HTTPScert"`
	HTTPSkey          string `json:"HTTPSkey"`
	HTTPXreadTimeout  int    `json:"HTTPXreadTimeout"`
	HTTPXwriteTimeout int    `json:"HTTPXwriteTimeout"`
	LogMaxSize        int    `json:"LogMaxSize"`
	LogMaxAge         int    `json:"LogMaxAge"`
	LogMaxBackups     int    `json:"LogMaxBackups"`
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding config file : ")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return config
}

func (C Config) String() string {

	v := reflect.ValueOf(C)
	ret := ""
	ret = ret + "\n[Config::\n"
	for i := 0; i < v.NumField(); i++ {
		ret = ret + fmt.Sprintf("%22s", v.Type().Field(i).Name)
		ret = ret + ":"
		ret = ret + fmt.Sprint(v.Field(i).Interface())
		ret = ret + "\n"
	}
	ret = ret + "\n]\n"
	return ret

}
