<style lang="less" scoped>
.model-list {
  display: flex;
  flex-flow: row wrap;
  gap: 24px;
  margin: 16px 24px;

  .model-item {
    display: flex;
    position: relative;
    justify-content: space-between;
    width: 100%;
    padding: 16px;
    border-radius: 2px;
    overflow: hidden;
    border: 1px solid #f0f0f0;
    background-color: #fff;
    transition: all 0.2s;

    &:hover {
      box-shadow:
        0 6px 30px 5px rgba(0, 0, 0, 0.05),
        0 16px 24px 2px rgba(0, 0, 0, 0.04),
        0 8px 10px -5px rgba(0, 0, 0, 0.08);
    }

    .add-btn {
      font-size: 14px;
      cursor: pointer;
    }

    .model-item-left {
      display: flex;
      flex-direction: column;
      align-items: flex-start;
      flex-basis: calc((100% - 320px));
    }

    .model-item-left-top {
      display: flex;
      justify-content: center;
      align-items: center;
      margin-bottom: 10px;
    }

    .line {
      display: inline-block;
      width: 2px;
      height: 15px;
      background-color: #dedede;
      margin: 5px 15px;
    }

    .model-item-right {
      display: flex;
      flex-flow: row wrap;
      gap: 8px;
      align-items: center;
      border-radius: 2px;
    }

    .btn-red {
      color: red;
      border-color: red;
    }

    .item-body {
      flex: 1;
    }

    .model-logo {
      display: block;
      height: 24px;
    }

    .model-name {
      color: #000000;
      font-family: 'PingFang SC';
      font-size: 16px;
      font-style: normal;
      font-weight: 600;
      line-height: 24px;
      margin-left: 8px;
    }

    .model-item-left-bottom {
      display: flex;
      width: 100%;
    }

    .model-desc {
      display: flex;
      flex-basis: calc((100% - 120px));
      line-height: 22px;
      font-size: 14px;
      color: #8c8c8c;
    }

    .model-tags {
      display: flex;
      flex-flow: row wrap;
      gap: 8px;

      .tag {
        height: 22px;
        line-height: 20px;
        padding: 0 8px;
        border-radius: 2px;
        font-size: 12px;
        border: 1px solid rgba(0, 0, 0, 0.15);
        background: rgba(0, 0, 0, 0.04);
      }
    }
  }

  .item-actions {
    display: flex;
    margin-top: 16px;
    border-top: 1px solid #f0f0f0;

    .action-item {
      position: relative;
      flex: 1;
      height: 54px;
      font-size: 14px;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      color: #2475fc;

      &::after {
        position: absolute;
        right: 0;
        display: block;
        content: '';
        width: 1px;
        height: 16px;
        background-color: #e4e6eb;
      }

      &:last-child::after {
        display: none;
      }

      .action-icon {
        margin-right: 4px;
      }
    }
  }
}
</style>

<template>
  <div class="model-list">
    <div class="model-item" v-for="item in props.list" :key="item.model_define">
      <div class="model-item-left">
        <div class="model-item-left-top">
          <img class="model-logo" :src="item.model_icon_url" alt="" />
          <div class="model-name">{{ item.model_name }}</div>
          <div class="line"></div>
          <div class="model-tags">
            <div v-for="(tag, index) in item.support_list" :key="index">
              <a-tooltip>
                <template #title v-if="tag == 'LLM'"
                  >大语言模型，添加机器人时需要选择一个大语言模型。用户的提问将交由大语言模型，结合知识库给出答案。</template
                >
                <template #title v-else-if="tag == 'TEXT EMBEDDING'"
                  >嵌入模型，添加知识库时需要选择嵌入模型，用于将知识库分段向量化。用户提问时也会将问题生成向量，通过比对向量相似度的方式四配跟提问语义相近的知识库分段。向量检索具备很好的语义分析能力。</template
                >
                <template #title v-else-if="tag == 'RERANK'"
                  >重排序模型,通过重排序模型,将从知识库中检索出来的分段进行重新排序。再将重排序的结果前top-K传递个大语言模型生成答案。重排序有助于优化检索结果，提高大语言模型回答的准确性。</template
                >
                <template #title v-else-if="tag == 'SPEECH2TEXT'"
                  >一种用于将音视频中的语音转化为文字的技术解决方案</template
                >
                <template #title v-else-if="tag == 'TTS'">语音合成模型</template>
                <template #title v-else-if="tag == 'IMAGE'">图像生成模型</template>
                <span class="tag">{{ tag }}</span>
              </a-tooltip>
            </div>
          </div>
        </div>
        <div class="model-item-left-bottom">
          <div class="model-desc">{{ item.introduce }}</div>
        </div>
      </div>
      <div class="model-item-right">
        <a-button v-if="item.model_define != 'azure'" class="add-btn" @click="handleSee(item)"
          >查看模型</a-button
        >
        <a-button class="add-btn" @click="handleAdd(item)">添加配置</a-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { SettingOutlined, PlusOutlined, EyeOutlined, FormOutlined } from '@ant-design/icons-vue'
import { toRaw } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n()

const emit = defineEmits(['edit', 'new', 'see', 'remove', 'add'])

const modelDefine = ['azure', 'ollama', 'xinference', 'openaiAgent']

const props = defineProps({
  type: {
    type: Number,
    default: () => 1
  },
  list: {
    type: Array,
    default: () => []
  }
})

const handleSee = (item) => {
  emit('see', toRaw(item))
}

const handleEdit = (item) => {
  let data = toRaw(item)

  emit('edit', data, data.config_list[0])
}

const handleNew = (item) => {
  emit('new', toRaw(item))
}

const handleAdd = (item) => {
  emit('add', toRaw(item))
}

const handleDel = (item) => {
  let data = toRaw(item)

  emit('remove', data.config_list[0])
}
</script>
