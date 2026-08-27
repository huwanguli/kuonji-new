import type { RouteRecordRaw } from 'vue-router'

export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/pages/HomePage.vue')
  },
  {
    path: '/articles/:id',
    name: 'article',
    component: () => import('@/pages/ArticleDetail.vue')
  },
  {
    path: '/categories',
    name: 'categories',
    component: () => import('@/pages/CategoriesPage.vue')
  },
  {
    path: '/categories/:id',
    name: 'category-detail',
    component: () => import('@/pages/CategoryDetail.vue')
  },
  {
    path: '/tags',
    name: 'tags',
    component: () => import('@/pages/TagsPage.vue')
  },
  {
    path: '/tags/:id',
    name: 'tag-detail',
    component: () => import('@/pages/TagDetail.vue')
  },
  {
    path: '/series',
    name: 'series',
    component: () => import('@/pages/SeriesPage.vue')
  },
  {
    path: '/series/:id',
    name: 'series-detail',
    component: () => import('@/pages/SeriesDetail.vue')
  },
  {
    path: '/search',
    name: 'search',
    component: () => import('@/pages/SearchPage.vue')
  },
  // ── 后台管理 ──
  {
    path: '/admin/login',
    name: 'admin-login',
    component: () => import('@/pages/admin/LoginPage.vue')
  },
  {
    path: '/admin',
    name: 'admin',
    component: () => import('@/pages/admin/AdminLayout.vue'),
    children: [
      { path: '', name: 'admin-dashboard', component: () => import('@/pages/admin/DashboardPage.vue') },
      { path: 'articles', name: 'admin-articles', component: () => import('@/pages/admin/ArticlesManage.vue') },
      { path: 'articles/new', name: 'admin-article-new', component: () => import('@/pages/admin/ArticleEditor.vue') },
      { path: 'articles/:id/edit', name: 'admin-article-edit', component: () => import('@/pages/admin/ArticleEditor.vue') },
      { path: 'categories', name: 'admin-categories', component: () => import('@/pages/admin/CategoriesManage.vue') },
      { path: 'tags', name: 'admin-tags', component: () => import('@/pages/admin/TagsManage.vue') },
      { path: 'series', name: 'admin-series', component: () => import('@/pages/admin/SeriesManage.vue') },
      { path: 'comments', name: 'admin-comments', component: () => import('@/pages/admin/CommentsManage.vue') },
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]
