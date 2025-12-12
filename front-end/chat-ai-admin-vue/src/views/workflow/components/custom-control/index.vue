<style lang="less" scoped>
.custom-control-warpper {
  position: fixed;
  left: 50%;
  transform: translateX(-50%);
  bottom: 24px;
  z-index: 100;

  .custom-control-body {
    position: relative;
  }
  .custom-control {
    display: flex;
    flex-flow: row nowrap;
    align-items: center;
    gap: 8px;
    padding: 4px 12px;
    border-radius: 8px;
    background-color: #fff;
    box-shadow: 0 4px 16px 0 #0000001a;

    .control-line {
      width: 1px;
      height: 24px;
      background-color: #d9d9d9;
    }
  }
  .zoom-control {
    display: flex;
    flex-flow: row nowrap;
    align-items: center;

    .action-btn {
      border-radius: 6px;
      transition: all 0.2s;
      &:hover {
        background-color: #e4e6eb;
        cursor: pointer;
      }
    }

    .zoom-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 32px;
      height: 32px;
      font-size: 16px;
      color: #595959;
    }
  }

  .node-list-fix {
    position: absolute;
    bottom: 48px;
    left: 50%;
    transform: translateX(-50%);
  }
}
</style>

<template>
  <div class="custom-control-warpper">
    <div class="custom-control-body">
      <div class="custom-control">
        <div class="control-item zoom-control">
          <div class="action-btn zoom-btn" @click="handleReduce">
            <svg-icon name="minus" size="16" />
          </div>
          <zoom-select v-model="zoom" @zoom-change="chagneZoom" @fitView="handleFitView" />
          <div class="action-btn zoom-btn" @click="handleAmplify">
            <svg-icon name="plus" size="16" />
          </div>
        </div>

        <i class="control-line"></i>

        <div class="control-item">
          <a-button type="primary" @click.stop="handleShowNodeList">
            <template #icon>
              <PlusOutlined />
            </template>
            <span>新建节点</span>
          </a-button>
        </div>

        <div class="control-item">
          <a-button @click="handleRunTest" style="background-color: #00ad3a" type="primary">
            <CaretRightOutlined /><span>运行测试</span>
          </a-button>
        </div>

        <div class="control-item">
          <a-button @click="handleAutoLayout" type="default">
            <PartitionOutlined /><span>整理画布</span>
          </a-button>
        </div>
      </div>
      <!-- 使用 v-show 保证拖拽添加节点时，NodeListPopup 组件不会被销毁，确保 addNode 事件能正常触发 -->
      <div class="node-list-fix" ref="nodeListRef" v-show="isShowMenu">
        <NodeListPopup @addNode="handleAddNode" type="float-btn" @mouseMove="handleMouseMove" v-model:active="nodeListTabActive" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import ZoomSelect from './zoom-select.vue'
import { PlusOutlined, CaretRightOutlined, PartitionOutlined } from '@ant-design/icons-vue'
import NodeListPopup from '../node-list-popup/index.vue'

const emit = defineEmits(['runTest', 'addNode', 'zoomChange', 'autoLayout'])

const props = defineProps({
  lf: {
    type: Object,
    default: () => ({}),
    required: true
  }
})

const { eventCenter } = props.lf.graphModel

const zoom = ref(100)

const handleFitView = () => {
  props.lf.fitView()
}

const setZoom = () => {
  let value = zoom.value / 100
  props.lf.zoom(value)

  emit('zoomChange', value)
}

const chagneZoom = (value) => {
  zoom.value = value
  setZoom()
}

const handleReduce = () => {
  // 四舍五入到整数
  let value = zoom.value - 10

  if (value < 1) {
    value = 1
  }

  zoom.value = value
  setZoom()
}

const handleAmplify = () => {
  let value = zoom.value + 10

  if (value > 800) {
    value = 800
  }

  zoom.value = value
  setZoom()
}

const onGraphTransform = (args) => {
  let value = Math.floor(args.transform.SCALE_X * 100)
  zoom.value = value
}

const isShowMenu = ref(false)
const nodeListRef = ref(null)
const nodeListTabActive = ref(1)

const handleShowNodeList = () => {
  nodeListTabActive.value = 1
  isShowMenu.value = true
}

// let miniMap = null

// const handleMouseenter = () => {
//   props.lf.extension.miniMap.show()
//   miniMap = document.querySelector('.lf-mini-map')
//   starListenEvent()
// }

// const handleMouseleave = (e) => {
//   if (e.offsetY > 0) {
//     props.lf.extension.miniMap.hide()
//   }
// }

// function starListenEvent() {
//   if (!miniMap) {
//     return
//   }
//   miniMap.addEventListener('mouseenter', () => {
//     props.lf.extension.miniMap.show()
//   })
//   miniMap.addEventListener('mouseleave', () => {
//     props.lf.extension.miniMap.hide()
//   })
// }

const handleRunTest = () => {
  emit('runTest')
}

const handleAutoLayout = () => {
  emit('autoLayout')
}

const documentClick = (e) => {
  if (isShowMenu.value) {
    const menus = nodeListRef.value
    if (!menus.contains(e.target)) {
      isShowMenu.value = false
    }
  }
}

const handleMouseMove = () => {
  if (isShowMenu.value) {
    isShowMenu.value = false
  }
}

const handleAddNode = (node) => {
  emit('addNode', node)
  isShowMenu.value = false
}

const showTriggerLit = () => {
  nodeListTabActive.value = 4
  isShowMenu.value = true
}

onMounted(() => {
  document.addEventListener('click', documentClick)
  eventCenter.on('graph:transform', onGraphTransform)
  eventCenter.on('custom:showTriggerLit', showTriggerLit)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', documentClick)
  eventCenter.off('graph:transform', onGraphTransform)
  eventCenter.off('custom:showTriggerLit', showTriggerLit)
})
</script>
