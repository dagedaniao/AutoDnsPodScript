package main

import (
	"fmt"
	"log"
	"os"
        "time"
	"errors"
)

const TomorrowToContinueStr = "Need tomorrow to scale!"

func init() {
	file := "./" + "message" + ".log"
        logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
        if err != nil {
        	panic(err)
        }
        log.SetOutput(logFile)
        log.SetPrefix("[qSkipTool]")
        log.SetFlags(log.LstdFlags | log.Lshortfile | log.Ldate)
	return 
}

func runOneTime(isScaleIn bool, config Config, rules RuleCnames) (err error) {
        fmt.Printf("Step3: Get Scale style: scale-to-targetCdn[%v] scale-immidately[%v] \n", isScaleIn, !config.Timely)
        log.Printf("Step3: Get Scale style: scale-to-targetCdn[%v] scale-immidately[%v] \n", isScaleIn, !config.Timely)

	hour, min, sec := config.StartTime.Hour, config.StartTime.Min, config.StartTime.Sec
        if !isScaleIn {
		hour, min, sec = config.EndTime.Hour, config.EndTime.Min, config.EndTime.Sec
	}
	if config.Timely {
		durtionT, err_ := GetDurtionTime(hour, min, sec)        
    		if err_ != nil {
			panic("GetDurtionTime failed,please check time format!!")
		}
		if durtionT < 0 {
			return
		}
        	fmt.Printf("wait time[%+v] \n", durtionT)
        	log.Printf("wait time[%+v] \n", durtionT)
		time.Sleep(durtionT)
	}

        fmt.Printf("Step4: Handle scale: partly-scale[%v], skip-time[%v]",config.Ps, config.TimeRate)
	fmt.Println("Rules:", rules)
        log.Printf("Step4: Handle scale: partly-scale[%v], skip-time[%v]",config.Ps, config.TimeRate)
	log.Println("Rules:", rules)
        err = HandleScale(config.Ps, isScaleIn, rules.Cnames, config.TimeRate, config.ApiToken)
        if err != nil {
		fmt.Printf("some domains scaled failed! err:[%v]",err)
		log.Printf("some domains scaled failed! err:[%v]",err)
        }
	return
}

func main() {
	fmt.Println("Tool started!")
	log.Println("Tool started!")
        confName := "service.conf"
        ruleName := "rule.conf"

        config, err := GetServiceConfig(confName)
        if err != nil {
		fmt.Println("Read file[%v] err[%v]!",confName, err)
		log.Println("Read file[%v] err[%v]!",confName, err)
		return
	}
        fmt.Printf("Step1: Load Common config[%+v] Success! \n", confName)
        log.Printf("Step1: Load Common config[%+v] Success! \n", confName)
        log.Printf("%+v \n", config)

        if !isApiTokenValid(config.ApiToken) {
		err = errors.New("ApiToken is wrong, the format is 'id,token'!")
		fmt.Println("ApiToken[%v]!", config.ApiToken)
		log.Println("ApiToken[%v]!", config.ApiToken)
		return
		
	}

        rules, err := GetRuleConfig(ruleName)
        if err != nil {
		fmt.Println("Read file[%v] err[%v]!", ruleName, err)
		log.Println("Read file[%v] err[%v]!", ruleName, err)
		return
	}
        fmt.Printf("Step2: Load Rule config[%+v] Success! \n", ruleName)
        log.Printf("Step2: Load Rule config[%+v] Success! \n", ruleName)
        log.Printf("%+v \n", rules)
        
        isScaleIn := false
	if !config.Timely {
		if config.ScaleIn {
			isScaleIn = true
		}
		runOneTime(isScaleIn, config, rules)
	} else { 
		for {
			isScaleIn, err = IsScaleIn(config.StartTime, config.EndTime)
        		if err != nil {
				if err.Error() ==  TomorrowToContinueStr {
					durtionT, err := GetDurtionTime(24, 0, 0)        
    					if err != nil {
						log.Println("GetDurtionTime failed,please check time format!!err:",err)
						fmt.Println("GetDurtionTime failed,please check time format!!err:",err)
					}
					fmt.Printf("wait[%v] for tomorrow to continue", durtionT)
					time.Sleep(durtionT)
					continue
				}
				fmt.Println("Get ScaleIn flag failed!", err)
				log.Println("Get ScaleIn flag failed!", err)
				return
        		}
			runOneTime(isScaleIn, config, rules)
		}
	}
	fmt.Println("Tool completed!")
	log.Println("Tool completed!")
}
