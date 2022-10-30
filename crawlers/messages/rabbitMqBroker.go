package messages

import (
	"dealScraper/lib/data/dto"
	"dealScraper/lib/messages"
	"dealScraper/lib/network"
	"time"
)

var rabbitMqBroker *messages.RabbitMqJsonMultiplexBroker[dto.SearchProductDto]
var messageProcessingTimeout = 1 * time.Minute
var inboundQueues = map[string]messages.AmqpJsonMultiplexBrokerInboundQueueMetadata{
	network.PriorityCrawlQueue: {
		QueueName:      network.PriorityCrawlQueue,
		ProcessTimeout: &messageProcessingTimeout,
	},
	network.CrawlQueue: {
		QueueName:      network.CrawlQueue,
		ProcessTimeout: &messageProcessingTimeout,
	},
}

const queueSwitchInterval = 10

func PickInboundQueue(currentQueueName string, queueMetadata messages.AmqpJsonMultiplexBrokerSelectQueueMetadataMap[dto.SearchProductDto]) string {
	if currentQueueName == "" {
		return network.CrawlQueue
	}

	if currentQueueName == network.PriorityCrawlQueue {
		if queueMetadata[network.PriorityCrawlQueue].MessageCount == 0 {
			return network.CrawlQueue
		}

		if queueMetadata[network.PriorityCrawlQueue].DeltaProcessedCount >= queueSwitchInterval {
			return network.CrawlQueue
		}
	}

	if currentQueueName == network.CrawlQueue {
		if queueMetadata[network.PriorityCrawlQueue].MessageCount == 0 {
			return network.CrawlQueue
		}

		if queueMetadata[network.CrawlQueue].MessageCount == 0 {
			return network.PriorityCrawlQueue
		}

		if queueMetadata[network.CrawlQueue].DeltaProcessedCount >= queueSwitchInterval {
			return network.PriorityCrawlQueue
		}
	}

	return currentQueueName
}
