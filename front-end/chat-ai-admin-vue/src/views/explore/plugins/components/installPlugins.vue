<template>
  <div class="_plugin-box">
    <LoadingBox v-if="loading"/>
    <div v-else-if="allList.length" class="plugin-list">
      <div v-for="item in allList"
           @click="linkDetail(item)"
           :key="item.name"
           class="plugin-item"
           :class="{'trigger-plugin-item': item.local.type == 'trigger'}"
           >
        <div class="type-tag">{{item.filter_type_title}}</div>
        <div class="base-info">
          <img class="avatar" :src="item.icon"/>
          <div class="info">
            <div class="head">
              <span class="name zm-line1">{{ item.title }}</span>
            </div>
            <div class="source">{{ item.author }}</div>
          </div>
        </div>
        <a-tooltip :title="getTooltipTitle(item.description, item)" placement="top">
          <div class="desc zm-line1" :ref="el => setDescRef(el, item)">{{ item.description }}</div>
        </a-tooltip>
        <div class="version">版本：v{{ item.local.version }} <span v-if="item.has_update" class="tag">有更新</span></div>
        <div class="action-box">
          <div class="left" @click.stop="">
            <a-switch v-model:checked="item.local.has_loaded" @change="openChange(item)"/>
          </div>
          <div class="right">
            <span v-if="item.local.type != 'trigger'" class="zm-link-pointer" @click.stop="delPlugin(item)">删除</span>
            <template v-if="item.help_url">
              <a-divider type="vertical"/>
              <a @click.stop class="c595959" :href="item.help_url" target="_blank">使用说明</a>
            </template>
            <template v-if="showConifgPlugins.includes(item.name) || item.local.multiNode">
              <a-divider type="vertical"/>
              <a @click.stop="showConfigModal(item)">配置</a>
            </template>
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
        <div>更多功能可到 <a  @click="emit('tabChange', 3)">探索>插件广场</a> 中去添加</div>
        <div class="text-center mt24">
          <a-button class="btn" type="primary" @click="emit('tabChange', 3)">去添加</a-button>
        </div>
      </template>
    </EmptyBox>

    <UpdateModal ref="updateRef" @ok="loadData"/>
    <FeishuConfigBox ref="feishuRef"/>
    <ConfigBox ref="configRef"/>
  </div>
</template>

<script setup>
import {onMounted, ref, watch, h, computed} from 'vue';
import {useRouter} from 'vue-router';
import {message, Modal} from 'ant-design-vue';
import EmptyBox from "@/components/common/empty-box.vue";
import UpdateModal from "./update-modal.vue";
import LoadingBox from "@/components/common/loading-box.vue";
import {closePlugin, getInstallPlugins, getPluginConfig, openPlugin, uninstallPlugin, triggerSwitch, triggerConfigList} from "@/api/plugins/index.js";
import {jsonDecode} from "@/utils/index.js";
import FeishuConfigBox from "@/views/explore/plugins/components/feishu-config-box.vue";
import ConfigBox from "@/views/explore/plugins/components/config-box.vue";

const emit = defineEmits(['installReport', 'tabChange'])
const props = defineProps({
  filterData: {
    type: Object,
    default: null
  }
})
const router = useRouter()

const updateRef = ref(null)
const feishuRef = ref(null)
const configRef = ref(null)
const loading = ref(true)
const list = ref([])
const showConifgPlugins = [
  'feishu_bitable',
  'official_account_profile',
  'official_batch_tag',
  'official_send_template_message',
  'official_send_message',
  'official_intelligent_api',
  'official_account_batch_send',
  'official_account_comment'
]

onMounted(() => {
  loadData()
})

// 获取 tooltip 标题
function getTooltipTitle(text, record) {
  if (!text) return null
  const canvas = document.createElement('canvas')
  const context = canvas.getContext('2d')
  context.font = '14px -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif'
  const textWidth = context.measureText(text).width
  const maxWidth = record?.title_width || 120
  return textWidth > maxWidth ? text : null
}

// watch(() => props.filterData, () => {
//   loadData()
// }, {
//   immediate: true,
//   deep: true
// })

function search() {
  list.value = []
  loadData()
}

const triggerList = ref([])

const allList = computed(()=>{
  if(props.filterData.filter_type == 0 || props.filterData.filter_type == 5){
    return [...triggerList.value, ...list.value]
  }

  return list.value
})


function loadData() {
  triggerConfigList(props.filterData).then(res=>{
    let _list = res?.data || []
    _list = _list.map(item => {
      return {
        ...item.remote,
        local: item.local,
        has_update: item.remote?.latest_version != item.local?.version,
        installing: false
      }
    })
    triggerList.value = _list
  })
  loading.value = true
  getInstallPlugins(props.filterData).then(res => {
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

async function openChange(item) {
  const cancel = () => item.local.has_loaded = !item.local.has_loaded
  if (!item.local.has_loaded) {
    Modal.confirm({
      title: '确认关闭该插件？',
      content: '关闭后，其他应用的位置都不可使用！确认关闭？',
      okText: '确定',
      cancelText: '取消',
      onOk: () => {
        if(item.local.type == 'trigger'){
          triggerSwitch({
            id: item.trigger_config_id,
            switch_status: 0
          }).then(() => {
            message.success('已关闭')
          }).catch(cancel)
        }else{
          closePlugin({name: item.name}).then(() => {
            message.success('已关闭')
          }).catch(cancel)
        }
      },
      onCancel(){
        cancel()
      }
    })
  } else {
    if(item.local.type == 'trigger'){
      triggerSwitch({
        id: item.trigger_config_id,
        switch_status: 1
      }).then(()=>{
        message.success('已开启')
      })
      return
    }
    if (item.name == 'feishu_bitable') {
      let {data} = await getPluginConfig({name: item.name})
      data = jsonDecode(data, {})
      if (!Object.keys(data).length) {
        cancel()
        return Modal.confirm({
          title: '授权飞书应用接口权限',
          content: h('div', {style:{color:' red'}}, '请先完成信息配置，获取接口权限'),
          okText: '确定',
          cancelText: '取消',
          onOk: () => {
            showFeishuConfig(item)
          },
        })
      }
    }
    openPlugin({name: item.name}).then(() => {
      message.success('已开启')
    }).catch(cancel)
  }
}

function showFeishuConfig(item) {
  feishuRef.value.show(item)
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
  if(item.local.type == 'trigger'){
    return
  }
  router.push({
    path: '/plugins/detail',
    query: {
      name: item.name
    }
  })
}

function showConfigModal(item) {
  if (item.name === 'feishu_bitable') {
    showFeishuConfig(item)
  } else {
    configRef.value.show(item)
  }
}

function update(item) {
  updateRef.value.show(item, item.latest_version_detail, item.local)
}

function setDescRef(el, item) {
  if (el && item) {
    item.title_width = el.offsetWidth
  }
}

defineExpose({
  search,
})
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

.c595959 {
  color: #595959;
}
</style>
