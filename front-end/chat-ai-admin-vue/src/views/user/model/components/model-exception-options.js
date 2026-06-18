export function getBackupModelOptionGroups(options = []) {
  return options.map((group) => {
    const models = (group.models || []).map((model) => {
      const name = model.show_model_name || model.use_model
      const value = `${group.model_config_id}::${model.use_model}`

      return {
        value,
        label: `${group.corp} ${name}`,
        model_config_id: group.model_config_id,
        corp: group.corp,
        model_icon_url: group.model_icon_url,
        use_model: model.use_model,
        name
      }
    })

    return {
      model_config_id: group.model_config_id,
      corp: group.corp,
      model_icon_url: group.model_icon_url,
      models
    }
  })
}

export function getBackupModelOptionByValue(groups = [], value) {
  if (!value) {
    return null
  }

  for (const group of groups) {
    const model = (group.models || []).find((item) => item.value === value)
    if (model) {
      return model
    }
  }

  return null
}

export function getModelExceptionAppUrl(robot = {}) {
  const robotId = Number(robot.robot_id || 0)
  if (!robotId) {
    return ''
  }

  if (Number(robot.application_type) === 1) {
    if (!robot.robot_key) {
      return ''
    }
    return `/#/robot/config/workflow?id=${encodeURIComponent(robotId)}&robot_key=${encodeURIComponent(robot.robot_key)}`
  }

  return `/#/robot/config/basic-config?id=${encodeURIComponent(robotId)}`
}
