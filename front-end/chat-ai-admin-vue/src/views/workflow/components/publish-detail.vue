<template>
  <div>
    <a-drawer v-model:open="open" :title="t('title_publish_detail')" placement="right" :width="400" :bodyStyle="{ padding: '12px 24px' }">
      <div class="version-detail">
        <div class="version-list" @click="handlePreviewVersion('')" v-if="robotInfo.draft_save_time && robotInfo.draft_save_time != 0">
          <div class="version-header">
            <div class="version-title">{{ t('label_current_draft') }}</div>
          </div>
          <div class="version-desc flex-center">
            {{ t('label_saved_at') }}{{ formatTime(robotInfo.draft_save_time, 'MM/DD HH:mm:ss') }}
            <div class="version-file-box">
              <a-popover placement="topRight" :overlay-style="{
                'max-width': '372px'
              }">
                <template #content>
                  <div class="version-text">
                    IP：{{ currentIp || '--' }}
                  </div>
                  <div class="version-text">
                    User Agent：{{ currentUA || '--' }}
                  </div>
                </template>
                <svg-icon
                  class="file-icon"
                  name="file-icon"
                  style="font-size: 14px; color: #333"
                ></svg-icon>
              </a-popover>
            </div>
          </div>
        </div>
        <div
          class="version-list"
          @click="handlePreviewVersion(item.version_id)"
          v-for="(item, index) in list"
          :key="item.version"
        >
          <div class="version-header">
            <div class="version-title">
              v{{ item.version }}
              <span class="time-text flex-center">
                {{ formatTime(item.create_time) }}
                <div class="version-file-box">
                  <a-popover placement="topRight" :overlay-style="{
                    'max-width': '372px'
                  }">
                    <template #content>
                      <div class="version-text">
                        IP：{{ item.last_edit_ip || '--' }}
                      </div>
                      <div class="version-text">
                        User Agent：{{ item.last_edit_user_agent || '--' }}
                      </div>
                    </template>
                    <svg-icon
                      class="file-icon"
                      name="file-icon"
                      style="font-size: 14px; color: #333"
                    ></svg-icon>
                  </a-popover>
                </div>
              </span>
            </div>
            <a-dropdown>
              <div class="hover-btn-box" @click.stop="">
                <EllipsisOutlined />
              </div>

              <template #overlay>
                <a-menu>
                  <a-menu-item :disabled="props.isLockedByOther" @click="setVersion(item)"> {{ t('btn_restore_version') }} </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
          <div class="version-desc">
            {{ item.version_desc }}
          </div>
        </div>
      </div>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, h, computed } from 'vue'
import { workFlowVersions, workFlowVersionDetail, getNodeList, getDraftKey } from '@/api/robot/index'
import { message, Modal } from 'ant-design-vue'
import { EllipsisOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import { useRobotStore } from '@/stores/modules/robot'
import { useRoute } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.publish-detail')

const query = useRoute().query
const open = ref(false)
const currentIp = ref('')
const currentUA = ref('')

const emit = defineEmits(['setVersion', 'preview'])
const props = defineProps({
  isLockedByOther: { type: Boolean, default: false }
})
const robotStore = useRobotStore()
const robotInfo = computed(() => {
  return robotStore.robotInfo
})

const list = ref([])
const showDrawer = () => {
  open.value = true
  getDetailList()
  // 获取当前客户端信息
  getDraftKey({
    robot_key: query.robot_key
  }).then((res) => {
    // const res = {"msg":"success","res":0,"data":{"is_self":true,"lock_res":true,"lock_ttl":955,"remote_addr":"171.83.17.34","robot_key":"yw5BnxX80G","staff_id":3432,"user_agent":"Mozilla\/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit\/537.36 (KHTML, like Gecko) Chrome\/142.0.0.0 Safari\/537.36"}}
    const data = res?.data || {}
    currentIp.value = data.remote_addr || ''
    currentUA.value = data.user_agent || ''
  })
}

const getDetailList = () => {
  workFlowVersions({
    robot_key: query.robot_key
  }).then((res) => {
    list.value = res.data.versions || []
  })
}

const setVersion = (item) => {
  Modal.confirm({
    title: t('msg_confirm_restore_version', { version: item.version }),
    icon: null,
    content: h(
      'div',
      { style: 'color:red;' },
      t('msg_restore_warning')
    ),
    onOk() {
      workFlowVersionDetail({
        robot_key: query.robot_key,
        version_id: item.version_id
      }).then((res) => {
        emit('setVersion', res.data)
        open.value = false
        if (res.res == 0) {
          message.success(t('msg_set_success'))
        }
      })
    },
    onCancel() {}
  })
}

const handlePreviewVersion = (version) => {
  // if (props.isLockedByOther) {
  //   message.warning('当前已有其他用户在编辑中，无法预览')
  //   return
  // }
  if (version == '') {
    getNodeList({
      robot_key: query.robot_key,
      data_type: 1
    }).then((res) => {
      emit('preview', res.data, version)
      open.value = false
    })
  } else {
    workFlowVersionDetail({
      robot_key: query.robot_key,
      version_id: version
    }).then((res) => {
      emit('preview', res.data, version)
      open.value = false
    })
  }
}

function formatTime(time, formatType = 'YYYY-MM-DD HH:mm:ss') {
  if (time <= 0) {
    return '--'
  }
  return dayjs(time * 1000).format(formatType)
}

defineExpose({
  showDrawer
})
</script>

<style lang="less" scoped>
.hover-btn-box {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 24px;
  height: 24px;
  border-radius: 6px;
  cursor: pointer;
  &:hover {
    background-color: #e4e6eb;
  }
}

.version-detail {
  display: flex;
  flex-direction: column;
  gap: 12px;
  .version-list {
    padding: 8px 16px;
    border-radius: 6px;
    background-color: #F2F4F7;
    cursor: pointer;
    .version-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      color: #000;
      font-weight: 600;
      margin-bottom: 8px;
      .time-text {
        color: #8c8c8c;
        margin-left: 12px;
      }
    }
  }
}

.flex-center {
  display: flex;
  align-items: center;
}

.version-file-box {
  margin-left: 8px;
}

.version-title {
  display: flex;
  align-items: center;
}

.version-desc {
  // 保持返回的格式
  white-space: pre-wrap;
  color: #595959;
}
</style>
