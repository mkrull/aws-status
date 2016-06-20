package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const (
	awsOk = "Service is operating normally"
)

var (
	regions = []string{"NA", "EU", "SA", "AP"}
)

type status struct {
	Region  string
	Service string
	Details string
}

type statusMap struct {
	m        []status
	n        []status
	notifier func(status) error
}

func (s *statusMap) Notify() []error {
	var errs []error
	for _, st := range s.n {
		err := s.notifier(st)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (s *statusMap) Log() {
	for _, s := range s.m {
		log.Printf("%s in %s: %s", s.Service, s.Region, s.Details)
	}
}

func dieOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func get(url string) *bytes.Buffer {
	resp, err := http.Get(url)
	dieOnError(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	dieOnError(err)

	return bytes.NewBuffer(body)
}

func parse(raw *bytes.Buffer, sm *statusMap) {
	doc, err := goquery.NewDocumentFromReader(raw)
	dieOnError(err)

	(*sm).m = []status{}
	for _, r := range regions {
		doc.Find(fmt.Sprintf("div#%s_block.pad8 tbody tr", r)).Each(func(i int, s *goquery.Selection) {
			service := s.Find("td.bb.top.pad8").Text()
			details := s.Find("td.bb.pad8").Not(".top").Text()

			status := status{
				Region:  r,
				Service: service,
				Details: details,
			}

			(*sm).m = append((*sm).m, status)
			if details != awsOk {
				(*sm).n = append((*sm).n, status)
			}
		})
	}
}

func defaultNotifier(s status) error {
	log.Printf("Service %s in region %s not fully operational: %s\n", s.Service, s.Region, s.Details)
	return nil
}
