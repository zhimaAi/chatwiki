<template>
  <PluginFormRender
    :node="node"
    :params="params"
    :output="action.output"
    :variableOptions="variableOptions"
    :keepFormatField="['tagid', 'tag_type']"
    @updateVar="emit('updateVar')"
  >
    <template #tag_name></template>
    <template #tag_type></template>

    <template #app_name></template>
    <template #app_secret></template>
    <template #app_id="{ state, item, keyName}">
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">{{ t('label_official_account') }}</div>
        </div>
        <div>
          <a-select
            v-model:value="item.value"
            :placeholder="t('ph_select_official_account')"
            style="width: 100%;"
            @change="(value, option) => appChange(state, option)"
          >
            <a-select-option
              v-for="app in apps"
              :key="app.app_id"
              :value="app.app_id"
              :name="app.app_name"
              :secret="app.app_secret"
            >
              {{ app.app_name }}
            </a-select-option>
          </a-select>
        </div>
      </div>
    </template>
    <template #tagid="{ state, item, keyName}">
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">{{ t('label_official_account_tag') }}</div>
          <a-tooltip :title="t('tooltip_sync_tags')">
            <a @click="syncTags">{{ t('btn_sync') }} <a-spin v-if="syncing" size="small"/></a>
          </a-tooltip>
        </div>
        <div class="tag-box">
          <a-select
            v-model:value="state.tag_type.value"
            @change="tagTypeChange(state)"
            style="width: 120px;"
          >
            <a-select-option :value="1">{{ t('btn_select_tag') }}</a-select-option>
            <a-select-option :value="2">{{ t('btn_insert_variable') }}</a-select-option>
          </a-select>
          <a-select
            v-if="state.tag_type.value == 1"
            v-model:value="item.value"
            :placeholder="t('ph_select_tag')"
            show-search
            :filter-option="filterOption"
            style="width: 100%;"
            @change="(val,opt) => tagChange(state,val,opt)"
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
            :defaultSelectedList="item.tags || []"
            :defaultValue="item.value"
            ref="atInputRef"
            @open="emit('updateVar')"
            @change="(text, selectedList) => changeValue(item, text, selectedList)"
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
    </template>
  </PluginFormRender>
</template>

<script setup>
import {ref, onMounted} from 'vue';
import { useI18n } from '@/hooks/web/useI18n';
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import {getWechatAppList} from "@/api/robot/index.js";
import PluginFormRender from "../pluginFormRender.vue";
import {runPlugin} from "@/api/plugins/index.js";
import {jsonDecode, sortObjectKeys} from "@/utils/index.js";
import {message} from 'ant-design-vue'

const { t } = useI18n('views.workflow.components.node-form-drawer.components.official-account.official-account-box');

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

const apps = ref([])
const tags = ref([])
const params = ref({})
const syncing = ref(false)
const config = ref({
  app_id: '',
  app_secret: '',
})

onMounted(() => {
  loadWxApps()
  getParams()
  if (props.actionName === 'getTagFans') {
    let node_params = jsonDecode(props.node?.node_params)
    if (node_params) {
      let args = node_params?.plugin?.params?.arguments || {}
      const {app_id, app_secret} = args
      if (app_id && app_secret) {
        config.value = {app_id, app_secret}
        loadTags()
      }
    }
  }
})

function getParams() {
  const baseParams = {
    ...props.action.params,
    app_id: {
      ...props.action.params.app_id,
      default: null,
    },
    app_name: {
      desc: t('label_official_account_name'),
      type: 'string',
    },
  }

  if (props.actionName === 'getTagFans') {
    Object.assign(baseParams, {
      tagid: {
        ...props.action.params.tagid,
        default: null,
      },
      tag_name: {
        desc: t('label_tag_name'),
        type: 'string',
      },
      tag_type: {
        desc: t('label_tag_type'),
        type: 'number',
        default: 1,
      }
    })
    params.value = sortObjectKeys(baseParams, ['app_id', 'tagid'])
  } else {
    params.value = baseParams
  }
}


function loadWxApps() {
  getWechatAppList({app_type: 'official_account'}).then(res => {
    apps.value = res?.data || []
  })
}

function syncTags() {
  if (!config.value.app_secret || !config.value.app_id) {
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
    name: 'official_account_profile',
    action: "default/exec",
    params: JSON.stringify({
      business: 'getTags',
      arguments:  config.value
    })
  }).then(res => {
    tags.value = res?.data?.tags || []
    return res
  })
}

function appChange(state, option) {
  state.app_secret.value = option.secret
  state.app_name.value = option.name
  if (props.actionName === 'getTagFans') {
    config.value = {
      app_id: option.key,
      app_secret:  option.secret
    }
    loadTags()
  }
}

function tagTypeChange(state) {
  state.tag_name.value = ''
  if (state.tag_type.value == 1) {
    state.tagid.value = null
  } else {
    state.tagid.value = ''
  }
}

function tagChange(state, val, opt) {
  state.tag_name.value = opt.name
  console.log('state', state)
}

function changeValue(item, text, selectedList) {
  item.value = text
  item.tags = selectedList
}

const filterOption = (input, option) => {
  return option.name.toLowerCase().indexOf(input.toLowerCase()) >= 0;
}
</script>

<style scoped lang="less">
.tag-box {
  display: flex;
  align-items: center;
  :deep(.mention-input-warpper) {
    height: 33px;
  }
}
</style>
