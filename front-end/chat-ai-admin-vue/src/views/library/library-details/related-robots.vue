<template>
  <div class="robot-page">
    <cu-scroll>
      <div class="associated-robot">
        <div class="associated-robot-title">{{ t('title', { count: lists.length }) }}</div>
        <a-button type="primary" @click="handleAddData">
          <template #icon>
            <PlusOutlined />
          </template>
          {{ t('btn_associate_robot') }}
        </a-button>
      </div>
      <div class="list-box">
        <div
          class="list-item-wrapper"
          v-for="item in lists"
          :key="item.id"
          @click="toEditRobot(item)"
        >
          <div class="list-item">
            <div class="robot-info">
              <img class="robot-avatar" :src="item.robot_avatar" alt="" />
              <div class="robot-info-content">
                <div class="robot-name">{{ item.robot_name }}</div>
                <div class="robot-type-tag">
                  {{ item.application_type == 0 ? t('type_chat_bot') : t('type_workflow') }}
                </div>
              </div>
            </div>
            <div class="robot-desc">{{ item.robot_intro }}</div>
          </div>
        </div>
      </div>
      <div class="empty-box" v-if="lists.length == 0 && !isLoading">
        <img src="@/assets/img/library/detail/empty.png" alt="" />
        <div class="title">{{ t('empty_no_data') }}</div>
      </div>
    </cu-scroll>
    <!-- 新增弹出，选择数据 -->
    <SeeModelAlert
      ref="seeModelAlertRef"
      :currentTitle="t('modal_title')"
      :robotList="robotList"
      @save="onSave"
    />
  </div>
</template>

<script setup>
import {
  PlusOutlined,
} from '@ant-design/icons-vue'
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'
import { getLibraryRobotInfo, relationRobot } from '@/api/library'
import SeeModelAlert from '@/components/see-model-alert/see-model-alert.vue'
import { getRobotList } from '@/api/robot/index.js'
import { message } from 'ant-design-vue'
import { useUserStore } from '@/stores/modules/user'

const { t } = useI18n('views.library.library-details.related-robots')

const userStore = useUserStore()
// 查看模型
const seeModelAlertRef = ref(null)
const robotList = ref([])
const rotue = useRoute()
const query = rotue.query
const lists = ref([])
const isLoading = ref(false)
const getLists = () => {
  isLoading.value = true
  getLibraryRobotInfo({
    id: query.id
  })
    .then((res) => {
      lists.value = res.data || []
      userStore.setRobotNums(lists.value.length);
    })
    .finally(() => {
      isLoading.value = false
    })
}

onMounted(() => {
  getLists()
})

// 适配组件的回显
const record = reactive({
  managed_robot_list: []
})

function formatRecord (array) {
  record.managed_robot_list = []
  for (let i = 0; i < array.length; i++) {
    const element = array[i];
    record.managed_robot_list.push(element)
  }
}

const handleAddData = async() => {
  await getList()

  // if (lists.value.length) {
  //   formatRecord(lists.value)
  // }

  // seeModelAlertRef.value.open('robot', 'edit', record) // 回显勾选
  seeModelAlertRef.value.open('robot')
}

const filteringFn = (datas, arr) => {
  // 提取id
  const listsIds = new Set(arr.map(item => item.id))

  // 过滤工作流机器人
  const array = datas.filter(item => item.application_type === '0')

  // 过滤已经关联的机器人
  return array.filter(item => !listsIds.has(item.id))
}

// 获取机器人列表
const getList = async () => {
  await getRobotList()
    .then((res) => {
      robotList.value = filteringFn(res.data, lists.value)
    })
    .catch(() => {})
}

const onSave = (ids) => {
  let params = {
    library_id: query.id,
    robot_ids: ids.join(',')
  }

  relationRobot(params)
    .then((res) => {
      message.success(t('msg_operation_success'))
      getLists()
    })
    .catch(() => {})
}

const toEditRobot = ({ id, robot_key, application_type }) => {
  if (application_type == 0) {
    window.open(`#/robot/config/basic-config?id=${id}&robot_key=${robot_key}`)
    // router.push({ path: '/robot/config/basic-config', query: { id: id, robot_key: robot_key } })
  } else {
    window.open(`#/robot/config/workflow?id=${id}&robot_key=${robot_key}`)
    // router.push({ path: '/robot/config/workflow', query: { id: id, robot_key: robot_key } })
  }
}
</script>

<style lang="less" scoped>
.robot-page {
  height: 100%;
  overflow: hidden;
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
    margin-bottom: 4px;
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

  .robot-type-tag {
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
      background: #e4e6eb;
    }

    .action-icon {
      font-size: 16px;
    }
  }
}

.empty-box {
  text-align: center;
  height: 100%;
  padding-top: 148px;
  img {
    width: 200px;
    height: 200px;
  }
  .title {
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
    color: #262626;
  }
}

.associated-robot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;

  .associated-robot-title {
    line-height: 24px;
    font-size: 16px;
    font-weight: 600;
    color: #262626;
  }
}

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
