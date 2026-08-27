<template>
  <SidebarLayout>
    <template #hero>
      <PageHero />
    </template>

    <header class="content-header">
      <h1 class="page-title" style="display:none">首页</h1>
    </header>

    <ArticleCard
      v-for="article in articles"
      :key="article.id"
      :article="article"
    />

    <p v-if="!loading && !articles.length" class="empty-state">
      这里还是空的 —— 等作者写下第一篇。
    </p>

    <!-- PAGINATION -->
    <nav v-if="totalPages > 1" class="pagination" aria-label="分页">
      <button :disabled="page <= 1" class="page-btn" @click="changePage(page - 1)">
        ←&nbsp;&nbsp;Newer Posts
      </button>
      <span class="page-info">Page {{ page }} of {{ totalPages }}</span>
      <button :disabled="page >= totalPages" class="page-btn" @click="changePage(page + 1)">
        Older Posts&nbsp;&nbsp;→
      </button>
    </nav>
  </SidebarLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api, type Article } from '@/utils/api'
import SidebarLayout from '@/components/SidebarLayout.vue'
import ArticleCard from '@/components/ArticleCard.vue'
import PageHero from '@/components/PageHero.vue'

const route = useRoute()
const router = useRouter()

const articles = ref<Article[]>([])
const total = ref(0)
const loading = ref(true)
const page = ref(parseInt(route.query.page as string) || 1)
const pageSize = 10
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))

async function load() {
  loading.value = true
  try {
    const res = await api.getArticles({ page: page.value, page_size: pageSize })
    articles.value = res.data.articles
    total.value = res.data.total
  } catch (e) {
    console.error('Failed to load articles:', e)
  } finally {
    loading.value = false
  }
}

function changePage(p: number) {
  if (p < 1 || p > totalPages.value) return
  page.value = p
  router.push({ query: p > 1 ? { page: String(p) } : {} })
}

watch(() => route.query.page, (val) => {
  page.value = parseInt(val as string) || 1
  load()
})

onMounted(load)
</script>

<style scoped>
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 80px;
  padding: 36px 0;
  font-size: 13px;
}
.page-btn {
  background: none;
  border: none;
  color: var(--muted-xlight);
  font-family: inherit;
  font-size: 13px;
  cursor: pointer;
  padding: 0;
  transition: color 0.15s ease;
}
.page-btn:hover:not(:disabled) { color: var(--ink); }
.page-btn:disabled { opacity: 0.3; cursor: default; }
.page-info { color: var(--muted); }

@media (max-width: 640px) {
  .pagination { gap: 24px; flex-wrap: wrap; }
}
</style>
