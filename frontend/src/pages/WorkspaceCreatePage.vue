<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useSession } from '@/features/auth/composables/useSession'
import { useWorkspaces } from '@/features/workspace/composables/useWorkspaces'
import BaseInput from '@/shared/components/BaseInput.vue'

const router = useRouter()
const session = useSession()
const workspaces = useWorkspaces()

const name = ref('')
const loading = ref(false)

async function submit() {
  loading.value = true
  const response = await workspaces.createGym(name.value)
  loading.value = false

  if (response.ok) {
    await router.push('/')
  }
}
</script>

<template>
  <div class="container stack-lg">
    <section v-if="!session.isEmailVerified.value" class="workspace-hero card">
      <div class="workspace-hero__content">
        <p class="eyebrow">Create workspace</p>
        <h1 class="page-title">Confirm your email first</h1>
        <p class="page-subtitle">
          We sent a confirmation link to
          <strong>{{ session.currentUser.value?.email }}</strong>.
          Verify your email before creating a gym workspace.
        </p>
        <div class="empty-state empty-state--compact">
          <p class="empty-state__text">Use the banner above to resend the confirmation email.</p>
        </div>
      </div>
    </section>

    <section v-else class="workspace-hero card">
      <div class="workspace-hero__content">
        <p class="eyebrow">Create workspace</p>
        <h1 class="page-title">Set up your gym</h1>
        <p class="page-subtitle">
          A gym workspace keeps programs, coaches, and athletes scoped to the right community.
          You will become the owner and can invite your team next.
        </p>
      </div>

      <form class="workspace-hero__form stack" @submit.prevent="submit">
        <BaseInput v-model="name" label="Gym name" placeholder="Downtown CrossFit" />
        <div v-if="workspaces.workspaceError.value" class="alert alert--error" role="alert">
          {{ workspaces.workspaceError.value }}
        </div>
        <button type="submit" class="btn-full" :disabled="loading">
          {{ loading ? 'Creating...' : 'Create Gym Workspace' }}
        </button>
      </form>
    </section>
  </div>
</template>
