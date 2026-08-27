<template>
  <div v-if="!isLoggedIn()" class="auth-redirect">
    <p>请先登录</p>
    <RouterLink to="/admin/login" class="btn btn-primary">去登录</RouterLink>
  </div>
  <div v-else class="admin-shell">
    <!-- 侧边导航 -->
    <aside class="admin-sidebar">
      <div class="admin-brand">
        <RouterLink to="/admin" class="admin-brand-link">Kuonji</RouterLink>
        <span class="admin-brand-sub">Admin</span>
      </div>
      <nav class="admin-nav">
        <RouterLink to="/admin" class="admin-nav-item" :class="{ active: route.path === '/admin' }">
          <span class="admin-nav-icon">📊</span> 仪表盘
        </RouterLink>
        <RouterLink to="/admin/articles" class="admin-nav-item" :class="{ active: route.path.startsWith('/admin/articles') }">
          <span class="admin-nav-icon">📝</span> 文章管理
        </RouterLink>
        <RouterLink to="/admin/categories" class="admin-nav-item" :class="{ active: route.path === '/admin/categories' }">
          <span class="admin-nav-icon">📂</span> 分类管理
        </RouterLink>
        <RouterLink to="/admin/tags" class="admin-nav-item" :class="{ active: route.path === '/admin/tags' }">
          <span class="admin-nav-icon">🏷️</span> 标签管理
        </RouterLink>
        <RouterLink to="/admin/series" class="admin-nav-item" :class="{ active: route.path === '/admin/series' }">
          <span class="admin-nav-icon">📚</span> 系列管理
        </RouterLink>
        <RouterLink to="/admin/comments" class="admin-nav-item" :class="{ active: route.path === '/admin/comments' }">
          <span class="admin-nav-icon">💬</span> 评论管理
        </RouterLink>
      </nav>
      <div class="admin-sidebar-footer">
        <RouterLink to="/" class="admin-nav-item">← 返回前台</RouterLink>
        <button class="admin-nav-item admin-logout" @click="logout">退出登录</button>
      </div>
    </aside>

    <!-- 主内容 -->
    <main class="admin-main">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { isLoggedIn, clearAuth } from '@/utils/auth'

const route = useRoute()
const router = useRouter()

function logout() {
  clearAuth()
  router.push('/admin/login')
}
</script>

<style scoped>
.auth-redirect {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  gap: 16px;
  color: var(--muted);
}

.admin-shell {
  display: grid;
  grid-template-columns: 220px 1fr;
  min-height: 100vh;
}

/* 侧边栏 */
.admin-sidebar {
  background: #1f1f1f;
  color: #ccc;
  padding: 20px 0;
  display: flex;
  flex-direction: column;
  position: sticky;
  top: 0;
  height: 100vh;
  overflow-y: auto;
}

.admin-brand {
  padding: 0 20px 20px;
  border-bottom: 1px solid #333;
  margin-bottom: 12px;
}
.admin-brand-link {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
}
.admin-brand-sub {
  font-size: 11px;
  color: #666;
  margin-left: 8px;
}

.admin-nav {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 0 8px;
}

.admin-nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 6px;
  font-size: 13px;
  color: #aaa;
  transition: background 0.15s, color 0.15s;
  cursor: pointer;
  border: none;
  background: none;
  font-family: inherit;
  text-decoration: none;
  width: 100%;
  text-align: left;
}
.admin-nav-item:hover { background: #2a2a2a; color: #fff; }
.admin-nav-item.active { background: #333; color: #fff; }

.admin-nav-icon { font-size: 15px; }

.admin-sidebar-footer {
  padding: 12px 8px 0;
  border-top: 1px solid #333;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.admin-logout { color: #e55; }
.admin-logout:hover { background: #3a1a1a; }

/* 主内容 */
.admin-main {
  padding: 32px 40px;
  background: var(--bg);
  overflow-y: auto;
}

@media (max-width: 768px) {
  .admin-shell { grid-template-columns: 1fr; }
  .admin-sidebar {
    position: relative;
    height: auto;
    flex-direction: row;
    flex-wrap: wrap;
    padding: 12px;
  }
  .admin-nav { flex-direction: row; flex-wrap: wrap; gap: 4px; }
  .admin-sidebar-footer { display: none; }
  .admin-main { padding: 20px 16px; }
}
</style>
