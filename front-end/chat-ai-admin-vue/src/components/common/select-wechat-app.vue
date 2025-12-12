<template>
  <a-modal
    v-model:open="visible"
    :title="title"
    :confirm-loading="saving"
    width="626px"
  >
    <slot name="alert">
      <a-alert type="info" class="zm-alert-info mt12">
        <template #message>
          公众号列表数据来源：系统管理>公众号管理。<a @click="linkOfficial">前往编辑公众号、添加公众号</a>
        </template>
      </a-alert>
    </slot>
    <LoadingBox v-if="loading" style="margin: 10vh 0;"/>
    <template v-else-if="list.length">
      <div style="margin-bottom: 12px;">
        <a-checkbox
          v-model:checked="allChcked"
          :disabled="!selectRows.length"
          @change="allChckedChange">全选
        </a-checkbox>
      </div>
      <a-checkbox-group v-model:value="selectedKeys">
        <div class="wechat-app-list">
          <div v-for="item in list"
               :key="item.id"
               :class="['wechat-app-item', {disabled: disabledAppIds.includes(item.app_id)}]"
               @click="select(item)"
          >
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
              <span @click.stop>
                <a-checkbox :value="item.app_id" :disabled="disabledAppIds.includes(item.app_id)" @change="checkedChange"/>
              </span>
            </div>
          </div>
        </div>
      </a-checkbox-group>
    </template>
    <EmptyBox v-else title="暂未绑定公众号">
      <template #desc>
        <a-button type="primary" @click="linkOfficial">绑定公众号</a-button>
      </template>
    </EmptyBox>

    <template #footer>
      <a-button type="primary" :disabled="!selectedKeys.length" @click="save">{{ okText }}</a-button>
    </template>
  </a-modal>
</template>
<script setup>
import {ref, computed, toRaw, nextTick} from 'vue';
import {useRouter} from 'vue-router';
import {ExclamationCircleFilled, CheckCircleFilled} from '@ant-design/icons-vue';
import LoadingBox from "@/components/common/loading-box.vue";
import EmptyBox from "@/components/common/empty-box.vue";
import {getWechatAppList} from "@/api/robot/index.js";
import {usePermissionStore} from "@/stores/modules/permission.js";

const emit = defineEmits(['ok'])
const props = defineProps({
  title: {
    type: String,
    default: '选择公众号'
  },
  okText: {
    type: String,
    default: '确 定'
  },
  disabledAppIds: {
    type: Array,
    default: () => ([])
  }
})
const router = useRouter()
const permissionStore = usePermissionStore()

const visible = ref(false)
const loading = ref(false)
const saving = ref(false)
const allChcked = ref(false)
const robotInfo = ref({})
const list = ref([])
const selectedKeys = ref([])

// const officialPerm = computed(() => {
//   let {role_permission, role_type} = permissionStore
//   return role_type == 1 || role_permission.includes('OfficialAccountMange')
// })

const selectRows = computed(() => {
  return list.value.filter(item => !props.disabledAppIds.includes(item.app_id))
})

function open(r, ids) {
  robotInfo.value = r
  selectedKeys.value = Array.isArray(ids) ? ids : []
  loadData()
  visible.value = true
}

function close() {
  visible.value = false
}

function loadData() {
  loading.value = true
  getWechatAppList({
    app_type: 'official_account',
    app_name: ''
  }).then((res) => {
    list.value = res.data || []
  }).finally(() => {
    loading.value = false
  })
}

function save() {
  let rows = list.value.filter(item => selectedKeys.value.includes(item.app_id))
  emit('ok', toRaw(selectedKeys.value), toRaw(rows))
}

function linkOfficial() {
  const url = router.resolve({path: '/user/official-account'})
  window.open(url.href, '_blank')
}

function select(item) {
  if (props.disabledAppIds.includes(item.app_id)) return
  let _idx = selectedKeys.value.indexOf(item.app_id)
  if (_idx > -1) {
    selectedKeys.value.splice(_idx, 1)
  } else {
    selectedKeys.value.push(item.app_id)
  }
  checkedChange()
}

function checkedChange() {
  nextTick(() => {
    allChcked.value = (selectedKeys.value.length === selectRows.value.length)
  })
}

function allChckedChange() {
  if (allChcked.value) {
    selectedKeys.value = []
    list.value.map(i => {
      if (!props.disabledAppIds.includes(i.app_id)) {
        selectedKeys.value.push(i.app_id)
      }
    })
  } else {
    selectedKeys.value = []
  }
}

defineExpose({
  open,
  close
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

    &.disabled {
      color: rgba(0, 0, 0, .25);
      background-color: #f5f5f5;
      border-color: #d9d9d9;
      text-shadow: none;
      box-shadow: none;
      cursor: not-allowed;
    }

    &:not(.disabled):hover {
      box-shadow: 0 3px 14px 2px #0000000d,
      0 8px 10px 1px #0000000f,
      0 5px 5px -3px #0000001a;
    }

    .app-info-block {
      display: flex;
      align-items: center;
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

.mt12 {
  margin-bottom: 12px;
}
</style>
