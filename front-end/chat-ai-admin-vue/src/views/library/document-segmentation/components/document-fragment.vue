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

    .fragment-content-status {
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
  .fragment-img {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-top: 8px;
    img {
      width: 80px;
      height: 80px;
      border-radius: 6px;
      cursor: pointer;
    }
  }
}
</style>

<template>
  <div class="document-fragment">
    <div class="fragment-header">
      <div class="fragment-info">
        <span class="fragment-number">#<span v-if="props.chunk_type == 4">{{ props.father_chunk_paragraph_number }}-</span>{{ props.number }}</span>
        <span class="fragment-title" v-if="props.title">{{ props.title }}</span>
        <span class="fragment-content-lenght">共{{ props.total }}个字符</span>
        <span class="fragment-content-status" v-if="props.status === 'paragraphsSegmented'">
          {{ props.currentData?.status_text || '-' }}<LoadingOutlined v-if="props.currentData?.status == 3" />
          <a-tooltip v-if="props.currentData?.status == 2 && props.currentData.errmsg" :title="props.currentData.errmsg">
            <strong class="cfb363f"
              >原因<ExclamationCircleOutlined class="err-icon cfb363f"
            /></strong>
          </a-tooltip>
        </span>
      </div>

      <div class="fragment-action">
        <a @click="handleEdit">编辑</a>
        <a-divider type="vertical" />
        <a @click="handleDelete">删除</a>
      </div>
    </div>

    <div class="fragment-content">{{ props.content }}</div>
    <div class="fragment-content" v-if="props.question">Q：{{ props.question }}</div>
    <div
      class="fragment-content"
      v-if="props.similar_question_list && props.similar_question_list.join('')"
    >
      相似问法：{{ props.similar_question_list.join('/') }}
    </div>
    <div class="fragment-content" v-if="props.answer">A：{{ props.answer }}</div>
    <div class="fragment-img" v-viewer>
      <img v-for="(item, index) in props.images" :key="index" :src="item" alt="" />
    </div>
  </div>
</template>

<script setup>
import {
  ExclamationCircleOutlined,
  LoadingOutlined,
} from '@ant-design/icons-vue'
const emit = defineEmits(['edit', 'delete'])
const props = defineProps({
  number: {
    type: [Number, String]
  },
  father_chunk_paragraph_number: {
    type: [Number, String]
  },
  chunk_type:{
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
  },
  images: {
    type: [Array, String]
  },
  similar_question_list: {
    type: [Array, String]
  },
  status: {
    type: String,
    default: ''
  },
  currentData: {
    type: Object,
    default: () => {}
  }
})

const handleEdit = () => {
  emit('edit')
}

const handleDelete = () => {
  emit('delete')
}
</script>
