<template>
  <node-common
    :properties="properties"
    :title="properties.node_name"
    :iconUrl="properties.node_icon || info.icon"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    style="min-width: 320px;min-height: 157px;"
  >
    <div v-if="isMultiNode" class="ding-group-node">
      <DynamicRender
        :nodeParams="nodeParams"
        :business="business"
        :valueOptions="valueOptions"
        @updateSize="updateSize"
      />
    </div>
    <div v-else class="ding-group-node">
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
            <at-text :options="valueOptions" :defaultSelectedList="item.tags" :defaultValue="item.value" ref="atTextRef" />
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
import DynamicRender from "./components/dynamic-render.vue";
import AtText from "@/views/workflow/components/at-input/at-text.vue";
import { pluginOutputToTree } from "@/constants/plugin.js";
import { haveOutKeyNode } from '@/views/workflow/components/util.js'
import { useI18n } from '@/hooks/web/useI18n'

export default {
  name: 'QuestionNode',
  components: {
    AtText,
    FeishuTableRender,
    OfficialSendMessageRender,
    DynamicRender,
    NodeCommon
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
  setup() {
    const { t } = useI18n('views.workflow.components.nodes.zm-plugins-node.index')
    return {
      t
    }
  },
  data() {
    return {
      atTextRef: null,
      menus: [{name: 'btn_delete', key: 'delete', color: '#fb363f'}],
      formState: {
        // app_id: {
        //   type: 'string',
        //   desc: '微信公众号app_id'
        //   name: 'app_id'
        // }
      },
      nodeParams: {},
      info: {},
      isMultiNode: false,
      pluginData: {},
      valueOptions: []
    }
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
  mounted() {
    const graphModel = this.getGraph()
    graphModel.eventCenter.on('custom:setNodeName', this.onUpatateNodeName)
    this.init()
  },
  onBeforeUnmount() {
    const graphModel = this.getGraph()
    graphModel.eventCenter.off('custom:setNodeName', this.onUpatateNodeName)
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
    onUpatateNodeName(data) {
      if (!haveOutKeyNode.includes(data.node_type)) {
        return
      }

      this.valueOptions = this.getNode().getAllParentVariable() || []

      const patchTags = (tags) => {
        if (!Array.isArray(tags) || tags.length === 0) {
          return
        }

        tags.forEach((tag) => {
          if (tag && tag.node_id == data.node_id) {
            if (tag.label) {
              let arr = tag.label.split('/')
              arr[0] = data.node_name
              tag.label = arr.join('/')
            }
            tag.node_name = data.node_name
          }
        })
      }

      if (this.formState && typeof this.formState === 'object') {
        Object.keys(this.formState).forEach((key) => {
          patchTags(this.formState[key]?.tags)
        })
      }

      const tagMaps = []

      if (this.nodeParams?.plugin?.tag_map) {
        tagMaps.push(this.nodeParams.plugin.tag_map)
      }
      if (this.nodeParams?.plugin?.params?.arguments?.tag_map) {
        tagMaps.push(this.nodeParams.plugin.params.arguments.tag_map)
      }

      tagMaps.forEach((tagMap) => {
        if (!tagMap || typeof tagMap !== 'object') {
          return
        }
        Object.keys(tagMap).forEach((k) => {
          patchTags(tagMap[k])
        })
      })

      this.$nextTick(() => {
        const atTextRef = this.$refs.atTextRef
        const refList = Array.isArray(atTextRef) ? atTextRef : atTextRef ? [atTextRef] : []
        refList.forEach((ref) => {
          ref && ref.refresh && ref.refresh()
        })

        this.setData({
          node_params: JSON.stringify(this.nodeParams)
        })

        this.updateSize()
      })
    },
    async init() {
      this.valueOptions = this.getNode().getAllParentVariable()
      this.nodeParamsFormat()
      await this.loadPluginInfo()
      await this.loadPluginParams()
      this.updateSize()
    },
    loadPluginInfo() {
      return getPluginInfo({name: this.pluginName}).then(res => {
        let data = res?.data || {}
        this.info = data.remote || {}
        this.isMultiNode = data?.local?.multiNode || false
      }).catch(() => {
        this.info = { icon: this.properties.node_icon || '' }
        this.isMultiNode = false
      })
    },
    loadPluginParams() {
      return runPlugin({
        name: this.pluginName,
        action: "default/get-schema",
        params: {}
      }).then(res => {
        this.pluginData = res?.data || {}
        this.formStateFormat()
        this.update(this.pluginData)
      }).catch(() => {
        this.pluginData = {}
        this.formStateFormat()
        this.update(this.pluginData)
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
    update (data) {
      if (typeof data === 'string') {
        data = JSON.parse(data)
      }
      const nodeParams = JSON.parse(this.properties.node_params)
      const schemaAction = (this.business && data?.[this.business]) ? data[this.business] : data
      const isAction = pluginHasAction(this.pluginName) || this.isMultiNode

      const oldArgs = isAction
        ? (nodeParams?.plugin?.params?.arguments || {})
        : (nodeParams?.plugin?.params || {})

      const updatedArgs = { ...oldArgs }
      for (const key in (this.formState || {})) {
        const meta = this.formState[key] || {}
        const isDefaultableField = meta?.enum_component || meta?.radio_component
        const def = meta?.enum_default
        if (isDefaultableField && (updatedArgs[key] == null || String(updatedArgs[key]).trim() === '')) {
          if (def != null && def !== '') {
            updatedArgs[key] = def
          }
        }
        const v = meta?.setting_cache_component
        const isCache = v === true || v === 1 || v === 'true' || v === '1'
        if (isCache) {
          const val = String(meta?.value ?? '')
          if (val !== '') {
            updatedArgs[key] = val
          } else {
            const cacheVal = localStorage.getItem(this.getCacheKey(key))
            if (cacheVal != null) {
              updatedArgs[key] = String(cacheVal)
            }
          }
        }
      }

      nodeParams.plugin = nodeParams.plugin || {}
      nodeParams.plugin.params = nodeParams.plugin.params || {}
      const currentUse = nodeParams?.plugin?.params?.use_plugin_config
      nodeParams.plugin.params.use_plugin_config = (schemaAction && Object.prototype.hasOwnProperty.call(schemaAction, 'use_plugin_config'))
        ? !!schemaAction.use_plugin_config
        : !!currentUse

      if (isAction) {
        nodeParams.plugin.params.arguments = {
          ...oldArgs,
          ...updatedArgs,
          tag_map: oldArgs.tag_map || {}
        }
      } else {
        nodeParams.plugin.params = {
          ...nodeParams.plugin.params,
          ...updatedArgs
        }
      }
      nodeParams.plugin.output_obj = pluginOutputToTree(schemaAction?.output || {})
      const newNodeParamsStr = JSON.stringify(nodeParams)
      if (this.properties.node_params !== newNodeParamsStr) {
        this.setData({
          ...this.properties,
          node_params: newNodeParamsStr
        })
      }
    },
    defaultHandle(data) {
      let args, tag_map
      if (pluginHasAction(this.pluginName) || this.isMultiNode) {
        this.formState = data?.[this.business]?.params || {}
        args = this.nodeParams?.plugin?.params?.arguments || {}
        tag_map = this.nodeParams?.plugin?.params?.arguments?.tag_map || {}
      } else {
        args = this.nodeParams?.plugin?.params || {}
        tag_map = this.nodeParams?.plugin?.tag_map || {}
        this.formState = data
      }
      for (let key in this.formState) {
        this.formState[key].value = String(args[key] || '')
        this.formState[key].tags = tag_map[key] || []
        const v = this.formState[key]?.setting_cache_component
        const enabled = v === true || v === 1 || v === 'true' || v === '1'
        if (enabled) {
          const cur = String(this.formState[key].value || '').trim()
          if (!cur) {
            const cacheVal = localStorage.getItem(this.getCacheKey(key))
            if (cacheVal != null) {
              this.formState[key].value = String(cacheVal)
            }
          }
        }
      }
      if (this.isMultiNode) {
        this.ensureMultiNodeRendering()
      }
    },
    getCacheKey(fieldKey) {
      let id = ''
      let robotKey = ''
      const hash = typeof window !== 'undefined' ? String(window.location.hash || '') : ''
      if (hash.includes('?')) {
        const query = hash.split('?')[1] || ''
        query.split('&').forEach((pair) => {
          const [k, v] = pair.split('=')
          if (k === 'id') id = decodeURIComponent(v || '')
          if (k === 'robot_key') robotKey = decodeURIComponent(v || '')
        })
      }
      if (!id) {
        id = String(localStorage.getItem('last_local_robot_id') || '')
      }
      const wfId = id || robotKey || ''
      return `${wfId}:${this.pluginName}:${fieldKey}`
    },
    ensureMultiNodeRendering() {
      const current = this.nodeParams?.plugin?.params?.rendering
      if (Array.isArray(current) && current.length) {
        const tag_map = this.nodeParams?.plugin?.params?.arguments?.tag_map || this.nodeParams?.plugin?.tag_map || {}
        let changed = false
        const list = current.map((item) => {
          const key = item?.key
          const originTags = Array.isArray(item?.tags) ? item.tags : []
          if (originTags.length) return item

          const nextTags = Array.isArray(tag_map?.[key])
            ? tag_map[key]
            : Array.isArray(this.formState?.[key]?.tags)
              ? this.formState[key].tags
              : []

          if (nextTags.length) {
            changed = true
            return { ...item, tags: nextTags }
          }
          return item
        })

        if (changed) {
          const nodeParams = JSON.parse(this.properties.node_params)
          nodeParams.plugin = nodeParams.plugin || {}
          nodeParams.plugin.params = nodeParams.plugin.params || {}
          nodeParams.plugin.params.rendering = list

          if (this.nodeParams?.plugin?.params) {
            this.nodeParams.plugin.params.rendering = list
          }

          this.setData({
            ...this.properties,
            node_params: JSON.stringify(nodeParams)
          })
        }
        return
      }

      const keys = Object.keys(this.formState || {})
        .filter((key) => key !== 'app_secret')
        .filter((key) => !this.formState?.[key]?.hide_official_component)
        .filter((key) => !this.formState?.[key]?.advanced_settings)
        .filter((key) => {
          const v = this.formState?.[key]?.plugin_config_component
          const enabled = v === true || v === 1 || v === 'true' || v === '1'
          return !enabled
        })
        .filter((key) => {
          const v = this.formState?.[key]?.setting_cache_component
          const enabled = v === true || v === 1 || v === 'true' || v === '1'
          return !enabled
        })
        .sort((a, b) => (+(this.formState?.[a]?.sort_num || 0)) - (+(this.formState?.[b]?.sort_num || 0)))

      const args = this.nodeParams?.plugin?.params?.arguments || {}

      const list = keys.map((key) => {
        const meta = this.formState?.[key] || {}
        let label = meta?.name || key
        let value = String(meta?.value ?? '')
        const tags = Array.isArray(meta?.tags) ? meta.tags : []

        if (key === 'app_id') {
          label = this.t('label_official_account_name')
          value = String(args.app_name ?? this.formState?.app_name?.value ?? '')
        }

        if (meta?.enum_component || meta?.radio_component) {
          const opt = (meta?.enum || []).find((e) => String(e?.value) === value) || (meta?.enum || []).find((e) => String(e?.value) === meta?.enum_default)
          if (opt && opt.name) {
            value = String(opt.name)
          }
        }

        return { label, value, tags, key }
      })
      if (!list.length) return

      const nodeParams = JSON.parse(this.properties.node_params)
      nodeParams.plugin = nodeParams.plugin || {}
      nodeParams.plugin.params = nodeParams.plugin.params || {}
      nodeParams.plugin.params.rendering = list

      if (this.nodeParams?.plugin?.params) {
        this.nodeParams.plugin.params.rendering = list
      }

      this.setData({
        ...this.properties,
        node_params: JSON.stringify(nodeParams)
      })
    },
    officialAccountHandle(data) {
      delete data?.[this.business]?.params.app_id
      delete data?.[this.business]?.params.app_secret
      if (data?.[this.business]?.params) {
        if (this.business === 'getTagFans') {
          let args = this.nodeParams?.plugin?.params?.arguments || {}
          let tag_map = this.nodeParams?.plugin?.params?.arguments?.tag_map || {}
          this.formState = {
            app_name: {type: 'string', name: this.t('label_official_account_name'), value: args.app_name, tags: []},
          }
          if (args.tag_type == 1) {
            this.formState.tag_name = {type: 'string', name: this.t('label_tag_name'), value: args.tag_name, tags: []}
          } else if (tag_map?.tagid?.length) {
            this.formState.tagid = {type: 'string', name: this.t('label_tag_name'), value: args.tagid, tags: tag_map.tagid || []}
          } else {
            this.formState.tagid = {type: 'string', name: this.t('label_tag_id'), value: args.tagid, tags: []}
          }
          this.formState.next_openid = {type: 'string', name: this.t('label_last_batch_openid'), value: args.next_openid, tags: tag_map.next_openid || []}
          this.updateSize()
          return
        } else {
          // 为了让app_name显示在前面
          data[this.business].params = {
            app_name: {
              type: 'string',
              name: this.t('label_official_account_name'),
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
        app_name: {type: 'string', name: this.t('label_official_account_name'), value: args.app_name, tags: []},
        template_name: {type: 'string', name: this.t('label_template_name'), value: args.template_name, tags: []},
        touser: {type: 'string', name: this.t('label_receiver'), value: args.touser, tags: tag_map.touser || []},
      }
      switch (args.link_type) {
        case 0:
          this.formState.link_type = {type: 'string', name: this.t('label_template_jump'), value: this.t('text_no_jump'), tags: []}
          break
        case 1:
          this.formState.link_type = {type: 'string', name: this.t('label_template_jump'), value: this.t('text_jump_link'), tags: []}
          this.formState.url = {type: 'string', name: this.t('label_jump_link'), value: args.url, tags: tag_map.url || []}
          break
        case 2:
          this.formState.link_type = {type: 'string', name: this.t('label_template_jump'), value: this.t('text_jump_miniprogram'), tags: []}
          this.formState.miniprogram_app = {type: 'string', name: this.t('label_miniprogram_appid'), ...args.miniprogram.appid }
          this.formState.miniprogram_path = {type: 'string', name: this.t('label_miniprogram_path'), ...args.miniprogram.pagepath }
          break
      }
      this.updateSize()
    },
    officialTagHandle() {
      let args = this.nodeParams?.plugin?.params?.arguments || {}
      let tag_map = this.nodeParams?.plugin?.params?.arguments?.tag_map || {}
      this.formState = {
        app_name: {type: 'string', name: this.t('label_official_account_name'), value: args.app_name, tags: []},
        openid_list: {type: 'string', name: this.t('label_fans_openid_list'), value: args.openid_list, tags: []},
      }
      if (args.tag_type == 1) {
        this.formState.tag_name = {type: 'string', name: this.t('label_tag_name'), value: args.tag_name, tags: []}
      } else {
        this.formState.tagid = {type: 'string', name: this.t('label_tag_name'), value: args.tagid, tags: tag_map.tagid || []}
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
      let article_type_text = !args.article_type ? '--' : args.article_type == 'news' ? this.t('text_news_message') : this.t('text_image_message')
      switch (this.business) {
        case 'add_draft':
          this.formState = {
            app_name: {type: 'string', name: this.t('label_official_account_name'), value: args.app_name, tags: []},
            article_type: {type: 'string', name: this.t('label_article_type'), value: article_type_text, tags: []},
          }
          break
        case 'update_draft':
        case 'publish_draft':
        case 'delete_draft':
        case 'preview_message':
        case 'get_draft':
        default:
          this.formState = {
            app_name: {type: 'string', name: this.t('label_official_account_name'), value: args.app_name, tags: []},
            media_id: {type: 'string', name: this.t('label_article_id'), value: args.media_id, tags: tag_map.media_id || []},
          }
          if (this.business === 'preview_message') {
            this.formState.account = {type: 'string', name: this.t('label_receive_user'), value: args.account, tags: tag_map.account || []}
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
