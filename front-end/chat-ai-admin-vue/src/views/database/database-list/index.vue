<template>
  <div class="library-page">
    <div class="library-page-body">
      <div class="list-box">
        <div class="list-item-wrapper">
          <div class="list-item add-library" @click="toAdd()">
            <PlusCircleOutlined class="add-library-icon" />
            <span class="add-library-text">新增数据表</span>
          </div>
        </div>
        <div class="list-item-wrapper" v-for="item in list" :key="item.id">
          <div class="list-item" @click.stop="toEdit(item)">
            <div class="item-header">
              <img class="library-icon" src="@/assets/img/database/base-icon.png" alt="" />
              <div class="library-info">
                <div class="library-title">{{ item.name }}</div>
                <div class="library-num">
                  <span>数据：{{ item.entry_count }}条</span>
                </div>
              </div>
              <span class="item-action" @click.stop>
                <a-dropdown>
                  <span class="menu-btn" @click.prevent>
                    <MoreOutlined />
                  </span>
                  <template #overlay>
                    <a-menu>
                      <a-menu-item>
                        <a href="javascript:;" @click.stop="toAdd(item)">编 辑</a>
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
            <div class="library-desc">{{ item.description }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <AddDataSheet @ok="getList" ref="addDataSheetRef"></AddDataSheet>
</template>

<script setup>
import { ref, createVNode } from 'vue'
import { useRouter } from 'vue-router'
import { Modal, message } from 'ant-design-vue'
import { MoreOutlined, PlusCircleOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { getFormList, delForm } from '@/api/database'
import AddDataSheet from './components/add-data-sheet.vue'
const router = useRouter()

const list = ref([])

const getList = () => {
  getFormList({}).then((res) => {
    list.value = res.data || []
  })
}

getList()

const addDataSheetRef = ref(null)
const toAdd = (data = {}) => {
  addDataSheetRef.value.show(data)
}

const toEdit = (data) => {
  router.push({
    path: '/database/details',
    query: {
      form_id: data.id,
      name: data.name
    }
  })
}

const handleDelete = (data) => {
  let secondsToGo = 3

  let modal = Modal.confirm({
    title: `删除确认`,
    icon: createVNode(ExclamationCircleOutlined),
    content: `删除数据表后，表中所有数据将一并被删除，不可恢复。确认删除数据表${data.name}吗?`,
    okText: secondsToGo + ' 确 定',
    okType: 'danger',
    okButtonProps: {
      disabled: true
    },
    cancelText: '取 消',
    onOk() {
      onDelete(data)
    },
    onCancel() {}
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
  delForm({ id }).then(() => {
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
    height: 158px;
    padding: 16px;
    border-radius: 2px;
    border: 1px solid #f0f0f0;
    background-color: #fff;
    transition: all 0.25s;
    cursor: pointer;

    &:hover {
      box-shadow: 0 4px 16px 0 #1b3a6929;
    }

    .item-header {
      position: relative;
      display: flex;
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
        width: 40px;
        height: 40px;
        margin-right: 8px;

      }
      .library-info{
        width: calc(100% - 68px);
      }
      .library-title {
        height: 22px;
        line-height: 22px;
        font-size: 14px;
        font-weight: 600;
        color: #262626;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
      .library-num{
        color: #8c8c8c;
        font-size: 14px;
        line-height: 22px;
        font-weight: 400;
      }
    }

    .library-desc {
      height: 66px;
      line-height: 22px;
      margin-top: 12px;
      font-size: 14px;
      font-weight: 400;
      color: #595959;
      // 超出2行显示省略号
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 3;
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
