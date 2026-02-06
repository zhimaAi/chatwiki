<template>
  <div class="tpl-container">
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
          <div class="option-label">{{ t('label_official_template') }}</div>
          <a-tooltip :title="t('tooltip_sync_templates')">
            <a @click="syncTemplates">{{ t('btn_sync') }} <a-spin v-if="syncing" size="small"/></a>
          </a-tooltip>
        </div>
        <div>
          <a-select
            v-model:value="formState.template_id"
            :placeholder="t('ph_select_template')"
            style="width: 100%;"
            @change="templateChange"
          >
            <a-select-option
              v-for="item in templates"
              :key="item.template_id"
              :name="item.title"
              :content="item.content"
            >
              {{ item.title }}
            </a-select-option>
          </a-select>
        </div>
      </div>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">{{ t('label_receiver_openid') }}</div>
          <div class="option-type">string</div>
        </div>
        <div>
          <AtInput
            type="textarea"
            :options="variableOptions"
            :defaultSelectedList="formState.tag_map?.touser || []"
            :defaultValue="formState.touser"
            ref="atInputRef"
            @open="emit('updateVar')"
            @change="(text, selectedList) => changeFieldValue('touser', text, selectedList)"
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
        <div class="desc">{{ t('desc_receiver_openid') }}</div>
      </div>
      <div
        v-for="(opt, key) in formState.data"
        class="options-item is-required"
        :key="key"
      >
        <div class="options-item-tit">
          <div class="option-label">{{ opt.title }}</div>
        </div>
        <div>
          <AtInput
            type="textarea"
            inputStyle="height: 64px;"
            :options="variableOptions"
            :defaultSelectedList="opt.tags"
            :defaultValue="opt.value"
            ref="atInputRef"
            @open="emit('updateVar')"
            @change="(text, selectedList) => changeValue(opt, text, selectedList, key)"
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
        <div v-if="opt.desc" :class="['desc', {error: opt.format_err}]">{{opt.desc}}</div>
      </div>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">{{ t('label_jump_type') }}</div>
        </div>
        <div>
          <a-radio-group
            v-model:value="formState.link_type"
            @change="linkTypeChange"
          >
            <a-radio :value="0">{{ t('opt_jump_none') }}</a-radio>
            <a-radio :value="1">{{ t('opt_jump_link') }}</a-radio>
            <a-radio :value="2">{{ t('opt_jump_miniprogram') }}</a-radio>
          </a-radio-group>
        </div>
      </div>
      <div v-if="formState.link_type == 1" class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">{{ t('label_link') }}</div>
        </div>
        <div>
          <AtInput
            type="textarea"
            inputStyle="height: 64px;"
            :options="variableOptions"
            :defaultSelectedList="formState.tag_map?.url || []"
            :defaultValue="formState.url"
            ref="atInputRef"
            @open="emit('updateVar')"
            @change="(text, selectedList) => changeFieldValue('url', text, selectedList)"
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
      <template v-else-if="formState.link_type == 2">
        <div class="options-item is-required">
          <div class="options-item-tit">
            <div class="option-label">{{ t('label_miniprogram_appid') }}</div>
          </div>
          <div>
            <AtInput
              type="textarea"
              inputStyle="height: 64px;"
              :options="variableOptions"
              :defaultValue="formState.miniprogram.appid.value"
              :defaultSelectedList="formState.miniprogram.appid.tags"
              ref="atInputRef"
              @open="emit('updateVar')"
              @change="(text, selectedList) => changeValue(formState.miniprogram.appid, text, selectedList)"
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
          <div class="desc">{{ t('desc_miniprogram_appid') }}</div>
        </div>
        <div class="options-item is-required">
          <div class="options-item-tit">
            <div class="option-label">{{ t('label_miniprogram_path') }}</div>
          </div>
          <div>
            <AtInput
              type="textarea"
              inputStyle="height: 64px;"
              :options="variableOptions"
              :defaultValue="formState.miniprogram.pagepath.value"
              :defaultSelectedList="formState.miniprogram.pagepath.tags"
              ref="atInputRef"
              @open="emit('updateVar')"
              @change="(text, selectedList) => changeValue(formState.miniprogram.pagepath, text, selectedList)"
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
const pluginName = 'official_send_template_message'
const outputData = ref([])
const apps = ref([])
const syncing = ref(false)
const templates = ref([])
const formState = reactive({
  app_id: undefined,
  template_id: undefined,
  app_secret: '',
  touser: '',
  url: '',
  data: {},
  miniprogram: {
    appid: {
      value: '',
      tags: [],
    },
    pagepath: {
      value: '',
      tags: [],
    }
  },
  link_type: 0,
  tag_map: {}
})

// Validation rules object
const getRegObj = () => ({
  thing: {
    reg: /^[\u4e00-\u9fa5a-zA-Z0-9\W_]{1,20}$/,
    txt: t('msg_validation_thing')
  },
  short_thing: {
    reg: /^[\u4e00-\u9fa5a-zA-Z0-9\W_]{1,5}$/,
    txt: t('msg_validation_short_thing')
  },
  number: {
    reg: /^\d{1,32}(.\d{1,2})?$/,
    txt: t('msg_validation_number')
  },
  const: {
    reg: /^[\u4e00-\u9fa5a-zA-Z0-9\W_]{1,20}$/,
    txt: t('msg_validation_const')
  },
  letter: {
    reg: /^[a-zA-Z]{1,32}$/,
    txt: t('msg_validation_letter')
  },
  symbol: {
    reg: /^[\W_]{1,5}$/,
    txt: t('msg_validation_symbol')
  },
  character_string: {
    reg: /^[a-zA-Z0-9\W_]{1,32}$/,
    txt: t('msg_validation_character_string')
  },
  time: {
    reg: /^\d{4}年\d{1,2}月\d{1,2}日\s+\d{1,2}:\d{1,2}(:\d{1,2})?(~\d{4}年\d{1,2}月\d{1,2}日\s+\d{1,2}:\d{1,2}(:\d{1,2})?)?$/,
    txt: t('msg_validation_time')
  },
  date: {
    reg: /^\d{4}年\d{1,2}月\d{1,2}日(\s+\d{1,2}:\d{1,2}(:\d{1,2})?)?(~\d{4}年\d{1,2}月\d{1,2}日\s+\d{1,2}:\d{1,2}(:\d{1,2})?)?$/,
    txt: t('msg_validation_date')
  },
  amount: {
    reg: /^[$¥]?\d{1,10}(.\d{1,2})?[元]?$/,
    txt: t('msg_validation_amount')
  },
  phone_number: {
    reg: /^[0-9-+*# ]{5,17}$/,
    txt: t('msg_validation_phone_number')
  },
  car_number: {
    reg: /^[京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼使领A-Z]{1}[a-zA-Z0-9]{6,7}[\u4e00-\u9fa5]{0,1}$/,
    txt: t('msg_validation_car_number')
  },
  name: {
    reg: /^(?:[\u4e00-\u9fa5]{1,10}|[a-zA-Z\p{P}]{1,20})$/,
    txt: t('msg_validation_name')
  },
  phrase: {
    reg: /^[\u4e00-\u9fa5]{1,5}$/,
    txt: t('msg_validation_phrase')
  },
})

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
  if (formState.template_id) {
    loadTemplates()
  }
}

function appChange(_, option) {
  const {key, secret, name} = option
  formState.app_secret = secret
  formState.app_name = name
  loadTemplates()
}

function templateChange(_, option) {
  const {name, content} = option
  formState.template_name = name
  formState.data = parseWxTemplateToObj(content)
}

function syncTemplates() {
  if (!formState.app_secret || !formState.app_id) {
    return message.warning(t('msg_please_select_official_account'))
  }
  if (syncing.value) return
  syncing.value = true
  loadTemplates().then(() => {
    message.success(t('msg_sync_completed'))
  }).finally(() => {
    syncing.value = false
  })
}

function loadTemplates() {
  return runPlugin({
    name: pluginName,
    action: "default/exec",
    params: JSON.stringify({
      business: 'getAllTemplates',
      arguments: {
        app_id: formState.app_id,
        app_secret: formState.app_secret,
      }
    })
  }).then(res => {
    templates.value = res?.data?.template_list || []
    return
  })
}

function changeValue(item, text, selectedList, field=null) {
  item.value = text;
  item.tags = selectedList
  if (field) {
    // 模板内容进行正则校验
    if (selectedList.length) {
      item.format_err = false
    } else {
      wxTemplateContCheck(field, item)
    }
  }
  update()
}

function changeFieldValue(field, text, selectedList) {
  formState[field] = text
  formState.tag_map[field] = selectedList
  update()
}

function linkTypeChange() {
  switch (formState.link_type) {
    case 0:
      formState.miniprogram.appid = {value: '', tags: []}
      formState.miniprogram.pagepath = {value: '', tags: []}
      formState.url = ''
      break
    case 1:
      formState.miniprogram.appid = {value: '', tags: []}
      formState.miniprogram.pagepath = {value: '', tags: []}
      break
    case 2:
      formState.url = ''
      break
  }
  update()
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

function parseWxTemplateToObj(str) {
  const obj = {};
  const lines = str.split('\n').map(s => s.trim()).filter(Boolean);
  const regObj = getRegObj();

  for (const line of lines) {
    // 情况1：有标题，例如：标题:{{content.DATA}}
    let match = line.match(/^(.*?):\s*\{\{([\w]+)\.DATA\}\}$/);
    if (match) {
      let key = match[2].replace(/\d/g, '')
      obj[match[2]] = { title: match[1], value: '', tags: [], desc: regObj[key]?.txt || '', format_err: false};
      continue;
    }

    // 情况2：无标题，例如：{{content.DATA}}
    match = line.match(/^\{\{([\w]+)\.DATA\}\}$/);
    if (match) {
      obj[match[1]] = { title: match[1], value: '', tags: []};
      continue;
    }
  }
  return obj;
}

function wxTemplateContCheck(field, item) {
  let key = field.replace(/\d/g, '')
  const regObj = getRegObj()
  if(regObj[key] && !regObj[key].reg.test(item.value)) {
    item.format_err = true
  } else {
    item.format_err = false
  }
}
</script>

<style scoped lang="less">
@import "../node-options";

.tpl-container {
  :deep(.mention-input-warpper) {
    height: 48px;
    .type-textarea {
      height: 44px;
      min-height: 44px;
    }
  }
}

.desc.error {
  color: red !important;
}
</style>
