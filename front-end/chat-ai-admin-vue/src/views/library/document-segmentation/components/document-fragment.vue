<style lang="less" scoped>
.document-fragment {
  padding: 16px;
  border-radius: 2px;
  background-color: #ffffff;

  .fragment-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 22px;
    line-height: 22px;
    font-size: 14px;

    .fragment-number,
    .fragment-title {
      font-weight: 600;
      color: #000000;
    }

    .fragment-title {
      padding-left: 4px;
    }

    .fragment-content-lenght {
      padding-left: 8px;
      color: #8c8c8c;
    }
  }

  .fragment-content {
    margin-top: 8px;
    line-height: 22px;
    font-size: 14px;
    color: #595959;
  }
}
</style>

<template>
  <div class="document-fragment">
    <div class="fragment-header">
      <div class="fragment-info">
        <span class="fragment-number">#{{ props.number }}</span><span class="fragment-title" v-if="props.title">{{
          props.title }}</span><span class="fragment-content-lenght">共{{ props.total }}个字符</span>
      </div>

      <div class="fragment-action">
        <a @click="handleEdit">编辑</a>
        <a-divider type="vertical" />
        <a @click="handleDelete">删除</a>
      </div>
    </div>

    <div class="fragment-content">{{ props.content }}</div>
    <div class="fragment-content" v-if="props.question">Q：{{ props.question }}</div>
    <div class="fragment-content" v-if="props.answer">A：{{ props.answer }}</div>
  </div>
</template>

<script setup>
const emit = defineEmits(['edit', 'delete'])
const props = defineProps({
  number: {
    type: [Number, String]
  },
  total: {
    type: [Number, String]
  },
  title: {
    type: [Number, String]
  },
  content: {
    type: [Number, String]
  },
  question: {
    type: [Number, String]
  },
  answer: {
    type: [Number, String]
  }
})

const handleEdit = () => {
  emit('edit')
}

const handleDelete = () => {
  emit('delete')
}
</script>
