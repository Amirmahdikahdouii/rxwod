<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useSession } from '@/features/auth/composables/useSession'
import { acceptInvitation, getInvitationPreview } from '@/features/workspace/api/workspaceApi'
import type { InvitationPreviewResponse } from '@/features/workspace/model/workspaceTypes'
import { ROLE_LABELS } from '@/features/workspace/model/workspaceTypes'

const route = useRoute()
const router = useRouter()
const session = useSession()

const token = computed(() => String(route.params.token ?? ''))
const preview = ref<InvitationPreviewResponse | null>(null)
const loading = ref(true)
const accepting = ref(false)
const error = ref<string | null>(null)
const success = ref<string | null>(null)

const roleLabel = computed(() =>
  preview.value ? ROLE_LABELS[preview.value.role] : '',
)

const isAlreadyMember = computed(() =>
  preview.value
    ? session.workspaces.value.some((workspace) => workspace.id === preview.value?.gymId)
    : false,
)

const emailMismatch = computed(() =>
  preview.value && session.currentUser.value
    ? session.currentUser.value.email.toLowerCase() !== preview.value.email.toLowerCase()
    : false,
)

const inviteRedirect = computed(() => `/invite/${encodeURIComponent(token.value)}`)

async function loadPreview() {
  loading.value = true
  error.value = null
  const response = await getInvitationPreview(token.value)
  loading.value = false

  if (!response.ok) {
    error.value = response.error
    return
  }

  preview.value = response.value
}

async function accept() {
  if (!preview.value) {
    return
  }

  accepting.value = true
  error.value = null
  const response = await acceptInvitation(preview.value.gymId, token.value)
  accepting.value = false

  if (!response.ok) {
    error.value = response.error
    return
  }

  await session.loadMe()
  session.setActiveWorkspace(preview.value.gymId)
  success.value = `Welcome to ${preview.value.gymName}!`
  await router.push('/')
}

onMounted(async () => {
  if (!token.value) {
    error.value = 'Invalid invitation link'
    loading.value = false
    return
  }

  if (session.isAuthenticated.value && !session.sessionReady.value) {
    await session.loadMe()
  }

  await loadPreview()
})
</script>

<template>
  <div class="auth-page">
    <section class="auth-card card">
      <div class="auth-card__intro">
        <p class="eyebrow">Gym invitation</p>
        <h1 class="page-title">Join a workspace on RXWOD</h1>
        <p v-if="loading" class="page-subtitle">Loading invitation details...</p>
      </div>

      <div v-if="error" class="alert alert--error" role="alert">{{ error }}</div>
      <div v-if="success" class="alert alert--success" role="status">{{ success }}</div>

      <template v-if="!loading && preview">
        <div v-if="preview.status !== 'pending'" class="stack">
          <p class="page-subtitle">
            <template v-if="preview.status === 'expired'">
              This invitation to <strong>{{ preview.gymName }}</strong> has expired.
            </template>
            <template v-else-if="preview.status === 'accepted'">
              This invitation to <strong>{{ preview.gymName }}</strong> has already been accepted.
            </template>
            <template v-else-if="preview.status === 'revoked'">
              This invitation to <strong>{{ preview.gymName }}</strong> is no longer valid.
            </template>
          </p>
          <RouterLink to="/login" class="btn-full">Go to login</RouterLink>
        </div>

        <div v-else-if="!session.isAuthenticated.value" class="stack">
          <p class="page-subtitle">
            You've been invited to join <strong>{{ preview.gymName }}</strong> as
            <strong>{{ roleLabel }}</strong>.
          </p>
          <p class="helper-text">Invitation sent to {{ preview.email }}</p>
          <RouterLink
            :to="{ path: '/register', query: { redirect: inviteRedirect } }"
            class="btn-full"
          >
            Create account
          </RouterLink>
          <RouterLink
            :to="{ path: '/login', query: { redirect: inviteRedirect } }"
            class="empty-state__link empty-state__link--secondary"
          >
            Already have an account? Log in
          </RouterLink>
        </div>

        <div v-else-if="isAlreadyMember" class="stack">
          <p class="page-subtitle">
            You're already a member of <strong>{{ preview.gymName }}</strong>.
          </p>
          <button type="button" class="btn-full" @click="session.setActiveWorkspace(preview.gymId); router.push('/')">
            Open workspace
          </button>
        </div>

        <div v-else-if="emailMismatch" class="stack">
          <p class="page-subtitle">
            This invitation was sent to <strong>{{ preview.email }}</strong>, but you're logged in as
            <strong>{{ session.currentUser.value?.email }}</strong>.
          </p>
          <p class="helper-text">Log out and sign in with the invited email, or ask the gym owner to send a new invite.</p>
          <button type="button" class="btn-full" @click="session.logout(); router.push({ path: '/login', query: { redirect: inviteRedirect } })">
            Log out and switch account
          </button>
        </div>

        <div v-else class="stack">
          <p class="page-subtitle">
            Accept your invitation to join <strong>{{ preview.gymName }}</strong> as
            <strong>{{ roleLabel }}</strong>.
          </p>
          <p class="helper-text">Signed in as {{ session.currentUser.value?.email }}</p>
          <button type="button" class="btn-full" :disabled="accepting" @click="accept">
            {{ accepting ? 'Accepting...' : 'Accept invitation' }}
          </button>
        </div>
      </template>
    </section>
  </div>
</template>
