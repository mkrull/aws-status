package main

import "log"

const (
	awsStatus = "http://status.aws.amazon.com"
)

func main() {
	var sm statusMap
	sm.notifier = defaultNotifier

	site := get(awsStatus)

	parse(site, &sm)
	sm.Log()

	for _, e := range sm.Notify() {
		log.Println("Error notifying:", e)
	}
}
