<template>
  <SidebarLayout>
    <header class="content-header">
      <h1 class="page-title">#{{ tag?.name ?? '标签' }}</h1>
    </header>
    <div v-if="articles.length">
      <ArticleCard v-for="article in articles" :key="article.id" :article="article" />
    </div>
    <p v-else class="empty-state">该标签下还没有文章。</p>
  </SidebarLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { api, type Article, type Tag } from '@/utils/api'
import SidebarLayout from '@/components/SidebarLayout.vue'
import ArticleCard from '@/components/ArticleCard.vue'

const route = useRoute()
const articles = ref<Article[]>([])
const allTags = ref<Tag[]>([])
const tag = computed(() => allTags.value.find(t => t.id === Number(route.params.id)))

async function load() {
  try {
    const [tagRes, artRes] = await Promise.all([
      api.getTags(),
      api.getArticles({ tag_id: route.params.id as string, page_size: 50 })
    ])
    allTags.value = tagRes.data
    articles.value = artRes.data.articles
  } catch (e) { console.error(e) }
}

watch(() => route.params.id, load)
onMounted(load)
</script>

<style scoped>
.content-header { margin-bottom: 24px; }
.page-title { font-size: 32px; font-weight: 500; color: #1f1f1f; line-height: 1.25; letter-spacing: -0.3px; margin-bottom: 8px; }
</style>
