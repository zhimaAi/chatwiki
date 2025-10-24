<style lang="less" scoped>
.share-form-body {
  padding: 16px 16px 0;

  .share-url-box {
    display: flex;
    align-items: center;
  }
  .copy-btn {
    margin-left: 8px;
  }
}
.form-box-bottom {
  padding: 12px 16px;
  background-color: #f2f4f7;
}
</style>
<template>
  <div class="share-form" ref="formWrapRef">
    <div class="share-form-body">
      <a-form
        :model="formState"
        :layout="layout"
        :label-col="labelCol"
        :wrapper-col="wrapperCol"
        autocomplete="off"
      >
        <a-form-item label="访问权限">
          <a-radio-group
            v-model:value="formState.access_rights"
            @change="submit"
            :disabled="operate_rights == 2"
          >
            <a-radio value="0">私有，仅能与应用关联，作为应用知识库</a-radio>
            <a-radio value="1">公开，发布到互联网，所有人可访问</a-radio>
          </a-radio-group>
        </a-form-item>

        <a-form-item label="复制链接" v-if="formState.access_rights == 1">
          <div class="share-url-box" style="width: 100; overflow: hidden">
            <a-input-group compact>
              <a-select
                style="width: 50%"
                v-model:value="formState.share_url"
                placeholder="请选择自定义域名"
                :getPopupContainer="() => formWrapRef"
                @change="submit"
                :disabled="operate_rights == 2"
              >
                <a-select-option :value="item.url" v-for="item in domainList" :key="item.id">
                  {{ item.url }}
                </a-select-option>
              </a-select>
              <a-input :disabled="true" :value="shareUrl" style="width: 50%" />
            </a-input-group>

            <a-button class="copy-btn" @click="copyText()" v-if="!hideCopy">复制</a-button>
          </div>
          <div>
            <a target="_blank" href="/#/user/domain" v-if="operate_rights == 4">添加自定义域名</a>
          </div>
        </a-form-item>
      </a-form>
    </div>

    <div class="form-box-bottom" v-if="hideCopy">
      <a-button type="primary" ghost @click="copyText()">
        <svg-icon
          name="link-left"
          style="font-size: 16px; color: #2475fc; margin-right: 4px"
        ></svg-icon>
        <span>复制链接</span>
      </a-button>
    </div>
  </div>
</template>

<script setup>
import { OPEN_BOC_BASE_URL } from '@/constants/index'
import useClipboard from 'vue-clipboard3'
import { getDomainList } from '@/api/user/index'
import { getLibraryInfo, editLibrary } from '@/api/library/index'
import { ref, reactive, toRaw, computed } from 'vue'
import { message } from 'ant-design-vue'
import { useRoute } from 'vue-router'
import { usePublicLibraryStore } from '@/stores/modules/public-library'

const emit = defineEmits(['success'])

const libraryStore = usePublicLibraryStore()

const { operate_rights } =libraryStore

const props = defineProps({
  baseUrl: {
    type: String,
    default: OPEN_BOC_BASE_URL + '/doc'
  },
  hideCopy: {
    type: Boolean,
    default: false
  },
  layout: {
    type: String,
    default: 'horizontal'
  },
  labelCol: {
    type: Object,
    default: () => {}
  },
  wrapperCol: {
    type: Object,
    default: () => {}
  },
  libraryId: {
    type: [Number, String],
    default: ''
  },
  libraryKey: {
    type: String,
    default: ''
  }
})

const { toClipboard } = useClipboard()
const route = useRoute()

const id = computed(() => route.query.library_id)

const doc_key = ref('')

// 访问权限
const formWrapRef = ref(null)
const formState = reactive({
  type: 1,
  library_intro: '',
  library_name: '',
  model_config_id: '',
  model_define: '',
  use_model: '',
  is_offline: false,
  access_rights: '0',
  share_url: '',
  library_key: '',
  avatar: '',
  avatar_file: ''
})

const shareUrl = computed(() => {
  return props.baseUrl + '/' + doc_key.value
})
const getData = async () => {
  try {
    const res = await getLibraryInfo({ id: id.value })
    Object.assign(formState, res.data)
  } catch (error) {
    message.error('获取数据失败')
    console.error('获取知识库信息失败:', error)
  }
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

const copyText = async () => {
  if (!formState.share_url) {
    message.error('请选择自定义域名')
    return
  }

  let text = formState.share_url + shareUrl.value

  try {
    await toClipboard(text)
    message.success('复制成功')
  } catch (e) {
    message.error('复制失败')
  }
}

const submit = async () => {
  let data = { ...toRaw(formState) }

  delete data.library_key

  if (data.avatar_file) {
    data.avatar = data.avatar_file
    delete data.avatar_file
  } else {
    delete data.avatar
    delete data.avatar_file
  }
  try {
    let res = await editLibrary(data)

    libraryStore.getLibraryInfo()
    
    emit('success', res)
  } catch (error) {
    // console.error('修改失败:', error)
  }
}

const init = async (data) => {
  try {
    doc_key.value = data.doc_key

    await getData()
    getMyDomainList()
  } catch (error) {
    // message.error('修改失败')
  }
}

defineExpose({
  init
})
</script>
