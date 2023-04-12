package main

import (
	"log"
	"package-service/db"
	"package-service/http"
	"time"
)

func main() {
	for {
		log.Println("Connecting to database...")
		if err := db.Connect(); err != nil {
			log.Printf("Cannot connect to DB: %s\n\n", err)
			time.Sleep(time.Second * 3)
			continue
		}
	
		break
	}

	r := http.InitRouter()

	r.Run()
}
