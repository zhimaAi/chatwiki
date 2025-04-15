<style lang="less" scoped>
.form-box {
  overflow: hidden;
  border-radius: 6px;
  background: #f2f4f7;
}

.model-icon {
  height: 18px;
}

.form-item-tip {
  color: #999;
}

.card-box {
  display: flex;
  gap: 8px;
}

.use-model-item {
  position: relative;
  width: 226px;
  height: 124px;
  border-radius: 2px;
  border: 2px solid #d9d9d9;
  cursor: pointer;
  padding: 15px;
  margin-bottom: 10px;
}

.use-model-item-top {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  font-size: 14px;
  color: #595959;
}

.use-model-item-top-icon {
  margin-right: 5px;
}

.use-model-item-top.active {
  color: #2475fc;
  font-weight: bolder;
}

.use-model-item.active {
  border: 2px solid #2475fc;

  .check-arrow {
    position: absolute;
    display: block;
    right: -1px;
    bottom: -1px;
    width: 24px;
    height: 24px;
    font-size: 24px;
    color: #fff;
  }

  .retrieval-mode-title {
    color: #2475fc;
  }
}

.form-box-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  height: 56px;

  .form-box-label {
    display: flex;
    align-items: center;

    .form-box-label-icon {
      width: 16px;
      height: 16px;
      font-size: 16px;
      margin-right: 8px;
    }
    .form-box-label-text {
      font-size: 16px;
      font-weight: 600;
      color: #262626;
    }
  }
}
</style>

<template>
  <div class="form-box">
    <div class="form-box-header">
      <div class="form-box-label">
        <svg-icon
          class="form-box-label-icon"
          name="book"
          style="font-size: 16px; color: #333"
        ></svg-icon>
        <span class="form-box-label-text">知识库信息</span>
      </div>
      <div></div>
    </div>
    <a-form :label-col="{ span: 6 }" style="width: 600px">
      <a-form-item ref="name" label="知识库名称" v-bind="validateInfos.library_name">
        <a-input
          v-model:value="formState.library_name"
          placeholder="请输入知识库名称，最多20个字"
          show-count
          :maxlength="20"
        />
      </a-form-item>

      <a-form-item label="知识库简介">
        <a-textarea v-model:value="formState.library_intro" placeholder="请输入知识库介绍" />
      </a-form-item>

      <a-form-item ref="name" label="知识库封面" v-bind="validateInfos.avatar">
        <AvatarInput v-model:value="formState.avatar" @change="onAvatarChange" />
        <div class="form-item-tip">请上传知识库封面，建议尺寸为100*100px.大小不超过100kb</div>
      </a-form-item>
      <a-form-item :wrapper-col="{ offset: 6, span: 8 }">
        <a-button
          type="primary"
          style="margin-left: 8px"
          :loading="saveLoading"
          @click.prevent="onSubmit"
          >保存</a-button
        >
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted, toRaw, watch } from 'vue'
import { Form, message } from 'ant-design-vue'
import AvatarInput from '@/views/library/add-library/components/avatar-input.vue'
import { LIBRARY_OPEN_AVATAR } from '@/constants/index'

// 设置全局默认的duration为（2秒）
message.config({
  duration: 2
})

const emit = defineEmits(['submit', 'update:value'])

const props = defineProps({
  value: {
    type: Object,
    default: () => {
      return {}
    }
  },
  saveLoading: {
    type: Boolean,
    default: false
  }
})

const type = ref(1)
const useForm = Form.useForm

const isActive = ref(false)
const defaultAvatar = LIBRARY_OPEN_AVATAR

const formState = reactive({
  type: type.value,
  access_rights: 0,
  library_name: '',
  library_intro: '',
  use_model: '',
  avatar: defaultAvatar,
  avatar_file: '',
  is_offline: false,
  model_config_id: ''
})

const rules = reactive({
  library_name: [{ required: true, message: '请输入库名称', trigger: 'change' }],
  library_intro: [{ required: true, message: '请输入库简介', trigger: 'change' }]
  // use_model: [{ required: true, message: '请选择嵌入模型', trigger: 'change' }]
})

const onAvatarChange = (data) => {
  formState.avatar = data.imageUrl
  formState.avatar_file = data.file
}

const { validate, validateInfos } = useForm(formState, rules)
// const resetForm = () => {}

const onSubmit = () => {
  validate()
    .then(() => {
      saveForm()
    })
    .catch((err) => {
      console.log(err, 'err')
    })
}

const saveForm = () => {
  let data = toRaw(formState)

  emit('submit', data)
}

const setFormState = (val) => {
  isActive.value = +val.type
  formState.library_name = val.library_name
  formState.library_intro = val.library_intro
  formState.use_model = val.use_model
  formState.model_config_id = val.model_config_id
  formState.avatar = val.avatar ? val.avatar : defaultAvatar
  formState.avatar_file = val.avatar_file ? val.avatar_file : ''
  formState.is_offline = val.is_offline
}

watch(props.value, (val) => {
  setFormState({ ...val })
})

onMounted(() => {})
</script>
