<template>
  <a-drawer
    v-model:open="open"
    placement="right"
    width="400px"
    :headerStyle="{display: 'none'}"
    :bodyStyle="{padding: 0}"
  >
    <div v-if="checking" class="checking-box">
      <div class="cont">
        <a-spin/>
        {{ t('msg_checking') }}
      </div>
    </div>
    <div class="head-box">
      <div class="base-info-box">
        <img class="avatar" :src="detail.avatar"/>
        <div class="info-box">
          <div class="left">
            <div class="name zm-line1">{{ detail.name }}</div>
            <div>
              <div v-if="detail.has_auth == 1" class="auth-tag">{{ t('label_authorized') }}</div>
              <div v-else class="auth-tag fail">{{ t('label_unauthorized') }}</div>
            </div>
          </div>
          <div class="right">
            <a-dropdown>
              <EllipsisOutlined class="icon"/>
              <template #overlay>
                <a-menu>
                  <a-menu-item @click.stop="editApp">{{ t('btn_edit') }}</a-menu-item>
                  <a-menu-item @click.stop="delApp"><span class="cFB363F">{{ t('btn_delete') }}</span></a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
            <CloseOutlined @click="open = false" class="icon"/>
          </div>
        </div>
      </div>
      <div class="link">{{ detail.description }}</div>
      <div class="link">{{ detail.url }}</div>
      <a-button v-if="detail.has_auth == 1" class="auth-btn" @click="authConfig(2)">{{ t('btn_reauthorize') }}</a-button>
      <a-button v-else type="primary" class="auth-btn" @click="authConfig(1)">{{ t('btn_authorize') }}</a-button>
    </div>
    <div class="body-box">
      <div v-if="detail.has_auth != 1" class="no-auth-box">
        <img src="@/assets/no-permission.png"/>
        <div class="tit">{{ t('title_need_auth') }}</div>
        <div class="desc">{{ t('msg_need_auth_desc') }}</div>
        <a class="link" @click="authConfig(1)">{{ t('btn_authorize_now') }}</a>
      </div>
      <template v-else>
        <div class="main-tit">
          <span>{{ t('msg_available_tools', { count: detail.tools.length }) }}</span>
          <div class="extra-box">
            <div class="update-time">{{ t('msg_updated', { time: detail.up_time_text }) }}</div>
            <a class="refresh-btn" @click="refresh(true)">
              <SyncOutlined :spin="refreshing"/> {{ t('btn_refresh') }}
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
import {useI18n} from "@/hooks/web/useI18n";

const {t} = useI18n('views.robot.robot-list.components.third-mcp-detail');

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
        message.success(t('msg_verify_success'))
      }, 2000)
    }).catch(err => {
      detail.value.has_auth = 0
      message.warning(t('msg_auth_failed', { error: err.message }))
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
      title: t('title_confirm'),
      content: t('msg_reauth_confirm'),
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
      needAuth && message.success(t('msg_refresh_complete'))
    }).finally(() => {
      refreshing.value = false
    })
  }
  if (needAuth) {
    if (refreshing.value) return
    Modal.confirm({
      title: t('title_confirm'),
      content: t('msg_refresh_confirm'),
      onOk: () => {
        refreshing.value = true
        setShowReqError(false)
        authTMcpProvider({provider_id: detail.value.id}).then(res => {
          detail.value.has_auth = 1
          func()
        }).catch(err => {
          detail.value.has_auth = 0
          message.warning(t('msg_verify_auth_failed', { error: err.message }))
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
