<template>
  <div>
    <a-modal
      v-model:open="show"
      :footer="null"
      :width="820"
      wrapClassName="no-padding-modal"
    >
      <template #title>
        <div class="modal-title-block">
           <span>{{ t('title_run_test') }}</span>
        </div>
      </template>
      <div class="flex-content-box">
        <div class="test-model-box">
          <div class="auto-extract-tip"><a-alert :message="t('msg_auto_extract_tip')" type="info" show-icon /></div>
          <a-form
            :model="formState"
            ref="formRef"
            layout="vertical"
            :wrapper-col="{ span: 24 }"
            autocomplete="off"
          >
            <a-form-item
              v-for="item in test_params"
              :key="item.node_key"
              :rules="[{ required: item.field.required, message: t('msg_required_field', { label: item.field.label }) }]"
            >
              <template #label>
                <a-flex :gap="4" >
                    <span>{{ item.field.label }}</span>
                    <a-tag style="margin: 0">{{ item.field.typ }}</a-tag>
                </a-flex>
              </template>
              <template v-if="item.field.typ == 'string'">
                <a-input :placeholder="t('ph_input')" v-model:value="item.field.Vals" />
              </template>
              <template v-if="item.field.typ == 'number'">
                <a-input-number
                  style="width: 100%"
                  :placeholder="t('ph_input')"
                  v-model:value="item.field.Vals"
                />
              </template>
              <template v-if="item.field.typ.includes('array')">
                <div class="input-list-box">
                  <div class="input-list-item" v-for="(input, i) in item.field.Vals" :key="i">
                    <a-form-item-rest
                      ><a-input :placeholder="t('ph_input')" v-model:value="input.value"
                    /></a-form-item-rest>

                    <CloseCircleOutlined
                      v-if="item.field.Vals.length > 1"
                      @click="handleDelItem(item.field.Vals, i)"
                    />
                  </div>
                  <div class="add-btn-box">
                    <a-button @click="handleAddItem(item.field.Vals)" block type="dashed"
                      >{{ t('btn_add') }}</a-button
                    >
                  </div>
                </div>
              </template>
            </a-form-item>
          </a-form>

          <div class="save-btn-box">
            <a-button
              :loading="loading"
              @click="handleSubmit"
              style="background-color: #00ad3a"
              type="primary"
              >
              <span v-if="!loading"><CaretRightOutlined />{{ t('btn_run_test') }}</span>
              <span v-if="loading">{{ t('msg_generating') }}</span>
            </a-button>

            <a-button style="margin-left: 16px" type="primary" v-if="runTestStatus === 'success'" @click="handleAutoExtractOutput">{{ t('btn_auto_extract') }}</a-button>
          </div>
        </div>
        <div class="preview-box">
          <template v-if="runTestStatus === 'success'">
            <div class="preview-title">
              <div class="title-text">{{ t('label_log_details') }}</div>
            </div>
            <div class="preview-content-block">
              <div class="title-block">{{ t('label_raw_output') }}<CopyOutlined @click="handleCopy('result')" /></div>
              <div class="preview-code-box">
                <vue-json-pretty :data="testResult.result" />
              </div>
            </div>
            <div class="preview-content-block">
              <div class="title-block">{{ t('label_run_log') }}<CopyOutlined @click="handleCopy('output')" /></div>
              <div class="preview-code-box">
                <vue-json-pretty :data="testResult.output" />
              </div>
            </div>
          </template>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'
import {
  CaretRightOutlined,
  CloseCircleOutlined,
  CopyOutlined
} from '@ant-design/icons-vue'
import VueJsonPretty from 'vue-json-pretty'
import 'vue-json-pretty/lib/styles.css'
import { reactive, ref, computed } from 'vue'
import { useRobotStore } from '@/stores/modules/robot'
import { callWorkFlowHttpTest } from '@/api/robot/index'
import { message } from 'ant-design-vue'
import { copyText, generateRandomId } from '@/utils/index'

const { t } = useI18n('views.workflow.components.node-form-drawer.http-node.run-test')


const props = defineProps({
  node_key: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['save', 'autoExtract'])

const robotStore = useRobotStore()

const isLockedByOther = computed(() => {
  return robotStore.robotInfo.isLockedByOther
})


const show = ref(false)
const test_params = ref([])
const testResult = ref({
    fields: [],
    output: {},
    result: {}
})
const loading = ref(false)
const runTestStatus = ref('')
const formState = reactive({
  is_draft: true,
  robot_key: robotStore.robotInfo.robot_key,
})


const handleOpenTestModal = async (fields) => {
  if (isLockedByOther.value) {
    message.warning(t('msg_locked_by_other'))
    return
  }
  
  testResult.value = {
    fields: [],
    output: {},
    result: {}
  }
  runTestStatus.value = ''
  test_params.value = fields
  show.value = true

  test_params.value.forEach((item) => {
    if (item.field.typ.includes('array')) {
      handleAddItem(item.field.Vals)
    }
  })
}


const handleDelItem = (item, index) => {
  item.splice(index, 1)
}

const handleAddItem = (item) => {
  item.push({
    value: '',
    key: generateRandomId(16)
  })
}

const formRef = ref(null)

const handleSubmit = () => {
  formRef.value.validate().then(() => {
    const postData = { ...formState, curl_node_key: props.node_key }
    const test_params_data = JSON.parse(JSON.stringify(test_params.value))

    test_params_data.forEach((item) => {
      if (item.field.typ.includes('array')) {
        item.field.Vals = item.field.Vals.map((it) => it.value)
      }
    })

    postData.test_params = JSON.stringify(test_params_data)

    loading.value = true
    runTestStatus.value = 'loading'

    callWorkFlowHttpTest({
      ...postData
    })
      .then((res) => {
        message.success(t('msg_test_success'))
        testResult.value = res.data;
        runTestStatus.value = 'success'
      })
      .catch((err) => {
        message.error(t('msg_test_failed'))
        runTestStatus.value = 'error'
      })
      .finally(() => {
        loading.value = false
      })
  })
}


const handleCopy = (key) => {
  copyText(JSON.stringify(testResult.value[key]))

  message.success(t('msg_copy_success'))
}

const handleAutoExtractOutput = () => {
  emit('autoExtract', JSON.parse(JSON.stringify(testResult.value.fields)))
}

const open = (fields) => {
  handleOpenTestModal(fields)
}

defineExpose({
  open
})
</script>

<style lang="less" scoped>
.flex-content-box {
  display: flex;
  max-height: 600px;
  overflow: hidden;
}
.test-model-box {
  flex: 1;
  margin: 24px 0 0 0;
  padding-right: 16px;
  overflow-y: auto;
  .top-title {
    font-weight: 600;
    margin-bottom: 16px;
  }
  .auto-extract-tip{
    margin-bottom: 16px;
  }
  .save-btn-box {
    margin: 32px 0;
    margin-top: 50px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
.tooltip-content {
  white-space: pre-wrap;
}
.loading-box {
  height: 100px;
  justify-content: center;
}
.result-list-box {
  margin: 24px 0;
  width: 100%;
  border: 1px solid #ebebeb;
  border-radius: 6px;
  display: flex;
  flex-direction: column;
  padding: 8px;
  .list-item-block {
    display: flex;
    align-items: center;
    overflow: hidden;
    gap: 8px;
    padding: 8px;
    color: #333;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    .right-active-icon {
      margin-left: auto;
      color: #2475fc;
      opacity: 0;
    }
    &:hover {
      background: #f2f4f7;
      .right-active-icon {
        opacity: 1;
      }
    }
    &.active {
      color: #2475fc;
      background: #e6efff;
      .right-active-icon {
        opacity: 0;
      }
    }
    .status-block {
      font-size: 20px;
    }
    .icon-name-box {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 14px;
      font-weight: 600;
      img {
        width: 24px;
        height: 24px;
      }
    }
    .time-tag {
      width: fit-content;
      border-radius: 4px;
      height: 22px;
      background: #d2f1dc;
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 0 4px;
      font-size: 12px;
    }
    .out-put-box {
      flex: 1;
      margin-left: 24px;
      overflow: hidden;
      .out-text-box {
        background: #f2f2f2;
        border-radius: 6px;
        padding: 8px;
        width: 100%;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
}

.input-list-box {
  display: flex;
  flex-direction: column;
  gap: 8px;
  .input-list-item {
    display: flex;
    gap: 8px;
  }
}

.preview-box {
  flex: 1;
  border-left: 1px solid #d9d9d9;
  padding: 16px;
  overflow-y: auto;
  .preview-title {
    display: flex;
    align-items: center;
    gap: 8px;
    .title-text {
      font-size: 15px;
      font-weight: 600;
    }
    .icon-name-box {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 14px;
      margin-left: 12px;
      img {
        width: 16px;
        height: 16px;
      }
    }
    .time-tag {
      width: fit-content;
      border-radius: 4px;
      height: 22px;
      background: #d2f1dc;
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 0 4px;
      font-size: 12px;
    }
  }
  .preview-content-block {
    margin-top: 16px;
    .title-block {
      font-size: 15px;
      color: #262626;
      display: flex;
      align-items: center;
      gap: 4px;
      .anticon-copy {
        cursor: pointer;
        &:hover {
          color: #2475fc;
        }
      }
    }
    .preview-code-box {
      width: fit-content;
      margin-top: 16px;
      padding: 8px;
      border-radius: 8px;
      border: 1px solid #d9d9d9;

      &::v-deep(.vjs-tree) {
        width: fit-content;
      }

      &::v-deep(.vjs-tree-node) {
        width: calc(100% + 16px);
        padding-right: 16px;
      }
    }
  }
}

.modal-title-block{
  display: flex;
  align-items: center;
  gap: 12px;
  .run-detail{
    display: flex;
    align-items: center;
    gap: 16px;
    background: #BFFBD7;
    padding: 4px 16px;
    font-size: 14px;
    color: #595959;
    border-radius: 8px;
  }
}

</style>
