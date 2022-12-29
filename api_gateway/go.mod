module api

go 1.19

require (
	github.com/rabbitmq/amqp091-go v1.5.0
	gorm.io/gorm v1.24.2
	lib v0.0.0-00010101000000-000000000000
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
)

replace lib => ../lib
