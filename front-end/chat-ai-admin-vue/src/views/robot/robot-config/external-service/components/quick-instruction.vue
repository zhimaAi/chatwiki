<style lang="less" scoped>
.form-box {
  margin-top: 12px;
}
.quick-instructions-box {
  // word-break: break-all;
}

.quick-action-right {
  display: flex;
  justify-content: center;
  align-items: center;

  .switch-item {
    margin-right: 10px;
  }
}

.quick-list-box {
  .list-box-wrap {
    position: relative;
  }
  .drag-btn {
    cursor: grab;
    margin-right: 8px;
  }
  .list-item {
    position: relative;
    display: flex;
    align-items: center;
    padding: 14px 16px;
    border: 1px solid #d9d9d9;
    background: #fff;
    border-radius: 2px;
    margin-bottom: 8px;
    cursor: pointer;
    &:hover {
      .opration-box {
        opacity: 1;
      }
    }
  }
  .left-box {
    word-break: break-all;
    width: calc(100% - 16px);
    flex: 1;
    .title {
      color: #262626;
      font-weight: 600;
      line-height: 22px;
      font-size: 14px;
      width: 100%;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }
    .content {
      font-size: 12px;
      color: #8c8c8c;
      line-height: 20px;
      margin-top: 2px;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }
  }
  .opration-box {
    top: 0;
    bottom: 0;
    background: #fff;
    padding-left: 8px;
    z-index: 99;
    position: absolute;
    opacity: 0;
    margin-left: 8px;
    display: flex;
    gap: 8px;
    color: #8c8c8c;
    font-size: 16px;
    cursor: pointer;
    right: 0px;
    padding-right: 16px;
  }
}
.help-icon {
  color: #8c8c8c;
  cursor: pointer;
  font-size: 16px;
}
.help-box {
  color: #262626;
  font-size: 14px;
  font-weight: 400;
  line-height: 22px;
  width: 272px;
  img {
    width: 272px;
    height: 160px;
    margin-top: 12px;
  }
}
</style>

<template>
  <div class="quick-instructions-box">
    <card-box>
      <template #title>
        快捷指令设置
        <a-popover placement="top">
          <template #content>
            <div class="help-box">
              <div>快捷指令是对话输入框上方的按钮,配置完成后,用户可以快速发起预设对话</div>
              <img src="@/assets/img/robot/quick-help.png" alt="" />
            </div>
          </template>
          <QuestionCircleOutlined class="help-icon" />
        </a-popover>
      </template>
      <template #icon>
        <svg-icon name="quick-Instruction" style="font-size: 16px; color: #262626"></svg-icon>
      </template>
      <template #action>
        <div class="quick-action-right">
          <a-switch
            @change="handleChangeShortcutSwitch"
            class="switch-item"
            :checkedValue="1"
            :unCheckedValue="0"
            v-model:checked="fastCommandSwitch"
            checked-children="开"
            un-checked-children="关"
          />
          <a-button @click="openAddModal" size="small" type="primary">添加</a-button>
        </div>
      </template>
      <div class="quick-list-box">
        <draggable
          v-model="quickLists"
          item-key="id"
          @end="(e) => handleDrag(e)"
          handle=".drag-btn"
        >
          <template #item="{ element, index }">
            <div class="list-box-wrap">
              <div class="list-item" :key="element.id">
                <div class="drag-btn"><svg-icon name="drag" /></div>
                <div class="left-box">
                  <div class="title">{{ element.title }}</div>
                  <div class="content">{{ element.content }}</div>
                </div>
                <div class="opration-box">
                  <EditOutlined @click="openAddModal(element)" />
                  <CloseCircleOutlined @click="handleDelete(element.id)" />
                </div>
              </div>
            </div>
          </template>
        </draggable>
      </div>
    </card-box>
    <a-modal v-model:open="show" :title="modalTitle" @ok="handleOk" width="576px">
      <div class="form-box">
        <a-form layout="vertical">
          <a-form-item label="指令标题" v-bind="validateInfos.title">
            <a-input
              placeholder="请输入指令标题"
              :maxLength="20"
              v-model:value="formState.title"
            ></a-input>
          </a-form-item>
          <a-form-item label="指令类型" required>
            <a-radio-group v-model:value="formState.typ" @change="handleChange">
              <a-radio :value="1">输入内容</a-radio>
              <a-radio :value="2">跳转网页</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item
            :label="formState.typ == 1 ? '指令内容' : '网页链接'"
            required
            v-bind="validateInfos.content"
          >
            <a-textarea
              style="height: 100px"
              v-if="formState.typ == 1"
              v-model:value="formState.content"
              placeholder="请输入指令内容"
            />
            <a-input
              v-else
              placeholder="请输入网页链接"
              v-model:value="formState.content"
            ></a-input>
          </a-form-item>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, toRaw, createVNode } from 'vue'
import { Form, message, Modal } from 'ant-design-vue'
import CardBox from './card-box.vue'
import { isValidURL } from '@/utils/validate.js'
import { useRobotStore } from '@/stores/modules/robot'
import { storeToRefs } from 'pinia'
import draggable from 'vuedraggable'
import {
  EditOutlined,
  CloseCircleOutlined,
  ExclamationCircleOutlined,
  QuestionCircleOutlined
} from '@ant-design/icons-vue'
import {
  updateFastCommandSwitch,
  getFastCommandList,
  saveFastCommand,
  deleteFastCommand,
  sortFastCommand
} from '@/api/robot/index'
import { useRoute } from 'vue-router'
const route = useRoute()
const robotStore = useRobotStore()
const { robotInfo } = storeToRefs(robotStore)
const quickLists = ref([])
const fastCommandSwitch = ref(parseInt(robotInfo.value.fast_command_switch) || 0)
const getLists = () => {
  getFastCommandList({
    robot_key: robotInfo.value.robot_key
  }).then((res) => {
    quickLists.value = res.data || []
    robotStore.setQuickCommand(res.data)
  })
}
getLists()
const useForm = Form.useForm
const show = ref(false)
const modalTitle = ref('添加快捷指令')
const formState = reactive({
  robot_id: route.query.id,
  title: '',
  typ: 1,
  content: '',
  id: ''
})

const handleChangeShortcutSwitch = (val) => {
  // 发送请求
  updateFastCommandSwitch({
    robot_key: robotInfo.value.robot_key
  }).then((res) => {
    fastCommandSwitch.value = val
    message.success(`修改成功`)
  }).catch((err) => {
    // 失败后回到原状态
    fastCommandSwitch.value = !fastCommandSwitch.value
  })
}

const handleChange = () => {
  formState.content = ''
  clearValidate(['content'])
}
const openAddModal = (data) => {
  show.value = true
  resetFields()
  modalTitle.value = data.id ? '编辑快捷指令' : '添加快捷指令'
  formState.title = data.title || ''
  formState.content = data.content || ''
  formState.id = data.id || ''
  formState.typ = +data.typ || 1
}
const formRules = reactive({
  title: [
    {
      message: '请输入指令标题',
      required: true
    }
  ],
  content: [
    {
      message: '',
      required: true
    },
    {
      validator: async (rule, value) => {
        if (value == '') {
          return Promise.reject(new Error(`请输入${formState.typ == 1 ? '内容' : '网址'}`))
        }
        if (formState.typ == 1) {
          // 内容
          return Promise.resolve()
        } else {
          if (isValidURL(value)) {
            return Promise.resolve()
          } else {
            return Promise.reject(new Error('请输入合法网页链接'))
          }
        }
      }
    }
  ]
})

const { resetFields, validate, clearValidate, validateInfos } = useForm(formState, formRules)

const handleOk = () => {
  validate().then(() => {
    saveFastCommand({
      ...toRaw(formState)
    }).then((res) => {
      message.success(`修改成功`)
      show.value = false
      getLists()
    })
  })
}
const handleDelete = (id) => {
  Modal.confirm({
    title: '提醒',
    icon: createVNode(ExclamationCircleOutlined),
    content: '确定要删除该快捷指令吗?',
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
    onOk() {
      deleteFastCommand({ id }).then((res) => {
        message.success(`删除成功`)
        getLists()
      })
    },
    onCancel() {}
  })
}

const handleDrag = (e) => {
  let sort = quickLists.value.map((item, index) => {
    return {
      id: +item.id,
      sort: index + 1
    }
  })

  sortFastCommand({
    data: sort,
  }).then((res) => {
    message.success(`保存成功`)
    getLists()
  })
}
</script>
