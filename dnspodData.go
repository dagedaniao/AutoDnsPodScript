package main

import (
	"strings"
)

const (
	DnspodHost = "https://dnsapi.cn/"
	ContentType = "application/x-www-form-urlencoded"
	RecordListUri = "Record.List"
	RecordModifyUri = "Record.Modify"
)

type Status struct {
	Code string `json:"code"`
	Message string `json:"message"`
	CreatedAt string `json:"create_at"`
}

type Domain struct {
	Id string `json:"id"`
	Name string `json:"name"`
	PunyCode string `json:"punycode"`
	Grade string `json:"grade"`
	Owner string `json:"owner"`
	ExStatus string `json:"ex_status"`
	Ttl int `json:"ttl"`
	MinTtl int `json:"min_ttl"`
	DnsPodns []string `json:"dnspod_ns"`
	Status string `json:"status"`
}

type Info struct {
	SubDomains string `json:"sub_domains"`
	RecordTotal string `json:"record_total"`
	RecordsNum string `json:"records_num"`
}

type Record struct {
	Id string `json:"id"`
	Ttl string `json:"ttl"`
	Value string `json:"value"`
	Enabled string `json:"enabled"`
	Status string `json:"status"`
	UpdatedOn string `json:"updated_on"`
	Name string `json:"name"`
	Line string `json:"line"`
	LineId string `json:"line_id"`
	Type string `json:"type"`
	Weight int `json:"weight"`
	MonitorStatus string `json:"monitor_status"`
	Remark string `json:"remark"`
	UseAqb string `json:"use_aqb"`
	Mx string `json:"mx"`
}

type RecordInfo struct {
	StatusInfo Status `json:"status"`
	DomainInfo Domain `json:"domain"`
	InfoDetail Info `json:"info"`
	Records []Record `json:"records"`
}

func isApiTokenValid(apiToken string) (ok bool) {
	if len(apiToken) != 0 {
		strs := strings.Split(apiToken,",")
		if len(strs) == 2 {
			ok = true
		}
	}
        return
}

func isZoneValid(zone string) (ok bool) {
	if len(zone) != 0 {
		strs := strings.Split(zone, ".")
		if len(strs) == 2 || len(strs) == 3 {
			ok = true
		}
	}
	return
}

func isDomainValid(domain, zone string) (ok bool) {
	if len(domain) == 0 || len(zone) == 0 {
		return 
	}

	domains := strings.Split(domain, ".")
	zones := strings.Split(zone, ".")
	if len(domains) < len(zones) {
			return 
	}
	for idx := len(zones); idx > 0; idx-- {
		if domains[len(domains) - idx] != zones[len(zones) - idx] {
			return 
		}
	}
	return true
}
