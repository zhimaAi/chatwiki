<template>
  <div class="form-box">
    <a-form ref="formRef" :model="formState" :labelCol="labelCol" :wrapperCol="wrapperCol">
      <div class="title-block">基础信息</div>
      <a-form-item
        label="节点名称"
        name="node_name"
        :rules="[{ required: true, message: '请输入节点名称' }]"
        :labelCol="{ span: 4 }"
      >
        <a-input
          v-model:value="formState.node_name"
          :maxLength="15"
          placeholder="请输入节点名称，最多15个字"
        ></a-input>
      </a-form-item>

      <a-form-item label="通知内容" name="notice_content" :labelCol="{ span: 4 }">
        <a-textarea
          ref="textareaRef"
          @click="getCursorPosition"
          @input="handleInput"
          v-model:value="formState.notice_content"
          placeholder="请输入通知内容"
          :maxLength="1000"
          :auto-size="{ minRows: 3, maxRows: 5 }"
        />
        <div class="insert-block">
          <span>插入 :</span>
          <div
            @click="handleInsertKey(item.field_key)"
            class="tag-item"
            v-for="item in defaultFiledLists"
            :key="item.field_key"
          >
            {{ item.field_name }}
          </div>
          <a-dropdown>
            <template #overlay>
              <a-menu>
                <a-menu-item v-for="item in customFiledLists" :key="item.field_key"
                  ><div @click="handleInsertKey(item.field_key)">
                    {{ item.field_name }}
                  </div></a-menu-item
                >
              </a-menu>
            </template>
            <div class="tag-item">其他字段<DownOutlined /></div>
          </a-dropdown>
        </div>
      </a-form-item>

      <div class="title-block mt24">通知方式</div>
      <div class="message-card-box">
        <div class="message-title">
          <a-checkbox v-model:checked="formState.dingtalk.enable">钉钉群通知</a-checkbox>
          <a href="https://www.kancloud.cn/wikizhima/shipinhaoxiaodian/3162861" target="_blank"
            >如何配置钉钉群机器人？</a
          >
        </div>
        <div class="content-body" v-if="formState.dingtalk.enable">
          <a-form-item
            label="Webhook地址"
            :name="['dingtalk', 'url']"
            :rules="{
              required: true,
              validator: (rule, value) => checkedFields(rule, value, 'dingtalk'),
            }"
          >
            <a-input v-model:value="formState.dingtalk.url" placeholder="请输入"></a-input>
          </a-form-item>
          <a-form-item
            label="加签Secret"
            :name="['dingtalk', 'secret']"
            :rules="{
              required: true,
              validator: (rule, value) => checkedFields(rule, value, 'dingtalk'),
            }"
          >
            <a-input v-model:value="formState.dingtalk.secret" placeholder="请输入"></a-input>
          </a-form-item>
          <a-form-item class="mb8" label="@人：" :name="['dingtalk', 'at_type']">
            <a-radio-group v-model:value="formState.dingtalk.at_type">
              <a-radio :value="1">不@人</a-radio>
              <a-radio :value="2">@所有人</a-radio>
              <a-radio :value="3">指定群成员</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item
            class="hide-label mb2"
            :colon="false"
            :wrapperCol="{ span: 24 }"
            :name="['dingtalk', 'at']"
            :rules="[{ required: true, message: '请设置需要@的人' }]"
            v-if="formState.dingtalk.at_type == 3"
          >
            <template #label></template>
            <div class="tag-box">
              <template v-for="(tag, index) in formState.dingtalk.at" :key="tag">
                <a-tag :closable="true" @close="handleClose('dingtalk', tag)">
                  {{ tag }}
                </a-tag>
              </template>
              <a-input
                v-if="formState.dingtalk.inputVisible"
                v-model:value="formState.dingtalk.inputValue"
                ref="dingtalkInputRef"
                placeholder="请输入"
                type="text"
                size="small"
                :style="{ width: '78px' }"
                @blur="handleInputConfirm('dingtalk')"
                @keyup.enter="handleInputConfirm('dingtalk')"
              />
              <a-tag
                v-else
                style="background: #fff; border-style: dashed; cursor: pointer"
                @click="showInput('dingtalk')"
              >
                <plus-outlined />
                添加群成员
              </a-tag>
            </div>
          </a-form-item>
        </div>
      </div>
      <div class="message-card-box mt16">
        <div class="message-title">
          <a-checkbox v-model:checked="formState.work_wx.enable">企微群通知</a-checkbox>
          <a href="https://www.kancloud.cn/wikizhima/zmwk/3071448" target="_blank"
            >如何配置企微群机器人？</a
          >
        </div>
        <div class="content-body" v-if="formState.work_wx.enable">
          <a-form-item
            label="Webhook地址"
            :name="['work_wx', 'url']"
            :rules="{
              required: true,
              validator: (rule, value) => checkedFields(rule, value, 'work_wx'),
            }"
          >
            <a-input v-model:value="formState.work_wx.url" placeholder="请输入"></a-input>
          </a-form-item>
          <a-form-item class="mb8" label="@人：" :name="['work_wx', 'at_type']">
            <a-radio-group v-model:value="formState.work_wx.at_type">
              <a-radio :value="1">不@人</a-radio>
              <a-radio :value="2">@所有人</a-radio>
              <a-radio :value="3">指定群成员</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item
            class="hide-label mb2"
            :colon="false"
            :wrapperCol="{ span: 24 }"
            :name="['work_wx', 'at']"
            :rules="[{ required: true, message: '请设置需要@的人' }]"
            v-if="formState.work_wx.at_type == 3"
          >
            <template #label></template>
            <div class="tag-box">
              <template v-for="(tag, index) in formState.work_wx.at" :key="tag">
                <a-tag :closable="true" @close="handleClose('work_wx', tag)">
                  {{ tag }}
                </a-tag>
              </template>
              <a-input
                v-if="formState.work_wx.inputVisible"
                v-model:value="formState.work_wx.inputValue"
                ref="work_wxInputRef"
                placeholder="请输入"
                type="text"
                size="small"
                :style="{ width: '78px' }"
                @blur="handleInputConfirm('work_wx')"
                @keyup.enter="handleInputConfirm('work_wx')"
              />
              <a-tag
                v-else
                style="background: #fff; border-style: dashed; cursor: pointer"
                @click="showInput('work_wx')"
              >
                <plus-outlined />
                添加群成员
              </a-tag>
            </div>
          </a-form-item>
        </div>
      </div>
      <div class="message-card-box mt16">
        <div class="message-title">
          <a-checkbox v-model:checked="formState.feishu.enable">飞书通知</a-checkbox>
          <a href="https://www.kancloud.cn/wikizhima/zmwk/3071448" target="_blank"
            >如何配置飞书机器人？</a
          >
        </div>
        <div class="content-body" v-if="formState.feishu.enable">
          <a-form-item
            label="Webhook地址"
            :name="['feishu', 'url']"
            :rules="{
              required: true,
              validator: (rule, value) => checkedFields(rule, value, 'feishu'),
            }"
          >
            <a-input v-model:value="formState.feishu.url" placeholder="请输入"></a-input>
          </a-form-item>
          <a-form-item
            class="mb8"
            label="加签Secret"
            :name="['feishu', 'secret']"
            :rules="{
              required: true,
              validator: (rule, value) => checkedFields(rule, value, 'feishu'),
            }"
          >
            <a-input v-model:value="formState.feishu.secret" placeholder="请输入"></a-input>
          </a-form-item>
        </div>
      </div>
    </a-form>
  </div>
</template>

<script setup>
import { ref, reactive, inject, watch, nextTick } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  CloseCircleFilled,
  CloseCircleOutlined,
  LoadingOutlined,
  DownOutlined,
  PlusOutlined,
  EditOutlined,
} from '@ant-design/icons-vue'
import { getFieldsList } from '@/api/retention/index.js'
import { useRoute } from 'vue-router'
const query = useRoute().query

const { updateNodeItem, updateModifyNum } = inject('nodeInfo')

const props = defineProps({
  properties: {
    type: Object,
    default: () => ({}),
  },
})

const emit = defineEmits(['ok'])
const formRef = ref()

const labelCol = {
  span: 5,
}
const wrapperCol = {
  span: 19,
}

const formState = reactive({
  node_name: '',
  notice_content: '',
  dingtalk: {
    enable: true,
    url: '',
    secret: '',
    at_type: 1,
    at: [],
    inputVisible: false,
    inputValue: '',
  },
  work_wx: {
    enable: true,
    url: '',
    at_type: 1,
    at: [],
    inputVisible: false,
    inputValue: '',
  },
  feishu: {
    enable: true,
    url: '',
    secret: '',
  },
})
let updateNum = 0
watch(
  () => props.properties,
  (val) => {
    try {
      formState.node_name = val.node_name
      formState.notice_content = val.notice_content
      formState.dingtalk = {
        enable: true,
        url: '',
        secret: '',
        at_type: 1,
        at: [],
        inputVisible: false,
        inputValue: '',
      }

      formState.work_wx = {
        enable: true,
        url: '',
        at_type: 1,
        at: [],
        inputVisible: false,
        inputValue: '',
      }
      formState.feishu = {
        enable: true,
        url: '',
        secret: '',
      }
      let notice_modes = val.notice_modes

      if (notice_modes && notice_modes.length) {
        notice_modes.forEach((item) => {
          formState[item.type].enable = item.enable
          formState[item.type].url = item.url
          formState[item.type].secret = item.secret
          if (item.at == '') {
            formState[item.type].at_type = 1
          }
          if (item.at == '*') {
            formState[item.type].at_type = 2
          }
          if (item.at != '*' && item.at != '') {
            formState[item.type].at_type = 3
            formState[item.type].at = item.at.split(',')
          }
        })
      }
      getFields()
      updateNum = 0
    } catch (error) {
      console.log(error)
    }
  },
  { immediate: true, deep: true },
)

watch(
  () => formState,
  () => {
    updateNum++
    updateModifyNum(updateNum)
  },
  {
    deep: true,
  },
)

const customFiledLists = ref([])
const defaultFiledLists = ref([])
function getFields() {
  getFieldsList({
    page: 1,
    limit: 9999,
    robot_id: query.robotId,
  }).then((res) => {
    if (!res.data) {
      return
    }
    customFiledLists.value = res.data.data.filter((item) => item.is_default == 0)
    defaultFiledLists.value = res.data.data.filter((item) => item.is_default == 1)
  })
}

const checkedFields = (rule, value, key) => {
  if (formState[key].enable) {
    if (value == '') {
      return Promise.reject('请输入')
    }
  }
  return Promise.resolve()
}

const textareaRef = ref(null)
const cursorPosition = ref(0)

// 获取光标位置
const getCursorPosition = () => {
  const textarea = textareaRef.value?.resizableTextArea?.textArea // 获取实际的 DOM 元素
  if (textarea) {
    cursorPosition.value = textarea.selectionStart
  }
}

// 输入事件（可选，处理输入内容变化）
const handleInput = () => {
  getCursorPosition()
}

const handleInsertKey = (key) => {
  // 插入内容
  insertText(`【${key}】`)
}

const insertText = (insertValue) => {
  const textarea = textareaRef.value?.resizableTextArea?.textArea
  if (textarea) {
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    // 重新组装文本
    formState.notice_content =
      formState.notice_content.slice(0, start) + insertValue + formState.notice_content.slice(end)

    // 更新光标位置到插入文字之后
    const newCursorPosition = start + insertValue.length
    nextTick(() => {
      textarea.setSelectionRange(newCursorPosition, newCursorPosition)
    })
  }
}

const handleClose = (key, removedTag) => {
  const tags = formState[key].at.filter((tag) => tag !== removedTag)
  formState[key].at = tags
}
const dingtalkInputRef = ref(null)
const work_wxInputRef = ref(null)
const showInput = (key) => {
  formState[key].inputVisible = true
  nextTick(() => {
    if (key == 'dingtalk') {
      dingtalkInputRef.value.focus()
    }
    if (key == 'work_wx') {
      work_wxInputRef.value.focus()
    }
  })
}
const handleInputConfirm = (key) => {
  const inputValue = formState[key].inputValue
  let tags = formState[key].at
  if (inputValue && tags.indexOf(inputValue) === -1) {
    tags = [...tags, inputValue]
  }
  formState[key].at = tags
  formState[key].inputVisible = false
  formState[key].inputValue = ''
}

const saveForm = () => {
  // 保存更新表单数据
  let updateInfo = {
    node_name: formState.node_name,
    notice_content: formState.notice_content,
  }
  let notice_modes = []
  let keyMap = ['dingtalk', 'work_wx', 'feishu']
  keyMap.forEach((item) => {
    let dataItem = {
      type: item,
      enable: formState[item].enable,
      url: formState[item].url,
      secret: formState[item].secret,
      at: '',
    }
    if (formState[item].at_type == 2) {
      dataItem.at = '*'
    }
    if (formState[item].at_type == 3) {
      dataItem.at = formState[item].at.join(',')
    }
    notice_modes.push(dataItem)
  })
  updateInfo.notice_modes = notice_modes
  updateNodeItem({ ...updateInfo })
}

const onSave = () => {
  formRef.value
    .validate()
    .then(() => {
      if (!formState.dingtalk.enable && !formState.work_wx.enable && !formState.feishu.enable) {
        return message.error('至少设置一种通知方式')
      }
      saveForm()
    })
    .catch((err) => {
      console.log('error', err)
    })
}

defineExpose({
  onSave,
})
</script>

<style lang="less" scoped>
@import './common.less';

.insert-block {
  margin-top: 4px;
  display: flex;
  align-items: center;
  color: #262626;
  flex-wrap: wrap;
  gap: 8px;
  .tag-item {
    height: 24px;
    width: fit-content;
    display: flex;
    align-items: center;
    padding: 0 8px;
    border: 1px solid #d9d9d9;
    border-radius: 2px;
    color: #595959;
    cursor: pointer;
  }
}

.message-card-box {
  background: #f2f4f7;
  border-radius: 6px;
  .message-title {
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-bottom: 1px solid #e4e6eb;
    height: 32px;
    padding: 0 16px;
  }
  .content-body {
    padding: 16px;
  }
}

.tag-box {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  .ant-tag {
    margin-inline-end: 0;
  }
}
</style>
