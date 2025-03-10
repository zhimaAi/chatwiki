<template>
  <div>
    <a-modal
      v-model:open="open"
      title="导入数据"
      :confirm-loading="confirmLoading"
      @ok="handleOk"
      :footer="null"
      :width="746"
    >
      <a-alert show-icon class="alert-box">
        <template #message>
          <div>
            支持通过Excel或者json文件导入数据。
            <a @click="exportExcelTemplate" class="mr8">excel模板</a>
            <a @click="exportJsonTemplate">json模板</a>
          </div>
          <div>
            导入时，如果d为空，则新建数据。如果id有值，则更新d对应的数据。id通过导出数据获取。单次导入上限1万条数据。
          </div>
        </template>
      </a-alert>
      <div class="upload-box">
        <a-upload-dragger
          v-model:fileList="fileList"
          name="form_files"
          :multiple="false"
          accept=".json,.xlsx,.csv"
          :headers="{ token: userStore.getToken }"
          :data="{ form_id: query.form_id }"
          action="/manage/uploadFormFile"
          @change="handleChange"
          @drop="handleDrop"
        >
          <p class="ant-upload-drag-icon">
            <inbox-outlined></inbox-outlined>
          </p>
          <p class="ant-upload-text">点击或将文件拖拽到这里上传</p>
          <p class="ant-upload-hint">仅支持：json.xlsx.csv文件，请注意必填项</p>
        </a-upload-dragger>
      </div>
    </a-modal>
    <a-modal v-model:open="resultOpen" title="导入数据" :footer="null" :width="746">
      <div class="progress-box" v-if="percent < 100">
        <a-progress type="circle" :percent="percent" />
        <div class="tip">数据导入中，请勿关闭页面</div>
      </div>
      <a-result v-else status="success" title="导入完成">
        <template #subTitle
          >导入成功{{ resultInfo.success }}，导入失败{{ resultInfo.total - resultInfo.success }}
          <span v-if="resultInfo.total - resultInfo.success > 0">
            您可以下载失败数据，修改后重新导入
          </span>
        </template>
        <template #extra>
          <a-button
            v-if="resultInfo.total - resultInfo.success > 0"
            @click="downFailData"
            type="primary"
            >下载失败数据</a-button
          >
        </template>
      </a-result>
    </a-modal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { InboxOutlined, EyeOutlined, DeleteOutlined } from '@ant-design/icons-vue'
import { useUserStore } from '@/stores/modules/user'
import { tableToExcel, exportToJsonWithSaver } from '@/utils/index'
import { useRoute } from 'vue-router'
import dayjs from 'dayjs'
import { uploadFormFile, getUploadFormFileProc } from '@/api/database'
import { message } from 'ant-design-vue'
import { reactive } from 'vue'
const userStore = useUserStore()
const emit = defineEmits(['ok'])

const query = useRoute().query

const props = defineProps({
  column: {
    type: Array,
    default: () => []
  }
})

const open = ref(false)
const confirmLoading = ref(false)

const resultInfo = reactive({
  total: 0,
  success: 0,
  err_data: []
})

const fileList = ref([])
const handleChange = (info) => {
  const status = info.file.status
  if (status !== 'uploading') {
  }
  if (status === 'done') {
    if (info.file.response.res == 0) {
      let task_id = info.file.response.data.task_id
      getProgressInfo(task_id)
    } else {
      fileList.value = []
      return message.error(info.file.response.msg)
    }
  } else if (status === 'error') {
    fileList.value = []
    message.error(`${info.file.name} file upload failed.`)
  }
}
function handleDrop(e) {
  console.log(e)
}
let timer = null
const getProgressInfo = (task_id) => {
  timer && clearInterval(timer)
  open.value = false
  resultOpen.value = true
  percent.value = 10
  resultInfo.total = 0
  resultInfo.success = 0
  resultInfo.err_data = []
  timer = setInterval(() => {
    getUploadFormFileProc({ task_id })
      .then((res) => {
        let data = res.data
        percent.value = (data.processed / data.total).toFixed()
        if (data.finish) {
          emit('ok')
          clearInterval(timer)
          percent.value = 100
        }
        if (data.total > 0) {
          resultInfo.total = data.total
          resultInfo.success = data.success
          resultInfo.err_data = data.err_data
        }
      })

      .catch(() => {
        clearInterval(timer)
      })
  }, 2000)
}

const exportExcelTemplate = () => {
  let str
  str = props.column.map((item) => item.name)
  str = 'id,' + str.join(',') + '\n'
  tableToExcel(str, [], [], '导入数据模板.csv')
}

const exportJsonTemplate = () => {
  let datas = {
    id: ''
  }
  props.column.forEach((item) => {
    datas[item.name] = ''
  })
  exportToJsonWithSaver([datas], '导入数据模板.json')
}

const show = () => {
  fileList.value = []
  open.value = true
}
const handleOk = () => {}

const resultOpen = ref(false)

const percent = ref(10)

const downFailData = () => {
  let str = 'id,'
  let fieds = ['id']
  for (let key in resultInfo.err_data[0]) {
    if (key != 'id' && key != 'err_msg') {
      str += key + ','
      fieds.push(key)
    }
  }
  str = str + 'err_msg\n'
  fieds.push('err_msg')
  let jsonData = resultInfo.err_data
  let name = '失败数据_' + dayjs().format('YYYY/MM/DD HH:mm') + '.xlsx'
  tableToExcel(str, jsonData, fieds, name)
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.alert-box {
  align-items: baseline;
}
.mr8 {
  margin-right: 8px;
}

.upload-box {
  padding: 24px;
}

.progress-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  margin: 100px 0;
  color: #8c8c8c;
}
</style>
