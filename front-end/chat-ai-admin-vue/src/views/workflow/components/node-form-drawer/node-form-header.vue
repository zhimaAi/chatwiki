<style lang="less" scoped>
.node-form-header {
  margin-bottom: 16px;
  
  .node-title {
    display: flex;
    align-items: center;
    & > img{
      width: 20px;
      height: 20px;
    }
  }
  .title-left {
    margin-right: 8px;
  }
  .title-conttent {
    flex: 1;
  }
  .title-right {
    display: flex;
    align-items: center;

    .action-btn {
      display: flex;
      justify-content: center;
      align-items: center;
      width: 24px;
      height: 24px;
      border-radius: 6px;
      cursor: pointer;
      transition: background-color 0.2s;
      margin-left: 8px;

      &:hover {
        background-color: #e4e6eb;
      }

      .btn-icon {
        display: block;
        width: 16px;
        height: 16px;
        font-size: 16px;
        line-height: 16px;
        text-align: center;
        color: #333;
      }
    }
  }
  .node-icon {
    display: block;
    width: 20px;
    height: 20px;
    border-radius: 6px;
  }
  .node-name-input {
    line-height: 24px;
    padding: 0 8px;
    font-size: 16px;
    font-weight: 600;
    border-radius: 6px;
    color: #262626;
    &::placeholder {
      font-weight: 400;
      font-size: 14px;
    }
  }

  .node-desc {
    margin-top: 8px;
    font-size: 14px;
    font-weight: 400;
    line-height: 22px;
    color: var(--wf-color-text-2);
  }
}
</style>

<template>
  <div class="node-form-header">
    <div class="node-title">
      <div class="title-left">
        <slot name="node-icon">
          <img class="node-icon" :src="nodeIconUrl" v-if="nodeIconUrl">
        </slot>
      </div>
      <div class="title-conttent">
        <a-input class="node-name-input" v-model:value="localTitle" placeholder="请输入节点名称" @change.stop="handleTitleChange" />
      </div>
      <div class="title-right">
        <slot name="runBtn"></slot>
        <a-dropdown v-if="!hideMenus">
          <div class="action-btn delete-btn" @click.prevent>
            <svg-icon name="point-h" class="btn-icon"></svg-icon>
          </div>
          <template #overlay>
            <a-menu style="width: 100px">
              <a-menu-item @click="handleDeleteNode()">
                <a href="javascript:;">删除</a>
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>

        <div class="action-btn close-btn" @click="handleClose">
          <CloseOutlined class="btn-icon" style="font-size: 12px" />
        </div>
      </div>
    </div>

    <div class="node-desc"><slot name="desc">{{ props.desc }}</slot></div>
  </div>
</template>

<script setup>
import { ref,  computed, inject } from 'vue'
import { CloseOutlined } from '@ant-design/icons-vue'

const changeTitle = inject('changeTitle')
const deleteNode = inject('deleteNode')
const close = inject('close')

const props = defineProps({
  iconName: {
    type: String,
    default: ''
  },
  title: {
    type: String,
    default: ''
  },
  desc: {
    type: String,
    default: ''
  },
  hideMenus: {
    type: Boolean,
    default: false
  }
})

let chagneTimer = null;

const getImggeUrl = (name) => {
  let url = new URL(`../../../../assets/svg/${name}.svg`, import.meta.url)
  return url.href
}

const nodeIconUrl = computed(() => getImggeUrl(props.iconName))

const localTitle = ref(props.title)

const handleTitleChange = () => {
  if (chagneTimer) {
    clearTimeout(chagneTimer)
  }
  chagneTimer = setTimeout(() => {
    changeTitle(localTitle.value)
  }, 250)
}

const handleDeleteNode = () => {
  deleteNode()
}

const handleClose = () => {
  close()
}
</script>
