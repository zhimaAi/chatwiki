<template>
  <div class="add-data-sheet">
    <a-modal v-model:open="open" :title="modalTitle" @ok="handleOk" :width="540">
      <a-form layout="vertical">
        <a-form-item v-bind="validateInfos.name">
          <template #label
            >字段名称&nbsp;
            <a-tooltip>
              <template #title>定义数据表的表头，可以在对应表头下存储相关数据。</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>
          <a-input :maxLength="64" v-model:value="formState.name" placeholder="请输入字段名称" />
        </a-form-item>
        <a-form-item v-bind="validateInfos.description">
          <template #label
            >字段描述&nbsp;
            <a-tooltip>
              <template #title>表头字段的说明，帮助用户或大模型理解表头字段</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>
          <a-input :maxLength="64" v-model:value="formState.description" placeholder="请输入字段描述" />
        </a-form-item>
        <a-form-item v-bind="validateInfos.type">
          <template #label
            >数据类型&nbsp;
            <a-tooltip>
              <template #title>选择存储字段对应的数据类型</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>
          <a-select v-model:value="formState.type" style="width: 100%" placeholder="请选择数据类型">
            <a-select-option v-for="item in typeOption" :value="item.value" :key="item.value">{{
              item.label
            }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <template #label
            >是否必要&nbsp;
            <a-tooltip>
              <template #title>
                <div>必要字段：在保存一行数据时，必须提供对应字段信息，否则无法保存该行数据</div>
                <div>非必要字段：缺失该字段信息时，一行数据仍可被保存在表中</div>
              </template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>
          <a-switch
            v-model:checked="formState.required"
            checked-children="开"
            un-checked-children="关"
            checkedValue="true"
            unCheckedValue="false"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { PlusOutlined, QuestionCircleOutlined } from '@ant-design/icons-vue'
import { reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Form, message } from 'ant-design-vue'
import { addFormField } from '@/api/database'
const modalTitle = ref('添加字段')
const rotue = useRoute()
const query = rotue.query
const open = ref(false)
const emit = defineEmits(['ok'])
const typeOption = [
  {
    label: 'string',
    value: 'string'
  },
  {
    label: 'integer',
    value: 'integer'
  },
  {
    label: 'number',
    value: 'number'
  },
  {
    label: 'boolean',
    value: 'boolean'
  }
]
const formState = reactive({
  name: '',
  description: '',
  type: void 0,
  required: 'false',
  form_id: query.form_id,
  id: ''
})
const show = (data) => {
  open.value = true
  resetFields()
  formState.name = data.name || ''
  formState.description = data.description || ''
  formState.type = data.type
  formState.required = data.required
  formState.id = data.id || ''
  modalTitle.value = data.id ? '编辑字段' : '添加字段'
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
      message: '请输入字段描述'
    }
  ],
  type: [
    {
      required: true,
      message: '请选择数据类型'
    }
  ]
})
const useForm = Form.useForm
const { resetFields, validate, validateInfos } = useForm(formState, formRules)
const handleOk = () => {
  validate().then(() => {
    addFormField(formState).then((res) => {
      let tip = formState.id ? '修改成功' : '新增成功'
      message.success(tip)
      open.value = false
      emit('ok')
    })
  })
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
