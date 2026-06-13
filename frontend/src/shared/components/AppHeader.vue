<script setup lang="ts">
import { RouterLink, useRouter } from 'vue-router'
import { useSession } from '@/features/auth/composables/useSession'
import { ROLE_LABELS } from '@/features/workspace/model/workspaceTypes'
import ThemeToggle from '@/shared/components/ThemeToggle.vue'

const router = useRouter()
const session = useSession()

async function logout() {
  session.logout()
  await router.push('/login')
}

async function updateWorkspace(value: string) {
  session.setActiveWorkspace(value)
  await router.push('/')
}
</script>

<template>
  <header class="app-header">
    <div class="app-header__inner">
      <RouterLink to="/" class="app-header__brand">RXWOD</RouterLink>

      <nav class="app-header__nav" aria-label="Main navigation">
        <template v-if="session.isAuthenticated.value">
          <RouterLink to="/" class="app-header__link">Create</RouterLink>
          <RouterLink to="/wods" class="app-header__link">Saved</RouterLink>
          <RouterLink to="/workspace" class="app-header__link">Workspace</RouterLink>
          <RouterLink to="/profile" class="app-header__link">Profile</RouterLink>
        </template>
        <template v-else>
          <RouterLink to="/login" class="app-header__link">Login</RouterLink>
          <RouterLink to="/register" class="app-header__link">Register</RouterLink>
        </template>
      </nav>

      <div class="app-header__actions">
        <div v-if="session.isAuthenticated.value" class="workspace-switcher">
          <select
            v-if="session.workspaces.value.length > 0"
            class="workspace-switcher__select"
            :value="session.activeWorkspaceId.value ?? ''"
            aria-label="Active workspace"
            @change="updateWorkspace(($event.target as HTMLSelectElement).value)"
          >
            <option v-for="workspace in session.workspaces.value" :key="workspace.id" :value="workspace.id">
              {{ workspace.name }}
            </option>
          </select>
          <RouterLink v-else to="/workspace/new" class="app-header__link app-header__link--accent">
            Create Gym
          </RouterLink>
          <span v-if="session.activeWorkspaceRole.value" class="role-pill">
            {{ ROLE_LABELS[session.activeWorkspaceRole.value] }}
          </span>
        </div>
        <ThemeToggle />
        <button v-if="session.isAuthenticated.value" type="button" class="secondary app-header__button" @click="logout">
          Logout
        </button>
      </div>
    </div>
  </header>
</template>
