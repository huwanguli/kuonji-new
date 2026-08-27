<template>
  <div>
    <div class="page-header">
      <h1 class="page-heading">文章管理</h1>
      <RouterLink to="/admin/articles/new" class="btn btn-primary">+ 写新文章</RouterLink>
    </div>
    <div class="table-wrap">
      <table class="admin-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>标题</th>
            <th>分类</th>
            <th>状态</th>
            <th>发布时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="a in articles" :key="a.id">
            <td class="td-id">{{ a.id }}</td>
            <td class="td-title">{{ a.title }}</td>
            <td>{{ a.category?.name ?? '-' }}</td>
            <td>
              <span class="status-badge" :class="statusClass(a.status)">{{ statusText(a.status) }}</span>
            </td>
            <td class="td-date">{{ a.published_at ? formatDate(a.published_at) : '-' }}</td>
            <td class="td-actions">
              <RouterLink :to="`/admin/articles/${a.id}/edit`" class="action-link">编辑</RouterLink>
              <button class="action-link action-publish" v-if="a.status === 0" @click="publish(a.id)">发布</button>
              <button class="action-link action-delete" @click="remove(a.id)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <p v-if="!articles.length" class="empty-state">暂无文章</p>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type Article } from '@/utils/api'
import { formatDate } from '@/utils/format'

const articles = ref<Article[]>([])

async function load() {
  try {
    const res = await api.getArticles({ page: 1, page_size: 100 })
    articles.value = res.data.articles
  } catch (e) { console.error(e) }
}

function statusText(s: number) { return ['草稿', '已发布', '定时发布'][s] ?? '未知' }
function statusClass(s: number) { return ['draft', 'published', 'scheduled'][s] ?? '' }

async function publish(id: number) {
  try {
    await api.adminUpdateArticleStatus(id, 1, new Date().toISOString())
    load()
  } catch (e) { console.error(e) }
}

async function remove(id: number) {
  if (!confirm('确定删除这篇文章？')) return
  try {
    await api.adminDeleteArticle(id)
    load()
  } catch (e) { console.error(e) }
}

onMounted(load)
</script>

<style scoped>
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px; }
.page-heading { font-size: 22px; font-weight: 500; color: #1f1f1f; }
.table-wrap { overflow-x: auto; }
.admin-table {
  width: 100%;
  border-collapse: collapse;
  background: var(--white);
  border: 1px solid var(--line);
  border-radius: var(--r-md);
  font-size: 13px;
}
.admin-table th {
  background: #f8f8f6;
  padding: 10px 14px;
  text-align: left;
  font-weight: 500;
  color: var(--muted);
  border-bottom: 1px solid var(--line);
}
.admin-table td {
  padding: 12px 14px;
  border-bottom: 1px solid var(--line);
  color: var(--ink);
}
.td-id { color: var(--muted-light); width: 50px; }
.td-title { max-width: 300px; }
.td-date { font-size: 12px; color: var(--muted-light); white-space: nowrap; }
.td-actions { white-space: nowrap; }

.status-badge {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
}
.status-badge.draft { background: #f5f0e0; color: #9a8530; }
.status-badge.published { background: #e4f4f1; color: #0e8f80; }
.status-badge.scheduled { background: #efebfb; color: #6b4de0; }

.action-link {
  font-size: 12px;
  color: var(--cat-purple);
  background: none;
  border: none;
  cursor: pointer;
  font-family: inherit;
  padding: 2px 6px;
  margin-right: 8px;
}
.action-link:hover { text-decoration: underline; }
.action-delete { color: #c53f2c; }
.action-publish { color: #0e8f80; }
</style>
