<template>
  <div class="_plugin-box">
    <LoadingBox v-if="loading"/>
    <div v-else-if="list.length" class="plugin-list">
      <div v-for="item in list"
           @click="linkDetail(item)"
           :key="item.name"
           class="plugin-item">
        <div class="type-tag">{{item.filter_type_title}}</div>
        <div class="base-info">
          <img class="avatar" :src="item.icon"/>
          <div class="info">
            <div class="head">
              <span class="name zm-line1">{{ item.title }}</span>
            </div>
            <div class="source">{{ item.author }}</div>
          </div>
        </div>
        <div class="desc zm-line1">{{ item.description }}</div>
        <div class="version">版本：v{{ item.latest_version }}</div>
        <div class="action-box">
          <div class="left">
            <svg-icon name="download"/>
            {{ item.installed_count || 0 }}
          </div>
          <div class="right">
            <template v-if="item.help_url">
              <a @click.stop class="c595959" :href="item.help_url" target="_blank">使用说明</a>
              <a-divider type="vertical"/>
            </template>
            <template v-if="item.local">
              <a v-if="item.has_update" @click.stop="install(item)">更新</a>
              <span v-else>已安装</span>
            </template>
            <a v-else @click.stop="install(item)">安装</a>
          </div>
        </div>
      </div>
    </div>
    <EmptyBox v-else title="暂无可用插件"/>

    <UpdateModal ref="updateRef" @ok="loadData"/>
  </div>
</template>

<script setup>
import {onMounted, ref, watch} from 'vue';
import {useRouter} from 'vue-router';
import EmptyBox from "@/components/common/empty-box.vue";
import {getInstallPlugins, getRemotePlugins} from "@/api/plugins/index.js";
import LoadingBox from "@/components/common/loading-box.vue";
import UpdateModal from "./update-modal.vue";

const emit = defineEmits(['installReport'])
const props = defineProps({
  filterData: {
    type: Object,
    default: null
  }
})
const router = useRouter()

const updateRef = ref(null)
const loading = ref(true)
const list = ref([])
const localMap = ref({})

onMounted(() => {
  loadData()
})

// watch(() => props.filterData, () => {
//   loadData()
// }, {
//   immediate: true,
//   deep: true
// })

function search() {
  list.value = []
  loadData()
}

async function loadData() {
  loading.value = true
  await loadInstalls()
  getRemotePlugins(props.filterData).then(res => {
    let _list = res?.data || []
    _list.forEach(item => {
      let {local, has_update} = localMap.value?.[item.name] || {}
      item.local = local
      item.has_update = has_update
    })
    list.value = _list
  }).finally(() => {
    loading.value = false
  })
}

async function loadInstalls() {
  await getInstallPlugins().then(res => {
    let _list = res?.data || []
    emit('installReport', _list.length)
    let map = {}
    _list.forEach(item => {
      map[item.local.name] = {
        local: item.local,
        has_update: item.remote.latest_version != item.local.version,
      }
    })
    localMap.value = map
  })
}

function linkDetail(item) {
  router.push({
    path: '/plugins/detail',
    query: {
      name: item.name
    }
  })
}

function install(item) {
  updateRef.value.show(item, item.latest_version_detail, item.local || null)
}

defineExpose({
  search,
})
</script>

<style scoped lang="less">
@import "plugins";

.text-center {
  text-align: center;
}

.mt24 {
  margin-top: 24px;
}

.c595959 {
  color: #595959;
}
</style>
