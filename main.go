package main

import (
	"log"
	"net/http"
	"time"

	"github.com/pamelag/pika/cab"
	"github.com/pamelag/pika/insight"
	"github.com/pamelag/pika/mysql"
	"github.com/pamelag/pika/server"
)

func main() {
	// Setup repository
	var (
		trips cab.TripRepository
	)

	// get the connection
	connection, err := getConnection()
	if err != nil {
		panic(err)
	}

	// inject the connection into the repository
	trips = mysql.NewTripRepository(connection)

	var is insight.Service
	// inject the repository into the service
	is = insight.NewService(trips)
	// decorate the service with logging
	is = insight.NewLoggingService(is)

	// Create the server with the services
	srv := server.New(is)

	httpServer := &http.Server{Addr: httpAddr,
		Handler:      srv,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second}

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Printf("HTTP server ListenAndServe: %v", err)
	}

}
