<template>
  <div>
    <a-modal v-model:open="open" title="保存到提示词库" :width="580" @ok="handleOk" :confirmLoading="isLoading">
      <a-form
        style="margin-top: 24px"
        :label-col="{ span: 4 }"
        :wrapper-col="{ span: 18 }"
        ref="formRef"
        :model="formState"
      >
        <a-form-item
          name="title"
          label="提示词标题"
          :rules="[{ required: true, message: '请输入提示词标题' }]"
        >
          <a-input
            v-model:value="formState.title"
            placeholder="请输入提示词标题"
            :maxLength="10"
          ></a-input>
        </a-form-item>
        <a-form-item name="group_id" label="分组">
          <a-select v-model:value="formState.group_id" style="width: 100%">
            <a-select-option v-for="item in groupList" :value="item.id">{{
              item.group_name
            }}</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, h, reactive } from 'vue'
import {} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { getPromptLibraryGroup, savePromptLibraryItems } from '@/api/user/index.js'
const open = ref(false)

const formState = reactive({
  title: '',
  group_id: 0,
  prompt_type: 0,
  prompt: '',
  prompt_struct: {}
})

const formRef = ref(null)
const show = (data) => {
  if(data.prompt_type == 0 && !data.prompt){
    return message.error('自定义提示为空,无法上传')
  }
  formRef.value && formRef.value.resetFields()
  let newData = {
    ...data
  }
  formState.title = ''
  formState.group_id = 0
  formState.prompt = ''
  formState.prompt_struct = {}
  formState.prompt_type = newData.prompt_type
  if (newData.prompt_type == 0) {
    formState.prompt = newData.prompt
  } else {
    formState.prompt_struct = JSON.stringify(newData.prompt_struct)
  }
  open.value = true
}

const groupList = ref([])
const getGroupList = () => {
  getPromptLibraryGroup().then((res) => {
    groupList.value = [
      {
        id: 0,
        group_name: '默认分组'
      },
      ...res.data
    ]
  })
}

getGroupList()

const isLoading = ref(false)

const handleOk = () => {
  formRef.value.validate().then(() => {
    let parmas = {
      ...formState
    }
    isLoading.value = true
    savePromptLibraryItems({
      ...parmas
    })
      .then(() => {
        message.success(`保存成功`)
        open.value = false
      })
      .finally(() => {
        isLoading.value = false
      })
  })
}

defineExpose({
  show
})
</script>
<style lang="less" scoped>

</style>
