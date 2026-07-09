<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { verifyEmail } from '@/features/auth/api/authApi'
import { useSession } from '@/features/auth/composables/useSession'

const route = useRoute()
const session = useSession()

const token = computed(() => String(route.params.token ?? ''))
const loading = ref(true)
const success = ref(false)
const error = ref<string | null>(null)
const invalidLink = ref(false)

async function confirmEmail() {
  loading.value = true
  error.value = null
  const response = await verifyEmail({ token: token.value })
  loading.value = false

  if (!response.ok) {
    error.value = response.error
    return
  }

  success.value = true
  if (session.isAuthenticated.value) {
    await session.loadMe()
  }
}

onMounted(async () => {
  if (!token.value) {
    invalidLink.value = true
    loading.value = false
    return
  }
  await confirmEmail()
})
</script>

<template>
  <div class="auth-page">
    <section class="auth-card card">
      <div class="auth-card__intro">
        <p class="eyebrow">Account setup</p>
        <h1 class="page-title">Confirm your email</h1>
        <p v-if="loading" class="page-subtitle">Verifying your email address...</p>
      </div>

      <div v-if="invalidLink" class="stack">
        <div class="alert alert--error" role="alert">Invalid confirmation link</div>
        <RouterLink to="/login" class="btn-full">Back to login</RouterLink>
      </div>

      <div v-else-if="success" class="stack">
        <div class="alert alert--success" role="status">
          Your email is confirmed. You can manage gym workspaces now.
        </div>
        <RouterLink :to="session.isAuthenticated.value ? '/profile' : '/login'" class="btn-full">
          {{ session.isAuthenticated.value ? 'Go to profile' : 'Sign in' }}
        </RouterLink>
      </div>

      <div v-else-if="error" class="stack">
        <div class="alert alert--error" role="alert">{{ error }}</div>
        <button type="button" class="btn-full" :disabled="loading" @click="confirmEmail">
          Try again
        </button>
        <RouterLink to="/login" class="empty-state__link empty-state__link--secondary">
          Back to login
        </RouterLink>
      </div>
    </section>
  </div>
</template>
