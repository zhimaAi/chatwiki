<template>
  <div class="tpl-container">
    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/input.svg" class="title-icon"/>输入</div>
      </div>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">公众号</div>
        </div>
        <div>
          <a-select
            v-model:value="formState.app_id"
            placeholder="请选择公众号"
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
          <div class="option-label">公众号模板</div>
          <a-tooltip title="同步最新的公众号模板">
            <a @click="syncTemplates">同步 <a-spin v-if="syncing" size="small"/></a>
          </a-tooltip>
        </div>
        <div>
          <a-select
            v-model:value="formState.template_id"
            placeholder="请选择模板"
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
          <div class="option-label">接收者（用户）openid</div>
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
            placeholder="请输入内容，键入“/”可以插入变量"
          >
            <template #option="{ label, payload }">
              <div class="field-list-item">
                <div class="field-label">{{ label }}</div>
                <div class="field-type">{{ payload.typ }}</div>
              </div>
            </template>
          </AtInput>
        </div>
        <div class="desc">普通用户的标识，对当前公众号唯一</div>
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
            placeholder="请输入内容，键入“/”可以插入变量"
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
          <div class="option-label">跳转方式</div>
        </div>
        <div>
          <a-radio-group
            v-model:value="formState.link_type"
            @change="linkTypeChange"
          >
            <a-radio :value="0">不跳转</a-radio>
            <a-radio :value="1">链接</a-radio>
            <a-radio :value="2">小程序</a-radio>
          </a-radio-group>
        </div>
      </div>
      <div v-if="formState.link_type == 1" class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">链接</div>
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
            placeholder="请输入内容，键入“/”可以插入变量"
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
            <div class="option-label">小程序APPID</div>
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
              placeholder="请输入内容，键入“/”可以插入变量"
            >
              <template #option="{ label, payload }">
                <div class="field-list-item">
                  <div class="field-label">{{ label }}</div>
                  <div class="field-type">{{ payload.typ }}</div>
                </div>
              </template>
            </AtInput>
          </div>
          <div class="desc">小程序appid必须与发模板消息的公众号是关联关系</div>
        </div>
        <div class="options-item is-required">
          <div class="options-item-tit">
            <div class="option-label">小程序路径</div>
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
              placeholder="请输入内容，键入“/”可以插入变量"
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
        <div><img src="@/assets/img/workflow/output.svg" class="title-icon"/>输出</div>
      </div>
      <div class="options-item">
        <OutputFields :tree-data="outputData"/>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref, reactive, onMounted, watch, inject} from 'vue';
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import {getWechatAppList} from "@/api/robot/index.js";
import {runPlugin} from "@/api/plugins/index.js";
import {pluginOutputToTree} from "@/constants/plugin.js";
import OutputFields from "@/views/workflow/components/feishu-table/output-fields.vue";
import {message} from 'ant-design-vue';


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
const RegObj = {
  thing: {
    reg: /^[\u4e00-\u9fa5a-zA-Z0-9\W_]{1,20}$/,
    txt: "20个以内字符,可汉字、数字、字母或符号组合"
  },
  short_thing: {
    reg: /^[\u4e00-\u9fa5a-zA-Z0-9\W_]{1,5}$/,
    txt: "5个以内字符,可汉字、数字、字母或符号组合"
  },
  number: {
    reg: /^\d{1,32}(.\d{1,2})?$/,
    txt: "32位以内数字,只能数字，可带小数"
  },
  const: {
    reg: /^[\u4e00-\u9fa5a-zA-Z0-9\W_]{1,20}$/,
    txt: "20位以内字符，超过无法下发注"
  },
  letter: {
    reg: /^[a-zA-Z]{1,32}$/,
    txt: "32位以内字母，只能字母"
  },
  symbol: {
    reg: /^[\W_]{1,5}$/,
    txt: "5位以内符号，只能符号"
  },
  character_string: {
    reg: /^[a-zA-Z0-9\W_]{1,32}$/,
    txt: "32位以内数字、字母或符号，可数字、字母或符号组合"
  },
  time: {
    reg: /^\d{4}年\d{1,2}月\d{1,2}日\s+\d{1,2}:\d{1,2}(:\d{1,2})?(~\d{4}年\d{1,2}月\d{1,2}日\s+\d{1,2}:\d{1,2}(:\d{1,2})?)?$/,
    txt: "24小时制时间格式（支持+年月日），支持填时间段，两个时间点之间用“~”符号连接 例如：2019年10月1日 15:01~2019年10月1日 18:01"
  },
  date: {
    reg: /^\d{4}年\d{1,2}月\d{1,2}日(\s+\d{1,2}:\d{1,2}(:\d{1,2})?)?(~\d{4}年\d{1,2}月\d{1,2}日\s+\d{1,2}:\d{1,2}(:\d{1,2})?)?$/,
    txt: "年月日格式（支持+24小时制时间），支持填时间段，两个时间点之间用“~”符号连接 例如：2019年10月1日 15:01~2019年10月1日 18:01"
  },
  amount: {
    reg: /^[$¥]?\d{1,10}(.\d{1,2})?[元]?$/,
    txt: "金额，1个币种符号+10位以内纯数字，可带小数，结尾可带“元”	可带小数 例如：¥200元、$200"
  },
  phone_number: {
    reg: /^[0-9-+*# ]{5,17}$/,
    txt: "电话，17位以内，数字、符号"
  },
  car_number: {
    reg: /^[京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼使领A-Z]{1}[a-zA-Z0-9]{6,7}[\u4e00-\u9fa5]{0,1}$/,
    txt: "车牌，8位以内，第一位与最后一位可为汉字，其余为字母或数字"
  },
  name: {
    reg: /^(?:[\u4e00-\u9fa5]{1,10}|[a-zA-Z\p{P}]{1,20})$/,
    txt: "10个以内纯汉字或20个以内纯字母或符号"
  },
  phrase: {
    reg: /^[\u4e00-\u9fa5]{1,5}$/,
    txt: "汉字，5个以内汉字"
  },
}

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
    return message.warning('请先选择公众号')
  }
  if (syncing.value) return
  syncing.value = true
  loadTemplates().then(() => {
    message.success('同步完成')
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

  for (const line of lines) {
    // 情况1：有标题，例如：标题:{{content.DATA}}
    let match = line.match(/^(.*?):\s*\{\{([\w]+)\.DATA\}\}$/);
    if (match) {
      let key = match[2].replace(/\d/g, '')
      obj[match[2]] = { title: match[1], value: '', tags: [], desc: RegObj[key]?.txt || '', format_err: false};
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
  if(RegObj[key] && !RegObj[key].reg.test(item.value)) {
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
