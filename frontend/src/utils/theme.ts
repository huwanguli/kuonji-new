// 暗色模式管理
import { ref } from 'vue'

const THEME_KEY = 'kuonji_theme'

const isDark = ref(false)

export function useTheme() {
  function init() {
    const saved = localStorage.getItem(THEME_KEY)
    if (saved === 'dark') {
      isDark.value = true
    } else if (saved === 'light') {
      isDark.value = false
    } else {
      isDark.value = window.matchMedia('(prefers-color-scheme: dark)').matches
    }
    apply()
  }

  function toggle() {
    isDark.value = !isDark.value
    localStorage.setItem(THEME_KEY, isDark.value ? 'dark' : 'light')
    apply()
  }

  function apply() {
    document.documentElement.classList.toggle('dark', isDark.value)
  }

  return { isDark, init, toggle }
}
