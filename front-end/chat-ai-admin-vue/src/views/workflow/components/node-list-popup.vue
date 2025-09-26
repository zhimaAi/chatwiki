<style lang="less" scoped>
.node-list-popup {
  .node-list-popup-content {
    width: 302px;
    max-height: 600px;
    overflow-y: auto;
    overflow-x: hidden;
    padding: 2px;
    border-radius: 6px;
    background: #fff;
    box-shadow:
      0 6px 30px 5px #0000000d,
      0 16px 24px 2px #0000000a,
      0 8px 10px -5px #00000014;
  }

  .node-list {
    .node-type {
      height: 30px;
      padding-left: 16px;
      font-size: 12px;
      color: #8c8c8c;
      display: flex;
      align-items: center;
    }
    .node-flex-box {
      display: flex;
      align-items: center;
      flex-wrap: wrap;
      gap: 4px 0;
      .node-item {
        height: 100%;
        height: 32px;
        cursor: pointer;
        width: 50%;
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 0 16px;
        &:hover {
          background: #e4e6eb;
          border-radius: 6px;
        }
        img {
          width: 20px;
          height: 20px;
        }
      }
    }
  }
}
</style>

<template>
  <div class="node-list-popup">
    <div class="node-list-popup-content">
      <template v-for="group in allGroupNodes"  :key="group.key">
        <div class="node-list" v-if="!group.hidden">
          <div class="node-type">{{ group.name }}</div>
          <div class="node-flex-box">
            <template v-for="node in group.nodes" :key="node.type">
              <div
                class="node-item"
                @click="handleAddNode(node)"
                v-if="!node.hidden"
              >
                <img :src="node.properties.node_icon" :alt="node.properties.node_name" />
                <div class="node-name">{{ node.properties.node_name }}</div>
              </div>
            </template>
            
          </div>
        </div>
      </template>
      
    </div>
  </div>
</template>

<script setup>
import { generateUniqueId } from '@/utils/index'
import { getAllGroupNodes } from './node-list'

const emit = defineEmits(['addNode'])

const props = defineProps({
  type: {
    type: String,
    default: ''
  }
})

const allGroupNodes = getAllGroupNodes(props.type)

let nodeNameMap = {}
const handleAddNode = (node) => {
  node = JSON.parse(JSON.stringify(node))
  node.id = generateUniqueId(node.type)
  // 同一类型的节点多次添加时，从第二次添加开始，默认名称后面加上序号
  if(nodeNameMap[node.type]){
    nodeNameMap[node.type] = nodeNameMap[node.type] + 1
    node.properties.node_name = node.properties.node_name + nodeNameMap[node.type]
  }else{
    nodeNameMap[node.type] = 1
  }
  emit('addNode', node)
}
</script>
