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
    padding: 0 16px 16px 16px;

    .library-item {
      position: relative;
      width: 336px;
      padding: 14px 16px;
      border-radius: 2px;
      border: 1px solid #d8dde5;
      background-color: #fff;

      .library-name {
        width: 100%;
        line-height: 22px;
        font-size: 14px;
        font-weight: 600;
        color: #262626;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
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
        top: 0;
        right: 6px;
        font-size: 16px;
        color: #8c8c8c;
        cursor: pointer;
      }
    }
  }
}
</style>

<template>
  <edit-box
    class="setting-box"
    title="数据库"
    icon-name="guanlianzhishiku"
    v-model:isEdit="isEdit"
    :bodyStyle="{ padding: 0 }"
  >
    <template #tip>
      <a-tooltip placement="top">
        <template #title>
          <span
            >关联数据库表，可以在机器人对话时，搜集用户数据进行存储，或者调用存储的数据回答用户提问。</span
          >
        </template>
        <QuestionCircleOutlined />
      </a-tooltip>
    </template>
    <template #extra>
      <div class="actions-box">
        <a-flex :gap="8">
          <a-button size="small" @click="handleOpenSelectLibraryAlert">关联数据表</a-button>
        </a-flex>
      </div>
    </template>
    <div class="library-list" v-if="selectedLibraryRows.length > 0">
      <div class="library-item" v-for="item in selectedLibraryRows" :key="item.id">
        <span class="close-btn" @click="handleRemoveCheckedLibrary(item)">
          <CloseCircleOutlined />
        </span>
        <div class="library-name">{{ item.name }}</div>
        <div class="library-intro">{{ item.description }}</div>
      </div>
    </div>
    <LibrarySelectAlert ref="librarySelectAlertRef" @change="onChangeLibrarySelected" />
  </edit-box>
</template>

<script setup>
import { getFormList } from '@/api/database/index'
import { ref, reactive, inject, watchEffect, computed, toRaw } from 'vue'
import { CloseCircleOutlined, QuestionCircleOutlined } from '@ant-design/icons-vue'
import EditBox from '../edit-box.vue'
import LibrarySelectAlert from './library-select-alert.vue'
const isEdit = ref(false)

const { robotInfo, updateRobotInfo } = inject('robotInfo')

const formState = reactive({
  form_ids: [],
})

// 知识库
const libraryList = ref([])
const librarySelectAlertRef = ref(null)
const selectedLibraryRows = computed(() => {
  return libraryList.value.filter((item) => {
    return formState.form_ids.includes(item.id)
  })
})

// 移除知识库
const handleRemoveCheckedLibrary = (item) => {
  let index = formState.form_ids.indexOf(item.id)

  formState.form_ids.splice(index, 1)

  onSave()
}

const onChangeLibrarySelected = (checkedList) => {
  formState.form_ids = [...checkedList]

  onSave()
}

const handleOpenSelectLibraryAlert = () => {
  librarySelectAlertRef.value.open([...formState.form_ids])
}



const onSave = () => {
  let formData = { ...toRaw(formState) }

  formData.form_ids = formData.form_ids.join(',')

  updateRobotInfo({ ...formData })
}

// 获取知识库
const getList = async () => {
  const res = await getFormList()
  if (res) {
    libraryList.value = res.data || []
  }
}

getList()

watchEffect(() => {
  formState.form_ids = robotInfo.form_ids.split(',')
})
</script>
