<template>
  <div class="_main">
    <div>
      <span class="label">配置：</span>
      <a-tag>{{ actionState?.config_name || '--' }}</a-tag>
    </div>
    <div v-if="!['create_bitable', 'create_tables', 'table_add_members', 'update_advanced', 'roles_list'].includes(business)">
      <span class="label">数据表：</span>
      <a-tag>
        <at-text
          v-if="actionState?.tag_map?.table_id && actionState?.tag_map?.table_id.length"
          class="batch-render-at"
          :options="valueOptions"
          :defaultSelectedList="actionState.tag_map.table_id"
          :defaultValue="actionState?.table_id"
        />
        <template v-else>{{actionState?.table_name || actionState?.table_id || '--'}}</template>
      </a-tag>
    </div>
    <div v-if="business == 'create_bitable'">
      <span class="label">多维表名称：</span>
      <a-tag>
        <at-text
          class="batch-render-at"
          :options="valueOptions"
          :defaultSelectedList="actionState?.tag_map?.name || []"
          :defaultValue="actionState?.name || '--'"
        />
      </a-tag>
    </div>
    <template v-else-if="business == 'create_tables'">
      <div>
        <span class="label">数据表名称：</span>
        <a-tag>
          <at-text
            class="batch-render-at"
            :options="valueOptions"
            :defaultSelectedList="actionState?.tag_map?.name || []"
            :defaultValue="actionState?.name || '--'"
          />
        </a-tag>
      </div>
    </template>
    <template v-else-if="business == 'create_views'">
      <div>
        <span class="label">视图名称：</span>
        <a-tag>
          <at-text
            class="batch-render-at"
            :options="valueOptions"
            :defaultSelectedList="actionState?.tag_map?.view_name || []"
            :defaultValue="actionState?.view_name || '--'"
          />
        </a-tag>
      </div>
    </template>
    <template v-else-if="business == 'create_record'">
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
    <template v-else-if="business == 'update_advanced'">
      <div>
        <span class="label">多维表名称：</span>
        <a-tag v-if="actionState.name">
          <at-text :options="valueOptions" :defaultSelectedList="actionState?.tag_map?.name || []" :defaultValue="actionState.name"/>
        </a-tag>
        <a-tag v-else>--</a-tag>
      </div>
      <div>
        <span class="label">高级权限：</span>
        <a-tag>{{actionState.is_advanced == 0 ? '不设置' : actionState.is_advanced == 1 ? '开启' : '关闭'}}</a-tag>
      </div>
    </template>
    <template v-else-if="business == 'table_add_members'">
      <div>
        <span class="label">多维表：</span>
        <a-tag>
          <at-text
            class="batch-render-at"
            :options="valueOptions"
            :defaultSelectedList="actionState?.tag_map?.app_token || []"
            :defaultValue="actionState?.app_token || '--'"
          />
        </a-tag>
      </div>
      <div>
        <span class="label">协作者角色：</span>
        <a-tag v-if="actionState.role_id">
          <at-text :options="valueOptions" :defaultSelectedList="actionState?.tag_map?.role_id || []" :defaultValue="actionState.role_id"/>
        </a-tag>
        <a-tag v-else>--</a-tag>
      </div>
      <div>
        <span class="label">协作者ID：</span>
        <a-tag v-if="actionState.member_id">
          <at-text :options="valueOptions" :defaultSelectedList="actionState?.tag_map?.member_id || []" :defaultValue="actionState.member_id"/>
        </a-tag>
        <a-tag v-else>--</a-tag>
      </div>
    </template>
    <template v-else-if="business == 'roles_list'">
      <div>
        <span class="label">多维表：</span>
        <a-tag>
          <at-text
            class="batch-render-at"
            :options="valueOptions"
            :defaultSelectedList="actionState?.tag_map?.app_token || []"
            :defaultValue="actionState?.app_token || '--'"
          />
        </a-tag>
      </div>
    </template>
    <template v-else-if="BatchActions.includes(business)">
      <div v-for="(item, key) in batchActionParam" :key="key">
        <span class="label">{{ item.name }}：</span>
        <div v-if="actionState[key]" class="field-box">
          <a-tag>
            <at-text
              class="batch-render-at"
              :options="valueOptions"
              :defaultSelectedList="actionState?.tag_map[key] || []"
              :defaultValue="actionState[key].toString()"
            />
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
import {BatchActions, feiShuOperatorMap, getBatchActionParams} from "@/constants/feishu-table.js";
import AtText from "@/views/workflow/components/at-input/at-text.vue";
import {runPlugin} from "@/api/plugins/index.js";

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
      BatchActions,
      batchActionParam: {},
      operatorMap: feiShuOperatorMap,
      actionState: {},
    }
  },
  async mounted() {
    if (BatchActions.includes(this.business)) {
      await this.initParams()
    }
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
    async initParams() {
      await runPlugin({
        name: 'feishu_bitable',
        action: "default/get-schema",
        params: {}
      }).then(res => {
        let data = res?.data || {}
        this.batchActionParam = getBatchActionParams(data[this.business]?.params || {})
      })
    },
    dataFormat() {
      let arg = this.nodeParams?.plugin?.params?.arguments || {}
      let {config_name, table_name} = this.nodeParams?.plugin?.params || {}
      this.actionState = {
        ...arg,
        config_name,
        table_name,
      }
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

.batch-render-at {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: pre-wrap;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  line-clamp: 2;
}
</style>
