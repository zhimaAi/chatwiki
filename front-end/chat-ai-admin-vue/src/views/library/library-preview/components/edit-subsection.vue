<template>
  <div>
    <a-modal v-model:open="open" title="编辑分段" @ok="handleOk" width="746px">
      <div class="form-box-wrapper">
        <div class="form-item">
          <div class="form-label">分段标题：</div>
          <div class="form-content">
            <a-input :maxLength="25" v-model:value="title" placeholder="请输入分段标题" />
          </div>
        </div>
        <template v-if="isQaDocment">
          <div class="form-item">
            <div class="form-label required">分段问题：</div>
            <div class="form-content">
              <a-textarea
                placeholder="请输入分段问题"
                v-model:value="question"
                style="height: 100px"
              ></a-textarea>
            </div>
          </div>
          <div class="form-item">
            <div class="form-label required">分段答案：</div>
            <div class="form-content">
              <a-textarea
                placeholder="请输入分段答案"
                v-model:value="answer"
                style="height: 150px"
              ></a-textarea>
              <div v-if="answer.length > 10000" class="error-tip">分段答案最多支持10000个字符</div>
            </div>
          </div>
        </template>
        <div class="form-item" v-else>
          <div class="form-label required">分段内容：</div>
          <div class="form-content">
            <a-textarea
              style="height: 150px"
              v-model:value="content"
              placeholder="请输入分段内容"
            />
            <div v-if="content.length > 10000" class="error-tip">分段内容最多支持10000个字符</div>
          </div>
        </div>
        <div class="form-item">
          <div class="form-label">附件</div>
          <div class="form-content">
            <div class="upload-box-wrapper">
              <a-tabs v-model:activeKey="activeKey" size="small">
                <a-tab-pane key="1">
                  <template #tab>
                    <span>
                      <svg-icon name="img-icon" style="font-size: 14px; color: #2475fc"></svg-icon>
                      图片
                      <span v-if="images.length">({{ images.length }})</span>
                    </span>
                  </template>
                </a-tab-pane>
              </a-tabs>
              <UploadImg v-model:value="images"></UploadImg>
            </div>
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { editParagraph } from '@/api/library'
import { useRoute } from 'vue-router'
import UploadImg from '@/components/upload-img/index.vue'
import { isArray } from 'ant-design-vue/lib/_util/util.js'
const route = useRoute()
const props = defineProps({
  detailsInfo: {
    type: Object
  }
})
const emit = defineEmits(['handleEdit'])
const activeKey = ref('1')
const open = ref(false)
const title = ref('')
const content = ref('')
const answer = ref('')
const question = ref('')
const images = ref([])
const id = ref('')
const isQaDocment = ref(false)
const showModal = (data) => {
  title.value = data.title || ''
  content.value = data.content || ''
  id.value = data.id || ''
  answer.value = data.answer || ''
  images.value = data.images || []
  question.value = data.question || ''
  isQaDocment.value = props.detailsInfo.is_qa_doc == '1'

  open.value = true
}

const handleOk = () => {
  if (!content.value && !isQaDocment.value) {
    return message.error('请输入分段内容')
  }
  if (!question.value && isQaDocment.value) {
    return message.error('请输入分段问题')
  }
  if (!answer.value && isQaDocment.value) {
    return message.error('请输入分段答案')
  }
  if (isQaDocment.value && answer.value.length > 10000) {
    return
  }
  if (!isQaDocment.value && content.value.length > 10000) {
    return
  }
  let data = {
    title: title.value,
    content: content.value,
    question: question.value,
    answer: answer.value,
    images: images.value
  }
  if (id.value) {
    data.id = id.value
  } else {
    data.file_id = route.query.id
  }
  let formData = new FormData()
  for (let key in data) {
    if (isArray(data[key])) {
      data[key].forEach((v) => {
        formData.append(key, v)
      })
    } else {
      formData.append(key, data[key])
    }
  }
  editParagraph(formData).then((res) => {
    message.success(data.id ? '修改成功' : '添加成功')
    open.value = false
    emit('handleEdit', {
      ...data
    })
  })
}
defineExpose({ showModal })
</script>
<style lang="less" scoped>
.form-box-wrapper {
  .form-item {
    margin-top: 16px;
  }
  .form-label {
    color: #262626;
    font-size: 14px;
    line-height: 22px;
    padding-top: 5px;
    &.required::before {
      content: '*';
      display: inline-block;
      color: #fb363f;
      margin-right: 2px;
    }
  }
  .form-content {
    margin-top: 8px;
  }
  .upload-box-wrapper {
    background: #f2f4f7;
    border-radius: 6px;
    &::v-deep(.ant-tabs-nav::before) {
      border-color: #f2f4f7;
    }
    &::v-deep(.ant-tabs-nav) {
      margin: 0;
      margin-left: 16px;
    }
  }
}
.error-tip {
  margin-top: 4px;
  color: #fb363f;
}
</style>
