# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a FigureYa LLM recommendation system written in Go that provides intelligent chart module recommendations based on user queries. The system uses OpenAI-compatible APIs to analyze user requests and recommend suitable visualization modules from a database of 200+ FigureYa chart modules.

## Architecture

The application follows a single-file architecture with these key components:

- **Main Service**: `main.go` contains the HTTP server, LLM integration, and recommendation logic
- **Environment Loading**: `load_env.go` handles `.env` file parsing with an `init()` function
- **Data Source**: `figureya_docs_llm.json` contains module metadata (需求描述, 实用场景, 图片类型)
- **LLM Integration**: Uses OpenAI-compatible chat completion API for semantic understanding

The system loads module data at startup, filters for modules with `llm_status: "ok"`, and uses LLM prompting to match user queries with relevant modules.

## Common Commands

### Development
```bash
# Install dependencies
go mod tidy

# Run the server (requires .env or environment variables)
go run main.go load_env.go

# Build production binary
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o figureya-recommend main.go load_env.go

# Test the API
curl -X POST http://localhost:8080/recommend \
  -H "Content-Type: application/json" \
  -d '{"query": "我想绘制生存曲线图"}'
```

### Environment Configuration
The system requires these environment variables (can be set via `.env` file):
- `OPENAI_API_KEY`: API key for LLM service
- `BASE_URL` or `OPENAI_URL`: LLM API endpoint
- `MODEL`: LLM model name (e.g., "deepseek-chat", "gpt-3.5-turbo")
- `PORT`: Server port (default: 8080)

The system automatically appends `/v1/chat/completions` to base URLs that don't already contain it.

## API Endpoints

- `POST /recommend`: Takes a query and returns recommended modules with scores and explanations
- `GET /modules`: Lists all available modules
- `GET /health`: Health check with module count

## Key Implementation Details

- Uses Gin web framework for HTTP routing
- LLM prompts are constructed by concatenating all module information into context
- Response parsing extracts JSON from LLM output and maps module names back to full module data
- CORS headers are enabled for cross-origin requests
- Only modules with `llm_status: "ok"` are included in recommendations