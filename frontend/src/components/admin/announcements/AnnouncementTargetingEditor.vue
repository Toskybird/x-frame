<template>
  <div class="rounded-2xl border border-gray-200 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-800/50">
    <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <div class="text-sm font-medium text-gray-900 dark:text-white">
          {{ t('admin.announcements.form.targetingMode') }}
        </div>
        <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">
          {{ t('admin.announcements.form.targetingAll') }}
        </div>
      </div>

      <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
        <input
          type="radio"
          name="announcement-targeting-mode"
          value="all"
          checked
          class="h-4 w-4"
          @change="setAllUsers"
        />
        {{ t('admin.announcements.form.targetingAll') }}
      </label>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import type { AnnouncementTargeting } from '@/types'

const { t } = useI18n()

const props = defineProps<{
  modelValue: AnnouncementTargeting
  /** Kept for AnnouncementsView compatibility; unused after example domains removed */
  groups?: unknown[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: AnnouncementTargeting): void
}>()

function setAllUsers() {
  emit('update:modelValue', { any_of: [] })
}

onMounted(() => {
  // Force all-users targeting when example subscription/group targeting was removed
  if ((props.modelValue?.any_of?.length ?? 0) > 0) {
    setAllUsers()
  }
})
</script>
