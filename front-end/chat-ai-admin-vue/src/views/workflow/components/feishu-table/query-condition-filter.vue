<template>
  <div class="query-condition-filter-warpper" :class="{ 'is-multiple': isMultiple }">
    <div class="query-condition-filter">
      <div class="auxiliary-line"></div>
      <div class="condition-select-box">
        <div class="select-wrapper">
          <a-dropdown>
            <a class="ant-dropdown-link" @click.prevent>
              {{ state.conjunction == 'and' ? t('label_and') : t('label_or') }}
              <DownOutlined />
            </a>
            <template #overlay>
              <a-menu style="width: 100px" @click="changeType">
                <a-menu-item key="and">
                  <a href="javascript:;">{{ t('label_and') }}</a>
                </a-menu-item>
                <a-menu-item key="or">
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
              v-model:value="row.field_name"
              @change="(value, option) => selectField(index, option)"
            >
              <a-select-option
                :value="item.field_name"
                :type="item.ui_type"
                :key="item.field_id"
                v-for="item in props.fieldOptions"
              >
                {{ item.field_name }}
              </a-select-option>
            </a-select>
          </div>
          <div class="operator-select-box">
            <!-- 绑定条件 -->
            <a-select
              v-model:value="row.operator"
              style="width: 180px"
              @change="() => changeOperator(row, index)"
            >
              <a-select-option
                style="width: 220px;"
                v-for="item in getFieldOperator(row.field_type)"
                :value="item.key"
                :key="item.key"
              >{{ item.label }}</a-select-option>
            </a-select>
          </div>
          <div class="field-value-box">
            <!-- 绑定条件值 -->

            <!-- 空值检查、布尔类型（无需输入） -->
            <a-input
              v-if="['isEmpty', 'isNotEmpty'].includes(row.operator)"
              :placeholder="t('ph_no_value_required')"
              disabled
            />
            <!-- 文本类型输入 -->
            <AtInput
              v-else
              :options="atInputOptions"
              :defaultValue="row.value"
              :defaultSelectedList="row.atTags"
              @open="showAtList"
              @change="(text, selectedList) => changeAtInputValue(text, selectedList, row, index)"
              :placeholder="t('ph_input_value_variable')"
            />
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
import { getFilterRulesByType } from '@/constants/database'
import AtInput from '../at-input/at-input.vue'
import { getFieldOperator } from "@/constants/feishu-table.js";
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.feishu-table.query-condition-filter')

const emit = defineEmits(['change'])

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
  conjunction: {
    type: String,
    default: 'and'
  }
})

const getNode = inject('getNode')

const state = reactive({
  list: [],
  conjunction: 'and'
})

watch(
  () => props.where,
  (newVal) => {
    state.list = newVal
  },
  {immediate: true}
)

watch(
  () => props.conjunction,
  (newVal) => {
    state.conjunction = newVal
  },
  {immediate: true}
)

const isMultiple = computed(() => {
  return props.where.length > 1
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
    conjunction: state.conjunction,
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

const changeAtInputValue = (text, selectedList, item, index) => {
  // item.value = text
  state.list[index].value = text
  state.list[index].atTags = selectedList

  change()
}

const changeOperator = (item, index) => {
  state.list[index].value = ''

  change()
}

const addCondition = () => {
  let rowData = {
    field_id: '',
    field_name: '',
    field_type: '',
    operator: undefined,
    value: '',
    atTags: [],
  }

  state.list.push(rowData)

  change()
}

const removeCondition = (index) => {
  state.list.splice(index, 1)
  change()
}

const selectField = (index, option) => {
  let data = {
    field_id: option.key,
    field_name: option.value,
    field_type: option.type,
    operator: undefined,
    value: '',
    atTags: [],
  }

  state.list[index] = data
  change()
}

const changeType = ({ key }) => {
  state.conjunction = key
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
    padding-left: 20px;
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
      width: 210px;
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