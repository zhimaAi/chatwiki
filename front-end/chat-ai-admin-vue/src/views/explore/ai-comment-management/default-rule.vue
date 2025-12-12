<template>
  <div class="default-rule-page">
    <div class="form-row">
      <div class="form-label">规则名称</div>
      <a-input v-model:value="ruleForm.rule_name" placeholder="请输入规则名称" style="width: 300px" disabled />
    </div>

    <div class="form-row">
      <div class="form-label">模型</div>
      <ModelSelect
        modelType="LLM"
        v-model:modeName="ruleForm.use_model"
        v-model:modeId="ruleForm.model_config_id"
        style="width: 300px;"
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
          <span>AI自动删评</span>
        </div>
        <a-switch v-model:checked="ruleForm.delete_comment_switch" :checkedValue="1" :unCheckedValue="0" />
      </div>
      <div v-if="ruleForm.delete_comment_switch==1" class="module-card-body">
        <div class="field">
          <div class="field-label">规则</div>
          <a-checkbox-group v-model:value="ruleForm.delete_type">
            <a-checkbox :value="1">触发敏感词</a-checkbox>
            <a-checkbox :value="2">AI检测</a-checkbox>
          </a-checkbox-group>
        </div>
        <div class="field" v-if="(ruleForm.delete_type || []).includes(1)">
          <div class="field-label">敏感词</div>
          <div class="tag-list">
            <a-tag v-for="(kw, idx) in ruleForm.delete_keywords" :key="kw" closable @close="removeKeyword(idx)">{{ kw }}</a-tag>
          </div>
          <a-textarea v-model:value="deleteKeywordInput" placeholder="请输入" @blur="onAddKeyword" :auto-size="{ minRows: 2 }" />
        </div>
        <div class="field" v-if="(ruleForm.delete_type || []).includes(2)">
          <div class="field-label">AI检测规则</div>
          <a-textarea v-model:value="ruleForm.delete_prompt" placeholder="请输入" :auto-size="{ minRows: 3 }" />
        </div>
        <div class="field" v-if="(ruleForm.delete_type || []).includes(2)">
          <div class="field-label">条件</div>
          <a-radio-group v-model:value="ruleForm.delete_condition">
            <a-radio :value="1">全部满足</a-radio>
            <a-radio :value="2">满足其中一个条件</a-radio>
          </a-radio-group>
        </div>
        <div class="field" v-if="(ruleForm.delete_type || []).includes(2)">
          <div class="field-label">优先级</div>
          <a-radio-group v-model:value="ruleForm.delete_priority">
            <a-radio :value="1">敏感词优先</a-radio>
            <a-radio :value="2">AI检测优先</a-radio>
          </a-radio-group>
        </div>
      </div>
    </div>

    <div class="module-card">
      <div class="module-card-header">
        <div class="module-card-title">
          <svg-icon name="comment-reply" style="font-size: 16px; margin-right: 8px;" />
          <span>AI自动回复</span>
        </div>
        <a-switch v-model:checked="ruleForm.reply_comment_switch" :checkedValue="1" :unCheckedValue="0" />
      </div>
      <div v-if="ruleForm.reply_comment_switch==1" class="module-card-body">
        <div class="field">
          <div class="field-label">自动回复规则</div>
          <a-textarea v-model:value="ruleForm.reply_check_prompt" placeholder="请输入" :auto-size="{ minRows: 3 }" />
        </div>
        <div class="field">
          <div class="field-label">回复内容</div>
          <a-radio-group v-model:value="ruleForm.reply_type">
            <a-radio :value="1">固定回复内容</a-radio>
            <a-radio :value="2">AI生成回复内容</a-radio>
          </a-radio-group>
        </div>
        <div class="field">
          <div class="field-label">{{ replyContentLabel }}</div>
          <a-textarea v-model:value="ruleForm.reply_prompt" placeholder="请输入" :auto-size="{ minRows: 3 }"/>
        </div>
      </div>
    </div>

    <div class="module-card">
      <div class="module-card-header">
        <div class="module-card-title">
          <svg-icon name="comment-selected" style="font-size: 16px; margin-right: 8px;" />
          <span>AI自动精选</span>
        </div>
        <a-switch v-model:checked="ruleForm.elect_comment_switch" :checkedValue="1" :unCheckedValue="0" />
      </div>
      <div v-if="ruleForm.elect_comment_switch==1" class="module-card-body">
        <div class="field">
          <div class="field-label">自动精选规则</div>
          <a-textarea v-model:value="ruleForm.elect_prompt" placeholder="请输入" :auto-size="{ minRows: 3 }" />
        </div>
      </div>
    </div>

    <div class="page-right-footer">
      <a-button type="primary" @click="onSave">保存并应用</a-button>
    </div>

    <a-modal v-model:open="taskModalOpen" title="选择群发任务" :width="720" @ok="confirmSelectTasks" @cancel="taskModalOpen=false">
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
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { message } from 'ant-design-vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import { getCommentRuleList, getBatchSendTaskList, saveCommentRule } from '@/api/robot'

const delete_prompt_default = "审核规则（请严格遵循）：\n\n"+
"违法违规： 包含明显违法、煽动暴力、恐怖、极端主义，或传播非法信息的内容。\n\n"+
"仇恨与攻击： 针对种族、民族、宗教、性别、性取向、残疾等个人或群体的侮辱、歧视、人身攻击或恶意嘲讽。\n\n"+
"严重不雅与污秽： 包含极度污言秽语、色情露骨描述或令人极度不适的暴力细节。\n\n"+
"垃圾广告与导流： 明显的、与讨论主题无关的商业广告、重复刷屏、推广联系方式或恶意引流链接。\n\n"+
"严重误导与欺诈： 传播已被官方证伪的、可能造成公共危害的虚假信息（如健康谣言、金融骗局），或试图进行欺诈。\n\n"+
"侵犯隐私： 公开他人的个人信息（电话、住址、身份证号等）。\n\n"+
"恶意破坏： 纯粹无意义的字符、重复内容，或旨在干扰社区正常讨论的恶意刷屏。"

const reply_check_prompt_default = "特征（满足以下所有条件）：\n\n"+
"意图明确：问题或需求清晰可识别\n\n"+
"问题标准化：属于常见问题或标准业务流程\n\n"+
"情绪稳定：用户情绪中立或积极（无明显愤怒/失望）\n\n"+
"信息完整：提供了足够判断信息\n\n"+
"解决方案已知：有预设的解决方案或知识库条目"

const reply_prompt_default = "解决问题导向： 针对问题提出具体的解决方案或下一步行动。\n\n"+
"保持专业与友善： 语气礼貌、积极，即使面对负面评论也保持冷静和专业。\n\n"+
"呼吁行动： 在适当时，引导客户进行下一步（如私信、联系客服、再次购买等）。\n\n"+
"符合品牌人设： 语言风格应与品牌形象一致（如：亲切的、专业的、活泼的、高端的)"

const reply_fixed_default = "感谢您的评论，您的建议是我们进步的基石，祝愿您开心每一天"

const elect_prompt_default = "精选标准\n\n"+
"好评： 包含具体、细节的使用体验、效果对比、场景描述，能为其他潜在客户提供真实参考。\n\n"+
"中评建议： 客观、理性地指出产品或服务的具体问题，并提出了可行的改进建议。这类评论对品牌成长极具价值。"

const ruleForm = reactive({
  id: '',
  rule_name: '默认规则',
  use_model: '',
  model_config_id: '',
  delete_comment_switch: 1,
  delete_type: [1, 2],
  delete_keywords: [],
  delete_condition: 2,
  delete_priority: 1,
  delete_prompt: delete_prompt_default,
  reply_comment_switch: 1,
  reply_check_prompt: reply_check_prompt_default,
  reply_type: 1,
  reply_prompt: '',
  elect_comment_switch: 1,
  elect_prompt: '',
})

const replyContentLabel = computed(() => ruleForm.reply_type === 1 ? '回复内容' : '回复内容规则')

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
  { title: '群发名称', dataIndex: 'task_name', key: 'task_name' },
]

const onSelectChange = (keys) => {
  taskSelectedRowKeys.value = keys.map((k) => String(k))
}

// const openTaskModal = () => {
//   taskSelectedRowKeys.value = [...selectedTaskIds.value].map((k) => String(k))
//   taskModalOpen.value = true
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
//       selectedTaskIds.value = selectedTaskIds.value.filter(x => x !== id)
//       selectedTasks.value = selectedTasks.value.filter(it => it.id !== id)
//     },
//     onCancel: () => {},
//   })
// }

const loadDefaultRule = () => {
  getCommentRuleList({ page: 1, size: 1, is_default: 1 }).then((res) => {
    const item = res?.data?.list?.[0] || (Array.isArray(res?.data) ? res.data[0] : res?.data) || null
    if (!item) return
    ruleForm.id = item.id || ''
    ruleForm.rule_name = item.rule_name || '默认规则'
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
      ruleForm.delete_prompt = r.prompt || delete_prompt_default
    }
    if (item.reply_comment_rule) {
      let r = item.reply_comment_rule
      if (typeof r === 'string') { try { r = JSON.parse(r) } catch {} }
      ruleForm.reply_check_prompt = r.check_reply_prompt || reply_check_prompt_default
      ruleForm.reply_type = +r.reply_type || ruleForm.reply_type
      ruleForm.reply_prompt = r.reply_prompt
      if(!ruleForm.reply_prompt) {
        if (ruleForm.reply_type == 1) {
          ruleForm.reply_prompt = reply_fixed_default
        } else {
          ruleForm.reply_prompt = reply_prompt_default
        }
      }
    }
    if (item.elect_comment_rule) {
      let r = item.elect_comment_rule
      if (typeof r === 'string') { try { r = JSON.parse(r) } catch {} }
      ruleForm.elect_prompt = r.prompt || elect_prompt_default
    }
  })
}

const onSave = () => {
  if (!String(ruleForm.rule_name || '').trim()) {
    return message.error('请输入规则名称')
  }
  if (!ruleForm.use_model) {
    return message.error('请选择模型')
  }
  // if (taskScope.value === 'specific' && selectedTaskIds.value.length === 0) {
  //   return message.error('请选择群发任务')
  // }
  if (ruleForm.delete_comment_switch == 1) {
    const types = Array.isArray(ruleForm.delete_type) ? ruleForm.delete_type : []
    const hasSensitive = types.includes(1)
    const hasDetection = types.includes(2)
    if (types.length === 0) {
      return message.error('AI自动删评，请至少勾选一个选项')
    }
    if (hasSensitive && (!ruleForm.delete_keywords || ruleForm.delete_keywords.length === 0)) {
      return message.error('请输入敏感词')
    }
    if (hasDetection && !String(ruleForm.delete_prompt || '').trim()) {
      return message.error('请输入AI检测规则')
    }
  }
  if (ruleForm.reply_comment_switch == 1) {
    if (!String(ruleForm.reply_check_prompt || reply_check_prompt_default)) {
      return message.error('请输入自动回复规则')
    }
    if (!String(ruleForm.reply_prompt)) {
      if (ruleForm.reply_type == 1) {
        ruleForm.reply_prompt = reply_fixed_default
      } else {
        ruleForm.reply_prompt = reply_prompt_default
      }
    }
  }
  if (ruleForm.elect_comment_switch == 1) {
    if (!String(ruleForm.elect_prompt || elect_prompt_default)) {
      return message.error('请输入自动精选规则')
    }
  }
  // const taskIds = taskScope.value === 'all' ? 0 : selectedTaskIds.value.join(',')
  const payload = {
    id: ruleForm.id,
    is_default: 1,
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
  saveCommentRule(payload).then(() => { message.success('保存成功') })
}

onMounted(() => { loadDefaultRule() })

watch(() => ruleForm.reply_type, (val) => {
  const cur = String(ruleForm.reply_prompt || '')
  if (val === 1) {
    if (!cur.trim() || cur === reply_prompt_default) {
      ruleForm.reply_prompt = reply_fixed_default
    }
  } else {
    if (!cur.trim() || cur === reply_fixed_default) {
      ruleForm.reply_prompt = reply_prompt_default
    }
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
