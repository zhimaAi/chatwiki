<template>
  <a-modal
    :open="open"
    :title="t('title')"
    :width="1080"
    :footer="null"
    :maskClosable="false"
    centered
    class="prompt-optimize-modal"
    @cancel="close"
  >
    <a-spin :spinning="loading">
      <div class="optimizer" ref="optimizerEl">
        <section class="left">
          <div class="history" ref="historyEl">
            <a-button v-if="hasMore" type="link" :loading="loadingHistory" :disabled="generating" @click="loadHistory">{{
              t('more')
            }}</a-button>
            <template v-for="block in blocks" :key="block.id">
              <div v-if="block.guide" class="guide">
                <h3><SparkleIcon filled />{{ t('guide') }}</h3>
                <p>{{ t('description') }}</p>
                <button
                  v-for="(label, i) in suggestions"
                  :key="i"
                  :disabled="generating"
                  @click="fillSuggestion(label)"
                >
                  <span>{{ label }}</span><ArrowRightOutlined />
                </button>
              </div>
              <article v-else>
                <div class="question">{{ block.demand }}</div>
                <div class="assistant"><SparkleIcon />{{ t('assistant') }}</div>
                <details
                  v-if="block.reasoning_content?.trim()"
                  :open="block.reasoningOpen"
                  class="reasoning"
                  @toggle="block.reasoningOpen = $event.currentTarget.open"
                >
                  <summary>{{ t('reasoning') }}<DownOutlined class="reasoning-chevron" /></summary>
                  <div class="reasoning-text">{{ block.reasoning_content }}</div>
                </details>
                <div v-if="block.reply || (block.status === 'success' && block.summary)" class="reply-bubble">
                  <ReplyText :content="block.reply || block.summary" />
                </div>
                <a-button
                  v-if="block.status === 'success' && Number(block.version_id) > 0"
                  class="view-version"
                  :disabled="!canSwitchVersion"
                  @click="selectVersion(block.version_id)"
                >
                  {{ versionNumber(block.version_id) ? t('viewVersion', { version: versionNumber(block.version_id) }) : t('viewSavedVersion') }}<RightOutlined />
                </a-button>
                <div v-if="block.status === 'pending'" class="reply-status" role="status">
                  <a-spin size="small" />{{ t('replying') }}
                </div>
                <a-alert v-else-if="block.status !== 'success'" type="error" :message="block.error || t('failed')" />
              </article>
            </template>
          </div>
          <div class="composer">
            <div v-if="slashOpen" class="names">
              <div class="names-title">{{ t('insert') }}</div>
              <div class="names-list">
              <button
                v-for="item in filteredNames"
                :key="item.type + item.name"
                @mousedown.prevent="insertName(item.name)"
              >
                <span>{{ item.name }}</span
                ><a-tag :color="item.type === 'Skill' ? 'blue' : 'green'">{{ item.type }}</a-tag>
              </button>
              <span v-if="!filteredNames.length">{{ t('noNames') }}</span>
              </div>
            </div>
            <textarea
              ref="inputEl"
              v-model="demand"
              :disabled="generating"
              :placeholder="t('placeholder')"
              @input="updateSlash"
              @click="updateSlash"
              @keydown="onKey"
            />
            <div class="controls">
              <a-radio-group
                class="mode-switch"
                v-model:value="mode"
                :disabled="generating"
                size="small"
              >
                <a-tooltip :title="t('optimizeTooltip')" placement="top">
                  <a-radio-button value="optimize">{{ t('optimize') }}</a-radio-button>
                </a-tooltip>
                <a-tooltip :title="t('rewriteTooltip')" placement="top">
                  <a-radio-button value="rewrite">{{ t('rewrite') }}</a-radio-button>
                </a-tooltip>
              </a-radio-group>
              <div class="model" :inert="loading || generating ? '' : null">
                <ModelSelect
                  modelType="LLM"
                  :dropdown-match-select-width="426"
                  :dropdown-style="{ minWidth: 0 }"
                  :mode-id="modelId"
                  :mode-name="modelName"
                  @loaded="onModelsLoaded"
                  @change="onModelChange"
                />
              </div>
              <div v-if="supportsReasoning" class="effort" :inert="loading || generating ? '' : null">
                <span>{{ t('reasoningEffort') }}</span>
                <ReasoningEffortSelect compact v-model:value="reasoningEffort" @change="saveModelPreference" />
              </div>
              <a-button
                class="send-button"
                shape="circle"
                type="primary"
                :loading="generating"
                :disabled="!canSend"
                :aria-label="t('send')"
                @click="send"
                ><ArrowUpOutlined /></a-button
              >
            </div>
          </div>
          <div class="tip">{{ t('tipBefore') }} <kbd>/</kbd> {{ t('tipAfter') }}</div>
        </section>
        <section class="right">
          <header>
            <strong>{{ t(streamingPrompt !== null || selectedId !== CURRENT_VERSION_ID ? 'optimized' : 'current') }}</strong>
            <a-select
              class="version-select"
              :value="pendingVersionId ?? selectedId"
              :loading="versionLoading"
              :disabled="!canSwitchVersion"
              :aria-label="t('versions')"
              @change="selectVersion"
            >
              <a-select-option :value="CURRENT_VERSION_ID">{{ t('currentVersion') }}</a-select-option>
              <a-select-option v-for="v in versions" :key="v.id" :value="String(v.id)">
                {{ t('versionLabel', { version: v.version, total: versions.length }) }}
              </a-select-option>
            </a-select>
          </header>
          <pre ref="previewEl" class="preview" :aria-busy="versionLoading">{{ streamingPrompt ?? preview }}</pre>
          <a-alert v-if="canUseVersion && adopted && saveMode === 'draft'" class="adopt-feedback" type="success" show-icon :message="t('draftApplied')" />
          <footer v-if="canUseVersion">
            <a-button :disabled="saving || savingLibrary" @click="openLibrary"><ReadOutlined />{{ t('library') }}</a-button>
            <a-button
              :class="{ adopted }"
              type="primary"
              :disabled="savingLibrary || adopted"
              :loading="saving"
              @click="adopt"
            ><CheckOutlined />{{ adopted ? t('adopted') : t('adopt') }}</a-button>
          </footer>
        </section>
      </div>
    </a-spin>
  </a-modal>
  <a-modal
    v-model:open="libraryOpen"
    :title="t('library')"
    :width="472"
    centered
    class="prompt-library-modal"
    wrap-class-name="prompt-optimize-library-wrap"
    :getContainer="getLibraryContainer"
    :confirmLoading="savingLibrary"
    @ok="saveLibrary"
    ><label>{{ t('templateTitle') }}</label
    ><a-input v-model:value="title" :maxLength="10" :placeholder="t('titlePlaceholder')"
  /></a-modal>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, reactive, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { ArrowRightOutlined, ArrowUpOutlined, CheckOutlined, DownOutlined, ReadOutlined, RightOutlined } from '@ant-design/icons-vue'
import ReplyText from './reply-text.vue'
import SparkleIcon from './sparkle-icon.vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useStorage } from '@/hooks/web/useStorage'
import { useUserStore } from '@/stores/modules/user'
import { useLocaleStoreWithOut } from '@/stores/modules/locale'
import { fetchEventSource } from '@/libs/fetch-event-source'
import ModelSelect from '@/components/model-select/model-select.vue'
import ReasoningEffortSelect from '@/components/reasoning-effort-select/index.vue'
import {
  getOptimizeContext,
  getOptimizeRecords,
  getOptimizeVersions,
  getOptimizeVersion
} from '@/api/prompt-optimize'
import { savePromptLibraryItems } from '@/api/user'
import { getUuid } from '@/utils'
import { isValidReasoningEffort } from '@/utils/reasoning-effort'
import { composeTimeline, findSlash, insertPlainName, mergeOlderRecords } from './state'

const props = defineProps({
  savePrompt: { type: Function, required: true },
  saveMode: { type: String, default: 'persist' }
})
const emit = defineEmits(['adopted'])
const { t } = useI18n('components.prompt-optimize')
const userStore = useUserStore()
const { getStorage, setStorage } = useStorage()
const modelPreferenceKey = computed(() => {
  const { admin_user_id, user_id } = userStore.userInfo || {}
  return user_id ? `prompt-optimize-model:${admin_user_id || user_id}:${user_id}` : ''
})
const CURRENT_VERSION_ID = '0'
const open = ref(false),
  loading = ref(false),
  loadingHistory = ref(false),
  generating = ref(false)
const oldRecords = ref([]),
  newRecords = ref([]),
  versions = ref([]),
  names = ref([])
const demand = ref(''),
  mode = ref('optimize'),
  modelId = ref(''),
  modelName = ref(''),
  reasoningEffort = ref('none')
const preview = ref(''),
  streamingPrompt = ref(null),
  selectedId = ref(CURRENT_VERSION_ID),
  pendingVersionId = ref(null),
  adoptedId = ref(null),
  saving = ref(false)
const libraryOpen = ref(false),
  savingLibrary = ref(false),
  title = ref(''),
  hasMore = ref(false)
const slashOpen = ref(false),
  slashQuery = ref(''),
  inputEl = ref(null),
  historyEl = ref(null)
const optimizerEl = ref(null)
const previewEl = ref(null)
const getLibraryContainer = () => optimizerEl.value?.closest('.ant-modal-content')
let robotId,
  scene,
  messageId,
  initialModel,
  beforeId,
  controller,
  generation = 0,
  versionRequest = 0,
  messageSequence = 0,
  slashStart = 0
const suggestions = computed(() => [t('suggestion1'), t('suggestion2'), t('suggestion3')])
const blocks = computed(() => composeTimeline(oldRecords.value, newRecords.value))
const filteredNames = computed(() =>
  names.value.filter((n) => n.name.toLowerCase().includes(slashQuery.value.toLowerCase()))
)
const availableModels = ref(null)
const selectedModel = computed(() => findModel(modelId.value, modelName.value))
const supportsReasoning = computed(() => Number(selectedModel.value?.thinking_type) === 2)
const versionLoading = computed(() => pendingVersionId.value !== null)
const adopted = computed(() => adoptedId.value === selectedId.value)
const canSwitchVersion = computed(() =>
  open.value && !loading.value && !generating.value && !saving.value && !savingLibrary.value && !libraryOpen.value
)
const canUseVersion = computed(() =>
  open.value && !loading.value && !generating.value && !versionLoading.value && selectedId.value !== CURRENT_VERSION_ID
)
const canSend = computed(() =>
  canSwitchVersion.value && !versionLoading.value && !!selectedModel.value && !!demand.value.trim()
)
watch(selectedId, () => {
  if (previewEl.value) previewEl.value.scrollTop = 0
}, { flush: 'post' })

function findModel(id, name) {
  return availableModels.value
    ?.find((config) => String(config.model_config_id) === String(id))
    ?.children.find((model) => model.use_model_name === name)
}

function readModelPreference() {
  try {
    return modelPreferenceKey.value ? getStorage(modelPreferenceKey.value) : null
  } catch {
    return null
  }
}

function saveModelPreference() {
  if (!modelPreferenceKey.value || !selectedModel.value) return
  try {
    setStorage(modelPreferenceKey.value, {
      id: modelId.value,
      name: modelName.value,
      reasoningEffort: supportsReasoning.value ? reasoningEffort.value : 'none'
    })
  } catch {
    // The current selection is still usable if browser storage is unavailable.
  }
}

function initializeModel() {
  if (!open.value || !initialModel || !availableModels.value) return
  const preference = readModelPreference()
  const cachedModel = findModel(preference?.id, preference?.name)
  const model = cachedModel || findModel(initialModel.id, initialModel.name)
  initialModel = null
  modelId.value = model?.model_config_id ?? ''
  modelName.value = model?.use_model_name ?? ''
  const cachedEffort = cachedModel && preference.reasoningEffort
  reasoningEffort.value = supportsReasoning.value && cachedEffort && isValidReasoningEffort(cachedEffort)
    ? cachedEffort
    : 'none'
}

function onModelsLoaded(models) {
  availableModels.value = models
  initializeModel()
}

function onModelChange(_, { modelId: id, modelName: name }) {
  if (String(modelId.value) !== String(id) || modelName.value !== name) {
    modelId.value = id
    modelName.value = name
    reasoningEffort.value = 'none'
  }
  saveModelPreference()
}
const versionNumber = (id) => versions.value.find((v) => String(v.id) === String(id))?.version || ''
function rememberVersion({ id, version, record_id, create_time }) {
  if (versions.value.some((entry) => String(entry.id) === String(id))) return
  versions.value = [...versions.value, { id, version, record_id, create_time }]
    .sort((a, b) => Number(b.version) - Number(a.version))
}
async function show(options) {
  close()
  const seq = ++generation
  robotId = options.robotId
  scene = options.scene || 'config'
  messageId = options.messageId || 0
  oldRecords.value = []
  newRecords.value = []
  versions.value = []
  names.value = []
  beforeId = undefined
  demand.value = ''
  mode.value = 'optimize'
  initialModel = null
  modelId.value = ''
  modelName.value = ''
  reasoningEffort.value = 'none'
  selectedId.value = CURRENT_VERSION_ID
  adoptedId.value = null
  preview.value = ''
  open.value = true
  loading.value = true
  try {
    const [context, list, history] = await Promise.all([
      getOptimizeContext({ robot_id: robotId }),
      getOptimizeVersions({ robot_id: robotId }),
      getOptimizeRecords({ robot_id: robotId })
    ])
    if (seq !== generation) return
    preview.value = context.data.prompt
    names.value = context.data.names || []
    initialModel = { id: context.data.model_config_id, name: context.data.use_model }
    initializeModel()
    versions.value = list.data || []
    acceptHistory(history.data)
    await scrollBottom()
  } catch (e) {
    if (seq === generation) {
      message.error(e.message || e.msg || t('failed'))
      close()
    }
  } finally {
    if (seq === generation) loading.value = false
  }
}
function acceptHistory(data) {
  const rows = data.list || []
  oldRecords.value = mergeOlderRecords(oldRecords.value, rows)
  beforeId = rows.at(-1)?.id
  hasMore.value = data.has_more
}
async function loadHistory() {
  if (loadingHistory.value || generating.value) return
  const seq = generation
  loadingHistory.value = true
  try {
    const r = await getOptimizeRecords({ robot_id: robotId, before_id: beforeId })
    if (seq === generation) {
      const viewport = historyEl.value
      const height = viewport?.scrollHeight
      const top = viewport?.scrollTop
      acceptHistory(r.data)
      await nextTick()
      if (viewport && viewport === historyEl.value && seq === generation && !generating.value) {
        viewport.scrollTop = top + viewport.scrollHeight - height
      }
    }
  } catch (error) {
    if (seq === generation) message.error(error.msg || error.message || t('failed'))
  } finally {
    if (seq === generation) loadingHistory.value = false
  }
}
async function selectVersion(id) {
  if (!canSwitchVersion.value) return
  const target = String(id)
  const seq = ++versionRequest
  const life = generation
  pendingVersionId.value = target
  try {
    const r = target === CURRENT_VERSION_ID
      ? await getOptimizeContext({ robot_id: robotId })
      : await getOptimizeVersion({ robot_id: robotId, version_id: target })
    if (seq !== versionRequest || life !== generation) return
    if (target !== CURRENT_VERSION_ID) rememberVersion(r.data)
    selectedId.value = target
    preview.value = r.data.prompt
  } catch (error) {
    if (seq === versionRequest && life === generation)
      message.error(error.msg || error.message || t('failed'))
  } finally {
    if (seq === versionRequest && life === generation) pendingVersionId.value = null
  }
}
function close() {
  generation++
  versionRequest++
  controller?.abort()
  controller = null
  generating.value = false
  loading.value = false
  loadingHistory.value = false
  pendingVersionId.value = null
  streamingPrompt.value = null
  slashOpen.value = false
  open.value = false
  libraryOpen.value = false
}
async function scrollBottom() {
  await nextTick()
  if (historyEl.value) historyEl.value.scrollTop = historyEl.value.scrollHeight
}
function fillSuggestion(label) {
  demand.value = label
  slashOpen.value = false
  slashQuery.value = ''
  nextTick(() => inputEl.value?.focus())
}
function updateSlash() {
  const match = findSlash(demand.value, inputEl.value?.selectionStart || 0)
  slashOpen.value = !!match
  slashQuery.value = match?.query || ''
  slashStart = match?.start || 0
}
function insertName(name) {
  const el = inputEl.value
  const result = insertPlainName(demand.value, slashStart, el.selectionStart, name)
  demand.value = result.text
  slashOpen.value = false
  nextTick(() => {
    el.focus()
    el.setSelectionRange(result.caret, result.caret)
  })
}
function onKey(event) {
  if (event.key === 'Escape') slashOpen.value = false
  if (event.key === 'Enter' && !event.shiftKey && !event.isComposing) {
    event.preventDefault()
    send()
  }
}
async function send() {
  if (!canSend.value) return
  const life = generation
  const row = reactive({
    id: `local-${++messageSequence}`,
    demand: demand.value.trim(),
    reply: '',
    reasoning_content: '',
    reasoningOpen: true,
    status: 'pending'
  })
  newRecords.value.push(row)
  demand.value = ''
  slashOpen.value = false
  generating.value = true
  void scrollBottom()
  const request = new AbortController()
  controller = request
  let completed = false
  let bodyStarted = false
  const body = new FormData()
  Object.entries({
    robot_id: robotId,
    scene,
    mode: mode.value,
    demand: row.demand,
    request_id: getUuid(32),
    model_config_id: modelId.value,
    use_model: modelName.value,
    reasoning_effort: supportsReasoning.value ? reasoningEffort.value : 'none',
    base_version_id: selectedId.value,
    message_id: messageId
  }).forEach(([k, v]) => body.append(k, v))
  try {
    await fetchEventSource(`${import.meta.env.VITE_BASE_API_URL}/manage/optimizePrompt`, {
      method: 'POST',
      body,
      signal: request.signal,
      openWhenHidden: true,
      // The local wrapper only propagates errors through onerror when this is true.
      // onerror always throws below, so this enables error handling, NOT retries.
      retryOnError: true,
      headers: {
        token: userStore.getToken,
        lang: useLocaleStoreWithOut().getCurrentLocale.lang
      },
      async onopen(response) {
        if (!response.ok || !response.headers.get('content-type')?.includes('text/event-stream')) {
          const data = await response.json().catch(() => ({}))
          throw new Error(data.msg || t('failed'))
        }
      },
      onmessage(event) {
        if (life !== generation || completed) return
        const viewport = historyEl.value
        const follow = viewport && viewport.scrollHeight - viewport.scrollTop - viewport.clientHeight <= 32
        switch (event.event) {
          case 'prompt_delta':
          case 'reply_delta': {
            const text = JSON.parse(event.data).text
            if (!bodyStarted && text.trim()) {
              bodyStarted = true
              row.reasoningOpen = false
            }
            if (event.event === 'prompt_delta') {
              if (streamingPrompt.value === null && previewEl.value) previewEl.value.scrollTop = 0
              streamingPrompt.value = (streamingPrompt.value ?? '') + text
              return
            }
            row.reply += text
            break
          }
          case 'reasoning_content':
            row.reasoning_content += JSON.parse(event.data).text
            break
          case 'error':
            throw new Error(JSON.parse(event.data).message || t('failed'))
          case 'done': {
            const result = JSON.parse(event.data)
            if (!result.record_id || typeof result.reply !== 'string' || !result.reply.trim()) {
              throw new Error(t('interrupted'))
            }
            const version = result.version
            if (version) {
              if (!version.id || typeof version.prompt !== 'string' || !version.prompt.trim()) {
                throw new Error(t('interrupted'))
              }
              rememberVersion(version)
              selectedId.value = String(version.id)
              preview.value = version.prompt
            }
            streamingPrompt.value = null
            if (!bodyStarted) row.reasoningOpen = false
            Object.assign(row, {
              record_id: result.record_id,
              reply: result.reply,
              reasoning_content: result.reasoning_content || '',
              version_id: version?.id,
              status: 'success'
            })
            completed = true
            break
          }
          default:
            return
        }
        if (follow) void scrollBottom()
      },
      onerror(error) {
        throw error
      },
      onclose() {
        if (!completed) throw new Error(t('interrupted'))
      }
    })
  } catch (e) {
    if (life === generation && !completed) {
      const error = e.message || t('failed')
      Object.assign(row, { status: 'failed', reply: '', error })
      message.error(error)
    }
  } finally {
    if (controller === request) controller = null
    if (life === generation) {
      streamingPrompt.value = null
      generating.value = false
    }
  }
}
async function adopt() {
  if (!canUseVersion.value || saving.value || savingLibrary.value || libraryOpen.value || adopted.value) return
  const life = generation
  const text = preview.value
  const savedVersion = selectedId.value
  saving.value = true
  try {
    const result = await props.savePrompt(text)
    if (result === false || life !== generation || savedVersion !== selectedId.value) return
    adoptedId.value = savedVersion
    emit('adopted', text)
    message.success(t(props.saveMode === 'draft' ? 'draftApplied' : 'adoptSuccess'))
  } catch (e) {
    if (life === generation) message.error(e.message || e.msg || t('failed'))
  } finally {
    saving.value = false
  }
}
async function saveLibrary() {
  if (!canUseVersion.value || !libraryOpen.value || !title.value.trim() || saving.value || savingLibrary.value) return
  savingLibrary.value = true
  try {
    await savePromptLibraryItems({
      title: title.value.trim(),
      prompt_type: 0,
      prompt: preview.value,
      group_id: 0
    })
    libraryOpen.value = false
    message.success(t('saved'))
  } finally {
    savingLibrary.value = false
  }
}
function openLibrary() {
  if (!canUseVersion.value || saving.value || savingLibrary.value) return
  title.value = ''
  libraryOpen.value = true
}
onBeforeUnmount(close)
defineExpose({ show })
</script>

<style scoped lang="less">
:global(.ant-modal.prompt-optimize-modal) {
  padding-bottom: 0;
}
:global(.ant-modal.prompt-optimize-modal .ant-modal-content) {
  padding: 0;
  border-radius: 16px;
  overflow: hidden;
}
:global(.ant-modal.prompt-optimize-modal .ant-modal-header) {
  box-sizing: border-box;
  height: 56px;
  padding: 16px 24px;
  margin: 0;
}
:global(.ant-modal.prompt-optimize-modal .ant-modal-body) {
  padding: 0;
}
.optimizer {
  display: grid;
  grid-template-columns: minmax(0, 3fr) minmax(0, 2fr);
  box-sizing: border-box;
  height: 638px;
  margin: 0;
  border-top: 1px solid #eee;
}
.left,
.right {
  min-width: 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
  padding: 20px 24px 16px;
}
.right {
  border-left: 1px solid #eee;
}
.history {
  flex: 1;
  overflow: auto;
  overscroll-behavior: contain;
  min-height: 0;
}
article {
  margin-bottom: 24px;
}
.guide h3 {
  font-size: 16px;
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 4px 0 14px;
}
.guide p,
.tip {
  font-size: 12px;
  color: #8c8c8c;
}
.guide p { line-height: 20px; margin-bottom: 16px; }
.guide button {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  max-width: 100%;
  font: inherit;
  border: 0;
  border-radius: 12px;
  background: #f5f5f5;
  padding: 14px 12px;
  text-align: left;
  margin: 8px 0;
  cursor: pointer;
  .anticon { color: #999; flex-shrink: 0; }
}
.question {
  background: #edf1ff;
  border-radius: 12px 12px 0 12px;
  padding: 14px;
  margin: 0 0 20px auto;
  width: fit-content;
  max-width: 90%;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}
.assistant {
  color: #8c8c8c;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.reply-bubble {
  width: fit-content;
  max-width: 100%;
  background: #f6f7f9;
  padding: 12px 14px;
  border-radius: 10px;
  margin: 0 0 8px;
}
.reasoning {
  margin-bottom: 10px;
  border-radius: 8px;
  background: #f6f7f9;
  color: #667085;
}
.reasoning summary {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 32px;
  padding: 7px 12px;
  border-radius: 10px;
  list-style: none;
  cursor: pointer;
  font-size: 12px;
  line-height: 18px;
  &::-webkit-details-marker { display: none; }
}
.reasoning-chevron { margin-left: auto; font-size: 10px; flex-shrink: 0; transition: transform .15s; }
.reasoning[open] .reasoning-chevron { transform: rotate(180deg); }
.reasoning-text { padding: 4px 12px 12px; line-height: 22px; white-space: pre-wrap; overflow-wrap: anywhere; }
.reply-status { display: flex; align-items: center; gap: 8px; color: #8c8c8c; font-size: 12px; }
.view-version.ant-btn { color: #2475fc; border-color: #2475fc; margin: 2px 0 4px; }
.composer {
  position: relative;
  border: 1px solid #e4e7ed;
  border-radius: 16px;
  padding: 12px;
  margin-top: 16px;
  box-shadow: 0 4px 20px #00000008;
}
textarea {
  border: 0;
  outline: 0;
  resize: none;
  width: 100%;
  height: 65px;
  display: block;
  padding: 0;
  font: inherit;
  color: #262626;
  background: #fff;
  line-height: 22px;
  margin-bottom: 12px;
  &::placeholder { color: #999; }
}
.controls {
  display: flex;
  align-items: center;
  gap: 8px;
}
.model {
  flex: 1;
  min-width: 0;
  :deep(.ant-select.ant-select-single .ant-select-selector) { border: 0; box-shadow: none !important; background: transparent; padding-inline: 8px; }
  :deep(.modal-seclet-new.ant-select .ant-select-clear .model-show-block .icon) {
    width: 16px;
    height: 16px;
    flex: 0 0 16px;
    object-fit: contain;
  }
}
.effort {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  gap: 2px;
  color: #667085;
  font-size: 12px;
}
.version-select.ant-select.ant-select-single {
  width: 146px;
  flex-shrink: 0;
  :deep(.ant-select-selector) { border: 1px solid #e5eaf3; border-radius: 12px; background: #f5f7fb; padding-inline: 12px; }
  :deep(.ant-select-selection-item), :deep(.ant-select-arrow) { color: #66758e; font-size: 12px; }
}
.mode-switch.ant-radio-group {
  flex-shrink: 0;
  padding: 2px;
  border-radius: 12px;
  background: #f0f0f0;
  :deep(.ant-radio-button-wrapper) { height: 28px; line-height: 28px; padding: 0 12px; border: 0; border-radius: 10px; background: transparent; color: #666; box-shadow: none; }
  :deep(.ant-radio-button-wrapper::before) { display: none; }
  :deep(.ant-radio-button-wrapper-checked:not(.ant-radio-button-wrapper-disabled)) { color: #2475fc; background: #fff; box-shadow: none; }
  :deep(.ant-radio-button-wrapper:not(.ant-radio-button-wrapper-disabled):hover),
  :deep(.ant-radio-button-wrapper:not(.ant-radio-button-wrapper-disabled):active) { color: #666; background: transparent; box-shadow: none; }
  :deep(.ant-radio-button-wrapper-checked:not(.ant-radio-button-wrapper-disabled):hover),
  :deep(.ant-radio-button-wrapper-checked:not(.ant-radio-button-wrapper-disabled):active) { color: #2475fc; background: #fff; box-shadow: none; }
}
.send-button.ant-btn {
  flex-shrink: 0;
  width: 24px;
  min-width: 24px;
  height: 24px;
  padding: 0;
  border: 0;
  box-shadow: none;
  font-size: 16px;
  &:disabled { background: #d8dde5; color: #fff; }
}
.tip {
  text-align: center;
  margin-top: 14px;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 4px;
  kbd { display: inline-flex; justify-content: center; align-items: center; width: 18px; height: 18px; border: 1px solid #ccc; border-radius: 4px; font: inherit; line-height: 1; background: #fafafa; color: #999; }
}
.names {
  position: absolute;
  bottom: 100%;
  left: 0;
  width: 100%;
  background: white;
  border: 1px solid #eee;
  border-radius: 12px;
  padding: 0;
  overflow: hidden;
  box-shadow: 0 4px 20px #0001;
  z-index: 2;
}
.names-title { padding: 8px 16px; color: #7a8699; font-size: 12px; line-height: 16px; border-bottom: 1px solid #e4e7ed; }
.names-list { padding: 4px; max-height: 170px; overflow-y: auto; }
.names-list > span { display: block; padding: 8px 12px; color: #7a8699; }
.names button {
  display: flex;
  justify-content: space-between;
  width: 100%;
  border: 0;
  background: white;
  align-items: center;
  gap: 12px;
  padding: 6px 12px;
  border-radius: 10px;
  color: #262626;
  font: inherit;
  line-height: 20px;
  text-align: left;
  > span:first-child { overflow: hidden; white-space: nowrap; text-overflow: ellipsis; }
  :deep(.ant-tag) { border: 0; border-radius: 12px; margin: 0; padding: 0 10px; font-size: 12px; line-height: 20px; flex-shrink: 0; }
  cursor: pointer;
}
.names button:hover {
  background: #f0f5ff;
  color: #2475fc;
}
header,
footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}
.preview {
  flex: 1;
  white-space: pre-wrap;
  overflow: auto;
  overscroll-behavior: contain;
  font-family: inherit;
  font-size: 14px;
  line-height: 1.8;
  margin: 20px 0;
  word-break: break-word;
  color: #666;
  min-height: 0;
}
footer {
  justify-content: flex-end;
}
.adopt-feedback { margin-bottom: 12px; }
.right footer .ant-btn.adopted,
.right footer .ant-btn.adopted:disabled {
  background: #1fa566;
  border-color: #1fa566;
  color: #fff;
  opacity: 1;
}
:global(.ant-modal.prompt-library-modal .ant-modal-content) {
  padding: 0;
  border-radius: 16px;
  overflow: hidden;
}
:global(.ant-modal-root .ant-modal-wrap.prompt-optimize-library-wrap) { position: absolute; inset: 0; overflow: auto; }
:global(.ant-modal.prompt-library-modal .ant-modal-header) {
  padding: 20px 24px 16px;
  margin: 0;
}
:global(.ant-modal.prompt-library-modal .ant-modal-body) { padding: 24px 32px 32px; }
:global(.ant-modal.prompt-library-modal .ant-modal-body label) { display: block; margin-bottom: 6px; }
:global(.ant-modal.prompt-library-modal .ant-modal-footer) { margin: 0; padding: 0 24px 16px; }
</style>
