<style lang="less" scoped>
.user-model-page {
  width: 100%;
  height: 100%;
  border-bottom: 1px solid #fff;
  border-right: 1px solid #fff;
  background-color: #f2f4f7;

  .list-wrapper {
    height: calc(100% - 47px);
    overflow-x: hidden;
    overflow-y: auto;
  }
}
.alert-content-box{
  padding: 16px 24px 0 24px;
  .ant-alert{
    align-items: baseline;
  }
  .title{
    font-size: 15px;
    font-weight: 600;
    color: #333;
  }
}
</style>

<template>
  <div class="user-model-page">
    <PageTabs v-model:value="activeTab" @change="onChangeTab" />
    <div class="list-wrapper">
      <div class="alert-content-box" v-if="activeTab == 0">
        <a-alert show-icon>
          <template #message>
            <div class="title">使用说明</div>
            <div>1、模型主要分为嵌入模型（TEXT EMBEDDING）和大模型（LLM），嵌入模型用于将上传知识库的文本内容进行向量化便于后续检索，大模型用于将检索到的相关知识进行理解后返回答案。</div>
            <div>2、配置模型时尽量选择同时支持LLM和TEXT EMBEDDING的模型，否则要配置多个模型。</div>
            <div>3、国内模型厂商一般都需要认证后才可使用，所以在配置模型前，确保在模型服务厂商那里已认证，否则配置模型时会提示配置参数错误。</div>
          </template>
        </a-alert>
      </div>

      <ModelList
        :list="addedModelList"
        :type="1"
        key="addedModelList"
        v-if="activeTab == 1"
        @edit="handleEditModel"
        @new="handleNewModel"
        @see="handleSeeModel"
        @remove="handleRemoveModel"
      />
      <ModelList
        :list="canAddModelList"
        :type="2"
        key="canAddModelList"
        @add="handleAddModel"
        @see="handleSeeModel"
        v-if="activeTab == 0"
      />
    </div>
  </div>
  <!-- 查看模型 -->
  <SeeModelAlert
    ref="seeModelAlertRef"
    @new="handleNewModel"
    @edit="handleEditModel"
    @remove="handleRemoveModel"
  />
  <!-- 设置模型 -->
  <SetModelAlert ref="setModelAlertRef" @ok="saveModelSuccess" />
</template>

<script setup>
import { ref, computed, createVNode } from 'vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { message, Modal } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { getModelConfigList, delModelConfig } from '@/api/model/index'
import PageTabs from './components/page-tabs.vue'
import ModelList from './components/model-list.vue'
import SetModelAlert from './components/set-model-alert.vue'
import SeeModelAlert from './components/see-model-alert.vue'

const { t } = useI18n()
const activeTab = ref(1)

const onChangeTab = () => {}
// 获取模型列表
const modelList = ref([])

const addedModelList = computed(() => {
  return modelList.value.filter((item) => item.listType == 1)
})

const canAddModelList = computed(() => {
  return modelList.value
})

const getModelList = () => {
  getModelConfigList().then((res) => {
    let list = res.data || []

    list.forEach((item) => {
      item.listType = item.config_list.length > 0 ? 1 : 0
    })

    // 添加成功后移动到底部
    list.sort((a, b) => a.listType - b.listType)
    modelList.value = list
  })
}

getModelList()

// 查看模型
const seeModelAlertRef = ref(null)
const handleSeeModel = (model) => {
  seeModelAlertRef.value.open(model)
}

// 新增模型
const handleNewModel = (model) => {
  setModelAlertRef.value.open(model)
}

// 添加模型
const handleAddModel = (model) => {
  setModelAlertRef.value.open(model)
}

// 修改模型配置
const setModelAlertRef = ref(null)
const handleEditModel = (model, record) => {
  setModelAlertRef.value.open(model, record)
}

const saveModelSuccess = () => {
  getModelList()
}

// 删除模型
const handleRemoveModel = (data) => {
  Modal.confirm({
    title: t('views.user.model.deleteModel'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('views.user.model.deleteModelContent'),
    onOk() {
      delModelConfig({ id: data.id }).then(() => {
        message.success(t('common.deleteSuccessful'))
        getModelList()
      })
    }
  })
}
</script>
