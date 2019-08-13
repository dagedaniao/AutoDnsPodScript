package main 

import (
	"time"
	"errors"
        "log"	
)

func IsHourOk(hour int) bool {
	return hour >= 0 && hour <= 24
}

func IsMinOk(min int) bool {
	return min >= 0 && min < 60
}

func IsSecOk(sec int) bool {
	return sec >= 0 && sec < 60
}

func GetDurtionTime(hour, min, sec int) (duration time.Duration, err error) {
	if !IsHourOk(hour) || !IsMinOk(min) || !IsSecOk(sec) {
		err = errors.New("time value is invalid!")
		log.Println("H:M:S", hour, min, sec)
		return
	} 	
	const baseTimeFormat = "2006-01-02 15:04:05"
	t := time.Now()
        year, month, day := t.Year(), t.Month(), t.Day()
	if hour == 24 || hour == 0 {
		hour = 0
		day += 1
	}
	tm1 := time.Date(year, month, day, hour, min, sec, 0, t.Location())
	duration = tm1.Sub(t)
	log.Println("Now: Target, durtionTime!", t, tm1, duration)
	return
}

func IsScaleIn(start, end TimeFormat) (ok bool, err error) {
	startD, err := GetDurtionTime(start.Hour, start.Min, start.Sec)
	if err != nil {
		log.Println("IsScale: startTime Err!", start, err)
		return
	}
	endD, err := GetDurtionTime(end.Hour, end.Min, end.Sec)
	if err != nil {
		log.Println("IsScale: endTime Err!", end, err)
		return
	}
	if startD < 0 && endD < 0 {
        	err = errors.New("Need tomorrow to scale!")
		return
	}
	if (endD < 0 && startD > 0) || (endD >=0 && startD >= 0 && endD > startD) {
		ok = true
		return
        }
        return 
}
