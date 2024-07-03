<template>
  <div bg-color="#f5f9ff" class="library-page">
    <div class="library-page-body">
      <div class="list-box">
        <div class="list-item-wrapper">
          <div class="list-item add-library" @click="toAdd">
            <PlusCircleOutlined class="add-library-icon" />
            <span class="add-library-text">新增知识库</span>
          </div>
        </div>
        <div class="list-item-wrapper" v-for="item in list" :key="item.id">
          <div class="list-item" @click.stop="toEdit(item)">
            <div class="item-header">
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
            <div class="library-desc">{{ item.library_intro }}</div>
            <div class="library-size">
              <span>文档数：{{ item.file_total }}</span>
              <span class="file-size">文档大小：{{ formatFileSize(item.file_size) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, createVNode } from 'vue'
import { useRouter } from 'vue-router'
import { Modal, message } from 'ant-design-vue'
import { MoreOutlined, PlusCircleOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { getLibraryList, deleteLibrary } from '@/api/library'
import { formatFileSize } from '@/utils/index'
const router = useRouter()

const list = ref([])

const getList = () => {
  getLibraryList({}).then((res) => {
    list.value = res.data || []
  })
}

getList()

const toAdd = () => {
  router.push({
    name: 'addLibrary'
  })
}

const toEdit = (data) => {
  router.push({
    name: 'libraryDetails',
    query: {
      id: data.id
    }
  })
}

const handleDelete = (data) => {
  let secondsToGo = 3

  let modal = Modal.confirm({
    title: `删除${data.library_name}`,
    icon: createVNode(ExclamationCircleOutlined),
    content: '您确定要删除此知识库吗？',
    okText: secondsToGo + ' 确 定',
    okType: 'danger',
    okButtonProps: {
      disabled: true
    },
    cancelText: '取 消',
    onOk() {
      onDelete(data)
    },
    onCancel() {
      // console.log('Cancel')
    }
  })

  let interval = setInterval(() => {
    if (secondsToGo == 1) {
      modal.update({
        okText: '确 定',
        okButtonProps: {
          disabled: false
        }
      })

      clearInterval(interval)
      interval = undefined
    } else {
      secondsToGo -= 1

      modal.update({
        okText: secondsToGo + ' 确 定',
        okButtonProps: {
          disabled: true
        }
      })
    }
  }, 1000)
}

const onDelete = ({ id }) => {
  deleteLibrary({ id }).then(() => {
    message.success('删除成功')
    getList()
  })
}
</script>

<style lang="less" scoped>
.library-page {
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
    height: 136px;
    padding: 16px;
    border-radius: 2px;
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

      .library-title {
        height: 22px;
        line-height: 22px;
        font-size: 14px;
        font-weight: 600;
        color: #262626;
      }
    }

    .library-desc {
      height: 44px;
      line-height: 22px;
      margin-top: 4px;
      font-size: 14px;
      font-weight: 400;
      color: #8c8c8c;
      // 超出2行显示省略号
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
    }
    .library-size {
      display: flex;
      margin-top: 12px;
      line-height: 22px;
      font-size: 14px;
      color: #595959;

      .file-size {
        padding-left: 20px;
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
}
// 大于1440px
@media screen and (min-width: 1440px) {
  .library-page {
    .list-box {
      .list-item-wrapper {
        width: 20%;
      }
    }
  }
}
</style>
