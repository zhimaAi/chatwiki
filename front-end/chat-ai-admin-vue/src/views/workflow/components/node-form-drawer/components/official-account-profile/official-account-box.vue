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
  </PluginFormRender>
</template>

<script setup>
import {ref, reactive, onMounted} from 'vue';
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import {getWechatAppList} from "@/api/robot/index.js";
import PluginFormRender from "../pluginFormRender.vue";

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

const formState = reactive({})
const apps = ref([])
const params = ref({})

onMounted(() => {
  loadWxApps()
  getParams()
})

function getParams() {
  params.value = {
    ...props.action.params,
    app_id: {
      ...props.action.params.app_id,
      default: null,
    },
    app_name: {
      desc: '公众号名称',
      type: 'string'
    }
  }
}

function loadWxApps() {
  getWechatAppList({app_type: 'official_account'}).then(res => {
    apps.value = res?.data || []
  })
}

function appChange(state, option) {
  state.app_secret.value = option.secret
  state.app_name.value = option.name
}
</script>

<style scoped lang="less">

</style>
