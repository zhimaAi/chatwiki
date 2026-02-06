<style lang="less" scoped>
.message-input-box {
  margin: 12px;
  padding: 10px;
  border-radius: 16px;
  border: 1px solid #ddd;
  box-shadow: 0 2px 6px 0 rgba(0, 0, 0, 0.08);
  transition: all 0.2s;
  background: #fff;

  &.is-focus {
    border: 1px solid #2475fc;
  }

  .message-input-body {
    display: flex;
    justify-content: center;
    align-items: flex-end;
  }

  .message-input {
    position: relative;
    flex: 1;
    padding: 4px 4px 4px 0;
    overflow: hidden;
  }

  

  .message-action{
    display: flex;
    align-items: center;
    
    .send-btn {
      width: 32px;
      height: 32px;
      padding: 0;
      margin: 0;
      font-size: 32px;
      border: none;
      outline: none;
      background: none;
      transition: all 0.2s;
      cursor: pointer;
      color: #2475fc;

      &:hover {
          opacity: 0.8;
      }

      &:disabled {
        opacity: 0.5;
      }
    }

    .file-action{
      position: relative;
      display: flex;
      align-items: center;
      height: 32px;
      padding: 0 8px;
      margin-right: 8px;
      border-radius: 32px;
      border: 1px solid #F0F0F0;

      .line{
        width: 1px;
        height: 14px;
        margin: 0 8px;
        background: #D9D9D9;
      }

      .action-btn{
        width: 16px;
        height: 16px;
        font-size: 16px;
        color: #595959;
        cursor: pointer;
        &:hover{
          color: #2475FC;
        }
      }

      .show-file{
        color: #2475FC;
      }

      .hide-file{
        color: #595959;
      }

      .file-number{
        position: absolute;
        left: 16px;
        top: -8px;
        width: 16px;
        height: 16px;
        border-radius: 50%;
        background: #f00;
        color: #fff;
        font-size: 12px;
        font-weight: 400;
        display: flex;
        align-items: center;
        justify-content: center;

        &.big{
          width: auto;
          padding: 0 4px;
          border-radius: 12px;
        }
      }
    }
  }
}
</style>

<template>
  <div class="message-input-box" :class="{ 'is-focus': isFocus }">
    <FileToolbar :file-list="props.fileList" @delete="deleteFile" v-if="props.fileList.length > 0 && showFiletoolbar" />
    <div class="message-input-body">
      <div class="message-input">
        <AutoSizeTextarea
          :value="value"
          @change="onChange"
          @focus="onFocus"
          @blur="onBlur"
          @enter="sendMessage"
        />
      </div>
      <div class="message-action">
        <div class="file-action" v-if="props.showUpload">
          <span class="file-number" :class="{ big: fileList.length > 9 }" v-if="fileList.length > 0">{{ fileList.length }}</span>
          <svg-icon class="action-btn select-file" name="circularNeedle" @click="openFileDialog"></svg-icon>
          <i class="line"></i>
          <Tippy :content="showFiletoolbar ? t('msg_hide_images') : t('msg_show_images')" placement="top">
            <svg-icon class="action-btn show-file" name="eye-open" v-if="showFiletoolbar" @click="showFiletoolbar = false"></svg-icon>
            <svg-icon class="action-btn hide-file" name="eye-close" v-else @click="showFiletoolbar = true"></svg-icon>
          </Tippy>
        </div>

        <button class="send-btn" @click="sendMessage" :disabled="disabled">
          <svg-icon name="paper-airplane-new-active" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, toRefs, computed } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useUpload } from '@/hooks/web/useUpload.js'
import { useChatStore } from '@/stores/modules/chat'
import AutoSizeTextarea from './auto-size-textarea.vue'
import FileToolbar from './file-toolbar.vue'
import { Tippy } from 'vue-tippy'

const { t } = useI18n('views.chat.components.message-input')


const emit = defineEmits(['update:value', 'send', 'update:fileList'])

const props = defineProps({
  value: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  },
  fileList: {
    type: Array,
    default: () => []
  },
  showUpload: {
    type: Boolean,
    default: false
  },
})

const chatStore = useChatStore()

const { robot } = chatStore

const isFocus = ref(false)
const { fileList } = toRefs(props)
const showFiletoolbar = ref(true)

const { openFileDialog } = useUpload({
  limit: 10,
  maxSize: 10,
  category: 'chat_image',
  fileList: fileList,
  multiple: true,
  accept: 'image/bmp,image/jpeg,image/png,image/tiff,image/heic,image/gif,image/webp',
  extraData: {
    robot_key: robot.robot_key,
    openid: robot.openid
  }
})

const deleteFile = (index) => {
  const newFileList = props.fileList.filter((_, i) => i !== index);
  emit('update:fileList', newFileList);
}

const disabled = computed(() => {
  if(props.fileList.length > 0){
    return false
  }

  return props.loading || props.value.trim().length === 0
})

const onChange = (val: string) => {
  emit('update:value', val)
}

const sendMessage = () => {
  emit('send', props.value)
}

const onFocus = () => {
  isFocus.value = true
}

const onBlur = () => {
  isFocus.value = false
}
</script>
