<template>
  <a-popover class="icon-selector" placement="bottomLeft" :arrow="false" @openChange="openChange">
    <template #content>
      <div class="icon-selector-container">
        <div class="search-box">
          <a-input 
            placeholder="请输入图标名称" 
            v-model:value="searchValue" 
            @pressEnter="search"
            allowClear
          >
            <template #suffix>
              <SearchOutlined style="color: rgba(0, 0, 0, 0.45)" />
            </template>
          </a-input>
        </div>
        <!-- 最近使用 -->
        <div v-if="recentIcons.length > 0 && !searchValue" class="recent-section">
          <div class="icon-list-label">最近使用</div>
          <div class="icon-list">
            <div 
              class="icon-item" 
              v-for="(item, index) in recentIcons" 
              :key="`recent-${index}`"
              :title="item.keywords.join(', ')"
              @click="selectIcon(item)"
            >
              {{ item.content }}
            </div>
          </div>
        </div>
        <!-- 列表 -->
        <div class="icon-list-label">{{ searchValue ? '搜索结果' : '表情与角色' }}</div>
        <div class="icon-list">
          <div 
            class="icon-item" 
            v-for="(item, index) in filteredIcons" 
            :key="index"
            :title="item.keywords.join(', ')"
            @click="selectIcon(item)"
          >
            {{ item.content }}
          </div>
          <div v-if="filteredIcons.length === 0" class="no-result">
            暂无匹配的图标
          </div>
        </div>
      </div>
    </template>

    <span class="icon-selector-handler" @click.prevent><slot></slot></span>
  </a-popover>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { SearchOutlined } from '@ant-design/icons-vue'
import docIcons from '@/config/open-doc/doc-icons'

const searchValue = ref('')
const recentIcons = ref([])

const openChange = (visible) => {
  if(visible){
    nextTick(() => {
      searchValue.value = '';
    })
  }
}

// 从本地存储加载最近使用的图标
const loadRecentIcons = () => {
  try {
    const stored = localStorage.getItem('recent-icons')
    if (stored) {
      recentIcons.value = JSON.parse(stored)
    }
  } catch (error) {
    console.warn('加载最近使用图标失败:', error)
    recentIcons.value = []
  }
}

// 保存最近使用的图标到本地存储
const saveRecentIcons = () => {
  try {
    localStorage.setItem('recent-icons', JSON.stringify(recentIcons.value))
  } catch (error) {
    console.warn('保存最近使用图标失败:', error)
  }
}

// 添加图标到最近使用列表
const addToRecentIcons = (icon) => {
  // 移除已存在的相同图标
  const existingIndex = recentIcons.value.findIndex(item => item.content === icon.content)
  if (existingIndex > -1) {
    recentIcons.value.splice(existingIndex, 1)
  }
  
  // 添加到列表开头
  recentIcons.value.unshift({
    ...icon,
    lastUsed: Date.now()
  })
  
  // 只保留最近10个
  if (recentIcons.value.length > 10) {
    recentIcons.value = recentIcons.value.slice(0, 10)
  }
  
  saveRecentIcons()
}

// 计算过滤后的图标列表
const filteredIcons = computed(() => {
  if (!searchValue.value.trim()) {
    return docIcons
  }
  
  const searchTerm = searchValue.value.toLowerCase().trim()
  return docIcons.filter(icon => {
    // 在关键词中进行半匹配搜索
    return icon.keywords.some(keyword => 
      keyword.toLowerCase().includes(searchTerm)
    )
  })
})

const search = () => {
  // 搜索功能通过计算属性实时响应，这里可以预留其他逻辑
}

// 定义事件
const emit = defineEmits(['select'])

// 选择图标
const selectIcon = (icon) => {
  addToRecentIcons(icon)
  emit('select', icon)

  // open.value = false;
}

// 组件挂载时加载最近使用的图标
onMounted(() => {
  loadRecentIcons()
})

</script>

<style lang="less" scoped>
.icon-selector-handler{
  display: inline-flex;
  cursor: pointer;
}

.icon-selector-container {
  width: 350px;

  .search-box {
    margin-bottom: 8px;
  }

  .recent-section {
    margin-bottom: 12px;
  }

  .icon-list-label {
    line-height: 20px;
    font-size: 12px;
    color: #8c8c8c;
  }

  .icon-list {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    margin-bottom: 8px;
    max-height: 200px;

    .icon-item {
      width: 28px;
      height: 28px;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      border-radius: 4px;
      font-size: 20px;
      transition: background-color 0.2s;
      
      &:hover {
        background-color: #f5f5f5;
      }
      
      &:active {
        background-color: #e6f7ff;
      }
    }
    
    .no-result {
      width: 100%;
      text-align: center;
      color: #999;
      font-size: 12px;
      padding: 20px 0;
    }
  }
}
</style>
