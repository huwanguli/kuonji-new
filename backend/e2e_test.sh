#!/bin/bash
# ============================================================
# Kuonji Blog E2E 测试脚本
# 针对运行中的后端服务 (localhost:8080)
# 使用方式: bash e2e_test.sh
# ============================================================

BASE="http://localhost:8080/api/v1"
PASS=0
FAIL=0
TOKEN=""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

assert_status() {
  local desc="$1" expected="$2" actual="$3" body="$4"
  if [ "$actual" = "$expected" ]; then
    echo -e "  ${GREEN}✓${NC} $desc (HTTP $actual)"
    PASS=$((PASS + 1))
  else
    echo -e "  ${RED}✗${NC} $desc: expected $expected, got $actual"
    echo "    Response: $body" | head -c 200
    echo
    FAIL=$((FAIL + 1))
  fi
}

assert_json() {
  local desc="$1" field="$2" expected="$3" body="$4"
  actual=$(echo "$body" | python3 -c "import sys,json; print(json.loads(sys.stdin.read())$field)" 2>/dev/null)
  if [ "$actual" = "$expected" ]; then
    echo -e "  ${GREEN}✓${NC} $desc ($field = $actual)"
    PASS=$((PASS + 1))
  else
    echo -e "  ${RED}✗${NC} $desc: expected $field=$expected, got $actual"
    FAIL=$((FAIL + 1))
  fi
}

echo "========================================="
echo " Kuonji Blog E2E 测试"
echo "========================================="
echo

# ── 1. 文章列表 ──
echo "▸ GET /articles"
resp=$(curl -s -w "\n%{http_code}" "$BASE/articles?page=1&page_size=5")
status=$(echo "$resp" | tail -1)
body=$(echo "$resp" | head -n -1)
assert_status "获取文章列表" "200" "$status" "$body"

total=$(echo "$body" | python3 -c "import sys,json; print(json.loads(sys.stdin.read())['data']['total'])" 2>/dev/null)
echo "  ℹ 文章总数: $total"

# ── 2. 标签过滤 ──
echo
echo "▸ GET /articles?tag_id=22 (go)"
resp=$(curl -s -w "\n%{http_code}" "$BASE/articles?tag_id=22&page_size=50")
status=$(echo "$resp" | tail -1)
body=$(echo "$resp" | head -n -1)
assert_status "按标签过滤文章" "200" "$status" "$body"

# ── 3. 分类过滤 ──
echo
echo "▸ GET /articles?category_id=11"
resp=$(curl -s -w "\n%{http_code}" "$BASE/articles?category_id=11")
status=$(echo "$resp" | tail -1)
body=$(echo "$resp" | head -n -1)
assert_status "按分类过滤文章" "200" "$status" "$body"

# ── 4. 文章详情 ──
echo
echo "▸ GET /articles/9"
resp=$(curl -s -w "\n%{http_code}" "$BASE/articles/9")
status=$(echo "$resp" | tail -1)
body=$(echo "$resp" | head -n -1)
assert_status "获取文章详情" "200" "$status" "$body"
assert_json "文章标题正确" "['data']['title']" "用 Gin 搭建 RESTful API 的最佳实践" "$body"

# ── 5. 文章 404 ──
echo
echo "▸ GET /articles/99999"
resp=$(curl -s -w "\n%{http_code}" "$BASE/articles/99999")
status=$(echo "$resp" | tail -1)
assert_status "不存在的文章返回 404" "404" "$status" ""

# ── 6. 分类列表 ──
echo
echo "▸ GET /categories"
resp=$(curl -s -w "\n%{http_code}" "$BASE/categories")
status=$(echo "$resp" | tail -1)
body=$(echo "$resp" | head -n -1)
assert_status "获取分类列表" "200" "$status" "$body"

# ── 7. 标签列表 ──
echo
echo "▸ GET /tags"
resp=$(curl -s -w "\n%{http_code}" "$BASE/tags")
status=$(echo "$resp" | tail -1)
body=$(echo "$resp" | head -n -1)
assert_status "获取标签列表" "200" "$status" "$body"

# ── 8. 系列列表 ──
echo
echo "▸ GET /series"
resp=$(curl -s -w "\n%{http_code}" "$BASE/series")
status=$(echo "$resp" | tail -1)
body=$(echo "$resp" | head -n -1)
assert_status "获取系列列表" "200" "$status" "$body"

# ── 9. 系列详情 ──
echo
echo "▸ GET /series/4"
resp=$(curl -s -w "\n%{http_code}" "$BASE/series/4")
status=$(echo "$resp" | tail -1)
body=$(echo "$resp" | head -n -1)
assert_status "获取系列详情" "200" "$status" "$body"

# ── 10. 搜索 ──
echo
echo "▸ GET /search?q=Go"
resp=$(curl -s -w "\n%{http_code}" "$BASE/search?q=Go")
status=$(echo "$resp" | tail -1)
body=$(echo "$resp" | head -n -1)
assert_status "搜索文章" "200" "$status" "$body"

echo
echo "▸ GET /search?q= (空关键词)"
resp=$(curl -s -w "\n%{http_code}" "$BASE/search?q=")
status=$(echo "$resp" | tail -1)
assert_status "空搜索返回 400" "400" "$status" ""

# ── 11. 评论 ──
echo
echo "▸ POST /articles/9/comments (创建评论)"
resp=$(curl -s -w "\n%{http_code}" -X POST "$BASE/articles/9/comments" \
  -H "Content-Type: application/json" \
  -d '{"nickname":"E2E Test","email":"e2e@test.com","content":"E2E test comment"}')
status=$(echo "$resp" | tail -1)
assert_status "创建评论" "201" "$status" ""

echo
echo "▸ GET /articles/9/comments"
resp=$(curl -s -w "\n%{http_code}" "$BASE/articles/9/comments")
status=$(echo "$resp" | tail -1)
body=$(echo "$resp" | head -n -1)
assert_status "获取评论列表" "200" "$status" "$body"

# ── 12. 认证 ──
echo
echo "▸ POST /auth/login"
resp=$(curl -s -w "\n%{http_code}" -X POST "$BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')
status=$(echo "$resp" | tail -1)
body=$(echo "$resp" | head -n -1)
assert_status "管理员登录" "200" "$status" "$body"
TOKEN=$(echo "$body" | python3 -c "import sys,json; print(json.loads(sys.stdin.read())['data']['token'])" 2>/dev/null)

if [ -n "$TOKEN" ]; then
  echo "  ℹ Token 获取成功"

  # ── 13. 管理接口 ──
  echo
  echo "▸ GET /admin/stats"
  resp=$(curl -s -w "\n%{http_code}" "$BASE/admin/stats" \
    -H "Authorization: Bearer $TOKEN")
  status=$(echo "$resp" | tail -1)
  body=$(echo "$resp" | head -n -1)
  assert_status "获取统计数据" "200" "$status" "$body"

  echo
  echo "▸ POST /admin/articles (创建文章)"
  resp=$(curl -s -w "\n%{http_code}" -X POST "$BASE/admin/articles" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"title":"E2E Test Article","content":"# E2E\n\nThis is an e2e test article.","summary":"E2E test","category_id":11,"tag_ids":[22,23],"status":1}')
  status=$(echo "$resp" | tail -1)
  body=$(echo "$resp" | head -n -1)
  assert_status "创建文章（直接发布）" "201" "$status" "$body"

  # 获取新文章 ID
  NEW_ID=$(echo "$body" | python3 -c "import sys,json; print(json.loads(sys.stdin.read())['data']['id'])" 2>/dev/null)
  if [ -n "$NEW_ID" ]; then
    echo "  ℹ 新文章 ID: $NEW_ID"

    echo
    echo "▸ PUT /admin/articles/$NEW_ID (更新文章)"
    resp=$(curl -s -w "\n%{http_code}" -X PUT "$BASE/admin/articles/$NEW_ID" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{"title":"E2E Test Article (Updated)"}')
    status=$(echo "$resp" | tail -1)
    assert_status "更新文章" "200" "$status" ""

    echo
    echo "▸ DELETE /admin/articles/$NEW_ID"
    resp=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE/admin/articles/$NEW_ID" \
      -H "Authorization: Bearer $TOKEN")
    status=$(echo "$resp" | tail -1)
    assert_status "删除文章" "204" "$status" ""
  fi

  # ── 14. 未认证访问管理接口 ──
  echo
  echo "▸ GET /admin/stats (无 Token)"
  resp=$(curl -s -w "\n%{http_code}" "$BASE/admin/stats")
  status=$(echo "$resp" | tail -1)
  assert_status "无认证访问管理接口返回 401" "401" "$status" ""
fi

# ── 15. CORS ──
echo
echo "▸ OPTIONS /articles (CORS preflight)"
resp=$(curl -s -w "\n%{http_code}" -X OPTIONS "$BASE/articles" \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: GET")
status=$(echo "$resp" | tail -1)
assert_status "CORS preflight 返回 204" "204" "$status" ""

# ── 结果 ──
echo
echo "========================================="
echo " 结果: ${GREEN}$PASS passed${NC}, ${RED}$FAIL failed${NC}"
echo "========================================="
if [ "$FAIL" -gt 0 ]; then
  exit 1
fi
