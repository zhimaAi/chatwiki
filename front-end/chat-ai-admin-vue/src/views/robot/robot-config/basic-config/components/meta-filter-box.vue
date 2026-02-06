<template>
  <div class="query-condition-filter-warpper" :class="{ 'is-multiple': state.list.length > 1 }">
    <div class="query-condition-filter">
      <div class="auxiliary-line"></div>
      <div class="condition-select-box">
        <div class="select-wrapper">
          <a-dropdown @click="changeType">
            <a class="ant-dropdown-link" @click.prevent>
              {{ state.type == 1 ? t('label_and') : t('label_or') }}
              <DownOutlined/>
            </a>
            <template #overlay>
              <a-menu style="width: 100px">
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
        <div class="field-item" v-for="(row, index) in state.list" :key="index">
          <div class="field-select-box">
            <!-- 绑定字段 -->
            <a-select
              v-model:value="row.key"
              :placeholder="t('ph_select')"
              style="width: 200px"
              @change="(val, opt) => fieldChange(row, val, opt)"
            >
              <a-select-option
                v-for="item in props.metaData"
                :value="item.key"
                :key="item.key"
                :type="item.type"
              >{{ item.name }}
              </a-select-option>
            </a-select>
          </div>
          <div class="operator-select-box">
            <!-- 绑定条件 -->
            <a-select
              v-model:value="row.op"
              style="width: 120px"
              :placeholder="t('ph_select')"
              :options="getFieldOptions(row.type)"
              @change="opChange(row)"
            />
          </div>
          <div class="field-value-box">
            <template v-if="![5,6].includes(row.op)">
              <a-date-picker v-if="row.type == 1" v-model:value="row.value" @change="change" format="YYYY-MM-DD HH:mm" style="width: 100%;"/>
              <a-input-number v-else-if="row.type == 2" v-model:value="row.value" :placeholder="t('ph_input')" @blur="change" style="width: 100%;"/>
              <a-input v-else v-model:value.trim="row.value" :placeholder="t('ph_input')" :maxlength="20" @blur="change" style="width: 100%;"/>
            </template>
          </div>
          <span class="field-del-btn" @click="removeCondition(index)">
            <svg-icon class="del-icon" name="close-circle"></svg-icon>
          </span>
        </div>
      </div>
    </div>

    <div v-if="state.list.length < 10" class="add-btn-box">
      <a-button class="add-btn" type="dashed" block @click="addCondition">
        <PlusOutlined/>
        {{ t('btn_add_condition') }}
      </a-button>
    </div>
  </div>
</template>

<script setup>
import {computed, inject, ref, watch, reactive, onMounted} from 'vue'
import dayjs from 'dayjs'
import {PlusOutlined, DownOutlined} from '@ant-design/icons-vue'
import {message} from 'ant-design-vue'
import {jsonDecode} from "@/utils/index.js"
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.basic-config.components.meta-filter-box');

const emit = defineEmits(['change', 'update:rule', 'update:type'])

const props = defineProps({
  rule: {
    type: String,
  },
  type: {
    type: Number,
  },
  metaData: {
    type: Array,
    default: () => []
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

const Operators = computed(() => [
  {
    value: 1,
    label: t('op_is'),
    use_types: [0, 1, 2]
  },
  {
    value: 2,
    label: t('op_is_not'),
    use_types: [0, 1, 2]
  },
  {
    value: 3,
    label: t('op_contains'),
    use_types: [0]
  },
  {
    value: 4,
    label: t('op_not_contains'),
    use_types: [0]
  },
  {
    value: 5,
    label: t('op_is_empty'),
    use_types: [0, 1, 2]
  },
  {
    value: 6,
    label: t('op_is_not_empty'),
    use_types: [0, 1, 2]
  },
  {
    value: 7,
    label: t('op_greater_than'),
    use_types: [1, 2]
  },
  {
    value: 8,
    label: t('op_equals'),
    use_types: [1, 2]
  },
  {
    value: 9,
    label: t('op_less_than'),
    use_types: [1, 2]
  },
  {
    value: 10,
    label: t('op_greater_than_or_equal'),
    use_types: [1, 2]
  },
  {
    value: 11,
    label: t('op_less_than_or_equal'),
    use_types: [1, 2]
  }
])

const state = reactive({
  list: [],
  type: 1
})

watch(() => props.rule, () => {
  let rule = jsonDecode(props.rule, [])
  rule.forEach(item => {
    if (item.type == 1 && item.value > 0) {
      item.value = dayjs(item.value * 1000)
    }
  })
  state.list = rule
}, {
  immediate: true
})

watch(() => props.type, () => {
  state.type = props.type
}, {
  immediate: true
})

const change = () => {
  let data = {
    type: state.type,
    list: state.list
  }
  let rule = JSON.parse(JSON.stringify(state.list))
  rule.forEach(item => {
    if (item.type == 1 && item.value) {
      item.value = dayjs(item.value).startOf('minute').unix()
    }
    item.value = item.value.toString()
  })
  emit('update:rule', JSON.stringify(rule))
  emit('update:type', state.type)
  emit('change', data)
}

const addCondition = () => {
  let rowData = {
    key: undefined,
    type: undefined,
    op: undefined,
    value: '',
  }

  state.list.push(rowData)
  change()
}

const removeCondition = (index) => {
  state.list.splice(index, 1)
  change()
}

const changeType = ({ key }) => {
  state.type = key
  change()
}

function fieldChange(row, _, opt) {
  row.type = opt.type
  row.value = ''
  change()
}

function opChange(row) {
  if ([5,6].includes(row.op)) row.value = ""
  change()
}

function getFieldOptions(type) {
  return Operators.value.filter(item =>
    item.use_types.includes(type)
  )
}

function verify() {
  try {
    for (let item of state.list) {
      if (!item.key) throw t('msg_select_filter_metadata')
      if (!item.op) throw t('msg_select_filter_operator')
      if (item.value && ![5,6].includes(item.op) && !item.value) throw t('msg_select_filter_metadata_value')
    }
    if (!state.list.length) throw t('msg_complete_at_least_one_filter_rule')
    return true
  } catch (e) {
    message.error(e)
    return false
  }
}

defineExpose({
  verify,
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
    padding-left: 30px;

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
      width: 120px;
      margin-right: 8px;
    }

    .field-select-box {
      width: 200px;
    }


    .field-value-box {
      flex: 1;
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
