<template>
  <div>
    <div class="page-header">
      <h1 class="page-heading">{{ isEdit ? '编辑文章' : '写新文章' }}</h1>
      <div class="header-actions">
        <button class="btn btn-ghost" @click="saveDraft">保存草稿</button>
        <button class="btn btn-primary" @click="publish">发布</button>
      </div>
    </div>

    <div class="editor-layout">
      <!-- 左：表单 + 编辑器 -->
      <div class="editor-main">
        <div class="form-group">
          <label class="form-label">标题</label>
          <input v-model="form.title" class="input input-lg" placeholder="文章标题" />
        </div>
        <div class="form-group">
          <label class="form-label">摘要</label>
          <textarea v-model="form.summary" class="input" rows="2" placeholder="一句话描述文章内容"></textarea>
        </div>

        <!-- Markdown 编辑器 + 预览 双栏 -->
        <div class="form-group">
          <div class="editor-toolbar">
            <label class="form-label" style="margin-bottom:0">正文（Markdown）</label>
            <div class="toolbar-actions">
              <div class="toolbar-group">
                <button class="toolbar-btn" @click="insertHeading(1)" title="一级标题">H1</button>
                <button class="toolbar-btn" @click="insertHeading(2)" title="二级标题">H2</button>
                <button class="toolbar-btn" @click="insertHeading(3)" title="三级标题">H3</button>
              </div>
              <span class="toolbar-sep"></span>
              <div class="toolbar-group">
                <button class="toolbar-btn" @click="insertBold" title="粗体 (Ctrl+B)"><b>B</b></button>
                <button class="toolbar-btn" @click="insertItalic" title="斜体 (Ctrl+I)"><i>I</i></button>
                <button class="toolbar-btn" @click="insertStrikethrough" title="删除线"><s>S</s></button>
              </div>
              <span class="toolbar-sep"></span>
              <div class="toolbar-group">
                <button class="toolbar-btn" @click="insertCode" title="行内代码">&lt;/&gt;</button>
                <button class="toolbar-btn" @click="insertCodeBlock" title="代码块">▤</button>
              </div>
              <span class="toolbar-sep"></span>
              <div class="toolbar-group">
                <button class="toolbar-btn" @click="insertLink" title="链接">🔗</button>
                <button class="toolbar-btn" @click="insertQuote" title="引用">❝</button>
                <button class="toolbar-btn" @click="insertUnorderedList" title="无序列表">• —</button>
                <button class="toolbar-btn" @click="insertOrderedList" title="有序列表">1.</button>
                <button class="toolbar-btn" @click="insertHr" title="分割线">——</button>
              </div>
              <span class="toolbar-sep"></span>
              <div class="toolbar-group">
                <label class="toolbar-btn upload-btn" title="上传图片">
                  🖼️
                  <input type="file" accept="image/*" style="display:none" @change="handleUpload" />
                </label>
                <div class="img-size-group" v-if="false"><!-- reserved --></div>
              </div>
              <span class="toolbar-sep"></span>
              <button class="toolbar-btn" :class="{ active: showPreview }" @click="showPreview = !showPreview" title="切换预览">👁</button>
            </div>
          </div>
          <div class="editor-split" :class="{ 'preview-only': false }">
            <textarea
              ref="editorRef"
              v-model="form.content"
              class="input editor-textarea"
              :class="{ 'half': showPreview }"
              rows="25"
              placeholder="用 Markdown 写作…"
              @scroll="syncScroll"
            ></textarea>
            <div
              v-if="showPreview"
              ref="previewRef"
              class="editor-preview markdown-body"
              @scroll="syncScrollReverse"
              v-html="previewHtml"
            ></div>
          </div>
          <!-- 图片大小选择（上传后弹出） -->
          <div v-if="showImgSizePicker" class="img-size-picker">
            <div class="img-size-inner">
              <span class="img-size-label">📐 选择图片插入大小：</span>
              <div class="img-size-btns">
                <button v-for="opt in imgSizeOptions" :key="opt.label" class="img-size-btn" @click="insertImageWithSize(opt.width)">
                  {{ opt.label }}{{ opt.width > 0 ? ` · ${opt.width}px` : '' }}
                </button>
                <button class="img-size-btn img-size-cancel" @click="showImgSizePicker = false; pendingImageUrl = ''">取消</button>
              </div>
            </div>
          </div>
          <p v-if="uploading" class="upload-status">⏳ 上传中…</p>
        </div>
      </div>

      <!-- 右：设置面板 -->
      <div class="editor-settings">
        <div class="form-group">
          <label class="form-label">分类</label>
          <select v-model="form.category_id" class="input">
            <option :value="0" disabled>选择分类</option>
            <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">系列（可选）</label>
          <select v-model="form.series_id" class="input">
            <option :value="null">无</option>
            <option v-for="s in series" :key="s.id" :value="s.id">{{ s.name }}</option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">标签</label>
          <div class="tag-picker">
            <label v-for="t in allTags" :key="t.id" class="tag-option">
              <input type="checkbox" :value="t.id" v-model="form.tag_ids" />
              <span>{{ t.name }}</span>
            </label>
          </div>
        </div>
        <div class="form-group">
          <label class="form-label">封面图</label>
          <div class="cover-upload">
            <div v-if="form.cover_image" class="cover-preview">
              <img :src="form.cover_image" alt="封面" />
              <button class="cover-remove" @click="form.cover_image = ''">×</button>
            </div>
            <label v-else class="cover-upload-btn">
              <input type="file" accept="image/*" style="display:none" @change="handleCoverUpload" />
              <span>+ 上传封面图</span>
            </label>
            <input v-model="form.cover_image" class="input" placeholder="或输入图片 URL" style="margin-top:8px" />
          </div>
        </div>
      </div>
    </div>

    <p v-if="notice" class="editor-notice" :class="{ error: noticeError }">{{ notice }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api, type Category, type Tag, type Series } from '@/utils/api'
import { renderMarkdown } from '@/utils/markdown'

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => !!route.params.id)

const categories = ref<Category[]>([])
const allTags = ref<Tag[]>([])
const series = ref<Series[]>([])
const notice = ref('')
const noticeError = ref(false)
const showPreview = ref(true)
const uploading = ref(false)
const editorRef = ref<HTMLTextAreaElement | null>(null)
const previewRef = ref<HTMLDivElement | null>(null)
let scrollSyncing = false

const form = reactive({
  title: '',
  summary: '',
  content: '',
  category_id: 0 as number,
  series_id: null as number | null,
  tag_ids: [] as number[],
  cover_image: ''
})

const previewHtml = computed(() => renderMarkdown(form.content || ''))

// ── 滚动同步 ──
function syncScroll() {
  if (scrollSyncing || !editorRef.value || !previewRef.value) return
  scrollSyncing = true
  const e = editorRef.value
  const p = previewRef.value
  const ratio = e.scrollTop / (e.scrollHeight - e.clientHeight || 1)
  p.scrollTop = ratio * (p.scrollHeight - p.clientHeight)
  requestAnimationFrame(() => { scrollSyncing = false })
}
function syncScrollReverse() {
  if (scrollSyncing || !editorRef.value || !previewRef.value) return
  scrollSyncing = true
  const e = editorRef.value
  const p = previewRef.value
  const ratio = p.scrollTop / (p.scrollHeight - p.clientHeight || 1)
  e.scrollTop = ratio * (e.scrollHeight - e.clientHeight)
  requestAnimationFrame(() => { scrollSyncing = false })
}

// ── 工具栏快捷插入 ──
function insertAtCursor(before: string, after: string = '') {
  const el = editorRef.value
  if (!el) return
  const start = el.selectionStart
  const end = el.selectionEnd
  const selected = form.content.substring(start, end)
  const replacement = before + selected + after
  form.content = form.content.substring(0, start) + replacement + form.content.substring(end)
  setTimeout(() => {
    el.focus()
    el.selectionStart = start + before.length
    el.selectionEnd = start + before.length + selected.length
  }, 0)
}

// 在当前行前插入前缀（标题、引用、列表等）
function insertLinePrefix(prefix: string) {
  const el = editorRef.value
  if (!el) return
  const start = el.selectionStart
  const content = form.content
  // 找到当前行的起始位置
  const lineStart = content.lastIndexOf('\n', start - 1) + 1
  form.content = content.substring(0, lineStart) + prefix + content.substring(lineStart)
  setTimeout(() => {
    el.focus()
    el.selectionStart = el.selectionEnd = start + prefix.length
  }, 0)
}

function insertHeading(level: number) {
  const prefix = '#'.repeat(level) + ' '
  insertLinePrefix(prefix)
}
function insertBold() { insertAtCursor('**', '**') }
function insertItalic() { insertAtCursor('*', '*') }
function insertStrikethrough() { insertAtCursor('~~', '~~') }
function insertCode() { insertAtCursor('`', '`') }
function insertCodeBlock() { insertAtCursor('\n```\n', '\n```\n') }
function insertLink() { insertAtCursor('[', '](url)') }
function insertQuote() { insertLinePrefix('> ') }
function insertUnorderedList() { insertLinePrefix('- ') }
function insertOrderedList() { insertLinePrefix('1. ') }
function insertHr() { insertAtCursor('\n---\n') }

// 图片大小选项
const imgSizeOptions = [
  { label: '小', width: 300 },
  { label: '中', width: 600 },
  { label: '大', width: 900 },
  { label: '原始', width: 0 },
]
const pendingImageUrl = ref('')
const showImgSizePicker = ref(false)

function handleUpload(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  uploading.value = true
  api.adminUploadImage(file).then(res => {
    pendingImageUrl.value = res.data.url
    showImgSizePicker.value = true
  }).catch((err: any) => {
    alert(err?.data?.message || '上传失败')
  }).finally(() => {
    uploading.value = false
    ;(e.target as HTMLInputElement).value = ''
  })
}

function insertImageWithSize(width: number) {
  const url = pendingImageUrl.value
  if (width > 0) {
    insertAtCursor(`<img src="${url}" width="${width}" />\n`)
  } else {
    insertAtCursor(`![](${url})\n`)
  }
  showImgSizePicker.value = false
  pendingImageUrl.value = ''
}

async function handleCoverUpload(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  uploading.value = true
  try {
    const res = await api.adminUploadImage(file)
    form.cover_image = res.data.url
  } catch (err: any) {
    alert(err?.data?.message || '上传失败')
  } finally {
    uploading.value = false
    ;(e.target as HTMLInputElement).value = ''
  }
}

// ── 数据加载 ──
onMounted(async () => {
  try {
    const [catRes, tagRes, serRes] = await Promise.all([
      api.getCategories(), api.getTags(), api.getSeries()
    ])
    const flat: Category[] = []
    function walk(list: Category[]) { for (const c of list) { flat.push(c); if (c.children?.length) walk(c.children) } }
    walk(catRes.data)
    categories.value = flat
    allTags.value = tagRes.data
    series.value = serRes.data
  } catch (e) { console.error(e) }

  if (isEdit.value) {
    try {
      const res = await api.getArticle(route.params.id as string)
      const a = res.data
      form.title = a.title
      form.summary = a.summary
      form.content = a.content ?? ''
      form.category_id = a.category?.id ?? 0
      form.series_id = a.series?.id ?? null
      form.tag_ids = a.tags?.map((t: Tag) => t.id) ?? []
      form.cover_image = a.cover_image
    } catch (e) { console.error(e) }
  }
})

// ── 保存 ──
async function saveDraft() { await doSave(0) }
async function publish() { await doSave(1) }

async function doSave(status: number) {
  notice.value = ''
  noticeError.value = false
  if (!form.title.trim()) { notice.value = '请输入标题'; noticeError.value = true; return }
  if (!form.content.trim()) { notice.value = '请输入正文'; noticeError.value = true; return }
  if (!form.category_id) { notice.value = '请选择分类'; noticeError.value = true; return }

  const body = {
    title: form.title, summary: form.summary, content: form.content,
    category_id: form.category_id, series_id: form.series_id,
    tag_ids: form.tag_ids, cover_image: form.cover_image, status
  }

  try {
    if (isEdit.value) {
      await api.adminUpdateArticle(Number(route.params.id), body)
    } else {
      await api.adminCreateArticle(body)
    }
    notice.value = status === 1 ? '文章已发布 ✓' : '草稿已保存 ✓'
    setTimeout(() => router.push('/admin/articles'), 800)
  } catch (e: any) {
    notice.value = e?.data?.message || '保存失败'
    noticeError.value = true
  }
}
</script>

<style scoped>
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px; }
.page-heading { font-size: 22px; font-weight: 500; color: #1f1f1f; }
.header-actions { display: flex; gap: 8px; }

.editor-layout {
  display: grid;
  grid-template-columns: 1fr 280px;
  gap: 24px;
  align-items: start;
}

.form-group { margin-bottom: 16px; }
.form-label { display: block; font-size: 12px; font-weight: 500; color: var(--muted); margin-bottom: 6px; letter-spacing: 0.5px; }
.input-lg { font-size: 18px; font-weight: 500; padding: 12px 14px; }

/* 工具栏 */
.editor-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}
.toolbar-actions { display: flex; gap: 2px; align-items: center; flex-wrap: wrap; }
.toolbar-group { display: flex; gap: 2px; }
.toolbar-sep { width: 1px; height: 20px; background: var(--line); margin: 0 4px; }
.toolbar-btn {
  background: var(--white);
  border: 1px solid var(--line);
  border-radius: var(--r-xs);
  padding: 4px 8px;
  font-size: 12px;
  cursor: pointer;
  color: var(--muted);
  transition: all 0.1s;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 28px;
  height: 28px;
}
.toolbar-btn:hover { border-color: var(--cat-purple); color: var(--ink); background: rgba(208,192,232,0.1); }
.toolbar-btn.active { background: var(--cat-purple); color: #fff; border-color: var(--cat-purple); }
.upload-btn { cursor: pointer; }

/* 图片大小选择 */
.img-size-picker {
  margin-top: 12px;
  padding: 14px 18px;
  background: linear-gradient(135deg, rgba(208,192,232,0.15), rgba(184,212,227,0.15));
  border: 2px solid var(--cat-purple);
  border-radius: var(--r-md);
  animation: fade-in 0.2s ease;
}
.img-size-inner { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
.img-size-label { font-size: 13px; font-weight: 500; color: var(--ink); white-space: nowrap; }
.img-size-btns { display: flex; gap: 6px; flex-wrap: wrap; }
.img-size-btn {
  padding: 6px 16px;
  border: 1px solid var(--cat-purple);
  border-radius: var(--r-sm);
  background: var(--white);
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  color: var(--ink);
  transition: all 0.15s;
}
.img-size-btn:hover { background: var(--cat-purple); color: #fff; }
.img-size-cancel { color: var(--muted); border-color: var(--line); font-weight: 400; }
.img-size-cancel:hover { background: var(--line); color: var(--ink); }

/* 编辑器分栏 */
.editor-split {
  display: grid;
  grid-template-columns: 1fr;
  gap: 0;
}
.editor-split:has(.editor-preview) {
  grid-template-columns: 1fr 1fr;
}
.editor-textarea {
  resize: vertical;
  min-height: 500px;
  font-family: var(--font-mono);
  font-size: 13px;
  line-height: 1.7;
  border-radius: var(--r-sm) 0 0 var(--r-sm);
}
.editor-textarea.half { border-radius: var(--r-sm) 0 0 var(--r-sm); }
.editor-preview {
  border: 1px solid var(--line);
  border-left: none;
  border-radius: 0 var(--r-sm) var(--r-sm) 0;
  padding: 16px 20px;
  overflow-y: auto;
  min-height: 500px;
  background: var(--white);
  font-size: 14px;
  line-height: 1.75;
}

.upload-status { font-size: 12px; color: var(--muted); margin-top: 4px; }

/* 设置面板 */
.editor-settings {
  background: var(--white);
  border: 1px solid var(--line);
  border-radius: var(--r-md);
  padding: 20px;
  position: sticky;
  top: 20px;
}

.tag-picker { display: flex; flex-wrap: wrap; gap: 6px; }
.tag-option {
  display: flex; align-items: center; gap: 4px;
  font-size: 12px; color: var(--ink); cursor: pointer;
  padding: 3px 8px; border: 1px solid var(--line); border-radius: var(--r-xs);
  transition: background 0.1s;
}
.tag-option:has(input:checked) { background: var(--cat-purple); color: #fff; border-color: var(--cat-purple); }
.tag-option input { display: none; }

/* 封面图上传 */
.cover-upload { }
.cover-preview {
  position: relative;
  border-radius: var(--r-sm);
  overflow: hidden;
  border: 1px solid var(--line);
  margin-bottom: 8px;
}
.cover-preview img { width: 100%; height: 120px; object-fit: cover; display: block; }
.cover-remove {
  position: absolute; top: 4px; right: 4px;
  background: rgba(0,0,0,0.6); color: #fff; border: none;
  border-radius: 50%; width: 24px; height: 24px;
  cursor: pointer; font-size: 14px; display: flex;
  align-items: center; justify-content: center;
}
.cover-upload-btn {
  display: flex; align-items: center; justify-content: center;
  height: 80px; border: 2px dashed var(--line); border-radius: var(--r-sm);
  cursor: pointer; color: var(--muted); font-size: 13px;
  transition: border-color 0.15s;
}
.cover-upload-btn:hover { border-color: var(--cat-purple); color: var(--ink); }

.editor-notice { margin-top: 16px; font-size: 13px; color: #0e8f80; }
.editor-notice.error { color: #c53f2c; }

select.input { cursor: pointer; }

@media (max-width: 1100px) {
  .editor-layout { grid-template-columns: 1fr; }
  .editor-settings { position: static; }
}
@media (max-width: 768px) {
  .editor-split:has(.editor-preview) { grid-template-columns: 1fr; }
  .editor-textarea { border-radius: var(--r-sm); }
  .editor-preview { border-left: 1px solid var(--line); border-top: none; border-radius: 0 0 var(--r-sm) var(--r-sm); }
}
</style>
