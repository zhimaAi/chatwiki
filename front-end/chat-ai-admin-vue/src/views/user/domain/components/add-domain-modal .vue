<style lang="less" scoped>
.add-domain-form {
  padding: 16px 0 48px 0;
  .form-item {
    display: flex;
    align-items: center;
    .form-label {
      margin-right: 8px;
    }
    .form-input-box {
      flex: 1;
    }
  }

  .form-tip {
    padding: 10px 0;
    color: #999;
    font-size: 13px;
  }
}
</style>

<template>
  <a-modal v-model:open="show" :title="title" @ok="handleOk" :confirmLoading="props.confirmLoading">
    <div class="add-domain-form">
      <div class="form-item">
        <div class="form-label">{{ t('domain_label') }}</div>
        <div class="form-input-box">
          <a-input-group compact style="width: 100%">
            <a-select v-model:value="form.protocol" style="width: 20%">
              <a-select-option :value="item.value" v-for="item in protocolList" :key="item.value">
                <span>{{ item.label }}</span>
              </a-select-option>
            </a-select>
            <a-input
              v-model:value="form.url"
              style="width: 80%"
              :placeholder="t('domain_placeholder')"
              @blur="onUrlBlur"
            />
          </a-input-group>
        </div>
      </div>
      <div class="form-tip">
        {{ t('domain_tip') }}
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const props = defineProps({
  confirmLoading: {
    type: Boolean,
    default: false
  }
})

function handleUrl(url) {
  url = url.replace('。', '.').toLowerCase()
  url = url.replace(/\/$/, '') //清除末尾/;
  url = url.replace(/(^\s*)|(\s*$)/g, '')
  url = url.replace('http://', '')
  url = url.replace('https://', '')
  return url
}

const { t } = useI18n('views.user.domain.components.add-domain-modal')

const emit = defineEmits(['ok'])
const show = ref(false)
const title = ref(t('add_title'))
const protocolList = ref([
  {
    label: 'http',
    value: 'http:'
  },
  {
    label: 'https',
    value: 'https:'
  }
])

const form = reactive({
  id: null,
  url: '',
  protocol: 'http:'
})

const onUrlBlur = () => {
  try {
    if (form.url) {
      const url = new URL(form.url)

      if (url.protocol == 'http:' || url.protocol == 'https:') {
        form.protocol = url.protocol
      } else {
        form.protocol = 'http:'
      }

      form.url = url.hostname
    }
  } catch (e) {
    // console.log(e)
    form.url = handleUrl(form.url)
  }
}

const handleOk = () => {
  if (!form.url) {
    return message.error(t('url_required'))
  }

  form.url = handleUrl(form.url)

  emit('ok', { ...form })
}

const open = (record) => {
  if (record) {
    form.url = record.url
    form.protocol = record.protocol
    form.id = record.id
    title.value = t('edit_title')
  } else {
    form.id = null
    form.url = ''
    form.protocol = 'http:'
    title.value = t('add_title')
  }

  show.value = true
}

const close = () => {
  show.value = false
}

defineExpose({
  open,
  close
})
</script>
