<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        desc="支持中间向用户提问问题，支持预置选项提问和开放式问题提问两种方式"
        @close="handleClose"
      >
      </NodeFormHeader>
    </template>
    <div class="problem-optimization-form">
      <div class="node-form-content">
        <a-form ref="formRef" layout="vertical" :model="formState">
          <div class="gray-block">
            <a-flex justify="space-between">
              <div class="gray-block-title">提问内容</div>
              <div class="btn-hover-wrap" @click="handleOpenFullAtModal">
                <FullscreenOutlined />
              </div>
            </a-flex>
            <div class="array-form-box">
              <AtInput
                inputStyle="height: 100px;"
                :options="valueOptions"
                :defaultValue="formState.answer_text"
                ref="atInputRef"
                placeholder="请输入消息内容，键入“/”插入变量"
                input-style="height: 130px"
                type="textarea"
                @open="showAtList"
                @change="(text, selectedList) => changeValue(text, selectedList)"
              />
            </div>

            <a-form-item label="回答类型">
              <a-radio-group v-model:value="formState.answer_type">
                <a-radio value="text">直接回答</a-radio>
                <a-radio value="menu"
                  >智能菜单回答

                  <a-popover :title="null">
                    <template #content>
                      <div>
                        <img style="width: 330px" src="@/assets/img/robot/menu-guide.png" alt="" />
                      </div>
                    </template>
                    <QuestionCircleOutlined />
                  </a-popover>
                </a-radio>
              </a-radio-group>
            </a-form-item>
            <div class="menu-list-box" v-if="formState.answer_type == 'menu'">
              <draggable
                v-model="formState.menu_content"
                item-key="key"
                @end="onDragEnd"
                group="table-rows"
                handle=".drag-btn"
              >
                <template #item="{ element, index }">
                  <div :key="element.id" class="menu-item">
                    <div class="menu-border-box">
                      <span class="drag-btn"><svg-icon name="drag" /></span>
                      <span class="index-number">{{ index + 1 }}</span>
                      <div class="input-box">
                        <AtInput
                          inputStyle="height: 22px;"
                          :bordered="false"
                          :options="valueOptions"
                          :defaultValue="element.content"
                          placeholder="请输入消息内容，键入“/”插入变量"
                          @open="showAtList"
                          @change="
                            (text, selectedList) => changeMenuValue(text, selectedList, element)
                          "
                        />
                      </div>
                    </div>
                    <span class="action-btn">
                      <CloseCircleOutlined @click="onDelete(index)" />
                    </span>
                  </div>
                </template>
              </draggable>
              <div>
                <a-button @click="handleAddMenu()" type="dashed" block
                  ><PlusOutlined />添加菜单</a-button
                >
              </div>
            </div>
          </div>

          <div class="gray-block mt16">
            <div class="gray-block-title">输出</div>
            <div class="options-item">
              <div class="option-label">question</div>
              <div class="option-type">string</div>
            </div>
            <div class="options-item">
              <div class="option-label">question_multiple</div>
              <div class="option-type" v-text="'array<object>'"></div>
              <a-tooltip>
                <template #title>
                  <pre>
                    {{ question_multiple_tip }}
                  </pre>
                </template>
                <QuestionCircleOutlined />
              </a-tooltip>
            </div>
          </div>
        </a-form>
      </div>
      <FullAtInput
        :options="valueOptions"
        :defaultValue="formState.answer_text"
        placeholder="请输入消息内容，键入'/'可以插入变量"
        type="textarea"
        @open="showAtList"
        @change="(text, selectedList) => changeValue(text, selectedList)"
        @ok="handleRefreshAtInput"
        ref="fullAtInputRef"
      />
    </div>
  </NodeFormLayout>
</template>

<script setup>
import { getUuid } from '@/utils/index'
import NodeFormLayout from '../node-form-layout.vue'
import NodeFormHeader from '../node-form-header.vue'
import { ref, reactive, watch, h, onMounted } from 'vue'
import {
  CloseCircleOutlined,
  PlusOutlined,
  PlusCircleOutlined,
  FullscreenOutlined,
  QuestionCircleOutlined
} from '@ant-design/icons-vue'
import draggable from 'vuedraggable'
import AtInput from '../../at-input/at-input.vue'
import FullAtInput from '../../at-input/full-at-input.vue'
import { message } from 'ant-design-vue'

const emit = defineEmits(['update-node'])
const props = defineProps({
  lf: {
    type: Object,
    default: null
  },
  nodeId: {
    type: String,
    default: ''
  },
  node: {
    type: Object,
    default: () => ({})
  }
})

const atInputRef = ref(null)
const fullAtInputRef = ref(null)

const valueOptions = ref([])

const question_multiple_tip = `
用户输入的文字和图片消息示例如下：
"question_multiple": [
{
"type": "text",
"text": "这是什么",
},{
"type": "image",
"image_url": "https://url"
}]"
`

function getOptions() {
  const nodeModel = props.lf.getNodeModelById(props.nodeId)
  if (nodeModel) {
    let list = nodeModel.getAllParentVariable()

    valueOptions.value = list
  }
}

const showAtList = (val) => {
  if (val) {
    getOptions()
  }
}

const defaultMenu = ref({
  menu_type: '-1',
  serial_no: '-1',
  content: '其他未点击智能菜单的情况（用户不可见）'
})

const formState = reactive({
  answer_text: '',
  answer_type: 'text',
  menu_content: [
    {
      menu_type: '1',
      serial_no: '1',
      content: '',
      key: getUuid(16)
    },
    {
      menu_type: '1',
      serial_no: '2',
      content: '',
      key: getUuid(16)
    }
  ],

  outputs: [
    {
      key: 'question',
      typ: 'string',
      subs: []
    },
    {
      key: 'question_multiple',
      typ: 'array<object>',
      subs: []
    }
  ]
})

const update = () => {
  let menu_content = formState.menu_content.map((item, index) => {
    return {
      ...item,
      serial_no: (index + 1).toString()
    }
  })
  menu_content.push(defaultMenu.value)
  let reply_content_list = [
    {
      reply_type: 'smartMenu',
      smart_menu: {
        menu_content: menu_content
      }
    }
  ]
  const data = JSON.stringify({
    question: {
      ...formState,
      reply_content_list
    }
  })

  emit('update-node', {
    ...props.node,
    node_params: data
  })
}

const changeValue = (text, selectedList) => {
  formState.answer_text = text
}

const changeMenuValue = (text, selectedList, item) => {
  item.content = text
}

const init = () => {
  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'
    let question = JSON.parse(dataRaw).question || {}
    question = JSON.parse(JSON.stringify(question))
    formState.answer_text = question.answer_text || ''
    formState.answer_type = question.answer_type || 'text'
    let reply_content_list = question.reply_content_list || []
    if (reply_content_list.length > 0) {
      let menu_content = reply_content_list[0].smart_menu?.menu_content || []
      if (menu_content.length > 0) {
        let findItem = menu_content.find((item) => item.menu_type == '-1')
        if (findItem) {
          defaultMenu.value = findItem
        }
        formState.menu_content = menu_content
          .filter((item) => item.menu_type != '-1')
          .map((item) => ({
            ...item,
            key: getUuid(16)
          }))
      }
    }

    getOptions()
  } catch (error) {
    console.log(error)
  }
}

watch(
  () => formState,
  () => {
    update()
  },
  { deep: true }
)

const handleRefreshAtInput = () => {
  atInputRef.value.refresh()
}
const handleOpenFullAtModal = () => {
  fullAtInputRef.value.show()
}

const handleClose = () => {
  emit('close')
}

const onDragEnd = () => {}

const onDelete = (index) => {
  if (formState.menu_content.length <= 1) {
    return message.warning('最少添加1条菜单')
  }
  formState.menu_content.splice(index, 1)
}

const handleAddMenu = () => {
  formState.menu_content.push({
    menu_type: '1',
    serial_no: (formState.menu_content.length + 1).toString(),
    content: '',
    key: getUuid(16)
  })
}

onMounted(() => {
  init()
})
</script>

<style lang="less" scoped>
@import '../form-block.less';
.code-edit-box {
  margin-top: 12px;
  .title-block {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
}
.action-btn {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease-in;
  &:hover {
    background: #e4e6eb;
  }
}

.output-box {
  .output-block {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 4px;
    color: #262626;
    .output-item {
      width: 214px;
    }
  }
  .flex-block-item .btn-hover-wrap {
    width: 24px;
    height: 24px;
  }
}

.menu-list-box {
  display: flex;
  flex-direction: column;
  .menu-item {
    display: flex;
    align-items: center;
    gap: 4px;
    margin-bottom: 4px;
    .input-box {
      flex: 1;
    }
    .drag-btn {
      cursor: grab;
    }
    .index-number {
      border-radius: 6px;
      background: var(--06, #d8dde5);
      font-size: 12px;
      line-height: 16px;
      color: #242933;
      padding: 0px 6px;
      margin-left: 6px;
      height: 18px;
      line-height: 18px;
    }
    .menu-border-box {
      flex: 1;
      display: flex;
      align-items: center;
      background: #fff;
      border: 1px solid var(--06, #d9d9d9);
      border-radius: 6px;
      padding-left: 4px;
      overflow: hidden;
    }
  }
}
</style>
