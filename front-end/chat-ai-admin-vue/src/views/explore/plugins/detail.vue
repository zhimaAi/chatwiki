<template>
  <div class="_container _plugin-detail-cont">
    <div class="header">
      <LeftOutlined @click="router.back()" class="back-icon"/>
      <span class="title">插件详情</span>
    </div>
    <LoadingBox v-if="loading"/>
    <div v-else class="content">
      <div class="base-info">
        <div class="top">
          <img class="avatar" :src="info.icon"/>
          <div class="info">
            <div class="left">
              <div class="name zm-line1">{{ info.title }} v{{info.latest_version}}</div>
              <div class="source">{{ info.author }}</div>
            </div>
            <div class="right">
              <template v-if="info.local">
                <a-button v-if="info.has_update" type="primary" @click="install">更 新</a-button>
                <a-button v-else>已安装</a-button>
              </template>
              <a-button v-else type="primary" @click="install">安 装</a-button>
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
            <div class="main-tit">插件概叙</div>
            <pre class="text-cont">{{info.overview}}</pre>
          </div>
          <div class="main-block">
            <div class="main-tit">使用说明</div>
            <pre class="text-cont" v-html="info.introduction"></pre>
          </div>
        </div>
        <div class="right">
          <div class="main-block">
            <div class="main-tit">版本</div>
            <div class="version-info">
              v{{ info.latest_version }}
              <a-divider type="vertical"/>
              {{ info.author }} {{ dayjs.unix(info?.latest_version_detail?.update_time).format('YYYY-MM-DD HH:mm') }}
            </div>
            <div class="text-cont">{{info?.latest_version_detail?.upgrade_description}}</div>
            <div>
              <a @click="showVersionHistory">版本历史
                <RightOutlined/>
              </a>
            </div>
          </div>
          <div class="main-block">
            <div class="main-tit">权限列表</div>
            <div class="text-cont">最大内存 {{ info?.latest_version_detail?.memory }}MB</div>
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
        }
      }

      .version-info {
        color: #262626;
      }
    }
  }
}
</style>
