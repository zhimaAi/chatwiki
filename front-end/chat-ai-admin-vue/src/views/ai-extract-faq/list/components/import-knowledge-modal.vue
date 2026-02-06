<template>
  <div>
    <a-modal
      v-model:open="open"
      :confirm-loading="confirmLoading"
      :title="t('modal_title')"
      width="746px"
      wrapClassName="no-padding-modal"
      :bodyStyle="{ 'padding-right': '24px', 'max-height': '540px', 'overflow-y': 'auto' }"
      @ok="handleSave"
    >
      <div class="modal-box">
        <a-radio-group v-model:value="tabs" @change="handleChangeTabs">
          <a-radio-button :value="1">{{ t('import_to_existing_library') }}</a-radio-button>
          <a-radio-button :value="2">{{ t('create_and_import_library') }}</a-radio-button>
        </a-radio-group>
        <template v-if="tabs == 1">
          <a-divider
            style="font-size: 14px; color: #8c8c8c; font-weight: 400"
            orientation="left"
            orientation-margin="0px"
          >
            {{ t('select_library_title') }}
          </a-divider>
          <div class="filter-block">
            <a-input-search
              v-model:value="library_name"
              :placeholder="t('search_library_placeholder')"
              style="width: 240px"
              @search="getLibraryData()"
            />
            <a-button @click="getLibraryData(true)" :icon="createVNode(SyncOutlined)">{{ t('refresh_button') }}</a-button>
          </div>
          <div class="empty-block" v-if="libraryLists.length == 0">
            <a-empty>{{ t('no_data') }}</a-empty>
          </div>
          <div class="list-box">
            <div
              class="list-item"
              @click="handleClickItem(item)"
              v-for="item in libraryLists"
              :key="item.id"
            >
              <a-radio :checked="selectIds == item.id"></a-radio>
              <div class="info-block">
                <div class="title">{{ item.library_name }}</div>
                <div class="desc">{{ item.library_intro }}</div>
              </div>
            </div>
          </div>
        </template>
        <AddQaLibrary ref="addQaLibraryRef" @ok="createLibrary" v-else />
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, createVNode, nextTick } from 'vue'
import { getLibraryList, importParagraph } from '@/api/library/index'
import { SyncOutlined } from '@ant-design/icons-vue'
import AddQaLibrary from './add-qa-library.vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'


const { t } = useI18n('views.ai-extract-faq.list.components.import-knowledge-modal')

const emit = defineEmits(['ok'])
const formState = reactive({})

const open = ref(false)
const tabs = ref(1)
const library_name = ref('')
const selectIds = ref('')

const file_id = ref('')
const ids = ref('')
const show = (data) => {
  file_id.value = data.id
  ids.value = data.ids || ''
  tabs.value = 1
  library_name.value = ''
  selectIds.value = ''
  open.value = true
  getLibraryData()
}
const confirmLoading = ref(false)

const handleSave = () => {
  if (tabs.value == 1) {
    importLibrary()
  } else {
    addQaLibraryRef.value.handleOk()
  }
}

const handleSaveForm = (library_id) => {
  confirmLoading.value = true
  importParagraph({
    library_id,
    file_id: file_id.value,
    ids: ids.value,
  })
    .then((res) => {
      message.success(t('import_success'))
      open.value = false
      emit('ok')
    })
    .finally(() => {
      confirmLoading.value = false
    })
}
const importLibrary = () => {
  if (!selectIds.value) {
    return message.error(t('select_library_error'))
  }
  handleSaveForm(selectIds.value)
}
const createLibrary = (id) => {
  handleSaveForm(id)
}

const libraryLists = ref([])
const getLibraryData = (isResh) => {
  getLibraryList({
    library_name: library_name.value,
    show_open_docs: 1
  }).then((res) => {
    libraryLists.value = res.data.filter(item => item.type == 2)
    isResh && message.success(t('refresh_success'))
  })
}

const handleClickItem = (item) => {
  selectIds.value = item.id
}

const addQaLibraryRef = ref(null)
const handleChangeTabs = () => {
  if (tabs.value == 2) {
    nextTick(() => {
      addQaLibraryRef.value && addQaLibraryRef.value.show()
    })
  }
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.modal-box {
  margin-top: 40px;
}
.filter-block {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.empty-block {
  margin: 32px 0;
}
.list-box {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-top: 8px;
  .list-item {
    width: calc(50% - 8px);
    border: 1px solid var(--07, #e4e6eb);
    padding: 14px 16px;
    display: flex;
    align-items: center;
    font-size: 14px;
    line-height: 22px;
    cursor: pointer;
    border-radius: 6px;
    transition: all 0.3s ease;
    &:hover {
      border: 1px solid var(---, #2475fc);
      box-shadow: 0 4px 16px 0 #1b3a6929;
    }
    .info-block {
      flex: 1;
      .title {
        color: #262626;
        font-weight: 600;
      }
      .desc {
        color: #8c8c8c;
        font-size: 12px;
        font-weight: 400;
        line-height: 20px;
        margin-top: 2px;
      }
    }
  }
}
.mt8 {
  margin-top: 8px;
}
</style>
