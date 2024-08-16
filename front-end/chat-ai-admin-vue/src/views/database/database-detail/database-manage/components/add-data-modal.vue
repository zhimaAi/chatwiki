<template>
  <div class="add-data-sheet">
    <a-modal
      v-model:open="open"
      :title="modalTitle"
      @ok="handleOk"
      :width="580"
      :bodyStyle="{
        'max-height': '590px',
        'overflow-y': 'auto',
        width: 'calc(100% + 24px)',
        'padding-right': '32px'
      }"
      class="add-data-sheet-modal"
    >
      <a-form style="margin-top: 16px" :model="formState" ref="formRef" layout="vertical">
        <a-form-item
          v-for="item in fieldLists"
          :name="item.name"
          :key="item.id"
          :rules="item.rules"
        >
          <template #label>
            <a-flex :gap="4" style="width: 100%">
              <div class="field-name">{{ item.name }}</div>
              <a-tooltip>
                <template #title>{{ item.description }}</template>
                <QuestionCircleOutlined />
              </a-tooltip>
            </a-flex>
          </template>
          <template v-if="item.type == 'boolean'">
            <a-switch
              v-model:checked="formState[item.name]"
              checked-children="开"
              un-checked-children="关"
              checkedValue="true"
              unCheckedValue="false"
            />
          </template>
          <template v-else-if="item.type == 'number' || item.type == 'integer'">
            <a-input-number
              style="width: 100%"
              v-model:value.number="formState[item.name]"
              :placeholder="`请输入内容`"
            />
          </template>
          <template v-else>
            <a-input v-model:value="formState[item.name]" :placeholder="`请输入内容`" />
          </template>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { PlusOutlined, QuestionCircleOutlined } from '@ant-design/icons-vue'
import { reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { addFormEntry } from '@/api/database'
const modalTitle = ref('添加数据')
const rotue = useRoute()
const query = rotue.query
const emit = defineEmits('ok')
const props = defineProps({
  column: {
    type: Array,
    default: () => []
  }
})
watch(
  () => props.column,
  (val) => {
    formatField()
  }
)

const fieldLists = ref([])
const formatField = () => {
  let lists = []
  props.column.forEach((item) => {
    formState[item.name] = ''
    item.rules = []
    if (item.required == 'true') {
      item.rules.push({ required: true, message: `请输入${item.name}` })
    }
    if (item.type == 'number') {
      item.rules.push({ type: 'number', message: `${item.name}必须为数字` })
    }
    if (item.type == 'integer') {
      // item.rules.push({ type: 'integer', message: `${item.name}必须为整数` })
      item.rules.push({ asyncValidator: isValidatorInteger })
    }
    if (item.type == 'boolean') {
      formState[item.name] = 'false'
    }
    lists.push({
      ...item
    })
    fieldLists.value = lists
  })
}
const isValidatorInteger = (rule, value) => {
  return new Promise((resolve, reject) => {
    if(value == '' || Number.isInteger(+value)){
      return resolve()
    }
    return reject('请输入整数')
  })
}
const open = ref(false)
const formState = reactive({
  form_id: query.form_id
})
const show = (data) => {
  open.value = true
  formRef.value && formRef.value.resetFields()
  formatField()
  formState.id = ''
  Object.assign(formState, convertNumberStringsToObjectNumbers(data))
  if (formState.id) {
    modalTitle.value = '修改数据'
  } else {
    modalTitle.value = '添加数据'
  }
}

const formRef = ref(null)
const handleOk = () => {
  formRef.value.validate().then((res) => {
    addFormEntry({ ...formState }).then((res) => {
      let tip = formState.id ? '修改成功' : '添加成功'
      message.success(tip)
      open.value = false
      emit('ok')
    })
  })
}
function convertNumberStringsToObjectNumbers(obj) {
  // 创建一个新对象来存放转换后的属性
  const newObj = {}

  // 遍历传入对象的所有可枚举属性
  for (let key in obj) {
    if (obj.hasOwnProperty(key)) {
      // 确保属性是对象自身的属性，而不是继承来的
      const value = obj[key]

      // 使用Number()函数尝试将属性值转换为数字
      // 如果转换成功（即不是NaN），并且原始值是字符串类型
      // 则将转换后的数字赋值给新对象的同名属性
      if (!isNaN(Number(value)) && typeof value === 'string' && value != '') {
        newObj[key] = Number(value)
      } else {
        // 如果转换失败或者原始值不是字符串，则直接复制原始值到新对象
        newObj[key] = value
      }
    }
  }

  // 返回新对象
  return newObj
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.form-tip {
  color: #8c8c8c;
  margin-top: 8px;
}
::v-deep(.ant-form-item) {
  margin-bottom: 16px;
}
.field-name {
  width: calc(100%);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
