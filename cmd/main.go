package main

import (
	"daemon/app/services"
	"daemon/persistence"
	"daemon/presentation"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	serveOnPort := getEnv("SERVE_ON_PORT", "8000")

	router := gin.Default()

	URLRepository := persistence.NewURLRepository()
	URLService := services.NewURLService(URLRepository)

	go func() {
		for {
			URLService.CheckUrls()
			time.Sleep(60 * time.Second)
		}
	}()

	URLroutes := presentation.CreateRoutes(URLService)
	for routePath, routeMethod := range URLroutes {
		if routeMethod.HandlerFunc != nil {
			router.Handle(routeMethod.Method, routePath, routeMethod.HandlerFunc)
		}
	}

	// Start the HTTP server
	err := http.ListenAndServe(":"+serveOnPort, router)
	if err != nil {
		log.Fatal("Error starting HTTP server: ", err)
	}
}
