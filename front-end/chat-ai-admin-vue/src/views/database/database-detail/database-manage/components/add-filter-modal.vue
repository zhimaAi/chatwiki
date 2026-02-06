<template>
  <div class="add-data-sheet">
    <a-modal v-model:open="open" :title="modalTitle" @ok="handleOk" :width="790">
      <a-form
        class="form-box"
        style="margin-top: 16px"
        :model="formState"
        ref="filterFormRef"
        layout="vertical"
      >
        <a-form-item
          :label="t('label_category_name')"
          name="name"
          :rules="[{ required: true, message: t('validator_category_name_required') }]"
        >
          <a-input :maxLength="64" v-model:value="formState.name" :placeholder="t('ph_category_name')" />
        </a-form-item>
        <div class="form-required-item">{{ t('label_filter_condition') }}</div>
        <div class="filter-list-box">
          <a-space
            v-for="(item, index) in formState.condition"
            :key="index"
            style="display: flex; margin-bottom: 8px"
            align="baseline"
          >
            <div>{{ t('condition_label', { index: index + 1 }) }}</div>
            <a-form-item>
              <a-select
                v-model:value="item.form_field_id"
                style="width: 150px"
                :placeholder="t('ph_select')"
                @change="handleChangeFiled(item, index)"
              >
                <a-select-option
                  :disabled="isDisabled(filed, item)"
                  v-for="filed in props.column"
                  :key="filed.id"
                  :value="filed.id"
                  >{{ filed.name }}</a-select-option
                >
              </a-select>
            </a-form-item>
            <a-form-item
              :name="['condition', index, 'rule']"
              :rules="{
                required: true,
                message: t('validator_rule_required')
              }"
            >
              <a-select v-model:value="item.rule" style="width: 120px" :placeholder="t('validator_rule_required')">
                <a-select-option
                  v-for="filed in ruleOptions[item.form_filed_type]"
                  :key="filed.value"
                  :value="filed.value"
                  >{{ filed.label }}</a-select-option
                >
              </a-select>
            </a-form-item>
            <a-form-item
              v-show="is_show_rule_value1(item)"
              :name="['condition', index, 'rule_value1']"
              :rules="{
                asyncValidator: (rule, value) => validatorItem1(rule, value, item),
              }"
            >
              <a-input v-model:value="item.rule_value1" :placeholder="t('ph_input')"></a-input>
            </a-form-item>
            <a-form-item
              v-show="is_show_rule_value2(item)"
              :name="['condition', index, 'rule_value2']"
              :rules="{
                asyncValidator: (rule, value) => validatorItem2(rule, value, item),
                trigger: 'blur'
              }"
            >
              <a-input v-model:value="item.rule_value2" :placeholder="t('ph_input')"></a-input>
            </a-form-item>
            <MinusCircleOutlined v-if="formState.condition.length > 1" @click="removeItem(index)" />
          </a-space>
          <div style="margin-left: 48px">
            <a-button @click="handleAddItem">{{ t('btn_add_condition') }}</a-button>
          </div>
        </div>
        <a-form-item :label="t('label_condition_relation')" required>
          <a-radio-group v-model:value="formState.type">
            <a-radio value="1">{{ t('radio_and') }}</a-radio>
            <a-radio value="2">{{ t('radio_or') }}</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { PlusOutlined, QuestionCircleOutlined, MinusCircleOutlined } from '@ant-design/icons-vue'
import { reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { addFormFilter, editFormFilter, getFormFilterInfo } from '@/api/database'
import { isNumberOrNumberString } from '@/utils/validate'

const { t } = useI18n('views.database.database-detail.database-manage.components.add-filter-modal')

const modalTitle = ref(t('modal_title_add'))
const rotue = useRoute()
const query = rotue.query
const emit = defineEmits(['ok'])
const props = defineProps({
  column: {
    type: Array,
    default: () => []
  }
})

const fieldLists = ref([])

const open = ref(false)
const formState = reactive({
  form_id: query.form_id,
  name: '',
  condition: [],
  type: '1',
  id: ''
})

const formatConditions = () => {
  let condition = [
    {
      form_field_id: props.column[0].id,
      form_filed_type: props.column[0].type,
      rule: void 0,
      rule_value1: '',
      rule_value2: ''
    }
  ]
  formState.condition = condition
}
const removeItem = (index) => {
  formState.condition.splice(index, 1)
}
const handleChangeFiled = (item, index) => {
  // 改变字段类型
  item.rule = void 0
  item.rule_value1 = ''
  item.rule_value2 = ''
  item.form_filed_type = props.column.filter((i) => i.id == item.form_field_id)[0].type
}
const isDisabled = (filed, item) => {
  if (item.form_field_id == filed.id) {
    return false
  }
  let hasAddId = formState.condition.map((i) => i.form_field_id)
  return hasAddId.includes(filed.id)
}

const handleAddItem = () => {
  let hasAddId = formState.condition.map((i) => i.form_field_id)
  let canAddItem = props.column.filter((i) => !hasAddId.includes(i.id))
  if (!canAddItem.length) {
    message.error(t('msg_field_all_added'))
    return
  }
  formState.condition.push({
    form_field_id: canAddItem[0].id,
    form_filed_type: canAddItem[0].type,
    rule: void 0,
    rule_value1: '',
    rule_value2: ''
  })
}

const is_show_rule_value1 = (item) => {
  // 第一个输入框是否显示出来
  if (!item.rule || item.form_filed_type == 'boolean') {
    return false
  }
  let showMap = ['eq', 'neq', 'contain', 'not_contain', 'gt', 'lt', 'gte', 'lte', 'eq', 'between']
  return showMap.includes(item.rule)
}
const is_show_rule_value2 = (item) => {
  // 第二个输入框是否显示出来
  if (!item.rule || item.form_filed_type == 'boolean') {
    return false
  }
  let showMap = ['between']
  return showMap.includes(item.rule)
}
const validatorItem1 = (rule, value, item) => {
  return new Promise((resolve, reject) => {
    if (is_show_rule_value1(item)) {
      if (!value) {
        return reject(t('validator_input_required'))
      }
      if (item.form_filed_type == 'integer') {
        if (!Number.isInteger(+value)) {
          return reject(t('validator_integer'))
        }
        return resolve()
      }
      if (item.form_filed_type == 'number') {
        if (!isNumberOrNumberString(value)) {
          return reject(t('validator_number'))
        }
        return resolve()
      }
      return resolve()
    } else {
      return resolve()
    }
  })
}

const validatorItem2 = (rule, value, item) => {
  return new Promise((resolve, reject) => {
    if (is_show_rule_value2(item)) {
      if (!value) {
        return reject(t('validator_input_required'))
      }
      if (item.form_filed_type == 'integer') {
        if (!Number.isInteger(+value)) {
          return reject(t('validator_integer'))
        }
        if (+value < +item.rule_value1) {
          return reject(t('validator_value2_greater'))
        }
        return resolve()
      }
      if (item.form_filed_type == 'number') {
        if (!isNumberOrNumberString(value)) {
          return reject(t('validator_number'))
        }
        if (+value < +item.rule_value1) {
          return reject(t('validator_value2_greater'))
        }
        return resolve()
      }
      return resolve()
    } else {
      return resolve()
    }
  })
}

const show = () => {
  open.value = true
  modalTitle.value = t('modal_title_add')
  formatConditions()
  formState.name = ''
  formState.type = '1'
  formState.id = ''
}
const edit = (data) => {
  open.value = true
  modalTitle.value = t('modal_title_edit')
  getFormFilterInfo({
    form_id: query.form_id,
    id: data.id
  }).then((res) => {
    let conditions = res.data.conditions || '[]'
    conditions = JSON.parse(conditions)
    conditions = conditions.map((item) => {
      return {
        form_field_id: item.form_field_id + '',
        form_filed_type: props.column.filter((i) => i.id == item.form_field_id)[0].type,
        rule: removePrefix(item.rule),
        rule_value1: item.rule_value1,
        rule_value2: item.rule_value2
      }
    })
    formState.condition = conditions
    formState.name = res.data.name
    formState.type = res.data.type
    formState.id = res.data.id
  })
}
function removePrefix(inputString) {
  // 定义要移除的前缀列表
  const prefixes = ['string_', 'integer_', 'number_', 'boolean_']

  // 循环检查每个前缀
  for (const prefix of prefixes) {
    if (inputString.startsWith(prefix)) {
      // 如果找到匹配的前缀，返回移除前缀后的字符串
      return inputString.slice(prefix.length)
    }
  }

  // 如果没有匹配的前缀，则返回原始字符串
  return inputString
}
const filterFormRef = ref(null)
const handleOk = () => {
  filterFormRef.value.validate().then((res) => {
    let condition = formState.condition.map((item) => {
      return {
        form_field_id: +item.form_field_id,
        rule: item.form_filed_type + '_' + item.rule,
        rule_value1: item.rule_value1,
        rule_value2: item.rule_value2
      }
    })
    let parmas = {
      ...formState,
      condition: JSON.stringify(condition)
    }
    let methodUrl = addFormFilter
    let tip = t('msg_add_success')
    if (formState.id) {
      methodUrl = editFormFilter
      tip = t('msg_edit_success')
    }
    methodUrl(parmas).then((res) => {
      message.success(tip)
      open.value = false
      emit('ok')
    })
  }).catch((err)=>{
  })
}

const ruleOptions = ref({
  string: [
    {
      label: t('rule_eq'),
      value: 'eq'
    },
    {
      label: t('rule_neq'),
      value: 'neq'
    },
    {
      label: t('rule_contain'),
      value: 'contain'
    },
    {
      label: t('rule_not_contain'),
      value: 'not_contain'
    },
    {
      label: t('rule_empty'),
      value: 'empty'
    },
    {
      label: t('rule_not_empty'),
      value: 'not_empty'
    }
  ],
  integer: [
    {
      label: t('rule_gt'),
      value: 'gt'
    },
    {
      label: t('rule_lt'),
      value: 'lt'
    },
    {
      label: t('rule_gte'),
      value: 'gte'
    },
    {
      label: t('rule_lte'),
      value: 'lte'
    },
    {
      label: t('rule_eq'),
      value: 'eq'
    },
    {
      label: t('rule_between'),
      value: 'between'
    }
  ],
  number: [
    {
      label: t('rule_gt'),
      value: 'gt'
    },
    {
      label: t('rule_lt'),
      value: 'lt'
    },
    {
      label: t('rule_gte'),
      value: 'gte'
    },
    {
      label: t('rule_lte'),
      value: 'lte'
    },
    {
      label: t('rule_eq'),
      value: 'eq'
    },
    {
      label: t('rule_between'),
      value: 'between'
    }
  ],
  boolean: [
    {
      label: t('rule_true'),
      value: 'true'
    },
    {
      label: t('rule_false'),
      value: 'false'
    }
  ]
})
defineExpose({
  show,
  edit
})
</script>

<style lang="less" scoped>
.form-box {
  margin-top: 24px;
  .form-required-item {
    margin: 16px 0 8px 0;
    &::before {
      display: inline-block;
      margin-right: 4px;
      color: #ff4d4f;
      font-size: 14px;
      font-family: SimSun, sans-serif;
      line-height: 1;
      content: '*';
    }
  }
  .ant-form-item {
    margin-bottom: 8px;
  }
  .filter-list-box {
    margin-bottom: 16px;
  }
}
</style>
