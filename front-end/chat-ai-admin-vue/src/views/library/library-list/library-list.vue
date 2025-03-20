<template>
  <div bg-color="#f5f9ff" class="library-page">
    <div class="library-page-body">
      <a-tabs v-model:activeKey="activeKey" @change="onChangeTab">
        <a-tab-pane key="all" tab="全部"> </a-tab-pane>
        <a-tab-pane key="0" tab="普通知识库" force-render> </a-tab-pane>
        <a-tab-pane key="2" tab="问答知识库"> </a-tab-pane>
      </a-tabs>

      <div>
        <LibraryList
          :show-create="libraryCreate"
          :list="list"
          @add="handleAdd"
          @edit="toEdit"
          @delete="handleDelete"
        />
      </div>
    </div>
    <AddLibrayPopup ref="addLibrayPopup" @ok="toAdd" />
    <AddLibraryModel ref="addLibraryModelRef" />
  </div>
</template>

<script setup>
import { ref, createVNode, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Modal, message } from 'ant-design-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import LibraryList from './components/libray-list/index.vue'
import AddLibrayPopup from './components/add-libray-popup.vue'
import { getLibraryList, deleteLibrary } from '@/api/library'
import { formatFileSize } from '@/utils/index'
import { usePermissionStore } from '@/stores/modules/permission'
import { DEFAULT_LIBRARY_AVATAR, DEFAULT_LIBRARY_AVATAR3 } from '@/constants/index'
import AddLibraryModel from '@/views/library/add-library/add-library-model.vue'

const permissionStore = usePermissionStore()
let { role_permission } = permissionStore
const libraryCreate = computed(() => role_permission.includes('LibraryCreate'))

const router = useRouter()

const addLibraryModelRef = ref(null)
const addLibrayPopup = ref(null)

const handleAdd = () => {
  addLibrayPopup.value.show()
}

const toAdd = (val) => {
  addLibraryModelRef.value.show({ type: val })
}

const activeKey = ref('all')

const list = ref([])

const getList = () => {
  let type = activeKey.value === 'all' ? '' : activeKey.value

  getLibraryList({ type }).then((res) => {
    let data = res.data || []

    data.forEach((item) => {
      item.file_size_str = formatFileSize(item.file_size)

      if (!item.avatar) {
        item.avatar = item.type == 0 ? DEFAULT_LIBRARY_AVATAR : DEFAULT_LIBRARY_AVATAR3
      }
    })

    list.value = data
  })
}

getList()

const onChangeTab = () => {
  getList()
}

const toEdit = (data) => {
  if (data.type == '1') {
    router.push({
      path: '/public-library/config',
      query: {
        library_id: data.id
      }
    })
  } else {
    router.push({
      name: 'libraryDetails',
      query: {
        id: data.id
      }
    })
  }
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
