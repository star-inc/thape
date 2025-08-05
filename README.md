# Thape

Casting container images to gzipped tarballs.

Thape is a lightweight HTTP service that allows you to download container images as gzipped tarballs (.tgz) directly from your browser or any HTTP client.

The service is designed to be simple and efficient, providing an easy way to export container images for backup, distribution, or analysis purposes.

## Prerequisites

Thape is required working with `Go 1.24` or later.

The service is designed to work with:

• Container registries supporting the OCI Distribution Specification
• Docker Hub and other public registries
• Private registries with HTTP Basic Authentication

## Get Started

To build and run Thape:

```bash
# Using Makefile
make
cp .env.sample .env
./build/thape

# Or manually
go build ./cmd/thape
APP_MODE=development HTTP_HOST=localhost HTTP_PORT=8080 ./thape
```

## System Architecture

The service is recommended to be used for light to medium loading tasks.

Thape provides a simple HTTP API for downloading container images as gzipped tarballs. The service pulls images using the go-containerregistry library and streams them directly to the client as compressed archives.

The service supports:
- Public container images from any OCI-compliant registry
- Private images with HTTP Basic Authentication
- Multi-architecture image selection
- Custom filename specification

## Configuration

The service uses environment variables for configuration:

```bash
APP_MODE=development
HTTP_HOST=localhost
HTTP_PORT=8080
```

The configuration is managed using the `nui.go` package, which reads environment variables and provides a clean interface for configuration management.

You can copy the `.env.sample` file to `.env` and modify it according to your needs:

```bash
cp .env.sample .env
```

**Note:** All environment variables are required. The service will not start without them.

## Dependencies

Install the Go module dependencies:

```bash
go mod download
```

## Development Environment

For development with hot-reload capabilities:

```bash
# Using Makefile (recommended)
make dev

# Or manually install air
go install github.com/air-verse/air@latest
air
```

## Production Environment

Build and run the service for production:

```bash
# Using Makefile
make
./build/thape

# Or manually
go build -o thape ./cmd/thape
./thape
```

For containerized deployment:

```bash
docker build -t thape .
docker run -p 8080:8080 thape
```

## API Documentation

### GET /

> Service information and usage instructions

Returns information about how to use the Thape service, including examples for public and private images.

### GET /{image_name}:{tag}

> Download container image as gzipped tarball

Downloads the specified container image as a gzipped tarball (.tgz).

**Parameters:**
- `image_name`: The name of the container image
- `tag`: The image tag or digest

**Query Parameters:**
- `arch`: Architecture specification (e.g., `linux/arm64`)
- `name`: Custom filename for the download

**Examples:**

Public image:
```
GET /alpine:latest
```

Public image with custom filename:
```
GET /alpine:latest?name=my-alpine
```

Multi-architecture image:
```
GET /alpine:latest?arch=linux/arm64
```

**Authentication:**

For private images, use HTTP Basic Authentication:

```bash
curl -u username:password http://localhost:8080/private/image:tag
```

**Response:**

The service returns the image as a gzipped tarball with appropriate headers:
- `Content-Type: application/x-gzip`
- `Content-Disposition: attachment; filename="{image_name}_{tag}.tgz"`

## License

Thape is the container image export service with [BSD-3-Clause licensed](LICENSE).

> (c) 2025 [Star Inc.](https://starinc.xyz/)
