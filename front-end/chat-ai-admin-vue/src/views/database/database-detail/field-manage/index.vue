<template>
  <div class="field-manage-page">
    <div class="page-title">
      字段管理
      <a-button type="primary" @click="handleAddField()">
        <template #icon>
          <PlusOutlined />
        </template>
        添加字段
      </a-button>
    </div>
    <div class="table-wrapper customize-scroll-style">
      <a-table :pagination="false" :data-source="data">
        <a-table-column key="name" data-index="name" :width="150">
          <template #title>
            字段名称
            <a-tooltip>
              <template #title>定义数据表的表头，可以在对应表头下存储相关数据。</template>
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
            字段描述
            <a-tooltip>
              <template #title>表头字段的说明，帮助用户或大模型理解表头字段</template>
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
            数据类型
            <a-tooltip>
              <template #title>选择存储字段对应的数据类型</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>
        </a-table-column>
        <a-table-column key="required" data-index="required" :width="120">
          <template #title>
            是否必要
            <a-tooltip>
              <template #title>
                <div>必要字段：在保存一行数据时，必须提供对应字段信息，否则无法保存该行数据</div>
                <div>非必要字段：缺失该字段信息时，一行数据仍可被保存在表中</div>
              </template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>
          <template #default="{ record }">
            <a-switch
              @change="handleChangeRequired(record)"
              :checked="record.required == 'true'"
              checked-children="开"
              un-checked-children="关"
            />
          </template>
        </a-table-column>
        <a-table-column key="action" title="操作" :width="90">
          <template #default="{ record }">
            <span>
              <a @click="handleAddField(record)">编辑</a>
              <a-divider type="vertical" />
              <a @click="onDelete(record)">删除</a>
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
import { getFormFieldList, delFormField, updateFormRequired } from '@/api/database'
import AddFieldModal from './components/add-field-modal.vue'

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
    message.success('修改成功')
    getData()
  })
}

const onDelete = (record) => {
  Modal.confirm({
    title: `删除确认`,
    icon: createVNode(ExclamationCircleOutlined),
    content: `删除后，字段下所有数据将一并删除，删除后不可恢复。确认删除字段${record.name}吗?`,
    okText: '确 定',
    okType: 'danger',
    cancelText: '取 消',
    onOk() {
      delFormField({ id: record.id }).then((res) => {
        message.success('删除成功')
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
