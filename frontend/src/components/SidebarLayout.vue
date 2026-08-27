<template>
  <div class="container">
    <slot name="hero" />

    <div class="main-grid">
      <!-- SIDEBAR -->
      <aside class="sidebar">
        <!-- CATEGORIES -->
        <div class="sidebar-section">
          <h3 class="section-label">CATEGORIES</h3>
          <div class="sidebar-list">
            <RouterLink
              v-for="cat in sidebarData.categories"
              :key="cat.id"
              :to="`/categories/${cat.id}`"
              class="sidebar-item"
            >
              <div class="sidebar-item-left">
                <span class="dot" :class="getDotColor(cat.name)"></span>
                <span class="sidebar-item-name">{{ cat.name }}</span>
              </div>
              <span class="sidebar-item-count">{{ cat.article_count ?? 0 }}</span>
            </RouterLink>
          </div>
        </div>

        <!-- SERIES -->
        <div class="sidebar-section">
          <h3 class="section-label">SERIES</h3>
          <div class="sidebar-list">
            <RouterLink
              v-for="s in sidebarData.series"
              :key="s.id"
              :to="`/series/${s.id}`"
              class="sidebar-series"
            >
              <div class="sidebar-series-left">
                <span class="series-dot" :class="getDotColor(s.name)"></span>
                <span class="sidebar-item-name">{{ s.name }}</span>
              </div>
              <span class="sidebar-item-count">{{ s.articleCount }}</span>
            </RouterLink>
          </div>
        </div>

        <!-- TAGS -->
        <div class="sidebar-section">
          <h3 class="section-label">TAGS</h3>
          <div class="sidebar-tags">
            <RouterLink
              v-for="tag in sidebarData.tags"
              :key="tag.id"
              :to="`/tags/${tag.id}`"
              class="tag"
            >
              #{{ tag.name }}
            </RouterLink>
          </div>
        </div>
      </aside>

      <!-- MAIN CONTENT SLOT -->
      <main class="content reading">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type Category, type Tag } from '@/utils/api'
import { getDotColor } from '@/utils/colors'

interface SeriesItem { id: number; name: string; articleCount: number }

const sidebarData = ref<{
  categories: Category[]
  tags: Tag[]
  series: SeriesItem[]
}>({ categories: [], tags: [], series: [] })

onMounted(async () => {
  try {
    const [catRes, tagRes, seriesRes] = await Promise.all([
      api.getCategories(),
      api.getTags(),
      api.getSeries()
    ])
    sidebarData.value.categories = catRes.data
    sidebarData.value.tags = tagRes.data

    // 为每个系列拉取文章数
    const seriesCounts = await Promise.all(
      seriesRes.data.map(async (s) => {
        try {
          const detail = await api.getSeriesDetail(s.id)
          return { id: s.id, name: s.name, articleCount: detail.data.articles?.length ?? 0 }
        } catch {
          return { id: s.id, name: s.name, articleCount: 0 }
        }
      })
    )
    sidebarData.value.series = seriesCounts
  } catch (e) {
    console.error('Failed to load sidebar data:', e)
  }
})
</script>

<style scoped>
.main-grid {
  display: grid;
  grid-template-columns: var(--sidebar-w) 1fr;
  gap: var(--gap);
  padding-top: 8px;
}

.sidebar { min-width: 0; }
.sidebar-section { margin-bottom: 48px; }
.sidebar-section:last-child { margin-bottom: 0; }
.sidebar-list { display: flex; flex-direction: column; gap: 14px; }

.sidebar-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 13px;
  color: var(--ink);
}
.sidebar-item-left { display: flex; align-items: center; gap: 10px; }
.sidebar-item-name { color: var(--ink); }
.sidebar-item-count { color: var(--muted-xlight); font-size: 12px; }

.sidebar-series {
  background: rgba(213,208,232,0.22);
  border-radius: var(--r-sm);
  padding: 12px 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  transition: background 0.15s ease;
}
.sidebar-series:nth-child(2) { background: rgba(184,212,227,0.22); }
.sidebar-series:nth-child(3) { background: rgba(200,230,208,0.30); }
.sidebar-series:nth-child(4) { background: rgba(245,213,213,0.32); }
.sidebar-series:nth-child(5) { background: rgba(232,216,232,0.32); }
.sidebar-series:hover { background: rgba(213,208,232,0.35); }
.sidebar-series-left { display: flex; align-items: center; gap: 10px; }

.series-dot { width: 10px; height: 10px; border-radius: 3px; }
.series-dot.dot-blue   { background: var(--cat-blue); }
.series-dot.dot-pink   { background: var(--cat-pink); }
.series-dot.dot-purple { background: var(--cat-purple); }
.series-dot.dot-green  { background: var(--cat-green); }
.series-dot.dot-amber  { background: var(--cat-amber); }
.series-dot.dot-rose   { background: var(--cat-rose); }

.sidebar-tags { display: flex; flex-wrap: wrap; gap: 8px; }

.content { min-width: 0; }

@media (max-width: 960px) {
  .main-grid { grid-template-columns: 1fr; }
  .sidebar {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 24px;
  }
  .sidebar-section { margin-bottom: 0; }
}
@media (max-width: 640px) {
  .sidebar { grid-template-columns: 1fr; }
}
</style>
