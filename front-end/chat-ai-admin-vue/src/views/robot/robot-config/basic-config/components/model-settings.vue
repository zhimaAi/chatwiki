<template>
  <edit-box class="setting-box" title="模型设置" icon-name="moxingshezhi" v-model:isEdit="isEdit">
    <template #extra>
      <div class="actions-box">
        <template v-if="isEdit">
          <a-flex :gap="8">
            <a-button @click="handleSave" size="small" type="primary">保存</a-button>
            <a-button @click="handleEdit(false)" size="small">取消</a-button>
          </a-flex>
        </template>
        <template v-else>
          <span class="model-name">{{ robotInfo.use_model }}</span>
          <a-divider type="vertical" />
          <a-button @click="handleEdit(true)" size="small">修改</a-button>
        </template>
      </div>
    </template>
    <div class="form-box" v-if="isEdit">
      <a-row :gutter="[32, 16]">
        <a-col v-bind="grid">
          <div class="form-item">
            <div class="form-item-label">
              <span>LLM模型&nbsp;</span>
            </div>
            <div class="form-item-body">
              <a-select
                v-model:value="formState.use_model"
                placeholder="请选择LLM模型"
                @change="handleChangeModel"
                style="width: 100%"
              >
                <a-select-opt-group v-for="item in modelList" :key="item.id">
                  <template #label>
                    <a-flex align="center" :gap="6">
                      <img class="model-icon" :src="item.icon" alt="" />{{item.name}}
                    </a-flex>
                  </template>
                  <a-select-option
                    :value="val"
                    :model_config_id="item.id"
                    v-for="val in item.children"
                    :key="val"
                  >
                    <span v-if="item.name == 'Azure OpenAI Service' && item.deployment_name">{{
                      item.deployment_name
                    }}</span>
                    <span v-else>{{ val }}</span>
                  </a-select-option>
                </a-select-opt-group>
              </a-select>
            </div>
          </div>
        </a-col>

        <a-col v-bind="grid">
          <div class="form-item">
            <div class="form-item-label">
              <span>温度&nbsp;</span>
              <a-tooltip>
                <template #title>温度越低，回答越严谨。温度越高，回答越发散。</template>
                <QuestionCircleOutlined class="question-icon" />
              </a-tooltip>
            </div>
            <div class="form-item-body">
              <div class="number-box">
                <div class="number-slider-box">
                  <a-slider
                    class="custom-slider"
                    v-model:value="formState.temperature"
                    :min="0"
                    :max="2"
                    :step="0.1"
                  />
                </div>
                <div class="number-input-box">
                  <a-input-number
                    v-model:value="formState.temperature"
                    :min="0"
                    :max="2"
                    :step="0.1"
                  />
                </div>
              </div>
            </div>
          </div>
        </a-col>

        <a-col v-bind="grid">
          <div class="form-item">
            <div class="form-item-label">
              <span>最大token&nbsp;</span>
              <a-tooltip>
                <template #title>问题+答案的最大token数，如果出现回答被截断，可调高此值</template>
                <QuestionCircleOutlined class="question-icon" />
              </a-tooltip>
            </div>
            <div class="form-item-body">
              <div class="number-box">
                <div class="number-slider-box">
                  <a-slider
                    class="custom-slider"
                    v-model:value="formState.max_token"
                    :min="0"
                    :max="100 * 1024"
                  />
                </div>
                <div class="number-input-box">
                  <a-input-number v-model:value="formState.max_token" :min="0" :max="100 * 1024" />
                </div>
              </div>
            </div>
          </div>
        </a-col>

        <a-col v-bind="grid">
          <div class="form-item">
            <div class="form-item-label">
              <span>上下文数量&nbsp;</span>
              <a-tooltip>
                <template #title
                  >提示词中携带的历史聊天记录轮次。设置为0则不携带聊天记录。最多设置10轮。注意，携带的历史聊天记录越多，消耗的token相应也就越多。</template
                >
                <QuestionCircleOutlined class="question-icon" />
              </a-tooltip>
            </div>
            <div class="form-item-body">
              <div class="number-box">
                <div class="number-slider-box">
                  <a-slider
                    class="custom-slider"
                    v-model:value="formState.context_pair"
                    :min="0"
                    :max="10"
                  />
                </div>
                <div class="number-input-box">
                  <a-input-number v-model:value="formState.context_pair" :min="0" :max="10" />
                </div>
              </div>
            </div>
          </div>
        </a-col>
      </a-row>
    </div>
  </edit-box>
</template>

<script setup>
import { getModelConfigOption } from '@/api/model/index'
import { ref, reactive, inject, toRaw } from 'vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import EditBox from './edit-box.vue'

const grid = reactive({ sm: 24, md: 24, lg: 12, xl: 8, xxl: 8 })
// 获取LLM模型
const modelList = ref([])

const isEdit = ref(false)
const { robotInfo, updateRobotInfo } = inject('robotInfo')
const formState = reactive({
  use_model: undefined,
  model_config_id: '',
  temperature: 0,
  max_token: 0,
  context_pair: 0
})

const handleChangeModel = (val, option) => {
  formState.model_config_id = option.model_config_id
}

const handleSave = () => {
  updateRobotInfo({ ...toRaw(formState) })
  isEdit.value = false;
}

const handleEdit = (val) => {
  formState.use_model = robotInfo.use_model
  formState.model_config_id = robotInfo.model_config_id
  formState.temperature = robotInfo.temperature
  formState.max_token = robotInfo.max_token
  formState.context_pair = robotInfo.context_pair
  isEdit.value = val
}

const getModelList = () => {
  getModelConfigOption({
    model_type: 'LLM'
  }).then((res) => {
    let list = res.data || []

    modelList.value = list.map((item) => {
      return {
        id: item.model_config.id,
        name: item.model_info.model_name,
        icon: item.model_info.model_icon_url,
        children: item.model_info.llm_model_list,
        deployment_name: item.model_config.deployment_name
      }
    })
  })
}

getModelList()
</script>

<style lang="less" scoped>
.setting-box {
  ::v-deep(.edit-box-body) {
    padding: 0;
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

    .model-name {
      font-size: 14px;
      line-height: 22px;
      color: #8c8c8c;
    }
  }

  .form-box {
    padding: 0 16px 16px 16px;

    .form-item-label {
      line-height: 22px;
      margin-bottom: 4px;
      font-size: 14px;
      color: #262626;

      .question-icon {
        color: #8c8c8c;
      }
    }

    .number-box {
      display: flex;
      align-items: center;

      .number-slider-box {
        flex: 1;
      }

      .number-input-box {
        margin-left: 20px;
      }
    }
  }
}

.model-icon {
  height: 18px;
}
</style>
