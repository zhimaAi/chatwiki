<template>
  <div class="guess-you-want">
    <div class="message-tabs">
      <div
        v-if="props.item.guess_you_want"
        @click="changeTabkey(props.item, 1)"
        class="tab-item"
        :class="{ active: props.item.question_tabkey == 1 }"
      >
        猜你想问
      </div>
      <div
        v-if="
          props.item.guess_you_want && props.common_question_list && props.enable_common_question
        "
        class="v-line"
      ></div>
      <div
        v-if="props.common_question_list && props.enable_common_question"
        @click="changeTabkey(props.item, 2)"
        class="tab-item"
        :class="{ active: props.item.question_tabkey == 2 }"
      >
        常见问题
      </div>
    </div>
    <div class="message-menus" v-if="props.item.question_tabkey == 1">
      <div
        @click="onClickMeun(guess)"
        class="menu-item"
        v-for="(guess, guessIndex) in props.item.guess_you_want"
        :key="guessIndex"
      >
        {{ guess }}
      </div>
    </div>
    <div class="message-menus" v-else>
      <div
        @click="onClickMeun(guess)"
        class="menu-item"
        v-for="(guess, guessIndex) in common_question_list"
        :key="guessIndex"
      >
        {{ guess }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const emit = defineEmits(['clickMeun'])
const props = defineProps({
  item: {
    type: Object,
    default: () => {}
  },
  common_question_list: {
    type: Object,
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
const onClickMeun = (data) => {
  emit('clickMeun', data)
}
</script>

<style lang="less" scoped>
.guess-you-want {
  border-radius: 4px 16px 16px 16px;
  margin-top: 8px;
  padding: 16px 12px;
  background: #fff;
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
        font-weight: 600;
      }
    }
    .v-line {
      width: 1px;
      height: 14px;
      background: #d9d9d9;
      margin: 0 16px;
    }
  }
  .message-content {
    color: #1a1a1a;
    font-size: 14px;
    line-height: 20px;
  }
  .message-menus {
    .menu-item {
      line-height: 22px;
      margin-top: 8px;
      font-size: 14px;
      border-radius: 4px;
      color: rgb(22, 71, 153);
      cursor: pointer;
      background: #e6efff;
      padding: 6px 12px;
    }
  }
}
</style>