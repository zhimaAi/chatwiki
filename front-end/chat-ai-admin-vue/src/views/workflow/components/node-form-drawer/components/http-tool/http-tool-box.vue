<style scoped lang="less">
@import "../node-options.less";
.node-icon{
  display: block;
  width: 20px;
  height: 20px;
  border-radius: 6px;
}
.kv {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 4px 0;
  .k {
    width: 80px;
    color: var(--wf-color-text-2);
  }
  .v {
    flex: 1;
    color: var(--wf-color-text-1);
    word-break: break-all;
  }
}
</style>

<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader :title="node.node_name" :desc="descText">
        <template #node-icon>
          <img v-if="node.properties?.node_icon" class="node-icon" :src="node.properties.node_icon"/>
          <a-spin v-else size="small"/>
        </template>
      </NodeFormHeader>
    </template>
    <div>
      <div class="node-options">
        <div class="options-title">
          <div><img src="@/assets/img/workflow/input.svg" class="title-icon"/>输入</div>
        </div>
        <div class="options-item">
          <div class="kv"><div class="k">请求方法</div><div class="v">{{ state.method }}</div></div>
          <div class="kv"><div class="k">请求地址</div><div class="v">{{ state.url }}</div></div>
          <div class="kv" v-if="state.query"><div class="k">查询参数</div><div class="v">{{ state.query }}</div></div>
          <div class="kv" v-if="state.headers"><div class="k">请求头</div><div class="v">{{ state.headers }}</div></div>
          <div class="kv" v-if="state.body"><div class="k">请求体</div><div class="v">{{ state.body }}</div></div>
        </div>
      </div>
      <div class="node-options">
        <div class="options-title">
          <div><img src="@/assets/img/workflow/output.svg" class="title-icon"/>输出</div>
        </div>
        <div class="options-item">
          <OutputFields :tree-data="outputData"/>
        </div>
      </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import {ref, inject, onMounted, computed} from 'vue'
import NodeFormLayout from '../../node-form-layout.vue'
import NodeFormHeader from '../../node-form-header.vue'
import OutputFields from "@/views/workflow/components/feishu-table/output-fields.vue";

const setData = inject('setData')

const props = defineProps({
  node: {
    type: Object,
    default: () => ({})
  }
})

const outputData = ref([])
const state = ref({
  method: '',
  url: '',
  query: '',
  headers: '',
  body: ''
})

const descText = computed(() => {
  return `${state.value.method || ''} ${state.value.url || ''}`.trim()
})

function init() {
  const nodeParams = JSON.parse(props.node.node_params || '{}')
  const args = nodeParams?.plugin?.params?.arguments || {}
  state.value.method = args.method || ''
  state.value.url = args.url || ''
  state.value.query = args.query || ''
  state.value.headers = args.headers || ''
  state.value.body = args.body || ''
  outputData.value = nodeParams?.plugin?.output_obj || []
}

onMounted(() => {
  init()
})
</script>
