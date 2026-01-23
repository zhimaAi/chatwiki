<template>
  <div>
    <a-modal
      v-model:open="show"
      :title="t('title')"
      @ok="handleOk"
      width="746px"
      :confirmLoading="confirmLoading"
    >
      <div class="seo-setting">
        <a-form :model="formData" layout="vertical">
          <a-form-item :label="t('label_title')">
            <a-input v-model:value="formData.seo_title" :placeholder="t('placeholder_title')" />
          </a-form-item>
          <a-form-item :label="t('label_desc')">
            <a-textarea v-model:value="formData.seo_desc" :placeholder="t('placeholder_desc')" />
          </a-form-item>
          <a-form-item :label="t('label_keywords')">
            <a-input
              v-model:value="formData.seo_keywords"
              :placeholder="t('placeholder_keywords')"
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
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.public-library.components.seo-setting')

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
      message.success(t('save_success'))
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
