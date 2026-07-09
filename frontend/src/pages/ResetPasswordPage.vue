<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { resetPassword } from '@/features/auth/api/authApi'
import BaseInput from '@/shared/components/BaseInput.vue'

const route = useRoute()

const token = computed(() => String(route.params.token ?? ''))
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const success = ref(false)
const error = ref<string | null>(null)
const invalidLink = ref(false)

async function submit() {
  if (password.value !== confirmPassword.value) {
    error.value = 'Passwords do not match'
    return
  }

  loading.value = true
  error.value = null
  const response = await resetPassword({
    token: token.value,
    password: password.value,
  })
  loading.value = false

  if (!response.ok) {
    error.value = response.error
    return
  }

  success.value = true
}

onMounted(() => {
  if (!token.value) {
    invalidLink.value = true
  }
})
</script>

<template>
  <div class="auth-page">
    <section class="auth-card card">
      <div class="auth-card__intro">
        <p class="eyebrow">Password recovery</p>
        <h1 class="page-title">Choose a new password</h1>
        <p v-if="!invalidLink && !success" class="page-subtitle">
          Enter and confirm your new password to finish resetting your RXWOD account.
        </p>
      </div>

      <div v-if="invalidLink" class="stack">
        <div class="alert alert--error" role="alert">Invalid reset link</div>
        <RouterLink to="/forgot-password" class="btn-full">Request a new reset link</RouterLink>
        <RouterLink to="/login" class="empty-state__link empty-state__link--secondary">
          Back to login
        </RouterLink>
      </div>

      <div v-else-if="success" class="stack">
        <div class="alert alert--success" role="status">
          Your password has been reset. You can now log in with your new password.
        </div>
        <RouterLink to="/login" class="btn-full">Go to login</RouterLink>
      </div>

      <template v-else>
        <form class="stack" @submit.prevent="submit">
          <BaseInput
            v-model="password"
            label="New password"
            type="password"
            placeholder="At least 8 characters"
          />
          <BaseInput
            v-model="confirmPassword"
            label="Confirm password"
            type="password"
            placeholder="Re-enter your password"
          />
          <p class="helper-text">Use at least 8 characters.</p>
          <div v-if="error" class="alert alert--error" role="alert">{{ error }}</div>
          <button type="submit" class="btn-full" :disabled="loading">
            {{ loading ? 'Resetting...' : 'Reset password' }}
          </button>
        </form>

        <p class="auth-card__footer">
          Need a new link?
          <RouterLink to="/forgot-password">Request reset link</RouterLink>
        </p>
      </template>
    </section>
  </div>
</template>
