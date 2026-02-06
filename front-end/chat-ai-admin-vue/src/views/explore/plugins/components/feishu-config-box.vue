<template>
  <div>
    <a-drawer
      v-model:open="open"
      placement="right"
      width="400px"
      :headerStyle="{display: 'none'}"
      :bodyStyle="{padding: 0}"
    >
      <div class="head-box">
        <div class="base-info-box">
          <img class="avatar" :src="detail.icon"/>
          <div class="info-box">
            <div class="left">
              <div class="name zm-line1">{{ detail.title }}</div>
              <div>{{ detail.author }}</div>
            </div>
            <div class="right">
              <CloseOutlined @click="open = false" class="icon"/>
            </div>
          </div>
        </div>
        <div class="link">{{ detail.description }}</div>
        <a-popover v-if="configLen > 0" placement="bottom">
          <template #content>
            <div class="tools-box">
              <div v-for="(item, key) in configData" :key="key" class="tool-item">
                <img class="cover" src="@/assets/img/robot/tool-icon.svg"/>
                <div class="info">
                  <div class="tit zm-line1">{{ item.name }}</div>
                  <div class="btns">
                    <a-tag v-if="item.is_default" color="#f50">{{ t('label_default') }}</a-tag>
                    <a v-else @click="setConfigDef(item, key)">{{ t('btn_set_default') }}</a>
                    <a @click="delConfig(item, key)">{{ t('btn_delete') }}</a>
                  </div>
                </div>
              </div>
            </div>
          </template>
          <a-button type="primary" :icon="h(PlusOutlined)" class="auth-btn" @click="showAddModal">
            {{ t('msg_authorized_credentials', {count: configLen}) }}
          </a-button>
        </a-popover>
        <a-button v-else type="primary" class="auth-btn" @click="showAddModal">{{ t('btn_authorize_now') }}</a-button>
      </div>
      <div class="body-box">
        <div v-if="!configLen" class="no-auth-box">
          <img src="@/assets/no-permission.png"/>
          <div class="tit">{{ t('msg_need_authorization') }}</div>
          <div class="desc">{{ t('msg_config_after_auth') }}</div>
          <a class="link" @click="showAddModal">{{ t('btn_authorize_now') }}</a>
        </div>
        <div v-else class="action-box">
          <div v-for="(item, i) in actionData" :key="i" class="action-item">
            <div class="tit">{{ item.title }}</div>
            <div class="desc">{{ item.desc }}</div>
            <div class="params">
              <a-tag
                v-for="(field, j) in Object.keys(item.params).slice(0, 3)"
                :key="j"
                :bordered="false">{{ field }}
              </a-tag>
              <a-popover placement="right">
                <template #content>
                  <div class="params-box">
                    <div v-for="(field, key) in item.params" :key="key" class="param-item">
                      <div class="field">
                        <span class="name">{{ key }}</span>
                        <span class="type">{{ field.type }}</span>
                        <span v-if="field.required" class="required">{{ t('label_required') }}</span>
                      </div>
                      <div class="desc">{{ field.desc }}</div>
                    </div>
                  </div>
                </template>
                <a>{{ t('label_params') }}</a>
              </a-popover>
            </div>
          </div>
        </div>
      </div>
    </a-drawer>

    <FeishuConfigModal ref="configModalRef" @change="loadConfig"/>
  </div>
</template>

<script setup>
import {ref, reactive, h, computed} from 'vue';
import {EllipsisOutlined, CloseOutlined, PlusOutlined} from '@ant-design/icons-vue';
import {message, Modal} from 'ant-design-vue';
import {authTMcpProvider, getTMcpProviderInfo} from "@/api/robot/thirdMcp";
import {jsonDecode, timeNowGapFormat} from "@/utils/index.js";
import {setShowReqError} from "@/utils/http/axios/config.js";
import {getPluginConfig, runPlugin, setPluginConfig} from "@/api/plugins/index.js";
import FeishuConfigModal from "@/views/explore/plugins/components/feishu-config-modal.vue";
import { useI18n } from '@/hooks/web/useI18n';

const { t } = useI18n('views.explore.plugins.components.feishu-config-box')
const emit = defineEmits(['del', 'edit', 'auth'])

const configModalRef = ref(null)
const open = ref(false)
const detail = ref({})
const configData = ref({})
const actionData = ref({})

const configLen = computed(() => {
  return Object.keys(configData.value).length
})

function show(info) {
  detail.value = info
  loadConfig()
  loadAction()
  open.value = true
}

function hide() {
  open.value = false
}

function loadConfig() {
  getPluginConfig({name: detail.value.name}).then(res => {
    configData.value = jsonDecode(res?.data, {})
  })
}

function loadAction() {
  runPlugin({
    name: detail.value.name,
    action: "default/get-schema",
    params: {}
  }).then(res => {
    let _data = res?.data || {}
    for (let key in _data) {
      if (_data[key].type != 'node') delete _data[key]
    }
    actionData.value = _data
  })
}

function showAddModal() {
  configModalRef.value.show(configData.value)
}

function setConfigDef(item, key) {
  for (let _key in configData.value) {
    configData.value[_key].is_default = (_key == key)
  }
  setPluginConfig({
    name: detail.value.name,
    data: JSON.stringify(configData.value)
  }).then(res => {
    loadConfig()
  })
}

function delConfig(item, key) {
  Modal.confirm({
    title: t('title_tip'),
    content: t('msg_confirm_delete_config'),
    onOk: () => {
      delete configData.value[key]
      setPluginConfig({
        name: detail.value.name,
        data: JSON.stringify(configData.value)
      }).then(res => {
        loadConfig()
        message.success(t('msg_deleted'))
      })
    }
  })
}

defineExpose({
  show,
  hide,
})
</script>

<style scoped lang="less">
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
}

.tools-box {
  width: 330px;

  .tool-item {
    display: flex;
    justify-content: space-between;
    gap: 12px;
    margin-top: 16px;

    &:first-child {
      margin-top: 0;
    }

    .cover {
      width: 32px;
      height: 32px;
      border-radius: 6px;
      flex-shrink: 0;
    }

    .info {
      flex: 1;
      display: flex;
      justify-content: space-between;
      align-items: center;

      .tit {
        max-width: 200px;
        color: #000000;
        font-weight: 400;
      }

      .btns {
        display: flex;
        align-items: center;
        gap: 8px;
      }
    }
  }
}

.action-box {
  display: flex;
  flex-direction: column;
  gap: 16px;

  .action-item {
    border-bottom: 1px solid #D9D9D9;
    padding-bottom: 16px;
    .tit {
      color: #262626;
      font-size: 14px;
      font-weight: 600;
    }

    .desc {
      color: #595959;
      font-size: 14px;
      font-weight: 400;
      margin-top: 8px;
    }

    .params {
      margin-top: 8px;
    }
  }
}

.params-box {
  max-height: 80vh;
  overflow-y: auto;

  .param-item {
    max-width: 500px;

    &:not(:last-child) {
      margin-bottom: 16px;
    }

    .field {
      color: #262626;
      font-size: 14px;
      display: flex;
      align-items: center;
      gap: 12px;

      .name {
        font-weight: 600;
      }

      .type {
        color: #595959;
      }

      .required {
        color: #ED744A;
        font-weight: 500;
      }
    }

    .desc {
      color: #8c8c8c;
      font-size: 14px;
      margin-top: 4px;
    }
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
</style>
