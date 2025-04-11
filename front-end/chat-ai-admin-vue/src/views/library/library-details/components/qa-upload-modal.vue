<template>
  <div>
    <a-modal
      :width="680"
      v-model:open="open"
      :confirm-loading="confirmLoading"
      :maskClosable="false"
      title="上传文档"
      @ok="handleSaveFiles"
      @cancel="handleCloseFileUploadModal"
    >
      <div class="alert-box">
        <a-alert show-icon>
          <template #message>
            <div>
              通过word或者EXce形式导入问答对，建议您参考<a
                href="https://xkf-upload-oss.xiaokefu.com.cn/static/chat-wiki/FAQ导入模板.docx"
                >word导入模板</a
              >&nbsp;
              <a href="https://xkf-upload-oss.xiaokefu.com.cn/static/chat-wiki/FAQ 导入模板.xlsx"
                >Excel导入模板</a
              >创建导入文档
            </div>
            <div>如果您不是使用导入模板创建的文档，可以点击高级设置指定提取规则</div>
          </template>
        </a-alert>
      </div>
      <div class="upload-file-box">
        <UploadFilesInput :type="2" v-model:value="fileList" @change="onFilesChange" />
      </div>
      <div class="hight-set" v-if="fileList.length > 0">
        <a @click="isHide = !isHide"
          >高级设置：

          <DownOutlined v-if="!isHide" />
          <UpOutlined v-else />
        </a>
        <template v-if="isHide">
          <div class="set-box" v-if="fileType == 1">
            <div class="form-item">
              <div class="form-item-label">问题所在列：</div>
              <div class="form-item-body">
                <a-select
                  v-model:value="formState.question_column"
                  placeholder="请选择列名"
                  style="width: 100%"
                >
                  <a-select-option
                    v-for="item in excellQaLists"
                    :value="item.value"
                    :key="item.value"
                    >{{ item.lable }}</a-select-option
                  >
                </a-select>
              </div>
            </div>
            <div class="form-item">
              <div class="form-item-label">答案所在列：</div>
              <div class="form-item-body">
                <a-select
                  v-model:value="formState.answer_column"
                  placeholder="请选择答案所在列"
                  style="width: 100%"
                >
                  <a-select-option
                    v-for="item in excellQaLists"
                    :value="item.value"
                    :key="item.value"
                    >{{ item.lable }}</a-select-option
                  >
                </a-select>
              </div>
            </div>
          </div>
          <div class="set-box" v-if="fileType == 2">
            <div class="form-item">
              <div class="form-item-label">问题开始标识符：</div>
              <div class="form-item-body">
                <a-input placeholder="请输入标识符" v-model:value="formState.question_lable" />
              </div>
            </div>
            <div class="form-item">
              <div class="form-item-label">答案开始标识符：</div>
              <div class="form-item-body">
                <a-input placeholder="请输入标识符" v-model:value="formState.answer_lable" />
              </div>
            </div>
          </div>
        </template>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import UploadFilesInput from '../../add-library/components/upload-input.vue'
import { DownOutlined, UpOutlined } from '@ant-design/icons-vue'
import { readLibFileExcelTitle, addLibraryFile } from '@/api/library/index'
import { message } from 'ant-design-vue'
import { useRoute, useRouter } from 'vue-router'
const router = useRouter()
const rotue = useRoute()
const query = rotue.query
const open = ref(false)
const confirmLoading = ref(false)
const fileList = ref([])

const emit = defineEmits(['ok'])
const excellQaLists = ref([])

const formState = reactive({
  question_column: void 0,
  answer_column: void 0,
  question_lable: '问题：',
  answer_lable: '答案：'
})

const isHide = ref(true)
const show = () => {
  open.value = true
}

const getExcelQaTitle = () => {
  // 获取excel的QA  问题所在列 下拉列表
  let formData = new FormData()
  fileList.value.forEach((file) => {
    formData.append('library_files', file)
  })
  readLibFileExcelTitle(formData).then((res) => {
    let datas = []
    for (let key in res.data) {
      datas.push({
        lable: res.data[key],
        value: key
      })
    }
    excellQaLists.value = datas
  })
}

const handleSaveFiles = () => {
  if (fileList.value.length == 0) {
    message.error('请选择文件')
    return
  }

  confirmLoading.value = true

  let formData = new FormData()

  formData.append('library_id', query.id)
  let isTableType = false
  fileList.value.forEach((file) => {
    if (file.name.includes('.xlsx') || file.name.includes('.csv')) {
      isTableType = true
    }
    formData.append('library_files', file)
  })
  if (isTableType) {
    formData.append('question_column', formState.question_column)
    formData.append('answer_column', formState.answer_column)
  } else {
    formData.append('question_lable', formState.question_lable)
    formData.append('answer_lable', formState.answer_lable)
  }
  formData.append('is_qa_doc', 1)
  addLibraryFile(formData)
    .then((res) => {
      emit('ok')
      open.value = false
      fileList.value = []
      confirmLoading.value = false
      if (isTableType) {
        router.push('/library/document-segmentation?document_id=' + res.data.file_ids[0])
      }
    })
    .catch(() => {
      confirmLoading.value = false
    })
}

const handleCloseFileUploadModal = () => {
  fileList.value = []
}

const fileType = ref(1) // 1表格  2doc
const onFilesChange = (files) => {
  formState.question_column = void 0
  formState.answer_column = void 0
  formState.question_lable = '问题：'
  formState.answer_lable = '答案：'
  fileList.value = files
  if (files && files.length > 0) {
    if (files[0].type.includes('.document')) {
      fileType.value = 2
    } else {
      fileType.value = 1
      getExcelQaTitle()
    }
  }
}
defineExpose({
  show
})
</script>

<style lang="less" scoped>
.alert-box {
  margin: 16px 0;
  ::v-deep(.ant-alert) {
    align-items: baseline;
  }
  ::v-deep(.anticon.anticon-info-circle) {
    position: relative;
    top: 2px;
  }
}

.hight-set {
  margin-top: 16px;
}
.set-box {
  padding: 12px 16px;
  background: #f2f2f2;
  border-radius: 6px;
  margin-top: 6px;
  display: flex;
  align-items: center;
  gap: 12px;
  .form-item {
    flex: 1;
    display: flex;
    align-items: center;
    .form-item-body {
      flex: 1;
    }
  }
}
</style>
