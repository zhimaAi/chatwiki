<template>
  <div class="guide-wrapper">
    <div class="guide-content">
      <div class="guide-header">
        <img src="@/assets/img/guide/header.png" alt="" />
      </div>
      <div class="guide-body">
        <div class="guide-left">
          <div class="process-box">
            <a-progress :percent="total_process" />
          </div>
          <div class="menu-list-box">
            <div
              class="menu-item"
              :class="{ active: item.type == current_process }"
              @click="handleChangeType(item)"
              v-for="(item, index) in guideMenus"
              :key="item.type"
            >
              <div class="icon-box">
                <CheckCircleFilled
                  v-if="item.is_finish == 1"
                  style="color: #21a665; font-size: 26px"
                />
                <svg-icon v-else :name="item.type" style="font-size: 40px"></svg-icon>
              </div>
              <div class="desc-box">
                <div class="title-text">
                  <span>0{{ index + 1 }}</span
                  >{{ item.name }}
                </div>
                <div class="desc-text">{{ item.desc }}</div>
              </div>
            </div>
          </div>
        </div>
        <div class="guide-right">
          <div class="step-list" v-if="current_item.steps">
            <div class="step-item" v-for="(item, index) in current_item.steps" :key="item.type">
              <div class="step-icon">
                <svg-icon
                  style="color: #2475fc"
                  v-if="item.is_finish == 1"
                  name="steps-success-icon"
                ></svg-icon>
                <span v-else class="num-text">{{ index + 1 }}</span>
              </div>
              <div class="step-text">{{ item.name }}</div>
              <div class="right-icon" v-if="index != current_item.steps.length - 1">
                <RightOutlined />
              </div>
            </div>
          </div>
          <div class="to-config-btn" v-if="current_item.show_config_btn">
            <a-button @click="handleToConfig()" type="primary" ghost>立即前往配置</a-button>
          </div>
          <div class="tag-list" v-if="current_item.tags">
            <div class="tag-item" v-for="(item, index) in current_item.tags" :key="index">
              <div class="tag-text-box">{{ item }}</div>
              <div class="right-icon" v-if="index != current_item.tags.length - 1">
                <RightOutlined />
              </div>
            </div>
          </div>
          <div class="tag-desc-box" v-if="current_item.step_desc">{{ current_item.step_desc }}</div>
          <div
            class="btn-block-wrapper"
            v-if="current_process == 'test_robot'"
            style="margin-top: 32px"
          >
            <a-button @click="handleTest" type="primary" style="background-color: #21a665">
              <CaretRightFilled />开始测试
            </a-button>
            <a-button @click="handlePre">上一步</a-button>
          </div>
          <div v-else class="btn-block-wrapper">
            <a-button @click="handlePre" v-if="current_item.index > 0">上一步</a-button>
            <a-button @click="handleNext" type="primary">下一步</a-button>
          </div>
          <div class="guide-preview-box" v-if="current_item.preview_imgs">
            <a-image
              width="100%"
              v-for="(img, index) in current_item.preview_imgs"
              :index="index"
              :src="img"
            />
            <div class="img-item"></div>
          </div>

          <div class="chat-test-box" v-if="current_process == 'test_robot'">
            <a-form :model="formState" ref="formRef" autocomplete="off" layout="vertical">
              <a-form-item
                label="机器人"
                name="robot_key"
                :rules="[{ required: true, message: '请选择机器人' }]"
              >
                <a-select
                  v-model:value="formState.robot_key"
                  placeholder="请选择机器人"
                  style="width: 364px"
                >
                  <a-select-option v-for="item in robotLists" :value="item.robot_key">{{
                    item.robot_name
                  }}</a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item
                label="测试问题"
                name="robot_key"
                :rules="[{ required: true, message: '请输入问题' }]"
              >
                <a-input
                  style="width: 364px"
                  v-model:value="formState.question"
                  placeholder="请输入问题"
                />
              </a-form-item>
            </a-form>

            <div class="test-result" v-if="hasTest">
              <div class="test-result-title">
                测试结果
                <div class="result-tag success" v-if="resultContent"><CheckCircleFilled />通过</div>
                <div class="result-tag error" v-else><CloseCircleFilled />错误</div>
              </div>
              <div class="result-content">{{ resultContent }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, reactive, onMounted, onUnmounted } from 'vue'
import { getRobotList } from '@/api/robot/index.js'
import { requestNotStream } from '@/api/guide/index'
import {
  RightOutlined,
  CaretRightFilled,
  CheckCircleFilled,
  CloseCircleFilled
} from '@ant-design/icons-vue'
import { useGuideStore } from '@/stores/modules/guide'
const guideStore = useGuideStore()

const robotLists = ref([])

const total_process = computed(() => {
  return +guideStore.total_process
})

const current_process = ref('set_model')

const menuMap = {
  set_model: {
    desc: '配置基础大模型',
    tags: ['设置', '模型管理', '可添加模型', '添加指定类型模型', '完成设置'],
    step_desc: '配置好的大模型可在机器人/知识库中引用',
    show_config_btn: true,
    preview_imgs: ['set_model_1'].map((name) => getImgUrl(name)),
    config_url: '/#/user/model?activeTab=0'
  },
  create_library: {
    desc: '创建知识库并导入知识',
    tags: ['知识库模块', '新增知识库', '上传指定类型文档', '完成设置'],
    step_desc: '配置好的知识库可在机器人中引用',
    show_config_btn: true,
    preview_imgs: ['create_library_1', 'create_library_2'].map((name) => getImgUrl(name)),
    config_url: '/#/library/list'
  },
  create_robot: {
    desc: '创建AI机器人并关联知识库',
    tags: ['机器人模块', '新增机器人', '关联知识库文档', '根据业务编辑关键词', '完成设置'],
    step_desc: '设置好的机器人即可测试发布',
    show_config_btn: true,
    preview_imgs: ['create_robot_1', 'create_robot_2', 'create_robot_3'].map((name) =>
      getImgUrl(name)
    ),
    config_url: '/#/robot/list'
  },
  test_robot: {
    desc: '测试机器人并分享给用户使用'
  }
}

const guideMenus = computed(() => {
  let list = guideStore.process_list || []
  list.forEach((item, index) => {
    item.desc = menuMap[item.type]?.desc
    item.tags = menuMap[item.type]?.tags
    item.step_desc = menuMap[item.type]?.step_desc
    item.show_config_btn = menuMap[item.type]?.show_config_btn
    item.preview_imgs = menuMap[item.type]?.preview_imgs
    item.config_url = menuMap[item.type]?.config_url
    item.index = index
  })
  return list
})

const current_item = computed(() => {
  return guideMenus.value.find((item) => item.type == current_process.value) || {}
})

const handleChangeType = (item) => {
  current_process.value = item.type
}

const handleToConfig = () => {
  window.open(current_item.value.config_url)
}

const handlePre = () => {
  let index = current_item.value.index
  current_process.value = guideMenus.value[index - 1].type
}

const handleNext = () => {
  let index = current_item.value.index
  current_process.value = guideMenus.value[index + 1].type
}

const formState = reactive({
  robot_key: void 0,
  openid: '1',
  question: '你好',
  form_ids: '',
  dialogue_id: 0,
  global: ''
})

const resultContent = ref('')
const hasTest = ref(false)

const formRef = ref()
const handleTest = () => {
  formRef.value.validate().then(() => {
    let parmas = {
      ...formState
    }
    let id = robotLists.value.find((item) => item.robot_key == formState.robot_key).id
    parmas.global = JSON.stringify({
      robot_key: formState.robot_key,
      id
    })

    requestNotStream(parmas)
      .then((res) => {
        resultContent.value = res.data.content
      })
      .catch(() => {
        resultContent.value = ''
      })
      .finally(() => {
        hasTest.value = true
        guideStore.getUseGuideProcess()
      })
  })
}

const getAllRobot = () => {
  getRobotList().then((res) => {
    robotLists.value = (res.data || []).filter((item) => item.application_type == 0)
    if (!formState.robot_key && robotLists.value.length) {
      formState.robot_key = robotLists.value[0].robot_key
    }
  })
}

function getImgUrl(name) {
  // 请注意，这不包括子目录中的文件
  return new URL(`../../assets/img/guide/${name}.png`, import.meta.url).href
}

const guideProcessTimer = ref(null)

onMounted(() => {
  getAllRobot()

  // 设置定时器，每5秒执行一次
  guideProcessTimer.value = setInterval(() => {
    guideStore.getUseGuideProcess()
  }, 5000)
})

// 组件卸载时清除定时器
onUnmounted(() => {
  if (guideProcessTimer.value) {
    clearInterval(guideProcessTimer.value)
    guideProcessTimer.value = null
  }
})
</script>

<style lang="less" scoped>
.guide-wrapper {
  width: 100%;
  height: 100%;
  overflow-y: auto;
  padding: 24px 0;
  .guide-content {
    width: 1000px;
    margin: 0 auto;
    background: #fff;
    border-radius: 16px 16px 12px 12px;
    height: fit-content;
    overflow: hidden;
  }

  .guide-header {
    height: 98px;
    width: 100%;
    background: linear-gradient(90deg, #469fff 0%, #3672ff 65.54%, #9c75ff 100%);
    img {
      width: 100%;
      height: 100%;
    }
  }
  .guide-body {
    background: #fff;
    display: flex;
    .guide-left {
      width: 256px;
      padding: 24px;
      border-right: 1px solid #d9d9d9;
    }
    .guide-right {
      flex: 1;
      padding: 24px;
    }
  }
}

.guide-left {
  .process-box {
    margin-bottom: 28px;
    .ant-progress-line {
      margin: 0;
    }
  }
  .menu-list-box {
    display: flex;
    flex-direction: column;
    gap: 8px;
    .menu-item {
      padding: 12px 16px;
      border-radius: 6px;
      display: flex;
      align-items: center;
      gap: 12px;
      cursor: pointer;
      transition: all 0.2s ease-in-out;
      &.active {
        background: #e5efff;
      }
      &:hover {
        background: #e5efff;
      }
      .icon-box {
        width: 40px;
        display: flex;
        align-items: center;
        justify-content: center;
      }
      .desc-box {
        display: flex;
        flex-direction: column;
        gap: 2px;
        .title-text {
          font-size: 14px;
          color: #000000;
          font-weight: 600;
          line-height: 22px;
          gap: 8px;
        }
        .desc-text {
          color: #8c8c8c;
          font-size: 12px;
          line-height: 20px;
        }
      }
    }
  }
}

.guide-right {
  .step-list {
    display: flex;
    align-items: center;
    gap: 8px;
    .step-item {
      display: flex;
      align-items: center;
      gap: 8px;
      .step-icon {
        font-size: 24px;
        .num-text {
          display: flex;
          align-items: center;
          justify-content: center;
          width: 24px;
          height: 24px;
          border-radius: 50%;
          background: var(---, #2475fc);
          color: #fff;
          font-size: 16px;
        }
      }
      .step-text {
        color: #000000d9;
        font-size: 14px;
      }
    }
  }
  .to-config-btn {
    margin-top: 32px;
  }
  .tag-list {
    display: flex;
    align-items: center;
    gap: 4px;
    color: #595959;
    margin-top: 16px;
    .tag-item {
      display: flex;
      align-items: center;
      gap: 4px;
      .tag-text-box {
        border-radius: 4px;
        padding: 0 6px;
        background: #e4e6eb;
        color: #164799;
        font-size: 14px;
        display: flex;
        align-items: center;
      }
    }
  }

  .tag-desc-box {
    margin-top: 8px;
    color: #8c8c8c;
    font-size: 14px;
    line-height: 22px;
  }

  .btn-block-wrapper {
    margin-top: 16px;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .guide-preview-box {
    margin-top: 16px;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .chat-test-box {
    margin-top: 16px;
    &::v-deep(.ant-form-item) {
      margin-bottom: 16px;
      .ant-form-item-label {
        padding-bottom: 4px;
      }
    }
  }

  .test-result {
    padding: 16px 16px 141px 16px;
    width: 455px;
    height: fit-content;
    min-height: 260px;
    background: var(--09, #f2f4f7);
    border-radius: 6px;
    .test-result-title {
      display: flex;
      align-items: center;
      gap: 8px;
      color: #262626;
      font-size: 14px;
      .result-tag {
        padding: 0 6px;
        border-radius: 6px;
        display: flex;
        align-items: center;
        gap: 2px;
        &.success {
          background: #acfcd5;
          color: #21a665;
        }
        &.error {
          background: #fbddde;
          color: #fb363f;
        }
      }
    }
    .result-content {
      margin-top: 16px;
      color: #3a4559;
      font-size: 14px;
      line-height: 22px;
      white-space: pre-wrap;
    }
  }
}
</style>
