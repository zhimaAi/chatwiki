<template>
  <div class="details-library-page">
    <cu-scroll :scrollbar="false">
      <div class="list-tools">
        <div class="tools-items">
          <div class="tool-item">
            <a-dropdown :trigger="['hover']" overlayClassName="add-dropdown-btn">
              <template #overlay>
                <a-menu @click="handleMenuClick">
                  <a-menu-item :key="1">
                    <div class="dropdown-btn-menu">
                      <a-flex class="title-block" :gap="4">
                        <svg-icon name="doc-icon"></svg-icon>
                        <div class="title">本地文档</div>
                      </a-flex>
                      <div class="desc" v-if="libraryInfo.type == 2">
                        上传本地 docx/csv/xlsx 等格式文件
                      </div>
                      <div class="desc" v-else>
                        上传本地 pdf/docx/txt/md/xlsx/csv/html 等格式文件
                      </div>
                    </div>
                  </a-menu-item>
                  <a-menu-item :key="2" v-if="libraryInfo.type != 2">
                    <div class="dropdown-btn-menu">
                      <a-flex class="title-block" :gap="4">
                        <svg-icon name="link-icon"></svg-icon>
                        <div class="title">在线数据</div>
                      </a-flex>
                      <div class="desc">获取在线网页内容</div>
                    </div>
                  </a-menu-item>
                  <a-menu-item :key="3">
                    <div class="dropdown-btn-menu">
                      <a-flex class="title-block" :gap="4">
                        <svg-icon name="cu-doc-icon"></svg-icon>
                        <div class="title">自定义文档</div>
                      </a-flex>
                      <div class="desc">自定义一个空文档，手动添加内容</div>
                    </div>
                  </a-menu-item>
                </a-menu>
              </template>
              <a-button type="primary">
                <template #icon>
                  <PlusOutlined />
                </template>
                <span>添加内容</span>
              </a-button>
            </a-dropdown>
          </div>
          <!-- <div class="tool-item">
            <a-button @click="handleBatchDelete" danger>批量删除</a-button>
          </div> -->

          <!-- <a-flex align="center" class="tool-item custom-select-box" v-show="false">
            <span>嵌入模型：</span>
            <model-select
              modelType="TEXT EMBEDDING"
              :isOffline="false"
              :modeName="modelForm.use_model"
              :modeId="modelForm.model_config_id"
              style="width: 240px"
              @change="onChangeModel"
              @loaded="onVectorModelLoaded"
            />
          </a-flex> -->
          <a-flex align="center" v-if="neo4j_status" class="tool-item custom-select-box pd-5-8">
            <ModelSelect
              modelType="LLM"
              v-show="false"
              @loaded="onVectorModelLoaded"
            />
            <span>生成知识图谱：</span>
            <a-switch
              v-model:checked="createGraphSwitch"
              @change="createGraphSwitchChange"
              checked-children="开"
              un-checked-children="关"
            />
          </a-flex>
        </div>
        <div class="right-box">
          <div class="status-box" v-if="libraryInfo.type == 0">
            <div class="status-item">
              <div class="status-label">已学习：</div>
              <div class="status-content">{{count_data.learned_count}}</div>
            </div>
            <div class="status-item">
              <div class="status-label">，未学习：</div>
              <div class="status-content content-tip">{{count_data.learned_wait_count}}</div>
            </div>
            <div class="status-item">
              <div class="status-label">，学习失败：</div>
              <div class="status-content content-tip">{{count_data.learned_err_count}}</div>
            </div>
          </div>
          <div class="status-select-box" v-if="libraryInfo.type == 0">
            <a-select
              v-model:value="queryParams.status"
              style="width: 150px"
              placeholder="请选择"
              @change="handleChangeStatus"
            >
              <a-select-option
                v-for="item in statusList"
                :key="item.id"
                :value="item.id"
                >{{ item.name }}</a-select-option
              >
            </a-select>
          </div>
          <div class="tool-item">
            <a-input
              style="width: 282px"
              v-model:value="queryParams.file_name"
              placeholder="请输入文档名称搜索"
              @change="onSearch"
            >
              <template #suffix>
                <SearchOutlined @click="onSearch" style="color: rgba(0, 0, 0, 0.25)" />
              </template>
            </a-input>
          </div>
        </div>
      </div>
      <page-alert style="margin-bottom: 16px;" title="使用说明" v-if="libraryInfo.type == 0">
        <div>
          <p>
            1、如果单次上传一个文档，上传成功后，系统会自动学习；如果单次上传多个文档，上传成功后，需要手动点击文档后面"学习"进行学习；如果解析失败，支持重新学习。
          </p>
          <p>2、未学习的文档数据不会被检索到。</p>
        </div>
      </page-alert>
      <!-- :row-selection="{ selectedRowKeys: state.selectedRowKeys, onChange: onSelectChange }" -->
      <div class="list-content">
        <a-table
          :columns="columns"
          :data-source="fileList"
          :scroll="{x: 1000}"
          row-key="id"
          :pagination="{
            current: queryParams.page,
            total: queryParams.total,
            pageSize: queryParams.size,
            showQuickJumper: true,
            showSizeChanger: true,
            pageSizeOptions: ['10', '20', '50', '100']
          }"
          @change="onTableChange"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'file_name'">
              <div class="doc-name-td">
                <a-popover :title="null" v-if="record.doc_type == 2">
                  <template #content>
                    原链接：<a :href="record.doc_url" target="_blank">{{ record.doc_url }} </a>
                    <CopyOutlined
                      v-copy="`${record.doc_url}`"
                      style="margin-left: 4px; cursor: pointer"
                    />
                  </template>
                  <a @click="handlePreview(record)">
                    <span v-if="['5', '6', '7'].includes(record.status)">{{ record.doc_url }}</span>
                    <span v-else>{{ record.file_name }}</span>
                  </a>
                </a-popover>
                <a @click="handlePreview(record)" v-else>
                  <span v-if="['5', '6', '7'].includes(record.status)">{{ record.doc_url }}</span>
                  <span v-else>{{ record.file_name }}</span>
                </a>
                <div v-if="record.doc_type == 2 && record.remark" class="url-remark">
                  备注：{{ record.remark }}
                </div>
              </div>
            </template>
            <template v-if="column.key === 'graph_entity_count'">
              <a @click="toGraph(record)">{{ record.graph_entity_count }}</a>
            </template>
            <template v-if="column.key === 'status'">
              <template v-if="record.file_ext == 'pdf' && record.pdf_parse_type >= 2">
                <div class="pdf-progress-box" v-if="record.status == 0">
                  <div class="progress-title">
                    <span class="status-box"><LoadingOutlined />文档解析中</span>
                    <a @click="handleCancelOcrPdf(record)">取消</a>
                  </div>
                  <div class="progress-bar">
                    <a-progress
                      size="small"
                      class="progress-bar-box"
                      :percent="parseInt((record.ocr_pdf_index / record.ocr_pdf_total) * 100)"
                      :show-info="false"
                    />
                    <div class="num-box">
                      {{ record.ocr_pdf_index }} / {{ record.ocr_pdf_total }}
                    </div>
                  </div>
                </div>
              </template>
              <template v-else>
                <span class="status-tag running" v-if="record.status == 0"
                  ><a-spin size="small" /> 转换中</span
                >
              </template>

              <span class="status-tag running" v-if="record.status == 1"
                ><a-spin size="small" /> 学习中</span
              >

              <span class="status-tag complete" v-if="record.status == 2"
                ><CheckCircleFilled /> 学习完成</span
              >

              <a-tooltip placement="top" v-if="record.status == 3">
                <template #title>
                  <span>{{ record.errmsg }}</span>
                </template>
                <span>
                  <span class="status-tag status-error"><CloseCircleFilled /> 转换失败</span>
                  <a class="ml8" v-if="libraryInfo.type == 2" @click="handlePreview(record)"
                    >学习</a
                  >
                </span>
              </a-tooltip>
              <a-tooltip placement="top" v-if="record.status == 8">
                <template #title>
                  <span>{{ record.errmsg }}</span>
                </template>
                <span>
                  <span class="status-tag status-error"><CloseCircleFilled /> 转化异常</span>
                </span>
              </a-tooltip>
              <template v-if="record.status == 4">
                <span class="status-tag"><ClockCircleFilled /> 待学习</span>
                <a class="ml8" @click="handlePreview(record)">学习</a>
              </template>
              <template v-if="record.status == 5">
                <span class="status-tag"><ClockCircleFilled /> 待获取</span>
              </template>
              <span class="status-tag running" v-if="record.status == 6"
                ><a-spin size="small" /> 获取中</span
              >
              <a-tooltip placement="top" v-if="record.status == 7">
                <template #title>
                  <span>{{ record.errmsg }}</span>
                </template>
                <span class="status-tag error"><CloseCircleFilled /> 获取失败</span>
              </a-tooltip>
              <template v-if="record.status == 9">
                <span class="status-tag cancel"><ExclamationCircleOutlined /> 取消解析</span>
              </template>

              <span class="status-tag subning" v-if="record.status == 10">
                <a-spin size="small" />
                正在分段
              </span>
            </template>
            <template v-if="column.key === 'graph_status'">
              <!--0待生成 1排队中 2生成完成 3生成失败 4生成中 5部分成功-->
              <template v-if="record.graph_status == 0">
                <span class="status-tag"><ClockCircleFilled /> 待生成</span>
                <a class="ml8" @click="createGraphTask(record)">生成</a>
              </template>
              <span v-else-if="record.graph_status == 1" class="status-tag running"
                ><HourglassFilled /> 排队中</span
              >
              <span v-else-if="record.graph_status == 2" class="status-tag complete"
                ><CheckCircleFilled /> 生成完成</span
              >
              <template v-else-if="record.graph_status == 3">
                <span class="status-tag error"><CloseCircleFilled /> 生成失败</span>
                <a class="ml8" @click="createGraphTask(record)">生成</a>
                <a-tooltip v-if="record.graph_err_msg" :title="record.graph_err_msg">
                  <div class="zm-line1 reason-text">原因：{{ record.graph_err_msg }}</div>
                </a-tooltip>
              </template>
              <span v-else-if="record.graph_status == 4" class="status-tag running"
                ><a-spin size="small" /> 生成中</span
              >
              <template v-else-if="record.graph_status == 5">
                <span class="status-tag warning"><CheckCircleFilled /> 部分成功</span>
                <div class="reason-text">
                  失败数：{{ record.graph_err_count || 0 }}
                  <a class="ml8" @click="handlePreview(record, { graph_status: 3 })">详情</a>
                </div>
              </template>
            </template>
            <template v-if="column.key === 'file_size'">
              <span v-if="record.doc_type == 3">-</span>
              <span v-else>{{ record.file_size_str }}</span>
            </template>
            <template v-if="column.key === 'doc_auto_renew_frequency'">
              <div class="text-block" v-if="record.doc_type == 2">
                <div class="time-content-box">
                  <span v-if="record.doc_auto_renew_frequency == 1">不自动更新</span>
                  <span v-if="record.doc_auto_renew_frequency == 2">每天</span>
                  <span v-if="record.doc_auto_renew_frequency == 3">每3天</span>
                  <span v-if="record.doc_auto_renew_frequency == 4">每7天</span>
                  <span v-if="record.doc_auto_renew_frequency == 5">每30天</span>
                  <span
                    class="ml4"
                    v-if="record.doc_auto_renew_frequency > 1 && record.doc_auto_renew_minute > 0"
                  >
                    {{ convertTime(record.doc_auto_renew_minute) }}
                  </span>
                  <span v-if="record.doc_auto_renew_frequency > 1">更新</span>
                  <a class="ml4 btn-hover-block" @click="handleEditOnlineDoc(record)">修改</a>
                </div>
                <a-flex align="center">
                  最近更新: 
                  {{ record.doc_last_renew_time_desc }}
                  <a-popconfirm title="确认更新?" @confirm="handleUpdataDoc(record)">
                    <a class="ml4 btn-hover-block">更新</a>
                  </a-popconfirm>
                </a-flex>
              </div>
              <div v-else>--</div>
            </template>
            <template v-if="column.key === 'paragraph_count'">
              <span v-if="record.status == 0 || record.status == 1">-</span>
              <span v-else>{{ record.paragraph_count }}</span>
            </template>
            <template v-if="column.key === 'action'">
              <a-flex :gap="8">
                <a-tooltip>
                  <template #title>重新分段</template>
                  <SyncOutlined class="btn-hover-block" @click="toReSegmentationPage(record)" />
                </a-tooltip>
                <a-dropdown>
                  <div class="table-btn" @click.prevent>
                    <MoreOutlined class="btn-hover-block" />
                  </div>
                  <template #overlay>
                    <a-menu>
                      <a-menu-item v-if="record.doc_type == 2">
                        <div @click="handleEditOnlineDoc(record)">编辑</div>
                      </a-menu-item>
                      <a-menu-item
                        :disabled="record.status == 6 || record.status == 7 || record.status == 0"
                      >
                        <div @click="handlePreview(record)">预览</div>
                      </a-menu-item>
                      <a-menu-item>
                        <div @click="handleOpenRenameModal(record)">重命名</div>
                      </a-menu-item>
                      <a-menu-item>
                        <a-popconfirm title="确定要删除吗?" @confirm="onDelete(record)">
                          <span>删除</span>
                        </a-popconfirm>
                      </a-menu-item>
                    </a-menu>
                  </template>
                </a-dropdown>
              </a-flex>
            </template>
          </template>
        </a-table>
      </div>
    </cu-scroll>
    <a-modal
      v-model:open="addFileState.open"
      :confirm-loading="addFileState.confirmLoading"
      :maskClosable="false"
      title="上传文档"
      @ok="handleSaveFiles"
      @cancel="handleCloseFileUploadModal"
      width="740px"
    >
      <a-form class="mt24">
        <a-form-item required label="">
          <div class="upload-file-box">
            <UploadFilesInput
              :type="libraryInfo.type"
              :maxCount="20"
              v-model:value="addFileState.fileList"
              @change="onFilesChange"
            />
          </div>
        </a-form-item>
        <a-form-item v-if="existPdfFile" required label="PDF解析模式" class="select-card-main">
          <div class="select-card-box">
            <div
              v-for="item in PDF_PARSE_MODE"
              :key="item.key"
              :class="['select-card-item', { active: addFileState.pdf_parse_type == item.key }, { disabled: item.key == 4 && !ali_ocr_switch }]"
              @click.stop="pdfParseTypeChange(item)"
            >
              <svg-icon class="check-arrow" name="check-arrow-filled"></svg-icon>
              <div class="card-title">
                <svg-icon :name="item.icon" style="font-size: 16px"></svg-icon>
                {{ item.title }}
                <div class="card-switch-box" v-if="item.key == 4 && !ali_ocr_switch">未开启</div>
              </div>
              <div class="card-desc">{{ item.desc }}</div>
              <div class="card-switch" v-if="item.key == 4 && !ali_ocr_switch">
                未开启阿里云OCR
                <div class="card-switch-btn" @click.stop="onGoSwitch">去开启</div>
              </div>
            </div>
          </div>
        </a-form-item>
      </a-form>
    </a-modal>
    <QaUploadModal @ok="getData" ref="qaUploadModalRef" />
    <a-modal
      v-model:open="addUrlState.open"
      :confirm-loading="addUrlState.confirmLoading"
      :maskClosable="false"
      title="添加在线数据"
      width="746px"
      @ok="handleSaveUrl"
      @cancel="handleCloseUrlModal"
    >
      <a-form
        class="url-add-form"
        layout="vertical"
        ref="urlFormRef"
        :model="addUrlState"
        :rules="addUrlState.rules"
      >
        <a-form-item name="urls" label="网页链接">
          <a-textarea
            style="height: 120px"
            v-model:value="addUrlState.urls"
            placeholder="请输入网页链接,形式：一行标题一行网页链接"
          />
        </a-form-item>
        <a-form-item name="doc_auto_renew_frequency" label="更新频率" required>
          <a-select v-model:value="addUrlState.doc_auto_renew_frequency" style="width: 100%">
            <a-select-option :value="1">不自动更新</a-select-option>
            <a-select-option :value="2">每天</a-select-option>
            <a-select-option :value="3">每3天</a-select-option>
            <a-select-option :value="4">每7天</a-select-option>
            <a-select-option :value="5">每30天</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item
          v-if="addUrlState.doc_auto_renew_frequency > 1"
          name="doc_auto_renew_minute"
          label="更新时间"
        >
          <a-time-picker
            valueFormat="HH:mm"
            v-model:value="addUrlState.doc_auto_renew_minute"
            format="HH:mm"
          />
        </a-form-item>
      </a-form>
    </a-modal>
    <AddCustomDocument
      :libraryInfo="libraryInfo"
      @ok="onSearch"
      ref="addCustomDocumentRef"
    ></AddCustomDocument>
    <RenameModal @ok="onSearch" ref="renameModalRef" />
    <OpenGrapgModal @ok="handleOpenGrapgOk" ref="openGrapgModalRef" />
    <GuideLearningTips ref="guideLearningTipsRef" />
    <EditOnlineDoc @ok="getData" ref="editOnlineDocRef" />
  </div>
</template>

<script setup>
import { useStorage } from '@/hooks/web/useStorage'
import { reactive, ref, toRaw, onUnmounted, onMounted, computed, createVNode } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { useRoute, useRouter } from 'vue-router'
import {
  PlusOutlined,
  SearchOutlined,
  CheckCircleFilled,
  CloseCircleFilled,
  ClockCircleFilled,
  MoreOutlined,
  HourglassFilled,
  ExclamationCircleOutlined,
  SyncOutlined,
  CopyOutlined,
  LoadingOutlined,
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import {
  getLibraryFileList,
  delLibraryFile,
  addLibraryFile,
  editLibrary,
  createGraph,
  manualCrawl,
  cancelOcrPdf
} from '@/api/library'
import { formatFileSize } from '@/utils/index'
import UploadFilesInput from '../add-library/components/upload-input.vue'
import { transformUrlData } from '@/utils/validate.js'
import AddCustomDocument from './components/add-custom-document.vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import OpenGrapgModal from '@/views/library/library-details/components/open-grapg-modal.vue'
import QaUploadModal from './components/qa-upload-modal.vue'
import RenameModal from './components/rename-modal.vue'
import PageAlert from '@/components/page-alert/page-alert.vue'
import GuideLearningTips from './components/guide-learning-tips.vue'
import EditOnlineDoc from './components/edit-online-doc.vue'
import { useUserStore } from '@/stores/modules/user'
import { useLibraryStore } from '@/stores/modules/library'
import { convertTime } from '@/utils/index'
import { useCompanyStore } from '@/stores/modules/company'

const { setStorage } = useStorage('localStorage')

const libraryStore = useLibraryStore()
const { changeGraphSwitch } = libraryStore

const companyStore = useCompanyStore()
const neo4j_status = computed(() => {
  return companyStore.companyInfo?.neo4j_status == 'true'
})
const ali_ocr_switch = computed(() => {
  return companyStore.companyInfo?.ali_ocr_switch == '1'
})

const userStore = useUserStore()
const PDF_PARSE_MODE = [
  {
    key: 2,
    title: '图文OCR解析',
    icon: 'pdf-ocr',
    desc: '通过OCR文字识别提取pdf文件内容，可以兼容扫描件，但是解析速度较慢。'
  },
  {
    key: 1,
    title: '纯文本解析',
    icon: 'pdf-text',
    desc: '只提取pdf中的文字内容，如果文档为扫描件可能提取不到内容。'
  },
  {
    key: 3,
    title: '图文解析',
    icon: 'pdf-img',
    desc: '提取PDF文档中的图片和文字'
  },
  {
    key: 4,
    title: '阿里云OCR解析',
    icon: 'ali-ocr',
    desc: '通过阿里云文档智能接口解析提取图片和文字',
  }
]
const rotue = useRoute()
const router = useRouter()
const query = rotue.query

const isLoading = ref(false)
const guideLearningTipsRef = ref(null)
const openGrapgModalRef = ref(null)
const createGraphSwitch = ref(false)
const libraryInfo = ref({
  library_intro: '',
  library_name: '',
  use_model: '',
  is_offline: null,
  type: 0
})
const state = reactive({
  selectedRowKeys: []
})

const modelForm = reactive({
  use_model: '',
  model_config_id: ''
})

const onSelectChange = (selectedRowKeys) => {
  state.selectedRowKeys = selectedRowKeys
}

const onChangeModel = (val, option) => {
  let new_use_model = option.modelName
  let new_model_config_id = option.modelId

  if (fileList.value.length > 0) {
    Modal.confirm({
      title: `确定切换模型为${new_use_model}吗？`,
      content: '切换后，所有学习文档将自动重新学习。',
      onOk() {
        modelForm.use_model = new_use_model
        modelForm.model_config_id = new_model_config_id
        saveLibraryConfig()
      },
      onCancel() {}
    })
  } else {
    modelForm.use_model = new_use_model
    modelForm.model_config_id = new_model_config_id
    saveLibraryConfig()
  }
}

const saveLibraryConfig = (showSuccessTip = true, callback = null) => {
  editLibrary({
    ...toRaw(libraryInfo.value),
    use_model: modelForm.use_model,
    model_config_id: modelForm.model_config_id
  }).then(() => {
    libraryInfo.value.use_model = modelForm.use_model
    libraryInfo.value.model_config_id = modelForm.model_config_id

    typeof callback === 'function' && callback()
    if (showSuccessTip) {
      message.success('保存成功')
    }
  })
}

const vectorModelList = ref([])

const onVectorModelLoaded = (list) => {
  vectorModelList.value = list

  // setDefaultModel()
}

const setDefaultModel = () => {
  // 防止没有数据时，切换模型报错
  if (!libraryInfo.value.id) {
    setTimeout(() => {
      setDefaultModel()
    }, 100)

    return
  }

  if (vectorModelList.value.length > 0 && !libraryInfo.value.use_model) {
    // 遍历查找chatwiki模型
    let modelConfig = null
    let model = null

    for (let item of vectorModelList.value) {
      if (item.model_define === 'chatwiki') {
        modelConfig = item
        for (let child of modelConfig.children) {
          if (child.name === 'text-embedding-v2') {
            model = child
            break
          }
        }
        break
      }
    }

    if (!modelConfig) {
      modelConfig = vectorModelList.value[0]
      model = modelConfig.children[0]
    }

    if (modelConfig) {
      modelForm.use_model = model.name
      modelForm.model_config_id = model.model_config_id

      saveLibraryConfig(false)
    }
  }
}

const fileList = ref([])

const count_data = reactive({
  "learned_err_count": 0, // 失败数
  "learned_count": 0, // 成功数
  "learned_wait_count": 0 // 待学习
})

const queryParams = reactive({
  library_id: query.id,
  file_name: undefined,
  status: '', // 2:成功 3:全部失败 8:部分失败 4:待学习
  page: 1,
  size: 10,
  total: 0
})

const statusList = [
  {
    id: '',
    name: '全部'
  },
  {
    id: '4',
    name: '待学习'
  },
  {
    id: '8',
    name: '部分失败'
  },
  {
    id: '3',
    name: '全部失败'
  },
  {
    id: '2',
    name: '成功'
  }
]

const columns = ref([])
const columnsDefault = [
  {
    title: '文档名称',
    dataIndex: 'file_name',
    key: 'file_name',
    width: 300
  },

  {
    title: '知识图谱',
    dataIndex: 'graph_status',
    key: 'graph_status',
    width: 200
  },
  {
    title: '文档格式',
    dataIndex: 'file_ext',
    key: 'file_ext',
    width: 100
  },
  {
    title: '文档大小',
    dataIndex: 'file_size_str',
    key: 'file_size',
    width: 100
  },
  {
    title: '分段',
    dataIndex: 'paragraph_count',
    key: 'paragraph_count',
    width: 120
  },
  {
    title: '更新频率/时间',
    dataIndex: 'doc_auto_renew_frequency',
    key: 'doc_auto_renew_frequency',
    width: 260
  },
  {
    title: '知识图谱实体数',
    dataIndex: 'graph_entity_count',
    key: 'graph_entity_count',
    width: 160
  },
  {
    title: '文档状态',
    dataIndex: 'status',
    key: 'status',
    width: 200
  },
  {
    title: '更新时间',
    dataIndex: 'update_time',
    key: 'update_time',
    width: 150
  },
  {
    title: '操作',
    dataIndex: 'action',
    key: 'action',
    fixed: 'right',
    width: 100
  }
]

const handleChangeStatus = (item) => {
  onSearch()
}

const onTableChange = (pagination) => {
  queryParams.page = pagination.current
  queryParams.size = pagination.pageSize
  getData()
}

const onSearch = () => {
  queryParams.page = 1
  getData()
}

let isLast = false

const onDelete = ({ id }) => {
  if (fileList.value.length == 1) {
    isLast = true
  }

  delLibraryFile({ id }).then(() => {
    if (isLast && queryParams.page > 1) {
      queryParams.page--
    }

    getData()
    message.success('删除成功')
  })
}

const handleBatchDelete = () => {
  if (state.selectedRowKeys.length == 0) {
    return message.error('请选择要删除的文档')
  }
  Modal.confirm({
    title: '删除确认',
    icon: createVNode(ExclamationCircleOutlined),
    content: createVNode(
      'div',
      {
        style: 'color: #8c8c8c;'
      },
      '确定要删除选择的文档吗?'
    ),
    onOk() {
      delLibraryFile({ id: state.selectedRowKeys.join(',') }).then(() => {
        getData()
        message.success('批量删除成功')
      })
    },
    onCancel() {}
  })
}

const handlePreview = (record, params = {}) => {
  if (record.status == '4') {
    return router.push({
      path: '/library/document-segmentation',
      query: { document_id: record.id, page: queryParams.page }
    })
  }
  if (record.status == '3' && libraryInfo.value.type != 2) {
    return message.error('学习失败,不可预览')
  }
  if (record.status == '0') {
    return message.error('转换中,稍候可预览')
  }
  if (record.status == '1') {
    return message.error('学习中,不可预览')
  }
  if (record.status == '6') {
    return message.error('获取中,不可预览')
  }
  if (record.status == '7') {
    return message.error('获取失败,不可预览')
  }
  if (record.status == '10') {
    return message.error('正在分段,不可预览')
  }

  router.push({ name: 'libraryPreview', query: { id: record.id, ...params } })
}

const getData = () => {
  let params = toRaw(queryParams)
  if (params.status == 0) {
    params.status = ''
  }
  isLoading.value = true
  getLibraryFileList(params).then((res) => {
    let info = res.data.info

    if (!modelForm.use_model && modelForm.use_model != info.use_model) {
      modelForm.use_model = info.use_model
      modelForm.model_config_id = info.model_config_id
    }

    libraryInfo.value = { ...info }
    createGraphSwitch.value = info.graph_switch == 1
    if (info.graph_switch == '0' || !neo4j_status.value) {
      columns.value = columnsDefault.filter(
        (item) => !['graph_status', 'graph_entity_count'].includes(item.key)
      )
    } else {
      columns.value = columnsDefault
    }

    let list = res.data.list || []
    let countData = res.data.count_data || {}

    queryParams.total = res.data.total

    count_data.learned_count = countData.learned_count
    count_data.learned_err_count = countData.learned_err_count
    count_data.learned_wait_count = countData.learned_wait_count

    let needRefresh = false
    fileList.value = list.map((item) => {
      // , '4' 是待学习，如果加进去会一直刷新状态不会改变
      if (['1', '6', '0', '5'].includes(item.status)) {
        needRefresh = true
      }
      item.file_size_str = formatFileSize(item.file_size)
      item.update_time = dayjs(item.update_time * 1000).format('YYYY-MM-DD HH:mm')
      item.doc_last_renew_time_desc = item.doc_last_renew_time > 0 ? dayjs(item.doc_last_renew_time * 1000).format(
        'YYYY-MM-DD HH:mm'
      ) : '--'
      return item
    })
    needRefresh && timingRefreshStatus()
    !needRefresh && clearInterval(timingRefreshStatusTimer.value)
  }).finally(() => {
    isLoading.value = false
  })
}

const timingRefreshStatusTimer = ref(null)
const timingRefreshStatus = () => {
  clearInterval(timingRefreshStatusTimer.value)
  timingRefreshStatusTimer.value = setInterval(() => {
    getData()
  }, 1000 * 5)
}

const addCustomDocumentRef = ref(null)

const handleMenuClick = (e) => {
  if (vectorModelList.value.length == 0) {
    Modal.confirm({
      title: `请先到模型管理中添加嵌入模型？`,
      content:
        '知识库学习需要使用到嵌入模型，请在系统管理-模型管理中添加。推荐使用通义千问、openai或者火山引擎的嵌入模型。',
      okText: '去添加',
      onOk() {
        router.push({ path: '/user/model' })
      }
    })
    return
  }

  let { key } = e

  if (key == 1) {
    handleOpenFileUploadModal()
  }
  if (key == 2) {
    handleOpenUrlModal()
  }
  if (key == 3) {
    addCustomDocumentRef.value.add()
  }
}
const addUrlState = reactive({
  open: false,
  urls: '',
  library_id: query.id,
  doc_auto_renew_frequency: 1,
  confirmLoading: false,
  doc_auto_renew_minute: '',
  rules: {
    urls: [
      {
        message: '请输入网页地址',
        required: true
      },
      {
        validator: (_rule, value) => {
          if (transformUrlData(value) === false) {
            return Promise.reject(new Error('网页地址不合法'))
          }
          return Promise.resolve()
        }
      }
    ]
  }
})

const handleOpenUrlModal = () => {
  addUrlState.open = true
  addUrlState.confirmLoading = false
  addUrlState.urls = ''
  addUrlState.doc_auto_renew_minute = '02:00'
  addUrlState.doc_auto_renew_frequency = 1
}
const urlFormRef = ref(null)
const handleSaveUrl = () => {
  // 保存本地内容
  urlFormRef.value
    .validate()
    .then(() => {
      addUrlState.confirmLoading = true
      addLibraryFile({
        library_id: addUrlState.library_id,
        urls: JSON.stringify(transformUrlData(addUrlState.urls)),
        doc_auto_renew_frequency: addUrlState.doc_auto_renew_frequency,
        doc_auto_renew_minute: convertTime(addUrlState.doc_auto_renew_minute),
        doc_type: 2
      }).then(() => {
        addUrlState.open = false
        addUrlState.confirmLoading = false
        onSearch()
      })
    })
    .catch(() => {
      addUrlState.confirmLoading = false
    })
}

const handleCloseUrlModal = () => {
  addUrlState.open = false
}

const addFileState = reactive({
  open: false,
  fileList: [],
  confirmLoading: false,
  pdf_parse_type: 1 //1纯文本解析，2ocr解析
})

const existPdfFile = computed(
  () => addFileState.fileList.filter((i) => i.type === 'application/pdf').length > 0
)

const qaUploadModalRef = ref(null)
const handleOpenFileUploadModal = () => {
  if (libraryInfo.value.type == 2) {
    qaUploadModalRef.value.show()
    return
  }
  addFileState.fileList = []
  addFileState.open = true
}

const handleCloseFileUploadModal = () => {
  addFileState.fileList = []
}

const onFilesChange = (files) => {
  addFileState.fileList = files
}

const handleSaveFiles = () => {
  // 提交后，需要显示引导学习的弹窗提示。如果勾选了“不再显示”，以后批量上传后不再显示弹窗提示。
  if (guideLearningTipsRef.value && !userStore.getGuideLearningTips && addFileState.fileList.length > 1) {
    guideLearningTipsRef.value.show()
  }

  if (addFileState.fileList.length == 0) {
    message.error('请选择文件')
    return
  }

  addFileState.confirmLoading = true

  let formData = new FormData()

  formData.append('library_id', queryParams.library_id)
  let isTableType = false
  addFileState.fileList.forEach((file) => {
    if (file.name.includes('.xlsx') || file.name.includes('.csv')) {
      isTableType = true
    }
    formData.append('library_files', file)
  })
  if (existPdfFile.value) {
    formData.append('pdf_parse_type', addFileState.pdf_parse_type)
  }
  addLibraryFile(formData)
    .then((res) => {
      getData()
      addFileState.open = false
      addFileState.fileList = []
      addFileState.confirmLoading = false
      if (isTableType && res.data.file_ids.length == 1) {
        router.push('/library/document-segmentation?document_id=' + res.data.file_ids[0])
      }
    })
    .catch(() => {
      addFileState.confirmLoading = false
    })
}

const renameModalRef = ref(null)
const handleOpenRenameModal = (record) => {
  renameModalRef.value.show(record)
}

const toReSegmentationPage = (record) => {
  router.push({
    path: '/library/document-segmentation',
    query: {
      document_id: record.id
    }
  })
}

const createGraphTask = (record) => {
  createGraph({ id: record.id }).then(() => {
    getData()
  })
}

const createGraphSwitchChange = () => {
  if (createGraphSwitch.value) {
    createGraphSwitch.value = false
    let data = {
      graph_model_config_id: libraryInfo.value.graph_model_config_id,
      graph_use_model: libraryInfo.value.graph_use_model
    }
    if (
      (!data.graph_model_config_id || !data.graph_use_model) &&
      vectorModelList.value.length > 0
    ) {
      let modelConfig = vectorModelList.value[0]
      let model = modelConfig.children[0]
      data.graph_use_model = model.name
      data.graph_model_config_id = model.model_config_id
    }
    openGrapgModalRef.value.show(data)
  } else {
    libraryInfo.value.graph_switch = 0
    saveLibraryConfig(false, () => {
      getData()
      changeGraphSwitch(0)
    })
  }
}

const handleOpenGrapgOk = (data) => {
  createGraphSwitch.value = true
  libraryInfo.value.graph_switch = 1
  libraryInfo.value.graph_model_config_id = data.graph_model_config_id
  libraryInfo.value.graph_use_model = data.graph_use_model
  saveLibraryConfig(false, () => {
    getData()
    changeGraphSwitch(1)
  })
}

const pdfParseTypeChange = (item) => {
  if (item.key == 4 && !ali_ocr_switch.value) {
    return false
  }
  addFileState.pdf_parse_type = item.key
}

const onGoSwitch = () => {
  window.open('#/user/aliocr')
}

const editOnlineDocRef = ref(null)
const handleEditOnlineDoc = (record) => {
  editOnlineDocRef.value.show({
    ...record
  })
}
const handleUpdataDoc = (record) => {
  manualCrawl({
    id: record.id
  }).then((res) => {
    message.success('更新成功')
    getData()
  })
}

const handleCancelOcrPdf = (record) => {
  Modal.confirm({
    title: `取消确认？`,
    content: '确认取消该文档解析',
    okText: '确定',
    onOk() {
      cancelOcrPdf({
        id: record.id
      }).then((res) => {
        message.success('取消成功')
        getData()
      })
    }
  })
}

const toGraph = (record) => {
  setStorage('graph:autoOpenFileId', record.id)
  
  router.push({
    path: '/library/details/knowledge-graph',
    query: {
      id: query.id,
    }
  })
}

onMounted(() => {
  if (query.page) {
    queryParams.page = +query.page
  }
  getData()
})

onUnmounted(() => {
  timingRefreshStatusTimer.value && clearInterval(timingRefreshStatusTimer.value)
})
</script>

<style lang="less" scoped>
.details-library-page {
  height: 100%;
  display: flex;
  flex-direction: column;
}
.doc-name-td {
  word-break: break-all;
}
.url-remark {
  color: #8c8c8c;
  margin-top: 2px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}
.padding-0 {
  padding: 0;
}
.test-menu-icon {
  color: #fff;
}
.library-name {
  height: 38px;
  line-height: 38px;
  padding: 0 16px;
  margin-bottom: 16px;
  font-size: 14px;
  font-weight: 600;
  color: #262626;
  border-radius: 2px;
  background-color: #f2f4f7;
  display: flex;
  align-items: center;
  .anticon-edit {
    margin-left: 8px;
    color: #8c8c8c;
    cursor: pointer;
  }
}
.between-content-box {
  display: flex;
  flex: 1;
  overflow: hidden;
  .left-menu-box {
    width: 232px;
    margin-right: 24px;
  }
  .right-content-box {
    flex: 1;
    overflow: hidden;
  }
}

.menu-item {
  width: 232px;
  height: 56px;
  display: flex;
  align-items: center;
  padding: 0 24px;
  background: #f5f5f5;
  border-radius: 2px;
  margin-bottom: 16px;
  cursor: pointer;
  &.active {
    background: #e6efff;
    border: 1px solid #2475fc;
  }
  .title {
    color: #242933;
    font-size: 14px;
    font-weight: 600;
    line-height: 22px;
    margin-left: 8px;
  }
}

.list-tools {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.list-tools {
  margin-bottom: 16px;
  .tools-items {
    display: flex;
    align-items: center;
    .tool-item {
      margin-right: 16px;
    }
  }
}

.list-content {
  .text-block{
    color: #595959;
  }
  .c-gray{
    color: #8c8c8c;
  }
  .time-content-box{
    display: flex;
    color: #8c8c8c;
  }
  .btn-hover-block{
    height: 24px;
    display: flex;
    align-items: center;
    padding: 0 8px;
    cursor: pointer;
    width: fit-content;
    border-radius: 6px;
    transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
    &:hover {
      background: #E4E6EB;
    }
  }
  .status-tag {
    display: inline-block;
    height: 24px;
    line-height: 24px;
    padding: 0 6px;
    border-radius: 2px;
    font-size: 14px;
    font-weight: 500;
    text-align: center;
    color: #595959;
    background-color: #edeff2;

    &.running {
      color: #2475fc;
      background-color: #e8effc;
    }

    &.subning {
      color: #6524FC;
      background-color: #EAE0FF;
    }

    &.complete {
      color: #21a665;
      background: #e8fcf3;
    }

    &.error {
      cursor: pointer;
      color: #fb363f;
      background-color: #f5c6c8;
    }

    &.warning {
      cursor: pointer;
      background: #faebe6;
      color: #ed744a;
    }

    &.status-learning {
      color: #2475fc;
      // background-color: #e8effc;
    }

    &.status-complete {
      color: #3a4559;
      background-color: #edeff2;
    }

    &.status-error {
      cursor: pointer;
      color: #fb363f;
      // background-color: #f5c6c8;
    }
    &.warning {
      cursor: pointer;
      // background: #faebe6;
      color: #ed744a;
    }
    &.cancel{
      background: #FAE4DC;
      color: #ed744a;
    }
  }
}
.pdf-progress-box {
  .progress-title {
    display: flex;
    align-items: center;
    gap: 8px;
    white-space: nowrap;
    .status-box {
      width: fit-content;
      padding: 0 6px;
      display: flex;
      align-items: center;
      gap: 2px;
      height: 22px;
      border-radius: 6px;
      background: #e8effc;
      color: #2475fc;
      font-weight: 500;
    }
  }
  .progress-bar {
    display: flex;
    align-items: center;
    gap: 8px;
    line-height: 20px;
    .ant-progress-line {
      margin: 0;
    }
    .progress-bar-box {
      flex: 1;
    }
    .num-box {
      font-size: 12px;
      color: #8c8c8c;
    }
  }
}
//.upload-file-box {
//  padding: 30px 0;
//}
.ml8 {
  margin-left: 8px;
}
.ml4 {
  margin-left: 4px;
}
.url-add-form {
  margin-top: 24px;
}

.add-dropdown-btn.ant-dropdown {
  .ant-dropdown-menu {
    padding: 0;
    border-radius: 0;
    ::v-deep(.ant-dropdown-menu-item) {
      padding: 12px 16px;
    }
  }
}

.dropdown-btn-menu {
  .title-block {
    color: #262626;
    font-size: 14px;
    font-weight: 600;
    line-height: 22px;
  }
  .desc {
    color: #8c8c8c;
    font-size: 14px;
    line-height: 22px;
  }
}
.table-btn {
  cursor: pointer;
  &:hover {
    color: #2475fc;
  }
}
.custom-select-box {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  > span {
    white-space: nowrap;
  }
  :deep(.ant-select-selector) {
    border: none !important;
    padding-left: 0 !important;
    height: unset !important;
  }
  padding-left: 8px;
}
.pd-5-8 {
  padding: 5px 8px;
}
.reason-text {
  color: #8c8c8c;
  font-size: 12px;
  line-height: 24px;
}

.select-card-box {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
  margin-top: 8px;
  .select-card-item {
    min-height: 105px;
    width: calc(50% - 8px);
    position: relative;
    padding: 16px;
    border-radius: 6px;
    border: 1px solid #d9d9d9;
    cursor: pointer;
    .check-arrow {
      position: absolute;
      display: block;
      right: -1px;
      bottom: -1px;
      width: 24px;
      height: 24px;
      font-size: 24px;
      color: #fff;
      opacity: 0;
      transition: all 0.2s cubic-bezier(0.8, 0, 1, 1);
    }
    .card-title {
      display: flex;
      align-items: center;
      gap: 4px;
      line-height: 22px;
      margin-bottom: 4px;
      color: #262626;
      font-weight: 600;
      font-size: 14px;
    }
    .title-icon {
      margin-right: 4px;
      font-size: 16px;
    }
    .card-desc {
      line-height: 22px;
      font-size: 14px;
      color: #595959;
    }

    .card-switch {
      display: flex;
      gap: 4px;
      color: #8c8c8c;
      font-size: 14px;
      line-height: 22px;

      .card-switch-btn {
        cursor: pointer;
        color: #2475fc;

        &:hover {
          opacity: 0.8;
        }
      }
    }

    .card-switch-box {
      width: 52px;
      height: 22px;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 10px;
      border-radius: 6px;
      border: 1px solid #EDB8A6;
      background: #FFECE5;
      color: #ed744a;
      font-size: 12px;
      font-style: normal;
      font-weight: 400;
      line-height: 20px;
    }

    &.active {
      background: var(--01-, #e5efff);
      border: 2px solid #2475fc;
      .check-arrow {
        opacity: 1;
      }
      .card-title {
        color: #2475fc;
      }
    }
  }
}
.mt24 {
  margin-top: 24px;
}

.right-box {
  display: flex;
  gap: 16px;
}

.status-box {
  display: flex;
  align-items: center;

  .status-item {
    display: flex;
    align-items: center;

    .status-label {
      color: #595959;
    }

    .content-tip {
      color: red;
    }
  }
}

.select-card-main {
  :deep(.ant-row) {
    display: block;
  }
}

.subning {
  :deep(.ant-spin-dot-item) {
    background-color: #6524FC !important;
  }
}
</style>
