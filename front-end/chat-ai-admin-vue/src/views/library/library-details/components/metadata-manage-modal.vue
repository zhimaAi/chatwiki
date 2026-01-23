<template>
  <a-modal
    v-model:open="visible"
    width="746px"
    :confirm-loading="saving"
  >
    <template #title>
      <div>{{ isEdit ? '修改元数据' : '元数据管理' }} <span
        class="desc">元数据用于描述文档的属性。应用召回知识时可根据元数据进行筛选</span></div>
    </template>
    <template #footer>
      <template v-if="isEdit">
        <a-button @click="visible = false">取 消</a-button>
        <a-button type="primary" @click="saveValues">保 存</a-button>
      </template>
    </template>
    <div class="_main">
      <template v-if="isEdit">
        <div class="meta-card edit-box">
          <div v-for="(item, index) in customList" :key="index" class="meta-item">
            <div class="left">
              <a-checkbox v-model:checked="item.checked"></a-checkbox>
              <a-tooltip :title="item.is_show == 1 ? '列表显示' : '列表隐藏'">
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
            <a-input-number v-else-if="item.type === 2" v-model:value="item.value" placeholder="请输入" :maxlength="20"
                            class="input-box" @change="item.checked = true"/>
            <a-input v-else v-model:value.trim="item.value" placeholder="请输入" :maxlength="20" class="input-box"
                     @change="item.checked = true"/>
          </div>
        </div>
      </template>
      <template v-else>
        <div class="mt16">
          <a-button type="primary" ghost @click="addData">
            <PlusOutlined/>
            新增元数据
          </a-button>
        </div>
        <div v-if="customList.length" class="meta-card">
          <template v-for="(item, index) in customList" :key="index">
            <div v-if="item.is_edit" class="meta-edit-item">
              <a-input v-model:value.trim="item.name" placeholder="请输入" :maxlength="20"/>
              <a-select v-model:value="item.type" :options="typeData" class="min" :disabled="item.id > 0"/>
              <a-select v-model:value="item.is_show" :options="showData" class="min"/>
              <div class="action-box">
                <a-button @click="cancelEdit(item, index)">取 消</a-button>
                <a-button type="primary" @click="dataChange(item)">保 存</a-button>
              </div>
            </div>
            <div v-else class="meta-item">
              <div class="left">
                <a-tooltip :title="item.is_show == 1 ? '列表显示' : '列表隐藏'">
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
          <div class="meta-tit">内置</div>
          <div v-for="(item, index) in builtList" :key="index" class="meta-item">
            <div class="left">
              <a-tooltip :title="item.is_show == 1 ? '列表显示' : '列表隐藏'">
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
import {
  delLibraryMetaSchema,
  getLibraryMetaSchemaList,
  saveLibraryMetaSchema,
  saveMetadata,
  saveQaMetadata
} from "@/api/library/index.js";

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
  {value: 1, label: '列表显示'},
  {value: 0, label: '列表隐藏'},
]
const typeData = [
  {value: 0, label: 'string'},
  {value: 1, label: 'time'},
  {value: 2, label: 'number'}
]
const typeMap = {
  0: 'string',
  1: 'time',
  2: 'number'
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
    title: '提示',
    content: '确认删除该元数据？',
    okText: '删除',
    cancelText: '取消',
    onOk: () => {
      delLibraryMetaSchema({id: record.id}).then(() => {
        customList.value.splice(index, 1)
        message.success('已删除')
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
  if (!record.name) return message.warning('请输入元数据名称')
  saveLibraryMetaSchema({
    ...record,
    is_show: Number(record.is_show),
    library_id: props.libraryId
  }).then(res => {
    const {name, type, is_show} = record
    record.id = res?.data
    record.is_edit = false
    record.backup = {name, type, is_show}
    message.success('已保存')
    emit('change')
  })
}

function saveValues() {
  // 修改元数据
  if (!customList.value.length) return message.error('暂无可用元数据')
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
  if (!_data.length) return message.error('请选择修改数据')
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
    return message.error('缺少知识库信息')
  }
  req.then(() => {
    message.success('已保存')
    emit('change')
    visible.value = false
  })
}

defineExpose({
  show,
})
</script>
<style lang="less">
.desc {
  color: #8c8c8c;
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
