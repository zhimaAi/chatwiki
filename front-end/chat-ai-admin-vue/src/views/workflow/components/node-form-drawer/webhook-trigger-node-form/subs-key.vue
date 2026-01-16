<template>
  <div
    :class="['sub-field-item-box', `sub-field-item-box_${props.level}`]"
    v-for="(item, index) in props.data"
    :key="item.cu_key"
  >
    <div class="sub-field-item">
      <img class="fork-icon" src="@/assets/svg/right-fork-icon.svg" alt="" />
      <div class="key-title-item">{{ item.key }}</div>
      <div class="typ-title-item">{{ item.typ }}</div>
      <div class="value-title-item">
        <a-select
          v-if="item.typ != 'object'"
          placeholder="请输入选择变量"
          v-model:value="item.desc"
          style="width: 90%"
        >
          <a-select-option :value="opt.value" v-for="opt in handleFilterOptions(item)" :key="opt.key">
            <span>{{ opt.label }}</span>
          </a-select-option>
        </a-select>
      </div>
      <div
        class="btn-hover-wrap"
        @click="handleAddOutPut(index, props.data)"
        v-if="showAddBtn(item)"
      >
        <PlusCircleOutlined />
      </div>
      <div class="btn-hover-wrap" @click="handleEditOutput(item, index)">
        <EditOutlined />
      </div>

      <div class="btn-hover-wrap" @click="onDelOutput(index, props.data)">
        <CloseCircleOutlined />
      </div>
    </div>
    <template v-if="item.subs && item.subs.length > 0">
      <SubsKey
        :request_content_type="props.request_content_type"
        :level="props.level + 1"
        :data="item.subs"
        :globalOptions="props.globalOptions"
      />
    </template>
  </div>
  <AddKeyModal
    @add="onOutputAdd"
    :level="props.level"
    :request_content_type="props.request_content_type"
    @edit="onOutputEdit"
    ref="addParamsModalRef"
  />
</template>

<script setup>
import { computed, ref } from 'vue'
import { CloseCircleOutlined, PlusCircleOutlined, EditOutlined } from '@ant-design/icons-vue'
import AddKeyModal from './add-key-modal.vue'
const props = defineProps({
  data: {},
  level: {
    type: Number
  },
  globalOptions: {
    type: Array
  },
  request_content_type: {
    type: String,
    default: ''
  }
})
const filterTypOptions = computed(() => {
  if (props.level > 2) {
    return typOptions.filter((item) => item.value != 'object')
  }
  return typOptions
})

let currentIndex = 0

const onDelOutput = (index, data) => {
  data.splice(index, 1)
}

const addParamsModalRef = ref(null)

const handleAddOutPut = (index) => {
  currentIndex = index
  addParamsModalRef.value.add()
}
const handleEditOutput = (data, index) => {
  addParamsModalRef.value.edit(data, index)
}
const onOutputAdd = (data) => {
  props.data[currentIndex].subs.push(data)
}

const onOutputEdit = (data, index) => {
  if(props.data[index].typ != data.typ){
    data.desc = void 0
  }
  props.data.splice(index, 1, data)
}

const showAddBtn = (item) => {
  return item.typ == 'object' && props.level <= 2
}

const handleFilterOptions = (item) => {
  let typ = item.typ
  if (typ == 'file') {
    typ = 'string'
  }
  if (typ) {
    return props.globalOptions.filter((item) => item.typ == typ)
  }
  return props.globalOptions
}

let typOptions = [
  {
    lable: 'string',
    value: 'string'
  },
  {
    lable: 'number',
    value: 'number'
  },
  {
    lable: 'boole',
    value: 'boole'
  },
  {
    lable: 'float',
    value: 'float'
  },
  {
    lable: 'object',
    value: 'object'
  },
  {
    lable: 'array\<string>',
    value: 'array\<string>'
  },
  {
    lable: 'array\<number>',
    value: 'array\<number>'
  },
  {
    lable: 'array\<boole>',
    value: 'array\<boole>'
  },
  {
    lable: 'array\<float>',
    value: 'array\<float>'
  },
  {
    lable: 'array\<object>',
    value: 'array\<object>'
  }
]
</script>

<style lang="less" scoped>
.sub-field-item {
  display: flex;
  align-items: center;
  margin-top: 6px;
  .btn-hover-wrap {
    margin-right: 4px;
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
.sub-field-item-box_3 .key-title-item {
  width: 114px;
}

.key-title-item {
  width: 148px;
}

.typ-title-item {
  width: 90px;
}
.value-title-item {
  width: 160px;
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
}
</style>
