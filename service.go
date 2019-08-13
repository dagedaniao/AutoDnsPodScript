package main 

import (
	"log"
	"time"
        "sync"
	"errors"
)

var toCovers = []int{5, 15, 30, 50, 75, 100}
var fromCovers = []int{10, 30, 60, 100}

func handleScaleImmediately(isScaleIn bool, rules []RuleCname, apitoken string) (err error) {
	log.Println("handleScaleImmediately,ScaleIn[%v], rules[%v]!", isScaleIn, rules)
	var wg sync.WaitGroup
	for idx, rule :=  range rules {
		source := rule.From
                target := rule.To
		if !isScaleIn {
 			source = rule.To
			target = rule.From
                } 
		log.Printf("run goroutine idx:%v,ScaleIn[%v], Domain:%v, View:%v, SourceCname:%v ---> TargetCname:%v !\n", idx, isScaleIn, rule.Domain, rule.View, source, target)
 		wg.Add(1)
                go func(apitoken, zone, domain, view, source, target string) {
			defer wg.Done()
		        log.Printf("run goroutine Domain:%v, View:%v, SourceCname:%v ---> TargetCname:%v !\n", domain, view, source, target)
			err_ := UpdateRecordWeight(apitoken, zone, domain, view, source, target, 100)
			if err_ != nil {
				err = errors.New("ScaleCname failed!")
				log.Println("domain | source | target | err\n", domain, source, target, err_)
			}
		}(apitoken, rule.Zone, rule.Domain, rule.View, source, target)
	}
	wg.Wait()
	return 
}

func handleScaleOnTime(isScaleIn bool, rules []RuleCname, skiptime int, apitoken string) (err error) {	
        var wg sync.WaitGroup
	covers := toCovers
        if !isScaleIn {
		covers = fromCovers
        }
        log.Println("handleScaleOnTime: isScaleIn | covers", isScaleIn, covers)
	for idx, rule :=  range rules {
		source := rule.From
                target := rule.To
		if !isScaleIn {
 			source = rule.To
			target = rule.From
                }
 		log.Println("handleScaleOntime: idx | zone | domain | view | source | target!\n", idx, rule.Zone, rule.Domain, rule.View, source, target) 
		wg.Add(1)
                go func(apitoken, zone, domain, view, source, target string, covers []int, skiptime int) {
			defer wg.Done()
		 	scaleCname(apitoken, zone, domain,view, source, target, covers, skiptime)
		}(apitoken, rule.Zone, rule.Domain, rule.View, source, target, covers, skiptime)
	}
	wg.Wait()        
        return
}

func scaleCname(apitoken, zone, domain, view, source, target string, covers []int, skiptime int) (err error)  {
        log.Println("scaleCname: Begin skiptime | covers", skiptime, covers)
	var errs []error
	for _, weight := range covers {
		log.Printf("run goroutine Domain:%v,View:%v, SourceCname:%v|cover[%v] ---> TargetCname:%v|cover[%v] !\n", domain, view, source, weight, target, 100-weight)
		err_ := UpdateRecordWeight(apitoken, zone, domain, view, source, target, weight)
		errs = append(errs, err_)
		time.Sleep(time.Duration(skiptime)*time.Second)
	}
	if len(errs) == len(covers) {
		err = errors.New("ScaleCname failed!")
		log.Println("domain | source | target | covers | err\n", domain, source, target, covers, err)
		return
	}
	return
}


func HandleScale(isPs, isScaleIn bool, rules []RuleCname, time int, apitoken string) (err error) {
        if isPs == false {
       		err = handleScaleImmediately(isScaleIn, rules, apitoken)
 		return
        }
        err = handleScaleOnTime(isScaleIn, rules, time, apitoken)
	return
}
