<template>
  <div class="add-data-sheet">
    <a-modal v-model:open="open" :title="modalTitle" @ok="handleOk" :width="540">
      <a-form layout="vertical">
        <a-form-item label="数据表名称" v-bind="validateInfos.name">
          <a-input :maxLength="64" v-model:value="formState.name" placeholder="请输入数据表名称" />
          <div class="form-tip">必须以英文字母开头，只能包含小写字母、数字、下划线</div>
        </a-form-item>
        <a-form-item label="数据表描述" v-bind="validateInfos.description">
          <a-textarea
            :maxLength="500"
            style="height: 80px"
            v-model:value="formState.description"
            placeholder="请输入数据表主要用途，让大模型更加深入理解此表功能"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { Form, message } from 'ant-design-vue'
import { addForm, editForm } from '@/api/database'
import { useRouter } from 'vue-router'
const router = useRouter()
const emit = defineEmits(['ok'])
const open = ref(false)
const modalTitle = ref('新建数据表')
const formState = reactive({
  name: '',
  description: '',
  id: ''
})
const show = (data = {}) => {
  open.value = true
  resetFields()
  formState.name = data.name || ''
  formState.description = data.description || ''
  formState.id = data.id || ''
  modalTitle.value = data.id ? '编辑数据表' : '新建数据表'
}
const formRules = reactive({
  name: [
    {
      required: true,
      validator: async (rule, value) => {
        if (!/^[a-z][a-z0-9_]*$/.test(value)) {
          return Promise.reject('必须以英文字母开头，只能包含小写字母、数字、下划线')
        }
        return Promise.resolve()
      }
    }
  ],
  description: [
    {
      required: true,
      message: '请输入数据表描述'
    }
  ]
})
const useForm = Form.useForm
const { resetFields, validate, validateInfos } = useForm(formState, formRules)
const handleOk = () => {
  validate().then((res) => {
    if (formState.id) {
      editForm(formState).then(() => {
        postCallback('编辑成功')
      })
    } else {
      addForm(formState).then((res) => {
        postCallback('创建成功')
        router.push({
          path: '/database/details/field-manage',
          query: {
            form_id: res.data.id,
          }
        })
      })
    }
  })
}
function postCallback(tip) {
  message.success(tip)
  open.value = false
  emit('ok')
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
</style>
