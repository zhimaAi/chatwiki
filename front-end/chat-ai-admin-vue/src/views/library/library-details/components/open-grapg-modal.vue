<template>
  <div>
    <a-modal v-model:open="open" :title="t('title_select_graph_model')" @ok="handleOk">
      <div>
        <div class="desc-text">{{ t('desc_graph_function') }}</div>
        <ModelSelect
          modelType="LLM"
          v-model:modeName="formState.graph_use_model"
          v-model:modeId="formState.graph_model_config_id"
          style="width: 300px"
          @change="onChangeModel"
        />
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, nextTick, createVNode } from 'vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import { Modal } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.library-details.components.open-grapg-modal')
const open = ref(false)
const running = ref(false)
const params = ref({})
const formState = reactive({
  graph_use_model: '',
  graph_model_config_id: ''
})

const emit = defineEmits(['ok'])
const show = (data) => {
  if (running.value) {
    return
  }
  running.value = true
  openConfirm()
  params.value = JSON.parse(JSON.stringify(data))
}

const onChangeModel = () => {}

const openConfirm = () => {
  Modal.confirm({
    title: t('title_confirm_enable_graph'),
    icon: null,
    content: createVNode('div', { style: 'color: #262626;' }, [
      createVNode('span', {}, t('msg_graph_llm_dependency')),
      createVNode('span', { style: 'color:red;' }, t('warning_token_consumption')),
      createVNode(
        'span',
        {},
        t('desc_graph_capability')
      )
    ]),
    okText: t('btn_confirm_enable'),
    cancelText: t('btn_later'),
    onOk() {
      open.value = true
      running.value = false
      nextTick(() => {
        formState.graph_use_model = params.value.graph_use_model
        formState.graph_model_config_id = params.value.graph_model_config_id
      })
    },
    onCancel: () => running.value = false
  })
}

const handleOk = () => {
  open.value = false
  emit('ok', {...formState})
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.desc-text {
  color: #8c8c8c;
  margin-bottom: 12px;
}
</style>
