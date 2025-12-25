<style lang="less" scoped>
.ding-group-node {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-height: 52px;

  > div {
    display: flex;
    word-break: break-all;

    .label {
      flex-shrink: 0;
      text-align: right;
      white-space: normal;
    }

    :deep(.ant-tag) {
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }
  }
}
</style>

<template>
  <node-common
    :title="properties.node_name"
    :iconUrl="info.icon"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    style="min-width: 320px;min-height: 157px;"
  >
    <div class="ding-group-node">
      <FeishuTableRender
        v-if="isFeishuTbPlugin"
        :nodeParams="nodeParams"
        :business="business"
        :valueOptions="valueOptions"
        @updateSize="updateSize"
      />
      <OfficialSendMessageRender
        v-if="isOfficialSendMessagePlugin"
        :nodeParams="nodeParams"
        :business="business"
        :valueOptions="valueOptions"
        @updateSize="updateSize"
      />
      <template v-else>
        <div v-for="(item, key) in formState" :key="key">
          <span class="label">{{ item.name || key }}：</span>
          <a-tag v-if="item.value">
            <at-text :options="valueOptions" :defaultSelectedList="item.tags" :defaultValue="item.value"/>
          </a-tag>
          <a-tag v-else>--</a-tag>
        </div>
      </template>
    </div>
  </node-common>
</template>

<script>
import NodeCommon from '../base-node.vue'
import {getPluginInfo, runPlugin} from "@/api/plugins/index.js";
import {pluginHasAction} from "@/constants/plugin.js";
import FeishuTableRender from "./components/feishu-table-render.vue";
import OfficialSendMessageRender from "./components/official-send-message-render.vue";
import AtText from "@/views/workflow/components/at-input/at-text.vue";

export default {
  name: 'QuestionNode',
  components: {
    AtText,
    FeishuTableRender,
    OfficialSendMessageRender,
    NodeCommon,
  },
  inject: ['getNode', 'getGraph', 'setData', 'resetSize'],
  props: {
    properties: {
      type: Object,
      default() {
        return {}
      },
    },
    isSelected: {type: Boolean, default: false},
    isHovered: {type: Boolean, default: false},
  },
  data() {
    return {
      menus: [{name: '删除', key: 'delete', color: '#fb363f'}],
      formState: {
        // app_id: {
        //   type: 'string',
        //   desc: '微信公众号app_id'
        //   name: 'app_id'
        // }
      },
      nodeParams: {},
      info: {},
      pluginData: {},
      valueOptions: []
    }
  },
  mounted() {
    this.init()
  },
  computed: {
    pluginName() {
      return this.nodeParams?.plugin?.name || null
    },
    business() {
      return this.nodeParams?.plugin?.params?.business || null
    },
    isFeishuTbPlugin() {
      return this.pluginName === 'feishu_bitable'
    },
    isOfficialSendMessagePlugin() {
      return this.pluginName === 'official_send_message'
    }
  },
  watch: {
    properties: {
      handler(newVal, oldVal) {
        if (newVal.dataRaw !== oldVal.dataRaw) {
          this.nodeParamsFormat()
          this.formStateFormat()
        }
      },
      deep: true
    }
  },
  methods: {
    init() {
      this.valueOptions = this.getNode().getAllParentVariable()
      this.nodeParamsFormat()
      this.loadPluginParams()
      this.loadPluginInfo()
      this.updateSize()
    },
    loadPluginInfo() {
      getPluginInfo({name: this.pluginName}).then(res => {
        let data = res?.data || {}
        this.info = data.remote
      })
    },
    loadPluginParams() {
      runPlugin({
        name: this.pluginName,
        action: "default/get-schema",
        params: {}
      }).then(res => {
        console.log('res?.data', res?.data)
        this.pluginData = res?.data || {}
        this.formStateFormat()
      })
    },
    nodeParamsFormat() {
      this.nodeParams = JSON.parse(this.properties.node_params)
    },
    formStateFormat() {
      if (this.isFeishuTbPlugin || this.isOfficialSendMessagePlugin) return
      let data = JSON.parse(JSON.stringify(this.pluginData))
      const customHandleFunc = {
        official_account_profile: this.officialAccountHandle,
        official_batch_tag: this.officialTagHandle,
        official_send_template_message: this.officialTemplateHandle,
        official_article: this.officialArticleHandle,
        official_draft: this.officialDraftHandle,
      }
      if (customHandleFunc[this.pluginName] ) {
        customHandleFunc[this.pluginName](data)
      } else {
        this.defaultHandle(data)
      }
      this.updateSize()
    },
    defaultHandle(data) {
      let args, tag_map
      if (pluginHasAction(this.pluginName)) {
        this.formState = data?.[this.business]?.params || {}
        args = this.nodeParams?.plugin?.params?.arguments || {}
        tag_map = this.nodeParams?.plugin?.params?.arguments?.tag_map || {}
      } else {
        args = this.nodeParams?.plugin?.params
        tag_map = this.nodeParams?.plugin?.tag_map || {}
        this.formState = data
      }
      for (let key in this.formState) {
        this.formState[key].value = String(args[key] || '')
        this.formState[key].tags = tag_map[key] || []
      }
    },
    officialAccountHandle(data) {
      delete data?.[this.business]?.params.app_id
      delete data?.[this.business]?.params.app_secret
      if (data?.[this.business]?.params) {
        if (this.business === 'getTagFans') {
          let args = this.nodeParams?.plugin?.params?.arguments || {}
          let tag_map = this.nodeParams?.plugin?.params?.arguments?.tag_map || {}
          this.formState = {
            app_name: {type: 'string', name: '公众号名称', value: args.app_name, tags: []},
          }
          if (args.tag_type == 1) {
            this.formState.tag_name = {type: 'string', name: '标签名称', value: args.tag_name, tags: []}
          } else if (tag_map?.tagid?.length) {
            this.formState.tagid = {type: 'string', name: '标签名称', value: args.tagid, tags: tag_map.tagid || []}
          } else {
            this.formState.tagid = {type: 'string', name: '标签ID', value: args.tagid, tags: []}
          }
          this.formState.next_openid = {type: 'string', name: '上批尾用户openid', value: args.next_openid, tags: tag_map.next_openid || []}
          this.updateSize()
          return
        } else {
          // 为了让app_name显示在前面
          data[this.business].params = {
            app_name: {
              type: 'string',
              name: '公众号名称',
            },
            ...data[this.business].params
          }
        }
      }
      this.defaultHandle(data)
    },
    officialTemplateHandle() {
      let args = this.nodeParams?.plugin?.params?.arguments || {}
      let tag_map = this.nodeParams?.plugin?.params?.arguments?.tag_map || {}
      this.formState = {
        app_name: {type: 'string', name: '公众号名称', value: args.app_name, tags: []},
        template_name: {type: 'string', name: '模板名称', value: args.template_name, tags: []},
        touser: {type: 'string', name: '接收者', value: args.touser, tags: tag_map.touser || []},
      }
      switch (args.link_type) {
        case 0:
          this.formState.link_type = {type: 'string', name: '模板跳转', value: '不跳转', tags: []}
          break
        case 1:
          this.formState.link_type = {type: 'string', name: '模板跳转', value: '跳转链接', tags: []}
          this.formState.url = {type: 'string', name: '跳转链接', value: args.url, tags: tag_map.url || []}
          break
        case 2:
          this.formState.link_type = {type: 'string', name: '模板跳转', value: '跳转小程序', tags: []}
          this.formState.miniprogram_app = {type: 'string', name: '小程序APPID', ...args.miniprogram.appid }
          this.formState.miniprogram_path = {type: 'string', name: '小程序路径', ...args.miniprogram.pagepath }
          break
      }
      this.updateSize()
    },
    officialTagHandle() {
      let args = this.nodeParams?.plugin?.params?.arguments || {}
      let tag_map = this.nodeParams?.plugin?.params?.arguments?.tag_map || {}
      this.formState = {
        app_name: {type: 'string', name: '公众号名称', value: args.app_name, tags: []},
        openid_list: {type: 'string', name: '粉丝openid列表', value: args.openid_list, tags: []},
      }
      if (args.tag_type == 1) {
        this.formState.tag_name = {type: 'string', name: '标签名称', value: args.tag_name, tags: []}
      } else {
        this.formState.tagid = {type: 'string', name: '标签名称', value: args.tagid, tags: tag_map.tagid || []}
      }
      this.updateSize()
    },
    officialArticleHandle() {
      let args = this.nodeParams?.plugin?.params?.arguments || {}
      let tag_map = this.nodeParams?.plugin?.params?.arguments?.tag_map || {}
      this.formState = {
        url: {type: 'string', name: '', value: args.url, tags: tag_map.url || []},
        number: {type: 'string', name: '', value: args.number, tags: tag_map.number || []},
      }
      this.updateSize()
    },
    officialDraftHandle() {
      let args = this.nodeParams?.plugin?.params?.arguments || {}
      let tag_map = this.nodeParams?.plugin?.params?.arguments?.tag_map || {}
      let article_type_text = !args.article_type ? '--' : args.article_type == 'news' ? '图文消息' : '图片消息'
      switch (this.business) {
        case 'add_draft':
          this.formState = {
            app_name: {type: 'string', name: '公众号名称', value: args.app_name, tags: []},
            article_type: {type: 'string', name: '文章类型', value: article_type_text, tags: []},
          }
          break
        case 'update_draft':
        case 'publish_draft':
        case 'delete_draft':
        case 'preview_message':
          this.formState = {
            app_name: {type: 'string', name: '公众号名称', value: args.app_name, tags: []},
            media_id: {type: 'string', name: '文章ID', value: args.media_id, tags: tag_map.media_id || []},
          }
          if (this.business === 'preview_message') {
            this.formState.account = {type: 'string', name: '接收用户', value: args.account, tags: tag_map.account || []}
          }
          break
      }
      this.updateSize()
    },
    updateSize() {
      this.$nextTick(() => {
        this.resetSize()
      })
    }
  },
}
</script>
