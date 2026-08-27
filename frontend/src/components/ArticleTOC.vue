<template>
  <nav v-if="headings.length" class="toc" aria-label="文章目录">
    <h4 class="toc-title">ON THIS PAGE</h4>
    <ul class="toc-list">
      <li
        v-for="h in headings"
        :key="h.id"
        class="toc-item"
        :class="{ 'toc-h3': h.level === 3 }"
      >
        <a
          :href="'#' + h.id"
          class="toc-link"
          :class="{ active: activeId === h.id }"
          @click.prevent="scrollTo(h.id)"
        >
          {{ h.text }}
        </a>
      </li>
    </ul>
  </nav>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue'

interface Heading {
  id: string
  text: string
  level: number
}

const props = defineProps<{ content: string }>()

const headings = ref<Heading[]>([])
const activeId = ref('')

function slugify(text: string): string {
  return text
    .toLowerCase()
    .trim()
    .replace(/[^\w\u4e00-\u9fa5]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

function extractHeadings() {
  if (!props.content) { headings.value = []; return }
  const regex = /^(#{2,3})\s+(.+)$/gm
  const result: Heading[] = []
  let match
  while ((match = regex.exec(props.content)) !== null) {
    const level = match[1].length
    const text = match[2].replace(/[*_`]/g, '').trim()
    const id = slugify(text)
    result.push({ id, text, level })
  }
  headings.value = result
}

function scrollTo(id: string) {
  const el = document.getElementById(id)
  if (el) {
    const offset = 80
    const y = el.getBoundingClientRect().top + window.scrollY - offset
    window.scrollTo({ top: y, behavior: 'smooth' })
  }
}

// 滚动监听高亮当前标题
let observer: IntersectionObserver | null = null

function setupObserver() {
  // 先断开旧的
  observer?.disconnect()

  observer = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (entry.isIntersecting) {
          activeId.value = entry.target.id
        }
      }
    },
    { rootMargin: '-80px 0px -70% 0px', threshold: 0 }
  )

  nextTick(() => {
    for (const h of headings.value) {
      const el = document.getElementById(h.id)
      if (el) observer!.observe(el)
    }
  })
}

function rebuild() {
  activeId.value = ''
  extractHeadings()
  setupObserver()
}

// 监听 content 变化（前后章切换时触发）
watch(() => props.content, () => {
  rebuild()
})

onMounted(() => {
  rebuild()
})

onUnmounted(() => {
  observer?.disconnect()
})
</script>

<style scoped>
.toc {
  position: sticky;
  top: 80px;
  max-height: calc(100vh - 120px);
  overflow-y: auto;
  padding-left: 20px;
  border-left: 1px solid var(--line);
}

.toc::-webkit-scrollbar { width: 3px; }
.toc::-webkit-scrollbar-thumb { background: var(--line); border-radius: 2px; }

.toc-title {
  font-size: 11px;
  color: var(--muted);
  letter-spacing: 2px;
  font-weight: 500;
  text-transform: uppercase;
  margin-bottom: 16px;
}

.toc-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.toc-h3 {
  padding-left: 14px;
}

.toc-link {
  font-size: 13px;
  color: var(--muted);
  line-height: 1.6;
  display: block;
  padding: 3px 0;
  transition: color 0.15s ease;
  text-decoration: none;
}

.toc-link:hover {
  color: var(--ink);
}

.toc-link.active {
  color: var(--ink);
  font-weight: 500;
}
</style>
