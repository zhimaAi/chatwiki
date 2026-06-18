import { describe, expect, it } from 'vitest'
import {
  getModelExceptionAppUrl,
  getBackupModelOptionByValue,
  getBackupModelOptionGroups
} from './model-exception-options'

describe('backup model option groups', () => {
  it('keeps provider icon and builds encoded values for model options', () => {
    const groups = getBackupModelOptionGroups([
      {
        model_config_id: 12,
        corp: 'DeepSeek',
        model_icon_url: '/upload/model_icon/deepseek.png',
        models: [
          {
            use_model: 'deepseek-chat',
            show_model_name: 'DeepSeek Chat'
          }
        ]
      }
    ])

    expect(groups).toEqual([
      {
        model_config_id: 12,
        corp: 'DeepSeek',
        model_icon_url: '/upload/model_icon/deepseek.png',
        models: [
          {
            value: '12::deepseek-chat',
            label: 'DeepSeek DeepSeek Chat',
            model_config_id: 12,
            corp: 'DeepSeek',
            model_icon_url: '/upload/model_icon/deepseek.png',
            use_model: 'deepseek-chat',
            name: 'DeepSeek Chat'
          }
        ]
      }
    ])
  })

  it('finds selected model display data by encoded value', () => {
    const groups = getBackupModelOptionGroups([
      {
        model_config_id: 7,
        corp: 'OpenAI',
        model_icon_url: '/upload/model_icon/openai.png',
        models: [{ use_model: 'gpt-4o', show_model_name: '' }]
      }
    ])

    expect(getBackupModelOptionByValue(groups, '7::gpt-4o')).toMatchObject({
      corp: 'OpenAI',
      model_icon_url: '/upload/model_icon/openai.png',
      name: 'gpt-4o'
    })
  })
})

describe('model exception app url', () => {
  it('includes robot_key when linking to workflow config', () => {
    expect(getModelExceptionAppUrl({
      application_type: 1,
      robot_id: 9,
      robot_key: 'wfRobot001'
    })).toBe('/#/robot/config/workflow?id=9&robot_key=wfRobot001')
  })
})
