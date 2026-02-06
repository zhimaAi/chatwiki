<template>
  <div>
    <a-modal
      v-model:open="open"
      :title="modalTitle"
      @ok="handleOk"
      :width="472"
      :ok-text="t('btn_confirm')"
      :cancel-text="t('btn_cancel')"
    >
      <a-form ref="formRef" layout="vertical" :model="formState" style="margin: 24px 0">
        <a-form-item
          :label="t('label_param_key')"
          name="key"
          :rules="[
            {
              required: true,
              validator: (rule, value) => checkKey(rule, value)
            }
          ]"
        >
          <a-input
            v-model:value="formState.key"
            :maxlength="20"
            placeholder="英文字母和下划线“_”组成，最多20字符"
          />
        </a-form-item>
        <a-form-item
          :label="t('label_param_type')"
          name="typ"
          :rules="[{ required: true, message: t('msg_select_type') }]"
        >
          <a-select v-model:value="formState.typ" :placeholder="t('ph_select')" style="width: 100%">
            <a-select-option v-for="op in typOptions" :value="op.value">{{
              op.value
            }}</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'
import { message } from 'ant-design-vue'
import { getUuid } from '@/utils/index'
import { ref, reactive, computed } from 'vue'

const { t } = useI18n('views.workflow.components.node-form-drawer.webhook-trigger-node-form.add-key-modal')

const emit = defineEmits(['add', 'edit', 'addSub'])

const props = defineProps({
  request_content_type: {
    type: String,
    default: ''
  },
  level: {
    type: Number,
    default: 1
  }
})

const typOptions = computed(() => {
  if (props.request_content_type == 'multipart/form-data') {
    return [
      {
        lable: 'string',
        value: 'string'
      },
      {
        lable: 'file',
        value: 'file'
      }
    ]
  }
  if (props.request_content_type == 'application/x-www-form-urlencoded') {
    return [
      {
        lable: 'string',
        value: 'string'
      }
    ]
  }

  if (props.request_content_type == 'application/json') {
    let baseList = [
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
        lable: 'object',
        value: 'object'
      },
      {
        lable: 'array<string>',
        value: 'array<string>'
      },
      {
        lable: 'array<number>',
        value: 'array<number>'
      },
      {
        lable: 'array\<object>',
        value: 'array\<object>'
      }
    ]
    if (props.level >= 2) {
      return baseList.filter((i) => i.value != 'object')
    } else {
      return baseList
    }
  }

  return [
    {
      lable: 'string',
      value: 'string'
    }
  ]
})

const open = ref(false)
const modalTitle = ref(t('title_add_param'))

let oprateType = 'add'

const checkKey = (rule, value) => {
  if (!value) {
    return Promise.reject(t('msg_enter_key'))
  }
  // 校验是否只包含英文字母和下划线
  const regex = /^[a-zA-Z_]+$/
  if (!regex.test(value)) {
    return Promise.reject(t('msg_only_letters_underscore'))
  }
  return Promise.resolve()
}
let INDEX = 0
const formState = reactive({
  key: '',
  typ: 'string',
  cu_key: '',
  desc: void 0,
  subs: []
})
const show = () => {
  formState.typ = 'string'
  open.value = true
}

const add = () => {
  Object.assign(formState, {
    key: '',
    typ: 'string',
    cu_key: getUuid(16),
    subs: []
  })
  oprateType = 'add'
  modalTitle.value = t('title_add_param')
  open.value = true
}

const addSub = (index) => {
  INDEX = index
  Object.assign(formState, {
    key: '',
    typ: 'string',
    cu_key: getUuid(16),
    subs: []
  })
  oprateType = 'addSub'
  modalTitle.value = t('title_add_param')
  open.value = true
}

const edit = (data, index) => {
  INDEX = index
  Object.assign(formState, {
    key: '',
    typ: 'string',
    cu_key: getUuid(16),
    subs: []
  })
  Object.assign(formState, {
    ...data
  })
  oprateType = 'edit'
  modalTitle.value = t('title_edit_param')
  open.value = true
}

const formRef = ref(null)
const handleOk = () => {
  formRef.value.validate().then(() => {
    if (oprateType == 'add') {
      emit('add', {
        ...formState
      })
    }
    if (oprateType == 'addSub') {
      emit(
        'addSub',
        {
          ...formState
        },
        INDEX
      )
    }
    if (oprateType == 'edit') {
      emit(
        'edit',
        {
          ...formState
        },
        INDEX
      )
    }
    open.value = false
  })
}

defineExpose({
  show,
  add,
  edit,
  addSub
})
</script>

<style lang="less" scoped></style>
