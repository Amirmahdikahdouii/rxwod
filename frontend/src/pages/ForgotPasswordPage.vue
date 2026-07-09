<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import { forgotPassword } from '@/features/auth/api/authApi'
import BaseInput from '@/shared/components/BaseInput.vue'

const email = ref('')
const loading = ref(false)
const submitted = ref(false)
const error = ref<string | null>(null)

async function submit() {
  if (!email.value.trim()) {
    error.value = 'Email is required'
    return
  }

  loading.value = true
  error.value = null
  const response = await forgotPassword({ email: email.value.trim() })
  loading.value = false

  if (!response.ok) {
    error.value = response.error
    return
  }

  submitted.value = true
}
</script>

<template>
  <div class="auth-page">
    <section class="auth-card card">
      <div class="auth-card__intro">
        <p class="eyebrow">Password recovery</p>
        <h1 class="page-title">Reset your password</h1>
        <p class="page-subtitle">
          Enter the email address for your RXWOD account and we'll send you a link to choose a new password.
        </p>
      </div>

      <div v-if="submitted" class="stack">
        <div class="alert alert--success" role="status">
          If an account exists for that email, we sent a reset link. Check your inbox.
        </div>
        <RouterLink to="/login" class="btn-full">Back to login</RouterLink>
      </div>

      <template v-else>
        <form class="stack" @submit.prevent="submit">
          <BaseInput v-model="email" label="Email" placeholder="owner@gym.com" />
          <div v-if="error" class="alert alert--error" role="alert">{{ error }}</div>
          <button type="submit" class="btn-full" :disabled="loading">
            {{ loading ? 'Sending...' : 'Send reset link' }}
          </button>
        </form>

        <p class="auth-card__footer">
          Remember your password?
          <RouterLink to="/login">Back to login</RouterLink>
        </p>
      </template>
    </section>
  </div>
</template>
