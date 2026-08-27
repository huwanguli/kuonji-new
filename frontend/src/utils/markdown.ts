// Markdown 渲染工具（模块顶层执行，保证 marked.use 只注册一次）
import { marked } from 'marked'
import { markedHighlight } from 'marked-highlight'
import hljs from 'highlight.js'

function slugify(text: string): string {
  return text
    .toLowerCase()
    .trim()
    .replace(/[^\w\u4e00-\u9fa5]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

// 自定义 renderer：给 h2/h3 添加 id
const renderer = new marked.Renderer()
renderer.heading = function ({ text, depth }: { text: string; depth: number; raw: string }) {
  const plainText = text.replace(/<[^>]*>/g, '')
  const id = slugify(plainText)
  return `<h${depth} id="${id}">${text}</h${depth}>`
}

marked.use(
  markedHighlight({
    langPrefix: 'hljs language-',
    highlight(code: string, lang: string) {
      const language = hljs.getLanguage(lang) ? lang : 'plaintext'
      return hljs.highlight(code, { language }).value
    }
  })
)

marked.use({ renderer })

export function renderMarkdown(content: string): string {
  return marked.parse(content) as string
}
