<style lang="less" scoped>
.setting-box {
  position: relative;
  .robot-info-box {
    .robot-prompt {
      line-height: 22px;
      font-size: 14px;
      white-space: pre-wrap;
      word-break: break-all;
      color: #595959;
    }

    .dialog-bg-set {
      margin-top: 16px;
      cursor: pointer;
      display: inline-flex;
      padding: 5px 16px;
      justify-content: center;
      align-items: center;
      gap: 10px;
      border-radius: 6px;
      border: 1px solid #D9D9D9;
      background: #FFF;
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 22px;
    }
  }
  .switch-item {
    position: absolute;
    right: 16px;
    top: calc(50% - 8px);
  }

  .modal-item-box {
    padding: 24px 0;
    display: flex;
    flex-direction: column;
    gap: 24px;

    .modal-item {
      .label {
        color: #262626;
        font-size: 14px;
        font-style: normal;
        font-weight: 400;
        line-height: 22px;
        margin-bottom: 4px;
      }
    }
  }
}
</style>

<template>
  <edit-box class="setting-box" title="问题优化" icon-name="problem-optimization">
    <template #extra>
      <span></span>
    </template>
    <div class="robot-info-box">
      <div class="robot-prompt">
        开启后，进行知识库搜索时，会根据对话记录，利用AI补全问题缺失的信息
      </div>
      <a-tooltip>
        <template #title>
          <span>填写对话背景，可以帮助大模型在上下文不全时也能更好的补全用户问题</span>
        </template>
        <div class="dialog-bg-set" @click="isBgSet = true">
          对话背景设置
        </div>
      </a-tooltip>
    </div>

    <a-switch
      @change="handleEdit"
      class="switch-item"
      checkedValue="true"
      unCheckedValue="false"
      v-model:checked="robotInfo.enable_question_optimize"
      checked-children="开"
      un-checked-children="关"
    />

    <div class="modal-box" ref="modalBoxRef">
      <a-modal :getContainer="() => $refs.modalBoxRef" v-model:open="isBgSet" :width="472" title="用户问题建议设置" @ok="handleSave">
        <div class="modal-item-box">
          <div class="modal-item">
            <div class="label">模型</div>
            <a-radio-group v-model:value="modelStatus" @change="handleChangeModel">
              <a-radio value="0">跟随机器人</a-radio>
              <a-radio value="1">指定模型</a-radio>
            </a-radio-group>
            <ModelSelect
              v-if="modelStatus == '1'"
              modelType="LLM"
              v-model:modeName="formState.optimize_question_use_model"
              v-model:modeId="formState.optimize_question_model_config_id"
              style="width: 100%; margin-top: 8px;"
              @loaded="onVectorModelLoaded"
            />
          </div>
          <div class="modal-item">
            <div class="label">对话背景</div>
            <a-textarea
              style="height: 120px"
              v-model:value="robotInfo.optimize_question_dialogue_background"
              placeholder="描述当前对话的背景，便于大模型补全用户问题，比如：当前对话是关于chatwiki的使用问题和功能介绍" />
          </div>
        </div>
      </a-modal>
    </div>
  </edit-box>
</template>
<script setup>
import { ref, reactive, inject, toRaw, nextTick, onMounted } from 'vue'
import EditBox from './edit-box.vue'
import ModelSelect from '@/components/model-select/model-select.vue'

const isBgSet = ref(false)
const isEdit = ref(false)
const { robotInfo, updateRobotInfo } = inject('robotInfo')

// optimize_question_model_config_id 接口返回0 modelStatus设置成0，有id modelStatus设置成1
const modelStatus = ref('0')
const formState = reactive({
  optimize_question_use_model: '',
  optimize_question_model_config_id: '',

})

onMounted(() => {
  modelStatus.value = robotInfo.optimize_question_model_config_id != 0 ? '1' : '0'
  formState.optimize_question_use_model = robotInfo.optimize_question_use_model
  formState.optimize_question_model_config_id = robotInfo.optimize_question_model_config_id
})

const onSave = () => {
  updateRobotInfo({ ...toRaw(formState) })
  isEdit.value = false
}

const vectorModelList = ref([])
const onVectorModelLoaded = (list) => {
  vectorModelList.value = list

  nextTick(() => {
    if (!formState.optimize_question_use_model || !Number(formState.optimize_question_model_config_id)) {
      setDefaultModel()
    }
  })
}
const setDefaultModel = () => {
  if (vectorModelList.value.length > 0) {
    // 遍历查找chatwiki模型
    let modelConfig = null
    let model = null

    // 云版默认选中qwen-max
    for (let item of vectorModelList.value) {
      if (item.model_define === 'tongyi') {
        modelConfig = item
        for (let child of modelConfig.children) {
          if (child.name === 'qwen-max') {
            model = child
            break
          }
        }
        break
      }
    }

    if (!modelConfig) {
      modelConfig = vectorModelList.value[0]
      model = modelConfig.children[0]
    }

    if (modelConfig && model) {
      formState.optimize_question_use_model = model.name
      formState.optimize_question_model_config_id = model.model_config_id
    }
  }
}
const handleSave = () => {
  isBgSet.value = false
  updateRobotInfo({ ...toRaw(formState) })
}

const handleChangeModel = (val) => {
  if (val != 1) {
    formState.optimize_question_use_model = ''
    formState.optimize_question_model_config_id = 0
  }
}

const handleEdit = (val) => {
  updateRobotInfo({});
}
</script>
