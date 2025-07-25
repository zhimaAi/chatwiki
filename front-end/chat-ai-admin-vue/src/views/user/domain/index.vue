<style lang="less" scoped>
.domain-page {
  width: 100%;
  height: 100%;
  padding: 24px;
  overflow-y: auto;
  background-color: #fff;

  .page-title {
    line-height: 24px;
    margin-bottom: 24px;
    font-size: 16px;
    font-weight: 600;
  }

  .domain-tips {
    margin-bottom: 16px;
    .tip-text {
      font-size: 14px;
    }
  }

  .list-action-box {
    margin-bottom: 8px;
  }

  .list-box {
    .action-box {
      display: flex;
      .action-btn {
        margin-right: 16px;
      }
    }
  }
}
</style>

<template>
  <div class="domain-page">
    <div class="page-container">
      <div class="page-title">自定义域名设置</div>

      <div class="domain-tips">
        <a-alert type="info" show-icon>
          <template #icon>
            <div>
              <ExclamationCircleFilled style="font-size: 16px" />
            </div>
          </template>
          <template #description>
            <div class="tip-text">
              <div>将域名解析到部署服务器，可以通过自由域名访问对外文档和对外服务web APP。</div>
              <!-- <div>
                1、域名解析请参考帮助文档，目前只支持阿里云备案的域名，将域名cname到wiki.aishipinhao.com。
              </div> -->
              <div>1、如需使用https,请在添加域名后，上传证书文件，包括公钥和私钥(ky文件)。</div>
            </div>
          </template>
        </a-alert>
      </div>

      <div class="list-action-box">
        <div>
          <a-button type="primary" @click="handleAddDomain">添加自定义域名</a-button>
        </div>
      </div>

      <div class="list-box">
        <a-table :columns="columns" :data-source="domainList">
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'url'">
              <span>{{ record.url }}</span>
            </template>
            <template v-else-if="column.key === 'is_upload'">
              <span> {{ record.is_upload ? '已上传' : '未上传' }}</span>
            </template>
            <template v-else-if="column.key === 'action'">
              <span class="action-box">
                <a class="action-btn" @click="handleUploadSSL(record)">上传证书</a>
                <a class="action-btn" @click="handleUploadValidationFile(record)">上传验证文件</a>
                <a class="action-btn" @click="handleEditDomain(record)">编辑</a>
                <a class="action-btn" @click="handleDelDomain(record)">删除</a>
              </span>
            </template>
          </template>
        </a-table>
      </div>
    </div>
    <AddDomainModal
      ref="addDomainModalRef"
      :confirmLoading="confirmLoading"
      @ok="handleSaveDomain"
    />
    <UploadSSL
      ref="uploadSSLRef"
      :confirm-loading="uploadSSLLoading"
      @ok="handleUploadSSLConfirm"
    />
    <UploadValidationFile
      ref="uploadValidationFileRef"
      :confirm-loading="uploadValidationFileLoading"
      @ok="saveValidationFile"
    />
  </div>
</template>

<script setup>
import {
  getDomainList,
  saveDomain,
  deleteDomain,
  uploadCertificate,
  uploadCheckFile
} from '@/api/user'
import { ref, createVNode } from 'vue'
import { ExclamationCircleFilled, ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { message, Modal } from 'ant-design-vue'
import AddDomainModal from './components/add-domain-modal .vue'
import UploadSSL from './components/upload-ssl.vue'
import UploadValidationFile from './components/upload-validation-file.vue'

const addDomainModalRef = ref(null)
const confirmLoading = ref(false)
const columns = [
  {
    title: '域名',
    dataIndex: 'url',
    key: 'url'
  },
  {
    title: '证书',
    dataIndex: 'is_upload',
    key: 'is_upload',
    width: 300
  },
  {
    title: '操作',
    key: 'action'
  }
]

const handleAddDomain = () => {
  addDomainModalRef.value.open()
}

const handleEditDomain = (record) => {
  let formData = {
    id: record.id,
    url: record.url,
    protocol: 'http:'
  }

  try {
    const url = new URL(record.url)
    if (url.protocol == 'http:' || url.protocol == 'https:') {
      formData.protocol = url.protocol
    } else {
      formData.protocol = 'http:'
    }

    formData.url = url.hostname
  } catch (e) {
    // console.log(e)
  }

  addDomainModalRef.value.open(formData)
}

const handleSaveDomain = (formData) => {
  confirmLoading.value = true

  let url = `${formData.protocol}//${formData.url}`

  saveDomain({ url: url, id: formData.id || undefined })
    .then(() => {
      message.success('保存成功')
      confirmLoading.value = false
      addDomainModalRef.value.close()
      getList()
    })
    .catch(() => {
      confirmLoading.value = false
    })
}

const handleDelDomain = (record) => {
  Modal.confirm({
    title: '删除',
    icon: createVNode(ExclamationCircleOutlined),
    content: '确定要删除吗？',
    onOk() {
      deleteDomain({ id: record.id }).then(() => {
        message.success('删除成功')
        getList()
      })
    },
    onCancel() {}
  })
}

const uploadSSLRef = ref(null)
const uploadSSLLoading = ref(false)

const handleUploadSSL = (record) => {
  uploadSSLRef.value.open(record)
}

const handleUploadSSLConfirm = (record) => {
  uploadSSLLoading.value = true

  let data = {
    id: record.id,
    ssl_certificate: record.ssl_certificate,
    ssl_certificate_key: record.ssl_certificate_key
  }

  uploadCertificate(data)
    .then(() => {
      uploadSSLLoading.value = false
      uploadSSLRef.value.close()
      message.success('上传成功')
      getList()
    })
    .catch(() => {
      uploadSSLLoading.value = false
    })
}

const uploadValidationFileRef = ref(null)
const uploadValidationFileLoading = ref(false)

const handleUploadValidationFile = (record) => {
  uploadValidationFileRef.value.open(record)
}

const saveValidationFile = (record) => {
  uploadValidationFileLoading.value = true
  uploadCheckFile({ ...record })
    .then(() => {
      uploadValidationFileLoading.value = false
      uploadValidationFileRef.value.close()

      message.success('上传成功')

      getList()
    })
    .catch(() => {
      uploadValidationFileLoading.value = false
    })
}

const domainList = ref([])

const getList = async () => {
  const res = await getDomainList()
  domainList.value = res.data || []
}

getList()
</script>
