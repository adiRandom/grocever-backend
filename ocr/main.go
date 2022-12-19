package main

import (
	"context"
	"github.com/joho/godotenv"
	"ocr/gateways/api"
	"ocr/gateways/events"
	"os"
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
	//file, err := os.Open(homeDir + "/Desktop/Reciepts/Mega.jpeg")
	//if err != nil {
	//	panic(err)
	//}
	//defer file.Close()
	//
	//str, err := s.ProcessImage(file)
	//print(str)
	//text := "MEGA IMAGE SRL\nALEEA BUHUSI, NR. 2, SPATIUL COM. 2A, BL. 3\nMUNICIPIUL BUCURESTI, SECTOR 3\nCOD IDENTIFICARE FISCALA: R06719278\n1.000 BUC. X 4.99\n365 LAPTE 1.5% 1L\n1.000 BUC. X 5.09\nVIVA PERN.CACA0200\nCARD\nING\nDET\nTOTAL\nTOTAL TVA\nTVA B 9.00%\n0.0011010\n0020008001\nVANZARE\nCONTACTLESS CIP\nNUMAR BATCH 2734\nDATA 09.11.22\nRRN 231306587658\nRASP 00/000\nJ.08\nLei\n4.99 B\n5.09 B\n10.08\n0.83\n0.83\n10.08\nard\nCOMERCIANT 90352300\nAPN DEBIT MASTERCARD\nSUMA 10.08 RON\nCVMR 3F0002\n************9631\nNUMAR CHITANTA 007357\nORA 08:46:10\n(C1) COD AUT. EHYRKQ\n"
	//s := services.GetParseService()
	//products, err := s.GetOcrProducts(text)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%v", products)
	println(os.Getenv("API_PORT"))

	router := api.GetRouter()
	broker := events.GetRabbitMqBroker()
	ctx := context.Background()
	go broker.Start(ctx)
	router.Run(os.Getenv("API_PORT"))
}
