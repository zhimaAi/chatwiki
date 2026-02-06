<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        desc=""
        @close="handleClose"
      >
      </NodeFormHeader>
    </template>
    <div class="problem-optimization-form">
      <div class="node-form-content">
        <a-form ref="formRef" layout="vertical" :model="formState">
          <div class="gray-block">
            <div class="gray-block-title"></div>
            <div class="model-setting-form">
              <a-textarea
                style="min-height: 300px"
                v-model:value.trim="formState.content"
                auto-size
                :placeholder="t('ph_input_comment')"
              />
            </div>
          </div>
        </a-form>
      </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'
import NodeFormLayout from '../node-form-layout.vue'
import NodeFormHeader from '../node-form-header.vue'
import { reactive, watch, onMounted, h } from 'vue'

const { t } = useI18n('views.workflow.components.node-form-drawer.explain-node-form.index')

const emit = defineEmits(['update-node'])
const props = defineProps({
  lf: {
    type: Object,
    default: null
  },
  nodeId: {
    type: String,
    default: ''
  },
  node: {
    type: Object,
    default: () => ({})
  }
})

const formState = reactive({
  content: ''
})

const update = () => {
  const data = JSON.stringify({
    ...formState
  })

  emit('update-node', {
    ...props.node,
    node_params: data
  })
}

const init = () => {
  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'
    let data = JSON.parse(dataRaw)
    formState.content = data.content
  } catch (error) {
    console.log(error)
  }
}

watch(
  () => formState,
  () => {
    update()
  },
  { deep: true }
)

const handleClose = () => {
  emit('close')
}

onMounted(() => {
  init()
})
</script>

<style lang="less" scoped>
@import '../form-block.less';
.model-setting-form {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}
</style>
