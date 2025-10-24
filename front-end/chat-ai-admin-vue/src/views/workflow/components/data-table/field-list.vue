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
        <div class="field-list-col field-list-col-head field-name-head">字段名</div>
        <div class="field-list-col field-list-col-head field-type-head">类型</div>
        <div class="field-list-col field-list-col-head field-value-head">字段值</div>
        <div
          class="field-list-col field-list-col-head field-del-head"
          v-if="props.showDelete"
        ></div>
      </div>

      <div class="field-list-row" v-if="showEmptyFieldRow">
        <div class="field-list-col field-name-col">--</div>
        <div class="field-list-col field-type-col">--</div>
        <div class="field-list-col field-value-col">
          <a-tooltip title="请先选择数据库">
            <a-input :disabled="true" placeholder="请输入参数值，键入/插入变量" />
          </a-tooltip>
        </div>
      </div>

      <div class="field-list-row">
        <div class="field-list-col field-name-col">jz_firstname</div>
        <div class="field-list-col field-type-col">string</div>
        <div class="field-list-col field-value-col">
          <a-input placeholder="请输入参数值，键入“/”插入变量" />
        </div>
        <div class="field-list-col field-del-col" v-if="props.showDelete">
          <span class="del-btn"><svg-icon class="del-icon" name="close-circle"></svg-icon></span>
        </div>
      </div>
    </div>

    <div class="add-btn-box" v-if="props.showAdd">
      <a-button class="add-btn" type="dashed" block @click="handleAddField"><PlusOutlined /> 添加更新字段</a-button>
    </div>

    <FieldSelectAlert :formId="formId" ref="fieldSelectAlertRef" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import FieldSelectAlert from './field-select-alert.vue'

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
