<style lang="less" scoped>
.chat-monitor-page {
  display: flex;
  overflow: hidden;
  width: 100%;
  height: 100%;
  border-radius: 6px;
}

.page-left {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  overflow-y: auto;
  background-color: #fff;
  border-right: 1px solid var(--07, #E4E6EB);

  .toolbar-box {
    margin: 16px 16px 0;
  }

  .app-list-box {
    padding: 16px 16px 16px 24px;
  }

  .chat-list-wrapper {
    flex: 1;
    overflow: hidden;
  }
}

.page-body {
  display: flex;
  flex: 1;
  flex-direction: column;
  overflow: hidden;
  width: 100%;
  height: 100%;
}

.library-checkbox-box {

  .list-tools {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 16px;
    margin-bottom: 8px;
  }

  .list-box {
    display: flex;
    flex-direction: column;
    height: auto;
    margin: 0 16px;

    .list-item-wraapper {
      display: flex;
      width: 256px;
      height: 40px;
      padding: 9px 8px;
      align-items: center;
      gap: 8px;
      border-radius: 6px;

      &:hover {
        background: #F2F4F7;
      }
    }

    .list-item {
      display: flex;
      align-items: center;
      gap: 8px;
      width: 100%;

      &:hover {
        cursor: pointer;
      }

      .list-item-box {
        display: flex;
        align-items: center;
        gap: 4px;

        .library-icon {
          display: flex;
          width: 16px;
          height: 16px;
          line-height: 22px;
          justify-content: center;
          align-items: center;
          border-radius: 3px;
        }

        .library-name {
          height: 22px;
          flex: 1 0 0;
          overflow: hidden;
          color: #262626;
          text-overflow: ellipsis;
          white-space: nowrap;
          font-size: 14px;
          font-style: normal;
          font-weight: 400;
          line-height: 22px;
        }
      }
    }

    .list-item :deep(span:last-child) {
      flex: 1;
      overflow: hidden;
    }
  }
}
</style>

<template>
  <div class="chat-monitor-page">
    <div class="page-left">
      <div class="toolbar-box">
        <ContainTabs :tabs="tabs" v-model:value="activeKey" @change="onChangeTab" />
      </div>
      <div class="library-checkbox-box">
        <div class="app-list-box">
        <a-checkbox v-model:checked="state.checkAll" :indeterminate="state.indeterminate" @change="onCheckAllChange">
            {{ t('select_all') }}
          </a-checkbox>
        </div>
        <a-checkbox-group v-model:value="state.checkedList" @change="onCheckListChange">
          <div class="list-box" ref="scrollContainer">
            <div class="list-item-wraapper" v-for="item in plainOptions" :key="item.id">
              <div class="list-item">
                <a-checkbox :value="item.id"></a-checkbox>
                <div class="list-item-box">
                  <div class="library-icon">
                    <a-image :preview="false" :width="32" :src="item.avatar" />
                  </div>
                  <div class="library-name" @click="onGoLibraryInfo(item.id)">{{ item.library_name }}</div>
                </div>
              </div>
            </div>
          </div>
        </a-checkbox-group>
      </div>
    </div>
    <div class="page-body">
      <SearchBox ref="searchBoxRef" @defaultParmas="onDefaultParmas" @search="onSearch" :isError="isError" :messageObj="messageObj"
        :librarySearchData="librarySearchData" :library_ids="state.checkedList" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, reactive, watch } from 'vue'
import { storeToRefs } from 'pinia'
import SearchBox from './components/search-box.vue'
import ContainTabs from '@/components/cu-tabs/contain-tabs.vue'
import { useSearchLiraryStore } from '@/stores/modules/search-lirary'
import { getLibraryList } from '@/api/library'
import { formatFileSize } from '@/utils/index'
import { LIBRARY_NORMAL_AVATAR, LIBRARY_QA_AVATAR } from '@/constants/index'
import { useSearchStore } from '@/stores/modules/search-lirary'
import { showErrorMsg } from '@/utils/index'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library-search.index')

const searchStore = useSearchStore()

let {
  getIndeterminate,
  getCheckAll,
  getCheckedList,
  getActiveKey
} = storeToRefs(searchStore)

const searchLiraryStore = useSearchLiraryStore()
const {
  getLibrarySearchFn,
  searchMessage
} = searchLiraryStore

const { messageObj } = storeToRefs(searchLiraryStore)

const searchBoxRef = ref(null)
const librarySearchData = ref(null)
const activeKey = ref('all')
const tabs = ref([
  {
    title: t('all'),
    value: 'all'
  }, {
    title: t('normal'),
    value: '0'
  },
  {
    title: t('qa'),
    value: '2'
  }
])

const defaultParmas = ref({})
const plainOptions = ref([]);

const state = reactive({
  indeterminate: false,
  checkAll: false,
  checkedList: [],
});

const onDefaultParmas = (obj) => {
  defaultParmas.value = obj
}

const isError = ref(false)
const onSearch = async (keyword) => {
  // 搜索之前获取最新的配置数据
  isError.value = false
  let res = await getLibrarySearchFn()
  librarySearchData.value = res.data

  if (!librarySearchData.value.id) {
    librarySearchData.value = Object.assign(defaultParmas.value)
  }
  console.log(librarySearchData.value,'---')

  if (!librarySearchData.value.model_config_id || !librarySearchData.value.use_model) {
    isError.value = true
    return showErrorMsg(t('configure_before_use'))
  }

  searchBoxRef.value?.handleRecallTest(librarySearchData.value)
  searchMessage(keyword, state.checkedList, librarySearchData.value)
}

const onCheckAllChange = (e) => {
  Object.assign(state, {
    checkedList: e.target.checked ? plainOptions.value.map(item => item.id) : [],
    indeterminate: !!state.checkedList.length && state.checkedList.length < plainOptions.value.length,
  });

  searchStore.setCheckedList(state.checkedList)
  searchStore.setIndeterminate(state.indeterminate)
};

const onCheckListChange = (checkeds) => {
  searchStore.setCheckedList(checkeds)
}

const onGoLibraryInfo = (id) => {
  window.open(`#/library/details/knowledge-document?id=${id}`)
}

watch(
  () => state.checkedList,
  val => {
    state.indeterminate = !!val.length && val.length < plainOptions.value.length;
    state.checkAll = val.length >= plainOptions.value.length;

    searchStore.setIndeterminate(state.indeterminate)
    searchStore.setCheckAll(state.checkAll)
  },
);


const updateTabNumber = () => {
  // let all = 0
  // let normal = 0
  // let qa = 0

  // plainOptions.value.forEach(item => {
  //   if (item.type == 0) {
  //     normal += 1
  //   }
  //   if (item.type == 2) {
  //     qa += 1
  //   }
  //   all += 1
  // })

  tabs.value = [
    {
      title: t('all'),
      value: 'all'
    },
    {
      title: t('normal'),
      value: '0'
    },
    {
      title: t('qa'),
      value: '2'
    }
  ]
  // tabs.value = [
  //   {
  //     title: '全部 (' + all + ')',
  //     value: 'all' 
  //   },
  //   {
  //     title: '普通知识库 (' + normal + ')',
  //     value: '0' 
  //   },
  //   {
  //     title: '问答知识库 (' + qa + ')',
  //     value: '2'
  //   }
  // ]
}

// 判断一个数组里面是否包含另外一个数组的所有数据
const containsAll = (main, sub) => {
  if (sub.length === 0) return false; // 空子集不认为包含所有
  const mainSet = new Set(main);
  return sub.every(item => mainSet.has(item));
};

// 判断一个数组里面是否不包含另一个数组的任何值
const hasIntersection = (mainArray, subArray) => {
  return !subArray.some(item => mainArray.includes(item));
}

const getList = async () => {
  let type = activeKey.value === 'all' ? '' : activeKey.value

  await getLibraryList({ type }).then((res) => {
    let data = res.data || []

    data.forEach((item) => {
      item.file_size_str = formatFileSize(item.file_size)

      if (!item.avatar) {
        item.avatar = item.type == 0 ? LIBRARY_NORMAL_AVATAR : LIBRARY_QA_AVATAR
      }
    })

    plainOptions.value = data;

    const plainOptionIds = plainOptions.value.map(item => item.id)

    // 主逻辑优化
    if (getCheckedList.value?.length && plainOptionIds?.length) {
      const checked = getCheckedList.value;
      const options = plainOptionIds;
      
      // 核心状态计算
      const isAllIncluded = containsAll(checked, options);
      const hasCommonItems = hasIntersection(checked, options);
      const isLonger = checked.length >= options.length;

      // 修正后的状态判断
      const checkAll = isLonger && isAllIncluded;
      const indeterminate = !checkAll && !hasCommonItems; // 关键修正点

      // 统一状态更新
      const stateUpdate = {
        checkAll,
        indeterminate
      };
      
      searchStore.$patch(stateUpdate);
      Object.assign(state, stateUpdate);
    } else {
      // 处理边界情况：无选中项或选项为空
      const stateUpdate = {
        checkAll: false,
        indeterminate: false
      };
      searchStore.$patch(stateUpdate);
      Object.assign(state, stateUpdate);
    }

    state.checkedList = getCheckedList.value
    state.indeterminate = getIndeterminate.value

    // 如果用户没有操作，第一次使用，默认全选
    if (!getCheckAll.value && getCheckedList.value.length == 0 && !getIndeterminate.value && getActiveKey.value == 'all') {
      searchStore.setCheckAll(true)
      state.checkedList = plainOptionIds
      searchStore.setCheckedList(plainOptionIds)
    }

    if (activeKey.value === 'all') {
      updateTabNumber()
    }
  })
}

const onChangeTab = (val) => {
  searchStore.setActiveKey(val)
  getList()
}

onMounted(async () => {
  if (getActiveKey.value) {
    activeKey.value = getActiveKey.value
  }
  await getList()
  let res = await getLibrarySearchFn()
  librarySearchData.value = res.data
})

onUnmounted(() => {

})
</script>
