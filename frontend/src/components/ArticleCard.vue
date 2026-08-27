<template>
  <article class="post-card">
    <RouterLink :to="`/articles/${article.id}`" class="post-link">
      <div class="post-category">
        <span class="dot" :class="dotColor"></span>
        <span class="post-category-label">{{ categoryName }} · {{ primaryTag }}</span>
      </div>
      <h2 class="post-title">{{ article.title }}</h2>
      <p v-if="article.summary" class="post-summary">{{ article.summary }}</p>
      <div class="post-meta">
        <span>{{ formatDate(article.published_at) }}</span>
        <span class="meta-sep">·</span>
        <span>{{ readInfo(article.content ?? '') }}</span>
      </div>
      <div class="post-tags">
        <span v-for="tag in article.tags" :key="tag.id" class="tag tag-sm">#{{ tag.name }}</span>
      </div>
    </RouterLink>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Article } from '@/utils/api'
import { formatDate, readInfo } from '@/utils/format'
import { getDotColor } from '@/utils/colors'

const props = defineProps<{ article: Article }>()

const categoryName = computed(() => props.article.category?.name ?? 'UNCATEGORIZED')
const primaryTag = computed(() => props.article.tags?.[0]?.name?.toUpperCase() ?? 'NOTE')
const dotColor = computed(() => getDotColor(props.article.category?.name ?? ''))
</script>

<style scoped>
.post-card {
  padding: 18px 0 22px;
  border-bottom: 1px solid var(--line);
}
.post-card:last-child { border-bottom: none; }

.post-link { display: block; color: inherit; }
.post-link:hover { color: inherit; }

.post-category {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 10px;
  color: var(--muted-light);
  letter-spacing: 1.8px;
  margin-bottom: 10px;
}
.post-category-label { text-transform: uppercase; }

.post-title {
  font-size: 19px;
  font-weight: 500;
  color: #2a2a2a;
  line-height: 1.35;
  margin-bottom: 10px;
  transition: color 0.15s ease;
}
.post-card:hover .post-title { color: #1a1a1a; }

.post-summary {
  font-size: 13px;
  color: var(--muted);
  line-height: 1.65;
  margin-bottom: 14px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.post-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: var(--muted-xlight);
  margin-bottom: 12px;
}
.meta-sep { color: var(--line); }

.post-tags { display: flex; gap: 6px; flex-wrap: wrap; }
</style>
