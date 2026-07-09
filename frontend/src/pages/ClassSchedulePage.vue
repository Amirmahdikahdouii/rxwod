<script setup lang="ts">
import { useSession } from '@/features/auth/composables/useSession'
import ClassSchedule from '@/features/classbooking/components/ClassSchedule.vue'
import { ROLE_LABELS } from '@/features/workspace/model/workspaceTypes'

const session = useSession()
</script>

<template>
  <div class="container stack-lg">
    <header class="page-header">
      <p class="eyebrow">Class Schedule</p>
      <h1 class="page-title">Schedule</h1>
      <p class="page-subtitle">
        Book classes at {{ session.activeWorkspace.value?.name ?? 'your gym' }}.
        <span v-if="session.activeWorkspaceRole.value" class="inline-role">
          {{ ROLE_LABELS[session.activeWorkspaceRole.value] }}
        </span>
      </p>
    </header>

    <ClassSchedule />
  </div>
</template>
