<template>
  <div>
    <div class="options-item is-required">
      <div class="options-item-tit">
        <div class="option-label">数据表名称</div>
        <div class="option-type">string</div>
      </div>
      <div class="min-input">
        <AtInput
          type="textarea"
          inputStyle="height: 64px;"
          :options="variableOptions"
          :defaultSelectedList="state?.tag_map?.name || []"
          :defaultValue="state.name"
          ref="atInputRef"
          @open="emit('updateVar')"
          @change="(val, tags) => changeValue('name', val, tags)"
          placeholder="请输入内容，键入“/”可以插入变量"
        >
          <template #option="{ label, payload }">
            <div class="field-list-item">
              <div class="field-label">{{ label }}</div>
              <div class="field-type">{{ payload.typ }}</div>
            </div>
          </template>
        </AtInput>
      </div>
      <div class="desc">数据表名称；名称中的首尾空格将会被默认去除，不可以包含 / \\ ? * : [ ] 等特殊字符"；1 字符 ～ 100
        字符
      </div>
    </div>
    <div class="options-item">
      <div class="options-item-tit">
        <div class="option-label">表格视图的名称</div>
        <div class="option-type">string</div>
      </div>
      <div class="min-input">
        <AtInput
          type="textarea"
          inputStyle="height: 64px;"
          :options="variableOptions"
          :defaultSelectedList="state?.tag_map?.default_view_name || []"
          :defaultValue="state.default_view_name"
          ref="atInputRef"
          @open="emit('updateVar')"
          @change="(val, tags) => changeValue('default_view_name', val, tags)"
          placeholder="请输入内容，键入“/”可以插入变量"
        >
          <template #option="{ label, payload }">
            <div class="field-list-item">
              <div class="field-label">{{ label }}</div>
              <div class="field-type">{{ payload.typ }}</div>
            </div>
          </template>
        </AtInput>
      </div>
      <div class="desc">表格视图的名称；名称中的首尾空格将会被默认去除，名称中不允许包含 [ ] 两个字符</div>
    </div>
    <div class="options-item">
      <div class="options-item-tit flex-between">
        <div class="flex-between">
          <div class="option-label">数据表的初始字段</div>
          <div class="option-type">jsonString</div>
        </div>
        <div class="btn-hover-wrap" @click="handleOpenFullAtModal">
          <FullscreenOutlined/>
        </div>
      </div>
      <div>
        <AtInput
          type="textarea"
          inputStyle="height: 64px;"
          :options="variableOptions"
          :defaultSelectedList="state?.tag_map?.fields || []"
          :defaultValue="state.fields"
          ref="atInputRef"
          @open="emit('updateVar')"
          @change="(val, tags) => changeValue('fields', val, tags)"
          placeholder="请输入内容，键入“/”可以插入变量"
        >
          <template #option="{ label, payload }">
            <div class="field-list-item">
              <div class="field-label">{{ label }}</div>
              <div class="field-type">{{ payload.typ }}</div>
            </div>
          </template>
        </AtInput>
      </div>
      <div class="desc">如果传入了 表格视图的名称，则必须传入
        数据表的初始字段。示例：[{"field_name":"索引字段","type":1},{"field_name":"单选","type":3,"ui_type":"SingleSelect","property":{"options":[{"name":"Enabled"},{"name":"Disabled"}]}}]
      </div>
    </div>

    <FullAtInput
      :options="variableOptions"
      :defaultSelectedList="state?.tag_map?.fields || []"
      :defaultValue="state.fields"
      placeholder="请输入内容，键入“/”可以插入变量"
      type="textarea"
      @open="emit('updateVar')"
      @change="(val, tags) => changeValue('fields', val, tags)"
      @ok="handleRefreshAtInput"
      ref="fullAtInputRef"
    />
  </div>
</template>

<script setup>
import {ref, reactive} from 'vue'
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import FullAtInput from "@/views/workflow/components/at-input/full-at-input.vue";
import {FullscreenOutlined} from '@ant-design/icons-vue';

const props = defineProps({
  variableOptions: {
    type: Array,
  }
})
const emit = defineEmits(['update', 'updateVar'])
const fullAtInputRef = ref(null)
const atInputRef = ref(null)
const state = reactive({
  name: '',
  default_view_name: '',
  fields: '',
  tag_map: {},
})

function init(nodeParams = null) {
  if (nodeParams && nodeParams?.plugin?.params?.arguments) {
    state.name = nodeParams.plugin.params.arguments?.name || ""
    state.default_view_name = nodeParams.plugin.params.arguments?.default_view_name || ""
    state.fields = nodeParams.plugin.params.arguments?.fields || ""
    state.tag_map = nodeParams.plugin.params.arguments?.tag_map || {}
  }
}

function changeValue(field, val, tags) {
  state[field] = val
  state.tag_map[field] = tags
  update()
}

function handleOpenFullAtModal() {
  fullAtInputRef.value.show()
}

function handleRefreshAtInput() {
  atInputRef.value.refresh()
}

function update() {
  emit('update', JSON.parse(JSON.stringify(state)))
}

defineExpose({
  init,
  update
})
</script>

<style scoped lang="less">
@import "common";

.flex-between {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
