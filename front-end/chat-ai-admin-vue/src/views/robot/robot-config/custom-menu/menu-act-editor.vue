<template>
  <div class="menu-act-editor">
    <div class="radio-row">
      <a-radio-group v-model:value="state.choose_act_item" @change="emitChange">
        <a-radio :value="1">发送消息</a-radio>
        <a-radio :value="2">跳转网页</a-radio>
        <a-radio :value="3">跳转小程序</a-radio>
        <a-radio :value="5">推送事件</a-radio>
      </a-radio-group>
    </div>

    <div v-if="state.choose_act_item === 1" class="block">
      <div class="label">回复内容</div>
      <div class="reply-list">
        <MultiReply v-for="(it, idx) in replyList" :key="idx" ref="replyRefs" v-model:value="replyList[idx]" :reply_index="idx"
          @change="onContentChange" @del="onDelItem" />
        <a-button type="dashed" style="width: 100%" :disabled="replyList.length >= 5" @click="addReply">
          <PlusOutlined /> 添加回复内容({{ replyList.length }}/5)
        </a-button>
      </div>
      <div class="label">回复方式</div>
      <a-radio-group v-model:value="reply_num" @change="emitChange"><a-radio :value="0">全部回复</a-radio><a-radio
          :value="1">随机回复一条</a-radio></a-radio-group>
    </div>

    <div v-else-if="state.choose_act_item === 2" class="block block-box">
      <div class="label"><span class="required">*</span>跳转链接</div>
      <a-form class="link-form" :model="linkModel" :rules="rules.link" ref="linkForm" layout="horizontal">
        <a-form-item name="linkUrl">
          <a-input v-model:value="linkModel.linkUrl" placeholder="请输入链接" @blur="emitChange" />
        </a-form-item>
      </a-form>
      <div class="tip">订阅者点击该菜单会跳到以下链接。必须以http://或https://开头</div>
    </div>

    <div v-else-if="state.choose_act_item === 3" class="block block-box">
      <div class="form-grid">
        <a-form :model="miniCard" :rules="rules.card" ref="cardForm" layout="horizontal" :labelCol="cardLabelCol"
          :wrapperCol="cardWrapperCol">
          <a-alert type="info" style="margin-bottom: 16px;"
            message="小程序卡片仅支持在公众号，微信客服和微信小程序中发送，其他渠道会发送失败；公众号内回复必须是关联小程序，微信小程序内回复必须是当前的小程序" />
          <a-form-item label="小程序appID" name="appid">
            <a-input v-model:value="miniCard.appid" placeholder="请输入小程序appID" />
            <div style="margin-top: 2px; color: #FB363F;">小程序右上角三个点>名字>更多资料>Appid;公众号内回复时必须跟小程序是关联关系</div>
          </a-form-item>
          <a-form-item label="小程序路径" name="page_path">
            <a-input v-model:value="miniCard.page_path" placeholder="请输入小程序路径" />
            <div style="margin-top: 2px; color: #8C8C8C;">请联系小程序开发者获取路径,比如/pages/index/index,注意，路径建议以/开头</div>
          </a-form-item>
          <a-form-item label="备用网页" name="standbyUrl">
            <a-input v-model:value="miniCard.standbyUrl" placeholder="请输入备用网页" />
            <div style="margin-top: 2px; color: #8C8C8C;">当不填时将使用默认的网页，旧版微信客户端无法支持小程序，用户点击菜单时将会打开备用网页</div>
          </a-form-item>
        </a-form>
      </div>
    </div>

    <div v-else-if="state.choose_act_item === 5" class="block block-box">
      <div class="label"><span class="required">*</span>设置key值</div>
      <a-input v-model:value="eventKey" placeholder="请输入key值" @blur="emitChange" />
      <div class="tip">设定key值，可以用于跟其他系统做交互。便于同时使用第三方系统开发跟小客服后台的菜单功能。如需使用小客服后台的自动打标签，请勿在多个菜单内设置相同的推送事件值。</div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { message } from 'ant-design-vue'
import MultiReply from '@/components/replay-card/multi-reply.vue'
import { PlusOutlined } from '@ant-design/icons-vue'

const replyRefs = ref([])
const props = defineProps({
  value: {
    type: Object,
    default: () => ({
      choose_act_item: 0,
      act_params: {}
    })
  }
})
const emit = defineEmits(['update:value', 'change'])

const state = reactive({
  choose_act_item: Number(props.value?.choose_act_item || 0)
})
const replyList = ref(((props.value?.act_params || {}).replyList && (props.value?.act_params || {}).replyList.length) ? (props.value?.act_params || {}).replyList : [{ type: 'text', description: '' }])
const reply_num = ref(Number((props.value?.act_params || {}).reply_num || 0))
const linkModel = reactive({
  linkUrl: String((props.value?.act_params || {}).linkUrl || '')
})
const eventKey = ref(String((props.value?.act_params || {}).key || ''))
const miniCard = reactive({
  appid: (props.value?.act_params || {}).appid || '',
  page_path: (props.value?.act_params || {}).page_path || '',
  standbyUrl: (props.value?.act_params || {}).standbyUrl || ''
})
const cardForm = ref(null)
const linkForm = ref(null)
const cardLabelCol = { span: 4 }
const cardWrapperCol = { span: 20 }

const rules = {
  card: {
    appid: [
      {
        required: true,
        message: '请输入小程序appID'
      }
    ],
    page_path: [
      {
        required: true,
        message: '请输入小程序路径'
      }
    ],
    standbyUrl: [
      {
        required: true,
        pattern: /^https?:\/\//,
        message: '备用网页需以http://或https://开头'
      }
    ]
  },
  link: {
    linkUrl: [
      {
        required: true,
        message: '请输入链接'
      },
      {
        pattern: /^https?:\/\//,
        message: '链接需以http://或https://开头'
      }
    ]
  }
}

watch(() => props.value, (val) => {
  state.choose_act_item = Number(val?.choose_act_item || 0)
  const ap = val?.act_params || {}
  replyList.value = (Array.isArray(ap.replyList) && ap.replyList.length) ? ap.replyList : [{ type: 'text', description: '' }]
  reply_num.value = Number(ap.reply_num || 0)
  linkModel.linkUrl = String(ap.linkUrl || '')
  eventKey.value = String(ap.key || '')
  miniCard.appid = ap.appid || ''
  miniCard.page_path = ap.page_path || ''
  miniCard.standbyUrl = ap.standbyUrl || ''
  replyRefs.value = []
}, { immediate: false })

watch(() => miniCard, () => {
  emitChange()
}, { deep: true })

function addReply () {
  if (replyList.value.length >= 5) return;
  replyList.value.push({ type: 'text', description: '' })
}
function onContentChange (payload) {
  const { reply_index, ...rest } = payload || {};
  if (reply_index >= 0 && reply_index < replyList.value.length) {
    replyList.value[reply_index] = rest;
    emitChange()
  }
}
function onDelItem (index) {
  replyList.value.splice(index, 1);
  emitChange()
}
async function emitChange () {
  emit('update:value', getValue());
  emit('change', getValue())
}
function getValue () {
  const act_params = {};
  if (state.choose_act_item === 1) {
    Object.assign(act_params, {
      replyList: replyList.value,
      reply_num: reply_num.value
    });
  }
  if (state.choose_act_item === 2) {
    Object.assign(act_params, {
      linkUrl: linkModel.linkUrl
    });
  }
  if (state.choose_act_item === 3) {
    Object.assign(act_params, {
      appid: miniCard.appid,
      page_path: miniCard.page_path,
      standbyUrl: miniCard.standbyUrl
    });
  }
  if (state.choose_act_item === 5) {
    Object.assign(act_params, {
      key: eventKey.value
    });
  }
  return {
    choose_act_item: state.choose_act_item,
    act_params
  }
}


async function validate () {
  const item = Number(state.choose_act_item || 0);
  if (item === 1) {
    if (!replyList.value.length) {
      message.error('请添加回复内容');
      return false
    }
    for (const comp of replyRefs.value) {
      if (comp && comp.validate) {
        const ok = await comp.validate();
        if (!ok) {
          return false
        }
      }
    }
    return true
  }
  if (item === 2) {
    try {
      await linkForm.value?.validate();
      return true
    } catch (_e) {
      const s = String(linkModel.linkUrl || '').trim();
      if (!s) {
        message.error('请输入链接')
      } else {
        message.error('链接需以http://或https://开头')
      }
      return false
    }
  }
  if (item === 3) {
    try {
      await cardForm.value?.validate();
      return true
    } catch (_e) {
      message.error('请完善小程序必填项');
      return false
    }
  } if (item === 5) {
    const s = String(eventKey.value || '').trim();
    if (!s) {
      message.error('请输入推送事件key');
      return false
    }
    return true
  }
  return true
}

defineExpose({
  getValue,
  validate
})
</script>

<style scoped lang="less">
.radio-row {
  margin-bottom: 24px;
}

.label {
  color: #262626;
  font-weight: 400;
  margin-bottom: 4px;
}

.required {
  color: #FB363F;
  margin-right: 4px;
}

.tip {
  color: #8c8c8c;
  font-size: 12px;
  margin-top: 4px;
}

.reply-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 24px;
}

.block-box {
  display: flex;
  flex-direction: column;
  align-items: self-start;
  width: 694px;
  padding: 16px 24px;
  border-radius: 6px;
  background: #F2F4F7;

  .link-form {
    width: 100%;

    .ant-form-item {
      margin-bottom: 0px;
    }
  }
}

.form-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
