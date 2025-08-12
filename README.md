<center>

# AkashGen API Go

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-AGPL_v3-green.svg)](LICENSE)
[![API](https://img.shields.io/badge/API-REST-orange.svg)](https://restfulapi.net/)

A high-performance Go-based API proxy service for the Akash Network's image generation capabilities. This service provides a clean REST API interface to interact with Akash Network's decentralized AI image generation infrastructure.

[中文文档](README_CN.md) | [English](README.md)

</center>

## Features

- 🚀 **High Performance**: Built with Go and Gin framework for optimal performance
- 🔄 **Asynchronous Processing**: Non-blocking job submission with status polling
- 🛡️ **Concurrency Control**: Built-in rate limiting to prevent resource exhaustion
- 📊 **Structured Logging**: Comprehensive logging with Zap for production monitoring
- 🏗️ **Clean Architecture**: Modular design following Go best practices
- 🌐 **REST API**: Simple and intuitive HTTP endpoints
- ⚡ **Graceful Shutdown**: Proper cleanup and connection handling
- 🎯 **GPU Preference**: Automatic GPU selection from available options

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Internet connection to access Akash Network

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/006lp/akashgen-api-go.git
   cd akashgen-api-go
   ```

2. **Initialize Go modules**
   ```bash
   go mod init akashgen-api-go
   go mod tidy
   ```

3. **Build and run**
   ```bash
   go build -o akashgen-api .
   ./akashgen-api
   ```

Or run directly:
```bash
go run main.go
```

The server will start on port `6571` by default.

## API Documentation

### Generate Image

Generate an AI image using Akash Network's decentralized infrastructure.

**Endpoint:** `POST /api/generate`

**Request Body:**
```json
{
    "prompt": "A beautiful sunset over mountains",
    "negative": "blurry, low quality",
    "sampler": "DPM++ 2M Karras",
    "scheduler": "karras"
}
```

**Parameters:**
- `prompt` (required): Text description of the desired image
- `negative` (optional): Text describing what to avoid in the image
- `sampler` (required): The sampling method to use
- `scheduler` (required): The scheduling algorithm

**Supported Samplers:**
```
euler, euler_cfg_pp, euler_ancestral, euler_ancestral_cfg_pp, heun, heunpp2, 
dpm_2, dpm_2_ancestral, lms, dpm_fast, dpm_adaptive, dpmpp_2s_ancestral, 
dpmpp_2s_ancestral_cfg_pp, dpmpp_sde, dpmpp_sde_gpu, dpmpp_2m, dpmpp_2m_cfg_pp, 
dpmpp_2m_sde, dpmpp_2m_sde_gpu, dpmpp_3m_sde, dpmpp_3m_sde_gpu, ddpm, lcm, 
ipndm, ipndm_v, deis, ddim, uni_pc, uni_pc_bh2
```

**Supported Schedulers:**
```
normal, karras, exponential, sgm_uniform, simple, ddim_uniform, beta, linear_quadratic
```

**Response:**
- **Success (200)**: Returns the generated image as binary data
- **Error (400)**: Invalid request format
- **Error (500)**: Generation failed or timeout

**Example using curl:**
```bash
curl -X POST http://localhost:6571/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "A serene lake with mountains in the background",
    "negative": "ugly, blurry, low resolution",
    "sampler": "DPM++ 2M Karras",
    "scheduler": "karras"
  }' \
  --output generated_image.png
```

### Health Check

Check if the service is running properly.

**Endpoint:** `GET /health`

**Response:**
```json
{
    "status": "ok"
}
```

## Configuration

The service can be configured by modifying the constants in `config/config.go`:

```go
const (
    // Timeout configurations
    GenerateTimeout    = 30 * time.Second  // Max time for generation request
    StatusTimeout      = 10 * time.Second  // Max time for status check
    ImageFetchTimeout  = 30 * time.Second  // Max time for image download
    PollingInterval    = 1 * time.Second   // How often to check job status
    MaxPollingDuration = 5 * time.Minute   // Max time to wait for completion
    
    // Concurrency
    MaxConcurrentRequests = 10  // Max simultaneous requests
    
    // Server
    ServerPort = ":6571"  // Port to listen on
)
```

## Architecture

```
akashgen-api-go/
├── main.go              # Application entry point
├── config/
│   └── config.go       # Configuration constants
├── handlers/
│   └── generate.go     # HTTP request handlers
├── models/
│   └── types.go        # Data structures and types
├── services/
│   ├── akash.go        # Akash Network API client
│   └── image.go        # Image fetching service
├── middleware/
│   └── logger.go       # HTTP logging middleware
└── utils/
    └── http.go         # HTTP utilities
```

### Components

- **Handlers**: Process HTTP requests and responses
- **Services**: Business logic for interacting with external APIs
- **Models**: Data structures for requests and responses
- **Middleware**: Cross-cutting concerns like logging
- **Config**: Centralized configuration management
- **Utils**: Reusable utility functions

## Development

### Project Structure

This project follows Go best practices with clear separation of concerns:

- Clean architecture with dependency injection
- Modular design for easy testing and maintenance
- Proper error handling and logging
- Context-based timeout management

### Building for Production

```bash
# Build optimized binary
go build -ldflags="-w -s" -o akashgen-api .

# Or build for different platforms
GOOS=linux GOARCH=amd64 go build -o akashgen-api-linux .
```

### Docker Support

Create a `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o akashgen-api .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/akashgen-api .
EXPOSE 6571
CMD ["./akashgen-api"]
```

Build and run:
```bash
docker build -t akashgen-api .
docker run -p 6571:6571 akashgen-api
```

## GPU Support

The service automatically selects from preferred GPU types in order:
1. RTX4090
2. A10
3. A100
4. V100-32Gi
5. H100

## Error Handling

The API provides detailed error responses:

```json
{
    "error": "Failed to generate image",
    "details": "specific error description"
}
```

Common error scenarios:
- Invalid JSON format
- Missing required fields
- Network timeouts
- Upstream service failures
- Job execution failures

## Performance

- **Concurrency**: Configurable request limiting prevents overload
- **Timeouts**: Multiple timeout layers prevent hanging requests
- **Connection Pooling**: HTTP client reuse for efficiency
- **Memory Management**: Proper cleanup of resources
- **Graceful Shutdown**: Clean termination handling

## Monitoring

The service provides structured JSON logs suitable for aggregation:

```json
{
    "level": "info",
    "ts": 1640995200.123456,
    "msg": "HTTP Request",
    "method": "POST",
    "path": "/api/generate",
    "status": 200,
    "duration": "2.5s",
    "client_ip": "192.168.1.100"
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go conventions and best practices
- Add tests for new functionality
- Update documentation for API changes
- Use meaningful commit messages
- Ensure code passes `go fmt` and `go vet`

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run race condition detection
go test -race ./...
```

## Troubleshooting

### Common Issues

1. **Port already in use**: Change the port in `config/config.go`
2. **Timeout errors**: Increase timeout values for slower networks
3. **Memory issues**: Reduce `MaxConcurrentRequests` if experiencing memory pressure
4. **Connection refused**: Ensure Akash Network endpoints are accessible

### Debug Mode

For development, you can enable debug logging by modifying the logger initialization in `main.go`:

```go
// Replace zap.NewProduction() with:
logger, err = zap.NewDevelopment()
```

## License

This project is licensed under the AGPL v3 License - see the [LICENSE](LICENSE) file for details.

## Support

- Create an [issue](https://github.com/yourusername/akashgen-api-go/issues) for bugs or feature requests
- Check [Akash Network documentation](https://docs.akash.network) for network-related questions
- Join the discussion in [Akash Community](https://akash.network/community)

## Acknowledgments

- [Akash Network](https://akash.network) for providing decentralized cloud infrastructure
- [Gin Web Framework](https://github.com/gin-gonic/gin) for HTTP routing
- [Zap](https://github.com/uber-go/zap) for structured logging

---

**Note**: This is an unofficial proxy service for Akash Network's image generation API. Please refer to Akash Network's official documentation for the most up-to-date information about their services.