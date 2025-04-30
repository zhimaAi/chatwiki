<style lang="less" scoped>
.variable-node {
  position: relative;
  z-index: 2;
  .node-desc {
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    margin-bottom: 16px;
    color: var(--wf-color-text-2);
  }
  .field-items{
    .field-item{
      display: flex;
      align-items: center;
      margin-bottom: 8px;
    }

    .field-name-box{
      width: 180px;
      margin-right: 8px;
    }
    .field-value-box{
      flex: 1;
      overflow: hidden;
    }
    .field-actions{
      margin-left: 12px;
      .action-btn{
        cursor: pointer;
        font-size: 16px;
        color: #595959;
      }
    }
  }
}

.field-value-option{
  display: flex;
  align-items: center;
  color: #595959;
  font-size: 14px;
  font-weight: 400;
  .option-label{
    font-weight: 400;
    margin-right: 6px;
  }
  .option-type{
    width: fit-content;
    padding: 1px 8px;
    border-radius: 6px;
    border: 1px solid #00000026;
    background: var(--10, #FFF);
    display: flex;
    align-items: center;
    font-size: 12px;
  }
}
</style>

<template>
  <node-common
    :title="properties.node_name"
    :menus="menus"
    :icon-name="properties.node_icon_name"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    :node_type="properties.node_type"
    @handleMenu="handleMenu"
  >
    <div class="variable-node">
      <div class="node-desc">通过本节点可以给全局变量赋值</div>
      <div class="field-items">
        <div class="field-item" v-for="(item, index) in list" :key="index"> 
          <div class="field-name-box">
            <a-select style="width: 100%;" placeholder="请输入选择变量" v-model:value="item.variable" @dropdownVisibleChange="dropdownVisibleChange" @change="update">
              <a-select-option :value="opt.value" v-for="opt in options" :key="opt.key">
                <span>{{ opt.label }}</span>
              </a-select-option>
            </a-select>
          </div>
          <div class="field-value-box">
            <at-input :ref="'atInput' + index" inputStyle="overflow-y: hidden; overflow-x: scroll; height: 32px;" :options="valueOptions" :defaultSelectedList="item.tags" :defaultValue="item.value" @open="showAtList"  @change="(text, selectedList) => changeValue(text, selectedList, item, index)" placeholder="请输入变量值，键入“/”插入变量">
              <template #option="{ label, payload }">
                <div class="field-value-option">
                  <div class="option-label">{{ label }}</div>
                  <div class="option-type">{{ payload.typ }}</div>
                </div>
              </template>
            </at-input>
          </div>
          <div class="field-actions">
            <CloseCircleOutlined class="action-btn" @click="handleDel(index)" />
          </div>
        </div>
      </div>

      <div>
        <a-button style="width: 100%;" class="add-btn" type="dashed" @click="handleAdd">
          <template #icon>
            <PlusOutlined />
          </template>
          添加赋值变量
        </a-button>
      </div>
    </div>
  </node-common>
</template>

<script>
import { PlusOutlined, CloseCircleOutlined } from '@ant-design/icons-vue'
import NodeCommon from '../base-node.vue'
import AtInput from '../at-input/at-input.vue'

export default {
  name: 'VariableAssignmentNode',
  inject: ['getNode', 'getGraph', 'setData'],
  components: {
    PlusOutlined,
    CloseCircleOutlined,
    NodeCommon,
    AtInput
  },
  props: {
    properties: {
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
      menus: [{ name: '删除', key: 'delete', color: '#fb363f' }],
      list: [],
      options: [],
      valueOptions: [],
    }
  },
  computed: {
   
  },
  mounted() {
    this.getOptions();
    this.getValueOptions()

    const graphModel = this.getGraph()
    
    let dataRaw = this.properties.dataRaw || this.properties.node_params || '{}'
    let node_params = JSON.parse(dataRaw)

    let fields  = node_params.assign || []
    
    fields.forEach(item => {
      item.tags = item.tags || []
    });

    this.list = fields

    this.updateHeight()
    graphModel.eventCenter.on('custom:setNodeName', this.onUpatateNodeName)
  },
  onBeforeUnmount() {
    const graphModel = this.getGraph()
    graphModel.eventCenter.off('custom:setNodeName', this.onUpatateNodeName)
  },
  methods: {
    onUpatateNodeName(data){
      if(data.node_type !== 'http-node'){
        return;
      }
      
      this.getValueOptions()

      this.$nextTick(() => {
        this.list.forEach((item, index) => {
          if(item.tags && item.tags.length > 0){
            item.tags.forEach(tag => {
              if(tag.node_id == data.node_id){
                let arr = tag.label.split('/')
                arr[0] = data.node_name
                tag.label = arr.join('/')
                tag.node_name = data.node_name
              }
            })
          }

          this.$refs[`atInput${index}`][0].refresh()
        })
      })
      
    },
    getOptions(){
      let globalVariable = this.getNode().getGlobalVariable()
      let diy_global = globalVariable.diy_global || []
      diy_global.forEach(item => {
        item.label = item.key
        item.value = 'global.'+item.key
      });
      this.options = diy_global || []
    },
    getValueOptions(){
      let options = this.getNode().getAllParentVariable();

      this.valueOptions = options || []
    },
    getHeight() {
      let length = this.list.length == 0 ? 1 : this.list.length;

      return 130 + 40 * length
    },
    updateHeight() {
      let height = this.getHeight()

      this.setData({
        height: height
      })
    },
    update() {
      let height = this.getHeight()
      let node_params = JSON.parse(this.properties.node_params)
      
      node_params.assign = [...this.list]

      this.setData({
        height: height,
        node_params: JSON.stringify(node_params)
      })
    },
    showAtList(val){
      if(val){
        this.getValueOptions()
      }
    },
    changeValue(text, selectedList, item){
      item.tags = selectedList
      item.value = text

      this.update();
    },
    handleAdd(){
      this.list.push({
        variable: void 0,
        value: ''
      })
      this.update()
    },
    handleDel(index){
      this.list.splice(index, 1)
      this.update()
    },
    handleMenu(item) {
      if (item.key === 'delete') {
        let node = this.getNode()
        this.getGraph().deleteNode(node.id)
      }
    },
    dropdownVisibleChange(visible){
      if(visible){
        this.getOptions()
      }
    },
  }
}
</script>
