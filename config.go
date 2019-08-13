package main

import (
"errors"
"io/ioutil"
"log"
"encoding/json"
)

type TimeFormat struct {
	Hour int `json:"hour"`
	Min  int `json:"min"`
	Sec  int `json:"sec"`
}

type Config struct {
        ApiToken string  `json:"apitoken"`
	Ps bool  `json:"partlyscale"`
	Timely bool  `json:"timely"`
	ScaleIn bool `json:"scalein"`
        StartTime TimeFormat `json:"start-time"`
        EndTime TimeFormat `json:"end-time"`
        TimeRate  int `json:"time-rate"`
}

type RuleCname struct {
        Zone string `json:"zone"`
        Domain string `json:"domain"`
        View string `json:"view"`
	From string `json:"from"`
        To string `json:"to"`
}

type RuleCnames struct {
	Cnames []RuleCname `json:"cnames"`
}

func LoadFile( fileName string) (content []byte, err error) {
	if len(fileName) == 0 {
                err = errors.New("fileName is empty!!")
                log.Printf("fileName is empty", fileName)
        	return 
        }
	content, err = ioutil.ReadFile(fileName)
	if err != nil {
        	log.Fatal(err)
                return 
        }
        return 
}

func GetServiceConfig(file string) (config Config, err error) {
	content, err := LoadFile(file)
        if err != nil {
        	log.Printf("Config file[%s] \nConfig content: %s", file, content, err)
   		return
        }
        err = json.Unmarshal(content, &config)
        if err != nil {
		log.Println("Service config content:", string(content))
            	log.Fatal(err)
            return
        }
        return
}

func GetRuleConfig(file string) (config RuleCnames, err error) {
	content, err := LoadFile(file)
        if err != nil {
        	log.Printf("Config file[%s] \nConfig content: %s", file, content, err)
   		return
        }
        err = json.Unmarshal(content, &config)
        if err != nil {
            log.Fatal(err)
            return
        }
        return
}
