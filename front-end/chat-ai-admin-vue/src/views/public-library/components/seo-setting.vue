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
        <a-form :model="formData" layout="vertical">
          <a-form-item label="标题">
            <a-input v-model:value="formData.seo_title" placeholder="请输入页面标题" />
          </a-form-item>
          <a-form-item label="描述">
            <a-textarea v-model:value="formData.seo_desc" placeholder="请输入页面描述" />
          </a-form-item>
          <a-form-item label="关键词">
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
import { ref, reactive, toRaw } from 'vue'
import { message } from 'ant-design-vue'

const emit = defineEmits(['ok'])

const confirmLoading = ref(false)
// 表单数据
const formData = reactive({
  library_key: '',
  seo_title: '',
  seo_desc: '',
  seo_keywords: ''
})

const show = ref(false)

const handleOk = () => {
  confirmLoading.value = true

  saveLibDocSeo({
    ...toRaw(formData)
  })
    .then(() => {
      confirmLoading.value = false
      show.value = false
      message.success('保存成功')
      emit('ok', { ...toRaw(formData) })
    })
    .catch(() => {
      confirmLoading.value = false
    })
}

const open = (val) => {
  formData.library_key = val.library_key || ''
  formData.doc_id = val.doc_id || ''
  formData.seo_title = val.seo_title || ''
  formData.seo_desc = val.seo_desc || ''
  formData.seo_keywords = val.seo_keywords || ''

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
