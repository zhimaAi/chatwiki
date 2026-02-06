<template>
  <div>
    <a-modal v-model:open="open" :title="t('title_selection_settings')" @ok="handleOk" :width="586">
      <div class="modal-box">
        <!-- <a-alert show-icon message="添加了分类标记的分段，重新分段时将会会被保留。"></a-alert> -->
        <div class="list-box">
          <a-table :data-source="tableData" :pagination="false">
            <a-table-column key="tags" :title="t('label_icon')" data-index="tags" :width="200">
              <template #default="{ record }">
                <StarFilled class="start-icon" :style="{ color: colorLists[record.type] }" />
              </template>
            </a-table-column>
            <a-table-column key="name" :title="t('label_remark')" data-index="name" :width="200">
              <template #default="{ record }">
                <a-input v-model:value="record.name" :maxLength="4" :placeholder="t('ph_set_max_4_chars')" />
              </template>
            </a-table-column>
          </a-table>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { StarOutlined, StarFilled } from '@ant-design/icons-vue'
import { getCategoryList, saveCategory } from '@/api/library'
import colorLists from '@/utils/starColors.js'
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.library-preview.components.classification-mark-modal')

const open = ref(false)
const emit = defineEmits(['ok'])

const tableData = ref([])
const show = () => {
  getCategoryLists()
  open.value = true
}

const getCategoryLists = () => {
  getCategoryList().then((res) => {
    tableData.value = res.data
  })
}

const handleOk = () => {
  let data = tableData.value.map((item) => {
    return {
      id: +item.id,
      name: item.name
    }
  })
  saveCategory({ data: JSON.stringify(data) }).then((res) => {
    message.success(t('msg_save_success'))
    emit('ok')
    open.value = false
  })
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.modal-box {
  margin-top: 24px;
}
.list-box {
  margin: 16px 0;
}
.start-icon {
  font-size: 20px;
}
</style>
