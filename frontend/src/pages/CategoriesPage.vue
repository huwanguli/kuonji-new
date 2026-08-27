<template>
  <SidebarLayout>
    <header class="content-header">
      <h1 class="page-title">分类</h1>
      <p class="page-desc">{{ total }} 个分类</p>
    </header>
    <div v-if="categories.length" class="category-tree">
      <CategoryNode v-for="cat in categories" :key="cat.id" :category="cat" :depth="0" />
    </div>
    <p v-else class="empty-state">还没有分类 —— 等作者来建立知识体系。</p>
  </SidebarLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { api, type Category } from '@/utils/api'
import { treeCount } from '@/utils/format'
import SidebarLayout from '@/components/SidebarLayout.vue'
import CategoryNode from '@/components/CategoryNode.vue'

const categories = ref<Category[]>([])
const total = computed(() => treeCount(categories.value))

onMounted(async () => {
  try {
    const res = await api.getCategories()
    categories.value = res.data
  } catch (e) { console.error(e) }
})
</script>

<style scoped>
.content-header { margin-bottom: 24px; }
.page-title { font-size: 32px; font-weight: 500; color: #1f1f1f; line-height: 1.25; letter-spacing: -0.3px; margin-bottom: 8px; }
.page-desc { font-size: 13px; color: var(--muted); }
.category-tree { display: flex; flex-direction: column; }
</style>
