<template>
  <div class="page-container">
    <div class="page-title">API Key管理</div>
    <a-alert class="mt8" type="info" show-icon>
      <template #icon>
        <div>
          <ExclamationCircleFilled />
        </div>
      </template>
      <template #message>
        <div class="alert-content">
          <div>
            请妥善保管您的API Key,注意不要外泄。API
            Key外泄可能导致其他人使用此Key访问您的机器人，消耗您模型资源。
          </div>
        </div>
      </template>
    </a-alert>
    <div class="opration-box">
      <a-button type="primary" @click="handleOpenAddModal">新增API Key</a-button>
      <a-input v-model:value="end_point" style="width: 500px" readonly>
        <template #addonBefore>
          <span @click="copyText(end_point)" style="cursor: pointer">复制API Endpoint</span>
        </template>
      </a-input>

      <a
        style="margin-left: auto"
        href="https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/bg1ol40fo68pgdae"
        target="_blank"
      >
        <a-button>API 文档</a-button>
      </a>
    </div>
    <div class="list-box">
      <div class="list-item" v-for="item in lists" :key="item.id">
        <a-flex justify="space-between" align="center">
          <div class="key-content">
            <span>API Key：</span>
            <span>{{ item.key }}</span>
          </div>
          <div class="right-box">
            <a-switch
              @change="handleChangeStatus(item)"
              :checked="item.status == 1"
              checked-children="开"
              un-checked-children="关"
            />

            <a-button @click="copyText(item.key)">复制</a-button>
            <a-button danger @click="handleDelApiKey(item)">删除</a-button>
          </div>
        </a-flex>
      </div>
    </div>
  </div>
</template>

<script setup>
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
    message.success('更新成功')
    getApiKeyList()
  })
}

const handleOpenAddModal = () => {
  Modal.confirm({
    title: h('div', { style: 'color:#262626; font-weight: 600;' }, '新增API Key'),
    icon: null,
    content: h(
      'div',
      { style: ' color: #f10; line-height: 22px;' },
      '请妥善保管您的API Key,注意不要外泄。API Key外泄可能导致其他人使用此Key访问您的机器人，消耗您模型资源。'
    ),
    bodyStyle: {
      padding: '12px 8px 4px 12px'
    },
    okText: '创建',
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
    message.success('创建成功')
    getApiKeyList()
  })
}

const handleDelApiKey = (item) => {
  let secondsToGo = 5
  let key = formatString(item.key)
  let modal = Modal.confirm({
    title: `确定删除API Key：${key}`,
    icon: createVNode(ExclamationCircleOutlined),
    content: h('div', { style: { color: '#FB363F' } }, '删除后，将无法使用该Key调用接口'),
    okText: secondsToGo + ' 确 定',
    okType: 'danger',
    cancelText: '取 消',
    okButtonProps: {
      disabled: true
    },
    onOk() {
      deleteRobotApikey({
        robot_key: robotInfo.value.robot_key,
        id: item.id
      }).then((res) => {
        message.success('删除成功')
        getApiKeyList()
      })
    },
    onCancel() {}
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
function formatString(str) {
  let start = str.slice(0, 4) // 取前四位
  let end = str.slice(-4) // 取后四位，注意这里使用了负数索引

  // 检查字符串长度，如果长度小于8，则直接返回原字符串或做适当调整
  if (str.length <= 4) {
    return str // 如果原字符串长度小于等于4，直接返回原字符串
  } else if (str.length <= 8) {
    // 如果原字符串长度在5到8之间，只取前四位并添加"..."
    return start + '...'
  } else {
    // 否则，返回前四位，中间加"..."，再加后四位
    return start + '...' + end
  }
}
const copyText = async (text) => {
  try {
    await toClipboard(text)
    message.success('复制成功')
  } catch (e) {
    message.error('复制失败')
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
