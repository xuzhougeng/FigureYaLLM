# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a FigureYa LLM recommendation system written in Go that provides intelligent chart module recommendations based on user queries. The system features:

- **Backend**: Go HTTP server with OpenAI-compatible LLM integration
- **Frontend**: Modern web interface with Google-like search experience
- **Data**: 317+ biomedical visualization modules from FigureYa project
- **Deployment**: Docker containerization with production-ready configuration

The system analyzes user queries using LLM semantic understanding and recommends suitable visualization modules with detailed explanations and direct links to module documentation.

## Architecture

### Core Components
- **main.go**: HTTP server, LLM integration, recommendation logic
- **load_env.go**: Environment configuration with `.env` support
- **static/**: Frontend assets (HTML, CSS, JavaScript)
- **figureya_docs_llm.json**: Module metadata with Chinese field names

### Frontend Structure
- **index.html**: Google-like search interface with recommendation cards
- **script.js**: Interactive JavaScript with API integration and URL generation
- **styles.css**: Modern responsive design with card layouts and animations

### Data Model
```go
type Module struct {
    Module      string `json:"module"`
    需求描述        string `json:"需求描述"`    // Requirement description
    实用场景        string `json:"实用场景"`    // Use cases
    图片类型        string `json:"图片类型"`    // Chart types
    LLMStatus   string `json:"llm_status"`
}
```

**Important**: JSON parsing uses `map[string]interface{}` approach due to Chinese field names in struct tags not being properly handled by Go's JSON decoder.

## Common Commands

### Development
```bash
# Install dependencies
go mod tidy

# Run development server
go run main.go load_env.go

# Run on specific port
PORT=8080 go run main.go load_env.go

# Test API functionality
curl -X POST http://localhost:8080/recommend \
  -H "Content-Type: application/json" \
  -d '{"query": "我想绘制生存曲线图"}'
```

### Docker Deployment
```bash
# Quick deployment (recommended)
./deploy.sh start

# Manual Docker commands
docker build -t figureya-recommend .
docker run -p 8080:8080 --env-file .env figureya-recommend

# Docker Compose
docker-compose up -d

# Production with nginx
./deploy.sh production
```

### Testing
```bash
# Health check
curl http://localhost:8080/health

# Frontend test
open http://localhost:8080

# Module data check
curl http://localhost:8080/modules | jq '.modules[0]'
```

## Environment Configuration

Required variables (via `.env` file or environment):
- `OPENAI_API_KEY`: LLM service API key
- `BASE_URL`: LLM API base URL (auto-appends `/v1/chat/completions`)
- `MODEL`: Model name (e.g., "deepseek-chat", "gpt-3.5-turbo")
- `PORT`: Server port (default: 8080)
- `PROVIDER`: Provider identifier (default: "openai")

Example `.env`:
```bash
OPENAI_API_KEY=sk-xxxxx
BASE_URL=https://api.deepseek.com
MODEL=deepseek-chat
PROVIDER=openai
PORT=8080
```

## API Endpoints

### Core API
- `GET /`: Frontend web interface
- `POST /recommend`: Get module recommendations
- `GET /modules`: List all available modules
- `GET /health`: System health check with module count
- `GET /static/*`: Static assets (CSS, JS, images)

### API Response Format
```json
{
  "query": "user query",
  "recommendations": [
    {
      "module": "FigureYa1survivalCurve_update",
      "description": "模块需求描述",
      "useCase": "实用场景说明",
      "chartType": "图片类型",
      "score": 0.95,
      "reason": "推荐理由"
    }
  ],
  "explanation": "整体推荐说明"
}
```

## Key Implementation Details

### JSON Parsing Issue
**Critical**: Chinese field names in struct tags don't work with Go's JSON decoder. Use map-based parsing:

```go
// ❌ This doesn't work
type Module struct {
    需求描述 string `json:"需求描述"`  // Won't parse correctly
}

// ✅ Use this instead
var data map[string]interface{}
json.Decode(&data)
description := getString(moduleMap, "需求描述")
```

### Frontend Integration
- Static files use absolute paths (`/static/styles.css`, `/static/script.js`)
- Module URLs generated with pattern: `https://ying-ge.github.io/FigureYa/{module}/{module}.html`
- Cards display full module information with clickable documentation links
- Responsive design supports desktop and mobile

### LLM Integration
- Constructs context by concatenating all module metadata
- Uses structured JSON prompts for consistent responses
- Handles various OpenAI-compatible APIs (DeepSeek, OpenAI, Ollama, etc.)
- Temperature set to 0.3 for consistent recommendations

### Docker Optimization
- Multi-stage build produces 9MB final image
- Uses scratch base image for minimal attack surface
- Includes health checks and proper signal handling
- Supports both development and production configurations

## Deployment Configurations

### Development
```bash
go run main.go load_env.go
# Access: http://localhost:8080
```

### Docker Container
```bash
docker build -t figureya-recommend .
docker run -p 8080:8080 --env-file .env figureya-recommend
# Access: http://localhost:8080
```

### Production (with Nginx)
```bash
docker-compose --profile production up -d
# Access: http://localhost:80 (HTTP) or https://localhost:443 (HTTPS)
```

## Troubleshooting

### Common Issues
1. **Static files 404**: Ensure paths use `/static/` prefix in HTML
2. **Empty module fields**: Check JSON parsing uses map approach for Chinese fields
3. **LLM API errors**: Verify API key and endpoint configuration
4. **Port conflicts**: Use `PORT=xxxx` environment variable
5. **Docker build failures**: Check network connectivity for Go module downloads

### Debug Commands
```bash
# Check module loading
curl http://localhost:8080/modules | jq '.modules[0]'

# Test recommendation
curl -X POST http://localhost:8080/recommend \
  -H "Content-Type: application/json" \
  -d '{"query": "测试"}' | jq .

# Monitor logs
docker-compose logs -f figureya-recommend
```

## Performance Notes

- **Startup time**: ~3 seconds
- **Memory usage**: ~30MB runtime
- **Image size**: 9MB (Docker)
- **API response**: 15-20 seconds (LLM processing)
- **Module loading**: 317 modules filtered from JSON
- **Frontend response**: <100ms (static files)