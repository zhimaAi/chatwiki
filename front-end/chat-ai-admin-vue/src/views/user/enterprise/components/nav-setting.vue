<template>
  <div class="content-box">
    <div class="custom-table">
      <div class="t-head">
        <div class="td">
          <span style="padding-left: 22px">{{ t('sort') }}</span>
          <a-tooltip :title="t('sortTip')">
            <QuestionCircleOutlined style="margin-left: 4px" />
          </a-tooltip>
        </div>
        <div class="td">
          {{ t('showInTop') }}
          <a-tooltip
            :title="t('showInTopTip')"
          >
            <QuestionCircleOutlined style="margin-left: 4px" />
          </a-tooltip>
        </div>
      </div>
      <div class="t-body">
        <draggable
          v-model="navMenuList"
          item-key="id"
          @end="onDragEnd"
          group="table-rows"
          handle=".drag-btn"
        >
          <template #item="{ element, index }">
            <div :key="element.id" class="t-row">
              <div class="t-item">
                <span class="drag-btn">
                  <img src="@/assets/svg/drag.svg" alt="" />
                </span>
                <span v-if="element.name.length > 30">
                  <a-tooltip>
                    <template #title>{{ element.name }}</template>
                    {{ element.name.slice(0, 30) + '...' }}
                  </a-tooltip>
                </span>
                <span v-else>{{ element.name }} </span>
              </div>
              <div class="t-item flex-center">
                <a-switch
                  :disabled="element.isDisabled"
                  @change="handleSave"
                  v-model:checked="element.open"
                  :checked-children="t('on')"
                  :un-checked-children="t('off')"
                />
              </div>
            </div>
          </template>
        </draggable>
      </div>
    </div>
    <a-empty v-if="navMenuList.length == 0" style="margin: 108px 0 24px 0" :image="simpleImage" />
  </div>
</template>

<script setup>
import { reactive, ref, computed, watch } from 'vue'
import { message, Empty } from 'ant-design-vue'
import draggable from 'vuedraggable'
import { saveTopNavigate } from '@/api/user/index.js'
import { QuestionCircleOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { useCompanyStore } from '@/stores/modules/company'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.enterprise')
const companyStore = useCompanyStore()
const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE

const emit = defineEmits(['handleGetCompany'])

const navMenuList = ref([])

const top_navigate = computed(() => {
  return companyStore.top_navigate
})

let topNavigateDefaultData = companyStore.topNavigateDefaultData

watch(
  () => top_navigate.value,
  (val) => {
    let lists = val || []
    navMenuList.value = lists.map((item) => {
      let currentItem = topNavigateDefaultData.find((it) => it.id == item.id)
      return {
        ...currentItem,
        ...item
      }
    })
  },
  {
    deep: true,
    immediate: true
  }
)

const onDragEnd = () => {
  handleSave()
}

const handleSave = () => {
  saveTopNavigate({
    id: companyStore.id,
    top_navigate: JSON.stringify(navMenuList.value)
  }).then((res) => {
    message.success(t('saveSuccess'))
    emit('handleGetCompany')
  })
}
</script>

<style lang="less" scoped>
.content-box {
  width: 600px;
}
.custom-table {
  width: 100%;
  .drag-btn {
    margin-right: 8px;
    cursor: grab;
    margin-top: 2px;
    transition: opacity 0.2s;
    img {
      width: 15px;
      height: 15px;
    }
  }

  .t-head {
    border-radius: 8px 8px 0 0;
    background: #f5f5f5;
    border-bottom: 1px solid #f0f0f0;
    font-weight: 600;
    text-align: left;
    display: flex;
    align-items: center;
    color: #262626;
    .td {
      display: flex;
      flex: 1;
      font-weight: 600;
      padding: 12px 16px;
    }
  }
  .t-body {
    color: #595959;
    .t-row {
      display: flex;
      align-items: center;
      border-bottom: 1px solid #f0f0f0;
      &:hover {
        background: #fafafa;
      }
      .t-item {
        display: flex;
        flex: 1;
        padding: 12px 16px;
      }
      .flex-center {
        display: flex;
        align-items: center;
      }
    }
  }
}
</style>
