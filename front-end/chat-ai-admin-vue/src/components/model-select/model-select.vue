<template>
  <a-select
    :value="value"
    :placeholder="placeholder"
    @change="handleChangeModel"
    style="width: 100%"
    :allowClear="true"
    class="modal-seclet-new"
  >
    <template #clearIcon>
      <div class="model-show-block" v-if="selectItem">
        <img class="icon" :src="selectItem.model_icon_url" alt="" />
        <div class="name-text">
          {{ selectItem.config_name || selectItem.model_name }}
          {{ selectChildItem.show_model_name || selectChildItem.use_model_name }}
        </div>
      </div>
    </template>
    <a-select-opt-group v-for="item in modelList" :key="item.key">
      <template #label>
        <a-flex align="center" :gap="6">
          <img class="model-icon" :src="item.model_icon_url" alt="" />
          <span>{{ item.config_name || item.model_name }}</span>
        </a-flex>
      </template>
      <a-select-option
        v-for="sub in item.children"
        :value="sub.value"
        :modelId="sub.model_config_id"
        :modelName="sub.name"
        :key="sub.key"
      >
        {{ sub.show_model_name || sub.name }}
      </a-select-option>
    </a-select-opt-group>
  </a-select>
</template>

<script setup>
import { getModelConfigOption } from '@/api/model/index'
import { ref, onMounted, watch, computed } from 'vue'
import { getModelOptionsList } from '@/components/model-select/index.js'

const emit = defineEmits(['change', 'update:modeName', 'update:modeId', 'loaded'])
const props = defineProps({
  modelType: {
    type: String,
    validator: (value) => {
      return ['TEXT EMBEDDING', 'RERANK', 'LLM'].includes(value)
    },
    required: true
  },
  modeName: {
    type: [String, Number],
    default: ''
  },
  modeId: {
    type: [String, Number],
    default: ''
  },
  placeholder: {
    type: String,
    default: '请选择嵌入模型'
  }
})

const value = ref()

watch(
  [() => props.modeId, () => props.modeName],
  ([newModeId, newModeName]) => {
    if (!newModeId || !newModeName) {
      value.value = undefined
    } else {
      value.value = newModeId + '-' + newModeName
    }
  },
  {
    immediate: true
  }
)

const modelList = ref([])

const selectItem = computed(() => {
  if (props.modeId) {
    return modelList.value.find((item) => item.key == props.modeId)
  }
  return {}
})

const selectChildItem = computed(() => {
  if (selectItem.value.children && selectItem.value.children.length) {
    return selectItem.value.children.find((item) => item.use_model_name == props.modeName)
  }
  return {}
})

const handleChangeModel = (val, option) => {
  if (!option) {
    return
  }
  emit('update:modeName', option.modelName)
  emit('update:modeId', option.modelId)
  emit('change', val, option)
}
const getModelList = () => {
  getModelConfigOption({
    model_type: props.modelType
  }).then((res) => {
    let list = res.data || []
    let { newList, choosableThinking } = getModelOptionsList(list)
    modelList.value = newList
    emit('loaded', modelList.value, choosableThinking)
  })
}

onMounted(() => {
  getModelList(false)

  if (props.modeId && props.modeName) {
    value.value = props.modeId + '-' + props.modeName
  }
})
</script>

<style lang="less" scoped>
.model-icon {
  width: 18px;
}
.modal-seclet-new.ant-select {
  &::v-deep(.ant-select-selection-item) {
    opacity: 0;
    position: relative;
    z-index: 9;
  }
  &::v-deep(.ant-select-clear) {
    opacity: 1 !important;
    width: fit-content;
    height: 22px;
    left: 12px;
    top: 5px;
    margin: 0;
    left: 0;
    right: 0;
    width: 100%;
    text-align: left;
    background: unset;
    font-size: 14px;
    padding: 0 12px;
    .model-show-block {
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      gap: 6px;
      .icon {
        height: 100%;
        width: auto;
      }
      .name-text {
        flex: 1;
        font-size: 14px;
        line-height: 22px;
        color: #000;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
}
</style>
