<style lang="less" scoped>
.list-box {
  display: flex;
  flex-flow: row wrap;
  margin: 0 -8px 0 -8px;
  flex: 1;
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

  .library-type {
    position: absolute;
    top: 0;
    right: 0;
    display: flex;
    padding: 1px 8px;
    flex-direction: column;
    justify-content: center;
    align-items: flex-start;
    gap: 10px;
    position: absolute;
    right: 0;
    border-radius: 0 8px;
    background: var(--09, #F2F4F7);

    .type-box {
      color: #595959;
      font-size: 12px;
      font-style: normal;
      font-weight: 400;
      line-height: 16px;
    }
  }

  .library-info {
    position: relative;
    display: flex;
    align-items: center;

    .library-icon {
      width: 48px;
      height: 48px;
      overflow: hidden;
      position: relative;
      > img {
        width: 100%;
        height: 100%;
        border-radius: 12px;
      }

      .sync-tag {
        position: absolute;
        left: 0;
        bottom: 0;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 100%;
        height: 18px;
        background: rgba(0,0,0,0.5);
        color: #FFF;
        font-size: 10px;
        border-radius: 0 0 14px 14px;
        &.run {
          background: #2475FC;
        }
        &.fail {
          background: #FB363F;
        }
      }
    }

    .library-info-content{
      flex: 1;
      padding-left: 12px;
      overflow: hidden;
    }
  }

  .library-title {
    height: 24px;
    line-height: 24px;
    margin-bottom: 4px;
    font-size: 16px;
    font-weight: 600;
    color: #262626;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  // .library-type {
  //   display: flex;
  //   .type-tag {
  //     height: 22px;
  //     line-height: 20px;
  //     padding: 0 8px;
  //     font-size: 12px;
  //     font-weight: 400;
  //     border-radius: 6px;
  //     color: rgb(36, 117, 252);
  //     border: 1px solid #CDE0FF;
  //   }
  //   .graph-tag {
  //     margin-left: 4px;
  //     &.gray-tag {
  //       border: 1px solid #00000026;
  //       background: #0000000a;
  //       color: #bfbfbf;
  //     }
  //   }
  // }
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
      display: flex;
      align-items: center;
      margin-right: 12px;
      &:last-child{
        margin-right: 0;
      }

      .icon {
        font-size: 14px;
        margin-right: 4px;
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

// 大于1920px
@media screen and (max-width: 1500px) {
  .list-box {
    .list-item-wrapper {
      width: 33.3%;
    }
  }
}
</style>

<template>
  <div class="list-box">
    <div v-if="props.list.length" class="list-item-wrapper" v-for="item in props.list" :key="item.id">
      <div class="list-item" @click.stop="toEdit(item)">
        <div class="library-type">
          <span class="type-box" v-if="item.type == 0">{{ t('normal_knowledge_base') }}</span>
          <span class="type-box" v-if="item.type == 1">{{ t('external_knowledge_base') }}</span>
          <span class="type-box" v-if="item.type == 2">{{ t('qa_knowledge_base') }}</span>
          <span class="type-box" v-if="item.type == 3">{{ t('official_account_knowledge_base') }}</span>
        </div>
        <div class="library-info">
          <div class="library-icon">
            <img :src="item.avatar"/>
            <template v-if="item.type == 3">
              <span v-if="item.sync_official_content_status == 2" class="sync-tag run">{{ t('syncing') }}</span>
              <span v-else-if="item.sync_official_content_status == 3" class="sync-tag fail">{{ t('sync_failed') }}</span>
            </template>
          </div>
          <div class="library-info-content">
            <div class="library-title">{{ item.library_name }}</div>
            <!-- <div class="library-type">
              <span class="type-tag" v-if="item.type == 0">{{ t('normal_knowledge_base') }}</span>
              <span class="type-tag" v-if="item.type == 1">{{ t('external_knowledge_base') }}</span>
              <span class="type-tag" v-if="item.type == 2">{{ t('qa_knowledge_base') }}</span>
              <span class="type-tag" v-if="item.type == 3">{{ t('official_account_knowledge_base') }}</span> -->
              <!-- <a-tooltip v-if="neo4j_status">
                <template #title>{{ item.graph_switch == 0 ? '未' : '已' }}开启知识图谱生成</template>
                <span class="type-tag graph-tag" :class="{ 'gray-tag': item.graph_switch == 0 }"
                  >Graph</span
                >
              </a-tooltip> -->
            <!-- </div> -->
          </div>
        </div>
        <div class="item-body">
          <a-tooltip :title="getTooltipTitle(item.library_intro, item, 14, 2, 0)" placement="top">
            <div class="library-desc" :ref="el => setDescRef(el, item)">{{ item.library_intro }}</div>
          </a-tooltip>
        </div>

        <div class="item-footer">
          <div class="library-size">
            <a-tooltip title="文档数量">
              <span class="text-item">
                <svg-icon class="icon" name="document-icon"></svg-icon>
                {{ item.file_total }}
              </span>
            </a-tooltip>
            <a-tooltip title="存储大小">
              <span class="text-item">
              <svg-icon class="icon" name="storage-icon"></svg-icon>
              {{ item.file_size_str }}</span>
            </a-tooltip>
            <a-tooltip title="关联应用">
              <span class="text-item">
                <svg-icon class="icon" name="relevance-icon"></svg-icon>
                {{ item.robot_nums || 0 }}
              </span>
            </a-tooltip>
          </div>

          <div class="action-box" @click.stop>
            <a-dropdown>
              <div class="action-item" @click.stop>
                <svg-icon class="action-icon" name="point-h"></svg-icon>
              </div>
              <template #overlay>
                <a-menu>
                  <a-menu-item>
                    <a class="delete-text-color" href="javascript:;" @click.stop="handleDelete(item)"
                      >{{ t('delete') }}</a
                    >
                  </a-menu-item>
                  <a-menu-item>
                    <div @click.stop="openEditGroupModal(item)">{{ t('modify_group') }}</div>
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
        </div>


      </div>
    </div>
    <EmptyBox v-else style="margin-top: 80px;"/>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { setDescRef, getTooltipTitle } from '@/utils/index'
import { useI18n } from '@/hooks/web/useI18n'
const emit = defineEmits(['add', 'edit', 'delete', 'openEditGroupModal'])

import { useCompanyStore } from '@/stores/modules/company'
import EmptyBox from "@/components/common/empty-box.vue";
const { t } = useI18n('views.library.library-list.components.library-list.index')

const companyStore = useCompanyStore()
const neo4j_status = computed(()=>{
  return companyStore.companyInfo?.neo4j_status == 'true'
})

const props = defineProps({
  list: {
    type: Array,
    default: () => []
  },
})

const toEdit = (item) => {
  emit('edit', item)
}

const handleDelete = (item) => {
  emit('delete', item)
}

const openEditGroupModal = (item)=>{
  emit('openEditGroupModal', item)
}
</script>
