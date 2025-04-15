<style lang="less" scoped>
.robot-page {
  padding-bottom: 8px;
  .page-title{
    line-height: 28px;
    margin: 16px 0;
    font-size: 20px;
    font-weight: 600;
    color: #000000;
  }
  .robot-item {
    padding: 8px;
    border-bottom: 1px solid #ccc;
  }

  ::v-deep(.ant-tabs-nav::before){
    border-bottom: 0;
  }

  .list-toolbar{
    display: flex;
    justify-content: space-between;
    margin-bottom: 8px;
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

  .list-item {
    position: relative;
    width: 100%;
    padding: 24px;
    border: 1px solid #E4E6EB;
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
      width: 52px;
      height: 52px;
      border-radius: 14px;
      overflow: hidden;
    }

    .robot-info-content {
      flex: 1;
      padding-left: 12px;
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

    .robot-type-tag{
      display: inline-block;
      height: 22px;
      line-height: 20px;
      padding: 0 8px;
      font-size: 12px;
      font-weight: 400;
      border-radius: 6px;
      color: rgb(36, 117, 252);
      border: 1px solid #CDE0FF;
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
      }
      .robot-action-item:hover {
        background: #E4E6EB;
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

.create-action{
  display: flex;
  align-items: center;
 .icon{
    width: 20px;
    height: 20px;
    margin-right: 8px;
 }
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
</style>

<template>
  <div class="robot-page">
    <h3 class="page-title">应用</h3>
    <page-alert class="mb-16" title="使用说明">
      <div>
        <p>
          1、应用包括两种类型：聊天机器人、工作流。聊天机器人适合新手用户，关联知识库后只需简单配置，即可创建一个基于私有知识库的问答机器人。工作流适合高级用户，利用系统预制节点自定义任务流程，适合解决复杂任务。
        </p>
        <p>2、发现模块也提供了丰富模版，可以根据需要选择适合模版快速创建应用。</p>
      </div>
    </page-alert>

    <div class="list-toolbar">
      <div class="toolbar-box">
        <ListTabs :tabs="tabs" v-model:value="activeKey" />
      </div>
      <div class="toolbar-box">
        <a-dropdown v-if="robotCreate">
          <a-button type="primary" @click.prevent="" >
            <template #icon>
              <PlusOutlined />
            </template>
            创建应用
          </a-button>
          <template #overlay>
            <a-menu>
              <a-menu-item @click.prevent="toAddRobot(0)">
                <span class="create-action">
                  <img class="icon" :src="DEFAULT_ROBOT_AVATAR" alt="">
                  <span>聊天机器人</span>
                </span>
              </a-menu-item>
              <a-menu-item @click.prevent="toAddRobot(1)">
                <span class="create-action">
                  <img class="icon" :src="DEFAULT_WORKFLOW_AVATAR" alt="">
                  <span>工作流</span>
                </span>
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
        
      </div>
    </div>

    <div class="list-box">
      <div
        class="list-item-wrapper"
        v-for="item in filterRobotList"
        :key="item.id"
        @click="toEditRobot(item)"
      >
        <div class="list-item">
          <div class="robot-info">
            <img class="robot-avatar" :src="item.robot_avatar" alt="" />
            <div class="robot-info-content">
              <div class="robot-name">{{ item.robot_name }}</div>
              <div class="robot-type-tag">
              {{ item.application_type == 0 ? '聊天机器人' : '工作流' }}
            </div>
            </div>
          </div>
          <div class="robot-desc">{{ item.robot_intro }}</div>
          <div class="robot-action" @click.stop>
            <!-- <div class="robot-action-item" @click="toEditRobot(item)"><svg-icon class="action-icon" name="jibenpeizhi" /></div> -->
            <!-- <div class="robot-action-item" @click="toTestPage(item)"><svg-icon class="action-icon" name="cmd" /></div> -->
            <a-tooltip title="聊天测试">
              <div class="robot-action-item" @click.stop="toTestPage(item)">
                <svg-icon class="action-icon" name="chat"></svg-icon>
              </div>
            </a-tooltip>
            <a-tooltip title="会话记录">
              <div class="robot-action-item" @click.stop="toSessionRecordPage(item)">
                <svg-icon class="action-icon" name="session"></svg-icon>
              </div>
            </a-tooltip>
            <a-tooltip title="统计分析">
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
                    <a class="delete-text-color" href="javascript:;" @click="handleDelete(item)"
                      >删 除</a
                    >
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
            
          </div>
        </div>
      </div>
    </div>

    <AddRobotAlert ref="addRobotAlertRef" />
  </div>
</template>

<script setup>
import { DEFAULT_ROBOT_AVATAR, DEFAULT_WORKFLOW_AVATAR} from '@/constants/index.js'
import { usePermissionStore } from '@/stores/modules/permission'
import { getRobotList, deleteRobot } from '@/api/robot/index.js'
import { ref, onMounted, createVNode, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ExclamationCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import AddRobotAlert from './components/add-robot-alert.vue'
import PageAlert from '@/components/page-alert/page-alert.vue'
import ListTabs from '@/components/cu-tabs/list-tabs.vue'

const tabs = ref([
  {
    title: '全部 (0)',
    value: '2'
  },{
    title: '聊天机器人 (0)',
    value: '0'
  },
  {
    title: '工作流 (0)',
    value: '1'
  }
])

const permissionStore = usePermissionStore()
let { role_permission } = permissionStore

const robotCreate = computed(() => role_permission.includes('RobotCreate'))

const router = useRouter()

const activeKey = ref('2')

const robotList = ref([])

const filterRobotList = computed(()=>{
  if(activeKey.value == 2){
    return robotList.value
  }
  return robotList.value.filter(item => item.application_type == activeKey.value)
})

const getList = () => {
  getRobotList()
    .then((res) => {
      let allNumber = 0;
      let chatNumber = 0;
      let workflowNumber = 0;

      res.data.forEach(item => {
        if(item.application_type == 0){
          chatNumber++;
        }else{
          workflowNumber++;
        }
        allNumber++;
      });

      tabs.value = [
        {
          title: '全部 ('+allNumber+')',
          value: '2' 
        },{
          title: '聊天机器人 ('+chatNumber+')',
          value: '0'
        },{
          title: '工作流 ('+workflowNumber+')',
          value: '1'
         },
      ]
      robotList.value = res.data
    })
    .catch(() => {})
}

const addRobotAlertRef = ref(null)
const toAddRobot = (val) => {
  // router.push({ name: 'addRobot' })
  addRobotAlertRef.value.open(val)
}

const toEditRobot = ({id, robot_key, application_type}) => {
  if(application_type == 0){
    router.push({ path: '/robot/config/basic-config', query: { id: id, robot_key: robot_key  } })
  }else{
    router.push({ path: '/robot/config/workflow', query: { id: id, robot_key: robot_key  } })
  }
}

const handleDelete = (data) => {
  console.log(data, '===')
  let secondsToGo = 3
  let modal = Modal.confirm({
    title: `删除机器人${data.robot_name}`,
    icon: createVNode(ExclamationCircleOutlined),
    content: '您确定要删除此机器人吗？',
    okText: secondsToGo + ' 确 定',
    okType: 'danger',
    cancelText: '取 消',
    okButtonProps: {
      disabled: true
    },
    onOk() {
      onDelete(data)
    },
    onCancel() {
      // console.log('Cancel')
    }
  })

  let interval = setInterval(() => {
    if (secondsToGo == 1) {
      modal.update({
        okText: '确 定',
        okButtonProps: {
          disabled: false
        }
      })

      clearInterval(interval)
      interval = undefined
    } else {
      secondsToGo -= 1

      modal.update({
        okText: secondsToGo + ' 确 定',
        okButtonProps: {
          disabled: true
        }
      })
    }
  }, 1000)
}

const onDelete = ({ id }) => {
  deleteRobot({ id }).then(() => {
    message.success('删除成功')
    getList()
  })
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

onMounted(() => {
  getList()
})
</script>
