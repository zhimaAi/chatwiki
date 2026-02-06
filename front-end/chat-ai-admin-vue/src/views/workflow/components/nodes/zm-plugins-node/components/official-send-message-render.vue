<template>
  <div class="_main">
    <div>
      <span class="label">公众号：</span>
      <a-tag>{{ actionState?.app_name || '--' }}</a-tag>
    </div>
    <div>
      <span class="label">接收人
        <span class="hint">
          <QuestionCircleOutlined />
          <span class="hint-pop">接收人openid</span>
        </span>：
      </span>
      <span v-if="actionState?.touser">
        <at-text
          :options="valueOptions"
          :defaultSelectedList="[]"
          :defaultValue="(actionState?.touser || '').toString()"/>
      </span>
      <span v-else>
        <a-tag>--</a-tag>
      </span>
    </div>
    <div><span class="label">消息类型：</span>
      <a-tag>{{ actionState?.msg_type_name || '--' }}</a-tag>
    </div>
    <div>
      <span class="label">回复内容：</span>
      <div v-if="Array.isArray(actionState.fields)" class="field-box">
        <a-tag v-for="field in actionState.fields" class="field-tag" :key="field.field_name">
          <span>{{ getFieldLabel(field) }}</span>
          <span class="arrow-icon">
            <img src="@/assets/img/workflow/arrow-right.svg"/>
          </span>
          <span v-if="getFieldValue(field)">
            <at-text :options="valueOptions" :defaultSelectedList="getFieldTags(field)" :defaultValue="getFieldValue(field).toString()"/>
          </span>
          <span v-else></span>
        </a-tag>
      </div>
      <div v-else>
        <a-tag>--</a-tag>
      </div>
    </div>
  </div>
</template>

<script>
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
import {feiShuOperatorMap} from "@/constants/feishu-table.js";
import AtText from "@/views/workflow/components/at-input/at-text.vue";

export default {
  name: "feishu-table-render",
  components: {AtText, QuestionCircleOutlined},
  props: {
    nodeParams: {
      type: Object
    },
    valueOptions: {
      type: Array
    },
  },
  data() {
    return {
      operatorMap: feiShuOperatorMap,
      actionState: {},
    }
  },
  mounted() {
    this.dataFormat()
  },
  watch: {
    nodeParams: {
      handler(newVal, oldVal) {
        this.dataFormat()
        this.$emit('updateSize')
      },
      deep: true
    }
  },
  methods: {
    dataFormat() {
      let params = this.nodeParams?.plugin?.params || {}
      let arg = params?.arguments || {}
      let {config_name, table_name} = params || {}
      this.actionState = {
        ...params,
        ...arg,
        config_name,
        table_name,
      }
      console.log('this.actionState', this.actionState)
    },
    getFieldLabel(field) {
      return field?.properties?.[field?.current_properties_key]?.name || field?.name
    },
    getFieldValue(field) {
      if (field?.media_component) {
        const p = field?.properties?.[field?.current_properties_key] || {}
        return p?.value ?? ''
      }
      return field?.value ?? ''
    },
    getFieldTags(field) {
      if (field?.media_component) {
        const p = field?.properties?.[field?.current_properties_key] || {}
        return p?.atTags ?? []
      }
      return field?.atTags ?? []
    }
  }
}
</script>

<style scoped lang="less">
._main {
  display: flex;
  flex-direction: column;
  gap: 8px;
  > div {
    display: flex;

    > .label {
      flex-shrink: 0;
      width: 100px;
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

.hint {
  position: relative;
  cursor: pointer;
}
.hint-pop {
  position: absolute;
  bottom: 30px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 12px;
  opacity: 0;
  pointer-events: none;
  min-width: max-content;
  min-height: 32px;
  padding: 6px 8px;
  color: #fff;
  text-align: start;
  text-decoration: none;
  word-wrap: break-word;
  background-color: rgba(0, 0, 0, 0.85);
  border-radius: 6px;
  box-shadow: 0 6px 16px 0 rgba(0, 0, 0, 0.08), 0 3px 6px -4px rgba(0, 0, 0, 0.12), 0 9px 28px 8px rgba(0, 0, 0, 0.05);

  &::after {
    content: '';
    position: absolute;
    bottom: -6px;
    left: 50%;
    transform: translateX(-50%);
    width: 0;
    height: 0;
    border-left: 8px solid transparent;
    border-right: 8px solid transparent;
    border-top: 8px solid rgba(0, 0, 0, 0.85);
  }
}
.hint:hover .hint-pop {
  opacity: 1;
  pointer-events: auto;
}
</style>
