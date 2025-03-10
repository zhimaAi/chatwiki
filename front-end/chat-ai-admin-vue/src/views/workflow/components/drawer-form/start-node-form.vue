<template>
  <div class="form-box">
    <a-form :labelCol="labelCol" :wrapperCol="wrapperCol">
      <div class="title-block">目标用户</div>
      <a-form-item ref="name" label="生效应用或渠道">
        <!-- <a-button type="dashed">
          <template #icon>
            <PlusOutlined />
          </template>
          选择应用或渠道
        </a-button> -->
        <a-cascader
          v-model:value="formState.wechatapp_channel"
          multiple
          @change="onTest"
          :show-search="{ filter }"
          :options="appsOptions"
          placeholder="请选择应用或渠道"
        />
        <div class="form-tip">同一个应用或渠道，只能在一个机器人中</div>
      </a-form-item>
      <a-form-item ref="name" label="目标用户">
        <a-radio-group v-model:value="formState.target_user_type">
          <a-radio :value="0">所有访客</a-radio>
          <a-radio :value="1"> </a-radio>
        </a-radio-group>
        <a-select v-model:value="formState.target_user_select" style="width: 114px">
          <a-select-option
            :value="option.value"
            v-for="option in targetUserOptions"
            :key="option.value"
            >{{ option.label }}</a-select-option
          >
        </a-select>
        <span class="ml4">未留资访客</span>
      </a-form-item>
      <div class="title-block">生效时间</div>
      <a-form-item ref="name" label="生效时间">
        <a-radio-group v-model:value="formState.effective_time_type">
          <a-radio :value="1">始终生效</a-radio>
          <a-radio :value="2">指定时间生效</a-radio>
        </a-radio-group>
      </a-form-item>
      <a-form-item
        v-if="formState.effective_time_type == 2"
        class="hide-label mb2"
        v-bind="validateInfos.tableData"
        :colon="false"
        :wrapperCol="{ span: 24 }"
      >
        <div class="time-list-box">
          <a-table :data-source="formState.tableData" :pagination="false">
            <a-table-column key="time_type" title="生效日期" data-index="time_type">
              <template #default="{ record }">
                <div v-if="record.time_type == 1">每周：{{ record.week_numbers }}</div>
                <div v-if="record.time_type == 2">
                  <div>指定日期：</div>
                  {{ record.date_start }} - {{ record.date_end }}
                </div>
              </template>
            </a-table-column>
            <a-table-column key="timeRang" title="自动回复时间段" data-index="timeRang">
              <template #default="{ record }">
                <div v-for="(item, index) in record.time_list" :key="index">
                  {{ item.start }} - {{ item.end }}
                </div>
              </template>
            </a-table-column>
            <a-table-column key="action" title="操作" data-index="action" :width="88">
              <template #default="{ record, index }">
                <a-flex :gap="8">
                  <EditOutlined @click="onOpenAddTimeModal(record, index)" />
                  <a-popconfirm
                    title="确认删除该条数据?"
                    ok-text="确定"
                    cancel-text="取消"
                    @confirm="onDelTimes(index)"
                  >
                    <CloseCircleOutlined />
                  </a-popconfirm>
                </a-flex>
              </template>
            </a-table-column>
          </a-table>
          <div class="mt8">
            <a-button @click="onOpenAddTimeModal({})" type="dashed" block>
              <template #icon>
                <PlusOutlined />
              </template>
              添加时间段
            </a-button>
          </div>
        </div>
      </a-form-item>
    </a-form>
    <AddTimes :tableData="formState.tableData" ref="addTimesRef" @ok="onSaveTimes" />
  </div>
</template>

<script setup>
import { getTargetUserOptions } from '../util.js'
import { ref, reactive, inject, watch } from 'vue'
import { Form, message, Modal } from 'ant-design-vue'
import AddTimes from './components/add-times.vue'
import {
  CloseCircleFilled,
  CloseCircleOutlined,
  LoadingOutlined,
  PlusOutlined,
  EditOutlined,
} from '@ant-design/icons-vue'
import { getWechatappChannel } from '@/api/robot/robot.js'
import { useRoute } from 'vue-router'
const query = useRoute().query

const { updateNodeItem, updateModifyNum } = inject('nodeInfo')

const props = defineProps({
  properties: {
    type: Object,
    default: () => ({}),
  },
})

const emit = defineEmits(['ok'])

const targetUserOptions = getTargetUserOptions()

const useForm = Form.useForm

const labelCol = {
  span: 6,
}
const wrapperCol = {
  span: 18,
}

const saveLoading = ref(false)

const formState = reactive({
  wechatapp_channel: [],
  target_user_type: 0,
  target_user_select: 1,
  effective_time_type: 1,
  effective_time: [],
  tableData: [],
})
const appsOptions = ref([])
let updateNum = 0
watch(
  () => props.properties,
  (val) => {
    try {
      val = JSON.parse(JSON.stringify(val))
      let wechatapp_channel = val.wechatapp_channel ? val.wechatapp_channel.split(',') : []
      formState.wechatapp_channel = []
      wechatapp_channel.forEach((item) => {
        let wechatItems = item.split('-')
        formState.wechatapp_channel.push(wechatItems)
      })

      let target_user = +val.target_user || 0
      if (target_user == 0) {
        formState.target_user_type = 0
      } else {
        formState.target_user_type = 1
        formState.target_user_select = target_user
      }

      let effective_time = val.effective_time || []
      if (effective_time.length) {
        formState.effective_time_type = 2
        formState.tableData = effective_time
      } else {
        formState.effective_time_type = 1
      }
      updateNum = 0
    } catch (error) {
      console.log(error)
    }
  },
  { immediate: true, deep: true },
)

watch(
  () => formState,
  () => {
    updateNum++
    updateModifyNum(updateNum)
  },
  {
    deep: true,
  },
)

const rules = reactive({
  tableData: [
    {
      validator: async (rule, value) => {
        if (formState.effective_time_type == 2) {
          if (value.length == 0) {
            return Promise.reject('至少需要设置一条时间段')
          }
        }
        return Promise.resolve()
      },
    },
  ],
  options: [
    {
      required: true,
      validator: async (rule, value) => {
        if (formState.field_type != 1) {
          let options = formState.options.map((item) => item.value).join('')
          if (options == '') {
            return Promise.reject('请至少输入一个选项')
          }
        }
        return Promise.resolve()
      },
    },
  ],
})

const { validate, validateInfos, resetFields } = useForm(formState, rules)

const saveForm = () => {
  let updateInfo = {}
  let wechatapp_channel = []
  formState.wechatapp_channel.forEach((item) => {
    if (item.length == 1) {
      // 判断一下这个有没有渠道 如果有的话 说明全选了
      let currentApp = appsOptions.value.filter((it) => it.value == item[0])
      if (currentApp.length > 0) {
        if (currentApp[0].children && currentApp[0].children.length) {
          // 确实是有渠道的
          currentApp[0].children.forEach((it) => {
            wechatapp_channel.push([currentApp[0].value, it.value])
          })
        } else {
          wechatapp_channel.push(item)
        }
      } else {
        wechatapp_channel.push(item)
      }
    } else {
      wechatapp_channel.push(item)
    }
  })
  wechatapp_channel = wechatapp_channel.map((item) => item.join('-'))
  updateInfo.wechatapp_channel = wechatapp_channel.join(',')
  if (formState.target_user_type == 0) {
    updateInfo.target_user = 0
  } else {
    updateInfo.target_user = formState.target_user_select
  }
  let effective_time = []
  if (formState.effective_time_type == 2) {
    formState.tableData.forEach((item) => {
      let timeItem = {
        time_type: item.time_type,
        week_numbers: item.week_numbers,
        date_start: item.date_start,
        date_end: item.date_end,
        time_list: item.time_list,
      }
      effective_time.push(timeItem)
    })
  }
  updateInfo.effective_time = effective_time
  updateNodeItem({ ...updateInfo })
}

const onSave = () => {
  validate()
    .then(() => {
      saveForm()
    })
    .catch((err) => {
      console.log('error', err)
    })
}

const fetchWechatappChannel = () => {
  getWechatappChannel().then((res) => {
    let apps = res.data.apps
    appsOptions.value = []
    apps.forEach((item) => {
      let children = null
      if (item.channels && item.channels.length > 0) {
        children = item.channels.map((it) => {
          return {
            ...it,
            label: it.channel_name,
            value: it.channel_id + '',
          }
        })
      }
      appsOptions.value.push({
        ...item,
        label: item.app_name,
        value: item.wechatapp_id,
        children,
      })
    })
  })
}

fetchWechatappChannel()
const filter = (inputValue, path) => {
  return path.some((option) => option.label.toLowerCase().indexOf(inputValue.toLowerCase()) > -1)
}

const onTest = () => {
  // console.log(formState.wechatapp_channel)
}

const addTimesRef = ref(null)
let timeIndex = null
const onOpenAddTimeModal = (data, index) => {
  timeIndex = index
  addTimesRef.value.onShow(JSON.parse(JSON.stringify(data)), index)
}

const onSaveTimes = (data) => {
  if (timeIndex >= 0) {
    // 编辑
    formState.tableData.splice(timeIndex, 1, data)
  } else {
    // 新增
    formState.tableData.push(data)
  }
}

const onDelTimes = (index) => {
  formState.tableData.splice(index, 1)
}

defineExpose({
  onSave,
})
</script>

<style lang="less" scoped>
@import './common.less';
.time-list-box {
  margin-top: -16px;
}
</style>
