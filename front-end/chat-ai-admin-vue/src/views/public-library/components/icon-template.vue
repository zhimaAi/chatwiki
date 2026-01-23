<template>
  <a-modal
    v-model:open="visible"
    :title="t('title')"
    :width="760"
    @cancel="handleCancel"
    @ok="handleConfirm"
    :confirm-loading="saveLoading"
  >
    <div class="icon-template-content">
      <!-- 模板选项 -->
      <div class="template-options">
        <a-row :gutter="16">
          <!-- 模板1 -->
          <a-col :span="8" v-for="item in templateList" :key="item.id">
            <div
              class="template-card"
              :class="{ active: selectedTemplate == item.id }"
              
            >
              <div class="template-preview" @click="selectTemplate(item.id)">
                <div class="folder-item" :class="['level-' + demo.level]" v-for="demo in item.preview" :key="demo.level">
                  <span class="folder-icon">{{ demo.icon }}</span>
                  <span>{{ demo.text }}</span>
                </div>
              </div>
              <div class="template-footer">
                <a-radio :checked="selectedTemplate == item.id" @click="selectTemplate(item.id)">
                  <span class="template-name">{{ item.name }}</span>
                </a-radio>
              </div>
            </div>
          </a-col>
        </a-row>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { onMounted, ref, computed, toRaw } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import {usePublicLibraryStore} from '@/stores/modules/public-library';
import { getIconTemplateList } from '../../../config/open-doc/icon-template';

const { t } = useI18n('views.public-library.components.icon-template')

const emit = defineEmits(['update:open', 'ok', 'cancel'])
const publicLibraryStore = usePublicLibraryStore()
const { saveEditLibrary } = publicLibraryStore;

const libraryInfo = computed(() => {
  return publicLibraryStore.libraryInfo
})

const visible = ref(false)
const saveLoading = ref(false)
const templateList = ref([])

const selectedTemplate = ref('1')
const open = () => {
  selectedTemplate.value = libraryInfo.value.icon_template_config_id
  visible.value = true;
}
// 选择模板
const selectTemplate = (templateId) => {
  selectedTemplate.value = templateId
}

// 取消
const handleCancel = () => {
  visible.value = false
  emit('cancel')
}

// 确定
const handleConfirm = () => {
  let data = toRaw(libraryInfo.value)

  data.icon_template_config_id = selectedTemplate.value;

  saveLoading.value = true;
  saveEditLibrary(data).then(() => {
    visible.value = false;
  }).finally(() => {
    saveLoading.value = false;

    emit('ok', selectedTemplate.value)
  })
}

defineExpose({
  open
})

onMounted(() => {
  templateList.value = getIconTemplateList();
})
</script>

<style lang="less" scoped>
.icon-template-content {
  padding: 20px 0;
}
.template-card {
  display: flex;
  flex-direction: column;
  

  .template-preview {
    position: relative;
    flex: 1;
    padding: 4px 16px;
    border: 1px solid #D9D9D9;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s ease;

    &:hover {
      border-color: #1890ff;
    }
  }

  &.active {
    .template-preview{
      border-color: #1890ff;
      box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
    }

    /* 编辑角标样式 */
    .template-preview::after {
      content: '';
      position: absolute;
      bottom: -2px;
      right: -2px;
      width: 32px;
      height: 32px;
      background: url('../../../assets/svg/check-arrow-filled.svg') 0 0 no-repeat;
    }
  }

  .folder-item {
    display: flex;
    align-items: center;
    line-height: 22px;
    padding: 5px 0;
    margin-bottom: 4px;
    font-size: 14px;

    &:last-child {
      margin-bottom: 0;
    }
  }

  .folder-item.level-0 {
    margin-left: 0;
  }

  .folder-item.level-1 {
    margin-left: 16px;
  }

  .folder-item.level-2 {
    margin-left: 32px;
  }

  .folder-item.level-3 {
    margin-left: 48px;
  }

  .folder-icon {
    margin-right: 4px;
    font-size: 16px;
  }

  .folder-icon.blue {
    color: #1890ff;
  }

  .folder-icon.yellow {
    color: #faad14;
  }

  .folder-icon.red {
    color: #ff4d4f;
  }

  .folder-icon.gray {
    color: #8c8c8c;
  }

  .template-footer {
    display: flex;
    align-items: center;
    justify-content: center;
    padding-top: 16px;
    border-top: 1px solid #f0f0f0;
  }

  .template-name {
    font-size: 14px;
    color: #666;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
  }
}
</style>
