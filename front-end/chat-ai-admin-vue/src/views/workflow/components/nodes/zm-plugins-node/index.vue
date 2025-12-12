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
import AtText from "@/views/workflow/components/at-input/at-text.vue";

export default {
  name: 'QuestionNode',
  components: {
    AtText,
    FeishuTableRender,
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
        this.pluginData = res?.data || {}
        this.formStateFormat()
      })
    },
    nodeParamsFormat() {
      this.nodeParams = JSON.parse(this.properties.node_params)
    },
    formStateFormat() {
      if (this.isFeishuTbPlugin) return
      let data = JSON.parse(JSON.stringify(this.pluginData))
      const customHandleFunc = {
        official_account_profile: this.officialAccountHandle
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
        // 为了让app_name显示在前面
        data[this.business].params = {
          app_name: {
            type: 'string',
            name: '公众号名称',
          },
          ...data[this.business].params
        }
      }
      this.defaultHandle(data)
    },
    updateSize() {
      this.$nextTick(() => {
        this.resetSize()
      })
    }
  },
}
</script>
