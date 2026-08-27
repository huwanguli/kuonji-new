<template>
  <div>
    <div class="page-header">
      <h1 class="page-heading">系列管理</h1>
      <form class="inline-form" @submit.prevent="addSeries">
        <input v-model="newName" class="input" placeholder="新系列名称" required />
        <button type="submit" class="btn btn-primary">添加</button>
      </form>
    </div>
    <div class="table-wrap">
      <table class="admin-table">
        <thead><tr><th>ID</th><th>名称</th><th>描述</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-for="s in seriesList" :key="s.id">
            <td class="td-id">{{ s.id }}</td>
            <td>{{ s.name }}</td>
            <td class="td-desc">{{ s.description || '-' }}</td>
            <td><button class="action-link action-delete" @click="remove(s.id)">删除</button></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type Series } from '@/utils/api'

const seriesList = ref<Series[]>([])
const newName = ref('')

async function load() {
  try { seriesList.value = (await api.getSeries()).data } catch (e) { console.error(e) }
}

async function addSeries() {
  if (!newName.value.trim()) return
  try { await api.adminCreateSeries({ name: newName.value.trim() }); newName.value = ''; load() } catch (e: any) { alert(e?.data?.message || '添加失败') }
}

async function remove(id: number) {
  if (!confirm('确定删除此系列？')) return
  try { await api.adminDeleteSeries(id); load() } catch (e: any) { alert(e?.data?.message || '删除失败') }
}

onMounted(load)
</script>

<style scoped>
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px; flex-wrap: wrap; gap: 12px; }
.page-heading { font-size: 22px; font-weight: 500; color: #1f1f1f; }
.inline-form { display: flex; gap: 8px; }
.table-wrap { overflow-x: auto; }
.admin-table { width: 100%; border-collapse: collapse; background: var(--white); border: 1px solid var(--line); border-radius: var(--r-md); font-size: 13px; }
.admin-table th { background: #f8f8f6; padding: 10px 14px; text-align: left; font-weight: 500; color: var(--muted); border-bottom: 1px solid var(--line); }
.admin-table td { padding: 12px 14px; border-bottom: 1px solid var(--line); color: var(--ink); }
.td-id { color: var(--muted-light); width: 50px; }
.td-desc { color: var(--muted); }
.action-link { font-size: 12px; background: none; border: none; cursor: pointer; font-family: inherit; padding: 2px 6px; }
.action-delete { color: #c53f2c; }
.action-delete:hover { text-decoration: underline; }
</style>
