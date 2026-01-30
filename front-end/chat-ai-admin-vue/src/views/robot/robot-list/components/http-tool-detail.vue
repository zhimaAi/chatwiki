<template>
  <div class="_http_detail_box">
    <div class="base-info-box">
      <div class="app-info">
        <img class="avatar" :src="avatarUrl"/>
        <div class="info-box">
          <div class="left-box">
            <div class="name">
              {{ tool.name }}
            </div>
            <div class="desc">{{ tool.description }}</div>
          </div>
          <div class="right-box">
            <a-button @click="$emit('addTool', tool)" type="primary" ghost :icon="h(PlusOutlined)">添加工具</a-button>
            <a-dropdown>
              <a-button :icon="h(EllipsisOutlined)"/>
              <template #overlay>
                <a-menu>
                  <a-menu-item @click="$emit('editBase', tool)">编辑</a-menu-item>
                  <a-menu-item @click="$emit('deleteBase', tool)"><span class="cFB363F">删除</span></a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
        </div>
      </div>
      <div class="app-data">
        <div class="data-item">
          <div class="title">创建时间</div>
          <div class="value">{{ formatTs(tool.create_time) }}</div>
        </div>
        <a-divider type="vertical" style="height: 24px;"/>
        <div class="data-item">
          <div class="title">更新时间</div>
          <div class="value">{{ formatTs(tool.update_time) }}</div>
        </div>
      </div>
    </div>
    <div class="tools-box">
      <div class="head-box">
        <AppstoreOutlined class="icon"/>
        工具列表（{{ nodesList.length || 0 }}）
      </div>
      <div v-if="nodesList.length" class="tools-list">
        <div v-for="(item, i) in nodesList" :key="i" class="tools-item">
          <div class="left">
            <div class="tit">{{ item.displayTitle }}</div>
            <div class="key">请求地址：{{ item.rawurl }}</div>
            <a-tooltip :title="getTooltipTitle(item.displayDesc, item)" placement="top">
              <div class="desc" :class="`titleRef_${item.id}`">{{item.displayDesc}}</div>
            </a-tooltip>
            <div class="params">
              <a-tag v-for="field in (item.displayParams || []).slice(0, 5)" :key="field.key" :bordered="false">{{ field.key }}</a-tag>
              <a-popover placement="right">
                <template #content>
                  <div class="params-box">
                    <div v-for="field in (item.displayParams || [])" :key="field.key" class="param-item">
                      <div class="field">
                        <span class="name">{{ field.key }}</span>
                        <span class="type">{{ field.typ }}</span>
                        <span v-if="field.required" class="required">必填</span>
                      </div>
                      <div class="desc">{{ field.desc }}</div>
                    </div>
                  </div>
                </template>
                <a>参数</a>
              </a-popover>
            </div>
          </div>
          <a-dropdown>
            <a-button size="small" :icon="h(EllipsisOutlined)"/>
            <template #overlay>
              <a-menu>
                <a-menu-item @click="$emit('editTool', item.originNode)">编辑</a-menu-item>
                <a-menu-item @click="delItem(item.originNode)"><span class="cFB363F">删除</span></a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </div>
      <div v-else class="empty-box">
        <img src="@/assets/empty.png"/>
        <div class="title">暂未添加工具</div>
        <a-button @click="$emit('addTool', tool)" class="btn" type="primary">立即添加</a-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, h } from 'vue'
import { AppstoreOutlined, EllipsisOutlined } from '@ant-design/icons-vue';
import { PlusOutlined } from '@ant-design/icons-vue';
import { delHttpToolItem } from '@/api/robot/http_tool.js'
import { Modal } from 'ant-design-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'

const props = defineProps({
  tool: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['addTool', 'editTool', 'ok', 'editBase', 'deleteBase'])

function delItem(item) {
  Modal.confirm({
    title: '删除确认',
    icon: h(ExclamationCircleOutlined),
    content: '确定删除该工具吗？删除后不可恢复',
    okText: '删 除',
    cancelText: '取 消',
    onOk() {
      return delHttpToolItem({ id: item.id }).then(() => {
        emit('ok')
      })
    }
  })
}

const avatarUrl = computed(() => {
  const s = String(props.tool?.avatar || '').trim().replace(/[`"']/g, '')
  return s || '/src/assets/img/workflow/node-service.svg'
})
function formatTs(ts) {
  const n = Number(ts || 0)
  if (!n) return '—'
  const ms = n > 1e12 ? n : n * 1000
  return dayjs(ms).format('YYYY-MM-DD HH:mm')
}

const nodesList = computed(() => {
  const arr = Array.isArray(props.tool?.nodes) ? props.tool.nodes : []
  return arr.map(n => {
    const info = n.http_tool_info || {}
    const np = n.node_params || {}
    const curl = np.curl || {}
    const paramsArr = Array.isArray(curl.params) ? curl.params : []
    const bodyArr = Array.isArray(curl.body) ? curl.body : []
    const prs = [...paramsArr, ...bodyArr].filter(x => String(x?.key || '').trim().length > 0)
    const displayParams = prs.map(x => ({ key: x.key, typ: x.value }))
    return {
      id: Number(n.id || 0),
      displayTitle: info.http_tool_node_name || n.node_name || '',
      displayIdent: n.node_key || info.http_tool_node_key || '',
      displayDesc: n.node_description || info.http_tool_node_description || '',
      rawurl: curl.rawurl || '',
      displayParams,
      originNode: n
    }
  })
})

// 获取 tooltip 标题
function getTooltipTitle(text, record) {
  if (!text) return null
  
  // 创建临时元素来测量文本宽度
  const canvas = document.createElement('canvas')
  const context = canvas.getContext('2d')
  // 14px 根据实际字体大小修改
  context.font = '14px -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif'
  
  const textWidth = context.measureText(text).width

  const titleRef = document.querySelector(`.titleRef_${record.create_time}`)
  if (titleRef) {
    record.title_width = titleRef.offsetWidth
  }
  const maxWidth = record?.title_width || 1250
  return textWidth > maxWidth ? text : null
}
</script>

<style scoped lang="less">
._http_detail_box {
  padding-right: 32px;
  .base-info-box {
    padding: 24px;
    border-radius: 6px;
    border: 1px solid transparent;
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
            max-width: calc((100vw / 2) - 170px);
            color: #595959;
            font-size: 14px;
            margin-top: 8px;
            height: 44px;
            line-height: 22px;
            font-weight: 400;
            word-break: break-all;
            // 超出2行显示省略号
            overflow: hidden;
            text-overflow: ellipsis;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            line-clamp: 2;
            -webkit-box-orient: vertical;
          }
          .params {
            margin-top: 8px;
          }
        }
      }
      .bind-list {
        margin-top: 24px;
        .head-box {
          color: #262626;
          font-size: 14px;
          font-weight: 600;
          margin-bottom: 8px;
        }
        .row {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 12px 16px;
          border: 1px solid #E4E6EB;
          border-radius: 8px;
          &:not(:last-child) {
            margin-bottom: 8px;
          }
          .info {
            display: flex;
            align-items: center;
            gap: 12px;
            .avatar {
              width: 32px;
              height: 32px;
              border-radius: 8px;
            }
            .name {
              color: #262626;
              font-size: 14px;
              font-weight: 600;
            }
            .desc {
              color: #8c8c8c;
              font-size: 12px;
            }
          }
          .actions {
            display: flex;
            align-items: center;
            gap: 8px;
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
