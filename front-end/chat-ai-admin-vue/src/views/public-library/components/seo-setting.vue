<template>
  <div>
    <a-modal
      v-model:open="show"
      title="SEO设置"
      @ok="handleOk"
      width="746px"
      :confirmLoading="confirmLoading"
    >
      <div class="seo-setting">
        <a-alert
          message="修改后，也需要发布才会生效。"
          type="info"
          show-icon
          style="margin-bottom: 16px"
        />
        <a-form :model="formData" layout="vertical">
          <a-form-item label="Title">
            <a-input v-model:value="formData.seo_title" placeholder="请输入页面标题" />
          </a-form-item>
          <a-form-item label="Description">
            <a-textarea v-model:value="formData.seo_desc" placeholder="请输入页面描述" />
          </a-form-item>
          <a-form-item label="Keyword">
            <a-input
              v-model:value="formData.seo_keywords"
              placeholder="请输入关键词,多个关键词用英文逗号分隔"
            />
          </a-form-item>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { saveLibDocSeo } from '@/api/public-library/index'
import { ref, toRaw } from 'vue'
import { message } from 'ant-design-vue'

const emit = defineEmits(['ok'])

const confirmLoading = ref(false)
// 表单数据
const formData = ref({
  library_key: '',
  seo_title: '',
  seo_desc: '',
  seo_keywords: ''
})

const show = ref(false)

const handleOk = () => {
  confirmLoading.value = true

  saveLibDocSeo({
    ...toRaw(formData.value)
  })
    .then(() => {
      confirmLoading.value = false
      show.value = false
      message.success('保存成功')
      emit('ok')
    })
    .catch(() => {
      confirmLoading.value = false
    })
}

const open = (params) => {
  let defaultData = {
    library_key: '',
    seo_title: '',
    seo_desc: '',
    seo_keywords: ''
  }

  formData.value = { ...defaultData, ...params }

  show.value = true
}

defineExpose({
  open
})
</script>

<style lang="less" scoped>
.seo-setting {
  padding: 20px;
}
</style>
