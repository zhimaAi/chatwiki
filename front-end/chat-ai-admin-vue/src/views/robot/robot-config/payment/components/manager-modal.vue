<template>
  <a-modal
    :title="t('title_add_manager')"
    v-model:open="visible"
    width="669px"
  >
    <template #footer>
      <a-button type="primary" @click="close">{{ t('btn_confirm') }}</a-button>
    </template>
    <div class="main">
      <span>{{ t('msg_reply_in_account') }}</span>
      <span class="code-box" @click="copy">{{authCode}} <CopyOutlined/></span>
      <span>{{ t('msg_set_admin_tip') }}</span>
    </div>
    <div class="example">
      <img src="@/assets/img/robot/app-charging/example-01.png"/>
      <img src="@/assets/img/robot/app-charging/example-02.png"/>
    </div>
  </a-modal>
</template>

<script setup>
import {ref, reactive, computed} from 'vue';
import {message} from 'ant-design-vue';
import {CopyOutlined} from '@ant-design/icons-vue'
import {addAuthCodeManager} from "@/api/robot/payment.js";
import {copyText} from "@/utils/index.js";
import { useI18n } from '@/hooks/web/useI18n';

const { t } = useI18n('views.robot.robot-config.payment.components.manager-modal');

const props = defineProps({
  robotId: {
    type: [Number, String]
  }
})
const visible = ref(false)
const loading = ref(false)
const authCode = ref('')

function show() {
  loadCode()
  visible.value = true
}

function close() {
  visible.value = false
}

function loadCode() {
  loading.value = true
  addAuthCodeManager({robot_id: props.robotId}).then(res => {
    authCode.value = res?.data?.code
  }).finally(() => {
    loading.value = false
  })
}

function copy() {
  copyText(authCode.value)
  message.success(t('msg_copy_success'))
}

defineExpose({
  show
})
</script>

<style scoped lang="less">
.main {
  display: flex;
  align-items: center;
  color: #595959;
  font-size: 14px;
  gap: 4px;
  margin-bottom: 8px;
  flex-wrap: wrap;

  .code-box {
    display: flex;
    padding: 1px 8px;
    align-items: center;
    gap: 4px;
    border-radius: 6px;
    border: 1px solid #D9D9D9;
    cursor: pointer;
  }
}

.example {
  display: flex;
  gap: 16px;
   > img  {
     width: 280px;
     height: 270px;
     border-radius: 16px;
   }
}
</style>
