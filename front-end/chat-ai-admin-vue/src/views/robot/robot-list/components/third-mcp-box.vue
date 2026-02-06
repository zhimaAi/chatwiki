<template>
  <div class="_mcp-box">
    <div v-if="loading" class="loading-box">
      <a-spin/>
    </div>
    <div v-else-if="list.length" class="list-box">
      <div
        v-for="(item, i) in list"
        @click="showDetail(item)"
        :key="i"
        class="app-item"
      >
        <div class="base-info-box">
          <img class="avatar" :src="item.avatar"/>
          <div class="info-box">
            <div class="name zm-line1">{{ item.name }}</div>
            <div>
              <div v-if="item.has_auth == 1" class="auth-tag">{{ t('tag_authorized') }}</div>
              <div v-else class="auth-tag fail">{{ t('tag_unauthorized') }}</div>
            </div>
          </div>
        </div>
        <div class="zm-line1">{{ item.url }}</div>
        <div class="bottom">
          <div class="extra-box">
            <div class="extra-item">{{ t('label_available_tools') }}{{ item.tools.length }}</div>
            <div class="extra-item">{{ item.up_time_text }} {{ t('label_updated') }}</div>
          </div>
          <a-dropdown>
            <a-button @click.stop size="small" :icon="h(EllipsisOutlined)"/>
            <template #overlay>
              <a-menu>
                <a-menu-item @click.stop="editApp(item)">{{ t('btn_edit') }}</a-menu-item>
                <a-menu-item @click.stop="delApp(item)"><span class="cFB363F">{{ t('btn_delete') }}</span></a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </div>
    </div>
    <div v-else class="empty-box">
      <img src="@/assets/empty.png"/>
      <div class="title">{{ t('title_no_mcp_added') }}</div>
      <div class="desc-box"><div class="link" @click="goMcpSquare">{{ t('link_more_mcp') }}</div></div>
      <a-button @click="showMcpModal" class="btn" type="primary">{{ t('btn_add_now') }}</a-button>
    </div>

    <ThirdMcpDetail ref="mcpDetailRef" @del="delApp" @edit="editApp" @auth="loadData"/>
    <ThirdMcpStore ref="mcpStoreRef" @ok="update"/>
  </div>
</template>

<script setup>
import {h, ref, onMounted} from 'vue';
import {message, Modal} from 'ant-design-vue';
import {EllipsisOutlined} from '@ant-design/icons-vue';
import ThirdMcpDetail from "@/views/robot/robot-list/components/third-mcp-detail.vue";
import ThirdMcpStore from "@/views/robot/robot-list/components/third-mcp-store.vue";
import {delTMcpProvider, getTMcpProviders} from "@/api/robot/thirdMcp.js";
import { jsonDecode, timeNowGapFormat } from "@/utils/index.js";
import { useRouter } from "vue-router";
import { useI18n } from '@/hooks/web/useI18n';

const { t } = useI18n('views.robot.robot-list.components.third-mcp-box')

const router = useRouter()

const mcpDetailRef = ref(null)
const mcpStoreRef = ref(null)
const loading = ref(true)
const list = ref([])
const emit = defineEmits(['listLoaded'])

onMounted(() => {
  init()
})

function init() {
  loading.value = true
  loadData()
}

function update() {
  loadData()
  mcpDetailRef.value && mcpDetailRef.value.refresh()
}

function loadData() {
  getTMcpProviders().then(res => {
    let _list = res?.data || []
    _list.forEach(item => {
      item.tools = jsonDecode(item.tools, [])
      item.up_time_text = timeNowGapFormat(item.update_time)
    })
    list.value = _list
    emit('listLoaded', list.value.length)
  }).finally(() => {
    loading.value = false
  })
}

function showDetail(item) {
  mcpDetailRef.value.show(item)
}

function showMcpModal() {
  mcpStoreRef.value.show()
}

function editApp(item) {
  mcpStoreRef.value.show(item)
}

function goMcpSquare () {
  const url = router.resolve({
    path: '/mcp/index'
  }).href
  window.open(url, '_blank')
}

function delApp(item) {
  Modal.confirm({
    title: t('modal_confirm_delete_title'),
    content: t('modal_confirm_delete_content'),
    okText: t('btn_confirm'),
    cancelText: t('btn_cancel'),
    onOk: () => {
      delTMcpProvider({provider_id: item.id}).then(() => {
        message.success(t('msg_deleted'))
        loadData()
        mcpDetailRef.value && mcpDetailRef.value.hide()
      })
    }
  })
}

defineExpose({
  update,
})
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
        flex-wrap: wrap;
        gap: 2px 12px;
        color: #8c8c8c;
        font-size: 12px;
        font-weight: 400;
      }
    }
  }

  .empty-box {
    display: flex;
    flex-direction: column;
    align-items: center;

    img {
      width: 200px;
      height: 200px;
    }

    .title {
      color: #262626;
      font-size: 16px;
      font-style: normal;
      font-weight: 600;
      line-height: 24px;
    }

    .desc-box {
      margin-top: 12px;
      display: flex;
      align-items: center;
      gap: 4px;
      color: #8c8c8c;
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 22px;

      .link {
        cursor: pointer;
        color: #2475fc;
        font-size: 14px;
        font-style: normal;
        font-weight: 400;
        line-height: 22px;

        &:hover {
          opacity: 0.8;
        }
      }
    }

    .btn {
      margin-top: 24px;
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

.cFB363F {
  color: #FB363F;
}
</style>
