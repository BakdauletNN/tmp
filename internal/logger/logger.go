package main


import (
	"log"
	"os"
)

func main(){
	file, err := os.OpenFile("logger.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(file)

}