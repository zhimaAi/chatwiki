<style lang="less" scoped>
.add-model-alert {
  .form-wrapper {
    display: flex;
    flex-direction: column;
    margin-top: 14px;
  }

  .select {
    padding-bottom: 2px;
  }

  .model-logo {
    display: block;
    height: 26px;
  }

  .list-box {
    width: 100%;
    max-height: 400px;
    overflow: auto;
  }

  .list-item {
    width: 100%;
    border-bottom: 1px solid #f0f0f0;
    line-height: 36px;
  }

  .tools-wrapper {
    display: flex;
    align-items: center;
    justify-content: flex-start;
    gap: 8px;
    padding: 8px;
  }
}

/* 滚动条样式 */
.list-box::-webkit-scrollbar {
    width: 4px; /*  设置纵轴（y轴）轴滚动条 */
    height: 4px; /*  设置横轴（x轴）轴滚动条 */
}
/* 滚动条滑块（里面小方块） */
.list-box::-webkit-scrollbar-thumb {
    border-radius: 0px;
    background: transparent;
}
/* 滚动条轨道 */
.list-box::-webkit-scrollbar-track {
    border-radius: 0;
    background: transparent;
}

/* hover时显色 */
.list-box:hover::-webkit-scrollbar-thumb {
    background: rgba(0,0,0,0.2);
}
.list-box:hover::-webkit-scrollbar-track {
    background: rgba(0,0,0,0.1);
}
</style>

<template>
  <a-modal class="add-model-alert" width="800px" v-model:open="show" :title="currentTitle" @ok="handleOk" @cancel="handleClose">
    <div class="form-wrapper" v-if="activeKey === 'robot'">
        <div class="select">
          {{ t('text_selected', { count: checkedList.length }) }}
        </div>
        <a-checkbox-group class="list-box" v-model:value="checkedList">
          <a-checkbox class="list-item" v-for="item in dataList" :key="item.id" :value="item.id">
            <div class="tools-wrapper">
              <img class="model-logo" :src="item.robot_avatar" alt="" />
              <div class="item-name">{{ item.robot_name }}</div>
            </div>
          </a-checkbox>
        </a-checkbox-group>
    </div>
    <div class="form-wrapper" v-else-if="activeKey === 'library'">
        <div class="select">
          {{ t('text_selected', { count: checkedList.length }) }}
        </div>
        <a-checkbox-group class="list-box" v-model:value="checkedList">
          <a-checkbox class="list-item" v-for="item in dataList" :key="item.id" :value="item.id">
            <div class="tools-wrapper">
              <div class="item-name">{{ item.library_name }}</div>
            </div>
          </a-checkbox>
        </a-checkbox-group>
    </div>
    <div class="form-wrapper" v-else-if="activeKey === 'form'">
        <div class="select">
          {{ t('text_selected', { count: checkedList.length }) }}
        </div>
        <a-checkbox-group class="list-box" v-model:value="checkedList">
          <a-checkbox class="list-item" v-for="item in dataList" :key="item.id" :value="item.id">
            <div class="tools-wrapper">
              <div class="item-name">{{ item.name }}</div>
            </div>
          </a-checkbox>
        </a-checkbox-group>
    </div>
    <div class="empty-box" v-if="!dataList.length">
      <a-empty :image="simpleImage" />
    </div>
  </a-modal>
</template>
<script setup>
import { ref, reactive, markRaw, toRaw, watch } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import { Empty } from 'ant-design-vue'

const { t } = useI18n('components.see-model-alert.see-model-alert')

const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE
const emit = defineEmits(['save'])
const props = defineProps({
  robotList: {
    type: Array,
    default: null
  },
  libraryList: {
    type: Array,
    default: null
  },
  formList: {
    type: Array,
    default: null
  },
  currentTitle: {
    type: String,
    default: '关联机器人'
  }
})
const checkedList = ref([])
const activeKey = ref('')
const show = ref(false)
const dataList = ref([])

const formatId = (arr) => {
  const newArr = []
  arr.map((item) => {
    newArr.push(item.id)
  })
  return newArr
}

const open = (key, status, record) => {
  checkedList.value = []
  activeKey.value = key
  if (activeKey.value === 'robot') {
    dataList.value = props.robotList

    // 编辑
    if (status && status === 'edit') {
      if (record) {
        checkedList.value = formatId(record.managed_robot_list)
      }
    }
  } else if (activeKey.value === 'library') {
    dataList.value = props.libraryList

    // 编辑
    if (status && status === 'edit') {
      if (record) {
        checkedList.value = formatId(record.managed_library_list)
      }
    }
  } else if (activeKey.value === 'form') {
    dataList.value = props.formList

    // 编辑
    if (status && status === 'edit') {
      if (record) {
        checkedList.value = formatId(record.managed_form_list)
      }
    }
  }
  
  show.value = true
}

const handleOk = () => {
  show.value = false
  if (!dataList.value.length) return
  emit('save', checkedList.value)
}

const handleClose = () => {
 
}

// watch(
//   () => checkedList.value,
//   (val) => {
//   },
// );

defineExpose({
  open
})
</script>
