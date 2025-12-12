<template>
  <div class="_main">
    <div>
      <span class="label">配置：</span>
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
            <span v-if="field.value">
              <at-text :options="valueOptions" :defaultSelectedList="field.atTags" :defaultValue="field.value.toString()"/>
            </span>
            <span v-else></span>
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
        <a-tag v-if="actionState.record_id">
          <at-text :options="valueOptions" :defaultSelectedList="actionState.record_tags" :defaultValue="actionState.record_id"/>
        </a-tag>
        <a-tag v-else>--</a-tag>
      </div>
    </template>
    <template v-else-if="business == 'search_records'">
      <div>
        <span class="label">查询条件：</span>
        <div v-if="Array.isArray(actionState?.filter?.conditions) && actionState.filter.conditions.length" class="field-box">
          <a-tag v-for="item in actionState.filter.conditions" class="field-tag" :key="item.field_id">
            <span>{{ item.field_name }}</span>
            <span class="arrow-icon">{{ operatorMap[item.operator]?.label }}</span>
            <span v-if="item.value">
              <at-text :options="valueOptions" :defaultSelectedList="item.atTags" :defaultValue="item.value.toString()"/>
            </span>
            <span v-else></span>
          </a-tag>
        </div>
        <div v-else>
          <a-tag>--</a-tag>
        </div>
      </div>
      <div>
        <span class="label">查询字段：</span>
        <div v-if="Array.isArray(actionState.field_names) && actionState.field_names.length" class="field-box">
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
        <a-tag v-if="actionState.record_id">
          <at-text :options="valueOptions" :defaultSelectedList="actionState.record_tags" :defaultValue="actionState.record_id"/>
        </a-tag>
        <a-tag v-else>--</a-tag>
      </div>
      <div>
        <span class="label">更新数据：</span>
        <div v-if="Array.isArray(actionState.fields)" class="field-box">
          <a-tag v-for="field in actionState.fields" class="field-tag" :key="field.field_name">
            <span>{{ field.field_name }}</span>
            <span class="arrow-icon">
              <img src="@/assets/img/workflow/arrow-right.svg"/>
            </span>
            <span v-if="field.value">
              <at-text :options="valueOptions" :defaultSelectedList="field.atTags" :defaultValue="field.value.toString()"/>
            </span>
            <span v-else></span>
          </a-tag>
        </div>
        <div v-else>
          <a-tag>--</a-tag>
        </div>
      </div>
    </template>
  </div>
</template>

<script>
import {FeiShuOperatorMap} from "@/constants/feishu-table.js";
import AtText from "@/views/workflow/components/at-input/at-text.vue";

export default {
  name: "feishu-table-render",
  components: {AtText},
  props: {
    business: {
      type: String,
    },
    nodeParams: {
      type: Object
    },
    valueOptions: {
      type: Array
    },
  },
  data() {
    return {
      operatorMap: FeiShuOperatorMap(),
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
      let arg = this.nodeParams?.plugin?.params?.arguments || {}
      let {config_name, table_name} = this.nodeParams?.plugin?.params || {}
      this.actionState = {
        ...arg,
        config_name,
        table_name,
      }
      console.log(arg)
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
