<template>
  <div class="_mcp_container">
    <div class="tabs-box">
      <div :class="['tab-item', {active: active == 1}]" @click="tabChange(1)">{{ t('tab_workflow_mcp') }}</div>
      <div :class="['tab-item', {active: active == 2}]" @click="tabChange(2)">{{ t('tab_external_mcp') }}</div>
    </div>
    <McpBox v-if="active == 1" ref="mcpRef"/>
    <ThirdMcpBox v-else ref="mcpRef" @listLoaded="onThirdListLoaded"/>
  </div>
</template>

<script setup>
import {onMounted, ref} from 'vue';
import {useRoute} from 'vue-router'
import {useI18n} from '@/hooks/web/useI18n'
import McpBox from "@/views/robot/robot-list/components/mcp-box.vue";
import ThirdMcpBox from "@/views/robot/robot-list/components/third-mcp-box.vue";

const {t} = useI18n('views.robot.robot-list.components.mcp-panel')

const route = useRoute()

const mcpRef = ref(null)
const active = ref(Number(localStorage.getItem('zm:chatwiki:mcp:active') || 1))
const emit = defineEmits(['thirdMcpListEmpty'])

onMounted(() => {
  if (route.query.mcp > 0) {
    tabChange(Number(route.query.mcp))
  }
})

function tabChange(key) {
  active.value = key
  localStorage.setItem('zm:chatwiki:mcp:active', key)
}

function update() {
  if (active.value == 2) {
    mcpRef.value.update()
  }
}

function onThirdListLoaded(len) {
  emit('thirdMcpListEmpty', len === 0)
}

defineExpose({
  update
})
</script>

<style scoped lang="less">
._mcp_container {
  height: 100%;
  overflow-y: auto;

  .tabs-box {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 16px;

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
</style>
