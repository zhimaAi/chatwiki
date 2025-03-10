<style lang="less" scoped>
.common-problem-box {
  position: relative;
  .left-content-box {
    width: 640px;
  }
  .preview-box {
    position: absolute;
    left: calc(640px + 48px);
    top: 0;
    width: 345px;
    height: 720px;
    box-shadow: 0 4px 16px 0 #00000033;
    background: #f0f2f5;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    .phone-body {
      padding: 10px;
      flex: 1;
      overflow-x: hidden;
      overflow-y: auto;
      -ms-overflow-style: none;
      scrollbar-width: none;
      &::-webkit-scrollbar {
        display: none;
      }
      .question-box {
        width: 232px;
        border-radius: 3.66px 14.63px 14.63px 14.63px;
        background: #fff;
        padding: 14px 10px;
        margin-left: 48px;
        .name {
          font-size: 12.81px;
          font-weight: 600;
          line-height: 20.12px;
          color: #8c8c8c;
        }
        .list-item {
          cursor: pointer;
          margin-top: 8px;
          padding: 5.49px 10.98px;
          border-radius: 3.66px;
          background: #e5efff;
          color: #164799;
          font-size: 12.81px;
          font-weight: 400;
          line-height: 18.29px;
          word-break: break-all;
        }
      }
    }
    .head-img {
      width: 343px;
    }
    .body-img {
      width: 321px;
    }
    .footer-img {
      width: 343px;
    }
  }
}
.form-box {
  margin-top: 12px;
}
.quess-icon {
  color: #8c8c8c;
  font-size: 16px;
  cursor: pointer;
}
.drag-btn {
  display: flex;
  align-items: center;
  margin-top: 2px;
  margin-right: 8px;
  cursor: grab;
}
.lists-box {
  min-height: 600px;
  margin-top: 16px;
  .list-item {
    position: relative;
    display: flex;
    align-items: center;
    padding: 14px 16px;
    background: #fff;
    border: 1px solid #d9d9d9;
    justify-content: space-between;
    margin-bottom: 8px;
    .title {
      flex: 1;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      color: #262626;
      font-size: 14px;
    }
    .opration-box {
      background: #fff;
      padding-left: 8px;
      z-index: 999;
      position: absolute;
      right: 16px;
      opacity: 0;
      margin-left: 8px;
      display: flex;
      gap: 8px;
      cursor: pointer;
      color: #8c8c8c;
      font-size: 16px;
    }
    &:hover {
      .opration-box {
        opacity: 1;
      }
    }
  }
}

.empty-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-bottom: 40px;
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
  .desc {
    margin-top: 4px;
    margin-bottom: 16px;
    color: #7a8699;
    font-size: 14px;
    line-height: 22px;
  }
}
</style>

<template>
  <div class="common-problem-box">
    <div class="left-content-box">
      <edit-box class="setting-box" title="常见问题设置" icon-name="common-quession">
        <template #icon>
          <a-tooltip>
            <template #title>最多10条。开启后，常见问题会显示在每条问答后面</template>
            <QuestionCircleOutlined class="quess-icon" />
          </a-tooltip>
        </template>
        <template #extra>
          <a-switch
            @change="handleChangeStatus"
            class="switch-item"
            checkedValue="true"
            unCheckedValue="false"
            v-model:checked="robotInfo.enable_common_question"
            checked-children="开"
            un-checked-children="关"
          />
          <a-divider type="vertical" />
          <a-button type="primary" size="small" @click="open">添加</a-button>
        </template>
        <div class="empty-box" v-if="common_question_list.length == 0">
          <img src="@/assets/empty.png" alt="" />
          <div class="title">暂无数据</div>
          <div class="desc">常见问题会显示在每条对话后面,供用户快速提问</div>
          <a-button @click="open" type="primary">立即添加</a-button>
        </div>
        <div class="lists-box" v-else>
          <draggable
            v-model="common_question_list_show"
            handle=".drag-btn"
            item-key="index"
            @end="onDragEnd"
          >
            <template #item="{ element, index }">
              <div class="list-item" :key="index">
                <span class="drag-btn"><svg-icon name="drag" /></span>
                <div class="title">{{ element }}</div>
                <div class="opration-box">
                  <EditOutlined @click="edit(element, index)" />
                  <CloseCircleOutlined @click="delQuestion(index)" />
                </div>
              </div>
            </template>
          </draggable>
        </div>
      </edit-box>
    </div>
    <div class="preview-box">
      <div class="phone-head">
        <img class="head-img" src="@/assets/img/robot/preview/phone-head.png" alt="" />
      </div>
      <div class="phone-body">
        <div><img class="body-img" src="@/assets/img/robot/preview/phone-body.png" alt="" /></div>
        <div
          class="question-box"
          v-if="common_question_list.length && robotInfo.enable_common_question == 'true'"
        >
          <div class="name">常见问题</div>
          <div class="list-item" v-for="item in common_question_list" :key="item">
            {{ item }}
          </div>
        </div>
      </div>
      <div class="phone-footer">
        <img class="footer-img" src="@/assets/img/robot/preview/phone-footer.png" alt="" />
      </div>
    </div>
    <a-modal v-model:open="show" :title="modalTitle" @ok="handleOk" width="476px">
      <a-form class="form-box" layout="vertical">
        <a-form-item label="问题名称" v-bind="validateInfos.question">
          <a-input
            :maxLength="100"
            v-model:value="formState.question"
            placeholder="请输入问题"
          ></a-input>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>
<script setup>
import { ref, reactive, inject, toRaw, computed, createVNode, watch } from 'vue'
import { Form, message, Modal } from 'ant-design-vue'
import draggable from 'vuedraggable'
import {
  EditOutlined,
  CloseCircleOutlined,
  ExclamationCircleOutlined,
  QuestionCircleOutlined
} from '@ant-design/icons-vue'
import EditBox from './edit-box.vue'
const isEdit = ref(false)
const { robotInfo, updateRobotInfo } = inject('robotInfo')
const common_question_list = computed(() => {
  return robotInfo.common_question_list ? JSON.parse(robotInfo.common_question_list) : []
})
const common_question_list_show = ref([])
common_question_list_show.value = common_question_list.value
watch(common_question_list, () => {
  common_question_list_show.value = common_question_list.value
})
const useForm = Form.useForm
const show = ref(false)
const modalTitle = ref('添加问题')
const formState = reactive({
  question: '',
  index: -1
})

const formRules = reactive({
  question: [
    {
      message: '请输入问题名称',
      required: true
    }
  ]
})

const { resetFields, validate, validateInfos } = useForm(formState, formRules)

const open = () => {
  if (common_question_list.value.length >= 10) {
    return message.error('最多设置10条常用问题')
  }
  show.value = true
  resetFields()
  formState.question = ''
  formState.index = -1
}
const edit = (question, index) => {
  show.value = true
  resetFields()
  formState.question = question
  formState.index = index
}
const delQuestion = (index) => {
  let commonQuestionList = []
  commonQuestionList = common_question_list.value.filter((item, i) => index != i)
  Modal.confirm({
    title: '提醒',
    icon: createVNode(ExclamationCircleOutlined),
    content: '确定要删除该问题吗?',
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
    onOk() {
      updateRobotInfo({
        common_question_list: JSON.stringify(commonQuestionList)
      })
    },
    onCancel() {}
  })
}
const handleOk = () => {
  validate().then(() => {
    let commonQuestionList = []
    if (formState.index >= 0) {
      // 编辑
      commonQuestionList = common_question_list.value.map((item, index) =>
        index === formState.index ? formState.question : item
      )
    } else {
      // 新增
      commonQuestionList = [formState.question, ...common_question_list.value]
    }
    updateRobotInfo({
      common_question_list: JSON.stringify(commonQuestionList)
    })
    show.value = false
  })
}
const handleChangeStatus = (val) => {
  updateRobotInfo({})
}

const onDragEnd = (e) => {
  updateRobotInfo({
    common_question_list: JSON.stringify(common_question_list_show.value)
  })
}
</script>
