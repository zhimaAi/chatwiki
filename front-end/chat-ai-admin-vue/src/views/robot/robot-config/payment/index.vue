<template>
  <div class="_container">
    <div class="main-tit">
      <div>{{ t('title_app_charging') }}
        <a-switch :checked="switchStatus" :checked-children="t('switch_on')" :un-checked-children="t('switch_off')" @change="switchChange"/>
      </div>
      <router-link :to="{path: '/robot/ability/payment-guide', query: route.query}">
        <a class="link">{{ t('link_view_feature_intro') }}</a>
      </router-link>
    </div>
    <div class="tab-box">
      <ListTabs :tabs="tabs" v-model:value="active" @change="tabChange"/>
    </div>
    <LoadingBox v-if="loading"/>
    <SettingsBox v-else-if="active == 1" :config="config" @update="updateConfig"/>
    <AuthCodeBox v-else :config="config"/>
  </div>
</template>

<script setup>
import {ref, onMounted, computed} from 'vue';
import {useRoute, useRouter} from 'vue-router';
import { useI18n } from '@/hooks/web/useI18n';
import ListTabs from "@/components/cu-tabs/list-tabs.vue";
import SettingsBox from "./components/settings-box.vue";
import AuthCodeBox from "./components/auth-code-box.vue";
import {getPaymentSetting} from "@/api/robot/payment.js";
import LoadingBox from "@/components/common/loading-box.vue";
import {jsonDecode} from "@/utils/index.js";
import {saveRobotAbilitySwitchStatus} from "@/api/explore/index.js";
import {Modal, message} from 'ant-design-vue';
import {useRobotStore} from "@/stores/modules/robot.js";

const { t } = useI18n('views.robot.robot-config.payment.index')

const route = useRoute()
const router = useRouter()
const robotStore = useRobotStore()
const switchStatus = computed(() => robotStore.paymentSwitchStatus === '1')
//const aiReplyStatus = computed(() => robotStore.paymentAiReplyStatus === '1')

const loading = ref(true)
const config = ref({})
const active = ref(2)
const tabs = ref([
  {
    title: t('tab_billing_settings'),
    value: 1
  },
  {
    title: t('tab_auth_code_management'),
    value: 2
  },
])

onMounted(() => {
  if (route.query.form_guide == 1) {
    active.value = 1
  } else {
    active.value = Number(localStorage.getItem('zm:payment:tab:active') || 1)
  }
  loadConfig()
})

function loadConfig() {
  loading.value = true
  getPaymentSetting({robot_id: route.query.id}).then(res => {
    let data = res?.data || {}
    if (!data || !Object.keys(data).length && !route.query.form_guide) {
      router.push({
        path: '/robot/ability/payment-guide',
        query: route.query
      })
    } else {
      updateConfig(data)
    }
  }).finally(() => {
    loading.value = false
  })
}

function updateConfig(data) {
  data.count_package = jsonDecode(data.count_package, [])
  data.duration_package = jsonDecode(data.duration_package, [])
  config.value = data
}

function tabChange() {
  localStorage.setItem('zm:payment:tab:active', active.value)
}

function switchChange(checked) {
  const switch_status = checked ? '1' : '0'

  const submit = () => {
    return saveRobotAbilitySwitchStatus({
      robot_id: route.query.id,
      ability_type: 'robot_payment',
      switch_status
    }).then(() => {
      robotStore.setPaymentSwitchStatus(switch_status)
      message.success(t('msg_operation_success'))
      window.dispatchEvent(
        new CustomEvent('robotAbilityUpdated', {
          detail: { robotId: route.query.id }
        })
      )
    })
  }

  if (switch_status === '0') {
    Modal.confirm({
      title: t('modal_title_tip'),
      content: t('modal_confirm_close'),
      onOk: submit
    })
  } else {
    submit()
  }
}

</script>

<style scoped lang="less">
._container {
  height: 100%;
  overflow-y: auto;
  padding: 16px 24px;

  .main-tit {
    color: #262626;
    font-size: 16px;
    font-weight: 600;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .link {
      font-weight: 400;
      font-size: 14px;
    }
  }

  .tab-box {
    display: flex;
    margin-top: 16px;
    margin-bottom: 8px;
  }
}
</style>
