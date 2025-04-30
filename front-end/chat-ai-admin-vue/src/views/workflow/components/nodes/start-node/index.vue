<style lang="less" scoped>
.start-node {
  position: relative;
  .node-desc {
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: var(--wf-color-text-2);
  }

  .start-node-options {
    background: #f2f4f7;
    border-radius: 6px;
    padding: 12px;
    margin-top: 16px;
    .options-title {
      color: var(--wf-color-text-1);
      display: flex;
      align-items: center;
      justify-content: space-between;
      font-weight: 600;
      height: 22px;;
      line-height: 22px;
      font-size: 14px;

      .title-icon{
        width: 16px;
        height: 16px;
        vertical-align: -3px;
        margin-right: 8px;;
      }
      .acton-box{
        font-weight: 400;
      }
    }
    .options-item{
      display: flex;
      align-items: center;
      margin-top: 12px;
      height: 22px;
      line-height: 22px;
     
      .options-item-body{
        flex: 1;
        display: flex;
        align-items: center;
      }
   
      .option-label{
        color: var(--wf-color-text-1);
        font-size: 14px;
        margin-right: 8px;
      }
      

      &.is-required .option-label::before{
        content: '*';
        color: #FB363F;
        display: inline-block;
        margin-right: 2px;
      }
      .option-type{
        height: 22px;
        line-height: 18px;
        padding: 0 8px;
        border-radius: 6px;
        border: 1px solid rgba(0, 0, 0, 0.15);
        background-color: #fff;
        color: var(--wf-color-text-3);
        font-size: 12px;
      }
      .item-actions-box{
        display: flex;
        align-items: center;

        .action-btn{
          margin-left: 12px;
          font-size: 16px;
          color: #595959;
          cursor: pointer;
        }
      }
    }
    .diy-options-item{
      height: 38px;
      padding: 8px;
    }
    .diy-options-item:hover{
      border-radius: 6px;
      background: #E4E6EB;
    }
  }

  .action-items-box{
    margin-top: 16px;
  }
}
</style>

<template>
  <node-common
    :title="properties.node_name"
    :icon-name="properties.node_icon_name"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    :node_type="properties.node_type"
  >
    <div class="start-node">
      <div class="node-desc">流程开始节点，点击可以设置全局变量、触发方式等</div>

      <div class="start-node-options">
        <div class="options-title">
          <div><img src="@/assets/img/workflow/output.svg" alt="" class="title-icon" />输出</div>
        </div>
        <div class="options-item" :class="{'is-required': item.required}" v-for="item in sys_global" :key="item.key">
          <div class="options-item-body">
            <div class="option-label">{{ item.key }}</div>
            <div class="option-type">{{ item.typ }}</div>
          </div>
        </div>
      </div>

      <div class="start-node-options">
        <div class="options-title">
          <div>
            <img src="@/assets/img/workflow/global-variable.svg" alt="" class="title-icon" />
            <span>自定义全局变量</span> 
          </div>
          <div>
            <a-tooltip arrow-point-at-center :overlayInnerStyle="{width: '600px'}" placement="top">
              <template #title>
                <div style="width: 600px;">
                  定义工作流的全局变量。通过全局变量可以临时缓存数据。全局变量可以通过以下方式赋值：<br>
                  1、通过webapp访问时，可以在连接中拼接参数给全局变量赋值<br>
                  2、通过api访问时，可以再API中传递参数值<br>
                  3、可以通过变量赋值节点，给全局变量赋值
                </div>
              </template>
              <QuestionCircleOutlined style="color: #595959;font-size: 16px;" />
            </a-tooltip>
          </div>
        </div>
        <div class="options-item diy-options-item"  :class="{'is-required': item.required}" v-for="(item, index) in diy_global" :key="item.key">
          <div class="options-item-body">
            <div class="option-label">{{ item.key }}</div>
            <div class="option-type">{{ item.typ }}</div>
          </div>
          <div class="item-actions-box">
            <svg-icon class="action-btn" name="edit-02" @click="handleEdit(item, index)"></svg-icon>
            <svg-icon class="action-btn" name="close-circle" @click="handleDel(index)"></svg-icon>
          </div>
        </div>

        <div class="action-items-box">
          <a-button style="width: 100%;" class="add-btn" type="dashed" @click="handleAdd">
            <template #icon>
              <PlusOutlined />
            </template>
            添加全局变量
          </a-button>
        </div>
      </div>

      <a-modal v-model:open="show" title="新增变量" okText="保 存" cancelText="取 消" @ok="handleOk">
        <a-form
          ref="formRef"
          :model="formState"
          :rules="formRules"
          label-align="right"
          :label-col="{ span: 4 }"
          :wrapper-col="{ span: 19 }"
        >
          <a-form-item label="变量名" name="key">
            <a-input v-model:value="formState.key" placeholder="请输入变量名" />
          </a-form-item>
          <a-form-item label="类型" name="typ">
            <a-select v-model:value="formState.typ">
              <a-select-option :value="item.value" v-for="item in typOptions" :key="item.value">
                <span>{{ item.lable }}</span>
              </a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="必填" name="required">
            <a-switch v-model:checked="formState.required" />
          </a-form-item>
          <a-form-item label="描述" name="desc">
            <a-textarea v-model:value="formState.desc" placeholder="请输入描述" />
          </a-form-item>
        </a-form>
      </a-modal>
    </div>
  </node-common>
</template>

<script>
import { QuestionCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import NodeCommon from '../base-node.vue'

const typOptions = [
  {
    lable: 'string',
    value: 'string'
  },
  {
    lable: 'number',
    value: 'number'
  },
  {
    lable: 'array<string>',
    value: 'array<string>'
  },
]

export default {
  name: 'StartNode',
  components: {
    NodeCommon,
    QuestionCircleOutlined,
    PlusOutlined
  },
  inject: ['getNode', 'getGraph', 'setData'],
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
      sys_global: [],
      diy_global: [],
      show: false,
      typOptions: [...typOptions],
      formState: {
        key: '',
        value: '',
        typ: 'string',
        required: false,
        desc: ''
      },
      editIndex: -1,
      formRules: {
        key: [
          { required: true, message: '请输入变量名', trigger: 'blur' },
          { pattern: /[a-zA-Z_][a-zA-Z0-9_\-.]*/, message: '英文字母和下划线组成', trigger: 'blur' },
          { min: 1, max: 20, message: '长度在1到20个字符', trigger: 'blur' }
        ],
        typ: [
          { required: true, message: '请选择变量类型', trigger: 'change' }
        ],
        required: [
          { required: true, message: '请选择是否必填', trigger: 'change' }
        ],
        desc: [
          { min: 1, max: 50, message: '长度在1到50个字符', trigger: 'blur' }
        ]
      }
    }
  },
  computed: {},
  mounted() {
    let node_params = JSON.parse(this.properties.node_params)

    this.sys_global = node_params.start.sys_global
    this.diy_global = node_params.start.diy_global

    this.updateHeight()
  },
  methods: {
    getHeight() {
      let length = this.diy_global.length;

      return 322 + (38 + 12) * length
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

      node_params.start.diy_global = [...this.diy_global]

      this.setData({
        height: height,
        node_params: JSON.stringify(node_params)
      })
    },
    handleAdd(){
      this.formState = {
        key: '',
        value: '',
        typ:'string',
        required: false,
        desc: ''
      }
      this.editIndex = -1
      this.show = true
    },
    handleEdit(item, index){
      this.formState = {...item}
      this.editIndex = index
      this.show = true
    },
    handleDel(index){
      this.diy_global.splice(index, 1)
      this.update()
    },
    async handleOk() {
      try {
        await this.$refs.formRef.validate()
        if(this.editIndex > -1){
          this.diy_global.splice(this.editIndex, 1, {...this.formState})
          this.editIndex = -1
        }else{
          this.diy_global.push({...this.formState})
        }
        
        this.show = false
        this.formState = {
          key: '',
          value: '',
          typ: 'string',
          required: false,
          desc: ''
        }

        this.update()
      } catch (error) {
        console.log('验证失败')
      }
    }
  }
}
</script>
