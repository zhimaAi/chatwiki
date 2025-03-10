<style lang="less" scoped>
.setting-box {
  .form-box {
    .question-title {
      margin-bottom: 8px;
    }

    .question-option {
      display: flex;
      align-items: center;
      margin-bottom: 8px;

      .question-option-content {
        flex: 1;
      }

      .action-box {
        padding-left: 10px;
      }

      .del-btn {
        font-size: 14px;
        color: #595959;
      }

      .drag-btn {
        width: 24px;
        font-size: 16px;
        color: #8c8c8c;
        cursor: pointer;
      }
    }

    .add-btn {
      height: 32px;
      line-height: 32px;
      padding: 0 12px;
      margin: 0 24px;
      font-size: 14px;
      color: #595959;
      background-color: #fff;
      border: 1px dashed #d9d9d9;
      cursor: pointer;
    }
  }

  .robot-unknown-content {
    white-space: pre-wrap;
    word-break: break-all;
  }

  .robot-unknown-question {
    .question-item {
      line-height: 22px;
      padding: 4px 12px;
      margin-top: 8px;
      font-size: 14px;
      border-radius: 2px;
      color: #595959;
      border: 1px solid #d9d9d9;
    }
  }
}
</style>

<template>
  <edit-box
    class="setting-box"
    title="未知问题提示语"
    icon-name="unknown-prompt"
    v-model:isEdit="isEdit"
    @save="onSave"
    @edit="handleEdit"
  >
    <template #tip>
      <a-tooltip placement="top" :overlayInnerStyle="{ width: '400px' }">
        <template #title>
          <span
            >仅知识库模式下，用户提问没有在知识库中大于score阈值的分段时，会直接回复未知问题提示语。</span
          >
        </template>
        <QuestionCircleOutlined />
      </a-tooltip>
    </template>
    <div class="form-box" v-show="isEdit">
      <div class="question-title">
        <a-textarea
          v-model:value="formState.unknown_question_prompt.content"
          placeholder="请输入未知问题提示语"
        />
      </div>
      <div class="question-options">
        <draggable
          v-model="formState.unknown_question_prompt.question"
          item-key="id"
          handle=".drag-btn"
          @start="drag = true"
          @end="drag = false"
        >
          <template #item="{ element, index }">
            <div class="question-option">
              <span class="drag-btn"><svg-icon name="drag" /></span>
              <a-input
                class="question-option-content"
                v-model:value="element.content"
                placeholder="请输入问题"
              />
              <div class="action-box">
                <CloseCircleOutlined class="del-btn" @click="deleteOption(index)" />
              </div>
            </div>
          </template>
        </draggable>

        <div class="add-btn" @click="addQuestion">
          <PlusOutlined class="add-btn-icon" />添加引导问题
        </div>
      </div>
    </div>
    <div class="robot-info-box" v-show="!isEdit">
      <div class="robot-unknown-content">{{ robotInfo.unknown_question_prompt.content }}</div>
      <div class="robot-unknown-question">
        <div
          class="question-item"
          v-for="(question, index) in robotInfo.unknown_question_prompt.question"
          :key="index"
        >
          {{ question.content }}
        </div>
      </div>
    </div>
  </edit-box>
</template>
<script setup>
import { ref, reactive, inject, toRaw, nextTick, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined, CloseCircleOutlined } from '@ant-design/icons-vue'
import draggable from 'vuedraggable'
import EditBox from './edit-box.vue'
import { getUuid } from '@/utils/index'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'

const isEdit = ref(false)
const drag = ref(false)
const { robotInfo, updateRobotInfo, scrollBoxToBottom } = inject('robotInfo')

const formState = reactive({
  unknown_question_prompt: {
    content: '',
    question: []
  }
})

const checkUnknownQuestion = () => {
  let isEmity = false

  if (formState.unknown_question_prompt.question.length > 0) {
    for (let i = 0; i < formState.unknown_question_prompt.question.length; i++) {
      if (!formState.unknown_question_prompt.question[i].content) {
        isEmity = i + 1
        break
      }
    }
  }

  return isEmity
}

const addQuestion = () => {
  formState.unknown_question_prompt.question.push({ content: '', id: getUuid(8) })
  nextTick(() => {
    scrollBoxToBottom()
  })
}

const deleteOption = (index) => {
  formState.unknown_question_prompt.question.splice(index, 1)
}

const onSave = () => {
  // if (!formState.unknown_question_prompt.content) {
  //   return message.error('请输入未知问题提示语')
  // }

  // if (checkUnknownQuestion()) {
  //   return message.error('引导问题内容不能为空')
  // }
  updateRobotInfo({ ...toRaw(formState) })
  isEdit.value = false;
}

const handleEdit = () => {
  let data = { ...robotInfo.unknown_question_prompt }
  let question = []

  if (data.question) {
    data.question.forEach(({ content, id }) => {
      let item = {
        id,
        content
      }

      if (!item.id) {
        item.id = getUuid(6)
      }

      question.push(item)
    })
  }

  if (data.content) {
    formState.unknown_question_prompt.content = data.content
  }

  if (question) {
    formState.unknown_question_prompt.question = question
  }
}

onMounted(() => {})
</script>
