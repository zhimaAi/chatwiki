<style lang="less" scoped>
.setting-box {
  position: relative;
  .robot-info-box {
    padding-right: 30px;
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 8px;
    .set-block {
    }
  }
  .switch-item {
    position: absolute;
    right: 16px;
    top: 16px;
    display: flex;
    align-items: center;
    gap: 8px;
  }
}
</style>

<template>
  <edit-box class="setting-box" :title="t('title_chat_cache')" icon-name="chat-cache">
    <template #extra>
      <span></span>
    </template>
    <div class="robot-info-box">
      {{ t('msg_cache_description') }}
      <div class="set-block">
        <a-input-group compact>
          <a-input-number
            v-if="state.unix == 'h'"
            @blur="onSave"
            :precision="0"
            v-model:value="state.hValue"
            :min="1"
            :max="72"
            style="width: 70px"
          />
          <a-input-number
            v-else
            @blur="onSave"
            :precision="0"
            v-model:value="state.mValue"
            :min="1"
            :max="4320"
            style="width: 70px"
          />
          <a-select v-model:value="state.unix" style="width: 70px" @change="handleChangeUnix">
            <a-select-option value="h">{{ t('unit_hour') }}</a-select-option>
            <a-select-option value="m">{{ t('unit_minute') }}</a-select-option>
          </a-select>
        </a-input-group>
      </div>
    </div>
    <div class="switch-item">
      <a-button @click="handleCleanCache" size="small">{{ t('btn_clean_cache') }}</a-button>
      <a-switch
        @change="handleEdit"
        :checkedValue="1"
        :unCheckedValue="0"
        v-model:checked="formState.cache_switch"
        :checked-children="t('switch_on')"
        :un-checked-children="t('switch_off')"
      />
    </div>
  </edit-box>
</template>
<script setup>
import { ref, reactive, inject, toRaw, watchEffect, watch, createVNode } from 'vue'
import EditBox from './edit-box.vue'
import { cleanRobotChatCache } from '@/api/robot/index'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { message, Modal } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.basic-config.components.chat-cache')
const { robotInfo, updateRobotInfo } = inject('robotInfo')

const formState = reactive({
  valid_time: '',
  cache_switch: 0
})

const state = reactive({
  unix: 'h',
  // 小时数
  hValue: 24,
  mValue: 1440
})

const onSave = () => {
  if (state.unix === 'h') {
    formState.valid_time = state.hValue * 3600
  } else {
    formState.valid_time = state.mValue * 60
  }
  updateRobotInfo({
    cache_config: {
      ...formState
    }
  })
}

const handleEdit = () => {
  onSave()
}

const handleChangeUnix = () => {
  if (state.unix == 'h') {
    state.hValue = parseInt(state.mValue / 60)
  } else {
    state.mValue = parseInt(state.hValue * 60)
  }
  // onSave()
}

const handleCleanCache = () => {
  Modal.confirm({
    title: t('msg_confirm_clear_cache'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('msg_confirm_clear_cache_content'),
    onOk() {
      cleanRobotChatCache({
        id: robotInfo.id,
        robot_key: robotInfo.robot_key
      }).then((res) => {
        message.success(t('msg_clear_success'))
      })
    }
  })
}
watch(
  () => robotInfo.cache_config,
  () => {
    let cache_config = robotInfo.cache_config || {}
    cache_config = JSON.parse(JSON.stringify(cache_config))

    formState.cache_switch = cache_config.cache_switch
    formState.valid_time = cache_config.valid_time
    // 判断 cache_config.valid_time 能不能整除24 能的话 state.unix = h 否则 state.unix = m
    if (cache_config.valid_time % 3600 == 0) {
      state.unix = 'h'
      state.hValue = cache_config.valid_time / 3600
    } else {
      state.unix = 'm'
      state.mValue = parseInt(cache_config.valid_time / 60) || 1440
    }
  },
  { immediate: true, deep: true } // immediate: true 表示立即执行一次
)
watchEffect(() => {})
</script>
