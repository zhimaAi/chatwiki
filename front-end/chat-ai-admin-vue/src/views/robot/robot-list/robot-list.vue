<style lang="less" scoped>
.robot-page {
  padding-bottom: 8px;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  .page-title {
    line-height: 28px;
    padding: 16px 0;
    font-size: 20px;
    font-weight: 600;
    color: #000000;
  }
  .list-group-box {
    display: flex;
    gap: 16px;
    flex: 1;
    overflow: hidden;
  }

  .group-list-box {
    width: 256px;
    height: calc(100% - 8px);
    border: 1px solid #d9d9d9;
    border-radius: 6px;
    margin-top: 8px;
    padding: 16px;
    padding-right: 0;
    position: relative;
    &.hide-group {
      width: 0;
      padding: 0;
      border-left: 0;
      border-top: 0;
      border-bottom: 0;
    }
    .hide-group-box {
      position: absolute;
      right: -8px;
      top: 40%;
      height: 50px;
      width: 13px;
      border-radius: 8px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #e5e5ea;
      color: #8c8c8c;
      cursor: pointer;
      opacity: 0.78;
      &:hover {
        opacity: 1;
      }
    }
    .head-title {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 4px;
      font-size: 16px;
      font-weight: 600;
      color: #262626;
    }

    .search-box {
      margin-top: 16px;
    }

    .classify-box {
      flex: 1;
      overflow: hidden;
      font-size: 14px;
      .classify-item {
        height: 32px;
        padding: 0 8px;
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-top: 4px;
        cursor: pointer;
        border-radius: 6px;
        color: #595959;
        transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);

        .classify-title {
          flex: 1;
          text-overflow: ellipsis;
          overflow: hidden;
          white-space: nowrap;
        }
        .num {
          display: block;
        }

        .btn-box {
          display: none;
        }
        &:hover {
          background: #f2f4f7;
          .num {
            display: none;
          }
          .num.num-block {
            display: block;
          }
          .btn-box {
            display: block;
          }
        }
        &.active {
          color: #2475fc;
          background: #e6efff;
        }
      }
    }

    .hover-btn-wrap {
      width: fit-content;
      height: 24px;
      border-radius: 6px;
      padding: 0 4px;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
      &:hover {
        background: #e4e6eb;
      }
    }
  }

  .robot-item {
    padding: 8px;
    border-bottom: 1px solid #ccc;
  }

  ::v-deep(.ant-tabs-nav::before) {
    border-bottom: 0;
  }

  .list-toolbar {
    display: flex;
    justify-content: space-between;
    margin-bottom: 8px;
  }
  .robot-scroll-box{
    flex: 1;
    overflow: hidden;
    height: 100%;
  }

  .list-box {
    flex: 1;
    width: 100%;
    display: flex;
    flex-flow: row wrap;
    margin: 0 -8px 0 -8px;
    height: fit-content;
  }

  .list-item-wrapper {
    padding: 8px;
    width: 25%;
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

    .robot-info {
      display: flex;
      // align-items: start;
      align-items: center;
    }

    .robot-avatar {
      width: 48px;
      height: 48px;
      border-radius: 14px;
      overflow: hidden;
    }

    .robot-info-content {
      flex: 1;
      padding-left: 12px;
      overflow: hidden;
    }

    .robot-name {
      height: 24px;
      line-height: 24px;
      margin-bottom: 4px;
      font-size: 16px;
      font-weight: 600;
      color: rgb(38, 38, 38);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .robot-desc {
      height: 44px;
      line-height: 22px;
      margin-top: 12px;
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

    .robot-type-tag {
      display: inline-block;
      height: 22px;
      line-height: 20px;
      padding: 0 8px;
      font-size: 12px;
      font-weight: 400;
      border-radius: 6px;
      color: rgb(36, 117, 252);
      border: 1px solid #cde0ff;
    }

    .robot-action {
      margin-top: 12px;
      font-size: 14px;
      height: 24px;
      color: #2475fc;
      display: flex;
      align-items: center;

      .robot-action-item {
        display: flex;
        align-items: center;
        height: 100%;
        margin-right: 12px;
        padding: 4px;
        border-radius: 6px;
        cursor: pointer;
        color: #595959;
        transition: all 0.2s;
        &.drag-btn {
          cursor: move;
        }
      }
      .robot-action-item:hover {
        background: #e4e6eb;
      }

      .action-icon {
        font-size: 16px;
      }
    }
  }

  .add-robot {
    display: flex;
    align-items: center;
    justify-content: center;
    line-height: 22px;
    color: #3a4559;
    cursor: pointer;

    .add-icon {
      font-size: 16px;
    }

    .add-text {
      padding-left: 4px;
      font-size: 14px;
    }
  }
}

.empty-box {
  width: 100%;
  padding-top: 100px;
}
.create-action {
  display: flex;
  align-items: center;
  .icon {
    width: 20px;
    height: 20px;
    margin-right: 8px;
  }
}
.toolbar-box {
  position: relative;
  padding-right: 32px;
}
.mr-16 {
  margin-right: 16px;
}

// 大于1920px
@media screen and (min-width: 1920px) {
  .robot-page {
    .list-box {
      .list-item-wrapper {
        width: 20%;
      }
    }
  }
}

.desc-box {
  width: 200px;
  position: absolute;
  bottom: -29px;
  left: -35px;
  margin-top: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
  color: #8c8c8c;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;

  .link {
    cursor: pointer;
    color: #2475fc;
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;

    &:hover {
      opacity: 0.8;
    }
  }
}
</style>

<template>
  <div class="robot-page">
<!--    <h3 class="page-title">应用</h3>-->
    <page-alert class="mt16" :title="t('title_usage_guide')">
      <div>
        <p>
          {{ t('usage_guide_content_1') }}
        </p>
        <p>{{ t('usage_guide_content_2') }}</p>
      </div>
    </page-alert>

    <div class="list-toolbar mt16">
      <div class="toolbar-box">
        <ListTabs :tabs="tabs" @change="initData" v-model:value="activeKey" />
      </div>
      <div class="toolbar-box" style="text-align: right;">
        <a-dropdown v-if="robotCreate">
          <a-button type="primary" @click.prevent="">
            <template #icon>
              <PlusOutlined />
            </template>
            {{ t('btn_create_app') }}
          </a-button>
          <template #overlay>
            <a-menu>
              <a-menu-item @click.prevent="toAddRobot(0)">
                <span class="create-action">
                  <img class="icon" :src="DEFAULT_ROBOT_AVATAR" alt="" />
                  <span>{{ t('label_chat_bot') }}</span>
                </span>
              </a-menu-item>
              <a-menu-item @click.prevent="toAddRobot(1)">
                <span class="create-action">
                  <img class="icon" :src="DEFAULT_WORKFLOW_AVATAR" alt="" />
                  <span>{{ t('label_workflow') }}</span>
                </span>
              </a-menu-item>
              <a-menu-item @click.prevent="toAddThirdMcp">
                <span class="create-action">
                  <img class="icon" :src="DEFAULT_MCP_AVATAR" alt="" />
                  <span>{{ t('label_external_mcp') }}</span>
                </span>
              </a-menu-item>
              <a-menu-item @click.prevent="toAddHttpTool">
                <span class="create-action">
                  <img class="icon" :src="DEFAULT_MCP_AVATAR" alt="" />
                  <span>{{ t('label_http_tool') }}</span>
                </span>
              </a-menu-item>
              <a-menu-item @click="toImportCsl">
                <span class="create-action">
                  <img class="icon" :src="DEFAULT_IMPORT_CSL_AVATAR" alt="">
                  <span>{{ t('label_import_csl_create') }}</span>
                </span>
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
        <div v-if="robotCreate && activeKey == '3'" class="desc-box"><div class="link" @click="goMcpSquare">{{ t('link_more_mcp_square') }}</div></div>
      </div>
    </div>

    <McpPanel v-if="activeKey == 3" ref="mcpPanelRef"/>
    <HttpToolBox v-else-if="activeKey == 6" ref="httpToolBoxRef" />
    <div v-else class="list-group-box">
      <div class="group-list-box" :class="{ 'hide-group': isHideGroup }">
        <cu-scroll style="padding-right: 16px">
          <div class="group-head-box" v-if="!isHideGroup">
            <div class="head-title">
              <div>{{ t('title_robot_group') }}</div>
              <a-tooltip :title="t('tooltip_new_group')">
                <div class="hover-btn-wrap" @click="openGroupModal({})"><PlusOutlined /></div>
              </a-tooltip>
            </div>
            <div class="search-box">
              <a-input
                v-model:value="groupSearchKey"
                allowClear
                :placeholder="t('ph_search_group')"
                style="width: 100%"
              >
                <template #suffix>
                  <SearchOutlined @click.stop="" />
                </template>
              </a-input>
            </div>
          </div>
          <div class="classify-box">
            <div
              class="classify-item"
              @click="handleChangeGroup(item)"
              :class="{ active: item.id == group_id }"
              v-for="item in filterGroupLists"
              :key="item.id"
            >
              <div class="classify-title">{{ item.group_name }}</div>
              <div class="right-content">
                <div class="num" :class="{ 'num-block': item.id <= 0 }">{{ item.total }}</div>
                <div class="btn-box" v-if="item.id > 0">
                  <a-dropdown placement="bottomRight">
                    <div class="hover-btn-wrap">
                      <EllipsisOutlined />
                    </div>
                    <template #overlay>
                      <a-menu>
                        <a-menu-item>
                          <div @click.stop="openGroupModal(item)">{{ t('btn_rename') }}</div>
                        </a-menu-item>
                        <a-menu-item>
                          <div @click.stop="handleDelGroup(item)">{{ t('btn_delete') }}</div>
                        </a-menu-item>
                      </a-menu>
                    </template>
                  </a-dropdown>
                </div>
              </div>
            </div>
          </div>
        </cu-scroll>
        <a-tooltip placement="right" :title="isHideGroup ? t('tooltip_expand_group') : t('tooltip_collapse_group')">
          <div class="hide-group-box" @click="handleChangeHideGroup">
            <LeftOutlined v-if="!isHideGroup" />
            <RightOutlined v-else />
          </div>
        </a-tooltip>
      </div>
      <div class="robot-scroll-box">
        <div class="empty-box" v-if="filterRobotList.length == 0">
          <a-empty></a-empty>
        </div>
        <cu-scroll style="padding-right: 16px">
          <draggable
            v-model="filterRobotList"
            handle=".drag-btn"
            item-key="id"
            :disabled="moving"
            class="list-box"
            @end="onDragEnd"
          >
            <template #item="{ element: item, index }">
              <div
                class="list-item-wrapper"
                :key="item.id"
                @click="toEditRobot(item)"
              >
                <div class="list-item">
                  <div class="robot-info">
                    <img class="robot-avatar" :src="item.robot_avatar" alt="" />
                    <div class="robot-info-content">
                      <div class="robot-name">{{ item.robot_name }}</div>
                      <div class="robot-type-tag">
                        {{ item.application_type == 0 ? t('label_chat_bot') : t('label_workflow') }}
                      </div>
                    </div>
                  </div>
                  <a-tooltip :title="getTooltipTitle(item.robot_intro, item, 14, 2, 10)" placement="top">
                    <div class="robot-desc" :ref="el => setDescRef(el, item)">{{ item.robot_intro }}</div>
                  </a-tooltip>
                  <div class="robot-action" @click.stop>
                    <!-- <div class="robot-action-item" @click="toEditRobot(item)"><svg-icon class="action-icon" name="jibenpeizhi" /></div> -->
                    <!-- <div class="robot-action-item" @click="toTestPage(item)"><svg-icon class="action-icon" name="cmd" /></div> -->
                    <a-tooltip :title="t('tooltip_drag_sort')">
                      <div class="robot-action-item drag-btn" @click.stop>
                        <svg-icon class="action-icon" name="dragv2"></svg-icon>
                      </div>
                    </a-tooltip>
                    <a-tooltip :title="t('tooltip_chat_test')">
                      <div class="robot-action-item" @click.stop="toTestPage(item)">
                        <svg-icon class="action-icon" name="chat"></svg-icon>
                      </div>
                    </a-tooltip>
                    <a-tooltip :title="t('tooltip_session_record')">
                      <div class="robot-action-item" @click.stop="toSessionRecordPage(item)">
                        <svg-icon class="action-icon" name="session"></svg-icon>
                      </div>
                    </a-tooltip>
                    <a-tooltip :title="t('tooltip_statistics')">
                      <div class="robot-action-item" @click.stop="toAnalysisPage(item)">
                        <svg-icon class="action-icon" name="analysis"></svg-icon>
                      </div>
                    </a-tooltip>

                    <a-dropdown>
                      <div class="robot-action-item" @click.stop>
                        <svg-icon class="action-icon" name="point-h"></svg-icon>
                      </div>
                      <template #overlay>
                        <a-menu>
                          <a-menu-item>
                            <a
                              class="delete-text-color"
                              href="javascript:;"
                              @click="handleDelete(item)"
                            >{{ t('btn_delete') }}</a
                            >
                          </a-menu-item>
                          <a-menu-item>
                            <div @click.stop="openEditGroupModal(item)">{{ t('btn_edit_group') }}</div>
                          </a-menu-item>
                          <a-menu-item>
                            <a class="copy-text-color" href="javascript:;" @click="handleCopy(item)"
                            >{{ t('btn_copy') }}</a
                            >
                          </a-menu-item>
                          <a-menu-item @click="handleExportCsl(item)">
                            <span>{{ t('btn_export_csl') }}</span>
                          </a-menu-item>
                        </a-menu>
                      </template>
                    </a-dropdown>
                  </div>
                </div>
              </div>
            </template>
          </draggable>
        </cu-scroll>
      </div>
    </div>

    <AddRobotAlert ref="addRobotAlertRef" @ok="initData" />
    <AddGroup ref="addGroupRef" @ok="initData" />
    <EditGroup ref="editGroupRef" @ok="initData" />
    <ImportCslAlert ref="importCslAlertRef" @ok="handleImportOk" />
    <ThirdMcpStore ref="thirdMcpRef" @ok="thirdMcpUpdate"/>
    <WorkflowExportModal ref="exportRef" @ok="(r, d) => runExportCsl(r, d)"/>
  </div>
</template>

<script setup>
import { DEFAULT_ROBOT_AVATAR, DEFAULT_WORKFLOW_AVATAR, DEFAULT_IMPORT_CSL_AVATAR, DEFAULT_MCP_AVATAR } from '@/constants/index.js'
import { usePermissionStore } from '@/stores/modules/permission'
import draggable from 'vuedraggable'
import {
  getRobotList,
  deleteRobot,
  robotCopy,
  getRobotGroupList,
  deleteRobotGroup,
  moveRobotSort,
} from '@/api/robot/index.js'
import { ref, onMounted, createVNode, computed, nextTick } from 'vue'
import {useRoute} from 'vue-router'
import {
  ExclamationCircleOutlined,
  PlusOutlined,
  SearchOutlined,
  EllipsisOutlined,
  LeftOutlined,
  RightOutlined
} from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import AddRobotAlert from './components/add-robot-alert.vue'
import ImportCslAlert from './components/import-csl-alert.vue'
import PageAlert from '@/components/page-alert/page-alert.vue'
import ListTabs from '@/components/cu-tabs/list-tabs.vue'
import { getRobotPermission } from '@/utils/permission'
import AddGroup from '@/views/robot/robot-list/components/add-group.vue'
import EditGroup from '@/views/robot/robot-list/components/edit-group.vue'
import { useOpenUrlWithToken } from '@/hooks/web'
import McpPanel from "@/views/robot/robot-list/components/mcp-panel.vue";
import ThirdMcpStore from "@/views/robot/robot-list/components/third-mcp-store.vue";
import HttpToolBox from "@/views/robot/robot-list/components/http-tool-box.vue";
import { useRouter } from 'vue-router'
import { setDescRef, getTooltipTitle, objectToQueryString } from '@/utils/index'
import WorkflowExportModal from "@/views/robot/robot-list/components/workflow-export-modal.vue";
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-list.robot-list')

const router = useRouter()

const route = useRoute()

const tabs = ref([
  {
    title: `${t('tab_all')} (0)`,
    value: '2'
  },
  {
    title: `${t('label_chat_bot')} (0)`,
    value: '0'
  },
  {
    title: `${t('label_workflow')} (0)`,
    value: '1'
  },
  {
    title: t('tab_mcp_plugin'),
    value: '3'
  },
  {
    title: t('tab_http_tool'),
    value: '6'
  }
])

const permissionStore = usePermissionStore()

let { role_permission, role_type } = permissionStore
const robotCreate = computed(() => role_type == 1 || role_permission.includes('RobotCreate'))


const exportRef = ref(null)
const mcpPanelRef = ref(null)
const activeKey = ref('2')
const httpToolBoxRef = ref(null)

const group_id = ref('')

const robotList = ref([])
const filterRobotList = ref([])
const moving = ref(false)

let hideGroupLocalKey = 'robot-list-hide-group-key'

const isHideGroup = ref(localStorage.getItem(hideGroupLocalKey) == 1)

const handleChangeHideGroup = () => {
  isHideGroup.value = !isHideGroup.value
  localStorage.setItem(hideGroupLocalKey, isHideGroup.value ? 1 : 0)
}

const initData = () => {
  getList()
  getGroupList()
}

const filterRobotInit = () => {
  let lists = robotList.value
  if (group_id.value != '') {
    lists = lists.filter((item) => item.group_id == group_id.value)
  }
  if (activeKey.value != 2) {
    lists = lists.filter((item) => item.application_type == activeKey.value)
  }
  filterRobotList.value = lists
}

function goMcpSquare () {
  const url = router.resolve({
    path: '/mcp/index'
  }).href
  window.open(url, '_blank')
}

const getList = () => {
    getRobotList()
    .then((res) => {
    let allNumber = 0
    let chatNumber = 0
    let workflowNumber = 0

    res.data.forEach((item) => {
      if (item.application_type == 0) {
        chatNumber++
      } else {
        workflowNumber++
      }
      allNumber++
    })

    tabs.value = [
      {
        title: `${t('tab_all')} (${allNumber})`,
        value: '2'
      },
      {
        title: `${t('label_chat_bot')} (${chatNumber})`,
        value: '0'
      },
      {
        title: `${t('label_workflow')} (${workflowNumber})`,
        value: '1'
      },
      {
        title: t('tab_mcp_plugin'),
        value: '3'
      },
      {
        title: t('tab_http_tool'),
        value: '6'
      }
    ]
    robotList.value = res.data
    filterRobotInit()
    })
    .catch(() => {})
}

const addRobotAlertRef = ref(null)
const importCslAlertRef = ref(null)
const thirdMcpRef = ref(null)
const toAddRobot = (val) => {
  // router.push({ name: 'addRobot' })
  addRobotAlertRef.value.open(val)
  addRobotAlertRef.value.setGroupId(group_id.value)
}

const toAddThirdMcp = () => {
  thirdMcpRef.value.show()
}

const toImportCsl = () => {
  importCslAlertRef.value.open()
}

const handleImportOk = (data) => {
  getList()
 toEditRobot(data)
}

const toEditRobot = ({ id, robot_key, application_type }) => {
  localStorage.setItem('last_local_robot_id', id)
  if (application_type == 0) {
    // router.push({ path: '/robot/config/basic-config', query: { id: id, robot_key: robot_key  } })
    window.open(
      `/#/robot/config/basic-config?id=${id}&robot_key=${robot_key}`,
      '_blank',
      'noopener'
    ) // 建议添加 noopener 防止安全漏洞
  } else {
    // router.push({ path: '/robot/config/workflow', query: { id: id, robot_key: robot_key  } })
    window.open(`/#/robot/config/workflow?id=${id}&robot_key=${robot_key}`, '_blank', 'noopener') // 建议添加 noopener 防止安全漏洞
  }
}
const toAddHttpTool = () => {
  activeKey.value = '6'
  nextTick(() => {
    // 直接调用子组件暴露的方法
    const httpComp = (typeof httpToolBoxRef.value !== 'undefined') ? httpToolBoxRef.value : null
    if (httpComp && httpComp.showStore) {
      httpComp.showStore()
    }
  })
}

const handleDelete = (data) => {
  let key = getRobotPermission(data.id)
  if (key != 4) {
    return message.error(t('msg_no_delete_permission'))
  }
  let secondsToGo = 3
  let modal = Modal.confirm({
    title: `${t('title_delete_robot')}${data.robot_name}`,
    icon: createVNode(ExclamationCircleOutlined),
    content: t('msg_confirm_delete_robot'),
    okText: secondsToGo + ' ' + t('btn_confirm'),
    okType: 'danger',
    cancelText: t('btn_cancel'),
    okButtonProps: {
      disabled: true
    },
    onOk() {
      onDelete(data)
    },
    onCancel() {}
  })

  let interval = setInterval(() => {
    if (secondsToGo == 1) {
      modal.update({
        okText: t('btn_confirm'),
        okButtonProps: {
          disabled: false
        }
      })

      clearInterval(interval)
      interval = undefined
    } else {
      secondsToGo -= 1

      modal.update({
        okText: secondsToGo + ' ' + t('btn_confirm'),
        okButtonProps: {
          disabled: true
        }
      })
    }
  }, 1000)
}

const onDelete = ({ id }) => {
  deleteRobot({ id }).then(() => {
    message.success(t('msg_delete_success'))
    getList()
  })
}

const handleCopy = (data) => {
  let secondsToGo = 3
  let modal = Modal.confirm({
    title: `${t('title_copy_robot')}${data.robot_name}`,
    icon: createVNode(ExclamationCircleOutlined),
    content: t('msg_confirm_copy_robot'),
    okText: secondsToGo + ' ' + t('btn_confirm'),
    okType: 'danger',
    cancelText: t('btn_cancel'),
    okButtonProps: {
      disabled: true
    },
    onOk() {
      onCopy(data)
    },
    onCancel() {}
  })

  let interval = setInterval(() => {
    if (secondsToGo == 1) {
      modal.update({
        okText: t('btn_confirm'),
        okButtonProps: {
          disabled: false
        }
      })

      clearInterval(interval)
      interval = undefined
    } else {
      secondsToGo -= 1

      modal.update({
        okText: secondsToGo + ' ' + t('btn_confirm'),
        okButtonProps: {
          disabled: true
        }
      })
    }
  }, 1000)
}

const onCopy = ({ id }) => {
  robotCopy({ id }).then(() => {
    message.success(t('msg_copy_success'))
    getList()
  })
}

// 导出为csl
const handleExportCsl = (data) => {
  exportRef.value.handle(data)
}

const runExportCsl = (data, params=null) => {
  const { openUrlWithToken } = useOpenUrlWithToken()
  params = params || {}
  params.id = data.id
  console.log('---', `/manage/robotExport?${objectToQueryString(params)}`)
  openUrlWithToken(`/manage/robotExport?${objectToQueryString(params)}`)
}

const toTestPage = (item) => {
  window.open(`#/robot/test?robot_key=${item.robot_key}&id=${item.id}`)
}

const toSessionRecordPage = (item) => {
  window.open(`#/robot/config/session-record?robot_key=${item.robot_key}&id=${item.id}`)
}

const toAnalysisPage = (item) => {
  window.open(`#/robot/config/statistical_analysis?robot_key=${item.robot_key}&id=${item.id}`)
}

const groupSearchKey = ref('')
const groupLists = ref([])

const filterGroupLists = computed(() => {
  return groupLists.value.filter((item) => item.group_name.includes(groupSearchKey.value))
})

const getGroupList = () => {
  getRobotGroupList({
    application_type: activeKey.value == 2 ? null : activeKey.value
  }).then((res) => {
    let lists = res.data || []
    let totalNumber = 0
    // 计算每个分组的机器人数量
    lists.forEach((group) => {
      totalNumber += +group.total
    })
    groupLists.value = [
      {
        group_name: t('tab_all'),
        total: totalNumber,
        id: ''
      },
      ...lists
    ]
  })
}

const handleChangeGroup = (item) => {
  group_id.value = item.id
  getList()
}

const addGroupRef = ref(null)
const openGroupModal = (data) => {
  addGroupRef.value.show({
    ...data
  })
}

const editGroupRef = ref(null)
const openEditGroupModal = (item) => {
  editGroupRef.value.show({
    ...item
  })
}

const handleDelGroup = (item) => {
  Modal.confirm({
    title: `${t('title_delete_group')}${item.group_name}`,
    icon: createVNode(ExclamationCircleOutlined),
    content: '',
    okText: t('btn_confirm_delete_group'),
    okType: 'danger',
    cancelText: t('btn_cancel'),
    onOk() {
      deleteRobotGroup({
        id: item.id
      }).then(() => {
        message.success(t('msg_delete_success'))
        getGroupList()
        if (group_id.value == item.id) {
          group_id.value = ''
          getList()
        }
      })
    }
  })
}

const thirdMcpUpdate = () => {
  if (mcpPanelRef.value) {
    mcpPanelRef.value.update()
  }
}

const onDragEnd = (e) => {
  if (e.newIndex === e.oldIndex) return
  moving.value = true
  let params = {
    move_robot_id: filterRobotList.value[e.newIndex].id,
    to_robot_id: 0
  }
  if (e.newIndex > 0) params.to_robot_id = filterRobotList.value[e.newIndex - 1].id
  moveRobotSort(params).then(res => {
    message.success(t('msg_operation_complete'))
  }).finally(() => {
    moving.value = false
  })
}

onMounted(() => {
  activeKey.value = route.query.active || '2'
  initData()
})
</script>
