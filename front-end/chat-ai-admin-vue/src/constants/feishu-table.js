import { useI18n } from '@/hooks/web/useI18n'

// 获取带翻译的操作符列表
export const getFeiShuOperator = () => {
  const { t } = useI18n('constants.database')
  return [
    {key: "is", label: t('label_is')},
    {key: "isNot", label: t('label_is_not')},
    {key: "contains", label: t('label_contain')},
    {key: "doesNotContain", label: t('label_not_contain')},
    {key: "isEmpty", label: t('label_empty')},
    {key: "isNotEmpty", label: t('label_not_empty')},
    {key: "isGreater", label: t('label_greater_than')},
    {key: "isGreaterEqual", label: t('label_greater_than_or_equal')},
    {key: "isLess", label: t('label_less_than')},
    {key: "isLessEqual", label: t('label_less_than_or_equal')},
    // {key: "like", label: "LIKE"},
    // {key: "in", label: "IN"}
  ]
}

// 带翻译的操作符列表常量
export const FeiShuOperator = getFeiShuOperator()

// 保持向后兼容的映射函数（使用原始中文 label）
export function FeiShuOperatorMap() {
  let map = {}
  for (let item of FeiShuOperator) {
    map[item.key] = item
  }
  return map
}

// 获取带翻译的操作符映射
export const getFeiShuOperatorMap = () => {
  let map = {}
  for (let item of getFeiShuOperator()) {
    map[item.key] = item
  }
  return map
}

// 带翻译的操作符映射常量
export const feiShuOperatorMap = getFeiShuOperatorMap()

export const ShowFieldTypes = ['Text', 'SingleSelect', 'MultiSelect', 'DateTime', 'Number', 'Checkbox']

export const FieldTypesOperatorMap = {
  Text: ['is', 'isNot', 'isEmpty', 'isNotEmpty', 'contains', 'doesNotContain'],
  Number: ['is', 'isNot', 'isEmpty', 'isNotEmpty', 'isGreater', 'isGreaterEqual', 'isLess', 'isLessEqual'],
  SingleSelect: ['is', 'isNot', 'isEmpty', 'isNotEmpty', 'contains', 'doesNotContain'],
  MultiSelect: ['is', 'isNot', 'isEmpty', 'isNotEmpty', 'contains', 'doesNotContain'],
  DateTime: ['is', 'isEmpty', 'isNotEmpty', 'isGreater', 'isLess'],
  Checkbox: ['is']
}

export function getFieldOperator(type) {
  let arr = FieldTypesOperatorMap[type] || []
  return FeiShuOperator.filter(i => arr.includes(i.key))
}

export const BatchActions = [
  'batch_create_record',
  'batch_delete_record',
  'batch_get_record',
  'batch_update_record',
]

export function getBatchActionParams(params) {
  const filterKeys = ['app_id', 'app_secret', 'app_token', 'table_id']
  const formState = {}
  for (let key in params) {
    if (!filterKeys.includes(key) && params[key].required) formState[key] = JSON.parse(JSON.stringify(params[key]))
  }
  return formState
}
