<template>
  <div
    :class="['sub-field-item-box', `sub-field-item-box_${props.level}`]"
    v-for="(item, index) in props.data"
    :key="item.cu_key"
  >
    <div class="sub-field-item">
      <a-flex>
        <img class="fork-icon" src="@/assets/svg/right-fork-icon.svg" alt="" />
        <a-input
          :style="{ width: widthPx }"
          v-model:value="item.key"
          placeholder="请输入"
        ></a-input>
      </a-flex>
      <!-- <img class="fork-icon" src="@/assets/svg/right-fork-icon.svg" alt="" /> -->
      <a-select
        @change="onTypeChange(item)"
        v-model:value="item.typ"
        placeholder="请选择"
        style="width: 114px"
      >
        <a-select-option v-for="op in filterTypOptions" :value="op.value">{{
          op.value
        }}</a-select-option>
      </a-select>
      <div style="width: 200px">
        <at-input
          v-if="item.typ != 'object'"
          :options="variableOptions"
          :defaultSelectedList="item.tags"
          :defaultValue="item.desc"
          @change="(text, selectedList) => changeOutputValue(text, selectedList, item, index)"
          placeholder="请输入变量值，键入“/”插入变量"
        >
          <template #option="{ label, payload }">
            <div class="field-list-item">
              <div class="field-label">{{ label }}</div>
              <div class="field-type">{{ payload.typ }}</div>
            </div>
          </template>
        </at-input>
      </div>
      <div class="btn-hover-wrap" @click="onAddSubs(index, props.data)" v-if="showAddBtn(item)">
        <PlusCircleOutlined />
      </div>
      <div class="btn-hover-wrap" v-else style="opacity: 0">
        <PlusCircleOutlined />
      </div>

      <div class="btn-hover-wrap" @click="onDelOutput(index, props.data)">
        <CloseCircleOutlined />
      </div>
    </div>
    <template v-if="item.subs && item.subs.length > 0">
      <SubsKey
        :level="props.level + 1"
        :typOptions="props.typOptions"
        :variableOptions="props.variableOptions"
        :data="item.subs"
      />
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { CloseCircleOutlined, PlusCircleOutlined } from '@ant-design/icons-vue'
import AtInput from '../../at-input/at-input.vue'
const props = defineProps({
  data: {},
  level: {
    type: Number
  },
  typOptions: {
    type: Array
  },
  variableOptions: {
    type: Array,
    default: () => []
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
    return '165px'
  }
  return '130px'
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

const changeOutputValue = (text, selectedList, item) => {
  item.tags = selectedList
  item.desc = text
}
</script>

<style lang="less" scoped>
.sub-field-item {
  display: flex;
  align-items: center;
  margin-top: 4px;
  gap: 8px;
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
.mr8 {
  margin-right: 8px;
}
.sub-field-item-box_3 .fork-icon {
  margin-left: 34px;
}
.fork-icon {
  width: 32px;
  height: 32px;
}
</style>
