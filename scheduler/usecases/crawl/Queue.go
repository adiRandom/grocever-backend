package crawl

import (
	"lib/data/dto/scheduling"
	"lib/network/amqp"
)

func GetQueueForPriority(priority string) string {
	if priority == scheduling.NORMAL {
		return amqp.CrawlQueue
	} else if priority == scheduling.PRIORITIZED {
		return amqp.PriorityCrawlQueue
	}
	return ""
}
