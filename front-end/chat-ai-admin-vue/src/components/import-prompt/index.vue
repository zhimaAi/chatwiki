<style>
.add-prommpt-modal .ant-modal .ant-modal-content {
  padding: 0;
}
</style>

<template>
  <div>
    <a-modal
      v-model:open="open"
      :title="null"
      :width="950"
      :footer="null"
      wrapClassName="add-prommpt-modal"
    >
      <div class="prompt-wrapper">
        <div class="group-box">
          <div class="main-title">{{ t('title_import_prompt') }}</div>
          <cu-scroll style="height: 480px">
            <div class="group-list-box">
              <div
                class="group-list"
                :class="{ active: item.id == groupId }"
                @click="handleChangeGroup(item)"
                v-for="item in groupList"
                :key="item.id"
              >
                <div class="group-name">
                  <div class="name-text">{{ item.group_name }}</div>
                  <a-tooltip v-if="item.group_desc">
                    <template #title>{{ item.group_desc }}</template>
                    <InfoCircleOutlined />
                  </a-tooltip>
                </div>
              </div>
            </div>
          </cu-scroll>
        </div>
        <div class="content-box">
          <div class="opration-btn-block">
            <a-button :loading="btnLoading" @click="refreshList"
              ><SyncOutlined v-if="!btnLoading" />{{ t('btn_refresh_list') }}</a-button
            >
            <a-button @click="toAddPage">{{ t('btn_manage_knowledge_base') }}</a-button>
          </div>
          <div class="content-scroll-box">
            <cu-scroll style="padding: 0 24px">
              <div class="loading-box" v-if="isLoading">
                <a-spin></a-spin>
              </div>
              <div class="prompt-list-box">
                <div class="prompt-list" v-for="item in lists" :key="item.id">
                  <div class="prompt-header">
                    <div class="prompt-title">
                      <div class="prompt-type diy" v-if="item.prompt_type == 0">{{ t('label_custom') }}</div>
                      <div class="prompt-type" v-if="item.prompt_type == 1">{{ t('label_structured') }}</div>
                      <div class="title">{{ item.title }}</div>
                    </div>
                    <div class="right-btn-box">
                      <a-tooltip :title="item.isHide ? t('tooltip_expand') : t('tooltip_collapse')">
                        <div class="hover-btn-box" @click="handleHide(item)">
                          <DownOutlined v-if="item.isHide" />
                          <UpOutlined v-else />
                        </div>
                      </a-tooltip>
                      <div class="hover-btn-box primary" @click="handleImport(item)">
                        <CheckCircleOutlined />{{ t('btn_import') }}
                      </div>
                    </div>
                  </div>
                  <div class="prompt-content-box">
                    <div v-if="item.prompt_type == 0" class="prompt-content">
                      <template v-if="item.isHide">
                        {{ item.prompt.slice(0, 200) }}
                        <span
                          @click="handleHide(item)"
                          style="cursor: pointer"
                          v-if="item.prompt.length > 200"
                          >...</span
                        >
                      </template>
                      <template v-else>{{ item.prompt }}</template>
                    </div>
                    <div v-else class="structure-list-box">
                      <div class="structure-list">
                        <div class="structure-title">{{ item.prompt_struct.role.subject }}</div>
                        <div class="structure-content">{{ item.prompt_struct.role.describe }}</div>
                      </div>
                      <div class="structure-list">
                        <div class="structure-title">{{ item.prompt_struct.task.subject }}</div>
                        <div class="structure-content">
                          {{ item.prompt_struct.task.describe }}
                          <span @click="handleHide(item)" style="cursor: pointer" v-if="item.isHide"
                            >...</span
                          >
                        </div>
                      </div>
                      <template v-if="!item.isHide">
                        <div class="structure-list">
                          <div class="structure-title">
                            {{ item.prompt_struct.constraints.subject }}
                          </div>
                          <div class="structure-content">
                            {{ item.prompt_struct.constraints.describe }}
                          </div>
                        </div>
                        <div class="structure-list" v-if="item.prompt_struct.skill">
                          <div class="structure-title">
                            {{ item.prompt_struct.skill.subject }}
                          </div>
                          <div class="structure-content">
                            {{ item.prompt_struct.skill.describe }}
                          </div>
                        </div>
                        <div class="structure-list">
                          <div class="structure-title">{{ item.prompt_struct.output.subject }}</div>
                          <div class="structure-content">
                            {{ item.prompt_struct.output.describe }}
                          </div>
                        </div>
                        <div class="structure-list">
                          <div class="structure-title">{{ item.prompt_struct.tone.subject }}</div>
                          <div class="structure-content">
                            {{ item.prompt_struct.tone.describe }}
                          </div>
                        </div>
                        <div
                          class="structure-list"
                          v-for="custom in item.prompt_struct.custom"
                          :key="custom.subject"
                        >
                          <div class="structure-title">{{ custom.subject }}</div>
                          <div class="structure-content">{{ custom.describe }}</div>
                        </div>
                      </template>
                    </div>
                  </div>
                </div>
              </div>
              <a-empty style="padding-top: 100px" v-if="lists.length == 0">
                <template #description>
                  <div class="empty-content">
                    <div class="title">{{ t('msg_no_prompts') }}</div>
                    <div class="desc">
                      {{ t('msg_no_prompts_desc') }}
                    </div>
                    <div>
                      <a-button type="primary" @click="toAddPage">{{ t('btn_add_now') }}</a-button>
                    </div>
                  </div>
                </template>
              </a-empty>
            </cu-scroll>
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, h } from 'vue'
import {
  InfoCircleOutlined,
  DownOutlined,
  UpOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  SyncOutlined
} from '@ant-design/icons-vue'
import { getPromptLibraryGroup, getPromptLibraryItems } from '@/api/user/index.js'
import { Modal } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('components.import-prompt.index')
const open = ref(false)

const emit = defineEmits(['ok'])
const show = () => {
  open.value = true
  getGroupList()
  getLists()
}
const handleImport = (item) => {
  Modal.confirm({
    title: t('title_import_confirm'),
    icon: h(ExclamationCircleOutlined),
    content: t('msg_import_confirm'),
    onOk() {
      emit('ok', item)
      open.value = false
    },
    onCancel() {
      console.log('Cancel')
    }
  })
}

const groupId = ref(-1)
const groupList = ref([])
const getGroupList = () => {
  getPromptLibraryGroup().then((res) => {
    groupList.value = [
      {
        id: -1,
        group_name: t('label_all')
      },
      {
        id: 0,
        group_name: t('label_default_group')
      },
      ...res.data
    ]
  })
}

const handleChangeGroup = (item) => {
  groupId.value = item.id
  getLists()
}

const lists = ref([])
const isLoading = ref(false)
const btnLoading = ref(false)
const refreshList = () => {
  btnLoading.value = true
  getLists()
  getGroupList()
}
const getLists = () => {
  isLoading.value = true
  getPromptLibraryItems({
    group_id: groupId.value,
    page: 1,
    size: 9999
  })
    .then((res) => {
      let data = res.data.list || []
      data = data.map((item) => {
        return {
          ...item,
          isHide: true,
          prompt_struct: item.prompt_struct ? JSON.parse(item.prompt_struct) : {}
        }
      })
      lists.value = data
    })
    .finally(() => {
      isLoading.value = false
      btnLoading.value = false
    })
}

const handleHide = (item) => {
  item.isHide = !item.isHide
}

const toAddPage = () => {
  window.open('#/user/prompt-library')
}

defineExpose({
  show
})
</script>
<style lang="less" scoped>
.prompt-wrapper {
  display: flex;
  height: 550px;
  overflow: hidden;
  .group-box {
    width: 220px;
    height: 100%;
    border-right: 1px solid var(--06, #d9d9d9);
    .main-title {
      height: 72px;
      display: flex;
      align-items: center;
      padding-left: 24px;
      color: #262626;
      font-size: 16px;
      font-weight: 600;
      line-height: 24px;
    }
    .group-list-box {
      padding: 0 16px;
    }
    .group-list {
      margin-bottom: 4px;
      height: 32px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 5px 8px;
      font-size: 14px;
      color: #595959;
      cursor: pointer;
      border-radius: 6px;
      transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
      .group-name {
        display: flex;
        align-items: center;
        gap: 4px;
      }
      .name-text {
        max-width: 160px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
      &:hover {
        background: #f2f4f7;
        .hover-btn-box {
          display: flex;
        }
      }
      &.active {
        background: #e6efff;
        color: #2475fc;
      }
    }
  }
}
.opration-btn-block {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  padding-right: 24px;
}
.content-box {
  flex: 1;
  overflow: hidden;
  height: 100%;
  padding: 24px 0;
  padding-top: 48px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  .content-scroll-box {
    flex: 1;
    overflow: hidden;
  }
}

.loading-box {
  position: absolute;
  top: 160px;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.prompt-list-box {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 8px;
  .prompt-list {
    padding: 16px;
    background: var(--09, #f2f4f7);
    border-radius: 6px;
  }
  .prompt-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    color: #8c8c8c;
    .prompt-title {
      display: flex;
      align-items: center;
      gap: 4px;
      line-height: 22px;
      .ant-divider-vertical {
        margin: 0;
      }
    }
    .title {
      color: #262626;
      font-weight: 600;
    }
    .prompt-type {
      display: flex;
      align-items: center;
      height: 22px;
      border-radius: 6px;
      padding: 0 6px;
      width: fit-content;
      background: #e4d2fa;
      color: #7000ff;
      &.diy {
        background: #d4e3fc;
        color: #2475fc;
      }
    }
    .right-btn-box {
      display: flex;
      align-items: center;
      gap: 8px;
    }
  }
}

.prompt-content-box {
  margin-top: 6px;
  color: #3a4559;
  line-height: 22px;
}
.prompt-content {
  white-space: pre-line;
}

.structure-list-box {
  .structure-list {
    margin-bottom: 6px;
    border-bottom: 1px solid #d9d9d9;
    .structure-title {
      color: #242933;
      line-height: 22px;
    }
    .structure-content {
      margin-top: 4px;
      margin-bottom: 6px;
      color: #3a4559;
      font-size: 14px;
      line-height: 22px;
      white-space: pre-line;
    }
  }
}

.hover-btn-box {
  display: flex;
  align-items: center;
  gap: 4px;
  height: 24px;
  width: fit-content;
  border-radius: 6px;
  padding: 0 6px;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);

  &:hover {
    background: #e3e5ea;
  }
  &.primary:hover {
    color: #2475fc;
  }
}

.empty-content {
  .title {
    font-size: 16px;
    font-weight: 600;
    color: #262626;
  }
  .desc {
    color: #8c8c8c;
    font-size: 14px;
    margin: 8px 0;
  }
}
</style>
