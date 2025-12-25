<style lang="less" scoped>
.field-list {
  .field-list-row {
    position: relative;
    display: flex;
    flex-wrap: nowrap;
    align-items: center;
    margin-bottom: 4px;

    &:last-child {
      margin-bottom: 0;
    }
  }

  .field-list-col {
    padding: 0 4px;
  }

  .field-name-col,
  .field-name-head {
    flex: 1;
  }

  .field-type-col,
  .field-type-head {
    flex-basis: 90px;
  }

  .field-value-col,
  .field-value-head {
    width: 256px;
  }

  .field-list-col-head {
    line-height: 22px;
    margin-bottom: 4px;
    font-size: 14px;
    color: #262626;
  }

  .field-name-col,
  .field-type-col {
    line-height: 22px;
    font-size: 14px;
    color: #595959;
  }

  .field-del-head,
  .field-del-col {
    width: 24px;
    display: flex;
    align-items: center;
  }

  .field-del-col {
    text-align: right;

    .del-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 16px;
      height: 16px;
      font-size: 16px;
      color: #595959;
      cursor: pointer;
    }
  }
}
.wide-value-col {
  .field-value-col,
  .field-value-head {
    width: 300px;
  }
}
.field-value-inline {
  display: flex;
  align-items: center;
  gap: 2px;
  min-width: 0;
}
.inline-input {
  flex: 1;
  min-width: 0;
}
.add-btn-box {
  margin-top: 8px;
}
</style>

<template>
  <div>
    <div class="field-list" :class="{'wide-value-col': hasObjectPropRows}">
      <div class="field-list-row">
        <div class="field-list-col field-list-col-head field-name-head">字段名</div>
        <div class="field-list-col field-list-col-head field-type-head">类型</div>
        <div class="field-list-col field-list-col-head field-value-head" v-if="showInput">字段值</div>
        <div
          class="field-list-col field-list-col-head field-del-head"
          v-if="props.showDelete"
        ></div>
      </div>

      <div class="field-list-row" v-if="showEmptyFieldRow">
        <div class="field-list-col field-name-col">--</div>
        <div class="field-list-col field-type-col">--</div>
        <div class="field-list-col field-value-col">
          <a-tooltip title="请先选择数据库">
            <a-input :disabled="true" placeholder="请输入参数值，键入/插入变量"/>
          </a-tooltip>
        </div>
      </div>

      <div class="field-list-row" v-for="(item, index) in state.list" :key="item.field_name">
        <div class="field-list-col field-name-col">
          <span v-if="item.__required" style="color:#FB363F; margin-right:2px;">*</span>
          {{ item.properties?.[item.current_properties_key]?.name ||  item.name }}
          <a-tooltip v-if="item.__functionalTip" :title="item.__functionalTip">
            <QuestionCircleOutlined style="cursor: pointer;" />
          </a-tooltip>
        </div>
        <div class="field-list-col field-type-col">{{ item.ui_type }}</div>
        <div class="field-list-col field-value-col" v-if="showInput">
          <template v-if="isObjectWithProps(item)">
            <div class="field-value-inline">
              <a-select
                v-model:value="item.current_properties_key"
                @change="(val) => onSelectObjectProp(item, val)"
              >
                <a-select-option
                  v-for="(prop, key) in (item.properties || {})"
                  :key="key"
                  :value="key"
                >{{ prop.name || key }}</a-select-option>
              </a-select>
              <div class="inline-input">
                <AtInput
                  :options="atInputOptions"
                  :defaultValue="(item.properties?.[item.current_properties_key]?.value ?? '')"
                  :defaultSelectedList="(item.properties?.[item.current_properties_key]?.atTags ?? [])"
                  @open="showAtList"
                  @change="(text, selectedList) => changeAtInputValue(text, selectedList, item, index)"
                  placeholder='请输入参数，键入"/"插入变量' />
              </div>
            </div>
          </template>
          <template v-else>
            <AtInput
              :options="atInputOptions"
              :defaultValue="item.value"
              :defaultSelectedList="item.atTags"
              @open="showAtList"
              @change="(text, selectedList) => changeAtInputValue(text, selectedList, item, index)"
              placeholder='请输入参数，键入"/"插入变量' />
          </template>
          <div v-if="item.__error" style="color:#FB363F; margin-top:4px;">{{ item.__error }}</div>
        </div>
        <div class="field-list-col field-del-col" v-if="props.showDelete">
          <span class="del-btn" @click="handleDel(index)">
            <svg-icon class="del-icon" name="close-circle"></svg-icon>
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import {ref, reactive, watch, inject, onMounted, computed} from 'vue'
import AtInput from '../../at-input/at-input.vue'

const emit = defineEmits(['change'])

const props = defineProps({
  showDelete: {
    type: Boolean,
    default: true
  },
  showInput: {
    type: Boolean,
    default: true
  },
  showEmptyFieldRow: {
    type: Boolean,
    default: false
  },
  fields: {
    type: Array,
    default: () => ([])
  },
  list: {
    type: Array,
    default: () => []
  },
  msgtype: {
    type: String,
    default: ''
  }
})

const getNode = inject('getNode')

const state = reactive({
  list: [],
  keys: [],
})

const hasObjectPropRows = computed(() => {
  return state.list.some(it => isObjectWithProps(it))
})

function getReplyMode () {
  const m = String(props.msgtype || '').toLowerCase()
  if (m === 'link') return 'imageText'
  if (m === 'miniprogrampage') return 'card'
  if (['text','image','url'].includes(m)) return m
  return 'text'
}

const fieldRuleMap = {
  imageText: {
    url: {
      required: true,
      validator: (v) => {
        // `【([a-f0-9]{32}\.)?[a-zA-Z_][a-zA-Z0-9_\-.]*】` 如果是这个正则表示是引用 就不用校验
        const refRegex = /^【([a-f0-9]{32}\.)?[a-zA-Z_][a-zA-Z0-9_\-.]*】$/
        if (refRegex.test(String(v||''))) return true
        return /^https?:\/\//.test(String(v||''))
      },
      validationTip: '链接需以http或https开头',
      functionalTip: '跳转到外链，支持插入变量，链接需以http或https开头。'
    },
    title: {
      required: true,
      validator: (v) => !!String(v||'').trim(),
      validationTip: '请输入标题',
      functionalTip: '展示在链接卡片的标题位置'
    },
    description: {
      required: true,
      validator: (v) => String(v||'').trim().length > 0 && String(v||'').length <= 300,
      validationTip: '请输入描述，最多300字',
      functionalTip: '展示在链接卡片的描述位置'
    },
    thumb_url: {
      required: true,
      validator: (v) => !!String(v||'').trim(),
      validationTip: '请上传链接图片',
      functionalTip: '支持png/jpg/jpeg，建议≤2M'
    }
  },
  text: {
    content: {
      required: true,
      validator: (v) => String(v||'').trim().length > 0 && String(v||'').length <= 300,
      validationTip: '请输入文本内容，最多300字',
      functionalTip: '支持插入变量，最多300字'
    }
  },
  image: {
    thumb_url: {
      required: true,
      validator: (v) => !!String(v||'').trim(),
      validationTip: '请上传图片',
      functionalTip: '支持png/jpg/jpeg，建议≤2M'
    }
  },
  url: {
    title: {
      required: true,
      validator: (v) => !!String(v||'').trim(),
      validationTip: '请输入标题',
      functionalTip: '展示在卡片标题位置'
    },
    url: {
      required: true,
      validator: (v) => {
        // `【([a-f0-9]{32}\.)?[a-zA-Z_][a-zA-Z0-9_\-.]*】` 如果是这个正则表示是引用 就不用校验
        const refRegex = /^【([a-f0-9]{32}\.)?[a-zA-Z_][a-zA-Z0-9_\-.]*】$/
        if (refRegex.test(String(v||''))) return true
        return /^https?:\/\//.test(String(v||''))
      },
      validationTip: '链接需以http或https开头',
      functionalTip: '点击后跳转的链接，链接需以http或https开头。'
    }
  },
  card: {
    title: {
      required: true,
      validator: (v) => !!String(v||'').trim(),
      validationTip: '请输入小程序标题',
      functionalTip: '展示在卡片标题位置'
    },
    appid: {
      required: true,
      validator: (v) => !!String(v||'').trim(),
      validationTip: '请输入小程序appID',
      functionalTip: '小程序右上角>名字>更多资料>AppID；公众号需与小程序关联'
    },
    page_path: {
      required: true,
      validator: (v) => /^\//.test(String(v||'')),
      validationTip: '请输入以/开头的小程序路径',
      functionalTip: '向开发者获取路径，如/pages/index/index'
    },
    thumb_url: {
      required: true,
      validator: (v) => !!String(v||'').trim(),
      validationTip: '请上传小程序封面',
      functionalTip: '支持png/jpg/jpeg，建议≤2M'
    }
  }
}

function applyFunctionalTips () {
  const mode = getReplyMode()
  const rules = fieldRuleMap[mode] || {}
  state.list.forEach((it, idx) => {
    const key = (it.field_name || '').toLowerCase()
    const rule = rules[key] || {}
    let backendOuter = {}
    if (typeof it.description !== 'undefined' || typeof it.desc !== 'undefined' || typeof it.required !== 'undefined') {
      backendOuter = { functionalTip: it.description || it.desc, required: it.required }
    }
    let propMeta = {}
    if (isObjectWithProps(it)) {
      const p = it.properties?.[it.current_properties_key]
      if (p && (typeof p.description !== 'undefined' || typeof p.desc !== 'undefined' || typeof p.required !== 'undefined')) {
        propMeta = { functionalTip: p.description || p.desc, required: p.required }
      }
    }
    const merged = {
      functionalTip: propMeta.functionalTip ?? backendOuter.functionalTip ?? rule.functionalTip ?? '',
      required: (propMeta.required ?? backendOuter.required ?? rule.required ?? it.required) ? true : false
    }
    state.list[idx].__functionalTip = merged.functionalTip
    state.list[idx].__required = !!merged.required
    state.list[idx].__error = null
  })
}

function validateField (item) {
  const mode = getReplyMode()
  const rules = fieldRuleMap[mode] || {}
  const key = (item.field_name || '').toLowerCase()
  let r = rules[key]
  if (!r && isObjectWithProps(item)) {
    const p = item.properties?.[item.current_properties_key]
    if (p) {
      r = { required: !!p.required }
    }
  }
  const val = isObjectWithProps(item) ? (item.properties?.[item.current_properties_key]?.value ?? '') : item.value
  if (!r) {
    item.__error = null
    return true
  }
  const tags = isObjectWithProps(item) ? (item.properties?.[item.current_properties_key]?.atTags ?? []) : item.atTags
  if (r.required && !String(val || '').trim() && (!Array.isArray(tags) || tags.length === 0)) {
    item.__error = r.validationTip
    return false
  }
  if (r.validator && !r.validator(val)) {
    item.__error = r.validationTip
    return false
  }
  item.__error = null
  return true
}

watch(() => props.list, (newVal) => {
  state.list = newVal
  state.keys = newVal.map(i => i.field_name)
  state.list.forEach(initializeObjectField)
  applyFunctionalTips()
}, {
  immediate: true
})

const atInputOptions = ref([])

const getAtInputOptions = () => {
  let options = getNode().getAllParentVariable();

  atInputOptions.value = options || []
}

const showAtList = () => {
  getAtInputOptions()
}

const changeAtInputValue = (text, selectedList, item, index) => {
  if (isObjectWithProps(item)) {
    const k = item.current_properties_key
    if (!item.properties[k]) item.properties[k] = {}
    item.properties[k].value = text
    item.properties[k].atTags = selectedList
  } else {
    state.list[index].value = text
    state.list[index].atTags = selectedList
  }
  change()
}

const handleDel = (index) => {
  state.list.splice(index, 1)

  change()
}

const change = () => {
  state.keys = state.list.map(i => i.field_name)
  emit('change', [...state.list], [...state.keys])
}

onMounted(() => {
  getAtInputOptions()
})

function isObjectWithProps (item) {
  return item.media_component
}

function onSelectObjectProp (item, key) {
  const prop = item.properties?.[key]
  item.current_properties_key = key
  // 保持父字段名不变，仅更新显示类型与提示
  item.ui_type = mapPropType(prop?.type)
  item.__functionalTip = prop?.desc || ''
  item.__required = !!prop?.required
  // 重置当前子字段的值
  if (!item.properties[key]) item.properties[key] = {}
  item.properties[key].value = item.properties[key].value ?? ''
  item.properties[key].atTags = item.properties[key].atTags ?? []
  change()
}

function initializeObjectField (item) {
  if (!isObjectWithProps(item)) return
  const keys = Object.keys(item.properties)
  if (!item.current_properties_key) {
    const presetKey = keys[0]
    onSelectObjectProp(item, presetKey)
  }
}

function mapPropType(t) {
  switch ((t || '').toLowerCase()) {
    case 'string': return 'Text'
    case 'number': return 'Number'
    case 'integer': return 'Number'
    case 'boolean': return 'Checkbox'
    case 'array<string>': return 'MultiSelect'
    default: return t || 'Text'
  }
}

function validateAll () {
  let ok = true
  const errors = []
  state.list.forEach(it => {
    if (!validateField(it)) {
      ok = false
      const prop = isObjectWithProps(it) ? (it.properties?.[it.current_properties_key] || {}) : null
      errors.push({
        field_name: (prop?.name || it.name || it.field_name),
        message: prop?.error || it.__error,
        functionalTip: prop?.desc || it.__functionalTip
      })
    }
  })
  return { ok, errors }
}

defineExpose({ validateAll })

</script>
