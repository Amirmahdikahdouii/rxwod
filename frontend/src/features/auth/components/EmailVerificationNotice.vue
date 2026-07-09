<script setup lang="ts">
import { ref } from 'vue'
import { resendVerificationEmail } from '@/features/auth/api/authApi'
import { useSession } from '@/features/auth/composables/useSession'

const session = useSession()

const loading = ref(false)
const feedback = ref<string | null>(null)
const error = ref<string | null>(null)

async function resend() {
  loading.value = true
  feedback.value = null
  error.value = null
  const response = await resendVerificationEmail()
  loading.value = false

  if (!response.ok) {
    error.value = response.error
    return
  }

  feedback.value = 'If your account is unverified, we sent a new confirmation link.'
}
</script>

<template>
  <div
    v-if="session.isAuthenticated.value && session.sessionReady.value && !session.isEmailVerified.value"
    class="email-verification-notice alert alert--warning"
    role="status"
  >
    <div class="email-verification-notice__content">
      <strong>Email confirmation pending</strong>
      <p>
        Check your inbox for a link to confirm
        <strong>{{ session.currentUser.value?.email }}</strong>
        before creating or managing gym workspaces.
      </p>
      <p v-if="feedback" class="email-verification-notice__feedback">{{ feedback }}</p>
      <p v-if="error" class="email-verification-notice__error">{{ error }}</p>
    </div>
    <button type="button" class="secondary" :disabled="loading" @click="resend">
      {{ loading ? 'Sending...' : 'Resend confirmation email' }}
    </button>
  </div>
</template>

<style scoped>
.email-verification-notice {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
  margin: 0;
  border-radius: 0;
  border-left: none;
  border-right: none;
}

.email-verification-notice__content p {
  margin: 0.35rem 0 0;
}

.email-verification-notice__feedback {
  color: var(--color-success, #166534);
}

.email-verification-notice__error {
  color: var(--color-danger, #991b1b);
}

@media (max-width: 640px) {
  .email-verification-notice {
    flex-direction: column;
  }
}
</style>
