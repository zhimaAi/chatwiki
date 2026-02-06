<template>
  <div>
    <a-modal
      v-model:open="open"
      :title="modalTitle"
      @cancel="handleClose"
      :footer="null"
      :width="760"
    >
      <div class="custom-table">
        <div class="t-head">
          <div class="td">
            <div style="width: 24px"></div>
            {{ t('column_category_name') }}
          </div>
          <div class="td">{{ t('column_type') }}</div>
          <div class="td">{{ t('column_enabled') }}</div>
          <div class="td flex-none" style="width: 120px">{{ t('column_action') }}</div>
        </div>
        <div class="t-body">
          <div class="t-row">
            <div class="t-item">
              <div style="width: 24px"></div>
              {{ t('text_all') }}
            </div>
            <div class="t-item">{{ t('type_system') }}</div>
            <div class="t-item">
              <a-switch :checked="true" disabled :checked-children="t('switch_on')" :un-checked-children="t('switch_off')" />
            </div>
            <div class="t-item flex-none" style="width: 120px"></div>
          </div>
          <draggable
            v-model="data"
            item-key="id"
            @end="onDragEnd"
            group="table-rows"
            handle=".drag-btn"
          >
            <template #item="{ element, index }">
              <div :key="element.id" class="t-row">
                <div class="t-item">
                  <span class="drag-btn"><svg-icon name="drag" /></span>
                  <span v-if="element.name.length > 10">
                    <a-tooltip>
                      <template #title>{{ element.name }}</template>
                      {{ element.name.slice(0, 10) + '...' }}
                    </a-tooltip>
                  </span>
                  <span v-else>{{ element.name }} </span>
                </div>
                <div class="t-item">{{ t('type_custom') }}</div>
                <div class="t-item">
                  <a-switch
                    @change="handleChangeEnabled(element)"
                    :checked="element.enabled"
                    :checked-children="t('switch_on')"
                    :un-checked-children="t('switch_off')"
                    checkedValue="true"
                    unCheckedValue="false"
                  />
                </div>
                <div class="t-item flex-none" style="width: 120px">
                  <a @click="handleOpenAddModal(element)">{{ t('action_edit') }}</a>
                  <a-divider type="vertical" />
                  <a @click="onDelete(element, index)">{{ t('action_delete') }}</a>
                </div>
              </div>
            </template>
          </draggable>
        </div>
      </div>
    </a-modal>
    <AddFilrerModal
      @ok="handleEditBack"
      :column="props.column"
      ref="addFilrerModalRef"
    ></AddFilrerModal>
  </div>
</template>

<script setup>
import { ref, createVNode } from 'vue'
import draggable from 'vuedraggable'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { message, Modal } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import {
  getFormFilterList,
  delFormFilter,
  updateFormFilterEnabled,
  updateFormFilterSort
} from '@/api/database'
import AddFilrerModal from './add-filter-modal.vue'
import { useRoute } from 'vue-router'

const { t } = useI18n('views.database.database-detail.database-manage.components.filter-manage-modal')

const rotue = useRoute()
const query = rotue.query
const emit = defineEmits(['ok', 'change'])
const modalTitle = ref(t('modal_title'))
const open = ref(false)

const props = defineProps({
  column: {
    type: Array,
    default: () => []
  }
})
const getSortLists = () => {
  getFormFilterList({ form_id: query.form_id }).then((res) => {
    data.value = res.data.map((item) => {
      return {
        ...item
      }
    })
  })
}
const handleEditBack = () => {
  getSortLists()
  emit('ok')
  emit('change')
}
const show = () => {
  getSortLists()
  open.value = true
}

const handleClose = () => {
  open.value = false
}
const data = ref()
const onDragEnd = () => {
  let filter_sort = []
  data.value.reverse().forEach((item, index) => {
    filter_sort.push({
      id: +item.id,
      sort: +index + 1
    })
  })
  updateFormFilterSort({
    form_id: query.form_id,
    filter_sort: JSON.stringify(filter_sort)
  }).then((res) => {
    message.success(t('msg_sort_success'))
    getSortLists()
    emit('change')
  })
}

const handleChangeEnabled = (record) => {
  updateFormFilterEnabled({
    form_id: query.form_id,
    id: record.id,
    enabled: record.enabled == 'true' ? 'false' : 'true'
  }).then((res) => {
    message.success(t('msg_set_success'))
    getSortLists()
    emit('change')
  })
}
const onDelete = (record, index) => {
  Modal.confirm({
    title: t('modal_delete_title'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('modal_delete_content'),
    okText: t('modal_delete_ok'),
    okType: 'danger',
    cancelText: t('modal_delete_cancel'),
    onOk() {
      delFormFilter({ id: record.id, form_id: query.form_id }).then((res) => {
        message.success(t('msg_delete_success'))
        data.value.splice(index, 1)
        emit('change')
      })
    },
    onCancel() {}
  })
}

const addFilrerModalRef = ref(null)
const handleOpenAddModal = (data) => {
  addFilrerModalRef.value.edit(data)
}
const handleOk = () => {}
defineExpose({
  show
})
</script>

<style lang="less" scoped>
.custom-table {
  .drag-btn {
    margin-right: 8px;
    cursor: grab;
    opacity: 0;
    transition: opacity 0.2s;
  }
  margin-bottom: 24px;
  margin-top: 24px;
  .t-head {
    border-radius: 8px 8px 0 0;
    background: #f5f5f5;
    border-bottom: 1px solid #f0f0f0;
    font-weight: 600;
    text-align: left;
    display: flex;
    align-items: center;
    .td {
      display: flex;
      flex: 1;
      font-weight: 600;
      padding: 12px 16px;
    }
    .td:first-of-type {
      flex: 1.3;
    }
    .flex-none {
      flex: none;
    }
  }
  .t-body {
    .t-row {
      display: flex;
      align-items: center;
      border-bottom: 1px solid #f0f0f0;
      &:hover {
        background: #fafafa;
        .drag-btn {
          opacity: 1;
        }
      }
      .t-item {
        display: flex;
        flex: 1;
        padding: 12px 16px;
      }
      .t-item:first-of-type {
        flex: 1.3;
      }
      .flex-none {
        flex: none;
      }
    }
  }
}
</style>
