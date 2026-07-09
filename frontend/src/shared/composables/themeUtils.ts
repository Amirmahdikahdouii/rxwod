export type ThemeMode = 'dark'

export const THEME_STORAGE_KEY = 'rxwod-theme'

export function getSystemTheme(): ThemeMode {
  return 'dark'
}

export function getStoredTheme(): ThemeMode | null {
  if (typeof window === 'undefined') {
    return null
  }
  const stored = localStorage.getItem(THEME_STORAGE_KEY)
  if (stored === 'dark') {
    return stored
  }
  return null
}

export function resolveInitialTheme(): ThemeMode {
  return 'dark'
}

export function applyTheme(theme: ThemeMode = 'dark'): void {
  document.documentElement.setAttribute('data-theme', theme)
}

export function persistTheme(theme: ThemeMode = 'dark'): void {
  localStorage.setItem(THEME_STORAGE_KEY, theme)
}

export function initTheme(): ThemeMode {
  const theme = resolveInitialTheme()
  applyTheme(theme)
  persistTheme(theme)
  return theme
}
