<template>
  <div>
    <div class="flex-center min-input">
      <a-select v-model:value="state.set_type" @change="setTypeChange" style="width: 112px;flex-shrink: 0;">
        <a-select-option value="1">选择数据表</a-select-option>
        <a-select-option value="2">插入变量</a-select-option>
      </a-select>
      <a-select
        v-if="state.set_type == 1"
        v-model:value="state.sheetIdOrName"
        @change="tableChange"
        placeholder="请选择数据表"
        style="width: 100%;">
        <a-select-option
          v-for="item in list"
          :key="item.id"
          :name="item.name"
          :value="item.id">
          {{ item.name }}
        </a-select-option>
      </a-select>
      <AtFullInput
        v-else
        type="textarea"
        inputStyle="height: 33px;"
        :options="variableOptions"
        :defaultSelectedList="state.tags || []"
        :defaultValue="state.sheetIdOrName"
        :canFull="canFull"
        ref="atInputRef"
        @open="emit('open')"
        @change="changeValue"
        placeholder="请输入表格ID或名称，键入“/”可以插入变量"
      >
        <template #option="{ label, payload }">
          <div class="field-list-item">
            <div class="field-label">{{ label }}</div>
            <div class="field-type">{{ payload.typ }}</div>
          </div>
        </template>
      </AtFullInput>
    </div>
    <div class="desc">数据表ID或名称，例如：tblz2wmlWiB1JGxS</div>
  </div>
</template>

<script setup>
import {watch, ref, reactive, computed, nextTick} from 'vue'
import {runPlugin} from "@/api/plugins/index.js";
import AtFullInput from "@/views/workflow/components/at-input/at-full-input.vue";
import {isValidURL} from "@/utils/validate.js";

const props = defineProps({
  formState: {type: Object, default: () => null},
  variableOptions: {type: Array, default: () => []},
  canFull:{type: Boolean, default: false}
})
const emit = defineEmits(['open', 'change'])

const atInputRef = ref(null)
const list = ref([])
const state = reactive({
  set_type: '1',
  sheetIdOrName: undefined,
  tags: [],
  table_name: '',
})

const baseConfigReady = computed(() => {
  const {
    baseId,
    operatorId,
    dingtalk_app_key,
    dingtalk_app_secret,
  } = props.formState

  return {
    ok: Boolean(baseId && operatorId && dingtalk_app_key && dingtalk_app_secret),
    config: {
      baseId,
      operatorId,
      dingtalk_app_key,
      dingtalk_app_secret,
    }
  }
})

watch(
  baseConfigReady,
  (val) => {
    list.value = []
    if (val.ok) {
      if(isValidURL(val.config.baseId)) {
        loadRemoteData()
      } else if (state.set_type == 1) {
        setTable()
      }
    } else {
      list.value = []
      // 插入变量时不清空
      if (state.set_type == 1) setTable()
      state.tags = []
      nextTick(() => {
        update()
      })
    }
  },
  { immediate: true }
)

watch(
  () => ({
    set_type: props.formState.extras?.set_type,
    table_name: props.formState.extras?.table_name,
    sheetIdOrName: props.formState.sheetIdOrName,
    sheetIdOrName_tags: props.formState.sheetIdOrName_tags,
  }),
  (val) => {
    state.set_type = val.set_type || '1'
    state.sheetIdOrName = val.sheetIdOrName || undefined
    state.tags = val.sheetIdOrName_tags || []
    state.table_name = val.table_name || ''
  },
  {
    immediate: true,
  }
)

function loadRemoteData() {
  const {dingtalk_app_key, dingtalk_app_secret, operatorId, baseId} = props.formState
  runPlugin({
    name: "dingtalk_ai_table",
    action: "default/exec",
    params: JSON.stringify({
      business: 'AllSheets',
      arguments: {
        dingtalk_app_key,
        dingtalk_app_secret,
        operatorId,
        baseId,
      }
    })
  }).then(res => {
    list.value = res?.data || []
    if (state.set_type == 1) {
      let tbInfo = null
      if (!state.sheetIdOrName) {
        tbInfo = list.value?.[0] || undefined
      } else {
        tbInfo = list.value.find(i => state.sheetIdOrName == i.id)
      }
      setTable(tbInfo)
    }
  }).catch(() => {
    // 插入变量时不清空
    if (state.set_type == 1) setTable()
  }).finally(() => {
    update()
  })
  return true
}

function setTypeChange() {
  let tbInfo = null
  if (state.set_type == '1') tbInfo = list.value[0]
  setTable(tbInfo)
  update()
}


function changeValue(val, tags) {
  state.sheetIdOrName = val
  state.tags = tags
  update()
}

function tableChange(_, opt) {
  state.table_name = opt.name
  update()
}

function setTable(info=null) {
  if (info?.id) {
    state.sheetIdOrName = info.id
    state.table_name = info.name
  } else {
    state.sheetIdOrName = state.set_type == '1' ? undefined : ''
    state.table_name = ''
  }
}

function update() {
  emit('change', 'sheetIdOrName', state.sheetIdOrName, state.tags, {
    set_type: state.set_type,
    table_name: state.table_name,
  })
}
</script>

<style scoped lang="less">
.flex-center {
  display: flex;
  align-items: center;
}

.min-input {
  :deep(.mention-input-warpper) {
    height: 32px;
    word-break: break-all;

    .type-textarea {
      height: 32px;
      min-height: 32px;
    }
  }
}
</style>
