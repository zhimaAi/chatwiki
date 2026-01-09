<template>
  <div class="_container">
    <div class="header">
      <div class="main-tab-box">
        <MainTab ref="tabRef"/>
      </div>
    </div>
    <div class="content">
      <div class="left-box hide-scrollbar">
        <div class="title">模板分类</div>
        <div class="cate-box">
          <div :class="['cate-item', {active: !filterData.category}]" @click="selectCate(null)">
            <div>全部</div>
            <span v-if="allCateCount > 0" class="count">{{ allCateCount }}</span>
          </div>
          <div
            v-for="cate in cates"
            :key="cate.id"
            :class="['cate-item', {active: cate.id == filterData.category?.id}]"
            @click="selectCate(cate)"
          >
            <div>{{ cate.name }}</div>
            <span class="count">{{ cate.template_count || 0}}</span>
          </div>
        </div>
      </div>
      <div class="right-box" @scroll="loadMore">
        <div class="filter-box">
          <a-input
            v-model:value.trim="filterData.keyword"
            @change="search"
            style="width: 360px;"
            allowClear
            placeholder="输入名称、描述搜索">
            <template #suffix>
              <SearchOutlined/>
            </template>
          </a-input>
          <div>
<!--            <a-button class="mr8">模板审核结果</a-button>-->
<!--            <a-button type="primary" :icon="h(PlusOutlined)">上架模板</a-button>-->
          </div>
        </div>
        <div class="_plugin-box">
          <div v-if="list.length" class="plugin-list">
            <router-link
              v-for="item in list"
              class="plugin-item"
              :key="item.id"
              :to="`/templates/detail?key=${item.robot_id}`" target="_blank"
            >
              <div class="base-info">
                <img class="avatar" :src="item.avatar || DEFAULT_TEMPLATE_AVATAR"/>
                <div class="info">
                  <div class="head">
                    <span class="name zm-line1">{{ item.name }}</span>
                  </div>
                  <div class="source">
                    <img class="author-avatar" :src="item.account_avatar"/>
                    <span class="zm-line1">{{ item.author }}</span>
                  </div>
                </div>
              </div>
              <a-tooltip :title="getTooltipTitle(item.description, item, 14, 2.5, 60)" placement="top">
                <div class="desc zm-line" :ref="el => setDescRef(el, item)">{{ item.description }}</div>
              </a-tooltip>
              <div class="action-box">
                <div class="left">
                  <TeamOutlined/>
                  {{ item.use_count || 0 }}
                </div>
                <div class="right">
                  <a @click.stop.prevent="useTemplate(item)"><svg-icon name="icon-rocket"/> 使用模板</a>
                </div>
              </div>
            </router-link>
          </div>
          <EmptyBox v-else title="暂无可用插件"/>
          <LoadingBox v-if="loading"/>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {onMounted, ref, reactive, h, computed} from 'vue';
import {useRoute, useRouter} from 'vue-router';
import {message, Modal} from 'ant-design-vue';
import {SearchOutlined, PlusOutlined, TeamOutlined, ExclamationCircleOutlined} from '@ant-design/icons-vue';
import MainTab from "@/views/explore/components/main-tab.vue";
import EmptyBox from "@/components/common/empty-box.vue";
import LoadingBox from "@/components/common/loading-box.vue";
import UpdateModal from "@/views/explore/plugins/components/update-modal.vue";
import {getTemplateCates, getTemplates, useRobotTemplate} from "@/api/explore/template.js";
import {DEFAULT_TEMPLATE_AVATAR} from "@/constants/index.js";
import {useCompanyStore} from "@/stores/modules/company.js";
import { setDescRef, getTooltipTitle, listScrollPullLoad } from '@/utils/index'

const route = useRoute()
const router = useRouter()
const tabRef = ref(null)
const active = ref(localStorage.getItem('zm:explore:active') || '4')
const loading = ref(true)
const finished = ref(false)
const list = ref([])
const cates = ref([])
const allCateCount = ref(0)
const filterData = reactive({
  keyword: '',
  category: null
})
const pagination = reactive({
  current: 1,
  pageSize: 50,
  total: 0
})
const companyStore = useCompanyStore()
const {companyInfo} = companyStore

const sysVersion = computed(() => {
  return companyInfo?.version
})

onMounted(() => {
  if (route.query.active > 1) {
    tabRef.value.change(route.query.active)
  }
  init()
})

function init() {
  loadCates()
  loadData()
}

function loadCates() {
  getTemplateCates().then(res => {
    let _list = res?.data || []
    allCateCount.value = 0
    _list.map(item => allCateCount.value += Number(item?.template_count || 0))
    cates.value = _list
  })
}

async function loadData() {
  loading.value = true
  let params = {
    page: pagination.current,
    size: pagination.pageSize,
  }
  if (filterData.keyword) params.keyword = filterData.keyword
  if (filterData.category) params.category_id = filterData.category?.id
  getTemplates(params).then(res => {
    let _list = res?.data?.data || []
    if (!_list.length || _list.length < pagination.pageSize) {
      finished.value = true
    }
    list.value.push(..._list)
    pagination.current += 1
    pagination.total = Number(res?.data?.total || 0)
  }).finally(() => {
    loading.value = false
  })
}

function loadMore(e) {
  if (loading.value || finished.value) return
  listScrollPullLoad(e, loadData)
}

function search() {
  pagination.current = 1
  pagination.total = 0
  list.value = []
  loadData()
}

function selectCate(cate) {
  filterData.category = cate
  search()
}

function checkVersion(sys_v, tpl_v) {
  sys_v = Number(sys_v.replace(/\D/g, ''))
  tpl_v = Number(tpl_v.replace(/\D/g, ''))
  return sys_v >= tpl_v
}

function useTemplate(item) {
  const run = () => {
    useRobotTemplate({template_id: item.id, csl_url: item.csl_url}).then(res => {
      message.success('使用成功')
      loadData()
      const {id, robot_key} = res.data
      const url = router.resolve({ path: '/robot/config/workflow', query: {id, robot_key} })
      window.open(url.href, '_blank')
    })
  }
  if (!checkVersion(sysVersion.value, item.version)) {
    Modal.confirm({
      title: '提示',
      content: '当前系统版本过低，可能无法使用此模板；请您升级到最新版本后使用！',
      icon: h(ExclamationCircleOutlined),
      okText: '继续使用',
      cancelText: '取 消',
      onOk: run
    })
  } else {
    Modal.confirm({
      title: '提示',
      content: `确定使用模板【${item.name}】创建应用吗?`,
      icon: h(ExclamationCircleOutlined),
      okText: '确 认',
      cancelText: '取 消',
      onOk: run
    })
  }
}

</script>

<style scoped lang="less">
._container {
  height: 100%;
  padding: 16px 24px 0 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.header {
  .main-tab-box {
    display: flex;
  }

  .tabs-box {
    display: flex;
    align-items: center;
    gap: 8px;

    .tab-item {
      padding: 5px 16px;
      border-radius: 6px;
      background: #EDEFF2;
      color: #595959;
      font-size: 14px;
      font-weight: 400;
      cursor: pointer;

      &.active {
        color: #2475fc;
        background: #D6E6FF;
      }
    }
  }
}

.content {
  flex: 1;
  display: flex;
  overflow: hidden;
  margin-top: 10px;

  .left-box {
    width: 208px;
    flex-shrink: 0;
    padding: 24px 24px 0 0;
    border-right: 1px solid #F0F0F0;
    overflow-y: auto;

    .title {
      color: #262626;
      font-size: 16px;
      font-style: normal;
      font-weight: 600;
      margin-bottom: 4px;
    }

    .cate-item {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 5px 8px;
      border-radius: 6px;
      color: #595959;
      cursor: pointer;

      &.active,
      &:hover {
        background: #E5EFFF;
        color: #2475FC;
      }

      .count {
        color: #8c8c8c;
        font-size: 12px;
      }
    }
  }

  .right-box {
    flex: 1;
    overflow-y: auto;
  }

  .filter-box {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    margin: 0 0 16px 24px;
  }
}

._plugin-box {
  padding-left: 24px;
}

.plugin-list {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;

  .plugin-item {
    flex: 0 0 calc((100% - 3 * 16px) / 4);
    //min-width: 200px;
    padding: 24px;
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    border-radius: 12px;
    border: 1px solid #E4E6EB;
    cursor: pointer;
    position: relative;

    &:hover {
      box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.12);
    }

    .base-info {
      display: flex;
      gap: 12px;

      .avatar {
        width: 62px;
        height: 62px;
        border-radius: 16px;
        border: 1px solid #F0F0F0;
      }

      .head {
        display: flex;
        align-items: center;
        gap: 4px;
      }

      .name {
        color: #262626;
        font-size: 16px;
        font-weight: 600;
      }

      .source {
        color: #8C8C8C;
        font-size: 12px;
        font-weight: 400;
        display: flex;
        align-items: center;
        gap: 4px;

        .author-avatar {
          width: 16px;
          height: 16px;
          border-radius: 16px;
        }
      }
    }

    .desc {
      color: #595959;
      font-size: 14px;
      font-weight: 400;
      height: 66px;
      width: 100%;
      -webkit-line-clamp: 3;
      line-clamp: 3;
      padding-left: 74px;
      margin-top: -16px;
    }

    .action-box {
      width: 100%;
      display: flex;
      align-items: center;
      justify-content: space-between;
      color: #595959;
      font-size: 14px;
      font-weight: 400;
      border-top: 1px solid #E4E6EB;
      padding-top: 20px;
      margin-top: 20px;
    }
  }

  .trigger-plugin-item {
    cursor: auto;

    &:hover {
      box-shadow: none;
    }
  }
}

.mr8 {
  margin-right: 8px;
}

/* 大屏幕时：5 列 */
@media (min-width: 1900px) {
  .plugin-list .plugin-item {
    flex: 0 0 calc((100% - 4 * 16px) / 5);
  }
}
</style>
