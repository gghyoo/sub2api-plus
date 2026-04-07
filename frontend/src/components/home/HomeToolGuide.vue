<template>
  <div class="mt-16">
    <!-- Section Title -->
    <div class="mb-8 text-center">
      <h2 class="mb-3 text-2xl font-bold text-gray-900 dark:text-white">
        {{ t('home.toolGuide.title') }}
      </h2>
      <p class="text-sm text-gray-600 dark:text-dark-400">
        {{ t('home.toolGuide.subtitle') }}
      </p>
    </div>

    <!-- Tool Cards Stacked Vertically -->
    <div class="mx-auto max-w-4xl space-y-8">
      <!-- Tool Card -->
      <div
        v-for="tool in tools"
        :key="tool.id"
        class="rounded-2xl border border-gray-200/50 bg-white/60 p-6 backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-800/60"
      >
        <!-- Tool Header -->
        <div class="mb-4 flex items-center gap-3">
          <div
            :class="[
              'flex h-10 w-10 items-center justify-center rounded-xl shadow-lg',
              tool.iconBg
            ]"
          >
            <component :is="tool.icon" class="h-5 w-5 text-white" />
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ tool.name }}
            </h3>
            <p class="text-xs text-gray-500 dark:text-dark-400">
              {{ tool.description }}
            </p>
          </div>
        </div>

        <!-- Platform Tabs -->
        <div v-if="tool.platforms.length > 1" class="border-b border-gray-200 dark:border-dark-700">
          <nav class="-mb-px flex space-x-6" aria-label="Platform">
            <button
              v-for="platform in tool.platforms"
              :key="platform.id"
              @click="setActivePlatform(tool.id, platform.id)"
              :class="[
                'whitespace-nowrap py-2.5 px-1 border-b-2 font-medium text-sm transition-colors',
                activePlatforms[tool.id] === platform.id
                  ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
              ]"
            >
              <span class="flex items-center gap-2">
                <component :is="platform.icon" class="w-4 h-4" />
                {{ platform.label }}
              </span>
            </button>
          </nav>
        </div>

        <!-- OS/Shell Tabs -->
        <div v-if="showShellTabs(tool)" class="border-b border-gray-200 dark:border-dark-700">
          <nav class="-mb-px flex space-x-4" aria-label="Shell">
            <button
              v-for="shell in shellTabs"
              :key="shell.id"
              @click="setActiveShell(tool.id, shell.id)"
              :class="[
                'whitespace-nowrap py-2.5 px-1 border-b-2 font-medium text-sm transition-colors',
                activeShells[tool.id] === shell.id
                  ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
              ]"
            >
              <span class="flex items-center gap-2">
                <component :is="shell.icon" class="w-4 h-4" />
                {{ shell.label }}
              </span>
            </button>
          </nav>
        </div>

        <!-- Code Blocks -->
        <div class="mt-4 space-y-4">
          <div
            v-for="(file, index) in getFiles(tool)"
            :key="index"
            class="relative"
          >
            <!-- File Hint -->
            <p v-if="file.hint" class="mb-1.5 flex items-center gap-1 text-xs text-amber-600 dark:text-amber-400">
              <Icon name="exclamationCircle" size="sm" class="flex-shrink-0" />
              {{ file.hint }}
            </p>
            <div class="overflow-hidden rounded-xl bg-gray-900 dark:bg-dark-900">
              <!-- Code Header -->
              <div class="flex items-center justify-between border-b border-gray-700 bg-gray-800 px-4 py-2 dark:border-dark-700 dark:bg-dark-800">
                <span class="font-mono text-xs text-gray-400">{{ file.path }}</span>
                <button
                  @click="copyContent(file.content, `${tool.id}-${activePlatforms[tool.id]}-${index}`)"
                  class="flex items-center gap-1.5 rounded-lg px-2.5 py-1 text-xs font-medium transition-colors"
                  :class="copiedId === `${tool.id}-${activePlatforms[tool.id]}-${index}`
                    ? 'bg-green-500/20 text-green-400'
                    : 'bg-gray-700 text-gray-300 hover:bg-gray-600 hover:text-white'"
                >
                  <svg v-if="copiedId === `${tool.id}-${activePlatforms[tool.id]}-${index}`" class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>
                  <svg v-else class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15.666 3.888A2.25 2.25 0 0013.5 2.25h-3c-1.03 0-1.9.693-2.166 1.638m7.332 0c.055.194.084.4.084.612v0a.75.75 0 01-.75.75H9a.75.75 0 01-.75-.75v0c0-.212.03-.418.084-.612m7.332 0c.646.049 1.288.11 1.927.184 1.1.128 1.907 1.077 1.907 2.185V19.5a2.25 2.25 0 01-2.25 2.25H6.75A2.25 2.25 0 014.5 19.5V6.257c0-1.108.806-2.057 1.907-2.185a48.208 48.208 0 011.927-.184" />
                  </svg>
                  {{ copiedId === `${tool.id}-${activePlatforms[tool.id]}-${index}` ? t('home.toolGuide.copied') : t('home.toolGuide.copy') }}
                </button>
              </div>
              <!-- Code Content -->
              <pre class="overflow-x-auto p-4 font-mono text-sm text-gray-100"><code v-if="file.highlighted" v-html="file.highlighted"></code><code v-else v-text="file.content"></code></pre>
            </div>
          </div>
        </div>

        <!-- Usage Note -->
        <div v-if="getNote(tool)" class="mt-4 flex items-start gap-3 rounded-lg border border-blue-100 bg-blue-50 p-3 dark:border-blue-800 dark:bg-blue-900/20">
          <Icon name="infoCircle" size="md" class="mt-0.5 flex-shrink-0 text-blue-500" />
          <p class="text-sm text-blue-700 dark:text-blue-300">
            {{ getNote(tool) }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, computed, type Component } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import { useClipboard } from '@/composables/useClipboard'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const { copyToClipboard: clipboardCopy } = useClipboard()
const appStore = useAppStore()

const copiedId = ref<string | null>(null)

// --- Icon Components ---
const TerminalIcon: Component = {
  render() {
    return h('svg', {
      fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', 'stroke-width': '1.5', class: 'w-5 h-5 text-white'
    }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'm6.75 7.5 3 2.25-3 2.25m4.5 0h3m-9 8.25h13.5A2.25 2.25 0 0 0 21 17.25V6.75A2.25 2.25 0 0 0 18.75 4.5H5.25A2.25 2.25 0 0 0 3 6.75v10.5A2.25 2.25 0 0 0 5.25 20.25Z' })
    ])
  }
}

const CodexIcon: Component = {
  render() {
    return h('svg', {
      fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', 'stroke-width': '1.5', class: 'w-5 h-5 text-white'
    }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M17.25 6.75 22.5 12l-5.25 5.25m-10.5 0L1.5 12l5.25-5.25m7.5-3-4.5 16.5' })
    ])
  }
}

const OpenCodeIcon: Component = {
  render() {
    return h('svg', {
      fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', 'stroke-width': '1.5', class: 'w-5 h-5 text-white'
    }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M9.813 15.904 9 18.75l-.813-2.846a4.5 4.5 0 0 0-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 0 0 3.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 0 0 3.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 0 0-3.09 3.09ZM18.259 8.715 18 9.75l-.259-1.035a3.375 3.375 0 0 0-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 0 0 2.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 0 0 2.456 2.456L21.75 6l-1.035.259a3.375 3.375 0 0 0-2.456 2.456Z' })
    ])
  }
}

const SmallTerminalIcon: Component = {
  render() {
    return h('svg', {
      fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', 'stroke-width': '1.5', class: 'w-4 h-4'
    }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'm6.75 7.5 3 2.25-3 2.25m4.5 0h3m-9 8.25h13.5A2.25 2.25 0 0 0 21 17.25V6.75A2.25 2.25 0 0 0 18.75 4.5H5.25A2.25 2.25 0 0 0 3 6.75v10.5A2.25 2.25 0 0 0 5.25 20.25Z' })
    ])
  }
}

const SparkleIcon: Component = {
  render() {
    return h('svg', {
      fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', 'stroke-width': '1.5', class: 'w-4 h-4'
    }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M9.813 15.904 9 18.75l-.813-2.846a4.5 4.5 0 0 0-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 0 0 3.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 0 0 3.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 0 0-3.09 3.09Z' })
    ])
  }
}

const AntigravityIcon: Component = {
  render() {
    return h('svg', {
      fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', 'stroke-width': '1.5', class: 'w-4 h-4'
    }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M15.59 14.37a6 6 0 0 1-5.84 7.38v-4.8m5.84-2.58a14.98 14.98 0 0 0 6.16-12.12A14.98 14.98 0 0 0 9.631 8.41m5.96 5.96a14.926 14.926 0 0 1-5.841 2.58m-.119-8.54a6 6 0 0 0-7.381 5.84h4.8m2.581-5.84a14.927 14.927 0 0 0-2.58 5.841m2.699 2.7c-.103.021-.207.041-.311.06a15.09 15.09 0 0 1-2.448-2.448 14.9 14.9 0 0 1 .06-.312m-2.24 2.39a4.493 4.493 0 0 0-1.757 4.306 4.493 4.493 0 0 0 4.306-1.758M16.5 9a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Z' })
    ])
  }
}

const AppleIcon: Component = {
  render() {
    return h('svg', { fill: 'currentColor', viewBox: '0 0 24 24', class: 'w-4 h-4' }, [
      h('path', { d: 'M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z' })
    ])
  }
}

const WindowsIcon: Component = {
  render() {
    return h('svg', { fill: 'currentColor', viewBox: '0 0 24 24', class: 'w-4 h-4' }, [
      h('path', { d: 'M3 12V6.75l6-1.32v6.48L3 12zm17-9v8.75l-10 .15V5.21L20 3zM3 13l6 .09v6.81l-6-1.15V13zm7 .25l10 .15V21l-10-1.91v-5.84z' })
    ])
  }
}

// --- Types ---
interface FileConfig {
  path: string
  content: string
  hint?: string
  highlighted?: string
}

interface PlatformConfig {
  id: string
  label: string
  icon: Component
}

interface ToolConfig {
  id: string
  name: string
  description: string
  icon: Component
  iconBg: string
  platforms: PlatformConfig[]
}

// --- Shell Tabs ---
const shellTabs = [
  { id: 'unix', label: 'macOS / Linux', icon: AppleIcon },
  { id: 'cmd', label: 'Windows CMD', icon: WindowsIcon },
  { id: 'powershell', label: 'PowerShell', icon: WindowsIcon }
]

// --- Tool Definitions ---
const tools: ToolConfig[] = [
  {
    id: 'claude-code',
    name: 'Claude Code',
    description: t('home.toolGuide.claudeCodeDesc'),
    icon: TerminalIcon,
    iconBg: 'bg-gradient-to-br from-orange-400 to-orange-600 shadow-orange-500/30',
    platforms: [
      { id: 'anthropic', label: 'Anthropic', icon: SmallTerminalIcon },
      { id: 'antigravity', label: 'Antigravity', icon: AntigravityIcon }
    ]
  },
  {
    id: 'codex',
    name: 'Codex CLI',
    description: t('home.toolGuide.codexDesc'),
    icon: CodexIcon,
    iconBg: 'bg-gradient-to-br from-green-500 to-green-600 shadow-green-500/30',
    platforms: [
      { id: 'openai', label: 'OpenAI', icon: SmallTerminalIcon }
    ]
  },
  {
    id: 'opencode',
    name: 'OpenCode',
    description: t('home.toolGuide.opencodeDesc'),
    icon: OpenCodeIcon,
    iconBg: 'bg-gradient-to-br from-purple-500 to-purple-600 shadow-purple-500/30',
    platforms: [
      { id: 'anthropic', label: 'Anthropic', icon: SmallTerminalIcon },
      { id: 'openai', label: 'OpenAI', icon: SmallTerminalIcon },
      { id: 'gemini', label: 'Gemini', icon: SparkleIcon },
      { id: 'antigravity', label: 'Antigravity', icon: AntigravityIcon }
    ]
  }
]

// --- Active State ---
const activePlatforms = reactive<Record<string, string>>({
  'claude-code': 'anthropic',
  'codex': 'openai',
  'opencode': 'anthropic'
})

const activeShells = reactive<Record<string, string>>({
  'claude-code': 'unix',
  'codex': 'unix',
  'opencode': 'unix'
})

function setActivePlatform(toolId: string, platformId: string) {
  activePlatforms[toolId] = platformId
}

function setActiveShell(toolId: string, shellId: string) {
  activeShells[toolId] = shellId
}

function showShellTabs(tool: ToolConfig): boolean {
  // OpenCode doesn't need shell tabs (it uses a JSON config file)
  return tool.id !== 'opencode'
}

// --- Code Generation ---
const BASE_URL = computed(() => appStore.apiBaseUrl || window.location.origin)
const API_KEY = 'sk-your-api-key'

function getFiles(tool: ToolConfig): FileConfig[] {
  const platform = activePlatforms[tool.id]
  const shell = activeShells[tool.id]

  switch (tool.id) {
    case 'claude-code':
      return getClaudeCodeFiles(platform, shell)
    case 'codex':
      return getCodexFiles(shell)
    case 'opencode':
      return getOpenCodeFiles(platform)
    default:
      return []
  }
}

function getClaudeCodeFiles(platform: string, shell: string): FileConfig[] {
  const baseUrl = platform === 'antigravity' ? `${BASE_URL.value}/antigravity` : BASE_URL.value

  let envPath: string
  let envContent: string
  switch (shell) {
    case 'unix':
      envPath = 'Terminal'
      envContent = `export ANTHROPIC_BASE_URL="${baseUrl}"
export ANTHROPIC_AUTH_TOKEN="${API_KEY}"
export CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1`
      break
    case 'cmd':
      envPath = 'Command Prompt'
      envContent = `set ANTHROPIC_BASE_URL=${baseUrl}
set ANTHROPIC_AUTH_TOKEN=${API_KEY}
set CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1`
      break
    case 'powershell':
      envPath = 'PowerShell'
      envContent = `$env:ANTHROPIC_BASE_URL="${baseUrl}"
$env:ANTHROPIC_AUTH_TOKEN="${API_KEY}"
$env:CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1`
      break
    default:
      envPath = 'Terminal'
      envContent = ''
  }

  const settingsPath = shell === 'unix' ? '~/.claude/settings.json' : '%userprofile%\\.claude\\settings.json'
  const settingsContent = `{
  "env": {
    "ANTHROPIC_BASE_URL": "${baseUrl}",
    "ANTHROPIC_AUTH_TOKEN": "${API_KEY}",
    "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1",
    "CLAUDE_CODE_ATTRIBUTION_HEADER": "0"
  }
}`

  return [
    { path: envPath, content: envContent },
    { path: settingsPath, content: settingsContent, hint: 'VSCode Claude Code' }
  ]
}

function getCodexFiles(shell: string): FileConfig[] {
  const isWindows = shell === 'cmd' || shell === 'powershell'
  const configDir = isWindows ? '%userprofile%\\.codex' : '~/.codex'

  const configContent = `model_provider = "OpenAI"
model = "gpt-5.4"
review_model = "gpt-5.4"
model_reasoning_effort = "xhigh"
disable_response_storage = true
network_access = "enabled"
windows_wsl_setup_acknowledged = true
model_context_window = 1000000
model_auto_compact_token_limit = 900000

[model_providers.OpenAI]
name = "OpenAI"
base_url = "${BASE_URL}/v1"
wire_api = "responses"
requires_openai_auth = true`

  const authContent = `{
  "OPENAI_API_KEY": "${API_KEY}"
}`

  return [
    {
      path: `${configDir}/config.toml`,
      content: configContent,
      hint: t('home.toolGuide.codexConfigHint')
    },
    {
      path: `${configDir}/auth.json`,
      content: authContent
    }
  ]
}

function getOpenCodeFiles(platform: string): FileConfig[] {
  const baseRoot = BASE_URL

  if (platform === 'antigravity') {
    return [
      generateOpenCodeAntigravityConfig(baseRoot)
    ]
  }

  const providerKey = platform
  const baseUrl = platform === 'gemini' ? `${baseRoot}/v1beta` : `${baseRoot}/v1`

  const provider: Record<string, any> = {
    [providerKey]: {
      options: {
        baseURL: baseUrl,
        apiKey: API_KEY
      }
    }
  }

  if (platform === 'gemini') {
    provider[providerKey].npm = '@ai-sdk/google'
  } else if (platform === 'anthropic') {
    provider[providerKey].npm = '@ai-sdk/anthropic'
  }

  const content = JSON.stringify({
    provider,
    $schema: 'https://opencode.ai/config.json'
  }, null, 2)

  return [{ path: 'opencode.json', content, hint: t('home.toolGuide.opencodeHint') }]
}

function generateOpenCodeAntigravityConfig(baseRoot: string): FileConfig {
  const provider: Record<string, any> = {
    'antigravity-claude': {
      npm: '@ai-sdk/anthropic',
      name: 'Antigravity (Claude)',
      options: {
        baseURL: `${baseRoot}/antigravity/v1`,
        apiKey: API_KEY
      },
      models: {
        'claude-opus-4-6-thinking': {
          name: 'Claude 4.6 Opus (Thinking)',
          limit: { context: 200000, output: 128000 },
          modalities: { input: ['text', 'image', 'pdf'], output: ['text'] },
          options: { thinking: { budgetTokens: 24576, type: 'enabled' } }
        },
        'claude-sonnet-4-6': {
          name: 'Claude 4.6 Sonnet',
          limit: { context: 200000, output: 64000 },
          modalities: { input: ['text', 'image', 'pdf'], output: ['text'] },
          options: { thinking: { budgetTokens: 24576, type: 'enabled' } }
        }
      }
    },
    'antigravity-gemini': {
      npm: '@ai-sdk/google',
      name: 'Antigravity (Gemini)',
      options: {
        baseURL: `${baseRoot}/antigravity/v1beta`,
        apiKey: API_KEY
      },
      models: {
        'gemini-2.5-flash': {
          name: 'Gemini 2.5 Flash',
          limit: { context: 1048576, output: 65536 },
          modalities: { input: ['text', 'image', 'pdf'], output: ['text'] },
          options: { thinking: { budgetTokens: 24576, type: 'disable' } }
        },
        'gemini-2.5-pro': {
          name: 'Gemini 2.5 Pro',
          limit: { context: 2097152, output: 65536 },
          modalities: { input: ['text', 'image', 'pdf'], output: ['text'] },
          options: { thinking: { budgetTokens: 24576, type: 'enabled' } }
        }
      }
    }
  }

  const content = JSON.stringify({
    provider,
    $schema: 'https://opencode.ai/config.json'
  }, null, 2)

  return { path: 'opencode.json', content, hint: t('home.toolGuide.opencodeHint') }
}

// --- Notes ---
function getNote(tool: ToolConfig): string {
  const platform = activePlatforms[tool.id]
  switch (tool.id) {
    case 'claude-code':
      return platform === 'antigravity'
        ? t('home.toolGuide.claudeCodeAntigravityNote')
        : t('home.toolGuide.claudeCodeNote')
    case 'codex':
      return t('home.toolGuide.codexNote')
    case 'opencode':
      return ''
    default:
      return ''
  }
}

// --- Copy ---
async function copyContent(content: string, id: string) {
  const success = await clipboardCopy(content, t('home.toolGuide.copied'))
  if (success) {
    copiedId.value = id
    setTimeout(() => {
      copiedId.value = null
    }, 2000)
  }
}
</script>
