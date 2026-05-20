#!/bin/bash

# ═══════════════════════════════════════════════════════════
#  Job Finder Platform — Seed Script
#  Дайын деректерді толтырады: admin, employer, jobseeker, 10 жұмыс
# ═══════════════════════════════════════════════════════════

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; BLUE='\033[0;34m'; NC='\033[0m'
API="http://localhost:8080/api/v1"

echo -e "${BLUE}"
echo "  ╔══════════════════════════════════════════╗"
echo "  ║     Job Finder — Деректерді толтыру      ║"
echo "  ╚══════════════════════════════════════════╝"
echo -e "${NC}"

# API дайын ба тексер
echo -e "${YELLOW}⏳ API Gateway күтілуде...${NC}"
for i in $(seq 1 15); do
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health 2>/dev/null)
  if [ "$STATUS" = "200" ]; then
    echo -e "${GREEN}✅ API Gateway дайын!${NC}\n"
    break
  fi
  if [ $i -eq 15 ]; then
    echo -e "${RED}❌ API Gateway қосылмады. docker compose up -d іске қосыңыз.${NC}"
    exit 1
  fi
  sleep 2
done

post() {
  local url=$1; local data=$2; local token=$3
  local auth=""
  if [ -n "$token" ]; then auth="-H \"Authorization: Bearer $token\""; fi
  curl -s -X POST "$API$url" \
    -H "Content-Type: application/json" \
    ${token:+-H "Authorization: Bearer $token"} \
    -d "$data"
}

get_token() {
  local email=$1; local pass=$2
  echo $(curl -s -X POST "$API/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$email\",\"password\":\"$pass\"}" | \
    python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('access_token',''))" 2>/dev/null)
}

# ── 1. EMPLOYER тіркелу ───────────────────────────────────
echo -e "${BLUE}👔 Employer аккаунты жасалуда...${NC}"
RESP=$(post "/auth/register" '{
  "email":"employer@jobfinder.kz",
  "password":"Password123",
  "first_name":"Айдар",
  "last_name":"Сейткали",
  "role":"employer"
}')
echo "   $RESP"

# ── 2. JOBSEEKER тіркелу ─────────────────────────────────
echo -e "\n${BLUE}👤 Jobseeker аккаунты жасалуда...${NC}"
RESP=$(post "/auth/register" '{
  "email":"user@jobfinder.kz",
  "password":"Password123",
  "first_name":"Арман",
  "last_name":"Бекжанов",
  "role":"jobseeker"
}')
echo "   $RESP"

# ── 3. DEMO аккаунт ──────────────────────────────────────
echo -e "\n${BLUE}🎓 Demo jobseeker жасалуда...${NC}"
RESP=$(post "/auth/register" '{
  "email":"demo@jobfinder.kz",
  "password":"Password123",
  "first_name":"Динара",
  "last_name":"Қасымова",
  "role":"jobseeker"
}')
echo "   $RESP"

# ── 4. EMPLOYER токен алу ────────────────────────────────
echo -e "\n${YELLOW}🔑 Employer токені алынуда...${NC}"
EMP_TOKEN=$(get_token "employer@jobfinder.kz" "Password123")
if [ -z "$EMP_TOKEN" ]; then
  echo -e "${RED}❌ Employer кіру сәтсіз. Бірақ жалғасамыз...${NC}"
else
  echo -e "${GREEN}✅ Employer токені алынды${NC}"
fi

# ── 5. 10 ЖҰМЫС ЖАРИЯЛАУ ────────────────────────────────
echo -e "\n${BLUE}💼 10 жұмыс жарияланудa...${NC}"

create_job() {
  local title="$1" desc="$2" loc="$3" type="$4" cat="$5" exp="$6" smin=$7 smax=$8
  RESULT=$(curl -s -X POST "$API/employer/jobs" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $EMP_TOKEN" \
    -d "{
      \"title\": \"$title\",
      \"description\": \"$desc\",
      \"location\": \"$loc\",
      \"job_type\": \"$type\",
      \"category\": \"$cat\",
      \"experience_level\": \"$exp\",
      \"salary_min\": $smin,
      \"salary_max\": $smax,
      \"currency\": \"KZT\"
    }")
  echo -e "   ${GREEN}✓${NC} $title"
}

create_job \
  "Senior Go Developer" \
  "Микросервис архитектурасын жасайтын тәжірибелі Go разработчик іздейміз. gRPC, PostgreSQL, Docker тәжірибесі міндетті." \
  "Алматы" "full-time" "Technology" "senior" 400000 700000

create_job \
  "Python Backend Engineer" \
  "FastAPI немесе Django тәжірибесі бар Python разработчик. ML pipeline жасаумен айналысасыз." \
  "Астана" "full-time" "Technology" "middle" 300000 500000

create_job \
  "React Frontend Developer" \
  "React, TypeScript, Redux тәжірибесі бар frontend разработчик. UI/UX дизайнерлермен тығыз жұмыс." \
  "Алматы" "full-time" "Technology" "middle" 280000 450000

create_job \
  "DevOps Engineer" \
  "Kubernetes, Docker, CI/CD pipeline жасау тәжірибесі бар DevOps инженер іздейміз. AWS немесе GCP білімі артықшылық." \
  "Remote" "remote" "Technology" "senior" 500000 800000

create_job \
  "Data Scientist" \
  "ML модельдерін жасайтын Data Scientist. Python, TensorFlow/PyTorch, SQL тәжірибесі міндетті." \
  "Алматы" "full-time" "Data & AI" "senior" 450000 750000

create_job \
  "iOS Developer (Swift)" \
  "Swift тілінде iOS қосымшалар жасайтын мобильді разработчик. App Store-да жарияланған қосымшалары болуы керек." \
  "Астана" "full-time" "Mobile" "middle" 350000 600000

create_job \
  "Junior Java Developer" \
  "Spring Boot тәжірибесі бар немесе үйренуге дайын Junior Java разработчик. Ментор қолдауы бар." \
  "Алматы" "full-time" "Technology" "junior" 150000 250000

create_job \
  "Product Manager" \
  "IT өнімдерін басқаратын Product Manager. Agile/Scrum әдістемесін білу міндетті. Roadmap жасау тәжірибесі." \
  "Алматы" "full-time" "Management" "middle" 400000 650000

create_job \
  "UX/UI Designer" \
  "Figma тәжірибесі бар UX/UI дизайнер. Пайдаланушы зерттеулері жүргізе алу, прототип жасай алу." \
  "Remote" "remote" "Design" "middle" 250000 400000

create_job \
  "QA Engineer (Automation)" \
  "Selenium, Playwright немесе Cypress тәжірибесі бар QA автоматизация инженері. CI/CD интеграциясы." \
  "Шымкент" "full-time" "Quality Assurance" "middle" 200000 350000

echo -e "\n${GREEN}═══════════════════════════════════════════════════${NC}"
echo -e "${GREEN}✅ Деректер толтырылды!${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════${NC}"
echo ""
echo -e "  📧 ${YELLOW}Employer:${NC}   employer@jobfinder.kz  / Password123"
echo -e "  📧 ${YELLOW}Jobseeker:${NC}  user@jobfinder.kz      / Password123"
echo -e "  📧 ${YELLOW}Demo:${NC}       demo@jobfinder.kz      / Password123"
echo ""
echo -e "  🌐 ${BLUE}Фронт:${NC} frontend/index.html файлын браузерде ашыңыз"
echo -e "  🌐 ${BLUE}API:${NC}   http://localhost:8080/api/v1/jobs"
echo ""
