<style lang="less" scoped>
.ding-group-node {
  display: flex;
  flex-direction: column;
  gap: 8px;

  > div {
    display: flex;
    align-items: center;

    > .label {
      flex-shrink: 0;
      width: 90px;
      text-align: right;
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
  >
    <div class="ding-group-node">
      <div v-for="(item, key) in formState" :key="key">
        <span class="label">{{ key }}：</span>
        <a-tag>{{ item.value || '--' }}</a-tag>
      </div>
    </div>
  </node-common>
</template>

<script>
import {jsonDecode} from '@/utils/index'
import NodeCommon from '../base-node.vue'
import {getPluginInfo, runPlugin} from "@/api/plugins/index.js";

export default {
  name: 'QuestionNode',
  components: {
    NodeCommon,
  },
  inject: ['getNode', 'getGraph', 'setData'],
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
      formState: {},
      nodeParams: {},
      info: {},
    }
  },
  mounted() {
    this.init()
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
      this.nodeParamsFormat()
      this.loadPluginParams()
      this.loadPluginInfo()
    },
    loadPluginInfo() {
      getPluginInfo({name: this.nodeParams?.plugin?.name}).then(res => {
        let data = res?.data || {}
        this.info = data.remote
      })
    },
    loadPluginParams() {
      runPlugin({
        name: this.nodeParams?.plugin?.name || '',
        action: "default/get-schema",
        params: {}
      }).then(res => {
        this.formState = res?.data || {}
        this.formStateFormat()
      })
    },
    nodeParamsFormat() {
      this.nodeParams = JSON.parse(this.properties.node_params)
    },
    formStateFormat() {
      for (let key in this.formState) {
        this.formState[key].value = String(this.nodeParams?.plugin?.params[key] || '')
        this.formState[key].tags = this.nodeParams?.plugin?.tag_map || []
      }
    },
  },
}
</script>
