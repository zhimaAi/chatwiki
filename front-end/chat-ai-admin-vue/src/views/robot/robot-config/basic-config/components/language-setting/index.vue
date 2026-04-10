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
import { ref, reactive, inject, toRaw, nextTick, computed, watch, watchEffect } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined, CloseCircleOutlined } from '@ant-design/icons-vue'
import WelcomeWords from './welcome-words.vue'
import UnknownProblemPrompt from './unknown-problem-prompt.vue'
import AnswerPrompt from './answer-prompt.vue'
import { languageMapList } from './languageMap'
import { useRobotStore } from '@/stores/modules/robot'
import { saveRobotLangConfigs } from '@/api/robot/index'

const lang_key = ref('zh-CN')

const robotStore = useRobotStore()

const multiLangConfigs = computed(() => {
  return robotStore.robotInfo.multi_lang_configs || []
})

const multi_lang_configs = ref([])

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

const isEdit = ref(false)

const handleUpdataConfig = (data) => {
  let newItem = {
    ...backCurrentMultiLangConfig.value,
    ...data
  }
  let index = multi_lang_configs.value.findIndex((item) => item.lang_key === lang_key.value)
  multi_lang_configs.value.splice(index, 1, newItem)
  console.log(multi_lang_configs.value, data, 'multi_lang_configs==')
}

const handleSave = () => {
  let parmas = {
    id: robotStore.robotInfo.id,
    multi_lang_configs: JSON.stringify(multi_lang_configs.value)
  }
  saveRobotLangConfigs(parmas).then((res) => {
    message.success('保存成功')
    isEdit.value = false
    setTimeout(() => {
      robotStore.getRobot(robotStore.robotInfo.id)
    }, 600)
  })
}

const handleEdit = (val) => {
  isEdit.value = val
  if(!val){
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
    padding: 0 16px 16px 16px;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }
}
</style>
