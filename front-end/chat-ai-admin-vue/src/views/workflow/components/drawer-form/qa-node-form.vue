<template>
  <div class="form-box">
    <a-form ref="formRef" :model="formState" :labelCol="labelCol" :wrapperCol="wrapperCol">
      <div class="title-block">基础信息</div>
      <a-form-item
        label="节点名称"
        name="node_name"
        :rules="[{ required: true, message: '请输入节点名称' }]"
      >
        <a-input
          v-model:value="formState.node_name"
          :maxLength="15"
          placeholder="请输入节点名称，最多15个字"
        ></a-input>
      </a-form-item>

      <a-form-item label="兜底回复" class="mb0" :name="['ears_reply', 'content']" required>
        <a-textarea
          placeholder="请输入兜底回复，最多一千字。在知识库匹配不到答案时，将会回复该文案。"
          v-model:value="formState.ears_reply.content"
          :maxLength="1000"
          :auto-size="{ minRows: 3, maxRows: 5 }"
        />
      </a-form-item>
      <a-form-item
        :name="['ears_reply', 'question']"
        :colon="false"
        :rules="[{
          trigger: 'change',
          validator: (rule, value) => checkedQuestion(rule, value),
        }]"
      >
        <template #label></template>
        <div class="question-guide-box">
          <draggable handle=".drag-btn" v-model="formState.ears_reply.question" item-key="key">
            <template #item="{ element, index }">
              <div class="question-guide-item" :key="element.key">
                <HolderOutlined class="drag-btn" />
                <div class="input-box">
                  <a-form-item-rest>
                    <a-input
                      :maxLength="14"
                      placeholder="请输入问题引导"
                      v-model:value="element.content"
                    ></a-input>
                  </a-form-item-rest>
                </div>
                <div class="btn-hover-wrap">
                  <CloseCircleOutlined class="del-btn" @click="onDelGuideItem(index)" />
                </div>
              </div>
            </template>
          </draggable>

          <div class="btn-box">
            <a-button
              :disabled="formState.ears_reply.question.length >= 6"
              @click="addGuideItem()"
              type="dashed"
              block
            >
              <template #icon>
                <PlusOutlined />
              </template>
              添加问题引导
            </a-button>
          </div>
        </div>
      </a-form-item>

      <div class="title-block">跳出逻辑</div>
      <a-form-item
        label="访客超过"
        name="no_resp_jump_time"
        :rules="[{
          trigger: 'change',
          validator: (rule, value) => checkedNoRespJumpTime(rule, value),
        }]"
      >
        <a-input
          v-model:value="formState.no_resp_jump_time"
          :maxLength="15"
          placeholder="请输入"
          type="number"
          style="width: 120px"
        ></a-input>
        <span> 分钟无响应，自动跳出当前节点进入下一节点</span>
        <div class="form-tip">最多30分钟，最少5分钟</div>
      </a-form-item>

      <a-form-item label=" 关键词转人工" class="mb0" name="keyword_switch_open" required>
        <a-switch
          v-model:checked="formState.keyword_switch_open"
          :checkedValue="1"
          :unCheckedValue="0"
          checked-children="开"
          un-checked-children="关"
        />
        <div class="form-tip">开启后，在问答节点内，访客触发指定关键词将自动转人工</div>
      </a-form-item>

      <a-form-item
        name="keywords"
        :colon="false"
        :rules="{
          trigger: 'change',
          validator: (rule, value) => checkedKeywords(rule, value),
        }"
      >
        <template #label></template>
        <a-form-item-rest>
          <div class="keyword-items">
            <div class="keyword-item" v-for="(item, index) in formState.keywords" :key="item.key">
              <div class="item-left">
                <a-select v-model:value="item.type" class="keyword-type-select" style="width: 86px">
                  <a-select-option value="1">全匹配</a-select-option>
                  <a-select-option value="2">半匹配</a-select-option>
                </a-select>
              </div>
              <div class="item-body">
                <a-input
                  class="keyword-text"
                  v-model:value="item.content"
                  style="width: 260px"
                  placeholder="请输入关键词"
                />
              </div>
              <div class="item-right">
                <div class="btn-hover-wrap">
                  <CloseCircleOutlined class="del-btn" @click="handleDelKeyword(index)" />
                </div>
              </div>
            </div>

            <div class="btn-box">
              <a-button
                :disabled="formState.keywords.length >= 100"
                @click="handleAddKeyword()"
                type="dashed"
                block
              >
                <template #icon>
                  <PlusOutlined />
                </template>
                添加转人工关键词
              </a-button>
            </div>
          </div>
        </a-form-item-rest>
      </a-form-item>
      <a-form-item label="转人工提示语" name="switch_content" required>
        <a-textarea
          placeholder="请输入转人工提示语"
          v-model:value="formState.switch_content"
          :maxLength="1000"
          :auto-size="{ minRows: 3, maxRows: 5 }"
        />
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup>
import { minutesToSeconds, secondsToMinutes} from '@/utils/index'
import { ref, reactive, inject, watch } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  CloseCircleOutlined,
  PlusOutlined,
  HolderOutlined,
} from '@ant-design/icons-vue'
import draggable from 'vuedraggable'
import { useRoute } from 'vue-router'
const query = useRoute().query

const { updateNodeItem, updateModifyNum } = inject('nodeInfo')

const node_sub_type = ref(61)

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
  ears_reply: {
    msg_type: 'text',
    content: '',
    question: [],
  },
  keyword_switch_open: 0, // 关键词转人工开关:0关,1开
  no_resp_jump_time: 5, // 访客无响应就自动跳出到下一节点的时间(秒)
  keywords: [
    // {
    //   type: '1', // 1全匹配,2半匹配
    //   content: '', // 关键词内容
    // },
  ],
  switch_content: '',
})

const setForm = (val) => {
  let data = JSON.parse(JSON.stringify(val))

  let ears_reply = data.ears_reply || {
    msg_type: 'text',
    content: '',
    question: [],
  }

  let form = {
    node_name: data.node_name,
    ears_reply: {
      msg_type: ears_reply.msg_type,
      content: ears_reply.content,
      question: ears_reply.question,
    },
    keyword_switch_open: data.keyword_switch_open || 0,
    no_resp_jump_time: secondsToMinutes(data.no_resp_jump_time) || 5,
    keywords: [],
    switch_content: data.switch_content,
  }

  if (data.ears_reply.question) {
    form.ears_reply.question.forEach((item) => {
      item.key = Math.random() * 10000000
    })
  }


  if (data.keyword_switch_full) {
    let keyword_switch_full = data.keyword_switch_full.split(',')

    keyword_switch_full.forEach((item) => {
      form.keywords.push({
        type: '1',
        key: Math.random() * 10000000,
        content: item,
      })
    })
  }

  if (data.keyword_switch_half) {
    let keyword_switch_half = data.keyword_switch_half.split(',')

    keyword_switch_half.forEach((item) => {
      form.keywords.push({
        type: '2',
        key: Math.random() * 10000000,
        content: item,
      })
    })
  }

  Object.assign(formState, form)
}

const checkedQuestion = (rule, value) => {
  if (value.length > 0) {
    let isEmpty =
      value
        .map((item) => item.content)
        .join('')
        .trim() == ''

    if (isEmpty) {
      return Promise.reject('请输入至少一个问题引导选项')
    }
  }
  return Promise.resolve()
}

const checkedNoRespJumpTime = (rule, value) => {
  if(!value){
    return Promise.reject('请输入访客无响应时长')
  }

  if (value < 5) {
    return Promise.reject('访客无响应时长最少5分钟')
  }
  if (value > 30) {
    return Promise.reject('访客无响应时长最多30分钟')
  }
  return Promise.resolve()
}

const checkedKeywords = (rule, value) => {
  if (formState.keyword_switch_open == 1) {
    if (formState.keywords.length == 0) {
      return Promise.reject('请添加至少一个转人工关键词')
    }

    let isEmpty = false
    for (let i = 0; i < formState.keywords.length; i++) {
      if (!formState.keywords[i].content) {
        isEmpty = true
        break
      }
    }

    if (isEmpty) {
      return Promise.reject('请输入转人工关键词')
    }
  }
  return Promise.resolve()
}

const addGuideItem = (index) => {
  formState.ears_reply.question.push({
    content: '',
    key: new Date().getTime(),
    node_key: '',
  })
}

const onDelGuideItem = (index, i) => {
  formState.ears_reply.question.splice(i, 1)
}

const handleAddKeyword = () => {
  formState.keywords.push({
    key: new Date().getTime(),
    type: '1', // 1全匹配,2半匹配
    content: '', // 关键词内容
  })
}
const handleDelKeyword = (index) => {
  formState.keywords.splice(index, 1)
}

const saveForm = () => {
  // 保存更新表单数据
  let updateInfo = {
    node_name: formState.node_name,
    switch_content: formState.switch_content,
    ears_reply: {
      msg_type: formState.ears_reply.question.length > 0 ? 'menu' : 'text',
      content: formState.ears_reply.content,
      question: [],
    },
    keyword_switch_open: formState.keyword_switch_open,
    no_resp_jump_time: minutesToSeconds(formState.no_resp_jump_time),
    keyword_switch_full: [],
    keyword_switch_half: [],
  }

  updateInfo.ears_reply.question = formState.ears_reply.question.map((item) => {
    return {
      content: item.content,
      node_key: item.node_key || '',
    }
  })

  formState.keywords.forEach((item) => {
    if (item.type == 1) {
      updateInfo.keyword_switch_full.push(item.content)
    }
    if (item.type == 2) {
      updateInfo.keyword_switch_half.push(item.content)
    }
  })

  updateInfo.keyword_switch_full = updateInfo.keyword_switch_full.join(',')
  updateInfo.keyword_switch_half = updateInfo.keyword_switch_half.join(',')

  updateNodeItem({ ...updateInfo })
}

const onSave = () => {
  formRef.value
    .validate()
    .then(() => {
      saveForm()
    })
    .catch((err) => {
      console.log('error', err)
    })
}

watch(
  () => props.properties,
  (val) => {
    setForm(val)
  },
  { immediate: true, deep: true },
)

defineExpose({
  onSave,
})
</script>

<style lang="less" scoped>
@import './common.less';

.mb0 {
  margin-bottom: 0;
}
.question-guide-box {
  display: flex;
  flex-direction: column;
  margin-top: 8px;

  .question-guide-item {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
    .input-box {
      flex: 1;
      margin-left: 4px;
    }
    .btn-hover-wrap {
      margin-left: 8px;
    }
    .del-btn {
      color: #8c8c8c;
      cursor: pointer;
    }
  }
}

.keyword-items {
  margin-top: 8px;
  .keyword-item {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
    .item-right {
      padding-left: 8px;
    }
    .del-btn {
      color: #8c8c8c;
      cursor: pointer;
    }
    .keyword-type-select {
      ::v-deep(.ant-select-selector) {
        border-right: 0;
        border-top-right-radius: 0;
        border-bottom-right-radius: 0;
      }
    }
    .keyword-text {
      border-top-left-radius: 0;
      border-bottom-left-radius: 0;
    }
  }
  .btn-box {
    margin-top: 8px;
  }
}
</style>
