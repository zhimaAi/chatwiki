<template>
  <a-select
    v-model:value="value"
    :placeholder="placeholder"
    @change="handleChangeModel"
    style="width: 100%"
  >
    <a-select-opt-group v-for="item in modelList" :key="item.key">
      <template #label>
        <a-flex align="center" :gap="6">
          <img class="model-icon" :src="item.icon" alt="" />
          <span>{{ item.name }}</span>
        </a-flex>
      </template>
      <a-select-option
        v-for="sub in item.children"
        :value="sub.value"
        :modelId="sub.model_config_id"
        :modelName="sub.name"
        :key="sub.key"
      >
        {{ sub.label }}
      </a-select-option>
    </a-select-opt-group>
  </a-select>
</template>

<script setup>
import { duplicateRemoval, removeRepeat } from '@/utils/index'
import { getModelConfigOption } from '@/api/model/index'
import { ref, onMounted, watch } from 'vue'

function uniqueArr(arr, arr1, key) {
  const keyVals = new Set(arr.map((item) => item.model_define))
  arr1.filter((obj) => {
    let val = obj[key]
    if (keyVals.has(val)) {
      arr.filter((obj1) => {
        if (obj1.model_define == val) {
          obj1.children = removeRepeat(obj1.children, obj.children)
          return false
        }
      })
    }
  })
  return arr
}

const emit = defineEmits(['change', 'update:modeName', 'update:modeId', 'loaded'])

const props = defineProps({
  modelType: {
    type: String,
    validator: (value) => {
      return ['TEXT EMBEDDING', 'RERAN', 'LLM'].includes(value)
    },
    required: true
  },
  isOffline: {
    type: Boolean,
    default: false
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

watch([() => props.modeId, () => props.modeName], ([newModeId, newModeName]) => {
  if (!newModeId || !newModeName) {
    return
  }
  value.value = newModeId + '-' + newModeName
})

const modelList = ref([])
const modelDefine = ['azure', 'ollama', 'xinference', 'openaiAgent']

const handleChangeModel = (val, option) => {
  emit('update:modeName', option.modelName)
  emit('update:modeId', option.modelId)
  emit('change', val, option)
}

const getModelList = () => {
  getModelConfigOption({
    model_type: props.modelType,
    is_offline: props.isOffline // 0 1 区分线上线下
  }).then((res) => {
    let list = res.data || []

    let newList = list.map((item) => {
      let { model_define, deployment_name, id } = item.model_config
      let key = `${model_define}-${deployment_name || ''}-${id}`

      let subModelList = []
      let children = []

      if (props.modelType == 'TEXT EMBEDDING') {
        subModelList = item.model_info.vector_model_list
      } else if (props.modelType == 'RERANK') {
        subModelList = item.model_info.rerank_model_list
      } else if (props.modelType == 'LLM') {
        subModelList = item.model_info.llm_model_list
      }

      for (let i = 0; i < subModelList.length; i++) {
        let ele = subModelList[i]
        let label = ele
        let value = id + '-' + ele // 这里的值使用父级id加子模型名称的方式来避免重复（因为有多个叫’默认的模型‘）
        // 部分模型显示的模型名称要特殊处理
        if (
          modelDefine.indexOf(item.model_config.model_define) > -1 &&
          item.model_config.deployment_name
        ) {
          label = item.model_config.deployment_name
        }

        children.push({
          key: id + '_' + ele,
          label: label,
          name: ele,
          deployment_name: deployment_name,
          model_config_id: id,
          model_define: model_define,
          value: value
        })
      }

      return {
        key: key,
        model_define: model_define,
        children: children,
        model_config_id: id,
        deployment_name: deployment_name,
        name: item.model_info.model_name,
        icon: item.model_info.model_icon_url
      }
    })

    // 如果modelList存在两个相同model_define情况就合并到一个对象的children中去
    newList = uniqueArr(duplicateRemoval(newList, 'model_define'), newList, 'model_define')

    if (newList.length === 0) {
      value.value = '无可用模型'
    }

    modelList.value = [...newList]

    emit('loaded', JSON.parse(JSON.stringify(newList)))
  })
}

watch(
  () => props.isOffline,
  () => {
    getModelList()
  }
)

onMounted(() => {
  getModelList(false)
})
</script>

<style lang="less" scoped>
.model-icon {
  width: 18px;
}
</style>
