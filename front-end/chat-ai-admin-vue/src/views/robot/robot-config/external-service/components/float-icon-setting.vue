<style lang="less" scoped>
.float-icon-setting {
  .form-box {
    .form-item {
      display: flex;
      margin-bottom: 24px;
    }
    .form-item-label {
      width: 120px;
      line-height: 32px;
      margin-right: 4px;
      text-align: right;
    }
    .number-input {
      display: flex;
      align-items: center;
      .input {
        width: 120px;
      }
      .unit {
        font-size: 16px;
        margin-left: 8px;
      }
    }
    .display-types {
      display: flex;
      .display-type {
        position: relative;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        width: 100px;
        height: 100px;
        margin-right: 8px;
        font-size: 14px;
        border-radius: 2px 2px 2px 2px;
        border: 1px solid rgba(0, 0, 0, 0.06);
        background-color: #fff;
        overflow: hidden;
        cursor: pointer;
        &:hover,
        &.is-checked {
          border: 1px solid #2475fc;
        }
        &.is-checked::after {
          content: '';
          position: absolute;
          bottom: -1px;
          right: -1px;
          width: 18px;
          height: 18px;
          background: url('@/assets/svg/check-arrow-filled.svg') no-repeat center;
          border-radius: 0 0 0 2px;
        }

        .type-icon {
          display: flex;
          flex-direction: column;
          justify-content: center;
          align-items: center;
          width: 56px;
          height: 56px;
          img {
            width: auto;
            height: auto;
            max-width: 100%;
            max-height: 100%;
          }
        }
        .type-name {
          color: rgba(0, 0, 0, 0.65);
          margin-top: 4px;
        }
        .checked-icon {
          position: absolute;
          bottom: 0;
          right: 0;
        }
      }
    }
  }
}
.float-icon-setting.is-workspace {
  .form-item {
    align-items: flex-start;
    margin-bottom: 28px;

    &:last-child {
      margin-bottom: 0;
    }
  }
  .form-item-label {
    width: 112px;
    margin-right: 24px;
    text-align: left;
    color: var(--app-text-primary);
  }
  .form-content {
    flex: 1;
    min-width: 0;
  }
  .form-box .display-types {
    gap: 12px;
  }
  .form-box .display-types .display-type {
    width: 88px;
    height: 92px;
    margin-right: 0;
    border-radius: 8px;
    border-color: #d8d8d8;
    background: #fff;

    &:hover,
    &.is-checked {
      border-color: var(--app-primary);
    }
    &.is-checked {
      background: var(--app-selected);
    }
    &.is-checked::after {
      display: none;
    }
    &.is-checked .type-name {
      color: var(--app-primary);
    }
    .type-name {
      margin-top: 10px;
      color: #222;
      font-size: 14px;
      line-height: 22px;
    }
    .type-icon {
      width: 28px;
      height: 28px;
      overflow: hidden;
    }
    .type-icon-text {
      width: 56px;
      height: 28px;
    }
    .type-icon img {
      width: 56px;
      height: 56px;
      max-width: none;
      max-height: none;
      flex-shrink: 0;
    }
    .type-icon-custom img {
      width: 28px;
      height: 28px;
      object-fit: cover;
      border-radius: 50%;
    }
  }
  .button-text-input {
    width: 100%;
    max-width: 280px;
    margin-top: 12px;
  }
  .metric-row,
  .check-row {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 8px 16px;
    min-height: 32px;
  }
  .metric-row .unit {
    margin-left: 0;
    color: var(--app-text-tertiary);
    font-size: 12px;
  }
  .metric-row :deep(.ant-input-number) {
    width: 72px;
  }
}
</style>

<template>
  <CardBox :title="t('title_float_icon_setting')">
    <template #icon>
      <MessageOutlined name="quick-Instruction" style="font-size: 16px; color: #262626" />
    </template>
    <div class="float-icon-setting" :class="{ 'is-workspace': isWorkspace }">
      {{ props.title }}
      <div class="form-box">
        <div class="form-item">
          <div class="form-item-label">{{ isWorkspace ? t('label_display_type_plain') : t('label_display_type') }}</div>
          <div class="form-content">
            <div v-if="isWorkspace" class="display-types">
              <div
                class="display-type"
                :class="{ 'is-checked': formData.displayType == 1 }"
                @click="handleDisplayType(1)"
              >
                <div class="type-icon">
                  <img src="@/assets/img/robot/display_type_1.svg" alt="" />
                </div>
                <div class="type-name">{{ t('type_simple') }}</div>
              </div>
              <div
                class="display-type"
                :class="{ 'is-checked': formData.displayType == 2 }"
                @click="handleDisplayType(2)"
              >
                <div class="type-icon type-icon-text">
                  <img src="@/assets/img/robot/display_type_2.svg" alt="" />
                </div>
                <div class="type-name">{{ t('type_with_text') }}</div>
              </div>
              <div
                class="display-type"
                :class="{ 'is-checked': formData.displayType == 3 }"
                @click="handleDisplayType(3)"
              >
                <cu-upload @change="uploadButtonIcon">
                  <div class="type-icon type-icon-custom">
                    <img v-if="formData.buttonIcon" :src="formData.buttonIcon" alt="" />
                    <PlusOutlined v-else style="font-size: 24px; color: #222" />
                  </div>
                </cu-upload>
                <div class="type-name">{{ t('type_custom') }}</div>
              </div>
            </div>
            <a-input
              v-if="isWorkspace && formData.displayType == 2"
              v-model:value="formData.buttonText"
              class="button-text-input"
              :maxlength="15"
              :placeholder="t('ph_input_button_text')"
            />
            <div v-if="!isWorkspace" class="display-types">
              <div
                class="display-type"
                :class="{ 'is-checked': formData.displayType == 1 }"
                @click="handleDisplayType(1)"
              >
                <div class="type-icon">
                  <img src="@/assets/img/robot/display_type_1.svg" alt="" />
                </div>
                <div class="type-name">{{ t('type_simple') }}</div>
              </div>
              <div
                class="display-type"
                :class="{ 'is-checked': formData.displayType == 2 }"
                @click="handleDisplayType(2)"
              >
                <div class="type-icon">
                  <img src="@/assets/img/robot/display_type_2.svg" alt="" />
                </div>
                <div class="type-name">{{ t('type_with_text') }}</div>
              </div>

              <div
                class="display-type"
                :class="{ 'is-checked': formData.displayType == 3 }"
                @click="handleDisplayType(3)"
                v-if="formData.buttonIcon"
              >
                <div class="type-icon">
                  <img
                    :src="formData.buttonIcon"
                    alt=""
                    title=""
                  />
                </div>
                <div class="type-name">{{ t('type_custom') }}</div>
              </div>

              <div
                class="display-type"
                v-if="uploadBoxShow"
              >
                <cu-upload @change="uploadButtonIcon">
                  <div class="type-icon">
                    <PlusOutlined style="font-size: 20px"  />
                  </div>
                </cu-upload>

                <div class="type-name">{{ formData.buttonIcon ? t('label_upload_image') : t('type_custom') }}</div>
              </div>
            </div>
          </div>
        </div>
        <div class="form-item" v-if="!isWorkspace && formData.displayType == 2">
          <div class="form-item-label">{{ t('label_button_text') }}</div>
          <div class="form-content">
            <a-input
              :placeholder="t('ph_input_button_text')"
              :maxlength="15"
              style="width: 315px"
              v-model:value="formData.buttonText"
            />
          </div>
        </div>
        <div v-if="isWorkspace" class="form-item">
          <div class="form-item-label">{{ t('label_icon_position') }}</div>
          <div class="form-content">
            <div class="metric-row">
              <span>{{ t('label_right_plain') }}</span>
              <a-input-number v-model:value="formData.rightMargin" :min="0" :max="2000" :precision="0" />
              <span class="unit">px</span>
              <span>{{ t('label_bottom_plain') }}</span>
              <a-input-number v-model:value="formData.bottomMargin" :min="0" :max="2000" :precision="0" />
              <span class="unit">px</span>
            </div>
          </div>
        </div>
        <div v-if="isWorkspace" class="form-item">
          <div class="form-item-label">{{ t('label_icon_display') }}</div>
          <div class="form-content">
            <div class="check-row">
              <a-checkbox
                :checked="Number(formData.showUnreadCount) === 1"
                @change="(e) => (formData.showUnreadCount = e.target.checked ? 1 : 0)"
              >
                {{ t('label_unread_plain') }}
              </a-checkbox>
              <a-checkbox
                :checked="Number(formData.showNewMessageTip) === 1"
                @change="(e) => (formData.showNewMessageTip = e.target.checked ? 1 : 0)"
              >
                {{ t('label_notice_plain') }}
              </a-checkbox>
            </div>
          </div>
        </div>
        <div v-if="!isWorkspace" class="form-item">
          <div class="form-item-label">{{ t('label_bottom_margin') }}</div>
          <div class="form-content">
            <div class="number-input">
              <a-input-number class="input" v-model:value="formData.bottomMargin" />
              <span class="unit">PX</span>
            </div>
          </div>
        </div>
        <div v-if="!isWorkspace" class="form-item">
          <div class="form-item-label">{{ t('label_right_margin') }}</div>
          <div class="form-content">
            <div class="number-input">
              <a-input-number class="input" v-model:value="formData.rightMargin" />
              <span class="unit">PX</span>
            </div>
          </div>
        </div>
        <div v-if="!isWorkspace" class="form-item">
          <div class="form-item-label">{{ t('label_show_unread_count') }}</div>
          <div class="form-content">
            <a-switch
              v-model:checked="formData.showUnreadCount"
              :unCheckedValue="0"
              :checkedValue="1"
            />
          </div>
        </div>
        <div v-if="!isWorkspace" class="form-item">
          <div class="form-item-label">{{ t('label_show_new_message_tip') }}</div>
          <div class="form-content">
            <a-switch
              v-model:checked="formData.showNewMessageTip"
              :unCheckedValue="0"
              :checkedValue="1"
            />
          </div>
        </div>
        <div v-if="!isWorkspace" class="form-item">
          <div class="form-item-label"></div>
          <div class="form-content">
            <a-button type="primary" :loading="props.saveLoading" @click="handleSubmit"
              >{{ t('btn_save') }}</a-button
            >
          </div>
        </div>
      </div>
    </div>
  </CardBox>
</template>

<script setup>
import { uploadFile } from '@/api/app'
import { computed, nextTick, onMounted, reactive, toRaw, watch } from 'vue'
import { MessageOutlined, PlusOutlined } from '@ant-design/icons-vue'
import CardBox from './card-box.vue'
import CuUpload from '@/components/cu-upload/cu-upload.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.external-service.components.float-icon-setting')

// 定义组件属性
const props = defineProps({
  saveLoading: {
    type: Boolean,
    default: false
  },
  form: {
    type: Object,
    default: () => {}
  },
  appearance: {
    type: String,
    default: 'card'
  }
})
const isWorkspace = computed(() => props.appearance === 'workspace')

// 定义组件事件
const emit = defineEmits(['save'])

// 响应式数据
const formData = reactive({
  displayType: 1,
  buttonText: t('default_button_text'),
  buttonIcon: '',
  bottomMargin: 32,
  rightMargin: 32,
  showUnreadCount: 1,
  showNewMessageTip: 1
})

const handleDisplayType = (type) => {
  formData.displayType = type
}

const uploadBoxShow = computed(() => {
  if(!formData.buttonIcon){
    return true
  }else{
    return formData.displayType == 3
  }
})

const uploadButtonIcon = (files) => {
  let file = files[0]

  handleDisplayType(3)

  uploadFile({
    file: file,
    category: 'icon'
  }).then((res) => {
    formData.buttonIcon = res.data.link
  })
}
// 方法
const handleSubmit = () => {
  console.log(formData.rightMargin)
  if (formData.rightMargin == '' || formData.rightMargin == null) {
    formData.rightMargin = 0
  }

  if (formData.bottomMargin == '' || formData.rightMargin == null) {
    formData.bottomMargin = 0
  }

  emit('save', toRaw(formData))
}

let ignoreFloatSave = 0
let floatReady = false

watch(
  () => props.form,
  (newVal) => {
    ignoreFloatSave += 1
    Object.assign(formData, newVal?.floatBtn || {})
    nextTick(() => {
      ignoreFloatSave = Math.max(0, ignoreFloatSave - 1)
    })
  },
  {
    immediate: true
  }
)

watch(
  formData,
  () => {
    if (!floatReady || ignoreFloatSave > 0 || !isWorkspace.value) return
    emit('save', { ...toRaw(formData) })
  },
  { deep: true }
)

onMounted(() => {
  floatReady = true
})

defineExpose({
  handleSubmit
})
</script>
