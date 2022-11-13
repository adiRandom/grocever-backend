package scheduling

import "lib/data/dto"

type CrawlDto struct {
	CrawlSource dto.CrawlSourceDto `json:"crawlSource"`
	Type        string             `json:"type"`
}

const NORMAL = "normal"
const PRIORITIZED = "prioritized"
