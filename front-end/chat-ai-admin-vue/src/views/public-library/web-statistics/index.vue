<template>
  <div class="ai-config-page">
    <ConfigPageMenu />
    <div class="page-container">
      <a-alert
        class="statistics-alert"
        message="在下方输入百度统计等平台的统计代码，文档公开访问页面会加载这段代码，这样您就可以在百度统计等平台查看页面的统计数据"
        type="info"
        show-icon
      />
      <div class="code-form-box">
        <a-textarea v-model:value="formState.statistics_set" :rows="16" />
        <div class="btn-box">
          <a-button type="primary" @click="handleSubmit">保存</a-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { getLibraryInfo, editLibrary } from '@/api/library/index'
import { useRoute } from 'vue-router'
import { reactive, computed, onMounted, toRaw } from 'vue'
import { message } from 'ant-design-vue'
import ConfigPageMenu from '../components/config-page-menu.vue'
import { usePublicLibraryStore } from '@/stores/modules/public-library'

const libraryStore = usePublicLibraryStore()

const route = useRoute()
const library_id = computed(() => route.query.library_id)
const formState = reactive({
  statistics_set: ''
})

const handleSubmit = async () => {
  try {
    const res = await editLibrary({
      ...toRaw(formState)
    })

    if (res.code === 0) {
      libraryStore.getLibraryInfo()
      message.success('保存成功')
    }
  } catch (error) {
    console.error('保存失败:', error)
  }
}

const getData = () => {
  getLibraryInfo({ id: library_id.value }).then((res) => {
    let data = res.data || {}

    Object.assign(formState, data)
  })
}

onMounted(() => {
  getData()
})
</script>

<style lang="less" scoped>
.page-container {
  padding: 8px 24px 24px;
  .code-form-box {
    margin-top: 8px;
  }
  .btn-box {
    margin-top: 8px;
  }
}
</style>
