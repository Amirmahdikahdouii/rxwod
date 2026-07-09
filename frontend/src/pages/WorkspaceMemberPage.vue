<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useSession } from '@/features/auth/composables/useSession'
import { useWorkspaces } from '@/features/workspace/composables/useWorkspaces'
import {
  canManageMemberTarget,
  canManageMembers,
  ROLE_LABELS,
  type WorkspaceRole,
} from '@/features/workspace/model/workspaceTypes'
import BaseSelect from '@/shared/components/BaseSelect.vue'

const route = useRoute()
const router = useRouter()
const session = useSession()
const workspaces = useWorkspaces()

const selectedRole = ref<'coach' | 'athlete'>('coach')
const pendingRoleChange = ref(false)
const pendingRemoval = ref(false)

const userId = computed(() => route.params.userId as string)
const canManage = computed(() => canManageMembers(session.activeWorkspaceRole.value))

const member = computed(() =>
  workspaces.members.value.find((entry) => entry.userId === userId.value),
)

const memberLabel = computed(() => member.value?.displayName || member.value?.email || 'Member')

const roleOptions = computed(() => [
  { value: 'coach', label: ROLE_LABELS.coach },
  { value: 'athlete', label: ROLE_LABELS.athlete },
])

const roleChanged = computed(
  () => member.value && selectedRole.value !== member.value.role,
)

async function loadMember() {
  if (!session.isEmailVerified.value || !canManage.value) {
    await router.replace('/workspace')
    return
  }

  if (workspaces.members.value.length === 0) {
    await workspaces.refreshMembers()
  }

  const current = workspaces.members.value.find((entry) => entry.userId === userId.value)
  if (!current || !canManageMemberTarget(current.role)) {
    await router.replace('/workspace')
    return
  }

  selectedRole.value = current.role as 'coach' | 'athlete'
  pendingRoleChange.value = false
  pendingRemoval.value = false
}

function requestRoleChange() {
  if (!roleChanged.value) {
    return
  }
  pendingRemoval.value = false
  pendingRoleChange.value = true
}

function cancelRoleChange() {
  pendingRoleChange.value = false
  if (member.value) {
    selectedRole.value = member.value.role as 'coach' | 'athlete'
  }
}

async function confirmRoleChange() {
  if (!member.value || !roleChanged.value) {
    return
  }

  const response = await workspaces.updateMemberRole(userId.value, selectedRole.value)
  pendingRoleChange.value = false

  if (response.ok) {
    selectedRole.value = response.value.role as 'coach' | 'athlete'
  }
}

function requestRemoval() {
  pendingRoleChange.value = false
  pendingRemoval.value = true
}

function cancelRemoval() {
  pendingRemoval.value = false
}

async function confirmRemoval() {
  const response = await workspaces.removeMember(userId.value)
  if (response.ok) {
    await router.push('/workspace')
  }
}

onMounted(loadMember)

watch(
  () => [session.activeWorkspaceId.value, userId.value],
  loadMember,
)
</script>

<template>
  <div class="container stack-lg">
    <header class="page-header">
      <RouterLink to="/workspace" class="empty-state__link empty-state__link--secondary">
        Back to workspace
      </RouterLink>
      <p class="eyebrow">Member</p>
      <h1 class="page-title">{{ memberLabel }}</h1>
      <p class="page-subtitle">Review member details and manage their workspace access.</p>
    </header>

    <p v-if="workspaces.loadingMembers.value" class="loading-state">Loading member...</p>

    <section v-else-if="member" class="dashboard-grid">
      <article class="card stack">
        <h2 class="section-title">Member Information</h2>
        <div class="profile-row">
          <span class="profile-avatar">{{ memberLabel.charAt(0) }}</span>
          <div>
            <strong>{{ member.displayName || member.email }}</strong>
            <p class="page-subtitle page-subtitle--flush">{{ member.email }}</p>
          </div>
        </div>
        <div class="member-detail-meta">
          <span class="role-pill">{{ ROLE_LABELS[member.role as WorkspaceRole] }}</span>
          <span class="status-pill">{{ member.status }}</span>
        </div>
      </article>

      <article class="card stack">
        <h2 class="section-title">Role</h2>
        <p class="page-subtitle page-subtitle--flush">
          Assign this member as a coach or athlete in {{ session.activeWorkspace.value?.name }}.
        </p>
        <BaseSelect
          v-model="selectedRole"
          label="Workspace role"
          :options="roleOptions"
          :disabled="workspaces.updatingMemberId.value === userId || pendingRoleChange"
          @update:model-value="pendingRoleChange = false"
        />
        <button
          type="button"
          :disabled="!roleChanged || workspaces.updatingMemberId.value === userId || pendingRoleChange"
          @click="requestRoleChange"
        >
          Update role
        </button>

        <div v-if="pendingRoleChange" class="confirmation-panel" role="alertdialog" aria-live="polite">
          <p class="confirmation-panel__text">
            Change {{ memberLabel }} from {{ ROLE_LABELS[member.role as WorkspaceRole] }} to
            {{ ROLE_LABELS[selectedRole] }}?
          </p>
          <div class="confirmation-panel__actions">
            <button type="button" class="secondary" @click="cancelRoleChange">Cancel</button>
            <button
              type="button"
              :disabled="workspaces.updatingMemberId.value === userId"
              @click="confirmRoleChange"
            >
              {{ workspaces.updatingMemberId.value === userId ? 'Updating...' : 'Confirm role change' }}
            </button>
          </div>
        </div>
      </article>

      <article class="card stack member-danger-zone">
        <h2 class="section-title">Remove from workspace</h2>
        <p class="page-subtitle page-subtitle--flush">
          Removing a member revokes their access to this gym. They can be invited again later.
        </p>
        <button
          v-if="!pendingRemoval"
          type="button"
          class="secondary danger-button"
          :disabled="workspaces.removingMemberId.value === userId"
          @click="requestRemoval"
        >
          Remove from workspace
        </button>

        <div v-else class="confirmation-panel confirmation-panel--danger" role="alertdialog" aria-live="polite">
          <p class="confirmation-panel__text">
            Remove {{ memberLabel }} from {{ session.activeWorkspace.value?.name }}? This action cannot be
            undone.
          </p>
          <div class="confirmation-panel__actions">
            <button type="button" class="secondary" @click="cancelRemoval">Cancel</button>
            <button
              type="button"
              class="danger-button"
              :disabled="workspaces.removingMemberId.value === userId"
              @click="confirmRemoval"
            >
              {{ workspaces.removingMemberId.value === userId ? 'Removing...' : 'Confirm removal' }}
            </button>
          </div>
        </div>
      </article>
    </section>

    <div v-if="workspaces.workspaceError.value" class="alert alert--error" role="alert">
      {{ workspaces.workspaceError.value }}
    </div>
    <div v-if="workspaces.workspaceSuccess.value" class="alert alert--success" role="status">
      {{ workspaces.workspaceSuccess.value }}
    </div>
  </div>
</template>
