<template>
  <cu-modal
    class="select-skill-modal"
    :open="visible"
    :width="948"
    :closable="false"
    :maskClosable="false"
    @cancel="handleClose"
  >
    <div class="modal-shell">
      <div class="modal-hero">
        <img class="hero-bg" src="@/assets/img/clawbot/gd-hbimg.png" alt="" />
        <span class="modal-close" @click="handleClose"><CloseOutlined /></span>
        <div class="modal-title-box">
          <div class="modal-title">{{ t('title_add_skill') }}</div>
          <div class="modal-subtitle">{{ t('subtitle_add_skill') }}</div>
        </div>
      </div>

      <div class="modal-body">
        <div class="select-skill-tip">
          <InfoCircleOutlined class="tip-icon" />
          <span>{{ t('select_skill_tip') }}</span>
        </div>

        <div class="select-skill-actions">
          <a-button type="primary" ghost class="create-btn" @click="handleCreateSkill">
            <template #icon>
              <PlusOutlined />
            </template>
            {{ t('btn_create_skill') }}
          </a-button>
          <a-button class="refresh-btn" :loading="loading" @click="loadSkillList">
            <template #icon>
              <ReloadOutlined />
            </template>
            {{ t('btn_refresh') }}
          </a-button>
        </div>

        <div v-if="loading" class="select-skill-loading">
          <a-spin />
        </div>
        <div v-else-if="skillList.length === 0" class="select-skill-empty">
          <a-empty :description="t('empty_skill')" />
        </div>
        <div v-else class="select-skill-grid">
          <div
            v-for="item in skillList"
            :key="item.id"
            class="select-skill-card"
            :class="{ selected: isSelected(item.skillId) }"
            @click="toggleSelect(item.skillId)"
          >
            <a-checkbox
              :checked="isSelected(item.skillId)"
              @click.stop
              @change="toggleSelect(item.skillId)"
            />
            <div class="select-skill-card-main">
              <div class="select-skill-card-title">{{ item.title }}</div>
              <a-tooltip :title="item.desc">
                <div class="select-skill-card-desc">{{ item.desc }}</div>
              </a-tooltip>
            </div>
          </div>
        </div>
      </div>

      <div class="modal-footer">
        <a-button class="footer-btn cancel-btn" @click="handleClose">{{ t('btn_cancel') }}</a-button>
        <a-button type="primary" class="footer-btn confirm-btn" :loading="submitting" @click="handleConfirm">
          {{ t('btn_confirm') }}
        </a-button>
      </div>
    </div>
  </cu-modal>
</template>

<script setup>
import { ref, watch } from 'vue'
import { CloseOutlined, InfoCircleOutlined, PlusOutlined, ReloadOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { getClawbotSkillList, saveClawbotRobotSkills } from '@/api/clawbot'
import CuModal from '@/components/common/cu-modal.vue'

const { t } = useI18n('views.clawbot.skills.index')

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  robotId: {
    type: [String, Number],
    default: ''
  },
  refreshKey: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['update:visible', 'confirm', 'create'])

const loading = ref(false)
const submitting = ref(false)
const skillList = ref([])
const selectedSkillIds = ref([])

watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      loadSkillList()
    }
  }
)

watch(
  () => props.refreshKey,
  () => {
    if (props.visible) {
      loadSkillList()
    }
  }
)

const loadSkillList = async () => {
  if (!props.robotId) {
    skillList.value = []
    selectedSkillIds.value = []
    return
  }

  loading.value = true
  try {
    const res = await getClawbotSkillList({
      id: props.robotId,
      source: 'mine'
    })
    if (res?.res === 0) {
      skillList.value = (res.data || []).map((item, index) => ({
        id: `${item.source_type || item.source || 'skill'}-${item.skill_id || 0}-${item.skill_name || index}`,
        skillId: String(item.skill_id || ''),
        title: item.remark_name || item.skill_name || '—',
        desc: item.intro || item.description || '—',
        raw: item
      })).filter((item) => item.skillId)
      selectedSkillIds.value = skillList.value
        .filter((item) => Number(item.raw?.is_selected) === 1)
        .map((item) => item.skillId)
    } else {
      message.error(res?.msg || t('msg_fetch_skill_failed'))
    }
  } catch (err) {
    console.error('获取 Skill 列表失败', err)
    message.error(err?.msg || t('msg_fetch_skill_failed'))
  } finally {
    loading.value = false
  }
}

const isSelected = (skillId) => selectedSkillIds.value.includes(String(skillId))

const toggleSelect = (skillId) => {
  const value = String(skillId)
  const index = selectedSkillIds.value.indexOf(value)
  if (index > -1) {
    selectedSkillIds.value.splice(index, 1)
  } else {
    selectedSkillIds.value.push(value)
  }
}

const handleCreateSkill = () => {
  emit('create')
}

const handleConfirm = async () => {
  if (!props.robotId) {
    return
  }

  submitting.value = true
  try {
    const res = await saveClawbotRobotSkills({
      id: props.robotId,
      skill_ids: selectedSkillIds.value.join(',')
    })
    if (res?.res === 0) {
      message.success(t('msg_save_success'))
      emit('confirm')
      handleClose()
    } else {
      message.error(res?.msg || t('msg_save_failed'))
    }
  } catch (err) {
    console.error('保存 Skill 选择失败', err)
    message.error(err?.msg || t('msg_save_failed'))
  } finally {
    submitting.value = false
  }
}

const handleClose = () => {
  emit('update:visible', false)
}
</script>

<style lang="less" scoped>
.modal-shell {
  border-radius: 16px;
}

.modal-hero {
  position: relative;
  height: 118px;

  .modal-close {
    position: absolute;
    top: 8px;
    right: 8px;
    z-index: 20;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    color: #595959;
    cursor: pointer;
  }

  .hero-bg {
    position: absolute;
    right: 0;
    bottom: -12px;
    left: 0;
    z-index: 10;
    display: block;
    width: 100%;
    height: 150px;
  }

  .modal-title-box {
    position: relative;
    z-index: 11;
    width: 100%;
    height: 100%;
    padding: 20px 24px 0;
  }

  .modal-title {
    margin-bottom: 4px;
    color: #000;
    font-size: 20px;
    font-weight: 600;
    line-height: 28px;
  }

  .modal-subtitle {
    color: #7a8699;
    font-size: 14px;
    font-weight: 400;
    line-height: 22px;
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

.select-skill-tip {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border: 1px solid #99bffd;
  border-radius: 8px;
  background: #e5efff;
  color: #3a4559;
  font-size: 14px;
  line-height: 22px;

  .tip-icon {
    flex-shrink: 0;
    color: #2475fc;
    font-size: 14px;
  }
}

.select-skill-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 20px;

  :deep(.ant-btn) {
    height: 36px;
    border-radius: 6px;
  }
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

.select-skill-loading,
.select-skill-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 248px;
}

.select-skill-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px 16px;
  max-height: 304px;
  margin-top: 20px;
  overflow-y: auto;
  padding-right: 2px;
}

.select-skill-card {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  min-height: 88px;
  padding: 16px;
  border: 1px solid #e5eaf2;
  border-radius: 8px;
  background: #fff;
  cursor: pointer;
  transition:
    border-color 0.2s ease,
    box-shadow 0.2s ease;

  &:hover,
  &.selected {
    border-color: #80aaff;
    box-shadow: 0 0 0 1px #80aaff inset;
  }

  :deep(.ant-checkbox) {
    top: 2px;
  }
}

.select-skill-card-main {
  min-width: 0;
  flex: 1;
}

.select-skill-card-title {
  color: #262626;
  font-size: 14px;
  font-weight: 600;
  line-height: 22px;
}

.select-skill-card-desc {
  display: -webkit-box;
  margin-top: 6px;
  overflow: hidden;
  color: #8c8c8c;
  font-size: 13px;
  line-height: 20px;
  word-break: break-word;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
}

.modal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  padding: 10px 24px;
  border-top: 1px solid #f0f0f0;
  border-radius: 0 0 16px 16px;
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
  color: #595959;
  border-color: #d9d9d9;
}

.confirm-btn {
  border: none;
  background: #2475fc;
  box-shadow: none;
}

@media (max-width: 768px) {
  .select-skill-grid {
    grid-template-columns: 1fr;
  }
}
</style>
