<template>
  <div>
    <a-modal
      v-model:open="show"
      :title="modalTitle"
      @ok="handleOk"
      width="472px"
      :confirmLoading="saveLoading"
    >
      <div class="form-box">
        <a-form layout="vertical">
          <a-form-item label="文档名称" v-bind="validateInfos.file_name">
            <a-input
              type="text"
              placeholder="请输入文档名称"
              v-model:value="formState.file_name"
            ></a-input>
          </a-form-item>
          <a-form-item label="文档类型" required v-if="false">
            <a-radio-group v-model:value="formState.is_qa_doc">
              <a-radio :value="0">普通文档</a-radio>
              <a-radio :value="1">QA文档</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item label="索引方式" required v-if="formState.is_qa_doc == 1">
            <div class="upload-document-type-box">
              <div
                class="type-item"
                :class="{ active: formState.qa_index_type == item.value }"
                v-for="item in qaIndexTypeList"
                :key="item.value"
                @click="handleChangeQaIndexType(item.value)"
              >
                <div class="right-block">
                  <div class="title-block">
                    <svg-icon
                      :name="
                        formState.qa_index_type == item.value ? item.iconNameActive : item.iconName
                      "
                    ></svg-icon>
                    <div class="title-text">{{ item.title }}</div>
                  </div>
                  <div class="desc">{{ item.desc }}</div>
                </div>
                <svg-icon
                  class="check-arrow"
                  name="check-arrow-filled"
                  style="font-size: 24px; color: #fff"
                  v-if="formState.qa_index_type == item.value"
                ></svg-icon>
              </div>
            </div>
          </a-form-item>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { validatePassword } from '@/utils/validate.js'
import { ref, reactive, toRaw } from 'vue'
import { Form, message } from 'ant-design-vue'
import { addLibraryFile } from '@/api/library'
import { useRouter, useRoute } from 'vue-router'
const emit = defineEmits(['ok'])
const router = useRouter()
const rotue = useRoute()
const useForm = Form.useForm
const show = ref(false)
const modalTitle = ref('添加自定义文档')

const props = defineProps({
  libraryInfo:{
    type: Object,
    default: () => {}
  },
  group_id: {
    type: [Number, String],
    default: () => 0
  },
})

const formState = reactive({
  library_id: rotue.query.id,
  doc_type: 3,
  file_name: '',
  is_qa_doc: 0,
  qa_index_type: 1
})
const qaIndexTypeList = ref([
  {
    iconName: 'file-search',
    iconNameActive: 'file-search',
    title: '问题与答案一起生成索引',
    value: 1,
    desc: '回答用户提问时，将用户提问与导入的问题和答案一起对比相似度，根据相似度高的问题和答案回复'
  },
  {
    iconName: 'comment-search',
    iconNameActive: 'comment-search',
    title: '仅对问题生成索引',
    value: 2,
    desc: '回答用户提问时，将用户提问与导入的问题一起对比相似度，再根据相似度高的问题和对应的答案来回复'
  }
])
const handleChangeQaIndexType = (type) => {
  formState.qa_index_type = type
}
const formRules = reactive({
  file_name: [
    {
      message: '请输入昵称',
      required: true
    }
  ]
})

const { resetFields, validate, validateInfos } = useForm(formState, formRules)

const add = () => {
  show.value = true
  resetFields()
  formState.file_name = ''
  formState.is_qa_doc = 0
  formState.qa_index_type = 1
  formState.library_id = props.libraryInfo.id || rotue.query.id
  if(props.libraryInfo.type == 2){
    formState.is_qa_doc = 1
  }
}

const saveLoading = ref(false)
const handleOk = () => {
  validate().then(() => {
    let formData = {
      ...toRaw(formState),
      group_id: props.group_id
    }
    saveLoading.value = true
    addLibraryFile(formData)
      .then((res) => {
        message.success(`${modalTitle.value}成功`)
        show.value = false
        router.push('/library/preview?id=' + res.data.file_ids[0])
        // emit('ok')
      })
      .finally(() => {
        saveLoading.value = false
      })
  })
}

defineExpose({
  add
})
</script>

<style lang="less" scoped>
.form-box {
  margin-top: 24px;
}

.upload-document-type-box {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  .type-item {
    position: relative;
    width: 100%;
    cursor: pointer;
    padding: 16px;
    display: flex;
    border: 1px solid #d9d9d9;
    border-radius: 2px;
    box-shadow: none;
    transition: box-shadow 1s;
    &:hover {
      box-shadow: 0px 5px 5px -3px rgba(0, 0, 0, 0.1), 0px 8px 10px 1px rgba(0, 0, 0, 0.06),
        0px 3px 14px 2px rgba(0, 0, 0, 0.05);
    }
    &.active {
      border: 2px solid #2475fc;
      .svg-action {
        color: #2475fc;
      }
      .right-block .title-text {
        color: #2475fc;
      }
    }
    .check-arrow {
      position: absolute;
      bottom: 0;
      right: -1px;
    }
  }
  .right-block {
    .title-block {
      display: flex;
      font-size: 14px;
      line-height: 22px;
      .title-text {
        margin-left: 2px;
        color: #262626;
        font-weight: 600;
      }
    }
    .desc {
      color: #595959;
      margin-top: 4px;
      line-height: 22px;
      word-break: break-all;
    }
  }
}
</style>
