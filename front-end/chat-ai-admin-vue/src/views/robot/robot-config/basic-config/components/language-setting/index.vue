<template>
  <div class="setting-box">
    <div class="edit-box-header">
      <div class="left-box">
        <svg-icon name="sys-promt"></svg-icon>
        <div>系统提示语</div>
      </div>
      <div class="right-box">
        <a-popover
          v-model:open="moreLanguageOpen"
          trigger="hover"
          placement="bottomRight"
          overlay-class-name="language-setting-popover"
          :overlayInnerStyle="{ padding: 0 }"
        >
          <template #content>
            <div class="more-language-box">
              <div class="more-language-header">
                <span>语言</span>
                <span>隐藏/显示</span>
              </div>
              <div class="more-language-list">
                <div
                  class="more-language-item"
                  v-for="item in moreLanguageList"
                  :key="item.value"
                  @click="toggleLanguageVisible(item.value)"
                >
                  <span class="language-name">{{ item.label }}</span>
                  <span class="language-switch" @click.stop>
                    <a-switch
                      size="small"
                      :checked="enabledMoreLanguageKeys.includes(item.value)"
                      @change="(checked) => setLanguageVisible(item.value, checked)"
                    />
                  </span>
                </div>
              </div>
            </div>
          </template>
          <a-button class="more-language-btn" size="small">
            <template #icon><SettingOutlined /></template>
            更多语言
          </a-button>
        </a-popover>
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
    <div class="tabs-box">
      <a-tabs v-model:activeKey="lang_key" size="small">
        <a-tab-pane
          v-for="item in visibleLanguageList"
          :key="item.value"
          :tab="item.label"
        ></a-tab-pane>
      </a-tabs>
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
        @save="handleUpdataConfig"
        :isEdit="isEdit"
      />
    </div>
  </div>
</template>
<script setup>
import { ref, computed, watch, watchEffect } from 'vue'
import { message } from 'ant-design-vue'
import { SettingOutlined } from '@ant-design/icons-vue'
import WelcomeWords from './welcome-words.vue'
import UnknownProblemPrompt from './unknown-problem-prompt.vue'
import AnswerPrompt from './answer-prompt.vue'
import { languageMapList } from './languageMap'
import { useRobotStore } from '@/stores/modules/robot'
import { saveRobotLangConfigs } from '@/api/robot/index'

const DEFAULT_LANGUAGE_KEYS = ['zh-CN', 'en-US']
const MORE_LANGUAGE_VISIBLE_KEY = 'robot_config_more_language_visible_keys'

const lang_key = ref('zh-CN')
const moreLanguageOpen = ref(false)
const enabledMoreLanguageKeys = ref(getLocalEnabledMoreLanguageKeys())

const robotStore = useRobotStore()
const isWorkflowRobot = computed(() => Number(robotStore.robotInfo.application_type || 0) === 1)

const multiLangConfigs = computed(() => {
  return robotStore.robotInfo.multi_lang_configs || []
})

const defaultLanguageList = computed(() => {
  return languageMapList.filter((item) => DEFAULT_LANGUAGE_KEYS.includes(item.value))
})

const moreLanguageList = computed(() => {
  return languageMapList.filter((item) => !DEFAULT_LANGUAGE_KEYS.includes(item.value))
})

const visibleLanguageList = computed(() => {
  return [
    ...defaultLanguageList.value,
    ...moreLanguageList.value.filter((item) => enabledMoreLanguageKeys.value.includes(item.value))
  ]
})

const multi_lang_configs = ref([])

watchEffect(() => {
  setMultiLangConfigs()
})

watch(
  visibleLanguageList,
  (list) => {
    if (!list.some((item) => item.value === lang_key.value)) {
      lang_key.value = 'zh-CN'
    }
  },
  { immediate: true }
)

function setMultiLangConfigs() {
  multi_lang_configs.value = multiLangConfigs.value.map((item) => {
    return {
      ...item
    }
  })
}

function getLocalEnabledMoreLanguageKeys() {
  try {
    const localValue = localStorage.getItem(MORE_LANGUAGE_VISIBLE_KEY)
    const keys = localValue ? JSON.parse(localValue) : []
    return Array.isArray(keys) ? keys.filter((key) => !DEFAULT_LANGUAGE_KEYS.includes(key)) : []
  } catch (error) {
    return []
  }
}

function saveLocalEnabledMoreLanguageKeys(keys) {
  localStorage.setItem(MORE_LANGUAGE_VISIBLE_KEY, JSON.stringify(keys))
}

function updateEnabledMoreLanguageKeys(nextKeys) {
  enabledMoreLanguageKeys.value = nextKeys
  saveLocalEnabledMoreLanguageKeys(nextKeys)
}

function setLanguageVisible(key, checked) {
  const nextKeys = checked
    ? Array.from(new Set([...enabledMoreLanguageKeys.value, key]))
    : enabledMoreLanguageKeys.value.filter((item) => item !== key)

  updateEnabledMoreLanguageKeys(nextKeys)
}

function toggleLanguageVisible(key) {
  setLanguageVisible(key, !enabledMoreLanguageKeys.value.includes(key))
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

const isEdit = ref(false)

const handleUpdataConfig = (data) => {
  let newItem = {
    ...backCurrentMultiLangConfig.value,
    ...data
  }
  let index = multi_lang_configs.value.findIndex((item) => item.lang_key === lang_key.value)
  multi_lang_configs.value.splice(index, 1, newItem)
}

const handleSave = () => {
  let parmas = {
    id: robotStore.robotInfo.id,
    multi_lang_configs: JSON.stringify(multi_lang_configs.value)
  }
  saveRobotLangConfigs(parmas).then(() => {
    message.success('保存成功')
    isEdit.value = false
    setTimeout(() => {
      robotStore.getRobot(robotStore.robotInfo.id)
    }, 600)
  })
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
      flex: 1;
      display: flex;
      align-items: center;
      gap: 8px;
      color: var(--02, #262626);
      font-size: 16px;
      font-style: normal;
      font-weight: 600;
    }

    .right-box {
      display: flex;
      align-items: center;
      gap: 8px;
    }

    .more-language-btn {
      display: flex;
      align-items: center;
      color: #595959;
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
  .tabs-box {
    overflow: hidden;
    padding: 0 16px 8px 16px;
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
  .edit-box-body {
    padding: 0 16px 16px 16px;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }
}
</style>

<style lang="less">
.language-setting-popover {
  .ant-popover-inner {
    padding: 0;
  }

  .more-language-box {
    width: 232px;
    padding: 8px 4px 8px 12px;
  }

  .more-language-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 8px 6px 4px;
    color: #8c8c8c;
    font-size: 12px;
    line-height: 20px;
  }

  .more-language-list {
    max-height: 196px;
    overflow-y: auto;
    padding-right: 4px;
  }

  .more-language-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 30px;
    padding: 0 8px 0 4px;
    cursor: pointer;
    border-radius: 4px;
    color: #262626;
    font-size: 13px;

    &:hover {
      background: #f5f7fa;
    }
  }

  .language-name {
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  .language-switch {
    display: inline-flex;
    align-items: center;
  }
}
</style>
