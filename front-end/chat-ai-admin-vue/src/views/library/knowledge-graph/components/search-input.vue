

<template>
  <div class="search-input-box" :class="{ 'focused': isFocused }">
    <input 
      ref="searchInput"
      class="search-input" 
      v-model="localValue" 
      placeholder="输入关键词搜索" 
      type="text"
      @keyup.enter="handleSearch"
      @focus="isFocused = true"
      @blur="isFocused = false"
      @input="onInput"
    />
    <div class="search-action-box">
      <span class="action-btn" v-if="value.length > 0" @click="clearValue">
        <CloseCircleFilled class="clear-icon" />
      </span>
      <span class="action-btn" @click="handleSearch">
        <SearchOutlined class="search-icon" />
      </span>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick, watch } from 'vue'
import { SearchOutlined, CloseCircleFilled } from '@ant-design/icons-vue'

const emit = defineEmits(['search', 'clear', 'input', 'update:value'])

const props = defineProps({
  value: {
    type: String,
    default: ''
  }
})

const localValue = ref(props.value)
const isFocused = ref(false)
const searchInput = ref(null)

const onInput = () => {
  emit('update:value', localValue.value)
}

// 清空输入框内容
const clearValue = () => {
  localValue.value = ''
  nextTick(() => {
    searchInput.value?.focus()
  })

  onInput()

  handleSearch()
}

const handleSearch = () => {
  emit('search', localValue.value)
}

watch(() => props.value, (val) => {
  if(localValue.value !== val){
    localValue.value = val
  }
})

</script>

<style lang="less" scoped>
.search-input-box {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 24px;
  border-radius: 12px;
  border: 1px solid #BFBFBF;
  background-color: #fff;
  transition: all 0.2s ease;
  filter: drop-shadow(0 4px 4px rgba(0, 89, 255, 0.16));

  &:hover,
  &.focused {
    border: 1px solid #2475FC;
    filter: drop-shadow(0 4px 8px rgba(0, 89, 255, 0.16));
  }


  .search-icon{
    color: #595959;
    font-size: 24px;
  }
  .clear-icon{
    color: #BFBFBF;
    font-size: 20px;
  }
  .search-input {
    flex: 1;
    line-height: 24px;
    font-size: 16px;
    font-weight: 400;
    border-radius: 8px;
    color: rgb(38, 38, 38);
    border: none;
    outline: none;
    background-color: #fff;
    &::placeholder {
      color: rgb(140, 140, 140);
    }
  }
  .search-action-box{
    display: flex;
    align-items: center;

    .action-btn{
      width: 32px;
      height: 32px;
      margin-left: 12px;
      border-radius: 6px;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      &:hover{
        background-color: rgba(228, 230, 235, 1);
      }
    }
  }
}


</style>