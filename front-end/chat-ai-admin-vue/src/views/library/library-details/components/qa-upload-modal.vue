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
              通过word或者Excel形式导入问答，建议您参考<a
                href="https://xkf-upload-oss.xiaokefu.com.cn/static/chat-wiki/template/FAQ导入模板.docx"
                >word导入模板</a
              >&nbsp;
              <a href="https://xkf-upload-oss.xiaokefu.com.cn/static/chat-wiki/template/FAQ 导入模板.xlsx"
                >Excel导入模板</a
              >创建导入文档
            </div>
            <div>如果您不是使用导入模板创建的文档，可以点击高级设置指定提取规则</div>
            <div>知识库已有相同的问题时，会覆盖原有问答</div>
          </template>
        </a-alert>
      </div>
      <a-flex align="center" style="margin: 16px 0">
        <div>所属分组：</div>
        <a-select v-model:value="formState.group_id" style="flex: 1" placeholder="请选择分组">
          <a-select-option
            v-for="item in props.groupLists.filter((item) => item.id >= 0)"
            :value="item.id"
            >{{ item.group_name }}</a-select-option
          >
        </a-select>
      </a-flex>
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
              <div class="form-item-label required">问题所在列：</div>
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
              <div class="form-item-label">相似问法所在列：</div>
              <div class="form-item-body">
                <a-select
                  v-model:value="formState.similar_column"
                  placeholder="请选择相似问题所在列"
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
              <div class="form-item-label required">答案所在列：</div>
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
              <div class="form-item-label">相似问法开始标识符：</div>
              <div class="form-item-body">
                <a-input placeholder="请输入标识符" v-model:value="formState.similar_label" />
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
import { ref, reactive, computed } from 'vue'
import UploadFilesInput from '../../add-library/components/upload-input.vue'
import { DownOutlined, UpOutlined } from '@ant-design/icons-vue'
import { readLibFileExcelTitle, addLibraryFile } from '@/api/library/index'
import { message } from 'ant-design-vue'
import { useRoute, useRouter } from 'vue-router'
import { useLibraryStore } from '@/stores/modules/library'
const libraryStore = useLibraryStore()
const { qa_index_type } = libraryStore
console.log(qa_index_type,'qa_index_type==')
const router = useRouter()
const rotue = useRoute()
const query = rotue.query
const open = ref(false)
const confirmLoading = ref(false)
const fileList = ref([])

const emit = defineEmits(['ok'])
const excellQaLists = ref([])

const props = defineProps({
  library_id: {
    type: [Number, String],
    default: () => ''
  },
  groupLists: {
    type: Array,
    default: () => []
  },
})

const libraryId = computed(() => {
  return props.library_id || query.id
})


const formState = reactive({
  question_column: void 0,
  answer_column: void 0,
  similar_column: void 0,
  question_lable: '问题：',
  answer_lable: '答案：',
  similar_label: '相似问法：',
  group_id: void 0,
})

const isHide = ref(true)
const show = (groupId) => {
  formState.group_id = void 0
  if(groupId >= 0){
    formState.group_id = groupId
  }
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

  formData.append('library_id', libraryId.value)
  let isTableType = false
  fileList.value.forEach((file) => {
    if (file.name.includes('.xlsx') || file.name.includes('.csv')) {
      isTableType = true
    }
    formData.append('library_files', file)
  })
  if (isTableType) {
    if(!formState.question_column){
      confirmLoading.value = false
      return message.error('请选择问题所在列')
    }
    if(!formState.answer_column){
      confirmLoading.value = false
      return message.error('请选择答案所在列')
    }
    formData.append('question_column', formState.question_column)
    formData.append('answer_column', formState.answer_column)
    formState.similar_column && formData.append('similar_column', formState.similar_column)
  } else {
    formData.append('question_lable', formState.question_lable)
    formData.append('answer_lable', formState.answer_lable)
    formData.append('similar_label', formState.similar_label)
  }
  formData.append('group_id', formState.group_id >= 0 ? formState.group_id : 0)
  formData.append('is_qa_doc', 1)
  formData.append('qa_index_type', qa_index_type)
  addLibraryFile(formData)
    .then((res) => {
      emit('ok')
      open.value = false
      fileList.value = []
      confirmLoading.value = false
      // if (isTableType) {
      //   router.push('/library/document-segmentation?document_id=' + res.data.file_ids[0])
      // }
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
  formState.similar_column = void 0
  formState.question_lable = '问题：'
  formState.answer_lable = '答案：'
  formState.similar_label = '相似问法：'
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
  flex-wrap: wrap;
  .form-item {
    width: calc(50% - 6px);
    display: flex;
    align-items: center;
    .form-item-body {
      flex: 1;
    }
  }
  .form-item-label.required{

    &::before{
      content: '*';
      margin-right: 2px;
      color: red;
      font-weight: 600;
      font-size: 14px;
    }
  }
}
</style>
