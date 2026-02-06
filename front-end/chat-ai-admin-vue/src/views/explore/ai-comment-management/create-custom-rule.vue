<template>
  <div class="default-rule-page">
    <div class="form-row">
      <div class="form-label">{{ t('rule_name_label') }}</div>
      <a-input
        v-model:value="ruleForm.rule_name"
        :placeholder="t('rule_name_placeholder')"
        style="width: 300px"
      />
    </div>

    <div class="form-row">
      <div class="form-label">{{ t('model_label') }}</div>
      <ModelSelect
        modelType="LLM"
        v-model:modeName="ruleForm.use_model"
        v-model:modeId="ruleForm.model_config_id"
        style="width: 300px;"
        :placeholder="t('select_model_placeholder')"
        @change="onModelChange"
      />
    </div>

    <!-- <div class="form-row">
      <div class="form-label">选择群发</div>
      <a-radio-group v-model:value="taskScope">
        <a-radio value="all">全部群发</a-radio>
        <a-radio value="specific">指定群发</a-radio>
      </a-radio-group>
    </div> -->

    <!-- <div class="radio-button">
      <a-button v-if="taskScope==='specific'" type="dashed" @click="openTaskModal"><PlusOutlined />选择群发</a-button>
    </div> -->

    <!-- <div v-if="taskScope==='specific' && selectedTasks.length" class="selected-task-table">
      <a-table :columns="selectedTaskColumns" :data-source="selectedTasks" row-key="id" :pagination="false">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key==='actions'">
            <a @click="removeSelectedTask(record.id)">删除</a>
          </template>
        </template>
      </a-table>
    </div> -->

    <div class="module-card">
      <div class="module-card-header">
        <div class="module-card-title">
          <svg-icon name="comment-delete" style="font-size: 16px; margin-right: 8px;" />
          <span>{{ t('module_delete_title') }}</span>
        </div>
        <a-switch v-model:checked="ruleForm.delete_comment_switch" :checkedValue="1" :unCheckedValue="0" />
      </div>
      <div v-if="ruleForm.delete_comment_switch==1" class="module-card-body">
        <div class="field">
          <div class="field-label">{{ t('delete_rule_label') }}</div>
          <a-checkbox-group v-model:value="ruleForm.delete_type">
            <a-checkbox :value="1">{{ t('delete_rule_sensitive_word') }}</a-checkbox>
            <a-checkbox :value="2">{{ t('delete_rule_ai_detection') }}</a-checkbox>
          </a-checkbox-group>
        </div>
        <div class="field" v-if="(ruleForm.delete_type || []).includes(1)">
          <div class="field-label">{{ t('delete_sensitive_word_label') }}</div>
          <div class="tag-list">
            <a-tag v-for="(kw, idx) in ruleForm.delete_keywords" :key="kw" closable @close="removeKeyword(idx)">{{ kw }}</a-tag>
          </div>
          <a-textarea
            v-model:value="deleteKeywordInput"
            :placeholder="t('common_input_placeholder')"
            @blur="onAddKeyword"
            :auto-size="{ minRows: 2 }"
          />
        </div>
        <div class="field" v-if="(ruleForm.delete_type || []).includes(2)">
          <div class="field-label">{{ t('delete_ai_rule_label') }}</div>
          <a-textarea
            v-model:value="ruleForm.delete_prompt"
            :placeholder="t('common_input_placeholder')"
            :auto-size="{ minRows: 3 }"
          />
        </div>
        <div class="field" v-if="(ruleForm.delete_type || []).includes(2)">
          <div class="field-label">{{ t('condition_label') }}</div>
          <a-radio-group v-model:value="ruleForm.delete_condition">
            <a-radio :value="1">{{ t('condition_all') }}</a-radio>
            <a-radio :value="2">{{ t('condition_any') }}</a-radio>
          </a-radio-group>
        </div>
        <div class="field" v-if="(ruleForm.delete_type || []).includes(2)">
          <div class="field-label">{{ t('priority_label') }}</div>
          <a-radio-group v-model:value="ruleForm.delete_priority">
            <a-radio :value="1">{{ t('priority_sensitive_first') }}</a-radio>
            <a-radio :value="2">{{ t('priority_ai_first') }}</a-radio>
          </a-radio-group>
        </div>
      </div>
    </div>

    <div class="module-card">
      <div class="module-card-header">
        <div class="module-card-title">
          <svg-icon name="comment-reply" style="font-size: 16px; margin-right: 8px;" />
          <span>{{ t('module_reply_title') }}</span>
        </div>
        <a-switch v-model:checked="ruleForm.reply_comment_switch" :checkedValue="1" :unCheckedValue="0" />
      </div>
      <div v-if="ruleForm.reply_comment_switch==1" class="module-card-body">
        <div class="field">
          <div class="field-label">{{ t('auto_reply_rule_label') }}</div>
          <a-textarea
            v-model:value="ruleForm.reply_check_prompt"
            :placeholder="t('common_input_placeholder')"
            :auto-size="{ minRows: 3 }"
          />
        </div>
        <div class="field">
          <div class="field-label">{{ t('reply_content_label') }}</div>
          <a-radio-group v-model:value="ruleForm.reply_type">
            <a-radio :value="1">{{ t('reply_type_fixed') }}</a-radio>
            <a-radio :value="2">{{ t('reply_type_ai') }}</a-radio>
          </a-radio-group>
        </div>
        <div class="field">
          <div class="field-label">{{ replyContentLabel }}</div>
          <a-textarea
            v-model:value="ruleForm.reply_prompt"
            :placeholder="t('common_input_placeholder')"
            :auto-size="{ minRows: 3 }"
          />
        </div>
      </div>
    </div>

    <div class="module-card">
      <div class="module-card-header">
        <div class="module-card-title">
          <svg-icon name="comment-selected" style="font-size: 16px; margin-right: 8px;" />
          <span>{{ t('module_elect_title') }}</span>
        </div>
        <a-switch v-model:checked="ruleForm.elect_comment_switch" :checkedValue="1" :unCheckedValue="0" />
      </div>
      <div v-if="ruleForm.elect_comment_switch==1" class="module-card-body">
        <div class="field">
          <div class="field-label">{{ t('auto_elect_rule_label') }}</div>
          <a-textarea
            v-model:value="ruleForm.elect_prompt"
            :placeholder="t('common_input_placeholder')"
            :auto-size="{ minRows: 3 }"
          />
        </div>
      </div>
    </div>

    <div class="page-right-footer">
      <a-button type="primary" @click="onSave">{{ t('btn_save_and_apply') }}</a-button>
    </div>

    <a-modal
      v-model:open="taskModalOpen"
      :title="t('task_modal_title')"
      :width="720"
      @ok="confirmSelectTasks"
      @cancel="taskModalOpen=false"
    >
      <div class="task-modal">
        <a-table
          :columns="taskColumns"
          :data-source="taskList"
          :row-key="record => String(record.id)"
          :loading="taskLoading"
          :pagination="{ current: taskPager.page, pageSize: taskPager.size, total: taskPager.total, showSizeChanger: true, pageSizeOptions: ['10','20','50','100'] }"
          @change="onTaskTableChange"
          :row-selection="{selectedRowKeys: taskSelectedRowKeys, onChange: onSelectChange, preserveSelectedRowKeys: true}"
        >
          <template #headerCell="{ column }">
            <span v-if="typeof column.title === 'string'">{{ column.title }}</span>
          </template>
        </a-table>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
// import { ExclamationCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import { getCommentRuleInfo, getBatchSendTaskList, saveCommentRule } from '@/api/robot'
import { useI18n } from '@/hooks/web/useI18n'

const route = useRoute()
const router = useRouter()
const { t } = useI18n('views.explore.ai-comment-management.create-custom-rule')
// const isCreateCustomRule = computed(() => route.name === 'exploreAiCommentManagementCreateCustomRule')
const isCopyCustomRule = computed(() => route.name === 'exploreAiCommentManagementCreateCustomRule' && route.query.copy === '1')
const isEditCustomRule = computed(() => route.name === 'exploreAiCommentManagementCreateCustomRule' && route.query.id)

const ruleForm = reactive({
  rule_name: '',
  use_model: '',
  model_config_id: '',
  delete_comment_switch: 1,
  delete_type: [1, 2],
  delete_keywords: [],
  delete_condition: 2,
  delete_priority: 1,
  delete_prompt: '',
  reply_comment_switch: 1,
  reply_check_prompt: '',
  reply_type: 1,
  reply_prompt: '',
  elect_comment_switch: 1,
  elect_prompt: '',
})

const replyContentLabel = computed(() =>
  ruleForm.reply_type === 1 ? t('reply_content_label') : t('reply_content_rule_label')
)

// const taskScope = ref('all')
// const selectedTaskIds = ref([])
// const selectedTasks = ref([])

// const selectedTaskColumns = [
//   { title: '群发名称', dataIndex: 'task_name', key: 'task_name' },
//   { title: '操作', dataIndex: 'actions', key: 'actions' },
// ]

const deleteKeywordInput = ref('')

const onModelChange = (val, option) => {
  ruleForm.use_model = option.modelName
  ruleForm.model_config_id = option.modelId
}

const removeKeyword = (idx) => {
  ruleForm.delete_keywords.splice(idx, 1)
}
const onAddKeyword = () => {
  const text = (deleteKeywordInput.value || '').trim()
  if (!text) return
  if (!ruleForm.delete_keywords.includes(text)) {
    ruleForm.delete_keywords.push(text)
  }
  deleteKeywordInput.value = ''
}

const taskModalOpen = ref(false)
const taskList = ref([])
const taskLoading = ref(false)
const taskPager = reactive({ page: 1, size: 10, total: 0 })
const taskSelectedRowKeys = ref([])

const taskColumns = [
  { title: t('task_column_task_name'), dataIndex: 'task_name', key: 'task_name' },
]

const onSelectChange = (keys) => {
  taskSelectedRowKeys.value = keys.map((k) => String(k))
}

// const openTaskModal = () => {
//   taskModalOpen.value = true
//   taskSelectedRowKeys.value = [...selectedTaskIds.value].map((k) => String(k))
//   loadTaskList()
// }

const loadTaskList = () => {
  taskLoading.value = true
  getBatchSendTaskList({ page: taskPager.page, size: taskPager.size }).then((res) => {
    const data = res?.data || { list: [], total: 0 }
    taskList.value = data.list || []
    taskPager.total = +data.total || 0
  }).finally(() => { taskLoading.value = false })
}

// const loadSelectedTasksForIds = () => {
//   const ids = selectedTaskIds.value.slice()
//   if (!ids.length) { selectedTasks.value = []; return }
//   taskLoading.value = true
//   getBatchSendTaskList({ page: 1, size: 1000 }).then((res) => {
//     const data = res?.data || { list: [], total: 0 }
//     const all = data.list || []
//     const setIds = new Set(ids.map((x) => String(x)))
//     selectedTasks.value = all.filter((it) => setIds.has(String(it.id)))
//   }).finally(() => { taskLoading.value = false })
// }

const onTaskTableChange = (pagination) => {
  taskPager.page = pagination.current
  taskPager.size = pagination.pageSize
  loadTaskList()
}

// const confirmSelectTasks = () => {
//   selectedTaskIds.value = [...taskSelectedRowKeys.value].map((k) => String(k))
//   const idSet = new Set(selectedTaskIds.value)
//   selectedTasks.value = (taskList.value || []).filter(it => idSet.has(String(it.id)))
//   taskModalOpen.value = false
// }

// const removeSelectedTask = (id) => {
//   Modal.confirm({
//     title: '提示',
//     icon: createVNode(ExclamationCircleOutlined),
//     content: '确认删除选中的群发任务吗？',
//     onOk: () => {
//       selectedTaskIds.value = selectedTaskIds.value.filter(x => String(x) !== String(id))
//       selectedTasks.value = selectedTasks.value.filter(it => String(it.id) !== String(id))
//     },
//     onCancel: () => {},
//   })
// }

const loadDefaultRule = () => {
  const params = {}
  if (isCopyCustomRule.value || isEditCustomRule.value) {
    params.id = route.query.id
  }
  getCommentRuleInfo(params).then((res) => {
    const item = res?.data || null
    if (!item) return
    ruleForm.rule_name = item.rule_name || t('rule_name_default')
    ruleForm.use_model = item.use_model || ruleForm.use_model
    ruleForm.model_config_id = item.model_config_id || ruleForm.model_config_id
    ruleForm.delete_comment_switch = +item.delete_comment_switch || 0
    ruleForm.reply_comment_switch = +item.reply_comment_switch || 0
    ruleForm.elect_comment_switch = +item.elect_comment_switch || 0
    // if (item.task_ids) {
    //   const ids = String(item.task_ids).split(',').map(s => s.trim()).filter(Boolean)
    //   if (ids.length === 0 || ids.includes('0')) {
    //     taskScope.value = 'all'
    //   } else {
    //     taskScope.value = 'specific'
    //     selectedTaskIds.value = ids.map((x) => String(x))
    //     loadSelectedTasksForIds()
    //   }
    // }
    if (item.delete_comment_rule) {
      let r = item.delete_comment_rule
      if (typeof r === 'string') { try { r = JSON.parse(r) } catch {} }
      ruleForm.delete_type = Array.isArray(r.type) ? r.type : []
      ruleForm.delete_keywords = Array.isArray(r.keywords) ? r.keywords : []
      ruleForm.delete_condition = +r.condition || ruleForm.delete_condition
      ruleForm.delete_priority = +r.priority || ruleForm.delete_priority
      ruleForm.delete_prompt = r.prompt || ''
    }
    if (item.reply_comment_rule) {
      let r = item.reply_comment_rule
      if (typeof r === 'string') { try { r = JSON.parse(r) } catch {} }
      ruleForm.reply_check_prompt = r.check_reply_prompt || ''
      ruleForm.reply_type = +r.reply_type || ruleForm.reply_type
      ruleForm.reply_prompt = r.reply_prompt || ''
    }
    if (item.elect_comment_rule) {
      let r = item.elect_comment_rule
      if (typeof r === 'string') { try { r = JSON.parse(r) } catch {} }
      ruleForm.elect_prompt = r.prompt || ''
    }
  })
}

const onSave = () => {
  if (!String(ruleForm.rule_name || '').trim()) {
    return message.error(t('error_enter_rule_name'))
  }
  if (!ruleForm.use_model) {
    return message.error(t('error_select_model'))
  }
  // if (taskScope.value === 'specific' && selectedTaskIds.value.length === 0) {
  //   return message.error('请选择群发任务')
  // }
  if (ruleForm.delete_comment_switch == 1) {
    const types = Array.isArray(ruleForm.delete_type) ? ruleForm.delete_type : []
    const hasSensitive = types.includes(1)
    const hasDetection = types.includes(2)
    if (types.length === 0) {
      return message.error(t('error_delete_at_least_one_option'))
    }
    if (hasSensitive && (!ruleForm.delete_keywords || ruleForm.delete_keywords.length === 0)) {
      return message.error(t('error_enter_sensitive_words'))
    }
    if (hasDetection && !String(ruleForm.delete_prompt || '').trim()) {
      return message.error(t('error_enter_ai_detect_rule'))
    }
  }
  if (ruleForm.reply_comment_switch == 1) {
    if (!String(ruleForm.reply_check_prompt || '').trim()) {
      return message.error(t('error_enter_auto_reply_rule'))
    }
    if (!String(ruleForm.reply_prompt || '').trim()) {
      return message.error(t('error_enter_reply_content', { label: replyContentLabel.value }))
    }
  }
  if (ruleForm.elect_comment_switch == 1) {
    if (!String(ruleForm.elect_prompt || '').trim()) {
      return message.error(t('error_enter_auto_elect_rule'))
    }
  }
  // const taskIds = taskScope.value === 'all' ? 0 : selectedTaskIds.value.join(',')
  const payload = {
    rule_name: ruleForm.rule_name,
    use_model: ruleForm.use_model,
    model_config_id: ruleForm.model_config_id,
    // task_ids: taskIds,
    delete_comment_switch: ruleForm.delete_comment_switch,
    delete_comment_rule: JSON.stringify((() => {
      const types = Array.isArray(ruleForm.delete_type) ? ruleForm.delete_type : []
      const hasDetection = types.includes(2)
      const rule = {
        type: types,
        keywords: ruleForm.delete_keywords,
        prompt: ruleForm.delete_prompt,
      }
      if (hasDetection) {
        rule.condition = ruleForm.delete_condition
        rule.priority = ruleForm.delete_priority
      }
      return rule
    })()),
    reply_comment_switch: ruleForm.reply_comment_switch,
    reply_comment_rule: JSON.stringify({
      check_reply_prompt: ruleForm.reply_check_prompt,
      reply_type: ruleForm.reply_type,
      reply_prompt: ruleForm.reply_prompt,
    }),
    elect_comment_switch: ruleForm.elect_comment_switch,
    elect_comment_rule: JSON.stringify({
      prompt: ruleForm.elect_prompt,
    }),
  }
  if (isEditCustomRule.value && !isCopyCustomRule.value) {
    payload.id = route.query.id
  }
  saveCommentRule(payload).then(() => {
    message.success(t('save_success'))
  })
  // 返回
  router.back()
}

onMounted(() => {
  if (isEditCustomRule.value || isCopyCustomRule.value) {
    loadDefaultRule()
  }
})
</script>

<style lang="less" scoped>
.default-rule-page {
  width: 824px;
  padding: 2px 0px 24px;
}
.form-row {
  display: flex;
  flex-direction: column;
  align-items: self-start;
  gap: 4px;
  margin-bottom: 16px;
}
.radio-button {
  margin: -8px 0 8px 0px;
  display: block;
}
.form-label {
  width: 80px;
  color: #262626;
}
.selected-task-table {
  margin-bottom: 16px;
}
.module-card {
  margin-top: 16px;
  border: 1px solid #e6e6e6;
  border-radius: 8px;
  background: #f7f9fc;
}
.module-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  font-weight: 600;
  color: #262626;

  .module-card-title {
    font-size: 16px;
  }
}
.module-card-body {
  padding: 0px 16px 16px 16px;
}
.field {
  margin-bottom: 16px;
}
.field-label {
  margin-bottom: 4px;
  color: #262626;
}
.tag-list {
  margin-bottom: 8px;
}
.page-right-footer {
  position: fixed;
  bottom: 0;
  right: 16px;
  display: flex;
  width: 100%;
  padding: 16px 1055px 16px 64px;
  align-items: center;
  border-radius: 0 0 2px 2px;
  background: #FFF;
  box-shadow: 0 -8px 4px 0 #0000000a;
}
</style>
