import { beforeEach, describe, expect, it } from 'vitest'
import {
  applyTheme,
  getStoredTheme,
  getSystemTheme,
  initTheme,
  persistTheme,
  resolveInitialTheme,
  THEME_STORAGE_KEY,
} from '@/shared/composables/themeUtils'

describe('themeUtils', () => {
  beforeEach(() => {
    localStorage.clear()
    document.documentElement.removeAttribute('data-theme')
  })

  it('persists and reads dark theme from localStorage', () => {
    persistTheme('dark')
    expect(getStoredTheme()).toBe('dark')
  })

  it('always resolves to dark theme', () => {
    persistTheme('dark')
    expect(resolveInitialTheme()).toBe('dark')
    localStorage.clear()
    expect(resolveInitialTheme()).toBe('dark')
  })

  it('always reports dark as system theme', () => {
    expect(getSystemTheme()).toBe('dark')
  })

  it('applies data-theme attribute to html', () => {
    applyTheme('dark')
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark')
  })

  it('ignores invalid stored values', () => {
    localStorage.setItem(THEME_STORAGE_KEY, 'invalid')
    expect(getStoredTheme()).toBeNull()
  })

  it('initTheme forces dark and persists preference', () => {
    localStorage.setItem(THEME_STORAGE_KEY, 'light')
    expect(initTheme()).toBe('dark')
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark')
    expect(localStorage.getItem(THEME_STORAGE_KEY)).toBe('dark')
  })
})

describe('useTheme', () => {
  beforeEach(() => {
    localStorage.clear()
    document.documentElement.removeAttribute('data-theme')
  })

  it('always exposes dark theme', async () => {
    const { useTheme } = await import('@/shared/composables/useTheme')
    const { theme, isDark, toggleTheme } = useTheme()

    expect(theme.value).toBe('dark')
    expect(isDark.value).toBe(true)

    toggleTheme()

    expect(theme.value).toBe('dark')
    expect(isDark.value).toBe(true)
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark')
  })
})
