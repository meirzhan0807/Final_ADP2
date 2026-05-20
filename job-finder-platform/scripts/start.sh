#!/bin/bash
set -e

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; BLUE='\033[0;34m'; NC='\033[0m'

echo -e "${BLUE}"
echo "  ╔═══════════════════════════════════════╗"
echo "  ║      Job Finder Platform v1.0         ║"
echo "  ║   Microservices Architecture (Go)     ║"
echo "  ╚═══════════════════════════════════════╝"
echo -e "${NC}"

command -v docker >/dev/null 2>&1 || { echo -e "${RED}❌ Docker not found. Install Docker first.${NC}"; exit 1; }
(docker compose version >/dev/null 2>&1 || docker-compose version >/dev/null 2>&1) || { echo -e "${RED}❌ Docker Compose not found.${NC}"; exit 1; }

echo -e "${GREEN}✅ Docker detected${NC}"

if [ ! -f .env ]; then
  cp .env.example .env
  echo -e "${YELLOW}⚠️  .env created from .env.example"
  echo -e "   Edit .env to add SMTP credentials for email (optional).${NC}"
fi

echo -e "\n${BLUE}🔨 Building and starting all services...${NC}"
docker compose up -d --build 2>&1 | tail -20

echo -e "\n${YELLOW}⏳ Waiting for services to become healthy (30s)...${NC}"
sleep 30

echo -e "\n${BLUE}📊 Service Status:${NC}"
docker compose ps

echo -e "\n${GREEN}═══════════════════════════════════════════════════${NC}"
echo -e "${GREEN}✅ Job Finder Platform is RUNNING!${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════${NC}"
echo ""
echo -e "  🌐 ${BLUE}API Gateway:${NC}   http://localhost:8080"
echo -e "  📊 ${BLUE}Grafana:${NC}       http://localhost:3000  (admin/admin123)"
echo -e "  🔥 ${BLUE}Prometheus:${NC}    http://localhost:9090"
echo -e "  🔍 ${BLUE}Jaeger:${NC}        http://localhost:16686"
echo -e "  📨 ${BLUE}NATS Monitor:${NC}  http://localhost:8222"
echo ""
echo -e "${YELLOW}📖 Quick API Test:${NC}"
echo ""
echo -e "  # Health check"
echo -e "  curl http://localhost:8080/health"
echo ""
echo -e "  # Register"
echo -e "  curl -X POST http://localhost:8080/api/v1/auth/register \\"
echo -e "    -H 'Content-Type: application/json' \\"
echo -e "    -d '{\"email\":\"test@example.com\",\"password\":\"password123\",\"first_name\":\"Test\",\"last_name\":\"User\",\"role\":\"jobseeker\"}'"
echo ""
echo -e "  # Login"
echo -e "  curl -X POST http://localhost:8080/api/v1/auth/login \\"
echo -e "    -H 'Content-Type: application/json' \\"
echo -e "    -d '{\"email\":\"test@example.com\",\"password\":\"password123\"}'"
echo ""
echo -e "  # Search jobs"
echo -e "  curl 'http://localhost:8080/api/v1/jobs/search?q=developer&location=Almaty'"
echo ""
echo -e "${BLUE}🛑 To stop: ${NC}docker compose down"
echo -e "${BLUE}🗑️  To reset: ${NC}docker compose down -v"
