<template>
  <div class="_plugin-box">
    <LoadingBox v-if="loading"/>
    <div v-else-if="list.length" class="plugin-list">
      <div v-for="item in list"
           @click="linkDetail(item)"
           :key="item.name"
           class="plugin-item">
        <div class="base-info">
          <img class="avatar" :src="item.icon"/>
          <div class="info">
            <div class="head">
              <span class="name zm-line1">{{ item.title }}</span>
            </div>
            <div class="source">{{ item.author }}</div>
          </div>
        </div>
        <div class="desc zm-line2">{{ item.description }}</div>
        <div class="version">版本：v{{ item.local.version }} <span v-if="item.has_update" class="tag">有更新</span></div>
        <div class="action-box">
          <div class="left" @click.stop="">
            <a-switch v-model:checked="item.local.has_loaded" @change="openChange(item)"/>
          </div>
          <div class="right">
            <span class="zm-link-pointer" @click.stop="delPlugin(item)">删除</span>
            <template v-if="item.has_update">
              <a-divider type="vertical"/>
              <a @click.stop="update(item)">更新</a>
            </template>
          </div>
        </div>
      </div>
    </div>
    <EmptyBox v-else title="暂未安装插件">
      <template #desc>
        <div>更多功能可到 <a>探索>插件广场</a> 中去添加</div>
        <div class="text-center mt24">
          <a-button class="btn" type="primary" @click="emit('tabChange', 3)">去添加</a-button>
        </div>
      </template>
    </EmptyBox>

    <UpdateModal ref="updateRef" @ok="loadData"/>
  </div>
</template>

<script setup>
import {onMounted, ref} from 'vue';
import {useRouter} from 'vue-router';
import {message, Modal} from 'ant-design-vue';
import EmptyBox from "@/components/common/empty-box.vue";
import UpdateModal from "./update-modal.vue";
import LoadingBox from "@/components/common/loading-box.vue";
import {closePlugin, getInstallPlugins, openPlugin, uninstallPlugin} from "@/api/plugins/index.js";

const emit = defineEmits(['installReport', 'tabChange'])
const router = useRouter()

const updateRef = ref(null)
const loading = ref(true)
const list = ref([])

onMounted(() => {
  loadData()
})

function loadData() {
  loading.value = true
  getInstallPlugins().then(res => {
    let _list = res?.data || []
    emit('installReport', _list.length)
    _list = _list.map(item => {
      return {
        ...item.remote,
        local: item.local,
        has_update: item.remote?.latest_version != item.local?.version,
        installing: false
      }
    })
    list.value = _list
  }).finally(() => {
    loading.value = false
  })
}

function openChange(item) {
  const cancel = () => item.local.has_loaded = !item.local.has_loaded
  if (!item.local.has_loaded) {
    Modal.confirm({
      title: '确认关闭该插件？',
      content: '关闭后，其他应用的位置都不可使用！确认关闭？',
      okText: '确定',
      cancelText: '取消',
      onOk: () => {
        closePlugin({name: item.name}).then(() => {
          message.success('已关闭')
        }).catch(() => cancel())
      },
      onCancel: () => {
        cancel()
      }
    })
  } else {
    openPlugin({name: item.name}).then(() => {
      message.success('已开启')
    }).catch(() => cancel())
  }
}


function delPlugin(item) {
  Modal.confirm({
    title: '确认删除该插件？',
    content: '删除后，其他应用的位置都不可使用！确认删除？',
    okText: '确定',
    cancelText: '取消',
    onOk: () => {
      uninstallPlugin({name: item.name}).then(() => {
        message.success('已删除')
        loadData()
      })
    }
  })
}

function linkDetail(item) {
  router.push({
    path: '/plugins/detail',
    query: {
      name: item.name
    }
  })
}

function update(item) {
  updateRef.value.show(item, item.latest_version_detail, item.local)
}
</script>

<style scoped lang="less">
@import "plugins";

.tag {
  padding: 1px 8px;
  border-radius: 6px;
  border: 1px solid #FB363F;
  color: #fb363f;
  font-size: 12px;
  font-weight: 400;
}

.text-center {
  text-align: center;
}

.mt24 {
  margin-top: 24px;
}
</style>
