<style lang="less" scoped>
.add-library-page {
  padding: 24px;

  .form-box {
    width: 554px;
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
</style>

<template>
  <div class="add-library-page">
    <div class="form-box">
      <a-form :label-col="{ span: 4 }">
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
      </a-form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Form, message } from 'ant-design-vue'

import { getLibraryInfo, editLibrary } from '@/api/library'
import { getModelConfigOption } from '@/api/model/index'
import { duplicateRemoval, removeRepeat } from '@/utils/index'
import { DEFAULT_LIBRARY_AVATAR2 } from '@/constants/index'
import AvatarInput from '@/views/library/add-library/components/avatar-input.vue'

const rotue = useRoute()
const query = rotue.query
const defaultAvatar = DEFAULT_LIBRARY_AVATAR2
const formState = reactive({
  library_name: '',
  library_intro: '',
  use_model: '',
  is_offline: '',
  model_config_id: '',
  avatar: defaultAvatar,
  avatar_file: ''
})
const currentModelDefine = ref('')
const isActive = ref(0)

const libraryInfo = ref({})
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

const handleEdit = () => {
  if (!formState.library_name) {
    return message.error('请输入知识库名称')
  }
  let data = {
    library_name: formState.library_name,
    library_intro: formState.library_intro,
    use_model: formState.use_model,
    model_config_id: formState.model_config_id,
    is_offline: formState.is_offline,
    id: rotue.query.id
  }
  if (oldModelDefineList.indexOf(currentModelDefine.value) > -1) {
    // 传给后端的是默认，渲染的是真实名称
    data.use_model = '默认'
  }
  if(formState.avatar_file){
    data.avatar = formState.avatar_file
  }
  editLibrary(data).then((res) => {
    message.success('修改成功')
  })
}
</script>
