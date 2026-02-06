<template>
  <a-modal v-model:open="visible" :width="500" @ok="handleOk" @cancel="handleCancel" :footer="null">
    <template #title>
      <div class="modal-title-box">
        {{ t('titleCookieAuthSetting') }}
        <img src="@/assets/enterprise-icon.svg" alt="" />
      </div>
    </template>
    <a-form class="modal-form" :model="form" layout="vertical">
      <!-- 提醒IP范围 -->
      <a-form-item :label="t('labelIpRange')">
        <a-checkbox-group v-model:value="form.cookie_tip_ip_locations">
          <a-checkbox value="domestic">{{ t('labelDomesticIp') }}</a-checkbox>
          <a-checkbox value="overseas">{{ t('labelOverseasIp') }}</a-checkbox>
        </a-checkbox-group>
      </a-form-item>

      <!-- 提醒页面 -->
      <a-form-item :label="t('labelReminderPages')">
        <a-checkbox-group v-model:value="form.cookie_tip_positions">
          <!-- <a-checkbox value="home">{{ t('labelHome') }}</a-checkbox> -->
          <a-checkbox value="login">{{ t('labelLoginPage') }}</a-checkbox>
          <a-checkbox value="webapp">{{ t('labelWebapp') }}</a-checkbox>
        </a-checkbox-group>
      </a-form-item>
    </a-form>

    <!-- 底部按钮 -->
    <div class="modal-footer">
      <a-button @click="handleCancel">{{ t('btnCancel') }}</a-button>
      <a-button type="primary" @click="handleOk">{{ t('btnConfirm') }}</a-button>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, reactive, createVNode } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { ExclamationCircleFilled } from '@ant-design/icons-vue'
import { saveCookieTip } from '@/api/user/index'
import { useCompanyStore } from '@/stores/modules/company'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.enterprise.components.cookie-setting')

const companyStore = useCompanyStore()

// 模态框可见性
const visible = ref(false)

// 表单数据
const form = reactive({
  cookie_tip_positions: [], // 默认选中国内IP
  cookie_tip_ip_locations: [] // 默认选中登录页
})

// 显示模态框的方法
const show = () => {
  Modal.confirm({
    title: t('msgWarmTip'),
    icon: createVNode(ExclamationCircleFilled),
    content: t('msgTrialVersion'),
    onOk() {}
  })
  return
  visible.value = true
  let companyInfo = companyStore.companyInfo
  form.cookie_tip_positions = companyInfo.cookie_tip_positions
    ? companyInfo.cookie_tip_positions.split(',')
    : []

  form.cookie_tip_ip_locations = companyInfo.cookie_tip_ip_locations
    ? companyInfo.cookie_tip_ip_locations.split(',')
    : []
}

// 隐藏模态框的方法
const hide = () => {
  visible.value = false
}

// 取消按钮处理
const handleCancel = () => {
  hide()
}

// 确定按钮处理
const handleOk = () => {
  saveCookieTip({
    cookie_tip_positions: form.cookie_tip_positions.join(','),
    cookie_tip_ip_locations: form.cookie_tip_ip_locations.join(',')
  }).then((res) => {
    message.success(t('msgSaveSuccess'))
    hide()
    companyStore.getCompanyInfo()
  })
}

// 暴露方法给父组件调用
defineExpose({
  show,
  hide
})
</script>

<style lang="less" scoped>
.modal-title-box {
  display: flex;
  align-items: center;
  gap: 4px;
  img {
    height: 18px;
  }
}
.modal-form {
  margin-top: 40px;
  &::v-deep(.ant-form-item) {
    margin-bottom: 10px;
    .ant-form-item-label {
      margin-bottom: 0;
      padding: 0;
    }
  }
}
.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 24px;
}
</style>
