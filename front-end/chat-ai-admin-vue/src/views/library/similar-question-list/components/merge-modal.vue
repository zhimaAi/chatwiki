<template>
  <a-modal
    :open="visible"
    :title="t('title')"
    @ok="handleConfirm"
    @cancel="handleCancel"
    :width="600"
  >
    <div class="merge-modal-content">
      <div class="merge-modal-subtitle">{{ t('subtitle') }}</div>
      <a-radio-group v-model:value="selectedOption" class="merge-options">
        <div
          v-for="item in options"
          :key="item.id"
          class="merge-option"
          @click="selectedOption = item.id"
        >
          <a-radio class="merge-option-radio" :value="item.id"></a-radio>
          <div class="merge-option-content">
              <div class="merge-option-answer">{{ item.answer || item.question }}</div>
              <div v-if="item.images && item.images.length" class="fragment-img" v-viewer>
                <img v-for="(img, index) in item.images" :key="index" :src="img" alt="" />
              </div>
            </div>
        </div>
      </a-radio-group>
    </div>
  </a-modal>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.similar-question-list.components.merge-modal')

const emit = defineEmits(['confirm', 'cancel'])

const visible = ref(false)
const options = ref([])
const selectedOption = ref(null)

const open = ({ options: opts, defaultSelected }) => {
  options.value = opts || []
  selectedOption.value = defaultSelected || null
  visible.value = true
}

const handleConfirm = async () => {
  if (!selectedOption.value) {
    return
  }
  emit('confirm', selectedOption.value, options.value)
  visible.value = false
}

const handleCancel = () => {
  emit('cancel')
  visible.value = false
}

defineExpose({
  open
})
</script>

<style lang="less" scoped>
.merge-modal-content {
  .merge-modal-subtitle {
    font-size: 14px;
    color: #262626;
    margin-bottom: 20px;
    line-height: 22px;
  }

  .merge-options {
    width: 100%;
    .merge-option {
      display: flex;
      align-items: center;
      margin-bottom: 8px;
      cursor: pointer;

      .merge-option-radio{
        margin-right: 0;
      }

      .merge-option-content {
        flex: 1;
        border: 1px solid #d9d9d9;
        border-radius: 6px;
        padding: 8px 12px;
        margin-left: 12px;
        
        .merge-option-answer {
          font-size: 14px;
          color: #595959;
          line-height: 22px;
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
            object-fit: cover;
          }
        }
      }
    }
  }
}
</style>
