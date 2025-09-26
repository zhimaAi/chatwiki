<template>
  <div class="form-block" @mousedown.stop="">
    <a-form ref="formRef" layout="vertical" :model="formState">
      <div class="gray-block">
        <div class="gray-block-title">输入</div>
        <a-form-item label="关联知识库">
          <div class="konwledge-list-item" v-for="item in selectedLibraryRows" :key="item.id">
            <div class="avatar-box">
              <img :src="item.avatar" alt="" />
            </div>
            <div class="content-box">
              <div class="list-name">{{ item.library_name }}</div>
              <div class="list-intro">{{ item.library_intro }}</div>
            </div>
            <div class="btn-hover-wrap" @click="handleDelKonwledge(item)">
              <CloseCircleOutlined />
            </div>
          </div>
          <div class="btn-block">
            <div>
              <a-button
                @click="handleOpenSelectLibraryAlert"
                :icon="h(PlusOutlined)"
                block
                type="dashed"
                >添加知识库</a-button
              >
            </div>
            <div>
              <a-button
                @click="handleOpenRecallSettingsAlert"
                :icon="h(SettingOutlined)"
                block
                type="dashed"
                >召回设置</a-button
              >
            </div>
          </div>
        </a-form-item>

        <div class="diy-form-item">
          <div class="form-label">用户问题</div>
          <div class="form-content">流程开始>用户问题</div>
        </div>
      </div>
      <div class="gray-block mt16">
        <div class="gray-block-title">输出</div>
        <div class="options-item">
          <div class="option-label">知识库引用</div>
          <div class="option-type">知识库引用</div>
        </div>
      </div>
    </a-form>
    <LibrarySelectAlert
      ref="librarySelectAlertRef"
      @close="getList"
      @change="onChangeLibrarySelected"
    />
    <RecallSettingsAlert ref="recallSettingsAlertRef" @change="onChangeRecallSettings" />
  </div>
</template>

<script setup>
import { ref, reactive, watch, h, computed, toRaw, nextTick } from 'vue'
import LibrarySelectAlert from './library-select-alert.vue'
import RecallSettingsAlert from './recall-settings-alert.vue'
import { getLibraryList } from '@/api/library/index'
import { message, Modal } from 'ant-design-vue'
import {
  CloseCircleFilled,
  CloseCircleOutlined,
  QuestionCircleOutlined,
  UpOutlined,
  DownOutlined,
  LoadingOutlined,
  PlusOutlined,
  SettingOutlined,
  EditOutlined,
  SyncOutlined,
  ExclamationCircleOutlined
} from '@ant-design/icons-vue'

const props = defineProps({
  properties: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['setData'])
const formRef = ref()

const formState = reactive({
  library_ids: [],
  rerank_status: 0,
  rerank_use_model: undefined,
  rerank_model_config_id: void 0,
  top_k: 5,
  similarity: 0.5,
  search_type: 1
})
let lock = false
watch(
  () => props.properties,
  (val) => {
    try {
      if (lock) {
        return
      }
      let libs = JSON.parse(val.node_params).libs || {}
      libs = JSON.parse(JSON.stringify(libs))
      for (let key in libs) {
        if (key == 'library_ids') {
          formState[key] = libs[key] ? libs[key].split(',') : []
        } else {
          formState[key] = libs[key]
        }
      }
      setTimeout(() => {
        updata()
      }, 1000)
      lock = true
    } catch (error) {}
  },
  { immediate: true, deep: true }
)

watch(
  () => formState,
  (val) => {
    updata()
  },
  { deep: true }
)

const updata = () => {
  emit('setData', {
    ...formState,
    library_ids: formState.library_ids.join(','),
    node_params: JSON.stringify({
      libs: {
        ...formState,
        rerank_model_config_id: formState.rerank_model_config_id
          ? +formState.rerank_model_config_id
          : void 0,
        library_ids: formState.library_ids.join(',')
      }
    }),
    height: 500,
  })
}

const libraryList = ref([])
const librarySelectAlertRef = ref(null)
const selectedLibraryRows = computed(() => {
  return libraryList.value.filter((item) => {
    return formState.library_ids.includes(item.id)
  })
})

// 移除知识库
const handleDelKonwledge = (item) => {
  let index = formState.library_ids.indexOf(item.id)
  formState.library_ids.splice(index, 1)
}

const checkedHeader = (rule, value) => {
  // if (value == null) {
  //   return Promise.reject('请输入延迟发送时间')
  // }
  // if (!Number.isInteger(value / 0.5)) {
  //   return Promise.reject('必须为0.5秒的倍数')
  // }
  return Promise.resolve()
}

const onChangeLibrarySelected = (checkedList) => {
  getList()
  formState.library_ids = [...checkedList]
}

const handleOpenSelectLibraryAlert = () => {
  librarySelectAlertRef.value.open([...formState.library_ids])
}

// 召回设置
const recallSettingsAlertRef = ref(null)

const handleOpenRecallSettingsAlert = () => {
  recallSettingsAlertRef.value.open(toRaw(formState))
}

const onChangeRecallSettings = (data) => {
  formState.rerank_status = data.rerank_status
  formState.rerank_use_model = data.rerank_use_model
  formState.rerank_model_config_id = data.rerank_model_config_id
  formState.top_k = data.top_k
  formState.similarity = data.similarity
  formState.search_type = data.search_type
}

// 获取知识库
const getList = async () => {
  const res = await getLibraryList({ type: '' })
  if (res) {
    libraryList.value = res.data || []
  }
}

getList()

defineExpose({})
</script>

<style lang="less" scoped>
@import '../form-block.less';

.options-item {
  margin-top: 12px;
  height: 22px;
  line-height: 22px;
  display: flex;
  align-items: center;
  gap: 8px;
  .option-label {
    color: var(--wf-color-text-1);
    font-size: 14px;
    &::before {
      content: '*';
      color: #fb363f;
      display: inline-block;
      margin-right: 2px;
    }
  }
  .option-type {
    height: 22px;
    width: fit-content;
    padding: 0 8px;
    border-radius: 6px;
    border: 1px solid rgba(0, 0, 0, 0.15);
    background-color: #fff;
    color: var(--wf-color-text-3);
    font-size: 12px;
    display: flex;
    align-items: center;
  }
}

.konwledge-list-item {
  padding: 14px 16px;
  height: 72px;
  border: 1px solid #d8dde5;
  border-radius: 6px;
  display: flex;
  align-items: center;
  background-color: #fff;
  font-size: 14px;
  line-height: 22px;
  margin-top: 4px;
  .avatar-box {
    width: 40px;
    height: 40px;
    border: 1px solid var(--07, #f0f0f0);
    border-radius: 6px;
    margin-right: 8px;
    img {
      width: 100%;
      height: 100%;
    }
  }
  .list-name {
    color: #262626;
    font-weight: 600;
  }
  .list-intro {
    font-size: 12px;
    margin-top: 2px;
    line-height: 20px;
    color: #8c8c8c;
    max-width: 400px;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
  .content-box {
    flex: 1;
  }
}

.btn-block {
  margin-top: 4px;
  display: flex;
  gap: 4px;
  & > div {
    flex: 1;
  }
}
</style>
