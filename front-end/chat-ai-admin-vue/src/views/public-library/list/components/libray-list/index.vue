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
  padding: 24px;
  border: 1px solid #E4E6EB;
  border-radius: 12px;
  background-color: #fff;
  transition: all 0.25s;
  cursor: pointer;

  &:hover {
    box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.12);
  }

  .library-info {
    position: relative;
    display: flex;
    align-items: center;
    .item-action {
      .menu-btn {
        position: absolute;
        right: 0;
        top: 0;
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
    .library-icon{
      width: 52px;
      height: 52px;
      border-radius: 14px;
      overflow: hidden;
    }
    .library-info-content{
      margin-left: 12px;
      flex: 1;
      overflow: hidden;
    }
    .library-title {
      height: 24px;
      line-height: 24px;
      font-size: 16px;
      font-weight: 600;
      color: #262626;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }
  .item-body{
    margin-top: 12px;
  }
  .library-desc {
    height: 44px;
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: rgb(89, 89, 89);
    // 超出2行显示省略号
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
  }
  
  .item-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 14px;
    color: #7a8699;
  }
  .library-size {
    display: flex;
    line-height: 20px;
    font-size: 12px;
    font-weight: 400;
    color: #7a8699;

    .text-item {
      margin-right: 12px;
      &:last-child{
        margin-right: 0;
      }
    }
  }

  .action-box {
    font-size: 14px;
    height: 24px;
    color: #2475fc;
    display: flex;
    align-items: center;

    .action-item {
      display: flex;
      align-items: center;
      height: 100%;
      padding: 4px;
      border-radius: 6px;
      cursor: pointer;
      color: #595959;
      transition: all 0.2s;
    }
    .action-item:hover {
      background: #E4E6EB;
    }

    .action-icon {
      font-size: 16px;
    }
  }
}

// 大于1920px
@media screen and (min-width: 1920px) {
  .list-box {
    .list-item-wrapper {
      width: 20%;
    }
  }
}
</style>

<template>
  <div class="list-box">
    <div class="list-item-wrapper" v-for="item in props.list" :key="item.id">
      <div class="list-item" @click.stop="toEdit(item)">
        <div class="library-info">
          <img class="library-icon" :src="item.avatar" alt="" />
          <div class="library-info-content">
            <div class="library-title">{{ item.library_name }}</div>
          </div>
        </div>
        <div class="item-body">
          <a-tooltip :title="getTooltipTitle(item.library_intro, item, 14, 2, 0)" placement="top">
            <div class="library-desc" :ref="el => setDescRef(el, item)">{{ item.library_intro }}</div>
          </a-tooltip>
        </div>

        <div class="item-footer">
          <div class="library-size">
            <span class="text-item">文档：{{ item.file_total }}</span>
            <span class="text-item">大小：{{ item.file_size_str }}</span>
          </div>

          <div class="action-box" @click.stop>
            <a-dropdown>
              <div class="action-item" @click.stop>
                <svg-icon class="action-icon" name="point-h"></svg-icon>
              </div>
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
          </div>

        </div>
      </div>
    </div>
  </div>

</template>

<script setup>
import { setDescRef, getTooltipTitle } from '@/utils/index'

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

const handleDelete = (item) => {
  emit('delete', item)
}
</script>
