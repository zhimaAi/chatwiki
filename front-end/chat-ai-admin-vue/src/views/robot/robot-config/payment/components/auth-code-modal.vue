<template>
  <a-modal
    :title="t('title_add_auth_code')"
    v-model:open="visible"
    :confirm-loading="saving"
    :cancel-text="finished ? t('btn_copy') : t('btn_cancel')"
    @ok="save"
    @cancel="close"
  >
    <template v-if="!finished">
      <a-alert type="info" class="zm-alert-info">
        <template #message>
          {{ t('msg_add_manager_tip') }}<a @click="emit('addManager')">{{ t('msg_add_manager_link') }}</a>
        </template>
      </a-alert>
      <a-form layout="vertical" class="mt16">
        <a-form-item :label="t('label_package')" required>
          <a-select v-model:value="formState.package_id" :placeholder="t('ph_select_package')">
            <a-select-option v-for="item in packages" :value="item.id">{{ item.name }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item :label="t('label_count')" required>
          <div class="count-box">
            <ZmRadioGroup v-model:value="formState.count" :options="options"/>
            <a-input-number style="width: 130px" v-model:value="formState.count" :min="1" :max="1000" :placeholder="t('ph_input')"/>
          </div>
        </a-form-item>
        <a-form-item :label="t('label_remark')">
          <a-input v-model:value.trim="formState.remark" :placeholder="t('ph_input')" style="width: 100%"/>
        </a-form-item>
      </a-form>
    </template>
    <div v-else class="success-box">
      <div class="tit">
        <CheckCircleFilled class="icon"/>
        {{ t('msg_add_success') }}
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
import { useI18n } from '@/hooks/web/useI18n';

const { t } = useI18n('views.robot.robot-config.payment.components.auth-code-modal');

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
    label: t('label_custom'),
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
    message.success(t('msg_copy_success'))
  }
  visible.value = false
}

function save() {
  if (finished.value) {
    emit('ok')
    visible.value = false
  } else {
    if (!formState.package_id) {
      return message.error(t('msg_select_package'))
    }
    if (!formState.count) {
      return message.error(t('msg_select_count'))
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
