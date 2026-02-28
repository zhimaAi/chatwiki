<template>
  <div>
    <a-modal
      v-model:open="open"
      title="添加配置"
      @ok="handleOk"
      okText="下一步"
      :width="670"
      wrapClassName="no-padding-modal"
      :bodyStyle="{ 'max-height': '70vh', 'overflow-y': 'auto' }"
    >
      <div class="model-list-box">
        <div
          class="list-item"
          @click="handleSelectItem(item)"
          v-for="item in canAddModelList"
          :key="item.model_define"
        >
          <div class="radio-box">
            <a-radio :checked="item.model_define == currentDefie"></a-radio>
          </div>
          <div class="model-info-box">
            <div class="name-block">
              <img class="icons-box" :src="item.model_icon_url" alt="" />
              <div class="name">{{ item.model_name }}</div>
              <a-divider type="vertical" />
              <div class="model-tags">
                <a-tooltip v-for="(tag, index) in item.support_list" :key="index">
                  <template #title v-if="tag == 'LLM'"
                    >大语言模型，添加机器人时需要选择一个大语言模型。用户的提问将交由大语言模型，结合知识库给出答案。</template
                  >
                  <template #title v-else-if="tag == 'TEXT EMBEDDING'"
                    >嵌入模型，添加知识库时需要选择嵌入模型，用于将知识库分段向量化。用户提问时也会将问题生成向量，通过比对向量相似度的方式匹配跟提问语义相近的知识库分段。向量检索具备很好的语义分析能力。</template
                  >
                  <template #title v-else-if="tag == 'RERANK'"
                    >重排序模型，通过重排序模型，将从知识库中检索出来的分段进行重新排序。再将重排序的结果前top-K传递给大语言模型生成答案。重排序有助于优化检索结果，提高大语言模型回答的准确性。</template
                  >
                  <template #title v-else-if="tag == 'SPEECH2TEXT'"
                    >一种用于将音视频中的语音转化为文字的技术解决方案</template
                  >
                  <template #title v-else-if="tag == 'TTS'">语音合成模型</template>
                  <template #title v-else-if="tag == 'IMAGE'">图像生成模型</template>
                  <span class="tag-item">{{ tag }}</span>
                </a-tooltip>
              </div>
            </div>
            <div class="desc-block">
              {{ item.introduce
              }}<a :href="item.help_links" target="_blank">在{{ item.model_name }}申请API KEY</a>
            </div>
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { showModelConfigList } from '@/api/model/index'
import { message } from 'ant-design-vue'
const open = ref(false)

const emit = defineEmits('add')

const currentDefie = ref('')
const selectItem = ref(null)

const canAddModelList = ref([])
const show = () => {
  currentDefie.value = ''
  selectItem.value = null
  getCanAddModelList()
  open.value = true
}

const handleSelectItem = (item) => {
  currentDefie.value = item.model_define
  selectItem.value = item
}
const handleOk = () => {
  if (!currentDefie.value) {
    return message.error('请选择模型')
  }
  open.value = false
  emit('add', selectItem.value)
}

const getCanAddModelList = () => {
  showModelConfigList().then((res) => {
    canAddModelList.value = res.data || []
  })
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.model-list-box {
  margin-top: 24px;
  display: flex;
  flex-direction: column;
  .list-item {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px;
    border-bottom: 1px solid var(--07, #f0f0f0);
    cursor: pointer;
    .model-info-box {
      display: flex;
      flex-direction: column;
      justify-content: space-between;
      .name-block {
        display: flex;
        align-items: center;
        .icons-box {
          height: 22px;
          margin-right: 8px;
        }
        .name {
          color: #000000;
          font-size: 16px;
          line-height: 24px;
          font-weight: 600;
        }
      }
      .model-tags {
        display: flex;
        align-items: center;
        flex-wrap: wrap;
        gap: 8px;
        .tag-item {
          height: 22px;
          padding: 0 8px;
          width: fit-content;
          display: flex;
          align-items: center;
          border: 1px solid #00000026;
          border-radius: 2px;
          background: #0000000a;
          color: #000000a6;
          font-size: 12px;
          font-weight: 400;
        }
      }
    }
    .desc-block {
      color: #8c8c8c;
      font-size: 12px;
      font-weight: 400;
      line-height: 20px;
      display: inline-flex;
      align-items: center;
      flex-wrap: wrap;
      gap: 0 8px;
    }
  }
}
</style>
