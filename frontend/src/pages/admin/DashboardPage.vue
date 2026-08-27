<template>
  <div>
    <h1 class="page-heading">仪表盘</h1>
    <div class="stats-grid">
      <div class="stat-card" v-for="s in statItems" :key="s.label">
        <div class="stat-value">{{ s.value }}</div>
        <div class="stat-label">{{ s.label }}</div>
      </div>
    </div>
    <div class="quick-actions">
      <h2 class="section-sub">快捷操作</h2>
      <div class="action-row">
        <RouterLink to="/admin/articles/new" class="btn btn-primary">写新文章</RouterLink>
        <RouterLink to="/admin/articles" class="btn btn-ghost">管理文章</RouterLink>
        <RouterLink to="/admin/comments" class="btn btn-ghost">审核评论</RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { api } from '@/utils/api'

const stats = ref<any>(null)

onMounted(async () => {
  try { stats.value = (await api.getStats()).data } catch (e) { console.error(e) }
})

const statItems = computed(() => stats.value ? [
  { label: '已发布文章', value: stats.value.total_articles },
  { label: '总访问量', value: stats.value.total_views },
  { label: '评论数', value: stats.value.total_comments },
  { label: '分类数', value: stats.value.total_categories },
  { label: '标签数', value: stats.value.total_tags },
  { label: '系列数', value: stats.value.total_series },
] : [])
</script>

<style scoped>
.page-heading { font-size: 22px; font-weight: 500; color: #1f1f1f; margin-bottom: 24px; }
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 16px;
  margin-bottom: 40px;
}
.stat-card {
  background: var(--white);
  border: 1px solid var(--line);
  border-radius: var(--r-md);
  padding: 20px;
  text-align: center;
}
.stat-value { font-size: 28px; font-weight: 600; color: #1f1f1f; }
.stat-label { font-size: 12px; color: var(--muted); margin-top: 4px; }
.section-sub { font-size: 16px; font-weight: 500; color: #1f1f1f; margin-bottom: 16px; }
.action-row { display: flex; gap: 12px; flex-wrap: wrap; }
</style>
