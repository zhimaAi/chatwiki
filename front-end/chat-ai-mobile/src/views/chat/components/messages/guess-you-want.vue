<template>
  <div class="guess-you-want">
    <div class="message-content">
      <div class="question-list">
        <div class="message-tabs">
          <div
            v-if="props.msg.guess_you_want && props.msg.guess_you_want.length"
            @click="changeTabkey(props.msg, 1)"
            class="tab-item"
            :class="{ active: props.msg.question_tabkey == 1 }"
          >
            猜你想问
          </div>
          <div
            v-if="
              props.msg.guess_you_want &&
              props.msg.guess_you_want.length &&
              props.enable_common_question &&
              props.common_question_list &&
              props.common_question_list.length
            "
            class="v-line"
          ></div>
          <div
            v-if="
              props.common_question_list &&
              props.common_question_list.length &&
              props.enable_common_question
            "
            @click="changeTabkey(props.msg, 2)"
            class="tab-item"
            :class="{ active: props.msg.question_tabkey == 2 }"
          >
            常见问题
          </div>
        </div>
        <template v-if="props.msg.question_tabkey == 1">
          <div
            class="question-item"
            v-for="item in props.msg.guess_you_want"
            :key="item"
            @click="sendTextMessage(item)"
          >
            <span>{{ item }}</span>
          </div>
        </template>
        <template v-if="props.msg.question_tabkey == 2">
          <div
            class="question-item"
            v-for="item in props.common_question_list"
            :key="item"
            @click="sendTextMessage(item)"
          >
            <span>{{ item }}</span>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const emit = defineEmits(['sendTextMessage'])
const props = defineProps({
  msg: {
    type: Object,
    required: true
  },
  common_question_list: {
    type: Array,
    default: () => []
  },
  enable_common_question: {
    type: Boolean,
    default: false
  }
})
const changeTabkey = (item, key) => {
  item.question_tabkey = key
}
const sendTextMessage = (text) => {
  emit('sendTextMessage', text)
}
</script>

<style lang="less" scoped>
.guess-you-want {
  margin-top: 24px;
  .message-content {
    position: relative;
    // display: inline-block;
    padding: 12px;
    margin-right: 50px;
    border-radius: 4px 16px 16px 16px;
    background-color: #fff;
    white-space: pre-wrap;
    word-break: break-all;
  }
  .question-list {
    padding-bottom: 4px;
    .question-item {
      cursor: pointer;
      line-height: 20px;
      padding: 6px 12px;
      margin-top: 8px;
      font-size: 14px;
      font-weight: 400;
      border-radius: 4px;
      color: #164799;
      background-color: #e6efff;
    }
    
    .question-item:hover {
      background: #DEE8FB;
    }

    .question-item:active {
      background-color: #e6efff;
    }
  }
  .message-tabs {
    display: flex;
    align-items: center;
    color: #8c8c8c;
    font-size: 14px;
    font-weight: 400;
    line-height: 22px;
    .tab-item {
      cursor: pointer;
      &.active {
        font-weight: 800;
      }
    }
    .v-line {
      width: 1px;
      height: 14px;
      background: #d9d9d9;
      margin: 0 16px;
    }
  }
}
</style>