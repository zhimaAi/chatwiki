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

.use-model-item{
  position: relative;
  width: 226px;
  height: 124px;
  border-radius: 2px;
  border: 2px solid #D9D9D9;
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
  color: #BFBFBF;
  font-weight: bolder;
}

.use-model-item.active {
  border: 2px solid #BFBFBF;

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
    color: #BFBFBF;
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
          <a-textarea  @blur="handleEdit" v-model:value="formState.library_intro" placeholder="请输入知识库介绍" />
        </a-form-item>

        <a-form-item label="嵌入模型" v-bind="validateInfos.use_model">
          <div class="card-box">
            <div class="use-model-item"
              :class="{ active: isActive == item.value }"
              v-for="item in libraryModeList"
              :key="item.value"
              @click="handleSelectLibrary(item)"
            >
              <div
                class="use-model-item-top"
                :class="{ active: isActive == item.value }"
              >
                <svg-icon class="use-model-item-top-icon" style="color: red;" :name="item.iconName"></svg-icon>
                <p>{{item.title}}</p>
              </div>
              <p>{{item.desc}}</p>
              <svg-icon
                class="check-arrow"
                name="select-disabled"
                style="font-size: 24px;color: #fff;"
                v-if="isActive == item.value"
              ></svg-icon>
            </div>
          </div>
          <a-select
            disabled
            v-model:value="formState.use_model"
            placeholder="请选择嵌入模型"
            @change="handleChangeModel"
          >
            <a-select-opt-group v-for="item in modelList" :key="item.id">
              <template #label>
                <span><img class="model-icon" :src="item.icon" alt="" /></span>
              </template>
              <a-select-option
                :value="val"
                :model_config_id="item.id"
                v-for="val in item.children"
                :key="val"
              >
                <span>{{ val }}</span>
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
import { useRoute} from 'vue-router'
import { Form, message } from 'ant-design-vue'

import { getLibraryInfo, editLibrary } from '@/api/library'

const rotue = useRoute()
const query = rotue.query
const formState = reactive({
  library_name: '',
  library_intro: '',
  use_model: '',
  is_offline: '',
})

const isActive = ref(0)

const libraryInfo = ref({})
const getInfo = () => {
  getLibraryInfo({ id: query.id }).then((res) => {
    libraryInfo.value = res.data
    isActive.value = libraryInfo.value.is_offline ? 2 : 1
    formState.library_name = res.data.library_name;
    formState.library_intro = res.data.library_intro;
    formState.use_model = res.data.use_model;
    formState.is_offline = res.data.is_offline;
  })
}
getInfo()
const libraryModeList = ref([
  {
    iconName: 'high',
    title: '高质量',
    value: 1,
    is_offline: false,
    desc: '使用在线的嵌入模型，在召回时具有更高的准确度，但需要花费token'
  },
  {
    iconName: 'economic',
    title: '经济',
    value: 2,
    is_offline: true,
    desc: '使用离线的向量模型，较在线模型准确度稍低，但是不需要消耗token'
  }
])

const useForm = Form.useForm


const rules = reactive({
  library_name: [{ required: true, message: '请输入库名称', trigger: 'blur' }],
  use_model: [{ required: true, message: '请选择嵌入模型', trigger: 'change' }],
})

const { validateInfos } = useForm(formState, rules)

const handleChangeModel = () => {
  return false
}

const handleSelectLibrary = () => {
  return false
}

// 获取嵌入模型列表
const modelList = ref([])

const handleEdit = () => {
  if (!formState.library_name) {
    return message.error('请输入知识库名称')
  }
  let data = {
    library_name: formState.library_name,
    library_intro: formState.library_intro,
    id: rotue.query.id
  }
  editLibrary(data).then((res) => {
    message.success('修改成功')
  })
}

</script>
