import { defineStore } from 'pinia'
import { getRobotList } from '@/api/robot'
import { getUuid, formatDisplayChatTime } from '@/utils/index'
import { useIM } from '@/hooks/event/useIM'
import { useUserStore } from '@/stores/modules/user'
import { useEventBus } from '@/hooks/event/useEventBus'
import { getReceiverList, getChatMessage, setReceiverRead } from '@/api/chat-monitor'

export const useChatMonitorStore = defineStore('chatMonitor', {
  state: () => {
    return {
      im: null,
      selectedRobotId: '',
      robotList: [],
      selectedReceiverId: '',
      receiverList: [],
      receiverListPage: {
        page: 0,
        size: 10
      },
      activeChat: null,
      chatMessagePageSize: 10, // 消息列表分页大小
      messageList: [],
      messageListScrollTop: 0,
      chatMessageLoadCompleted: false, // 消息加载完了
      messageListLoading: false // 消息列表加载中
    }
  },
  getters: {},
  actions: {
    async init() {
      this.selectedRobotId = ''
      this.selectedReceiverId = ''
      this.robotList = []
      this.receiverList = []
      this.messageList = []
      this.receiverListPage = {
        page: 0,
        size: 10
      }
      this.messageListScrollTop = 0
      this.chatMessageLoadCompleted = false
      this.messageListLoading = false
      this.activeChat = null
      await this.getRobotList()
      await this.getReceiverList()

      if (this.receiverList.length > 0) {
        await this.switchChat(this.receiverList[0])
      }

      this.initIM()
    },
    closeIM() {
      if (this.im) {
        this.im.close()
      }
    },
    initIM() {
      const im = useIM()
      const user = useUserStore()

      this.im = im

      im.connect(user.user_id)
      im.on('message', (res) => {
        let msg = res.data

        if (!msg || res.msg_type !== 'receiver_notify') {
          return
        }

        if (res.change == 'create' || res.change == 'update') {
          this.onAddReceiver(res)
        }

        if (res.change == 'delete') {
          this.onDeleteReceiver(res)
        }

        if (res.change == 'c_message' || res.change == 'ai_message') {
          this.onAddMessage(res)
        }
      })
    },
    onAddMessage(res) {
      if(!this.activeChat){
        return
      }

      const emitter = useEventBus()

      let msg = res.data

      if (this.activeChat.session_id != msg.session_id) {
        return
      }

      msg.displayName = msg.name || msg.nickname
      msg.dispayTime = formatDisplayChatTime(msg.create_time)

      msg.uid = getUuid(32)
      msg.loading = false
      // msg.robot_avatar = robot.robot_avatar

      if (msg.menu_json && typeof msg.menu_json === 'string') {
        msg.menu_json = JSON.parse(msg.menu_json)
      }

      if (msg.quote_file && typeof msg.quote_file === 'string') {
        msg.quote_file = JSON.parse(msg.quote_file)
      }
      if (msg.reply_content_list && typeof msg.reply_content_list === 'string') {
        try { msg.reply_content_list = JSON.parse(msg.reply_content_list) } catch (_) { msg.reply_content_list = [] }
      }
      this.messageList.push(msg)

      emitter.emit('onAddMessage', msg)
    },
    onAddReceiver(res) {
      let data = res.data

      data.displayName = data.name || data.nickname
      data.come_from = JSON.parse(data.come_from)

      let item = this.receiverList.find((item) => item.session_id == data.session_id)

      if (item) {
        Object.assign(item, data)
      }else{
        if(this.selectedRobotId == '' || this.selectedRobotId == data.robot_id) {
          this.receiverList.unshift(data)
        }
      }
    },
    onDeleteReceiver(res) {
      let data = res.data
      let index = this.receiverList.findIndex((item) => item.id == data)

      if (index == -1) {
        return
      }

      if (this.activeChat && this.activeChat.id == data) {
        this.activeChat = null
      }

      this.receiverList.splice(index, 1)
    },
    async getRobotList(params) {
      return getRobotList(params).then((res) => {
        this.robotList = res.data || []
        return res
      })
    },
    async getReceiverList (params) {
      if (params?.page) {
        this.receiverListPage.page = params?.page
      } else {
        this.receiverListPage.page++
      }

      return getReceiverList({ ...this.receiverListPage, robot_id: this.selectedRobotId, ...params }).then(
        (res) => {
          let list = res.data.list || []

          for (let i = 0; i < list.length; i++) {
            list[i].displayName = list[i].name || list[i].nickname

            if (list[i].come_from) {
              list[i].come_from = JSON.parse(list[i].come_from)
            }
          }


          if (this.receiverListPage.page == 1) {
            this.receiverList = list
          } else {
            this.receiverList = this.receiverList.concat(list)
          }

          return res
        }
      )
    },
    changeRobot(params) {
      this.activeChat = null
      this.resetReceiverList(params)
    },
    resetReceiverList(params) {
      this.receiverListPage.page = 0
      this.receiverList = []

      return this.getReceiverList(params)
    },
    async getChatMessage() {
      if (this.messageListLoading) {
        return
      }

      this.messageListLoading = true

      let min_id = 0
      let list = this.messageList.filter((item) => !item.isWelcome)

      if (list.length > 0) {
        min_id = list[0].id
      }

      let params = {
        robot_key: this.activeChat.robot_key,
        openid: this.activeChat.openid,
        min_id: min_id,
        size: this.chatMessagePageSize,
        rel_user_id: this.activeChat.rel_user_id
      }

      return getChatMessage(params)
        .then((res) => {
          this.messageListLoading = false

          let list = res.data.list || []
          const _customer = res?.data?.customer || {}
          const _robot = res?.data?.robot || {}

          list.sort((a, b) => {
            return a.id - b.id
          })

          for (let i = 0; i < list.length; i++) {
            list[i].displayName = list[i].name || list[i].nickname
            if (list[i].is_customer == 1) {
              list[i].displayName = list[i].displayName || _customer.name
              list[i].avatar = list[i].avatar || _customer.avatar
            } else {
              list[i].displayName = list[i].displayName || _robot.robot_name
              list[i].avatar = list[i].avatar || _robot.robot_avatar
            }
            list[i].uid = getUuid(32)

            if (list[i].menu_json && typeof list[i].menu_json === 'string') {
              list[i].menu_json = JSON.parse(list[i].menu_json)
            }


            if (list[i].quote_file && typeof list[i].quote_file === 'string') {
              list[i].quote_file = JSON.parse(list[i].quote_file)
            }

            if (list[i].reply_content_list && typeof list[i].reply_content_list === 'string') {
              try { list[i].reply_content_list = JSON.parse(list[i].reply_content_list) } catch (_) { list[i].reply_content_list = [] }
            }

            list[i].dispayTime = formatDisplayChatTime(list[i].create_time)
          }

          this.messageList = [...list, ...this.messageList]

          // 消息加载完了
          if (list.length < this.chatMessagePageSize) {
            this.chatMessageLoadCompleted = true
          }

          return res
        })
        .catch((err) => {
          this.messageListLoading = false
          console.log(err)
        })
    },
    // 清除红点
    claerUnread() {
      return setReceiverRead({ id: this.selectedReceiverId })
    },
    async switchChat(receiver) {
      this.chatMessageLoadCompleted = false
      this.messageListLoading = false
      this.messageList = []
      this.messageListScrollTop = 0
      this.selectedReceiverId = receiver.id
      this.activeChat = receiver

      this.claerUnread()

      await this.getChatMessage()
    }
  }
})
