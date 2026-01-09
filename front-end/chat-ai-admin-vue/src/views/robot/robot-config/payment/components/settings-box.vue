<template>
  <div>
    <a-alert
      type="info"
      show-icon
      message="当用户无可用次数时，会返回给用户套餐购买说明，引导用户购买。本功能仅管控微信公众号、微信小程序、微信客服对外服务，其他渠道不受管控"
      class="zm-alert-info"
    />
    <div class="content">
      <div class="left">
        <a-button @click="showCopyModal">复制设置到其他应用</a-button>
        <div class="form-box">
          <div class="form-item">
            <div class="form-item-tit">允许试用次数</div>
            <a-select v-model:value="formState.try_count" style="width: 320px;">
              <a-select-option v-for="i in 50" :key="i" :value="i">{{ i }}</a-select-option>
            </a-select>
          </div>
          <div class="form-item">
            <div class="form-item-tit">套餐方案</div>
            <a-radio-group v-model:value="formState.package_type" @change="packageTypeChange">
              <a-radio value="1">按次数收费</a-radio>
              <a-radio value="2">按时长收费</a-radio>
            </a-radio-group>
          </div>
          <div class="form-item">
            <div class="form-item-tit">套餐方案</div>
            <div>
              <a-button type="primary" ghost :disabled="currentPackage.length > 9" @click="showAddModal()">
                <PlusOutlined/>
                新增套餐 ({{ currentPackage.length }}/10)
              </a-button>
              <table class="combo-table mt8">
                <thead>
                <tr>
                  <th class="ant-table-cell" width="120">套餐名</th>
                  <th class="ant-table-cell" width="100">时长 (天)</th>
                  <th class="ant-table-cell" width="120">总次数 (次)</th>
                  <th class="ant-table-cell" width="100">费用 (元)</th>
                  <th class="ant-table-cell" width="140">操作</th>
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
                        <a @click="showAddModal(element, index)">编辑</a>
                        <a @click="delPackage(element, index)" class="ml16">删除</a>
                      </td>
                    </tr>
                  </template>
                </draggable>
              </table>
              <EmptyBox v-if="!currentPackage.length"/>
            </div>
          </div>
          <div class="form-item">
            <div class="form-item-tit">联系二维码（推荐上传微信二维码方便用户联系）</div>
            <div>
              <a-upload
                accept=".png,.jpg,.jpeg"
                :show-upload-list="false"
                :before-upload="handleBeforeUpload"
              >
                <a-button type="primary" ghost>
                  <PlusOutlined/>
                  上传图片
                </a-button>
              </a-upload>
              <div class="qrcode-box mt8" v-if="formState.contact_qrcode">
                <img class="qrcode" :src="formState.contact_qrcode"/>
                <a @click="delQrcode">删除</a>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="right">
        <div class="tit">无使用次数用户，对话时默认返回图片</div>
        <div class="preview-box" id="previewBox">
          <div :class="['wrap',{'no-qrcode': !formState.contact_qrcode}]">
            <div class="tip">
              <ExclamationCircleFilled class="icon"/>
              您的使用次数已耗尽，请扫码联系客服购买
            </div>
            <img v-if="formState.contact_qrcode" class="qrcode" :src="formState.contact_qrcode" crossorigin="anonymous"/>
            <div class="title">次数套餐</div>
            <div class="table-box">
              <table>
                <tr v-for="(item, i) in currentPackage" :key="i">
                  <td>{{ item.name }}</td>
                  <td v-if="formState.package_type == 2">{{ item.duration }}天</td>
                  <td>{{ item.count }}次</td>
                  <td>{{ item.price }}元</td>
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
        content="配置修改后需保存后生效"
        :getPopupContainer="triggerNode => triggerNode">
        <a-button type="primary" style="width: 120px;" :loading="saving" @click.stop="save">保 存</a-button>
      </a-popover>
    </div>

    <CombStoreModal ref="packageRef" :type="formState.package_type" @ok="packageChange"/>
    <CopyConfigModal ref="copyRef" :robot-id="route.query.id"/>
  </div>
</template>

<script setup>
import {ref, reactive, onMounted, nextTick, computed, h} from 'vue';
import {useRoute} from 'vue-router';
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
    title: `提示`,
    icon: h(ExclamationCircleOutlined),
    content: `确认删除该套餐?`,
    okText: '确 定',
    cancelText: '取 消',
    onOk() {
      currentPackage.value.splice(index, 1)
    }
  })
}

function handleBeforeUpload(file) {
  try {
    const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png'
    if (!isJpgOrPng) throw '只支持JPG、PNG格式的图片'
    const isLt2M = file.size / 1024 < 1024 * 2
    if (!isLt2M) throw '图片大小不能超过2M'
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
    title: `提示`,
    icon: h(ExclamationCircleOutlined),
    content: `确认删除联系二维码?`,
    okText: '确 定',
    cancelText: '取 消',
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
      message.success('已保存')
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
