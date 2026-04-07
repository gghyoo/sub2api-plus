<template>
  <div class="mt-12">
    <!-- Section Title -->
    <div class="mb-6 text-center">
      <h2 class="mb-2 text-2xl font-bold text-gray-900 dark:text-white">
        {{ t('home.modelsTable.title') }}
      </h2>
      <p class="text-sm text-gray-600 dark:text-dark-400">
        {{ t('home.modelsTable.subtitle') }}
      </p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
      <div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-200 border-t-primary-600"></div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="mx-auto max-w-4xl rounded-lg border border-red-200 bg-red-50 p-4 text-center text-sm text-red-600 dark:border-red-800 dark:bg-red-900/20 dark:text-red-400">
      {{ error }}
    </div>

    <!-- Table -->
    <div v-else-if="models.length > 0" class="mx-auto max-w-6xl overflow-hidden rounded-2xl border border-gray-200/50 bg-white/60 backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-800/60">
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-gray-200 bg-gray-50/80 dark:border-dark-700 dark:bg-dark-800/80">
              <th class="px-4 py-3 text-left font-semibold text-gray-700 dark:text-gray-300">{{ t('home.modelsTable.model') }}</th>
              <th class="px-4 py-3 text-right font-semibold text-gray-700 dark:text-gray-300">{{ t('home.modelsTable.inputPrice') }}</th>
              <th class="px-4 py-3 text-right font-semibold text-gray-700 dark:text-gray-300">{{ t('home.modelsTable.outputPrice') }}</th>
              <th class="px-4 py-3 text-left font-semibold text-gray-700 dark:text-gray-300">{{ t('home.modelsTable.endpoints') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(model, index) in models"
              :key="model.id"
              :class="[
                'border-b border-gray-100 transition-colors last:border-b-0 dark:border-dark-700/50',
                index % 2 === 0 ? 'bg-white/40 dark:bg-dark-800/30' : 'bg-gray-50/40 dark:bg-dark-900/20'
              ]"
            >
              <td class="px-4 py-3">
                <span class="font-mono text-xs font-medium text-gray-900 dark:text-white">{{ model.id }}</span>
              </td>
              <td class="px-4 py-3 text-right">
                <span v-if="model.input_price_per_mtok" class="text-gray-700 dark:text-gray-300">${{ model.input_price_per_mtok.toFixed(2) }}</span>
                <span v-else class="text-gray-400 dark:text-dark-500">-</span>
              </td>
              <td class="px-4 py-3 text-right">
                <span v-if="model.output_price_per_mtok" class="text-gray-700 dark:text-gray-300">${{ model.output_price_per_mtok.toFixed(2) }}</span>
                <span v-else class="text-gray-400 dark:text-dark-500">-</span>
              </td>
              <td class="px-4 py-3">
                <div class="flex flex-wrap gap-1.5">
                  <span
                    v-for="ep in model.endpoints"
                    :key="ep"
                    class="inline-flex items-center rounded-md bg-primary-50 px-2 py-0.5 font-mono text-xs text-primary-700 dark:bg-primary-900/30 dark:text-primary-400"
                  >
                    {{ ep }}
                  </span>
                  <span v-if="!model.endpoints || model.endpoints.length === 0" class="text-gray-400 dark:text-dark-500">-</span>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="border-t border-gray-200 bg-gray-50/50 px-4 py-2 text-center text-xs text-gray-500 dark:border-dark-700 dark:bg-dark-800/50 dark:text-dark-400">
        {{ t('home.modelsTable.priceUnit') }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

interface ModelInfo {
  id: string
  type: string
  created_at: string
  input_price_per_mtok?: number
  output_price_per_mtok?: number
  endpoints?: string[]
}

const models = ref<ModelInfo[]>([])
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    const baseUrl = window.location.origin
    const resp = await fetch(`${baseUrl}/v1/models`)
    if (!resp.ok) throw new Error(`HTTP ${resp.status}`)
    const data = await resp.json()
    if (Array.isArray(data?.data)) {
      models.value = data.data
    }
  } catch (e: any) {
    error.value = t('home.modelsTable.loadError')
  } finally {
    loading.value = false
  }
})
</script>
