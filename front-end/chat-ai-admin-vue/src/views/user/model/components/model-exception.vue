<style lang="less" scoped>
.model-exception-page {
  width: 100%;
  height: 100%;
  padding: 24px;
  overflow-y: auto;
  background-color: #fff;

  .exception-tips {
    margin-bottom: 24px;
    .tip-text {
      font-size: 14px;
    }
  }

  .backup-block {
    margin-bottom: 32px;
    .block-label {
      line-height: 22px;
      margin-bottom: 8px;
      font-size: 14px;
      font-weight: 600;
      color: #333;
    }

    .backup-select-wrapper {
      display: flex;
      align-items: center;
      gap: 12px;
    }

    .backup-model-select.ant-select {
      :deep(.ant-select-selection-item) {
        display: flex;
        align-items: center;
      }

      .backup-model-option {
        width: 100%;
      }
    }
  }

  .model-icon {
    width: 18px;
    height: 18px;
    object-fit: contain;
    flex: none;
  }

  .backup-model-option {
    min-width: 0;
    line-height: 22px;
    color: #262626;

    .model-name {
      flex: 1;
      min-width: 0;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }

  :global(.backup-model-select-popup .model-icon) {
    width: 18px;
    height: 18px;
    object-fit: contain;
    flex: none;
  }

  :global(.backup-model-select-popup .backup-model-option) {
    min-width: 0;
    line-height: 22px;
    color: #262626;
  }

  :global(.backup-model-select-popup .backup-model-option .model-name) {
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .error-logs-block {
    .block-label {
      line-height: 22px;
      margin-bottom: 8px;
      font-size: 14px;
      font-weight: 600;
      color: #333;
    }
    .app-box {
      display: flex;
      flex-wrap: wrap;
      gap: 8px 16px;
      .app-link {
        color: #262626;
        cursor: pointer;
        &:hover {
          color: #2475fc;
          text-decoration: underline;
        }
      }
      .app-deleted {
        color: #bfbfbf;
        cursor: not-allowed;
      }
    }
  }
}
</style>

<template>
  <div class="model-exception-page">
    <div class="exception-tips">
      <a-alert type="info" show-icon>
        <template #icon>
          <div>
            <ExclamationCircleFilled style="font-size: 16px" />
          </div>
        </template>
        <template #description>
          <div class="tip-text">{{ t('exception_alert_tip') }}</div>
        </template>
      </a-alert>
    </div>

    <div class="backup-block">
      <div class="block-label">{{ t('backup_model_label') }}</div>
      <div class="backup-select-wrapper">
        <a-select
          allowClear
          class="backup-model-select"
          popupClassName="backup-model-select-popup"
          v-model:value="backupValue"
          :placeholder="t('backup_model_placeholder')"
          style="width: 320px"
          @change="handleBackupChange"
        >
          <a-select-opt-group
            v-for="group in options"
            :key="group.model_config_id"
          >
            <template #label>
              <a-flex align="center" :gap="6">
                <img v-if="group.model_icon_url" class="model-icon" :src="group.model_icon_url" alt="" />
                <span>{{ group.corp }}</span>
              </a-flex>
            </template>
            <a-select-option
              v-for="model in group.models"
              :key="model.value"
              :value="model.value"
              :label="model.label"
            >
              <a-flex class="backup-model-option" align="center" :gap="6">
                <img v-if="model.model_icon_url" class="model-icon" :src="model.model_icon_url" alt="" />
                <span class="model-name">{{ model.name }}</span>
              </a-flex>
            </a-select-option>
          </a-select-opt-group>
        </a-select>
      </div>
    </div>

    <div class="error-logs-block">
      <div class="block-label">{{ t('error_logs_title') }}</div>
      <a-table
        :columns="columns"
        :data-source="errorLogs"
        :loading="logsLoading"
        :pagination="false"
        :scroll="{ y: 500 }"
        row-key="rowKey"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'app'">
            <div class="app-box">
              <template
                v-for="(robot, index) in record.robots"
                :key="`${robot.robot_id}-${robot.robot_name}-${index}`"
              >
                <span
                  v-if="robot.robot_id > 0"
                  class="app-link"
                  @click="handleJumpApp(robot)"
                >{{ robot.robot_name }}</span>
                <span v-else class="app-deleted">
                  {{ robot.robot_name || t('app_deleted') }}
                  <template v-if="robot.robot_name">({{ t('app_deleted') }})</template>
                </span>
              </template>
            </div>
          </template>
        </template>
      </a-table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { ExclamationCircleFilled } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { getBackupModelConfig, setBackupModelConfig, getModelErrorLogs } from '@/api/model/index'
import {
  getModelExceptionAppUrl,
  getBackupModelOptionGroups
} from './model-exception-options'

const { t } = useI18n('views.user.model')

const backupValue = ref(void 0)
const options = ref([])

const getBackupConfig = () => {
  getBackupModelConfig().then((res) => {
    const data = res.data || {}
    options.value = getBackupModelOptionGroups(data.options || [])
    const backup = data.backup
    backupValue.value = backup ? `${backup.model_config_id}::${backup.use_model}` : void 0
  })
}

const handleBackupChange = (val) => {
  let params = { model_config_id: 0, use_model: '' }
  if (val) {
    const idx = val.indexOf('::')
    params.model_config_id = val.slice(0, idx)
    params.use_model = val.slice(idx + 2)
  }
  setBackupModelConfig(params).then(() => {
    message.success(val ? t('backup_set_success') : t('backup_clear_success'))
    getBackupConfig()
  })
}

const errorLogs = ref([])
const logsLoading = ref(false)
const columns = [
  { title: t('error_col_date'), dataIndex: 'date', key: 'date', width: 140 },
  { title: t('error_col_model'), dataIndex: 'model_name', key: 'model_name', width: 240 },
  { title: t('error_col_app'), key: 'app' },
  { title: t('error_col_error_num'), dataIndex: 'error_num', key: 'error_num', width: 120 }
]

const getErrorLogs = () => {
  logsLoading.value = true
  getModelErrorLogs()
    .then((res) => {
      errorLogs.value = (res.data?.list || []).map((item, index) => ({
        ...item,
        rowKey: `${item.date}-${item.model_name}-${index}`
      }))
    })
    .finally(() => {
      logsLoading.value = false
    })
}

const handleJumpApp = (robot) => {
  const path = getModelExceptionAppUrl(robot)
  if (!path) return
  window.open(path, '_blank', 'noopener')
}

onMounted(() => {
  getBackupConfig()
  getErrorLogs()
})
</script>
