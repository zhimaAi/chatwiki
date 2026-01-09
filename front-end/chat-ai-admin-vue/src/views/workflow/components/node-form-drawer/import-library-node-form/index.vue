<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        desc="将相关内容导入到对应知识库"
        @close="handleClose"
      >
      </NodeFormHeader>
    </template>
    <div class="problem-optimization-form">
      <div class="node-form-content">
        <a-form ref="formRef" layout="vertical" :model="formState">
          <div class="gray-block">
            <div class="gray-block-title">
              <svg-icon name="nav-library" style="font-size: 14px; color: #333"></svg-icon>
              知识库
            </div>
            <a-form-item :label="null">
              <div class="konwledge-list-item" v-for="item in selectedLibraryRows" :key="item.id">
                <div class="avatar-box">
                  <img :src="item.avatar" alt="" />
                </div>
                <div class="content-box">
                  <div class="list-name">{{ item.library_name }}</div>
                  <div class="list-intro">{{ item.library_intro }}</div>
                </div>
                <div class="btn-hover-wrap" @click="handleDelKonwledge(item)">
                  <CloseCircleOutlined />
                </div>
              </div>
              <div class="btn-block">
                <div>
                  <a-button
                    @click="handleOpenSelectLibraryAlert"
                    :icon="h(PlusOutlined)"
                    block
                    type="dashed"
                    >添加知识库</a-button
                  >
                </div>
              </div>
            </a-form-item>

            <a-form-item label="知识分组">
              <a-select
                v-model:value="formState.library_group_id"
                style="width: 100%"
                placeholder="请选择"
              >
                <a-select-option v-for="item in groupLists" :value="item.id" :key="item.id">{{
                  item.group_name
                }}</a-select-option>
              </a-select>
            </a-form-item>

            <a-form-item
              label="导入类型"
              v-if="selectedLibraryType != 2 && selectedLibraryType >= 0"
            >
              <a-segmented
                class="customer-segmented"
                v-model:value="formState.import_type"
                :options="options1"
              />
            </a-form-item>

            <a-form-item label="分段问题重复" v-if="selectedLibraryType == 2">
              <a-segmented
                class="customer-segmented"
                v-model:value="formState.qa_repeat_op"
                :options="options2"
              />
            </a-form-item>

            <a-form-item
              label="文档URL重复"
              v-if="
                selectedLibraryType != 2 &&
                selectedLibraryType >= 0 &&
                formState.import_type == 'url'
              "
            >
              <a-segmented
                class="customer-segmented"
                v-model:value="formState.normal_url_repeat_op"
                :options="options2"
              />
            </a-form-item>
          </div>

          <div class="gray-block mt16">
            <div class="gray-block-title">
              <svg-icon name="input" style="font-size: 14px; color: #333"></svg-icon>
              输入
            </div>
            <template v-if="selectedLibraryType != 2 && selectedLibraryType >= 0">
              <div class="flex-block-item" v-if="formState.import_type == 'content'">
                <div class="form-item-label required">文档标题</div>
                <div class="form-content-box flex1">
                  <at-input
                    inputStyle="overflow-y: hidden; overflow-x: scroll; height: 22px;"
                    :options="variableOptions"
                    :defaultValue="formState.normal_title"
                    @open="getVlaueVariableList"
                    @change="
                      (text, selectedList) => changeValue(text, selectedList, 'normal_title')
                    "
                    placeholder="请输入变量值，键入“/”插入变量"
                  >
                    <template #option="{ label, payload }">
                      <div class="field-list-item">
                        <div class="field-label">{{ label }}</div>
                        <div class="field-type">{{ payload.typ }}</div>
                      </div>
                    </template>
                  </at-input>
                </div>
              </div>
              <div class="flex-block-item" v-if="formState.import_type == 'content'">
                <div class="form-item-label required">文档内容</div>
                <div class="form-content-box flex1">
                  <at-input
                    inputStyle="overflow-y: hidden; overflow-x: scroll; height: 22px;"
                    :options="variableOptions"
                    :defaultValue="formState.normal_content"
                    @open="getVlaueVariableList"
                    @change="
                      (text, selectedList) => changeValue(text, selectedList, 'normal_content')
                    "
                    placeholder="请输入变量值，键入“/”插入变量"
                  >
                    <template #option="{ label, payload }">
                      <div class="field-list-item">
                        <div class="field-label">{{ label }}</div>
                        <div class="field-type">{{ payload.typ }}</div>
                      </div>
                    </template>
                  </at-input>
                </div>
              </div>
              <div class="flex-block-item" v-if="formState.import_type == 'url'">
                <div class="form-item-label">文档URL</div>
                <div class="form-content-box flex1">
                  <at-input
                    inputStyle="overflow-y: hidden; overflow-x: scroll; height: 22px;"
                    :options="variableOptions"
                    :defaultValue="formState.normal_url"
                    @open="getVlaueVariableList"
                    @change="(text, selectedList) => changeValue(text, selectedList, 'normal_url')"
                    placeholder="请输入变量值，键入“/”插入变量"
                  >
                    <template #option="{ label, payload }">
                      <div class="field-list-item">
                        <div class="field-label">{{ label }}</div>
                        <div class="field-type">{{ payload.typ }}</div>
                      </div>
                    </template>
                  </at-input>
                </div>
              </div>
            </template>

            <template v-if="selectedLibraryType == 2">
              <div class="flex-block-item">
                <div class="form-item-label required">分段问题</div>
                <div class="form-content-box flex1">
                  <at-input
                    inputStyle="overflow-y: hidden; overflow-x: scroll; height: 22px;"
                    :options="variableOptions"
                    :defaultValue="formState.qa_question"
                    @open="getVlaueVariableList"
                    @change="(text, selectedList) => changeValue(text, selectedList, 'qa_question')"
                    placeholder="请输入变量值，键入“/”插入变量"
                  >
                    <template #option="{ label, payload }">
                      <div class="field-list-item">
                        <div class="field-label">{{ label }}</div>
                        <div class="field-type">{{ payload.typ }}</div>
                      </div>
                    </template>
                  </at-input>
                </div>
              </div>
              <div class="flex-block-item">
                <div class="form-item-label required">分段答案</div>
                <div class="form-content-box flex1">
                  <at-input
                    inputStyle="overflow-y: hidden; overflow-x: scroll; height: 22px;"
                    :options="variableOptions"
                    :defaultValue="formState.qa_answer"
                    @open="getVlaueVariableList"
                    @change="(text, selectedList) => changeValue(text, selectedList, 'qa_answer')"
                    placeholder="请输入变量值，键入“/”插入变量"
                  >
                    <template #option="{ label, payload }">
                      <div class="field-list-item">
                        <div class="field-label">{{ label }}</div>
                        <div class="field-type">{{ payload.typ }}</div>
                      </div>
                    </template>
                  </at-input>
                </div>
              </div>
              <div class="flex-block-item">
                <div class="form-item-label">
                  答案附图
                  <a-tooltip title="图片url，<array string>格式，支持传入多个">
                    <QuestionCircleOutlined />
                  </a-tooltip>
                </div>
                <div class="form-content-box flex1">
                  <at-input
                    inputStyle="overflow-y: hidden; overflow-x: scroll; height: 22px;"
                    :options="variableOptionsArr"
                    :defaultValue="formState.qa_images_variable"
                    @open="getVlaueVariableList"
                    @change="
                      (text, selectedList) => changeValue(text, selectedList, 'qa_images_variable')
                    "
                    placeholder="请输入变量值，键入“/”插入变量"
                  >
                    <template #option="{ label, payload }">
                      <div class="field-list-item">
                        <div class="field-label">{{ label }}</div>
                        <div class="field-type">{{ payload.typ }}</div>
                      </div>
                    </template>
                  </at-input>
                </div>
              </div>
              <div class="flex-block-item">
                <div class="form-item-label">
                  相似问法
                  <a-tooltip title="<array string>格式，支持传入多个">
                    <QuestionCircleOutlined />
                  </a-tooltip>
                </div>
                <div class="form-content-box flex1">
                  <at-input
                    inputStyle="overflow-y: hidden; overflow-x: scroll; height: 22px;"
                    :options="variableOptionsArr"
                    :defaultValue="formState.qa_similar_question_variable"
                    @open="getVlaueVariableList"
                    @change="
                      (text, selectedList) =>
                        changeValue(text, selectedList, 'qa_similar_question_variable')
                    "
                    placeholder="请输入变量值，键入“/”插入变量"
                  >
                    <template #option="{ label, payload }">
                      <div class="field-list-item">
                        <div class="field-label">{{ label }}</div>
                        <div class="field-type">{{ payload.typ }}</div>
                      </div>
                    </template>
                  </at-input>
                </div>
              </div>
            </template>
          </div>
          <div class="gray-block mt16">
            <div class="gray-block-title">
              <svg-icon name="output" style="font-size: 14px; color: #333"></svg-icon>
              输出
            </div>
            <div class="options-item">
              <div class="option-label">msg</div>
              <div class="desc">ok表示成功，其他输出为报错日志</div>
            </div>
          </div>
        </a-form>
        <LibrarySelectAlert
          ref="librarySelectAlertRef"
          :showWxType="!!wxAppLibary"
          select_type="radio"
          @change="onChangeLibrarySelected"
        />
      </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import NodeFormLayout from '../node-form-layout.vue'
import NodeFormHeader from '../node-form-header.vue'
import { ref, reactive, watch, computed, onMounted, h, nextTick } from 'vue'
import { CloseCircleOutlined, PlusOutlined, QuestionCircleOutlined } from '@ant-design/icons-vue'
import { getLibraryList, getLibraryGroup } from '@/api/library/index'
import LibrarySelectAlert from '../components/library-select-alert.vue'
import { getSpecifyAbilityConfig } from '@/api/explore/index.js'
import AtInput from '../../at-input/at-input.vue'

const emit = defineEmits(['update-node'])
const props = defineProps({
  lf: {
    type: Object,
    default: null
  },
  nodeId: {
    type: String,
    default: ''
  },
  node: {
    type: Object,
    default: () => ({})
  }
})

const options1 = [
  {
    label: '导入文档内容',
    value: 'content'
  },
  {
    label: '导入url',
    value: 'url'
  }
]

const options2 = [
  {
    label: '依然导入',
    value: 'import'
  },
  {
    label: '不导入',
    value: 'not import'
  },
  {
    label: '更新内容',
    value: 'update'
  }
]

const groupListMap = ref({})

const variableOptions = ref([])
const variableOptionsArr = ref([])

const getVlaueVariableList = () => {
  const nodeModel = props.lf.getNodeModelById(props.nodeId)
  if (nodeModel) {
    let list = nodeModel.getAllParentVariable()
    list.forEach((item) => {
      item.tags = item.tags || []
    })

    variableOptions.value = list
    variableOptionsArr.value = filterCalueOptions(list)
  }
}

function filterCalueOptions(list) {
  function traverse(items) {
    const result = []

    for (const item of items) {
      if (item.children && item.children.length > 0) {
        // 有子节点的情况，递归处理子节点
        const filteredChildren = traverse(item.children)

        if (filteredChildren.length > 0) {
          // 如果过滤后的子节点不为空，保留当前节点并更新其子节点
          const newItem = { ...item, children: filteredChildren }
          result.push(newItem)
        }
      } else {
        // 叶子节点，检查类型是否符合要求
        if (item.typ && item.typ === 'array<string>') {
          result.push({ ...item })
        }
      }
    }

    return result
  }

  return traverse(list)
}

const wxAppLibary = ref(null)
const formState = reactive({
  library_group_id: '0',
  library_id: [],
  import_type: 'content',
  normal_url: '',
  normal_title: '',
  normal_content: '',
  normal_url_repeat_op: 'import',
  qa_question: '',
  qa_answer: '',
  qa_images_variable: '',
  qa_similar_question_variable: '',
  qa_repeat_op: 'import',
  outputs: [
    {
      key: 'msg',
      typ: 'string'
    }
  ]
})

const update = () => {
  const data = JSON.stringify({
    library_import: {
      ...formState,
      library_id: formState.library_id.join(',')
    }
  })

  emit('update-node', {
    ...props.node,
    node_params: data
  })
}

const init = () => {
  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'
    let library_import = JSON.parse(dataRaw).library_import || {}

    getVlaueVariableList()

    library_import = JSON.parse(JSON.stringify(library_import))

    for (let key in library_import) {
      if (key == 'library_id') {
        formState.library_id = library_import[key] ? library_import[key].split(',') : []
      } else {
        formState[key] = library_import[key]
      }
    }
    // 公众号知识库是否开启
    getSpecifyAbilityConfig({ ability_type: 'library_ability_official_account' }).then((res) => {
      let _data = res?.data || {}
      if (_data?.user_config?.switch_status == 1) {
        wxAppLibary.value = _data
      }
    })

    setTimeout(() => {
      getGroupList()
    }, 300)
  } catch (error) {
    console.log(error)
  }
}

const groupLists = computed(() => {
  let library_id = formState.library_id[0]
  let defaultList = [
    {
      group_name: '未分组',
      id: '0'
    }
  ]
  if (library_id) {
    return groupListMap.value[library_id] || defaultList
  }

  return defaultList
})

const libraryList = ref([])
const librarySelectAlertRef = ref(null)
const selectedLibraryRows = computed(() => {
  return libraryList.value.filter((item) => {
    if (!wxAppLibary.value) {
      return formState.library_id.includes(item.id) && item.type != 3
    } else {
      return formState.library_id.includes(item.id)
    }
  })
})

const selectedLibraryType = computed(() => {
  if (selectedLibraryRows.value && selectedLibraryRows.value.length) {
    return selectedLibraryRows.value[0].type
  }
  return -1
})

// 移除知识库
const handleDelKonwledge = (item) => {
  let index = formState.library_id.indexOf(item.id)
  formState.library_id.splice(index, 1)
}

const onChangeLibrarySelected = (checkedList) => {
  getList()
  formState.library_id = [...checkedList]
  formState.library_group_id = '0'
  nextTick(() => {
    getGroupList()
  })
}

const handleOpenSelectLibraryAlert = () => {
  librarySelectAlertRef.value.open([...formState.library_id])
}

// 获取知识库
const getList = async () => {
  const res = await getLibraryList({ type: '' })
  if (res) {
    libraryList.value = res.data || []
  }
}

const getGroupList = () => {
  nextTick(() => {
    if (formState.library_id.length == 0) {
      return
    }
    getLibraryGroup({
      group_type: selectedLibraryType.value == 0 ? 1 : 0,
      library_id: formState.library_id[0]
    }).then((res) => {
      groupListMap.value = {
        ...groupListMap.value,
        [formState.library_id[0]]: res.data
      }
      // groupListMap.value[formState.library_id[0]] = res.data || []
    })
  })
}

const changeValue = (text, selectedList, key) => {
  formState[key] = text
}

watch(
  () => formState,
  () => {
    update()
  },
  { deep: true }
)

watch(
  () => formState.library_id,
  () => {
    getGroupList()
  },
  { deep: true }
)

const handleClose = () => {
  emit('close')
}

onMounted(() => {
  init()
  getList()
})
</script>

<style lang="less" scoped>
@import '../form-block.less';
.flex-block-item {
  margin-bottom: 12px;
  .form-item-label {
    width: 100px;
    display: flex;
    align-items: center;
    gap: 2px;
    position: relative;
    &.required {
      &::before {
        content: '*';
        color: #fb363f;
      }
    }
  }
  .form-content-box {
    overflow: hidden;
  }
}
.node-form-content {
  ::v-deep(.ant-form-item-label) {
    padding: 0 0 4px;
  }
}
.konwledge-list-item {
  padding: 14px 16px;
  height: 72px;
  border: 1px solid #d8dde5;
  border-radius: 6px;
  display: flex;
  align-items: center;
  background-color: #fff;
  font-size: 14px;
  line-height: 22px;
  margin-bottom: 4px;
  .avatar-box {
    width: 40px;
    height: 40px;
    border: 1px solid var(--07, #f0f0f0);
    border-radius: 6px;
    margin-right: 8px;
    img {
      width: 100%;
      height: 100%;
    }
  }
  .list-name {
    color: #262626;
    font-weight: 600;
  }
  .list-intro {
    font-size: 12px;
    margin-top: 2px;
    line-height: 20px;
    color: #8c8c8c;
    max-width: 400px;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
  .content-box {
    flex: 1;
  }
}

.desc {
  color: #8c8c8c;
  line-height: 22px;
}

.btn-block {
  display: flex;
  gap: 4px;
  & > div {
    flex: 1;
  }
}
</style>
