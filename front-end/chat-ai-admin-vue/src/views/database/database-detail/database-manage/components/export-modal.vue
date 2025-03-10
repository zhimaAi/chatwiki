<template>
  <div>
    <a-modal v-model:open="open" :title="modalTitle" @ok="handleOk" :width="450">
      <a-form class="mt16" layout="vertical">
        <a-form-item label="">
          <a-radio-group v-model:value="formState.type" name="radioGroup">
            <a-radio :value="1">全部数据</a-radio>
            <a-radio :value="2">自定义日期</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item
          required
          label="选择导出日期"
          v-show="formState.type == 2"
          v-bind="validateInfos.dates"
        >
          <a-range-picker
            style="width: 100%"
            v-model:value="formState.dates"
            valueFormat="YYYY-MM-DD"
            format="YYYY-MM-DD"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { PlusOutlined, QuestionCircleOutlined } from '@ant-design/icons-vue'
import { reactive, ref } from 'vue'
import { Form, message } from 'ant-design-vue'
import { useRoute } from 'vue-router'
import { exportFormEntry } from '@/api/database'
const rotue = useRoute()
const query = rotue.query
const modalTitle = ref('导出数据')
const open = ref(false)
const formState = reactive({
  dates: [],
  type: 1
})
const show = () => {
  open.value = true
  formState.type = 1
  formState.dates = []
}
const formRules = reactive({
  dates: [
    {
      validator: async (rule, value) => {
        if (formState.type == 2 && !value.length) {
          return Promise.reject('请选择导出日期')
        }
        return Promise.resolve()
      }
    }
  ]
})
const useForm = Form.useForm
const { resetFields, validate, validateInfos } = useForm(formState, formRules)
const handleOk = () => {
  validate().then(() => {
    let parmas = {
      form_id: query.form_id
    }
    if (formState.type == 2) {
      parmas.start_date = formState.dates[0]
      parmas.end_date = formState.dates[1]
    }
    exportFormEntry(parmas)
    open.value = false
  })
}
defineExpose({
  show
})
</script>

<style lang="less" scoped>
.mt16 {
  margin-top: 16px;
}
</style>
