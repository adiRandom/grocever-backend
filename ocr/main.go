package main

import (
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	//s := services.GetOcrService()
	//
	//// Open a file from desktop as Reader
	//homeDir, err := os.UserHomeDir()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//file, err := os.Open(homeDir + "/Desktop/Reciepts/3.jpeg")
	//if err != nil {
	//	panic(err)
	//}
	//defer file.Close()
	//
	//s.ProcessImage(file)

	Get

}
