#!/bin/bash

# FigureYa LLM推荐系统快速部署脚本
# Usage: ./deploy.sh [build|start|stop|restart|logs|clean]

set -e

PROJECT_NAME="figureya-recommend"
IMAGE_NAME="figureya-recommend:latest"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
check_docker() {
    if ! docker info >/dev/null 2>&1; then
        log_error "Docker is not running. Please start Docker first."
        exit 1
    fi
}

# Check if .env file exists
check_env() {
    if [ ! -f ".env" ]; then
        log_warning ".env file not found. Creating from .env.example..."
        if [ -f ".env.example" ]; then
            cp .env.example .env
            log_warning "Please edit .env file with your API credentials before starting the service."
        else
            log_error ".env.example file not found. Please create .env file manually."
            exit 1
        fi
    fi
}

# Build Docker image
build_image() {
    log_info "Building Docker image..."
    docker build -t $IMAGE_NAME .
    log_success "Docker image built successfully"
}

# Start services
start_services() {
    log_info "Starting services..."
    docker-compose up -d
    log_success "Services started successfully"

    # Wait for health check
    log_info "Waiting for service to be healthy..."
    sleep 10

    if docker-compose ps | grep -q "healthy\|Up"; then
        log_success "Service is running at http://localhost:8080"
    else
        log_error "Service failed to start properly"
        docker-compose logs
        exit 1
    fi
}

# Stop services
stop_services() {
    log_info "Stopping services..."
    docker-compose down
    log_success "Services stopped"
}

# Restart services
restart_services() {
    log_info "Restarting services..."
    docker-compose restart
    log_success "Services restarted"
}

# Show logs
show_logs() {
    docker-compose logs -f
}

# Clean up
clean_up() {
    log_info "Cleaning up..."
    docker-compose down -v --remove-orphans
    docker rmi $IMAGE_NAME 2>/dev/null || true
    docker system prune -f
    log_success "Cleanup completed"
}

# Production deployment
deploy_production() {
    log_info "Deploying to production with nginx..."
    docker-compose --profile production up -d
    log_success "Production deployment completed"
    log_info "Service available at:"
    log_info "  - HTTP: http://localhost:80"
    log_info "  - HTTPS: https://localhost:443 (if SSL configured)"
}

# Main script logic
case "${1:-start}" in
    build)
        check_docker
        build_image
        ;;
    start)
        check_docker
        check_env
        build_image
        start_services
        ;;
    stop)
        check_docker
        stop_services
        ;;
    restart)
        check_docker
        restart_services
        ;;
    logs)
        check_docker
        show_logs
        ;;
    clean)
        check_docker
        clean_up
        ;;
    production)
        check_docker
        check_env
        build_image
        deploy_production
        ;;
    *)
        echo "Usage: $0 [build|start|stop|restart|logs|clean|production]"
        echo ""
        echo "Commands:"
        echo "  build       - Build Docker image"
        echo "  start       - Build and start all services (default)"
        echo "  stop        - Stop all services"
        echo "  restart     - Restart services"
        echo "  logs        - Show service logs"
        echo "  clean       - Stop services and remove containers/images"
        echo "  production  - Deploy with nginx proxy"
        echo ""
        echo "Examples:"
        echo "  $0                  # Quick start"
        echo "  $0 start           # Same as above"
        echo "  $0 production      # Production deployment with nginx"
        echo "  $0 logs            # View logs"
        echo "  $0 clean           # Full cleanup"
        exit 1
        ;;
esac