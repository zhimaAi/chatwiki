<style lang="less" scoped>
.question-input {

  .add-question,
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
  }
}
</style>
<template>
  <div class="question-input">
    <a-button class="add-question" type="dashed" @click="addQuestion">
      <PlusOutlined /> 添加引导问题
    </a-button>
    <div class="question-title">
      <a-textarea v-model:value="content" placeholder="请输入欢迎语" @change="onChangeContent" />
    </div>
    <div class="question-options">
      <div class="question-option" v-for="(item, index) in question" :key="index">
        <a-input class="question-option-content" v-model:value="item.content" placeholder="请输入问题"
          @change="onChangeQuestion" />
        <div class="action-box">
          <CloseCircleOutlined class="del-btn" @click="deleteOption(index)" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { Form } from 'ant-design-vue'
import { PlusOutlined, CloseCircleOutlined } from '@ant-design/icons-vue'

const emit = defineEmits(['update:value'])

const props = defineProps({
  value: {
    type: Object,
    default: () => {
      return {
        content: '',
        question: []
      }
    }
  }
})

const formItemContext = Form.useInjectFormItemContext()

const content = ref('')
const question = ref([])

const addQuestion = () => {
  question.value.push({ content: '' })
}

const deleteOption = (index) => {
  question.value.splice(index, 1)
  triggerChange()
}

const triggerChange = () => {
  let data = {
    content: content.value,
    question: [...question.value]
  }
  emit('update:value', { ...data })

  formItemContext.onFieldChange()
}

const onChangeContent = () => {
  triggerChange()
}

const onChangeQuestion = () => {
  triggerChange()
}

watch(
  () => props.value,
  (newValue) => {
    content.value = newValue.content
    question.value = newValue.question
  },
  {
    immediate: true
  }
)
</script>
