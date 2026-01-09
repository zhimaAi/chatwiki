<template>
  <div class="_container">
    <div class="main-tit">
      <div>应用收费
        <a-switch :checked="switchStatus" checked-children="开" un-checked-children="关" @change="switchChange"/>
      </div>
      <router-link :to="{path: '/robot/ability/payment-guide', query: route.query}">
        <a class="link">查看功能介绍</a>
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
import ListTabs from "@/components/cu-tabs/list-tabs.vue";
import SettingsBox from "./components/settings-box.vue";
import AuthCodeBox from "./components/auth-code-box.vue";
import {getPaymentSetting} from "@/api/robot/payment.js";
import LoadingBox from "@/components/common/loading-box.vue";
import {jsonDecode} from "@/utils/index.js";
import {saveRobotAbilitySwitchStatus} from "@/api/explore/index.js";
import {Modal, message} from 'ant-design-vue';
import {useRobotStore} from "@/stores/modules/robot.js";

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
    title: '收费设置',
    value: 1
  },
  {
    title: '授权码管理',
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
      message.success('操作成功')
      window.dispatchEvent(
        new CustomEvent('robotAbilityUpdated', {
          detail: { robotId: route.query.id }
        })
      )
    })
  }

  if (switch_status === '0') {
    Modal.confirm({
      title: '提示',
      content: '关闭后，设置的收费策略将失效，确认关闭？',
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
