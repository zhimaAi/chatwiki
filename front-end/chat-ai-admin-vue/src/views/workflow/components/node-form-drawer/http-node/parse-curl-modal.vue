<template>
  <a-modal
    v-model:open="open"
    :title="t('title_parse_curl')"
    :width="680"
    :confirm-loading="loading"
    @ok="handleParse"
    @cancel="handleClose"
    :okText="t('btn_parse')"
    :cancelText="t('btn_cancel')"
  >
    <div class="parse-curl-modal">
      <a-alert
        :message="t('msg_parse_curl_tip')"
        type="info"
        show-icon
        class="mb16"
      />
      <a-form-item label="">
        <a-textarea
          class="curl-textarea"
          v-model:value="curlInput"
          :placeholder="placeholderText"
          :auto-size="{ minRows: 8, maxRows: 8 }"
          :maxlength="50000"
          show-count
        />
      </a-form-item>
      <div v-if="parseError" class="parse-error">
        <a-alert :message="parseError" type="error" show-icon closable @close="parseError = ''" />
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'
import { ref } from 'vue'
import parseCurl from '@bany/curl-to-json';

const { t } = useI18n('views.workflow.components.node-form-drawer.http-node.parse-curl-modal')

const emit = defineEmits(['parse', 'close'])

const open = ref(false)
const loading = ref(false)
const curlInput = ref('')
const parseError = ref('')
const placeholderText = ref(t('ph_curl_example', {data: "'{\"key\":\"value\"}'"}))

const show = () => {
  open.value = true
  curlInput.value = ''
  parseError.value = ''
  loading.value = false
}

const handleClose = () => {
  open.value = false
  curlInput.value = ''
  parseError.value = ''
  emit('close')
}

const hide = () => {
  open.value = false
}

const handleParse = async () => {
  if (!curlInput.value.trim()) {
    parseError.value = t('msg_input_curl')
    return
  }

  parseError.value = ''

  try {
    let curlCommand = curlInput.value.trim()
    // 1. 预处理：将 shell 转义的单引号 ('\'') 替换回普通的单引号 (')
    // 注意：这步操作是为了把 Shell 语法变成 JS 库容易理解的纯字符串
    // 很多简单的 parser 不支持 shell 的拼接语法
    const PLACEHOLDER = "__SINGLE_QUOTE_PLACEHOLDER__";
    // 匹配 Shell 语法的转义结构: ' + \ + ' + '
    let cleanCommand = curlCommand.replace(/'\\''/g, PLACEHOLDER);

    // Step 2: 移除换行符 (常规清洗)
    cleanCommand = cleanCommand.replace(/\\\n/g, " ").replace(/\\/g, "");
 
    let out = parseCurl(decodeURIComponent(cleanCommand));
   
    // 处理查询参数：解码被 parseCurl 编码的特殊字符
    if (out.params) {
      const decodedQuery = {};
      for (const [key, value] of Object.entries(out.params)) {
        // 对键进行解码，将 a%5B0%5D 转回 a[0]
        const decodedKey = decodeURIComponent(key);
        // 对值也进行解码（以防值中也有编码字符）
        const decodedValue = typeof value === 'string' ? decodeURIComponent(value) : value;
        decodedQuery[decodedKey] = decodedValue;
      }
      out.params = decodedQuery;
    }

    // 处理 POST 请求体数据：解码被 parseCurl 编码的特殊字符
    // 只对 application/x-www-form-urlencoded 类型的数据进行解码
    // JSON 类型不需要解码
    if (out.data) {
      const contentType = out.header && out.header['Content-Type'] 
        ? out.header['Content-Type'].toLowerCase() 
        : '';
      
      // 判断是否是 x-www-form-urlencoded 类型
      const isFormUrlEncoded = contentType.includes('application/x-www-form-urlencoded');
      
      if (isFormUrlEncoded && typeof out.data === 'object') {
        const decodedData = {};
        for (const [key, value] of Object.entries(out.data)) {
          // 对键进行解码
          const decodedKey = decodeURIComponent(key);
          // 对值也进行解码
          const decodedValue = typeof value === 'string' ? decodeURIComponent(value) : value;
          decodedData[decodedKey] = decodedValue;
        }
        out.data = decodedData;
      }
      // JSON 类型的数据不需要解码，保持原样
    }

    out = JSON.stringify(out).split(PLACEHOLDER).join("'");;

    out = JSON.parse(out);
    
    if (out.method !== 'GET' && out.method !== 'POST') {
      parseError.value = t('msg_only_get_post')
      return
    }

    emit('parse', out)
  } catch (error) {
    parseError.value = error.message
  }
}

defineExpose({
  show,
  hide,
  handleClose
})
</script>

<style lang="less" scoped>
.parse-curl-modal {
  .mb16 {
    margin-bottom: 16px;
  }

  .curl-textarea {
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 12px;
  }
}
</style>
