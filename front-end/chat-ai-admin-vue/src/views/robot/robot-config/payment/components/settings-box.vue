<template>
  <div>
    <a-alert
      type="info"
      show-icon
      :message="t('alert_info_message')"
      class="zm-alert-info"
    />
    <div class="content">
      <div class="left">
        <a-button @click="showCopyModal">{{ t('btn_copy_settings') }}</a-button>
        <div class="form-box">
          <div class="form-item">
            <div class="form-item-tit">{{ t('label_allow_trial_count') }}</div>
            <a-select v-model:value="formState.try_count" style="width: 320px;">
              <a-select-option v-for="i in 50" :key="i" :value="i">{{ i }}</a-select-option>
            </a-select>
          </div>
          <div class="form-item">
            <div class="form-item-tit">{{ t('label_package_plan') }}</div>
            <a-radio-group v-model:value="formState.package_type" @change="packageTypeChange">
              <a-radio value="1">{{ t('radio_charge_by_count') }}</a-radio>
              <a-radio value="2">{{ t('radio_charge_by_duration') }}</a-radio>
            </a-radio-group>
          </div>
          <div class="form-item">
            <div class="form-item-tit">{{ t('label_package_plan') }}</div>
            <div>
              <a-button type="primary" ghost :disabled="currentPackage.length > 9" @click="showAddModal()">
                <PlusOutlined/>
                {{ t('btn_add_package') }} ({{ currentPackage.length }}/10)
              </a-button>
              <table class="combo-table mt8">
                <thead>
                <tr>
                  <th class="ant-table-cell" width="120">{{ t('table_header_package_name') }}</th>
                  <th class="ant-table-cell" width="100">{{ t('table_header_duration') }}</th>
                  <th class="ant-table-cell" width="120">{{ t('table_header_total_count') }}</th>
                  <th class="ant-table-cell" width="100">{{ t('table_header_cost') }}</th>
                  <th class="ant-table-cell" width="140">{{ t('table_header_operation') }}</th>
                </tr>
                </thead>
                <draggable
                  v-model="currentPackage"
                  @end="packageFormat"
                  item-key="id"
                  group="table-rows"
                  handle=".drag-btn"
                  class="ant-table-tbody"
                  tag="tbody"
                >
                  <template #item="{ element, index }">
                    <tr class="ant-table-row">
                      <td class="ant-table-cell">
                        <span class="drag-btn"><svg-icon name="drag"/></span>
                        {{ element.name }}
                      </td>
                      <td class="ant-table-cell">{{ element.duration }}</td>
                      <td class="ant-table-cell">{{ element.count }}</td>
                      <td class="ant-table-cell">{{ element.price }}</td>
                      <td class="ant-table-cell">
                        <a @click="showAddModal(element, index)">{{ t('btn_edit') }}</a>
                        <a @click="delPackage(element, index)" class="ml16">{{ t('btn_delete') }}</a>
                      </td>
                    </tr>
                  </template>
                </draggable>
              </table>
              <EmptyBox v-if="!currentPackage.length"/>
            </div>
          </div>
          <div class="form-item">
            <div class="form-item-tit">{{ t('label_contact_qrcode') }}</div>
            <div>
              <a-upload
                accept=".png,.jpg,.jpeg"
                :show-upload-list="false"
                :before-upload="handleBeforeUpload"
              >
                <a-button type="primary" ghost>
                  <PlusOutlined/>
                  {{ t('btn_upload_image') }}
                </a-button>
              </a-upload>
              <div class="qrcode-box mt8" v-if="formState.contact_qrcode">
                <img class="qrcode" :src="formState.contact_qrcode"/>
                <a @click="delQrcode">{{ t('btn_delete_qrcode') }}</a>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="right">
        <div class="tit">{{ t('label_preview_title') }}</div>
        <div class="preview-box" id="previewBox">
          <div :class="['wrap',{'no-qrcode': !formState.contact_qrcode}]">
            <div class="tip">
              <ExclamationCircleFilled class="icon"/>
              {{ t('preview_tip_usage_exhausted') }}
            </div>
            <img v-if="formState.contact_qrcode" class="qrcode" :src="formState.contact_qrcode" crossorigin="anonymous"/>
            <div class="title">{{ t('preview_title_package') }}</div>
            <div class="table-box">
              <table>
                <tr v-for="(item, i) in currentPackage" :key="i">
                  <td>{{ item.name }}</td>
                  <td v-if="formState.package_type == 2">{{ item.duration }}{{ t('preview_unit_day') }}</td>
                  <td>{{ item.count }}{{ t('preview_unit_count') }}</td>
                  <td>{{ item.price }}{{ t('preview_unit_yuan') }}</td>
                </tr>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="footer">
      <a-popover
        :open="saveTipShow"
        :content="t('popover_save_tip')"
        :getPopupContainer="triggerNode => triggerNode">
        <a-button type="primary" style="width: 120px;" :loading="saving" @click.stop="save">{{ t('btn_save') }}</a-button>
      </a-popover>
    </div>

    <CombStoreModal ref="packageRef" :type="formState.package_type" @ok="packageChange"/>
    <CopyConfigModal ref="copyRef" :robot-id="route.query.id"/>
  </div>
</template>

<script setup>
import {ref, reactive, onMounted, nextTick, computed, h} from 'vue';
import {useRoute} from 'vue-router';
import { useI18n } from '@/hooks/web/useI18n';
import {message, Modal} from 'ant-design-vue';
import {PlusOutlined, ExclamationCircleFilled, ExclamationCircleOutlined} from '@ant-design/icons-vue';
import draggable from 'vuedraggable'
import CombStoreModal from "./comb-store-modal.vue";
import {getPaymentSetting, savePaymentSetting} from "@/api/robot/payment.js";
import EmptyBox from "@/components/common/empty-box.vue";
import {uploadFile} from "@/api/app/index.js";
import {canvasToFile, jsonDecode} from "@/utils/index.js";
import html2canvas from 'html2canvas'
import CopyConfigModal from "@/views/robot/robot-config/payment/components/copy-config-modal.vue";

const { t } = useI18n('views.robot.robot-config.payment.components.settings-box')

const emit = defineEmits(['update'])
const props = defineProps({
  config: {
    type: Object
  }
})
const route = useRoute()

const packageRef = ref(null)
const copyRef = ref(null)
const loading = ref(false)
const saving = ref(false)
const formStateBackup = ref('')
const formState = reactive({
  robot_id: route.query.id,
  try_count: 3,
  package_type: '2', // 1按次数收费 2按时长收费
  count_package: [],
  duration_package: [],
  contact_qrcode: '',
  package_poster: '',
})
const currentPackage = ref([])
const editPackageIndex = ref(-1)
const saveTipShow = computed(() => {
  if (formStateBackup.value == '') {
    return false
  }
  return formStateBackup.value != JSON.stringify(formState)
})

onMounted(() => {
  init()
})

function init() {
  dataAssign(JSON.parse(JSON.stringify(props.config)))
}

async function loadData() {
  loading.value = true
  getPaymentSetting({robot_id: route.query.id}).then(res => {
    let data = res?.data || {}
    emit('update', JSON.parse(JSON.stringify(data)))
    data.count_package = jsonDecode(data.count_package, [])
    data.duration_package = jsonDecode(data.duration_package, [])
    dataAssign(data)
  }).finally(() => {
    loading.value = false
  })
}

function dataAssign(data) {
  Object.assign(formState, data)
  formStateBackup.value = JSON.stringify(formState)
  packageTypeChange()
}

function showAddModal(record = null, index = -1) {
  editPackageIndex.value = index
  packageRef.value.show(record)
}

function packageTypeChange() {
  currentPackage.value = formState.package_type == 1 ? formState.count_package : formState.duration_package
}

function packageChange(info) {
  if (editPackageIndex.value > -1) {
    currentPackage.value[editPackageIndex.value] = info
  } else {
    currentPackage.value.push(info)
  }
}

function delPackage(record, index) {
  Modal.confirm({
    title: t('modal_title_tip'),
    icon: h(ExclamationCircleOutlined),
    content: t('modal_confirm_delete_package'),
    okText: t('modal_btn_confirm'),
    cancelText: t('modal_btn_cancel'),
    onOk() {
      currentPackage.value.splice(index, 1)
    }
  })
}

function handleBeforeUpload(file) {
  try {
    const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png'
    if (!isJpgOrPng) throw t('error_image_format')
    const isLt2M = file.size / 1024 < 1024 * 2
    if (!isLt2M) throw t('error_image_size')
    uploadFile({category: 'library_image', file}).then((res) => {
      formState.contact_qrcode = res.data.link
    })
  } catch (e) {
    message.error(e)
  }
  return false
}

function delQrcode() {
  Modal.confirm({
    title: t('modal_title_tip'),
    icon: h(ExclamationCircleOutlined),
    content: t('modal_confirm_delete_qrcode'),
    okText: t('modal_btn_confirm'),
    cancelText: t('modal_btn_cancel'),
    onOk() {
      formState.contact_qrcode = ""
    }
  })
}

function showCopyModal() {
  copyRef.value.show()
}

function createPoster() {
  return html2canvas(document.getElementById('previewBox'), {
    useCORS: true,
  }).then(async function (canvas) {
    const file = await canvasToFile(canvas)
    const {data} = await uploadFile({category: 'library_image', file})
    formState.package_poster = data.link
    return data
  });
}

function packageFormat() {
  const key = formState.package_type == 1 ? 'count_package' : 'duration_package'
  formState[key] = currentPackage.value
}

function save() {
  saving.value = true
  nextTick(async () => {
    await createPoster()
    packageFormat()
    let data = JSON.parse(JSON.stringify(formState))
    data.count_package = JSON.stringify(data.count_package)
    data.duration_package = JSON.stringify(data.duration_package)
    savePaymentSetting(data).then(res => {
      message.success(t('msg_saved'))
      loadData()
    }).finally(() => {
      saving.value = false
    })
  })
}
</script>

<style scoped lang="less">
.content {
  display: flex;
  justify-content: space-between;
  margin-top: 24px;
  padding-bottom: 80px;
}

.form-box {
  margin-top: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;

  .combo-table {
    border-collapse: collapse;
    border-spacing: 0;

    th {
      color: #262626;
      font-size: 14px;
      font-weight: 400;
      padding: 8px 16px;
      background: #F5F5F5;
      text-align: left;
    }

    td {
      padding: 8px 16px;
      border-bottom: 1px solid #E8E8E8;
    }

    .drag-btn {
      cursor: move;
    }
  }

  .form-item {
    .form-item-tit {
      margin-bottom: 4px;
      color: #262626;
    }

    .qrcode {
      width: 104px;
      height: 104px;
      border-radius: 2px;
      border: 1px solid #D9D9D9;
      padding: 8px;
    }
  }
}

.right {
  flex: 1;
  margin-left: 24px;
  display: flex;
  flex-direction: column;
  align-items: center;

  .preview-box {
    width: 375px;
    border-radius: 16px;
    border: 6px solid #FFF;
    background: #99E581;
    box-shadow: 0 4px 24px 0 #00000029;
    position: relative;

    .wrap {
      border-radius: 16px;
      background: url("@/assets/img/robot/app-charging/bg-v1.png") no-repeat;
      background-size: 100%;
      padding: 30px 8px 16px;
      min-height: 450px;
      overflow: hidden;
      overflow-y: auto;
      display: flex;
      flex-direction: column;

      &.no-qrcode {
        background: url("@/assets/img/robot/app-charging/bg-v2.png") no-repeat;
        background-size: 100%;
        min-height: 190px;

        .title {
          margin-top: 24px;
        }
      }
    }

    .tip {
      display: flex;
      align-items: center;
      gap: 4px;
      border-radius: 12px;
      padding: 8px 12px;
      background: #1a1a1ac7;
      box-shadow: 0 4px 8px 0 #14563c29;
      color: #FFF;
      font-size: 15px;
    }

    .title {
      color: #000000;
      text-align: center;
      font-size: 18px;
      font-weight: 600;
      margin-bottom: 12px;
    }

    .qrcode {
      width: 160px;
      height: 160px;
      margin: 44px auto 84px;
    }

    .table-box {
      overflow-y: auto;
      width: 100%;
      padding: 0 3.5px;
      color: #0c3b24;
      font-size: 14px;
      font-weight: 400;
      border-radius: 12px;

      table {
        width: 100%;
        border: none;
        border-collapse: collapse;
        border-spacing: 0;

        tr {
          td {
            padding: 4px 8px;
            border: none;
          }

          &:nth-child(odd) {
            background: #D6FFEB;
          }

          &:nth-child(even) {
            background: #FFF;
          }
        }
      }
    }
  }
}

.mt8 {
  margin-top: 8px;
}

.ml16 {
  margin-left: 16px;
}

.qrcode-box {
  display: flex;
  gap: 16px;
}

.footer {
  position: fixed;
  bottom: 0;
  left: 278px;
  padding: 24px 0;

  :deep(.ant-popover-content) {
    .ant-popover-inner {
      background: #2475fc;

      .ant-popover-inner-content {
        color: #FFF;
      }
    }

    .ant-popover-arrow {
      &:before {
        content: "";
        background: #2475FC;
      }
    }
  }
}
</style>
