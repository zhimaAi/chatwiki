<style lang="less" scoped>
@import '../form-block.less';

.node-box {
  .condition-box {
    display: flex;
    align-items: center;
    .condition-left-box {
      margin-right: 4px;
    }
    .connection-text {
      display: inline-block;
      height: 18px;
      line-height: 18px;
      padding: 0 8px;
      font-size: 12px;
      font-weight: 400;
      border-radius: 4px;
      color: #595959;
      background: #e4e6eb;
    }
    .condition-line {
      position: relative;
      height: 100%;
      width: 24px;
      margin-right: 4px;
      &::after {
        content: '';
        position: absolute;
        width: 24px;
        top: 12px;
        bottom: 12px;
        border: 1px solid #d9d9d9;
        border-right: 0;
        border-radius: 2px;
        border-top-right-radius: 0;
        border-bottom-right-radius: 0;
      }
    }

    .field-items {
      .field-item {
        display: flex;
        align-items: center;
        min-height: 24px;
        line-height: 16px;
        padding: 2px 4px;
        margin-bottom: 2px;
        border-radius: 4px;
        border: 1px solid #d9d9d9;

        &:last-child {
          margin-bottom: 0;
        }

        .field-name,
        .field-value {
          line-height: 16px;
          font-size: 12px;
          font-weight: 400;
          color: #595959;
        }
        .field-rule {
          height: 18px;
          line-height: 18px;
          padding: 0 8px;
          margin: 0 4px;
          font-size: 12px;
          font-weight: 400;
          border-radius: 4px;
          color: #595959;
          background: #e4e6eb;
        }
      }
    }
  }

  .options-list {
    flex: 1;
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }
  .options-item {
    display: flex;
    align-items: center;
    height: 22px;
    padding: 2px 2px 2px 4px;
    border-radius: 4px;
    border: 1px solid #d9d9d9;

    &.is-required .option-label::before {
      vertical-align: middle;
      content: '*';
      color: #fb363f;
      margin-right: 2px;
    }

    .option-label {
      color: var(--wf-color-text-3);
      font-size: 12px;
      margin-right: 4px;
    }

    .option-type {
      height: 18px;
      line-height: 18px;
      padding: 0 8px;
      border-radius: 4px;
      font-size: 12px;
      background-color: #e4e6eb;
      color: var(--wf-color-text-3);
    }
  }
}
</style>

<template>
  <node-common
    :properties="properties"
    :title="props.properties.node_name"
    :icon-name="properties.node_icon_name"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    :node_type="properties.node_type"
  >
    <div class="node-box">
      <div class="static-field-list">
        <div class="static-field-item">
          <div class="static-field-item-label">数据表</div>
          <div class="static-field-item-content">
            <div class="static-field-value">
              {{ state.formData.form_name || '--' }}
            </div>
          </div>
        </div>

        <div class="static-field-item" style="align-items: center">
          <div class="static-field-item-label">查询条件</div>
          <div class="static-field-item-content">
            <span class="static-field-value" v-if="state.formData.where.length == 0">--</span>
            <div class="condition-box">
              <div class="condition-left-box" v-if="state.formData.where.length > 1">
                <span class="connection-text">{{ state.formData.typ == 1 ? '且' : '或' }}</span>
              </div>
              <div class="condition-line" v-if="state.formData.where.length > 1"></div>
              <div class="condition-body">
                <div class="field-items">
                  <div class="field-item" v-for="item in state.formData.where" :key="item.id">
                    <span class="field-name">{{ item.field_name }}</span>
                    <span class="field-rule">{{ getRuleLabel(item.field_type, item.rule) }}</span>
                    <span class="field-value">
                      <AtText
                        :options="atInputOptions"
                        :default-value="item.rule_value1"
                        :defaultSelectedList="item.atTags"
                      />
                      <template v-if="item.rule_value2">
                        <span>-</span>
                        <AtText
                          :options="atInputOptions"
                          :default-value="item.rule_value2"
                          :defaultSelectedList="item.atTags2"
                        />
                      </template>
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="static-field-item">
          <div class="static-field-item-label">查询字段</div>
          <div class="static-field-item-content">
            <div class="options-list">
              <div
                class="options-item"
                v-for="item in state.formData.fields"
                :key="item.key"
              >
                <div class="option-label">{{ item.name }}</div>
                <div class="option-type">{{ item.type }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script setup>
import { getFilterRuleLabel } from '@/constants/database'
import { useDataTableStore } from '@/stores/modules/data-table'
import { ref, reactive, inject, onMounted, toRaw, nextTick, watch, onBeforeUnmount } from 'vue'
import NodeCommon from '../base-node.vue'
import AtText from '../../at-input/at-text.vue'


const props = defineProps({
  properties: {
    type: Object,
    default: () => ({})
  },
  isSelected: { type: Boolean, default: false },
  isHovered: { type: Boolean, default: false }
})

const setData = inject('setData')
const resetSize = inject('resetSize')
const getNode = inject('getNode')
const graphModel = inject('getGraph')

const dataTableStore = useDataTableStore()

let node_params = {}

const state = reactive({
  tableList: [],
  formData: {
    form_name: '',
    form_description: '',
    form_id: '',
    typ: 1,
    fields: [],
    where: [],
    order: [],
    limit: 100
  },
  fieldList: []
})

const atInputOptions = ref([])

const getAtInputOptions = () => {
  let options = getNode().getAllParentVariable()

  atInputOptions.value = options || []
}


const init = () => {
  let dataRaw = props.properties.dataRaw || props.properties.node_params || '{}'

  node_params = JSON.parse(dataRaw)

  getAtInputOptions()

  let where = whereDataConditions(node_params.form_select.where || [])

  node_params.form_select.where = where

  state.formData = node_params.form_select

  if (state.formData.form_id) {
    getFieldList(state.formData.form_id)
  }

  update()

  nextTick(() => {
    resetSize()
  })
}

const update = () => {
  let form = { ...toRaw(state.formData) }

  form.where = whereDataFormatter(form.where)

  node_params.form_select = form

  setData({
    ...props.node,
    node_params: JSON.stringify(node_params)
  })
}

const whereDataConditions = (where) => {
  let conditions = []

  where.forEach((item) => {
    let data = JSON.parse(JSON.stringify(item))
    let ruleArr = data.rule.split('_')

    data.form_field_id = data.form_field_id * 1
    data.field_type = ruleArr[0]
    // 删除数组的第一个
    ruleArr.shift()
    // 数组还原成字符串
    data.rule = ruleArr.join('_')

    conditions.push(data)
  })

  return conditions
}

const whereDataFormatter = (where) => {
  let conditions = []

  where.forEach((item) => {
    let data = JSON.parse(JSON.stringify(item))

    data.form_field_id = data.form_field_id * 1
    data.rule = item.field_type + '_' + item.rule

    conditions.push(data)
  })

  return conditions
}

const getTableList = async () => {
  const list = await dataTableStore.getFormList()
  if (list) {
    state.tableList = list
  }
}

const getFieldList = async (form_id) => {
  dataTableStore
    .getFormFieldList({ form_id: form_id })
    .then((list) => {
      list.forEach((item) => {
        item.id = item.id * 1
      })

      state.fieldList = list
    })
    .catch(() => {})
}

const getRuleLabel = (field_type, rule) => {
  return getFilterRuleLabel(rule, field_type)
}

const onUpatateNodeName = () => {
  // console.log(props.properties.node_name, 'onUpatateNodeName')
  init()
}

watch(
  () => props.properties.dataRaw,
  (newVal, oldVal) => {
    if (newVal !== oldVal) {
      init()
    }
  }
)

onMounted(() => {
  getTableList()
  init()

  const mode = graphModel()

  mode.eventCenter.on('custom:setNodeName', onUpatateNodeName)
})

onBeforeUnmount(() => {
  const mode = graphModel()

  mode.eventCenter.off('custom:setNodeName', onUpatateNodeName)
})
</script>
