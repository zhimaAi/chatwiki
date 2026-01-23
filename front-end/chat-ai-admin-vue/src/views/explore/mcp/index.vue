<template>
  <div class="_container">
    <div class="header">
      <div class="main-tab-box">
        <MainTab ref="tabRef" @change="tabChange"/>
      </div>
    </div>
    <PublicNetworkCheck v-if="!isPublicNetwork"/>
    <div v-else class="content" ref="contentRef">
      <div class="filter-box">
        <div class="radio-tabs">
          <div @click="filterTypeChange(0)" :class="['radio-tab', {active: filterData.filter_type == 0 }]">
            <!-- <svg-icon name="icon-all"/> -->
            <span>全部</span>
          </div>
          <div v-for="type in types"
               @click="filterTypeChange(type.id)"
               :key="type.id"
               :class="['radio-tab', {active: type.id == filterData.filter_type }]"
          >
            <img v-if="type.icon_path" class="icon" :src="type.icon_path">
            <span>{{ type.type_title }}</span>
          </div>
        </div>
        <a-input
          v-model:value.trim="filterData.title"
          @change="filterDataChange"
          style="width: 360px;"
          allowClear
          placeholder="搜索MCP">
          <template #suffix>
            <SearchOutlined/>
          </template>
        </a-input>
      </div>
      <Square @installReport="updateTab" :filterData="filterData" :scrollRoot="contentRef" ref="pluginRef"/>
    </div>
  </div>
</template>

<script setup>
import { onMounted , ref, reactive} from 'vue';
import {useRoute} from 'vue-router'
import {SearchOutlined} from '@ant-design/icons-vue';
import Square from "./components/square.vue";
import MainTab from "@/views/explore/components/main-tab.vue";
import {getMcpSquareTypeList} from "@/api/mcp/index.js";
import PublicNetworkCheck from "@/components/common/public-network-check.vue";
import {usePublicNetworkCheck} from "@/composables/usePublicNetworkCheck.js";

const {isPublicNetwork} = usePublicNetworkCheck()
const route = useRoute()
const tabRef = ref(null)
const pluginRef = ref(null)
const contentRef = ref(null)
const active = ref(localStorage.getItem('zm:explore:active') || '5')
const types = ref([])
const filterData = reactive({
  title: '',
  filter_type: 0
})

onMounted(() => {
  if (route.query.active > 1) {
    tabRef.value.change(route.query.active)
  }
  loadTypes()
})

function loadTypes() {
  getMcpSquareTypeList().then(res => {
    types.value = res?.data || []
  })
}

function filterDataChange() {
  pluginRef.value.search()
}

function filterTypeChange(key) {
  filterData.filter_type = key
  filterDataChange()
}

function tabChange(val){
  if (val > 1) active.value = val
}

function updateTab() {
  tabRef.value.update()
}
</script>

<style scoped lang="less">
._container {
  height: 100%;
  padding: 16px 24px 0 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.header {
  .main-tab-box {
    display: flex;
  }
  .tabs-box {
    display: flex;
    align-items: center;
    gap: 8px;

    .tab-item {
      padding: 5px 16px;
      border-radius: 6px;
      background: #EDEFF2;
      color: #595959;
      font-size: 14px;
      font-weight: 400;
      cursor: pointer;

      &.active {
        color: #2475fc;
        background: #D6E6FF;
      }
    }
  }
}

.content {
  flex: 1;
  overflow-y: auto;

  .filter-box {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 16px;
    margin: 10px 0 16px;

    .radio-tabs {
      display: flex;
      flex-wrap: wrap;
      align-items: center;
      gap: 8px;

      .radio-tab {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 14px;
        font-weight: 400;
        padding: 4px 16px;
        border-radius: 6px;
        cursor: pointer;
        &.active,
        &:hover {
          color: #2475FC;
          background: #D6E6FF;
        }

        .icon {
          width: 16px;
          height: 16px;
        }
      }
    }
  }
}
</style>
