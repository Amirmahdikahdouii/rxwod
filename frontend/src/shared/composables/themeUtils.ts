export type ThemeMode = 'light' | 'dark'

export const THEME_STORAGE_KEY = 'rxwod-theme'

export function getSystemTheme(): ThemeMode {
  if (typeof window === 'undefined') {
    return 'light'
  }
  if (typeof window.matchMedia !== 'function') {
    return 'light'
  }
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

export function getStoredTheme(): ThemeMode | null {
  if (typeof window === 'undefined') {
    return null
  }
  const stored = localStorage.getItem(THEME_STORAGE_KEY)
  if (stored === 'light' || stored === 'dark') {
    return stored
  }
  return null
}

export function resolveInitialTheme(): ThemeMode {
  return getStoredTheme() ?? getSystemTheme()
}

export function applyTheme(theme: ThemeMode): void {
  document.documentElement.setAttribute('data-theme', theme)
}

export function persistTheme(theme: ThemeMode): void {
  localStorage.setItem(THEME_STORAGE_KEY, theme)
}

export function initTheme(): ThemeMode {
  const theme = resolveInitialTheme()
  applyTheme(theme)
  return theme
}
