<template>
  <div class="article-container">
    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/setting-icon.svg" class="title-icon"/> 基础信息</div>
      </div>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">公众号</div>
        </div>
        <div>
          <a-select
            v-model:value="formState.app_id"
            placeholder="请选择公众号"
            style="width: 100%;"
            @change="appChange"
          >
            <a-select-option
              v-for="app in apps"
              :key="app.app_id"
              :name="app.app_name"
              :secret="app.app_secret"
            >
              {{ app.app_name }}
            </a-select-option>
          </a-select>
        </div>
      </div>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">文章类型</div>
        </div>
        <div>
          <a-radio-group v-model:value="formState.article_type" @change="articleTypeChange">
            <a-radio value="news">图文消息</a-radio>
            <a-radio value="newspic">图片消息</a-radio>
          </a-radio-group>
        </div>
      </div>
    </div>
    <div class="node-options article-form-box">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/input.svg" class="title-icon"/>输入</div>
        <a href="https://developers.weixin.qq.com/doc/subscription/api/draftbox/draftmanage/api_draft_add.html"
           style="font-weight: 400;"
           target="_blank">查看帮助文档</a>
      </div>
      <div v-if="actionName === 'add_draft'" class="options-item">
        <div class="news-tabs">
          <div
            v-for="(article, i) in formState.articles"
            :key="i"
            :class="['news-item', {active: i == currentNewsIndex}]"
            @click="articleChange(article, i)"
          >
            {{ i == 0 ? '主图文' : `副图文${i}` }}
            <CloseOutlined v-if="formState.articles.length > 1" class="del-icon" @click.stop="delNews(i)"/>
          </div>
          <div v-if="formState.articles.length < 8" class="news-item" @click="addNews">
            <PlusOutlined/>
          </div>
        </div>
      </div>
      <template v-else>
        <div class="options-item baseline is-required">
          <div class="options-item-tit">
            <div class="option-label">文章ID</div>
          </div>
          <div style="width: 100%;">
            <AtInput
              type="textarea"
              inputStyle="height: 33px;"
              :options="variableOptions"
              :defaultSelectedList="formState.tag_map?.media_id || []"
              :defaultValue="formState.media_id"
              ref="atInputRef"
              @open="emit('updateVar')"
              @change="(text, selectedList) => changeStateFieldValue('media_id', text, selectedList)"
              placeholder="请输入内容，键入“/”可以插入变量"
            >
              <template #option="{ label, payload }">
                <div class="field-list-item">
                  <div class="field-label">{{ label }}</div>
                  <div class="field-type">{{ payload.typ }}</div>
                </div>
              </template>
            </AtInput>
            <div class="desc">文章media_id</div>
          </div>
        </div>
        <div class="options-item baseline is-required">
          <div class="options-item-tit">
            <div class="option-label">文章位置</div>
          </div>
          <div>
            <AtInput
              type="textarea"
              inputStyle="height: 33px;"
              :options="variableOptions"
              :defaultSelectedList="formState.tag_map?.index || []"
              :defaultValue="formState.index"
              ref="atInputRef"
              @open="emit('updateVar')"
              @change="(text, selectedList) => changeStateFieldValue('index', text, selectedList)"
              placeholder="请输入内容，键入“/”可以插入变量"
            >
              <template #option="{ label, payload }">
                <div class="field-list-item">
                  <div class="field-label">{{ label }}</div>
                  <div class="field-type">{{ payload.typ }}</div>
                </div>
              </template>
            </AtInput>
            <div class="desc">要更新的文章在图文消息中的位置（多图文消息时，此字段才有意义），第一篇为0（输入0~7的正整数）</div>
          </div>
        </div>
      </template>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">标题</div>
        </div>
        <AtInput
          type="textarea"
          inputStyle="height: 33px;"
          :options="variableOptions"
          :defaultSelectedList="currentNews.tag_map?.title || []"
          :defaultValue="currentNews.title"
          ref="atInputRef"
          @open="emit('updateVar')"
          @change="(text, selectedList) => changeFieldValue('title', text, selectedList)"
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
      <template v-if="currentNews.article_type === 'news'">
        <div class="options-item">
          <div class="options-item-tit">
            <div class="option-label">作者</div>
          </div>
          <AtInput
            type="textarea"
            inputStyle="height: 33px;"
            :options="variableOptions"
            :defaultSelectedList="currentNews.tag_map?.author || []"
            :defaultValue="currentNews.author"
            ref="atInputRef"
            @open="emit('updateVar')"
            @change="(text, selectedList) => changeFieldValue('author', text, selectedList)"
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
        <div class="options-item">
          <div class="options-item-tit">
            <div class="option-label">摘要</div>
          </div>
          <AtInput
            type="textarea"
            inputStyle="height: 33px;"
            :options="variableOptions"
            :defaultSelectedList="currentNews.tag_map?.digest || []"
            :defaultValue="currentNews.digest"
            ref="atInputRef"
            @open="emit('updateVar')"
            @change="(text, selectedList) => changeFieldValue('digest', text, selectedList)"
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
      </template>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">正文内容</div>
        </div>
        <div class="cont-box">
          <AtInput
            type="textarea"
            inputStyle="height: 33px;"
            :options="variableOptions"
            :defaultSelectedList="currentNews.tag_map?.content || []"
            :defaultValue="currentNews.content"
            ref="atInputRef"
            @open="emit('updateVar')"
            @change="(text, selectedList) => changeFieldValue('content', text, selectedList)"
            placeholder="请输入内容，键入“/”可以插入变量"
          >
            <template #option="{ label, payload }">
              <div class="field-list-item">
                <div class="field-label">{{ label }}</div>
                <div class="field-type">{{ payload.typ }}</div>
              </div>
            </template>
          </AtInput>
          <a-tooltip title="编辑窗口">
            <FullscreenOutlined class="edit-icon" @click="showContModal"/>
          </a-tooltip>
        </div>
      </div>
      <template v-if="currentNews.article_type === 'news'">
        <div class="options-item">
          <div class="options-item-tit">
            <div class="option-label">原文地址</div>
          </div>
          <AtInput
            type="textarea"
            inputStyle="height: 33px;"
            :options="variableOptions"
            :defaultSelectedList="currentNews.tag_map?.content_source_url || []"
            :defaultValue="currentNews.content_source_url"
            ref="atInputRef"
            @open="emit('updateVar')"
            @change="(text, selectedList) => changeFieldValue('content_source_url', text, selectedList)"
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
        <div class="options-item is-required baseline">
          <div class="options-item-tit">
            <div class="option-label">封面素材</div>
          </div>
          <div style="flex: 1;">
            <AtInput
              type="textarea"
              inputStyle="height: 33px;"
              :options="variableOptions"
              :defaultSelectedList="currentNews.tag_map?.thumb_media_id || []"
              :defaultValue="currentNews.thumb_media_id"
              ref="atInputRef"
              @open="emit('updateVar')"
              @change="(text, selectedList) => changeFieldValue('thumb_media_id', text, selectedList)"
              placeholder="请输入内容，键入“/”可以插入变量"
            >
              <template #option="{ label, payload }">
                <div class="field-list-item">
                  <div class="field-label">{{ label }}</div>
                  <div class="field-type">{{ payload.typ }}</div>
                </div>
              </template>
            </AtInput>
            <ZmRadioGroup class="mt4" v-model:value="currentNews.thumb_type" :options="thumbTypeOpts" @change="update"/>
          </div>
        </div>
      </template>
      <div class="options-item">
        <div class="options-item-tit">
          <div class="option-label">文章评论</div>
        </div>
        <ZmRadioGroup v-model:value="currentNews.need_open_comment" :options="comOpts" @change="update"/>
      </div>
      <div class="options-item">
        <div class="options-item-tit">
          <div class="option-label">评论人群</div>
        </div>
        <ZmRadioGroup v-model:value="currentNews.only_fans_can_comment" :options="fansComOpts" @change="update"/>
      </div>
      <template v-if="currentNews.article_type === 'news'">
        <!--        <div class="options-item">-->
        <!--          <div class="options-item-tit">-->
        <!--            <div class="option-label">封面裁剪坐标 (2.35:1)</div>-->
        <!--          </div>-->
        <!--          <ZmRadioGroup v-model:value="currentNews.pic_crop_235_1" :options="positionOpts"/>-->
        <!--        </div>-->
        <!--        <div class="options-item">-->
        <!--          <div class="options-item-tit">-->
        <!--            <div class="option-label">封面裁剪坐标(1:1)</div>-->
        <!--          </div>-->
        <!--          <ZmRadioGroup v-model:value="currentNews.pic_crop_1_1" :options="positionOpts"/>-->
        <!--        </div>-->
      </template>
      <template v-else>
        <div class="options-item is-required baseline">
          <div class="options-item-tit">
            <div class="option-label">图片</div>
          </div>
          <div class="img-box">
            <div v-for="(item, index) in currentNews.image_info" :key="index" class="img-item">
              <AtInput
                type="textarea"
                inputStyle="height: 33px;"
                :options="variableOptions"
                :defaultSelectedList="item.tags || []"
                :defaultValue="item.value"
                ref="atInputRef"
                @open="emit('updateVar')"
                @change="(text, selectedList) => changeValue(item, text, selectedList)"
                placeholder="请输入内容，键入“/”可以插入变量"
              >
                <template #option="{ label, payload }">
                  <div class="field-list-item">
                    <div class="field-label">{{ label }}</div>
                    <div class="field-type">{{ payload.typ }}</div>
                  </div>
                </template>
              </AtInput>
              <CloseCircleOutlined @click="delItem('image_info', index)"/>
            </div>
            <a-button
              v-if="currentNews.image_info.length < 20"
              type="dashed"
              :icon="h(PlusOutlined)"
              @click="addItem('image_info')">新增图片
            </a-button>
          </div>
        </div>
        <div class="options-item baseline">
          <div class="options-item-tit">
            <div class="option-label">封面信息</div>
          </div>
          <div class="img-box">
            <div v-for="(item, index) in currentNews.cover_info" :key="index" class="img-item">
              <AtInput
                type="textarea"
                inputStyle="height: 33px;"
                :options="variableOptions"
                :defaultSelectedList="item.tags || []"
                :defaultValue="item.value"
                ref="atInputRef"
                @open="emit('updateVar')"
                @change="(text, selectedList) => changeValue(item, text, selectedList)"
                placeholder="请输入内容，键入“/”可以插入变量"
              >
                <template #option="{ label, payload }">
                  <div class="field-list-item">
                    <div class="field-label">{{ label }}</div>
                    <div class="field-type">{{ payload.typ }}</div>
                  </div>
                </template>
              </AtInput>
              <CloseCircleOutlined @click="delItem('cover_info', index)"/>
            </div>
            <a-button
              v-if="currentNews.cover_info.length < 5"
              type="dashed"
              :icon="h(PlusOutlined)"
              @click="addItem('cover_info')">新增封面
            </a-button>
          </div>
        </div>
<!--        <div class="options-item baseline">-->
<!--          <div class="options-item-tit">-->
<!--            <div class="option-label">商品信息</div>-->
<!--          </div>-->
<!--          <div class="img-box">-->
<!--            <div v-for="(item, index) in currentNews.product_info" :key="index" class="img-item">-->
<!--              <AtInput-->
<!--                type="textarea"-->
<!--                inputStyle="height: 33px;"-->
<!--                :options="variableOptions"-->
<!--                :defaultSelectedList="item.tags || []"-->
<!--                :defaultValue="item.value"-->
<!--                ref="atInputRef"-->
<!--                @open="emit('updateVar')"-->
<!--                @change="(text, selectedList) => changeValue(item, text, selectedList)"-->
<!--                placeholder="请输入内容，键入“/”可以插入变量"-->
<!--              >-->
<!--                <template #option="{ label, payload }">-->
<!--                  <div class="field-list-item">-->
<!--                    <div class="field-label">{{ label }}</div>-->
<!--                    <div class="field-type">{{ payload.typ }}</div>-->
<!--                  </div>-->
<!--                </template>-->
<!--              </AtInput>-->
<!--              <CloseCircleOutlined @click="delItem('product_info', index)"/>-->
<!--            </div>-->
<!--            <a-button-->
<!--              v-if="currentNews.cover_info.length < 5"-->
<!--              type="dashed"-->
<!--              :icon="h(PlusOutlined)"-->
<!--              @click="addItem('product_info')">新增商品-->
<!--            </a-button>-->
<!--          </div>-->
<!--        </div>-->
      </template>
    </div>
    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/output.svg" class="title-icon"/>输出</div>
      </div>
      <div class="options-item">
        <OutputFields :tree-data="outputData"/>
      </div>
    </div>

    <DraftNewsCont
      v-if="formState.article_type === 'news'"
      ref="newsContRef"
      :variable-options="variableOptions"
      :appData="{  app_id: formState.app_id, app_secret: formState.app_secret}"
      @change="changeContent"
    />
    <DraftNewspicCont
      v-else
      ref="newsContRef"
      :variable-options="variableOptions"
      @change="changeContent"
    />
  </div>
</template>

<script setup>
import {ref, reactive, onMounted, watch, inject, computed, h, nextTick} from 'vue';
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import {getWechatAppList} from "@/api/robot/index.js";
import {pluginOutputToTree} from "@/constants/plugin.js";
import {jsonDecode} from "@/utils/index.js";
import OutputFields from "@/views/workflow/components/feishu-table/output-fields.vue";
import {PlusOutlined, CloseOutlined, FullscreenOutlined, CloseCircleOutlined, QuestionCircleOutlined} from '@ant-design/icons-vue';
import ZmRadioGroup from "@/components/common/zm-radio-group.vue";
import DraftNewsCont from "./components/draft-news-cont.vue";
import DraftNewspicCont from "./components/draft-newspic-cont.vue";

const emit = defineEmits(['updateVar'])
const props = defineProps({
  node: {
    type: Object,
    default: () => ({})
  },
  action: {
    type: Object,
  },
  actionName: {
    type: String,
  },
  variableOptions: {
    type: Array,
  }
})

const setData = inject('setData')
const newsContRef = ref(null)
const outputData = ref([])
const apps = ref([])
const articleData = {
  title: "",
  article_type: "",
  author: "",
  digest: "",
  content: "",
  content_source_url: "",
  need_open_comment: 0,
  only_fans_can_comment: 0,
  pic_crop_235_1: '',
  pic_crop_1_1: '',
  thumb_type: 1,
  thumb_media_id: "",
  image_info: [{value: '', tags: []}],
  product_info: [{value: '', tags: []}],
  cover_info: [{value: '', tags: []}],
  tag_map: {},
}
const formState = reactive({
  app_id: undefined,
  app_secret: "",
  app_name: "",
  article_type: "news",
  articles: [],

  // 更新文章
  media_id: "",
  index: 0,
  tag_map: {},
})
const currentNewsIndex = ref(0)
const currentNews = computed(() => {
  return formState.articles[currentNewsIndex.value] || articleData
})

const positionOpts = [
  {label: '居中', value: 1},
  {label: '左对齐', value: 2},
  {label: '右对齐', value: 3},
]

const thumbTypeOpts = [
  {label: '素材url', value: 1},
  {label: '素材MedialD', value: 2},
]

const comOpts = [
  {label: '关闭', value: 0},
  {label: '打开', value: 1},
]

const fansComOpts = [
  {label: '所有人可评论', value: 0},
  {label: '粉丝才可评论', value: 1},
]

onMounted(() => {
  init()
})

watch(() => props.action, () => {
  outputData.value = pluginOutputToTree(JSON.parse(JSON.stringify(props.action.output || '{}')))
}, {
  deep: true,
  immediate: true
})

function init() {
  loadWxApps()
  addNews()
  nodeParamsAssign()
}

function loadWxApps() {
  getWechatAppList({app_type: 'official_account'}).then(res => {
    apps.value = res?.data || []
  })
}

function nodeParamsAssign() {
  let nodeParams = JSON.parse(props.node.node_params)
  let arg = nodeParams?.plugin?.params?.arguments || {}
  // 更新文章时 articles 为对象
  if (props.actionName === 'update_draft' && arg.articles && Object.keys(arg.articles).length) {
    arg.articles = [arg.articles]
  }
  Object.assign(formState, arg)
  if (!formState.app_id) {
    let app = jsonDecode( window.localStorage.getItem('zm:ai:draft:last:app'))
    if (app) {
      Object.assign(formState, app)
      update()
    }
  }
}

function appChange(_, option) {
  const {key, secret, name} = option
  formState.app_secret = secret
  formState.app_name = name
  window.localStorage.setItem('zm:ai:draft:last:app', JSON.stringify({
    app_id: key,
    app_secret: secret,
    app_name: name
  }))
  update()
}

function articleTypeChange() {
  formState.articles = []
  addNews()
}

function articleChange(article, i) {
  currentNewsIndex.value = i
}

function changeStateFieldValue(field, text, selectedList) {
  formState[field] = text
  formState.tag_map[field] = selectedList
  update()
}

function changeFieldValue(field, text, selectedList) {
  currentNews.value[field] = text
  currentNews.value.tag_map[field] = selectedList
  update()
}

function changeValue(item, text, selectedList) {
  item.value = text
  item.tags = selectedList
  update()
}

function showContModal() {
  newsContRef.value.show(currentNews.value.content, currentNews.value.tag_map.content || [])
}


function changeContent(value,tags) {
  currentNews.value.content = value
  currentNews.value.tag_map.content = tags || []
  update()
}

function addNews() {
  let data = JSON.parse(JSON.stringify(articleData))
  data.article_type = formState.article_type
  formState.articles.push(data)
}

function delNews(index) {
  formState.articles.splice(index, 1)
  update()
}

function addItem(field) {
  formState.articles[currentNewsIndex.value][field].push({
    value: "",
    tags: []
  })
}

function delItem(field, index) {
  formState.articles[currentNewsIndex.value][field].splice(index, 1)
}

function update() {
  nextTick(() => {
    let nodeParams = JSON.parse(props.node.node_params)
    nodeParams.plugin.output_obj = outputData.value
    let arg = {...formState}
    if (props.actionName === 'update_draft') {
      arg.articles = arg.articles[0]
    }
    Object.assign(nodeParams.plugin.params, {
      arguments: {
        ...arg
      }
    })
    setData({
      ...props.node,
      node_params: JSON.stringify(nodeParams)
    })
  })
}
</script>

<style scoped lang="less">
@import "../node-options";

.article-container {
  :deep(.mention-input-warpper) {
    height: 32px;
    word-break: break-all;

    .type-textarea {
      height: 32px;
      min-height: 32px;
    }
  }
}

.news-tabs {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  color: #262626;
  font-size: 14px;
  gap: 4px;

  .news-item {
    padding: 9px 16px;
    border-radius: 2px;
    border: 1px solid #F0F0F0;
    background: #FAFAFA;
    cursor: pointer;

    &.active,
    &:hover {
      color: #2475FC;
      background: #FFF;
    }

    &:hover {
      .del-icon {
        display: inline-block;
      }
    }

    .del-icon {
      display: none;
      margin-left: 4px;
    }
  }
}

.article-form-box.node-options {
  .options-item {
    display: flex;
    flex-direction: row;

    &.baseline {
      align-items: baseline;
    }

    .options-item-tit {
      width: 83px;
      text-align: left;
      flex-shrink: 0;
    }
  }
}

.cont-box {
  position: relative;
  flex: 1;

  .edit-icon {
    padding: 4px;
    border-radius: 6px;
    background: #E4E6EB;
    position: absolute;
    right: 4px;
    top: 5px;
  }
}

.img-box {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;

  .img-item {
    display: flex;
    align-items: center;
    gap: 8px;
  }
}

.mt4 {
  margin-top: 4px;
}
</style>
