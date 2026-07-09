<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useSession } from '@/features/auth/composables/useSession'
import BaseInput from '@/shared/components/BaseInput.vue'

const router = useRouter()
const route = useRoute()
const session = useSession()

const displayName = ref('')
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
  const response = await session.register({
    displayName: displayName.value,
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
        <p class="eyebrow">Start coaching smarter</p>
        <h1 class="page-title">Create your RXWOD account</h1>
        <p class="page-subtitle">
          Register with your email, then create your first gym workspace or accept an invite from a gym owner.
        </p>
      </div>

      <form class="stack" @submit.prevent="submit">
        <BaseInput v-model="displayName" label="Display name" placeholder="Alex Coach" />
        <BaseInput v-model="email" label="Email" placeholder="alex@gym.com" />
        <BaseInput v-model="password" label="Password" type="password" placeholder="At least 8 characters" />
        <p class="helper-text">Use at least 8 characters. You can switch gyms after signing in.</p>
        <div v-if="error" class="alert alert--error" role="alert">{{ error }}</div>
        <button type="submit" class="btn-full" :disabled="loading">
          {{ loading ? 'Creating account...' : 'Create Account' }}
        </button>
      </form>

      <p class="auth-card__footer">
        Already have an account?
        <RouterLink to="/login">Login</RouterLink>
      </p>
    </section>
  </div>
</template>
