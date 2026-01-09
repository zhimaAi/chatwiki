<template>
  <node-common
    :properties="properties"
    :title="props.properties.node_name"
    :menus="menus"
    :icon-name="props.properties.node_icon_name"
    :isSelected="props.isSelected"
    :isHovered="props.isHovered"
    :node-key="props.properties.node_key"
    :node_type="props.properties.node_type"
    style="width: 420px"
  >
    <div class="ai-dialogue-node">
      <div class="field-list">
        <div class="field-item">
          <div class="field-item-label">知识库</div>
          <div class="field-item-content">
            <div class="field-value">
              <span class="field-key"> {{ library_name }}</span>
            </div>
          </div>
        </div>
        <template v-if="selectedLibraryType == 2">
          <div class="field-item">
            <div class="field-item-label">分段问题</div>
            <div class="field-item-content">
              <div class="field-value">
                <span class="field-key">
                  <at-text :options="valueOptions" :defaultValue="formState.qa_question" />
                </span>
              </div>
            </div>
          </div>
          <div class="field-item">
            <div class="field-item-label">分段答案</div>
            <div class="field-item-content">
              <div class="field-value">
                <span class="field-key">
                  <at-text :options="valueOptions" :defaultValue="formState.qa_answer" />
                </span>
              </div>
            </div>
          </div>
        </template>

        <template v-if="selectedLibraryType != 2 && selectedLibraryType >= 0">
          <div class="field-item" v-if="formState.import_type == 'content'">
            <div class="field-item-label">文档标题</div>
            <div class="field-item-content">
              <div class="field-value">
                <span class="field-key">
                  <at-text :options="valueOptions" :defaultValue="formState.normal_title" />
                </span>
              </div>
            </div>
          </div>
          <div class="field-item" v-if="formState.import_type == 'content'">
            <div class="field-item-label">文档内容</div>
            <div class="field-item-content">
              <div class="field-value">
                <span class="field-key">
                  <at-text :options="valueOptions" :defaultValue="formState.normal_content" />
                </span>
              </div>
            </div>
          </div>
          <div class="field-item" v-if="formState.import_type == 'url'">
            <div class="field-item-label">文档URL</div>
            <div class="field-item-content">
              <div class="field-value">
                <span class="field-key">
                  <at-text :options="valueOptions" :defaultValue="formState.normal_url" />
                </span>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>
  </node-common>
</template>

<script setup>
import { ref, reactive, watch, onMounted, inject, nextTick, onBeforeUnmount, computed } from 'vue'
import NodeCommon from '../base-node.vue'
import AtText from '../../at-input/at-text.vue'
import { useWorkflowStore } from '@/stores/modules/workflow'

const workflowStore = useWorkflowStore()

const libraryLists = computed(() => workflowStore.libraryLists)

const props = defineProps({
  properties: {
    type: Object,
    default: () => ({})
  },
  isSelected: { type: Boolean, default: false },
  isHovered: { type: Boolean, default: false }
})

const getNode = inject('getNode')
const resetSize = inject('resetSize')

const valueOptions = ref([])

// --- State ---
const menus = ref([])
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

const selectedLibraryRows = computed(() => {
  return libraryLists.value.filter((item) => {
    return formState.library_id.includes(item.id)
  })
})

const library_name = computed(() => {
  return selectedLibraryRows.value[0]?.library_name
})

const selectedLibraryType = computed(() => {
  if (selectedLibraryRows.value && selectedLibraryRows.value.length) {
    return selectedLibraryRows.value[0].type
  }
  return -1
})

const reset = () => {
  const dataRaw = props.properties.dataRaw || props.properties.node_params || '{}'
  let library_import = {}
  try {
    getValueOptions()
    library_import = JSON.parse(dataRaw).library_import || {}
    for (let key in library_import) {
      if (key == 'library_id') {
        formState.library_id = library_import[key] ? library_import[key].split(',') : []
      } else {
        formState[key] = library_import[key]
      }
    }
  } catch (e) {
    image_generation = {}
  }

  nextTick(() => {
    resetSize()
  })
}

function getValueOptions() {
  let options = getNode().getAllParentVariable()
  valueOptions.value = options || []
}

watch(
  () => props.properties,
  (newVal, oldVal) => {
    const newDataRaw = newVal.dataRaw || newVal.node_params || '{}'
    const oldDataRaw = oldVal.dataRaw || oldVal.node_params || '{}'

    if (newDataRaw != oldDataRaw) {
      reset()
    }
  },
  { deep: true }
)

onMounted(() => {
  reset()
  resetSize()
})

onBeforeUnmount(() => {})
</script>

<style lang="less" scoped>
.ai-dialogue-node {
  .field-list {
    .field-item {
      display: flex;
      margin-bottom: 8px;
      &:last-child {
        margin-bottom: 0;
      }
    }
    .field-item-label {
      width: 60px;
      line-height: 22px;
      margin-right: 8px;
      font-size: 14px;
      font-weight: 400;
      color: #262626;
      text-align: right;
    }
    .field-item-content {
      flex: 1;
      overflow: hidden;
    }
    .category-list {
      width: 100%;
      overflow: hidden;
    }
    .field-value {
      width: 100%;
      line-height: 16px;
      padding: 3px 4px;
      height: 24px;
      margin-bottom: 8px;
      border-radius: 4px;
      font-size: 12px;
      border: 1px solid #d9d9d9;
      color: #595959;
      background: #fff;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      &:last-child {
        margin-bottom: 0;
      }
      .category-value {
        width: 100%;
      }

      .field-type {
        padding: 1px 8px;
        margin-left: 4px;
        border-radius: 4px;
        font-size: 12px;
        line-height: 16px;
        font-weight: 400;
        background: #e4e6eb;
      }
    }
  }

  .options-list {
    flex: 1;
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }
  .options-item {
    display: flex;
    align-items: center;
    height: 22px;
    padding: 2px 2px 2px 4px;
    border-radius: 4px;
    border: 1px solid #d9d9d9;

    &.is-required .option-label::before {
      vertical-align: middle;
      content: '*';
      color: #fb363f;
      margin-right: 2px;
    }

    .option-label {
      color: var(--wf-color-text-3);
      font-size: 12px;
      margin-right: 4px;
    }

    .option-type {
      height: 18px;
      line-height: 18px;
      padding: 0 8px;
      border-radius: 4px;
      font-size: 12px;
      background-color: #e4e6eb;
      color: var(--wf-color-text-3);
    }
  }
}
</style>
