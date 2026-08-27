<template>
  <SidebarLayout>
    <header class="content-header">
      <h1 class="page-title">系列</h1>
      <p class="page-desc">{{ seriesList.length }} 个系列</p>
    </header>
    <div v-if="seriesList.length" class="series-grid">
      <RouterLink
        v-for="s in seriesList"
        :key="s.id"
        :to="`/series/${s.id}`"
        class="series-card"
      >
        <div class="series-card-dot" :class="getDotColor(s.name)"></div>
        <div class="series-card-body">
          <span class="series-card-name">{{ s.name }}</span>
          <span v-if="s.description" class="series-card-desc">{{ s.description }}</span>
          <span class="series-card-count">{{ s.articleCount }} 篇文章</span>
        </div>
      </RouterLink>
    </div>
    <p v-else class="empty-state">还没有系列 —— 等作者组织教程。</p>
  </SidebarLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/utils/api'
import { getDotColor } from '@/utils/colors'
import SidebarLayout from '@/components/SidebarLayout.vue'

interface SeriesItem { id: number; name: string; description?: string; articleCount: number }
const seriesList = ref<SeriesItem[]>([])

onMounted(async () => {
  try {
    const res = await api.getSeries()
    seriesList.value = await Promise.all(
      res.data.map(async (s) => {
        try {
          const detail = await api.getSeriesDetail(s.id)
          return { id: s.id, name: s.name, description: s.description, articleCount: detail.data.articles?.length ?? 0 }
        } catch { return { id: s.id, name: s.name, description: s.description, articleCount: 0 } }
      })
    )
  } catch (e) { console.error(e) }
})
</script>

<style scoped>
.content-header { margin-bottom: 24px; }
.page-title { font-size: 32px; font-weight: 500; color: #1f1f1f; line-height: 1.25; letter-spacing: -0.3px; margin-bottom: 8px; }
.page-desc { font-size: 13px; color: var(--muted); }

.series-grid { display: flex; flex-direction: column; gap: 10px; }
.series-card {
  display: flex; align-items: center; gap: 12px;
  padding: 12px 14px; border-radius: var(--r-sm);
  background: rgba(213,208,232,0.22); transition: background 0.15s ease; color: var(--ink);
}
.series-card:nth-child(2) { background: rgba(184,212,227,0.22); }
.series-card:nth-child(3) { background: rgba(200,230,208,0.30); }
.series-card:nth-child(4) { background: rgba(245,213,213,0.32); }
.series-card:nth-child(5) { background: rgba(232,216,232,0.32); }
.series-card:hover { background: rgba(213,208,232,0.35); }

.series-card-dot { width: 10px; height: 10px; border-radius: 3px; flex-shrink: 0; }
.series-card-dot.dot-blue   { background: var(--cat-blue); }
.series-card-dot.dot-pink   { background: var(--cat-pink); }
.series-card-dot.dot-purple { background: var(--cat-purple); }
.series-card-dot.dot-green  { background: var(--cat-green); }
.series-card-dot.dot-amber  { background: var(--cat-amber); }
.series-card-dot.dot-rose   { background: var(--cat-rose); }

.series-card-body { display: flex; flex-direction: column; gap: 2px; }
.series-card-name { font-size: 14px; font-weight: 500; color: var(--ink); }
.series-card-desc { font-size: 12px; color: var(--muted); line-height: 1.5; }
.series-card-count { font-size: 12px; color: var(--muted-light); }
</style>
