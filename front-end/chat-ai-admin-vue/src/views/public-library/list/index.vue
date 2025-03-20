<template>
  <div bg-color="#f5f9ff" class="library-page">
    <div class="library-page-body">
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
import { Modal, message } from 'ant-design-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import LibraryList from './components/libray-list/index.vue'
import { getLibraryList, deleteLibrary } from '@/api/library'
import { formatFileSize } from '@/utils/index'
import { DEFAULT_LIBRARY_AVATAR, DEFAULT_LIBRARY_AVATAR2 } from '@/constants/index'
import AddLibraryModel from '../add/add-library-model.vue'

const router = useRouter()

const addLibraryModelRef = ref(null)

const list = ref([])

const getList = () => {
  getLibraryList({ type: 1 }).then((res) => {
    let data = res.data || []

    data.forEach((item) => {
      item.file_size_str = formatFileSize(item.file_size)

      if (!item.avatar) {
        item.avatar = item.type == 0 ? DEFAULT_LIBRARY_AVATAR : DEFAULT_LIBRARY_AVATAR2
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
  router.push({
    path: '/public-library/home',
    query: {
      library_id: data.id,
      library_key: data.library_key
    }
  })
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
  :deep(.ant-tabs-nav) {
    margin-bottom: 8px;
  }
}
// 大于1440px
@media screen and (min-width: 1440px) {
  .library-page {
  }
}
</style>
