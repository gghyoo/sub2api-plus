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
              <th class="px-4 py-3 text-right font-semibold text-gray-700 dark:text-gray-300">{{ t('home.modelsTable.cacheReadPrice') }}</th>
              <th class="px-4 py-3 text-right font-semibold text-gray-700 dark:text-gray-300">{{ t('home.modelsTable.cacheCreationPrice') }}</th>
              <th class="px-4 py-3 text-left font-semibold text-gray-700 dark:text-gray-300">{{ t('home.modelsTable.platform') }}</th>
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
                <span
                  class="cursor-pointer border-b border-dashed border-gray-400 font-mono text-xs font-medium text-gray-900 transition-colors hover:border-primary-500 hover:text-primary-600 dark:text-white dark:hover:border-primary-400 dark:hover:text-primary-400"
                  :title="copiedModel === model.id ? t('home.toolGuide.copied') : t('home.toolGuide.copy')"
                  @click="copyModelId(model.id)"
                >
                  {{ copiedModel === model.id ? '✓' : '' }}{{ model.id }}
                </span>
              </td>
              <td class="px-4 py-3 text-right">
                <span v-if="model.input_price_per_mtok" class="text-gray-700 dark:text-gray-300">${{ model.input_price_per_mtok.toFixed(2) }}</span>
                <span v-else class="text-gray-400 dark:text-dark-500">-</span>
              </td>
              <td class="px-4 py-3 text-right">
                <span v-if="model.output_price_per_mtok" class="text-gray-700 dark:text-gray-300">${{ model.output_price_per_mtok.toFixed(2) }}</span>
                <span v-else class="text-gray-400 dark:text-dark-500">-</span>
              </td>
              <td class="px-4 py-3 text-right">
                <span v-if="model.cache_read_price_per_mtok" class="text-gray-700 dark:text-gray-300">${{ model.cache_read_price_per_mtok.toFixed(2) }}</span>
                <span v-else class="text-gray-400 dark:text-dark-500">-</span>
              </td>
              <td class="px-4 py-3 text-right">
                <span v-if="model.cache_creation_price_per_mtok" class="text-gray-700 dark:text-gray-300">${{ model.cache_creation_price_per_mtok.toFixed(2) }}</span>
                <span v-else class="text-gray-400 dark:text-dark-500">-</span>
              </td>
              <td class="px-4 py-3">
                <div class="flex flex-wrap gap-1.5">
                  <span
                    v-for="p in getPlatforms(model.endpoints)"
                    :key="p.label"
                    class="inline-flex items-center rounded-md px-2 py-0.5 text-xs font-medium"
                    :class="p.color"
                  >
                    {{ p.label }}
                  </span>
                  <span v-if="getPlatforms(model.endpoints).length === 0" class="text-gray-400 dark:text-dark-500">-</span>
                </div>
              </td>
              <td class="px-4 py-3">
                <div class="flex flex-wrap gap-1.5">
                  <span
                    v-for="ep in model.endpoints"
                    :key="ep"
                    :class="epColor(ep)"
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
      <!-- Models Guide Tip -->
      <div
        v-if="renderedGuide"
        class="border-t border-green-100 bg-green-50/50 px-4 py-3 dark:border-green-800/50 dark:bg-green-900/20"
      >
        <div class="flex items-start gap-3 p-3">
          <svg class="mt-0.5 h-5 w-5 flex-shrink-0 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div
            class="markdown-body prose prose-sm max-w-none text-sm text-green-700 dark:text-green-300 dark:prose-invert"
            v-html="renderedGuide"
          ></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

const { t } = useI18n()
const appStore = useAppStore()

marked.setOptions({
  breaks: true,
  gfm: true,
})

const modelsGuide = computed(() => appStore.cachedPublicSettings?.models_guide || '')

const renderedGuide = computed(() => {
  if (!modelsGuide.value) return ''
  const html = marked.parse(modelsGuide.value) as string
  return DOMPurify.sanitize(html)
})

interface ModelInfo {
  id: string
  type: string
  created_at: string
  input_price_per_mtok?: number
  output_price_per_mtok?: number
  cache_read_price_per_mtok?: number
  cache_creation_price_per_mtok?: number
  endpoints?: string[]
}

const models = ref<ModelInfo[]>([])
const loading = ref(true)
const error = ref('')
const copiedModel = ref<string | null>(null)

function getPlatforms(endpoints?: string[]): { label: string; color: string }[] {
  if (!endpoints || endpoints.length === 0) return []
  const platforms: { label: string; color: string }[] = []
  const hasAnthropic = endpoints.some(ep => ep === '/v1/messages' && !ep.startsWith('/antigravity'))
  const hasOpenAI = endpoints.some(ep => ep === '/v1/responses' || ep === '/v1/chat/completions')
  const hasGemini = endpoints.some(ep => ep.startsWith('/v1beta'))
  const hasAntigravity = endpoints.some(ep => ep.startsWith('/antigravity'))
  if (hasAnthropic) platforms.push({ label: 'Anthropic', color: 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400' })
  if (hasOpenAI) platforms.push({ label: 'OpenAI', color: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400' })
  if (hasGemini) platforms.push({ label: 'Gemini', color: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400' })
  if (hasAntigravity) platforms.push({ label: 'Antigravity', color: 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400' })
  return platforms
}

function epColor(ep: string): string {
  const base = 'inline-flex items-center rounded-md px-2 py-0.5 font-mono text-xs'
  if (ep.startsWith('/antigravity')) return `${base} bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400`
  if (ep.startsWith('/v1beta')) return `${base} bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400`
  if (ep === '/v1/messages') return `${base} bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400`
  if (ep === '/v1/chat/completions') return `${base} bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400`
  if (ep === '/v1/responses') return `${base} bg-indigo-100 text-indigo-700 dark:bg-indigo-900/30 dark:text-indigo-400`
  return `${base} bg-gray-100 text-gray-700 dark:bg-gray-900/30 dark:text-gray-400`
}

async function copyModelId(modelId: string) {
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(modelId)
    } else {
      const textarea = document.createElement('textarea')
      textarea.value = modelId
      textarea.style.position = 'fixed'
      textarea.style.opacity = '0'
      document.body.appendChild(textarea)
      textarea.select()
      document.execCommand('copy')
      document.body.removeChild(textarea)
    }
    copiedModel.value = modelId
    setTimeout(() => { copiedModel.value = null }, 2000)
  } catch {}
}

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
