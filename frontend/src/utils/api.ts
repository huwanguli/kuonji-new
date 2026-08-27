// API 请求封装（纯 fetch，不依赖任何框架）

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface Article {
  id: number
  title: string
  summary: string
  content?: string
  cover_image: string
  status: number
  published_at: string | null
  view_count: number
  category: Category | null
  tags: Tag[]
  series: Series | null
  created_at: string
  updated_at: string
}

export interface Category {
  id: number
  name: string
  description: string
  parent_id: number | null
  sort: number
  children?: Category[]
  article_count?: number
}

export interface Tag {
  id: number
  name: string
}

export interface Series {
  id: number
  name: string
  description: string
  cover_image: string
}

export interface Comment {
  id: number
  article_id: number
  parent_id: number | null
  nickname: string
  email: string
  content: string
  status: number
  created_at: string
  children?: Comment[]
}

export interface Paginated<T> {
  articles: T[]
  total: number
  page: number
  page_size: number
}

import { authHeaders } from './auth'

const BASE_URL = import.meta.env.VITE_API_BASE || '/api/v1'

async function request<T>(url: string, options: RequestInit = {}): Promise<T> {
  const res = await fetch(`${BASE_URL}${url}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers || {})
    }
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ message: res.statusText }))
    throw { status: res.status, data: err }
  }
  return res.json()
}

export const api = {
  // 文章
  getArticles: (params: Record<string, string | number> = {}) => {
    const qs = new URLSearchParams()
    for (const [k, v] of Object.entries(params)) {
      if (v !== undefined && v !== null && v !== '') qs.set(k, String(v))
    }
    const query = qs.toString()
    return request<ApiResponse<Paginated<Article>>>(`/articles${query ? '?' + query : ''}`)
  },

  getArticle: (id: number | string) =>
    request<ApiResponse<Article>>(`/articles/${id}`),

  // 分类
  getCategories: () =>
    request<ApiResponse<Category[]>>('/categories'),

  // 标签
  getTags: () =>
    request<ApiResponse<Tag[]>>('/tags'),

  // 系列
  getSeries: () =>
    request<ApiResponse<Series[]>>('/series'),

  getSeriesDetail: (id: number | string) =>
    request<ApiResponse<Series & { articles: Article[] }>>(`/series/${id}`),

  // 评论
  getComments: (articleId: number | string) =>
    request<ApiResponse<Comment[]>>(`/articles/${articleId}/comments`),

  createComment: (articleId: number | string, body: { nickname: string; email: string; content: string }) =>
    request<ApiResponse<Comment>>(`/articles/${articleId}/comments`, {
      method: 'POST',
      body: JSON.stringify(body)
    }),

  // 搜索
  search: (keyword: string, params: Record<string, string | number> = {}) => {
    const qs = new URLSearchParams({ q: keyword })
    for (const [k, v] of Object.entries(params)) {
      if (v !== undefined && v !== null && v !== '') qs.set(k, String(v))
    }
    return request<ApiResponse<Paginated<Article>>>(`/search?${qs.toString()}`)
  },

  // ── Admin API ──

  // 登录
  login: (username: string, password: string) =>
    request<ApiResponse<{ token: string; user: any }>>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password })
    }),

  // 统计
  getStats: () =>
    request<ApiResponse<{ total_articles: number; total_views: number; total_comments: number; total_categories: number; total_tags: number; total_series: number }>>('/admin/stats', {
      headers: authHeaders()
    }),

  // 文章管理
  adminCreateArticle: (body: Record<string, any>) =>
    request<ApiResponse<Article>>('/admin/articles', {
      method: 'POST',
      body: JSON.stringify(body),
      headers: authHeaders()
    }),

  adminUpdateArticle: (id: number, body: Record<string, any>) =>
    request<ApiResponse<Article>>(`/admin/articles/${id}`, {
      method: 'PUT',
      body: JSON.stringify(body),
      headers: authHeaders()
    }),

  adminDeleteArticle: (id: number) =>
    request<ApiResponse<null>>(`/admin/articles/${id}`, {
      method: 'DELETE',
      headers: authHeaders()
    }),

  adminUpdateArticleStatus: (id: number, status: number, publishedAt?: string) =>
    request<ApiResponse<Article>>(`/admin/articles/${id}/status`, {
      method: 'PUT',
      body: JSON.stringify({ status, published_at: publishedAt }),
      headers: authHeaders()
    }),

  // 分类管理
  adminCreateCategory: (body: { name: string; description?: string; parent_id?: number; sort?: number }) =>
    request<ApiResponse<Category>>('/admin/categories', {
      method: 'POST',
      body: JSON.stringify(body),
      headers: authHeaders()
    }),

  adminUpdateCategory: (id: number, body: Record<string, any>) =>
    request<ApiResponse<Category>>(`/admin/categories/${id}`, {
      method: 'PUT',
      body: JSON.stringify(body),
      headers: authHeaders()
    }),

  adminDeleteCategory: (id: number) =>
    request<ApiResponse<null>>(`/admin/categories/${id}`, {
      method: 'DELETE',
      headers: authHeaders()
    }),

  // 标签管理
  adminCreateTag: (name: string) =>
    request<ApiResponse<Tag>>('/admin/tags', {
      method: 'POST',
      body: JSON.stringify({ name }),
      headers: authHeaders()
    }),

  adminDeleteTag: (id: number) =>
    request<ApiResponse<null>>(`/admin/tags/${id}`, {
      method: 'DELETE',
      headers: authHeaders()
    }),

  // 系列管理
  adminCreateSeries: (body: { name: string; description?: string; sort?: number }) =>
    request<ApiResponse<Series>>('/admin/series', {
      method: 'POST',
      body: JSON.stringify(body),
      headers: authHeaders()
    }),

  adminDeleteSeries: (id: number) =>
    request<ApiResponse<null>>(`/admin/series/${id}`, {
      method: 'DELETE',
      headers: authHeaders()
    }),

  // 评论管理
  adminUpdateCommentStatus: (id: number, status: number) =>
    request<ApiResponse<Comment>>(`/admin/comments/${id}/status`, {
      method: 'PUT',
      body: JSON.stringify({ status }),
      headers: authHeaders()
    }),

  adminDeleteComment: (id: number) =>
    request<ApiResponse<null>>(`/admin/comments/${id}`, {
      method: 'DELETE',
      headers: authHeaders()
    }),

  // 获取所有评论（含待审核）
  adminGetAllComments: () =>
    request<ApiResponse<Comment[]>>('/admin/comments', {
      headers: authHeaders()
    }),

  // 图片上传（multipart/form-data）
  adminUploadImage: async (file: File): Promise<ApiResponse<{ url: string; filename: string; size: number }>> => {
    const formData = new FormData()
    formData.append('file', file)
    const res = await fetch(`${BASE_URL}/admin/upload`, {
      method: 'POST',
      headers: authHeaders(),
      body: formData
    })
    if (!res.ok) {
      const err = await res.json().catch(() => ({ message: res.statusText }))
      throw { status: res.status, data: err }
    }
    return res.json()
  },
}
