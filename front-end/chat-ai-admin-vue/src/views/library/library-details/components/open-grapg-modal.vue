<template>
  <div>
    <a-modal v-model:open="open" title="请先选择生成知识图谱模型" @ok="handleOk">
      <div>
        <div class="desc-text">开启后，会使用指定模型抽取知识中实体和关系。</div>
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
const open = ref(false)
const formState = reactive({
  graph_use_model: '',
  graph_model_config_id: ''
})

const emit = defineEmits(['ok'])
const show = (data) => {
  open.value = true
  data = JSON.parse(JSON.stringify(data))
  nextTick(() => {
    formState.graph_use_model = data.graph_use_model
    formState.graph_model_config_id = data.graph_model_config_id
  })
}

const onChangeModel = () => {}

const handleOk = () => {
  open.value = false
  Modal.confirm({
    title: '确定开启生成知识图谱吗?',
    icon: null,
    content: createVNode('div', { style: 'color: #262626;' }, [
      createVNode('span', {}, '生成知识图谱依赖大模型能力'),
      createVNode('span', { style: 'color:red;' }, '可能会产生比较大的Token消耗'),
      createVNode(
        'span',
        {},
        '。知识图谱疸长精准的实体关系推理和逻辑验证（比如需要查询问题是X商品有哪些优患方案），但对非结构化文本和语义模糊查询支持较弱。建议您优先使用向量，当单纯向量检索无法满足时再尝试知识图谱'
      )
    ]),
    okText: '确定开启',
    cancelText: '暂不开启',
    onOk() {
      emit('ok', {
        ...formState
      })
    },
    onCancel() {
      emit('ok', {})
    }
  })
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
