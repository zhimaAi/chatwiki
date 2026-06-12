<template>
  <cu-modal
    class="add-tool-modal"
    :width="948"
    :open="visible"
    :closable="false"
    @cancel="handleClose"
  >
    <div class="modal-shell">
      <div class="modal-hero">
        <img class="hero-bg" src="@/assets/img/clawbot/gd-hbimg.png" alt="" />
        <span class="modal-close" @click="handleClose"><CloseOutlined /></span>
        <div class="modal-title-box">
          <div class="modal-title">{{ t('title_add_tool') }}</div>
          <div class="modal-subtitle">{{ t('subtitle') }}</div>
        </div>
      </div>

      <div class="modal-body">
        <div class="modal-info-tip">
          <InfoCircleOutlined class="tip-icon" />
          <span>{{ t('info_tip') }}</span>
        </div>

        <div class="modal-actions">
          <a-button type="primary" ghost class="create-btn" @click="handleCreateNew">
            <template #icon>
              <PlusOutlined />
            </template>
            {{ t('btn_create_tool') }}
          </a-button>
          <a-button class="refresh-btn" @click="handleRefresh">
            <template #icon>
              <ReloadOutlined />
            </template>
            {{ t('btn_refresh') }}
          </a-button>
        </div>

        <div v-if="loading" class="loading-wrap">
          <a-spin />
        </div>
        <div v-else-if="toolList.length === 0" class="empty-wrap">
          <a-empty :description="t('empty_data')" />
        </div>
        <div v-else class="item-grid">
          <div
            v-for="item in toolList"
            :key="item.id"
            :class="['item-card', { selected: isSelected(item.id) }]"
            @click="toggleSelect(item.id)"
          >
            <div class="card-header">
              <a-checkbox
                :checked="isSelected(item.id)"
                @click.stop
                @change="toggleSelect(item.id)"
              />
              <div class="card-main">
                <div class="card-title">
                  <span class="title-text">{{ item.name }}</span>
                  <span v-if="item.unpublished" class="unpublished-tag">{{ t('tag_unpublished') }}</span>
                </div>
                <div class="card-desc">{{ item.desc || '—' }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="modal-footer">
        <a-button class="footer-btn cancel-btn" @click="handleClose">{{ t('btn_cancel') }}</a-button>
        <a-button
          type="primary"
          class="footer-btn confirm-btn"
          :loading="submitLoading"
          :disabled="selectedIds.length === 0"
          @click="handleConfirm"
        >
          {{ confirmButtonText }}
        </a-button>
      </div>
    </div>
  </cu-modal>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { InfoCircleOutlined, PlusOutlined, ReloadOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { CloseOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { getRobotList, relationWorkFlow } from '@/api/robot/index'
import CuModal from '@/components/common/cu-modal.vue'

const { t } = useI18n('views.clawbot.skills.components.AddToolModal')
const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  robotId: {
    type: [String, Number],
    default: ''
  },
  workFlowIds: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:visible', 'confirm'])
const router = useRouter()

const submitLoading = ref(false)
const selectedIds = ref([])
const loading = ref(false)
const toolList = ref([])

const confirmButtonText = computed(() => {
  return selectedIds.value.length ? t('btn_confirm_with_count', { count: selectedIds.value.length }) : t('btn_confirm')
})

const fetchToolList = async () => {
  loading.value = true
  try {
    const res = await getRobotList({ application_type: 1 })
    if (res && res.res === 0) {
      const list = res.data || []
      toolList.value = list.map((item) => ({
        id: item.id,
        name: item.robot_name,
        desc: item.robot_intro || '',
        unpublished: item.has_published != 1
      }))
    } else {
      message.error(res?.msg || t('msg_fetch_tools_failed'))
    }
  } catch (err) {
    console.error('获取工具列表失败', err)
    message.error(t('msg_fetch_tools_failed'))
  } finally {
    loading.value = false
  }
}

watch(() => props.visible, (val) => {
  if (val) {
    fetchToolList().then(() => {
      selectedIds.value = props.workFlowIds ? props.workFlowIds.split(',').filter(Boolean) : []
    })
  }
})

const isSelected = (id) => selectedIds.value.includes(id)

const toggleSelect = (id) => {
  const index = selectedIds.value.indexOf(id)
  if (index >= 0) {
    selectedIds.value.splice(index, 1)
  } else {
    selectedIds.value.push(id)
  }
}

const handleRefresh = () => {
  fetchToolList()
}

const handleCreateNew = () => {
  handleClose()
  router.push('/robot/list')
}

const handleConfirm = () => {
  if (selectedIds.value.length === 0) {
    return
  }

  submitLoading.value = true
  const workFlowIds = selectedIds.value.join(',')

  relationWorkFlow({
    id: props.robotId,
    work_flow_ids: workFlowIds
  }).then((res) => {
    if (res && res.res === 0) {
      message.success(t('msg_add_success'))
      emit('confirm', {
        type: 'tool',
        selectedIds: [...selectedIds.value],
        work_flow_ids: workFlowIds
      })
      handleClose()
    } else {
      message.error(res?.msg || t('msg_add_failed'))
    }
  }).catch((err) => {
    console.error('添加工具失败', err)
    message.error(t('msg_add_failed'))
  }).finally(() => {
    submitLoading.value = false
  })
}

const handleClose = () => {
  selectedIds.value = []
  emit('update:visible', false)
}
</script>

<style lang="less" scoped>
.modal-shell {
  border-radius: 16px;
}

.modal-hero {
  height: 118px;
  position: relative;
  .modal-close {
    position: absolute;
    top: 8px;
    right: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    cursor: pointer;
    color: #595959;
    z-index: 20;
  }
  .hero-bg {
    display: block;
    position: absolute;
    bottom:  -12px;
    left: 0;
    right: 0;
    width: 100%;
    height: 150px;
    z-index: 10;
  }
  .modal-title-box{
    position: relative;
    height: 100%;
    width: 100%;
    z-index: 11;
    padding: 20px 24px 0 24px;
  }
  .modal-title { 
    line-height: 28px;
    margin-bottom: 4px;
    font-weight: 600;
    font-size: 20px;
    color: #000;
  }
  .modal-subtitle { 
    line-height: 22px;
    font-weight: 400;
    font-size: 14px;
    color: #7A8699;
  }
}


.modal-body,
.modal-footer {
  position: relative;
  z-index: 1;
  background: #fff;
}

.modal-body {
  padding: 20px 24px 18px;
}

.modal-info-tip {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 9px 16px;
  border: 1px solid #99bffd;
  border-radius: 8px;
  background: #e9f1fe;
  color: #3a4559;
  font-size: 14px;
  line-height: 22px;

  .tip-icon {
    color: #2475fc;
    font-size: 14px;
    flex-shrink: 0;
  }
}

.modal-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-top: 16px;
}

.refresh-btn,
.create-btn {
  height: 36px;
  padding: 0 16px;
  border-radius: 8px;
  font-size: 14px;
  line-height: 22px;
}

.refresh-btn {
  color: #595959;
  border-color: #d9d9d9;
}

.create-btn {
  color: #2475fc;
  border-color: #2475fc;
}

.loading-wrap,
.empty-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 248px;
}

.item-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin-top: 10px;
  max-height: 292px;
  overflow-y: auto;
  padding-right: 2px;
}

.item-card {
  padding: 13px 17px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: #fff;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    border-color: #2475fc;
  }

  &.selected {
    padding: 12px 16px;
    border: 2px solid #2475fc;
    background: #f5f9ff;
  }
}

.card-header {
  display: flex;
  align-items: flex-start;
  gap: 10px;

  :deep(.ant-checkbox-wrapper) {
    line-height: 1;
  }

  :deep(.ant-checkbox) {
    top: 2px;
  }
}

.card-main {
  flex: 1;
  min-width: 0;
}

.card-title {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
}

.title-text {
  color: #262626;
  font-size: 14px;
  font-weight: 600;
  line-height: 22px;
}

.unpublished-tag {
  color: #7a8699;
  font-size: 12px;
  line-height: 20px;
}

.card-desc {
  margin-top: 6px;
  color: #7a8699;
  font-size: 12px;
  line-height: 20px;
}

.modal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  padding: 10px 24px;
  border-radius: 0 0 16px 16px;
  border-top: 1px solid #f0f0f0;
  background: #fff;
}

.footer-btn {
  min-width: 80px;
  height: 32px;
  padding: 0 16px;
  border-radius: 6px;
  font-size: 14px;
  line-height: 22px;
}

.cancel-btn {
  border-color: #d9d9d9;
  color: #595959;
}

.confirm-btn {
  border: none;
  box-shadow: none;
  background: #2475fc;
}

@media (max-width: 960px) {
  .item-grid {
    grid-template-columns: 1fr;
  }
}
</style>
