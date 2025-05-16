<style lang="less" scoped>
.popup-wrapper {
  display: flex;
  flex-flow: column nowrap;
  width: 256px;
  height: 100%;
  border-radius: 6px;
  background-color: #fff;

  .popup-header {
    display: flex;
    flex-flow: row nowrap;
    justify-content: space-between;
    align-items: center;
    padding: 16px 16px;

    .popup-title {
      line-height: 24px;
      font-size: 16px;
      font-weight: 600;
      color: #000000;
    }

    .close-btn {
      font-size: 16px;
      color: #595959;
      cursor: pointer;
    }
  }

  .popup-body {
    flex: 1;
    padding-bottom: 16px;
    overflow: hidden;

    .popup-content {
      padding: 0 16px;
    }
  }
}

.entity-details {
  .to-details {
    display: flex;
    padding: 0 8px;
    height: 24px;
    line-height: 24px;
    font-size: 14px;
    font-weight: 400;
    border-radius: 6px;
    color: #595959;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      background: #E4E6EB;
    }

    .right-icon {
      font-size: 16px;
      color: #595959;
    }
  }

  .entity-item {
    margin-bottom: 16px;
  }

  .entity-item-header {
    display: flex;
    flex-flow: row nowrap;
    justify-content: space-between;
    margin-bottom: 4px;

    .entity-item-label {
      line-height: 22px;
      font-size: 14px;
      font-weight: 600;
      color: #262626;
    }
  }

  .entity-name,
  .fragment-content {
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: #595959;
  }

  .doc-item {
    margin-bottom: 4px;

    &:last-child {
      margin-bottom: 0;
    }
  }

  .subgraph-item {
    padding-bottom: 8px;
    margin-bottom: 8px;
    border-bottom: 1px solid #f0f0f0;

    &:last-child {
      margin-bottom: 0;
      border-bottom: 0;
    }
  }

  .file-info {
    display: flex;
    align-items: center;
    border-radius: 6px;
    color: #2475FC;
    transition: all 0.2s;
    cursor: pointer;

    .file-icon {
      font-size: 16px;
    }

    .file-name {
      font-size: 14px;
      font-weight: 400;
      margin-left: 4px;
    }
  }
}
</style>

<template>
  <div class="popup-wrapper entity-details-wrapper" v-if="state.show">
    <div class="popup-header">
      <div class="popup-title">实体详情</div>
      <CloseOutlined class="close-btn" @click="close" />
    </div>

    <div class="popup-body" v-if="state.node">
      <cu-scroll>
        <div class="popup-content">
          <div class="entity-details">
            <div class="entity-item">
              <div class="entity-item-header">
                <div class="entity-item-label">实体</div>
              </div>
              <div class="entity-name">{{ state.node.name }}</div>
            </div>


            <div class="entity-item" style="margin-bottom: 8px;">
              <div class="entity-item-header">
                <div class="entity-item-label">归属文档</div>
              </div>

              <div class="file-info">
                <svg-icon class="file-icon" name="doc-file" style=""></svg-icon>
                <span class="file-name">{{ state.node.file_name }}</span>
              </div>
            </div>


            <div class="entity-item">
              <div class="entity-item-header">
                <div class="entity-item-label">分段</div>
              </div>

              <div>
                <div class="doc-item">
                  <div class="fragment-content">{{ state.node.content }}</div>
                </div>
              </div>
            </div>


            <div class="entity-item">
              <div class="entity-item-header">
                <div class="entity-item-label">子图</div>
              </div>

              <div>
                <div class="" v-if="state.fromNodes.length == 0">无</div>
                <div class="subgraph-item" v-for="item in state.fromNodes" :key="item.id">
                  <div class="fragment-content">{{ item.name }}</div>
                </div>
              </div>
            </div>

          </div>
        </div>

      </cu-scroll>
    </div>
  </div>
</template>

<script setup>
import { reactive, nextTick } from 'vue'
import { CloseOutlined } from '@ant-design/icons-vue'
import CuScroll from '@/views/library/knowledge-graph/components/cu-scroll.vue';

const state = reactive({
  show: false,
  node: null,
  fromNodes: []
})

const close = () => {
  state.show = false
}

const open = (node, fromNodes) => {
  state.node = node
  state.fromNodes = fromNodes

  if(state.show){
    state.show = false
    nextTick(() => {
      state.show = true
    })
  }else{
    state.show = true
  }
}

defineExpose({
  open,
  close
})
</script>
