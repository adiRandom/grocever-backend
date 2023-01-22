package crawl

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"lib/data/dto/crawl"
	"lib/data/dto/scheduling"
	"log"
	"scheduler/data/database/entities"
	repositories "scheduler/data/database/repositories/productRequeue"
	"time"
)

type RequeueService struct {
}

var requeueService *RequeueService = nil

func GetRequeueService() *RequeueService {
	if requeueService == nil {
		requeueService = &RequeueService{}
	}
	return requeueService
}

func (s *RequeueService) Requeue(product crawl.ProductDto) error {
	repository := repositories.GetRepository()
	requeueEntities := entities.NewProductRequeueEntities(product)
	for _, entity := range requeueEntities {
		err := repository.Create(entity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *RequeueService) StartCronRequeue() {
	cron := gocron.NewScheduler(time.UTC)
	_, err := cron.Every(1).Day().At("00:00").Do(s.requeue)
	//_, err := cron.Every(30).Seconds().Do(s.requeue)
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not start crawl requeue cron: %s", err))
	}

	cron.StartAsync()
}

func (s *RequeueService) requeue() {
	repository := repositories.GetRepository()
	products, err := repository.GetProductsForRequeue()
	if err != nil {
		return
	}

	for _, product := range products {
		crawlDto := scheduling.NewRequeueCrawlScheduleDto(
			product.OcrProductName,
			[]crawl.SourceDto{product.CrawlSource.ToDto()},
			scheduling.Normal,
		)
		scheduler.ScheduleCrawl(crawlDto)
	}
}
