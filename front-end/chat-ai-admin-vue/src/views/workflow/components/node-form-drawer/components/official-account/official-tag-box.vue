<template>
  <div class="tag-container">
    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/input.svg" class="title-icon"/>{{ t('label_section_input') }}</div>
      </div>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">{{ t('label_official_account') }}</div>
        </div>
        <div>
          <a-select
            v-model:value="formState.app_id"
            :placeholder="t('ph_select_official_account')"
            style="width: 100%;"
            @change="appChange"
          >
            <a-select-option
              v-for="app in apps"
              :key="app.app_id"
              :name="app.app_name"
              :secret="app.app_secret"
            >
              {{ app.app_name }}
            </a-select-option>
          </a-select>
        </div>
      </div>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">{{ t('label_official_account_tag') }}</div>
          <a-tooltip :title="t('tooltip_sync_tags')">
            <a @click="syncTags">{{ t('btn_sync') }} <a-spin v-if="syncing" size="small"/></a>
          </a-tooltip>
        </div>
        <div class="tag-box">
          <a-select
            v-model:value="formState.tag_type"
            @change="tagTypeChange"
            style="width: 120px;"
          >
            <a-select-option :value="1">{{ t('btn_select_tag') }}</a-select-option>
            <a-select-option :value="2">{{ t('btn_insert_variable') }}</a-select-option>
          </a-select>
          <a-select
            v-if="formState.tag_type == 1"
            v-model:value="formState.tagid"
            :placeholder="t('ph_select_tag')"
            style="width: 100%;"
            @change="tagChange"
            show-search
            :filter-option="filterOption"
          >
            <a-select-option
              v-for="item in tags"
              :key="item.id"
              :value="item.id"
              :name="item.name"
            >
              {{ item.name }}
            </a-select-option>
          </a-select>
          <AtInput
            v-else
            type="textarea"
            inputStyle="height: 33px;"
            :options="variableOptions"
            :defaultSelectedList="formState.tag_map?.tagid || []"
            :defaultValue="formState.tagid"
            ref="atInputRef"
            @open="emit('updateVar')"
            @change="(text, selectedList) => changeFieldValue('tagid', text, selectedList)"
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
          <div class="option-label">{{ t('label_fans_openid_list') }}</div>
          <div class="option-type">string</div>
        </div>
        <div>
          <AtInput
            type="textarea"
            inputStyle="height: 64px;"
            :options="variableOptions"
            :defaultSelectedList="formState.tag_map?.openid_list || []"
            :defaultValue="formState.openid_list"
            ref="atInputRef"
            @open="emit('updateVar')"
            @change="(text, selectedList) => changeFieldValue('openid_list', text, selectedList)"
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
        <div class="desc">{{ t('desc_fans_openid_list') }}</div>
      </div>
    </div>
    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/output.svg" class="title-icon"/>{{ t('label_section_output') }}</div>
      </div>
      <div class="options-item">
        <OutputFields :tree-data="outputData"/>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref, reactive, onMounted, watch, inject} from 'vue';
import { useI18n } from 'vue-i18n';
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import {getWechatAppList} from "@/api/robot/index.js";
import {runPlugin} from "@/api/plugins/index.js";
import {pluginOutputToTree} from "@/constants/plugin.js";
import OutputFields from "@/views/workflow/components/feishu-table/output-fields.vue";
import {message} from 'ant-design-vue';

const { t } = useI18n();

const emit = defineEmits(['updateVar'])
const props = defineProps({
  node: {
    type: Object,
    default: () => ({})
  },
  action: {
    type: Object,
  },
  actionName: {
    type: String,
  },
  variableOptions: {
    type: Array,
  }
})

const setData = inject('setData')
const pluginName = 'official_batch_tag'
const outputData = ref([])
const apps = ref([])
const tags = ref([])
const formState = reactive({
  app_id: undefined,
  app_secret: '',
  openid_list: '',
  tagid: '',
  tag_type: 1,
  tag_name: 1,
  tag_map: {}
})
const syncing = ref(false)

onMounted(() => {
  init()
})

watch(() => props.action, () => {
  outputData.value = pluginOutputToTree(JSON.parse(JSON.stringify(props.action?.output || '{}')))
}, {
  deep: true,
  immediate: true
})

function init() {
  loadWxApps()
  nodeParamsAssign()
}

function loadWxApps() {
  getWechatAppList({app_type: 'official_account'}).then(res => {
    apps.value = res?.data || []
  })
}

function nodeParamsAssign() {
  let nodeParams = JSON.parse(props.node.node_params)
  let arg = nodeParams?.plugin?.params?.arguments || {}
  Object.assign(formState, arg)
  if (formState.tagid) {
    if (formState.tagid > 0) {
      formState.tag_type = 1
    } else {
      formState.tag_type = 2
    }
    loadTags()
  }
}

function appChange(_, option) {
  const {key, secret, name} = option
  formState.app_secret = secret
  formState.app_name = name
  loadTags()
}

function syncTags() {
  if (!formState.app_secret || !formState.app_id) {
    return message.warning(t('msg_please_select_official_account'))
  }
  if (syncing.value) return
  syncing.value = true
  loadTags().then(() => {
    message.success(t('msg_sync_completed'))
  }).finally(() => {
    syncing.value = false
  })
}


function loadTags() {
  return runPlugin({
    name: pluginName,
    action: "default/exec",
    params: JSON.stringify({
      business: 'getTags',
      arguments: {
        app_id: formState.app_id,
        app_secret: formState.app_secret,
      }
    })
  }).then(res => {
    tags.value = res?.data?.tags || []
    return res
  })
}

function changeFieldValue(field, text, selectedList) {
  if (field === 'openid_list') {
    text = text.trim()
    text = text.replace(/，/g, ',')
  }
  formState[field] = text
  formState.tag_map[field] = selectedList
  update()
}

function tagChange(val, opt) {
  formState.tag_name = opt.name
  update()
}

function tagTypeChange() {
  formState.tag_name = ''
  if (formState.tag_type == 1) {
    formState.tagid = null
  } else {
    formState.tagid = ''
  }
}

function update() {
  let nodeParams = JSON.parse(props.node.node_params)
  nodeParams.plugin.output_obj = outputData.value
  Object.assign(nodeParams.plugin.params, {
    arguments: {
      ...formState
    }
  })
  setData({
    ...props.node,
    node_params: JSON.stringify(nodeParams)
  })
}

const filterOption = (input, option) => {
  return option.name.toLowerCase().indexOf(input.toLowerCase()) >= 0;
}
</script>

<style scoped lang="less">
@import "../node-options";
.tag-box {
  display: flex;
  align-items: center;
  :deep(.mention-input-warpper) {
    height: 33px;
  }
}
</style>
