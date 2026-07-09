import { createRouter, createWebHistory } from 'vue-router'
import { useSession } from '@/features/auth/composables/useSession'
import LoginPage from '@/pages/LoginPage.vue'
import ForgotPasswordPage from '@/pages/ForgotPasswordPage.vue'
import ResetPasswordPage from '@/pages/ResetPasswordPage.vue'
import InviteAcceptPage from '@/pages/InviteAcceptPage.vue'
import ProfilePage from '@/pages/ProfilePage.vue'
import RegisterPage from '@/pages/RegisterPage.vue'
import WorkspaceCreatePage from '@/pages/WorkspaceCreatePage.vue'
import WorkspaceMemberPage from '@/pages/WorkspaceMemberPage.vue'
import WorkspacePage from '@/pages/WorkspacePage.vue'
import WODCreatePage from '@/pages/WODCreatePage.vue'
import WODListPage from '@/pages/WODListPage.vue'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: LoginPage, meta: { guestOnly: true } },
    { path: '/register', component: RegisterPage, meta: { guestOnly: true } },
    { path: '/forgot-password', component: ForgotPasswordPage, meta: { guestOnly: true } },
    { path: '/reset-password/:token', component: ResetPasswordPage },
    { path: '/invite/:token', component: InviteAcceptPage },
    { path: '/', component: WODCreatePage, meta: { requiresAuth: true, requiresWorkspace: true } },
    { path: '/wods', component: WODListPage, meta: { requiresAuth: true, requiresWorkspace: true } },
    {
      path: '/wods/:id/edit',
      component: WODCreatePage,
      meta: { requiresAuth: true, requiresWorkspace: true },
    },
    { path: '/profile', component: ProfilePage, meta: { requiresAuth: true } },
    { path: '/workspace', component: WorkspacePage, meta: { requiresAuth: true, requiresWorkspace: true } },
    {
      path: '/workspace/members/:userId',
      component: WorkspaceMemberPage,
      meta: { requiresAuth: true, requiresWorkspace: true },
    },
    { path: '/workspace/new', component: WorkspaceCreatePage, meta: { requiresAuth: true } },
  ],
})

router.beforeEach(async (to) => {
  const session = useSession()

  if (!session.sessionReady.value && session.isAuthenticated.value) {
    await session.loadMe()
  }

  if (to.meta.requiresAuth && !session.isAuthenticated.value) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }

  if (to.meta.guestOnly && session.isAuthenticated.value) {
    return session.activeWorkspace.value ? '/' : '/workspace/new'
  }

  if (to.meta.requiresWorkspace && session.workspaces.value.length === 0) {
    return '/workspace/new'
  }

  if (to.meta.requiresWorkspace && !session.activeWorkspace.value) {
    return '/profile'
  }

  return true
})
