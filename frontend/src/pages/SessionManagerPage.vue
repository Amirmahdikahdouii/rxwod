<script setup lang="ts">
import { useSession } from '@/features/auth/composables/useSession'
import SessionManager from '@/features/classbooking/components/SessionManager.vue'
import { ROLE_LABELS } from '@/features/workspace/model/workspaceTypes'

const session = useSession()
</script>

<template>
  <div class="container stack-lg">
    <header class="page-header">
      <p class="eyebrow">Session Management</p>
      <h1 class="page-title">Manage Schedule</h1>
      <p class="page-subtitle">
        Create class sessions and manage rosters for {{ session.activeWorkspace.value?.name ?? 'your gym' }}.
        <span v-if="session.activeWorkspaceRole.value" class="inline-role">
          {{ ROLE_LABELS[session.activeWorkspaceRole.value] }}
        </span>
      </p>
    </header>

    <SessionManager />
  </div>
</template>
