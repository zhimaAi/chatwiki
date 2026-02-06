<style lang="less" scoped>
.http-tool-list-box {
  .search-box {
    margin: 0 8px 8px;
  }
  .node-list {
    .node-item {
      .node-info {
        display: flex;
        align-items: center;
        padding: 4px 8px;
        border-radius: 6px;
        cursor: pointer;
        &:hover {
          background: #e4e6eb;
        }
        .avatar {
          width: 20px;
          height: 20px;
          flex-shrink: 0;
          border-radius: 4px;
          margin-right: 8px;
        }
        .info {
          flex: 1;
          display: flex;
          align-items: center;
          justify-content: space-between;
          .name {
            color: #262626;
            font-size: 14px;
            font-weight: 400;
          }
          .total {
            color: #8c8c8c;
            font-size: 12px;
            font-weight: 400;
          }
        }
      }
      .node-tools {
        margin-left: 28px;
        border-left: 1px solid #d9d9d9;
      }
      .node-tool-item {
        padding: 2px 8px;
        border-radius: 6px;
        cursor: pointer;
        &:hover {
          background: #e4e6eb;
        }
      }
    }
  }
  .empty-box {
    text-align: center;
  }
  .more-link {
    padding: 8px 16px 0;
    border-top: 1px solid #F0F0F0;
    color: #2475FC;
  }
}

.node-info-pop {
  display: flex;
  flex-direction: column;
  gap: 12px;
  color: #595959;
  font-size: 14px;
  font-weight: 400;
  min-width: 260px;

  .info {
    display: flex;
    align-items: center;
    color: #262626;
    font-size: 16px;

    .avatar {
      width: 62px;
      height: 62px;
      border-radius: 14.59px;
      flex-shrink: 0;
      margin-right: 12px;
    }
  }

  .extra {
    color: #8c8c8c;
    font-size: 12px;
  }
}
.params-box {
  max-height: 80vh;
  overflow-y: auto;

  .param-item {
    max-width: 500px;

    &:first-child {
      border-bottom: 1px solid #e4e6eb;
      padding-bottom: 16px;
    }
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
    }

    .desc {
      color: #8c8c8c;
      font-size: 14px;
      margin-top: 4px;
    }
  }
}
</style>

<template>
  <div class="http-tool-list-box">
    <div class="node-list" v-if="showList.length">
      <div v-for="(item, j) in showList" :key="j" class="node-item">
        <a-popover placement="right">
          <template #content>
            <div class="node-info-pop">
              <div class="info">
                <img class="avatar" :src="item.avatar || '/src/assets/img/workflow/node-service.svg'" />
                <div class="name">{{ item.name }}</div>
              </div>
              <div class="extra">{{ t('label_node_count') }}{{ (item.nodes || []).length }}</div>
              <div>{{ item.description }}</div>
            </div>
          </template>
          <div class="node-info" @click="item.expand = !item.expand">
            <img class="avatar" :src="item.avatar || '/src/assets/img/workflow/node-service.svg'" />
            <div class="info">
              <span class="name">{{ item.name }}</span>
              <span class="total">
                {{ (item.nodes || []).length }} <DownOutlined v-if="item.expand"/> <RightOutlined v-else/>
              </span>
            </div>
          </div>
        </a-popover>
        <div v-show="item.expand" class="node-tools">
          <a-popover
            v-for="(node, i) in (item.nodes || [])"
            :key="i"
            placement="right"
          >
            <template #content>
              <div class="params-box">
                <div class="param-item">
                  <div class="field">
                    <span class="name">{{ node.node_name }}</span>
                  </div>
                  <div class="desc">{{ node.node_description }}</div>
                </div>
                <div class="param-item">
                  <div class="field">
                    <span class="name">{{ formatMethod(node?.node_params?.curl?.method) }}</span>
                    <span class="type">{{ cleanupUrl(node?.node_params?.curl?.rawurl) }}</span>
                  </div>
                </div>
                <div class="param-item" v-if="Array.isArray(node?.node_params?.curl?.params) && node.node_params.curl.params.length">
                  <div class="field">
                    <span class="name">{{ t('label_query_params') }}</span>
                  </div>
                  <div class="desc">
                    {{ node.node_params.curl.params.map(p => p.key).filter(Boolean).join(', ') }}
                  </div>
                </div>
                <div class="param-item" v-if="Array.isArray(node?.node_params?.curl?.headers) && node.node_params.curl.headers.length">
                  <div class="field">
                    <span class="name">{{ t('label_request_headers') }}</span>
                  </div>
                  <div class="desc">
                    {{ node.node_params.curl.headers.map(h => h.key).filter(Boolean).join(', ') }}
                  </div>
                </div>
                <div class="param-item" v-if="Array.isArray(node?.node_params?.curl?.output) && node.node_params.curl.output.length">
                  <div class="field">
                    <span class="name">{{ t('label_output_fields') }}</span>
                  </div>
                  <div class="desc">
                    {{ node.node_params.curl.output.map(o => o.key).filter(Boolean).join(', ') }}
                  </div>
                </div>
              </div>
            </template>
            <div class="node-tool-item" @click="handleAddNode(item, node)">{{ node.node_name }}</div>
          </a-popover>
        </div>
      </div>
    </div>
    <div v-else class="empty-box">
      <img style="height: 200px" src="@/assets/empty.png" />
      <div>{{ t('msg_no_http_tools') }}</div>
      <a @click="handleOpenAddHttp">{{ t('btn_go_add') }}<RightOutlined /></a>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { RightOutlined, DownOutlined } from '@ant-design/icons-vue'
import { getHttpTools } from '@/api/robot/http_tool.js'
import { addHttpToolNode } from '../node-list'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-list-popup.http-tool-list')

const emit = defineEmits(['add'])

const props = defineProps({
  keyword: {
    type: String,
    default: ''
  }
})

const list = ref([])

onMounted(() => {
  loadData()
})

function loadData() {
  getHttpTools({ page: 1, size: 9999, with_nodes: 1 }).then(res => {
    const arr = res?.data?.list || []
    arr.forEach(i => i.expand = false)
    list.value = arr
  })
}

const showList = computed(() => {
  if (props.keyword) {
    const kw = props.keyword
    return list.value.filter((item) => {
      let info = (item.description || '') + (item.name || '')
      if (Array.isArray(item.nodes)) {
        item.nodes.forEach(n => {
          info += (n.node_name || '') + (n.node_description || '')
          const url = n?.node_params?.curl?.rawurl || ''
          info += url
        })
      }
      return info.indexOf(kw) > -1
    })
  }
  return list.value
})

const handleAddNode = (tool, node) => {
  const curl = node?.node_params?.curl || {}
  const newNode = addHttpToolNode({
    name: node?.node_name || tool?.name,
    avatar: node?.http_tool_avatar || tool?.avatar,
    method: curl?.method,
    url: curl?.rawurl,
    headers: Array.isArray(curl?.headers) ? curl.headers : [],
    params: Array.isArray(curl?.params) ? curl.params : [],
    body: Array.isArray(curl?.body) ? curl.body : [],
    timeout: curl?.timeout || 60,
    output: Array.isArray(curl?.output) ? curl.output : [],
    http_auth: Array.isArray(curl?.http_auth) ? curl.http_auth : [],
    http_tool_info: node?.http_tool_info || {
      http_tool_name: tool?.name || '',
      http_tool_key: tool?.tool_key || '',
      http_tool_avatar: tool?.avatar || '',
      http_tool_description: tool?.description || '',
      http_tool_node_key: node?.node_key || '',
      http_tool_node_name: node?.node_name || '',
      http_tool_node_description: node?.node_description || ''
    }
  })
  emit('add', newNode)
}

function handleOpenAddHttp() {
  window.open('/#/robot/list?active=6', '_blank')
}

function formatMethod(m) {
  return String(m || '').toUpperCase() || 'GET'
}
function cleanupUrl(url) {
  return String(url || '').replace(/[`"]/g, '').trim()
}
</script>
