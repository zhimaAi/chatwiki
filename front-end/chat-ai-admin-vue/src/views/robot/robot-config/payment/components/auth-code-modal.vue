<template>
  <a-modal
    title="新增授权码"
    v-model:open="visible"
    :confirm-loading="saving"
    :cancel-text="finished ? '复制' : '取消' "
    @ok="save"
    @cancel="close"
  >
    <template v-if="!finished">
      <a-alert type="info" class="zm-alert-info">
        <template #message>设置管理员后，管理员可在公众号内回复【授权码】快速新垲，<a @click="emit('addManager')">点此新增管理员</a></template>
      </a-alert>
      <a-form layout="vertical" class="mt16">
        <a-form-item label="套餐" required>
          <a-select v-model:value="formState.package_id" placeholder="选择套餐">
            <a-select-option v-for="item in packages" :value="item.id">{{ item.name }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="数量" required>
          <div class="count-box">
            <ZmRadioGroup v-model:value="formState.count" :options="options"/>
            <a-input-number v-model:value="formState.count" :min="1" :max="1000" placeholder="请输入"/>
          </div>
        </a-form-item>
        <a-form-item label="备注">
          <a-input v-model:value.trim="formState.remark" placeholder="请输入" style="width: 100%"/>
        </a-form-item>
      </a-form>
    </template>
    <div v-else class="success-box">
      <div class="tit">
        <CheckCircleFilled class="icon"/>
        新增成功
      </div>
      <div class="code-list">
        <div class="code-item" v-for="code in list" :key="code.content">{{code.content}}</div>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import {ref, reactive, computed} from 'vue';
import ZmRadioGroup from "@/components/common/zm-radio-group.vue";
import {CheckCircleFilled} from '@ant-design/icons-vue';
import {message} from 'ant-design-vue';
import {addAuthCode} from "@/api/robot/payment.js";
import {copyText} from "@/utils/index.js";

const emit = defineEmits(['addManager', 'ok'])
const props = defineProps({
  robotId: {
    type: [Number, String]
  },
  packages: {
    type: Array,
  }
})
const visible = ref(false)
const saving = ref(false)
const finished = ref(false)
const list = ref([])
const formState = reactive({
  package_id: undefined,
  count: 1,
  remark: ''
})

const options = computed(() => {
  const base = [
    {label: '1', value: 1},
    {label: '5', value: 5},
    {label: '10', value: 10},
    {label: '50', value: 50},
  ]
  base.push({
    label: '自定义',
    value: [1, 5, 10, 50].includes(formState.count) ? '' : formState.count
  })
  return base
})

function show() {
  list.value = []
  finished.value = false
  visible.value = true
}

function close() {
  if (finished.value) {
    emit('ok')
    let codes = list.value.map(i => i.content)
    copyText(codes.join('\n'))
    message.success('复制成功')
  }
  visible.value = false
}

function save() {
  if (finished.value) {
    emit('ok')
    visible.value = false
  } else {
    if (!formState.package_id) {
      return message.error('请选择套餐')
    }
    if (!formState.count) {
      return message.error('请选择数量')
    }
    saving.value = true
    addAuthCode({
      ...formState,
      robot_id: props.robotId
    }).then(res => {
      list.value = res?.data || []
      finished.value = true
    }).finally(() => {
      saving.value = false
    })
  }
}

defineExpose({
  show
})
</script>

<style scoped lang="less">
.mt16 {
  margin-top: 16px;
}

.count-box {
  display: flex;
  align-items: center;
  gap: 8px;
}

.success-box {
  .tit {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #262626;
    font-size: 20px;
    font-weight: 600;
    margin-bottom: 8px;

    .icon {
      color: #21A665;
    }
  }

  .code-list {
    max-height: 256px;
    border-radius: 6px;
    background: #EBF7FF;
    overflow-y: auto;
    padding: 8px 16px;
  }
}
</style>
