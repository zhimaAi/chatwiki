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
    width: 25%;
  }
  .field-value-col,
  .field-value-head {
    flex: 1;
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
        <div class="field-list-col field-list-col-head field-value-head">{{ t('label_field_value') }}</div>
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
            <a-input :disabled="true" :placeholder="t('ph_input_value_variable')" />
          </a-tooltip>
        </div>
      </div>

      <div class="field-list-row">
        <div class="field-list-col field-name-col">jz_firstname</div>
        <div class="field-list-col field-type-col">string</div>
        <div class="field-list-col field-value-col">
          <a-input :placeholder="t('ph_input_value_variable')" />
        </div>
        <div class="field-list-col field-del-col" v-if="props.showDelete">
          <span class="del-btn"><svg-icon class="del-icon" name="close-circle"></svg-icon></span>
        </div>
      </div>
    </div>

    <div class="add-btn-box" v-if="props.showAdd">
      <a-button class="add-btn" type="dashed" block @click="handleAddField"><PlusOutlined /> {{ t('btn_add_update_field') }}</a-button>
    </div>

    <FieldSelectAlert :formId="formId" ref="fieldSelectAlertRef" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import FieldSelectAlert from './field-select-alert.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.data-table.field-list')

const props = defineProps({
  showAdd: {
    type: Boolean,
    default: false
  },
  showDelete: {
    type: Boolean,
    default: true
  },
  showEmptyFieldRow: {
    type: Boolean,
    default: false
  },
  formId: {
    type: String,
    default: ''
  }
})

const fieldSelectAlertRef = ref()

const handleAddField = () => {
  fieldSelectAlertRef.value.open()
}
</script>