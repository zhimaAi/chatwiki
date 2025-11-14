<template>
  <a-drawer
    v-model:open="open"
    placement="right"
    width="400px"
    headerStyle="display:none;"
    bodyStyle="padding:0;"
  >
    <div v-if="checking" class="checking-box">
      <div class="cont">
        <a-spin/>
        正在验证...
      </div>
    </div>
    <div class="head-box">
      <div class="base-info-box">
        <img class="avatar" :src="detail.avatar"/>
        <div class="info-box">
          <div class="left">
            <div class="name zm-line1">{{ detail.name }}</div>
            <div>
              <div v-if="detail.has_auth == 1" class="auth-tag">已授权</div>
              <div v-else class="auth-tag fail">未授权</div>
            </div>
          </div>
          <div class="right">
            <a-dropdown>
              <EllipsisOutlined class="icon"/>
              <template #overlay>
                <a-menu>
                  <a-menu-item @click.stop="editApp">编辑</a-menu-item>
                  <a-menu-item @click.stop="delApp"><span class="cFB363F">删除</span></a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
            <CloseOutlined @click="open = false" class="icon"/>
          </div>
        </div>
      </div>
      <div class="link">{{ detail.description }}</div>
      <div class="link">{{ detail.url }}</div>
      <a-button v-if="detail.has_auth == 1" class="auth-btn" @click="authConfig(2)">重新授权</a-button>
      <a-button v-else type="primary" class="auth-btn" @click="authConfig(1)">授 权</a-button>
    </div>
    <div class="body-box">
      <div v-if="detail.has_auth != 1" class="no-auth-box">
        <img src="@/assets/no-permission.png"/>
        <div class="tit">需要授权</div>
        <div class="desc">授权后，MCP工具将显示在这里</div>
        <a class="link" @click="authConfig(1)">立即授权</a>
      </div>
      <template v-else>
        <div class="main-tit">
          <span>{{ detail.tools.length }}个可用工具</span>
          <div class="extra-box">
            <div class="update-time">{{ detail.up_time_text }}更新</div>
            <a class="refresh-btn" @click="refresh(true)">
              <SyncOutlined :spin="refreshing"/> 刷新
            </a>
          </div>
        </div>
        <div class="tools-box">
          <div v-for="(tool, i) in detail.tools" :key="i" class="tool-item">
            <img class="cover" src="@/assets/img/robot/tool-icon.svg"/>
            <div class="info">
              <div class="tit">{{ tool.name }}</div>
              <div class="desc">{{ tool.description}}</div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </a-drawer>
</template>

<script setup>
import {ref} from 'vue';
import {EllipsisOutlined, CloseOutlined, SyncOutlined} from '@ant-design/icons-vue';
import {message, Modal} from 'ant-design-vue';
import {authTMcpProvider, getTMcpProviderInfo} from "@/api/robot/thirdMcp";
import {jsonDecode, timeNowGapFormat} from "@/utils/index.js";
import {setShowReqError} from "@/utils/http/axios/config.js";

const emit = defineEmits(['del', 'edit', 'auth'])

const open = ref(false)
const checking = ref(false)
const refreshing = ref(false)
const detail = ref({})

function show(info) {
  detail.value = info
  open.value = true
}

function hide() {
  open.value = false
}

function authConfig(type) {
  const auth = () => {
    checking.value = true
    setShowReqError(false)
    authTMcpProvider({provider_id: detail.value.id}).then(res => {
      setTimeout(() => {
        checking.value = false
        detail.value.has_auth = 1
        refresh()
        message.success('已验证，授权完成')
      }, 2000)
    }).catch(err => {
      detail.value.has_auth = 0
      message.warning('授权失败：'+err.message)
      checking.value = false
    }).finally(() => {
      emit('auth', detail.value.has_auth)
      setTimeout(() => {
        setShowReqError(true)
      }, 500)
    })
  }
  if (type == 2) {
    Modal.confirm({
      title: '提示',
      content: '授权或者刷新应用列表都可能会影响已有工具的调用，确认继续操作吗？',
      onOk: () => {
        auth()
      }
    })
  } else {
    auth()
  }
}

function editApp() {
  emit('edit', detail.value)
}

function delApp() {
  emit('del', detail.value)
}

function refresh(needAuth=false) {
  const func = () => {
    getTMcpProviderInfo({provider_id: detail.value.id}).then(res => {
      let info = res?.data || {}
      info.tools = jsonDecode(info.tools, [])
      info.up_time_text = timeNowGapFormat(info.update_time)
      detail.value = info
      needAuth && message.success('刷新完成')
    }).finally(() => {
      refreshing.value = false
    })
  }
  if (needAuth) {
    if (refreshing.value) return
    Modal.confirm({
      title: '提示',
      content: '刷新应用列表或者授权都可能会影响已有工具的调用，确认继续操作吗？',
      onOk: () => {
        refreshing.value = true
        setShowReqError(false)
        authTMcpProvider({provider_id: detail.value.id}).then(res => {
          detail.value.has_auth = 1
          func()
        }).catch(err => {
          detail.value.has_auth = 0
          message.warning('验证授权失败：'+err.message)
          refreshing.value = false
        }).finally(() => {
          emit('auth', detail.value.has_auth)
          setShowReqError(true)
        })
      }
    })
  } else {
    func()
  }
}

defineExpose({
  show,
  hide,
  refresh,
})
</script>

<style scoped lang="less">
.checking-box {
  width: 100%;
  height: 100%;
  position: absolute;
  top: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 99;

  .cont {
    padding: 40px 80px;
    border-radius: 6px;
    background: #FFF;
    box-shadow: 0 4px 16px 0 #1b3a6929;
    display: flex;
    gap: 8px;
    font-size: 16px;
    font-weight: 500;
  }
}

.head-box {
  display: flex;
  flex-direction: column;
  gap: 12px;
  color: #595959;
  font-size: 14px;
  font-weight: 400;
  border-bottom: 1px solid #E4E6EB;
  padding: 24px;

  .base-info-box {
    display: flex;
    align-items: center;
    justify-content: space-between;

    .avatar {
      width: 62px;
      height: 62px;
      border-radius: 14.59px;
      flex-shrink: 0;
      margin-right: 12px;
    }

    .info-box {
      flex: 1;
      display: flex;
      justify-content: space-between;
      align-items: flex-start;

      .left {
        max-width: 70%;
      }

      .right {
        display: flex;
        align-items: center;
        gap: 12px;
      }

      .icon {
        cursor: pointer;
      }

      .name {
        color: #262626;
        font-size: 16px;
        font-weight: 600;
        margin-bottom: 4px;
      }
    }
  }

  .link {
    word-break: break-all;
  }

  .auth-btn {
    width: 100%;
  }
}

.body-box {
  padding: 24px;

  .main-tit {
    display: flex;
    align-items: center;
    justify-content: space-between;
    color: #262626;
    font-size: 14px;
    font-weight: 600;

    .extra-box {
      display: flex;
      align-items: center;
      gap: 12px;
    }

    .update-time {
      color: #8c8c8c;
      font-size: 12px;
      font-weight: 400;
    }

    .refresh-btn {
      font-weight: 400;
    }
  }

  .tools-box {
    .tool-item {
      display: flex;
      justify-content: space-between;
      gap: 12px;
      margin-top: 16px;

      .cover {
        width: 32px;
        height: 32px;
        border-radius: 6px;
        flex-shrink: 0;
      }

      .info {
        flex: 1;

        .tit {
          color: #000000;
          font-weight: 400;
        }

        .desc {
          color: #8c8c8c;
          font-size: 12px;
          font-weight: 400;
          margin-top: 4px;
        }
      }
    }
  }
}

.auth-tag {
  display: inline-block;
  padding: 1px 8px;
  border-radius: 6px;
  background: #C4F5DB;
  color: #17814e;
  font-size: 12px;
  font-weight: 400;

  &.fail {
    color: #ED744A;
    background: #FFECE6;
  }
}

.no-auth-box {
  display: flex;
  flex-direction: column;
  align-items: center;

  img {
    width: 200px;
    height: 200px
  }

  .tit {
    color: #262626;
    font-size: 16px;
    font-weight: 600;
  }

  .desc {
    color: #8c8c8c;
    font-size: 14px;
    font-weight: 400;
    margin: 8px 0 12px;
  }
}

.cFB363F {
  color: #FB363F;
}
</style>
