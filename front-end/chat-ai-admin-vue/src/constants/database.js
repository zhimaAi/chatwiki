// 基础过滤规则定义
const FILTER_RULES = {
  // 相等性规则
  EQUALITY: {
    EQ: { label: '等于', value: 'eq' },
    NEQ: { label: '不等于', value: 'neq' },
    IS: { label: '是', value: 'eq' },
    IS_NOT: { label: '不是', value: 'neq' }
  },
  
  // 比较规则
  COMPARISON: {
    GT: { label: '大于', value: 'gt' },
    LT: { label: '小于', value: 'lt' },
    GTE: { label: '大于等于', value: 'gte' },
    LTE: { label: '小于等于', value: 'lte' },
    BETWEEN: { label: '介于', value: 'between' }
  },
  
  // 文本规则
  TEXT: {
    CONTAIN: { label: '包含', value: 'contain' },
    NOT_CONTAIN: { label: '不包含', value: 'not_contain' }
  },
  
  // 空值规则
  NULLABILITY: {
    EMPTY: { label: '为空', value: 'empty' },
    NOT_EMPTY: { label: '不为空', value: 'not_empty' }
  },
  
  // 布尔规则
  BOOLEAN: {
    TRUE: { label: '是', value: 'true' },
    FALSE: { label: '不是', value: 'false' }
  }
}

// 字段类型与规则映射
export const FIELD_TYPE_RULES = {
  string: [
    FILTER_RULES.EQUALITY.IS,
    FILTER_RULES.EQUALITY.IS_NOT,
    FILTER_RULES.TEXT.CONTAIN,
    FILTER_RULES.TEXT.NOT_CONTAIN,
    FILTER_RULES.NULLABILITY.EMPTY,
    FILTER_RULES.NULLABILITY.NOT_EMPTY
  ],
  
  integer: [
    FILTER_RULES.COMPARISON.GT,
    FILTER_RULES.COMPARISON.LT,
    FILTER_RULES.COMPARISON.GTE,
    FILTER_RULES.COMPARISON.LTE,
    FILTER_RULES.EQUALITY.EQ,
    FILTER_RULES.COMPARISON.BETWEEN
  ],
  
  number: [
    FILTER_RULES.COMPARISON.GT,
    FILTER_RULES.COMPARISON.LT,
    FILTER_RULES.COMPARISON.GTE,
    FILTER_RULES.COMPARISON.LTE,
    FILTER_RULES.EQUALITY.EQ,
    FILTER_RULES.COMPARISON.BETWEEN
  ],
  
  boolean: [
    FILTER_RULES.BOOLEAN.TRUE,
    FILTER_RULES.BOOLEAN.FALSE
  ]
}

// 导出基础规则常量，便于其他模块使用
export { FILTER_RULES }

// 工具函数：根据字段类型获取可用规则
export const getFilterRulesByType = (fieldType) => {
  return FIELD_TYPE_RULES[fieldType] || []
}

// 工具函数：根据规则值获取规则标签
export const getFilterRuleLabel = (ruleValue, fieldType) => {
  const rules = getFilterRulesByType(fieldType)
  const rule = rules.find(r => r.value === ruleValue)
  return rule?.label || ruleValue
}