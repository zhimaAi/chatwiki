<template>
  <a-drawer
    v-model:open="visible"
    :closable="false"
    width="568px"
    body-style="padding: 16px;"
    header-style="padding: 16px 16px 0;border: none;">
    <template #title>
      <div class="title-block">
        <div class="left">
          <img class="avatar" src="http://dev1.zhima_chat_ai.applnk.cn/upload/default/mcp_avatar.svg"/>
          <a-input class="input" :value="properties.node_name"/>
        </div>
        <CloseOutlined @click="visible = false" class="icon"/>
      </div>
      <div class="title-desc">{{ properties.tool_info.description }}</div>
    </template>
    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/input.svg" class="title-icon"/>输入</div>
      </div>
      <div v-for="(item, key) in formState.params"
           :key="key"
           :class="['options-item', {'is-required': item.required}]">
        <div class="options-item-tit">
          <div class="option-label">{{ key }}</div>
          <div class="option-type">{{ item.type }}</div>
        </div>
        <div>
          <!--          <a-input v-model:value="item.value" @change="update" placeholder="键入/插入变量"/>-->
          <AtInput
            type="textarea"
            inputStyle="height: 64px;"
            :options="variableOptions"
            :defaultSelectedList="item.tags"
            :defaultValue="item.value"
            ref="atInputRef"
            @open="getValueVariableList"
            @change="(text, selectedList) => update(item, text, selectedList)"
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
        <div class="desc">{{ item.description }}</div>
      </div>
    </div>
    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/output.svg" class="title-icon"/>输出</div>
      </div>
      <div class="options-item">
        <div class="options-item-tit">
          <div class="option-label">输出</div>
        </div>
        <div class="options-item-tit">
          <div class="option-label">text</div>
          <div class="option-type">string</div>
        </div>
        <div class="desc">工具生成的内容</div>
      </div>
    </div>
  </a-drawer>
</template>

<script setup>
import {ref, reactive, inject} from 'vue';
import {CloseOutlined} from '@ant-design/icons-vue';
import AtInput from "@/views/workflow/components/nodes/at-input/at-input.vue";

const emit = defineEmits(['update'])

const getNode = inject('getNode')

const props = defineProps({
  properties: {
    type: Object,
    default: () => ({})
  }
})

const visible = ref(false)
const formState = reactive({
  params: {}
})
const variableOptions = ref([])

function open() {
  dataFormat()
  visible.value = true
}

function dataFormat() {
  let nodeParams = JSON.parse(props.properties?.node_params)
  let inputSchema = props.properties?.tool_info?.inputSchema || {}
  let params = inputSchema?.properties || {}
  let requireds = inputSchema?.required || []
  params = JSON.parse(JSON.stringify(params))
  for (let key in params) {
    params[key].value = nodeParams?.mcp?.arguments[key] || ''
    params[key].tags = nodeParams?.mcp?.tag_map[key] || []
    params[key].required = requireds.includes(key)
  }
  formState.params = params
}

function getValueVariableList() {
  variableOptions.value = getNode().getAllParentVariable()
}

function update(item, text, selectedList) {
  item.value = text
  item.tags = selectedList
  let data = {
    arguments: {},
    tag_map: {}
  }
  for (let key in formState.params) {
    data.arguments[key] = formState.params[key].value
    data.tag_map[key] = formState.params[key].tags
  }
  emit('update', data)
}

defineExpose({
  open,
})
</script>

<style scoped lang="less">
.node-options {
  background: #f2f4f7;
  border-radius: 6px;
  padding: 12px;
  margin-top: 16px;

  &:first-child {
    margin-top: 0;
  }

  .options-title {
    color: var(--wf-color-text-1);
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-weight: 600;
    height: 22px;;
    line-height: 22px;
    font-size: 14px;

    .title-icon {
      width: 16px;
      height: 16px;
      vertical-align: -3px;
      margin-right: 8px;;
    }

    .acton-box {
      font-weight: 400;
    }
  }

  .options-item {
    display: flex;
    flex-direction: column;
    margin-top: 12px;
    line-height: 22px;
    gap: 4px;

    .options-item-tit {
      display: flex;
      align-items: center;
    }

    .option-label {
      color: var(--wf-color-text-1);
      font-size: 14px;
      margin-right: 8px;
    }

    .desc {
      color: var(--wf-color-text-2);
    }


    &.is-required .option-label::before {
      content: '*';
      color: #FB363F;
      display: inline-block;
      margin-right: 2px;
    }

    .option-type {
      height: 22px;
      line-height: 18px;
      padding: 0 8px;
      border-radius: 6px;
      border: 1px solid rgba(0, 0, 0, 0.15);
      background-color: #fff;
      color: var(--wf-color-text-3);
      font-size: 12px;
    }

    .item-actions-box {
      display: flex;
      align-items: center;

      .action-btn {
        margin-left: 12px;
        font-size: 16px;
        color: #595959;
        cursor: pointer;
      }
    }
  }
}

.title-block {
  display: flex;
  align-items: center;
  justify-content: space-between;

  .left {
    flex: 1;
    display: flex;
    align-items: center;
    gap: 8px;
    margin-right: 24px;

    .input {
      font-weight: 600;
    }

    .avatar {
      width: 20px;
      height: 20px;
      border-radius: 4.71px;
      flex-shrink: 0;
    }
  }
}

.title-desc {
  margin-top: 8px;
  font-size: 14px;
  font-weight: 400;
  color: var(--wf-color-text-2);
}
</style>
