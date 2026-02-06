<style lang="less" scoped>
.field-list {
  .field-list-row {
    position: relative;
    display: flex;
    flex-wrap: nowrap;
    align-items: center;
    margin-bottom: 4px;

    &:last-child {
      margin-bottom: 0;
    }
  }

  .field-list-col {
    padding: 0 4px;
  }

  .field-name-col,
  .field-type-col,
  .field-name-head,
  .field-type-head {
    flex: 1;
  }

  .field-value-col,
  .field-value-head {
    width: 256px;
  }

  .field-list-col-head {
    line-height: 22px;
    margin-bottom: 4px;
    font-size: 14px;
    color: #262626;
  }

  .field-name-col,
  .field-type-col {
    line-height: 22px;
    font-size: 14px;
    color: #595959;
  }

  .field-del-head,
  .field-del-col {
    width: 24px;
    display: flex;
    align-items: center;
  }

  .field-del-col {
    text-align: right;

    .del-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 16px;
      height: 16px;
      font-size: 16px;
      color: #595959;
      cursor: pointer;
    }
  }
}

.add-btn-box {
  margin-top: 8px;
}
</style>

<template>
  <div>
    <div class="field-list">
      <div class="field-list-row">
        <div class="field-list-col field-list-col-head field-name-head">{{ t('label_field_name') }}</div>
        <div class="field-list-col field-list-col-head field-type-head">{{ t('label_type') }}</div>
        <div class="field-list-col field-list-col-head field-value-head" v-if="showInput">{{ t('label_field_value') }}</div>
        <div
          class="field-list-col field-list-col-head field-del-head"
          v-if="props.showDelete"
        ></div>
      </div>

      <div class="field-list-row" v-if="showEmptyFieldRow">
        <div class="field-list-col field-name-col">--</div>
        <div class="field-list-col field-type-col">--</div>
        <div class="field-list-col field-value-col">
          <a-tooltip :title="t('ph_select_database_first')">
            <a-input :disabled="true" :placeholder="t('ph_input_value_variable')"/>
          </a-tooltip>
        </div>
      </div>

      <div class="field-list-row" v-for="(item, index) in state.list" :key="item.field_name">
        <div class="field-list-col field-name-col">{{ item.field_name }}</div>
        <div class="field-list-col field-type-col">{{ item.ui_type || fieldMap[item.field_name]?.ui_type }}</div>
        <div class="field-list-col field-value-col" v-if="showInput">
<!--          <a-checkbox v-if="item.ui_type === 'Checkbox'" v-model:checked="item.value"/>-->
<!--          <a-select-->
<!--            v-else-if="item.ui_type === 'SingleSelect'"-->
<!--            :placeholder="t('ph_select')"-->
<!--            v-model:value="item.value"-->
<!--            style="width: 100%;"-->
<!--          >-->
<!--            <a-select-option v-for="opt in item.property.options" :value="opt.id" :key="opt.id">{{opt.name}}</a-select-option>-->
<!--          </a-select>-->
<!--          <a-select-->
<!--            v-else-if="item.ui_type === 'MultiSelect'"-->
<!--            :placeholder="请选择（多选）"-->
<!--            v-model:value="item.value"-->
<!--            style="width: 100%;"-->
<!--            mode="multiple"-->
<!--          >-->
<!--            <a-select-option v-for="opt in item.property.options" :value="opt.id" :key="opt.id">{{opt.name}}</a-select-option>-->
<!--          </a-select>-->
          <AtInput
            :options="atInputOptions"
            :defaultValue="item.value"
            :defaultSelectedList="item.atTags"
            @open="showAtList"
            @change="(text, selectedList) => changeAtInputValue(text, selectedList, item, index)"
            :placeholder="formatPlaceholder(item)"/>
        </div>
        <div class="field-list-col field-del-col" v-if="props.showDelete">
          <span class="del-btn" @click="handleDel(index)">
            <svg-icon class="del-icon" name="close-circle"></svg-icon>
          </span>
        </div>
      </div>
    </div>

    <div class="add-btn-box" v-if="props.showAdd">
      <a-button class="add-btn" type="dashed" block @click="handleAddField">
        <PlusOutlined/>
        {{ t('btn_add_field') }}
      </a-button>
    </div>

    <FieldSelectAlert ref="fieldSelectAlertRef" :fields="fields" @ok="onChangeSelect"/>
  </div>
</template>

<script setup>
import {ref, reactive, watch, inject, onMounted, computed} from 'vue'
import {PlusOutlined} from '@ant-design/icons-vue'
import FieldSelectAlert from './field-select-alert.vue'
import AtInput from '../../at-input/at-input.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.feishu-table.field-selector.index')

const emit = defineEmits(['change'])

const props = defineProps({
  showAdd: {
    type: Boolean,
    default: false
  },
  showDelete: {
    type: Boolean,
    default: true
  },
  showInput: {
    type: Boolean,
    default: true
  },
  showEmptyFieldRow: {
    type: Boolean,
    default: false
  },
  fields: {
    type: Array,
    default: () => ([])
  },
  list: {
    type: Array,
    default: () => []
  },
})

const getNode = inject('getNode')

const state = reactive({
  list: [],
  keys: [],
})

const fieldMap = computed(() => {
  return Object.fromEntries(
    props.fields.map(item => [item.field_name, item])
  )
})

watch(() => props.list, (newVal) => {
  state.list = newVal
  state.keys = newVal.map(i => i.field_name)
}, {
  immediate: true
})

const atInputOptions = ref([])

const getAtInputOptions = () => {
  let options = getNode().getAllParentVariable();

  atInputOptions.value = options || []
}

const showAtList = () => {
  getAtInputOptions()
}

const changeAtInputValue = (text, selectedList, item, index) => {
  // item.value = text
  state.list[index].value = text
  state.list[index].atTags = selectedList

  change()
}


const onChangeSelect = (selectedRowKeys, selectedRows) => {
  state.list = [...selectedRows]
  state.keys = [...selectedRowKeys]
  change()
}

const handleDel = (index) => {
  state.list.splice(index, 1)

  change()
}

const change = () => {
  emit('change', [...state.list], [...state.keys])
}

const fieldSelectAlertRef = ref()

const handleAddField = () => {
  fieldSelectAlertRef.value.open({
    selectedRows: props.list,
    selectedRowKeys: props.list.map(item => {
      return item.field_name
    })
  })
}

const formatPlaceholder = item => {
  switch (item.ui_type) {
    case 'Checkbox':
      return t('ph_input_true_false')
    case 'MultiSelect':
      return t('ph_input_multi_select')
      break
  }
  return t('ph_input_value_variable')
}

onMounted(() => {
  getAtInputOptions()
})
</script>