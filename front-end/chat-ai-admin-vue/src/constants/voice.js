import { useI18n } from '@/hooks/web/useI18n'

export const getDefaultRadioOpt = () => {
  const { t } = useI18n('constants.voice')
  return [
    {label: t('label_close'), value: false},
    {label: t('label_open'), value: true},
  ]
}

export const defaultRadioOpt = getDefaultRadioOpt()

export const getDefaultVoiceOpts = () => {
  const { t } = useI18n('constants.voice')
  return [
    { label: t('voice_male'), value: 'male-qn-jingying' },
    { label: t('voice_female'), value: 'female-chengshu' },
    { label: t('voice_boy'), value: 'clever_boy' },
    { label: t('voice_girl'), value: 'lovely_girl' },
  ]
}

export const defaultVoiceOpts = getDefaultVoiceOpts()

export const getEmotionMap = () => {
  const { t } = useI18n('constants.voice')
  return {
    happy: t('emotion_happy'),
    sad: t('emotion_sad'),
    angry: t('emotion_angry'),
    fearful: t('emotion_fearful'),
    disgusted: t('emotion_disgusted'),
    surprised: t('emotion_surprised'),
    calm: t('emotion_calm'),
    fluent: t('emotion_fluent'),
    whisper: t('emotion_whisper')
  }
}

export const emotionMap = getEmotionMap()


export const audioSampleRates = [8000, 16000, 22050, 24000, 32000, 44100];

export const audioBitrate = [32000, 64000, 128000, 256000];

export const audioFormats = ['mp3', 'pcm', 'flac', 'wav'];

export const getAudioChannels = () => {
  const { t } = useI18n('constants.voice')
  return [
    {label: t('channel_mono'), value: 1},
    {label: t('channel_stereo'), value: 2},
  ]
}

export const audioChannels = getAudioChannels()

export const formatMap = {
  url: 'url',
  hex: 'hex'
};

export const getEffectMap = () => {
  const { t } = useI18n('constants.voice')
  return {
    spacious_echo: t('effect_spacious_echo'),
    auditorium_echo: t('effect_auditorium_echo'),
    lofi_telephone: t('effect_lofi_telephone'),
    robotic: t('effect_robotic')
  }
}

export const effectMap = getEffectMap()

export const getLanguageMap = () => {
  const { t } = useI18n('constants.voice')
  return {
    auto: t('language_auto'),
    Chinese: t('language_chinese'),
    Yue: t('language_yue'),
    English: t('language_english'),
    Arabic: t('language_arabic'),
    Russian: t('language_russian'),
    Spanish: t('language_spanish'),
    French: t('language_french'),
    Portuguese: t('language_portuguese'),
    German: t('language_german'),
    Turkish: t('language_turkish'),
    Dutch: t('language_dutch'),
    Ukrainian: t('language_ukrainian'),
    Vietnamese: t('language_vietnamese'),
    Indonesian: t('language_indonesian'),
    Japanese: t('language_japanese'),
    Italian: t('language_italian'),
    Korean: t('language_korean'),
    Thai: t('language_thai'),
    Polish: t('language_polish'),
    Romanian: t('language_romanian'),
    Greek: t('language_greek'),
    Czech: t('language_czech'),
    Finnish: t('language_finnish'),
    Hindi: t('language_hindi'),
    Bulgarian: t('language_bulgarian'),
    Danish: t('language_danish'),
    Hebrew: t('language_hebrew'),
    Malay: t('language_malay'),
    Persian: t('language_persian'),
    Slovak: t('language_slovak'),
    Swedish: t('language_swedish'),
    Croatian: t('language_croatian'),
    Filipino: t('language_filipino'),
    Hungarian: t('language_hungarian'),
    Norwegian: t('language_norwegian'),
    Slovenian: t('language_slovenian'),
    Catalan: t('language_catalan'),
    Nynorsk: t('language_nynorsk'),
    Tamil: t('language_tamil'),
    Afrikaans: t('language_afrikaans'),
  }
}

export const languageMap = getLanguageMap()

export const getVoiceSettingSliderConfig = () => {
  const { t } = useI18n('constants.voice')
  return {
    'speed': {
      ui_type: 'slider',
      key: 'voice_setting.speed',
      label: t('label_speed'),
      type: 'number',
      required: true,
      tooltip: t('tooltip_speed'),
      default: 1.0,
      min: 0.5,
      max: 2,
      step: 0.1,
      marks: {0.5: '0.5', 2: '2'}
    },
    'vol': {
      ui_type: 'slider',
      key: 'voice_setting.vol',
      label: t('label_volume'),
      type: 'number',
      required: true,
      tooltip: t('tooltip_volume'),
      default: 1,
      min: 0,
      max: 10,
      step: 1,
      marks: {0: '0', 10: '10'}
    },
    'pitch': {
      ui_type: 'slider',
      key: 'voice_setting.pitch',
      label: t('label_pitch'),
      type: 'number',
      required: true,
      tooltip: t('tooltip_pitch'),
      default: 0,
      min: -12,
      max: 12,
      step: 1,
      marks: {'-12': '-12', 12: '12'}
    }
  }
}

export const voiceSettingSliderConfig = getVoiceSettingSliderConfig()


export const getVoiceModifySliderConfig = () => {
  const { t } = useI18n('constants.voice')
  return {
    'pitch': {
      ui_type: 'slider',
      key: 'voice_modify.pitch',
      label: t('label_pitch'),
      type: 'number',
      tooltip: t('tooltip_pitch_modify'),
      default: 0,
      min: -100,
      max: 100,
      step: 1,
      marks: {'-100': '-100', 100: '100'}
    },

    'intensity': {
      ui_type: 'slider',
      key: 'voice_modify.intensity',
      label: t('label_intensity'),
      type: 'number',
      tooltip: t('tooltip_intensity'),
      default: 0,
      min: -100,
      max: 100,
      step: 1,
      marks: {'-100': '-100', 100: '100'}
    },

    'timbre': {
      ui_type: 'slider',
      key: 'voice_modify.timbre',
      label: t('label_timbre'),
      type: 'number',
      tooltip: t('tooltip_timbre'),
      default: 0,
      min: -100,
      max: 100,
      step: 1,
      marks: {'-100': '-100', 100: '100'}
    },
  }
}

export const voiceModifySliderConfig = getVoiceModifySliderConfig()

export const getVoiceOutputObj = () => {
  const { t } = useI18n('constants.voice')
  return {
    "data": {
      "desc": t('desc_audio_synthesis_result'),
      "type": "object",
      "properties": {
        "audio": {
          "desc": t('desc_audio_data'),
          "type": "string"
        },
        "status": {
          "desc": t('desc_audio_synthesis_status'),
          "type": "number"
        }
      }
    },
    "extra_info": {
      "desc": t('desc_audio_extra_info'),
      "type": "object",
      "properties": {
        "audio_length": {
          "desc": t('desc_audio_length'),
          "type": "number"
        },
        "audio_sample_rate": {
          "desc": t('desc_audio_sample_rate'),
          "type": "number"
        },
        "audio_size": {
          "desc": t('desc_audio_size'),
          "type": "number"
        },
        "bitrate": {
          "desc": t('desc_bitrate'),
          "type": "number"
        },
        "word_count": {
          "desc": t('desc_word_count'),
          "type": "number"
        },
        "invisible_character_ratio": {
          "desc": t('desc_invisible_character_ratio'),
          "type": "number"
        },
        "usage_characters": {
          "desc": t('desc_usage_characters'),
          "type": "number"
        },
        "audio_format": {
          "desc": t('desc_audio_format'),
          "type": "string"
        },
        "audio_channel": {
          "desc": t('desc_audio_channel'),
          "type": "number"
        }
      }
    },
    "trace_id": {
      "desc": t('desc_request_trace_id'),
      "type": "string"
    },
    "base_resp": {
      "desc": t('desc_base_response'),
      "type": "object",
      "properties": {
        "status_code": {
          "desc": t('desc_status_code'),
          "type": "number"
        },
        "status_msg": {
          "desc": t('desc_status_msg'),
          "type": "string"
        }
      }
    }
  }
}

export const voiceOutputObj = getVoiceOutputObj()

export const getVoiceCloneOutputObj = () => {
  const { t } = useI18n('constants.voice')
  return {
    "input_sensitive": {
      "desc": t('desc_input_sensitive'),
      "type": "boolean"
    },
    "input_sensitive_type": {
      "desc": t('desc_sensitive_type'),
      "type": "number"
    },
    "demo_audio": {
      "desc": t('desc_demo_audio'),
      "type": "string"
    },
    "base_resp": {
      "desc": t('desc_base_response'),
      "type": "object",
      "properties": {
        "status_code": {
          "desc": t('desc_status_code'),
          "type": "number"
        },
        "status_msg": {
          "desc": t('desc_status_msg'),
          "type": "string"
        }
      }
    }
  }
}

export const voiceCloneOutputObj = getVoiceCloneOutputObj()


