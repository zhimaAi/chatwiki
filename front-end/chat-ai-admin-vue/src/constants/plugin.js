import {getPluginConfig} from "@/api/plugins/index.js";
import {jsonDecode} from "@/utils/index.js";

// 方法插件
export const HasActionPluginNames = [
  'feishu_bitable',
  'official_account_profile',
  'official_batch_tag',
  'official_send_template_message',
  'web_content_extraction',
  'official_article',
  'official_send_message',
  'official_draft'
]

// 存在输出内容的插件
export const HasOutputPluginNames = [
  'feishu_bitable',
  'official_account_profile',
  'official_batch_tag',
  'official_send_template_message',
  'web_content_extraction',
  'official_article',
  'official_send_message',
  'official_draft'
]

// 插件方法默认参数（非特殊处理无需在此添加-将通过action params 自动获取）
export const PluginActionDefaultArgumentsMap = {
  feishu_bitable: {
    create_record: {
      app_id: '',
      app_secret: '',
      app_token: '',
      table_id: '',
      fields: [],
    },
    delete_record: {
      app_id: '',
      app_secret: '',
      app_token: '',
      table_id: '',
      record_id: '',
      record_tags: []
    },
    update_record: {
      app_id: '',
      app_secret: '',
      app_token: '',
      table_id: '',
      fields: [],
      record_id: '',
      record_tags: [],
    },
    search_records: {
      app_id: '',
      app_secret: '',
      app_token: '',
      table_id: '',
      field_names: [],
      filter: {
        conjunction:'and',
        conditions: [],
      },
      sort: [],
      page_size: 100
    }
  },
  official_account_profile: {
    batch_get_user_info: {
      app_id: '',
      app_secret: '',
      user_list: [],
      tag_map: {},
    },
    get_fans: {
      app_id: '',
      app_secret: '',
      next_openid: '',
      tag_map: {},
    },
    get_user_info: {
      app_id: '',
      app_secret: '',
      openid: '',
      tag_map: {},
    },
  }
}

// 是否是方法插件
export const pluginHasAction = (name) => {
  return HasActionPluginNames.includes(name)
}

// 获取插件方法默认参数
export const getPluginActionDefaultArguments = (pluginName, actionName) => {
  return JSON.parse(JSON.stringify(PluginActionDefaultArgumentsMap?.[pluginName]?.[actionName] || {}))
}

export function pluginOutputToTree(obj) {
  function randomKey() {
    return Math.random() * 10000;
  }

  function walk(key, value) {
    const node = {
      key,
      name: value.name || '',
      desc: value.desc || '',
      title: key,
      typ: "object",
      subs: [],
      cu_key: randomKey()
    };

    // 如果 value 有 type，优先使用
    if (value.type) {
      node.typ = value.type;
    } else if (value.properties) {
      node.typ = "object";
    }

    if (value.properties && typeof value.properties === "object") {
      // 如果包含 properties，则递归处理
      node.subs = Object.entries(value.properties).map(([k, v]) =>
        walk(k, v)
      );
    } else if (value.type === "array<object>") {
      // 如果是对象数组
      node.subs = Object.entries(value.items.properties).map(([k, v]) =>
        walk(k, v)
      );
    } else {
      // 无 properties 且是基础类型，则 subs 留空
      node.subs = [];
    }

    return node;
  }

  return Object.entries(obj).map(([key, value]) => walk(key, value));
}

// 插件配置
const _pluginConfigMap = {}
export const getPluginConfigData = async (name) => {
  if (!_pluginConfigMap[name]) {
    await getPluginConfig({name: name}).then(res => {
      _pluginConfigMap[name] = jsonDecode(res?.data, {})
    })
  }
  return _pluginConfigMap[name]
}
