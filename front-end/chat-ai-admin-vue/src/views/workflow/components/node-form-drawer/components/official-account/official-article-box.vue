<template>
  <div class="_main-box">
    <div v-if="!loginStatus" class="login-info-box">
      <div class="tit">{{ t('title_login') }}</div>
      <a-button type="primary" class="link" @click="goLogin">{{ t('btn_login') }}</a-button>
      <div class="desc">{{ t('msg_logged_in') }}
        <a-button type="link" @click="refresh" :loading="loading">{{ t('btn_refresh') }}</a-button>
      </div>
    </div>
    <PluginFormRender
      ref="formRef"
      :node="node"
      :params="action.params"
      :output="action.output"
      :variableOptions="variableOptions"
      @updateVar="emit('updateVar')">
      <template v-if="loginStatus" #input-title-extra>
        <div class="official-info">
          <img class="avatar" :src="loginInfo?.headimgurl"/>
          <span>{{ loginInfo?.nickname }}</span>
          <CheckCircleFilled class="icon"/>
        </div>
      </template>
      <template #username="{state, keyName, item}"></template>
    </PluginFormRender>

    <QrcodeBox/>
  </div>
</template>

<script setup>
import {onMounted, ref, computed} from 'vue';
import { useI18n } from 'vue-i18n';
import {CheckCircleFilled} from '@ant-design/icons-vue';
import PluginFormRender from "../pluginFormRender.vue";
import QrcodeBox from "@/components/common/qrcode-box.vue";
import {useOfficialArticleLogin} from "@/composables/useOfficialArticleLogin.js";
import {useUserStore} from "@/stores/modules/user.js";

const { t } = useI18n();

const emit = defineEmits(['updateVar'])
const props = defineProps({
  node: {
    type: Object,
    default: () => ({})
  },
  action: {
    type: Object,
  },
  actionName: {
    type: String,
  },
  variableOptions: {
    type: Array,
  }
})

const {
  loginInfo,
  loginStatus,
  loading,
  refresh,
  goLogin
} = useOfficialArticleLogin()

const userStore = useUserStore()
const formRef = ref(null)

// 开源版本仅支持admin
// const username = computed(() => {
//   return userStore?.userInfo?.user_id
// })
const username = ref('admin')

onMounted(() => {
  formRef.value.setState({
    username: {
      value: username.value,
      tags: []
    },
  })
})
</script>

<style scoped lang="less">
.login-info-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  border-radius: 6px;
  padding: 24px;
  background: #F2F4F7;
  margin-bottom: 16px;

  .link {
    width: 168px;
    margin-top: 16px;
  }

  .tit {
    color: #262626;
    font-size: 16px;
    font-weight: 600;
  }

  .desc {
    color: #8c8c8c;
    font-size: 14px;
    font-weight: 400;
    margin: 12px 0;

    :deep(.ant-btn) {
      padding: 0;
    }
  }
}

.official-info {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #262626;
  font-size: 14px;

  .avatar {
    width: 16px;
    height: 16px;
    border-radius: 4px;
  }

  .icon {
    color: #21A665;
  }
}
</style>
