<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { useSession } from '@/features/auth/composables/useSession'
import { useWorkspaces } from '@/features/workspace/composables/useWorkspaces'
import {
  canCreateWOD,
  canManageMemberTarget,
  canManageMembers,
  ROLE_LABELS,
} from '@/features/workspace/model/workspaceTypes'
import BaseInput from '@/shared/components/BaseInput.vue'

const session = useSession()
const workspaces = useWorkspaces()

const coachEmail = ref('')
const athleteEmail = ref('')
const invitingCoach = ref(false)
const invitingAthlete = ref(false)
const lastInviteLink = ref<string | null>(null)
const copyFeedback = ref<string | null>(null)

function buildInviteLink(token: string) {
  return `${window.location.origin}/invite/${encodeURIComponent(token)}`
}

async function copyInviteLink() {
  if (!lastInviteLink.value) {
    return
  }
  await navigator.clipboard.writeText(lastInviteLink.value)
  copyFeedback.value = 'Invite link copied'
  setTimeout(() => {
    copyFeedback.value = null
  }, 2000)
}

async function inviteCoach() {
  invitingCoach.value = true
  lastInviteLink.value = null
  const response = await workspaces.invite('coach', coachEmail.value)
  invitingCoach.value = false
  if (response.ok) {
    coachEmail.value = ''
    if (response.value.token) {
      lastInviteLink.value = buildInviteLink(response.value.token)
    }
  }
}

async function inviteAthlete() {
  invitingAthlete.value = true
  lastInviteLink.value = null
  const response = await workspaces.invite('athlete', athleteEmail.value)
  invitingAthlete.value = false
  if (response.ok) {
    athleteEmail.value = ''
    if (response.value.token) {
      lastInviteLink.value = buildInviteLink(response.value.token)
    }
  }
}

const roleLabel = computed(() =>
  session.activeWorkspaceRole.value ? ROLE_LABELS[session.activeWorkspaceRole.value] : 'Member',
)
const canManage = computed(() => canManageMembers(session.activeWorkspaceRole.value))
const canCreate = computed(() => canCreateWOD(session.activeWorkspaceRole.value))

onMounted(async () => {
  if (canManage.value) {
    await workspaces.refreshMembers()
  }
})

watch(
  () => session.activeWorkspaceId.value,
  async () => {
    if (canManage.value) {
      await workspaces.refreshMembers()
    }
  },
)
</script>

<template>
  <div class="container stack-lg">
    <header class="page-header">
      <p class="eyebrow">Workspace</p>
      <h1 class="page-title">{{ session.activeWorkspace.value?.name ?? 'Gym Workspace' }}</h1>
      <p class="page-subtitle">
        You are working as <strong>{{ roleLabel }}</strong>. Workspace permissions control what you can do.
      </p>
    </header>

    <section class="dashboard-grid">
      <article class="card stack">
        <h2 class="section-title">Your Access</h2>
        <span class="role-pill role-pill--large">{{ roleLabel }}</span>
        <p v-if="canManage" class="page-subtitle page-subtitle--flush">
          You can manage members, invite coaches and athletes, and update gym programs.
        </p>
        <p v-else-if="canCreate" class="page-subtitle page-subtitle--flush">
          You can create new plans and view saved programs in this gym.
        </p>
        <p v-else class="page-subtitle page-subtitle--flush">
          You can view plans assigned to this gym. Creation and editing are owner or coach actions.
        </p>
        <div class="row row--align-center">
          <RouterLink v-if="canCreate" to="/" class="empty-state__link">Create Program</RouterLink>
          <RouterLink to="/wods" class="empty-state__link empty-state__link--secondary">
            View Programs
          </RouterLink>
        </div>
      </article>

      <article class="card stack">
        <h2 class="section-title">Workspace Switch</h2>
        <div class="workspace-list">
          <button
            v-for="workspace in session.workspaces.value"
            :key="workspace.id"
            type="button"
            class="workspace-option"
            :class="{ 'workspace-option--active': workspace.id === session.activeWorkspaceId.value }"
            @click="session.setActiveWorkspace(workspace.id)"
          >
            <span>{{ workspace.name }}</span>
            <span class="role-pill">{{ ROLE_LABELS[workspace.role] }}</span>
          </button>
        </div>
        <RouterLink to="/workspace/new" class="empty-state__link">Create another gym</RouterLink>
      </article>
    </section>

    <section v-if="canManage" class="dashboard-grid">
      <article class="card stack">
        <h2 class="section-title">Invite Team</h2>
        <form class="invite-form" @submit.prevent="inviteCoach">
          <BaseInput v-model="coachEmail" label="Coach email" placeholder="coach@gym.com" />
          <button type="submit" :disabled="invitingCoach">
            {{ invitingCoach ? 'Inviting...' : 'Invite Coach' }}
          </button>
        </form>
        <form class="invite-form" @submit.prevent="inviteAthlete">
          <BaseInput v-model="athleteEmail" label="Athlete email" placeholder="athlete@email.com" />
          <button type="submit" :disabled="invitingAthlete">
            {{ invitingAthlete ? 'Inviting...' : 'Invite Athlete' }}
          </button>
        </form>
        <div v-if="workspaces.workspaceError.value" class="alert alert--error" role="alert">
          {{ workspaces.workspaceError.value }}
        </div>
        <div v-if="workspaces.workspaceSuccess.value" class="alert alert--success" role="status">
          {{ workspaces.workspaceSuccess.value }}
        </div>
        <div v-if="lastInviteLink" class="stack">
          <p class="helper-text">Share this invite link with the recipient:</p>
          <div class="row row--align-center">
            <input class="invite-link-input" type="text" readonly :value="lastInviteLink" />
            <button type="button" class="secondary compact-button" @click="copyInviteLink">Copy link</button>
          </div>
          <p v-if="copyFeedback" class="helper-text">{{ copyFeedback }}</p>
        </div>
      </article>

      <article class="card stack">
        <div class="section-heading-row">
          <h2 class="section-title">Members</h2>
          <button type="button" class="secondary compact-button" @click="workspaces.refreshMembers">
            Refresh
          </button>
        </div>
        <p v-if="workspaces.loadingMembers.value" class="loading-state">Loading members...</p>
        <div v-else-if="workspaces.members.value.length === 0" class="empty-state empty-state--compact">
          <p class="empty-state__text">No members to show yet.</p>
        </div>
        <div v-else class="member-list">
          <article v-for="member in workspaces.members.value" :key="member.userId" class="member-card">
            <strong class="member-card__name">{{ member.displayName || member.email }}</strong>
            <div class="member-card__meta">
              <span class="member-card__email">{{ member.email }}</span>
              <span class="role-pill">{{ ROLE_LABELS[member.role] }}</span>
              <RouterLink
                v-if="canManage && canManageMemberTarget(member.role)"
                :to="`/workspace/members/${member.userId}`"
                class="member-card__manage"
              >
                Manage
              </RouterLink>
            </div>
          </article>
        </div>
      </article>
    </section>
  </div>
</template>
