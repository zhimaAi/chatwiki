<template>
  <div class="page-container">
    <div class="page-title">{{ t('page_title') }}</div>
    <a-alert class="mt8" type="info" show-icon>
      <template #icon>
        <div>
          <ExclamationCircleFilled />
        </div>
      </template>
      <template #message>
        <div class="alert-content">
          <div>
            {{ t('alert_message') }}
          </div>
        </div>
      </template>
    </a-alert>
    <div class="opration-box">
      <a-button type="primary" @click="handleOpenAddModal">{{ t('add_api_key_btn') }}</a-button>
      <a-input v-model:value="end_point" style="width: 500px" readonly>
        <template #addonBefore>
          <span @click="copyText(end_point)" style="cursor: pointer">{{ t('copy_api_endpoint') }}</span>
        </template>
      </a-input>

      <a
        style="margin-left: auto"
        href="https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/bg1ol40fo68pgdae"
        target="_blank"
      >
        <a-button>{{ t('api_documentation') }}</a-button>
      </a>
    </div>
    <div class="list-box">
      <div class="list-item" v-for="item in lists" :key="item.id">
        <a-flex justify="space-between" align="center">
          <div class="key-content">
            <span>{{ t('api_key_label') }}</span>
            <span>{{ item.key }}</span>
          </div>
          <div class="right-box">
            <a-switch
              @change="handleChangeStatus(item)"
              :checked="item.status == 1"
              :checked-children="t('switch_on')"
              :un-checked-children="t('switch_off')"
            />

            <a-button @click="copyText(item.key)">{{ t('copy_btn') }}</a-button>
            <a-button danger @click="handleDelApiKey(item)">{{ t('delete_btn') }}</a-button>
          </div>
        </a-flex>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'
import {
  ExclamationCircleFilled,
  CopyOutlined,
  ExclamationCircleOutlined
} from '@ant-design/icons-vue'
import { message, Modal } from 'ant-design-vue'
import { ref, h, createVNode } from 'vue'
import { storeToRefs } from 'pinia'
import { useRobotStore } from '@/stores/modules/robot'
import {
  listRobotApikey,
  updateRobotApikey,
  addRobotApikey,
  deleteRobotApikey
} from '@/api/robot/index'
import useClipboard from 'vue-clipboard3'
const { t } = useI18n('views.robot.api-key-manage.index')
const { toClipboard } = useClipboard()
const robotStore = useRobotStore()
const { robotInfo } = storeToRefs(robotStore)
const end_point = ref('')
const lists = ref([])
const getApiKeyList = () => {
  listRobotApikey({
    robot_key: robotInfo.value.robot_key
  }).then((res) => {
    end_point.value = res.data.end_point
    lists.value = res.data.list
  })
}
getApiKeyList()
const handleChangeStatus = (item) => {
  updateRobotApikey({
    robot_key: robotInfo.value.robot_key,
    id: item.id
  }).then((res) => {
    message.success(t('update_success'))
    getApiKeyList()
  })
}

const handleOpenAddModal = () => {
  Modal.confirm({
    title: h('div', { style: 'color:#262626; font-weight: 600;' }, t('add_api_key_title')),
    icon: null,
    content: h(
      'div',
      { style: ' color: #f10; line-height: 22px;' },
      t('add_api_key_warning')
    ),
    bodyStyle: {
      padding: '12px 8px 4px 12px'
    },
    okText: t('create_btn'),
    width: 480,
    onOk: () => {
      handleCrrate()
    },
    onCancel: () => {}
  })
}

const handleCrrate = () => {
  addRobotApikey({
    robot_key: robotInfo.value.robot_key
  }).then((res) => {
    message.success(t('create_success'))
    getApiKeyList()
  })
}

const handleDelApiKey = (item) => {
  let secondsToGo = 5
  let key = formatString(item.key)
  let modal = Modal.confirm({
    title: t('confirm_delete_title', { key }),
    icon: createVNode(ExclamationCircleOutlined),
    content: h('div', { style: { color: '#FB363F' } }, t('delete_warning')),
    okText: secondsToGo + ' ' + t('confirm_btn'),
    okType: 'danger',
    cancelText: t('cancel_btn'),
    okButtonProps: {
      disabled: true
    },
    onOk() {
      deleteRobotApikey({
        robot_key: robotInfo.value.robot_key,
        id: item.id
      }).then((res) => {
        message.success(t('delete_success'))
        getApiKeyList()
      })
    },
    onCancel() {}
  })

  let interval = setInterval(() => {
    if (secondsToGo == 1) {
      modal.update({
        okText: t('confirm_btn'),
        okButtonProps: {
          disabled: false
        }
      })

      clearInterval(interval)
      interval = undefined
    } else {
      secondsToGo -= 1

      modal.update({
        okText: secondsToGo + ' ' + t('confirm_btn'),
        okButtonProps: {
          disabled: true
        }
      })
    }
  }, 1000)
}
function formatString(str) {
  let start = str.slice(0, 4)
  let end = str.slice(-4)

  if (str.length <= 4) {
    return str
  } else if (str.length <= 8) {
    return start + '...'
  } else {
    return start + '...' + end
  }
}
const copyText = async (text) => {
  try {
    await toClipboard(text)
    message.success(t('copy_success'))
  } catch (e) {
    message.error(t('copy_fail'))
  }
}
</script>

<style lang="less" scoped>
.page-container {
  padding: 16px 24px;
  .page-title {
    line-height: 24px;
    font-size: 16px;
    font-weight: 600;
  }
  &::v-deep(.ant-alert-info) {
    margin-top: 16px;
    padding: 9px 16px;
    border-radius: 6px;
    background: #e9f1fe;
    border: 1px solid #99bffd;
    font-size: 14px;
    line-height: 22px;
    align-items: baseline;
    & > .anticon {
      color: #2475fc;
    }
    .alert-content {
      color: #3a4559;
      font-weight: 400;
    }
  }

  .opration-box {
    margin-top: 16px;
    display: flex;
    gap: 16px;
  }

  .list-box {
    margin-top: 16px;
    .list-item {
      padding: 14px 16px;
      border: 1px solid #d9d9d9;
      margin-bottom: 16px;
      border-radius: 4px;
      .key-content {
        flex: 1;
        word-break: break-all;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        span {
          color: #262626;
          font-weight: 500;
          line-height: 22px;
          font-size: 14px;
        }
      }
      .right-box {
        margin-left: 78px;
        display: flex;
        align-items: center;
        gap: 12px;
      }
    }
  }
}
</style>
