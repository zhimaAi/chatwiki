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

.field-box {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;

  .field-tag {
    display: flex;
    align-items: center;
    gap: 4px;
    margin: 0;
    word-break: break-all;
  }

  .arrow-icon {
    display: inline-flex;
    align-items: center;
    margin: 0 4px;
    padding: 0 4px;
    height: 18px;
    font-size: 12px;
    border-radius: 4px;
    background: #e4e6eb;

    img {
      width: 16px;
      height: 16px;
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
      <template v-if="business">
        <!--飞书多维表插件方法-->
        <div><span class="label">配置：</span>
          <a-tag>{{ actionState?.config_name || '--' }}</a-tag>
        </div>
        <div><span class="label">数据表：</span>
          <a-tag>{{ actionState?.table_name || '--' }}</a-tag>
        </div>
        <template v-if="business == 'create_record'">
          <div>
            <span class="label">插入数据：</span>
            <div v-if="Array.isArray(actionState.fields)" class="field-box">
              <a-tag v-for="field in actionState.fields" class="field-tag" :key="field.field_name">
                <span>{{ field.field_name }}</span>
                <span class="arrow-icon">
                  <img src="@/assets/img/workflow/arrow-right.svg"/>
                </span>
                <span>{{ field.value.toString() }}</span>
              </a-tag>
            </div>
            <div v-else>
              <a-tag>--</a-tag>
            </div>
          </div>
        </template>
        <template v-else-if="business == 'delete_record'">
          <div>
            <span class="label">删除ID：</span>
            <span>{{ actionState?.record_id || '--' }}</span>
          </div>
        </template>
        <template v-else-if="business == 'search_records'">
          <div>
            <span class="label">查询条件：</span>
            <div v-if="Array.isArray(actionState?.filter?.conditions)" class="field-box">
              <a-tag v-for="item in actionState.filter.conditions" class="field-tag" :key="item.field_id">
                <span>{{ item.field_name }}</span>
                <span class="arrow-icon">{{ operatorMap[item.operator]?.label }}</span>
                <span>{{ item.value }}</span>
              </a-tag>
            </div>
            <div v-else>
              <a-tag>--</a-tag>
            </div>
          </div>
          <div>
            <span class="label">查询字段：</span>
            <div v-if="Array.isArray(actionState.field_names)" class="field-box">
              <a-tag v-for="field in actionState.field_names" class="field-tag" :key="field.field">{{ field }}</a-tag>
            </div>
            <div v-else>
              <a-tag>--</a-tag>
            </div>
          </div>
        </template>
        <template v-else-if="business == 'update_record'">
          <div>
            <span class="label">更新ID：</span>
            <span>{{ actionState?.record_id || '--' }}</span>
          </div>
          <div>
            <span class="label">更新数据：</span>
            <div v-if="Array.isArray(actionState.fields)" class="field-box">
              <a-tag v-for="field in actionState.fields" class="field-tag" :key="field.field_name">
                <span>{{ field.field_name }}</span>
                <span class="arrow-icon">
                  <img src="@/assets/img/workflow/arrow-right.svg"/>
                </span>
                <span>{{ field.value.toString() }}</span>
              </a-tag>
            </div>
            <div v-else>
              <a-tag>--</a-tag>
            </div>
          </div>
        </template>
      </template>
      <template v-else>
        <div v-for="(item, key) in formState" :key="key">
          <span class="label">{{ key }}：</span>
          <a-tag>{{ item.value || '--' }}</a-tag>
        </div>
      </template>
    </div>
  </node-common>
</template>

<script>
import NodeCommon from '../base-node.vue'
import {getPluginInfo, runPlugin} from "@/api/plugins/index.js";
import {FeiShuOperatorMap} from "@/constants/feishu-table.js";

export default {
  name: 'QuestionNode',
  components: {
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
      operatorMap: FeiShuOperatorMap(),
      menus: [{name: '删除', key: 'delete', color: '#fb363f'}],
      formState: {},
      nodeParams: {},
      info: {},
      actionState: {},
    }
  },
  mounted() {
    this.init()
  },
  computed: {
    business() {
      return this.nodeParams?.plugin?.params?.business || null
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
      this.nodeParamsFormat()
      this.loadPluginParams()
      this.loadPluginInfo()

      this.$nextTick(() => {
        this.resetSize()
      })
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
      if (this.business) {
        let arg = this.nodeParams?.plugin?.params?.arguments || {}
        let {config_name, table_name} = this.nodeParams?.plugin?.params || {}
        this.actionState = {
          ...arg,
          config_name,
          table_name,
        }
      } else {
        for (let key in this.formState) {
          this.formState[key].value = String(this.nodeParams?.plugin?.params[key] || '')
          this.formState[key].tags = this.nodeParams?.plugin?.tag_map || []
        }
      }

      this.$nextTick(() => {
        this.resetSize()
      })
    },
  },
}
</script>

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