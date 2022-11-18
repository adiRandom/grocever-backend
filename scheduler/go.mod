module scheduler

go 1.19

require (
	github.com/go-co-op/gocron v1.18.0
	github.com/joho/godotenv v1.4.0
	github.com/rabbitmq/amqp091-go v1.5.0
	gorm.io/gorm v1.24.1
	lib v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	gorm.io/driver/mysql v1.4.4 // indirect
)

replace lib => ../lib
