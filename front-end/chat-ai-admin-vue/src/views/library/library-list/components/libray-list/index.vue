<style lang="less" scoped>
.list-box {
  display: flex;
  flex-flow: row wrap;
  margin: 0 -8px;
}
.list-item-wrapper {
  padding: 8px;
  width: 25%;
}
.list-item {
  position: relative;
  width: 100%;
  height: 200px;
  padding: 16px;
  border-radius: 6px;
  border: 1px solid #f0f0f0;
  background-color: #fff;
  transition: all 0.25s;
  cursor: pointer;

  &:hover {
    box-shadow: 0 4px 16px 0 #1b3a6929;
  }

  &::after {
    content: '';
    display: block;
    position: absolute;
    bottom: 0;
    right: 0;
    width: 80px;
    height: 80px;
    opacity: 0.5;
    background: url('../../../assets/img/library/library_item_bg.svg') 0 0 no-repeat;
    background-size: cover;
  }

  .item-header {
    position: relative;
    display: flex;
    align-items: center;

    .item-action {
      .menu-btn {
        width: 22px;
        height: 22px;
        text-align: center;
        line-height: 22px;
        font-size: 16px;
        cursor: pointer;
        &:hover {
          color: #2475fc;
        }
      }
    }

    .library-icon {
      width: 40px;
      height: 40px;
      border-radius: 6px;
      margin-right: 12px;
    }

    .library-title {
      flex: 1;

      height: 22px;
      line-height: 22px;
      font-size: 14px;
      font-weight: 600;
      color: #262626;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }

  .item-body {
    margin-top: 16px;
  }

  .library-desc {
    max-height: 44px;
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: #595959;
    // 超出2行显示省略号
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
  }
  .library-size {
    display: flex;
    margin-top: 16px;
    line-height: 22px;
    font-size: 14px;
    color: #7a8699;

    .file-size {
      padding-left: 20px;
    }
  }

  .library-type {
    margin-bottom: 8px;
    line-height: 1;
    font-size: 0;
    .type-tag {
      display: inline-block;
      line-height: 22px;
      padding: 0 8px;
      border-radius: 6px;
      font-size: 14px;
      border: 1px solid #99bffd;
      color: #2475fc;
    }
  }
}
.add-library {
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 22px;
  color: #3a4559;
  cursor: pointer;

  .add-library-icon {
    font-size: 16px;
  }
  .add-library-text {
    padding-left: 4px;
    font-size: 14px;
  }
}
// 大于1440px
@media screen and (min-width: 1440px) {
  .list-box {
    .list-item-wrapper {
      width: 20%;
    }
  }
}
</style>

<template>
  <div class="list-box">
    <div class="list-item-wrapper" v-if="props.showCreate">
      <div class="list-item add-library" @click="toAdd">
        <PlusCircleOutlined class="add-library-icon" />
        <span class="add-library-text">新增知识库</span>
      </div>
    </div>
    <div class="list-item-wrapper" v-for="item in props.list" :key="item.id">
      <div class="list-item" @click.stop="toEdit(item)">
        <div class="item-header">
          <img class="library-icon" :src="item.avatar" alt="" />
          <div class="library-title">{{ item.library_name }}</div>
          <span class="item-action" @click.stop>
            <a-dropdown>
              <span class="menu-btn" @click.prevent>
                <MoreOutlined />
              </span>
              <template #overlay>
                <a-menu>
                  <a-menu-item>
                    <a href="javascript:;" @click.stop="toEdit(item)">管 理</a>
                  </a-menu-item>
                  <a-menu-item>
                    <a class="delete-text-color" href="javascript:;" @click="handleDelete(item)"
                      >删 除</a
                    >
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </span>
        </div>
        <div class="item-body">
          <div class="library-type">
            <span class="type-tag" v-if="item.type == 0">普通知识库</span>
            <span class="type-tag" v-if="item.type == 1">对外知识库</span>
            <span class="type-tag" v-if="item.type == 2">问答知识库</span>
          </div>
          <div class="library-desc">{{ item.library_intro }}</div>
        </div>
        <div class="library-size">
          <span>文档数：{{ item.file_total }}</span>
          <span class="file-size">文档大小：{{ item.file_size_str }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { MoreOutlined, PlusCircleOutlined } from '@ant-design/icons-vue'

const emit = defineEmits(['add', 'edit', 'delete'])

const props = defineProps({
  list: {
    type: Array,
    default: () => []
  },
  showCreate: {
    type: Boolean,
    default: true
  }
})

const toEdit = (item) => {
  emit('edit', item)
}

const toAdd = (item) => {
  emit('add', item)
}

const handleDelete = (item) => {
  emit('delete', item)
}
</script>
