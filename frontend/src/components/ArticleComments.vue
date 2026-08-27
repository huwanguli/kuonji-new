<template>
  <section class="comments" aria-label="评论">
    <header class="comments-header">
      <h3 class="section-label">COMMENTS · {{ total }}</h3>
      <div class="comments-sort">Sort by newest ↓</div>
    </header>

    <div v-if="comments.length" class="comment-list">
      <CommentThread v-for="c in comments" :key="c.id" :comment="c" />
    </div>

    <!-- 提交表单 -->
    <div class="comment-form-wrap">
      <div class="avatar avatar-sm"></div>
      <div class="comment-form-box">
        <div class="form-row">
          <input v-model="form.nickname" class="input" placeholder="昵称" required maxlength="50" />
          <input v-model="form.email" class="input" type="email" placeholder="邮箱（仅用于回复通知）" required maxlength="100" />
        </div>
        <textarea
          v-model="form.content"
          class="input comment-input"
          placeholder="Write a thoughtful comment…"
          rows="3"
          required
          maxlength="2000"
        />
        <div class="form-actions">
          <span v-if="notice" class="notice" :class="{ error: noticeError }">{{ notice }}</span>
          <button type="button" class="btn btn-primary" :disabled="submitting" @click="submit">
            {{ submitting ? '提交中…' : 'Post' }}
          </button>
        </div>
      </div>
    </div>

    <p v-if="!comments.length && !form.content" class="empty-state">还没有评论，来抢沙发。</p>
  </section>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { api, type Comment } from '@/utils/api'
import CommentThread from './CommentThread.vue'

const props = defineProps<{ articleId: number }>()

const comments = ref<Comment[]>([])
const total = ref(0)
const submitting = ref(false)
const notice = ref('')
const noticeError = ref(false)

const form = reactive({ nickname: '', email: '', content: '' })

function countComments(list: Comment[]): number {
  let n = list.length
  for (const c of list) if (c.children?.length) n += countComments(c.children)
  return n
}

async function load() {
  try {
    const res = await api.getComments(props.articleId)
    comments.value = res.data
    total.value = countComments(res.data)
  } catch (e) {
    console.error('Failed to load comments:', e)
  }
}

async function submit() {
  if (!form.nickname.trim() || !form.email.trim() || !form.content.trim()) return
  submitting.value = true
  noticeError.value = false
  try {
    await api.createComment(props.articleId, {
      nickname: form.nickname.trim(),
      email: form.email.trim(),
      content: form.content.trim()
    })
    notice.value = '评论已提交，审核通过后展示 ✓'
    form.content = ''
    load()
  } catch (e: any) {
    noticeError.value = true
    notice.value = e?.data?.message || '提交失败，请稍后再试'
  } finally {
    submitting.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.comments {
  padding-top: 24px;
  border-top: 1px solid var(--line);
}

.comments-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}
.comments-sort { font-size: 11px; color: var(--muted-light); }

.comment-list { display: flex; flex-direction: column; margin-bottom: 28px; }

.comment-form-wrap {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  margin-bottom: 12px;
}
.comment-form-box {
  flex: 1;
  border: 1px solid var(--line);
  border-radius: var(--r-md);
  background: var(--white);
  padding: 12px 14px;
}
.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.75rem;
  margin-bottom: 0.75rem;
}
.comment-input {
  resize: vertical;
  min-height: 70px;
  font-family: inherit;
  line-height: 1.6;
  margin-bottom: 0.75rem;
}
.form-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}
.notice { font-size: 12px; color: var(--cat-green); }
.notice.error { color: #c53f2c; }
.comments-empty { padding: 2rem 0; }

@media (max-width: 640px) { .form-row { grid-template-columns: 1fr; } }
</style>
