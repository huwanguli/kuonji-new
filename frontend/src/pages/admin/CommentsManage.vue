<template>
  <div>
    <div class="page-header">
      <h1 class="page-heading">评论管理</h1>
    </div>
    <div class="comment-list">
      <div v-for="c in comments" :key="c.id" class="comment-card">
        <div class="comment-meta">
          <span class="comment-nick">{{ c.nickname }}</span>
          <span class="comment-email">{{ c.email }}</span>
          <span class="status-badge" :class="commentStatusClass(c.status)">{{ commentStatusText(c.status) }}</span>
        </div>
        <p class="comment-text">{{ c.content }}</p>
        <div class="comment-actions">
          <button v-if="c.status !== 1" class="action-link action-approve" @click="approve(c.id)">通过</button>
          <button v-if="c.status !== 2" class="action-link action-spam" @click="spam(c.id)">标记垃圾</button>
          <button class="action-link action-delete" @click="remove(c.id)">删除</button>
        </div>
      </div>
    </div>
    <p v-if="!comments.length" class="empty-state">暂无评论</p>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type Comment } from '@/utils/api'

const comments = ref<Comment[]>([])

async function load() {
  try {
    // 拉取所有文章的评论（用 search 拿到所有文章，再逐个拉评论 —— 简化方案：直接获取所有文章列表）
    const articlesRes = await api.getArticles({ page: 1, page_size: 100 })
    const allComments: Comment[] = []
    for (const a of articlesRes.data.articles) {
      try {
        const cRes = await api.getComments(a.id)
        allComments.push(...flattenComments(cRes.data))
      } catch {}
    }
    comments.value = allComments.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
  } catch (e) { console.error(e) }
}

function flattenComments(list: Comment[]): Comment[] {
  const result: Comment[] = []
  for (const c of list) {
    result.push(c)
    if (c.children?.length) result.push(...flattenComments(c.children))
  }
  return result
}

function commentStatusText(s: number) { return ['待审核', '已通过', '垃圾'][s] ?? '未知' }
function commentStatusClass(s: number) { return ['pending', 'approved', 'spam'][s] ?? '' }

async function approve(id: number) {
  try { await api.adminUpdateCommentStatus(id, 1); load() } catch (e) { console.error(e) }
}
async function spam(id: number) {
  try { await api.adminUpdateCommentStatus(id, 2); load() } catch (e) { console.error(e) }
}
async function remove(id: number) {
  if (!confirm('确定删除此评论？')) return
  try { await api.adminDeleteComment(id); load() } catch (e: any) { alert(e?.data?.message || '删除失败') }
}

onMounted(load)
</script>

<style scoped>
.page-header { margin-bottom: 24px; }
.page-heading { font-size: 22px; font-weight: 500; color: #1f1f1f; }
.comment-list { display: flex; flex-direction: column; gap: 12px; }
.comment-card {
  background: var(--white);
  border: 1px solid var(--line);
  border-radius: var(--r-md);
  padding: 16px;
}
.comment-meta { display: flex; align-items: center; gap: 10px; margin-bottom: 8px; }
.comment-nick { font-size: 13px; font-weight: 500; color: var(--ink); }
.comment-email { font-size: 12px; color: var(--muted-light); }
.comment-text { font-size: 13px; color: var(--ink-light); line-height: 1.6; margin-bottom: 10px; }
.comment-actions { display: flex; gap: 8px; }

.status-badge { font-size: 11px; padding: 2px 8px; border-radius: 10px; }
.status-badge.pending { background: #f5f0e0; color: #9a8530; }
.status-badge.approved { background: #e4f4f1; color: #0e8f80; }
.status-badge.spam { background: #fde8e8; color: #c53f2c; }

.action-link { font-size: 12px; background: none; border: none; cursor: pointer; font-family: inherit; padding: 2px 6px; }
.action-link:hover { text-decoration: underline; }
.action-approve { color: #0e8f80; }
.action-spam { color: #9a8530; }
.action-delete { color: #c53f2c; }
</style>
