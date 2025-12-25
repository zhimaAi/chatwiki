<template>
  <ListTabs :tabs="tabs" v-model:value="active" @change="change" />
</template>

<script setup>
import {ref, onMounted} from 'vue';
import {useRouter} from 'vue-router'
import ListTabs from "@/components/cu-tabs/list-tabs.vue";
import {getInstallPlugins, triggerConfigList} from "@/api/plugins/index.js";

const emit = defineEmits(['change'])

const router = useRouter()
const active = ref(localStorage.getItem('zm:explore:active') || '1')
const tabs = ref([
  {
    title: '模板广场',
    value: '4'
  },
  {
    title: '功能',
    value: '1'
  },
  {
    title: `插件(${localStorage.getItem('zm:explore:plugins:count') || '0'})`,
    value: '2'
  },
  {
    title: '插件广场',
    value: '3'
  },
  {
    title: 'MCP广场',
    value: '5'
  }
])

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
    let plugin = tabs.value.find(i => i.value == 2)
    plugin.title = `插件(${_list.length + triggerLength})`
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
