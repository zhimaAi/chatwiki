import { reactive, ref } from 'vue'
import { defineStore } from 'pinia'
import {
  searchLirary,
  getLibrarySearch
} from '@/api/library'
import { getUuid, devLog } from '@/utils/index'
import { useEventBus } from '@/hooks/event/useEventBus'

export const useSearchStore = defineStore('searchStore', {
  state: () => ({
    indeterminate: false,  // 半选
    checkAll: false,  // 是否全选
    checkedList: [],  // 知识搜索选中的模型
    activeKey: 'all', // 知识搜索选中的模型类型 all：全部 0：普通 2：问答
  }),
  getters: {
    getIndeterminate () {
      return this.indeterminate
    },
    getCheckAll () {
      return this.checkAll
    },
    getCheckedList () {
      return this.checkedList
    },
    getActiveKey () {
      return this.activeKey
    }
  },
  actions: {
    setIndeterminate (val) {
      this.indeterminate = val
    },
    setCheckAll (val) {
      this.checkAll = val
    },
    setCheckedList (val) {
      this.checkedList = val
    },
    setActiveKey (val) {
      this.activeKey = val
    }
  },
  persist: {
    enabled: true,
    strategies: [
      {
        key: 'searchStore', // 存储的键名
        storage: localStorage, // 存储位置
        paths: ['indeterminate', 'checkAll', 'checkedList', 'activeKey'] // 可以选择只持久化部分数据
      }
    ]
  }
})

export const useSearchLiraryStore = defineStore('search-lirary', () => {
  const emitter = useEventBus()
  const messageObj = reactive({
    reasoning_content: '',
    show_reasoning: false,
    content: '',
    quote_file: [],
    id: '',
    finish: false,
    debug: [],
    error: [],
    recall_time: 0,
    request_time: 0,
    loading: false,
  })

  let mySSE = null

  // 对话id
  const dialogue_id = ref(0)

  const openid = ref('')

  // 搜索配置的信息
  const searchFormState = reactive({
    context_pair: "",
    create_time: "",
    id: "",
    max_token: "",
    model_config_id: "",
    rerank_model_config_id: "",
    rerank_status: "",
    rerank_use_model: "",
    search_type: "",
    similarity: "",
    size: "",
    temperature: "",
    update_time: "",
    use_model: "",
    user_id: ""
  })

  // 更新AI的消息到列表
  const updateAiMessage = (type, content) => {
    if (type == 'reasoning_content') {
      let oldText = messageObj.reasoning_content || ''
      messageObj.reasoning_content = oldText + content

      // 推理开始
      // messageObj.reasoning_status = true
      messageObj.show_reasoning = true
    }

    if (type == 'sending') {
      // 推理结束
      // messageObj.reasoning_status = false

      let oldText = messageObj.content || ''
      // console.log('sending', content.length, content)
      messageObj.content = oldText + content
    }

    if (type == 'quote_file') {
      messageObj.quote_file = content.length > 0 ? content : []
    }

    if (type == 'ai_message') {
      // 【ID1017779】【芝麻AI】机器人直连
      // 在实时聊天的时候，也需要把 ai_message 中的菜单内容显示出来
      if (content.menu_json && content.msg_type == 2) {
        let menu_json_obj = JSON.parse(content.menu_json)
        //messageObj.content = menu_json_obj.content
        messageObj.question = menu_json_obj.question
      }

      messageObj.id = content.id

      messageObj.content = content.content
      if (content.quote_file && typeof content.quote_file === 'string') {
        messageObj.quote_file = JSON.parse(content.quote_file)
      }
    }

    // 回答完毕
    if (type == 'finish') {
      messageObj.finish = true
    }

    if (type == 'debug') {
      messageObj.debug = content.length > 0 ? content : []
    }

    if (type == 'error') {
      messageObj.error = content.length > 0 ? content : []
    }
    if (type == 'recall_time') {
      messageObj.recall_time = content || 0
    }

    if (type == 'request_time') {
      messageObj.request_time = content || 0
    }

    emitter.emit('updateAiMessage', messageObj)
  }

  // 关闭AI的消息加载状态
  const closeAiMessageLoading = () => {
    messageObj.loading = false
  }

  // 搜索
  const searchMessage = (keyword, library_ids, searchConfigData) => {
    let aiMsg = {
      uid: getUuid(32)
    }

    // 清空之前的内容
    messageObj.content = ''
    messageObj.finish = false

    let params = {
      model_config_id: searchConfigData.model_config_id,
      use_model: searchConfigData.use_model,
      rerank_status: searchConfigData.rerank_status,
      temperature: searchConfigData.temperature,
      max_token: searchConfigData.max_token,
      size: searchConfigData.size,
      similarity: searchConfigData.similarity,
      search_type: searchConfigData.search_type,
      question: keyword,
      id: library_ids.join(','),
      rerank_use_model: searchConfigData.rerank_use_model,
      recall_type: 1
    }

    if (searchConfigData.rerank_status == 1) {
      params.rerank_model_config_id = searchConfigData.rerank_model_config_id
    }
    // if (import.meta.env.DEV) {
    //   params.debug = 0
    // }

    mySSE = searchLirary(params)

    mySSE.onMessage = (res) => {
      if (res.event !== 'sending') {
        devLog(res)
      }
      // 更新对话id
      if (res.event == 'dialogue_id') {
        dialogue_id.value = res.data
      }

      // 更新机器人深度思考的内容
      if (res.event == 'reasoning_content') {
        updateAiMessage('reasoning_content', res.data, aiMsg.uid)
      }

      // 更新机器人的消息
      if (res.event == 'sending') {
        updateAiMessage('sending', res.data, aiMsg.uid)
      }

      // 更新机器人消息的消息id时间等
      if (res.event == 'ai_message') {
        let data = JSON.parse(res.data)

        updateAiMessage('ai_message', data, aiMsg.uid)
      }

      // 更新引用文件
      if (res.event == 'quote_file') {
        let data = JSON.parse(res.data)

        updateAiMessage('quote_file', data, aiMsg.uid)
      }

      // 更新机器人的消息
      if (res.event == 'finish') {
        updateAiMessage('finish', res.data, aiMsg.uid)
      }

      // 更新prompt日志
      if (res.event == 'debug') {
        let data = JSON.parse(res.data)

        updateAiMessage('debug', data, aiMsg.uid)
      }

      // 更新prompt错误日志
      if (res.event == 'error') {
        let data = res.data

        updateAiMessage('error', data, aiMsg.uid)
      }
      // 更新prompt recall_time
      if (res.event == 'recall_time') {
        let data = res.data

        updateAiMessage('recall_time', data, aiMsg.uid)
      }
      // 更新prompt request_time
      // if (res.event == 'request_time') {
      //   let data = res.data

      //   updateAiMessage('request_time', data, aiMsg.uid)
      // }
    }

    mySSE.onClose = () => {
      closeAiMessageLoading()

      mySSE = null
    }
  }

  // 获取搜索配置
  const getLibrarySearchFn = async () => {
    if (mySSE) {
      mySSE.abort()
      mySSE = null
    }

    const res = await getLibrarySearch()

    // "data": {
    //   "context_pair": "6",
    //   "create_time": "1744191057",
    //   "id": "2",
    //   "max_token": "2000",
    //   "model_config_id": "11",
    //   "rerank_model_config_id": "8",
    //   "rerank_status": "1",
    //   "rerank_use_model": "BAAI/bge-reranker-v2-m3",
    //   "search_type": "1",
    //   "similarity": "0.6",
    //   "size": "10",
    //   "temperature": "0.5",
    //   "update_time": "1744616719",
    //   "use_model": "qwen-max",
    //   "user_id": "1"
    // }

    try {
      let _data = res.data

      Object.assign(searchFormState, { ..._data });
    } catch (e) {
      Promise.reject(e)
    }
    return res
  }

  return {
    searchFormState,
    dialogue_id,
    openid,
    messageObj,
    getLibrarySearchFn,
    searchMessage
  }
})
