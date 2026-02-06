<template>
  <div>
    <a-modal v-model:open="open" :title="t('modal_title')" :width="624" :footer="null">
      <a-table :data-source="list" :loading="loading" :pagination="false" :scroll="{ y: 500 }">
        <a-table-column key="index" data-index="index" :title="t('table.rank')" :width="100">
          <template #default="{ index }">
            {{ index + 1 }}
          </template>
        </a-table-column>

        <a-table-column key="robot_name" :title="t('table.robot')" :width="140">
          <template #default="{ record }">
            {{ record.robot_name }}
          </template>
        </a-table-column>
        <a-table-column key="tip" :title="t('table.trigger_count')" :width="120">
          <template #default="{ record }">
            <a-flex :gap="12">
              <span>{{ record.tip }}</span>
            </a-flex>
          </template>
        </a-table-column>
      </a-table>
    </a-modal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import { statLibraryDataRobotDetail, statLibraryRobotDetail } from '@/api/library'

const { t } = useI18n('views.trigger-statics.list.components.detail-modal')

const open = ref(false)
const loading = ref(false)
const list = ref([])
const show = (item, type) => {
  open.value = true
  loading.value = true
  if (type == 1) {
    getContentDetail(item)
  } else {
    getLibraryDetail(item)
  }
}

const getContentDetail = (data) => {
  statLibraryDataRobotDetail({
    begin_date_ymd: data.begin_date_ymd,
    end_date_ymd: data.end_date_ymd,
    library_id: data.library_id,
    data_id: data.data_id,
  })
    .then((res) => {
      list.value = res.data || []
    })
    .finally(() => {
      loading.value = false
    })
}

const getLibraryDetail = (data) => {
  statLibraryRobotDetail({
    begin_date_ymd: data.begin_date_ymd,
    end_date_ymd: data.end_date_ymd,
    library_id: data.library_id,
    group_id: data.group_id
  })
    .then((res) => {
      list.value = res.data || []
    })
    .finally(() => {
      loading.value = false
    })
}

defineExpose({
  show
})
</script>

<style lang="less" scoped></style>
