<template>
  <div class="ai-config-page">
    <ConfigPageMenu />
    <div class="page-container">
      <div class="ai-config-form">
        <a-form
          ref="formRef"
          :model="formState"
          :label-col="{ span: 4 }"
          :wrapper-col="{ span: 8 }"
        >
          <a-form-item>
            <template #label>
              {{ t('embed') }}&nbsp;
              <a-tooltip>
                <template #title>{{ t('embed_tooltip') }}</template>
                <InfoCircleOutlined />
              </a-tooltip>
            </template>
            <a-switch
              v-model:checked="formState.use_model_switch"
              unCheckedValue="0"
              checkedValue="1"
            />
          </a-form-item>

          <a-form-item
            :label="t('embed_model')"
            name="use_model"
            :rules="[{ required: true, message: t('select_embed_model') }]"
            v-if="formState.use_model_switch == 1"
          >
            <ModelSelect
              modelType="TEXT EMBEDDING"
              :placeholder="t('select_embed_model')"
              v-model:modeName="formState.use_model"
              v-model:modeId="formState.model_config_id"
            />
          </a-form-item>

          <a-form-item>
            <template #label>
              {{ t('ai_summary') }}&nbsp;
              <a-tooltip>
                <template #title>{{ t('ai_summary_tooltip') }}</template>
                <InfoCircleOutlined />
              </a-tooltip>
            </template>
            <a-switch v-model:checked="formState.ai_summary" unCheckedValue="0" checkedValue="1" />
          </a-form-item>

          <a-form-item
            :label="t('summary_model')"
            name="ai_summary_model"
            :rules="[{ required: true, message: t('select_summary_model') }]"
            v-if="formState.ai_summary == 1"
          >
            <ModelSelect
              modelType="LLM"
              v-model:modeName="formState.ai_summary_model"
              v-model:modeId="formState.summary_model_config_id"
            />
          </a-form-item>

          <a-form-item :wrapper-col="{ offset: 4, span: 8 }">
            <a-button type="primary" :loading="loading" @click="handleSubmit">{{ t('save') }}</a-button>
          </a-form-item>
        </a-form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { getLibraryInfo, editLibrary } from '@/api/library/index'
import { useRoute } from 'vue-router'
import { ref, reactive, computed, onMounted, toRaw } from 'vue'
import { InfoCircleOutlined } from '@ant-design/icons-vue'
import { message, Modal } from 'ant-design-vue'
import ConfigPageMenu from '../components/config-page-menu.vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import { usePublicLibraryStore } from '@/stores/modules/public-library'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.public-library.ai.index')

const libraryStore = usePublicLibraryStore()

const route = useRoute()
const library_id = computed(() => route.query.library_id)
const formRef = ref()
const loading = ref(false)
let oldState = {}
const formState = reactive({
  library_name: '',
  library_intro: '',
  ai_summary: '',
  ai_summary_model: '',
  access_rights: '',
  model_config_id: '',
  use_model: '',
  use_model_switch: '',
  statistics_set: '',
  summary_model_config_id: ''
})

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    // TODO: 调用API保存数据
    let isChangeModel = oldState.use_model != formState.use_model
    if (isChangeModel) {
      Modal.confirm({
        title: t('confirm_switch_model', { model: formState.use_model }),
        content: t('switch_model_warning'),
        okText: t('confirm'),
        cancelText: t('cancel'),
        onOk() {
          saveData()
        }
      })
    } else {
      saveData()
    }
  } catch (error) {
    console.error(t('form_validation_failed'), error)
  }
}

const saveData = () => {
  let data = {
    ...toRaw(formState)
  }

  loading.value = true
  editLibrary(data)
    .then(() => {
      message.success(t('save_success'))
      loading.value = false
      oldState = JSON.parse(JSON.stringify(data))

      libraryStore.getLibraryInfo()
    })
    .catch(() => {
      loading.value = false
    })
}

const getData = () => {
  getLibraryInfo({ id: library_id.value }).then((res) => {
    let data = res.data || {}
    data.use_model = data.use_model || void 0
    data.ai_summary_model = data.ai_summary_model || void 0
    oldState = JSON.parse(JSON.stringify(data))
    Object.assign(formState, data)
  })
}

onMounted(() => {
  getData()
})
</script>

<style lang="less" scoped>
.ai-config-form {
  padding: 8px 24px 24px;
}
.model-icon {
  height: 18px;
}
</style>
