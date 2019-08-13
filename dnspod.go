package main

import (
	"fmt"
        "bytes"
	"net/http"
	"io/ioutil"
        "log"
	"strings"
	"errors"
	"strconv"	
	"encoding/json"
)

func checkCommonMsg(apiToken, zone, domain, view string) (err error) {
	if !isApiTokenValid(apiToken) {
		err = errors.New("apiToken is invalid, format is [id,token]!")
		return 
	}
	if !isZoneValid(zone) {
		err = errors.New("zone is invalid, need 2 or 3 info, example baidu.com or pub.edu.cn!")
		return
	}
	if !isDomainValid(domain, zone) {
		err = errors.New("domain is invalid, domain need include zone!")
		return
	}
	if len(view) == 0 {
		err = errors.New("view is empty, exmaple for view:默认，华东!")
		return
	}
	return 
}

func checkModifyAdditionMsg(recordId, recordLine, cname string, weight int) (err error) {
	if len(recordId) == 0 {
		err = errors.New("record_id is empty!")
		return
	}
	if len(recordLine) == 0 {
		err = errors.New("record_line is empty!")
		return
	}
	if len(cname) == 0 {
		err = errors.New("cname is empty!")
		return
	}
	if weight < 0 || weight > 100 {
		err = errors.New("weight is out of range, valid range is [0,100]!")
		return
	}
	return
}

func getSubDomain(domainName, zone string) (subdomain string) {
	strs := strings.Split(domainName, ".")
	zones := strings.Split(zone, ".")
	var subStrs []string
	for idx := 0; idx < (len(strs) - len(zones)); idx++ {
		subStrs = append(subStrs, strs[idx])
	}
	subdomain = strings.Join(subStrs, ".")
	return
}

func constructRecordBody(apiToken, domain, subdomain, view, recordId, recordLine, cname string, weight *int) (data *bytes.Buffer) {
        var bodyStr string
	if len(apiToken) != 0 {
		loginToken := "login_token="+apiToken
		bodyStr += loginToken
	}
	format := "format="+"json"
	bodyStr += "&"+format
	if len(domain) != 0 {
		domainStr := "domain="+domain
		bodyStr += "&"+domainStr
	}
	if len(subdomain) != 0 {
		subdomainStr := "sub_domain="+subdomain
		bodyStr += "&"+subdomainStr
	}
	if len(view) != 0{
		viewStr := "record_line="+view	
		bodyStr += "&"+viewStr
	}
	if len(recordId) != 0 {
		recordIdStr := "record_id="+recordId
		bodyStr += "&"+recordIdStr
	}
	if len(recordLine) != 0 {
		recordLineStr := "record_line="+recordLine
		bodyStr += "&"+recordLineStr
	}
	if len(cname) != 0 {
		cnameStr := "value="+cname
		bCname := []byte(cname)
		if bCname[len(bCname)-1] != '.' {
			cnameStr += "."
		}
		bodyStr += "&"+cnameStr
	}
	recordType := "record_type="+"CNAME"
	bodyStr += "&"+recordType
	if weight != nil {
		if *weight == 0 {
			status := "status=disable"
			bodyStr += "&"+status
		} else {
			status := "status=enable"
			bodyStr += "&"+status
			weightStr := "weight="+strconv.FormatInt(int64(*weight), 10)
			bodyStr += "&"+weightStr
		}
	}

	if len(bodyStr) == 0 {
		return
	}
	log.Println("Msg body:",bodyStr)
	body := []byte(bodyStr)
	data = bytes.NewBuffer(body)
	return
}


func GetDomainRecord(apiToken, zone, domainName, view string) (resbody []byte, err error){
	if err = checkCommonMsg(apiToken, zone, domainName, view); err != nil {
		fmt.Println("GetDomainRecord:", err)
		log.Println("GetDomainRecord:", err)
		return
	}
	uri := DnspodHost + RecordListUri
	domain := zone
	subdomain := getSubDomain(domainName, zone)
	body := constructRecordBody(apiToken, domain, subdomain, view, "", "", "", nil)
	if body == nil {
		err = errors.New("RecordList body is empty!")
		return
 	}
	resp, err := http.Post(uri, ContentType, body)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		return
	}
	defer resp.Body.Close()
        
        resbody, _ = ioutil.ReadAll(resp.Body)
	resStr := string(resbody)
	log.Println("GetDomainRecord:resp:\n",resStr)
	return
}

func ModifyDomainRecord(apiToken, zone, domainName, view, recordId, recordLine, cname string, weight int) (err error) {
	if err = checkCommonMsg(apiToken, zone, domainName, view); err != nil {
		fmt.Println("ModifyDomainRecord:",err)
		log.Println("ModifyDomainRecord:",err)
		return
	}
	if err = checkModifyAdditionMsg(recordId, recordLine, cname, weight); err != nil {
		fmt.Println("checkModifyAdditionMsg:",err)
		log.Println("checkModifyAdditionMsg:",err)
		return
	}
	uri := DnspodHost + RecordModifyUri
	domain := zone
	subdomain := getSubDomain(domainName, zone)
	body := constructRecordBody(apiToken, domain, subdomain, view, recordId, recordLine, cname, &weight)
	if body == nil {
		err = errors.New("RecordList body is empty!")
		return
 	}
	resp, err := http.Post(uri, ContentType, body)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	
        resbody, _ := ioutil.ReadAll(resp.Body)
	resStr := string(resbody)
	log.Println("GetDomainRecord:resp:\n",resStr)
	return
} 

func adpterCname(cname string) (str string) {
	str = cname
	if len(cname) > 0 {
		bCname := []byte(cname)
		if bCname[len(bCname)-1] != '.' {
			str += "."
		}
        }
	return
}

func UpdateRecordWeight(apiToken, zone, domainName, view, sourceCname, targetCname string, weight int) (err error) {
	log.Println("UpdateRecordWeight Begin!")
	
	resBody, err := GetDomainRecord(apiToken, zone, domainName, view)
	if err != nil {
		fmt.Println("GetDomainRecord failed!")
		log.Println("GetDomainRecord failed!")
		return
	}
	var record RecordInfo	
	var cnames []string
        weights := make(map[string]int)
	sourceCname = adpterCname(sourceCname)
	targetCname = adpterCname(targetCname)
	cnames = append(cnames, sourceCname)
	cnames = append(cnames, targetCname)
	weights[sourceCname] = 100 - weight
	weights[targetCname] = weight
	err = json.Unmarshal(resBody, &record)
	if err != nil {
		fmt.Println("Json failed!domain | sourceCname | targetCname | weight | err",domainName, sourceCname,targetCname, weight, err)
		log.Println("Json failed!domain | sourceCname | targetCname | weight | err",domainName, sourceCname,targetCname, weight, err)
		return
	}
	log.Println("Json:", record)
	if len(cnames) != len(record.Records) && (view != "" && view != "默认") {
		fmt.Println("Records num != 2!", domainName, len(record.Records))
		log.Println("Records num != 2!", domainName, len(record.Records))
		return
	}
        for _, cname := range cnames {
        	for _, r := range record.Records {
			if r.Value == cname && r.Line == view{
				recordId := r.Id
				recordLine := r.Line
				weight := weights[cname]
				err_ := ModifyDomainRecord(apiToken, zone, domainName, view, recordId, recordLine, cname, weight)
				if err_ != nil {
					fmt.Println("ModifyDomainRecord failed!")
					log.Println("ModifyDomainRecord failed!")
				}
				break
			}
		}
	}

	log.Println("UpdateRecordWeight End!")
	return
}
