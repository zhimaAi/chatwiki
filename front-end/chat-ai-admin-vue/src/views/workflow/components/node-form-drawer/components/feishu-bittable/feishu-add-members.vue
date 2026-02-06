<template>
  <div>
      <div class="options-item is-required">
      <div class="options-item-tit">
        <div class="option-label">{{ t('label_select_role') }}</div>
        <div class="option-type">string</div>
      </div>
      <div class="min-input">
        <AtInput
          type="textarea"
          inputStyle="height: 64px;"
          :options="variableOptions"
          :defaultSelectedList="state?.tag_map?.role_id || []"
          :defaultValue="state.role_id"
          ref="atInputRef"
          @open="emit('updateVar')"
          @change="(val, tags) => changeValue('role_id', val, tags)"
          :placeholder="t('ph_input_content')"
        >
          <template #option="{ label, payload }">
            <div class="field-list-item">
              <div class="field-label">{{ label }}</div>
              <div class="field-type">{{ payload.typ }}</div>
            </div>
          </template>
        </AtInput>
      </div>
    </div>
    <div class="options-item is-required">
      <div class="options-item-tit">
        <div class="option-label">{{ t('label_member_id_type') }}</div>
        <div class="option-type">string</div>
      </div>
      <div class="min-input">
        <a-select v-model:value="state.member_id_type" :placeholder="t('ph_select')" style="width: 100%;">
          <a-select-option v-for="key in roleIdType" :key="key">{{key}}</a-select-option>
        </a-select>
      </div>
      <div class="desc">{{ t('desc_member_id_type') }}</div>
    </div>
    <div class="options-item is-required">
      <div class="options-item-tit">
        <div class="option-label">{{ t('label_member_id') }}</div>
        <div class="option-type">string</div>
      </div>
      <div class="min-input">
        <AtInput
          type="textarea"
          inputStyle="height: 64px;"
          :options="variableOptions"
          :defaultSelectedList="state?.tag_map?.member_id || []"
          :defaultValue="state.member_id"
          ref="atInputRef"
          @open="emit('updateVar')"
          @change="(val, tags) => changeValue('member_id', val, tags)"
          :placeholder="t('ph_input_content')"
        >
          <template #option="{ label, payload }">
            <div class="field-list-item">
              <div class="field-label">{{ label }}</div>
              <div class="field-type">{{ payload.typ }}</div>
            </div>
          </template>
        </AtInput>
      </div>
      <div class="desc" v-html="t('desc_member_id_link')"></div>
    </div>
  </div>
</template>

<script setup>
import {ref, reactive} from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import {runPlugin} from "@/api/plugins/index.js";

const { t } = useI18n('views.workflow.components.node-form-drawer.components.feishu-bittable.feishu-add-members')

const props = defineProps({
  config: {
    type: Object,
    default: () => ({
      app_id: '',
      app_secret: '',
      app_token: '',
    })
  },
  variableOptions: {
    type: Array,
  }
})
const emit = defineEmits(['update', 'updateVar'])
//const roleIdType = ['open_id', 'union_id', 'user_id', 'chat_id', 'department_id', 'open_department_id']
const roleIdType = ['user_id']

const atInputRef = ref(null)
const roles = ref([])
const state = reactive({
  member_id: "",
  role_id: undefined,
  role_name: '',
  member_id_type: "user_id",
  tag_map: {},
})

function init(nodeParams = null) {
  if (nodeParams && nodeParams?.plugin?.params?.arguments) {
    Object.assign(state, nodeParams.plugin.params.arguments)
  }
}

function changeValue(field, val, tags) {
  state[field] = val
  state.tag_map[field] = tags
  update()
}

function update() {
  emit('update', JSON.parse(JSON.stringify(state)))
}

defineExpose({
  init,
  update
})
</script>

<style scoped lang="less">
@import "common";

.flex-between {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
