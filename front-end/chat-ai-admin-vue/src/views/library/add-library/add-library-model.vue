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

.indexing-methods-box {
  .list-item {
    margin-top: 8px;
    padding: 16px;
    border: 1px solid #d9d9d9;
    position: relative;
    cursor: pointer;
    &.active {
      border: 2px solid #2475fc;
      .list-title-block {
        color: #2475fc;
        .svg-action {
          font-size: 16px;
          color: #2475fc;
        }
      }
      .check-icon {
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

.main-title-block {
  margin: 16px 0;
  padding-bottom: 8px;
  font-size: 14px;
  font-weight: 600;
  border-bottom: 1px solid #d9d9d9;
}
</style>

<template>
  <a-modal
    v-model:open="visible"
    :title="t('add_library')"
    @ok="handleOk"
    :confirmLoading="saveLoading"
    wrapClassName="no-padding-modal"
    :bodyStyle="{ 'padding-right': '24px', 'max-height': '540px', 'overflow-y': 'auto' }"
    @cancel="handleCancel"
    width="700px"
  >
    <div class="form-box">
      <a-form layout="vertical" ref="formRef" :model="formState">
        <template v-if="formState.type != 3">
          <a-form-item ref="name"  name="library_name" :label="t('library_name')" :rules="[{ required: true, message: t('please_enter_library_name'), trigger: 'change' }]">
            <a-input
              v-model:value="formState.library_name"
              :placeholder="t('library_name_placeholder')"
              :maxlength="20"
            />
          </a-form-item>

          <a-form-item :label="t('library_intro')">
            <a-textarea
              :maxlength="1000"
              v-model:value="formState.library_intro"
              :placeholder="t('library_intro_placeholder')"
            />
          </a-form-item>

          <a-form-item ref="name" :label="t('library_cover')">
            <AvatarInput v-model:value="formState.avatar" @change="onAvatarChange" />
            <div class="form-item-tip">{{ t('upload_cover_tip') }}</div>
          </a-form-item>
        </template>
        <template v-else>
          <a-form-item :label="t('history_content_period')">
            <a-select v-model:value="formState.sync_official_history_type">
              <a-select-option :value="10">{{ t('all') }}</a-select-option>
              <a-select-option :value="1">{{ t('within_half_year') }}</a-select-option>
              <a-select-option :value="2">{{ t('within_one_year') }}</a-select-option>
              <a-select-option :value="3">{{ t('within_three_years') }}</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item :label="t('daily_fetch_latest')">
            <a-switch v-model:checked="formState.enable_cron_sync_official_content" :checked-children="t('on')" :un-checked-children="t('off')"></a-switch>
          </a-form-item>
        </template>
        <div class="hight-set-text" @click="isHide = !isHide">
          <CaretRightOutlined v-if="isHide" />
          <CaretDownOutlined v-else />
          {{ t('advanced_settings') }}
          <span style="color: #8c8c8c">{{ t('for_advanced_users') }}</span>
        </div>
        <a-form-item v-show="!isHide">
          <template #label
            >{{ t('vector_model') }}
            <a-tooltip>
              <template #title
                >{{ t('vector_model_tooltip') }}</template
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
        <a-form-item :label="t('indexing_method')" v-show="!isHide && formState.type == 2">
          <div class="indexing-methods-box">
            <div
              class="list-item"
              :class="{ active: formState.qa_index_type == 1 }"
              @click="handleChangeQaIndexType(1)"
            >
              <svg-icon class="check-icon" name="check-arrow-filled"></svg-icon>
              <div class="list-title-block">
                <svg-icon name="file-search"></svg-icon>
                {{ t('qa_indexing') }}
              </div>
              <div class="list-content">
                {{ t('qa_indexing_desc') }}
              </div>
            </div>
            <div
              class="list-item"
              :class="{ active: formState.qa_index_type == 2 }"
              @click="handleChangeQaIndexType(2)"
            >
              <svg-icon class="check-icon" name="check-arrow-filled"></svg-icon>
              <div class="list-title-block">
                <svg-icon name="comment-search"></svg-icon>
                {{ t('question_only_indexing') }}
              </div>
              <div class="list-content">
                {{ t('question_only_indexing_desc') }}
              </div>
            </div>
          </div>
        </a-form-item>
        <div v-show="!isHide && (formState.type == 0 || formState.type == 3)">
          <a-form-item :label="t('segmentation_method')" required>
            <div class="form-alert-tip" style="white-space: wrap;">
              {{ t('segmentation_tip') }}
            </div>
            <div class="select-card-box">
              <div
                class="select-card-item"
                @click="handleChangeSegmentationType(1)"
                :class="{ active: formState.chunk_type == 1 }"
              >
                <svg-icon class="check-arrow" name="check-arrow-filled"></svg-icon>
                <div class="card-title">
                  <svg-icon name="ordinary-segmentation" class="title-icon"></svg-icon>
                  {{ t('ordinary_segmentation') }}
                </div>
                <div class="card-desc">
                  {{ t('ordinary_segmentation_desc') }}
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
                  {{ t('semantic_segmentation') }}
                </div>
                <div class="card-desc">
                  {{ t('semantic_segmentation_desc') }}
                </div>
              </div>

              <div
                class="select-card-item"
                @click="handleChangeSegmentationType(4)"
                :class="{ active: formState.chunk_type == 4 }"
              >
                <svg-icon class="check-arrow" name="check-arrow-filled"></svg-icon>
                <div class="card-title">
                  <svg-icon name="semantic-segmentation" class="title-icon"></svg-icon>
                  {{ t('parent_child_segmentation') }}
                </div>
                <div class="card-desc">
                  {{ t('parent_child_segmentation_desc') }}
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
                  {{ t('ai_segmentation') }}
                </div>
                <div class="card-desc">
                  {{ t('ai_segmentation_desc') }}
                </div>
              </div>
            </div>
          </a-form-item>
          <template v-if="formState.chunk_type == 1">
            <a-form-item :label="t('segmentation_separator')" required>
              <a-select
                v-model:value="formState.normal_chunk_default_separators_no"
                mode="tags"
                style="width: 100%"
                :placeholder="t('select_segmentation_separator')"
                :options="segmentationTags"
              ></a-select>
            </a-form-item>
            <a-flex>
              <a-form-item :label="t('max_chunk_size')" required>
                <a-flex :gap="8" align="center">
                  <a-input-number
                    v-model:value="formState.normal_chunk_default_chunk_size"
                    style="width: 230px"
                    :precision="0"
                    :min="1"
                    :max="100000"
                  />
                  {{ t('characters') }}
                </a-flex>
              </a-form-item>
              <a-form-item :label="t('chunk_overlap')">
                <a-flex :gap="8" align="center">
                  <a-input-number
                    v-model:value="formState.normal_chunk_default_chunk_overlap"
                    style="width: 230px"
                    :precision="0"
                    :min="1"
                    :max="100000"
                  />
                  {{ t('characters') }}
                </a-flex>
              </a-form-item>
            </a-flex>
            <div class="graph-switch-box">
              <div class="lable-text">
                {{ t('auto_merge_small_chunks') }}
                <a-tooltip
                  :title="t('auto_merge_tooltip')"
                >
                  <QuestionCircleOutlined class="ml4" />
                </a-tooltip>
              </div>
              <a-switch
                checkedValue="false"
                unCheckedValue="true"
                v-model:checked="formState.normal_chunk_default_not_merged_text"
                :checked-children="t('on')"
                :un-checked-children="t('off')"
              />
            </div>
          </template>
          <template v-if="formState.chunk_type == 2">
            <a-form-item required v-if="formState.chunk_type == 2">
              <template #label>
                {{ t('segmentation_threshold') }}
                <a-tooltip>
                  <template #title
                    >{{ t('segmentation_threshold_tooltip') }}</template
                  >
                  <QuestionCircleOutlined style="cursor: pointer; margin-left: 2px" />
                </a-tooltip>
              </template>
              <a-input-number
                v-model:value="formState.semantic_chunk_default_threshold"
                style="width: 100%"
                :placeholder="t('please_enter_segmentation_threshold')"
                :precision="0"
                :min="0"
                :max="100"
              />
            </a-form-item>
            <a-flex>
              <a-form-item :label="t('max_chunk_size')" required>
                <a-flex :gap="8" align="center">
                  <a-input-number
                    v-model:value="formState.semantic_chunk_default_chunk_size"
                    style="width: 230px"
                    :precision="0"
                    :min="1"
                    :max="100000"
                  />
                  {{ t('characters') }}
                </a-flex>
              </a-form-item>
              <a-form-item :label="t('chunk_overlap')">
                <a-flex :gap="8" align="center">
                  <a-input-number
                    v-model:value="formState.semantic_chunk_default_chunk_overlap"
                    style="width: 230px"
                    :precision="0"
                    :min="1"
                    :max="100000"
                  />
                  {{ t('characters') }}
                </a-flex>
              </a-form-item>
            </a-flex>
          </template>

          <template v-if="formState.chunk_type == 3">
            <a-form-item required v-if="formState.chunk_type == 3">
              <template #label> {{ t('ai_llm_model') }} </template>
              <ModelSelect
                modelType="LLM"
                :placeholder="t('please_select_ai_model')"
                v-model:modeName="formState.ai_chunk_model"
                v-model:modeId="formState.ai_chunk_model_config_id"
                :modeName="formState.ai_chunk_model"
                :modeId="formState.ai_chunk_model_config_id"
                style="width: 100%"
                @loaded="onVectorModelLoaded"
              />
            </a-form-item>
            <a-form-item :label="t('prompt_setting')" required>
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
                {{ t('max_single_characters') }}
                <a-tooltip>
                  <template #title
                    >{{ t('max_single_characters_tooltip') }}</template
                  >
                  <QuestionCircleOutlined style="cursor: pointer; margin-left: 2px" />
                </a-tooltip>
              </template>
              <a-input-number
                class="form-item-inptu-numbner"
                v-model:value="formState.ai_chunk_size"
                :placeholder="t('please_enter_max_characters')"
                :precision="0"
                :min="0"
                :formatter="(value) => parseInt(value)"
                :parser="(value) => parseInt(value)"
              />
              {{ t('characters') }}
            </a-form-item>
          </template>

          <template v-if="formState.chunk_type == 4">
            <div class="main-title-block">{{ t('parent_block_context') }}</div>
            <a-form-item :label="t('chunk_type')">
              <a-radio-group v-model:value="formState.father_chunk_paragraph_type">
                <a-radio :value="1"
                  >{{ t('full_text') }}
                  <a-tooltip
                    :title="t('full_text_tooltip')"
                  >
                    <QuestionCircleOutlined />
                  </a-tooltip>
                </a-radio>
                <a-radio :value="2"
                  >{{ t('paragraph') }}
                  <a-tooltip
                    :title="t('paragraph_tooltip')"
                  >
                    <QuestionCircleOutlined />
                  </a-tooltip>
                </a-radio>
              </a-radio-group>
            </a-form-item>

            <a-form-item :label="t('segmentation_separator')" v-if="formState.father_chunk_paragraph_type == 2">
              <a-select
                :placeholder="t('please_select')"
                style="width: 100%"
                mode="tags"
                :options="segmentationTags"
                v-model:value="formState.father_chunk_separators_no"
              >
              </a-select>
            </a-form-item>
            <a-form-item :label="t('max_chunk_size')" v-if="formState.father_chunk_paragraph_type == 2">
              <a-flex align="center" :gap="8">
                <a-input-number
                  style="flex: 1"
                  v-model:value="formState.father_chunk_chunk_size"
                  :placeholder="t('max_chunk_size')"
                  :min="200"
                  :max="10000"
                  :precision="0"
                  :formatter="(value) => parseInt(value)"
                  :parser="(value) => parseInt(value)"
                /><span class="unit-text">{{ t('characters') }}</span>
              </a-flex>
            </a-form-item>
            <div class="main-title-block">{{ t('child_block_retrieval') }}</div>

            <a-form-item :label="t('segmentation_separator')">
              <a-select
                :placeholder="t('please_select')"
                style="width: 100%"
                mode="tags"
                :options="segmentationTags"
                v-model:value="formState.son_chunk_separators_no"
              >
              </a-select>
            </a-form-item>
            <a-form-item :label="t('max_chunk_size')">
              <a-flex align="center" :gap="8">
                <a-input-number
                  style="flex: 1"
                  v-model:value="formState.son_chunk_chunk_size"
                  :placeholder="t('max_chunk_size')"
                  :min="200"
                  :max="10000"
                  :precision="0"
                  :formatter="(value) => parseInt(value)"
                  :parser="(value) => parseInt(value)"
                /><span class="unit-text">{{ t('characters') }}</span>
              </a-flex>
            </a-form-item>
          </template>

          <div class="graph-switch-box" v-show="neo4j_status">
            <div class="lable-text">
              {{ t('generate_knowledge_graph') }}
              <a-tooltip>
                <template #title
                  >{{ t('generate_knowledge_graph_tooltip') }}</template
                >
                <QuestionCircleOutlined class="ml4" />
              </a-tooltip>
            </div>
            <a-switch
              @change="handleGraphSwitch"
              :checked="formState.graph_switch"
              :checked-children="t('on')"
              :un-checked-children="t('off')"
            />
          </div>
          <a-form-item :label="t('knowledge_graph_model')" v-show="formState.graph_switch && neo4j_status">
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
import {createLibrary, createOfficialLibrary, getSeparatorsList} from '@/api/library/index'
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
import { useI18n } from '@/hooks/web/useI18n'

import { useCompanyStore } from '@/stores/modules/company'

const { t } = useI18n('views.library.add-library.add-library-model')
const companyStore = useCompanyStore()
const neo4j_status = computed(() => {
  return companyStore.companyInfo?.neo4j_status == 'true'
})

// 设置全局默认的duration为（2秒）
message.config({
  duration: 2
})

const formRef = ref(null)
const isHide = ref(true)

const router = useRouter()
const visible = ref(false)

const saveLoading = ref(false)
const default_ai_chunk_prumpt = t('default_ai_chunk_prompt')
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
  normal_chunk_default_separators_no: [12, 11],
  normal_chunk_default_chunk_size: 512,
  normal_chunk_default_chunk_overlap: 50,
  semantic_chunk_default_chunk_size: 512,
  semantic_chunk_default_chunk_overlap: 50,
  semantic_chunk_default_threshold: 90,
  normal_chunk_default_not_merged_text: 'false',
  graph_switch: false,
  graph_model_config_id: void 0,
  graph_use_model: '',
  ai_chunk_size: 5000, // ai大模型分段最大字符数
  ai_chunk_model: '', // ai大模型分段模型名称
  ai_chunk_model_config_id: '', // ai大模型分段模型配置id
  ai_chunk_prumpt: default_ai_chunk_prumpt, // ai大模型分段提示词设置
  group_id: 0,
  father_chunk_paragraph_type: 2,
  father_chunk_separators_no: [],
  father_chunk_chunk_size: 1024,
  son_chunk_separators_no: [],
  son_chunk_chunk_size: 512,

  sync_official_history_type: 2,
  enable_cron_sync_official_content: true,
  app_id_list: '',
})

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
  formRef.value.validate()
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
  formData.append('qa_index_type', formState.qa_index_type)
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
    JSON.stringify(formState.normal_chunk_default_separators_no)
  )
  formData.append(
    'normal_chunk_default_not_merged_text',
    formState.normal_chunk_default_not_merged_text
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

  formData.append('father_chunk_paragraph_type', formState.father_chunk_paragraph_type)
  formData.append(
    'father_chunk_separators_no',
    JSON.stringify(formState.father_chunk_separators_no)
  )
  formData.append('son_chunk_separators_no', JSON.stringify(formState.son_chunk_separators_no))
  formData.append('father_chunk_chunk_size', formState.father_chunk_chunk_size)
  formData.append('son_chunk_chunk_size', formState.son_chunk_chunk_size)

  formData.append('group_id', formState.group_id)
  formData.append('is_default', 2)
  let request = createLibrary
  if (formState.type == 3) {
    request = createOfficialLibrary
    formData.append('app_id_list', formState.app_id_list)
    formData.append('sync_official_history_type', formState.sync_official_history_type)
    formData.append('enable_cron_sync_official_content', Number(formState.enable_cron_sync_official_content))
  }
  saveLoading.value = true
  request(formData)
    .then((res) => {
      message.success(t('create_success'))
      visible.value = false

      let path = '/library/details/knowledge-document'
      let query
      if (formState.type == 3) {
        query = {id: res?.data?.[0]}
      } else {
        query = {id: res?.data?.id}
        if (formState.doc_type == 3) {
          path = '/library/preview'
          query = {
            id: res.data.file_ids[0]
          }
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

const show = ({ type, group_id, wx_app_ids }) => {
  formState.type = type
  if (type == 3) {
    formState.app_id_list = wx_app_ids
  } else if (type == 0) {
    formState.avatar = LIBRARY_NORMAL_AVATAR
    formState.avatar_file = LIBRARY_NORMAL_AVATAR
  } else {
    formState.avatar = LIBRARY_QA_AVATAR
    formState.avatar_file = LIBRARY_QA_AVATAR
  }

  formState.library_name = ''
  formState.library_intro = ''
  formState.qa_index_type = 1
  formState.chunk_type = 1
  formState.normal_chunk_default_separators_no = [12, 11]
  formState.normal_chunk_default_chunk_size = 512
  formState.normal_chunk_default_chunk_overlap = 50
  formState.normal_chunk_default_not_merged_text = 'false'
  formState.semantic_chunk_default_chunk_size = 512
  formState.semantic_chunk_default_chunk_overlap = 50
  formState.semantic_chunk_default_threshold = 50
  formState.graph_switch = false

  formState.father_chunk_paragraph_type = 2
  formState.father_chunk_separators_no = [12, 11]
  formState.father_chunk_chunk_size = 1024
  formState.son_chunk_separators_no = [8, 10]
  formState.son_chunk_chunk_size = 512

  formState.ai_chunk_size = 5000
  formState.ai_chunk_model = ''
  formState.ai_chunk_model_config_id = ''
  formState.ai_chunk_prumpt = default_ai_chunk_prumpt
  formState.group_id = group_id || 0
  visible.value = true
}

const handleChangeQaIndexType = (type) => {
  if (type == formState.qa_index_type) {
    return
  }
  formState.qa_index_type = type
}
defineExpose({
  show
})

onMounted(() => {})
</script>
