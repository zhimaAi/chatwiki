<template>
  <div class="form-block" @mousedown.stop="">
    <a-form ref="formRef" layout="vertical" :model="formState">
      <draggable
        style="display: flex; flex-direction: column; gap: 8px"
        handle=".drag-btn"
        v-model="formState.term"
        item-key="key"
      >
        <template #item="{ element: item, index }">
          <div class="gray-block" :key="item.key">
            <div class="gray-block-title">
              <a-flex :gap="8"
                ><HolderOutlined class="icon drag-btn" />{{ index == 0 ? 'if' : 'else if' }}
              </a-flex>
              <div
                v-if="formState.term.length > 1"
                class="btn-hover-wrap"
                @click="handleDelBranch(index)"
              >
                <CloseCircleOutlined />
              </div>
            </div>
            <div class="condition-list-box">
              <div class="left-select-box">
                <a-select
                  size="small"
                  v-model:value="item.is_or"
                  :bordered="false"
                  style="width: 64px"
                >
                  <a-select-option :value="0">and</a-select-option>
                  <a-select-option :value="1">or</a-select-option>
                </a-select>
              </div>
              <div class="condition-body">
                <div class="condition-item" v-for="(term, i) in item.terms" :key="term.key">
                  <!-- <a-select
                    placeholder="请选择"
                    v-model:value="term.variable"
                    style="width: 100px"
                    @change="handleVariableChange(term)"
                  >
                    <a-select-option v-for="option in variableOptions" :value="option.key">{{
                      option.label || option.key
                    }}</a-select-option>
                  </a-select> -->
                  <a-cascader
                    v-model:value="term.variable"
                    @change="handleVariableChange(term)"
                    @dropdownVisibleChange="onDropdownVisibleChange"
                    style="width: 180px"
                    :options="variableOptions"
                    :allowClear="false"
                    :displayRender="({ labels }) => labels.join('/')"
                    :field-names="{ children: 'children' }"
                    placeholder="请选择"
                  />
                  <a-select v-model:value="term.type" style="width: 130px" placeholder="请选择">
                    <a-select-option v-for="option in getTypeOptions(term)" :value="option.value" :key="option.value">{{
                      option.label
                    }}</a-select-option>
                  </a-select>
                  <a-input
                    v-if="term.type != 5 && term.type != 6"
                    placeholder="请输入"
                    v-model:value="term.value"
                    style="width: 170px"
                  ></a-input>
                  <div class="btn-hover-wrap" @click="handleDelCondition(index, i)">
                    <CloseCircleOutlined />
                  </div>
                </div>
              </div>
            </div>
            <div class="btn-wrap">
              <a-button
                @click="handleAddCondition(index)"
                :icon="h(PlusOutlined)"
                block
                type="dashed"
                >添加条件</a-button
              >
            </div>
          </div>
        </template>
      </draggable>

      <div class="add-btn-block">
        <a-button @click="handleAddBranch" :icon="h(PlusOutlined)" block type="dashed"
          >添加分支</a-button
        >
      </div>
      <div class="gray-block mt8">
        <div class="gray-block-title">else</div>
        <div class="main-text">不符合上述所有分支的条件时，走默认分支</div>
      </div>
    </a-form>
  </div>
</template>

<script setup>
import { ref, reactive, watch, h, inject, onMounted } from 'vue'
import draggable from 'vuedraggable'
import {
  CloseCircleOutlined,
  PlusOutlined,
  HolderOutlined,
} from '@ant-design/icons-vue'
import { specialNodeList } from '@/views/workflow/components/util.js'

const props = defineProps({
  properties: {
    type: Object,
    default: () => ({})
  }
})

// const graphModel = inject('getGraph')
const getNode = inject('getNode')

const emit = defineEmits(['setData'])
const formRef = ref()

const formState = reactive({
  term: []
})

const variableOptions = ref([])


function getOptions(){
  let list = getNode().getAllParentVariable()

  variableOptions.value = handleOptions(list)
}

// 递归处理Options
function handleOptions(options) {
  options.forEach((item) => {

    if(item.typ == 'node'){
      if(item.node_type == 1){
        item.value = 'global'
      }else{
        item.value = item.node_id;
      }
    }else{
      item.value = item.key
    }

    if (item.children && item.children.length > 0) {
      item.children = handleOptions(item.children)
    }
  })

  return options
}

const onDropdownVisibleChange = (visible) => {
  if (!visible) {
    getOptions()
  }
}


let lock = false

watch(
  () => props.properties,
  (val) => {
    try {
      if (lock) {
        return
      }

      let term = JSON.parse(val.node_params).term || []

      term = term.map((item) => {
        let terms = item.terms.map((it) => {
          // 判断是不是特殊节点
          let specialKey = '';

          for(let i = 0; i < specialNodeList.length; i++){
            if(it.variable.indexOf(specialNodeList[i]) > -1){
              specialKey = specialNodeList[i]
              break;
            }
          }

          if(specialKey != ''){
            let arr = it.variable.split('.')
            it.variable = [arr[0], specialKey]
          }else{
            it.variable = it.variable.split('.')
          }
          return {
            ...it,
            type: it.type > 0 ? it.type : 1,
            key: Math.random() * 10000
          }
        })

        return {
          ...item,
          is_or: item.is_or ? 1 : 0,
          terms,
          key: Math.random() * 10000
        }
      })
      formState.term = term
      lock = true
      setTimeout(() => {
        emit('setData', {
          ...formState,
          height: getNodeHeight()
        })
      }, 100)
    } catch (error) {
      console.log(error)
    }
  },
  { immediate: true, deep: true }
)

watch(
  () => formState,
  () => {
    let term = JSON.parse(JSON.stringify(formState.term))

    term = term.map((item) => {
      return {
        is_or: item.is_or == 1,
        next_node_key: item.next_node_key,
        terms: item.terms.map((it) => {
          let variable = it.variable.join('.')
          return {
            variable: variable,
            is_mult: it.is_mult,
            type: it.type,
            value: it.value
          }
        })
      }
    })
    emit('setData', {
      ...formState,
      term,
      node_params: JSON.stringify({
        term: term
      }),
      height: getNodeHeight()
    })
  },
  { deep: true }
)

function getNodeHeight() {
  let len = 0
  let termLen = formState.term.length
  formState.term.forEach((item) => {
    len += item.terms.length
  })
  len = len - termLen
  return 240 + 150 * termLen + 46 * len
}
const handleAddCondition = (index) => {
  formState.term[index].terms.push({
    is_mult: false,
    type: void 0,
    value: '',
    variable: [],
    key: Math.random() * 10000
  })
}

const handleDelCondition = (index, i) => {
  formState.term[index].terms.splice(i, 1)
}

const handleAddBranch = () => {
  formState.term.push({
    is_or: 0,
    key: Math.random() * 10000,
    next_node_key: '',
    terms: [
      {
        is_mult: false,
        type: void 0,
        value: '',
        variable: [],
        key: Math.random() * 10000
      }
    ]
  })
}

const handleDelBranch = (index) => {
  formState.term.splice(index, 1)
}

const handleVariableChange = (term) => {
  let typ = getTypeByVariable(term)
  term.is_mult = typ.includes('array')
  term.type = void 0
}

let baseTypeOptions = [
  {
    label: '等于',
    value: 1
  },
  {
    label: '不等于',
    value: 2
  },
  {
    label: '包含',
    value: 3
  },
  {
    label: '不包含',
    value: 4
  },
  {
    label: '为空',
    value: 5
  },
  {
    label: '不为空',
    value: 6
  }
]

let baseTypeOptions2 = [
  {
    label: '包含其中一项',
    value: 3
  },
  {
    label: '不包含其中一项',
    value: 4
  },
  {
    label: '为空',
    value: 5
  },
  {
    label: '不为空',
    value: 6
  }
]

function getTypeByVariable(data) {
  let typ = ''
  if (data.variable && data.variable.length > 0) {
    let slectItem = variableOptions.value.filter((item) => item.key == data.variable[0])
    if (slectItem && slectItem.length) {
      typ = slectItem[0].typ
    }
    if (typ == 'object') {
      if (slectItem[0] && slectItem[0].children) {
        slectItem = slectItem[0].subs.filter((item) => item.key == data.variable[1])
      }
      if (slectItem && slectItem.length) {
        typ = slectItem[0].typ
      }

      if (typ == 'object') {
        if (slectItem[0] && slectItem[0].children) {
          slectItem = slectItem[0].subs.filter((item) => item.key == data.variable[2])
        }
        if (slectItem && slectItem.length) {
          typ = slectItem[0].typ
        }
      }
    }
  }
  return typ
}

function getTypeOptions(data) {
  if (!data.variable) {
    return []
  }
  let typ = getTypeByVariable(data)

  if (typ == '') {
    if (data.is_mult) {
      return baseTypeOptions2
    } else {
      return baseTypeOptions
    }
  }

  if (typ == 'boole') {
    return [
      {
        label: '等于',
        value: 1
      },
      {
        label: '不等于',
        value: 2
      }
    ]
  }

  if (typ.includes('array')) {
    return baseTypeOptions2
  }

  return baseTypeOptions
}

onMounted(() => {
  getOptions()
})

defineExpose({})
</script>

<style lang="less" scoped>
@import '../form-block.less';

.main-text {
  color: #595959;
}
.btn-wrap {
  margin-top: 8px;
  padding-left: 76px;
  padding-right: 40px;
}
.condition-list-box {
  display: flex;
  align-items: center;
  .left-select-box {
    width: 90px;
    ::v-deep(.ant-select) {
      border-radius: 6px;
      transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
      &:hover {
        background: #e4e6eb;
      }
    }
    ::v-deep(.ant-select-selector) {
      color: #2475fc;
    }
    ::v-deep(.ant-select-arrow) {
      color: #2475fc;
    }
  }
  .condition-body {
    display: flex;
    flex-direction: column;
    gap: 8px;
    .condition-item {
      display: flex;
      gap: 8px;
    }
  }
}
.gray-block {
  margin-top: 8px;
}
.gray-block-title {
  display: flex;
  justify-content: space-between;
}
.add-btn-block {
  margin-top: 8px;
}
</style>
