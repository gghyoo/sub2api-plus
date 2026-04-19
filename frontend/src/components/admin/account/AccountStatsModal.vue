<template>
  <BaseDialog
    :show="show"
    :title="t('admin.accounts.usageStatistics')"
    width="extra-wide"
    @close="handleClose"
  >
    <div class="space-y-6">
      <!-- Account Info Header -->
      <div
        v-if="account"
        class="flex items-center justify-between rounded-xl border border-primary-200 bg-gradient-to-r from-primary-50 to-primary-100 p-3 dark:border-primary-700/50 dark:from-primary-900/20 dark:to-primary-800/20"
      >
        <div class="flex items-center gap-3">
          <div
            class="flex h-10 w-10 items-center justify-center rounded-lg bg-gradient-to-br from-primary-500 to-primary-600"
          >
            <Icon name="chartBar" size="md" class="text-white" />
          </div>
          <div>
            <div class="font-semibold text-gray-900 dark:text-gray-100">{{ account.name }}</div>
            <div class="text-xs text-gray-500 dark:text-gray-400">
              {{ t('admin.accounts.last30DaysUsage') }}
            </div>
          </div>
        </div>
        <span
          :class="[
            'rounded-full px-2.5 py-1 text-xs font-semibold',
            account.status === 'active'
              ? 'bg-green-100 text-green-700 dark:bg-green-500/20 dark:text-green-400'
              : 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-400'
          ]"
        >
          {{ account.status }}
        </span>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <LoadingSpinner />
      </div>

      <template v-else-if="stats">
        <!-- Row 1: Main Stats Cards -->
        <div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
          <!-- 30-Day Total Cost -->
          <div
            class="card border-emerald-200 bg-gradient-to-br from-emerald-50 to-white p-4 dark:border-emerald-800/30 dark:from-emerald-900/10 dark:to-dark-700"
          >
            <div class="mb-2 flex items-center justify-between">
              <span class="text-xs font-medium text-gray-500 dark:text-gray-400">{{
                t('admin.accounts.stats.totalCost')
              }}</span>
              <div class="rounded-lg bg-emerald-100 p-1.5 dark:bg-emerald-900/30">
                <Icon name="dollar" size="sm" class="text-emerald-600 dark:text-emerald-400" />
              </div>
            </div>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">
              ${{ formatCost(stats.summary.total_cost) }}
            </p>
            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
              {{ t('admin.accounts.stats.accumulatedCost') }}
              <span class="text-gray-400 dark:text-gray-500">
                ({{ t('usage.userBilled') }}: ${{ formatCost(stats.summary.total_user_cost) }} ·
                {{ t('admin.accounts.stats.standardCost') }}: ${{
                  formatCost(stats.summary.total_standard_cost)
                }})
              </span>
            </p>
          </div>

          <!-- 30-Day Total Requests -->
          <div
            class="card border-blue-200 bg-gradient-to-br from-blue-50 to-white p-4 dark:border-blue-800/30 dark:from-blue-900/10 dark:to-dark-700"
          >
            <div class="mb-2 flex items-center justify-between">
              <span class="text-xs font-medium text-gray-500 dark:text-gray-400">{{
                t('admin.accounts.stats.totalRequests')
              }}</span>
              <div class="rounded-lg bg-blue-100 p-1.5 dark:bg-blue-900/30">
                <Icon name="bolt" size="sm" class="text-blue-600 dark:text-blue-400" />
              </div>
            </div>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ formatNumber(stats.summary.total_requests) }}
            </p>
            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
              {{ t('admin.accounts.stats.totalCalls') }}
            </p>
          </div>

          <!-- Daily Average Cost -->
          <div
            class="card border-amber-200 bg-gradient-to-br from-amber-50 to-white p-4 dark:border-amber-800/30 dark:from-amber-900/10 dark:to-dark-700"
          >
            <div class="mb-2 flex items-center justify-between">
              <span class="text-xs font-medium text-gray-500 dark:text-gray-400">{{
                t('admin.accounts.stats.avgDailyCost')
              }}</span>
              <div class="rounded-lg bg-amber-100 p-1.5 dark:bg-amber-900/30">
                <Icon
                  name="calculator"
                  size="sm"
                  class="text-amber-600 dark:text-amber-400"
                />
              </div>
            </div>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">
              ${{ formatCost(stats.summary.avg_daily_cost) }}
            </p>
             <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
              {{
                t('admin.accounts.stats.basedOnActualDays', {
                  days: stats.summary.actual_days_used
                })
              }}
              <span class="text-gray-400 dark:text-gray-500">
                ({{ t('usage.userBilled') }}: ${{ formatCost(stats.summary.avg_daily_user_cost) }})
              </span>
            </p>
          </div>

          <!-- Daily Average Requests -->
          <div
            class="card border-purple-200 bg-gradient-to-br from-purple-50 to-white p-4 dark:border-purple-800/30 dark:from-purple-900/10 dark:to-dark-700"
          >
            <div class="mb-2 flex items-center justify-between">
              <span class="text-xs font-medium text-gray-500 dark:text-gray-400">{{
                t('admin.accounts.stats.avgDailyRequests')
              }}</span>
              <div class="rounded-lg bg-purple-100 p-1.5 dark:bg-purple-900/30">
                <svg
                  class="h-4 w-4 text-purple-600 dark:text-purple-400"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z"
                  />
                </svg>
              </div>
            </div>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ formatNumber(Math.round(stats.summary.avg_daily_requests)) }}
            </p>
            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
              {{ t('admin.accounts.stats.avgDailyUsage') }}
            </p>
          </div>
        </div>

        <!-- Row 2: Today, Highest Cost, Highest Requests -->
        <div class="grid grid-cols-1 gap-4 lg:grid-cols-3">
          <!-- Today Overview -->
          <div class="card p-4">
            <div class="mb-3 flex items-center gap-2">
              <div class="rounded-lg bg-cyan-100 p-1.5 dark:bg-cyan-900/30">
                <Icon name="clock" size="sm" class="text-cyan-600 dark:text-cyan-400" />
              </div>
              <span class="text-sm font-semibold text-gray-900 dark:text-white">{{
                t('admin.accounts.stats.todayOverview')
              }}</span>
            </div>
            <div class="space-y-2">
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('usage.accountBilled') }}</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white"
                  >${{ formatCost(stats.summary.today?.cost || 0) }}</span
                >
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('usage.userBilled') }}</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white"
                  >${{ formatCost(stats.summary.today?.user_cost || 0) }}</span
                >
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{
                  t('admin.accounts.stats.requests')
                }}</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white">{{
                  formatNumber(stats.summary.today?.requests || 0)
                }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{
                  t('admin.accounts.stats.tokens')
                }}</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white">{{
                  formatTokens(stats.summary.today?.tokens || 0)
                }}</span>
              </div>
            </div>
          </div>

          <!-- Highest Cost Day -->
          <div class="card p-4">
            <div class="mb-3 flex items-center gap-2">
              <div class="rounded-lg bg-orange-100 p-1.5 dark:bg-orange-900/30">
                <Icon name="fire" size="sm" class="text-orange-600 dark:text-orange-400" />
              </div>
              <span class="text-sm font-semibold text-gray-900 dark:text-white">{{
                t('admin.accounts.stats.highestCostDay')
              }}</span>
            </div>
            <div class="space-y-2">
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{
                  t('admin.accounts.stats.date')
                }}</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white">{{
                  stats.summary.highest_cost_day?.label || '-'
                }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('usage.accountBilled') }}</span>
                <span class="text-sm font-semibold text-orange-600 dark:text-orange-400"
                  >${{ formatCost(stats.summary.highest_cost_day?.cost || 0) }}</span
                >
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('usage.userBilled') }}</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white"
                  >${{ formatCost(stats.summary.highest_cost_day?.user_cost || 0) }}</span
                >
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{
                  t('admin.accounts.stats.requests')
                }}</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white">{{
                  formatNumber(stats.summary.highest_cost_day?.requests || 0)
                }}</span>
              </div>
            </div>
          </div>

          <!-- Highest Request Day -->
          <div class="card p-4">
            <div class="mb-3 flex items-center gap-2">
              <div class="rounded-lg bg-indigo-100 p-1.5 dark:bg-indigo-900/30">
                <Icon
                  name="trendingUp"
                  size="sm"
                  class="text-indigo-600 dark:text-indigo-400"
                />
              </div>
              <span class="text-sm font-semibold text-gray-900 dark:text-white">{{
                t('admin.accounts.stats.highestRequestDay')
              }}</span>
            </div>
            <div class="space-y-2">
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{
                  t('admin.accounts.stats.date')
                }}</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white">{{
                  stats.summary.highest_request_day?.label || '-'
                }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{
                  t('admin.accounts.stats.requests')
                }}</span>
                <span class="text-sm font-semibold text-indigo-600 dark:text-indigo-400">{{
                  formatNumber(stats.summary.highest_request_day?.requests || 0)
                }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('usage.accountBilled') }}</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white"
                  >${{ formatCost(stats.summary.highest_request_day?.cost || 0) }}</span
                >
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('usage.userBilled') }}</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white"
                  >${{ formatCost(stats.summary.highest_request_day?.user_cost || 0) }}</span
                >
              </div>
            </div>
          </div>
        </div>

        <!-- Row 3: Token Stats -->
        <div class="grid grid-cols-1 gap-4 lg:grid-cols-3">
          <!-- Accumulated Tokens -->
          <div class="card p-4">
            <div class="mb-3 flex items-center gap-2">
              <div class="rounded-lg bg-teal-100 p-1.5 dark:bg-teal-900/30">
                <Icon name="cube" size="sm" class="text-teal-600 dark:text-teal-400" />
              </div>
              <span class="text-sm font-semibold text-gray-900 dark:text-white">{{
                t('admin.accounts.stats.accumulatedTokens')
              }}</span>
            </div>
            <div class="space-y-2">
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{
                  t('admin.accounts.stats.totalTokens')
                }}</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white">{{
                  formatTokens(stats.summary.total_tokens)
                }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{
                  t('admin.accounts.stats.dailyAvgTokens')
                }}</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white">{{
                  formatTokens(Math.round(stats.summary.avg_daily_tokens))
                }}</span>
              </div>
            </div>
          </div>

          <!-- Performance -->
          <div class="card p-4">
            <div class="mb-3 flex items-center gap-2">
              <div class="rounded-lg bg-rose-100 p-1.5 dark:bg-rose-900/30">
                <Icon name="bolt" size="sm" class="text-rose-600 dark:text-rose-400" />
              </div>
              <span class="text-sm font-semibold text-gray-900 dark:text-white">{{
                t('admin.accounts.stats.performance')
              }}</span>
            </div>
            <div class="space-y-2">
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{
                  t('admin.accounts.stats.avgResponseTime')
                }}</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white">{{
                  formatDuration(stats.summary.avg_duration_ms)
                }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{
                  t('admin.accounts.stats.daysActive')
                }}</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white"
                  >{{ stats.summary.actual_days_used }} / {{ stats.summary.days }}</span
                >
              </div>
            </div>
          </div>

          <!-- Recent Activity -->
          <div class="card p-4">
            <div class="mb-3 flex items-center gap-2">
              <div class="rounded-lg bg-lime-100 p-1.5 dark:bg-lime-900/30">
                <Icon
                  name="clipboard"
                  size="sm"
                  class="text-lime-600 dark:text-lime-400"
                />
              </div>
              <span class="text-sm font-semibold text-gray-900 dark:text-white">{{
                t('admin.accounts.stats.recentActivity')
              }}</span>
            </div>
            <div class="space-y-2">
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{
                  t('admin.accounts.stats.todayRequests')
                }}</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white">{{
                  formatNumber(stats.summary.today?.requests || 0)
                }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{
                  t('admin.accounts.stats.todayTokens')
                }}</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white">{{
                  formatTokens(stats.summary.today?.tokens || 0)
                }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{
                  t('admin.accounts.stats.todayCost')
                }}</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white"
                  >${{ formatCost(stats.summary.today?.cost || 0) }}</span
                >
              </div>
            </div>
          </div>
        </div>

        <!-- Usage Trend Chart -->
        <div class="card p-4">
          <h3 class="mb-4 text-sm font-semibold text-gray-900 dark:text-white">
            {{ t('admin.accounts.stats.usageTrend') }}
          </h3>
          <div class="h-64">
            <Line v-if="trendChartData" :data="trendChartData" :options="lineChartOptions" />
            <div
              v-else
              class="flex h-full items-center justify-center text-sm text-gray-500 dark:text-gray-400"
            >
              {{ t('admin.dashboard.noDataAvailable') }}
            </div>
          </div>
        </div>

        <!-- Model Distribution -->
        <ModelDistributionChart :model-stats="stats.models" :loading="false" />

        <EndpointDistributionChart
          :endpoint-stats="stats.endpoints || []"
          :loading="false"
          :title="t('usage.inboundEndpoint')"
        />

        <EndpointDistributionChart
          :endpoint-stats="stats.upstream_endpoints || []"
          :loading="false"
          :title="t('usage.upstreamEndpoint')"
        />

        <!-- GLM Coding Plan Usage (shown when account has bigmodel base URL) -->
        <template v-if="(glmUsage || glmLoading) && (codingPlanPlatform === 'glm' || !codingPlanPlatform)">
          <div class="flex items-center gap-2 border-t border-gray-200 pt-4 dark:border-gray-700">
            <div class="rounded-lg bg-violet-100 p-1.5 dark:bg-violet-900/30">
              <svg class="h-4 w-4 text-violet-600 dark:text-violet-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <h3 class="text-sm font-semibold text-gray-900 dark:text-white">
              {{ t('admin.accounts.stats.glmCodingPlan') }}
            </h3>
          </div>

          <!-- Loading state -->
          <div v-if="glmLoading" class="flex items-center justify-center py-8">
            <LoadingSpinner />
          </div>

          <!-- GLM data -->
          <template v-else-if="glmUsage">
            <!-- Summary: total calls & tokens -->
            <div class="grid grid-cols-2 gap-4">
              <div class="card border-blue-200 bg-gradient-to-br from-blue-50 to-white p-4 dark:border-blue-800/30 dark:from-blue-900/10 dark:to-dark-700">
                <span class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('admin.accounts.stats.glmTotalCalls') }}</span>
                <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ formatNumber(glmUsage.total_model_calls) }}</p>
                <p class="text-xs text-gray-400 dark:text-gray-500">{{ t('admin.accounts.stats.glm24h') }}</p>
              </div>
              <div class="card border-teal-200 bg-gradient-to-br from-teal-50 to-white p-4 dark:border-teal-800/30 dark:from-teal-900/10 dark:to-dark-700">
                <span class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('admin.accounts.stats.glmTotalTokens') }}</span>
                <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ formatTokens(glmUsage.total_tokens) }}</p>
                <p class="text-xs text-gray-400 dark:text-gray-500">
                  {{ t('admin.accounts.stats.glmLevel') }}: <span class="font-medium text-teal-600 dark:text-teal-400">{{ glmUsage.quota_level || '-' }}</span>
                </p>
              </div>
            </div>

            <!-- Row: Model Usage + Tool Usage -->
            <div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
              <!-- Model Usage -->
              <div class="card p-4">
                <h4 class="mb-3 text-sm font-semibold text-gray-900 dark:text-white">
                  {{ t('admin.accounts.stats.glmModelUsage') }}
                </h4>
                <div class="overflow-x-auto">
                  <table class="w-full text-sm">
                    <thead>
                      <tr class="border-b border-gray-200 dark:border-gray-700">
                        <th class="pb-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400">Model</th>
                        <th class="pb-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400">Tokens</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="item in glmUsage.models" :key="item.modelName" class="border-b border-gray-100 last:border-0 dark:border-gray-800">
                        <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ item.modelName }}</td>
                        <td class="py-1.5 text-right tabular-nums text-gray-900 dark:text-gray-100">{{ formatTokens(item.totalTokens) }}</td>
                      </tr>
                      <tr v-if="!glmUsage.models?.length">
                        <td colspan="2" class="py-3 text-center text-xs text-gray-400">{{ t('admin.dashboard.noDataAvailable') }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              <!-- Tool Usage -->
              <div class="card p-4">
                <h4 class="mb-3 text-sm font-semibold text-gray-900 dark:text-white">
                  {{ t('admin.accounts.stats.glmToolUsage') }}
                </h4>
                <div class="overflow-x-auto">
                  <table class="w-full text-sm">
                    <thead>
                      <tr class="border-b border-gray-200 dark:border-gray-700">
                        <th class="pb-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400">Tool</th>
                        <th class="pb-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400">Count</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="item in glmUsage.tools" :key="item.toolCode" class="border-b border-gray-100 last:border-0 dark:border-gray-800">
                        <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ item.toolName }}</td>
                        <td class="py-1.5 text-right tabular-nums text-gray-900 dark:text-gray-100">{{ formatNumber(item.totalUsageCount) }}</td>
                      </tr>
                      <tr v-if="!glmUsage.tools?.length">
                        <td colspan="2" class="py-3 text-center text-xs text-gray-400">{{ t('admin.dashboard.noDataAvailable') }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>

            <!-- Quota Limits -->
            <div v-if="glmUsage.quota_limits?.length" class="card p-4">
              <h4 class="mb-3 text-sm font-semibold text-gray-900 dark:text-white">
                {{ t('admin.accounts.stats.glmQuotaLimits') }}
              </h4>
              <div class="space-y-3">
                <div v-for="(limit, idx) in glmUsage.quota_limits" :key="idx">
                  <div class="mb-1 flex items-center justify-between">
                    <span class="text-xs font-medium text-gray-700 dark:text-gray-300">{{ formatQuotaLabel(limit) }}</span>
                    <span class="text-xs font-semibold tabular-nums" :class="limit.percentage > 80 ? 'text-red-600 dark:text-red-400' : 'text-gray-900 dark:text-white'">
                      {{ limit.percentage.toFixed(1) }}%
                    </span>
                  </div>
                  <div class="h-2 w-full rounded-full bg-gray-200 dark:bg-gray-700">
                    <div
                      class="h-2 rounded-full transition-all"
                      :class="limit.percentage > 80 ? 'bg-red-500' : limit.percentage > 50 ? 'bg-amber-500' : 'bg-emerald-500'"
                      :style="{ width: Math.min(limit.percentage, 100) + '%' }"
                    />
                  </div>
                  <div v-if="limit.currentValue != null && limit.remaining != null" class="mt-0.5 text-xs tabular-nums text-gray-400 dark:text-gray-500">
                    {{ limit.currentValue.toLocaleString() }} / {{ (limit.currentValue + limit.remaining).toLocaleString() }}
                  </div>
                  <div v-if="limit.nextResetTime" class="mt-0.5 text-xs tabular-nums text-gray-400 dark:text-gray-500">
                    {{ t('admin.accounts.stats.glmResetTime') }}：{{ formatResetTime(limit.nextResetTime) }}
                    <span class="ml-2">{{ formatResetDate(limit.nextResetTime) }}</span>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </template>

        <!-- MiniMax Coding Plan Usage (shown when account has minimax base URL) -->
        <template v-if="(minimaxUsage || minimaxLoading) && codingPlanPlatform === 'minimax'">
          <div class="flex items-center gap-2 border-t border-gray-200 pt-4 dark:border-gray-700">
            <div class="rounded-lg bg-indigo-100 p-1.5 dark:bg-indigo-900/30">
              <svg class="h-4 w-4 text-indigo-600 dark:text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <h3 class="text-sm font-semibold text-gray-900 dark:text-white">
              {{ t('admin.accounts.stats.minimaxCodingPlan') }}
            </h3>
          </div>

          <!-- Loading state -->
          <div v-if="minimaxLoading" class="flex items-center justify-center py-8">
            <LoadingSpinner />
          </div>

          <!-- MiniMax data -->
          <template v-else-if="minimaxUsage">
            <div class="card p-4">
              <h4 class="mb-3 text-sm font-semibold text-gray-900 dark:text-white">
                {{ t('admin.accounts.stats.minimaxModelQuota') }}
              </h4>
              <div class="space-y-4">
                <div v-for="model in minimaxUsage" :key="model.model_name" class="space-y-2 rounded-lg border border-gray-100 p-3 dark:border-gray-800">
                  <div class="text-sm font-medium text-gray-900 dark:text-white">{{ model.model_name }}</div>

                  <!-- 5h Window -->
                  <div v-if="model.current_interval_total_count > 0">
                    <div class="mb-1 flex items-center justify-between">
                      <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.accounts.stats.minimax5hWindow') }}</span>
                      <span class="text-xs tabular-nums text-gray-500 dark:text-gray-400">
                        {{ t('admin.accounts.stats.minimaxRemaining') }}：{{ model.current_interval_usage_count.toLocaleString() }} / {{ model.current_interval_total_count.toLocaleString() }}
                      </span>
                    </div>
                    <div class="h-2 w-full rounded-full bg-gray-200 dark:bg-gray-700">
                      <div
                        class="h-2 rounded-full transition-all"
                        :class="{
                          'bg-red-500': (model.current_interval_usage_count / model.current_interval_total_count * 100) < 20,
                          'bg-amber-500': (model.current_interval_usage_count / model.current_interval_total_count * 100) >= 20 && (model.current_interval_usage_count / model.current_interval_total_count * 100) < 50,
                          'bg-emerald-500': (model.current_interval_usage_count / model.current_interval_total_count * 100) >= 50
                        }"
                        :style="{ width: Math.min((1 - model.current_interval_usage_count / model.current_interval_total_count) * 100, 100) + '%' }"
                      />
                    </div>
                    <div class="mt-0.5 text-xs tabular-nums text-gray-400 dark:text-gray-500">
                      <span v-if="model.end_time">{{ t('admin.accounts.stats.minimaxResetTime') }}：{{ formatResetDate(model.end_time) }}</span>
                    </div>
                  </div>

                  <!-- Weekly -->
                  <div v-if="model.current_weekly_total_count > 0">
                    <div class="mb-1 flex items-center justify-between">
                      <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.accounts.stats.minimaxWeekly') }}</span>
                      <span class="text-xs tabular-nums text-gray-500 dark:text-gray-400">
                        {{ t('admin.accounts.stats.minimaxRemaining') }}：{{ model.current_weekly_usage_count.toLocaleString() }} / {{ model.current_weekly_total_count.toLocaleString() }}
                      </span>
                    </div>
                    <div class="h-2 w-full rounded-full bg-gray-200 dark:bg-gray-700">
                      <div
                        class="h-2 rounded-full transition-all"
                        :class="{
                          'bg-red-500': (model.current_weekly_usage_count / model.current_weekly_total_count * 100) < 20,
                          'bg-amber-500': (model.current_weekly_usage_count / model.current_weekly_total_count * 100) >= 20 && (model.current_weekly_usage_count / model.current_weekly_total_count * 100) < 50,
                          'bg-emerald-500': (model.current_weekly_usage_count / model.current_weekly_total_count * 100) >= 50
                        }"
                        :style="{ width: Math.min((1 - model.current_weekly_usage_count / model.current_weekly_total_count) * 100, 100) + '%' }"
                      />
                    </div>
                    <div class="mt-0.5 text-xs tabular-nums text-gray-400 dark:text-gray-500">
                      <span v-if="model.weekly_end_time">{{ t('admin.accounts.stats.minimaxResetTime') }}：{{ formatResetDate(model.weekly_end_time) }}</span>
                    </div>
                  </div>
                </div>

                <div v-if="!minimaxUsage.length" class="py-3 text-center text-xs text-gray-400">
                  {{ t('admin.accounts.stats.minimaxNoModels') }}
                </div>
              </div>
            </div>
          </template>
        </template>

        <!-- Kimi Coding Plan Usage (shown when account has kimi base URL) -->
        <template v-if="(kimiUsage || kimiLoading) && codingPlanPlatform === 'kimi'">
          <div class="flex items-center gap-2 border-t border-gray-200 pt-4 dark:border-gray-700">
            <div class="rounded-lg bg-rose-100 p-1.5 dark:bg-rose-900/30">
              <svg class="h-4 w-4 text-rose-600 dark:text-rose-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <h3 class="text-sm font-semibold text-gray-900 dark:text-white">
              {{ t('admin.accounts.stats.kimiCodingPlan') }}
            </h3>
          </div>

          <!-- Loading state -->
          <div v-if="kimiLoading" class="flex items-center justify-center py-8">
            <LoadingSpinner />
          </div>

          <!-- Kimi data -->
          <template v-else-if="kimiUsage">
            <!-- Summary: total limit & remaining -->
            <div class="grid grid-cols-2 gap-4">
              <div class="card border-rose-200 bg-gradient-to-br from-rose-50 to-white p-4 dark:border-rose-800/30 dark:from-rose-900/10 dark:to-dark-700">
                <span class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('admin.accounts.stats.kimiTotalLimit') }}</span>
                <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ formatNumber(parseInt(kimiUsage.usage.limit || '0', 10)) }}</p>
                <p class="text-xs text-gray-400 dark:text-gray-500">{{ t('admin.accounts.stats.kimiQuota') }}</p>
              </div>
              <div class="card border-emerald-200 bg-gradient-to-br from-emerald-50 to-white p-4 dark:border-emerald-800/30 dark:from-emerald-900/10 dark:to-dark-700">
                <span class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('admin.accounts.stats.kimiRemaining') }}</span>
                <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ formatNumber(parseInt(kimiUsage.usage.remaining || '0', 10)) }}</p>
                <p class="text-xs text-gray-400 dark:text-gray-500">
                  {{ t('admin.accounts.stats.kimiResetTime') }}: <span class="font-medium text-emerald-600 dark:text-emerald-400">{{ formatResetDate(new Date(kimiUsage.usage.resetTime).getTime()) }}</span>
                </p>
              </div>
            </div>

            <!-- Window Limits -->
            <div v-if="kimiUsage.limits?.length" class="card p-4">
              <h4 class="mb-3 text-sm font-semibold text-gray-900 dark:text-white">
                {{ t('admin.accounts.stats.kimiWindowLimits') }}
              </h4>
              <div class="space-y-3">
                <div v-for="(limit, idx) in kimiUsage.limits" :key="idx">
                  <div class="mb-1 flex items-center justify-between">
                    <span class="text-xs font-medium text-gray-700 dark:text-gray-300">{{ formatKimiWindowLabel(limit) }}</span>
                    <span class="text-xs tabular-nums text-gray-500 dark:text-gray-400">
                      {{ t('admin.accounts.stats.kimiRemaining') }}: {{ formatNumber(parseInt(limit.detail.remaining || '0', 10)) }} / {{ formatNumber(parseInt(limit.detail.limit || '0', 10)) }}
                    </span>
                  </div>
                  <div class="h-2 w-full rounded-full bg-gray-200 dark:bg-gray-700">
                    <div
                      class="h-2 rounded-full transition-all"
                      :class="getKimiBarColor(limit)"
                      :style="{ width: Math.min(kimiUtilization(limit), 100) + '%' }"
                    />
                  </div>
                  <div v-if="limit.detail.resetTime" class="mt-0.5 text-xs tabular-nums text-gray-400 dark:text-gray-500">
                    {{ t('admin.accounts.stats.kimiResetTime') }}: {{ formatResetDate(new Date(limit.detail.resetTime).getTime()) }}
                  </div>
                </div>
              </div>
            </div>
          </template>
        </template>
      </template>

      <!-- No Data State -->
      <div
        v-else-if="!loading"
        class="flex flex-col items-center justify-center py-12 text-gray-500 dark:text-gray-400"
      >
        <Icon name="chartBar" size="xl" class="mb-4 h-12 w-12" />
        <p class="text-sm">{{ t('admin.accounts.stats.noData') }}</p>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end">
        <button
          @click="handleClose"
          class="rounded-lg bg-gray-100 px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-200 dark:bg-dark-600 dark:text-gray-300 dark:hover:bg-dark-500"
        >
          {{ t('common.close') }}
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'
import { Line } from 'vue-chartjs'
import BaseDialog from '@/components/common/BaseDialog.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import ModelDistributionChart from '@/components/charts/ModelDistributionChart.vue'
import EndpointDistributionChart from '@/components/charts/EndpointDistributionChart.vue'
import Icon from '@/components/icons/Icon.vue'
import { adminAPI } from '@/api/admin'
import type { GLMUsageResponse, MiniMaxModelRemain, KimiUsageResponse } from '@/api/admin/accounts'
import type { Account, AccountUsageStatsResponse } from '@/types'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

const { t } = useI18n()

const props = defineProps<{
  show: boolean
  account: Account | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const loading = ref(false)
const stats = ref<AccountUsageStatsResponse | null>(null)

// GLM usage state
const glmUsage = ref<GLMUsageResponse | null>(null)
const glmLoading = ref(false)

// MiniMax usage state
const minimaxUsage = ref<MiniMaxModelRemain[] | null>(null)
const minimaxLoading = ref(false)

// Coding plan platform detection
const codingPlanPlatform = ref<'glm' | 'minimax' | 'kimi' | null>(null)

// Kimi usage state
const kimiUsage = ref<KimiUsageResponse | null>(null)
const kimiLoading = ref(false)

// Dark mode detection
const isDarkMode = computed(() => {
  return document.documentElement.classList.contains('dark')
})

// Chart colors
const chartColors = computed(() => ({
  text: isDarkMode.value ? '#e5e7eb' : '#374151',
  grid: isDarkMode.value ? '#374151' : '#e5e7eb'
}))

// Line chart data
const trendChartData = computed(() => {
  if (!stats.value?.history?.length) return null

  return {
    labels: stats.value.history.map((h) => h.label),
    datasets: [
      {
        label: t('usage.accountBilled') + ' (USD)',
        data: stats.value.history.map((h) => h.actual_cost),
        borderColor: '#3b82f6',
        backgroundColor: 'rgba(59, 130, 246, 0.1)',
        fill: true,
        tension: 0.3,
        yAxisID: 'y'
      },
      {
        label: t('usage.userBilled') + ' (USD)',
        data: stats.value.history.map((h) => h.user_cost),
        borderColor: '#10b981',
        backgroundColor: 'rgba(16, 185, 129, 0.08)',
        fill: false,
        tension: 0.3,
        borderDash: [5, 5],
        yAxisID: 'y'
      },
      {
        label: t('admin.accounts.stats.requests'),
        data: stats.value.history.map((h) => h.requests),
        borderColor: '#f97316',
        backgroundColor: 'rgba(249, 115, 22, 0.1)',
        fill: false,
        tension: 0.3,
        yAxisID: 'y1'
      }
    ]
  }
})

// Line chart options with dual Y-axis
const lineChartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  interaction: {
    intersect: false,
    mode: 'index' as const
  },
  plugins: {
    legend: {
      position: 'top' as const,
      labels: {
        color: chartColors.value.text,
        usePointStyle: true,
        pointStyle: 'circle',
        padding: 15,
        font: {
          size: 11
        }
      }
    },
    tooltip: {
      callbacks: {
        label: (context: any) => {
          const label = context.dataset.label || ''
          const value = context.raw
          if (label.includes('USD')) {
            return `${label}: $${formatCost(value)}`
          }
          return `${label}: ${formatNumber(value)}`
        }
      }
    }
  },
  scales: {
    x: {
      grid: {
        color: chartColors.value.grid
      },
      ticks: {
        color: chartColors.value.text,
        font: {
          size: 10
        },
        maxRotation: 45,
        minRotation: 0
      }
    },
    y: {
      type: 'linear' as const,
      display: true,
      position: 'left' as const,
      grid: {
        color: chartColors.value.grid
      },
      ticks: {
        color: '#3b82f6',
        font: {
          size: 10
        },
        callback: (value: string | number) => '$' + formatCost(Number(value))
      },
      title: {
        display: true,
        text: t('usage.accountBilled') + ' (USD)',
        color: '#3b82f6',
        font: {
          size: 11
        }
      }
    },
    y1: {
      type: 'linear' as const,
      display: true,
      position: 'right' as const,
      grid: {
        drawOnChartArea: false
      },
      ticks: {
        color: '#f97316',
        font: {
          size: 10
        },
        callback: (value: string | number) => formatNumber(Number(value))
      },
      title: {
        display: true,
        text: t('admin.accounts.stats.requests'),
        color: '#f97316',
        font: {
          size: 11
        }
      }
    }
  }
}))

// Load stats when modal opens
watch(
  () => props.show,
  async (newVal) => {
    if (newVal && props.account) {
      await Promise.allSettled([loadStats(), loadCodingPlanUsage()])
    } else {
      stats.value = null
      glmUsage.value = null
      minimaxUsage.value = null
      kimiUsage.value = null
      codingPlanPlatform.value = null
    }
  }
)

const loadStats = async () => {
  if (!props.account) return

  loading.value = true
  try {
    stats.value = await adminAPI.accounts.getStats(props.account.id, 30)
  } catch (error) {
    console.error('Failed to load account stats:', error)
    stats.value = null
  } finally {
    loading.value = false
  }
}

const loadCodingPlanUsage = async () => {
  if (!props.account) return

  glmLoading.value = true
  minimaxLoading.value = true
  kimiLoading.value = true
  try {
    const result = await adminAPI.accounts.getCodingPlanUsage(props.account.id)
    codingPlanPlatform.value = result.platform
    if (result.platform === 'glm' && result.glm) {
      glmUsage.value = result.glm
    } else if (result.platform === 'minimax' && result.minimax) {
      minimaxUsage.value = result.minimax.models
    } else if (result.platform === 'kimi' && result.kimi) {
      kimiUsage.value = result.kimi
    }
  } catch {
    // Not a coding plan account or API error — silently hide
    glmUsage.value = null
    minimaxUsage.value = null
    kimiUsage.value = null
    codingPlanPlatform.value = null
  } finally {
    glmLoading.value = false
    minimaxLoading.value = false
    kimiLoading.value = false
  }
}

const handleClose = () => {
  emit('close')
}

// Format helpers
const formatCost = (value: number): string => {
  if (value >= 1000) {
    return (value / 1000).toFixed(2) + 'K'
  } else if (value >= 1) {
    return value.toFixed(2)
  } else if (value >= 0.01) {
    return value.toFixed(3)
  }
  return value.toFixed(4)
}

const formatNumber = (value: number): string => {
  if (value >= 1_000_000) {
    return (value / 1_000_000).toFixed(2) + 'M'
  } else if (value >= 1_000) {
    return (value / 1_000).toFixed(2) + 'K'
  }
  return value.toLocaleString()
}

const formatTokens = (value: number): string => {
  if (value >= 1_000_000_000) {
    return `${(value / 1_000_000_000).toFixed(2)}B`
  } else if (value >= 1_000_000) {
    return `${(value / 1_000_000).toFixed(2)}M`
  } else if (value >= 1_000) {
    return `${(value / 1_000).toFixed(2)}K`
  }
  return value.toLocaleString()
}

const formatDuration = (ms: number): string => {
  if (ms >= 1000) {
    return `${(ms / 1000).toFixed(2)}s`
  }
  return `${Math.round(ms)}ms`
}

// GLM quota label: "Token usage (5 Hour)" or "MCP usage (1 Month)"
const formatQuotaLabel = (limit: { type: string; unit: number; number: number }) => {
  const unitMap: Record<number, string> = { 1: 'Year', 2: 'Month', 3: 'Hour', 4: 'Day', 5: 'Month', 6: 'Week' }
  const typeLabel = limit.type === 'TOKENS_LIMIT' ? t('admin.accounts.stats.glmTokenQuota') : t('admin.accounts.stats.glmMcpQuota')
  const unitLabel = unitMap[limit.unit] || `${limit.number}${limit.unit}`
  return `${typeLabel} (${limit.number} ${unitLabel})`
}

const formatResetTime = (ms: number): string => {
  const d = new Date(ms)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${pad(d.getHours())}:${pad(d.getMinutes())}`
}

// Kimi helpers
const formatKimiWindowLabel = (limit: { window: { duration: number; timeUnit: string } }) => {
  const unitMap: Record<string, string> = {
    'TIME_UNIT_MINUTE': 'min',
    'TIME_UNIT_SECOND': 's',
    'TIME_UNIT_HOUR': 'h',
    'TIME_UNIT_DAY': 'd'
  }
  const unit = unitMap[limit.window.timeUnit] || limit.window.timeUnit
  return `${limit.window.duration}${unit} window`
}

const kimiUtilization = (limit: { detail: { limit: string; remaining: string } }) => {
  const total = parseInt(limit.detail.limit, 10) || 0
  const remaining = parseInt(limit.detail.remaining, 10) || 0
  if (total <= 0) return 0
  const used = total - remaining
  return (used / total) * 100
}

const getKimiBarColor = (limit: { detail: { limit: string; remaining: string } }) => {
  const total = parseInt(limit.detail.limit, 10) || 0
  const remaining = parseInt(limit.detail.remaining, 10) || 0
  if (total <= 0) return 'bg-gray-400'
  const pct = (remaining / total) * 100
  if (pct < 20) return 'bg-red-500'
  if (pct < 50) return 'bg-amber-500'
  return 'bg-emerald-500'
}

const formatResetDate = (ms: number): string => {
  const d = new Date(ms)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}
</script>
