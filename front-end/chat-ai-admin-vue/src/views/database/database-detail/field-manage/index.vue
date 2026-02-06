<template>
  <div class="field-manage-page">
    <div class="page-title">
      {{ t('title') }}
      <a-button type="primary" @click="handleAddField()">
        <template #icon>
          <PlusOutlined />
        </template>
        {{ t('btn_add_field') }}
      </a-button>
    </div>
    <div class="table-wrapper customize-scroll-style">
      <a-table :pagination="false" :data-source="data">
        <a-table-column key="name" data-index="name" :width="150">
          <template #title>
            {{ t('column_name_title') }}
            <a-tooltip>
              <template #title>{{ t('column_name_tooltip') }}</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>
          <template #default="{ record }">
            <span v-if="record.name.length > 50">
              <a-tooltip>
                <template #title>{{ record.name }}</template>
                {{ record.name.slice(0, 50) + '...' }}
              </a-tooltip>
            </span>
            <span v-else>{{ record.name }} </span>
          </template>
        </a-table-column>
        <a-table-column key="description" data-index="description" :width="150">
          <template #title>
            {{ t('column_description_title') }}
            <a-tooltip>
              <template #title>{{ t('column_description_tooltip') }}</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>
          <template #default="{ record }">
            <span v-if="record.description.length > 50">
              <a-tooltip>
                <template #title>{{ record.description }}</template>
                {{ record.description.slice(0, 50) + '...' }}
              </a-tooltip>
            </span>
            <span v-else>{{ record.description }} </span>
          </template>
        </a-table-column>
        <a-table-column key="type" data-index="type" :width="120">
          <template #title>
            {{ t('column_type_title') }}
            <a-tooltip>
              <template #title>{{ t('column_type_tooltip') }}</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>
        </a-table-column>
        <a-table-column key="required" data-index="required" :width="120">
          <template #title>
            {{ t('column_required_title') }}
            <a-tooltip>
              <template #title>
                <div>{{ t('column_required_tooltip_required') }}</div>
                <div>{{ t('column_required_tooltip_optional') }}</div>
              </template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>
          <template #default="{ record }">
            <a-switch
              @change="handleChangeRequired(record)"
              :checked="record.required == 'true'"
              :checked-children="t('switch_on')"
              :un-checked-children="t('switch_off')"
            />
          </template>
        </a-table-column>
        <a-table-column key="action" :title="t('column_action_title')" :width="90">
          <template #default="{ record }">
            <span>
              <a @click="handleAddField(record)">{{ t('action_edit') }}</a>
              <a-divider type="vertical" />
              <a @click="onDelete(record)">{{ t('action_delete') }}</a>
            </span>
          </template>
        </a-table-column>
      </a-table>
    </div>
  </div>
  <AddFieldModal @ok="getData" ref="addFieldModalRef"></AddFieldModal>
</template>

<script setup>
import {
  PlusOutlined,
  QuestionCircleOutlined,
  ExclamationCircleOutlined
} from '@ant-design/icons-vue'
import { ref, reactive, createVNode } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { useRoute } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'
import { getFormFieldList, delFormField, updateFormRequired } from '@/api/database'
import AddFieldModal from './components/add-field-modal.vue'

const { t } = useI18n('views.database.database-detail.field-manage.index')

const rotue = useRoute()
const query = rotue.query

const data = ref([])
const getData = () => {
  getFormFieldList({ form_id: query.form_id }).then((res) => {
    data.value = res.data || []
  })
}
getData()

const handleChangeRequired = (record) => {
  updateFormRequired({
    id: record.id,
    form_id: query.form_id,
    required: record.required == 'true' ? 'false' : 'true'
  }).then((res) => {
    message.success(t('msg_modify_success'))
    getData()
  })
}

const onDelete = (record) => {
  Modal.confirm({
    title: t('modal_delete_title'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('modal_delete_content', { name: record.name }),
    okText: t('modal_delete_ok'),
    okType: 'danger',
    cancelText: t('modal_delete_cancel'),
    onOk() {
      delFormField({ id: record.id }).then((res) => {
        message.success(t('msg_delete_success'))
        getData()
      })
    },
    onCancel() {}
  })
}

const addFieldModalRef = ref(null)
const handleAddField = (data = {}) => {
  addFieldModalRef.value.show(JSON.parse(JSON.stringify(data)))
}
</script>

<style lang="less" scoped>
.field-manage-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  .page-title {
    display: flex;
    align-items: center;
    background-color: #fff;
    color: #000000;
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
    justify-content: space-between;
    padding-right: 24px;
  }

  .table-wrapper {
    margin-top: 24px;
    padding-right: 24px;
    flex: 1;
    overflow: auto;
  }
  ::v-deep(.ant-table-cell) {
    word-break: break-all;
  }
}
</style>
