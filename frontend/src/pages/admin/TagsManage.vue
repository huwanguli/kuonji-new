<template>
  <div>
    <div class="page-header">
      <h1 class="page-heading">标签管理</h1>
      <form class="inline-form" @submit.prevent="addTag">
        <input v-model="newName" class="input" placeholder="新标签名称" required />
        <button type="submit" class="btn btn-primary">添加</button>
      </form>
    </div>
    <div class="tag-list">
      <div v-for="t in tags" :key="t.id" class="tag-item">
        <span class="tag-name">#{{ t.name }}</span>
        <button class="tag-delete" @click="remove(t.id)">×</button>
      </div>
    </div>
    <p v-if="!tags.length" class="empty-state">暂无标签</p>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type Tag } from '@/utils/api'

const tags = ref<Tag[]>([])
const newName = ref('')

async function load() {
  try { tags.value = (await api.getTags()).data } catch (e) { console.error(e) }
}

async function addTag() {
  if (!newName.value.trim()) return
  try { await api.adminCreateTag(newName.value.trim()); newName.value = ''; load() } catch (e: any) { alert(e?.data?.message || '添加失败') }
}

async function remove(id: number) {
  if (!confirm('确定删除此标签？')) return
  try { await api.adminDeleteTag(id); load() } catch (e: any) { alert(e?.data?.message || '删除失败') }
}

onMounted(load)
</script>

<style scoped>
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px; flex-wrap: wrap; gap: 12px; }
.page-heading { font-size: 22px; font-weight: 500; color: #1f1f1f; }
.inline-form { display: flex; gap: 8px; }
.tag-list { display: flex; flex-wrap: wrap; gap: 8px; }
.tag-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 14px;
  background: var(--white);
  border: 1px solid var(--line);
  border-radius: var(--r-full);
  font-size: 13px;
}
.tag-name { color: var(--ink); }
.tag-delete {
  background: none;
  border: none;
  font-size: 16px;
  color: var(--muted-light);
  cursor: pointer;
  line-height: 1;
  padding: 0 2px;
}
.tag-delete:hover { color: #c53f2c; }
</style>
