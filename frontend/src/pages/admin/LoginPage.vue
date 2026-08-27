<template>
  <div class="login-page">
    <div class="login-box">
      <h1 class="login-title">Kuonji</h1>
      <p class="login-sub">后台管理登录</p>
      <form @submit.prevent="doLogin" class="login-form">
        <input v-model="username" class="input" placeholder="用户名" required />
        <input v-model="password" class="input" type="password" placeholder="密码" required />
        <p v-if="error" class="login-error">{{ error }}</p>
        <button type="submit" class="btn btn-primary login-btn" :disabled="loading">
          {{ loading ? '登录中…' : '登录' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '@/utils/api'
import { setToken, setUser } from '@/utils/auth'

const router = useRouter()
const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function doLogin() {
  loading.value = true
  error.value = ''
  try {
    const res = await api.login(username.value, password.value)
    setToken(res.data.token)
    setUser(res.data.user)
    router.push('/admin')
  } catch (e: any) {
    error.value = e?.data?.message || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg);
}
.login-box {
  width: 340px;
  padding: 40px;
  background: var(--white);
  border: 1px solid var(--line);
  border-radius: var(--r-md);
}
.login-title {
  font-size: 24px;
  font-weight: 500;
  text-align: center;
  color: #1f1f1f;
  margin-bottom: 4px;
}
.login-sub {
  font-size: 13px;
  color: var(--muted);
  text-align: center;
  margin-bottom: 28px;
}
.login-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.login-error {
  font-size: 12px;
  color: #c53f2c;
}
.login-btn {
  margin-top: 4px;
  width: 100%;
  padding: 10px;
}
</style>
