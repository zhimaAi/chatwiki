<template>
  <div>
    <a-drawer v-model:open="open" title="发布详情" placement="right" :width="400">
      <div class="version-detail">
        <div class="version-list" @click="handlePreviewVersion('')">
          <div class="version-header">
            <div class="version-title">当前版本</div>
          </div>
          <div class="version-desc" v-if="draftSaveTime && draftSaveTime.time">
            最近保存于{{ draftSaveTime.time }}
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
              <span class="time-text">{{ formatTime(item.create_time) }}</span>
            </div>
            <a-dropdown>
              <div class="hover-btn-box" @click.stop="">
                <EllipsisOutlined />
              </div>

              <template #overlay>
                <a-menu>
                  <a-menu-item @click="setVersion(item)"> 恢复到此版本 </a-menu-item>
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
import { workFlowVersions, workFlowVersionDetail, getNodeList } from '@/api/robot/index'
import { message, Modal } from 'ant-design-vue'
import { EllipsisOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import { useRobotStore } from '@/stores/modules/robot'
import { useRoute } from 'vue-router'
const query = useRoute().query
const open = ref(false)

const emit = defineEmits(['setVersion', 'preview'])

const robotStore = useRobotStore()
const draftSaveTime = computed(() => {
  return robotStore.draftSaveTime
})

const list = ref([])
const showDrawer = () => {
  open.value = true
  getDetailList()
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
    title: `确定要恢复到v${item.version}版吗`,
    icon: null,
    content: h(
      'div',
      { style: 'color:red;' },
      '恢复后，将覆盖当前草稿内容。如需要将此版本作为发布版本，需要您手动点击发布'
    ),
    onOk() {
      workFlowVersionDetail({
        robot_key: query.robot_key,
        version_id: item.version_id
      }).then((res) => {
        emit('setVersion', res.data)
        open.value = false
        message.success('设置成功')
      })
    },
    onCancel() {}
  })
}

const handlePreviewVersion = (version) => {
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

function formatTime(time) {
  if (time <= 0) {
    return '--'
  }
  return dayjs(time * 1000).format('YYYY-MM-DD HH:mm:ss')
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
  gap: 16px;
  .version-list {
    padding: 8px 16px;
    border-radius: 4px;
    background-color: #f2f2f2;
    cursor: pointer;
    .version-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      color: #000;
      font-weight: 600;
      .time-text {
        color: #8c8c8c;
        margin-left: 12px;
      }
    }
  }
}
</style>
