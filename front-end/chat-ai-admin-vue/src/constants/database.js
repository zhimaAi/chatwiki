import { useI18n } from '@/hooks/web/useI18n'

// 基础过滤规则定义
const FILTER_RULES = {
  // 相等性规则
  EQUALITY: {
    EQ: { key: 'label_equal', value: 'eq' },
    NEQ: { key: 'label_not_equal', value: 'neq' },
    IS: { key: 'label_is', value: 'eq' },
    IS_NOT: { key: 'label_is_not', value: 'neq' }
  },

  // 比较规则
  COMPARISON: {
    GT: { key: 'label_greater_than', value: 'gt' },
    LT: { key: 'label_less_than', value: 'lt' },
    GTE: { key: 'label_greater_than_or_equal', value: 'gte' },
    LTE: { key: 'label_less_than_or_equal', value: 'lte' },
    BETWEEN: { key: 'label_between', value: 'between' }
  },

  // 文本规则
  TEXT: {
    CONTAIN: { key: 'label_contain', value: 'contain' },
    NOT_CONTAIN: { key: 'label_not_contain', value: 'not_contain' }
  },

  // 空值规则
  NULLABILITY: {
    EMPTY: { key: 'label_empty', value: 'empty' },
    NOT_EMPTY: { key: 'label_not_empty', value: 'not_empty' }
  },

  // 布尔规则
  BOOLEAN: {
    TRUE: { key: 'label_is', value: 'true' },
    FALSE: { key: 'label_is_not', value: 'false' }
  }
}

// 获取带翻译的过滤规则
export const getFilterRules = () => {
  const { t } = useI18n('constants.database')
  return {
    EQUALITY: {
      EQ: { label: t('label_equal'), value: 'eq' },
      NEQ: { label: t('label_not_equal'), value: 'neq' },
      IS: { label: t('label_is'), value: 'eq' },
      IS_NOT: { label: t('label_is_not'), value: 'neq' }
    },
    COMPARISON: {
      GT: { label: t('label_greater_than'), value: 'gt' },
      LT: { label: t('label_less_than'), value: 'lt' },
      GTE: { label: t('label_greater_than_or_equal'), value: 'gte' },
      LTE: { label: t('label_less_than_or_equal'), value: 'lte' },
      BETWEEN: { label: t('label_between'), value: 'between' }
    },
    TEXT: {
      CONTAIN: { label: t('label_contain'), value: 'contain' },
      NOT_CONTAIN: { label: t('label_not_contain'), value: 'not_contain' }
    },
    NULLABILITY: {
      EMPTY: { label: t('label_empty'), value: 'empty' },
      NOT_EMPTY: { label: t('label_not_empty'), value: 'not_empty' }
    },
    BOOLEAN: {
      TRUE: { label: t('label_is'), value: 'true' },
      FALSE: { label: t('label_is_not'), value: 'false' }
    }
  }
}

export const filterRules = getFilterRules()

// 字段类型与规则映射
export const getFieldTypeRules = () => {
  const rules = getFilterRules()

  return {
    string: [
      rules.EQUALITY.IS,
      rules.EQUALITY.IS_NOT,
      rules.TEXT.CONTAIN,
      rules.TEXT.NOT_CONTAIN,
      rules.NULLABILITY.EMPTY,
      rules.NULLABILITY.NOT_EMPTY
    ],

    integer: [
      rules.COMPARISON.GT,
      rules.COMPARISON.LT,
      rules.COMPARISON.GTE,
      rules.COMPARISON.LTE,
      rules.EQUALITY.EQ,
      rules.COMPARISON.BETWEEN
    ],

    number: [
      rules.COMPARISON.GT,
      rules.COMPARISON.LT,
      rules.COMPARISON.GTE,
      rules.COMPARISON.LTE,
      rules.EQUALITY.EQ,
      rules.COMPARISON.BETWEEN
    ],

    boolean: [
      rules.BOOLEAN.TRUE,
      rules.BOOLEAN.FALSE
    ]
  }
}

export const fieldTypeRules = getFieldTypeRules()

// 导出基础规则常量，便于其他模块使用
export { FILTER_RULES }

// 工具函数：根据字段类型获取可用规则
export const getFilterRulesByType = (fieldType) => {
  const rules = getFieldTypeRules()
  return rules[fieldType] || []
}

// 工具函数：根据规则值获取规则标签
export const getFilterRuleLabel = (ruleValue, fieldType) => {
  const rules = getFilterRulesByType(fieldType)
  const rule = rules.find(r => r.value === ruleValue)
  return rule?.label || ruleValue
}