// thape - Casting container images to gzipped tarballs.
// (c) 2025 Star Inc.

package main

import (
	_ "github.com/joho/godotenv/autoload"

	"compress/gzip"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/star-inc/nui.go"
)

// Config constants for the Thape service.
const (
	configAddress nui.EnvKey = "THAPE_ADDRESS"
)

// setupRouter initializes the Gin router with the necessary routes and handlers.
// It sets up the root route for the welcome message and the image request handler.
//
// The root route provides information on how to use the service, including examples of public and private image requests.
// The image request handler processes requests for container images, pulling them and returning them as gzipped tarballs.
//
// The image request handler supports:
// - Public images in the format: /<image_name>:<tag>
// - Private images with HTTP Basic Auth in the format: <user>:<pass>@localhost/<your_server>/<image_name>:<tag>
func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Thape\n\n"+
			"For downloading a gzipped tarball (.tar.gz) of the container image.\n\n"+
			"Public Image: /<image_name>:<tag>\n"+
			"Example: /alpine:latest\n\n"+
			"Private Image (URL Auth): <user>:<pass>@localhost/<your_server>/<image_name>:<tag>\n"+
			"Example: user:pass@localhost/10.0.0.1/my-image:1.0\n\n"+
			"Optional Query Parameters:\n"+
			"?arch=<architecture>  (e.g., ?arch=linux/arm64)\n"+
			"?name=<custom_name>   (e.g., ?name=my-alpine-service)")
	})

	router.GET("/*imagePath", handleImageRequest)

	return router
}

// handleImageRequest processes the image request, pulling the image and returning it as a gzipped tarball.
// It supports optional HTTP Basic Auth, architecture specification, and custom filenames.
//
// The imagePath parameter should be in the format:
// - For public images: /<image_name>:<tag>
// - For private images: <user>:<pass>@localhost/<your_server>/<image_name>:<tag>
//
// Optional query parameters:
// - ?arch=<architecture>  (e.g., ?arch=linux/arm64)
// - ?name=<custom_name>   (e.g., ?name=my-alpine-service)
//
// Example usage:
// - Public image: /alpine:latest
// - Private image: user:pass@localhost/10.0.0.1/my-image:1.0
func handleImageRequest(c *gin.Context) {
	fullImagePath := strings.TrimPrefix(c.Param("imagePath"), "/")
	if fullImagePath == "" {
		c.String(http.StatusBadRequest, "Bad request: image name is required.")
		return
	}

	imageName := fullImagePath
	var craneOpts []crane.Option

	// Check for HTTP Basic Auth header
	if username, password, ok := c.Request.BasicAuth(); ok {
		log.Printf("Basic auth credentials detected for user: %s", username)
		basicAuth := authn.Basic{
			Username: username,
			Password: password,
		}
		craneOpts = append(craneOpts, crane.WithAuth(&basicAuth))
	}

	// Handle optional platform/architecture
	if archStr := c.Query("arch"); archStr != "" {
		platform, err := v1.ParsePlatform(archStr)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid arch format '%s': %v", archStr, err)
			return
		}
		log.Printf("Requesting specific platform: %s", archStr)
		craneOpts = append(craneOpts, crane.WithPlatform(platform))
	}

	ref, err := name.ParseReference(imageName)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid image name '%s': %v", imageName, err)
		return
	}

	log.Printf("Request received, processing image: %s", ref.Name())

	img, err := crane.Pull(ref.Name(), craneOpts...)
	if err != nil {
		log.Printf("Failed to pull image '%s': %v", ref.Name(), err)
		if strings.Contains(err.Error(), "UNAUTHORIZED") {
			c.String(http.StatusUnauthorized, "Authentication failed (UNAUTHORIZED). Please check your credentials.")
		} else {
			c.String(http.StatusInternalServerError, "Failed to pull image: %v", err)
		}
		return
	}

	// Handle optional custom filename
	fileName := ""
	if customName := c.Query("name"); customName != "" {
		fileName = customName + ".tar.gz"
	} else {
		fileName = strings.Replace(ref.Context().RepositoryStr(), "/", "_", -1) + "_" + ref.Identifier() + ".tar.gz"
	}

	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(fileName))
	c.Header("Content-Type", "application/x-gzip")

	gzipWriter := gzip.NewWriter(c.Writer)
	defer gzipWriter.Close()

	if err := tarball.Write(ref, img, gzipWriter); err != nil {
		log.Printf("Error streaming gzipped tarball to client: %v", err)
	}

	log.Printf("Successfully sent gzipped image: %s", ref.Name())
}

// main initializes the Gin router and starts the server.
func main() {
	// Initialize the router
	router := setupRouter()

	// Construct the listening address
	address := nui.String(configAddress)
	log.Printf("Server starting, listening on http://%s", address)

	// Start the server
	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
