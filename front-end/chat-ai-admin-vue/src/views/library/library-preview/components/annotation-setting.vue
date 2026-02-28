<template>
  <div class="annotation-setting-block" @click.stop="">
    <div class="title-block">
      {{ t('title_annotation') }}
      <a-tooltip>
        <template #title>{{ t('msg_annotation_description') }}</template>
        <QuestionCircleOutlined />
      </a-tooltip>
      :
    </div>
    <template v-for="(tag, index) in state.tags" :key="tag">
      <template v-if="state.editingTag === tag">
        <a-input
          :ref="(el) => setEditInputRef(el, tag)"
          v-model:value="state.editingValue"
          :maxlength="20"
          type="text"
          size="small"
          style="width: 120px;"
          @mousedown.stop
          @blur="handleEditConfirm"
          @pressEnter="handleEditConfirm"
          @keyup.esc="handleEditCancel"
        />
      </template>
      <template v-else>
        <a-tooltip v-if="tag.length > 10" :title="tag">
          <a-tag :closable="true" @close="handleClose(tag)" @click="handleEditTag(tag)">
            {{ `${tag.slice(0, 10)}...` }}
          </a-tag>
        </a-tooltip>
        <a-tag v-else :closable="true" @close="handleClose(tag)" @click="handleEditTag(tag)">
          {{ tag }}
        </a-tag>
      </template>
    </template>
    <template v-if="state.tags.length <= 50">
      <a-input
        v-if="state.inputVisible"
        ref="inputRef"
        v-model:value="state.inputValue"
        :maxLength="20"
        type="text"
        size="small"
        style="width: 78px;"
        @blur="handleInputConfirm"
        @pressEnter="handleInputConfirm"
      />
      <a-tag v-else style="background: #fff; border-style: dashed" @click="showInput">
        <plus-outlined />
        {{ t('btn_add_tag') }}
      </a-tag>
    </template>
  </div>
</template>

<script setup>
import { QuestionCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { ref, reactive, nextTick, onMounted, watch } from 'vue'
import { editParagraph } from '@/api/library'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.library-preview.components.annotation-setting')
const props = defineProps({
  currentItem: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['save'])

const inputRef = ref()
const editInputRefs = ref({})
const state = reactive({
  tags: [],
  inputVisible: false,
  inputValue: '',
  editingTag: '',
  editingValue: ''
})

const setEditInputRef = (el, tag) => {
  if (el) {
    editInputRefs.value[tag] = el
  }
}

const handleSave = () => {
  let parmas = {
    title: props.currentItem.title,
    content: props.currentItem.content,
    question: props.currentItem.question,
    answer: props.currentItem.answer,
    images: props.currentItem.images,
    category_id: props.currentItem.category_id,
    id: props.currentItem.id,
    similar_questions: JSON.stringify(state.tags)
  }
  editParagraph(parmas).then((res) => {
    emit('save', parmas)
  })
}

const handleClose = (removedTag) => {
  const tags = state.tags.filter((tag) => tag !== removedTag)
  state.tags = tags
  handleSave()
}

const handleEditTag = (tag) => {
  state.editingTag = tag
  state.editingValue = tag
  nextTick(() => {
    const inputEl = editInputRefs.value[tag]
    if (inputEl) {
      inputEl.focus()
    }
  })
}

const handleEditConfirm = () => {
  const editingValue = state.editingValue.trim()
  if (editingValue && editingValue !== state.editingTag) {
    const tagIndex = state.tags.indexOf(state.editingTag)
    if (tagIndex > -1) {
      state.tags[tagIndex] = editingValue
      handleSave()
    }
  }
  state.editingTag = ''
  state.editingValue = ''
}

const handleEditCancel = () => {
  state.editingTag = ''
  state.editingValue = ''
}

const showInput = () => {
  state.inputVisible = true
  nextTick(() => {
    inputRef.value.focus()
  })
}
const handleInputConfirm = () => {
  if(state.inputValue.trim() === '') {
    state.inputVisible = false
    return
  }
  const inputValue = state.inputValue
  let tags = state.tags
  if (inputValue && tags.indexOf(inputValue) === -1) {
    tags = [...tags, inputValue]
  }
  console.log(tags)
  Object.assign(state, {
    tags,
    inputVisible: false,
    inputValue: ''
  })
  handleSave()
}

onMounted(() => {
  state.tags = props.currentItem.similar_questions || []
})
</script>

<style lang="less" scoped>
.annotation-setting-block {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  .title-block {
    display: flex;
    align-items: center;
    font-size: 14px;
    font-weight: 600;
    color: #000000;
    span {
      margin-left: 2px;
    }
  }
  .ant-tag {
    margin: 0;
  }
}
</style>
