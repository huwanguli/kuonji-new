<template>
  <div class="article-page" v-if="article">
    <div class="container article-container">
      <!-- 文章正文区 -->
      <article class="article-main">
        <!-- 文章头部 -->
        <header class="article-header">
          <div class="article-category">
            <span class="dot" :class="getDotColor(article.category?.name ?? '')"></span>
            <span class="article-category-label">CATEGORY · {{ article.category?.name?.toUpperCase() ?? 'UNCATEGORIZED' }}</span>
          </div>
          <h1 class="article-title">{{ article.title }}</h1>
          <p v-if="article.summary" class="article-subtitle">{{ article.summary }}</p>
          <div class="article-meta-row">
            <span>{{ formatDate(article.published_at) }}</span>
            <span class="meta-sep">·</span>
            <span>{{ readInfo(article.content ?? '') }}</span>
          </div>
          <div class="article-tags-row">
            <RouterLink v-for="tag in article.tags" :key="tag.id" :to="`/tags/${tag.id}`" class="tag tag-sm">#{{ tag.name }}</RouterLink>
          </div>
        </header>

        <!-- 正文 -->
        <div class="markdown-body" v-html="renderedContent"></div>

        <!-- FILED UNDER -->
        <div v-if="article.tags.length" class="filed-under">
          <h3 class="section-label">FILED UNDER</h3>
          <div class="filed-tags">
            <RouterLink v-for="tag in article.tags" :key="tag.id" :to="`/tags/${tag.id}`" class="tag tag-sm">#{{ tag.name }}</RouterLink>
          </div>
        </div>

        <!-- 系列导航 -->
        <div v-if="seriesNav" class="series-nav">
          <RouterLink :to="`/series/${article.series!.id}`" class="series-nav-label">{{ article.series!.name }}</RouterLink>
          <div class="series-chapters">
            <RouterLink v-if="seriesNav.prev" :to="`/articles/${seriesNav.prev.id}`" class="chapter-nav chapter-prev">
              <span class="chapter-arrow">←</span>
              <div class="chapter-info">
                <span class="chapter-label">PREVIOUS</span>
                <span class="chapter-title">{{ seriesNav.prev.title }}</span>
              </div>
            </RouterLink>
            <RouterLink v-if="seriesNav.next" :to="`/articles/${seriesNav.next.id}`" class="chapter-nav chapter-next">
              <div class="chapter-info chapter-info-right">
                <span class="chapter-label">NEXT →</span>
                <span class="chapter-title">{{ seriesNav.next.title }}</span>
              </div>
              <span class="chapter-arrow">→</span>
            </RouterLink>
          </div>
        </div>

        <!-- 评论 -->
        <ArticleComments :article-id="article.id" />
      </article>

      <!-- 右侧目录 -->
      <aside class="article-toc-aside">
        <ArticleTOC :content="article.content ?? ''" />

        <!-- 系列文章目录 -->
        <nav v-if="seriesArticlesList.length" class="series-toc">
          <h4 class="toc-title">{{ article.series?.name ?? 'SERIES' }}</h4>
          <ol class="series-toc-list">
            <li
              v-for="(sa, idx) in seriesArticlesList"
              :key="sa.id"
              class="series-toc-item"
              :class="{ 'series-toc-current': sa.id === article.id }"
            >
              <RouterLink
                :to="`/articles/${sa.id}`"
                class="series-toc-link"
              >
                <span class="series-toc-num">{{ String(idx + 1).padStart(2, '0') }}</span>
                <span class="series-toc-text">{{ sa.title }}</span>
              </RouterLink>
            </li>
          </ol>
        </nav>
      </aside>
    </div>
  </div>

  <p v-else-if="!loading" class="empty-state container">文章不存在或已被删除。</p>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { api, type Article as ArticleType } from '@/utils/api'
import { formatDate, readInfo } from '@/utils/format'
import { renderMarkdown } from '@/utils/markdown'
import { getDotColor } from '@/utils/colors'
import ArticleComments from '@/components/ArticleComments.vue'
import ArticleTOC from '@/components/ArticleTOC.vue'

const route = useRoute()

const article = ref<ArticleType | null>(null)
const loading = ref(true)
const seriesArticles = ref<{ prev: any | null; next: any | null } | null>(null)
const seriesArticlesList = ref<any[]>([])

const renderedContent = computed(() => {
  if (!article.value?.content) return ''
  return renderMarkdown(article.value.content)
})

const seriesNav = computed(() => seriesArticles.value)

async function load() {
  loading.value = true
  seriesArticlesList.value = []
  seriesArticles.value = null
  try {
    const id = route.params.id as string
    const res = await api.getArticle(id)
    article.value = res.data

    if (article.value?.series) {
      try {
        const seriesRes = await api.getSeriesDetail(article.value.series.id)
        const articles = seriesRes.data.articles || []
        seriesArticlesList.value = articles
        const idx = articles.findIndex((a: any) => a.id === Number(id))
        seriesArticles.value = {
          prev: idx > 0 ? articles[idx - 1] : null,
          next: idx < articles.length - 1 ? articles[idx + 1] : null
        }
      } catch {
        seriesArticles.value = null
        seriesArticlesList.value = []
      }
    }
  } catch (e) {
    console.error('Failed to load article:', e)
    article.value = null
  } finally {
    loading.value = false
  }
}

watch(() => route.params.id, load)
onMounted(load)
</script>

<style scoped>
/* 文章页独立布局：正文 + 右侧 TOC */
.article-page {
  padding: 40px 0 2rem;
}

.article-container {
  display: grid;
  grid-template-columns: 1fr 200px;
  gap: 48px;
  align-items: start;
}

.article-main {
  min-width: 0;
  max-width: 840px;
  justify-self: center;
}

/* 右侧目录 */
.article-toc-aside {
  position: sticky;
  top: 80px;
  display: flex;
  flex-direction: column;
  gap: 32px;
}

/* 系列文章目录 */
.series-toc { }
.series-toc .toc-title {
  font-size: 11px;
  color: var(--muted);
  letter-spacing: 2px;
  font-weight: 500;
  text-transform: uppercase;
  margin-bottom: 12px;
}
.series-toc-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.series-toc-item { }
.series-toc-current .series-toc-link {
  color: var(--ink);
  font-weight: 500;
  background: rgba(208,192,232,0.12);
}
.series-toc-link {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 5px 8px;
  border-radius: var(--r-xs);
  font-size: 12px;
  color: var(--muted);
  line-height: 1.4;
  transition: background 0.1s, color 0.1s;
}
.series-toc-link:hover {
  background: rgba(208,192,232,0.08);
  color: var(--ink);
}
.series-toc-num {
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--muted-xlight);
  flex-shrink: 0;
  min-width: 1.4em;
  margin-top: 1px;
}
.series-toc-text {
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

/* 文章头部 */
.article-header { padding: 28px 32px; margin-bottom: 20px; background: var(--white); border-radius: var(--r-md); border: 1px solid var(--line); }
.article-category {
  display: flex; align-items: center; gap: 8px;
  font-size: 10px; color: var(--muted-light); letter-spacing: 1.8px; margin-bottom: 18px;
}
.article-category-label { text-transform: uppercase; }
.article-title {
  font-size: 32px; font-weight: 500; color: #1f1f1f;
  line-height: 1.25; letter-spacing: -0.3px; margin-bottom: 14px;
}
.article-subtitle { font-size: 15px; color: var(--muted); margin-bottom: 16px; line-height: 1.55; }
.article-meta-row { display: flex; align-items: center; gap: 8px; font-size: 13px; color: var(--muted-light); margin-bottom: 16px; }
.meta-sep { color: var(--line); }
.article-tags-row { display: flex; gap: 6px; flex-wrap: wrap; }

.filed-under {
  padding: 36px 0 0; margin-top: 32px; border-top: 1px solid var(--line);
}
.filed-under .section-label { margin-bottom: 14px; }
.filed-tags { display: flex; gap: 6px; flex-wrap: wrap; margin-bottom: 32px; }

.series-nav { margin-top: 8px; padding-top: 32px; border-top: 1px solid var(--line); }
.series-nav-label {
  display: block; font-size: 10px; color: var(--muted);
  letter-spacing: 2px; font-weight: 500; margin-bottom: 16px; text-transform: uppercase;
}
.series-chapters { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.chapter-nav {
  border: 1px solid var(--line); border-radius: var(--r-md); padding: 14px 16px;
  display: flex; align-items: center; gap: 14px;
  background: var(--bg-card); transition: border-color 0.15s ease; color: var(--ink);
}
.chapter-nav:hover { border-color: var(--cat-purple); color: var(--ink); }
.chapter-next { justify-content: flex-end; text-align: right; }
.chapter-arrow { font-size: 16px; color: var(--muted-light); flex-shrink: 0; }
.chapter-info { display: flex; flex-direction: column; gap: 3px; min-width: 0; }
.chapter-info-right { text-align: right; }
.chapter-label { font-size: 9px; color: var(--muted-light); letter-spacing: 1.8px; font-weight: 500; }
.chapter-title { font-size: 13px; color: var(--ink); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

/* 响应式 */
@media (max-width: 1100px) {
  .article-container {
    grid-template-columns: 1fr;
  }
  .article-toc-aside {
    display: none;
  }
}

@media (max-width: 640px) {
  .article-title { font-size: 26px; }
  .series-chapters { grid-template-columns: 1fr; }
}
</style>
