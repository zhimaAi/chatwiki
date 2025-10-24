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
    top: calc(50% - 8px);
  }
}
</style>

<template>
  <edit-box class="setting-box" title="聊天缓存" icon-name="chat-cache">
    <template #extra>
      <span></span>
    </template>
    <div class="robot-info-box">
      开启后，会缓存用户的问题和大模型的答案，如果之后有用户提出相同的问题，会直接回复缓存的答案。缓存有效期为
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
            <a-select-option value="h">小时</a-select-option>
            <a-select-option value="m">分钟</a-select-option>
          </a-select>
        </a-input-group>
      </div>
    </div>

    <a-switch
      @change="handleEdit"
      class="switch-item"
      :checkedValue="1"
      :unCheckedValue="0"
      v-model:checked="formState.cache_switch"
      checked-children="开"
      un-checked-children="关"
    />
  </edit-box>
</template>
<script setup>
import { ref, reactive, inject, toRaw, watchEffect, watch } from 'vue'
import EditBox from './edit-box.vue'

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
watch(
  () => robotInfo.cache_config,
  () => {
    let cache_config = robotInfo.cache_config || {}
    cache_config = JSON.parse(JSON.stringify(cache_config))
    console.log(robotInfo.cache_config, '==')
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
