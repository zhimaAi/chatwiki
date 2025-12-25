<template>
  <div class="_plugin-box">
    <LoadingBox v-if="loading"/>
    <div v-else-if="list.length" class="plugin-list">
      <div v-for="(item, idx) in list"
             :key="item.name"
             @click="linkDetail(item)"
             class="plugin-item">
        <div class="base-info">
          <img class="avatar" v-if="item.icon" :src="item.icon"/>
          <img class="avatar" v-else src="@/assets/img/default-mcp.png"/>
          <div class="info">
            <div class="head">
              <span class="name zm-line1">{{ item.title }}</span>
              <!-- 来源 -->
              <a class="source" :href="item.source_url" target="_blank">{{ item.source }}</a>
            </div>
            <div class="source">{{ item.author }}</div>
          </div>
        </div>
        <a-tooltip :title="getTooltipTitle(item.description, item)" placement="top">
          <div class="desc zm-line1" :ref="el => setDescRef(el, item)">{{ item.description }}</div>
        </a-tooltip>
        <div class="type-box" :ref="el => setTypeBoxRef(idx, el)">
          <div
            v-for="type in renderTypes(item, idx)"
            :key="type"
            class="type-tag"
          >{{type.type_title}}</div>
          <a-tooltip v-if="overflowMap[idx]" placement="top">
            <template #title>
              <div class="popover-type-box">
                <div
                  v-for="type in item.filter_type_list"
                  :key="type"
                  class="type-tag"
                >{{type.type_title}}</div>
              </div>
            </template>
            <div class="type-tag type-more" @click.stop>更多</div>
          </a-tooltip>
        </div>
      </div>
      <div ref="loadMoreRef" class="load-more-sentinel"></div>
    </div>
    <EmptyBox v-else title="暂无可用MCP"/>
  </div>
</template>

<script setup>
import {onMounted, onUnmounted, ref, watch, nextTick} from 'vue';
import {useRouter} from 'vue-router';
import EmptyBox from "@/components/common/empty-box.vue";
import {getMcpSquareList} from "@/api/mcp/index.js";
import LoadingBox from "@/components/common/loading-box.vue";

const props = defineProps({
  filterData: {
    type: Object,
    default: null
  },
  scrollRoot: {
    type: Object,
    default: null
  }
})
const router = useRouter()

const loading = ref(true)
const list = ref([])
const typeBoxRefs = ref({})
const overflowMap = ref({})
const visibleTagsMap = ref({})
const page = ref(1)
const size = ref(50)
const hasMore = ref(true)
const loadingMore = ref(false)
const loadMoreRef = ref(null)
let observer = null
let rafId = 0

onMounted(() => {
  loadData()
  window.addEventListener('resize', handleResize)
  initObserver()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (observer) {
    observer.disconnect()
    observer = null
  }
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
  page.value = 1
  hasMore.value = true
  loadData(true)
}

async function loadData(reset = false) {
  loading.value = true
  const params = { ...(props.filterData || {}), page: page.value, size: size.value }
  getMcpSquareList(params).then(res => {
    let _list = res?.data?.data || []
    if (reset || page.value === 1) {
      list.value = _list
    } else {
      list.value = [...list.value, ..._list]
    }
    hasMore.value = Array.isArray(_list) && _list.length >= size.value
    nextTick(() => {
      scheduleCalcOverflow()
    })
  }).finally(() => {
    loading.value = false
    loadingMore.value = false
  })
}

function linkDetail(item) {
  window.open(item.source_url, '_blank')
}

function renderTypes (item, idx) {
  const vis = visibleTagsMap.value[idx]
  return Array.isArray(vis) ? vis : item.filter_type_list
}

function setTypeBoxRef (index, el) {
  if (typeof index === 'undefined') return
  const prev = typeBoxRefs.value[index]
  if (el) {
    if (prev !== el) {
      typeBoxRefs.value[index] = el
      scheduleCalcOverflow()
    }
  } else {
    if (prev) {
      delete typeBoxRefs.value[index]
    }
  }
}

function calcOverflow() {
  const keys = Object.keys(typeBoxRefs.value)
  keys.forEach((idx) => {
    const el = typeBoxRefs.value[idx]
    const item = list.value[Number(idx)]
    if (!el || !item) return
    const containerWidth = el.clientWidth
    const style = getComputedStyle(el)
    const gapStr = style.gap || '0px'
    const gap = parseFloat(gapStr.split(' ')[0]) || 0

    const temp = document.createElement('div')
    temp.style.position = 'absolute'
    temp.style.visibility = 'hidden'
    temp.style.display = 'flex'
    temp.style.flexWrap = 'nowrap'
    temp.style.gap = gapStr
    el.appendChild(temp)

    const tagWidths = []
    item.filter_type_list.forEach((t) => {
      const d = document.createElement('div')
      d.className = 'type-tag'
      d.style.whiteSpace = 'nowrap'
      d.textContent = t.type_title
      temp.appendChild(d)
      tagWidths.push(d.getBoundingClientRect().width + gap)
    })

    const more = document.createElement('div')
    more.className = 'type-tag type-more'
    more.style.whiteSpace = 'nowrap'
    more.textContent = '更多'
    temp.appendChild(more)
    const moreWidth = more.getBoundingClientRect().width + gap

    el.removeChild(temp)

    const totalWidth = tagWidths.reduce((s, w, i) => s + w + (i > 0 ? gap : 0), 0)
    const prevVis = visibleTagsMap.value[idx]
    const prevOverflow = !!overflowMap.value[idx]
    if (totalWidth <= containerWidth) {
      const shouldUpdate = prevOverflow || !Array.isArray(prevVis) || prevVis.length !== item.filter_type_list.length
      if (shouldUpdate) {
        visibleTagsMap.value[idx] = item.filter_type_list
        overflowMap.value[idx] = false
      }
      return
    }

    let sum = 0
    let count = 0
    for (let i = 0; i < tagWidths.length; i++) {
      const require = (count > 0 ? gap : 0) + tagWidths[i]
      const reserve = (count > 0 ? gap : 0) + moreWidth
      if (sum + require + reserve <= containerWidth) {
        sum += require
        count++
      } else {
        break
      }
    }
    const prevLen = Array.isArray(prevVis) ? prevVis.length : 0
    const shouldUpdate = !prevOverflow || prevLen !== count
    if (shouldUpdate) {
      visibleTagsMap.value[idx] = item.filter_type_list.slice(0, count)
      overflowMap.value[idx] = true
    }
  })
}

function handleResize () {
  nextTick(() => scheduleCalcOverflow())
}

function initObserver() {
  observer = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        if (!loading.value && !loadingMore.value && hasMore.value && list.value.length > 0) {
          loadingMore.value = true
          page.value += 1
          loadData()
        }
      }
    })
  }, { root: props.scrollRoot || null, threshold: 0.1 })
  nextTick(() => {
    if (loadMoreRef.value) observer.observe(loadMoreRef.value)
  })
}

function scheduleCalcOverflow () {
  if (rafId) {
    cancelAnimationFrame(rafId)
    rafId = 0
  }
  rafId = requestAnimationFrame(() => {
    rafId = 0
    calcOverflow()
  })
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
._plugin-box {
  width: 100%;
  min-height: 100%;
}

.plugin-list {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
  position: relative;

  .plugin-item {
    flex: 0 0 calc((100% - 3 * 16px) / 4);
    min-width: 320px;
    padding: 24px;
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
    border-radius: 12px;
    border: 1px solid #E4E6EB;
    cursor: pointer;
    position: relative;

    &:hover {
      box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.12);
    }

    .type-tag {
      white-space: nowrap;
      display: flex;
      padding: 1px 8px;
      flex-direction: column;
      align-items: flex-start;
      gap: 10px;
      border-radius: 6px;
      background: #F2F4F7;
      color: #595959;
      font-size: 12px;
      font-style: normal;
      font-weight: 400;
      line-height: 16px;
    }

    .base-info {
      display: flex;
      align-items: center;
      gap: 12px;

      .avatar {
        width: 62px;
        height: 62px;
        border-radius: 16px;
        border: 1px solid #F0F0F0;
      }

      .head {
        display: flex;
        flex-direction: column;
        gap: 4px;

        .source {
          display: -webkit-box;
          -webkit-box-orient: vertical;
          -webkit-line-clamp: 1;
          overflow: hidden;
          color: #164799;
          text-overflow: ellipsis;
          font-size: 12px;
          font-style: normal;
          font-weight: 400;
          line-height: 20px;
        }
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
      }
    }

    .version {
      color: #8C8C8C;
      font-size: 12px;
      font-weight: 400;
      margin-top: -8px;
    }

    .desc {
      color: #595959;
      font-size: 14px;
      font-weight: 400;
      height: 22px;
      width: 100%;
    }

    .type-box {
      height: 18px;
      width: 100%;
      display: flex;
      gap: 4px;
      position: relative;
      max-height: 22px;
      overflow: hidden;
      flex-wrap: nowrap;
    }
  }
  
  .trigger-plugin-item {
    cursor: auto;
    
    &:hover {
      box-shadow: none;
    }
  }
  .type-more {}
  .popover-type-box {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    max-width: 480px;
  }
  .load-more-sentinel {
    width: 100%;
    height: 1px;
  }
} 


/* 大屏幕时：5 列 */
@media (min-width: 1900px) {
  .plugin-list .plugin-item {
    flex: 0 0 calc((100% - 4 * 16px) / 5);
  }
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
