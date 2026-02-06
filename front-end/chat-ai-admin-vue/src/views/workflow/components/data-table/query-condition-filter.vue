<template>
  <div class="query-condition-filter-warpper" :class="{ 'is-multiple': isMultiple }">
    <div class="query-condition-filter">
      <div class="auxiliary-line"></div>
      <div class="condition-select-box">
        <div class="select-wrapper">
          <a-dropdown>
            <a class="ant-dropdown-link" @click.prevent>
              {{ state.type == 1 ? t('label_and') : t('label_or') }}
              <DownOutlined />
            </a>
            <template #overlay>
              <a-menu style="width: 100px" @click="changeType">
                <a-menu-item :key="1">
                  <a href="javascript:;">{{ t('label_and') }}</a>
                </a-menu-item>
                <a-menu-item :key="2">
                  <a href="javascript:;">{{ t('label_or') }}</a>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </div>

      <div class="field-items">
        <div class="field-item" v-for="(row, index) in state.list" :key="row.field">
          <div class="field-select-box">
            <!-- 绑定字段 -->
            <a-select
              style="width: 100px"
              v-model:value="row.form_field_id"
              @change="(value, option) => selectField(index, option)"
            >
              <a-select-option :value="0" type="integer" name="ID" key="id">ID</a-select-option>
              <a-select-option
                :value="item.id"
                :type="item.type"
                :name="item.name"
                :key="item.id"
                v-for="item in props.fieldOptions"
              >
                {{ item.name }}
              </a-select-option>
            </a-select>
          </div>
          <div class="operator-select-box">
            <!-- 绑定条件 -->
            <a-select
              v-model:value="row.rule"
              style="width: 150px"
              @change="(value, option) => selectRule(index, option)"
            >
              <a-select-option
                :value="item.value"
                style="width: 220px;"
                v-for="item in operatorOptionsMap[row.field_type]"
                :key="item.value"
                >{{ item.label }}</a-select-option
              >
            </a-select>
          </div>
          <div class="field-value-box">
            <!-- 绑定条件值 -->
            <!-- 介于范围输入 -->
            <div class="site-input-group-wrapper" v-if="row.rule === 'between'">
              <div class="site-input-left"  :title="t('ph_input_number_variable')">
                <AtInput
                  :options="atInputOptions"
                  :defaultValue="row.rule_value1"
                  :defaultSelectedList="row.atTags"
                  @open="showAtList"
                  @change="
                    (text, selectedList) => changeAtInputValue(text, selectedList, row, index)
                  "
                  :placeholder="t('ph_input_number')"
                />
              </div>

              <span class="site-input-split">~</span>
              <div class="site-input-right" :title="t('ph_input_number_variable')">
                <AtInput
                  :options="atInputOptions"
                  :defaultValue="row.rule_value2"
                  :defaultSelectedList="row.atTags2"
                  @open="showAtList"
                  @change="
                    (text, selectedList) => changeAtInputValue2(text, selectedList, row, index)
                  "
                  :placeholder="t('ph_input_number')"
                />
              </div>
            </div>

            <!-- 布尔类型选择 -->
            <!-- <a-select
              v-else-if="row.field_type === 'boolean'"
              :value="row.rule_value1"
              @change="value => updateRuleValue(index, 'rule_value1', value)"
              style="width: 100%"
              :placeholder="t('ph_select')"
            >
              <a-select-option value="true">是/真</a-select-option>
              <a-select-option value="false">否/假</a-select-option>
            </a-select> -->

            <!-- 空值检查、布尔类型（无需输入） -->
            <a-input
              v-else-if="['empty', 'not_empty', 'boolean'].includes(row.field_type)"
              :placeholder="t('ph_no_value_required')"
              disabled
            />

            <!-- 数字类型输入 -->
            <!-- <a-input
              v-else-if="['integer', 'number'].includes(row.field_type)"
              :value="row.rule_value1"
              @input="e => updateRuleValueByNumber(index, 'rule_value1', e)"
              style="width: 100%"
              :placeholder="t('ph_input_number')"
            /> -->

            <AtInput
              v-else-if="['integer', 'number'].includes(row.field_type)"
              :options="atInputOptions"
              :defaultValue="row.rule_value1"
              :defaultSelectedList="row.atTags"
              @open="showAtList"
              @change="(text, selectedList) => changeAtInputValue(text, selectedList, row, index)"
              :placeholder="t('ph_input_number_variable')"
            />

            <!-- 文本类型输入 -->
            <AtInput
              v-else
              :options="atInputOptions"
              :defaultValue="row.rule_value1"
              :defaultSelectedList="row.atTags"
              @open="showAtList"
              @change="(text, selectedList) => changeAtInputValue(text, selectedList, row, index)"
              :placeholder="t('ph_input_value_variable')"
            />
            <!-- <a-input
              v-else
              :value="row.rule_value1"
              @input="e => updateRuleValue(index, 'rule_value1', e)"
              :placeholder="请输入文本内容"
            /> -->
          </div>
          <span class="field-del-btn" @click="removeCondition(index)">
            <svg-icon class="del-icon" name="close-circle"></svg-icon>
          </span>
        </div>
      </div>
    </div>

    <div class="add-btn-box">
      <a-tooltip :title="t('ph_select_database_first')" style="width: 100%" v-if="disabled">
        <span>
          <a-button class="add-btn" type="dashed" disabled block>
            <PlusOutlined /> {{ t('btn_add_condition') }}
          </a-button>
        </span>
      </a-tooltip>

      <a-button class="add-btn" type="dashed" block @click="addCondition" v-else>
        <PlusOutlined /> {{ t('btn_add_condition') }}
      </a-button>
    </div>
  </div>
</template>

<script setup>
import { computed, inject, ref, watch, reactive, onMounted } from 'vue'
import { PlusOutlined, DownOutlined } from '@ant-design/icons-vue'
import { getFieldTypeRules, getFilterRulesByType } from '@/constants/database'
import AtInput from '../at-input/at-input.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.data-table.query-condition-filter')

const emit = defineEmits(['change', 'changeType'])

const props = defineProps({
  where: {
    type: Array,
    default: () => []
  },
  fieldOptions: {
    type: Array,
    default: () => []
  },
  disabled: {
    type: Boolean,
    default: false
  },
  type: {
    type: [Number, String],
    default: 1
  }
})

const getNode = inject('getNode')

const state = reactive({
  list: props.where,
  type: props.type
})

watch(
  () => props.where,
  (newVal) => {
    state.list = newVal
  }
)

watch(
  () => props.type,
  (newVal) => {
    state.type = newVal
  }
)

const isMultiple = computed(() => {
  return props.where.length > 1
})

// 动态获取操作符选项的计算属性
const operatorOptionsMap = computed(() => {
  return getFieldTypeRules()
})

const atInputOptions = ref([])

const getAtInputOptions = () => {
  let options = getNode().getAllParentVariable();

  atInputOptions.value = options || []
}

const showAtList = () => {
  getAtInputOptions()
}

const change = () => {
  let data = {
    type: state.type,
    list: state.list
  }

  emit('change', data)
}

// eslint-disable-next-line no-unused-vars
const updateRuleValue = (index, key, e) => {
  let value = e.target.value

  state.list[index][key] = value
  change()
}

// eslint-disable-next-line no-unused-vars
const updateRuleValueByNumber = (index, key, e) => {
  let value = e.target.value
  // 格式化数据，允许输入数字、小数点和负号
  // 只保留数字、小数点和负号，其他字符都删除
  value = value.replace(/[^\d.-]/g, '')

  // 确保负号只能出现在开头
  if (value.indexOf('-') > 0) {
    value = value.replace(/-/g, '')
  }

  // 确保只有一个小数点
  const parts = value.split('.')
  if (parts.length > 2) {
    value = parts[0] + '.' + parts.slice(1).join('')
  }

  // 确保负号在最前面
  if (value.includes('-') && !value.startsWith('-')) {
    value = '-' + value.replace(/-/g, '')
  }

  state.list[index][key] = value
  change()
}

const changeAtInputValue = (text, selectedList, item, index) => {
  // item.value = text
  state.list[index].rule_value1 = text
  state.list[index].atTags = selectedList

  change()
}

const changeAtInputValue2 = (text, selectedList, item, index) => {
  // item.value = text
  state.list[index].rule_value2 = text
  state.list[index].atTags2 = selectedList

  change()
}

const addCondition = () => {
  let rowData = {
    form_field_id: '',
    field_name: '',
    field_type: '',
    rule: '',
    rule_value1: '',
    rule_value2: '',
    atTags: [],
    atTags2: []
  }

  state.list.push(rowData)

  change()
}

const removeCondition = (index) => {
  state.list.splice(index, 1)
  change()
}

const selectField = (index, option) => {
  let ruleRows = getFilterRulesByType(option.type)

  let data = {
    form_field_id: option.value,
    field_name: option.name,
    field_type: option.type,
    rule: ruleRows[0].value,
    rule_value1: '',
    rule_value2: '',
    atTags: [],
    atTags2: []
  }

  state.list[index] = data

  change()
}

const selectRule = (index, option) => {
  let data = state.list[index]
  data.rule = option.value
  data.rule_value1 = ''
  data.rule_value2 = ''
  data.atTags = []
  data.atTags2 = []

  change()
}

const changeType = ({ key }) => {
  state.type = key

  emit('changeType', key)
  change()
}

onMounted(() => {
  getAtInputOptions()
})
</script>

<style lang="less" scoped>
.query-condition-filter-warpper {
  position: relative;
  .query-condition-filter {
    display: flex;
    position: relative;
    width: 100%;

    .condition-select-box {
      display: none;
      position: absolute;
      left: 0;
      top: 50%;
      transform: translateY(-50%) translateX(-50%);
      background: #f2f4f7;
      z-index: 10;
    }
    .auxiliary-line {
      display: none;
      position: absolute;
      left: 0;
      top: 15px;
      bottom: 15px;
      width: 48px;
      border-radius: 8px;
      border: 1px solid #bfbfbf;
      border-right: 0;
      border-top-right-radius: 0;
      border-bottom-right-radius: 0;
    }
  }

  &.is-multiple {
    padding-left: 22px;
    .query-condition-filter {
      padding-left: 32px;
    }
    .auxiliary-line {
      display: block;
    }
    .condition-select-box {
      display: block;
    }
    .add-btn-box {
      padding-left: 32px;
    }
  }

  .field-items {
    flex: 1;
    .field-item {
      display: flex;
      align-items: center;
      flex-wrap: nowrap;
      margin-bottom: 8px;

      &:last-child {
        margin-bottom: 0;
      }
    }
    .field-select-box,
    .operator-select-box {
      flex: 1;
      margin-right: 8px;
    }
    .field-value-box {
      width: 178px;
      overflow: hidden;
    }
    .field-del-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 24px;
      height: 24px;
      margin-left: 8px;
      font-size: 16px;
      color: #595959;
      cursor: pointer;
    }
  }

  .add-btn-box {
    margin-top: 8px;
    .add-btn {
      width: 100%;
    }
  }

  .site-input-group-wrapper {
    display: flex;

    .site-input-split {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 30px;
      border-bottom: 1px solid #d9d9d9;
      border-top: 1px solid #d9d9d9;
      pointer-events: none;
      background-color: #fff;
    }

    .site-input-left,
    .site-input-right {
      flex: 1;
      text-align: center;
      overflow: hidden;
    }

    .site-input-left :deep(.mention-input-warpper) {
      border-top-right-radius: 0;
      border-bottom-right-radius: 0;
    }

    .site-input-right :deep(.mention-input-warpper) {
      border-top-left-radius: 0;
      border-bottom-left-radius: 0;
    }
  }
}
</style>