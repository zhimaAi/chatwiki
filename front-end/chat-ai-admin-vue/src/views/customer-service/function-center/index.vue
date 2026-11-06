<template>
  <section class="function-page">
    <CustomerPageHeader :title="pane === 'logs' ? t('invoke_logs') : pane === 'rate' ? t('rate_limit') : t('function_center')" :back="pane !== 'home'" :back-label="t('back')" @back="pane = 'home'" />
    <div class="content" :class="{ 'home-content': pane === 'home' }">
      <div v-show="pane === 'rate'" class="rate-pane"><RateLimit ref="rateLimit" :robot-key="robotKey" embedded @status-change="rateStatus = $event" /></div>
      <InvokeLogs v-if="pane === 'logs' && selected" :key="selected.id" :robot-id="selected.id" :robot-key="selected.robot_key" />
      <div v-else-if="pane === 'home'" class="function-home">
        <div class="workflow-tip"><img :src="icons.info" alt="" /><p>{{ t('workflow_tip') }}</p></div>
        <div class="workflow-toolbar">
          <a-dropdown :trigger="['click']"><a-button type="primary" class="add-workflow"><img :src="icons.add" alt="" />{{ t('add_workflow') }}</a-button>
            <template #overlay><a-menu><a-menu-item key="create" @click="creator?.open(1)">{{ t('create_workflow') }}</a-menu-item><a-menu-item key="import" @click="importer?.open()">{{ t('import_csl') }}</a-menu-item></a-menu></template>
          </a-dropdown>
        </div>
        <a-alert v-if="error" type="error" :message="error" show-icon><template #action><a @click="load">{{ t('retry') }}</a></template></a-alert>
        <a-alert v-if="rateStatus.error" class="rate-error" type="error" :message="rateStatus.error" show-icon><template #action><a @click="pane = 'rate'">{{ t('view') }}</a></template></a-alert>
        <a-spin :spinning="loading">
          <div class="cards">
            <article class="function-card rate-card" role="button" tabindex="0" @click="pane = 'rate'" @keydown.enter.self="pane = 'rate'" @keydown.space.self.prevent="pane = 'rate'">
              <div class="card-controls" @click.stop @keydown.stop><a-switch :checked="rateStatus.enabled" :loading="rateStatus.loading" :disabled="!rateStatus.ready || rateStatus.loading" :aria-label="t('rate_limit')" @change="toggleRate" /></div>
              <span class="feature-icon rate-icon"><img :src="icons.rate" alt="" /></span>
              <h3>{{ t('rate_limit') }}</h3><p :title="t('rate_limit_desc')">{{ t('rate_limit_desc') }}</p>
            </article>
            <article v-for="item in workflows" :key="item.id" class="function-card workflow-card" :class="{ active: isAssociated(item.id) }" role="button" tabindex="0" @click="openWorkflow(item)" @keydown.enter.self="openWorkflow(item)" @keydown.space.self.prevent="openWorkflow(item)">
              <div class="card-controls" @click.stop @keydown.stop>
                <a-dropdown :trigger="['click']" placement="bottomRight" :open="openWorkflowMenuId === String(item.id)" @open-change="open => updateWorkflowMenu(item.id, open)"><a-button type="text" class="more-button" :class="{ 'is-open': openWorkflowMenuId === String(item.id) }" :loading="editingId === String(item.id)" :aria-label="t('workflow_more')"><img :src="openWorkflowMenuId === String(item.id) ? icons.moreActive : icons.more" alt="" /></a-button><template #overlay><a-menu class="workflow-actions" @click="({ key }) => action(key, item)">
                <a-menu-item v-for="entry in workflowActions" :key="entry.key"><span class="workflow-menu-item"><span class="menu-icon"><img :src="entry.icon" alt="" /></span><span>{{ t(entry.label) }}</span></span></a-menu-item>
                </a-menu></template></a-dropdown>
                <a-switch :checked="isAssociated(item.id)" :loading="linkingId === String(item.id)" :disabled="!!linkingId || loading" :aria-label="t('workflow_association', { name: item.robot_name })" @change="checked => toggleWorkflow(item, checked)" />
              </div>
              <a-avatar class="feature-icon workflow-icon" :size="40" shape="square" :src="item.robot_avatar || item.robot_avatar_url"><template #icon><img :src="icons.workflow" alt="" /></template></a-avatar>
              <h3><span class="workflow-title" :title="item.robot_name">{{ item.robot_name }}</span><span class="workflow-badge">{{ t('custom_workflow') }}</span></h3>
              <p :title="item.robot_intro || '—'">{{ item.robot_intro || '—' }}</p>
            </article>
          </div>
          <a-empty v-if="!loading && !workflows.length" :description="t('empty_workflows')" />
        </a-spin>
      </div>
    </div>
    <AddRobotAlert ref="creator" create-redirect="none" @ok="created" @edit-robot="load" />
    <ImportCsl ref="importer" :application-type="1" @ok="imported" />
    <ExportModal ref="exporter" @ok="exportWorkflow" />
  </section>
</template>

<script setup>
import { computed, inject, onBeforeUnmount, ref } from 'vue'
import CustomerPageHeader from '../components/page-header.vue'
import { useRouter } from 'vue-router'
import { Modal, message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { getRobotList, getRobotInfo, robotCopy, deleteRobot, relationWorkFlow } from '@/api/robot'
import { useOpenUrlWithToken } from '@/hooks/web/useOpenUrlWithToken'
import AddRobotAlert from '@/views/robot/robot-list/components/add-robot-alert.vue'
import ImportCsl from '@/views/robot/robot-list/components/import-csl-alert.vue'
import ExportModal from '@/views/robot/robot-list/components/workflow-export-modal.vue'
import RateLimit from '@/views/robot/robot-config/rate-limit/index.vue'
import InvokeLogs from '@/views/robot/robot-config/invoke-logs/index.vue'
import infoIcon from '@/assets/img/customer-service/function-center/info.svg'
import addIcon from '@/assets/img/customer-service/function-center/add.svg'
import rateIcon from '@/assets/img/customer-service/function-center/rate.svg'
import workflowIcon from '@/assets/img/customer-service/function-center/workflow.svg'
import moreIcon from '@/assets/img/customer-service/knowledge-base/more.svg'
import moreActiveIcon from '@/assets/img/customer-service/function-center/more-active.svg'
import editIcon from '@/assets/img/customer-service/function-center/edit.svg'
import logsIcon from '@/assets/img/customer-service/function-center/logs.svg'
import copyIcon from '@/assets/img/customer-service/function-center/copy.svg'
import exportIcon from '@/assets/img/customer-service/function-center/export.svg'
import deleteIcon from '@/assets/img/customer-service/function-center/delete.svg'
const icons = { info: infoIcon, add: addIcon, rate: rateIcon, workflow: workflowIcon, more: moreIcon, moreActive: moreActiveIcon }
const workflowActions = [
  { key: 'edit', label: 'edit', icon: editIcon },
  { key: 'logs', label: 'invoke_logs', icon: logsIcon },
  { key: 'copy', label: 'copy', icon: copyIcon },
  { key: 'export', label: 'export_csl', icon: exportIcon },
  { key: 'delete', label: 'delete', icon: deleteIcon }
]
const { t } = useI18n('views.customer-service.index')
const router = useRouter()
const { openUrlWithToken } = useOpenUrlWithToken()
const context = inject('customerRobotContext')
const robotId = context.robotInfo.id
const robotKey = String(context.robotInfo.robot_key)
const pane = ref('home')
const selected = ref(null)
const workflows = ref([])
const loading = ref(false)
const error = ref('')
const creator = ref(null)
const importer = ref(null)
const exporter = ref(null)
const rateLimit = ref(null)
const rateStatus = ref({ ready: false, loading: true, enabled: false, error: '' })
const linkingId = ref('')
const editingId = ref('')
const openWorkflowMenuId = ref('')
const associatedIds = computed(() => String(context.robotInfo.work_flow_ids || '').split(',').filter(Boolean))
const isAssociated = id => associatedIds.value.includes(String(id))
let sequence = 0
let editSequence = 0
let alive = true
async function load() {
  const version = ++sequence
  loading.value = true
  error.value = ''
  try {
    const response = await getRobotList({ application_type: 1 })
    if (alive && version === sequence) workflows.value = (response.data || []).filter(item => Number(item.application_type) === 1)
  } catch (err) { if (alive && version === sequence) error.value = err.message || t('load_failed') }
  finally { if (alive && version === sequence) loading.value = false }
}
function openWorkflow(item, showTips = false) {
  if (!item?.id || !item.robot_key) { message.error(t('invalid_robot')); return }
  const query = { id: item.id, robot_key: item.robot_key, source: 'customer-service', ...(showTips ? { show_tips: 1 } : {}) }
  window.open(router.resolve({ path: '/robot/config/workflow', query }).href, '_blank', 'noopener,noreferrer')
}
function created(item) { if (!alive) return; void load(); openWorkflow(item, true) }
function imported(item) { if (!alive) return; void load(); openWorkflow(item) }
async function toggleRate(checked) {
  const success = await rateLimit.value?.setEnabled(checked)
  if (alive && !success) pane.value = 'rate'
}
async function toggleWorkflow(item, checked) {
  if (!alive || linkingId.value || loading.value) return
  linkingId.value = String(item.id)
  try {
    await context.waitForSave()
    if (!alive) return
    await context.refreshRobot()
    if (!alive || isAssociated(item.id) === checked) return
    const ids = checked ? [...associatedIds.value, String(item.id)] : associatedIds.value.filter(id => id !== String(item.id))
    const params = { id: robotId, work_flow_ids: [...new Set(ids)].join(',') }
    // 使用配置保存队列及详情刷新，避免后续全量配置保存覆盖工作流关联。
    const success = await context.saveAction(() => relationWorkFlow(params))
    if (alive && success) message.success(t('operation_success'))
  } catch (err) { if (alive) message.error(err.message || t('save_failed')) }
  finally { if (alive) linkingId.value = '' }
}
function exportWorkflow(item, options = {}) {
  const query = new URLSearchParams({ ...options, id: item.id })
  openUrlWithToken(`${import.meta.env.VITE_BASE_API_URL || ''}/manage/robotExport?${query}`)
}
async function editWorkflow(item) {
  const version = ++editSequence
  editingId.value = String(item.id)
  try {
    // 列表不返回 en_name，编辑时单独读工作流详情，不能借用当前客服 Store。
    const response = await getRobotInfo({ id: item.id })
    if (!alive || version !== editSequence) return
    const info = response.data
    if (String(info?.id) !== String(item.id) || info.robot_key !== item.robot_key || Number(info.application_type) !== 1) throw new Error(t('invalid_robot'))
    await creator.value.open(1, true, info)
  } catch (err) { if (alive && version === editSequence) message.error(err.message || t('load_failed')) }
  finally { if (alive && version === editSequence) editingId.value = '' }
}
function updateWorkflowMenu(id, open) {
  if (open) openWorkflowMenuId.value = String(id)
  else if (openWorkflowMenuId.value === String(id)) openWorkflowMenuId.value = ''
}
function action(key, item) {
  openWorkflowMenuId.value = ''
  if (key === 'edit') return editWorkflow(item)
  if (key === 'logs') { selected.value = item; pane.value = 'logs'; return }
  if (key === 'export') return exporter.value.handle(item)
  if (!['delete', 'copy'].includes(key)) return
  Modal.confirm({ title: t(key === 'delete' ? 'confirm_delete_workflow' : 'confirm_copy_workflow', { name: item.robot_name }), okType: key === 'delete' ? 'danger' : 'primary',
    onOk: async () => {
      if (!alive) return
      await (key === 'delete' ? deleteRobot({ id: item.id }) : robotCopy({ id: item.id }))
      if (alive) { message.success(t('operation_success')); await load() }
    }
  })
}
void load()
onBeforeUnmount(() => { alive = false; sequence++ })
</script>
<style lang="less" scoped>
.function-page { flex: 1; min-height: 0; display: flex; flex-direction: column; margin: 8px 8px 8px 0; overflow: hidden; border-radius: 12px; background: #fff; }
.content { flex: 1; min-height: 0; overflow: auto; padding: 20px 24px; }
.home-content { padding: 32px; }
.function-home { display: flex; flex-direction: column; gap: 24px; min-width: 0; container-type: inline-size; }
.rate-pane { height: 100%; }
.workflow-tip {
  display: flex; align-items: flex-start; gap: 12px; padding: 16px; border-radius: 8px; background: #f1f6fe;
  img { flex-shrink: 0; margin-top: 4px; }
  p { margin: 0; color: #363636; font-size: 14px; line-height: 24px; }
}
.workflow-toolbar { display: flex; align-items: center; }
.add-workflow { display: inline-flex; align-items: center; justify-content: center; gap: 6px; height: 32px; padding: 0 12px; border-radius: 6px; border-color: #146ef5; background: #146ef5; box-shadow: none; font-size: 14px; line-height: 22px; font-weight: 400; img { flex-shrink: 0; } }
.cards { display: grid; grid-template-columns: minmax(0, 1fr); gap: 16px; }
.function-card {
  display: grid; grid-template-columns: 40px minmax(0, 1fr) auto; grid-template-rows: 32px 22px; column-gap: 12px; row-gap: 8px;
  align-items: center; min-width: 0; min-height: 128px; padding: 32px 24px; border: 1px solid #e8e8e8; border-radius: 12px; cursor: pointer;
  background: #fff; transition: border-color .2s;
  &:hover { border-color: #d8d8d8; }
  &:focus-visible { outline: 2px solid #146ef5; outline-offset: 2px; }
  h3 { grid-column: 2; grid-row: 1; display: flex; align-items: center; gap: 8px; min-width: 0; margin: 0; font-size: 16px; line-height: 24px; font-weight: 500; color: #080808; }
  p { grid-column: 2 / -1; grid-row: 2; min-width: 0; margin: 0; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; color: #5a5a5a; font-size: 14px; line-height: 22px; }
}
.card-controls {
  grid-column: 3; grid-row: 1; display: flex; align-items: center; justify-content: flex-end; gap: 8px; height: 32px;
  :deep(.ant-switch) { min-width: 36px; height: 20px; background: #d8d8d8; }
  :deep(.ant-switch-checked) { background: #146ef5; }
  :deep(.ant-switch-handle) { width: 16px; height: 16px; top: 2px; inset-inline-start: 2px; }
  :deep(.ant-switch-checked .ant-switch-handle) { inset-inline-start: calc(100% - 18px); }
  :deep(.ant-switch-handle::before) { border-radius: 8px; box-shadow: 0 2px 8px rgba(8, 8, 8, .04); }
}
.more-button {
  display: inline-flex; align-items: center; justify-content: center; width: 32px; height: 32px; padding: 0; border: 0; border-radius: 6px; background: transparent; box-shadow: none;
  &.is-open, &:hover { background: #e8f0fe; }
}
.feature-icon { grid-column: 1; grid-row: 1 / 3; display: flex; align-items: center; justify-content: center; flex-shrink: 0; width: 40px; height: 40px; border-radius: 12px; }
.workflow-icon { background: transparent; :deep(img) { display: block; } }
.workflow-title { min-width: 0; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; }
.workflow-badge { flex-shrink: 0; padding: 2px 8px; border-radius: 4px; background: #f5f5f5; color: #5a5a5a; font-size: 12px; font-weight: 400; line-height: 20px; }
.workflow-actions {
  width: 144px; padding: 4px; border: 1px solid #e8e8e8; border-radius: 8px;
  box-shadow: 0 8px 24px 8px rgba(8, 8, 8, .04), 0 6px 12px rgba(8, 8, 8, .04), 0 4px 8px -8px rgba(8, 8, 8, .06);
  :deep(.ant-dropdown-menu-item) { min-height: 32px; margin: 0; padding: 4px 8px; border-radius: 4px; color: #222; font-size: 14px; line-height: 22px; }
}
.workflow-menu-item { display: flex; align-items: center; gap: 8px; }
.menu-icon { display: inline-flex; align-items: center; justify-content: center; flex-shrink: 0; width: 16px; height: 16px; }
@container (min-width: 736px) {
  .cards { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}
@container (min-width: 1112px) {
  .cards { grid-template-columns: repeat(3, minmax(0, 1fr)); }
}
@media (max-width: 500px) {
  .home-content { padding: 24px 16px; }
  .function-card { column-gap: 8px; padding-right: 16px; padding-left: 16px; }
}
</style>
