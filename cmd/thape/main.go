// thape - casting container images to gzipped tarballs.
// (c) 2025 Star Inc.

package main

import (
	_ "github.com/joho/godotenv/autoload"

	"fmt"
	"log"
	"time"

	"github.com/star-inc/thape/config"
	"github.com/star-inc/thape/kernel"
	"github.com/star-inc/thape/routes"
)

func init() {
	fmt.Println("Thape")
	fmt.Println("===")
	fmt.Println("Casting container images to gzipped tarballs.")
	fmt.Println("\n(c) 2025 Star Inc.")
}

func main() {
	// Define routes
	setupRouters := []kernel.SetupRouter{
		routes.SetupRouter,
	}

	// Create httpServer
	httpServer := kernel.NewHTTPd(setupRouters)

	// Run httpServer
	fmt.Println("Start:", time.Now().Format(time.RFC3339))
	addr := fmt.Sprintf("%s:%d", config.HttpHost, config.HttpPort)
	if err := httpServer.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
