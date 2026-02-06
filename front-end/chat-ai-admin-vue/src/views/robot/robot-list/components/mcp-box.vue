<template>
  <div class="_mcp-box">
    <div v-if="loading" class="loading-box">
      <a-spin/>
    </div>
    <template v-else-if="appInfo.id > 0">
      <div class="base-info-box">
        <div class="app-info">
          <img class="avatar" :src="appInfo.avatar"/>
          <div class="info-box">
            <div class="left-box">
              <div class="name">
                {{ appInfo.name }}
                <span v-if="appInfo.status_bool" class="status-tag finished"><CheckCircleFilled/> {{ t('status_published') }}</span>
                <span v-else class="status-tag waiting"><ExclamationCircleFilled/> {{ t('status_unpublished') }}</span>
              </div>
              <div class="desc">{{ appInfo.description }}</div>
            </div>
            <div class="right-box">
              <a-switch
                  v-model:checked="appInfo.status_bool"
                  @change="statusChange"
                  :loading="statusLoading"
                  :checked-children="t('switch_on')"
                  :un-checked-children="t('switch_off')"/>
              <a-button @click="showWorkflowModal" type="primary" ghost :icon="h(PlusOutlined)">{{ t('btn_add_workflow') }}</a-button>
              <a-dropdown>
                <a-button :icon="h(EllipsisOutlined)"/>
                <template #overlay>
                  <a-menu>
                    <a-menu-item @click="showMcpModal">{{ t('btn_edit') }}</a-menu-item>
                    <a-menu-item @click="delMcp"><span class="cFB363F">{{ t('btn_delete') }}</span></a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </div>
          </div>
        </div>
        <div class="app-data">
          <template v-for="(item, i) in appData" :key="i">
            <div class="data-item">
              <div class="title">{{ t(item.title) }}</div>
              <div class="value">{{ item.value }}</div>
            </div>
            <a-divider type="vertical" style="height: 24px;"/>
          </template>
        </div>
      </div>
      <div class="tools-box">
        <div class="head-box">
          <AppstoreOutlined class="icon"/>
          {{ t('label_tools_list') }}（{{ tools.length || 0 }}）
        </div>
        <div v-if="tools.length" class="tools-list">
          <div v-for="(item, i) in tools" :key="i" class="tools-item">
            <div class="left">
              <div class="tit">{{ item.robot_name }}</div>
              <div class="key">{{ t('label_unique_identifier') }}：{{ item.name }}</div>
              <div class="desc">{{ item.robot_intro }}</div>
              <div class="params">
                <a-tag
                    v-for="(field, j) in item.params.slice(0, 5)"
                    :key="j"
                    :bordered="false">{{ field.key }}
                </a-tag>
                <a-popover placement="right">
                  <template #content>
                    <div class="params-box">
                      <div v-for="(field, j) in item.params" :key="j" class="param-item">
                        <div class="field">
                          <span class="name">{{ field.key }}</span>
                          <span class="type">{{ field.typ }}</span>
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
            <a-dropdown>
              <a-button size="small" :icon="h(EllipsisOutlined)"/>
              <template #overlay>
                <a-menu>
                  <a-menu-item @click="showIdentModal(item)">{{ t('btn_edit_identifier') }}</a-menu-item>
                  <a-menu-item @click="delTool(item)"><span class="cFB363F">{{ t('btn_delete') }}</span></a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
        </div>
        <div v-else class="empty-box">
          <img src="@/assets/empty.png"/>
          <div class="title">{{ t('msg_no_workflow_bound') }}</div>
          <a-button @click="showWorkflowModal" class="btn" type="primary">{{ t('btn_bind_now') }}</a-button>
        </div>
      </div>
    </template>
    <div v-else class="empty-box">
      <img src="@/assets/empty.png"/>
      <div class="title">{{ t('msg_no_mcp_added') }}</div>
      <a-button @click="showMcpModal" class="btn" type="primary">{{ t('btn_add_now') }}</a-button>
    </div>

    <McpStore ref="mcpRef" @ok="loadData"/>
    <SelectWorkflowModal ref="workflowRef" @ok="saveTools"/>
    <a-modal
        v-model:open="identVisible"
        :confirm-loading="identSaving"
        :title="t('title_edit_identifier')"
        @ok="saveIdent">
      <div>{{ t('label_identifier') }}</div>
      <div class="mt4">
        <a-input v-model:value="identValue" @input="identInput" :placeholder="t('ph_input_identifier')" :maxlength="32"/>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import {h, ref, onMounted} from 'vue';
import {message, Modal} from 'ant-design-vue'
import {
  PlusOutlined,
  CheckCircleFilled,
  EllipsisOutlined,
  AppstoreOutlined,
  ExclamationCircleFilled
} from '@ant-design/icons-vue';
import McpStore from "@/views/robot/robot-list/components/mcp-store.vue";
import SelectWorkflowModal from "@/views/robot/robot-list/components/select-workflow-modal.vue";
import {
  saveMcpTool,
  delMcpTool,
  editMcpTool,
  getMcpServers,
  updateMcpSrvStatus,
  delMcpServer
} from "@/api/robot/mcp.js";
import {jsonDecode} from "@/utils/index.js";
import { useI18n } from '@/hooks/web/useI18n';

const { t } = useI18n('views.robot.robot-list.components.mcp-box');

const mcpRef = ref(null)
const workflowRef = ref(null)
const loading = ref(true)
const statusLoading = ref(false)
const appInfo = ref(null)
const appData = ref([])
const tools = ref([])
const identVisible = ref(false)
const identSaving = ref(false)
const identValue = ref('')
const identRecord = ref(null)

onMounted(() => {
  init()
})

function init() {
  loading.value = true
  loadData()
}

function loadData() {
  getMcpServers().then(res => {
    appInfo.value = res?.data?.info || {}
    tools.value = res?.data?.tool_list || {}

    appInfo.value.status_bool = (appInfo.value?.publish_status == 1)
    tools.value.forEach(item => {
      item.params = jsonDecode(item.params, [])
    })
    appDataFormat()
  }).finally(() => {
    loading.value = false
  })
}

function statusChange() {
  const cancel = () => {
    appInfo.value.status_bool = !appInfo.value.status_bool
  }
  const ok = () => {
    statusLoading.value = true
    updateMcpSrvStatus({
      server_id: appInfo.value.id,
      publish_status: Number(appInfo.value.status_bool)
    }).then(res => {
      message.success(t('msg_operation_completed'))
    }).finally(() => {
      statusLoading.value = false
    }).catch(() => {
      cancel()
    })
  }
  if (!appInfo.value.status_bool) {
    Modal.confirm({
      title: t('title_confirm_close_mcp'),
      content: t('msg_close_mcp_warning'),
      okText: t('btn_confirm'),
      cancelText: t('btn_cancel'),
      onOk: () => ok(),
      onCancel: () =>  cancel()
    })
  } else {
    ok()
  }
}

function appDataFormat() {
  appData.value = [
    {
      title: 'label_url',
      value: `${window.location.origin}/mcp`
    },
    {
      title: 'label_request_method',
      value: 'POST'
    },
    {
      title: 'label_auth_method',
      value: 'Service token / API key'
    },
    {
      title: 'label_request_param_name',
      value: 'Authorization'
    },
    {
      title: 'label_api_key',
      value: `Bearer ${appInfo.value?.api_key}`
    },
  ]
}

function showMcpModal() {
  let _info = null
  if (appInfo.value?.id > 0) {
    _info = {
      avatar: appInfo.value.avatar,
      name: appInfo.value.name,
      description: appInfo.value.description,
      server_id: appInfo.value.id,
    }
  }
  mcpRef.value.show(_info)
}

function delMcp() {
  Modal.confirm({
    title: t('title_confirm_delete_mcp'),
    content: t('msg_delete_mcp_warning'),
    okText: t('btn_confirm'),
    cancelText: t('btn_cancel'),
    onOk: () => {
      delMcpServer({server_id: appInfo.value.id}).then(() => {
        message.success(t('msg_deleted'))
        loadData()
      })
    }
  })
}

function showWorkflowModal() {
  workflowRef.value.show(tools.value.map(i => i.robot_id))
}

function showIdentModal(item) {
  identRecord.value = item
  identValue.value = item.name
  identVisible.value = true
}

const identInput = e => {
  e.target.value = e.target.value.replace(/[^a-zA-Z0-9_]/g, '')
  identValue.value = e.target.value
}

function saveIdent() {
  if (!identValue.value) return message.warning(t('msg_input_identifier_warning'))
  identSaving.value = true
  editMcpTool({
    tool_id: identRecord.value.id,
    name: identValue.value,
  }).then(res => {
    identRecord.value.name = identValue.value
    identVisible.value = false
    message.success(t('msg_saved'))
  }).finally(() => {
    identSaving.value = false
  })
}

function delTool(item) {
  Modal.confirm({
    title: t('title_confirm_delete_tool'),
    okText: t('btn_confirm'),
    cancelText: t('btn_cancel'),
    onOk: () => {
      delMcpTool({tool_id: item.id}).then(() => {
        message.success(t('msg_deleted'))
        loadData()
      })
    }
  })
}

function saveTools(ids) {
  saveMcpTool({
    server_id: appInfo.value.id,
    robot_id: ids.toString()
  }).then(res => {
    loadData()
    message.success(t('msg_saved'))
  })
}
</script>

<style scoped lang="less">
._mcp-box {
  padding-right: 32px;

  .loading-box {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 10vh 20vw;
  }

  .base-info-box {
    padding: 24px;
    border-radius: 6px;
    border: 1px solid transparent; /* 必须透明 */
    background: linear-gradient(#fff, #fff) padding-box, linear-gradient(90deg, #70D9FF, #3C01FF, #ECB3FF) border-box;

    .app-info {
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
        align-items: center;
        justify-content: space-between;

        .right-box {
          display: flex;
          align-items: center;
          gap: 12px;
        }

        .name {
          color: #262626;
          font-size: 16px;
          font-weight: 600;
        }

        .desc {
          color: #8c8c8c;
          font-size: 14px;
          font-weight: 400;
          margin-top: 4px;
        }
      }
    }

    .app-data {
      display: flex;
      align-items: center;
      gap: 16px;
      margin-top: 16px;
      color: #595959;
      font-size: 14px;
      font-weight: 400;

      .title {
        color: #8c8c8c;
        font-size: 12px;
      }
    }
  }

  .tools-box {
    margin-top: 24px;

    .head-box {
      color: #262626;
      font-size: 14px;
      font-weight: 600;
    }

    .tools-list {
      margin-top: 16px;
      display: flex;
      flex-wrap: wrap;
      gap: 16px;

      .tools-item {
        flex: 0 0 calc((100% - 16px) / 2);
        padding: 24px;
        border-radius: 12px;
        border: 1px solid #E4E6EB;
        display: flex;
        justify-content: space-between;
        align-items: flex-end;
        gap: 8px;

        .tit {
          color: #262626;
          font-size: 14px;
          font-weight: 600;
        }

        .key {
          color: #8c8c8c;
          font-size: 12px;
          font-weight: 400;
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
  }
}

.status-tag {
  display: inline-block;
  padding: 2px 6px;
  align-items: center;
  gap: 4px;
  border-radius: 6px;
  background: #CAFCE4;
  white-space: nowrap;
  font-size: 14px;
  font-weight: 500;

  &.finished {
    background-color: #CAFCE4;
    color: #21A665;
  }

  &.waiting {
    background-color: #F0F0F0;
    color: #595959;
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


.empty-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin: 40px 0;

  img {
    width: 200px;
    height: 200px;
  }

  .title {
    color: #262626;
    font-size: 16px;
    font-weight: 600;
    line-height: 24px;
  }

  .btn {
    margin-top: 24px;
  }
}

.cFB363F {
  color: #FB363F;
}

.mt4 {
  margin-top: 4px;
}
</style>
