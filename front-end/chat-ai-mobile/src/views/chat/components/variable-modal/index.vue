<template>
  <van-popup
    v-model:show="open"
    round
    :closeable="false"
    :close-on-click-overlay="false"
    position="bottom"
    class="variable-modal-class"
  >
    <div style="padding-top: 20px">
      <div style="background: #fff; border-radius: 24px 24px 0 0">
        <div class="modal-header">
          <div class="banner-bg">
            <img src="@/assets/img/form-modal-header.png" alt="" />
          </div>
          <div class="title">请先填写信息</div>
          <div class="close-btn" @click="handleClose" v-if="isEdit">
            <van-icon name="clear" />
          </div>
        </div>
        <div class="form-body">
          <van-form ref="formRef" @submit="handleSave">
            <!-- 字段名 -->
            <template v-for="item in wait_variables" :key="item.variable_key">
              <div
                class="form-label"
                :class="{ required: item.must_input == 1 }"
                v-if="item.variable_type != 'checkbox_switch'"
              >
                {{ item.variable_name }}
              </div>
              <div class="field-item-blcok" :class="[`field-item-${item.variable_type}`]">
                <van-field
                  v-if="item.variable_type === 'input_string'"
                  :name="item.variable_key"
                  :label="null"
                  :required="item.must_input == 1"
                  :placeholder="`请输入${item.variable_name}`"
                  :rules="getRules(item)"
                  v-model.trim="formState[item.variable_key]"
                  :maxlength="getLength(item)"
                />

                <van-field
                  v-else-if="item.variable_type === 'input_number'"
                  :name="item.variable_key"
                  :label="null"
                  :required="item.must_input == 1"
                  :placeholder="`请输入${item.variable_name}`"
                  :rules="getRules(item)"
                  v-model="formState[item.variable_key]"
                  type="number"
                />

                <van-field
                  v-else-if="item.variable_type === 'select_one'"
                  :name="item.variable_key"
                  :label="null"
                  :required="item.must_input == 1"
                  :placeholder="`请选择${item.variable_name}`"
                  :rules="getRules(item)"
                  v-model="formState[item.variable_key]"
                  is-link
                  readonly
                  @click="handleShowPick(item)"
                />

                <van-field
                  v-else-if="item.variable_type === 'checkbox_switch'"
                  :name="item.variable_key"
                  :label="null"
                  :required="item.must_input == 1"
                >
                  <template #input>
                    <van-checkbox v-model="formState[item.variable_key]" shape="square">{{
                      item.variable_name
                    }}</van-checkbox>
                  </template>
                </van-field>
              </div>
            </template>

            <div style="margin: 32px 16px">
              <van-button round block type="primary" native-type="submit"> 提交 </van-button>
            </div>
          </van-form>
        </div>
      </div>
    </div>

    <!-- 选择器弹窗 -->
    <van-popup v-model:show="showPicker" round position="bottom">
      <van-picker :columns="currentOptions" @confirm="onConfirm" @cancel="showPicker = false" />
    </van-popup>
  </van-popup>
</template>

<script setup>
import { ref, computed, watch, reactive } from 'vue'
import { useRoute } from 'vue-router'
import { useChatStore } from '@/stores/modules/chat'
import { showFailToast } from 'vant'

const query = useRoute().query
const chatStore = useChatStore()

const isEdit = ref(false)

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
  return chat_variables.value.need_fill_variable
})

// 弹窗控制
const open = ref(false)
const showPicker = ref(false)
const currentField = ref('')
const currentOptions = computed(() => {
  const field = wait_variables.value.find((item) => item.variable_key === currentField.value)
  return (
    field?.options?.map((option) => {
      return {
        text: option.label,
        value: option.label
      }
    }) || []
  )
})

// 获取当前选项的索引
const getCurrentOptionIndex = computed(() => {
  const currentValue = formState[currentField.value]
  const options = currentOptions.value
  const index = options.indexOf(currentValue)
  return index >= 0 ? index : 0
})

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
      formState[item.variable_key] = item.value || ''
    }

    if (item.variable_type == 'checkbox_switch') {
      // 开关
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
      formState[item.variable_key] = item.default_value || ''
    }

    if (item.variable_type == 'checkbox_switch') {
      // 开关
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

// 获取验证规则
const getRules = (item) => {
  const rules = []
  if (item.must_input == 1) {
    rules.push({
      required: true,
      message: `请输入${item.variable_name}`,
      validator: (val) => {
        if (item.must_input == 1) {
          return val !== undefined && val !== null && val !== ''
        }
        return true
      }
    })
  }
  return rules
}

watch(
  () => wait_variables.value,
  (val) => {
    val && initData(wait_variables.value)
  },
  {
    immediate: true,
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

const handleClose = () => {
  open.value = false
}

const handleShowPick = (item) => {
  showPicker.value = true
  currentField.value = item.variable_key
}

// 处理表单提交
const handleSave = async (values) => {
  try {
    // 手动触发表单验证
    await formRef.value?.validate()

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
  } catch (error) {
    console.error('验证失败', error)
    showFailToast(error.message || '提交失败，请检查输入内容')
  }
}

function getPostData() {
  let list = wait_variables.value

  return list.map((item) => {
    let value = formState[item.variable_key]
    if (item.variable_type == 'checkbox_switch') {
      value = value ? '1' : '0'
    }
    return {
      ...item,
      value
    }
  })
}

// 选择器确认事件
const onConfirm = ({ selectedValues }) => {
  formState[currentField.value] = selectedValues[0]
  showPicker.value = false
}

const formRef = ref(null)

defineExpose({
  handleEdit
})
</script>

<style>
.variable-modal-class.van-popup {
  background: none;
}
</style>
<style lang="less" scoped>
.modal-header {
  // padding: 20px 16px 10px 16px;
  position: relative;
  border-radius: 20px;
  height: 76px;
  // background-color: #f8f9fa;
  .banner-bg {
    position: absolute;
    top: -20px;
    left: 0;
    right: 0;
    height: 96px;
    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }
  .close-btn {
    font-size: 32px;
    position: absolute;
    right: 12px;
    top: -12px;
    cursor: pointer;
  }

  .title {
    color: #000000;
    font-weight: 600;
    position: relative;
    padding-left: 24px;
    display: flex;
    align-items: center;
    height: 100%;
    font-size: 24px;
  }
}

.form-body {
  padding-top: 16px;
  padding-left: 16px;
  padding-right: 12px;
  max-height: 70vh;
  overflow-y: auto;
  .form-label {
    color: #262626;
    font-size: 14px;
    font-weight: 600;
    line-height: 20px;
    margin-bottom: 4px;
    &.required::before {
      content: '*';
      color: red;
      margin-right: 4px;
      font-size: 12px;
    }
  }
  &::v-deep(.field-item-blcok) {
    margin-bottom: 16px;
  }
  .field-item-checkbox_switch {
    &::v-deep(.van-cell) {
      border: none;
      padding-left: 0;
    }
  }
  &::v-deep(.van-cell) {
    border: 1px solid var(--06, #d9d9d9);
    border-radius: 6px;
    padding: 7px 12px;
    &::after {
      border: none;
    }
  }
}

.from-footer {
  padding: 16px 0;
}
</style>
