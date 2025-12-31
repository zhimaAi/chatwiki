<style lang="less" scoped>
.node-list-popup {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
  overflow: hidden;
  border-radius: 6px 6px 0 0;
  box-shadow: 0 4px 16px 0 #0000001a;
  .tabs-box {
    display: flex;
    align-items: center;
    border-radius: 6px;
    background: #f3f3f3;

    .tab-item {
      display: flex;
      align-items: center;
      gap: 4px;
      color: #595959;
      font-size: 14px;
      font-weight: 400;
      padding: 12px 18px;
      cursor: pointer;

      &.active {
        color: #2475fc;
        font-weight: 600;
        border-radius: 9px 9px 0 0;
        background: #fff;
      }
    }
  }

  .node-list-popup-content {
    flex: 1;
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

  .node-box {
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 8px;
    overflow: hidden;

    .search-box {
      margin: 0 8px 8px;
    }

    .node-list {
      flex: 1;
      overflow-y: auto;
      .node-info {
        display: flex;
        align-items: center;
        padding: 4px 8px;
        border-radius: 6px;
        cursor: pointer;
        &:hover {
          background: #e4e6eb;
        }

        .avatar {
          width: 20px;
          height: 20px;
          flex-shrink: 0;
          border-radius: 4.71px;
          margin-right: 8px;
        }

        .info {
          flex: 1;
          display: flex;
          align-items: center;
          justify-content: space-between;

          .name {
            color: #262626;
            font-size: 14px;
            font-weight: 400;
          }

          .total {
            color: #8c8c8c;
            font-size: 12px;
            font-weight: 400;
          }
        }
      }

      .node-tools {
        margin-left: 28px;
        border-left: 1px solid #d9d9d9;
      }
      .node-tool-item {
        padding: 2px 8px;
        border-radius: 6px;
        cursor: pointer;
        &:hover {
          background: #e4e6eb;
        }
      }
    }

    .more-link {
      padding: 8px 16px 0;
      border-top: 1px solid #F0F0F0;
      color: #2475FC;
    }
  }
}

.node-info-pop {
  display: flex;
  flex-direction: column;
  gap: 12px;
  color: #595959;
  font-size: 14px;
  font-weight: 400;
  min-width: 260px;

  .info {
    display: flex;
    align-items: center;
    color: #262626;
    font-size: 16px;

    .avatar {
      width: 62px;
      height: 62px;
      border-radius: 14.59px;
      flex-shrink: 0;
      margin-right: 12px;
    }
  }

  .extra {
    color: #8c8c8c;
    font-size: 12px;
  }
}
.params-box {
  max-height: 80vh;
  overflow-y: auto;

  .param-item {
    max-width: 500px;

    &:first-child {
      border-bottom: 1px solid #e4e6eb;
      padding-bottom: 16px;
    }
    &:not(:last-child) {
      margin-bottom: 16px;
    }

    .field {
      color: #262626;
      font-size: 14px;
      display: flex;
      align-items: center;
      gap: 12px;

      .name {
        font-weight: 600;
      }

      .type {
        color: #595959;
      }

      .required {
        color: #ed744a;
        font-weight: 500;
      }
    }

    .desc {
      color: #8c8c8c;
      font-size: 14px;
      margin-top: 4px;
    }
  }
}

.empty-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin: 40px 0;

  img {
    width: 200px;
    height: 200px;
  }

  .title {
    color: #262626;
    font-size: 16px;
    font-weight: 600;
    line-height: 24px;
  }

  .btn {
    margin-top: 24px;
  }
}
</style>

<template>
  <div class="node-list-popup">
    <div class="tabs-box">
      <div :class="['tab-item', { active: tabActive == 1 }]" @click="tabChange(1)">
        <svg-icon name="base-node-type" /> <span>基础功能</span>
      </div>
      <div :class="['tab-item', {active: tabActive == 3}]" @click="tabChange(3)">
        <svg-icon name="plugin-node-type"/> <span>插件</span>
      </div>
      <div :class="['tab-item', { active: tabActive == 2 }]" @click="tabChange(2)">
        <svg-icon name="mcp-node-type" /> <span>MCP</span>
      </div>
      <div v-if="props.type == 'float-btn'" :class="['tab-item', { active: props.active == 4 }]" @click="tabChange(4)">
        <svg-icon name="mcp-node-type" /> <span>触发器</span>
      </div>
    </div>
    <div class="node-list-popup-content">
      <template v-if="tabActive == 1">
        <template v-for="group in allGroupNodes" :key="group.key">
          <div class="node-list" v-if="!group.hidden">
            <div class="node-type">{{ group.name }}</div>
            <div class="node-flex-box">
              <template v-for="node in group.nodes" :key="node.type">
                <div
                  class="node-item"
                  @mousedown="handleMouseDownOnNode($event, node)"
                  v-if="!node.hidden"
                >
                  <img :src="node.properties.node_icon" :alt="node.properties.node_name" />
                  <div class="node-name">{{ node.properties.node_name }}</div>
                </div>
              </template>
            </div>
          </div>
        </template>
      </template>
      <div v-else-if="tabActive == 2" class="node-box">
        <div class="search-box">
          <a-input-search v-model:value.trim="mcpKeyword" allowClear placeholder="请输入名称查询" />
        </div>
        <div v-if="showMcpNodes.length" class="node-list">
          <div v-for="(item, j) in showMcpNodes" :key="j" class="node-item">
            <a-popover placement="right">
              <template #content>
                <div class="node-info-pop">
                  <div class="info">
                    <img class="avatar" :src="item.avatar" />
                    <div class="name">{{ item.name }}</div>
                  </div>
                  <div>{{ item.url }}</div>
                  <div class="extra">可用工具：{{ item.tools.length }}</div>
                </div>
              </template>
              <div class="node-info" @click="item.expand = !item.expand">
                <img class="avatar" :src="item.avatar" />
                <div class="info">
                  <span class="name">{{ item.name }}</span>
                  <span class="total">
                    {{ item.tools.length }} <DownOutlined v-if="item.expand"/> <RightOutlined v-else/>
                  </span>
                </div>
              </div>
            </a-popover>
            <div v-show="item.expand" class="node-tools">
              <a-popover v-for="(tool, i) in item.tools" :key="i" placement="right">
                <template #content>
                  <div class="params-box">
                    <div class="param-item">
                      <div class="field">
                        <span class="name">{{ tool.name }}</span>
                      </div>
                      <div class="desc">{{ tool.description }}</div>
                    </div>
                    <div
                      v-for="(field, key) in tool.inputSchema.properties"
                      :key="key"
                      class="param-item"
                    >
                      <div class="field">
                        <span class="name">{{ key }}</span>
                        <span class="type">{{ field.type }}</span>
                        <span v-if="tool.inputSchema.required && tool.inputSchema.required.includes(key)" class="required"
                          >必填</span
                        >
                      </div>
                      <div class="desc">{{ field.description }}</div>
                    </div>
                  </div>
                </template>
                <div class="node-tool-item" @mousedown="addMcpNode($event, item, tool)">{{ tool.name }}</div>
              </a-popover>
            </div>
          </div>
        </div>
        <div v-else class="empty-box">
          <img style="height: 200px" src="@/assets/empty.png" />
          <div>暂无可用MCP插件</div>
          <a @click="handleOpenAddMcp">去添加<RightOutlined /></a>
        </div>
      </div>
      <div v-else-if="tabActive == 3" class="node-box">
        <template v-if="allPluginNodes.length">
          <div class="node-list">
            <template
              v-for="node in allPluginNodes"
              :key="node.type">
              <div v-if="!isActionsPlugin(node)"
                   @click="handleAddNode(node)"
                   class="node-item"
              >
                <div class="node-info">
                  <img class="avatar" :src="node.properties.node_icon"/>
                  <div class="info"><span class="name">{{ node.properties.node_name }}</span></div>
                </div>
              </div>
              <template v-else>
                <div v-if="getPluginActions(node.plugin_name).length > 1" class="node-item">
                  <a-popover placement="right">
                    <template #content>
                      <div class="node-info-pop">
                        <div class="info">
                          <img class="avatar" :src="node.properties.node_icon" />
                          <div class="name">{{ node.properties.node_name }}</div>
                        </div>
                        <div>{{ node.properties.node_desc }}</div>
                        <div class="extra">可用工具：{{ getPluginActions(node.plugin_name).length  }}</div>
                      </div>
                    </template>
                    <div class="node-info" @click="node.expand = !node.expand">
                      <img class="avatar" :src="node.properties.node_icon" />
                      <div class="info">
                        <span class="name">{{ node.properties.node_name }}</span>
                        <span class="total">
                        {{ getPluginActions(node.plugin_name).length }} <DownOutlined v-if="node.expand"/> <RightOutlined v-else/>
                      </span>
                      </div>
                    </div>
                  </a-popover>
                  <div v-show="node.expand" class="node-tools">
                    <a-popover v-for="action in getPluginActions(node.plugin_name)" :key="action.name" placement="right">
                      <template #content>
                        <div class="params-box">
                          <div class="param-item">
                            <div class="field">
                              <span class="name">{{ action.title }}</span>
                            </div>
                            <div class="desc">{{ action.desc }}</div>
                          </div>
                          <div
                            v-for="(field, key) in action.params"
                            :key="key"
                            class="param-item"
                          >
                            <div class="field">
                              <span class="name">{{ key }}</span>
                              <span class="type">{{ field.type }}</span>
                              <span v-if="field.required" class="required">必填</span>
                            </div>
                            <div class="desc">{{ field.desc }}</div>
                          </div>
                        </div>
                      </template>
                      <div class="node-tool-item" @mousedown="addPluginNode($event, node, action, action.name)">{{ action.title }}</div>
                    </a-popover>
                  </div>
                </div>
                <!--仅存在一个方法时-->
                <template v-else>
                  <div v-for="action in getPluginActions(node.plugin_name)"
                       @click="addPluginNode(null, node, action, action.name)"
                       class="node-item"
                  >
                    <div class="node-info">
                      <img class="avatar" :src="node.properties.node_icon"/>
                      <div class="info"><span class="name">{{ node.properties.node_name }}</span></div>
                    </div>
                  </div>
                </template>
              </template>
            </template>
          </div>
          <a class="more-link" href="/#/plugins/index?active=2" target="_blank">更多插件 <RightOutlined/></a>
        </template>
        <div v-else class="empty-box">
          <img style="height: 200px;" src="@/assets/empty.png"/>
          <div>暂无可用插件</div>
          <a href="/#/plugins/index?active=2" target="_blank">去添加<RightOutlined/></a>
        </div>
      </div>
      <div  v-else-if="tabActive == 4" class="node-box">
        <TriggerList @add="handleAddTrigger" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { RightOutlined, DownOutlined } from '@ant-design/icons-vue'
import TriggerList from './trigger-list.vue'
import {
  getAllGroupNodes,
  getAllPluginNodes,
  getAllMcpNodes,
  getMcpNode,
  getPluginActionNode,
  getPluginActions,
} from '../node-list'
import {jsonDecode} from "@/utils/index.js";
import {pluginHasAction} from "@/constants/plugin.js";

const emit = defineEmits(['addNode', 'addTrigger', 'mouseMove', 'update:active'])

const props = defineProps({
  type: {
    type: String,
    default: ''
  },
  active: {
    type: Number,
    default: 1
  },
  excludedNodeTypes: {
    type: Array,
    default: () => []
  },
})

const tabActive = ref(props.active)

watch(() => props.active, (newVal) => {
  tabActive.value = newVal
})

const allGroupNodes = ref([])
const allPluginNodes = ref([])
const allMcpNodes = ref([])
const mcpKeyword = ref('')

watch(() => props.excludedNodeTypes, () => {
  allGroupNodes.value = getAllGroupNodes({excludedNodeTypes: props.excludedNodeTypes})
}, {
  immediate: true
})


onMounted(() => {
  getAllMcpNodes().then((res) => {
    allMcpNodes.value = res
  })
  getAllPluginNodes().then(res => {
    allPluginNodes.value = res
  })
})

const showMcpNodes = computed(() => {
  if (mcpKeyword.value) {
    return allMcpNodes.value.filter((item) => {
      let info = item.description + item.name
      return info.indexOf(mcpKeyword.value) > -1
    })
  } else {
    return allMcpNodes.value
  }
})

const handleMouseDownOnNode = (event, node) => {
  // 阻止默认行为，例如文本选择
  event.preventDefault()

  let isDragged = false
  const startX = event.clientX
  const startY = event.clientY

  // 复制节点数据，为之后创建节点做准备
  const nodeData = JSON.parse(JSON.stringify(node))

  // 创建一个临时的、跟随鼠标的预览元素
  const ghost = document.createElement('div')
  ghost.classList.add('node-item-ghost')
  ghost.innerHTML = event.currentTarget.innerHTML // 复制外观

  document.body.appendChild(ghost)

  const moveGhost = (e) => {
    ghost.style.left = `${e.clientX + 5}px`
    ghost.style.top = `${e.clientY + 5}px`
  }

  moveGhost(event)

  const onMouseMove = (e) => {
    if (
      !isDragged &&
      (Math.abs(e.clientX - startX) > 5 || Math.abs(e.clientY - startY) > 5)
    ) {
      isDragged = true
    }
    moveGhost(e)
    emit('mouseMove')
  }

  const onMouseUp = (e) => {
    // 清理
    document.body.removeChild(ghost)
    window.removeEventListener('mousemove', onMouseMove)
    window.removeEventListener('mouseup', onMouseUp)

    if (isDragged) {
      // 如果拖动了，发出事件，通知画布创建节点
      // 传递节点数据和鼠标抬起时的事件对象（包含坐标）
      emit('addNode', { node: nodeData, event: e })
    } else {
      // 如果没有拖动，执行原有的 click 逻辑
      handleAddNode(nodeData)
    }
  }

  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', onMouseUp, { once: true })
}

const handleAddNode = (node) => {
  node = JSON.parse(JSON.stringify(node))

  emit('addNode', { node })
}

const addMcpNode = (event, mcp, tool) => {
  let node = getMcpNode(mcp, tool);

  handleMouseDownOnNode(event, node)
}

const handleAddTrigger = (node) => {
  emit('addNode', { node, isTriggerNode: true })
}

const addPluginNode = (event, node, action, name) => {
  node = getPluginActionNode(node, action, name)
  if (event) {
    handleMouseDownOnNode(event, node)
  } else {
    emit('addNode', { node })
  }
}

function tabChange(key) {
  tabActive.value = key
  emit('update:active', key)
}

function handleOpenAddMcp() {
  window.open('/#/robot/list?active=3&mcp=2', '_blank')
}

function isActionsPlugin(node) {
  let node_params = jsonDecode(node?.properties?.node_params, {})
  return pluginHasAction(node_params?.plugin?.name) || node_params?.plugin?.multiNode
}
</script>

<style lang="less">
.node-item-ghost {
  position: fixed;
  z-index: 10000;
  pointer-events: none;
  display: flex;
  align-items: center;
  height: 32px;
  width: 200px;
  gap: 8px;
  padding: 0 16px;
  border-radius: 6px;
  opacity: 0.7;
  background-color: #fff;
  box-shadow: 0 4px 16px 0 #0000001a;

  img {
    width: 20px;
    height: 20px;
  }
}
</style>
