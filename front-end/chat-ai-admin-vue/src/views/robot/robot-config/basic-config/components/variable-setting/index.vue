<style lang="less" scoped>
.setting-box {
  .variable-list-box {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-wrap: wrap;
    .variable-item {
      display: flex;
      align-items: center;
      color: #595959;
      font-size: 14px;
      line-height: 22px;
      gap: 4px;
      .variable-info-block {
        border: 1px solid var(--06, #d9d9d9);
        padding: 4px 8px;
        border-radius: 6px;
        padding-left: 12px;
        display: flex;
        align-items: center;
        background: #fff;
      }
      .type-tags {
        padding: 0px 8px;
        border-radius: 6px;
        border: 1px solid #00000026;
        color: #595959;
        font-size: 12px;
        line-height: 20px;
        height: 22px;
      }
      .btn-box {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 32px;
        height: 32px;
        border-radius: 6px;
        font-size: 16px;
        cursor: pointer;
        transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
        &:hover {
          background: var(--07, #e4e6eb);
        }
      }
      .is-required {
        &::before {
          content: '*';
          color: red;
          font-size: 12px;
          margin-right: 2px;
        }
      }
    }
  }
}
</style>

<template>
  <edit-box class="setting-box" title="变量" icon-name="variable-icon">
    <template #tip>
      <a-tooltip placement="top" :overlayInnerStyle="{ width: '400px' }">
        <template #title>
          <span
            >变量将以表单形式让用户在对话前填写，用户填写的表单内容将自动替换提示词中的变量。提示词中输入/可选择变量</span
          >
        </template>
        <QuestionCircleOutlined />
      </a-tooltip>
    </template>
    <template #extra>
      <div class="actions-box">
        <a-flex :gap="8">
          <a-button size="small" @click="handleOpenModal">新增变量</a-button>
        </a-flex>
      </div>
    </template>
    <div class="form-box">
      <div class="variable-list-box">
        <div class="variable-item" v-for="item in chatVariables" :key="item.id">
          <div class="variable-info-block">
            <span :class="{ 'is-required': item.must_input == 1 }">{{ item.variable_key }}</span>
            <span>（{{ item.variable_name }}）</span>
            <span class="type-tags">
              {{ typeNameMap[item.variable_type] }}
            </span>
          </div>
          <a-dropdown>
            <div class="btn-box">
              <svg-icon name="more-cycle"></svg-icon>
            </div>
            <template #overlay>
              <a-menu>
                <a-menu-item>
                  <a @click="handleEdit(item)">编 辑</a>
                </a-menu-item>
                <a-menu-item>
                  <a @click="handleDelVaribel(item)">删 除</a>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </div>
    </div>
    <AddVariableModal ref="addVariableModalRef" @ok="handleOk" />
  </edit-box>
</template>
<script setup>
import { ref, onMounted, computed, createVNode } from 'vue'
import { EditOutlined, CloseCircleOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue'
import EditBox from '../edit-box.vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import { deleteChatVariable } from '@/api/robot/index'
import AddVariableModal from './add-variable-modal.vue'
import { message, Modal } from 'ant-design-vue'
import { useRobotStore } from '@/stores/modules/robot'
const robotStore = useRobotStore()

let typeNameMap = {
  input_string: 'string',
  input_number: 'number',
  select_one: 'string',
  checkbox_switch: 'boolean'
}

const chatVariables = computed(() => {
  return robotStore.chatVariables
})

const addVariableModalRef = ref(null)
const handleOpenModal = () => {
  if(chatVariables.value.length >= 10){
    return message.error('最多添加10个变量')
  }
  addVariableModalRef.value.show()
}
const handleEdit = (item) => {
  addVariableModalRef.value.show(JSON.parse(JSON.stringify(item)))
}

const handleOk = () => {
  robotStore.fetchChatVariables()
}

const handleDelVaribel = (item) => {
  Modal.confirm({
    title: '删除确认?',
    icon: createVNode(ExclamationCircleOutlined),
    content: `确认删除该变量【${item.variable_name}】?`,
    okType: 'danger',
    onOk() {
      deleteChatVariable({
        robot_key: item.robot_key,
        id: item.id
      }).then((res) => {
        message.success('删除成功')
        handleOk()
      })
    }
  })
}

onMounted(() => {})
</script>
