<style lang="less" scoped>
.message-node {
  position: relative;
  .node-desc {
    line-height: 22px;
    margin-bottom: 4px;
    font-size: 14px;
    font-weight: 400;
    color: var(--wf-color-text-2);
  }
  .message-item {
    margin-top: 12px;
    &:first-child {
      margin-top: 0;
    }
  }
  .q-title {
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: #262626;
  }
  .q-image {
    width: 48px;
    height: 48px;

    & > img {
      display: block;
      width: 100%;
      height: 100%;
    }
  }
  .q-btn-box {
    margin-top: 12px;
    .q-btn {
      position: relative;
      line-height: 22px;
      padding: 8px 12px;
      margin-top: 12px;
      border-radius: 4px;
      font-size: 14px;
      color: #595959;
      background-color: #f2f4f7;

      &:first-child {
        margin-top: 0;
      }
    }
  }

  .no-data{
    height: 38px;
    line-height: 22px;
    font-size: 14px;
    color: var(--wf-color-text-2);
  }
}
</style>

<template>
  <node-common
    :title="properties.node_name"
    :menus="menus"
    :icon-name="properties.node_icon_name"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    @handleMenu="handleMenu"
  >
    <div class="message-node">
      <div class="no-data" v-if="properties.message_list.length == 0">点击设置消息内容</div>
      <div class="message-item" v-for="(item, index) in properties.message_list" :key="index">
        <div class="node-desc">延时{{ item.delay }}s</div>
        <div class="q-title" v-if="item.msg_type == 'text' || item.msg_type == 'menu'">
          {{ item.sort_content }}
        </div>
        <div class="q-image" v-if="item.msg_type == 'image'">
          <img class="image" :src="item.content" alt="" />
        </div>

        <div class="q-btn-box" v-if="item.msg_type == 'menu' && item.question.length">
          <div class="q-btn" v-for="(option, index) in item.question" :key="index">
            {{ option.content }}
          </div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script>
import NodeCommon from '../base-node.vue'

export default {
  name: 'MessageNode',
  components: {
    NodeCommon,
  },
  inject: ['getNode', 'getGraph'],
  props: {
    properties: {
      type: Object,
      default() {
        return {
          message_list: [],
        }
      },
    },
    isSelected: { type: Boolean, default: false },
    isHovered: { type: Boolean, default: false },
  },
  data() {
    return {
      menus: [{ name: '删除', key: 'delete', color: '#fb363f' }],
    }
  },
  methods: {
    handleMenu(item) {
      if (item.key === 'delete') {
        let node = this.getNode()

        this.getGraph().deleteNode(node.id)
      }
    },
  },
}
</script>
