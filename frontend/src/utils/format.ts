// 格式化工具

export function slugify(text: string): string {
  const slug = text
    .toLowerCase()
    .trim()
    .replace(/[^\w\u4e00-\u9fa5]+/g, '-')
    .replace(/^-+|-+$/g, '')
  return slug || 'post'
}

const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']

export function formatDate(d: string | null, withTime = false): string {
  if (!d) return ''
  const date = new Date(d)
  if (Number.isNaN(date.getTime())) return ''
  const y = date.getFullYear()
  const m = monthNames[date.getMonth()]
  const day = date.getDate()
  if (!withTime) return `${m} ${day}, ${y}`
  const h = String(date.getHours()).padStart(2, '0')
  const min = String(date.getMinutes()).padStart(2, '0')
  return `${m} ${day}, ${y} ${h}:${min}`
}

export function readTime(content: string): number {
  const chars = (content || '').replace(/\s+/g, '').length
  return Math.max(1, Math.round(chars / 400))
}

/** 返回字数和预估阅读时间，如 "1.2k 字 · 约 3 分钟" */
export function readInfo(content: string): string {
  const chars = (content || '').replace(/\s+/g, '').length
  const minutes = Math.max(1, Math.round(chars / 400))
  const charStr = chars >= 1000 ? `${(chars / 1000).toFixed(1)}k` : String(chars)
  return `${charStr} 字 · 约 ${minutes} 分钟`
}

export function treeCount(list: any[] = []): number {
  let count = list.length
  for (const item of list) {
    if (item.children?.length) {
      count += treeCount(item.children)
    }
  }
  return count
}

export function relativeTime(dateStr: string): string {
  const d = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - d.getTime()
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
  if (diffDays === 0) return 'today'
  if (diffDays === 1) return 'yesterday'
  if (diffDays < 7) return `${diffDays} days ago`
  if (diffDays < 30) return `${Math.floor(diffDays / 7)} weeks ago`
  return formatDate(dateStr)
}
