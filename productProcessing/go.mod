module productProcessing

go 1.19

require (
	github.com/rabbitmq/amqp091-go v1.5.0
	gorm.io/driver/mysql v1.4.4
	gorm.io/gorm v1.24.1
	lib v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
)

replace lib => ../lib
