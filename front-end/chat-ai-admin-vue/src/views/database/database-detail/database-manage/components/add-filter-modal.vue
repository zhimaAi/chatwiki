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
          label="分类名称"
          name="name"
          :rules="[{ required: true, message: '请输入分类名称' }]"
        >
          <a-input :maxLength="64" v-model:value="formState.name" placeholder="请输入分类名称" />
        </a-form-item>
        <div class="form-required-item">筛选条件</div>
        <div class="filter-list-box">
          <a-space
            v-for="(item, index) in formState.condition"
            :key="index"
            style="display: flex; margin-bottom: 8px"
            align="baseline"
          >
            <div>条件{{ index + 1 }}:</div>
            <a-form-item>
              <a-select
                v-model:value="item.form_field_id"
                style="width: 150px"
                placeholder="请选择"
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
                message: '请选择规则'
              }"
            >
              <a-select v-model:value="item.rule" style="width: 120px" placeholder="请选择规则">
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
              <a-input v-model:value="item.rule_value1" placeholder="请输入"></a-input>
            </a-form-item>
            <a-form-item
              v-show="is_show_rule_value2(item)"
              :name="['condition', index, 'rule_value2']"
              :rules="{
                asyncValidator: (rule, value) => validatorItem2(rule, value, item),
                trigger: 'blur'
              }"
            >
              <a-input v-model:value="item.rule_value2" placeholder="请输入"></a-input>
            </a-form-item>
            <MinusCircleOutlined v-if="formState.condition.length > 1" @click="removeItem(index)" />
          </a-space>
          <div style="margin-left: 48px">
            <a-button @click="handleAddItem">添加条件</a-button>
          </div>
        </div>
        <a-form-item label="条件之间关系" required>
          <a-radio-group v-model:value="formState.type">
            <a-radio value="1">与（And）</a-radio>
            <a-radio value="2">或（Or）</a-radio>
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
import { addFormFilter, editFormFilter, getFormFilterInfo } from '@/api/database'
import { isNumberOrNumberString } from '@/utils/validate'
const modalTitle = ref('添加分类')
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
    message.error('字段已全部添加')
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
        return reject('请输入')
      }
      if (item.form_filed_type == 'integer') {
        if (!Number.isInteger(+value)) {
          return reject('请输入整数')
        }
        return resolve()
      }
      if (item.form_filed_type == 'number') {
        if (!isNumberOrNumberString(value)) {
          return reject('请输入数字')
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
        return reject('请输入')
      }
      if (item.form_filed_type == 'integer') {
        if (!Number.isInteger(+value)) {
          return reject('请输入整数')
        }
        if (+value < +item.rule_value1) {
          return reject('右边的值需要大于左边的值')
        }
        return resolve()
      }
      if (item.form_filed_type == 'number') {
        if (!isNumberOrNumberString(value)) {
          return reject('请输入数字')
        }
        if (+value < +item.rule_value1) {
          return reject('右边的值需要大于左边的值')
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
  modalTitle.value = '添加分类'
  formatConditions()
  formState.name = ''
  formState.type = '1'
  formState.id = ''
}
const edit = (data) => {
  open.value = true
  modalTitle.value = '编辑分类'
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
    let tip = '添加成功'
    if (formState.id) {
      methodUrl = editFormFilter
      tip = '编辑成功'
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
      label: '是',
      value: 'eq'
    },
    {
      label: '不是',
      value: 'neq'
    },
    {
      label: '包含',
      value: 'contain'
    },
    {
      label: '不包含',
      value: 'not_contain'
    },
    {
      label: '为空',
      value: 'empty'
    },
    {
      label: '不为空',
      value: 'not_empty'
    }
  ],
  integer: [
    {
      label: '大于',
      value: 'gt'
    },
    {
      label: '小于',
      value: 'lt'
    },
    {
      label: '大于等于',
      value: 'gte'
    },
    {
      label: '小于等于',
      value: 'lte'
    },
    {
      label: '等于',
      value: 'eq'
    },
    {
      label: '介于',
      value: 'between'
    }
  ],
  number: [
    {
      label: '大于',
      value: 'gt'
    },
    {
      label: '小于',
      value: 'lt'
    },
    {
      label: '大于等于',
      value: 'gte'
    },
    {
      label: '小于等于',
      value: 'lte'
    },
    {
      label: '等于',
      value: 'eq'
    },
    {
      label: '介于',
      value: 'between'
    }
  ],
  boolean: [
    {
      label: '是',
      value: 'true'
    },
    {
      label: '不是',
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
