package main

import "scheduler/services/crawl"

func main() {
	crawlScheduler := crawl.GetScheduler()
	defer crawlScheduler.Close()
}
