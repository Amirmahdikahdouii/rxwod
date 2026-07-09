<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useSession } from '@/features/auth/composables/useSession'
import BaseInput from '@/shared/components/BaseInput.vue'

const router = useRouter()
const route = useRoute()
const session = useSession()

const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref<string | null>(null)

function redirectTarget() {
  const redirect = route.query.redirect
  if (typeof redirect === 'string' && redirect.startsWith('/')) {
    return redirect
  }
  return session.activeWorkspace.value ? '/' : '/workspace/new'
}

async function submit() {
  loading.value = true
  error.value = null
  const response = await session.login({
    email: email.value,
    password: password.value,
  })
  loading.value = false

  if (!response.ok) {
    error.value = response.error
    return
  }

  await router.push(redirectTarget())
}
</script>

<template>
  <div class="auth-page">
    <section class="auth-card card">
      <div class="auth-card__intro">
        <p class="eyebrow">Welcome back</p>
        <h1 class="page-title">Log in to your RXWOD workspace</h1>
        <p class="page-subtitle">
          Pick up where your gym left off: build plans, review saved workouts, and switch between workspaces.
        </p>
      </div>

      <form class="stack" @submit.prevent="submit">
        <BaseInput v-model="email" label="Email" placeholder="owner@gym.com" />
        <BaseInput v-model="password" label="Password" type="password" placeholder="Your password" />
        <RouterLink to="/forgot-password" class="empty-state__link empty-state__link--secondary">
          Forgot password?
        </RouterLink>
        <div v-if="error" class="alert alert--error" role="alert">{{ error }}</div>
        <button type="submit" class="btn-full" :disabled="loading">
          {{ loading ? 'Logging in...' : 'Login' }}
        </button>
      </form>

      <p class="auth-card__footer">
        New to RXWOD?
        <RouterLink to="/register">Create an account</RouterLink>
      </p>
    </section>
  </div>
</template>
