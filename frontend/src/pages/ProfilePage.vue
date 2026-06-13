<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useSession } from '@/features/auth/composables/useSession'
import { ROLE_LABELS } from '@/features/workspace/model/workspaceTypes'
import BaseSelect from '@/shared/components/BaseSelect.vue'

const router = useRouter()
const session = useSession()

const workspaceOptions = computed(() =>
  session.workspaces.value.map((workspace) => ({
    value: workspace.id,
    label: `${workspace.name} (${ROLE_LABELS[workspace.role]})`,
  })),
)

async function updateWorkspace(workspaceID: string) {
  session.setActiveWorkspace(workspaceID)
  await router.push('/')
}

async function logout() {
  session.logout()
  await router.push('/login')
}
</script>

<template>
  <div class="container stack-lg">
    <header class="page-header">
      <p class="eyebrow">Profile</p>
      <h1 class="page-title">Your account</h1>
      <p class="page-subtitle">Manage your session and choose the gym workspace you are working in.</p>
    </header>

    <section class="dashboard-grid">
      <article class="card stack">
        <h2 class="section-title">Account</h2>
        <div class="profile-row">
          <span class="profile-avatar">{{ session.currentUser.value?.displayName?.charAt(0) || 'U' }}</span>
          <div>
            <strong>{{ session.currentUser.value?.displayName || 'RXWOD User' }}</strong>
            <p class="page-subtitle page-subtitle--flush">{{ session.currentUser.value?.email }}</p>
          </div>
        </div>
        <button type="button" class="secondary" @click="logout">Logout</button>
      </article>

      <article class="card stack">
        <h2 class="section-title">Active Workspace</h2>
        <div v-if="session.workspaces.value.length > 0" class="stack">
          <BaseSelect
            :model-value="session.activeWorkspaceId.value ?? ''"
            label="Gym workspace"
            :options="workspaceOptions"
            @update:model-value="updateWorkspace"
          />
          <div v-if="session.activeWorkspace.value" class="workspace-summary">
            <span class="role-pill">{{ ROLE_LABELS[session.activeWorkspace.value.role] }}</span>
            <strong>{{ session.activeWorkspace.value.name }}</strong>
          </div>
        </div>
        <div v-else class="empty-state empty-state--compact">
          <h3 class="empty-state__title">No gym yet</h3>
          <p class="empty-state__text">Create a workspace to start saving programs.</p>
          <RouterLink to="/workspace/new" class="empty-state__link">Create Gym</RouterLink>
        </div>
      </article>
    </section>
  </div>
</template>
