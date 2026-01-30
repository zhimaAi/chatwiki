<template>
  <a-modal
    v-model:open="visible"
    width="720px"
    :confirm-loading="saving"
    :title="editingId ? '编辑工具' : '添加工具'"
    @ok="save">
    <div class="gray-block" style="background: #FFF;">
      <div class="gray-block-title-left" style="margin-bottom: 16px;">
        <svg-icon name="jibenpeizhi" size="16" /> <span>基本信息</span>
      </div>
      <a-form ref="formRef" :model="form" :rules="rules" layout="vertical">
        <a-form-item name="node_name_en" label="名称" required>
          <a-input
            v-model:value="form.node_name_en"
            placeholder="请输入名称"
            :maxlength="50"
          />
          <div class="form-tips">只能输入英文数字和字符"_"、"."、"-"，最多不超过50个字符</div>
        </a-form-item>
        <a-form-item name="title" label="备注名" required>
          <a-input
            v-model:value="form.title"
            placeholder="请输入备注名"
            :maxlength="20"
          />
          <div class="form-tips">用于显示，不超过20个字</div>
        </a-form-item>
        <a-form-item name="description" label="描述" required>
          <a-textarea
            v-model:value="form.description"
            :auto-size="{ minRows: 2, maxRows: 4 }"
            placeholder="请输入描述"
            :maxlength="500"
          />
          <div class="form-tips">还可以输入 {{ 500 - (form.description ? form.description.length : 0) }} 个字</div>
        </a-form-item>
      </a-form>
    </div>
    <div class="gray-block">
      <div class="gray-block-title">
        <div class="gray-block-title-left">
          <svg-icon name="input" size="16"/> <span>输入</span>
        </div>
        <a-button ghost type="primary" size="small" @click="showParseModal">
          <CodeOutlined /> 导入cURL
        </a-button>
      </div>
      <a-form class="form-box" ref="inputFormRef" :model="form" :rules="inputRules" layout="vertical">
        <a-form-item label="请求地址" required>
          <div class="flex-block-item">
            <a-select v-model:value="form.method" style="width: 120px">
              <a-select-option value="GET">GET</a-select-option>
              <a-select-option value="POST">POST</a-select-option>
              <a-select-option value="PUT">PUT</a-select-option>
              <a-select-option value="DELETE">DELETE</a-select-option>
              <a-select-option value="PATCH">PATCH</a-select-option>
            </a-select>
            <a-input class="flex1" v-model:value="form.url" placeholder="请输入请求地址"/>
          </div>
        </a-form-item>
        <div class="array-form-box">
          <div class="form-item-label">HEADERS</div>
          <div class="form-item-list" v-for="(item, i) in form.headers" :key="i">
            <div class="flex-block-item">
              <a-input style="width: 120px" v-model:value="item.key" placeholder="Key"/>
              <a-input class="flex1" v-model:value="item.value" placeholder="value"/>
              <div class="btn-hover-wrap" @click="form.headers.splice(i,1)"><CloseCircleOutlined /></div>
            </div>
          </div>
          <a-button block type="dashed" :icon="h(PlusOutlined)" @click="form.headers.push({key:'', value:''})">添加参数</a-button>
        </div>
        <div class="array-form-box">
          <div class="form-item-label">PARAMS</div>
          <div class="form-item-list" v-for="(item, i) in form.params" :key="i">
            <div class="flex-block-item">
              <a-input style="width: 120px" v-model:value="item.key" placeholder="Key"/>
              <a-input class="flex1" v-model:value="item.value" placeholder="value"/>
              <div class="btn-hover-wrap" @click="form.params.splice(i,1)"><CloseCircleOutlined /></div>
            </div>
          </div>
          <a-button block type="dashed" :icon="h(PlusOutlined)" @click="form.params.push({key:'', value:''})">添加参数</a-button>
        </div>
        <div class="array-form-box mt12" v-if="['POST','PUT','PATCH'].includes(form.method)">
          <a-form-item label="BODY" class="body-form-item">
            <a-radio-group v-model:value="form.body_type">
              <a-radio :value="0">none</a-radio>
              <a-radio :value="1">x-www-form-urlencoded</a-radio>
              <a-radio :value="2">JSON</a-radio>
            </a-radio-group>
          </a-form-item>
          <div v-if="form.body_type == 1">
            <div class="form-item-list" v-for="(item, i) in form.body" :key="i">
              <div class="flex-block-item">
                <a-input style="width: 120px" v-model:value="item.key" placeholder="Key"/>
                <a-input class="flex1" v-model:value="item.value" placeholder="value"/>
                <div class="btn-hover-wrap" @click="form.body.splice(i,1)"><CloseCircleOutlined /></div>
              </div>
            </div>
            <a-button block type="dashed" :icon="h(PlusOutlined)" @click="form.body.push({key:'', value:''})">添加参数</a-button>
          </div>
          <a-form-item v-if="form.body_type == 2">
            <a-textarea v-model:value="form.body_raw" :auto-size="{ minRows: 4, maxRows: 10 }" placeholder='{ "foo": "bar" }'/>
          </a-form-item>
        </div>
        <div class="array-form-box timeout-form-item">
          <a-form-item label="超时时间(秒)">
            <a-input-number v-model:value="form.timeout" :min="1" :max="600" style="width: 160px" />
          </a-form-item>
        </div>
      </a-form>
    </div>
    <div class="gray-block">
      <div class="gray-block-title">
        <div class="gray-block-title-left">
          <span>鉴权</span>
          <a-tooltip>
            <template #title>鉴权参数：导出模板CSL文件，或者上架模板时，自动清空参数值</template>
            <QuestionCircleOutlined />
          </a-tooltip>
        </div>
      </div>
      <div class="output-box">
        <div class="output-block">
          <div class="output-item" style="width: 33%">Key</div>
          <div class="output-item" style="width: 33%">Value</div>
          <div class="output-item" style="width: 33%">Add To</div>
        </div>
        <div class="array-form-box" @mousedown.stop="">
          <div class="form-item-list" v-for="(item, i) in form.auth" :key="i">
            <div class="flex-block-item" style="gap: 12px">
              <a-input style="width: 33%" v-model:value="item.key" placeholder="请输入"/>
              <a-input style="width: 33%" v-model:value="item.value" placeholder="请输入"/>
              <a-select v-model:value="item.add_to" style="width: 33%">
                <a-select-option value="HEADERS">HEADERS</a-select-option>
                <a-select-option value="PARAMS">PARAMS</a-select-option>
                <a-select-option value="BODY">BODY</a-select-option>
              </a-select>
              <div class="btn-hover-wrap" @click="handleDelAuthentication(i)">
                <CloseCircleOutlined />
              </div>
            </div>
          </div>
          <a-button @click="handleOpenAuthenticationModal" :icon="h(PlusOutlined)" block type="dashed">添加参数</a-button>
        </div>
      </div>
    </div>
    <div class="gray-block">
      <div class="gray-block-title">
        <div class="gray-block-title-left">
          <svg-icon name="output" size="16"/> <span>输出</span>
        </div>
      </div>
      <div class="output-box">
        <div class="output-block">
          <div class="output-item" style="width: 33%">参数Key</div>
          <div class="output-item" style="width: 33%">类型</div>
        </div>
        <div class="array-form-box">
          <div class="form-item-list" v-for="(item, i) in form.output" :key="i">
            <div class="flex-block-item" style="gap: 12px">
              <a-input style="width: 214px" v-model:value="item.key" placeholder="请输入"/>
              <a-select @change="onTypeChange(item)" v-model:value="item.typ" style="width: 214px" placeholder="请选择">
                <a-select-option v-for="op in typOptions" :key="op.value" :value="op.value">{{ op.value }}</a-select-option>
              </a-select>
              <div class="btn-hover-wrap" v-if="item.typ == 'object'" @click="onAddSubs(i)">
                <PlusCircleOutlined />
              </div>
              <div class="btn-hover-wrap" @click="form.output.splice(i,1)"><CloseCircleOutlined /></div>
            </div>
            <div class="sub-field-box" v-if="item.subs && item.subs.length > 0">
              <SubKey :data="item.subs" :level="2" :typOptions="typOptions" />
            </div>
          </div>
          <a-button block type="dashed" :icon="h(PlusOutlined)" @click="form.output.push({key:'', typ:'string', subs: []})">添加参数</a-button>
        </div>
      </div>
    </div>
    <ParseCurlModal ref="parseCurlModalRef" @parse="handleParseResult" />
    <AddAuthenticationModal ref="addAuthModalRef" @ok="handleSaveAuthentication" />
  </a-modal>
</template>

<script setup>
import { ref, reactive, h, nextTick } from 'vue'
import { CloseCircleOutlined, PlusOutlined, PlusCircleOutlined, CodeOutlined, QuestionCircleOutlined } from '@ant-design/icons-vue'
import ParseCurlModal from '@/views/workflow/components/node-form-drawer/http-node/parse-curl-modal.vue'
import { saveHttpToolItem, delHttpToolItem } from '@/api/robot/http_tool.js'
import AddAuthenticationModal from './add-authentication-modal.vue'
import SubKey from '@/views/workflow/components/node-form-drawer/http-node/subs-key.vue'
import { message, Modal } from 'ant-design-vue'

const emit = defineEmits(['ok'])
const props = defineProps({
  toolId: {
    type: Number,
    default: 0
  }
})

const visible = ref(false)
const saving = ref(false)
const parseCurlModalRef = ref(null)
const addAuthModalRef = ref(null)
const formRef = ref(null)
const inputFormRef = ref(null)

const form = reactive({
  title: '',
  description: '',
  node_name_en: '',
  method: 'GET',
  url: '',
  auth: [],
  headers: [],
  params: [],
  body_type: 0,
  body: [],
  body_raw: '',
  output: [{key:'', typ:'string', subs: []}],
  timeout: 30
})
const editingId = ref(0)

const typOptions = [
  { lable: 'string', value: 'string' },
  { lable: 'number', value: 'number' },
  { lable: 'boole', value: 'boole' },
  { lable: 'float', value: 'float' },
  { lable: 'object', value: 'object' },
  { lable: 'array<string>', value: 'array<string>' },
  { lable: 'array<number>', value: 'array<number>' },
  { lable: 'array<boole>', value: 'array<boole>' },
  { lable: 'array<float>', value: 'array<float>' },
  { lable: 'array<object>', value: 'array<object>' }
]

const rules = {
  node_name_en: [
    { required: true, message: '请输入名称' },
    {
      validator: (rule, value) => {
        const reg = /^[a-zA-Z0-9_.-]+$/
        if (!reg.test(value)) return Promise.reject('只能输入英文数字和字符"_"、"."、"-"')
        return Promise.resolve()
      }
    }
  ],
  title: [
    { required: true, message: '请输入备注名' }
  ],
  description: [
    { required: true, message: '请输入描述' }
  ]
}
const inputRules = {
  auth: [
    {
      validator: (_rule, value) => {
        const arr = Array.isArray(value) ? value : []
        if (!arr.length) return Promise.reject('请添加鉴权参数')
        const ok = arr.every(it => String(it.key || '').trim().length > 0)
        if (!ok) return Promise.reject('请填写鉴权参数Key')
        return Promise.resolve()
      }
    }
  ]
}

function show(item = null) {
  if (item && item.id) {
    editingId.value = item.id
    // 支持两种结构：旧的http_tool_item结构与节点结构
    if (item.node_params?.curl || item.data_raw) {
      const info = item.http_tool_info || {}
      form.node_name_en = info.http_tool_node_name_en || item.node_name_en || ''
      form.title = info.http_tool_node_name || item.node_name || ''
      form.description = item.node_description || info.http_tool_node_description || ''
      let curl = item.node_params?.curl
      if (!curl && item.data_raw) {
        try {
          const parsed = JSON.parse(item.data_raw || '{}')
          curl = parsed.curl || {}
        } catch {
          curl = {}
        }
      }
      form.method = curl.method || 'GET'
      form.url = String(curl.rawurl || '').replace(/[`"]/g, '').trim()
      form.headers = Array.isArray(curl.headers) ? JSON.parse(JSON.stringify(curl.headers)) : []
      form.params = Array.isArray(curl.params) ? JSON.parse(JSON.stringify(curl.params)) : []
      form.body_type = curl.type ?? 0
      form.body = Array.isArray(curl.body) ? JSON.parse(JSON.stringify(curl.body)) : []
      form.body_raw = curl.body_raw || ''
      form.output = Array.isArray(curl.output) ? JSON.parse(JSON.stringify(curl.output)) : [{key:'', typ:'string', subs: []}]
      const httpAuth = curl.http_auth || []
      form.auth = Array.isArray(httpAuth) ? httpAuth.map(a => ({ key: a.key, value: a.value, add_to: a.add_to || 'HEADERS' })) : []
      form.timeout = Number(curl.timeout ?? 30)
    } else {
      form.node_name_en = item.node_name_en || ''
      form.title = item.robot_name || item.name || ''
      form.description = item.robot_intro || ''
      const raw = item.raw || {}
      form.method = raw.method || 'GET'
      form.url = raw.url || ''
      form.headers = Array.isArray(raw.headers) ? JSON.parse(JSON.stringify(raw.headers)) : []
      form.params = Array.isArray(raw.params) ? JSON.parse(JSON.stringify(raw.params)) : []
      form.body_type = raw.body_type ?? 0
      form.body = Array.isArray(raw.body) ? JSON.parse(JSON.stringify(raw.body)) : []
      form.body_raw = raw.body_raw || ''
      form.output = Array.isArray(raw.output) ? JSON.parse(JSON.stringify(raw.output)) : [{key:'', typ:'string', subs: []}]
      form.timeout = 30
    }
  } else {
    editingId.value = 0
    Object.assign(form, {
      node_name_en: '',
      title: '',
      description: '',
      method: 'GET',
      url: '',
      headers: [],
      params: [],
      body_type: 0,
      body: [],
      body_raw: '',
      output: [{key:'', typ:'string', subs: []}],
      timeout: 30
    })
  }
  visible.value = true
}
function handleParseResult(parsed) {
  const handleValue = (v) => typeof v === 'string' ? v : JSON.stringify(v)
  let isJson = false
  form.method = parsed.method || form.method
  form.url = parsed.url || form.url
  if (parsed.header && typeof parsed.header === 'object') {
    const list = Object.keys(parsed.header).map(k => {
      const val = parsed.header[k]
      if (String(k).toLowerCase() === 'content-type' && String(val).toLowerCase().includes('application/json')) {
        isJson = true
      }
      return { key: k, value: handleValue(val) }
    }).filter(it => String(it.key).toLowerCase() !== 'content-type')
    form.headers = list
  }
  if (parsed.params && typeof parsed.params === 'object') {
    const list = Object.keys(parsed.params).map(k => ({ key: k, value: handleValue(parsed.params[k]) }))
    form.params = list
  }
  if (parsed.data && typeof parsed.data === 'object') {
    const body = Object.keys(parsed.data).map(k => ({ key: k, value: handleValue(parsed.data[k]) }))
    form.body = body
    form.body_raw = JSON.stringify(parsed.data)
  } else {
    form.body = []
    form.body_raw = ''
  }
  if (['POST', 'PUT', 'PATCH'].includes(form.method)) {
    form.body_type = isJson ? 2 : 1
  } else {
    form.body_type = 0
  }
  parseCurlModalRef.value && parseCurlModalRef.value.hide()
}
function showParseModal() {
  parseCurlModalRef.value && parseCurlModalRef.value.show()
}

function handleOpenAuthenticationModal() {
  addAuthModalRef.value && addAuthModalRef.value.show(form.auth)
}
function handleSaveAuthentication (list) {
  const data = (list || []).map(item => ({
    key: item.auth_key,
    value: item.auth_value,
    add_to: item.auth_value_addto
  }))
  form.auth = data
  nextTick(() => {
    inputFormRef.value?.validate()
  })
}
function handleDelAuthentication(index) {
  form.auth.splice(index, 1)
}
function onTypeChange(item) {
  item.subs = []
}
function onAddSubs(index) {
  const d = form.output[index]
  if (!Array.isArray(d.subs)) d.subs = []
  d.subs.push({ key: '', value: '', subs: [] })
}
async function save() {
  try {
    await formRef.value?.validate()
    await inputFormRef.value?.validate()
  } catch (e) {
    return
  }
  const headers = Array.isArray(form.headers) ? JSON.parse(JSON.stringify(form.headers)) : []
  const params = Array.isArray(form.params) ? JSON.parse(JSON.stringify(form.params)) : []
  const body = Array.isArray(form.body) ? JSON.parse(JSON.stringify(form.body)) : []
  const bodyRaw = String(form.body_raw || '')
  const httpAuth = (Array.isArray(form.auth) ? form.auth : [])
    .filter(ap => String(ap.key || '').trim())
    .map(ap => ({ key: String(ap.key || ''), value: String(ap.value || ''), add_to: String(ap.add_to || 'HEADERS') }))
  const dataRawObj = {
    curl: {
      method: form.method || 'GET',
      rawurl: String(form.url || '').replace(/[`"]/g, '').trim(),
      headers,
      params,
      type: form.body_type ?? 0,
      body,
      body_raw: bodyRaw,
      output: Array.isArray(form.output) ? JSON.parse(JSON.stringify(form.output)) : [],
      http_auth: httpAuth,
      timeout: Number(form.timeout || 30)
    }
  }
  const payload = {
    http_tool_id: props.toolId,
    node_name_en: form.node_name_en || '',
    node_name: form.title || '',
    node_description: form.description || '',
    node_remark: '',
    data_raw: JSON.stringify(dataRawObj)
  }
  const doSave = () => {
    saving.value = true
    const req = editingId.value ? saveHttpToolItem({ id: editingId.value, ...payload }) : saveHttpToolItem(payload)
    return req.then(() => {
      message.success('操作成功')
      emit('ok')
      visible.value = false
    }).finally(() => {
      saving.value = false
    })
  }
  if (editingId.value) {
    Modal.confirm({
      title: '确认保存',
      content: '修改工具后使用到该工具的工作流可能无法正常运行。确认保存编辑么？',
      okText: '确认',
      cancelText: '取消',
      onOk() {
        return doSave()
      }
    })
  } else {
    doSave()
  }
}

function del(item) {
  delHttpToolItem(item.id).then(() => {
    message.success('操作成功')
    emit('ok')
    visible.value = false
  }).finally(() => saving.value = false)
}

defineExpose({
  show,
  del
})
</script>

<style scoped lang="less">
.gray-block {
  margin-top: 24px;
}
.gray-block-title {
  display: flex;
  gap: 8px;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
  margin-bottom: 16px;
}
.gray-block-title-left {
  display: flex;
  gap: 8px;
  align-items: center;
}
.form-box {
  border-radius: 6px;
  background: #F2F4F7;
  padding: 16px;
}
.output-box {
  border-radius: 6px;
  background: #F2F4F7;
  padding: 16px;
}
.flex-block-item {
  display: flex;
  align-items: center;
  gap: 8px;
}
.array-form-box {
  margin-top: 16px;
}
.ant-form-item {
  margin-bottom: 16px;
}
.timeout-form-item .ant-form-item {
  margin-bottom: 0;
}
.form-item-list {
  margin-bottom: 8px;
}
.btn-hover-wrap {
  cursor: pointer;
}
.flex1 {
  flex: 1;
}
.output-block {
  display: flex;
  gap: 12px;
}
.output-item {
  font-size: 12px;
  color: #8c8c8c;
}
.output-box .flex-block-item .btn-hover-wrap {
  width: 24px;
  height: 24px;
}
.body-form-item {
  margin-bottom: 0;
}
.mt12 {
  margin-top: 12px;
}
.form-tips {
  align-self: stretch;
  color: #8c8c8c;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
}
</style>
