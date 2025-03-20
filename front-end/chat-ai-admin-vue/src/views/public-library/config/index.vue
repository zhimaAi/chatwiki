<style lang="less" scoped>
.add-library-page {
  .page-container {
    padding: 0 24px 24px;
  }
  .page-title {
    line-height: 24px;
    padding-bottom: 15px;
    font-size: 16px;
    font-weight: 600;
    color: #000000;
  }

  .form-box-wrapper {
    margin-bottom: 16px;
  }
}
</style>

<template>
  <div class="add-library-page">
    <ConfigPageMenu />
    <div class="page-container">
      <div class="form-box-wrapper">
        <LibraryForm
          v-model:value="state"
          :saveLoading="saveLoading1"
          @submit="handleSaveLibraryInfo"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { getLibraryInfo, editLibrary } from '@/api/library/index'
import { reactive, ref, toRaw, onMounted, computed } from 'vue'
import { message } from 'ant-design-vue'
import { useRoute } from 'vue-router'
import ConfigPageMenu from '../components/config-page-menu.vue'
import LibraryForm from './components/library-form.vue'

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

const saveLoading1 = ref(false)

const handleSaveLibraryInfo = (data) => {
  Object.assign(state, data)

  saveLoading1.value = true

  submit()
    .then(() => {
      saveLoading1.value = false
    })
    .catch(() => {
      saveLoading1.value = false
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
