<template>
  <AppLayout>
    <div class="space-y-6">
      <div class="card p-6 sm:p-8">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
          {{ t('dashboard.title') }}
        </h2>
        <p class="mt-2 text-sm text-gray-600 dark:text-gray-300">
          {{ t('dashboard.welcomeMessage') }}
          <span v-if="displayName" class="font-medium text-gray-900 dark:text-white">
            {{ displayName }}
          </span>
        </p>
        <div class="mt-6 flex flex-wrap gap-3">
          <router-link to="/profile" class="btn btn-primary">
            {{ t('nav.profile') }}
          </router-link>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import AppLayout from '@/components/layout/AppLayout.vue'

const { t } = useI18n()
const authStore = useAuthStore()

const displayName = computed(() => {
  const user = authStore.user
  if (!user) return ''
  return user.username || user.email?.split('@')[0] || ''
})
</script>
