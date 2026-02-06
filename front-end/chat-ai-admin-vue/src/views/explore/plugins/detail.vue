<template>
  <div class="_container _plugin-detail-cont">
    <div class="header">
      <LeftOutlined @click="router.back()" class="back-icon"/>
      <span class="title">{{ t('title_plugin_detail') }}</span>
    </div>
    <LoadingBox v-if="loading"/>
    <div v-else class="content">
      <div class="base-info">
        <div class="top">
          <img class="avatar" :src="info.icon"/>
          <div class="info">
            <div class="left">
              <div class="name zm-line1">{{ info.title }}</div>
              <div class="source">{{ info.author }}</div>
            </div>
            <div class="right">
              <a v-if="info.help_url" :href="info.help_url" target="_blank">
                <a-button type="primary" ghost>{{ t('btn_usage_guide') }}</a-button>
              </a>
              <template v-if="info.local">
                <a-button v-if="info.has_update" type="primary" @click="install">{{ t('btn_update') }}</a-button>
                <a-button v-else>{{ t('btn_installed') }}</a-button>
              </template>
              <a-button v-else type="primary" @click="install">{{ t('btn_install') }}</a-button>
            </div>
          </div>
        </div>
        <div class="desc">{{ info.description }}</div>
        <div class="action-box">
          <div class="left">
            <svg-icon name="download"/>
            {{ info.installed_count || 0 }}
          </div>
          <div class="right"></div>
        </div>
      </div>
      <div class="detail-info">
        <div class="left">
          <div v-if="info.overview" class="main-block">
            <div class="main-tit">{{ t('title_plugin_overview') }}</div>
            <pre class="text-cont">{{info.overview}}</pre>
          </div>
          <div class="main-block">
            <div class="main-tit">{{ t('btn_usage_guide') }}</div>
            <pre class="text-cont" v-viewer v-html="info.introduction"></pre>
          </div>
        </div>
        <div class="right">
          <div class="main-block">
            <div class="main-tit">{{ t('title_category') }}</div>
            <div class="text-cont flex-center">
              <img class="cate-icon" :src="info.filter_type_icon"/>
              {{info.filter_type_title}}
            </div>
          </div>
          <div class="main-block">
            <div class="main-tit">{{ t('title_version') }}</div>
            <div class="version-info">
              v{{ info.latest_version }}
              <a-divider type="vertical"/>
              {{ info.author }} {{ dayjs.unix(info?.latest_version_detail?.update_time).format('YYYY-MM-DD HH:mm') }}
            </div>
            <div class="text-cont">{{info?.latest_version_detail?.upgrade_description}}</div>
            <div>
              <a @click="showVersionHistory">{{ t('btn_version_history') }}
                <RightOutlined/>
              </a>
            </div>
          </div>
          <div class="main-block">
            <div class="main-tit">{{ t('title_permissions') }}</div>
            <div class="text-cont">{{ t('label_max_memory') }} {{ info?.latest_version_detail?.memory }}MB</div>
          </div>
        </div>
      </div>
    </div>

    <VersionModal ref="versionRef"/>
    <UpdateModal ref="updateRef" @ok="loadInfo"/>
  </div>
</template>

<script setup>
import {onMounted, ref} from 'vue';
import {useRouter, useRoute} from 'vue-router';
import dayjs from 'dayjs';
import {LeftOutlined, RightOutlined} from '@ant-design/icons-vue';
import VersionModal from "./components/version-modal.vue";
import UpdateModal from "./components/update-modal.vue";
import LoadingBox from "@/components/common/loading-box.vue";
import {getPluginInfo} from "@/api/plugins/index.js";
import { useI18n } from '@/hooks/web/useI18n';

const { t } = useI18n('views.explore.plugins.detail')
const router = useRouter()
const route = useRoute()

const updateRef = ref(true)
const loading = ref(true)
const pluginKey = ref(null)
const versionRef = ref(null)
const info = ref({})

onMounted(() => {
  pluginKey.value = route.query?.name
  loadInfo()
})

function loadInfo() {
  loading.value = true
  getPluginInfo({name: pluginKey.value}).then(res => {
    let data = res?.data || {}
    info.value = data.remote
    info.value.local = data.local
    info.value.has_update = (data.local?.version != data.remote?.latest_version)
    info.value.latest_version_history.map(item => {
      item.update_time_text = dayjs.unix(item.update_time).format('YYYY-MM-DD HH:mm')
      item.author = info.value.author
    })
  }).finally(() => {
    loading.value = false
  })
}

function showVersionHistory() {
  versionRef.value.show(info.value.latest_version_history || [])
}

function install() {
  let {latest_version_detail, local} = info.value
  updateRef.value.show(info.value, latest_version_detail, local)
}
</script>

<style lang="less">
._plugin-detail-cont .text-cont {
    img {
      max-width: 640px;
      cursor: pointer;
    }
}
</style>
<style scoped lang="less">
._container {
  display: flex;
  flex-direction: column;
  align-items: center;
  overflow: hidden;

  > div {
    width: 952px;
  }

  .header {
    padding: 16px 0;
    display: flex;
    align-items: center;
    gap: 10px;
    color: #000000;
    font-size: 20px;
    font-weight: 600;

    .back-icon {
      cursor: pointer;
    }
  }

  .base-info {
    display: flex;
    flex-direction: column;
    padding: 24px;
    gap: 12px;
    border-radius: 12px;
    border: 1px solid #E4E6EB;

    .top {
      display: flex;
      align-items: center;
      gap: 12px;

      .avatar {
        width: 62px;
        height: 62px;
        border-radius: 16px;
        border: 1px solid #F0F0F0;
      }

      .info {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: space-between;
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

    .right {
      display: flex;
      align-items: center;
      gap: 8px;
    }

    .desc {
      color: #595959;
      font-size: 14px;
      font-weight: 400;
    }

    .action-box {
      width: 100%;
      display: flex;
      align-items: center;
      justify-content: space-between;
      color: #595959;
      font-size: 14px;
      font-weight: 400;
    }
  }

  .detail-info {
    display: flex;
    justify-content: space-between;
    margin-top: 24px;

    .left {
      padding-right: 32px;
      border-right: 1px solid #F0F0F0;
      flex-shrink: 0;
      width: calc(100% - 280px);
    }

    .right {
      padding-left: 24px;
      width: 280px;
      flex-shrink: 0;
    }

    .main-block {
      display: flex;
      flex-direction: column;
      gap: 8px;
      margin-bottom: 32px;

      .main-tit {
        color: #000000;
        font-size: 16px;
        font-weight: 600;
        margin-bottom: 8px;
      }

      .text-cont {
        color: #595959;
        font-size: 14px;
        font-weight: 400;
        line-height: 22px;
        white-space: pre-wrap;
        word-wrap: break-word;
        overflow: hidden;

        img {
          max-width: 620px !important;
          cursor: pointer;
        }
      }

      .version-info {
        color: #262626;
      }
    }
  }
}

.cate-icon {
  width: 14px;
  height: 14px;
  margin-right: 4px;
}

.flex-center {
  display: flex;
  align-items: center;
}
</style>
