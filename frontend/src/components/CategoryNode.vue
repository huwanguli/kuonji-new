<template>
  <div class="category-node">
    <RouterLink :to="`/categories/${category.id}`" class="category-row">
      <div class="category-left">
        <span class="dot" :class="getDotColor(category.name)"></span>
        <span class="category-name">{{ category.name }}</span>
      </div>
      <span class="category-count">{{ totalArticles }} 篇文章</span>
    </RouterLink>
    <div v-if="category.children?.length" class="category-children">
      <CategoryNode
        v-for="child in category.children"
        :key="child.id"
        :category="child"
        :depth="depth + 1"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Category } from '@/utils/api'
import { getDotColor } from '@/utils/colors'

const props = defineProps<{ category: Category; depth: number }>()

const totalArticles = computed(() => {
  let n = 0
  if (props.category.children?.length) {
    for (const c of props.category.children) {
      n += 1
      if (c.children?.length) n += c.children.length
    }
  }
  return n || 0
})
</script>

<style scoped>
.category-node { display: flex; flex-direction: column; }

.category-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 0;
  border-bottom: 1px solid var(--line);
  color: var(--ink);
}
.category-row:hover { color: var(--ink); }

.category-left { display: flex; align-items: center; gap: 10px; }
.category-name { font-size: 14px; font-weight: 500; }
.category-count { font-size: 12px; color: var(--muted-xlight); }
.category-children { margin-left: 1.6rem; display: flex; flex-direction: column; }
</style>
