<template>
  <div class="flex items-center gap-3">
    <!-- 验证码图片 -->
    <div
      class="flex-shrink-0 overflow-hidden rounded-lg border border-gray-300 dark:border-dark-600"
      style="width: 140px; height: 50px"
    >
      <img
        v-if="imageData"
        :src="imageData"
        alt="captcha"
        class="h-full w-full cursor-pointer select-none"
        :class="{ 'opacity-50': loading }"
        @click="refresh"
        :title="t('auth.captchaRefresh')"
      />
      <div
        v-else
        class="flex h-full w-full items-center justify-center bg-gray-100 dark:bg-dark-700"
      >
        <svg
          class="h-5 w-5 animate-spin text-gray-400"
          fill="none"
          viewBox="0 0 24 24"
        >
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
        </svg>
      </div>
    </div>

    <!-- 输入框 + 刷新按钮 -->
    <div class="flex flex-1 items-center gap-2">
      <input
        ref="inputRef"
        v-model="inputCode"
        type="text"
        maxlength="6"
        autocomplete="off"
        class="input flex-1 tracking-widest"
        :placeholder="t('auth.captchaPlaceholder')"
        @input="onInput"
        @keydown.enter.prevent="onInput"
      />
      <button
        type="button"
        class="flex-shrink-0 rounded-lg p-2 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600 dark:hover:bg-dark-700 dark:hover:text-dark-300"
        :title="t('auth.captchaRefresh')"
        @click="refresh"
      >
        <svg class="h-5 w-5" :class="{ 'animate-spin': loading }" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import apiClient from '@/api/client'

const { t } = useI18n()

const emit = defineEmits<{
  verify: [payload: { id: string; code: string }]
  error: []
}>()

const captchaId = ref('')
const imageData = ref('')
const inputCode = ref('')
const loading = ref(false)
const inputRef = ref<HTMLInputElement | null>(null)

async function refresh() {
  loading.value = true
  inputCode.value = ''
  try {
    const { data } = await apiClient.get<{ captcha_id: string; image: string }>('/auth/captcha')
    captchaId.value = data.captcha_id
    imageData.value = data.image
  } catch {
    emit('error')
  } finally {
    loading.value = false
  }
}

function onInput() {
  const code = inputCode.value.trim()
  if (code.length >= 4) {
    emit('verify', { id: captchaId.value, code })
  }
}

function reset() {
  inputCode.value = ''
  refresh()
}

defineExpose({ reset })

onMounted(() => {
  refresh()
})
</script>
