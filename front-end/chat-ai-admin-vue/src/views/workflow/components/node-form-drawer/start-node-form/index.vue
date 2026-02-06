<style lang="less" scoped>
.add-trigger-box{
  margin-bottom: 16px;
  .add-trigger-desc{
    line-height: 22px;
    margin-top: 4px;
    font-size: 14px;
    color: #595959;
  }
}
.start-node {
  position: relative;

  .start-node-options {
    background: #f2f4f7;
    border-radius: 6px;
    padding: 12px;
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
      
      .option-label::before{
        content: '*';
        color: #FB363F;
        display: inline-block;
        margin-right: 2px;
        opacity: 0;
      }

      &.is-required .option-label::before{
        content: '*';
        color: #FB363F;
        display: inline-block;
        margin-right: 2px;
        opacity: 1;
      }
      .option-type {
        height: 22px;
        line-height: 18px;
        padding: 0 8px;
        border-radius: 6px;
        border: 1px solid rgba(0, 0, 0, 0.15);
        background-color: #fff;
        color: var(--wf-color-text-3);
        font-size: 12px;
      }
      .option-desc {
        max-width: 90px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        margin-left: 8px;
        font-size: 14px;
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
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        :hideMenus="true"
      >
        <template #desc>
          <span>{{ t('desc_start_node') }}</span>
        </template>
      </NodeFormHeader>
    </template>

    <div class="start-node-form">
      <div class="add-trigger-box">
        <a-button type="primary" ghost @click.stop="handleAddTrigger">
          <template #icon>
            <PlusOutlined />
          </template>
          <span>{{ t('btn_add_trigger') }}</span>
        </a-button>
        <div class="add-trigger-desc">{{ t('desc_add_trigger') }}</div>
      </div>
      <div class="node-form-content">
        <div class="start-node">
          <div class="start-node-options">
            <div class="options-title">
              <div>
                <img src="@/assets/img/workflow/global-variable.svg" alt="" class="title-icon" />
                <span>{{ t('title_custom_global_variable') }}</span>
              </div>
              <div>
                <a-tooltip
                  placement="topRight"
                  :overlayStyle="{ maxWidth: '600px' }"
                >
                  <template #title>
                    <div v-html="t('tip_global_variable_definition')"></div>
                  </template>
                  <QuestionCircleOutlined style="color: #595959; font-size: 16px" />
                </a-tooltip>
              </div>
            </div>
            <div
              class="options-item diy-options-item"
              :class="{ 'is-required': item.required }"
              v-for="(item, index) in diy_global"
              :key="item.key"
            >
              <div class="options-item-body">
                <div class="option-label">{{ item.key }}</div>
                <div class="option-type">{{ item.typ }}</div>
                <a-tooltip :title="item.desc">
                  <div class="option-desc">{{ item.desc }}</div>
                </a-tooltip>
              </div>
              <div class="item-actions-box">
                <svg-icon
                  class="action-btn"
                  name="edit-02"
                  @click="handleEdit(item, index)"
                ></svg-icon>
                <svg-icon class="action-btn" name="close-circle" @click="handleDel(index)"></svg-icon>
              </div>
            </div>

            <div class="action-items-box">
              <a-button style="width: 100%" class="add-btn" type="dashed" @click="handleAdd">
                <template #icon>
                  <PlusOutlined />
                </template>
                {{ t('btn_add_global_variable') }}
              </a-button>
            </div>
          </div>

          <a-modal
            v-model:open="show"
            :title="t('title_add_variable')"
            :okText="t('btn_save')"
            :cancelText="t('btn_cancel')"
            @ok="handleOk"
          >
            <a-form
              ref="formRef"
              :model="formState"
              :rules="formRules"
              label-align="right"
              :label-col="{ span: 6}"
              :wrapper-col="{ span: 18 }"
            >
              <a-form-item :label="t('label_variable_name')" name="key">
                <a-input v-model:value="formState.key" :placeholder="t('ph_input_variable_name')" />
              </a-form-item>
              <a-form-item :label="t('label_type')" name="typ">
                <a-select v-model:value="formState.typ">
                  <a-select-option :value="item.value" v-for="item in typOptions" :key="item.value">
                    <span>{{ item.lable }}</span>
                  </a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item :label="t('label_required')" name="required">
                <a-switch v-model:checked="formState.required" />
              </a-form-item>
              <a-form-item :label="t('label_description')" name="desc">
                <a-textarea v-model:value="formState.desc" :placeholder="t('ph_input_description')" />
              </a-form-item>
            </a-form>
          </a-modal>
        </div>
      </div>
    </div>
  </NodeFormLayout>
</template>

<script>
import { computed } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import NodeFormLayout from '../node-form-layout.vue'
import NodeFormHeader from '../node-form-header.vue'
import { QuestionCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import  {message } from 'ant-design-vue'
import { useWorkflowStore } from '@/stores/modules/workflow'

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
    lable: 'boole',
    value: 'boole'
  },
  {
    lable: 'array<string>',
    value: 'array<string>'
  },
  {
    lable: 'array<number>',
    value: 'array<number>'
  },
  {
    lable: 'array<object>',
    value: 'array<object>'
  },
]

export default {
  emits: ['close', 'update-node'],
  inject: ['getGraph'],
  components: {
    NodeFormLayout,
    NodeFormHeader,
    QuestionCircleOutlined,
    PlusOutlined
  },
  props: {
    nodeType: {
      type: String,
      default: 'start-node',
    },
    node: {
      type: Object,
      default: () => ({})
    }
  },
  setup() {
    const { t } = useI18n('views.workflow.components.node-form-drawer.start-node-form.index')
    const workflowStore = useWorkflowStore()

    const formRules = computed(() => ({
      key: [
        { required: true, message: t('msg_input_variable_name'), trigger: 'blur' },
        { pattern: /[a-zA-Z_][a-zA-Z0-9_\-.]*/, message: t('msg_variable_name_format'), trigger: 'blur' },
        { min: 1, max: 20, message: t('msg_variable_name_length'), trigger: 'blur' }
      ],
      typ: [
        { required: true, message: t('msg_select_type'), trigger: 'change' }
      ],
      required: [
        { required: true, message: t('msg_select_required'), trigger: 'change' }
      ],
      desc: [
        { min: 1, max: 50, message: t('msg_description_length'), trigger: 'blur' }
      ]
    }))

    return {
      t,
      workflowStore,
      formRules
    }
  },
  data() {
    return {
      sys_global: [],
      diy_global: [],
      trigger_list:[],
      show: false,
      typOptions: [...typOptions],
      formState: {
        key: '',
        value: '',
        typ: 'string',
        required: false,
        desc: ''
      },
      editIndex: -1
    }
  },
  mounted() {
    let node_params = JSON.parse(this.node.node_params)
    this.sys_global = node_params.start.sys_global
    this.diy_global = node_params.start.diy_global
    this.trigger_list = node_params.start.trigger_list
  },
  methods: {
    update() {
      let node_params = JSON.parse(this.node.node_params)

      node_params.start.diy_global = [...this.diy_global]

      let data = {
        ...this.node,
        node_params: JSON.stringify(node_params)
      }

      this.$emit('update-node', data)
    },
    handleAddTrigger(){
      this.workflowStore.getTriggerList(this.$route.query.robot_key, true);
      this.getGraph().eventCenter.emit('custom:showTriggerLit')
    },
    getTriggerVariables(){
      let variables = {}
      this.trigger_list.forEach(item => {
        let outputs = item.outputs || []
        outputs.forEach(output => {
          variables[output.variable] = {...output}
        })
      })

      return variables;
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
      this.getTriggerVariables()
      // 判断是不是被触发器关联的变量
      let triggerVariableMap = this.getTriggerVariables()
      let field = this.diy_global[index]

      if(triggerVariableMap['global.' + field.key]){
        return message.error(t('msg_variable_linked_to_trigger'))
      }
 
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
        console.log(t('msg_validation_failed'))
      }
    },
  }
}
</script>
