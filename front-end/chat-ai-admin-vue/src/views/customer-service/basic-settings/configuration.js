import { saveRobot } from '@/api/robot'
import { cloneConfig, serializeRobotConfig } from '../config-session'

export const recallFields = [
  'search_type', 'top_k', 'similarity', 'rerank_status', 'rerank_use_model', 'rerank_model_config_id',
  'meta_search_switch', 'meta_search_type', 'meta_search_condition_list', 'rrf_weight',
  'recall_neighbor_switch', 'recall_neighbor_top_k', 'recall_neighbor_before_num', 'recall_neighbor_after_num', 'library_search_type'
]

export function getRecallConfig(robot) {
  const config = Object.fromEntries(recallFields.map(key => [key, cloneConfig(robot[key])]))
  if (typeof config.rrf_weight === 'string') config.rrf_weight = JSON.parse(config.rrf_weight || '{}')
  config.rrf_weight = { vector: 0, search: 0, graph: 0, ...config.rrf_weight }
  return config
}

export function getPromptConfig(robot) {
  return cloneConfig({ prompt: robot.prompt || '', prompt_struct: robot.prompt_struct, prompt_type: robot.prompt_type })
}

export function saveConfigurationPatch(context, patch) {
  // 弹窗不将草稿写入共享表单；取消或保存失败后，其他设置不能顺带保存这些修改。
  const data = serializeRobotConfig({ ...cloneConfig(context.robotInfo), ...cloneConfig(patch) })
  const type = Number(context.robotInfo.application_type || 0)
  return context.saveAction(() => saveRobot(data, type))
}
