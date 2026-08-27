<template>
  <SidebarLayout>
    <header class="content-header">
      <h1 class="page-title">标签</h1>
      <p class="page-desc">{{ tags.length }} 个标签</p>
    </header>
    <div v-if="tags.length" class="tag-cloud">
      <RouterLink v-for="tag in tags" :key="tag.id" :to="`/tags/${tag.id}`" class="tag">#{{ tag.name }}</RouterLink>
    </div>
    <p v-else class="empty-state">还没有标签。</p>
  </SidebarLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type Tag } from '@/utils/api'
import SidebarLayout from '@/components/SidebarLayout.vue'

const tags = ref<Tag[]>([])
onMounted(async () => {
  try { tags.value = (await api.getTags()).data } catch (e) { console.error(e) }
})
</script>

<style scoped>
.content-header { margin-bottom: 24px; }
.page-title { font-size: 32px; font-weight: 500; color: #1f1f1f; line-height: 1.25; letter-spacing: -0.3px; margin-bottom: 8px; }
.page-desc { font-size: 13px; color: var(--muted); }
.tag-cloud { display: flex; flex-wrap: wrap; gap: 8px; }
</style>
