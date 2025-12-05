<style lang="less" scoped>
.add-library-page {
  padding: 24px;

  .form-box {
    width: 630px;
    margin: 0 auto;
  }
}

.model-icon {
  height: 18px;
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

.card-box {
  display: flex;
  justify-content: space-between;
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
    color: #bfbfbf;
  }
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
  width: 720px;
}
</style>

<template>
  <cu-scroll>
    <div class="add-library-page">
      <div class="form-box">
        <a-form :label-col="{ span: 6 }">
          <a-form-item ref="name" label="知识库名称" v-bind="validateInfos.library_name">
            <a-input
              @blur="handleEdit"
              v-model:value="formState.library_name"
              placeholder="请输入知识库名称，最多20个字"
            />
          </a-form-item>

          <a-form-item label="知识库简介">
            <a-textarea
              @blur="handleEdit"
              v-model:value="formState.library_intro"
              placeholder="请输入知识库介绍"
            />
          </a-form-item>

          <a-form-item ref="name" label="知识库封面" v-bind="validateInfos.avatar">
            <AvatarInput v-model:value="formState.avatar" @change="onAvatarChange" />
            <div class="form-item-tip">请上传知识库封面，建议尺寸为100*100px.大小不超过100kb</div>
          </a-form-item>

          <a-form-item label="嵌入模型" v-bind="validateInfos.use_model">
            <div class="card-box" v-if="false">
              <div
                class="use-model-item"
                :class="{ active: isActive == item.value }"
                v-for="item in libraryModeList"
                :key="item.value"
                @click="handleSelectLibrary(item)"
              >
                <div class="use-model-item-top" :class="{ active: isActive == item.value }">
                  <svg-icon
                    class="use-model-item-top-icon"
                    style="color: #2475fc"
                    :name="item.iconName"
                  ></svg-icon>
                  <p>{{ item.title }}</p>
                </div>
                <p>{{ item.desc }}</p>
                <svg-icon
                  class="check-arrow"
                  name="check-arrow-filled"
                  style="font-size: 24px; color: #fff"
                  v-if="isActive == item.value"
                ></svg-icon>
              </div>
            </div>
            <!-- 自定义选择器 -->
            <CustomSelector
              v-model="formState.use_model"
              label-key="use_model_name"
              value-key="value"
              :modelType="'TEXT EMBEDDING'"
              :model-config-id="formState.model_config_id"
              @change="handleModelChange"
              @loaded="onVectorModelLoaded"
            />
            <!-- <a-select
              @change="handleChangeModel"
              v-model:value="formState.use_model"
              placeholder="请选择嵌入模型"
            >
              <a-select-opt-group v-for="item in modelList" :key="item.id">
                <template #label>
                  <a-flex align="center" :gap="6">
                    <img class="model-icon" :src="item.icon" alt="" />{{ item.name }}
                  </a-flex>
                </template>
                <a-select-option
                  :value="
                    modelDefine.indexOf(item.model_define) > -1 && val.deployment_name
                      ? val.deployment_name
                      : val.name + val.id
                  "
                  :model_config_id="item.id"
                  :current_obj="val"
                  v-for="val in item.children"
                  :key="val.name + val.id"
                >
                  <span v-if="modelDefine.indexOf(item.model_define) > -1 && val.deployment_name">{{
                    val.deployment_name
                  }}</span>
                  <span v-else>{{ val.name }}</span>
                </a-select-option>
              </a-select-opt-group>
            </a-select> -->
          </a-form-item>
          <a-form-item label="生成知识图谱" v-show="neo4j_status">
            <a-switch
              @change="handleGraphSwitch"
              :checked="formState.graph_switch"
              checked-children="开"
              un-checked-children="关"
            />
            <div class="form-item-tip">开启后，可以文档列表手动点击知识图谱学习生成知识图谱</div>
          </a-form-item>
          <a-form-item label="知识图谱模型" v-show="formState.graph_switch && neo4j_status">
            <ModelSelect
              modelType="LLM"
              v-model:modeName="formState.graph_use_model"
              v-model:modeId="formState.graph_model_config_id"
              style="width: 300px"
              @change="onChangeModel"
              @loaded="onVectorModelLoaded"
            />
          </a-form-item>
          <a-form-item label="索引方式" v-if="isQaLibrary">
            <div class="indexing-methods-box">
              <div
                class="list-item"
                :class="{ active: formState.qa_index_type == 1 }"
                @click="handleChangeQaIndexType(1)"
              >
                <svg-icon class="check-icon" name="check-arrow-filled"></svg-icon>
                <div class="list-title-block">
                  <svg-icon name="file-search"></svg-icon>
                  问题与答案一起生成索引
                </div>
                <div class="list-content">
                  回答用户提问时，将用户提问与导入的问题和答案一起对比相似度，根据相似度高的问题和答案回复
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
                  仅对问题生成索引
                </div>
                <div class="list-content">
                  回答用户提问时，将用户提问与导入的问题一起对比相似度，再根据相似度高的问题和对应的答案来回复
                </div>
              </div>
            </div>
          </a-form-item>
          <template v-if="!isQaLibrary">
            <a-form-item label="分段方式" required>
              <div class="form-alert-tip">
                提示：语义分段更适合没有排版过的文章，即没有明显换行符号的文本，否则更推荐使用普通分段
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
                  v-if="!isQaLibrary"
                  class="select-card-item"
                  @click="handleChangeSegmentationType(4)"
                  :class="{ active: formState.chunk_type == 4 }"
                >
                  <svg-icon class="check-arrow" name="check-arrow-filled"></svg-icon>
                  <div class="card-title">
                    <svg-icon name="semantic-segmentation" class="title-icon"></svg-icon>
                    父子分段
                  </div>
                  <div class="card-desc">
                    基于文章中句号等符号进行分段，不会消耗模型token。父分段会拆分为若干子分段，子块用于检索，父块用作上下文
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
                  @change="handleEdit"
                  v-model:value="formState.normal_chunk_default_separators_no"
                  mode="tags"
                  style="width: 100%"
                  placeholder="请选择分段标识符"
                  :options="segmentationTags"
                ></a-select>
              </a-form-item>
              <a-form-item label="分段最大长度" required>
                <a-flex :gap="8" align="center">
                  <a-input-number
                    @blur="handleEdit"
                    v-model:value="formState.normal_chunk_default_chunk_size"
                    style="width: 220px"
                    :precision="0"
                    :min="1"
                    :max="10000"
                  />
                  字符
                </a-flex>
              </a-form-item>
              <a-form-item label="分段重叠长度">
                <a-flex :gap="8" align="center">
                  <a-input-number
                    @blur="handleEdit"
                    v-model:value="formState.normal_chunk_default_chunk_overlap"
                    style="width: 220px"
                    :precision="0"
                    :min="1"
                    :max="5000"
                  />
                  字符
                </a-flex>
              </a-form-item>
              <a-form-item>
                <template #label>
                  自动合并较小分段
                  <a-tooltip
                    title="开启后，如果分段长度不足设置的最大分段长度，会尝试与下一分段合并，直至合并后的分段字符数大于分段最大长度"
                  >
                    <QuestionCircleOutlined />
                  </a-tooltip>
                </template>
                <a-switch
                  @change="handleEdit"
                  checkedValue="false"
                  unCheckedValue="true"
                  v-model:checked="formState.normal_chunk_default_not_merged_text"
                  checked-children="开"
                  un-checked-children="关"
                />
              </a-form-item>
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
                  @blur="handleEdit"
                  v-model:value="formState.semantic_chunk_default_threshold"
                  style="width: 100%"
                  placeholder="请输入分段阈值"
                  :precision="0"
                  :min="0"
                  :max="100"
                />
              </a-form-item>
              <a-form-item label="分段最大长度" required>
                <a-flex :gap="8" align="center">
                  <a-input-number
                    @blur="handleEdit"
                    v-model:value="formState.semantic_chunk_default_chunk_size"
                    style="width: 220px"
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
                    @blur="handleEdit"
                    v-model:value="formState.semantic_chunk_default_chunk_overlap"
                    style="width: 220px"
                    :precision="0"
                    :min="1"
                    :max="100000"
                  />
                  字符
                </a-flex>
              </a-form-item>
            </template>

            <template v-if="formState.chunk_type == 3">
              <a-form-item required v-if="formState.chunk_type == 3">
                <template #label> AI大模型 </template>
                <ModelSelect
                  modelType="LLM"
                  placeholder="请选择AI大模型"
                  v-model:modeName="formState.ai_chunk_model"
                  v-model:modeId="formState.ai_chunk_model_config_id"
                  :modeName="formState.ai_chunk_model"
                  :modeId="formState.ai_chunk_model_config_id"
                  style="width: 300px"
                  @change="onChangeModel"
                  @loaded="onVectorModelLoaded"
                />
              </a-form-item>
              <a-form-item label="提示词设置" required>
                <a-flex :gap="8" align="center">
                  <a-textarea
                    @blur="handleEdit"
                    :maxLength="500"
                    style="height: 80px"
                    v-model:value="formState.ai_chunk_prumpt"
                    :placeholder="defaultAiChunkPrumpt"
                  />
                </a-flex>
              </a-form-item>
              <a-form-item>
                <template #label>
                  单次最大字符数
                  <a-tooltip>
                    <template #title
                      >由于大模型支持的上下文数量有限制，如果上传的文档较大，会按照最大字符数先拆分成多个分段，再提交给大模型进行分段。</template
                    >
                    <QuestionCircleOutlined style="cursor: pointer; margin-left: 2px" />
                  </a-tooltip>
                </template>
                <a-input-number
                  @blur="handleEdit"
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

            <template v-if="formState.chunk_type == 4 && !isQaLibrary">
              <div class="main-title-block">父块（用作上下文）</div>
              <a-form-item label="分段类型">
                <a-radio-group @change="handleEdit" v-model:value="formState.father_chunk_paragraph_type">
                  <a-radio :value="1"
                    >全文
                    <a-tooltip
                      title="整个文档用作父块并直接检索。请注意，出于性能原因，超过 10000 个标记的文本将被自动截断。"
                    >
                      <QuestionCircleOutlined />
                    </a-tooltip>
                  </a-radio>
                  <a-radio :value="2"
                    >段落
                    <a-tooltip
                      title="此模式根据分隔符和最大块长度将文本拆分为段落，使用拆分文本作为检索的父块"
                    >
                      <QuestionCircleOutlined />
                    </a-tooltip>
                  </a-radio>
                </a-radio-group>
              </a-form-item>

              <a-form-item label="分段标识符" v-if="formState.father_chunk_paragraph_type == 2">
                <a-select
                  placeholder="请选择"
                  @change="handleEdit"
                  style="width: 100%"
                  mode="tags"
                  :options="segmentationTags"
                  v-model:value="formState.father_chunk_separators_no"
                >
                </a-select>
              </a-form-item>
              <a-form-item label="分段最大长度" v-if="formState.father_chunk_paragraph_type == 2">
                <a-flex align="center" :gap="8">
                  <a-input-number
                    style="flex: 1"
                    @blur="handleEdit"
                    v-model:value="formState.father_chunk_chunk_size"
                    placeholder="分段最大长度"
                    :min="200"
                    :max="10000"
                    :precision="0"
                    :formatter="(value) => parseInt(value)"
                    :parser="(value) => parseInt(value)"
                  /><span class="unit-text">字符</span>
                </a-flex>
              </a-form-item>
              <div class="main-title-block">子块（用于检索）</div>

              <a-form-item label="分段标识符">
                <a-select
                  placeholder="请选择"
                  @change="handleEdit"
                  style="width: 100%"
                  mode="tags"
                  :options="segmentationTags"
                  v-model:value="formState.son_chunk_separators_no"
                >
                </a-select>
              </a-form-item>
              <a-form-item label="分段最大长度">
                <a-flex align="center" :gap="8">
                  <a-input-number
                    style="flex: 1"
                    @blur="handleEdit"
                    v-model:value="formState.son_chunk_chunk_size"
                    placeholder="分段最大长度"
                    :min="200"
                    :max="10000"
                    :precision="0"
                    :formatter="(value) => parseInt(value)"
                    :parser="(value) => parseInt(value)"
                  /><span class="unit-text">字符</span>
                </a-flex>
              </a-form-item>
            </template>
          </template>
        </a-form>
      </div>
      <OpenGrapgModal @ok="handleOpenGrapgOk" ref="openGrapgModalRef" />
    </div>
  </cu-scroll>
</template>

<script setup>
import { reactive, ref, h, nextTick, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Form, message, Modal } from 'ant-design-vue'
import { QuestionCircleOutlined, CheckCircleFilled } from '@ant-design/icons-vue'
import { getLibraryInfo, editLibrary, getSeparatorsList } from '@/api/library'
import { LIBRARY_OPEN_AVATAR } from '@/constants/index'
import AvatarInput from '@/views/library/add-library/components/avatar-input.vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import OpenGrapgModal from './components/open-grapg-modal.vue'
import CustomSelector from '@/components/custom-selector/index.vue'
import { useCompanyStore } from '@/stores/modules/company'
import { formatSeparatorsNo } from '@/utils/index'

const companyStore = useCompanyStore()
const neo4j_status = computed(() => {
  return companyStore.companyInfo?.neo4j_status == 'true'
})

const rotue = useRoute()
const router = useRouter()
const query = rotue.query
const defaultAvatar = LIBRARY_OPEN_AVATAR
const defaultAiChunkPrumpt =
  '你是一位文章分段助手，根据文章内容的语义进行合理分段，确保每个分段表述一个完整的语义，每个分段字数控制在500字左右，最大不超过1000字。请严格按照文章内容进行分段，不要对文章内容进行加工，分段完成后输出分段后的内容。'

const formState = reactive({
  library_name: '',
  library_intro: '',
  use_model: '',
  use_model_icon: '', // 新增图标字段
  use_model_name: '', // 新增系统名称
  is_offline: '',
  model_config_id: '',
  avatar: defaultAvatar,
  avatar_file: '',
  graph_switch: false,
  graph_model_config_id: void 0,
  graph_use_model: '',
  chunk_type: 1,
  normal_chunk_default_separators_no: [12, 11],
  normal_chunk_default_chunk_size: 512,
  normal_chunk_default_chunk_overlap: 50,
  semantic_chunk_default_chunk_size: 512,
  semantic_chunk_default_chunk_overlap: 50,
  semantic_chunk_default_threshold: 90,
  normal_chunk_default_not_merged_text: 'false',
  ai_chunk_size: 5000, // ai大模型分段最大字符数
  ai_chunk_model: '', // ai大模型分段模型名称
  ai_chunk_model_config_id: '', // ai大模型分段模型配置id
  ai_chunk_prumpt: defaultAiChunkPrumpt, // ai大模型分段提示词设置
  qa_index_type: 1,
  group_id: 0,
  father_chunk_paragraph_type: 2,
  father_chunk_separators_no: [],
  father_chunk_chunk_size: 1024,
  son_chunk_separators_no: [],
  son_chunk_chunk_size: 512
})
const currentModelDefine = ref('')
const isActive = ref(0)

const libraryInfo = ref({})

// 处理选择事件
const handleModelChange = (item) => {
  formState.use_model =
    modelDefine.includes(item.rawData.model_define) && item.rawData.deployment_name
      ? item.rawData.deployment_name
      : item.rawData.name
  formState.use_model_icon = item.icon
  formState.use_model_name = item.use_model_name
  formState.model_config_id = item.rawData.id
  currentModelDefine.value = item.rawData.model_define
  handleEdit()
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

const isQaLibrary = ref(true)
const getInfo = () => {
  getLibraryInfo({ id: query.id }).then((res) => {
    libraryInfo.value = res.data
    isActive.value = libraryInfo.value.is_offline ? 2 : 1
    formState.library_name = res.data.library_name
    formState.qa_index_type = res.data.qa_index_type
    formState.library_intro = res.data.library_intro
    formState.use_model = res.data.use_model
    formState.is_offline = res.data.is_offline
    formState.group_id = res.data.group_id

    formState.model_config_id = res.data.model_config_id
    formState.avatar = res.data.avatar ? res.data.avatar : defaultAvatar
    formState.avatar_file = res.data.avatar_file ? res.data.avatar_file : ''

    formState.graph_switch = res.data.graph_switch != '0'
    formState.graph_model_config_id =
      res.data.graph_model_config_id > 0 ? res.data.graph_model_config_id : void 0
    formState.graph_use_model = res.data.graph_use_model

    formState.chunk_type = +res.data.chunk_type
    formState.normal_chunk_default_separators_no = formatSeparatorsNo( res.data.normal_chunk_default_separators_no, [12, 11])
    formState.normal_chunk_default_chunk_size = res.data.normal_chunk_default_chunk_size
    formState.normal_chunk_default_not_merged_text = res.data.normal_chunk_default_not_merged_text
    formState.normal_chunk_default_chunk_overlap = res.data.normal_chunk_default_chunk_overlap
    formState.semantic_chunk_default_chunk_size = res.data.semantic_chunk_default_chunk_size
    formState.semantic_chunk_default_chunk_overlap = res.data.semantic_chunk_default_chunk_overlap
    formState.semantic_chunk_default_threshold = res.data.semantic_chunk_default_threshold
    formState.ai_chunk_size = res.data.ai_chunk_size || 5000
    formState.ai_chunk_model = res.data.ai_chunk_model
    formState.ai_chunk_model_config_id = res.data.ai_chunk_model_config_id
    formState.ai_chunk_prumpt = res.data.ai_chunk_prumpt || defaultAiChunkPrumpt
    formState.father_chunk_paragraph_type = +res.data.father_chunk_paragraph_type || 2
    formState.father_chunk_separators_no = formatSeparatorsNo(res.data.father_chunk_separators_no, [12, 11])
    formState.father_chunk_chunk_size = +res.data.father_chunk_chunk_size || 1024
    formState.son_chunk_separators_no = formatSeparatorsNo(res.data.son_chunk_separators_no, [8, 10])
    formState.son_chunk_chunk_size = +res.data.son_chunk_chunk_size || 512
    isQaLibrary.value = res.data.type == 2
  })
}
getInfo()
const libraryModeList = ref([
  {
    iconName: 'high-active',
    title: '高质量',
    value: 1,
    is_offline: false,
    desc: '使用在线的嵌入模型，在召回时具有更高的准确度，但需要花费token'
  }
  // {
  //   iconName: 'economic',
  //   title: '经济',
  //   value: 2,
  //   is_offline: true,
  //   desc: '使用离线的向量模型，较在线模型准确度稍低，但是不需要消耗token'
  // }
])

const useForm = Form.useForm

const rules = reactive({
  library_name: [{ required: true, message: '请输入库名称', trigger: 'blur' }],
  use_model: [{ required: true, message: '请选择嵌入模型', trigger: 'change' }]
})

const { validateInfos } = useForm(formState, rules)

const handleChangeModel = (val, option) => {
  const self = option.current_obj
  formState.use_model =
    modelDefine.indexOf(self.model_define) > -1 && self.deployment_name
      ? self.deployment_name
      : self.name
  currentModelDefine.value = self.model_define
  formState.model_config_id = self.id || option.model_config_id
  handleEdit()
}

const onAvatarChange = (data) => {
  formState.avatar = data.imageUrl
  formState.avatar_file = data.file
  handleEdit()
}

const handleSelectLibrary = () => {
  return false
}

const modelDefine = ['azure', 'ollama', 'xinference', 'openaiAgent']
const oldModelDefineList = ['azure']

const onChangeModel = () => {
  handleEdit()
}
const vectorModelList = ref([])
const onVectorModelLoaded = (list) => {
  vectorModelList.value = list

  nextTick(() => {
    if (!formState.ai_chunk_model || !Number(formState.ai_chunk_model_config_id)) {
      setDefaultModel()
    }
  })
  // handleEdit()
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

const openGrapgModalRef = ref(null)
const handleGraphSwitch = (val) => {
  if (val) {
    formState.graph_switch = false
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
    handleEdit()
  }
}

const handleChangeSegmentationType = (type) => {
  formState.chunk_type = type
  handleEdit()
}

const handleOpenGrapgOk = (data) => {
  if (data.graph_model_config_id) {
    formState.graph_switch = true
    formState.graph_model_config_id = data.graph_model_config_id
    formState.graph_use_model = data.graph_use_model
    handleEdit(() => {
      Modal.confirm({
        title: '已开启知识图谱',
        content: '您可以在文档列表中点击知识图谱学习，系統将在您手动操作后开始抽取知识图谱',
        cancelText: '知道了',
        okText: '去学习',
        icon: h(CheckCircleFilled, { style: { color: '#52c41a' } }),
        onOk: () =>
          router.push({
            path: '/library/details/knowledge-document',
            query: { id: query.id }
          })
      })
    })
  }
}

const handleEdit = (callback = null) => {
  if (!formState.library_name) {
    return message.error('请输入知识库名称')
  }
  let data = {
    library_name: formState.library_name,
    qa_index_type: formState.qa_index_type,
    library_intro: formState.library_intro,
    use_model: formState.use_model,
    model_config_id: formState.model_config_id,
    is_offline: formState.is_offline,
    graph_switch: formState.graph_switch ? 1 : 0,
    graph_model_config_id: formState.graph_model_config_id,
    graph_use_model: formState.graph_use_model,
    chunk_type: formState.chunk_type,
    normal_chunk_default_separators_no: JSON.stringify(formState.normal_chunk_default_separators_no),
    normal_chunk_default_chunk_size: formState.normal_chunk_default_chunk_size,
    normal_chunk_default_chunk_overlap: formState.normal_chunk_default_chunk_overlap,
    normal_chunk_default_not_merged_text: formState.normal_chunk_default_not_merged_text,
    semantic_chunk_default_chunk_size: formState.semantic_chunk_default_chunk_size,
    semantic_chunk_default_chunk_overlap: formState.semantic_chunk_default_chunk_overlap,
    semantic_chunk_default_threshold: formState.semantic_chunk_default_threshold,
    ai_chunk_size: formState.ai_chunk_size,
    ai_chunk_model: formState.ai_chunk_model,
    ai_chunk_model_config_id: formState.ai_chunk_model_config_id,
    ai_chunk_prumpt: formState.ai_chunk_prumpt,
    father_chunk_paragraph_type: formState.father_chunk_paragraph_type,
    father_chunk_separators_no: JSON.stringify(formState.father_chunk_separators_no),
    father_chunk_chunk_size: formState.father_chunk_chunk_size,
    son_chunk_separators_no: JSON.stringify(formState.son_chunk_separators_no),
    son_chunk_chunk_size: formState.son_chunk_chunk_size,
    group_id: formState.group_id,
    id: rotue.query.id
  }
  if (oldModelDefineList.indexOf(currentModelDefine.value) > -1) {
    // 传给后端的是默认，渲染的是真实名称
    data.use_model = '默认'
  }
  if (formState.avatar_file) {
    data.avatar = formState.avatar_file
  }
  editLibrary(data).then((res) => {
    typeof callback === 'function' ? callback() : message.success('修改成功')
  })
}

const handleChangeQaIndexType = (type) => {
  if (type == formState.qa_index_type) {
    return
  }
  formState.qa_index_type = type
  handleEdit()
}
</script>
