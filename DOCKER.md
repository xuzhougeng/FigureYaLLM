# Docker 部署指南

FigureYa LLM推荐系统的Docker容器化部署指南。

## 🚀 快速开始

### 方法1: 使用一键部署脚本

```bash
# 快速启动（推荐）
./deploy.sh

# 或者指定启动命令
./deploy.sh start
```

### 方法2: 使用Docker Compose

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 方法3: 手动Docker命令

```bash
# 构建镜像
docker build -t figureya-recommend .

# 运行容器
docker run -d \
  --name figureya-recommend \
  -p 8080:8080 \
  --env-file .env \
  figureya-recommend
```

## 📋 前置要求

- Docker 20.0+
- Docker Compose 2.0+
- `.env`文件配置（或环境变量）

## 🔧 环境配置

创建`.env`文件：

```bash
# 复制示例配置
cp .env.example .env

# 编辑配置文件
vi .env
```

必需的环境变量：
```bash
OPENAI_API_KEY=your_api_key_here
BASE_URL=https://api.deepseek.com
MODEL=deepseek-chat
PROVIDER=openai
PORT=8080
```

## 📦 镜像特性

- **极小体积**: ~9MB（基于scratch镜像）
- **多阶段构建**: 优化构建速度和镜像大小
- **安全设计**: 非root用户运行
- **健康检查**: 内置健康监控
- **生产就绪**: 包含完整的静态资源

## 🛠 部署选项

### 开发环境部署

```bash
# 简单启动
docker-compose up -d

# 实时查看日志
docker-compose logs -f figureya-recommend
```

### 生产环境部署（带Nginx）

```bash
# 启动带nginx代理的生产环境
./deploy.sh production

# 或者使用profile
docker-compose --profile production up -d
```

生产环境特性：
- Nginx反向代理
- 负载均衡
- 静态资源缓存
- 速率限制
- SSL支持（需配置证书）

## 📊 服务验证

部署成功后验证服务：

```bash
# 健康检查
curl http://localhost:8080/health

# Web界面
open http://localhost:8080

# API测试
curl -X POST http://localhost:8080/recommend \
  -H "Content-Type: application/json" \
  -d '{"query": "我想绘制生存曲线图"}'
```

## 🔍 监控和维护

### 查看容器状态

```bash
# 查看运行状态
docker-compose ps

# 查看资源使用
docker stats figureya-recommend

# 查看容器信息
docker inspect figureya-recommend
```

### 日志管理

```bash
# 实时日志
docker-compose logs -f

# 最近日志
docker-compose logs --tail=100

# 指定服务日志
docker-compose logs figureya-recommend
```

### 更新部署

```bash
# 重新构建和部署
./deploy.sh restart

# 或者手动更新
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## 🛡️ 安全配置

### 网络安全

```yaml
# docker-compose.yml 网络配置
networks:
  figureya-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

### 资源限制

```yaml
services:
  figureya-recommend:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 256M
        reservations:
          cpus: '0.1'
          memory: 64M
```

### 卷挂载安全

```yaml
volumes:
  # 只读挂载配置文件
  - ./.env:/app/.env:ro
  - ./figureya_docs_llm.json:/figureya_docs_llm.json:ro
```

## 🔧 故障排除

### 常见问题

1. **端口占用**
   ```bash
   # 检查端口使用
   netstat -tulpn | grep 8080

   # 修改端口
   export PORT=8081
   ```

2. **权限问题**
   ```bash
   # 检查文件权限
   ls -la .env

   # 修复权限
   chmod 644 .env
   ```

3. **网络连接**
   ```bash
   # 测试API连接
   docker exec figureya-recommend wget -q --spider http://localhost:8080/health
   ```

4. **内存不足**
   ```bash
   # 检查容器资源使用
   docker stats --no-stream

   # 增加内存限制
   docker-compose down
   # 编辑 docker-compose.yml 增加内存限制
   docker-compose up -d
   ```

### 调试模式

```bash
# 调试模式启动
docker run -it --rm \
  --env-file .env \
  -p 8080:8080 \
  figureya-recommend /bin/sh

# 查看容器内部
docker exec -it figureya-recommend /bin/sh
```

## 📈 性能优化

### 构建优化

```bash
# 使用构建缓存
docker build --cache-from figureya-recommend .

# 并行构建
docker buildx build --platform linux/amd64,linux/arm64 .
```

### 运行时优化

```yaml
# docker-compose.yml 优化配置
services:
  figureya-recommend:
    restart: unless-stopped
    environment:
      - GIN_MODE=release
    deploy:
      replicas: 2
      update_config:
        parallelism: 1
        delay: 10s
```

## 🎯 部署到云平台

### AWS ECS

```bash
# 推送到 ECR
aws ecr get-login-password | docker login --username AWS --password-stdin <account>.dkr.ecr.region.amazonaws.com
docker tag figureya-recommend:latest <account>.dkr.ecr.region.amazonaws.com/figureya-recommend:latest
docker push <account>.dkr.ecr.region.amazonaws.com/figureya-recommend:latest
```

### 云服务器

```bash
# 复制到服务器
scp -r . user@server:/app/figureya-recommend/

# 服务器上启动
ssh user@server "cd /app/figureya-recommend && ./deploy.sh"
```

## 📝 维护指南

### 定期维护

```bash
# 清理未使用的镜像
docker system prune -f

# 更新基础镜像
docker pull golang:1.21-alpine
docker build --no-cache -t figureya-recommend .

# 备份配置
tar -czf backup-$(date +%Y%m%d).tar.gz .env figureya_docs_llm.json
```

### 监控指标

- 内存使用: < 100MB
- CPU使用: < 10%
- 响应时间: < 2s
- 健康检查: 通过率 > 99%

通过以上配置，你可以轻松地在任何支持Docker的环境中部署FigureYa推荐系统！