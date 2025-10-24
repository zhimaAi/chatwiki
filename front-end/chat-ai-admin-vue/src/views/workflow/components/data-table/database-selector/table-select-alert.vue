<style lang="less" scoped>
.library-checkbox-box {
  padding-top: 16px;

  .list-tools {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 16px;
    margin-bottom: 8px;
  }

  .list-box {
    display: flex;
    flex-flow: row wrap;
    height: 388px;
    width: 100%;
    overflow-y: auto;
    align-content: flex-start;
    margin: 0 -8px;

    .list-item-wraapper {
      padding: 8px;
      width: 50%;
    }

    .list-item {
      width: 100%;
      padding: 14px 12px;
      border: 1px solid #f0f0f0;
      border-radius: 2px;

      &:hover {
        cursor: pointer;
        box-shadow: 0 4px 16px 0 #1b3a6929;
      }

      .library-name {
        line-height: 22px;
        font-size: 14px;
        font-weight: 600;
        color: #262626;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }

      .library-desc {
        line-height: 20px;
        margin-top: 2px;
        font-size: 12px;
        font-weight: 400;
        color: #8c8c8c;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }
    }

    .list-item :deep(span:last-child) {
      flex: 1;
      overflow: hidden;
    }
  }
}
.empty-box {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  padding-top: 40px;
  padding-bottom: 40px;
  color: #8c8c8c;
  img {
    width: 150px;
    height: 150px;
  }
}
</style>

<template>
  <a-modal width="746px" v-model:open="show" title="关联数据表" ok-text="确定" cancel-text="取消" @ok="saveCheckedList">
    <div class="library-checkbox-box">

      <a-radio-group v-model:value="loaclValue" style="width: 100%" v-if="options.length">
        <div class="list-box" ref="scrollContainer">
          <div class="list-item-wraapper" v-for="item in props.options" :key="item.id">
            <a-radio class="list-item" :value="item.id">
              <div class="library-name">{{ item.name }}</div>
              <div class="library-desc">{{ item.description }}</div>
            </a-radio>
          </div>
        </div>
      </a-radio-group>

      <div class="empty-box" v-if="!options.length">
        <img src="@/assets/img/library/preview/empty.png" alt="" />
        <div>
          暂无数据表, 请先去添加数据表
          <a @click="openAddLibrary"> 去添加</a>
        </div>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { ref } from 'vue'

const emit = defineEmits(['ok'])

const props = defineProps({
  value: {
    type: [String, Number],
    default: ''
  },
  options: {
    type: Array,
    default: () => []
  }
})

const loaclValue = ref()
const show = ref(false)

const open = (val) => {
  loaclValue.value = val
  
  show.value = true
}

const saveCheckedList = () => {
  show.value = false
  triggerChange()
}



const triggerChange = () => {
  let arr = props.options.filter(item => item.id === loaclValue.value)

  emit('ok', loaclValue.value, JSON.parse(JSON.stringify(arr[0])))
}

const openAddLibrary = () => {
  window.open('/#/database/list')
}

defineExpose({
  open
})
</script>
