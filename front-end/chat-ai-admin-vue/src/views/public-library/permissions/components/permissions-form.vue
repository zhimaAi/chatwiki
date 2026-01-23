<style lang="less" scoped>
.form-box {
  overflow: hidden;
  border-radius: 6px;
  background: #f2f4f7;
  .form-item-tip {
    color: #999;
  }

  .form-box-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 16px;
    height: 56px;

    .form-box-label {
      display: flex;
      align-items: center;

      .form-box-label-icon {
        width: 16px;
        height: 16px;
        font-size: 16px;
      }
      .form-box-label-text {
        padding-left: 4px;
        font-size: 16px;
        font-weight: 600;
        color: #262626;
      }
    }
  }
  .share-url-box {
    display: flex;
    align-items: center;
  }
  .copy-btn {
    margin-left: 8px;
  }
}
</style>

<template>
  <div class="form-box">
    <div class="form-box-header">
      <div class="form-box-label">
        <svg-icon name="auth" style="font-size: 14px; color: #333"></svg-icon>
        <span class="form-box-label-text">{{ t('access_rights') }}</span>
      </div>
      <div>
        <!-- <a-button type="primary" size="small" style="margin-left: 8px" @click.prevent="onSubmit"
          >保存</a-button
        > -->
      </div>
    </div>

    <a-form :label-col="{ span: 4 }" style="width: 750px">
      <a-form-item :label="t('access_rights')">
        <a-radio-group v-model:value="formState.access_rights" @change="onSubmit">
          <a-radio value="0">{{ t('private_option') }}</a-radio>
          <a-radio value="1">{{ t('public_option') }}</a-radio>
        </a-radio-group>
      </a-form-item>

      <a-form-item :label="t('copy_link')" v-if="formState.access_rights == 1">
        <div class="share-url-box">
          <a-input-group compact style="width: 600px;display: flex;">
            <a-select
              style="width: 250px"
              v-model:value="formState.share_url"
              :placeholder="t('select_domain_placeholder')"
              @change="onSubmit"
            >
              <a-select-option :value="item.url" v-for="item in domainList" :key="item.id">{{
                item.url
              }}</a-select-option>
            </a-select>
            <a-input
              :disabled="true"
              :value="OPEN_BOC_BASE_URL + '/home/' + formState.library_key"
              style="width: 350px"
            />
          </a-input-group>

          <a-button class="copy-btn" @click="copyText">{{ t('copy') }}</a-button>
        </div>
        <div><a href="/#/user/domain">{{ t('add_custom_domain') }}</a></div>
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup>
import useClipboard from 'vue-clipboard3'
import { getDomainList } from '@/api/user/index'
import { reactive, ref, onMounted, watch, computed } from 'vue'
import { message } from 'ant-design-vue'
import { OPEN_BOC_BASE_URL } from '@/constants/index'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.public-library.permissions.components.permissions-form')

const { toClipboard } = useClipboard()

const emit = defineEmits(['submit', 'update:value'])

const props = defineProps({
  value: {
    type: Object,
    default: () => {
      return {}
    }
  },
  saveLoading: {
    type: Boolean,
    default: false
  }
})

// 访问权限
const formState = reactive({
  access_rights: '0',
  share_url: undefined,
  library_key: ''
})

const shareUrl = computed(() => {
  return formState.share_url + OPEN_BOC_BASE_URL + '/home/' + formState.library_key
})

const onSubmit = () => {
  emit('submit', formState)
}

const domainList = ref([])
const getMyDomainList = () => {
  getDomainList().then((res) => {
    domainList.value = res.data || []
    if (!formState.share_url) {
      if (domainList.value.length > 0) {
        formState.share_url = domainList.value[0].url
      }
    }
  })
}

const setFormState = (val) => {
  formState.access_rights = val.access_rights
  formState.share_url = val.share_url
  formState.library_key = val.library_key

  if (domainList.value.length == 0 || !formState.share_url) {
    getMyDomainList()
  }
}

const copyText = async () => {
  if (!formState.share_url) {
    message.error(t('please_select_domain'))

    return
  }

  let text = shareUrl.value

  try {
    await toClipboard(text)
    message.success(t('copy_success'))
  } catch (e) {
    message.error(t('copy_failed'))
  }
}

watch(props.value, (val) => {
  setFormState({ ...val })
})

onMounted(() => {
  getMyDomainList()
})
</script>
