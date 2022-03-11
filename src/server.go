package main

import (
	"fmt"
	"log"
	"net/http"
	"pitstop-api/src/config"
	_ "pitstop-api/src/database"
)

func main() {
	r := initRoutes()

	fmt.Println("Démarrage du serveur sur le port", config.ServerPort, "...")
	if err := http.ListenAndServe(":"+config.ServerPort, r); err != nil {
		log.Fatal(err)
	}
}
