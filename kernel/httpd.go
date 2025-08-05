// thape - casting container images to gzipped tarballs.
// (c) 2025 Star Inc.

package kernel

import (
	"github.com/gin-gonic/gin"
)

type SetupRouter func(e *gin.Engine)

func NewHTTPd(routes []SetupRouter) *gin.Engine {
	// Create Gin engine
	engine := gin.Default()

	// Deploy routes
	for _, deploy := range routes {
		deploy(engine)
	}

	// Return engine
	return engine
}
