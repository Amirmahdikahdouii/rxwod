import { computed, ref } from 'vue'
import { applyTheme, persistTheme, type ThemeMode } from '@/shared/composables/themeUtils'

const theme = ref<ThemeMode>('dark')

applyTheme('dark')

export function useTheme() {
  const isDark = computed(() => true)

  function setTheme(_mode: ThemeMode) {
    applyTheme('dark')
    persistTheme('dark')
  }

  function toggleTheme() {
    // Dark-only theme; no-op for backwards compatibility.
  }

  return {
    theme,
    isDark,
    setTheme,
    toggleTheme,
  }
}
