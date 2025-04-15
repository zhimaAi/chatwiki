<template>
  <div class="robot-page">
    <cu-scroll>
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
              </div>
            </div>
            <div class="robot-desc">{{ item.robot_intro }}</div>
            <div class="robot-type-tag">
              {{ item.application_type == 0 ? '聊天机器人' : '工作流' }}
            </div>
          </div>
        </div>
      </div>
      <div class="empty-box" v-if="lists.length == 0 && !isLoading">
        <img src="@/assets/img/library/detail/empty.png" alt="" />
        <div class="title">暂无数据</div>
      </div>
    </cu-scroll>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { getLibraryRobotInfo } from '@/api/library'
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
    })
    .finally(() => {
      isLoading.value = false
    })
}

getLists()

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
  height: 165px;
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

  .robot-type-tag {
    width: fit-content;
    display: flex;
    align-items: center;
    padding: 0 8px;
    color: #2475fc;
    font-size: 14px;
    font-weight: 400;
    border-radius: 6px;
    border: 1px solid #99bffd;
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
