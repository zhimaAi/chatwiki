<style lang="less" scoped>
.robot-page {
  .robot-item {
    padding: 8px;
    border-bottom: 1px solid #ccc;
  }

  ::v-deep(.ant-tabs-nav::before){
    border-bottom: 0;
  }

  .top-banner {
    position: relative;
    line-height: 22px;
    padding: 16px;
    margin-bottom: 8px;
    font-size: 14px;
    color: #3a4559;
    border-radius: 2px;
    background-color: #e6efff;

    &::after {
      content: '';
      position: absolute;
      right: 0;
      top: 0;
      width: 552px;
      height: 76px;
      background: url('@/assets/img/robot/robot_top_banner.png') 0 0 no-repeat;
      background-size: cover;
    }
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
    height: 198px;
    padding: 14px 16px;
    border: 1px solid #f0f0f0;
    border-radius: 2px;
    background-color: #fff;
    transition: all 0.25s;

    &:hover {
      box-shadow: 0 4px 16px 0 #1b3a6929;
    }

    .item-action {
      position: absolute;
      top: 16px;
      right: 16px;

      .menu-btn {
        width: 16px;
        height: 16px;
        display: flex;
        justify-content: center;
        align-items: center;
        cursor: pointer;

        &:hover {
          color: #2475fc;
        }
      }
    }

    .robot-info {
      display: flex;
      // align-items: start;
      align-items: center;
    }

    .robot-avatar {
      width: 40px;
      height: 40px;
      border-radius: 2px;
      overflow: hidden;
    }

    .robot-info-content {
      flex: 1;
      padding-left: 12px;
    }

    .robot-name {
      line-height: 22px;
      margin-bottom: 2px;
      font-size: 14px;
      font-weight: 600;
      color: #262626;
    }

    .robot-desc {
      margin-top: 16px;
      height: 44px;
      line-height: 22px;
      font-size: 14px;
      font-weight: 400;
      color: #595959;
      // 超出2行显示省略号
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
    }

    .robot-type-tag{
      width: fit-content;
      display: flex;
      align-items: center;
      padding: 0 8px;
      color: #2475fc;
      font-size: 14px;
      font-weight: 400;
      border-radius: 6px;
      border: 1px solid #99BFFD;
      height: 24px;
    }

    .robot-action {
      margin-top: 16px;
      font-size: 14px;
      height: 24px;
      color: #2475fc;
      display: flex;
      align-items: center;
      gap: 8px;
      .robot-action-item {
        height: 100%;
        color: #595959;
        font-weight: 400;
        display: flex;
        width: fit-content;
        align-items: center;
        gap: 4px;
        border-radius: 6px;
        padding: 0 8px;
        cursor: pointer;
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

// 大于1440px
@media screen and (min-width: 1440px) {
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
    <!-- <div class="top-banner">
      <div>
        1.可以创建多个不同的机器人，不同机器人应用在不同场景中，不同机器人可以关联不同的知识库
      </div>
      <div>
        2.您可以复制链接将机器人提供给您的用户使用。在对外提供服务之前，建议您进行充分测试，并适当调整知识库
      </div>
    </div> -->

    <a-tabs v-model:activeKey="activeKey">
      <a-tab-pane key="2" tab="全部"></a-tab-pane>
      <a-tab-pane key="0" tab="聊天机器人"></a-tab-pane>
      <a-tab-pane key="1" tab="工作流"></a-tab-pane>
    </a-tabs>

    <div class="list-box">
      <div class="list-item-wrapper" v-if="robotCreate">
        <div class="list-item add-robot" @click="toAddRobot">
          <PlusCircleOutlined class="add-icon" />
          <span class="add-text">新增应用</span>
        </div>
      </div>

      <div
        class="list-item-wrapper"
        v-for="item in filterRobotList"
        :key="item.id"
        @click="toEditRobot(item)"
      >
        <div class="list-item">
          <span class="item-action" @click.stop>
            <a-dropdown>
              <span class="menu-btn" @click.stop>
                <MoreOutlined />
              </span>
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
          </span>

          <div class="robot-info">
            <img class="robot-avatar" :src="item.robot_avatar" alt="" />
            <div class="robot-info-content">
              <div class="robot-name">{{ item.robot_name }}</div>
            </div>
          </div>
          <div class="robot-desc">{{ item.robot_intro }}</div>
          <div class="robot-type-tag">
            {{ item.application_type == 0 ? '聊天机器人' : '工作流' }}
          </div>
          <div class="robot-action" @click.stop>
            <div class="robot-action-item" @click="toEditRobot(item)"><svg-icon class="action-icon" name="jibenpeizhi" /> 管理</div>
            <div class="robot-action-item" @click="toTestPage(item)"><svg-icon class="action-icon" name="cmd" /> 测试</div>
          </div>
        </div>
      </div>
    </div>

    <AddRobotAlert ref="addRobotAlertRef" />
  </div>
</template>

<script setup>
import { getRobotList, deleteRobot } from '@/api/robot/index.js'
import { ref, onMounted, createVNode, computed } from 'vue'
import { useRouter } from 'vue-router'
import { PlusCircleOutlined, MoreOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import AddRobotAlert from './components/add-robot-alert.vue'
import { usePermissionStore } from '@/stores/modules/permission'

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
      robotList.value = res.data
    })
    .catch(() => {})
}

const addRobotAlertRef = ref(null)
const toAddRobot = () => {
  // router.push({ name: 'addRobot' })
  addRobotAlertRef.value.open()
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
  // router.push({
  //   path: '/robot/test',
  //   query: {
  //     robot_key: item.robot_key,
  //     id: item.id
  //   }
  // })
}

onMounted(() => {
  getList()
})
</script>
