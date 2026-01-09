<template>
  <a-modal
    :title="info ? '编辑套餐' : '新增套餐'"
    v-model:open="visible"
    @ok="save"
  >
    <a-form layout="vertical">
      <a-form-item label="套餐名" required>
        <a-input v-model:value.trim="formState.name" placeholder="示例：单日套餐，最多6个字" :maxlength="6"/>
      </a-form-item>
      <a-form-item v-if="type == 2" label="时长（天）" required>
        <a-input-number v-model:value="formState.duration" :min="0" :precision="0" placeholder="请输入" style="width: 100%"/>
      </a-form-item>
      <a-form-item label="次数（次）" required>
        <a-input-number v-model:value="formState.count" :min="0" :precision="0" placeholder="请输入" style="width: 100%"/>
      </a-form-item>
      <a-form-item label="费用（元）" required>
        <a-input-number v-model:value="formState.price" :min="0" :precision="2" placeholder="请输入" style="width: 100%"/>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import {ref, reactive, toRaw} from 'vue'
import {message} from 'ant-design-vue'

const emit = defineEmits(['ok'])
const props = defineProps({
  type: {
    type: [Number, String],
    default: 1
  }
})
const info = ref(null)
const visible = ref(false)
const formState = reactive({})

function show(_info) {
  info.value = _info || null
  Object.assign(formState, _info || {
    id: 0,
    name: "",
    duration: "",
    count: "",
    price: ""
  })
  visible.value = true
}

function save() {
  try {
    if (!formState.name) throw '请输入套餐名称'
    if (props.type == 2 && !formState.duration) throw '请输入套餐时长'
    if (!formState.count) throw '请输入套餐次数'
    if (!formState.price) throw '请输入套餐费用'
    let res = toRaw(formState)
    if (props.type == 1) delete res.duration
    visible.value = false
    emit('ok', JSON.parse(JSON.stringify(res)))
  } catch (e) {
    message.error(e)
  }
}

defineExpose({
  show,
})
</script>

<style scoped>

</style>
