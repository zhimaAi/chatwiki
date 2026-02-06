<template>
  <div
    v-for="item in items"
    :key="item.key"
    class="options-item"
    :class="{ 'is-required': item.meta?.required }"
  >
    <div class="options-item-tit">
      <div class="option-label">
        {{ item.meta?.name || item.key }}
        <a-tooltip :title="t('tip_sync_tags')" v-if="item.meta?.tag_component" :overlayStyle="{ maxWidth: '600px' }" placement="top">
          <a @click="syncFieldTags && syncFieldTags(item.key, item.meta)">
            {{ t('btn_sync') }} <a-spin v-if="syncingTags && syncingTags[item.key]" size="small" />
          </a>
        </a-tooltip>
        <a-tooltip v-if="item.meta?.tip" :overlayStyle="{ maxWidth: '600px' }" placement="top">
          <template #title>
            <div class="tip-content">{{ item.meta?.tip }}</div>
          </template>
          <QuestionCircleOutlined />
        </a-tooltip>
      </div>
      <div class="option-type">{{ item.meta?.type || 'string' }}</div>
      <ZmRadioGroup
        v-if="item.meta?.switch_component"
        style="margin-left: 8px;"
        :value="String(formState[item.key] || '').toLowerCase() === 'true' ? 'true' : 'false'"
        :options="switchOptions"
        @change="(val) => {
          onFieldChange && onFieldChange(item.key, val);
          update && update();
        }"
      />
      <div v-if="item.meta?.show_advanced_settings_icon" style="margin-left: auto; margin-bottom: 4px;">
        <a-tooltip :title="t('tip_advanced_settings')">
          <span
            @click="onToggleAdvanced && onToggleAdvanced()"
            :style="{ color: '#2475FC', padding: '4px 8px', borderRadius: '24px', border: '1px solid #D9D9D9', background: '#E5EFFF', cursor: 'pointer' }"
          >
            <SettingOutlined/>
          </span>
        </a-tooltip>
      </div>
    </div>

    <div v-if="item.meta?.select_official_component">
      <a-select
        :value="formState.app_id"
        :placeholder="t('ph_select_official_account')"
        style="width: 100%;"
        @change="(value, option) => {
          onFieldChange && onFieldChange('app_id', value);
          appChange && appChange(value, option);
        }"
      >
        <a-select-option
          v-for="app in apps"
          :key="app.app_id"
          :value="app.app_id"
          :name="app.app_name"
          :app_secret="app.app_secret"
        >
          {{ app.app_name }}
        </a-select-option>
      </a-select>
      <div class="desc" v-html="rednerDesc(item.meta?.desc)"></div>
      <div v-if="errors[item.key]" class="desc" style="color:#FB363F;">{{ errors[item.key] }}</div>
    </div>

    <div v-else-if="item.meta?.enum_component">
      <a-select
        :value="formState[item.key]"
        :placeholder="item.meta?.required ? t('ph_please_select') : t('ph_please_select_optional')"
        style="width: 100%;"
        @change="(value) => {
          onFieldChange && onFieldChange(item.key, value);
          update && update();
        }"
      >
        <a-select-option
          v-for="opt in (sortedEnum ? sortedEnum(item.meta?.enum) : (item.meta?.enum || []))"
          :key="opt.value"
          :value="opt.value"
        >
          {{ opt.name }}
        </a-select-option>
      </a-select>
      <div class="desc" v-html="rednerDesc(item.meta?.desc)"></div>
      <div v-if="errors[item.key]" class="desc" style="color:#FB363F;">{{ errors[item.key] }}</div>
    </div>

    <div v-else-if="item.meta?.radio_component">
      <a-radio-group
        :value="formState[item.key]"
        @change="(e) => {
          const v = e?.target?.value ?? e
          onFieldChange && onFieldChange(item.key, v);
          update && update();
        }"
      >
        <a-radio
          v-for="opt in (sortedEnum ? sortedEnum(item.meta?.enum) : (item.meta?.enum || []))"
          :key="opt.value"
          :value="opt.value"
        >
          {{ opt.name }}
        </a-radio>
      </a-radio-group>
      <div class="desc" v-html="rednerDesc(item.meta?.desc)"></div>
      <div v-if="errors[item.key]" class="desc" style="color:#FB363F;">{{ errors[item.key] }}</div>
    </div>

    <div v-else-if="item.meta?.switch_component">
      <div class="desc" v-html="rednerDesc(item.meta?.desc)"></div>
      <div v-if="errors[item.key]" class="desc" style="color:#FB363F;">{{ errors[item.key] }}</div>
    </div>

    <div v-else-if="item.meta?.date_range_begin_component">
      <div class="tag-box">
        <a-select
          :value="dateRangeTypes[item.meta?.begin_date_key]"
          @change="(val) => dateRangeTypeChange && dateRangeTypeChange(item.meta, val)"
          style="width: 120px; padding-right: 4px;"
        >
          <a-select-option :value="1">{{ t('label_select_date') }}</a-select-option>
          <a-select-option :value="2">{{ t('label_insert_variable') }}</a-select-option>
        </a-select>
        <a-range-picker
          v-if="dateRangeTypes[item.meta?.begin_date_key] == 1"
          :value="getRangeValue ? getRangeValue(item.meta) : null"
          format="YYYY-MM-DD"
          style="width: 100%;"
          @calendarChange="(dates) => onDateRangeCalendarChange && onDateRangeCalendarChange(item.meta, dates)"
          @change="(dates, dateStrings) => onDateRangeChange && onDateRangeChange(item.meta, dates, dateStrings)"
          :disabled-date="(current) => getDisabledDateForRange ? getDisabledDateForRange(item.meta, current) : false"
        />
        <template v-else>
          <AtFullInput
            inputStyle="height: 33px; width: 214px;"
            :options="variableOptions"
            :defaultSelectedList="formState[(item.meta?.begin_date_key || '') + '_tags'] || []"
            :defaultValue="formState[item.meta?.begin_date_key] || ''"
            :canFull="item?.meta?.text_edit_component"
            @open="onOpenVar && onOpenVar()"
            @change="(text, selectedList) => onChangeDynamic && onChangeDynamic(item.meta?.begin_date_key, text, selectedList)"
            :placeholder="t('ph_input_start_date')"
          >
            <template #option="{ label, payload }">
              <div class="field-list-item">
                <div class="field-label">{{ label }}</div>
                <div class="field-type">{{ payload.typ }}</div>
              </div>
            </template>
          </AtFullInput>
          <div class="date-range-separator">-</div>
          <AtFullInput
            inputStyle="height: 33px; width: 214px;"
            :options="variableOptions"
            :defaultSelectedList="formState[(item.meta?.end_date_key || '') + '_tags'] || []"
            :defaultValue="formState[item.meta?.end_date_key] || ''"
            :canFull="item?.meta?.text_edit_component"
            @open="onOpenVar && onOpenVar()"
            @change="(text, selectedList) => onChangeDynamic && onChangeDynamic(item.meta?.end_date_key, text, selectedList)"
            :placeholder="t('ph_input_end_date')"
          >
            <template #option="{ label, payload }">
              <div class="field-list-item">
                <div class="field-label">{{ label }}</div>
                <div class="field-type">{{ payload.typ }}</div>
              </div>
            </template>
          </AtFullInput>
        </template>
      </div>
      <div class="desc" v-html="rednerDesc(item.meta?.desc)"></div>
      <div v-if="errors[item.key]" class="desc" style="color:#FB363F;">{{ errors[item.key] }}</div>
    </div>

    <div v-else-if="item.meta?.tag_component && showTagComponents">
      <div class="tag-box">
        <a-select
          :value="tagTypes[item.key]"
          @change="(val) => tagTypeChange && tagTypeChange(item.key, val, item.meta)"
          style="width: 120px;"
        >
          <a-select-option :value="1">{{ t('label_select_tag') }}</a-select-option>
          <a-select-option :value="2">{{ t('label_insert_variable') }}</a-select-option>
        </a-select>
        <a-select
          v-if="tagTypes[item.key] == 1"
          :value="formState[item.key]"
          :placeholder="t('ph_select_tag')"
          style="width: 100%;"
          @change="(val, opt) => tagChange && tagChange(item.key, val, opt)"
          show-search
          :filter-option="filterOption"
        >
          <a-select-option
            v-for="t in (tagLists[item.key] || [])"
            :key="t.id"
            :value="t.id"
            :name="t.name"
          >
            {{ t.name }}
          </a-select-option>
        </a-select>
        <AtFullInput
          v-else
          type="textarea"
          inputStyle="height: 33px;"
          :options="variableOptions"
          :defaultSelectedList="formState[item.key + '_tags'] || []"
          :defaultValue="formState[item.key] || ''"
          :canFull="item?.meta?.text_edit_component"
          @open="onOpenVar && onOpenVar()"
          @change="(text, selectedList) => onChangeDynamic && onChangeDynamic(item.key, text, selectedList)"
          :placeholder="t('ph_input_content')"
        >
          <template #option="{ label, payload }">
            <div class="field-list-item">
              <div class="field-label">{{ label }}</div>
              <div class="field-type">{{ payload.typ }}</div>
            </div>
          </template>
        </AtFullInput>
      </div>
      <div class="desc" v-html="rednerDesc(item.meta?.desc)"></div>
      <div v-if="errors[item.key]" class="desc" style="color:#FB363F;">{{ errors[item.key] }}</div>
    </div>

    <div v-else-if="item.meta?.dingtalk_user_info_component">
      <DingLoginBox :formState="formState" @change="onChangeDynamic"/>
    </div>

    <div v-else-if="item.meta?.select_remote_component">
      <DingTableSelect
        :formState="formState"
        :variableOptions="variableOptions"
        :canFull="item?.meta?.text_edit_component"
        @open="onOpenVar"
        @change="onChangeDynamic"
      />
    </div>

    <div v-else>
      <AtFullInput
        :type="item.meta?.in === 'body' ? 'textarea' : undefined"
        :inputStyle="item.meta?.in === 'body' ? 'height: 64px;' : undefined"
        :options="variableOptions"
        :defaultSelectedList="formState[item.key + '_tags'] || []"
        :defaultValue="formState[item.key] || ''"
        :canFull="item?.meta?.text_edit_component"
        @open="onOpenVar && onOpenVar()"
        @change="(text, selectedList) => onChangeDynamic && onChangeDynamic(item.key, text, selectedList)"
        :placeholder="t('ph_input_content')"
      >
        <template #option="{ label, payload }">
          <div class="field-list-item">
            <div class="field-label">{{ label }}</div>
            <div class="field-type">{{ payload.typ }}</div>
          </div>
        </template>
      </AtFullInput>
      <div class="desc" v-html="rednerDesc(item.meta?.desc)"></div>
      <div v-if="errors[item.key]" class="desc" style="color:#FB363F;">{{ errors[item.key] }}</div>
    </div>
  </div>
</template>

<script setup>
import { QuestionCircleOutlined, SettingOutlined } from '@ant-design/icons-vue'
import ZmRadioGroup from '@/components/common/zm-radio-group.vue'
import DingLoginBox from "./components/ding-login-box.vue";
import DingTableSelect from "./components/ding-table-select.vue";
import AtFullInput from "@/views/workflow/components/at-input/at-full-input.vue";
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-form-drawer.components.dynamic-api.field-list')

const switchOptions = [
  { label: t('label_on'), value: 'true' },
  { label: t('label_off'), value: 'false' }
]
defineProps({
  actionName: {type: String},
  items: { type: Array, default: () => [] },
  formState: { type: Object, default: () => ({}) },
  errors: { type: Object, default: () => ({}) },
  apps: { type: Array, default: () => [] },
  variableOptions: { type: Array, default: () => [] },
  tagTypes: { type: Object, default: () => ({}) },
  tagLists: { type: Object, default: () => ({}) },
  showTagComponents: { type: Boolean, default: false },
  sortedEnum: { type: Function },
  appChange: { type: Function },
  update: { type: Function },
  onFieldChange: { type: Function },
  onChangeDynamic: { type: Function },
  filterOption: { type: Function },
  tagTypeChange: { type: Function },
  tagChange: { type: Function },
  syncFieldTags: { type: Function },
  syncingTags: { type: Object, default: () => ({}) },
  dateRangeTypes: { type: Object, default: () => ({}) },
  dateRangeTypeChange: { type: Function },
  getRangeValue: { type: Function },
  onDateRangeCalendarChange: { type: Function },
  onDateRangeChange: { type: Function },
  getDisabledDateForRange: { type: Function },
  onOpenVar: { type: Function },
  onToggleAdvanced: { type: Function }
})

function rednerDesc(text) {
  if (!text || typeof text !== "string") {
    return ''
  }
  return text.replace(/\n/g, '<br/')
}
</script>

<style scoped lang="less">
.options-item {
  display: flex;
  flex-direction: column;
  margin-top: 12px;
  line-height: 22px;
  gap: 4px;

  .options-item-tit {
    display: flex;
    align-items: center;
  }

  .option-label {
    color: var(--wf-color-text-1);
    font-size: 14px;
    margin-right: 8px;
  }

  .desc {
    color: var(--wf-color-text-2);
    word-break: break-all;
  }


  &.is-required .option-label::before {
    content: '*';
    color: #FB363F;
    display: inline-block;
    margin-right: 2px;
  }

  .option-type {
    height: 22px;
    line-height: 18px;
    padding: 0 8px;
    border-radius: 6px;
    border: 1px solid rgba(0, 0, 0, 0.15);
    background-color: #fff;
    color: var(--wf-color-text-3);
    font-size: 12px;
  }

  .item-actions-box {
    display: flex;
    align-items: center;

    .action-btn {
      margin-left: 12px;
      font-size: 16px;
      color: #595959;
      cursor: pointer;
    }
  }
}

.tag-box {
  display: flex;
  align-items: center;
  :deep(.mention-input-warpper) {
    height: 33px;
  }
}

.date-range-separator {
  padding: 0 6px;
}
.tip-content {
  white-space: pre-wrap;
}
</style>
