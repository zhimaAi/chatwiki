<template>
  <div class="content-box">
    <div class="search-block" @click="handleOpenModal">
      <a-input-search
        v-model:value="searchState.name"
        placeholder="搜索名称"
        style="width: 220px"
        @search="onSearch"
      />
    </div>
    <div class="list-box">
      <a-table
        :data-source="tableData"
        :pagination="{
          current: searchState.page,
          total: total,
          pageSize: searchState.size,
          showQuickJumper: true,
          showSizeChanger: true,
          pageSizeOptions: ['10', '20', '50', '100']
        }"
        @change="onTableChange"
      >
        <a-table-column title="类型" data-index="IP" :width="100">
          <template #default="{ record }">{{ record.type }}</template>
        </a-table-column>
        <a-table-column title="应用名称" data-index="model" :width="100"> </a-table-column>
        <a-table-column title="今日消耗(k)" data-index="amount1" :width="140">
          <template #default="{ record }">{{ record.amount }}</template>
        </a-table-column>
        <a-table-column title="Token消耗(k)" data-index="amount" :width="140">
          <template #default="{ record }">{{ record.amount }}</template>
        </a-table-column>
        <a-table-column title="已消耗(k)" data-index="amount2" :width="140">
          <template #default="{ record }">{{ record.amount }}</template>
        </a-table-column>
        <a-table-column title="剩余(k)" data-index="amount3" :width="140">
          <template #default="{ record }">{{ record.amount }}</template>
        </a-table-column>
        <a-table-column title="备注" data-index="date" :width="200">
          <template #default="{ record }">{{ record.date }}</template>
        </a-table-column>
        <a-table-column title="操作" data-index="action" :width="120">
          <template #default="{ record }">{{ record.amount }}</template>
        </a-table-column>
      </a-table>
    </div>
    <TokenQuotaModal ref="tokenQuotaModalRef" @ok="getData" />
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import TokenQuotaModal from './components/token-quota-modal.vue'
const searchState = reactive({
  name: '',
  page: 1,
  size: 10
})

const total = ref(0)

const onSearch = () => {
  searchState.page = 1
}

const onTableChange = (pagination) => {
  searchState.page = pagination.current
  searchState.size = pagination.pageSize
  getData()
}
const getData = () => {}

const tokenQuotaModalRef = ref(null)
const handleOpenModal = ()=>{
  tokenQuotaModalRef.value.show()
}

</script>

<style lang="less" scoped>
.content-box {
  padding: 0 16px;
  .search-block {
    display: flex;
    align-items: center;
    gap: 16px;
  }
  .list-box {
    margin-top: 16px;
  }
}
</style>
