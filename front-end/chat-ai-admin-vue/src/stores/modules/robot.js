import { reactive, ref } from 'vue'
import { defineStore } from 'pinia'
import { getRobotInfo } from '@/api/robot/index'

// WebApp配置
const external_config_h5_default = {
  accessRestrictionsType: 1, // 访问限制
  logo: '',
  pageTitle: '',
  navbarShow: 2,
  lang: 'zh-CN',
  pageStyle: {
    navbarBackgroundColor: '#2475FC'
  }
}
// 嵌入网站配置
const external_config_pc_default = {
  headTitle: '',
  headSubTitle: 'Based on LLM, free and open-source.',
  headImage: '',
  lang: 'zh-CN',
  pageStyle: {
    headBackgroundColor: 'linear-gradient,to right,#2435E7,#01A0FB'
  },
  floatBtn: {
    displayType: 1,
    buttonText: '快来聊聊吧~',
    buttonIcon: '',
    bottomMargin: 32,
    rightMargin: 32,
    showUnreadCount: 1,
    showNewMessageTip: 1
  }
}

export const useRobotStore = defineStore('robot', () => {
  const robotInfo = reactive({
    id: undefined,
    robot_key: '',
    fast_command_switch: '', // webapp存储的开关
    yunpc_fast_command_switch: '', // 嵌入网站存储的开关
    app_id: -1, // webapp:-1
    app_id_embed: -2, // 嵌入网站:-2
    robot_name: '',
    robot_intro: '',
    robot_avatar: undefined,
    robot_avatar_url: '',
    use_model: undefined,
    model_config_id: '',
    optimize_question_dialogue_background: '',
    optimize_question_model_config_id: '',
    optimize_question_use_model: '',
    rerank_status: 0,
    rerank_use_model: undefined,
    rerank_model_config_id: undefined,
    temperature: 0,
    max_token: 0,
    context_pair: 0,
    top_k: 1,
    similarity: 0,
    search_type: 1,
    prompt: '',
    prompt_struct: {},
    prompt_struct_default: {},
    prompt_type: 1,
    library_ids: '',
    form_ids: '',
    work_flow_ids: '',
    welcomes: {
      content: '',
      question: []
    },
    h5_domain: '',
    h5_website: '',
    pc_domain: '',
    chat_type: '',
    library_qa_direct_reply_score: '',
    library_qa_direct_reply_switch: '',
    mixture_qa_direct_reply_score: '',
    mixture_qa_direct_reply_switch: '',
    unknown_question_prompt: {
      content: '',
      question: []
    },
    show_type: 1,
    enable_question_optimize: 'false',
    enable_question_guide: 'false',
    enable_common_question: 'false',
    common_question_list: [],
    answer_source_switch: 'false',
    feedback_switch: 1,
    application_type: 0,
    sensitive_words_switch: 0,
    question_guide_num: '',
    start_node_key: '',
  })

  // WebApp配置
  const external_config_h5 = reactive(JSON.parse(JSON.stringify(external_config_h5_default)))
  // 嵌入网站配置
  const external_config_pc = reactive(JSON.parse(JSON.stringify(external_config_pc_default)))

  const setRobotInfo = (data) => {
    let welcomes = JSON.parse(data.welcomes)

    welcomes.question = welcomes.question.map((content) => {
      return {
        content: content
      }
    })

    let unknown_question_prompt = {
      content: '',
      question: []
    }
    if (data.unknown_question_prompt) {
      unknown_question_prompt = JSON.parse(data.unknown_question_prompt)
    }

    unknown_question_prompt.question = unknown_question_prompt.question.map((content) => {
      return {
        content: content
      }
    })

    let prompt_struct = {}
    if (data.prompt_struct) {
      prompt_struct = JSON.parse(data.prompt_struct)
    }

    let prompt_struct_default = {}
    if (data.prompt_struct_default) {
      prompt_struct_default = JSON.parse(data.prompt_struct_default)
    }

    robotInfo.id = data.id
    robotInfo.robot_key = data.robot_key
    robotInfo.fast_command_switch = data.fast_command_switch
    robotInfo.yunpc_fast_command_switch = data.yunpc_fast_command_switch
    robotInfo.robot_name = data.robot_name
    robotInfo.robot_intro = data.robot_intro
    robotInfo.robot_avatar_url = data.robot_avatar
    robotInfo.use_model = data.use_model
    robotInfo.model_config_id = data.model_config_id
    robotInfo.optimize_question_dialogue_background = data.optimize_question_dialogue_background
    robotInfo.optimize_question_model_config_id = data.optimize_question_model_config_id
    robotInfo.optimize_question_use_model = data.optimize_question_use_model
    robotInfo.rerank_status = Number(data.rerank_status)
    robotInfo.rerank_use_model = data.rerank_use_model
    robotInfo.rerank_model_config_id = data.rerank_model_config_id
    robotInfo.prompt = data.prompt
    robotInfo.prompt_struct = prompt_struct
    robotInfo.prompt_struct_default = prompt_struct_default
    robotInfo.prompt_type = +data.prompt_type
    robotInfo.library_ids = data.library_ids || ''
    robotInfo.form_ids = data.form_ids || ''
    robotInfo.work_flow_ids = data.work_flow_ids || ''
    robotInfo.welcomes = welcomes
    robotInfo.unknown_question_prompt = unknown_question_prompt

    robotInfo.show_type = Number(data.show_type)
    robotInfo.temperature = Number(data.temperature)
    robotInfo.max_token = Number(data.max_token)
    robotInfo.context_pair = Number(data.context_pair)
    robotInfo.top_k = Number(data.top_k)
    robotInfo.similarity = Number(data.similarity)
    robotInfo.search_type = Number(data.search_type)
    robotInfo.think_switch = Number(data.think_switch)
    robotInfo.pc_domain = data.pc_domain
    robotInfo.chat_type = data.chat_type
    robotInfo.library_qa_direct_reply_score = data.library_qa_direct_reply_score
    robotInfo.library_qa_direct_reply_switch = data.library_qa_direct_reply_switch
    robotInfo.mixture_qa_direct_reply_score = data.mixture_qa_direct_reply_score
    robotInfo.mixture_qa_direct_reply_switch = data.mixture_qa_direct_reply_switch
    robotInfo.h5_domain = data.h5_domain
    robotInfo.h5_website = data.h5_domain + '/#/chat?robot_key=' + data.robot_key
    robotInfo.enable_question_optimize = data.enable_question_optimize
    robotInfo.enable_question_guide = data.enable_question_guide
    robotInfo.enable_common_question = data.enable_common_question
    robotInfo.common_question_list = data.common_question_list
    robotInfo.answer_source_switch = data.answer_source_switch
    robotInfo.feedback_switch = +data.feedback_switch 
    robotInfo.application_type = data.application_type
    robotInfo.sensitive_words_switch = +data.sensitive_words_switch
    robotInfo.question_guide_num = +data.question_guide_num

    robotInfo.start_node_key = data.start_node_key
    // h5配置
    if (data.external_config_h5 !== '') {
      Object.assign(external_config_h5, JSON.parse(data.external_config_h5))
    } else {
      Object.assign(external_config_h5, JSON.parse(JSON.stringify(external_config_h5_default)))
      external_config_h5.pageTitle = robotInfo.robot_name
      external_config_h5.logo = robotInfo.robot_avatar_url
    }

    // 嵌入网站配置
    if (data.external_config_pc !== '') {
      Object.assign(external_config_pc, JSON.parse(data.external_config_pc))
    } else {
      external_config_pc.headTitle = robotInfo.robot_name
      external_config_pc.headImage = robotInfo.robot_avatar_url
    }
  }

  const getRobot = async (id) => {
    const res = await getRobotInfo({ id })

    if (res.code !== 0) {
      return false
    }

    // 特殊处理，后端返回了library_ids: ",357,358,359,244" 处理成 正常的字符串数据
    if (res.data.library_ids) {
      res.data.library_ids = res.data.library_ids.split(',').filter(Boolean).join(',');
    }

    let data = res.data

    setRobotInfo(data)

    return res
  }
  const quickCommandLists = ref([])
  const setQuickCommand = async (data) => {
    quickCommandLists.value = data || []
  }

  const draftSaveTime = ref({})
  const setDrafSaveTime = (data) => {
    draftSaveTime.value = data
  }

  const modelList = ref([])

  const setModelList = (data) => {
    modelList.value = data
  }

  return {
    robotInfo,
    getRobot,
    external_config_h5,
    external_config_pc,
    quickCommandLists,
    setQuickCommand,
    draftSaveTime,
    setDrafSaveTime,
    modelList,
    setModelList
  }
})
