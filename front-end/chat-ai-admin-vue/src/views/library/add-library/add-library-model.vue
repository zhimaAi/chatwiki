<style lang="less" scoped>
.form-box {
  margin-top: 12px;
  margin-bottom: 24px;
}

.form-item-tip {
  color: #999;
}

.form-alert-tip {
  color: #8c8c8c;
  font-size: 12px;
  font-weight: 400;
  line-height: 14px;
  margin: 2px 0 6px;
  white-space: nowrap;
}

.hight-set-text {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  width: fit-content;
  color: rgba(0, 0, 0, 0.88);
  margin-bottom: 16px;
}
.ml4 {
  margin-left: 4px;
}
.graph-switch-box {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.select-card-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  width: 550px;
  .select-card-item {
    width: 100%;
    position: relative;
    padding: 16px;
    border-radius: 6px;
    border: 1px solid #d9d9d9;
    cursor: pointer;
    .check-arrow {
      position: absolute;
      display: block;
      right: -1px;
      bottom: -1px;
      width: 24px;
      height: 24px;
      font-size: 24px;
      color: #fff;
      opacity: 0;
      transition: all 0.2s cubic-bezier(0.8, 0, 1, 1);
    }
    .card-title {
      display: flex;
      align-items: center;
      line-height: 22px;
      margin-bottom: 4px;
      color: #262626;
      font-weight: 600;
      font-size: 14px;
    }
    .title-icon {
      margin-right: 4px;
      font-size: 16px;
    }
    .card-desc {
      line-height: 22px;
      font-size: 14px;
      color: #595959;
    }

    &.active {
      background: var(--01-, #e5efff);
      border: 2px solid #2475fc;
      .check-arrow {
        opacity: 1;
      }
      .card-title {
        color: #2475fc;
      }
    }
  }
}
</style>

<template>
  <a-modal
    v-model:open="visible"
    title="添加知识库"
    @ok="handleOk"
    :confirmLoading="saveLoading"
    wrapClassName="no-padding-modal"
    :bodyStyle="{ 'padding-right': '24px', 'max-height': '540px', 'overflow-y': 'auto' }"
    @cancel="handleCancel"
    width="620px"
  >
    <div class="form-box">
      <a-form layout="vertical">
        <a-form-item ref="name" label="知识库名称" v-bind="validateInfos.library_name">
          <a-input
            v-model:value="formState.library_name"
            placeholder="请输入知识库名称，最多20个字"
            :maxlength="20"
          />
        </a-form-item>

        <a-form-item label="知识库简介" v-bind="validateInfos.library_intro">
          <a-textarea
            :maxlength="1000"
            v-model:value="formState.library_intro"
            placeholder="请输入知识库介绍"
          />
        </a-form-item>

        <a-form-item ref="name" label="知识库封面" v-bind="validateInfos.avatar">
          <AvatarInput v-model:value="formState.avatar" @change="onAvatarChange" />
          <div class="form-item-tip">请上传知识库封面，建议尺寸为100*100px.大小不超过100kb</div>
        </a-form-item>
        <div class="hight-set-text" @click="isHide = !isHide">
          <CaretRightOutlined v-if="isHide" />
          <CaretDownOutlined v-else />
          高级设置
          <span style="color: #8c8c8c">适用于高级用户</span>
        </div>
        <a-form-item v-show="!isHide">
          <template #label
            >文档解析向量模型
            <a-tooltip>
              <template #title
                >向量模型可以将分段数据转化为向量形式，存储到向量数据库中，便于后续根据用户问题查询匹配。</template
              >
              <QuestionCircleOutlined class="ml4" />
            </a-tooltip>
          </template>
          <ModelSelect
            modelType="TEXT EMBEDDING"
            v-model:modeName="formState.use_model"
            v-model:modeId="formState.model_config_id"
            @loaded="onTextModelLoaded"
            style="width: 100%"
          />
        </a-form-item>
        <div v-show="!isHide && formState.type == 0">
          <a-form-item label="分段方式" required>
            <div class="form-alert-tip">提示：语义分段更适合没有排版过的文章，即没有明显换行符号的文本，否则更推荐使用普通分段</div>
            <div class="select-card-box">
              <div
                class="select-card-item"
                @click="handleChangeSegmentationType(1)"
                :class="{ active: formState.chunk_type == 1 }"
              >
                <svg-icon class="check-arrow" name="check-arrow-filled"></svg-icon>
                <div class="card-title">
                  <svg-icon name="ordinary-segmentation" class="title-icon"></svg-icon>
                  普通分段
                </div>
                <div class="card-desc">
                  基于文章中句号、空行，或者自定义符号进行分段，不会消耗模型token
                </div>
              </div>
              <div
                class="select-card-item"
                @click="handleChangeSegmentationType(2)"
                :class="{ active: formState.chunk_type == 2 }"
              >
                <svg-icon class="check-arrow" name="check-arrow-filled"></svg-icon>
                <div class="card-title">
                  <svg-icon name="semantic-segmentation" class="title-icon"></svg-icon>
                  语义分段
                </div>
                <div class="card-desc">
                  将文章拆分成句子后，通过语句向量相似度进行分段，会消耗模型token
                </div>
              </div>
              <div
                class="select-card-item"
                @click="handleChangeSegmentationType(3)"
                :class="{ active: formState.chunk_type == 3 }"
              >
                <svg-icon class="check-arrow" name="check-arrow-filled"></svg-icon>
                <div class="card-title">
                  <svg-icon name="semantic-segmentation" class="title-icon"></svg-icon>
                  AI分段
                </div>
                <div class="card-desc">
                  将文章提交给大模型，大模型基于设定的提示词进行分段，会消耗大量模型token
                </div>
              </div>
            </div>
          </a-form-item>
          <template v-if="formState.chunk_type == 1">
            <a-form-item label="分段标识符" required>
              <a-select
                v-model:value="formState.normal_chunk_default_separators_no"
                mode="multiple"
                style="width: 100%"
                placeholder="请选择分段标识符"
                :options="segmentationTags"
              ></a-select>
            </a-form-item>
            <a-flex>
              <a-form-item label="分段最大长度" required>
                <a-flex :gap="8" align="center">
                  <a-input-number
                    v-model:value="formState.normal_chunk_default_chunk_size"
                    style="width: 230px"
                    :precision="0"
                    :min="1"
                    :max="100000"
                  />
                  字符
                </a-flex>
              </a-form-item>
              <a-form-item label="分段重叠长度">
                <a-flex :gap="8" align="center">
                  <a-input-number
                    v-model:value="formState.normal_chunk_default_chunk_overlap"
                    style="width: 230px"
                    :precision="0"
                    :min="1"
                    :max="100000"
                  />
                  字符
                </a-flex>
              </a-form-item>
            </a-flex>
          </template>
          <template v-if="formState.chunk_type == 2">
            <a-form-item required v-if="formState.chunk_type == 2">
              <template #label>
                分段阈值
                <a-tooltip>
                  <template #title
                    >用于控制分段拆分的标准，数值0~100,数值越大，分段越少，数值越小，分段越多。</template
                  >
                  <QuestionCircleOutlined style="cursor: pointer; margin-left: 2px" />
                </a-tooltip>
              </template>
              <a-input-number
                v-model:value="formState.semantic_chunk_default_threshold"
                style="width: 100%"
                placeholder="请输入分段阈值"
                :precision="0"
                :min="0"
                :max="100"
              />
            </a-form-item>
            <a-flex>
              <a-form-item label="分段最大长度" required>
                <a-flex :gap="8" align="center">
                  <a-input-number
                    v-model:value="formState.semantic_chunk_default_chunk_size"
                    style="width: 230px"
                    :precision="0"
                    :min="1"
                    :max="100000"
                  />
                  字符
                </a-flex>
              </a-form-item>
              <a-form-item label="分段重叠长度">
                <a-flex :gap="8" align="center">
                  <a-input-number
                    v-model:value="formState.semantic_chunk_default_chunk_overlap"
                    style="width: 230px"
                    :precision="0"
                    :min="1"
                    :max="100000"
                  />
                  字符
                </a-flex>
              </a-form-item>
            </a-flex>
          </template>

          <template v-if="formState.chunk_type == 3">
            <a-form-item required v-if="formState.chunk_type == 3">
              <template #label>
                AI大模型
              </template>
              <ModelSelect
                modelType="LLM"
                placeholder="请选择AI大模型"
                v-model:modeName="formState.ai_chunk_model"
                v-model:modeId="formState.ai_chunk_model_config_id"
                :modeName="formState.ai_chunk_model"
                :modeId="formState.ai_chunk_model_config_id"
                style="width: 100%"
                @loaded="onVectorModelLoaded"
              />
            </a-form-item>
            <a-form-item label="提示词设置" required>
              <a-flex :gap="8" align="center">
                <a-textarea
                  :maxLength="500"
                  style="height: 80px"
                  v-model:value="formState.ai_chunk_prumpt"
                  :placeholder="default_ai_chunk_prumpt"
                />
              </a-flex>
            </a-form-item>
            <a-form-item>
              <template #label>
                单次最大字符数
                <a-tooltip>
                  <template #title>由于大模型支持的上下文数量有限制，如果上传的文档较大，会按照最大字符数先拆分成多个分段，再提交给大模型进行分段。</template>
                  <QuestionCircleOutlined style="cursor: pointer; margin-left: 2px" />
                </a-tooltip>
              </template>
              <a-input-number
                class="form-item-inptu-numbner"
                v-model:value="formState.ai_chunk_size"
                placeholder="请输入单次最大字符数"
                :precision="0"
                :min="0"
                :formatter="(value) => parseInt(value)"
                :parser="(value) => parseInt(value)"
              />
              字符
            </a-form-item>
          </template>

          <div class="graph-switch-box" v-show="neo4j_status">
            <div class="lable-text">
              生成知识图谱
              <a-tooltip>
                <template #title
                  >文件分块后，所有块将用于生成知识图谱，提高复杂问题推理能力。</template
                >
                <QuestionCircleOutlined class="ml4" />
              </a-tooltip>
            </div>
            <a-switch
              @change="handleGraphSwitch"
              :checked="formState.graph_switch"
              checked-children="开"
              un-checked-children="关"
            />
          </div>
          <a-form-item label="知识图谱模型" v-show="formState.graph_switch && neo4j_status">
            <ModelSelect
              modelType="LLM"
              v-model:modeName="formState.graph_use_model"
              v-model:modeId="formState.graph_model_config_id"
              style="width: 300px"
              @loaded="onVectorModelLoaded"
            />
          </a-form-item>
        </div>
      </a-form>
    </div>
  </a-modal>
  <OpenGrapgModal @ok="handleOpenGrapgOk" ref="openGrapgModalRef" />
</template>

<script setup>
import { reactive, ref, onMounted, computed, nextTick } from 'vue'
import { Form, message } from 'ant-design-vue'
import { useRouter } from 'vue-router'
import { createLibrary, getSeparatorsList } from '@/api/library/index'
import AvatarInput from './components/avatar-input.vue'
import { LIBRARY_NORMAL_AVATAR, LIBRARY_QA_AVATAR } from '@/constants/index'
import { transformUrlData } from '@/utils/validate.js'
import {
  CaretRightOutlined,
  CaretDownOutlined,
  QuestionCircleOutlined
} from '@ant-design/icons-vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import OpenGrapgModal from '@/views/library/library-details/components/open-grapg-modal.vue'

import { useCompanyStore } from '@/stores/modules/company'
const companyStore = useCompanyStore()
const neo4j_status = computed(()=>{
  return companyStore.companyInfo?.neo4j_status == 'true'
})


// 设置全局默认的duration为（2秒）
message.config({
  duration: 2
})

const isHide = ref(true)

const router = useRouter()
const visible = ref(false)

const useForm = Form.useForm
const saveLoading = ref(false)
const default_ai_chunk_prumpt = '你是一位文章分段助手，根据文章内容的语义进行合理分段，确保每个分段表述一个完整的语义，每个分段字数控制在500字左右，最大不超过1000字。请严格按照文章内容进行分段，不要对文章内容进行加工，分段完成后输出分段后的内容。'
const formState = reactive({
  type: '0',
  access_rights: 0,
  library_name: '',
  library_intro: '',
  use_model: '',
  model_config_id: '',
  library_files: undefined,
  avatar: '',
  avatar_file: undefined,
  is_offline: false,
  urls: '',
  doc_type: 1,
  file_name: '',
  is_qa_doc: 0, // 0 普通文档 1QA文档
  qa_index_type: 1, // 1问题与答案一起生成索引 2仅对问题生成索引
  doc_auto_renew_frequency: 1,
  chunk_type: 1,
  normal_chunk_default_separators_no: [10, 12],
  normal_chunk_default_chunk_size: 512,
  normal_chunk_default_chunk_overlap: 50,
  semantic_chunk_default_chunk_size: 512,
  semantic_chunk_default_chunk_overlap: 50,
  semantic_chunk_default_threshold: 90,
  graph_switch: false,
  graph_model_config_id: void 0,
  graph_use_model: '',
  ai_chunk_size: 5000, // ai大模型分段最大字符数
  ai_chunk_model:'', // ai大模型分段模型名称
  ai_chunk_model_config_id: '', // ai大模型分段模型配置id
  ai_chunk_prumpt: default_ai_chunk_prumpt // ai大模型分段提示词设置
})

const rules = reactive({
  library_name: [{ required: true, message: '请输入库名称', trigger: 'change' }]
})

const { validate, validateInfos } = useForm(formState, rules)

const onAvatarChange = (data) => {
  formState.avatar = data.imageUrl
  formState.avatar_file = data.file
}

const segmentationTags = ref([])
getSeparatorsList().then((res) => {
  segmentationTags.value = res.data.map((item) => {
    return {
      label: item.name,
      value: item.no
    }
  })
})

const handleChangeSegmentationType = (type) => {
  formState.chunk_type = type
}

const vectorModelList = ref([])
const onVectorModelLoaded = (list) => {
  vectorModelList.value = list

  nextTick(() => {
    if (!formState.ai_chunk_model || !Number(formState.ai_chunk_model_config_id)) {
      setDefaultModel()
    }
  })
}

const setDefaultModel = () => {
  if (vectorModelList.value.length > 0) {
    // 遍历查找chatwiki模型
    let modelConfig = null
    let model = null

    // 云版默认选中qwen-max
    for (let item of vectorModelList.value) {
      if (item.model_define === 'tongyi') {
        modelConfig = item
        for (let child of modelConfig.children) {
          if (child.name === 'qwen-max') {
            model = child
            break
          }
        }
        break
      }
    }

    if (!modelConfig) {
      modelConfig = vectorModelList.value[0]
      model = modelConfig.children[0]
    }

    if (modelConfig && model) {
      formState.ai_chunk_model = model.name
      formState.ai_chunk_model_config_id = model.model_config_id
    }
  }
}

const onTextModelLoaded = (list) => {
  if (list.length) {
    formState.use_model = list[0].children[0].name
    formState.model_config_id = list[0].model_config_id
  }
}

const openGrapgModalRef = ref(null)
const handleGraphSwitch = (val) => {
  if (val) {
    let data = {
      graph_model_config_id: formState.graph_model_config_id,
      graph_use_model: formState.graph_use_model
    }
    if (!formState.graph_model_config_id || !formState.graph_use_model) {
      if (vectorModelList.value.length > 0) {
        let modelConfig = vectorModelList.value[0]
        if (modelConfig) {
          let model = modelConfig.children[0]
          data.graph_use_model = model.name
          data.graph_model_config_id = model.model_config_id
        }
      }
    }
    openGrapgModalRef.value.show(data)
  } else {
    formState.graph_switch = false
  }
}

const handleOpenGrapgOk = (data) => {
  if (data.graph_model_config_id) {
    formState.graph_switch = true
    formState.graph_model_config_id = data.graph_model_config_id
    formState.graph_use_model = data.graph_use_model
  }
}

const handleCancel = () => {
  visible.value = false
}

const handleOk = () => {
  validate()
    .then(() => {
      saveForm()
    })
    .catch((err) => {
      console.log(err, 'err')
    })
}

const saveForm = () => {
  let formData = new FormData()

  let newFormState = JSON.parse(JSON.stringify(formState)) // 深拷贝，不能改变原对象

  formData.append('type', formState.type)
  formData.append('access_rights', formState.access_rights)
  formData.append('library_name', formState.library_name)
  formData.append('library_intro', formState.library_intro)
  formData.append('use_model', newFormState.use_model)
  formData.append('model_config_id', formState.model_config_id)
  formData.append('avatar', formState.avatar_file ? formState.avatar_file : '')
  formData.append('is_offline', formState.is_offline)
  formData.append('urls', JSON.stringify(transformUrlData(formState.urls)))
  formData.append('doc_type', formState.doc_type)
  formData.append('file_name', formState.file_name)
  formData.append('is_qa_doc', formState.is_qa_doc)
  formData.append('qa_index_type', formState.qa_index_type)
  formData.append('doc_auto_renew_frequency', formState.doc_auto_renew_frequency)
  formData.append('chunk_type', formState.chunk_type)
  formData.append(
    'normal_chunk_default_separators_no',
    formState.normal_chunk_default_separators_no.join(',')
  )
  formData.append('normal_chunk_default_chunk_size', formState.normal_chunk_default_chunk_size)
  formData.append(
    'normal_chunk_default_chunk_overlap',
    formState.normal_chunk_default_chunk_overlap
  )
  formData.append('semantic_chunk_default_chunk_size', formState.semantic_chunk_default_chunk_size)
  formData.append(
    'semantic_chunk_default_chunk_overlap',
    formState.semantic_chunk_default_chunk_overlap
  )
  formData.append('semantic_chunk_default_threshold', formState.semantic_chunk_default_threshold)
  formData.append('graph_switch', formState.graph_switch ? 1 : 0)
  formData.append('graph_model_config_id', formState.graph_model_config_id)
  formData.append('graph_use_model', formState.graph_use_model)

  formData.append('ai_chunk_size', formState.ai_chunk_size)
  formData.append('ai_chunk_model', formState.ai_chunk_model)
  formData.append('ai_chunk_model_config_id', formState.ai_chunk_model_config_id)
  formData.append('ai_chunk_prumpt', formState.ai_chunk_prumpt)

  saveLoading.value = true

  createLibrary(formData)
    .then((res) => {
      message.success('创建成功')
      visible.value = false

      let path = '/library/details/knowledge-document'
      let query = {
        id: res.data.id
      }

      if (formState.doc_type == 3) {
        path = '/library/preview'
        query = {
          id: res.data.file_ids[0]
        }
      }

      router.push({
        path,
        query
      })
      saveLoading.value = false
      visible.value = false
    })
    .catch(() => {
      saveLoading.value = false
    })
}

const show = ({ type }) => {
  formState.type = type
  if (type == 0) {
    formState.avatar = LIBRARY_NORMAL_AVATAR
    formState.avatar_file = LIBRARY_NORMAL_AVATAR
  } else {
    formState.avatar = LIBRARY_QA_AVATAR
    formState.avatar_file = LIBRARY_QA_AVATAR
  }

  formState.library_name = ''
  formState.library_intro = ''
  formState.chunk_type = 1
  formState.normal_chunk_default_separators_no = [10, 12]
  formState.normal_chunk_default_chunk_size = 512
  formState.normal_chunk_default_chunk_overlap = 50
  formState.semantic_chunk_default_chunk_size = 512
  formState.semantic_chunk_default_chunk_overlap = 50
  formState.semantic_chunk_default_threshold = 50
  formState.graph_switch = false

  formState.ai_chunk_size = 5000
  formState.ai_chunk_model = ''
  formState.ai_chunk_model_config_id = ''
  formState.ai_chunk_prumpt = default_ai_chunk_prumpt

  visible.value = true
}

defineExpose({
  show
})

onMounted(() => {})
</script>
