<template>
  <div>
    <a-modal
      v-model:open="open"
      :title="t('title_publish')"
      :ok-text="t('btn_publish_now')"
      @ok="handleOk"
      :confirmLoading="saveLoading"
      :width="600"
    >
      <div class="form-box">
        <div class="form-item">
          <div class="form-label">{{ t('label_version_number') }}</div>
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
          <div class="form-label">{{ t('label_version_description') }}</div>
          <div class="form-content">
            <a-textarea
              v-model:value="version_desc"
              :placeholder="t('ph_version_desc')"
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
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.version-model')

const query = useRoute().query
const open = ref(false)
let node_list = ''
const robotStore = useRobotStore()
const value1 = ref('')
const value2 = ref('')
const value3 = ref('')

const version_desc = ref('')

const saveLoading = ref(false)

const emit = defineEmits(['ok', 'handleOpenErrorNode'])

const show = async (list) => {
  let res = await getVersion()
  if(res.res == 1){
    if(res.data && res.data.err_node_key){
      emit('handleOpenErrorNode', res.data.err_node_key)
      return
    }
  }
  version_desc.value = ''
  node_list = list
  open.value = true
}

// 唯一标识存取
const UNI_STORAGE_KEY = 'wf_uni_identifier'
const getUniIdentifier = () => {
  try {
    let id = localStorage.getItem(UNI_STORAGE_KEY)
    if (!id) {
      id = `${Date.now()}_${Math.random().toString(36).slice(2, 10)}`
      localStorage.setItem(UNI_STORAGE_KEY, id)
    }
    return id
  } catch (e) {
    return `${Date.now()}_${Math.random().toString(36).slice(2, 10)}`
  }
}

// 组装 user_agent（操作系统 + 浏览器信息）
const buildUserAgent = () => {
  try {
    const ua = navigator.userAgent || ''
    const platform = navigator.platform || ''
    let os = 'Unknown'
    if (/Windows/i.test(ua)) os = 'Windows'
    else if (/Macintosh|Mac OS X/i.test(ua)) os = 'MacOS'
    else if (/Linux/i.test(ua)) os = 'Linux'
    else if (/Android/i.test(ua)) os = 'Android'
    else if (/iPhone|iPad|iPod/i.test(ua)) os = 'iOS'

    let browser = 'Unknown'
    let version = ''
    const m = ua.match(/Edg\/([\d\.]+)/) || ua.match(/Chrome\/([\d\.]+)/) || ua.match(/Firefox\/([\d\.]+)/) || ua.match(/Version\/([\d\.]+).*Safari/)
    if (m) {
      if (ua.includes('Edg/')) browser = 'Edge'
      else if (ua.includes('Chrome/')) browser = 'Chrome'
      else if (ua.includes('Firefox/')) browser = 'Firefox'
      else if (ua.includes('Safari') && !ua.includes('Chrome')) browser = 'Safari'
      version = m[1]
    }
    return `platform=${platform}; os=${os}; browser=${browser}/${version}; ua=${ua}`
  } catch (e) {
    return 'ua=unknown'
  }
}

const handleOk = () => {
  saveLoading.value = true
  workFlowPublishVersion({
    robot_key: query.robot_key,
    node_list,
    version: `${value1.value}.${value2.value}.${value3.value}`,
    version_desc: version_desc.value,
    uni_identifier: getUniIdentifier(),
    user_agent: buildUserAgent(),
    draft_save_type: 'automatic',
    draft_save_time: +robotStore.robotInfo.draft_save_time || 0
  })
    .then(async(res) => {
      // 刷新并同步最新草稿时间戳
      await robotStore.getRobot(query.id)
      robotStore.setDrafSaveTime({
        draft_save_type: 'automatic',
        draft_save_time: +robotStore.robotInfo.draft_save_time || 0,
        uni_identifier: getUniIdentifier(),
        user_agent: buildUserAgent()
      })
      open.value = false
      if (res && res.res == 0) {
        message.success(t('msg_publish_success'))
        emit('ok')
      }
    }).catch((res)=>{
      if(res.data && res.data.err_node_key){
        emit('handleOpenErrorNode', res.data.err_node_key)
      }
      open.value = false
    })
    .finally(() => {
      saveLoading.value = false
    })
}

const getVersion = async () => {
  return workFlowNextVersion({
    robot_key: query.robot_key
  }).then((res) => {
    let version_params = res.data.version_params || []
    value1.value = version_params[0]
    value2.value = version_params[1]
    value3.value = version_params[2]
    return res
  }).catch((res)=>{
    return res
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
