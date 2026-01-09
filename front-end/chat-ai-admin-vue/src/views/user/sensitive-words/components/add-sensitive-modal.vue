<template>
  <div>
    <a-modal v-model:open="open" :title="modalTitle" @ok="handleOk" :width="620">
      <div class="form-modal-box">
        <div class="form-item">
          <div class="form-label">{{ t('sensitive_words_label') }}<span>{{ t('sensitive_words_tip') }}</span></div>
          <div class="form-content">
            <a-textarea v-model:value="words" :placeholder="t('placeholder')" style="height: 160px" />
          </div>
        </div>
        <div class="form-item">
          <div class="form-label">{{ t('effective_robots') }}</div>
          <div class="form-content">
            <a-radio-group v-model:value="trigger_type">
              <a-radio :value="0">{{ t('all_robots') }}</a-radio>
              <a-radio :value="1">{{ t('specific_robots') }}</a-radio>
            </a-radio-group>
            <div class="robot-list" v-if="trigger_type == 1">
              <div class="checked-item" v-for="item in props.robotList" :key="item.id">
                <a-checkbox
                  @click.stop="handleClickItem(item)"
                  :checked="robot_ids.includes(item.id)"
                  ><div class="robot-name-text">{{ item.robot_name }}</div></a-checkbox
                >
              </div>
            </div>
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { saveSensitiveWords } from '@/api/robot'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.sensitive-words.components.add-sensitive-modal')

const emit = defineEmits(['ok'])

const props = defineProps({
  robotList: {
    type: Array,
    default: () => []
  }
})

const open = ref(false)
const modalTitle = ref(t('title_add'))

const id = ref('')
const words = ref('')
const trigger_type = ref(0)
const robot_ids = ref([])
const show = (data) => {
  open.value = true
  modalTitle.value = data.id ? t('title_edit') : t('title_add')
  id.value = data.id || ''
  words.value = data.words || ''
  trigger_type.value = +data.trigger_type || 0
  robot_ids.value = data.robot_ids || []
}

const handleClickItem = (item) => {
  if (robot_ids.value.includes(item.id)) {
    robot_ids.value = robot_ids.value.filter((id) => id !== item.id)
  } else {
    robot_ids.value.push(item.id)
  }
}

const handleOk = () => {
  if (!words.value) {
    return message.error(t('error_no_words'))
  }
  if (trigger_type.value == 1 && robot_ids.value.length == 0) {
    return message.error(t('error_no_robot'))
  }
  saveSensitiveWords({
    id: id.value,
    words: words.value,
    trigger_type: trigger_type.value,
    robot_ids: robot_ids.value.join(',')
  }).then((res) => {
    message.success(id.value ? t('success_edit') : t('success_add'))
    open.value = false
    emit('ok')
  })
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.form-modal-box {
  .form-item {
    margin-top: 16px;
  }
  .form-label {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
    span {
      color: #8c8c8c;
    }
  }
}
.robot-list {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  margin-top: 8px;

  .checked-item {
    width: 33.3%;
    height: 32px;
    display: flex;
    align-items: center;
    padding: 0 8px;
    cursor: pointer;
    transition: all 0.2s ease;
    &:hover {
      background: #e4e6eb;
      border-radius: 4px;
    }
    .robot-name-text {
      width: 140px;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }
  }
}
</style>
