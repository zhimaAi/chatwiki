<template>
  <div class="_container _plugin-detail-cont">
    <div class="header">
<!--      <LeftOutlined @click="router.back()" class="back-icon"/>-->
      <span class="title">{{ t('title_template_detail') }}</span>
    </div>
    <LoadingBox v-if="loading"/>
    <div v-else class="content">
      <div class="base-info">
        <div class="top">
          <img class="avatar" :src="tplInfo.avatar"/>
          <div class="info">
            <div class="left">
              <div class="name zm-line1">{{ tplInfo.name }}</div>
              <div class="source">
                <img class="avatar" :src="tplInfo.account_avatar"/>
                {{ tplInfo.author }}
              </div>
            </div>
            <div class="right">
              <a-button class="use-btn" type="primary" ghost @click="useTemplate">
                <svg-icon name="icon-rocket"/>
                {{ t('btn_use_template') }}
              </a-button>
            </div>
          </div>
        </div>
        <div class="desc">{{ tplInfo.description }}</div>
        <div class="action-box">
          <div class="left">
            <TeamOutlined/>
            {{ tplInfo.use_count || 0 }}
          </div>
          <div class="right"></div>
        </div>
      </div>
      <div class="detail-info">
        <div class="left">
<!--          <div v-if="info.overview" class="main-block">-->
<!--            <div class="main-tit">插件概叙</div>-->
<!--            <pre class="text-cont">{{ info.overview }}</pre>-->
<!--          </div>-->
          <div class="main-block">
            <div class="main-tit">{{ t('title_usage_instructions') }}</div>
            <pre class="text-cont" v-viewer v-html="tplInfo.instructions"></pre>
          </div>
        </div>
        <div class="right">
          <div class="main-block">
            <div class="main-tit">{{ t('title_template_category') }}</div>
            <div class="text-cont flex-center">
              {{ tplInfo.category_name }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {onMounted, ref} from 'vue';
import {useRoute} from 'vue-router';
import {TeamOutlined} from '@ant-design/icons-vue';
import LoadingBox from "@/components/common/loading-box.vue";
import {getTplDetailMain} from "@/api/explore/template.js";
import { useI18n } from '@/hooks/web/useI18n'
import {useTemplateCreateRobot} from "@/composables/useTemplateCreateRobot.js";

const { t } = useI18n('views.explore.templates.detail')
const {useTpl} = useTemplateCreateRobot()
const route = useRoute()

const loading = ref(true)
const tplKey = ref(null)
const tplInfo = ref({})

onMounted(() => {
  tplKey.value = route.query?.key
  loadInfo()
})

function loadInfo() {
  loading.value = true
  getTplDetailMain({robot_id: tplKey.value}).then(res => {
    tplInfo.value = res?.data || {}
  }).finally(() => {
    loading.value = false
  })
}

function useTemplate() {
  useTpl(tplInfo.value)
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
        display: flex;
        align-items: center;
        gap: 4px;
        .avatar {
          width: 16px;
          height: 16px;
          border-radius: 16px;
        }
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
      word-break: break-all;
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

.use-btn {
  display: flex;
  align-items: center;
  gap: 4px;
}
</style>
