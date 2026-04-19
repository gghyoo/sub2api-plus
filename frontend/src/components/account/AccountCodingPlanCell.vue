<template>
  <!-- GLM / MiniMax: progress bars -->
  <div v-if="bars.length > 0" class="space-y-1">
    <UsageProgressBar
      v-for="bar in bars"
      :key="bar.label"
      :label="bar.label"
      :utilization="bar.utilization"
      :resets-at="bar.resetsAt"
      :color="bar.color"
    />
  </div>

  <!-- Loading -->
  <div v-else-if="loading" class="space-y-1">
    <div class="flex items-center gap-1">
      <div class="h-3 w-[32px] animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
      <div class="h-1.5 w-8 animate-pulse rounded-full bg-gray-200 dark:bg-gray-700" />
      <div class="h-3 w-[32px] animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
    </div>
    <div class="flex items-center gap-1">
      <div class="h-3 w-[32px] animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
      <div class="h-1.5 w-8 animate-pulse rounded-full bg-gray-200 dark:bg-gray-700" />
      <div class="h-3 w-[32px] animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import type { Account } from '@/types'
import { adminAPI } from '@/api/admin'
import type { GLMQuotaLimit, MiniMaxModelRemain, KimiLimit } from '@/api/admin/accounts'
import UsageProgressBar from './UsageProgressBar.vue'

const props = defineProps<{ account: Account }>()

const platform = ref<'glm' | 'minimax' | 'kimi' | null>(null)
const loading = ref(false)
const glmLimits = ref<GLMQuotaLimit[]>([])
const minimaxModel = ref<MiniMaxModelRemain | null>(null)
const kimiLimits = ref<KimiLimit[]>([])

const unitLabels: Record<number, string> = { 1: 'Y', 2: 'M', 3: 'h', 4: 'D', 5: 'M', 6: 'W' }

interface Bar {
  label: string
  utilization: number
  resetsAt: string | null
  color: 'indigo' | 'emerald' | 'purple' | 'amber'
}

const colors: Bar['color'][] = ['indigo', 'emerald', 'purple', 'amber']

const bars = computed<Bar[]>(() => {
  if (platform.value === 'glm') {
    return glmLimits.value.map((limit, i) => ({
      label: `${limit.number}${unitLabels[limit.unit] || '?'}`,
      utilization: limit.percentage ?? 0,
      resetsAt: limit.nextResetTime ? new Date(limit.nextResetTime).toISOString() : null,
      color: colors[i % colors.length],
    }))
  }

  if (platform.value === 'minimax' && minimaxModel.value) {
    const m = minimaxModel.value
    const result: Bar[] = []

    // 5h interval: usage_count is REMAINING, total is TOTAL
    const total5h = m.current_interval_total_count
    const remaining5h = m.current_interval_usage_count
    const used5h = total5h - remaining5h
    result.push({
      label: '5h',
      utilization: total5h > 0 ? (used5h / total5h) * 100 : 0,
      resetsAt: m.end_time ? new Date(m.end_time).toISOString() : null,
      color: 'indigo',
    })

    // Weekly
    const totalW = m.current_weekly_total_count
    const remainingW = m.current_weekly_usage_count
    const usedW = totalW - remainingW
    result.push({
      label: '7d',
      utilization: totalW > 0 ? (usedW / totalW) * 100 : 0,
      resetsAt: m.weekly_end_time ? new Date(m.weekly_end_time).toISOString() : null,
      color: 'emerald',
    })

    return result
  }

  if (platform.value === 'kimi' && kimiLimits.value.length > 0) {
    return kimiLimits.value.map((limit, i) => {
      const total = parseInt(limit.detail.limit, 10) || 0
      const remaining = parseInt(limit.detail.remaining, 10) || 0
      const used = total - remaining
      const duration = limit.window.duration
      const unit = limit.window.timeUnit === 'TIME_UNIT_MINUTE' ? 'm' : 's'
      return {
        label: `${duration}${unit}`,
        utilization: total > 0 ? (used / total) * 100 : 0,
        resetsAt: limit.detail.resetTime || null,
        color: colors[i % colors.length],
      }
    })
  }

  return []
})

onMounted(async () => {
  loading.value = true
  try {
    const result = await adminAPI.accounts.getCodingPlanUsage(props.account.id)
    platform.value = result.platform

    if (result.platform === 'glm' && result.glm) {
      glmLimits.value = result.glm.quota_limits ?? []
    } else if (result.platform === 'minimax' && result.minimax) {
      minimaxModel.value = (result.minimax.models ?? []).find(m => /^MiniMax-M/i.test(m.model_name)) ?? null
    } else if (result.platform === 'kimi' && result.kimi) {
      kimiLimits.value = result.kimi.limits ?? []
    }
  } catch {
    // Not a coding plan account or API error
  } finally {
    loading.value = false
  }
})
</script>
