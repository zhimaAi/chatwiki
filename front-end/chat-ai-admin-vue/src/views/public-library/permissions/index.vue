<template>
  <div class="permissions-page">
    <ConfigPageMenu />
    <div class="page-container">
      <div class="form-box-wrapper">
        <PermissionsForm
          v-model:value="state"
          :saveLoading="saveLoading2"
          @submit="handleSavePermissions"
        />
      </div>
      <div class="collaborator-list-wrapper">
        <CollaboratorList v-if="state.library_key" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { getLibraryInfo, editLibrary } from '@/api/library/index'
import { useRoute } from 'vue-router'
import { reactive, ref, toRaw, onMounted, computed, provide } from 'vue'
import { message } from 'ant-design-vue'
import ConfigPageMenu from '../components/config-page-menu.vue'
import PermissionsForm from './components/permissions-form.vue'
import CollaboratorList from './components/collaborator-list.vue'
import { usePublicLibraryStore } from '@/stores/modules/public-library'

const libraryStore = usePublicLibraryStore()

const route = useRoute()

const id = computed(() => route.query.library_id)

const state = reactive({
  type: 1,
  library_intro: '',
  library_name: '',
  model_config_id: '',
  model_define: '',
  use_model: '',
  is_offline: false,
  access_rights: '0',
  share_url: '',
  library_key: '',
  avatar: '',
  avatar_file: ''
})

// 使用provide向子组件提供响应式state
provide('libraryState', state)

// 提供更新state的方法
provide('updateState', (newState) => {
  Object.assign(state, newState)
})

const saveLoading2 = ref(false)

const handleSavePermissions = (data) => {
  Object.assign(state, data)

  saveLoading2.value = true

  submit()
    .then(() => {
      saveLoading2.value = false
    })
    .catch(() => {
      saveLoading2.value = false
    })
}

const submit = async () => {
  let data = { ...toRaw(state) }

  delete data.library_key

  if (data.avatar_file) {
    data.avatar = data.avatar_file
    delete data.avatar_file
  } else {
    delete data.avatar
    delete data.avatar_file
  }

  await editLibrary(data)

  libraryStore.getLibraryInfo()
  
  message.success('修改成功')
}

const getData = () => {
  getLibraryInfo({ id: id.value }).then((res) => {
    Object.assign(state, res.data)
  })
}

onMounted(() => {
  getData()
})
</script>

<style lang="less" scoped>
.page-container {
  padding: 0 24px 24px;

  .form-box-wrapper {
    margin-bottom: 16px;
  }
}
</style>
