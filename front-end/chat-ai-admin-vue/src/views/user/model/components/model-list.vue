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
                <template #title v-if="tag == 'LLM'">{{ t('llm_tooltip') }}</template>
                <template #title v-else-if="tag == 'TEXT EMBEDDING'">{{ t('embedding_tooltip') }}</template>
                <template #title v-else-if="tag == 'RERANK'">{{ t('rerank_tooltip') }}</template>
                <template #title v-else-if="tag == 'SPEECH2TEXT'">{{ t('speech2text_tooltip') }}</template>
                <template #title v-else-if="tag == 'TTS'">{{ t('tts_tooltip') }}</template>
                <template #title v-else-if="tag == 'IMAGE'">{{ t('image_tooltip') }}</template>
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
        <a-button v-if="item.model_define != 'azure'" class="add-btn" @click="handleSee(item)">{{
          t('view_model')
        }}</a-button>
        <a-button class="add-btn" @click="handleAdd(item)">{{ t('add_config') }}</a-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { SettingOutlined, PlusOutlined, EyeOutlined, FormOutlined } from '@ant-design/icons-vue'
import { toRaw } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.model.components.model-list')

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
