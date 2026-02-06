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
        <slot name="head-extra"></slot>
      </div>
      <div class="body-box">
        <slot name="body">
          <div class="action-box">
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
        </slot>
      </div>
    </a-drawer>
    <PluginConfigBox ref="pluginConfigRef" @auth="emit('auth')" />
  </div>
</template>

<script setup>
import {ref, reactive, h, computed} from 'vue';
import {EllipsisOutlined, CloseOutlined, PlusOutlined} from '@ant-design/icons-vue';
import {message, Modal} from 'ant-design-vue';
import {jsonDecode, timeNowGapFormat} from "@/utils/index.js";
import {getPluginConfig, runPlugin, setPluginConfig} from "@/api/plugins/index.js";
import PluginConfigBox from "./plugin-config-box.vue";
import { useI18n } from '@/hooks/web/useI18n';

const { t } = useI18n('views.explore.plugins.components.config-box')
const emit = defineEmits(['del', 'edit', 'auth'])

const open = ref(false)
const detail = ref({})
const actionData = ref({})
const pluginConfigRef = ref(null)

function show(info) {
  detail.value = info
  loadAction()
}

function hide() {
  open.value = false
}

function loadAction() {
  runPlugin({
    name: detail.value.name,
    action: "default/get-schema",
    params: {}
  }).then(res => {
    let _data = res?.data || {}
    const usePluginConfig = Object.values(_data || {}).some(v => v && v.use_plugin_config === true)
    if (usePluginConfig) {
      pluginConfigRef.value && pluginConfigRef.value.show(detail.value, _data)
      open.value = false
      return
    }
    for (let key in _data) {
      if (_data[key].type != 'node') delete _data[key]
    }
    actionData.value = _data
    open.value = true
  })
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

.cFB363F {
  color: #FB363F;
}

.tip-desc {
  color: #8C8C8C;
  margin-top: 4px;
}

.mt16 {
  margin-top: 16px;
}
</style>
