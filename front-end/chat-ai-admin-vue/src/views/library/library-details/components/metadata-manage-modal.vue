<template>
  <a-modal
    v-model:open="visible"
    width="746px"
    :confirm-loading="saving"
  >
    <template #title>
      <div class="title">{{ isEdit ? t('title_edit_metadata') : t('title_metadata_manage') }} <span
        class="desc">{{ t('desc_metadata') }}</span></div>
    </template>
    <template #footer>
      <template v-if="isEdit">
        <a-button @click="visible = false">{{ t('btn_cancel') }}</a-button>
        <a-button type="primary" @click="saveValues">{{ t('btn_save') }}</a-button>
      </template>
    </template>
    <div class="_main">
      <template v-if="isEdit">
        <div class="meta-card edit-box">
          <div v-for="(item, index) in customList" :key="index" class="meta-item">
            <div class="left">
              <a-checkbox v-model:checked="item.checked"></a-checkbox>
              <a-tooltip :title="item.is_show == 1 ? t('label_list_show') : t('label_list_hide')">
                <span @click="showChange(item)">
                  <EyeOutlined v-if="item.is_show == 1" class="icon"/>
                  <EyeInvisibleOutlined v-else class="icon"/>
                </span>
              </a-tooltip>
              <span class="meta-name">{{ item.name }}</span>
              <span class="meta-type">{{ typeMap[item.type] }}</span>
            </div>
            <a-date-picker v-if="item.type == 1" v-model:value="item.value" format="YYYY-MM-DD HH:mm" class="input-box"
                           @change="item.checked = true"/>
            <a-input-number v-else-if="item.type === 2" v-model:value="item.value" :placeholder="t('ph_please_input')" :maxlength="20"
                            class="input-box" @change="item.checked = true"/>
            <a-input v-else v-model:value.trim="item.value" :placeholder="t('ph_please_input')" :maxlength="20" class="input-box"
                     @change="item.checked = true"/>
          </div>
        </div>
      </template>
      <template v-else>
        <div class="mt16">
          <a-button type="primary" ghost @click="addData">
            <PlusOutlined/>
            {{ t('btn_add_metadata') }}
          </a-button>
        </div>
        <div v-if="customList.length" class="meta-card">
          <template v-for="(item, index) in customList" :key="index">
            <div v-if="item.is_edit" class="meta-edit-item">
              <a-input v-model:value.trim="item.name" :placeholder="t('ph_please_input')" :maxlength="20"/>
              <a-select v-model:value="item.type" :options="typeData" class="min" :disabled="item.id > 0"/>
              <a-select v-model:value="item.is_show" :options="showData" class="min"/>
              <div class="action-box">
                <a-button @click="cancelEdit(item, index)">{{ t('btn_cancel') }}</a-button>
                <a-button type="primary" @click="dataChange(item)">{{ t('btn_save') }}</a-button>
              </div>
            </div>
            <div v-else class="meta-item">
              <div class="left">
                <a-tooltip :title="item.is_show == 1 ? t('label_list_show') : t('label_list_hide')">
                  <span @click="showChange(item)">
                    <EyeOutlined v-if="item.is_show == 1" class="icon"/>
                    <EyeInvisibleOutlined v-else class="icon"/>
                  </span>
                </a-tooltip>
                <span class="meta-name">{{ item.name }}</span>
                <span class="meta-type">{{ typeMap[item.type] }}</span>
              </div>
              <div class="right">
                <EditOutlined class="icon" @click="item.is_edit = true"/>
                <CloseCircleOutlined class="icon" @click="delData(item, index)"/>
              </div>
            </div>
          </template>
        </div>
        <div class="meta-card">
          <div class="meta-tit">{{ t('label_builtin') }}</div>
          <div v-for="(item, index) in builtList" :key="index" class="meta-item">
            <div class="left">
              <a-tooltip :title="item.is_show == 1 ? t('label_list_show') : t('label_list_hide')">
                <span @click="showChange(item)">
                  <EyeOutlined v-if="item.is_show == 1" class="icon"/>
                  <EyeInvisibleOutlined v-else class="icon"/>
                </span>
              </a-tooltip>
              <span class="meta-name">{{ item.name }}</span>
              <span class="meta-type">{{ typeMap[item.type] }}</span>
            </div>
          </div>
        </div>
      </template>
    </div>
  </a-modal>
</template>
<script setup>
import {ref} from 'vue'
import dayjs from 'dayjs'
import {PlusOutlined, EyeOutlined, EyeInvisibleOutlined, EditOutlined, CloseCircleOutlined} from '@ant-design/icons-vue'
import {Modal, message} from 'ant-design-vue'
import {useI18n} from '@/hooks/web/useI18n'
import {
  delLibraryMetaSchema,
  getLibraryMetaSchemaList,
  saveLibraryMetaSchema,
  saveMetadata,
  saveQaMetadata
} from "@/api/library/index.js";

const {t} = useI18n('views.library.library-details.components.metadata-manage-modal')

const emit = defineEmits(['change'])
const props = defineProps({
  libraryId: {
    type: [String, Number]
  },
  fileId: { // 普通知识库修改元数据
    type: [String, Number]
  },
  qaIds: { // 问答知识库修改问答元数据
    type: Array,
    default: () => ([])
  }
})
const showData = [
  {value: 1, label: t('label_list_show')},
  {value: 0, label: t('label_list_hide')},
]
const typeData = [
  {value: 0, label: t('label_type_string')},
  {value: 1, label: t('label_type_time')},
  {value: 2, label: t('label_type_number')}
]
const typeMap = {
  0: t('label_type_string'),
  1: t('label_type_time'),
  2: t('label_type_number')
}
const visible = ref(false)
const saving = ref(false)
const isEdit = ref(false)
const customList = ref([])
const builtList = ref([])

function show(edit = false, metas = null) {
  if (Array.isArray(metas)) {
    metaListFormat(metas)
  } else {
    loadBaseData()
  }
  isEdit.value = edit
  visible.value = true
}

function loadBaseData() {
  getLibraryMetaSchemaList({library_id: props.libraryId}).then(res => {
    metaListFormat(res?.data || [])
  })
}

function metaListFormat(metas) {
  customList.value = []
  builtList.value = []
  for (let meta of metas) {
    if (meta.is_builtin == 1) {
      builtList.value.push(meta)
    } else {
      meta.backup = JSON.parse(JSON.stringify(meta))
      if (meta.type == 1 && meta.value > 0) {
        meta.value = dayjs(meta.value * 1000)
      }
      meta.is_edit = false
      meta.checked = false
      customList.value.push(meta)
    }
  }
}

function addData() {
  customList.value.unshift({
    id: '',
    is_edit: true,
    is_builtin: 0,
    is_show: 0,
    name: "",
    type: 0,
    library_id: props.libraryId,
    backup: {},
  })
}

function delData(record, index) {
  Modal.confirm({
    title: t('title_tip'),
    content: t('msg_confirm_delete'),
    okText: t('btn_delete'),
    cancelText: t('btn_cancel'),
    onOk: () => {
      delLibraryMetaSchema({id: record.id}).then(() => {
        customList.value.splice(index, 1)
        message.success(t('msg_deleted'))
      })
    }
  })
}

function cancelEdit(record, index) {
  if (!record.id) {
    customList.value.splice(index, 1)
  } else {
    const {name, type, is_show} = record.backup
    record.name = name
    record.type = type
    record.is_show = is_show
    record.is_edit = false
  }
}

function showChange(record) {
  record.is_show = record.is_show == 1 ? 0 : 1
  dataChange(record)
}

function dataChange(record) {
  // 修改元数据字段
  if (!record.name) return message.warning(t('msg_please_enter_metadata_name'))
  saveLibraryMetaSchema({
    ...record,
    is_show: Number(record.is_show),
    library_id: props.libraryId
  }).then(res => {
    const {name, type, is_show} = record
    record.id = res?.data
    record.is_edit = false
    record.backup = {name, type, is_show}
    message.success(t('msg_saved'))
    emit('change')
  })
}

function saveValues() {
  // 修改元数据
  if (!customList.value.length) return message.error(t('msg_no_available_metadata'))
  let _data = []
  for (let record of customList.value) {
    if (!record.checked) continue
    let val = record.value
    if (record.type == 1) {
      if (record.value) {
        val = record.value.startOf('minute').unix()
      } else {
        val = ""
      }
    }
    _data.push({
      key: record.key,
      value: val
    })
  }
  if (!_data.length) return message.error(t('msg_please_select_data'))
  let req
  if (props.qaIds.length) {
    req = saveQaMetadata({
      ids: props.qaIds.toString(),
      list: JSON.stringify(_data)
    })
  } else if (props.fileId) {
    req = saveMetadata({
      file_id: props.fileId,
      list: JSON.stringify(_data)
    })
  } else {
    return message.error(t('msg_missing_library_info'))
  }
  req.then(() => {
    message.success(t('msg_saved'))
    emit('change')
    visible.value = false
  })
}

defineExpose({
  show,
})
</script>
<style lang="less" scoped>
.title{
  padding-right: 30px;
  line-height: 18px;
}
.desc {
  color: #8c8c8c;
  line-height: 18px;
  font-size: 14px;
  font-weight: 400;
  margin-left: 16px;
}

._main {
  min-height: 200px;
}

.meta-card {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 16px;

  &.edit-box {

    .meta-item {
      border-radius: 6px;
      border: 1px solid #D9D9D9;
      padding: 4px 12px;

      .left {
        border: none;
        padding: 0;
      }
    }

    .input-box {
      width: 260px;
    }
  }
}

.meta-tit {
  color: #262626;
}

.meta-edit-item {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #595959;

  .min {
    flex-shrink: 0;
    width: 140px;
  }

  .action-box {
    display: flex;
    align-items: center;
    gap: 8px;
  }
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
  justify-content: space-between;

  &:hover {
    .right {
      display: flex;
    }
  }

  .left {
    flex: 1;
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 4px 12px;
    border-radius: 6px;
    border: 1px solid #D9D9D9;
    color: #595959;

    .icon {
      color: #8c8c8c !important;
      cursor: pointer;
    }

    .meta-type {
      color: #8c8c8c;
    }
  }

  .right {
    display: none;
    align-items: center;
    gap: 8px;

    .icon {
      cursor: pointer;
      padding: 4px;
    }
  }
}

.mt16 {
  margin-top: 16px;
}

.link-btn {
  padding: 0 8px 0 !important;
}
</style>
