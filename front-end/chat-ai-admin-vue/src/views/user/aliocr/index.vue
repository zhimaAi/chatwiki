<template>
  <div class="aliocr-main">
    <div class="aliocr-title">阿里云OCR</div>
    <div class="aliocr-switch-box">
      <div class="aliocr-switch-nav-box">
        <svg-icon name="ali-ocr" style="font-size: 24px; color: #FF5915"></svg-icon>
        <div class="aliocr-switch-nav-title">阿里云OCR</div>
      </div>
      <div class="aliocr-info">开启后，Pdf文档解析支持阿里云OCR解析，结果更精准，高效</div>
      <div class="aliocr-switch-btn-box">
        <a-switch
          @change="onChangeSwitch"
          v-model:checked="aliocrSwitch"
          checked-children="开"
          un-checked-children="关"
          :checkedValue="true"
          :unCheckedValue="false"
        />
      </div>
    </div>
    <a-spin tip="加载中..." :spinning="spinning">
      <div class="aliocr-config-box">
        <div class="aliocr-config-nav">
          <SettingOutlined style="font-size: 16px;" />
          <div class="aliocr-config-nav-title">接入配置</div>
          <div class="aliocr-config-nav-line"></div>
          <div class="alicor-config-nav-url" @click="onGoUrl('https://www.aliyun.com/product/ai/docmind?spm=5176.28536895.J_kUfM_yzYYqU72woCZLHoY.6.28bd586cZUeowZ')">如何接入?</div>
        </div>
        <div class="aliocr-config-form">
          <div class="aliocr-config-form-item">
            <div class="aliocr-config-form-label">Accesskey ID</div>
            <a-input v-model:value="formState.ali_ocr_key" placeholder="请输入" />
          </div>
          <div class="aliocr-config-form-item">
            <div class="aliocr-config-form-label">Accesskey Secret</div>
            <a-input v-model:value="formState.ali_ocr_secret" placeholder="请输入" />
          </div>
          <div class="aliocr-config-form-item">
            <div class="aliocr-config-form-label">配置完成后，点击 <div class="alicor-config-nav-url" @click="handleTest">调用测试</div> </div>
            <a-button type="primary" @click="handleSave">保存</a-button>
          </div>
        </div>
      </div>
    </a-spin>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, createVNode } from 'vue'
import { useCompanyStore } from '@/stores/modules/company'
import { checkAliOcr, saveAliOcr } from '@/api/user/index.js'
import { message, Modal } from 'ant-design-vue'
import { SettingOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue'

const aliocrSwitch = ref(false)
const formState = reactive({
  ali_ocr_key: '',
  ali_ocr_secret: ''
})
const spinning = ref(false)
const companyStore = useCompanyStore()

const handleGetCompany = async() => {
  await companyStore.getCompanyInfo()
}

function checkFormState() {
  if (!formState.ali_ocr_key) {
    return '请输入Accesskey ID'
  }
  if (!formState.ali_ocr_secret) {
    return '请输入Accesskey Secret'
  }
}

const onChangeSwitch = (check) => {
  if (check) {
    const errMsg = checkFormState()
    if (errMsg) {
      aliocrSwitch.value = !aliocrSwitch.value
      Modal.confirm({
        title: `接入阿里云OCR配置`,
        icon: createVNode(ExclamationCircleOutlined),
        content: `开启前请先接入阿里云OCR`,
        okText: '知道了',
        cancelText: null,
        okCancel: false,
        onOk() {
        },
        onCancel() {}
      })
    } else {
      handleSave()
    }
  } else {
    handleSave()
  }
}

const handleTest = () => {
  const errMsg = checkFormState()
  if (errMsg) {
    return message.error(errMsg)
  }
  spinning.value = true
  checkAliOcr({
    ...formState
  }).then((res) => {
    message.success('测试通过')
  }).finally(() => {
    spinning.value = false
  })
}

const handleSave = () => {
  const errMsg = checkFormState()
  if (errMsg) {
    return message.error(errMsg)
  }
  spinning.value = true
  saveAliOcr({
    ali_ocr_switch: aliocrSwitch.value ? 1 : 2,
    ...formState
  }).then((res) => {
    message.success('保存成功')
    handleGetCompany()
  }).finally(() => {
    spinning.value = false
  })
}

const onGoUrl = (url) => {
  window.open(url)
}

onMounted(async() => {
  await handleGetCompany()
  if (companyStore.ali_ocr_key) {
    formState.ali_ocr_key = companyStore.ali_ocr_key
  }
  if (companyStore.ali_ocr_secret) {
    formState.ali_ocr_secret = companyStore.ali_ocr_secret
  }
  if (companyStore.ali_ocr_switch) {
    aliocrSwitch.value = companyStore.ali_ocr_switch == 1 ? true : false
  }
})
</script>

<style lang="less" scoped>
.aliocr-main {
  .aliocr-title {
    display: flex;
    width: 100%;
    padding: 12px 24px;
    flex-direction: column;
    align-items: flex-start;
    border-radius: 2px 2px 0 0;
    border-bottom: 1px solid #F0F0F0;
    background: #FFF;
    color: #262626;
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
  }

  .aliocr-switch-box {
    display: flex;
    flex-direction: column;
    width: 552px;
    height: 144px;
    flex-shrink: 0;
    border-radius: 6px;
    background: #FFF0EB;
    padding: 24px;
    margin: 16px 24px;
    background-image: url('@/assets/img/database/ali-ocr.png');
    background-size: 297px auto;
    background-repeat: no-repeat;
    background-position: right;

    .aliocr-switch-nav-box {
      display: flex;
      align-items: center;
      gap: 8px;

      .aliocr-switch-nav-title {
        color: #262626;
        font-size: 20px;
        font-style: normal;
        font-weight: 600;
        line-height: 28px;
      }
    }

    .aliocr-info {
      color: #3a4559;
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 22px;
      margin: 8px 0 16px;
    }
  }

  .aliocr-config-box {
    width: calc(100% - 48px);
    height: 130px;
    flex-shrink: 0;
    border-radius: 6px;
    background: #F2F4F7;
    margin-left: 24px;

    .aliocr-config-nav {
      display: flex;
      width: 100%;
      padding: 16px;
      align-items: center;
      gap: 8px;
      border-radius: 4px;

      .aliocr-config-nav-title {
        color: #262626;
        font-size: 16px;
        font-style: normal;
        font-weight: 600;
        line-height: 24px;
      }

      .aliocr-config-nav-line {
        width: 1px;
        height: 16px;
        border-radius: 1px;
        background: #D9D9D9;
      }
    }

    .aliocr-config-form {
      display: flex;
      align-items: center;
      gap: 16px;
      margin: 0 16px;

      .aliocr-config-form-item {
        display: flex;
        width: 360px;
        flex-direction: column;
        align-items: flex-start;
        gap: 4px;

        .aliocr-config-form-label {
          color: #262626;
          text-align: right;
          font-size: 14px;
          font-style: normal;
          font-weight: 400;
          line-height: 22px;
        }
      }
    }

    .alicor-config-nav-url {
      display: inline-block ;
      color: #2475fc;
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 22px;
      cursor: pointer;

      &:hover {
        opacity: 0.8;
      }
    }
  }
}
</style>