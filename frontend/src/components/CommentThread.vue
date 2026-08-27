<template>
  <div class="comment-thread">
    <article class="comment">
      <div class="avatar avatar-sm" :style="{ background: avatarGradient }"></div>
      <div class="comment-body">
        <div class="comment-head">
          <span class="comment-nickname">{{ comment.nickname }}</span>
          <span class="comment-meta-sep">·</span>
          <span class="comment-date">{{ relativeTime(comment.created_at) }}</span>
        </div>
        <p class="comment-content">{{ comment.content }}</p>
      </div>
    </article>
    <div v-if="comment.children?.length" class="comment-replies">
      <CommentThread
        v-for="child in comment.children"
        :key="child.id"
        :comment="child"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Comment } from '@/utils/api'
import { relativeTime } from '@/utils/format'

const props = defineProps<{ comment: Comment }>()

const gradients = [
  'linear-gradient(135deg, #F5D5D5, #D5D0E8)',
  'linear-gradient(135deg, #C8E6D0, #B8D4E3)',
  'linear-gradient(135deg, #D5D0E8, #F5D5D5)',
  'linear-gradient(135deg, #B8D4E3, #D5D0E8)',
]
const avatarGradient = computed(() => {
  let h = 0
  for (let i = 0; i < props.comment.nickname.length; i++) {
    h = ((h << 5) - h + props.comment.nickname.charCodeAt(i)) | 0
  }
  return gradients[Math.abs(h) % gradients.length]
})
</script>

<style scoped>
.comment-thread { display: flex; flex-direction: column; }

.comment {
  display: flex;
  gap: 12px;
  padding: 0 0 22px;
}

.comment-body { flex: 1; min-width: 0; }

.comment-head {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-bottom: 4px;
}
.comment-nickname { font-size: 13px; font-weight: 500; color: var(--ink); }
.comment-meta-sep { color: var(--muted-light); font-size: 12px; }
.comment-date { font-size: 12px; color: var(--muted-light); }

.comment-content {
  font-size: 13px;
  color: var(--muted);
  line-height: 1.65;
  word-break: break-word;
}

.comment-replies {
  margin-left: 44px;
  padding-left: 16px;
  border-left: 2px solid var(--line);
}
</style>
