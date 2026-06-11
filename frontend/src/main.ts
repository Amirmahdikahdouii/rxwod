import { initTheme } from '@/shared/composables/themeUtils'

initTheme()

import { createApp } from 'vue'
import App from '@/app/App.vue'
import { router } from '@/app/router'
import '@/style.css'

createApp(App).use(router).mount('#app')
