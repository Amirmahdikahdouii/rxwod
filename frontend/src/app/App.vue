<script setup lang="ts">
import AppHeader from '@/shared/components/AppHeader.vue'
import EmailVerificationNotice from '@/features/auth/components/EmailVerificationNotice.vue'
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="app-shell">
    <AppHeader />
    <div class="app-content">
      <EmailVerificationNotice />
      <main class="app-main">
        <RouterView />
      </main>
    </div>
  </div>
</template>
