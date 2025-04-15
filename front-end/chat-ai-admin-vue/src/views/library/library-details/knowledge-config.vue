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
  align-items: center;
  gap: 16px;
  width: 550px;
  .select-card-item {
    width: calc(50% - 8px);
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
      min-height: 44px;
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
  <cu-scroll>
    <div class="add-library-page">
      <div class="form-box">
        <a-form :label-col="{ span: 5 }">
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
            <a-select
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
            </a-select>
          </a-form-item>
          <a-form-item label="生成知识图谱">
            <a-switch
              @change="handleGraphSwitch"
              :checked="formState.graph_switch"
              checked-children="开"
              un-checked-children="关"
            />
            <div class="form-item-tip">开启后，可以文档列表手动点击知识图谱学习生成知识图谱</div>
          </a-form-item>
          <a-form-item label="知识图谱模型" v-show="formState.graph_switch">
            <ModelSelect
              modelType="LLM"
              v-model:modeName="formState.graph_use_model"
              v-model:modeId="formState.graph_model_config_id"
              style="width: 300px"
              @change="onChangeModel"
              @loaded="onVectorModelLoaded"
            />
          </a-form-item>
          <template v-if="!isQaLibrary">
            <a-form-item label="分段方式" required>
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
              </div>
              <a-alert
                v-if="formState.chunk_type == 2"
                style="margin-top: 12px; width: 650px"
                message="提示：语义分段更适合没有排版过的文章，即没有明显换行符号的文本，否则更推荐使用普通分段"
              ></a-alert>
            </a-form-item>
            <template v-if="formState.chunk_type == 1">
              <a-form-item label="分段标识符" required>
                <a-select
                  @change="handleEdit"
                  v-model:value="formState.normal_chunk_default_separators_no"
                  mode="multiple"
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
                    :max="100000"
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
                    :max="100000"
                  />
                  字符
                </a-flex>
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
          </template>
        </a-form>
      </div>
      <OpenGrapgModal @ok="handleOpenGrapgOk" ref="openGrapgModalRef" />
    </div>
  </cu-scroll>
</template>

<script setup>
import { reactive, ref, h} from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Form, message, Modal } from 'ant-design-vue'
import { QuestionCircleOutlined, CheckCircleFilled } from '@ant-design/icons-vue'
import { getLibraryInfo, editLibrary, getSeparatorsList } from '@/api/library'
import { getModelConfigOption } from '@/api/model/index'
import { duplicateRemoval, removeRepeat } from '@/utils/index'
import { LIBRARY_OPEN_AVATAR } from '@/constants/index'
import AvatarInput from '@/views/library/add-library/components/avatar-input.vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import OpenGrapgModal from './components/open-grapg-modal.vue'

const rotue = useRoute()
const router = useRouter()
const query = rotue.query
const defaultAvatar = LIBRARY_OPEN_AVATAR
const formState = reactive({
  library_name: '',
  library_intro: '',
  use_model: '',
  is_offline: '',
  model_config_id: '',
  avatar: defaultAvatar,
  avatar_file: '',
  graph_switch: false,
  graph_model_config_id: void 0,
  graph_use_model: '',
  chunk_type: 1,
  normal_chunk_default_separators_no: [10, 12],
  normal_chunk_default_chunk_size: 512,
  normal_chunk_default_chunk_overlap: 50,
  semantic_chunk_default_chunk_size: 512,
  semantic_chunk_default_chunk_overlap: 50,
  semantic_chunk_default_threshold: 90
})
const currentModelDefine = ref('')
const isActive = ref(0)

const libraryInfo = ref({})

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
    formState.library_intro = res.data.library_intro
    formState.use_model = res.data.use_model
    formState.is_offline = res.data.is_offline

    formState.model_config_id = res.data.model_config_id
    formState.avatar = res.data.avatar ? res.data.avatar : defaultAvatar
    formState.avatar_file = res.data.avatar_file ? res.data.avatar_file : ''

    formState.graph_switch = res.data.graph_switch != '0'
    formState.graph_model_config_id =
      res.data.graph_model_config_id > 0 ? res.data.graph_model_config_id : void 0
    formState.graph_use_model = res.data.graph_use_model

    formState.chunk_type = +res.data.chunk_type
    formState.normal_chunk_default_separators_no = res.data.normal_chunk_default_separators_no
      ? res.data.normal_chunk_default_separators_no.split(',').map((item) => +item)
      : []
    formState.normal_chunk_default_chunk_size = res.data.normal_chunk_default_chunk_size
    formState.normal_chunk_default_chunk_overlap = res.data.normal_chunk_default_chunk_overlap
    formState.semantic_chunk_default_chunk_size = res.data.semantic_chunk_default_chunk_size
    formState.semantic_chunk_default_chunk_overlap = res.data.semantic_chunk_default_chunk_overlap
    formState.semantic_chunk_default_threshold = res.data.semantic_chunk_default_threshold
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

// 获取嵌入模型列表
const modelList = ref([])
const modelDefine = ['azure', 'ollama', 'xinference', 'openaiAgent']
const oldModelDefineList = ['azure']

function uniqueArr(arr, arr1, key) {
  const keyVals = new Set(arr.map((item) => item.model_define))
  arr1.filter((obj) => {
    let val = obj[key]
    if (keyVals.has(val)) {
      arr.filter((obj1) => {
        if (obj1.model_define == val) {
          obj1.children = removeRepeat(obj1.children, obj.children)
          return false
        }
      })
    }
  })
  return arr
}

const getModelList = () => {
  getModelConfigOption({
    model_type: 'TEXT EMBEDDING',
    is_offline: false
  }).then((res) => {
    let list = res.data || []
    let children = []

    modelList.value = list.map((item) => {
      children = []
      for (let i = 0; i < item.model_info.vector_model_list.length; i++) {
        const ele = item.model_info.vector_model_list[i]
        children.push({
          name: ele,
          deployment_name: item.model_config.deployment_name,
          id: item.model_config.id,
          model_define: item.model_info.model_define
        })
      }

      return {
        id: item.model_config.id,
        name: item.model_info.model_name,
        model_define: item.model_info.model_define,
        icon: item.model_info.model_icon_url,
        children: children,
        deployment_name: item.model_config.deployment_name
      }
    })

    // 如果modelList存在两个相同model_define情况就合并到一个对象的children中去
    modelList.value = uniqueArr(
      duplicateRemoval(modelList.value, 'model_define'),
      modelList.value,
      'model_define'
    )
  })
}

getModelList()

const onChangeModel = () => {
  handleEdit()
}
const vectorModelList = ref([])
const onVectorModelLoaded = (list) => {
  vectorModelList.value = list
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
        icon: h(CheckCircleFilled, {style: {color: '#52c41a'}}),
        onOk: () => router.push({
          path: '/library/details/knowledge-document',
          query: {id: query.id}
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
    library_intro: formState.library_intro,
    use_model: formState.use_model,
    model_config_id: formState.model_config_id,
    is_offline: formState.is_offline,
    graph_switch: formState.graph_switch ? 1 : 0,
    graph_model_config_id: formState.graph_model_config_id,
    graph_use_model: formState.graph_use_model,
    chunk_type: formState.chunk_type,
    normal_chunk_default_separators_no: formState.normal_chunk_default_separators_no.join(','),
    normal_chunk_default_chunk_size: formState.normal_chunk_default_chunk_size,
    normal_chunk_default_chunk_overlap: formState.normal_chunk_default_chunk_overlap,
    semantic_chunk_default_chunk_size: formState.semantic_chunk_default_chunk_size,
    semantic_chunk_default_chunk_overlap: formState.semantic_chunk_default_chunk_overlap,
    semantic_chunk_default_threshold: formState.semantic_chunk_default_threshold,
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
</script>
