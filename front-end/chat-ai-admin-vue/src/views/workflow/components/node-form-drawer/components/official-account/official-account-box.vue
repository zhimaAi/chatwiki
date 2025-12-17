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
          <div class="option-label">公众号</div>
        </div>
        <div>
          <a-select
            v-model:value="item.value"
            placeholder="请选择公众号"
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
          <div class="option-label">公众号标签</div>
        </div>
        <div class="tag-box">
          <a-select
            v-model:value="state.tag_type.value"
            @change="tagTypeChange(state)"
            style="width: 120px;"
          >
            <a-select-option :value="1">选择标签</a-select-option>
            <a-select-option :value="2">插入变量</a-select-option>
          </a-select>
          <a-select
            v-if="state.tag_type.value == 1"
            v-model:value="item.value"
            placeholder="请选择标签"
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
  </PluginFormRender>
</template>

<script setup>
import {ref, onMounted} from 'vue';
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import {getWechatAppList} from "@/api/robot/index.js";
import PluginFormRender from "../pluginFormRender.vue";
import {runPlugin} from "@/api/plugins/index.js";
import {jsonDecode, sortObjectKeys} from "@/utils/index.js";

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

onMounted(() => {
  loadWxApps()
  getParams()
  if (props.actionName === 'getTagFans') {
    let node_params = jsonDecode(props.node?.node_params)
    if (node_params) {
      let args = node_params?.plugin?.params?.arguments || {}
      const {app_id, app_secret} = args
      if (app_id && app_secret) {
        loadTags(app_id, app_secret)
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
      desc: '公众号名称',
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
        desc: '标签名称',
        type: 'string',
      },
      tag_type: {
        desc: '标签类型',
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

function loadTags(app_id, app_secret) {
  runPlugin({
    name: 'official_account_profile',
    action: "default/exec",
    params: JSON.stringify({
      business: 'getTags',
      arguments: {
        app_id: app_id,
        app_secret: app_secret,
      }
    })
  }).then(res => {
    tags.value = res?.data?.tags || []
  })
}

function appChange(state, option) {
  state.app_secret.value = option.secret
  state.app_name.value = option.name
  if (props.actionName === 'getTagFans') {
    loadTags(option.key, option.secret)
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
