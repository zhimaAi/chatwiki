export const defaultRadioOpt = [
  {label: '关闭', value: false},
  {label: '开启', value: true},
]

export const defaultVoiceOpts = [
  { label: '男生', value: 'male-qn-jingying' },
  { label: '女生', value: 'female-chengshu' },
  { label: '男童', value: 'clever_boy' },
  { label: '女童', value: 'lovely_girl' },
]

export const emotionMap = {
  happy: '高兴',
  sad: '悲伤',
  angry: '愤怒',
  fearful: '害怕',
  disgusted: '厌恶',
  surprised: '惊讶',
  calm: '中性',
  fluent: '生动',
  whisper: '低语'
};


export const audioSampleRates = [8000, 16000, 22050, 24000, 32000, 44100];

export const audioBitrate = [32000, 64000, 128000, 256000];

export const audioFormats = ['mp3', 'pcm', 'flac', 'wav'];

export const audioChannels = [
  {label: '单声道', value: 1},
  {label: '双声道', value: 2},
]

export const formatMap = {
  url: 'url',
  hex: 'hex'
};

export const effectMap = {
  spacious_echo: '空旷回音',
  auditorium_echo: '礼堂广播',
  lofi_telephone: '电话失真',
  robotic: '电音'
};

export const languageMap = {
  auto: '自动检测',
  Chinese: '中文',
  Yue: '粤语',
  English: '英语',
  Arabic: '阿拉伯语',
  Russian: '俄语',
  Spanish: '西班牙语',
  French: '法语',
  Portuguese: '葡萄牙语',
  German: '德语',
  Turkish: '土耳其语',
  Dutch: '荷兰语',
  Ukrainian: '乌克兰语',
  Vietnamese: '越南语',
  Indonesian: '印度尼西亚语',
  Japanese: '日语',
  Italian: '意大利语',
  Korean: '韩语',
  Thai: '泰语',
  Polish: '波兰语',
  Romanian: '罗马尼亚语',
  Greek: '希腊语',
  Czech: '捷克语',
  Finnish: '芬兰语',
  Hindi: '印地语',
  Bulgarian: '保加利亚语',
  Danish: '丹麦语',
  Hebrew: '希伯来语',
  Malay: '马来语',
  Persian: '波斯语',
  Slovak: '斯洛伐克语',
  Swedish: '瑞典语',
  Croatian: '克罗地亚语',
  Filipino: '菲律宾语',
  Hungarian: '匈牙利语',
  Norwegian: '挪威语',
  Slovenian: '斯洛文尼亚语',
  Catalan: '加泰罗尼亚语',
  Nynorsk: '新挪威语',
  Tamil: '泰米尔语',
  Afrikaans: '南非荷兰语',
};

export const voiceSettingSliderConfig = {
  'speed': {
    ui_type: 'slider',
    key: 'voice_setting.speed',
    label: '语速',
    type: 'number',
    required: true,
    tooltip: '合成音频的语速，取值范围 [0.5, 2]，默认值为 1.0',
    default: 1.0,
    min: 0.5,
    max: 2,
    step: 0.1,
    marks: {0.5: '0.5', 2: '2'}
  },
  'vol': {
    ui_type: 'slider',
    key: 'voice_setting.vol',
    label: '音量',
    type: 'number',
    required: true,
    tooltip: '合成音频的音量，取值越大，音量越高。取值范围 (0,10]，默认值为 1',
    default: 1,
    min: 0,
    max: 10,
    step: 1,
    marks: {0: '0', 10: '10'}
  },
  'pitch': {
    ui_type: 'slider',
    key: 'voice_setting.pitch',
    label: '音高',
    type: 'number',
    required: true,
    tooltip: '合成音频的语调，取值范围 [-12,12]，默认值为 0，其中 0 为原音色输出',
    default: 0,
    min: -12,
    max: 12,
    step: 1,
    marks: {'-12': '-12', 12: '12'}
  }
}


export const voiceModifySliderConfig = {
  'pitch': {
    ui_type: 'slider',
    key: 'voice_modify.pitch',
    label: '音高',
    type: 'number',
    tooltip: '音高调整（低沉/明亮），范围 [-100,100]，数值接近 -100，声音更低沉；接近 100，声音更明亮',
    default: 0,
    min: -100,
    max: 100,
    step: 1,
    marks: {'-100': '-100', 100: '100'}
  },

  'intensity': {
    ui_type: 'slider',
    key: 'voice_modify.intensity',
    label: '强度',
    type: 'number',
    tooltip: '强度调整（力量感/柔和），范围 [-100,100]，数值接近 -100，声音更刚劲；接近 100，声音更轻柔',
    default: 0,
    min: -100,
    max: 100,
    step: 1,
    marks: {'-100': '-100', 100: '100'}
  },

  'timbre': {
    ui_type: 'slider',
    key: 'voice_modify.timbre',
    label: '音色',
    type: 'number',
    tooltip: '音色调整（磁性/清脆），范围 [-100,100]，数值接近 -100，声音更浑厚；数值接近 100，声音更清脆',
    default: 0,
    min: -100,
    max: 100,
    step: 1,
    marks: {'-100': '-100', 100: '100'}
  },
}

export const voiceOutputObj = {
  "data": {
    "desc": "音频合成结果",
    "type": "object",
    "properties": {
      "audio": {
        "desc": "音频数据（hex 编码）",
        "type": "string"
      },
      "status": {
        "desc": "音频合成状态",
        "type": "number"
      }
    }
  },
  "extra_info": {
    "desc": "音频附加信息",
    "type": "object",
    "properties": {
      "audio_length": {
        "desc": "音频时长（毫秒）",
        "type": "number"
      },
      "audio_sample_rate": {
        "desc": "音频采样率",
        "type": "number"
      },
      "audio_size": {
        "desc": "音频大小（字节）",
        "type": "number"
      },
      "bitrate": {
        "desc": "音频码率",
        "type": "number"
      },
      "word_count": {
        "desc": "文本字数",
        "type": "number"
      },
      "invisible_character_ratio": {
        "desc": "不可见字符比例",
        "type": "number"
      },
      "usage_characters": {
        "desc": "实际消耗字符数",
        "type": "number"
      },
      "audio_format": {
        "desc": "音频格式",
        "type": "string"
      },
      "audio_channel": {
        "desc": "音频声道数",
        "type": "number"
      }
    }
  },
  "trace_id": {
    "desc": "请求追踪 ID",
    "type": "string"
  },
  "base_resp": {
    "desc": "基础响应信息",
    "type": "object",
    "properties": {
      "status_code": {
        "desc": "状态码，0 表示成功",
        "type": "number"
      },
      "status_msg": {
        "desc": "状态描述",
        "type": "string"
      }
    }
  }
}

export const voiceCloneOutputObj = {
  "input_sensitive": {
    "desc": "输入内容是否命中敏感信息",
    "type": "boolean"
  },
  "input_sensitive_type": {
    "desc": "敏感类型标识（0 表示无敏感）",
    "type": "number"
  },
  "demo_audio": {
    "desc": "示例音频数据",
    "type": "string"
  },
  "base_resp": {
    "desc": "基础响应信息",
    "type": "object",
    "properties": {
      "status_code": {
        "desc": "状态码，0 表示成功",
        "type": "number"
      },
      "status_msg": {
        "desc": "状态描述",
        "type": "string"
      }
    }
  }
}


