<template>
  <section class="skill-page">
    <CustomerPageHeader :title="t('title')" />
    <div class="skill-content">
      <section class="agent-intro">
        <div class="intro-copy">
          <h2>{{ t('intro_title') }}</h2>
          <p>{{ t('intro_description') }}</p>
          <div class="intro-actions">
            <a-button type="primary" class="switch-agent" @click="emit('switch-agent')">{{ t('switch_agent') }}</a-button>
            <span class="agent-notice"><img :src="noticeIcon" alt="" />{{ t('agent_only') }}</span>
          </div>
        </div>
        <img class="agent-illustration" :src="agentIllustration" alt="" />
      </section>

      <section class="agent-benefits">
        <h2>{{ t('benefits_title') }}</h2>
        <p class="benefits-description">{{ t('benefits_description') }}</p>
        <div class="benefit-cards">
          <article v-for="benefit in benefits" :key="benefit.key" class="benefit-card">
            <span class="benefit-icon"><img :src="benefit.icon" alt="" /></span>
            <h3>{{ t(benefit.key + '_title') }}</h3>
            <p>{{ t(benefit.key + '_description') }}</p>
          </article>
        </div>
      </section>

      <footer v-if="canCreateAgent" class="create-agent">
        <div><h2>{{ t('create_title') }}</h2><p>{{ t('create_description') }}</p></div>
        <a-button type="primary" class="create-agent-button" @click="openCreator"><PlusOutlined />{{ t('create_agent') }}</a-button>
      </footer>
    </div>
    <AddAgentAlert v-if="canCreateAgent" ref="creator" @ok="emit('robot-created', $event)" />
  </section>
</template>

<script setup>
import { computed, ref } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { checkRole } from '@/utils/permission'
import CustomerPageHeader from '../components/page-header.vue'
import AddAgentAlert from '@/views/clawbot/components/add-agent-alert.vue'
import agentIllustration from '@/assets/img/customer-service/skill/agent-reasoning.png'
import noticeIcon from '@/assets/img/customer-service/skill/notice.svg'
import planningIcon from '@/assets/img/customer-service/skill/planning.svg'
import toolsIcon from '@/assets/img/customer-service/skill/tools.svg'
import memoryIcon from '@/assets/img/customer-service/skill/memory.svg'

const { t } = useI18n('views.customer-service.skill.index')
const emit = defineEmits(['switch-agent', 'robot-created'])
const creator = ref(null)
const canCreateAgent = computed(() => checkRole(['ClawBotCreate']))
const benefits = [
  { key: 'planning', icon: planningIcon },
  { key: 'tools', icon: toolsIcon },
  { key: 'memory', icon: memoryIcon }
]
function openCreator() {
  if (canCreateAgent.value) creator.value?.open()
}
</script>

<style lang="less" scoped>
.skill-page {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-height: 0;
  margin: 8px 8px 8px 0;
  overflow: hidden;
  border-radius: 12px;
  background: #fff;
}
.skill-content {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 24px;
  min-height: 0;
  padding: 24px;
  overflow-y: auto;
}
h2, h3, p { margin: 0; }
h2 { color: #080808; font-size: 20px; font-weight: 600; line-height: 28px; }
p { font-size: 14px; line-height: 22px; }
.agent-intro {
  position: relative;
  display: flex;
  flex-shrink: 0;
  align-items: center;
  min-height: 240px;
  overflow: hidden;
  border-radius: 12px;
  background: #f1f6fe;
}
.intro-copy {
  position: relative;
  z-index: 1;
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 12px;
  padding: 24px 32px;
  p { max-width: 480px; color: #363636; }
}
.intro-actions { display: flex; flex-wrap: wrap; align-items: center; gap: 12px; padding-top: 4px; }
.switch-agent { height: 40px; padding: 0 16px; border-radius: 6px; }
.agent-notice { display: inline-flex; align-items: center; gap: 4px; padding: 4px 8px; border-radius: 4px; color: #994000; background: #fff; font-size: 14px; line-height: 24px; }
.agent-illustration { display: block; align-self: stretch; flex-shrink: 0; width: 460px; height: 240px; object-fit: contain; }
.agent-benefits { flex-shrink: 0; }
.benefits-description { margin-top: 8px; color: #5a5a5a; }
.benefit-cards { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 16px; margin-top: 20px; }
.benefit-card {
  min-height: 208px;
  padding: 24px;
  border: 1px solid #d8d8d8;
  border-radius: 12px;
  h3 { margin-top: 20px; color: #080808; font-size: 16px; font-weight: 500; line-height: 24px; }
  p { margin-top: 8px; color: #5a5a5a; }
}
.benefit-icon { display: inline-flex; align-items: center; justify-content: center; width: 40px; height: 40px; border-radius: 6px; background: #f1f6fe; }
.create-agent {
  display: flex;
  flex-shrink: 0;
  flex-wrap: wrap;
  align-items: center;
  justify-content: center;
  gap: 24px;
  margin-top: auto;
  padding: 20px 24px;
  border: 1px solid #e8e8e8;
  border-radius: 12px;
  background: #f1f6fe;
  h2 { font-size: 18px; line-height: 26px; }
  p { margin-top: 4px; color: #5a5a5a; }
}
.create-agent-button { display: inline-flex; align-items: center; gap: 8px; height: 40px; padding: 0 20px; border-radius: 6px; }
@media (max-width: 1300px) {
  .intro-copy { padding: 24px; }
  .agent-illustration { width: 36%; }
}
@media (max-width: 1050px) {
  .agent-intro { flex-direction: column; align-items: stretch; }
  .agent-illustration { align-self: center; width: 100%; max-width: 460px; }
  .benefit-cards { grid-template-columns: 1fr; }
  .benefit-card { min-height: 0; }
}
</style>
