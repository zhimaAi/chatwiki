<template>
  <div class="main-content-block">
    <div class="btn-block">
      <a-button type="primary" @click="handleOpenSelectLibraryAlert">{{ t('btn_relate_library') }}</a-button>
    </div>
    <div class="list-box">
      <div class="list-item-wrapper" v-for="item in selectedLibraryRows" :key="item.id">
        <div class="list-item" @click.stop="toEdit(item)">
          <img
            class="default-icon"
            v-if="item.id == robotInfo.default_library_id"
            src="@/assets/img/robot/default-allow.svg"
            alt=""
          />
          <div class="library-info">
            <img class="library-icon" :src="item.avatar" alt="" />
            <div class="library-info-content">
              <div class="library-title">{{ item.library_name }}</div>
              <div class="library-type">
                <span class="type-tag" v-if="item.type == 0">{{ t('label_normal_library') }}</span>
                <span class="type-tag" v-if="item.type == 1">{{ t('label_external_library') }}</span>
                <span class="type-tag" v-if="item.type == 2">{{ t('label_qa_library') }}</span>
                <span class="type-tag" v-if="item.type == 3">{{ t('label_official_account_library') }}</span>
                <a-tooltip v-if="neo4j_status">
                  <template #title
                    >{{ item.graph_switch == 0 ? t('msg_graph_disabled') : t('msg_graph_enabled') }}</template
                  >
                  <span class="type-tag graph-tag" :class="{ 'gray-tag': item.graph_switch == 0 }"
                    >Graph</span
                  >
                </a-tooltip>
              </div>
            </div>
          </div>
          <div class="item-body">
            <div class="library-desc">{{ item.library_intro }}</div>
          </div>

          <div class="item-footer">
            <div class="library-size">
              <span class="text-item">{{ t('label_docs') }}：{{ item.file_total }}</span>
              <span class="text-item">{{ t('label_size') }}：{{ item.file_size_str }}</span>
              <span class="text-item">{{ t('label_related_apps') }}：{{ item.robot_nums || 0 }}</span>
            </div>

            <div class="action-box" @click.stop v-if="item.id != robotInfo.default_library_id">
              <a-dropdown>
                <div class="action-item" @click.stop>
                  <svg-icon class="action-icon" name="point-h"></svg-icon>
                </div>
                <template #overlay>
                  <a-menu>
                    <a-menu-item v-if="item.type != 1">
                      <a @click.stop="handleSetDefaultLibrary(item)">{{ t('btn_set_default') }}</a>
                    </a-menu-item>
                    <a-menu-item>
                      <a class="delete-text-color" @click.stop="handleRemoveCheckedLibrary(item)"
                        >{{ t('btn_cancel_relation') }}</a
                      >
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </div>
          </div>
        </div>
      </div>
    </div>
    <LibrarySelectAlert ref="librarySelectAlertRef" :showWxType="!!wxAppLibary" @change="onChangeLibrarySelected" />
  </div>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'
import { getLibraryList } from '@/api/library/index'
import { relationLibrary } from '@/api/robot/index'
import { storeToRefs } from 'pinia'
import { ref, reactive, watchEffect, computed, toRaw, onMounted, createVNode } from 'vue'
import { useRoute } from 'vue-router'
import { Modal } from 'ant-design-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import LibrarySelectAlert from '@/views/robot/robot-config/basic-config/components/associated-knowledge-base/library-select-alert.vue'
import { useCompanyStore } from '@/stores/modules/company'
const companyStore = useCompanyStore()
import { LIBRARY_NORMAL_AVATAR, LIBRARY_QA_AVATAR } from '@/constants/index'
import { formatFileSize } from '@/utils/index'
import { useRobotStore } from '@/stores/modules/robot'
import { message } from 'ant-design-vue'
import {getSpecifyAbilityConfig} from "@/api/explore/index.js";

const { t } = useI18n('views.robot.robot-config.library-config.related-library')
const neo4j_status = computed(() => {
  return companyStore.companyInfo?.neo4j_status == 'true'
})

const query = useRoute().query

const robotStore = useRobotStore()
const { getRobot } = robotStore
const { robotInfo } = storeToRefs(robotStore)

const formState = reactive({
  library_ids: [],
  default_library_id: ''
})

const libraryList = ref([])
const wxAppLibary = ref(null)

const librarySelectAlertRef = ref(null)
const selectedLibraryRows = computed(() => {
  let lists = libraryList.value.filter((item) => {
    return formState.library_ids.includes(item.id)
  })
  // robotInfo.value.default_library_id
  // 将默认知识库放在第一位
  if (robotInfo.value.default_library_id) {
    let defaultLibrary = lists.find((item) => item.id == robotInfo.value.default_library_id)
    if (defaultLibrary) {
      lists = lists.filter((item) => item.id != robotInfo.value.default_library_id)
      lists.unshift(defaultLibrary)
    }
  }

  if (!wxAppLibary.value) lists = lists.filter((item) => item.type != 3)
  return lists
})

// 移除知识库
const handleRemoveCheckedLibrary = (item) => {
  let index = formState.library_ids.indexOf(item.id)

  formState.library_ids.splice(index, 1)

  onSave()
}

const handleSetDefaultLibrary = (item) => {
  Modal.confirm({
    title: t('msg_confirm_set_default', { library_name: item.library_name }),
    // icon: createVNode(ExclamationCircleOutlined),
    icon: null,
    content: createVNode(
      'div',
      {
        style: 'color:red;'
      },
      t('msg_one_default_library')
    ),
    onOk() {
      formState.default_library_id = item.id
      onSave(formState.default_library_id)
    },
    onCancel() {}
  })
}

const handleOpenSelectLibraryAlert = () => {
  librarySelectAlertRef.value.open([...formState.library_ids])
}

const onChangeLibrarySelected = (checkedList) => {
  formState.library_ids = [...checkedList]

  onSave()
}

const onSave = (default_library_id) => {
  let formData = { ...toRaw(formState) }

  formData.library_ids = formData.library_ids.join(',')
  relationLibrary({
    library_ids: formData.library_ids,
    default_library_id: default_library_id || null,
    id: query.id
  }).then((res) => {
    message.success(t('msg_operation_success'))
    getRobot(query.id)
  })
}

watchEffect(() => {
  formState.library_ids = robotInfo.value.library_ids.split(',')
  formState.default_library_id = robotInfo.value.default_library_id || ''
})

const toEdit = (data) => {
  if (data.type == '1') {
    window.open(`/#/public-library/config?library_id=${data.id}`, '_blank', 'noopener')
  } else {
    window.open(`/#/library/details?id=${data.id}`, '_blank', 'noopener')
  }
}

const getList = async () => {
  const res = await getLibraryList({ type: '', show_open_docs: 1 })
  if (res) {
    let data = res.data || []

    data.forEach((item) => {
      item.file_size_str = formatFileSize(item.file_size)

      if (!item.avatar) {
        item.avatar = item.type == 0 ? LIBRARY_NORMAL_AVATAR : LIBRARY_QA_AVATAR
      }
    })

    libraryList.value = data
  }
}

const getWxLb = async () => {
  // 公众号知识库是否开启
  await getSpecifyAbilityConfig({ability_type: 'library_ability_official_account'}).then((res) => {
    let _data = res?.data || {}
    if (_data?.user_config?.switch_status == 1) {
      wxAppLibary.value = _data
    }
  })
}


onMounted(() => {
  getWxLb()
  getList()
})
</script>

<style lang="less" scoped>
.main-content-block {
  padding: 16px;
}

.list-box {
  display: flex;
  flex-flow: row wrap;
  margin: 0 -8px 0 -8px;
}
.list-item-wrapper {
  padding: 8px;
  width: 25%;
}
.default-icon {
  position: absolute;
  left: 0;
  top: 0;
  width: 38px;
}
.list-item {
  position: relative;
  width: 100%;
  padding: 24px;
  border: 1px solid #e4e6eb;
  border-radius: 12px;
  background-color: #fff;
  transition: all 0.25s;
  cursor: pointer;

  &:hover {
    box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.12);
  }

  .library-info {
    position: relative;
    display: flex;
    align-items: center;

    .library-icon {
      width: 52px;
      height: 52px;
      border-radius: 14px;
      overflow: hidden;
    }

    .library-info-content {
      flex: 1;
      padding-left: 12px;
      overflow: hidden;
    }
  }

  .library-title {
    height: 24px;
    line-height: 24px;
    margin-bottom: 4px;
    font-size: 16px;
    font-weight: 600;
    color: #262626;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .library-type {
    display: flex;
    .type-tag {
      height: 22px;
      line-height: 20px;
      padding: 0 8px;
      font-size: 12px;
      font-weight: 400;
      border-radius: 6px;
      color: rgb(36, 117, 252);
      border: 1px solid #cde0ff;
    }
    .graph-tag {
      margin-left: 4px;
      &.gray-tag {
        border: 1px solid #00000026;
        background: #0000000a;
        color: #bfbfbf;
      }
    }
  }
  .item-body {
    margin-top: 12px;
  }
  .library-desc {
    height: 44px;
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: rgb(89, 89, 89);
    // 超出2行显示省略号
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
  }
  .item-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 14px;
    color: #7a8699;
    min-height: 24px;
  }
  .library-size {
    display: flex;
    line-height: 20px;
    font-size: 12px;
    font-weight: 400;
    gap: 8px;
    color: #7a8699;
    flex-wrap: wrap;
  }

  .action-box {
    font-size: 14px;
    height: 24px;
    color: #2475fc;
    display: flex;
    align-items: center;

    .action-item {
      display: flex;
      align-items: center;
      height: 100%;
      padding: 4px;
      border-radius: 6px;
      cursor: pointer;
      color: #595959;
      transition: all 0.2s;
    }
    .action-item:hover {
      background: #e4e6eb;
    }

    .action-icon {
      font-size: 16px;
    }
  }
}

// 大于1920px
@media screen and (min-width: 1920px) {
  .list-box {
    .list-item-wrapper {
      width: 20%;
    }
  }
}
</style>
