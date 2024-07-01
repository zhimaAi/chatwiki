<style lang="less" scoped>
.chat-page {
  display: flex;
  width: 100vw;
  height: 100vh;

  .page-body {
    flex: 1;
    height: 100%;
  }
}
</style>

<template>
  <div class="chat-page">
    <div class="page-left">
      <RobotSidebar v-model:value="robot.id" :robot-list="robotList" @change="onChangeRobot" />
    </div>
    <div class="page-body">
      <NoRobot v-if="robotList.length == 0" />
      <ChatBox v-else :name="robot.robot_name" :src="robot.chat_link" />
    </div>
  </div>
</template>

<script setup>
import { getRobotList } from '@/api/robot'
import RobotSidebar from './components/robot-sidebar.vue'
import NoRobot from './components/no-robot.vue'
import ChatBox from './components/chat-box.vue'
import { reactive, ref } from 'vue'
import Storage from '@/utils/storage'

const robotList = ref([])

const init = () => {
  getRobotList({
    admin_user_id: 1
  }).then((res) => {
    let current_robot_id = Storage.get('current_robot_id')
    let list = res.data.list || []
    let index = 0

    list.forEach((item, i) => {
      if (item.id == current_robot_id) {
        index = i
      }

      item.avatar = res.data.h5_domain + item.robot_avatar
      item.chat_link = res.data.h5_domain + '/#/chat?robot_key=' + item.robot_key
    })

    robotList.value = list

    if (list.length > 0) {
      onChangeRobot(list[index])
    }
  })
}

const robot = reactive({
  robot_name: '',
  chat_link: '',
  id: ''
})

const onChangeRobot = (data) => {
  robot.chat_link = ''
  robot.id = data.id
  robot.robot_name = data.robot_name

  Storage.set('current_robot_id', data.id)

  setTimeout(() => {
    robot.chat_link = data.chat_link
  }, 20)
}

init()
</script>
