<template>
  <div>
    <a-modal
      v-model:open="open"
      title="发布"
      ok-text="立即发布"
      @ok="handleOk"
      :confirmLoading="saveLoading"
      :width="600"
    >
      <div class="form-box">
        <div class="form-item">
          <div class="form-label">版本号：</div>
          <div class="form-content">
            <div class="version-input-box">
              v
              <a-input-number v-model:value="value1" :precision="0" :min="0" :max="1000" />
              <span>.</span>
              <a-input-number v-model:value="value2" :precision="0" :min="0" :max="1000" />
              <span>.</span>
              <a-input-number v-model:value="value3" :precision="0" :min="0" :max="1000" />
            </div>
          </div>
        </div>
        <div class="form-item">
          <div class="form-label">版本描述：</div>
          <div class="form-content">
            <a-textarea
              v-model:value="version_desc"
              placeholder="请描述本次发布版本的更新内容"
              :auto-size="{ minRows: 4, maxRows: 7 }"
            />
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { workFlowNextVersion, workFlowPublishVersion } from '@/api/robot/index'
import { message } from 'ant-design-vue'
import { useRoute } from 'vue-router'
import { useRobotStore } from '@/stores/modules/robot'

const query = useRoute().query
const open = ref(false)
let node_list = ''
const robotStore = useRobotStore()
const value1 = ref('')
const value2 = ref('')
const value3 = ref('')

const version_desc = ref('')

const saveLoading = ref(false)

const emit = defineEmits(['ok'])

const show = (list) => {
  getVersion()
  version_desc.value = ''
  node_list = list
  open.value = true
}
const handleOk = () => {
  saveLoading.value = true
  workFlowPublishVersion({
    robot_key: query.robot_key,
    node_list,
    version: `${value1.value}.${value2.value}.${value3.value}`,
    version_desc: version_desc.value,
    draft_save_type: 'automatic',
    draft_save_time: +robotStore.robotInfo.draft_save_time || 0
  })
    .then(async(res) => {
      // 刷新并同步最新草稿时间戳
      await robotStore.getRobot(query.id)
      const ts = +robotStore.robotInfo.draft_save_time || 0
      robotStore.setDrafSaveTime({
        draft_save_type: 'automatic',
        draft_save_time: ts
      })
      open.value = false
      if (res && res.res == 0) {
        message.success('发布成功')
        emit('ok')
      }
    })
    .finally(() => {
      saveLoading.value = false
    })
}

const getVersion = () => {
  workFlowNextVersion({
    robot_key: query.robot_key
  }).then((res) => {
    let version_params = res.data.version_params || []
    value1.value = version_params[0]
    value2.value = version_params[1]
    value3.value = version_params[2]
  })
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.form-box {
  margin: 24px 0;
  display: flex;
  flex-direction: column;
  gap: 16px;
  .form-item {
    color: #262626;
    font-size: 14px;
    .form-label {
      margin-bottom: 6px;
    }
  }
  .version-input-box {
    display: flex;
    gap: 4px;
    align-items: self-end;
  }
}
</style>
