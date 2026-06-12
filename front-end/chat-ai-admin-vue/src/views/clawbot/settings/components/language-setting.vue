<template>
  <div class="setting-box">
    <div class="edit-box-header">
      <div class="left-box">
        <a-tabs v-model:activeKey="lang_key" size="small">
          <a-tab-pane
            v-for="item in languageMapList"
            :key="item.value"
            :tab="item.label"
          ></a-tab-pane>
        </a-tabs>
      </div>
      <div class="right-box">
        <div class="actions-box">
          <template v-if="isEdit">
            <a-flex :gap="8">
              <a-button @click="handleSave" size="small" type="primary">保存</a-button>
              <a-button @click="handleEdit(false)" size="small">取消</a-button>
            </a-flex>
          </template>
          <template v-else>
            <a-button @click="handleEdit(true)" size="small">修改</a-button>
          </template>
        </div>
      </div>
    </div>
    <div class="edit-box-body" v-if="currentMultiLangConfig">
      <WelcomeWords
        :isEdit="isEdit"
        :backCurrentMultiLangConfig="backCurrentMultiLangConfig"
        :currentMultiLangConfig="currentMultiLangConfig"
        @save="handleUpdataConfig"
      />
      <UnknownProblemPrompt
        v-if="!isWorkflowRobot"
        :isEdit="isEdit"
        :backCurrentMultiLangConfig="backCurrentMultiLangConfig"
        :currentMultiLangConfig="currentMultiLangConfig"
        @save="handleUpdataConfig"
      />
      <AnswerPrompt
        :currentMultiLangConfig="currentMultiLangConfig"
        :backCurrentMultiLangConfig="backCurrentMultiLangConfig"
        :isEdit="isEdit"
        @save="handleUpdataConfig"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watchEffect } from 'vue'
import { message } from 'ant-design-vue'
import { saveRobotLangConfigs } from '@/api/robot/index'
import { useClawbotStore } from '@/stores/modules/clawbot'
import { useRobotStore } from '@/stores/modules/robot'
import WelcomeWords from '@/views/robot/robot-config/basic-config/components/language-setting/welcome-words.vue'
import UnknownProblemPrompt from '@/views/robot/robot-config/basic-config/components/language-setting/unknown-problem-prompt.vue'
import AnswerPrompt from '@/views/robot/robot-config/basic-config/components/language-setting/answer-prompt.vue'
import { languageMapList } from '@/views/robot/robot-config/basic-config/components/language-setting/languageMap'

const lang_key = ref('zh-CN')

const clawbotStore = useClawbotStore()
const robotStore = useRobotStore()
const isWorkflowRobot = computed(() => Number(robotStore.robotInfo.application_type || 0) === 1)

const multiLangConfigs = computed(() => {
  return robotStore.robotInfo.multi_lang_configs || []
})

const multi_lang_configs = ref([])
const isEdit = ref(false)

watchEffect(() => {
  setMultiLangConfigs()
})

function setMultiLangConfigs() {
  multi_lang_configs.value = multiLangConfigs.value.map((item) => {
    return {
      ...item
    }
  })
}

const currentMultiLangConfig = computed(() => {
  return (
    multiLangConfigs.value.find((item) => item.lang_key === lang_key.value) ||
    multiLangConfigs.value[0]
  )
})

const backCurrentMultiLangConfig = computed(() => {
  return (
    multi_lang_configs.value.find((item) => item.lang_key === lang_key.value) ||
    multi_lang_configs.value[0]
  )
})

const handleUpdataConfig = (data) => {
  const newItem = {
    ...backCurrentMultiLangConfig.value,
    ...data
  }
  const index = multi_lang_configs.value.findIndex((item) => item.lang_key === lang_key.value)
  multi_lang_configs.value.splice(index, 1, newItem)
}

const handleSave = async () => {
  const params = {
    id: robotStore.robotInfo.id,
    multi_lang_configs: JSON.stringify(multi_lang_configs.value)
  }

  await saveRobotLangConfigs(params)
  message.success('保存成功')
  isEdit.value = false
  await robotStore.getRobot(robotStore.robotInfo.id)
  await clawbotStore.fetchRobotInfo(String(robotStore.robotInfo.id))
}

const handleEdit = (val) => {
  isEdit.value = val
  if (!val) {
    setMultiLangConfigs()
  }
}
</script>

<style lang="less" scoped>
.setting-box {
  border-radius: 6px;
  overflow: hidden;
  background-color: #f2f4f7;

  .edit-box-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 9px 16px;

    .left-box {
      width: calc(100% - 100px);
      overflow: hidden;

      &::v-deep(.ant-tabs-nav) {
        margin-bottom: 0;

        &::before {
          border: 0;
        }

        .ant-tabs-nav-wrap {
          padding-left: 0;
        }
      }
    }

    .right-box {
      display: flex;
      align-items: center;
    }

    .actions-box {
      display: flex;
      align-items: center;
      line-height: 22px;
      font-size: 14px;
      color: #595959;

      .action-btn {
        cursor: pointer;
      }

      .save-btn {
        color: #2475fc;
      }
    }
  }

  .edit-box-body {
    padding: 0 16px 16px;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }
}
</style>
