<style lang="less" scoped>
.segmentation-setting {
  .setting-item {
    margin-bottom: 16px;
    transition: all 0.2s;
    border: 1px solid #fff;

    .setting-item-header {
      display: flex;
      align-items: center;
      height: 72px;
      padding-left: 16px;
      border-radius: 2px;
      background-color: #f5f5f5;
      transition: all 0.2s;
      cursor: pointer;

      .setting-item-icon {
        width: 40px;
        height: 40px;
        background-repeat: no-repeat;
        background-size: cover;
      }

      .setting-item-info {
        padding-left: 8px;
      }

      .setting-item-name {
        line-height: 22px;
        font-size: 14px;
        font-weight: 600;
        color: #242933;
      }

      .setting-item-desc {
        line-height: 20px;
        padding-top: 2px;
        font-size: 12px;
        font-weight: 400;
        color: #3a4559;
      }
    }

    &:hover {
      border-radius: 4px;
      overflow: hidden;
      border: 1px solid #2475fc;

      .setting-item-header {
        background: #e5efff;
      }
    }

    &.active {
      border-radius: 4px;
      overflow: hidden;
      border: 1px solid #2475fc;

      .setting-item-header {
        background: #e5efff;
      }
    }

    .intelligent-icon {
      background-image: url('@/assets/img/library/document-segmentation/intelligent.svg');
    }

    .intelligent-icon.active {
      background-image: url('@/assets/img/library/document-segmentation/intelligent_active.svg');
    }

    .custom-icon {
      background-image: url('@/assets/img/library/document-segmentation/custom.svg');
    }

    .custom-icon.active {
      background-image: url('@/assets/img/library/document-segmentation/custom_active.svg');
    }
  }

  .setting-item-body {
    padding: 16px;
  }

  .sub-setting-item-name {
    height: 30px;
    line-height: 30px;
    padding-left: 8px;
    margin-bottom: 16px;
    border-radius: 4px;
    font-size: 14px;
    font-weight: 600;
    color: #242933;
    background-color: #f2f4f7;
  }

  .document-types {
    display: flex;
    margin-bottom: 16px;

    .document-type {
      display: inline-flex;
      justify-content: center;
      align-items: center;
      border-radius: 2px;
      padding: 5px 16px;
      margin-right: 8px;
      font-size: 14px;
      font-weight: 400;
      color: #595959;
      background: #fff;
      border: 1px solid #d9d9d9;
      cursor: pointer;
      transition: all 0.2s;

      &:last-child {
        margin-right: 0;
      }

      &.active {
        color: #2475fc;
        background: #f5f9ff;
        border: 1px solid #2475fc;
      }
    }
  }

  .custom-setting-form {
    .form-item{
      margin-bottom: 16px;
    }
    .form-item-label {
      line-height: 22px;
      margin-bottom: 4px;
      font-size: 14px;
      font-weight: 400;
      color: #262626;

      &::before {
        content: '*';
        font-size: 14px;
        line-height: 22px;
        padding-right: 2px;
        color: #fb363f;
      }
    }

    .unit-text {
      font-size: 14px;
      color: #595959;
    }
    .form-tip{
      color: #8c8c8c;
      font-size: 14px;
      font-weight: 400;
      line-height: 22px;
      margin-top: 2px;
    }
  }
  .excel-qa-form {
    .form-item{
      .form-item-label {
        padding-top: 5px;
      }
    }
  }

  .indexing-methods-box {
    .list-item {
      margin-top: 8px;
      padding: 16px;
      border: 1px solid #D9D9D9;
      position: relative;
      cursor: pointer;
      &.active{
        border: 2px solid #2475FC;
        .list-title-block {
          color: #2475fc;
          .svg-action {
            font-size: 16px;
            color: #2475fc;
          }
        }
        .check-icon{
          display: block;
        }
      }
      .check-icon {
        position: absolute;
        right: 0;
        bottom: 0;
        font-size: 18px;
        color: #fff;
        display: none;
      }
      .list-title-block {
        display: flex;
        align-items: center;
        font-size: 14px;
        font-weight: 600;
        line-height: 22px;
        color: #262626;
        .svg-action {
          font-size: 16px;
          margin-right: 4px;
          color: #262626;
        }
      }
      .list-content {
        margin-top: 4px;
        color: #595959;
        font-size: 14px;
        line-height: 22px;
      }
    }
  }

}
</style>

<template>
  <div class="segmentation-setting">
    <div class="setting-items">
      <div class="setting-item intelligent-setting" :class="{ active: formState.is_diy_split == 0 }">
        <div class="setting-item-header" @click="changeSettingType(0)">
          <span class="setting-item-icon intelligent-icon" :class="{ active: formState.is_diy_split == 0 }"></span>
          <div class="setting-item-info">
            <div class="setting-item-name">智能分段</div>
            <div class="setting-item-desc">系统自动分段，推荐不了解分段规则时使用</div>
          </div>
        </div>
        <div class="setting-item-body" v-if="formState.is_diy_split == 0">
          <div class="sub-setting-item-name">文档类型</div>

          <div class="document-types">
            <div class="document-type" :class="{ active: formState.is_qa_doc == 0 }" @click="changeDocumentType(0)">
              普通文档
            </div>
            <div class="document-type" :class="{ active: formState.is_qa_doc == 1 }" @click="changeDocumentType(1)">
              QA文档
            </div>
          </div>
          <a-space class="custom-setting-form" v-if="isHtmlOrDocx">
            <div class="form-item">
              <div class="form-item-label">提取图片</div>
              <div class="form-item-body">
                  <a-switch @change="onChange" v-model:checked="formState.enable_extract_image"></a-switch>
                  <div class="form-tip">开启提取图片功能，将自动从文档中提取图片，作为图片上方文本分段的附件。仅支持.docx文件和网页</div>
              </div>
            </div>
          </a-space>
          <template v-if="formState.is_qa_doc == 1">
            <div class="sub-setting-item-name">文件切分</div>
              <!-- 表格类型的QA文档 -->
            <div class="custom-setting-form excel-qa-form" v-if="props.mode == 1">
              <div class="form-item">
                <div class="form-item-label">问题所在列：</div>
                <div class="form-item-body">
                  <a-select v-model:value="formState.question_column" @change="onChagneFormInput" placeholder="请选择列名" style="width: 100%;">
                    <a-select-option v-for="item in props.excellQaLists" :value="item.value" :key="item.value">{{item.lable}}</a-select-option>
                  </a-select>
                </div>
              </div>
              <div class="form-item">
                <div class="form-item-label">答案所在列：</div>
                <div class="form-item-body">
                  <a-select v-model:value="formState.answer_column" @change="onChagneFormInput" placeholder="请选择答案所在列" style="width: 100%;">
                    <a-select-option v-for="item in props.excellQaLists" :value="item.value" :key="item.value">{{item.lable}}</a-select-option>
                  </a-select>
                </div>
              </div>
              <div class="sub-setting-item-name">索引方式</div>
              <div class="indexing-methods-box">
                <div class="list-item" :class="{active: formState.qa_index_type == 1}" @click="handleChangeQaIndexType(1)">
                  <svg-icon class="check-icon" name="check-arrow-filled"></svg-icon>
                  <div class="list-title-block">
                    <svg-icon name="file-search"></svg-icon>
                    问题与答案一起生成索引
                  </div>
                  <div class="list-content">回答用户提问时，将用户提问与导入的问题和答案一起对比相似度，根据相似度高的问题和答案回复</div>
                </div>
                <div class="list-item" :class="{active: formState.qa_index_type == 2}" @click="handleChangeQaIndexType(2)">
                  <svg-icon class="check-icon" name="check-arrow-filled"></svg-icon>
                  <div class="list-title-block">
                    <svg-icon name="comment-search"></svg-icon>
                    仅对问题生成索引
                  </div>
                  <div class="list-content">回答用户提问时，将用户提问与导入的问题一起对比相似度，再根据相似度高的问题和对应的答案来回复</div>
                </div>
              </div>
            </div>
            <div class="custom-setting-form" v-else>
              <a-space :size="16">
                <div class="form-item">
                  <div class="form-item-label">问题开始标识符：</div>
                  <div class="form-item-body">
                    <a-input placeholder="请输入标识符" v-model:value="formState.question_lable"
                      @change="onChagneFormInput" />
                  </div>
                </div>

                <div class="form-item">
                  <div class="form-item-label">答案开始标识符：</div>
                  <div class="form-item-body">
                    <a-input placeholder="请输入标识符" v-model:value="formState.answer_lable" @change="onChagneFormInput" />
                  </div>
                </div>
              </a-space>
              <div class="sub-setting-item-name">索引方式</div>
              <div class="indexing-methods-box">
                <div class="list-item" :class="{active: formState.qa_index_type == 1}" @click="handleChangeQaIndexType(1)">
                  <svg-icon class="check-icon" name="check-arrow-filled"></svg-icon>
                  <div class="list-title-block">
                    <svg-icon name="file-search"></svg-icon>
                    问题与答案一起生成索引
                  </div>
                  <div class="list-content">回答用户提问时，将用户提问与导入的问题和答案一起对比相似度，根据相似度高的问题和答案回复</div>
                </div>
                <div class="list-item" :class="{active: formState.qa_index_type == 2}" @click="handleChangeQaIndexType(2)">
                  <svg-icon class="check-icon" name="check-arrow-filled"></svg-icon>
                  <div class="list-title-block">
                    <svg-icon name="comment-search"></svg-icon>
                    仅对问题生成索引
                  </div>
                  <div class="list-content">回答用户提问时，将用户提问与导入的问题一起对比相似度，再根据相似度高的问题和对应的答案来回复</div>
                </div>
              </div>
            </div>
          </template>
        </div>
      </div>

      <div class="setting-item custom-setting" :class="{ active: formState.is_diy_split == 1 }" v-if="props.mode != 1">
        <div class="setting-item-header" @click="changeSettingType(1)">
          <span class="setting-item-icon custom-icon" :class="{ active: formState.is_diy_split == 1 }"></span>
          <div class="setting-item-info">
            <div class="setting-item-name">自定义分段</div>
            <div class="setting-item-desc">根据文档自行设置分段标识符、分段长度等</div>
          </div>
        </div>
        <div class="setting-item-body" v-if="formState.is_diy_split == 1">
          <div class="custom-setting-form subsection-form">
            <div class="form-item" style="margin-bottom: 18px">
              <div class="form-item-label">分段标识符：</div>
              <div class="form-item-body">
                <a-select placeholder="请选择" style="width: 100%" mode="multiple" v-model:value="formState.separators_no"
                  @change="onChagneFormInput">
                  <a-select-option :value="item.no" v-for="item in separatorsOptions" :key="item.no">{{ item.name
                    }}</a-select-option>
                </a-select>
              </div>
            </div>
            <a-space :size="16">
              <div class="form-item">
                <div class="form-item-label">分段最大长度：</div>
                <div class="form-item-body">
                  <a-space>
                    <a-input-number style="width: 140px" v-model:value="formState.chunk_size" placeholder="分段最大长度"
                      :min="200" :max="2000" :precision="0" :formatter="(value) => parseInt(value)"
                      :parser="(value) => parseInt(value)" @change="onChagneFormInput" /><span
                      class="unit-text">字符</span>
                  </a-space>
                </div>
              </div>

              <div class="form-item">
                <div class="form-item-label">分段重叠长度：</div>
                <div class="form-item-body">
                  <a-space>
                    <a-input-number style="width: 140px" v-model:value="formState.chunk_overlap" placeholder="分段重叠长度"
                      :min="2" :formatter="(value) => parseInt(value)" :parser="(value) => parseInt(value)"
                      @change="onChagneFormInput" /><span class="unit-text">字符</span>
                  </a-space>
                </div>
              </div>
            </a-space>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { getSeparatorsList } from '@/api/library/index'
import { reactive, ref, toRaw, computed } from 'vue'
import { Form } from 'ant-design-vue'
import { message } from 'ant-design-vue'

const useForm = Form.useForm
const emit = defineEmits(['change','validate'])

const props = defineProps({
  mode: {
    type: Number,
    default: 0
  },
  excellQaLists:{
    type: Array,
    default: () => []
  },
  libFileInfo:{
    type: Object,
    default: () => {}
  }
})
const isHtmlOrDocx = computed(()=>{
  return props.libFileInfo.doc_type == '1' && (props.libFileInfo.file_ext == 'docx' || props.libFileInfo.file_ext == 'html' )
})
const formState = reactive({
  is_diy_split: 0, // 0 智能分段 1 自定义分段
  separators_no: [], // 自定义分段-分隔符序号集
  chunk_size: 512, // 自定义分段-分段最大长度 默认512，最大值不得超过2000
  chunk_overlap: 50, // 自定义分段-分段重叠长度 默认为50，最小不得低于10，最大不得超过最大分段长度的50%
  is_qa_doc: 0, // 0 普通文档 1 QA文档
  question_lable: '', // QA文档-问题开始标识符
  answer_lable: '', // QA文档-答案开始标识符
  question_column: void 0, // excel QA文档 问题所在列
  answer_column: void 0, // excel QA文档 答案所在列
  qa_index_type: 1, // excel QA文档 索引方式
  enable_extract_image: true,
})

const formRules = reactive({
  question_lable: [
    {
      message: '请输入问题开始标识符',
      validator: async (rule, value) => {
        if (formState.is_diy_split == 0 && formState.is_qa_doc == 1 && props.mode == 0) {
          if (!value) {
            return Promise.reject('请输入问题开始标识符')
          }

          return Promise.resolve()
        }

        return Promise.resolve()
      }
    }
  ],
  answer_lable: [
    {
      message: '请输入答案开始标识符',
      validator: async (rule, value) => {
        if (formState.is_diy_split == 0 && formState.is_qa_doc == 1 && props.mode == 0) {
          if (!value) {
            return Promise.reject('请输入答案开始标识符')
          }

          return Promise.resolve()
        }

        return Promise.resolve()
      }
    }
  ],
  separators_no: [
    {
      message: '请选择分段标识符',
      validator: async (rule, value) => {
        if (formState.is_diy_split == 1 && value.length == 0) {
          return Promise.reject('请选择分段标识符')
        }

        return Promise.resolve()
      }
    }
  ],
  chunk_size: [
    {
      validator: async (rule, value) => {
        if (formState.is_diy_split == 1) {
          if (!value) {
            return Promise.reject('请输入分段最大长度')
          } else if (value > 2000) {
            return Promise.reject('最大分段长最大值不得超过2000')
          }

          return Promise.resolve()
        }

        return Promise.resolve()
      }
    }
  ],
  chunk_overlap: [
    {
      validator: async (rule, value) => {
        if (formState.is_diy_split == 1) {
          if (!value) {
            return Promise.reject('请输入分段重叠长度')
          } else if (value < 2) {
            return Promise.reject('分段重叠长度最小不得低于2')
          } else if (value > parseInt(formState.chunk_size / 2)) {
            return Promise.reject('分段重叠长度最大不得超过最大分段长度的50%')
          }

          return Promise.resolve()
        }

        return Promise.resolve()
      }
    }
  ],
  question_column: [
    {
      message: '请选择问题所在列',
      validator: async (rule, value) => {
        if (props.mode == 1 && formState.is_qa_doc == 1) {
          if (!value) {
            return Promise.reject('请选择问题所在列')
          }

          return Promise.resolve()
        }

        return Promise.resolve()
      }
    }
  ],
  answer_column: [
    {
      message: '请选择答案所在列',
      validator: async (rule, value) => {
        if (props.mode == 1 && formState.is_qa_doc == 1) {
          if (!value) {
            return Promise.reject('请选择答案所在列')
          }

          return Promise.resolve()
        }

        return Promise.resolve()
      }
    }
  ],
})

const { resetFields, validate, validateInfos } = useForm(formState, formRules)

// 切换分段模式
const changeSettingType = (val) => {
  formState.is_diy_split = val
  onChange(false)
}

// 切换文档类型
const changeDocumentType = (val) => {
  formState.is_qa_doc = val
  onChange(false)
}

const onChange = (showErrorAlert = true) => {
  let form = toRaw(formState)
  validate()
    .then(() => {
      emit('change', { ...form })
      emit('validate', '')
    })
    .catch((err) => {
      const { errorFields } = err
      if (errorFields.length > 0) {
        // console.log('error', err)
        if (showErrorAlert) {
          message.error(errorFields[0]['errors'][0])
        }
        emit('validate', errorFields[0]['errors'][0])
        return
      }
      emit('validate', '')
      emit('change', { ...form })
    })
}

// 查询分段节流
let formInEditTimer = null
const onChagneFormInput = () => {
  if (formInEditTimer) {
    clearTimeout(formInEditTimer)
    formInEditTimer = null
  }

  formInEditTimer = setTimeout(() => {
    onChange()
  }, 500)
}

// 分段标识符列表
const separatorsOptions = ref([])

const getSeparatorsOptions = () => {
  getSeparatorsList().then((res) => {
    separatorsOptions.value = res.data || []
  })
}

getSeparatorsOptions()

const handleChangeQaIndexType = (type) => {
  if(type == formState.qa_index_type){
    return
  }
  formState.qa_index_type = type;
  onChange();
}
</script>
