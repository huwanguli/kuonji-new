<template>
  <div class="shell">
    <!-- 顶部导航栏 -->
    <header class="top-nav">
      <div class="container nav-inner">
        <RouterLink to="/" class="nav-brand" aria-label="回到首页">
          <span class="nav-brand-name">Kuonji</span>
          <span class="nav-brand-sub">Blog</span>
        </RouterLink>
        <nav class="nav-links" aria-label="主导航">
          <RouterLink
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            class="nav-link"
            :class="{ active: isActive(item.to) }"
          >
            {{ item.label }}
          </RouterLink>
          <button class="theme-toggle" @click="toggle" :aria-label="isDark ? '切换亮色' : '切换暗色'">
            <svg v-if="!isDark" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>
            <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>
          </button>
        </nav>
      </div>
    </header>

    <!-- 主体 -->
    <main class="site-main">
      <RouterView />
    </main>

    <!-- 底部 -->
    <footer class="site-footer">
      <div class="container footer-inner">
        <div class="footer-copy">© {{ year }} Kuonji · 学习记录与技术分享</div>
        <div class="footer-icons">
          <a href="#" class="footer-icon" aria-label="GitHub">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"/></svg>
          </a>
          <a href="#" class="footer-icon" aria-label="Twitter">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4l11.7 16L20 4H4z"/><path d="M4 4l8 8"/><path d="M12 12l8-8"/></svg>
          </a>
          <a href="#" class="footer-icon" aria-label="RSS">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M4 11a9 9 0 0 1 9 9"/><path d="M4 4a16 16 0 0 1 16 16"/><circle cx="5" cy="19" r="1"/></svg>
          </a>
          <a href="#" class="footer-icon" aria-label="Email">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="4" width="20" height="16" rx="2"/><path d="M22 6l-10 7L2 6"/></svg>
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useTheme } from '@/utils/theme'
import { onMounted } from 'vue'

const route = useRoute()
const year = new Date().getFullYear()
const { isDark, init, toggle } = useTheme()

onMounted(init)

const navItems = [
  { to: '/', label: 'Home' },
  { to: '/categories', label: 'Archives' },
  { to: '/tags', label: 'Tags' },
  { to: '/series', label: 'About' }
]

function isActive(to: string): boolean {
  if (to === '/') return route.path === '/'
  return route.path.startsWith(to)
}
</script>

<style scoped>
.shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.top-nav {
  height: 48px;
  display: flex;
  align-items: center;
  background: var(--bg-nav);
  border-bottom: 1px solid var(--line);
  backdrop-filter: blur(12px);
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.nav-brand {
  display: flex;
  align-items: baseline;
  gap: 6px;
  flex-shrink: 0;
}
.nav-brand:hover .nav-brand-name { color: var(--ink); }
.nav-brand-name {
  font-size: 18px;
  font-weight: 600;
  color: var(--ink);
  letter-spacing: -0.3px;
  transition: color 0.15s;
}
.nav-brand-sub {
  font-size: 11px;
  color: var(--muted-light);
  letter-spacing: 0.5px;
}

.nav-links {
  display: flex;
  gap: 32px;
  font-size: 13px;
}

.nav-link {
  color: var(--muted-light);
  transition: color 0.15s ease;
}

.nav-link:hover,
.nav-link.active {
  color: var(--ink);
}

.theme-toggle {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--muted-light);
  padding: 4px;
  display: flex;
  align-items: center;
  transition: color 0.15s ease;
}
.theme-toggle:hover { color: var(--ink); }

.site-main {
  flex: 1;
}

.site-footer {
  border-top: 1px solid var(--line);
  padding: 22px 0;
  margin-top: 48px;
}

.footer-inner {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.footer-copy {
  font-size: 12px;
  color: var(--muted-light);
}

.footer-icons {
  display: flex;
  gap: 14px;
  align-items: center;
}

.footer-icon {
  color: var(--muted-light);
  display: flex;
  align-items: center;
  transition: color 0.15s ease;
}

.footer-icon:hover {
  color: var(--ink);
}

@media (max-width: 768px) {
  .nav-brand-sub { display: none; }
  .nav-links { gap: 20px; }
}

@media (max-width: 640px) {
  .nav-links { gap: 16px; font-size: 12px; }
  .nav-brand-name { font-size: 16px; }
}
</style>
