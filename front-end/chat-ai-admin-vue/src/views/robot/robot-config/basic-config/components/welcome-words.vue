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

  .robot-welcome-content {
    white-space: pre-wrap;
    word-break: break-all;
  }

  .robot-welcome-question {
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
    title="欢迎语"
    icon-name="huanyingyu"
    v-model:isEdit="isEdit"
    @save="onSave"
    @edit="handleEdit"
  >
    <div class="form-box" v-show="isEdit">
      <div class="question-title">
        <a-textarea v-model:value="formState.welcomes.content" placeholder="请输入欢迎语" />
      </div>
      <div class="question-options">
        <draggable
          v-model="formState.welcomes.question"
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
      <div class="robot-welcome-content">{{ robotInfo.welcomes.content }}</div>
      <div class="robot-welcome-question">
        <div
          class="question-item"
          v-for="(question, index) in robotInfo.welcomes.question"
          :key="index"
        >
          {{ question.content }}
        </div>
      </div>
    </div>
  </edit-box>
</template>
<script setup>
import { ref, reactive, inject, toRaw, nextTick } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined, CloseCircleOutlined } from '@ant-design/icons-vue'
import draggable from 'vuedraggable'
import EditBox from './edit-box.vue'
import { getUuid } from '@/utils/index'

const isEdit = ref(false)
const drag = ref(false)
const { robotInfo, updateRobotInfo, scrollBoxToBottom } = inject('robotInfo')

const formState = reactive({
  welcomes: {
    content: '',
    question: []
  }
})

const checkWelcomeQuestion = () => {
  let isEmity = false

  if (formState.welcomes.question.length > 0) {
    for (let i = 0; i < formState.welcomes.question.length; i++) {
      if (!formState.welcomes.question[i].content) {
        isEmity = i + 1
        break
      }
    }
  }

  return isEmity
}

const addQuestion = () => {
  formState.welcomes.question.push({ content: '', id: getUuid(8) })
  // nextTick(() => {
  //   scrollBoxToBottom()
  // })
}

const deleteOption = (index) => {
  formState.welcomes.question.splice(index, 1)
}

const onSave = () => {
  if (!formState.welcomes.content) {
    return message.error('请输入欢迎语')
  }

  if (checkWelcomeQuestion()) {
    return message.error('引导问题内容不能为空')
  }

  updateRobotInfo({ ...toRaw(formState) })
  isEdit.value = false;
}

const handleEdit = () => {
  let data = { ...robotInfo.welcomes }
  let question = []

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

  formState.welcomes.content = data.content
  formState.welcomes.question = question
}
</script>
