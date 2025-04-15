<style lang="less" scoped>
.setting-box {
  .actions-box {
    display: flex;
    align-items: center;
    line-height: 22px;
    font-size: 14px;
    color: #595959;

    .action-btn {
      cursor: pointer;
    }

    .save-btn {
      color: #2475fc;
    }
  }

  .library-list {
    display: flex;
    flex-flow: row wrap;
    gap: 16px;
    padding: 0 16px 0 16px;

    .library-item {
      position: relative;
      width: 336px;
      padding: 14px 16px;
      border-radius: 2px;
      border: 1px solid #d8dde5;
      background-color: #fff;
      cursor: pointer;

      .library-name {
        width: 100%;
        line-height: 22px;
        font-size: 14px;
        font-weight: 600;
        color: #262626;
        display: flex;
        align-items: center;
        gap: 4px;
      }
      .library-name-text {
        max-width: 240px;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }
      .graph-tag {
        width: 42px;
        height: 18px;
        display: flex;
        align-items: center;
        justify-content: center;
        border: 1px solid #00000026;
        background: #0000000a;
        border-radius: 6px;
        font-size: 12px;
        line-height: 16px;
        color: #bfbfbf;
        &.active {
          border: 1px solid #99bffd;
          color: #2475fc;
          background: #fff;
        }
      }

      .library-intro {
        width: 100%;
        line-height: 20px;
        font-size: 12px;
        color: #8c8c8c;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }

      .close-btn {
        position: absolute;
        top: 28px;
        right: 16px;
        font-size: 16px;
        color: #8c8c8c;
        cursor: pointer;
      }
    }
  }
  .setting-info-block {
    padding: 16px;
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    color: #595959;
    line-height: 22px;
    .set-item{
      display: flex;
      align-items: center;
    }
  }
}
</style>

<template>
  <edit-box
    class="setting-box"
    title="关联知识库"
    icon-name="guanlianzhishiku"
    v-model:isEdit="isEdit"
    :bodyStyle="{ padding: 0 }"
  >
    <template #extra>
      <div class="actions-box">
        <a-flex :gap="8">
          <a-button size="small" @click="handleOpenRecallSettingsAlert">召回设置</a-button>
          <a-button size="small" @click="handleOpenSelectLibraryAlert">关联知识库</a-button>
        </a-flex>
      </div>
    </template>
    <div class="library-list" v-if="selectedLibraryRows.length > 0">
      <div
        class="library-item"
        v-for="item in selectedLibraryRows"
        :key="item.id"
        @click="toLibraryDetail(item)"
      >
        <span class="close-btn" @click.stop="handleRemoveCheckedLibrary(item)">
          <CloseCircleOutlined />
        </span>
        <div class="library-name">
          <a-tooltip>
            <template #title>{{ item.graph_switch == 0 ? '未' : '已' }}开启知识图谱生成</template>
            <div class="graph-tag" :class="{ active: item.graph_switch == 1 }">Graph</div>
          </a-tooltip>
          <div class="library-name-text">{{ item.library_name }}</div>
        </div>
        <div class="library-intro">{{ item.library_intro }}</div>
      </div>
    </div>
    <div class="setting-info-block">
      <div class="set-item">
        检索模式：
        <span v-if="formState.search_type == 1">混合检索</span>
        <span v-if="formState.search_type == 2">向量检索</span>
        <span v-if="formState.search_type == 3">全文检索</span>
        <span v-if="formState.search_type == 4">知识图谱检索</span>
      </div>
      <div class="set-item">
        Top K：
        <span>{{ formState.top_k }}</span>
      </div>
      <div class="set-item" v-if="formState.search_type <= 2">
        相似度阈值：
        <span>{{ formState.similarity }}</span>
      </div>
      <div class="set-item">
        Rerank模型：
        <span v-if="formState.rerank_status == 0">关</span>
        <span v-else>{{ formState.rerank_use_model }}</span>
      </div>
    </div>
    <LibrarySelectAlert ref="librarySelectAlertRef" @change="onChangeLibrarySelected" />
    <RecallSettingsAlert ref="recallSettingsAlertRef" @change="onChangeRecallSettings" />
    <NoOpenGraphModal :list="noOpenLibraryList" @refreshList="getList" ref="noOpenGraphModalRef" />
  </edit-box>
</template>

<script setup>
import { useStorage } from '@/hooks/web/useStorage'
import { getLibraryList } from '@/api/library/index'
import { ref, reactive, inject, watchEffect, computed, toRaw, onMounted } from 'vue'
import { CloseCircleOutlined } from '@ant-design/icons-vue'
import { Modal } from 'ant-design-vue'
import EditBox from '../edit-box.vue'
import LibrarySelectAlert from './library-select-alert.vue'
import RecallSettingsAlert from './recall-settings-alert.vue'
import NoOpenGraphModal from './no-open-graph-modal.vue'

const { getStorage, removeStorage } = useStorage('localStorage')

const isEdit = ref(false)

const { robotInfo, updateRobotInfo } = inject('robotInfo')

const formState = reactive({
  library_ids: [],
  rerank_status: 0,
  rerank_use_model: undefined,
  rerank_model_config_id: '',
  top_k: 0,
  similarity: 0,
  search_type: 1
})

// 知识库
const libraryList = ref([])
const librarySelectAlertRef = ref(null)
const selectedLibraryRows = computed(() => {
  return libraryList.value.filter((item) => {
    return formState.library_ids.includes(item.id)
  })
})

const noOpenLibraryList = computed(() => {
  return selectedLibraryRows.value.filter((item) => item.graph_switch == 0)
})

// 移除知识库
const handleRemoveCheckedLibrary = (item) => {
  let index = formState.library_ids.indexOf(item.id)

  formState.library_ids.splice(index, 1)

  onSave()
}

const onChangeLibrarySelected = (checkedList) => {
  formState.library_ids = [...checkedList]

  onSave()
}

const handleOpenSelectLibraryAlert = () => {
  librarySelectAlertRef.value.open([...formState.library_ids])
}

// 召回设置
const recallSettingsAlertRef = ref(null)

const handleOpenRecallSettingsAlert = () => {
  recallSettingsAlertRef.value.open(toRaw(formState))
}
const noOpenGraphModalRef = ref(null)

const onChangeRecallSettings = (data) => {
  formState.rerank_status = data.rerank_status
  formState.rerank_use_model = data.rerank_use_model
  formState.rerank_model_config_id = data.rerank_model_config_id
  formState.top_k = data.top_k
  formState.similarity = data.similarity
  formState.search_type = data.search_type
  if (data.search_type == 1 || data.search_type == 4) {
    if (noOpenLibraryList.value.length > 0) {
      noOpenGraphModalRef.value.show()
    }
  }
  onSave()
}

const onSave = () => {
  let formData = { ...toRaw(formState) }

  formData.library_ids = formData.library_ids.join(',')

  updateRobotInfo({ ...formData })
}

// 获取知识库
const getList = async () => {
  const res = await getLibraryList({ type: '' })
  if (res) {
    libraryList.value = res.data || []
  }
}

// 显示未关联知识库提示
const handleShowNoLibraryTip = () => {
  Modal.confirm({
    title: '还未关联知识库',
    content: '关联知识库后，机器人会根据知识库作答。',
    okText: '立即关联',
    cancelText: '暂不关联',
    onOk() {
      handleOpenSelectLibraryAlert()
    }
  })
}

formState.rerank_use_model = robotInfo.rerank_use_model || undefined

watchEffect(() => {
  formState.library_ids = robotInfo.library_ids.split(',')
  formState.rerank_status = robotInfo.rerank_status || 0
  formState.rerank_use_model = robotInfo.rerank_use_model || undefined
  formState.rerank_model_config_id = robotInfo.rerank_model_config_id || ''
  formState.top_k = robotInfo.top_k
  formState.similarity = robotInfo.similarity
  formState.search_type = robotInfo.search_type
})

const toLibraryDetail = (item) => {
  window.open(`#/library/details/knowledge-document?id=${item.id}`)
}

onMounted(() => {
  getList()

  // 检查本地缓存中的showNoLibraryTip值
  const showNoLibraryTip = getStorage('showNoLibraryTip')
  if (showNoLibraryTip) {
    handleShowNoLibraryTip()

    removeStorage('showNoLibraryTip')
  }
})
</script>
