<template>
  <SidebarLayout>
    <header class="content-header">
      <h1 class="page-title">搜索</h1>
      <p v-if="searched" class="page-desc">{{ total }} 条结果</p>
    </header>
    <form class="search-form" role="search" @submit.prevent="doSearch">
      <input
        v-model="input"
        type="search"
        placeholder="输入关键词搜索文章…"
        class="input search-input"
        aria-label="搜索关键词"
      />
      <button type="submit" class="btn btn-primary">搜索</button>
    </form>
    <div v-if="articles.length">
      <ArticleCard v-for="article in articles" :key="article.id" :article="article" />
    </div>
    <p v-else-if="searched" class="empty-state">没有找到与「{{ keyword }}」相关的文章。</p>
  </SidebarLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api, type Article } from '@/utils/api'
import SidebarLayout from '@/components/SidebarLayout.vue'
import ArticleCard from '@/components/ArticleCard.vue'

const route = useRoute()
const router = useRouter()

const input = ref((route.query.q as string) || '')
const keyword = computed(() => input.value.trim())
const articles = ref<Article[]>([])
const total = ref(0)
const searched = ref(false)

async function load() {
  if (!keyword.value) { searched.value = false; articles.value = []; return }
  try {
    const res = await api.search(keyword.value)
    articles.value = res.data.articles
    total.value = res.data.total
    searched.value = true
  } catch (e) { console.error(e) }
}

function doSearch() {
  const q = input.value.trim()
  if (!q) return
  router.push({ path: '/search', query: { q } })
}

watch(() => route.query.q, (val) => {
  input.value = (val as string) || ''
  load()
})

onMounted(load)
</script>

<style scoped>
.content-header { margin-bottom: 16px; }
.page-title { font-size: 32px; font-weight: 500; color: #1f1f1f; line-height: 1.25; letter-spacing: -0.3px; margin-bottom: 8px; }
.page-desc { font-size: 13px; color: var(--muted); }

.search-form { display: flex; gap: 10px; margin-bottom: 24px; }
.search-input { flex: 1; font-size: 13px; padding: 10px 14px; }

@media (max-width: 640px) { .search-form { flex-direction: column; } }
</style>
