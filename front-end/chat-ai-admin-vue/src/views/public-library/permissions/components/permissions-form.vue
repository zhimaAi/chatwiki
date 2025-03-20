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
        <span class="form-box-label-text">访问权限</span>
      </div>
      <div>
        <!-- <a-button type="primary" size="small" style="margin-left: 8px" @click.prevent="onSubmit"
          >保存</a-button
        > -->
      </div>
    </div>

    <a-form :label-col="{ span: 4 }" style="width: 750px">
      <a-form-item label="访问权限">
        <a-radio-group v-model:value="formState.access_rights" @change="onSubmit">
          <a-radio value="0">私有，仅能与应用关联，作为应用知识库</a-radio>
          <a-radio value="1">公开，发布到互联网，所有人可访问</a-radio>
        </a-radio-group>
      </a-form-item>

      <a-form-item label="分享链接" v-if="formState.access_rights == 1">
        <div class="share-url-box">
          <a-input-group compact style="width: 460px">
            <a-select
              style="width: 45%"
              v-model:value="formState.share_url"
              placeholder="请选择自定义域名"
              @change="onSubmit"
            >
              <a-select-option :value="item.url" v-for="item in domainList" :key="item.id">{{
                item.url
              }}</a-select-option>
            </a-select>
            <a-input
              :disabled="true"
              :value="'/open/home/' + formState.library_key"
              style="width: 55%"
            />
          </a-input-group>

          <a-button class="copy-btn" @click="copyText">复制</a-button>
        </div>
        <div><a href="/#/user/domain">添加自定义域名</a></div>
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup>
import useClipboard from 'vue-clipboard3'
import { getDomainList } from '@/api/user/index'
import { reactive, ref, onMounted, watch, computed } from 'vue'
import { message } from 'ant-design-vue'

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
  return formState.share_url + '/open/home/' + formState.library_key
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
      } else {
        formState.share_url = window.location.origin
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
    message.error('请选择自定义域名')
  }

  let text = shareUrl.value

  try {
    await toClipboard(text)
    message.success('复制成功')
  } catch (e) {
    message.error('复制失败')
  }
}

watch(props.value, (val) => {
  setFormState({ ...val })
})

onMounted(() => {
  getMyDomainList()
})
</script>
