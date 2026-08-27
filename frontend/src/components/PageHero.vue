<template>
  <section class="page-hero fade-in">
    <!-- 头像/图片区域 -->
    <div class="hero-banner">
      <div class="hero-banner-inner">
        <div class="hero-greeting">
          <h1 class="hero-name">Kuonji</h1>
          <p class="hero-tagline">积跬步，至千里</p>
        </div>
        <div class="hero-decor" aria-hidden="true">
          <div class="hero-decor-dot" style="background:var(--cat-blue)"></div>
          <div class="hero-decor-dot" style="background:var(--cat-pink)"></div>
          <div class="hero-decor-dot" style="background:var(--cat-purple)"></div>
          <div class="hero-decor-dot" style="background:var(--cat-green)"></div>
        </div>
      </div>
    </div>

    <!-- 一句话 / 简介 -->
    <p class="hero-quote">Building in public. Notes, tutorials, and wandering thoughts.</p>

    <!-- 小统计条 -->
    <div class="hero-stats">
      <span class="hero-stat"><strong>{{ stats.posts }}</strong> 文章</span>
      <span class="hero-stat-sep">·</span>
      <span class="hero-stat"><strong>{{ stats.categories }}</strong> 分类</span>
      <span class="hero-stat-sep">·</span>
      <span class="hero-stat"><strong>{{ stats.series }}</strong> 系列</span>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/utils/api'

const stats = ref({ posts: 0, categories: 0, series: 0 })

onMounted(async () => {
  try {
    const [artRes, catRes, serRes] = await Promise.all([
      api.getArticles({ page_size: 1 }),
      api.getCategories(),
      api.getSeries()
    ])
    stats.value = {
      posts: artRes.data.total,
      categories: catRes.data.length,
      series: serRes.data.length
    }
  } catch {}
})
</script>

<style scoped>
.page-hero {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 48px 0 32px;
  gap: 20px;
}

/* 头图横幅 */
.hero-banner {
  width: 100%;
  max-width: 680px;
  border-radius: var(--r-lg);
  background: linear-gradient(135deg, var(--cat-blue) 0%, var(--cat-purple) 50%, var(--cat-pink) 100%);
  padding: 40px 36px;
  box-shadow: 0 8px 32px rgba(184,212,227,0.25);
  overflow: hidden;
  position: relative;
}

.hero-banner-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.hero-greeting { position: relative; z-index: 1; }

.hero-name {
  font-size: 28px;
  font-weight: 600;
  color: #fff;
  letter-spacing: -0.5px;
  margin-bottom: 6px;
  text-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.hero-tagline {
  font-size: 14px;
  color: rgba(255,255,255,0.8);
  letter-spacing: 1px;
}

/* 装饰圆点 */
.hero-decor {
  display: flex;
  gap: 10px;
}
.hero-decor-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  opacity: 0.5;
}
.hero-decor-dot:nth-child(1) { opacity: 0.7; transform: scale(1.2); }
.hero-decor-dot:nth-child(2) { opacity: 0.55; }
.hero-decor-dot:nth-child(3) { opacity: 0.4; }
.hero-decor-dot:nth-child(4) { opacity: 0.3; }

/* 一句话 */
.hero-quote {
  font-size: 15px;
  color: var(--ink-soft);
  text-align: center;
  max-width: 480px;
  line-height: 1.6;
}

/* 统计条 */
.hero-stats {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 13px;
  color: var(--muted);
}
.hero-stat strong {
  color: var(--ink);
  font-weight: 600;
}
.hero-stat-sep {
  color: var(--line);
}

@media (max-width: 768px) {
  .page-hero { padding: 32px 0 24px; gap: 16px; }
  .hero-banner { padding: 28px 24px; }
  .hero-name { font-size: 22px; }
}
</style>
