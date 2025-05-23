<template>
  <div bg-color="#f5f9ff" class="library-page">
    <div class="page-title">对外文档</div>

    <page-alert style="margin-bottom: 16px;" title="使用说明">
      <div>
        <p>
          1、对外文档作为在线文档创作管理工具，创作的文档既可作为知识库关联到机器人使用，也可分享给好友查看。
        </p>
        <p>2、支持设置每篇文档的title、description、keyword，提高SEO搜索权重。</p>
      </div>
    </page-alert>
    
    <div class="library-page-body">
      <div class="list-toolbar">
        <div class="toolbar-box">
          <h3 class="list-total">全部 ({{ list.length }})</h3>
        </div>

        <div class="toolbar-box">
          <a-button type="primary" @click="handleAdd()">
            <template #icon>
              <PlusOutlined />
            </template>
            新建对外文档
          </a-button>
        </div>
      </div>
      <div>
        <LibraryList :list="list" @add="handleAdd" @edit="toEdit" @delete="handleDelete" />
      </div>
    </div>
    <AddLibraryModel ref="addLibraryModelRef" />
  </div>
</template>

<script setup>
import { ref, createVNode } from 'vue'
import { useRouter } from 'vue-router'
import { getLibraryList, deleteLibrary } from '@/api/library'
import { formatFileSize } from '@/utils/index'
import { LIBRARY_NORMAL_AVATAR, LIBRARY_OPEN_AVATAR } from '@/constants/index'
import { Modal, message } from 'ant-design-vue'
import { ExclamationCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import LibraryList from './components/libray-list/index.vue'
import AddLibraryModel from '../add/add-library-model.vue'
import PageAlert from '@/components/page-alert/page-alert.vue'

const router = useRouter()

const addLibraryModelRef = ref(null)

const list = ref([])

const getList = () => {
  getLibraryList({ type: 1 }).then((res) => {
    let data = res.data || []

    data.forEach((item) => {
      item.file_size_str = formatFileSize(item.file_size)

      if (!item.avatar) {
        item.avatar = item.type == 0 ? LIBRARY_NORMAL_AVATAR : LIBRARY_OPEN_AVATAR
      }
    })

    list.value = data
  })
}

getList()

const handleAdd = () => {
  addLibraryModelRef.value.show()
}

const toEdit = (data) => {
  // router.push({
  //   path: '/public-library/home',
  //   query: {
  //     library_id: data.id,
  //     library_key: data.library_key
  //   }
  // })
  window.open(`/#/public-library/home?library_id=${data.id}&library_key=${data.library_key}`, "_blank", "noopener") // 建议添加 noopener 防止安全漏洞
}

const handleDelete = (data) => {
  let secondsToGo = 3

  let modal = Modal.confirm({
    title: `删除${data.library_name}`,
    icon: createVNode(ExclamationCircleOutlined),
    content: '您确定要删除此知识库吗？',
    okText: secondsToGo + ' 确 定',
    okType: 'danger',
    okButtonProps: {
      disabled: true
    },
    cancelText: '取 消',
    onOk() {
      onDelete(data)
    },
    onCancel() {
      // console.log('Cancel')
    }
  })

  let interval = setInterval(() => {
    if (secondsToGo == 1) {
      modal.update({
        okText: '确 定',
        okButtonProps: {
          disabled: false
        }
      })

      clearInterval(interval)
      interval = undefined
    } else {
      secondsToGo -= 1

      modal.update({
        okText: secondsToGo + ' 确 定',
        okButtonProps: {
          disabled: true
        }
      })
    }
  }, 1000)
}

const onDelete = ({ id }) => {
  deleteLibrary({ id }).then(() => {
    message.success('删除成功')
    getList()
  })
}
</script>

<style lang="less" scoped>
.library-page {
  .page-title{
    margin: 16px 0;
    font-size: 20px;
    font-weight: 600;
    line-height: 28px;
    color: #000000;
  }

  .list-toolbar{
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;

    .list-total{
      line-height: 24px;
      font-size: 16px;
      font-weight: 600;
      color: #262626;
    }
  }
}
// 大于1920px
@media screen and (min-width: 1920px) {
  .library-page {
  }
}
</style>
