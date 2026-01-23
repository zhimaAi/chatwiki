<template>
  <a-modal
    v-model:open="visible"
    title="关联公众号"
    :confirm-loading="saving"
    width="626px"
    @ok="save"
  >
    <LoadingBox v-if="loading" style="margin: 10vh 0;"/>
    <template v-else-if="list.length">
      <div class="action-box">
        <a-button type="primary" ghost :icon="h(PlusOutlined)" @click="linkOfficial">添加公众号</a-button>
      </div>
      <div class="wechat-app-list">
        <div v-for="item in list"
             :key="item.id"
             class="wechat-app-item">
          <div class="item-body">
            <img class="app-avatar" :src="item.app_avatar" alt=""/>
            <div class="app-info">
              <div class="app-name">{{ item.app_name }}</div>
              <div class="app-desc">Appid：{{ item.app_id }}</div>
            </div>
          </div>
          <div class="ext-info-list">
            <div class="status-block status-success" v-if="item.account_is_verify == 'true'">
              <CheckCircleFilled/>
              已认证
            </div>
            <div v-else class="status-block status-warning">
              <ExclamationCircleFilled/>
              未认证
            </div>
            <a-checkbox v-model:checked="item.checked"/>
          </div>
        </div>
      </div>
    </template>
    <EmptyBox v-else title="暂未绑定公众号">
      <template #desc>
        <a-button type="primary" @click="linkOfficial">绑定公众号</a-button>
      </template>
    </EmptyBox>
  </a-modal>
</template>
<script setup>
import {ref, h, computed} from 'vue';
import {useRouter} from 'vue-router';
import {message} from 'ant-design-vue';
import {PlusOutlined, CheckCircleFilled, ExclamationCircleFilled} from '@ant-design/icons-vue';
import LoadingBox from "@/components/common/loading-box.vue";
import EmptyBox from "@/components/common/empty-box.vue";
import {getWechatAppList, robotBindWxApp} from "@/api/robot/index.js";
import {usePermissionStore} from "@/stores/modules/permission.js";

const emit = defineEmits(['change'])
const router = useRouter()
const permissionStore = usePermissionStore()

const visible = ref(false)
const loading = ref(false)
const saving = ref(false)
const robotInfo = ref({})
const list = ref([])

const officialPerm = computed(() => {
  let {role_permission, role_type} = permissionStore
  return role_type == 1 || role_permission.includes('OfficialAccountManage')
})

function open(r, ids) {
  robotInfo.value = r
  loadData(ids)
  visible.value = true
}

function loadData(ids) {
  ids = Array.isArray(ids) ? ids : []
  loading.value = true
  getWechatAppList({
    app_type: 'official_account',
    app_name: ''
  }).then((res) => {
    let _list = res.data || []
    _list.forEach(item => {
      item.checked = ids.includes(item.app_id)
    })
    list.value = _list
  }).finally(() => {
    loading.value = false
  })
}

function save() {
  saving.value = true
  let row = list.value.filter(item => item.checked)
  let ids = row.map(i => i.app_id).toString()
  robotBindWxApp({
    robot_id: robotInfo.value.id,
    app_id_list: ids
  }).then(() => {
    visible.value = false
    emit('change', row, ids)
  }).finally(() => {
    saving.value = false
  })
}

function linkOfficial() {
  const url = router.resolve({path: '/user/official-account'})
  window.open(url.href, '_blank')
}

defineExpose({
  open
})
</script>

<style scoped lang="less">
.action-box {
  text-align: right;
  margin-bottom: 16px;
}

.wechat-app-list {
  display: flex;
  gap: 16px;
  flex-flow: row wrap;

  .wechat-app-item {
    width: 280px;
    height: fit-content;
    padding: 8px;
    border-radius: 8px;
    border: 1px solid #e4e6eb;
    overflow: hidden;
    cursor: pointer;
    position: relative;

    .app-info-block {
      display: flex;
      align-items: center;
    }

    &:hover {
      box-shadow: 0 3px 14px 2px #0000000d,
      0 8px 10px 1px #0000000f,
      0 5px 5px -3px #0000001a;
    }

    .item-body {
      display: flex;
      align-items: center;
      overflow: hidden;
    }

    .app-info {
      flex: 1;
      overflow: hidden;
    }

    .app-avatar {
      display: block;
      width: 32px;
      height: 32px;
      border-radius: 16px;
      margin-right: 12px;
      border: 1px solid #f0f0f0;
    }

    .app-name,
    .app-desc {
      width: 100%;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }

    .app-name {
      line-height: 24px;
      font-size: 16px;
      font-weight: 600;
      color: #262626;
    }

    .app-desc {
      line-height: 20px;
      font-size: 12px;
      font-weight: 400;
      color: #8c8c8c;
      margin-top: 4px;
    }

    .ext-info-list {
      margin-top: 12px;
      display: flex;
      flex-wrap: wrap;
      color: #595959;
      line-height: 22px;
      font-size: 14px;
      display: flex;
      align-items: center;
      justify-content: space-between;

      .status-block {
        display: flex;
        align-items: center;
        border-radius: 6px;
        padding: 0 6px;
        font-size: 14px;
        line-height: 22px;
        gap: 2px;

        &.status-success {
          background: #cafce4;
          color: #21a665;
        }

        &.status-warning {
          background: #fae4dc;
          color: #ed744a;
        }
      }
    }
  }
}
</style>
