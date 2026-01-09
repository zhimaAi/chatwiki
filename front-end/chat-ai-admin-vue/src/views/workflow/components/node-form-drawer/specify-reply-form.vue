<style lang="less" scoped>
.specify-reply-form {
  position: relative;

  .node-content {
    background: #f2f4f7;
    border-radius: 6px;
    padding: 12px;
    margin-top: 16px;
  }

  .form-label {
    margin-bottom: 12px;
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    color: #262626;
    line-height: 22px;
    font-weight: 600;
  }

  .form-item {
    .form-item-label {
      line-height: 22px;
      margin-bottom: 2px;
      font-size: 14px;
      color: #262626;
    }
  }
  .flex-between-box {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .btn-hover-wrap {
    width: 32px;
    height: 32px;
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
}
</style>

<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        desc="触发该节点，将生成一条固定消息"
      ></NodeFormHeader>
    </template>
    <div class="specify-reply-form">
      <div class="node-content">
          <div class="form-label">输出</div>
          <div class="form-item" @mousedown.stop="">
            <div class="form-item-label">
              <div class="flex-between-box">
                <div>消息内容</div>
                <div class="btn-hover-wrap" @click="handleOpenFullAtModal">
                  <FullscreenOutlined />
                </div>
              </div>
            </div>
            <div class="form-item-body">
              <AtInput
                inputStyle="height: 100px;"
                :options="valueOptions"
                :defaultSelectedList="formState.content_tags"
                :defaultValue="formState.content"
                ref="atInputRef"
                placeholder="请输入消息内容，键入“/”可以插入变量"
                input-style="height: 130px"
                type="textarea"
                @open="showAtList"
                @change="(text, selectedList) => changeValue(text, selectedList)"
              />
            </div>
          </div>
        </div>

        <FullAtInput
          :options="valueOptions"
          :defaultSelectedList="formState.content_tags"
          :defaultValue="formState.content"
          placeholder="请输入消息内容，键入“/”可以插入变量"
          type="textarea"
          @open="showAtList"
          @change="(text, selectedList) => changeValue(text, selectedList)"
          @ok="handleRefreshAtInput"
          ref="fullAtInputRef"
        />
    </div>
  </NodeFormLayout>
 
</template>

<script>
import NodeFormLayout from './node-form-layout.vue'
import NodeFormHeader from './node-form-header.vue'
import { haveOutKeyNode } from '@/views/workflow/components/util.js'
import { FullscreenOutlined } from '@ant-design/icons-vue'
import AtInput from '../at-input/at-input.vue'
import FullAtInput from '../at-input/full-at-input.vue'

export default {
  name: 'SpecifyReplyNode',
  inject: ['getNode', 'getGraph', 'setData'],
  components: {
    NodeFormHeader,
    AtInput,
    FullscreenOutlined,
    FullAtInput,
    NodeFormLayout
  },
  props: {
    node: {
      type: Object,
      default() {
        return {}
      }
    },
    isSelected: { type: Boolean, default: false },
    isHovered: { type: Boolean, default: false }
  },
  data() {
    return {
      valueOptions: [],
      formState: {
        content: '',
        content_tags: []
      }
    }
  },
  computed: {},
  watch:  {
    formState(){
      this.$nextTick(() => {
        this.update()
      })
    }
  },
  mounted() {
    this.getValueOptions()

    const graphModel = this.getGraph()

    let dataRaw = this.node.dataRaw || this.node.node_params || '{}'
    let node_params = JSON.parse(dataRaw)

    const reply = node_params.reply || {
      content: '',
      content_tags: []
    }

    this.formState.content = reply.content

    this.formState.content_tags = reply.content_tags || []

    this.$nextTick(() => {
      this.$refs[`atInputRef`].refresh()
    })

    this.update()

    graphModel.eventCenter.on('custom:setNodeName', this.onUpatateNodeName)
  },
  onBeforeUnmount() {
    const graphModel = this.getGraph()
    graphModel.eventCenter.off('custom:setNodeName', this.onUpatateNodeName)
  },
  methods: {
    onUpatateNodeName (data) {
      if (!this.$refs[`atInputRef`]) return
      if(!haveOutKeyNode.includes(data.node_type)){
        return
      }

      this.getValueOptions()

      this.$nextTick(() => {
        if (this.formState.content_tags && this.formState.content_tags.length > 0) {
          this.formState.content_tags.forEach((tag) => {
            if (tag.node_id == data.node_id) {
              let arr = tag.label.split('/')
              arr[0] = data.node_name
              tag.label = arr.join('/')
              tag.node_name = data.node_name
            }
          })
        }

        this.$refs[`atInputRef`].refresh()
      })
    },
    getValueOptions() {
      let options = this.getNode().getAllParentVariable()

      this.valueOptions = options || []
    },
    update() {
      let node_params = JSON.parse(this.node.node_params)

      node_params.reply = { ...this.formState }

      this.setData({
        ...this.node,
        node_params: JSON.stringify(node_params)
      })
    },
    showAtList(val) {
      if (val) {
        this.getValueOptions()
      }
    },
    changeValue(text, selectedList) {
      this.formState.content_tags = selectedList
      this.formState.content = text

      this.update()
    },
    handleOpenFullAtModal() {
      this.$refs.fullAtInputRef.show()
    },
    handleRefreshAtInput(){
      this.$refs[`atInputRef`].refresh()
    },
  }
}
</script>
