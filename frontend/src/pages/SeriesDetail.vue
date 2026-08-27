<template>
  <SidebarLayout>
    <header class="content-header">
      <h1 class="page-title">{{ series?.name ?? '系列' }}</h1>
      <p v-if="series?.description" class="page-desc">{{ series.description }}</p>
      <p class="page-count">{{ articles.length }} 篇文章</p>
    </header>
    <div v-if="articles.length" class="chapter-list">
      <RouterLink
        v-for="(article, index) in articles"
        :key="article.id"
        :to="`/articles/${article.id}`"
        class="chapter-row"
      >
        <span class="chapter-index">CH.{{ String(index + 1).padStart(2, '0') }}</span>
        <span class="chapter-info">
          <span class="chapter-title">{{ article.title }}</span>
          <span class="chapter-meta">{{ formatDate(article.published_at) }} · {{ readInfo(article.content ?? '') }}</span>
        </span>
        <span class="chapter-arrow">→</span>
      </RouterLink>
    </div>
    <p v-else class="empty-state">该系列还没有文章。</p>
  </SidebarLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { api, type Article, type Series } from '@/utils/api'
import { formatDate, readInfo } from '@/utils/format'
import SidebarLayout from '@/components/SidebarLayout.vue'

const route = useRoute()
const series = ref<Series | null>(null)
const articles = ref<Article[]>([])

async function load() {
  try {
    const res = await api.getSeriesDetail(route.params.id as string)
    series.value = res.data
    articles.value = res.data.articles || []
  } catch (e) { console.error(e) }
}

watch(() => route.params.id, load)
onMounted(load)
</script>

<style scoped>
.content-header { margin-bottom: 24px; }
.page-title { font-size: 32px; font-weight: 500; color: #1f1f1f; line-height: 1.25; letter-spacing: -0.3px; margin-bottom: 8px; }
.page-desc { font-size: 13px; color: var(--muted); margin-bottom: 4px; }
.page-count { font-size: 12px; color: var(--muted-light); }

.chapter-list { display: flex; flex-direction: column; }
.chapter-row {
  display: flex; align-items: center; gap: 14px;
  padding: 16px 0; border-bottom: 1px solid var(--line);
  color: var(--ink); transition: background 0.15s ease;
}
.chapter-row:hover { color: var(--ink); }
.chapter-index { font-size: 11px; font-family: var(--font-mono); color: var(--muted-light); letter-spacing: 1.8px; flex-shrink: 0; min-width: 3.5em; }
.chapter-info { flex: 1; display: flex; flex-direction: column; gap: 2px; min-width: 0; }
.chapter-title { font-size: 14px; font-weight: 500; color: #2a2a2a; }
.chapter-meta { font-size: 11px; color: var(--muted-xlight); }
.chapter-arrow { color: var(--muted-light); opacity: 0; transform: translateX(-6px); transition: opacity 0.15s ease, transform 0.15s ease; }
.chapter-row:hover .chapter-arrow { opacity: 1; transform: none; }
</style>
