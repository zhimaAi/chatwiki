<template>
  <div class="library-page">
    <PageTabs class="mb-16" :tabs="pageTabs" active="/database/list"></PageTabs>
    <page-alert style="margin-bottom: 16px" :title="t('usage_guide.title')">
      <div>
        <p>{{ t('usage_guide.point_1') }}</p>
        <p>{{ t('usage_guide.point_2') }}</p>
      </div>
    </page-alert>

    <div class="library-page-body">
      <div class="list-toolbar">
        <div class="toolbar-box">
          <h3 class="list-total">{{ t('list.all') }} ({{ list.length }})</h3>
        </div>

        <div class="toolbar-box">
          <a-button type="primary" @click="toAdd()" v-if="formCreate">
            <template #icon>
              <PlusOutlined />
            </template>
            {{ t('list.create_btn') }}
          </a-button>
        </div>
      </div>

      <div class="list-box">
        <div class="list-item-wrapper" v-for="item in list" :key="item.id">
          <div class="list-item" @click.stop="toEdit(item)">
            <div class="library-info">
              <img class="library-icon" src="@/assets/img/database/base-icon.svg" alt="" />
              <div class="library-info-content">
                <div class="library-title">{{ item.name }}</div>
              </div>
            </div>
            <div class="item-body">
              <div class="library-desc">{{ item.description }}</div>
            </div>

            <div class="item-footer">
              <div class="library-size">
                <span class="text-item">{{ t('item.data_count') }}{{ item.entry_count }}{{ t('item.count_unit') }}</span>
                <span class="text-item">{{ t('item.related_apps') }}{{ item.robot_nums || 0 }}</span>
              </div>

              <div class="action-box" @click.stop>
                <a-dropdown>
                  <div class="action-item" @click.stop>
                    <svg-icon class="action-icon" name="point-h"></svg-icon>
                  </div>
                  <template #overlay>
                    <a-menu>
                      <a-menu-item>
                        <a href="javascript:;" @click.stop="toAdd(item)">{{ t('item.edit') }}</a>
                      </a-menu-item>
                      <a-menu-item>
                        <a class="delete-text-color" href="javascript:;" @click="handleDelete(item)"
                          >{{ t('item.delete') }}</a
                        >
                      </a-menu-item>
                    </a-menu>
                  </template>
                </a-dropdown>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <AddDataSheet @ok="getList" ref="addDataSheetRef"></AddDataSheet>
</template>

<script setup>
import { usePermissionStore } from '@/stores/modules/permission'
import { ref, createVNode, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Modal, message } from 'ant-design-vue'
import { PlusOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { getFormList, delForm } from '@/api/database'
import AddDataSheet from './components/add-data-sheet.vue'
import PageTabs from '@/components/cu-tabs/page-tabs.vue'
import PageAlert from '@/components/page-alert/page-alert.vue'
import { getDatabasePermission } from '@/utils/permission'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.database.database-list.index')
const router = useRouter()

const pageTabs = computed(() => [
  {
    title: t('page_tabs.library'),
    path: '/library/list'
  },
  {
    title: t('page_tabs.database'),
    path: '/database/list'
  },
  {
    title: t('page_tabs.document_extract_faq'),
    path: '/ai-extract-faq/list'
  },
  {
    title: t('page_tabs.trigger_statistics'),
    path: '/trigger-statics/list'
  },
])

const permissionStore = usePermissionStore()
let { role_permission, role_type } = permissionStore
const formCreate = computed(() => role_type == 1 || role_permission.includes('FormCreate'))

const list = ref([])

const getList = () => {
  getFormList({}).then((res) => {
    list.value = res.data || []
  })
}

getList()

const addDataSheetRef = ref(null)
const toAdd = (data = {}) => {
  console.log(data)
  if (data.id) {
    let key = getDatabasePermission(data.id)
    if (!(key == 4 || key == 2)) {
      return message.error(t('message.no_edit_permission'))
    }
  }
  addDataSheetRef.value.show(data)
}

const toEdit = (data) => {
  // router.push({
  //   path: '/database/details',
  //   query: {
  //     form_id: data.id,
  //     name: data.name
  //   }
  // })
  window.open(`/#/database/details?form_id=${data.id}&name=${data.name}`, "_blank", "noopener") // 建议添加 noopener 防止安全漏洞
}

const handleDelete = (data) => {
  let key = getDatabasePermission(data.id)
  if (key != 4) {
    return message.error(t('message.no_delete_permission'))
  }
  let secondsToGo = 3

  let modal = Modal.confirm({
    title: t('delete_confirm.title'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('delete_confirm.content', { name: data.name }),
    okText: secondsToGo + ' ' + t('delete_confirm.confirm_btn'),
    okType: 'danger',
    okButtonProps: {
      disabled: true
    },
    cancelText: t('delete_confirm.cancel_btn'),
    onOk() {
      onDelete(data)
    },
    onCancel() {}
  })

  let interval = setInterval(() => {
    if (secondsToGo == 1) {
      modal.update({
        okText: t('delete_confirm.confirm_btn'),
        okButtonProps: {
          disabled: false
        }
      })

      clearInterval(interval)
      interval = undefined
    } else {
      secondsToGo -= 1

      modal.update({
        okText: secondsToGo + ' ' + t('delete_confirm.confirm_btn'),
        okButtonProps: {
          disabled: true
        }
      })
    }
  }, 1000)
}

const onDelete = ({ id }) => {
  delForm({ id }).then(() => {
    message.success(t('message.delete_success'))
    getList()
  })
}
</script>

<style lang="less" scoped>
.library-page {
  .list-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
  }

  .list-total {
    line-height: 24px;
    font-size: 16px;
    font-weight: 600;
    color: #262626;
  }
  .list-box {
    display: flex;
    flex-flow: row wrap;
    margin: 0 -8px;
  }
  .list-item-wrapper {
    padding: 8px;
    width: 25%;
  }
  .list-item {
    position: relative;
    width: 100%;
    padding: 24px;
    border: 1px solid #e4e6eb;
    border-radius: 12px;
    background-color: #fff;
    transition: all 0.25s;
    cursor: pointer;

    &:hover {
      box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.12);
    }

    .library-info {
      position: relative;
      display: flex;
      align-items: center;
      .item-action {
        .menu-btn {
          position: absolute;
          right: 0;
          top: 0;
          width: 22px;
          height: 22px;
          text-align: center;
          line-height: 22px;
          font-size: 16px;
          cursor: pointer;
          &:hover {
            color: #2475fc;
          }
        }
      }
      .library-icon {
        width: 52px;
        height: 52px;
        border-radius: 14px;
        overflow: hidden;
      }
      .library-info-content {
        margin-left: 12px;
        flex: 1;
        overflow: hidden;
      }
      .library-title {
        height: 24px;
        line-height: 24px;
        font-size: 16px;
        font-weight: 600;
        color: #262626;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
    }
    .item-body {
      margin-top: 12px;
    }
    .library-desc {
      height: 44px;
      line-height: 22px;
      font-size: 14px;
      font-weight: 400;
      color: rgb(89, 89, 89);
      // 超出2行显示省略号
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      line-clamp: 2;
      -webkit-box-orient: vertical;
    }

    .item-footer {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-top: 14px;
      color: #7a8699;
    }
    .library-size {
      display: flex;
      line-height: 20px;
      font-size: 12px;
      font-weight: 400;
      color: #7a8699;

      .text-item {
        margin-right: 12px;
        &:last-child {
          margin-right: 0;
        }
      }
    }

    .action-box {
      font-size: 14px;
      height: 24px;
      color: #2475fc;
      display: flex;
      align-items: center;

      .action-item {
        display: flex;
        align-items: center;
        height: 100%;
        padding: 4px;
        border-radius: 6px;
        cursor: pointer;
        color: #595959;
        transition: all 0.2s;
      }
      .action-item:hover {
        background: #e4e6eb;
      }

      .action-icon {
        font-size: 16px;
      }
    }
  }
}
// 大于1920px
@media screen and (min-width: 1920px) {
  .library-page {
    .list-box {
      .list-item-wrapper {
        width: 20%;
      }
    }
  }
}
</style>
