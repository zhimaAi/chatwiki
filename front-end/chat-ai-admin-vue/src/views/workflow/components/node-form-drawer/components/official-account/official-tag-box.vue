<template>
  <div class="tag-container">
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
          <div class="option-label">公众号标签</div>
        </div>
        <div class="tag-box">
          <a-select
            v-model:value="formState.tag_type"
            @change="tagTypeChange"
            style="width: 120px;"
          >
            <a-select-option :value="1">选择标签</a-select-option>
            <a-select-option :value="2">插入变量</a-select-option>
          </a-select>
          <a-select
            v-if="formState.tag_type == 1"
            v-model:value="formState.tagid"
            placeholder="请选择标签"
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
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">粉丝openid列表</div>
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
        <div class="desc">粉丝openid列表，最多50个，逗号分割</div>
      </div>
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

onMounted(() => {
  init()
})

watch(() => props.action, () => {
  outputData.value = pluginOutputToTree(JSON.parse(JSON.stringify(props.action.output || '{}')))
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

function loadTags() {
  runPlugin({
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
