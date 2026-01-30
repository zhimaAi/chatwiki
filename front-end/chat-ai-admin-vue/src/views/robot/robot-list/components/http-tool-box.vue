<template>
  <div class="_http_tool_box">
    <div v-if="loading" class="loading-box">
      <a-spin/>
    </div>
    <template v-else>
      <template v-if="!detailTool">
        <div class="list-head">
          <div class="title">HTTP工具</div>
        </div>
        <div v-if="tools.length" class="list-box">
          <div
            v-for="(item, i) in tools"
            :key="i"
            @click="openDetail(item)"
            class="app-item"
          >
            <div class="base-info-box">
              <img class="avatar" :src="formatAvatar(item.avatar)"/>
              <div class="info-box">
                <div class="name zm-line1">{{ item.name }}</div>
                <a-tooltip :title="getTooltipTitle(item.description, item, 14, 2, 10)" placement="top">
                  <div class="description" :ref="el => setDescRef(el, item)">{{ item.description }}</div>
                </a-tooltip>
              </div>
            </div>
            <div class="bottom">
              <div class="extra-box">
                <div class="extra-item" v-if="item.node_count > 0">{{ item.node_count }}个可用工具</div>
                <div class="extra-item" v-else>暂无可用的工具</div>
              </div>
              <div class="bottom-right">
                <div class="update-time">{{ item.up_time_text }}更新</div>
                <a-dropdown>
                  <a-button @click.stop size="small" :icon="h(EllipsisOutlined)"/>
                  <template #overlay>
                    <a-menu>
                      <a-menu-item @click.stop="editItem(item)">编辑</a-menu-item>
                      <a-menu-item @click.stop="delItem(item)"><span class="cFB363F">删除</span></a-menu-item>
                    </a-menu>
                  </template>
                </a-dropdown>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="empty-box">
          <img src="@/assets/empty.png"/>
          <div class="title">暂未添加HTTP工具</div>
          <a-button @click="showStore" class="btn" type="primary">立即添加</a-button>
        </div>
      </template>
      <template v-else>
        <div class="breadcrumb">
          <a @click="backToList">HTTP工具</a> <span>/</span> <span class="cur">{{ detailTool.name }}</span>
        </div>
        <HttpToolDetail
          :tool="detailTool"
          @back="backToList"
          @editTool="showEditToolModal"
          @ok="loadData"
          @addTool="showToolModal"
          @editBase="editItem"
          @deleteBase="delItem"
        />
      </template>
      <HttpToolStore ref="storeRef" @ok="loadData"/>
      <AddHttpToolModal :toolId="detailTool?.id || 0" ref="addToolRef" @ok="loadData" />
    </template>
  </div>
</template>

<script setup>
import {h, ref, onMounted} from 'vue';
import {message, Modal} from 'ant-design-vue'
import { EllipsisOutlined } from '@ant-design/icons-vue';
import HttpToolStore from './http-tool-store.vue'
import AddHttpToolModal from './add-http-tool-modal.vue'
import HttpToolDetail from './http-tool-detail.vue'
import { getHttpTools, delHttpTool } from '@/api/robot/http_tool.js'
import { setDescRef, getTooltipTitle, timeNowGapFormat } from '@/utils/index'

const storeRef = ref(null)
const addToolRef = ref(null)
const loading = ref(true)
const tools = ref([])
const detailTool = ref(null)

onMounted(() => {
  loadData()
})

function loadData() {
  loading.value = true
  const with_nodes = 1
  getHttpTools({ page: 1, size: 9999, with_nodes }).then(res => {
    const list = (res?.data?.list || []).map(it => ({
      ...it,
      up_time_text: timeNowGapFormat(it.update_time),
      avatar: formatAvatar(it.avatar)
    }))
    tools.value = list
    if (detailTool.value) {
      const cur = list.find(i => i.id === detailTool.value.id)
      if (cur) {
        detailTool.value = cur
      }
    }
  }).finally(() => {
    loading.value = false
  })
}

function showStore() {
  storeRef.value.show()
}

function openDetail(item) {
  detailTool.value = item
}
function backToList() {
  detailTool.value = null
}

function editItem (item) {
  storeRef.value.show({
    id: item.id,
    name: item.name,
    description: item.description,
    avatar: item.avatar
  })
}

function showToolModal(item) {
  // 改为添加工具
  addToolRef.value && addToolRef.value.show()
}

function showEditToolModal (item) {
  addToolRef.value && addToolRef.value.show(item)
}

function delItem(item) {
  Modal.confirm({
    title: '删除确认',
    content: '确认删除该HTTP工具么？删除工具后使用到该工具的工作流可能无法正常运行。确认删除么？',
    okText: '确定',
    cancelText: '取消',
    onOk: () => {
      delHttpTool({ id: item.id }).then(() => {
        message.success('已删除')
        if (detailTool.value && detailTool.value.id === item.id) {
          detailTool.value = null
        }
        loadData()
      })
    }
  })
}

function formatAvatar(avatar) {
  const s = String(avatar || '').trim().replace(/[`"']/g, '')
  return s || '/src/assets/img/workflow/node-service.svg'
}

defineExpose({
  showStore
})
</script>

<style scoped lang="less">
._http_tool_box {
  padding-right: 32px;
  height: 100%;
  overflow-y: auto;

  .loading-box {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 10vh 20vw;
  }
  .list-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;
    .title {
      color: #262626;
      font-size: 16px;
      font-weight: 600;
    }
  }
  .list-box {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 16px;
    .app-item {
      flex-shrink: 0;
      width: calc((100% - 16px * 3) / 4);
      max-width: 600px;
      display: flex;
      flex-direction: column;
      gap: 12px;
      padding: 24px;
      border-radius: 12px;
      border: 1px solid #E4E6EB;
      color: #595959;
      font-size: 14px;
      font-weight: 400;
      cursor: pointer;
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
          .name {
            color: #262626;
            font-size: 16px;
            font-weight: 600;
            margin-bottom: 4px;
          }
          .description {
            height: 44px;
            line-height: 22px;
            margin-top: 12px;
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
        }
      }
      .bottom {
        display: flex;
        align-items: center;
        justify-content: space-between;
      }
      .extra-box {
        display: flex;
        align-items: center;
        gap: 12px;
        color: #8c8c8c;
        font-size: 12px;
        font-weight: 400;
      }

      .bottom-right {
        display: flex;
        align-items: center;
        gap: 12px;
        color: #8c8c8c;
        font-size: 12px;
        font-weight: 400;
      }
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
.breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 4px 16px;
  color: #595959;
  .cur {
    color: #262626;
    font-weight: 600;
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
</style>
