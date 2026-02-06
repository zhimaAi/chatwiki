<template>
  <van-popup
    v-model:show="showCookieModal"
    round
    position="bottom"
    class="cookie-modal-popup"
    :close-on-click-overlay="false"
    :closeable="false"
  >
    <div class="cookie-modal-container">
      <h2 class="title">{{ t('title_privacy_consent') }}</h2>

      <div class="content">
        <p
          class="description"
          v-html="
            t('msg_cookie_description', {
              cookiePolicyUrl: 'https://cloud.chatwiki.com/#/privacy_policy',
              privacyPolicyUrl: 'https://cloud.chatwiki.com/#/privacy_policy'
            })
          "
        ></p>
      </div>

      <div class="button-group" v-if="props.isMobileDevice">
        <van-button type="primary" class="accept-btn" @click="handleAccept">
          {{ t('btn_accept_all') }}
        </van-button>
        <van-button class="decline-btn" @click="handleDecline"> {{ t('btn_decline') }} </van-button>
      </div>
      <div class="pc-button-group" v-else>
        <van-button class="decline-btn" @click="handleDecline"> {{ t('btn_decline') }} </van-button>
        <van-button type="primary" class="accept-btn" @click="handleAccept">
          {{ t('btn_accept_all') }}
        </van-button>
      </div>
    </div>
  </van-popup>
</template>

<script setup>
import { ref } from 'vue'
import { Popup, Button } from 'vant'
import { Storage } from '@/utils/Storage'
import { getCookieTip } from '@/api/user/index'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.chat.components.cookie-modal')

const props = defineProps({
  isMobileDevice: {
    type: Boolean,
    default: false
  }
})

const COOKIE_CONSENT_KEY = 'cookie_consent'

const showCookieModal = ref(false)
const emit = defineEmits(['onAccept', 'onDecline'])

const show = async () => {
  const consent = Storage.get(COOKIE_CONSENT_KEY)
  if (!consent) {
    let res = await getCookieTip({ position: 'webapp' })
    let is_show_cookie_tip = res.data.is_show_cookie_tip
    if (is_show_cookie_tip) {
      showCookieModal.value = true
    } else {
      emit('onAccept')
    }
  } else {
    emit('onAccept')
  }
}

const handleAccept = () => {
  Storage.set(COOKIE_CONSENT_KEY, 'accepted')
  showCookieModal.value = false
  emit('onAccept')
}

const handleDecline = () => {
  // Storage.set(COOKIE_CONSENT_KEY, 'declined')
  showCookieModal.value = false
  emit('onDecline')
}

const onGoLink = (url) => {
  window.open(url)
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.cookie-modal-popup {
  max-width: 100vw;
}

.cookie-modal-container {
  background: #fff;
  border-radius: 16px;
  margin: 0 auto;
  width: 60%;
  padding: 24px 0;
  a {
    color: #2475fc;
    text-decoration: none;
    cursor: pointer;
  }
}

@media only screen and (max-width: 500px) {
  .cookie-modal-container {
    width: 100%;
    padding: 16px 16px 24px 16px;
    .decline-btn {
      height: 40px;
      background: var(--07, #f0f2f5);
      border: none;
    }
    .accept-btn {
      height: 40px;
    }
  }
}

.title {
  color: #262626;
  font-size: 16px;
  font-weight: 600;
  line-height: 24px;
  margin-bottom: 8px;
}

.content {
  margin-bottom: 24px;
}

.description {
  color: #595959;
  font-size: 14px;
  line-height: 22px;
  &::v-deep(a) {
    color: #2475fc;
    text-decoration: none;
    cursor: pointer;
  }
}

.button-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.pc-button-group {
  display: flex;
  align-items: center;
  gap: 8px;
}
.decline-btn {
  height: 32px;
  border-radius: 8px;

  border: 1px solid var(--06, #d9d9d9);
  color: #1a1a1a;
  font-size: 14px;
  font-weight: 400;
  transition: all 0.3s;
  background: var(--10, #fff);
  &:hover {
    border-color: #2475fc;
    color: #2475fc;
  }
}

.accept-btn {
  height: 32px;
  border-radius: 8px;
  background: #2475fc;
  border: none;
  color: #fff;
  font-size: 14px;
  font-weight: 400;
  transition: all 0.3s;

  &:hover {
    opacity: 0.85;
  }
}

.accept-btn:active {
  opacity: 0.9;
  transform: scale(0.98);
}
</style>
