import { useModelStore } from '@/stores/modules/model'

export function getModelOptionsList(data) {
  let newList = data || []
  let choosableThinking = {}
  newList.forEach((item) => {
    item.key = item.config_info.id
    item.id = item.config_info.id
    item.config_name = item.config_info.config_name
    item.model_config_id = item.config_info.id
    item.children = item.use_model_configs.map((it) => {
      if (it.thinking_type == 2) {
        choosableThinking[item.config_info.id + '#' + it.use_model_name] = true
      }
      return {
        ...it,
        key: item.config_info.id + '_' + it.use_model_name,
        value: item.config_info.id + '_' + it.use_model_name,
        model_config_id: item.config_info.id,
        name: it.use_model_name,
        label: it.show_model_name
      }
    })
  })

  return {
    newList,
    choosableThinking
  }
}

export function getModelNameText(model_config_id, use_model_name, model_type = 'LLM') {
  if (!model_config_id || !use_model_name) {
    return ''
  }
  let modelStore = useModelStore()
  if (modelStore.allModelList.length == 0) {
    // modelStore.getAllmodelList()
    return ''
  }
  let allModelList = modelStore.allModelList
  allModelList = allModelList.map((item) => {
    return {
      ...item,
      use_model_configs: item.use_model_configs.filter((it) => it.model_type == model_type)
    }
  })

  let findConfig = allModelList.find((item) => item.config_info.id == model_config_id)
  let findModel = null
  if (findConfig) {
    findModel = findConfig.use_model_configs.find((item) => item.use_model_name == use_model_name)
  }
  if (findModel) {
    return `${findConfig.config_info.config_name || findConfig.config_info.show_model_name || findConfig.model_name}  ${findModel.show_model_name || findModel.use_model_name}`
  }
  return ''
}
