export const FeiShuOperator = [
  {key: "is", label: "是"},
  {key: "isNot", label: "不是"},
  {key: "contains", label: "包含"},
  {key: "doesNotContain", label: "不包含"},
  {key: "isEmpty", label: "为空"},
  {key: "isNotEmpty", label: "不为空"},
  {key: "isGreater", label: "大于"},
  {key: "isGreaterEqual", label: "大于等于"},
  {key: "isLess", label: "小于"},
  {key: "isLessEqual", label: "小于等于"},
  // {key: "like", label: "LIKE"},
  // {key: "in", label: "IN"}
]

export function FeiShuOperatorMap() {
  let map = {}
  for (let item of FeiShuOperator) {
    map[item.key] = item
  }
  return map
}

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
