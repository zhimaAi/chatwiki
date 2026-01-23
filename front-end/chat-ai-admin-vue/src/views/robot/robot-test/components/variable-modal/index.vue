<template>
  <a-modal
    wrapClassName="variable-modal-class"
    v-model:open="open"
    :title="null"
    :footer="null"
    :closable="isEdit"
    :maskClosable="false"
    :width="560"
  >
    <div class="modal-header">
      <div class="banner-bg">
        <img src="@/assets/img/form-modal-header.png" alt="" />
      </div>
      <div class="title">请先填写信息</div>
    </div>
    <div class="form-body">
      <a-form :model="formState" ref="formRef" layout="vertical">
        <!-- 字段名 -->
        <a-form-item
          v-for="item in wait_variables"
          :key="item.variable_key"
          :label="item.variable_type == 'checkbox_switch' ? null : item.variable_name"
          :name="item.variable_key"
          :rules="[{ required: item.must_input == 1, message: `请输入${item.variable_name}` }]"
        >
          <template v-if="item.variable_type == 'input_string'">
            <a-input
              v-model:value="formState[item.variable_key]"
              :maxLength="getLength(item)"
              placeholder="请输入"
            />
          </template>
          <template v-if="item.variable_type == 'input_number'">
            <a-input-number
              style="width: 100%"
              v-model:value="formState[item.variable_key]"
              stringMode
              :maxLength="50"
              placeholder="请输入"
            />
          </template>
          <template v-if="item.variable_type == 'select_one'">
            <a-select v-model:value="formState[item.variable_key]" style="width: 100%">
              <a-select-option
                v-for="option in item.options"
                :value="option.label"
                :key="option.label"
                >{{ option.label }}</a-select-option
              >
            </a-select>
          </template>
          <template v-if="item.variable_type == 'checkbox_switch'">
            <a-checkbox v-model:checked="formState[item.variable_key]">{{
              item.variable_name
            }}</a-checkbox>
          </template>
        </a-form-item>
      </a-form>
    </div>
    <div class="from-footer">
      <div class="btn-box">
        <a-button @click="handleSave" block type="primary">提 交</a-button>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, computed, watch, reactive } from 'vue'
import { useRoute } from 'vue-router'
import { useChatStore } from '@/stores/modules/chat'
const query = useRoute().query

const isEdit = ref(false)

const chatStore = useChatStore()
const chat_variables = computed(() => {
  return chatStore.chat_variables || {}
})

const wait_variables = computed(() => {
  let list = chat_variables.value.wait_variables || []
  return list.map((item) => {
    return {
      ...item,
      options: item.options ? JSON.parse(item.options) : []
    }
  })
})

const fill_variables = computed(() => {
  return chat_variables.value.fill_variables || []
})

const need_fill_variable = computed(() => {
  return chat_variables.value.need_fill_variable || false
})

const open = ref(false)

const formState = reactive({})

const handleEdit = () => {
  isEdit.value = true
  fill_variables.value.forEach((item) => {
    if (item.variable_type == 'input_string' || item.variable_type == 'input_number') {
      // 文本 数字
      formState[item.variable_key] = item.value || ''
    }

    if (item.variable_type == 'select_one') {
      // 下拉单选
      formState[item.variable_key] = item.value || void 0
    }

    if (item.variable_type == 'checkbox_switch') {
      // 文本
      formState[item.variable_key] = item.value == 1
    }
  })
  open.value = true
}

const initData = (list) => {
  list.forEach((item) => {
    if (item.variable_type == 'input_string' || item.variable_type == 'input_number') {
      // 文本 数字
      formState[item.variable_key] = item.default_value || ''
    }

    if (item.variable_type == 'select_one') {
      // 下拉单选
      formState[item.variable_key] = item.default_value || void 0
    }

    if (item.variable_type == 'checkbox_switch') {
      // 文本
      formState[item.variable_key] = item.default_value == 1
    }
  })
}

const getLength = (item) => {
  if (item.max_input_length > 0) {
    return +item.max_input_length
  }
  return 50
}

watch(
  () => wait_variables.value,
  (val) => {
    console.log(chat_variables.value, '===')
    val && initData(wait_variables.value)
  },
  {
    deep: true
  }
)

watch(
  () => need_fill_variable.value,
  (val) => {
    if (val) {
      isEdit.value = false
      open.value = true
    }
  },
  {
    deep: true,
    immediate: true
  }
)

const formRef = ref(null)

const handleSave = () => {
  formRef.value.validate().then(() => {
    let chat_prompt_variables = getPostData()
    let variables_key = `chat_prompt_variables_${query.robot_key}`

    if (isEdit.value) {
      chatStore.handleEditVariables({
        chat_prompt_variables
      })
    } else {
      localStorage.setItem(variables_key, JSON.stringify(chat_prompt_variables))
    }
    open.value = false
  })
}

function getPostData() {
  let list = wait_variables.value

  return list.map((item) => {
    let value = formState[item.variable_key]
    if (item.variable_type == 'checkbox_switch') {
      value = value ? 1 : 0
    }
    return {
      ...item,
      value
    }
  })
}

defineExpose({
  handleEdit
})
</script>
<style lang="less">
.variable-modal-class .ant-modal .ant-modal-content {
  padding: 0;
}
</style>

<style lang="less" scoped>
.modal-header {
  height: 76px;
  display: flex;
  align-items: center;
  position: relative;
  .title {
    position: relative;
    z-index: 9;
    padding-left: 24px;
    color: #000000;
    font-size: 24px;
    font-weight: 600;
  }
  .banner-bg {
    height: 96px;
    position: absolute;
    left: 0;
    right: 0;
    top: -20px;
    overflow: hidden;
    img {
      width: 100%;
      height: 100%;
    }
  }
}
.form-body {
  padding: 16px 24px 0 24px;
  max-height: 500px;
  overflow-y: auto;
  &::v-deep(.ant-form-item) {
    margin-bottom: 16px;
  }
}
.from-footer {
  padding: 8px 24px 24px 24px;
}
</style>
