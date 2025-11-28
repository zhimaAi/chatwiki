<template>
  <ListTabs :tabs="tabs" v-model:value="active" @change="change" />
</template>

<script setup>
import {ref, onMounted} from 'vue';
import {useRouter} from 'vue-router'
import ListTabs from "@/components/cu-tabs/list-tabs.vue";
import {getInstallPlugins} from "@/api/plugins/index.js";

const emit = defineEmits(['change'])

const router = useRouter()
const active = ref(localStorage.getItem('zm:explore:active') || '1')
const tabs = ref([
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

  }
}

function loadInstallPlugins() {
  getInstallPlugins().then(res => {
    let _list = res?.data || []
    tabs.value[1].title = `插件(${_list.length})`
    localStorage.setItem('zm:explore:plugins:count', _list.length)
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
