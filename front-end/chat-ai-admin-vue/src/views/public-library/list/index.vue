<template>
  <div bg-color="#f5f9ff" class="library-page">
    <div class="page-title">{{ t('page_title') }}</div>

    <page-alert style="margin-bottom: 16px;" :title="t('usage_instruction')">
      <div>
        <p>{{ t('usage_instruction_1') }}</p>
        <p>{{ t('usage_instruction_2') }}</p>
      </div>
    </page-alert>

    <div class="library-page-body">
      <div class="list-toolbar">
        <div class="toolbar-box">
          <h3 class="list-total">{{ t('all') }} ({{ list.length }})</h3>
        </div>

        <div class="toolbar-box" v-if="createOpenLibDoc">
          <a-button type="primary" @click="handleAdd()">
            <template #icon>
              <PlusOutlined />
            </template>
            {{ t('new_document') }}
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
import { ref, createVNode, computed } from 'vue'
import { getLibraryList, deleteLibrary } from '@/api/library'
import { formatFileSize } from '@/utils/index'
import { LIBRARY_NORMAL_AVATAR, LIBRARY_OPEN_AVATAR } from '@/constants/index'
import { Modal, message } from 'ant-design-vue'
import { ExclamationCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import LibraryList from './components/libray-list/index.vue'
import AddLibraryModel from '../add/add-library-model.vue'
import PageAlert from '@/components/page-alert/page-alert.vue'
import { usePermissionStore } from '@/stores/modules/permission'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.public-library.list')
let { role_permission, role_type } = usePermissionStore()
const createOpenLibDoc = computed(() => role_type == 1 || role_permission.includes('CreateOpenLibDoc'))

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
    title: t('delete_title', { library_name: data.library_name }),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('delete_confirm'),
    okText: secondsToGo + ' ' + t('confirm'),
    okType: 'danger',
    okButtonProps: {
      disabled: true
    },
    cancelText: t('cancel'),
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
        okText: t('confirm'),
        okButtonProps: {
          disabled: false
        }
      })

      clearInterval(interval)
      interval = undefined
    } else {
      secondsToGo -= 1

      modal.update({
        okText: secondsToGo + ' ' + t('confirm'),
        okButtonProps: {
          disabled: true
        }
      })
    }
  }, 1000)
}

const onDelete = ({ id }) => {
  deleteLibrary({ id }).then(() => {
    message.success(t('delete_success'))
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
