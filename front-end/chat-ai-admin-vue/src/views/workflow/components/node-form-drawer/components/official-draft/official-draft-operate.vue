<template>
  <PluginFormRender
    :node="node"
    :params="params"
    :output="action.output"
    :variableOptions="variableOptions"
    @updateVar="emit('updateVar')"
  >
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
        <div v-if="['publish_draft', 'preview_message'].includes(actionName)" class="desc">
          <ExclamationCircleFilled/>
          {{ t('msg_unsupported_account') }}
        </div>
      </div>
    </template>
    <template v-if="actionName === 'preview_message'" #type="{state, item}">
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">{{ t('label_message_user_type') }}</div>
        </div>
        <div>
          <a-radio-group v-model:value="item.value">
            <a-radio value="touser">{{ t('opt_official_user_openid') }}</a-radio>
            <a-radio value="towxname">{{ t('opt_wechat_id') }}</a-radio>
          </a-radio-group>
        </div>
      </div>
    </template>
  </PluginFormRender>
</template>

<script setup>
import {ref, reactive, computed, onMounted} from 'vue';
import {ExclamationCircleFilled} from '@ant-design/icons-vue';
import PluginFormRender from "@/views/workflow/components/node-form-drawer/components/pluginFormRender.vue";
import {getWechatAppList} from "@/api/robot/index.js";
import {jsonDecode, sortObjectKeys} from "@/utils/index.js";
import {useI18n} from "@/hooks/web/useI18n";

const { t } = useI18n('views.workflow.components.node-form-drawer.components.official-draft.official-draft-operate')
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
  },
})
const params = ref({})
const apps = ref([])

onMounted(() => {
  loadWxApps()
  getParams()
})

function getParams() {
  const baseParams = {
    ...props.action.params,
    app_id: {
      ...props.action.params.app_id,
      default: null,
    },
    app_name: {
      desc: 'label_official_account_name',
      type: 'string',
    },
  }
  let app = jsonDecode(window.localStorage.getItem('zm:ai:draft:last:app'))
  if (app) {
    baseParams.app_id.default = app.app_id
    baseParams.app_name.default = app.app_name
    baseParams.app_secret.default = app.app_secret
  }
  if (props.actionName === "preview_message") {
    baseParams.type = {
      ...props.action.params.type,
      default: 'touser'
    }
    params.value = sortObjectKeys(baseParams, ['app_id', 'media_id', 'type', 'account'])
  } else {
    params.value = baseParams
  }
}

function loadWxApps() {
  getWechatAppList({app_type: 'official_account'}).then(res => {
    apps.value = res?.data || []
  })
}

function appChange(state, option) {
  const {key, secret, name} = option
  state.app_secret.value = secret
  state.app_name.value = name
  window.localStorage.setItem('zm:ai:draft:last:app', JSON.stringify({
    app_id: key,
    app_secret: secret,
    app_name: name
  }))
  update()
}

function update() {

}
</script>

<style scoped>

</style>
