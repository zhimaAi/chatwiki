<template>
  <a-modal
    :title="info ? t('title_edit_package') : t('title_create_package')"
    v-model:open="visible"
    @ok="save"
  >
    <a-form layout="vertical">
      <a-form-item :label="t('label_package_name')" required>
        <a-input v-model:value.trim="formState.name" :placeholder="t('ph_package_name_example')" :maxlength="6"/>
      </a-form-item>
      <a-form-item v-if="type == 2" :label="t('label_duration')" required>
        <a-input-number v-model:value="formState.duration" :min="0" :precision="0" :placeholder="t('ph_input')" style="width: 100%"/>
      </a-form-item>
      <a-form-item :label="t('label_count')" required>
        <a-input-number v-model:value="formState.count" :min="0" :precision="0" :placeholder="t('ph_input')" style="width: 100%"/>
      </a-form-item>
      <a-form-item :label="t('label_price')" required>
        <a-input-number v-model:value="formState.price" :min="0" :precision="2" :placeholder="t('ph_input')" style="width: 100%"/>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import {ref, reactive, toRaw} from 'vue'
import {message} from 'ant-design-vue'
import {useI18n} from '@/hooks/web/useI18n'

const {t} = useI18n('views.robot.robot-config.payment.components.comb-store-modal')

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
    if (!formState.name) throw t('msg_input_package_name')
    if (props.type == 2 && !formState.duration) throw t('msg_input_package_duration')
    if (!formState.count) throw t('msg_input_package_count')
    if (!formState.price) throw t('msg_input_package_price')
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
