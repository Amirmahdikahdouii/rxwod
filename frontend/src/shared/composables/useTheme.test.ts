import { beforeEach, describe, expect, it } from 'vitest'
import {
  applyTheme,
  getStoredTheme,
  getSystemTheme,
  persistTheme,
  resolveInitialTheme,
  THEME_STORAGE_KEY,
} from '@/shared/composables/themeUtils'

describe('themeUtils', () => {
  beforeEach(() => {
    localStorage.clear()
    document.documentElement.removeAttribute('data-theme')
  })

  it('persists and reads theme from localStorage', () => {
    persistTheme('dark')
    expect(getStoredTheme()).toBe('dark')
  })

  it('resolves stored theme over system default', () => {
    persistTheme('light')
    expect(resolveInitialTheme()).toBe('light')
  })

  it('falls back to system theme when nothing stored', () => {
    expect(resolveInitialTheme()).toBe(getSystemTheme())
  })

  it('applies data-theme attribute to html', () => {
    applyTheme('dark')
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark')
  })

  it('ignores invalid stored values', () => {
    localStorage.setItem(THEME_STORAGE_KEY, 'invalid')
    expect(getStoredTheme()).toBeNull()
  })
})

describe('useTheme', () => {
  beforeEach(() => {
    localStorage.clear()
    document.documentElement.removeAttribute('data-theme')
  })

  it('toggles theme and persists preference', async () => {
    const { useTheme } = await import('@/shared/composables/useTheme')
    const { theme, toggleTheme } = useTheme()

    const initial = theme.value
    toggleTheme()

    expect(theme.value).not.toBe(initial)
    expect(localStorage.getItem(THEME_STORAGE_KEY)).toBe(theme.value)
    expect(document.documentElement.getAttribute('data-theme')).toBe(theme.value)
  })
})
