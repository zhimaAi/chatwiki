<style lang="less" scoped>
.seo-info-wrapper {
  padding: 11px 16px;
  border-radius: 6px;
  border: 1px solid #f2f4f7;
  background-color: #f2f4f7;

  .seo-info-header {
    display: flex;
    align-items: center;

    .header-body {
      flex: 1;
      padding-left: 4px;

      .seo-title {
        color: #595959;
        font-size: 14px;
        font-weight: 600;
      }

      .seo-edit-title{
        font-size: 14px;
        color: #8c8c8c;
        &::before {
          content: ' ';
          display: inline-block;
          width: 1px;
          height: 12px;
          background-color: #D8DDE6;
          margin: 0 8px;
        }
      }
    }

    .seo-icon {
      width: 33px;
      height: 16px;
      margin-right: 4px;
    }

    .seo-label {
      display: inline-block;
      line-height: 22px;
      font-size: 14px;
      font-weight: 600;
      color: #262626;
    }
    .actions-box {
      .toggle-btn {
        color: #2475fc;

        .toggle-icon {
          margin-left: 2px;
          font-size: 12px;
          color: #2475fc;
        }
      }
    }
  }

  .seo-info-preview {
    margin-top: 8px;

    .info-item{
      display: flex;
      line-height: 22px;
      margin-bottom: 4px;
      color: #3a4559;

      &:last-child {
        margin-bottom: 0;
      }

      .info-item-label{
        width: 70px;
        text-align: right;
      }

      .info-item-content{
        flex: 1;
        white-space: pre-wrap;
        word-break: break-all;
      }
    }
  }

  .seo-info-form{
    .form-row{
      display: flex;
      margin-top: 16px;
      gap: 16px;
    }

    .form-col{
      flex: 1;
    }

    .form-item-label{
      line-height: 22px;
      margin-bottom: 4px;
      font-size: 14px;
      color: #262626;

    }
  }

  &.seo-is-set {
    background: #f0f4ff;
    border: 1px solid #2475fc;
    .header-body {
      .seo-title {
        color: #14161a;
      }
    }
  }

  &.seo-is-expand {
  }
}
</style>

<template>
  <div class="seo-info-wrapper" :class="{ 'seo-is-expand': isOpen, 'seo-is-set': isSet }">
    <div class="seo-info-header">
      <img src="@/assets/svg/seo-gray.svg" class="seo-icon" v-if="!isSet" />
      <img src="@/assets/svg/seo-light.svg" class="seo-icon" v-else />
      <span class="seo-label" v-if="isEdit">SEO设置</span>
      <div class="header-body">
        <template v-if="isEdit">
          <span class="seo-edit-title" v-if="isEdit">SEO设置也要发布后才会生效</span>
        </template>
        <template v-else>
          <span class="seo-title">{{ isSet ? formData.seo_title : '还未添加SEO设置' }}</span>
        </template>
      </div>

      <div class="actions-box">
        <a-button type="primary" size="small" :loading="confirmLoading" @click="handleSave" v-if="isEdit">保存</a-button>
        <template v-else>
          <a-button class="toggle-btn" type="text" size="small" @click="handleEdit" v-if="!isSet">
            <a>去设置 <RightOutlined class="toggle-icon" /></a>
          </a-button>
          <a-button class="toggle-btn" type="text" size="small" @click="handleEdit" v-else>
            <a>修改</a>
          </a-button>

          <a-button class="toggle-btn" type="text" size="small" @click="handleToggle" v-if="isSet">
            <a>{{ isOpen ? '收起' : '展开' }}</a>
            <RightOutlined class="toggle-icon" v-if="!isOpen" />
            <UpOutlined class="toggle-icon" v-if="isOpen" />
          </a-button>
        </template>
      </div>
    </div>

    <div class="seo-info-preview" v-if="isOpen && !isEdit">
      <div class="info-item">
        <div class="info-item-label">描述：</div>
        <div class="info-item-content">
          {{ formData.seo_desc }}
        </div>
      </div>
      <div class="info-item">
        <div class="info-item-label">关键词：</div>
        <div class="info-item-content">
          {{ formData.seo_keywords }}
        </div>
      </div>
    </div>

    <div class="seo-info-form" v-if="isEdit && !isOpen">
      <div class="form-row">
        <div class="form-col">
          <div class="form-item">
            <div class="form-item-label">Title</div>
            <div class="form-item-content">
              <a-input v-model:value="formData.seo_title" placeholder="请输入" />
            </div>
          </div>
        </div>
      </div>
      <div class="form-row">
        <div class="form-col">
          <div class="form-item">
            <div class="form-item-label">Description</div>
            <div class="form-item-content">
              <a-textarea v-model:value="formData.seo_desc" placeholder="请输入" />
            </div>
          </div>
        </div>

        <div class="form-col">
          <div class="form-item">
            <div class="form-item-label">Keyword</div>
            <div class="form-item-content">
              <a-textarea v-model:value="formData.seo_keywords" placeholder="请输入" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, toRaw, reactive, watch } from 'vue'
import { saveLibDocSeo } from '@/api/public-library/index'
import { useStorage } from '@/hooks/web/useStorage'
import { message } from 'ant-design-vue'
import { RightOutlined, UpOutlined } from '@ant-design/icons-vue'

const { setStorage, getStorage } = useStorage('localStorage')
const emit = defineEmits(['save'])

const props = defineProps({
  form: {
    type: Object,
    default: () => {}
  }
})

const show = ref(false)
const ssoInofExpand = getStorage('sso_inof_expand') || false
const confirmLoading = ref(false)
const isSet = ref(false) // 是否设置
const isEdit = ref(false) // 是否编辑
const isOpen = ref(ssoInofExpand) // 是否展开

// 表单数据
const formData = reactive({
  library_key: '',
  doc_id: '',
  seo_title: '',
  seo_desc: '',
  seo_keywords: ''
})

watch(() => props.form, (val) => {
 formData.library_key = val.library_key || ''
 formData.doc_id = val.doc_id || ''
 formData.seo_title = val.seo_title || ''
 formData.seo_desc = val.seo_desc || ''
 formData.seo_keywords = val.seo_keywords || ''

  if(val.seo_title || val.seo_keywords || val.seo_desc){
    isSet.value = true
  }else{
    isSet.value = false
  }
}, {
  immediate: true,
  deep: true
})

const handleToggle = () => {
  isOpen.value = !isOpen.value
  isEdit.value = false

  setStorage('sso_inof_expand', isOpen.value)
}

const handleEdit = () => {
  isOpen.value = false
  isEdit.value = true
}

const handleSave = () => {
  if(confirmLoading.value) return

  confirmLoading.value = true

  saveLibDocSeo({
    ...toRaw(formData)
  })
    .then(() => {
      confirmLoading.value = false
      show.value = false
      message.success('保存成功')
      isEdit.value = false

      emit('save', {...toRaw(formData)})
    })
    .catch(() => {
      confirmLoading.value = false
    })
}

defineExpose({
  open
})
</script>
