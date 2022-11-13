package crawl

import (
	"lib/data/dto/scheduling"
	"lib/network/amqp"
)

func GetQueueForPriority(priority string) string {
	if priority == scheduling.Normal {
		return amqp.CrawlQueue
	} else if priority == scheduling.Prioritized {
		return amqp.PriorityCrawlQueue
	}
	return ""
}
