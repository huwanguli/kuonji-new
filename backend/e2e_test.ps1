# ============================================================
# Kuonji Blog E2E 测试脚本 (PowerShell)
# 针对运行中的后端服务 (localhost:8080)
# 使用方式: pwsh e2e_test.ps1
# ============================================================

$BASE = "http://localhost:8080/api/v1"
$script:PASS = 0
$script:FAIL = 0

function Assert-Status($desc, $expected, $actual, $body = "") {
    if ($actual -eq $expected) {
        Write-Host "  [PASS] $desc (HTTP $actual)" -ForegroundColor Green
        $script:PASS++
    } else {
        Write-Host "  [FAIL] $desc : expected $expected, got $actual" -ForegroundColor Red
        if ($body) { Write-Host "    Response: $($body.Substring(0, [Math]::Min(200, $body.Length)))" }
        $script:FAIL++
    }
}

function Invoke-Api($method, $path, $body = $null, $token = $null) {
    $uri = "$BASE$path"
    $headers = @{ "Content-Type" = "application/json" }
    if ($token) { $headers["Authorization"] = "Bearer $token" }
    try {
        $params = @{ Uri = $uri; Method = $method; Headers = $headers; ErrorAction = "Stop"; StatusCodeVariable = "sc" }
        if ($body) { $params.Body = $body }
        $resp = Invoke-RestMethod @params
        return @{ Status = $sc; Body = ($resp | ConvertTo-Json -Depth 10 -Compress) }
    } catch {
        $status = 0
        if ($_.Exception.Response) {
            $status = [int]$_.Exception.Response.StatusCode
        }
        $errBody = ""
        try {
            $stream = $_.Exception.Response.GetResponseStream()
            $reader = New-Object System.IO.StreamReader($stream)
            $errBody = $reader.ReadToEnd()
        } catch {}
        return @{ Status = $status; Body = $errBody }
    }
}

Write-Host "========================================="
Write-Host " Kuonji Blog E2E 测试 (PowerShell)"
Write-Host "========================================="
Write-Host

# 1. 文章列表
Write-Host "GET /articles"
$r = Invoke-Api "GET" "/articles?page=1&page_size=5"
Assert-Status "获取文章列表" 200 $r.Status $r.Body

# 2. 标签过滤
Write-Host "`nGET /articles?tag_id=22"
$r = Invoke-Api "GET" "/articles?tag_id=22&page_size=50"
Assert-Status "按标签过滤文章" 200 $r.Status $r.Body

# 3. 分类过滤
Write-Host "`nGET /articles?category_id=11"
$r = Invoke-Api "GET" "/articles?category_id=11"
Assert-Status "按分类过滤文章" 200 $r.Status $r.Body

# 4. 文章详情
Write-Host "`nGET /articles/9"
$r = Invoke-Api "GET" "/articles/9"
Assert-Status "获取文章详情" 200 $r.Status $r.Body

# 5. 文章 404
Write-Host "`nGET /articles/99999"
$r = Invoke-Api "GET" "/articles/99999"
Assert-Status "不存在的文章返回 404" 404 $r.Status

# 6. 分类列表
Write-Host "`nGET /categories"
$r = Invoke-Api "GET" "/categories"
Assert-Status "获取分类列表" 200 $r.Status $r.Body

# 7. 标签列表
Write-Host "`nGET /tags"
$r = Invoke-Api "GET" "/tags"
Assert-Status "获取标签列表" 200 $r.Status $r.Body

# 8. 系列列表
Write-Host "`nGET /series"
$r = Invoke-Api "GET" "/series"
Assert-Status "获取系列列表" 200 $r.Status $r.Body

# 9. 系列详情
Write-Host "`nGET /series/4"
$r = Invoke-Api "GET" "/series/4"
Assert-Status "获取系列详情" 200 $r.Status $r.Body

# 10. 搜索
Write-Host "`nGET /search?q=Go"
$r = Invoke-Api "GET" "/search?q=Go"
Assert-Status "搜索文章" 200 $r.Status $r.Body

Write-Host "`nGET /search?q= (空关键词)"
$r = Invoke-Api "GET" "/search?q="
Assert-Status "空搜索返回 400" 400 $r.Status

# 11. 评论
Write-Host "`nPOST /articles/9/comments"
$r = Invoke-Api "POST" "/articles/9/comments" '{"nickname":"E2E Test","email":"e2e@test.com","content":"E2E test comment"}'
Assert-Status "创建评论" 201 $r.Status $r.Body

Write-Host "`nGET /articles/9/comments"
$r = Invoke-Api "GET" "/articles/9/comments"
Assert-Status "获取评论列表" 200 $r.Status $r.Body

# 12. 登录
Write-Host "`nPOST /auth/login"
$r = Invoke-Api "POST" "/auth/login" '{"username":"admin","password":"admin123"}'
Assert-Status "管理员登录" 200 $r.Status $r.Body

$TOKEN = $null
if ($r.Status -eq 200) {
    $data = ($r.Body | ConvertFrom-Json).data
    $TOKEN = $data.token
    Write-Host "  Token 获取成功"

    # 13. 统计
    Write-Host "`nGET /admin/stats"
    $r = Invoke-Api "GET" "/admin/stats" $null $TOKEN
    Assert-Status "获取统计数据" 200 $r.Status $r.Body

    # 14. 创建文章
    Write-Host "`nPOST /admin/articles (创建文章)"
    $r = Invoke-Api "POST" "/admin/articles" '{"title":"E2E Test Article","content":"# E2E Test","summary":"E2E","category_id":11,"tag_ids":[22],"status":1}' $TOKEN
    Assert-Status "创建文章（直接发布）" 201 $r.Status $r.Body

    if ($r.Status -eq 201) {
        $NEW_ID = ($r.Body | ConvertFrom-Json).data.id
        Write-Host "  新文章 ID: $NEW_ID"

        Write-Host "`nPUT /admin/articles/$NEW_ID"
        $r = Invoke-Api "PUT" "/admin/articles/$NEW_ID" '{"title":"E2E Updated"}' $TOKEN
        Assert-Status "更新文章" 200 $r.Status $r.Body

        Write-Host "`nDELETE /admin/articles/$NEW_ID"
        $r = Invoke-Api "DELETE" "/admin/articles/$NEW_ID" $null $TOKEN
        Assert-Status "删除文章" 204 $r.Status
    }

    # 15. 标签管理
    $tagName = "e2e-test-$(Get-Random)"
    Write-Host "`nPOST /admin/tags (创建标签: $tagName)"
    $r = Invoke-Api "POST" "/admin/tags" "{`"name`":`"$tagName`"}" $TOKEN
    Assert-Status "创建标签" 201 $r.Status $r.Body

    if ($r.Status -eq 201) {
        $TAG_ID = ($r.Body | ConvertFrom-Json).data.id
        Write-Host "`nDELETE /admin/tags/$TAG_ID"
        $r = Invoke-Api "DELETE" "/admin/tags/$TAG_ID" $null $TOKEN
        Assert-Status "删除标签" 204 $r.Status
    }

    # 16. 上传接口
    Write-Host "`nPOST /admin/upload (无文件)"
    $r = Invoke-Api "POST" "/admin/upload" $null $TOKEN
    Assert-Status "无文件上传返回 400" 400 $r.Status

    Write-Host "`nPOST /admin/upload (无认证)"
    $r = Invoke-Api "POST" "/admin/upload"
    Assert-Status "无认证上传返回 401" 401 $r.Status

    # 17. 暗色模式（前端静态检查不需要后端测试）

    # 18. 无认证访问
    Write-Host "`nGET /admin/stats (无 Token)"
    $r = Invoke-Api "GET" "/admin/stats"
    Assert-Status "无认证返回 401" 401 $r.Status
}

# 结果
Write-Host
Write-Host "========================================="
Write-Host " 结果: $($script:PASS) passed, $($script:FAIL) failed" -ForegroundColor $(if ($script:FAIL -gt 0) { "Red" } else { "Green" })
Write-Host "========================================="

if ($script:FAIL -gt 0) { exit 1 }
