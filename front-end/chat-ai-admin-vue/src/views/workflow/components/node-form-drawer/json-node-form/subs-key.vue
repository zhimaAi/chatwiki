<template>
  <div
    :class="['sub-field-item-box', `sub-field-item-box_${props.level}`]"
    v-for="(item, index) in props.data"
    :key="item.cu_key"
  >
    <div class="sub-field-item">
      <img class="fork-icon" src="@/assets/svg/right-fork-icon.svg" alt="" />
      <a-input
        class="mr12"
        :style="{ width: widthPx }"
        v-model:value="item.key"
        :placeholder="t('ph_input')"
      ></a-input>
      <img class="fork-icon" src="@/assets/svg/right-fork-icon.svg" alt="" />
      <a-select
        class="mr12"
        @change="onTypeChange(item)"
        v-model:value="item.typ"
        :placeholder="t('ph_select')"
        :style="{ width: widthPx }"
      >
        <a-select-option v-for="op in filterTypOptions" :value="op.value">{{
          op.value
        }}</a-select-option>
      </a-select>
      <div
        class="btn-hover-wrap mr12"
        @click="onAddSubs(index, props.data)"
        v-if="showAddBtn(item)"
      >
        <PlusCircleOutlined />
      </div>

      <div class="btn-hover-wrap" @click="onDelOutput(index, props.data)">
        <CloseCircleOutlined />
      </div>
    </div>
    <template v-if="item.subs && item.subs.length > 0">
      <SubsKey :level="props.level + 1" :typOptions="props.typOptions" :data="item.subs" />
    </template>
  </div>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'
import { computed } from 'vue'
import { CloseCircleOutlined, PlusCircleOutlined } from '@ant-design/icons-vue'

const { t } = useI18n('views.workflow.components.node-form-drawer.json-node-form.subs-key')
const props = defineProps({
  data: {},
  level: {
    type: Number
  },
  typOptions: {
    type: Array
  }
})

const filterTypOptions = computed(() => {
  if (props.level > 2) {
    return props.typOptions.filter((item) => item.value != 'object')
  }
  return props.typOptions
})

const widthPx = computed(() => {
  if (props.level == 2) {
    return '180px'
  }
  return '146px'
})

const onAddSubs = (index, data) => {
  data[index].subs.push({
    key: '',
    value: '',
    subs: [],
    cu_key: Math.random() * 10000
  })
}

const onDelOutput = (index, data) => {
  data.splice(index, 1)
}

const onTypeChange = (data) => {
  data.subs = []
}

const showAddBtn = (item) => {
  return item.typ == 'object' && props.level <= 2
}
</script>

<style lang="less" scoped>
.sub-field-item {
  display: flex;
  align-items: center;
  margin-top: 4px;
  .btn-hover-wrap {
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s ease-in;
    width: 24px;
    height: 24px;
    &:hover {
      background: #e4e6eb;
    }
  }
}
.mr12 {
  margin-right: 12px;
}
.sub-field-item-box_3 .fork-icon {
  margin-left: 34px;
}
.fork-icon {
  width: 32px;
  height: 32px;
  margin-right: 2px;
}
</style>
