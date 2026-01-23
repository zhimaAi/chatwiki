<template>
  <div>
    <a-modal
      v-model:open="open"
      title="模板详情"
      wrapClassName="no-padding-modal"
      :footer="null"
      :width="800"
    >
      <div class="_container _plugin-detail-cont">
        <LoadingBox v-if="loading" />
        <div v-else class="content">
          <div class="base-info">
            <div class="top">
              <img class="avatar" :src="tplInfo.avatar" />
              <div class="info">
                <div class="left">
                  <div class="name zm-line1">{{ tplInfo.name }}</div>
                  <div class="source">
                    <img class="avatar" :src="tplInfo.account_avatar" />
                    {{ tplInfo.author }}
                  </div>
                </div>
                <div class="right">
                  <a-button class="use-btn" type="primary" ghost @click="useTemplate">
                    <svg-icon name="icon-rocket" />
                    使用模板
                  </a-button>
                </div>
              </div>
            </div>
            <div class="desc">{{ tplInfo.description }}</div>
          </div>
          <div class="detail-info">
            <div class="left">
              <div class="main-block">
                <div class="main-tit">使用说明</div>
                <pre class="text-cont" v-viewer v-html="tplInfo.instructions"></pre>
              </div>
            </div>
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, h, computed } from 'vue'
import { useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import LoadingBox from '@/components/common/loading-box.vue'
import { getTplDetailMain, useRobotTemplate } from '@/api/explore/template.js'
import { useCompanyStore } from '@/stores/modules/company.js'

const open = ref(false)

const show = (data) => {
  tplKey.value = data.robot_id
  loadInfo()
  open.value = true
}

const router = useRouter()

const companyStore = useCompanyStore()
const { companyInfo } = companyStore
const sysVersion = computed(() => {
  return companyInfo?.version
})

const loading = ref(true)
const tplKey = ref(null)
const tplInfo = ref({})

function loadInfo() {
  loading.value = true
  getTplDetailMain({ robot_id: tplKey.value })
    .then((res) => {
      tplInfo.value = res?.data || {}
    })
    .finally(() => {
      loading.value = false
    })
}

function useTemplate() {
  const item = tplInfo.value
  function checkVersion(sys_v, tpl_v) {
    sys_v = Number(sys_v.replace(/\D/g, ''))
    tpl_v = Number(tpl_v.replace(/\D/g, ''))
    return sys_v >= tpl_v
  }
  const run = () => {
    useRobotTemplate({ template_id: item.id, csl_url: item.csl_url }).then((res) => {
      message.success('使用成功')
      const { id, robot_key } = res.data
      const url = router.resolve({ path: '/robot/config/workflow', query: { id, robot_key } })
      window.open(url.href, '_blank')
      open.value = false
    })
  }
  if (!checkVersion(sysVersion.value, item.version)) {
    Modal.confirm({
      title: '提示',
      content: '当前系统版本过低，可能无法使用此模板；请您升级到最新版本后使用！',
      icon: h(ExclamationCircleOutlined),
      okText: '继续使用',
      cancelText: '取 消',
      onOk: run
    })
  } else {
    Modal.confirm({
      title: '提示',
      content: `确定使用模板【${item.name}】创建应用吗?`,
      icon: h(ExclamationCircleOutlined),
      okText: '确 认',
      cancelText: '取 消',
      onOk: run
    })
  }
}

defineExpose({
  show
})
</script>

<style scoped lang="less">
._plugin-detail-cont .text-cont {
  &::v-deep(img) {
    max-width: 640px;
    cursor: pointer;
  }
}
._container {
  display: flex;
  flex-direction: column;
  align-items: center;
  overflow: hidden;

  > div {
    width: 100%;
    padding-top: 8px;
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
    border: 1px solid #e4e6eb;
    margin-right: 24px;

    .top {
      display: flex;
      align-items: center;
      gap: 12px;

      .avatar {
        width: 62px;
        height: 62px;
        border-radius: 16px;
        border: 1px solid #f0f0f0;
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
        color: #8c8c8c;
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
    max-height: 500px;
    overflow-y: auto;
    padding-right: 12px;

    .left {
      flex-shrink: 0;
      width: 100%;
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
