<template>
  <ListTabs :tabs="tabs" v-model:value="active" @change="change">
    <template #extra="{key}">
      <ExclamationCircleFilled v-if="!isPublicNetwork && key != 1" style="color: #FB363F;"/>
    </template>
  </ListTabs>
</template>

<script setup>
import {ref, onMounted, computed} from 'vue';
import {useRouter} from 'vue-router'
import {ExclamationCircleFilled} from '@ant-design/icons-vue'
import ListTabs from "@/components/cu-tabs/list-tabs.vue";
import {getInstallPlugins, triggerConfigList} from "@/api/plugins/index.js";
import {usePublicNetworkCheck} from "@/composables/usePublicNetworkCheck.js";
import { useI18n } from '@/hooks/web/useI18n';

const { t } = useI18n('views.explore.components.main-tab')
const emit = defineEmits(['change'])
const router = useRouter()
const active = ref(localStorage.getItem('zm:explore:active') || '1')
const pluginCount = ref(localStorage.getItem('zm:explore:plugins:count') || '0')
const tabs = computed(() => [
  {
    title: t('title_template_marketplace'),
    value: '4'
  },
  {
    title: t('title_features'),
    value: '1'
  },
  {
    title: `${t('title_plugins')}(${pluginCount.value})`,
    value: '2'
  },
  {
    title: t('title_plugin_marketplace'),
    value: '3'
  },
  {
    title: t('title_mcp_marketplace'),
    value: '5'
  }
])
const {isPublicNetwork} = usePublicNetworkCheck()

onMounted(() => {
  loadInstallPlugins()
})

function change(val){
  active.value = val.toString()
  localStorage.setItem('zm:explore:active', val)
  emit('change', val)
  switch (Number(val)) {
    case 1:
      router.push({path: '/explore/index'})
      break
    case 2:
    case 3:
      router.push({path: '/plugins/index', query: {active: val}})
      break
    case 4:
      router.push({path: '/templates/index'})
      break
    case 5:
      router.push({path: '/mcp/index'})
      break

  }
}

async function loadInstallPlugins() {
  let res = await triggerConfigList()
  let triggerLength = res.data.length
  getInstallPlugins().then(res => {
    let _list = res?.data || []
    pluginCount.value = _list.length + triggerLength
    localStorage.setItem('zm:explore:plugins:count', _list.length + triggerLength)
  })
}

function update() {
  loadInstallPlugins()
}

defineExpose({
  change,
  update,
})

</script>

<style scoped>

</style>
