<template>
  <a-modal
    v-model:open="open"
    :title="statusInput.title"
    :ok-text="statusInput.text"
    :confirm-loading="installing"
    @ok="install"
  >
    <div class="plugin-box">
      <div class="base-info">
        <img class="avatar" :src="plug.icon"/>
        <div class="info">
          <div class="head">
            <span class="name zm-line1">
              {{ plug.title }}
              <span v-if="local" class="version-tag">v{{ local.version }} → v{{ plug.latest_version }}</span>
              <span v-else class="version-tag">v{{ plug.latest_version }}</span>
            </span>
          </div>
          <div class="source">{{ plug.author }}</div>
        </div>
      </div>
      <div class="desc zm-line2">{{ plug.description }}</div>
    </div>
  </a-modal>
</template>

<script setup>
import {ref, computed} from 'vue';
import {message} from 'ant-design-vue';
import {downloadPlugin} from "@/api/plugins/index.js";
import { useI18n } from '@/hooks/web/useI18n';

const { t } = useI18n('views.explore.plugins.components.update-modal');

const emit = defineEmits(['ok'])
const open = ref(false)
const installing = ref(false)
const plug = ref({})
const version = ref({})
const local = ref(null)

const statusInput = computed(() => {
  if (!local.value) {
    return {
      title: t('title_install_plugin'),
      text: t('btn_install'),
    }
  } else {
    return {
      title: t('title_update_plugin'),
      text: t('btn_update'),
    }
  }
})

function show(p, v, l = null) {
  plug.value = p
  version.value = v
  local.value = l
  open.value = true
}

function install() {
  installing.value = true
  downloadPlugin({
    download_data: JSON.stringify([{
      url: version.value.download_url,
      version_id: version.value.id
    }])
  }).then(() => {
    emit('ok')
    open.value = false
    message.success(t('msg_operation_completed', { operation: statusInput.value.text }))
  }).finally(() => {
    installing.value = false
  })
}

defineExpose({
  show,
})
</script>

<style scoped lang="less">
.plugin-box {
  padding: 24px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 12px;
  border-radius: 12px;
  border: 1px solid #E4E6EB;

  .base-info {
    display: flex;
    align-items: center;
    gap: 12px;

    .avatar {
      width: 62px;
      height: 62px;
      border-radius: 16px;
      border: 1px solid #F0F0F0;
    }

    .head {
      display: flex;
      align-items: center;
      gap: 4px;
    }

    .name {
      color: #262626;
      font-size: 16px;
      font-weight: 600;
    }

    .source {
      color: #8C8C8C;
      font-size: 12px;
      font-weight: 400;
    }
  }

  .desc {
    color: #595959;
    font-size: 14px;
    font-weight: 400;
  }
}

.version-tag {
  color: #ED744A;
  font-size: 16px;
  font-weight: 600;
}
</style>
