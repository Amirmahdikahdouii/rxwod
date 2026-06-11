import { computed, ref } from 'vue'
import {
  applyTheme,
  persistTheme,
  resolveInitialTheme,
  type ThemeMode,
} from '@/shared/composables/themeUtils'

function readThemeFromDOM(): ThemeMode | null {
  const attr = document.documentElement.getAttribute('data-theme')
  if (attr === 'light' || attr === 'dark') {
    return attr
  }
  return null
}

const theme = ref<ThemeMode>(readThemeFromDOM() ?? resolveInitialTheme())

export function useTheme() {
  const isDark = computed(() => theme.value === 'dark')

  function setTheme(mode: ThemeMode) {
    theme.value = mode
    applyTheme(mode)
    persistTheme(mode)
  }

  function toggleTheme() {
    setTheme(theme.value === 'dark' ? 'light' : 'dark')
  }

  return {
    theme,
    isDark,
    setTheme,
    toggleTheme,
  }
}
