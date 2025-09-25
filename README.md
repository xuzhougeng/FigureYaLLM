# FigureYa LLM 推荐系统

基于LLM的图表模块推荐系统，使用Go语言开发，支持OpenAI-compatible API协议。

> If you use this code in your work or research, in addition to complying with the license, we kindly request that you cite our publication:
> Xiaofan Lu, et al. (2025). FigureYa: A Standardized Visualization Framework for Enhancing Biomedical Data Interpretation and Research Efficiency. iMetaMed. https://doi.org/10.1002/imm3.70005


Visit <https://rvctmxhcmykh.sealosbja.site> for fun!

## 功能特性

- 🤖 基于LLM的智能推荐
- 📊 支持200+图表模块
- 🚀 轻量级Go实现
- 🔌 OpenAI-compatible API
- 💡 语义理解用户需求
- 📝 详细推荐理由

## 快速开始

### 1. 环境要求

- Go 1.21+
- OpenAI API key 或兼容的LLM API

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，设置API密钥
```

### 4. 启动服务

```bash
# 使用环境变量
export OPENAI_API_KEY="your_api_key"
export OPENAI_URL="https://api.openai.com/v1/chat/completions"
export MODEL="gpt-3.5-turbo"

go run main.go
```

服务将在 `http://localhost:8080` 启动

## API 接口

### 获取推荐

```bash
POST /recommend
Content-Type: application/json

{
  "query": "我想绘制生存曲线图"
}
```

响应：
```json
{
  "query": "我想绘制生存曲线图",
  "recommendations": [
    {
      "module": "FigureYa1survivalCurve_update",
      "需求描述": "该模块旨在帮助用户使用自己的数据绘制生存曲线...",
      "实用场景": "适用于展示分类样本的生存曲线...",
      "图片类型": "生存曲线",
      "score": 0.95,
      "reason": "完全匹配用户需求，专门用于绘制生存曲线"
    }
  ],
  "explanation": "基于您的需求，推荐以下生存分析相关模块..."
}
```

### 查看所有模块

```bash
GET /modules
```

### 健康检查

```bash
GET /health
```

## 支持的LLM API

### OpenAI
```bash
export OPENAI_URL="https://api.openai.com/v1/chat/completions"
export MODEL="gpt-3.5-turbo"
```

### Ollama (本地部署)
```bash
export OPENAI_URL="http://localhost:11434/v1/chat/completions"
export MODEL="llama2:7b"
```

### Together AI
```bash
export OPENAI_URL="https://api.together.xyz/v1/chat/completions"
export MODEL="mistralai/Mixtral-8x7B-Instruct-v0.1"
```

### 其他兼容API
任何支持OpenAI chat completions格式的API都可以使用。

## 部署

### Docker部署

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/figureya_docs_llm.json .
CMD ["./main"]
```

### 构建单文件可执行程序

```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o figureya-recommend main.go
```

## 使用示例

### 生存分析
```bash
curl -X POST http://localhost:8080/recommend \
  -H "Content-Type: application/json" \
  -d '{"query": "我想分析患者生存时间和生存率"}'
```

### 基因表达分析
```bash
curl -X POST http://localhost:8080/recommend \
  -H "Content-Type: application/json" \
  -d '{"query": "需要做RNA-seq差异表达基因的火山图"}'
```

### 相关性分析
```bash
curl -X POST http://localhost:8080/recommend \
  -H "Content-Type: application/json" \
  -d '{"query": "想看两个基因表达量的相关性"}'
```

## 项目结构

```
├── main.go                 # 主程序
├── go.mod                  # Go模块文件
├── figureya_docs_llm.json  # 模块数据
├── .env.example            # 环境变量示例
└── README.md              # 说明文档
```

## 性能特点

- 🚀 **轻量**: 单文件可执行，~10MB
- ⚡ **快速**: 毫秒级响应
- 💾 **低内存**: ~20MB内存占用
- 🔧 **易部署**: 无外部依赖

## License

MIT License


