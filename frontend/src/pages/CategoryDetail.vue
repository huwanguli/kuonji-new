<template>
  <SidebarLayout>
    <header class="content-header">
      <h1 class="page-title">{{ category?.name ?? '分类' }}</h1>
      <p v-if="category?.description" class="page-desc">{{ category.description }}</p>
    </header>
    <div v-if="articles.length">
      <ArticleCard v-for="article in articles" :key="article.id" :article="article" />
    </div>
    <p v-else class="empty-state">该分类下还没有文章。</p>
  </SidebarLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { api, type Article, type Category } from '@/utils/api'
import SidebarLayout from '@/components/SidebarLayout.vue'
import ArticleCard from '@/components/ArticleCard.vue'

const route = useRoute()
const articles = ref<Article[]>([])
const categories = ref<Category[]>([])
const category = computed(() => {
  const id = Number(route.params.id)
  function find(list: Category[]): Category | undefined {
    for (const c of list) {
      if (c.id === id) return c
      if (c.children?.length) { const r = find(c.children); if (r) return r }
    }
  }
  return find(categories.value)
})

async function load() {
  try {
    const [catRes, artRes] = await Promise.all([
      api.getCategories(),
      api.getArticles({ category_id: route.params.id as string, page_size: 50 })
    ])
    categories.value = catRes.data
    articles.value = artRes.data.articles
  } catch (e) { console.error(e) }
}

watch(() => route.params.id, load)
onMounted(load)
</script>

<style scoped>
.content-header { margin-bottom: 24px; }
.page-title { font-size: 32px; font-weight: 500; color: #1f1f1f; line-height: 1.25; letter-spacing: -0.3px; margin-bottom: 8px; }
.page-desc { font-size: 13px; color: var(--muted); }
</style>
