<template>
  <div class="subManage-edit">
    <a-breadcrumb class="subManage-breadcrumb">
      <a-breadcrumb-item><a :href="autoReplyUrl">关键词回复</a></a-breadcrumb-item>
      <a-breadcrumb-item>新增规则</a-breadcrumb-item>
    </a-breadcrumb>

    <div class="main">
      <a-form ref="formRef" :model="form" :rules="rules">
        <a-form-item label="规则名称" name="name" :rules="[{ required: true, message: '请输入规则名称' }]">
          <a-input v-model:value="form.name" placeholder="请输入" style="width: 512px;" />
        </a-form-item>

        <div class="nav-box" style="padding-top: 10px;">
          <svg-icon name="set-keywords" style="font-size: 16px;"></svg-icon>
          设置关键词
        </div>

        <!-- 精准匹配 -->
        <div class="item-box">
          <div class="item-title-box">
            <div class="item-title">精准匹配</div>
            <a-tooltip>
              <template #title>发送消息内容只提到某个关键词，才回复指定的内容</template>
              <QuestionCircleOutlined style="color: #8C8C8C;" />
            </a-tooltip>
          </div>
          <div class="flex" style="gap:8px;align-items:center;">
            <template v-if="addingFull">
              <a-input
                ref="fullInputRef"
                v-model:value="fullKeywordInput"
                placeholder="输入关键词后回车或失去焦点添加"
                style="width:260px;"
                @pressEnter="confirmAddFull"
                @blur="confirmAddFull"
              />
            </template>
            <template v-else>
              <a-button type="dashed" class="add-btn" @click="startAddFull">
                <template #icon>
                  <PlusOutlined />
                </template>
                添加关键词
              </a-button>
            </template>
          </div>
          <div class="tag-container" style="margin-top:8px;">
            <a-tag v-for="item in fullKeywords" :key="item" closable @close="removeFullKeyword(item)">
              {{ item }}
            </a-tag>
          </div>
        </div>

        <!-- 模糊匹配 -->
        <div class="item-box" style="margin-top: 12px;">
          <div class="item-title-box">
            <div class="item-title">模糊匹配</div>
            <a-tooltip>
              <template #title>发送消息内容只要包含某个关键词，才回复指定的内容</template>
              <QuestionCircleOutlined style="color: #8C8C8C;" />
            </a-tooltip>
          </div>
          <div class="flex" style="gap:8px;align-items:center;">
            <template v-if="addingHalf">
              <a-input
                ref="halfInputRef"
                v-model:value="halfKeywordInput"
                placeholder="输入关键词后回车或失去焦点添加"
                style="width:260px;"
                @pressEnter="confirmAddHalf"
                @blur="confirmAddHalf"
              />
            </template>
            <template v-else>
              <a-button type="dashed" class="add-btn" @click="startAddHalf">
                <template #icon>
                  <PlusOutlined />
                </template>
                添加关键词
              </a-button>
            </template>
          </div>
          <div class="tag-container" style="margin-top:8px;">
            <a-tag v-for="item in halfKeywords" :key="item" closable @close="removeHalfKeyword(item)">
              {{ item }}
            </a-tag>
          </div>
        </div>

        <!-- 回复内容 -->
        <div class="nav-box" style="margin-top: 24px;">
          <svg-icon name="reply-content" style="font-size: 16px;"></svg-icon>
          回复内容
        </div>
        <div class="item-box">
          <MultiReply v-for="(it, idx) in replyList" :key="idx" v-model:value="replyList[idx]" :reply_index="idx"
            @change="onContentChange" @del="onDelItem" />
          <a-button type="dashed" style="width: 694px;" :disabled="replyList.length >= 5" @click="addReplyItem">
            <template #icon>
              <PlusOutlined />
            </template>
            添加回复内容({{replyList.length}}/5)
          </a-button>
        </div>

        <!-- 回复设置 -->
        <div class="nav-box" style="margin-top: 24px;">
          <svg-icon name="reply-settings" style="font-size: 16px;"></svg-icon>
          回复设置
        </div>
        <div class="item-box">
          <div class="item-title-box">
            <div class="item-title">回复方式：</div>
          </div>
          <div class="radio-container">
            <a-radio-group v-model:value="form.reply_num" @change="handleReplyTypeChange">
              <a-radio value="0">全部回复</a-radio>
              <a-radio value="1">随机回复一条</a-radio>
            </a-radio-group>
          </div>
        </div>

        <!-- 保存 底部固定 -->
        <div class="btn-container">
          <a-button type="primary" @click="onSubmit">保存</a-button>
        </div>
      </a-form>
    </div>
  </div>
</template>
<script setup>
import { QuestionCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { onMounted, ref, reactive, computed, nextTick, toRaw } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import MultiReply from '@/components/replay-card/multi-reply.vue'
import { getRobotKeywordReply, saveRobotKeywordReply, checkKeyWordRepeat } from '@/api/explore/index.js'

const query = useRoute().query
const ruleId = ref(+query.rule_id || +query['rule-id'] || 0)
const router = useRouter()
const autoReplyUrl = computed(() => `/#/robot/ability/auto-reply?id=${query.id}&robot_key=${query.robot_key}`)
const formRef = ref(null)
const form = reactive({
  name: '',
  reply_num: '0'
})
const rules = {
  name: [
    { required: true, message: '请输入规则名称', trigger: 'blur' }
  ]
}

// 关键词集合
const fullKeywords = ref([])
const halfKeywords = ref([])
const fullKeywordInput = ref('')
const halfKeywordInput = ref('')
const addingFull = ref(false)
const addingHalf = ref(false)
const fullInputRef = ref(null)
const halfInputRef = ref(null)
const addingFullBusy = ref(false)
const startAddFull = () => {
  addingFull.value = true
  nextTick(() => {
    const comp = fullInputRef.value
    if (comp?.focus) comp.focus()
    else comp?.$el?.querySelector('input')?.focus()
  })
}
const confirmAddFull = async () => {
  if (addingFullBusy.value) return
  const v = (fullKeywordInput.value || '').trim()
  if (!v) { message.warning('请输入关键词'); return }
  if (fullKeywords.value.includes(v)) { message.warning('已存在该关键词'); return }
  try {
    addingFullBusy.value = true
    const res = await checkKeyWordRepeat({ id: ruleId.value, robot_id: query.id, keyword: v })
    const repeat = res?.data?.is_repeat
    const ruleName = res?.data?.rule_name || ''
    if (repeat) {
      message.error(`关键词与规则「${ruleName}」重复`)
      return
    }
    fullKeywords.value.push(v)
    fullKeywordInput.value = ''
    addingFull.value = false
  } catch (e) {
    message.error('校验失败，请稍后重试')
  } finally {
    addingFullBusy.value = false
  }
}
const removeFullKeyword = (k) => {
  fullKeywords.value = fullKeywords.value.filter((x) => x !== k)
}
const addingHalfBusy = ref(false)
const startAddHalf = () => {
  addingHalf.value = true
  nextTick(() => {
    const comp = halfInputRef.value
    if (comp?.focus) comp.focus()
    else comp?.$el?.querySelector('input')?.focus()
  })
}
const confirmAddHalf = async () => {
  if (addingHalfBusy.value) return
  const v = (halfKeywordInput.value || '').trim()
  if (!v) { message.warning('请输入关键词'); return }
  if (halfKeywords.value.includes(v)) { message.warning('已存在该关键词'); return }
  try {
    addingHalfBusy.value = true
    const res = await checkKeyWordRepeat({ id: ruleId.value, robot_id: query.id, keyword: v })
    const repeat = res?.data?.is_repeat
    const ruleName = res?.data?.rule_name || ''
    if (repeat) {
      message.error(`关键词与规则「${ruleName}」重复`)
      return
    }
    halfKeywords.value.push(v)
    halfKeywordInput.value = ''
    addingHalf.value = false
  } catch (e) {
    message.error('校验失败，请稍后重试')
  } finally {
    addingHalfBusy.value = false
  }
}
const removeHalfKeyword = (k) => {
  halfKeywords.value = halfKeywords.value.filter((x) => x !== k)
}

// 回复内容列表
const replyList = ref([{ type: 'text', description: '' }])
const addReplyItem = () => {
  if (replyList.value.length >= 5) {
    message.warning('最多添加5条回复内容')
    return
  }
  replyList.value.push({ type: 'text', description: '' })
}
const onContentChange = (payload) => {
  const { reply_index, ...rest } = payload
  if (reply_index >= 0 && reply_index < replyList.value.length) {
    replyList.value[reply_index] = rest
  }
}
const onDelItem = (index) => {
  replyList.value.splice(index, 1)
}

function mapApiTypeToTabType (t) {
  switch (t) {
    case 'text': return 1
    case 'image': return 2
    case 'card': return 6
    case 'imageText': return 1 // 兜底为文本
    default: return 1
  }
}

const handleReplyTypeChange = () => {

}

function mapApiReplyToItem (rc) {
  const typeStr = rc?.type || rc?.reply_type || 'text'
  const tabType = mapApiTypeToTabType(typeStr)
  let content = ''
  if (tabType === 1) content = rc?.description || ''
  else if (tabType === 2) content = rc?.thumb_url || rc?.pic || ''
  else if (tabType === 6) content = ''
  return {
    id: rc?.id || '',
    content,
    type: tabType,
    type_main_id: rc?.type_main_id || '',
    extra: rc?.extra || {},
    err_msg: rc?.err_msg || {},
    delay_time: rc?.delay_time || '',
    key: rc?.key || ''
  }
}

function serializeReplyContent (list) {
  return (list || []).map((it) => ({ ...it, status: '1' }))
}

function serializeReplyTypeCodes (list) {
  const map = { text: '2', image: '4', card: '3', imageText: '1', url: '5' }
  return list.map((it) => map[it.type] || '').filter(Boolean)
}

const onSubmit = () => {
  formRef.value?.validate().then(async () => {
    if (!fullKeywords.value.length && !halfKeywords.value.length) {
      message.warning('请至少添加一个关键词（精准或模糊）')
      return
    }
    if (!replyList.value.length) {
      message.warning('请至少添加一条回复内容')
      return
    }
    const payload = {
      robot_id: query.id,
      name: form.name,
      full_keyword: fullKeywords.value,
      half_keyword: halfKeywords.value,
      reply_content: JSON.stringify(serializeReplyContent(replyList.value)),
      reply_type: serializeReplyTypeCodes(replyList.value),
      reply_num: form.reply_num
    }
    if (ruleId.value) payload.id = ruleId.value
    try {
      const res = await saveRobotKeywordReply(payload)
      if (res && res.res == 0) {
        message.success('保存成功')
        router.push({ path: '/robot/ability/auto-reply', query: { id: query.id, robot_key: query.robot_key } })
      }
    } catch (e) {
      // message.error('保存失败，请稍后重试')
    }
  })
}


onMounted(async () => {
  const copyId = +(query.copy_id || 0)
  if (!ruleId.value && copyId) {
    try {
      const res = await getRobotKeywordReply({ id: copyId })
      const data = res?.data || {}
      form.name = `${(data?.name || '')}副本`
      // fullKeywords.value = Array.isArray(data?.full_keyword) ? data.full_keyword : []
      // halfKeywords.value = Array.isArray(data?.half_keyword) ? data.half_keyword : []
      const list = Array.isArray(data?.reply_content) ? data.reply_content : []
      replyList.value = list.map((rc) => ({
        type: (rc?.type || rc?.reply_type || 'text'),
        description: rc?.description || '',
        thumb_url: rc?.thumb_url || rc?.pic || '',
        title: rc?.title || '',
        url: rc?.url || '',
        appid: rc?.appid || '',
        page_path: rc?.page_path || ''
      }))
      form.reply_num = data.reply_num
    } catch (e) {
      message.error('加载规则失败，请稍后重试')
    }
    return
  }
  if (!ruleId.value) return
  try {
    const res = await getRobotKeywordReply({ id: ruleId.value })
    const data = res?.data || {}
    form.name = data?.name || ''
    fullKeywords.value = Array.isArray(data?.full_keyword) ? data.full_keyword : []
    halfKeywords.value = Array.isArray(data?.half_keyword) ? data.half_keyword : []
    const list = Array.isArray(data?.reply_content) ? data.reply_content : []
    replyList.value = list.map((rc) => ({
      type: (rc?.type || rc?.reply_type || 'text'),
      description: rc?.description || '',
      thumb_url: rc?.thumb_url || rc?.pic || '',
      title: rc?.title || '',
      url: rc?.url || '',
      appid: rc?.appid || '',
      page_path: rc?.page_path || ''
    }))
    form.reply_num = data.reply_num
  } catch (e) {
    message.error('加载规则失败，请稍后重试')
  }
})

</script>
<style lang="less" scoped>
.subManage-edit {
  padding: 16px 24px;
  width: 100%;
  height: 100%;
  border-bottom: 1px solid #fff;
  border-right: 1px solid #fff;
  background-color: #fff;
  overflow-x: hidden;
  overflow-y: auto;

  .subManage-breadcrumb {
    display: flex;
    align-items: center;
    color: #000000;
    font-family: "PingFang SC";
    font-size: 14px;
    font-style: normal;
    line-height: 22px;
    padding-bottom: 16px;
  }

  .main {
    padding: 0 8px;
    border-radius: 6px;
    background-color: white;
    padding-bottom: 24px;

    .title {
      border-radius: 6px;
      padding: 12px 0 12px 24px;
      align-items: flex-start;
      border-bottom: 1px solid var(--07, #F0F0F0);
      background: #FFF;
      display: flex;
      align-items: center;
      color: #262626;
      font-family: "PingFang SC";
      font-size: 14px;
      font-style: normal;
      font-weight: 600;
      line-height: 22px;
      gap: 8px;
      margin-bottom: 24px;
    }
  }
}

.mr-8 {
  margin-right: 8px;
}

.mr16 {
  margin-right: 16px;
}

.nav-box {
  color: #262626;
  font-size: 14px;
  font-style: normal;
  font-weight: 600;
  line-height: 22px;
  margin-bottom: 16px;
}

.tag-container .ant-tag {
  margin-right: 8px;
}

.flex {
  display: flex;
}

.item-title-box {
  display: flex;
  align-items: center;
  gap: 2px;
  color: #262626;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
  margin-bottom: 4px;
}

.btn-container {
  position: fixed;
  bottom: 0;
  right: 16px;
  display: flex;
  width: calc(100% - 270px);
  padding: 16px 1055px 16px 32px;
  align-items: center;
  border-radius: 0 0 2px 2px;
  background: #FFF;
  box-shadow: 0 -8px 4px 0 #0000000a;
}
</style>
