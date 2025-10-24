<template>
  <div>
    <a-modal
      v-model:open="open"
      :title="modalTitle"
      @ok="handleOk"
      :width="472"
      ok-text="确定"
      cancel-text="取消"
    >
      <a-form ref="formRef" layout="vertical" :model="formState" style="margin: 24px 0">
        <a-form-item
          label="参数key"
          name="key"
          :rules="[
            {
              required: true,
              validator: (rule, value) => checkKey(rule, value)
            }
          ]"
        >
          <a-input
            v-model:value="formState.key"
            :maxlength="20"
            placeholder="英文字母和下划线“_”组成，最多20字符"
          />
        </a-form-item>
        <a-form-item
          label="参数类型"
          name="typ"
          :rules="[{ required: true, message: '请选择参数类型' }]"
        >
          <a-select v-model:value="formState.typ" placeholder="请选择" style="width: 100%">
            <a-select-option v-for="op in typOptions" :value="op.value">{{
              op.value
            }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item
          label="是否必填"
          name="required"
          :rules="[{ required: true, message: '请选择是否必填' }]"
        >
          <a-switch v-model:checked="formState.required" />
        </a-form-item>
        <a-form-item
          label="默认值"
          name="default"
          :rules="[{ required: formState.required, message: '请输入默认值' }]"
        >
          <a-input v-model:value="formState.default" placeholder="必填时大模型未提取到会返回默认值" />
        </a-form-item>
        <a-form-item label="枚举值" name="enum">
          <a-textarea
            style="height: 80px"
            v-model:value="formState.enum"
            placeholder="列举该字段的值，一行一个。系统会要求大模型只能返回列举的值"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { message } from 'ant-design-vue'
import { ref, reactive } from 'vue'

const emit = defineEmits(['add', 'edit', 'addSub'])

const open = ref(false)
const modalTitle = ref('新增参数')

let oprateType = 'add'

const checkKey = (rule, value) => {
  if (!value) {
    return Promise.reject('请输入参数key')
  }
  // 校验是否只包含英文字母和下划线
  const regex = /^[a-zA-Z_]+$/
  if (!regex.test(value)) {
    return Promise.reject('只能包含英文字母和下划线')
  }
  return Promise.resolve()
}
let INDEX = 0
const formState = reactive({
  key: '',
  typ: 'string',
  required: false,
  default: '',
  enum: '',
  cu_key: '',
  subs: [],
})
const show = () => {
  open.value = true
}

const add = () => {
  Object.assign(formState, {
    key: '',
    typ: 'string',
    required: false,
    default: '',
    enum: '',
    cu_key: Math.random() * 10000,
    subs: []
  })
  oprateType = 'add'
  modalTitle.value = '新增参数'
  open.value = true
}

const addSub = (index) => {
  INDEX = index
  Object.assign(formState, {
    key: '',
    typ: 'string',
    required: false,
    default: '',
    enum: '',
    cu_key: Math.random() * 10000,
    subs: []
  })
  oprateType = 'addSub'
  modalTitle.value = '新增参数'
  open.value = true
}

const edit = (data, index) => {
  INDEX = index
  Object.assign(formState, {
    key: '',
    typ: 'string',
    required: false,
    default: '',
    enum: '',
    cu_key: Math.random() * 10000,
    subs: []
  })
  Object.assign(formState, {
    ...data
  })
  oprateType = 'edit'
  modalTitle.value = '编辑参数'
  open.value = true
}

const formRef = ref(null)
const handleOk = () => {
  formRef.value.validate().then(() => {
    if (oprateType == 'add') {
      emit('add', {
        ...formState
      })
    }
    if (oprateType == 'addSub') {
      emit(
        'addSub',
        {
          ...formState
        },
        INDEX
      )
    }
    if (oprateType == 'edit') {
      emit(
        'edit',
        {
          ...formState
        },
        INDEX
      )
    }
    open.value = false
  })
}

const typOptions = [
  {
    lable: 'string',
    value: 'string'
  },
  {
    lable: 'number',
    value: 'number'
  },
  {
    lable: 'boole',
    value: 'boole'
  },
  {
    lable: 'float',
    value: 'float'
  },
  {
    lable: 'object',
    value: 'object'
  },
  {
    lable: 'array\<string>',
    value: 'array\<string>'
  },
  {
    lable: 'array\<number>',
    value: 'array\<number>'
  },
  {
    lable: 'array\<boole>',
    value: 'array\<boole>'
  },
  {
    lable: 'array\<float>',
    value: 'array\<float>'
  },
  {
    lable: 'array\<object>',
    value: 'array\<object>'
  }
]

defineExpose({
  show,
  add,
  edit,
  addSub
})
</script>

<style lang="less" scoped></style>
