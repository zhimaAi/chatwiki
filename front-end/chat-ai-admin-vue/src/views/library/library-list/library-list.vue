<template>
  <div bg-color="#f5f9ff" class="library-page">
    <PageTabs class="mb-16" :tabs="pageTabs" active="/library/list"></PageTabs>

    <page-alert class="mb-16" title="使用说明">
      <div>
        <p>
          1、知识库可以关联到聊天机器人中使用，创建机器人之前请先创建知识库，然后去机器人设置中关联。
        </p>
        <p>2、知识库支持普通知识库和问答知识库。普通知识库适用于非结构化数据，支持text/doc/pdf/md/html等格式文件，上传后系统会自动分段处理，也支持自定义分段。问答知识库适用于一问一答形式结构化数据，支持通过excel/word批量上传或自定义添加问题和答案。</p>
      </div>
    </page-alert>
    
    <div class="library-page-body">

      <div class="list-toolbar">
        <div class="toolbar-box">
          <ListTabs :tabs="tabs" v-model:value="activeKey" @change="onChangeTab" />
        </div>

        <div class="toolbar-box">
          <a-dropdown v-if="libraryCreate">
            <a-button type="primary" @click.prevent="" >
              <template #icon>
                <PlusOutlined />
              </template>
              新建知识库
            </a-button>
            <template #overlay>
              <a-menu>
                <a-menu-item @click.prevent="handleAdd(0)">
                  <span class="create-action">
                    <img class="icon" :src="LIBRARY_NORMAL_AVATAR" alt="">
                    <span>普通知识库</span>
                  </span>
                </a-menu-item>
                <a-menu-item @click.prevent="handleAdd(2)">
                  <span class="create-action">
                    <img class="icon" :src="LIBRARY_QA_AVATAR" alt="">
                    <span>问答知识库</span>
                  </span>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
          
        </div>
      </div>

      <div>
        <LibraryList
          :show-create="libraryCreate"
          :list="list"
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
import { getLibraryList, deleteLibrary } from '@/api/library'
import { formatFileSize } from '@/utils/index'
import { usePermissionStore } from '@/stores/modules/permission'
import { LIBRARY_NORMAL_AVATAR, LIBRARY_QA_AVATAR } from '@/constants/index'
import { ExclamationCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import LibraryList from './components/libray-list/index.vue'
import AddLibrayPopup from './components/add-libray-popup.vue'
import AddLibraryModel from '@/views/library/add-library/add-library-model.vue'
import PageTabs from '@/components/cu-tabs/page-tabs.vue'
import PageAlert from '@/components/page-alert/page-alert.vue'
import ListTabs from '@/components/cu-tabs/list-tabs.vue'

const router = useRouter()

const pageTabs = ref([{
  title: '知识库',
  path: '/library/list'
}, {
  title: '数据库',
  path: '/database/list'
}])

const tabs = ref([
  {
    title: '全部 (0)',
    value: 'all'
  },{
    title: '普通知识库 (0)',
    value: '0'
  },
  {
    title: '问答知识库 (0)',
    value: '2'
  }
])

const permissionStore = usePermissionStore()
let { role_permission } = permissionStore
const libraryCreate = computed(() => role_permission.includes('LibraryCreate'))

const addLibraryModelRef = ref(null)
const addLibrayPopup = ref(null)

const handleAdd = (type) => {
  // addLibrayPopup.value.show(type)
  toAdd(type)
}

const toAdd = (val) => {
  addLibraryModelRef.value.show({ type: val })
}

const activeKey = ref('all')

const list = ref([])
const updateTabNumber = () => {
  let all = 0
    let normal = 0
    let qa =  0
  list.value.forEach(item => {
    if (item.type == 0) {
        normal += 1
      }
      if (item.type == 2) {
        qa += 1 
      }
      all += 1
  })

  tabs.value = [
      {
        title: '全部 (' + all + ')',
        value: 'all' 
      },
      {
        title: '普通知识库 (' + normal + ')',
        value: '0' 
      },
      {
        title: '问答知识库 (' + qa + ')',
        value: '2'
      }
    ]
}
const getList = () => {
  let type = activeKey.value === 'all' ? '' : activeKey.value

  getLibraryList({ type }).then((res) => {
    let data = res.data || []

    data.forEach((item) => {
      item.file_size_str = formatFileSize(item.file_size)

      if (!item.avatar) {
        item.avatar = item.type == 0 ? LIBRARY_NORMAL_AVATAR : LIBRARY_QA_AVATAR
      }
    })

    list.value = data;

    if(activeKey.value === 'all'){
      updateTabNumber()
    }
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
  .list-toolbar{
    display: flex;
    justify-content: space-between;
    margin-bottom: 8px;
  }
}
.create-action{
  display: flex;
  align-items: center;
 .icon{
    width: 20px;
    height: 20px;
    margin-right: 8px;
 }
}
// 大于1920px
@media screen and (min-width: 1920px) {
  .library-page {
  }
}
</style>
