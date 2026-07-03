<template>
  <a-modal
    v-model:open="open"
    :title="modalTitle"
    :ok-text="okText"
    :cancel-text="cancelText"
    :confirm-loading="submitting"
    @ok="handleOk"
    :destroyOnClose="true"
    :width="490"
  >
    <div class="form-box-wrapper" :class="currentSchema.submitType === 'goods' ? 'goods-editor-box' : 'field-editor-box'">
      <a-form
        :layout="isFieldMode ? 'vertical' : 'horizontal'"
        :label-col="isFieldMode ? undefined : { span: 4 }"
        :wrapper-col="isFieldMode ? undefined : { span: 20 }"
      >
        <a-form-item
          v-for="fieldName in currentSchema.fields"
          :key="fieldName"
          :label="isFieldMode ? '' : getFieldLabel(fieldName)"
          :required="isFieldRequired(fieldName)"
          :class="fieldName === 'cover_images' ? 'goods-image-item' : ''"
        >
          <template v-if="fieldDefs[fieldName].component === 'upload'">
            <a-upload
              v-model:fileList="fileList"
              list-type="picture-card"
              :multiple="true"
              :maxCount="5"
              accept="image/*"
              :beforeUpload="handleBeforeUpload"
              :custom-request="handleCustomUpload"
              @change="handleUploadChange"
              @preview="handlePreview"
            >
              <div v-if="fileList.length < 5" class="upload-trigger">
                <PlusOutlined class="upload-icon" />
                <div class="upload-text">{{ t('goods_modal.upload_action') }}</div>
              </div>
            </a-upload>
            <div class="upload-tip">{{ t('goods_modal.image_hint') }}</div>
          </template>

          <a-textarea
            v-else-if="fieldDefs[fieldName].component === 'textarea'"
            v-model:value="formModel[fieldName]"
            :rows="getTextareaRows(fieldName)"
            :placeholder="getFieldPlaceholder(fieldName)"
            :maxlength="getFieldMaxlength(fieldName)"
          />

          <a-tree-select
            v-else-if="fieldDefs[fieldName].component === 'select'"
            v-model:value="formModel[fieldName]"
            :placeholder="getFieldPlaceholder(fieldName)"
            :tree-data="groupTreeOptions"
            show-search
            tree-default-expand-all
            tree-node-filter-prop="title"
            style="width: 100%"
          />

          <a-input-number
            v-else-if="fieldDefs[fieldName].component === 'number'"
            v-model:value="formModel[fieldName]"
            :placeholder="getFieldPlaceholder(fieldName)"
            :min="fieldDefs[fieldName].min ?? 0"
            :precision="fieldDefs[fieldName].precision"
            style="width: 100%"
          />

          <a-input
            v-else
            v-model:value="formModel[fieldName]"
            :placeholder="getFieldPlaceholder(fieldName)"
            :maxlength="getFieldMaxlength(fieldName)"
            :show-count="fieldDefs[fieldName].showCount === true"
          />
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { generateRandomId } from '@/utils/index'
import { uploadGoodsImage } from '@/api/goods-library'
import { useI18n } from '@/hooks/web/useI18n'
import { api as viewerApi } from 'v-viewer'

defineProps({
  groupTreeOptions: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['submit-goods', 'submit-field'])

const { t } = useI18n('views.library.goods-library.index')

const FIELD_NAMES = {
  GOODS_ID: 'goods_id',
  GOODS_NAME: 'goods_name',
  GROUP_ID: 'group_id',
  CATEGORY: 'category',
  BRAND: 'brand',
  PRICE: 'price',
  STOCK: 'stock',
  IMAGES: 'images',
  DESCRIPTION: 'description',
  QA: 'qa',
  CUSTOM_INFO: 'custom_info',
  LINK: 'link'
}

const fieldDefs = {
  [FIELD_NAMES.GOODS_ID]: {
    labelKey: 'goods_modal.id_label',
    placeholderKey: 'goods_modal.id_placeholder',
    component: 'input',
    required: true,
    maxlength: 100,
    validate(value) {
      if (!String(value || '').trim()) {
        return t('validation.goods_id_required')
      }
      return ''
    }
  },
  [FIELD_NAMES.GOODS_NAME]: {
    labelKey: 'goods_modal.name_label',
    placeholderKey: 'goods_modal.name_placeholder',
    component: 'input',
    required: true,
    maxlength: 100,
    validate(value) {
      if (!String(value || '').trim()) {
        return t('validation.goods_name_required')
      }
      return ''
    }
  },
  [FIELD_NAMES.GROUP_ID]: {
    labelKey: 'goods_modal.group_label',
    placeholderKey: 'goods_modal.group_placeholder',
    component: 'select'
  },
  [FIELD_NAMES.CATEGORY]: {
    labelKey: 'goods_modal.category_label',
    placeholderKey: 'goods_modal.category_placeholder',
    component: 'input',
    maxlength: 100
  },
  [FIELD_NAMES.BRAND]: {
    labelKey: 'goods_modal.brand_label',
    placeholderKey: 'goods_modal.brand_placeholder',
    component: 'input',
    maxlength: 100
  },
  [FIELD_NAMES.PRICE]: {
    labelKey: 'goods_modal.price_label',
    placeholderKey: 'goods_modal.price_placeholder',
    component: 'number',
    min: 0,
    precision: 2,
    maxlength: 100
  },
  [FIELD_NAMES.STOCK]: {
    labelKey: 'goods_modal.stock_label',
    placeholderKey: 'goods_modal.stock_placeholder',
    component: 'number',
    min: 0,
    precision: 0,
    maxlength: 100
  },
  [FIELD_NAMES.IMAGES]: {
    labelKey: 'goods_modal.image_label',
    component: 'upload',
    getValue: () => getImageValues()
  },
  [FIELD_NAMES.DESCRIPTION]: {
    labelKey: 'table.description',
    placeholderKey: 'field_editor.text_placeholder',
    component: 'textarea',
    maxlength: 1000,
    rows: 7
  },
  [FIELD_NAMES.QA]: {
    labelKey: 'table.qa',
    placeholderKey: 'field_editor.text_placeholder',
    component: 'textarea',
    maxlength: 1000,
    rows: 7
  },
  [FIELD_NAMES.CUSTOM_INFO]: {
    labelKey: 'custom_info.title',
    placeholderKey: 'field_editor.text_placeholder',
    component: 'textarea',
    maxlength: 1000,
    rows: 7
  },
  [FIELD_NAMES.LINK]: {
    labelKey: 'goods_modal.link_label',
    placeholderKey: 'goods_modal.link_placeholder',
    component: 'input',
    maxlength: 1000
  }
}

const editorSchemas = {
  goods_full: {
    titleKey: 'goods_modal.title_edit',
    submitType: 'goods',
    fields: [
      FIELD_NAMES.GOODS_ID,
      FIELD_NAMES.GOODS_NAME,
      FIELD_NAMES.GROUP_ID,
      FIELD_NAMES.CATEGORY,
      FIELD_NAMES.BRAND,
      FIELD_NAMES.PRICE,
      FIELD_NAMES.STOCK,
      FIELD_NAMES.LINK,
      FIELD_NAMES.IMAGES
    ]
  },
  goods_basic: {
    titleKey: 'table.basic_info',
    submitType: 'goods',
    fields: [
      FIELD_NAMES.GOODS_ID,
      FIELD_NAMES.GOODS_NAME,
      FIELD_NAMES.GROUP_ID,
      FIELD_NAMES.CATEGORY,
      FIELD_NAMES.BRAND,
      FIELD_NAMES.PRICE,
      FIELD_NAMES.STOCK,
      FIELD_NAMES.LINK
    ]
  },
  goods_images: {
    titleKey: 'table.image',
    submitType: 'field',
    fieldKey: FIELD_NAMES.IMAGES,
    fields: [FIELD_NAMES.IMAGES]
  },
  goods_description: {
    titleKey: 'table.description',
    submitType: 'field',
    fieldKey: FIELD_NAMES.DESCRIPTION,
    fields: [FIELD_NAMES.DESCRIPTION]
  },
  goods_qa: {
    titleKey: 'table.qa',
    submitType: 'field',
    fieldKey: FIELD_NAMES.QA,
    fields: [FIELD_NAMES.QA]
  },
  goods_custom_info: {
    titleKey: 'custom_info.title',
    submitType: 'field',
    fieldKey: FIELD_NAMES.CUSTOM_INFO,
    fields: [FIELD_NAMES.CUSTOM_INFO]
  }
}

const open = ref(false)
const submitting = ref(false)
const currentSchemaName = ref('goods_full')
const currentRow = ref(null)
const fileList = ref([])

const formModel = reactive({
  id: '',
  goods_id: '',
  goods_name: '',
  group_id: undefined,
  category: '',
  brand: '',
  price: undefined,
  stock: undefined,
  description: '',
  qa: '',
  custom_info: '',
  link: ''
})

const currentSchema = computed(() => {
  return editorSchemas[currentSchemaName.value]
})

const isFieldMode = computed(() => {
  return currentSchema.value.submitType === 'field'
})

const modalTitle = computed(() => {
  if (currentSchemaName.value === 'goods_full' && !formModel.id) {
    return t('goods_modal.title_add')
  }

  if (isFieldMode.value) {
    return t(currentSchema.value.titleKey)
  }

  return t('goods_modal.title_edit_content')
})

const okText = computed(() => {
  return currentSchema.value.submitType === 'goods'
    ? (formModel.id ? t('goods_modal.save_edit') : t('goods_modal.save_add'))
    : t('field_editor.save')
})

const cancelText = computed(() => {
  return currentSchema.value.submitType === 'goods'
    ? t('confirm.cancel_btn')
    : t('field_editor.cancel')
})

const normalizeImages = (value) => {
  if (!value) {
    return []
  }

  const list = Array.isArray(value) ? value : [value]
  return list
    .filter(Boolean)
    .map((item, index) => {
      if (typeof item === 'string') {
        return {
          uid: `${generateRandomId(6)}-${index}`,
          name: item.split('/').pop() || `image-${index + 1}`,
          url: item,
          status: 'done'
        }
      }

      const url = item.url || item.thumbUrl || item.path || item.src || ''
      return {
        uid: item.uid || `${generateRandomId(6)}-${index}`,
        name: item.name || url.split('/').pop() || `image-${index + 1}`,
        url,
        thumbUrl: item.thumbUrl || url,
        status: 'done'
      }
    })
    .filter((item) => item.url || item.thumbUrl)
}

const getImageValues = () => {
  return fileList.value
    .map((item) => item.url || item.thumbUrl || '')
    .filter(Boolean)
}

const getFieldLabel = (fieldName) => {
  return t(fieldDefs[fieldName].labelKey)
}

const getFieldPlaceholder = (fieldName) => {
  const key = fieldDefs[fieldName].placeholderKey
  return key ? t(key) : ''
}

const getFieldMaxlength = (fieldName) => {
  return fieldDefs[fieldName].maxlength
}

const getTextareaRows = (fieldName) => {
  return fieldDefs[fieldName].rows || 4
}

const isFieldRequired = (fieldName) => {
  return fieldDefs[fieldName].required === true
}

const resetFormModel = () => {
  formModel.id = ''
  formModel.goods_id = ''
  formModel.goods_name = ''
  formModel.group_id = undefined
  formModel.category = ''
  formModel.brand = ''
  formModel.price = undefined
  formModel.stock = undefined
  formModel.description = ''
  formModel.qa = ''
  formModel.custom_info = ''
  formModel.link = ''
}

const fillGoodsModel = (data = {}) => {
  formModel.id = data.id || ''
  formModel.goods_id = data.goods_id || data.id || ''
  formModel.goods_name = data.goods_name || data.name || ''
  formModel.group_id = data.group_id ? String(data.group_id) : undefined
  formModel.category = data.category || ''
  formModel.brand = data.brand || ''
  formModel.price = data.price !== undefined && data.price !== null && data.price !== '' ? Number(data.price) : undefined
  formModel.stock = data.stock !== undefined && data.stock !== null && data.stock !== '' ? Number(data.stock) : undefined
  formModel.link = data.link || ''
  fileList.value = normalizeImages(data.images)
}

const fillFieldModel = (config = {}) => {
  const fieldKey = editorSchemas[currentSchemaName.value].fieldKey
  currentRow.value = config.row || null
  fileList.value = []

  if (fieldKey === FIELD_NAMES.IMAGES) {
    fileList.value = normalizeImages(config.value)
    return
  }

  formModel[fieldKey] = Array.isArray(config.value) ? config.value.join('\n') : (config.value ?? '')
}

const show = (config = {}) => {
  resetFormModel()
  currentRow.value = null
  fileList.value = []
  submitting.value = false
  currentSchemaName.value = config.name || 'goods_full'

  if (currentSchema.value.submitType === 'goods') {
    fillGoodsModel(config.data || {})
  } else {
    fillFieldModel(config)
  }

  open.value = true
}

const close = () => {
  submitting.value = false
  open.value = false
}

const handlePreview = (file) => {
  const images = fileList.value
    .map((item) => item.url || item.thumbUrl)
    .filter(Boolean)

  if (!images.length) return

  const currentIndex = images.indexOf(file.url || file.thumbUrl)

  viewerApi({
    images,
    options: {
      initialViewIndex: currentIndex >= 0 ? currentIndex : 0,
      toolbar: true,
      title: false,
      movable: true,
      zoomable: true,
      rotatable: true,
      scalable: true
    }
  })
}

const handleBeforeUpload = (file) => {
  if (file.size > 10 * 1024 * 1024) {
    message.error(t('validation.image_size_limit'))
    return false
  }

  return true
}

const handleCustomUpload = async ({ file, onSuccess, onError }) => {
  try {
    const res = await uploadGoodsImage({ file })
    const link = res?.data?.link || ''
    onSuccess({ link }, file)
  } catch (e) {
    onError(e)
  }
}

const handleUploadChange = (info) => {
  const nextList = info.fileList.slice(0, 5)

  nextList.forEach((item) => {
    if (item.status === 'done' && item.response?.link) {
      item.url = item.response.link
      item.thumbUrl = item.response.link
    }
  })

  fileList.value = nextList
}

const validateBySchema = () => {
  for (const fieldName of currentSchema.value.fields) {
    const validator = fieldDefs[fieldName].validate
    if (!validator) {
      continue
    }

    const value = fieldDefs[fieldName].getValue
      ? fieldDefs[fieldName].getValue()
      : formModel[fieldName]
    const errorMessage = validator(value)
    if (errorMessage) {
      message.error(errorMessage)
      return false
    }
  }

  return true
}

const buildGoodsPayload = () => {
  const images = getImageValues()
  return {
    id: formModel.id,
    goods_id: formModel.goods_id,
    goods_name: formModel.goods_name,
    group_id: formModel.group_id,
    category: formModel.category,
    brand: formModel.brand,
    price: formModel.price,
    stock: formModel.stock,
    link: formModel.link,
    images
  }
}

const buildFieldPayload = () => {
  const fieldKey = currentSchema.value.fieldKey
  const value = fieldDefs[fieldKey].getValue ? fieldDefs[fieldKey].getValue() : formModel[fieldKey]
  return {
    fieldKey,
    fieldLabel: t(fieldDefs[fieldKey].labelKey),
    value,
    row: currentRow.value
  }
}

const handleOk = () => {
  if (!validateBySchema()) {
    return
  }

  if (submitting.value) {
    return
  }

  submitting.value = true

  const submitActions = {
    close,
    setSubmitting: (value) => {
      submitting.value = value
    }
  }

  if (currentSchema.value.submitType === 'goods') {
    emit('submit-goods', buildGoodsPayload(), submitActions)
  } else {
    emit('submit-field', buildFieldPayload(), submitActions)
  }
}

defineExpose({
  show,
  close
})
</script>

<style lang="less" scoped>
.form-box-wrapper{
  padding-top: 8px;

  :deep(.ant-form-item) {
    margin-bottom: 16px;
  }
}

.upload-trigger {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
}

.upload-icon {
  font-size: 24px;
  color: rgba(0, 0, 0, 0.45);
}

.upload-text {
  margin-top: 14px;
  color: #595959;
  font-size: 14px;
  line-height: 22px;
}

.upload-tip {
  margin-top: 2px;
  color: #8c8c8c;
  font-size: 14px;
  line-height: 22px;
}

.goods-image-item {
  width: 336px;
}
</style>
