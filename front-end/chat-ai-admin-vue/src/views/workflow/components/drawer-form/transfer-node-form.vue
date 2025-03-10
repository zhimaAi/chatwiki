<template>
  <div class="form-box">
    <a-form ref="formRef" :model="formState" :labelCol="labelCol" :wrapperCol="wrapperCol">
      <div class="title-block">基础信息</div>
      <a-form-item
        label="节点名称"
        name="node_name"
        :rules="[{ required: true, message: '请输入节点名称' }]"
      >
        <a-input
          v-model:value="formState.node_name"
          :maxLength="15"
          placeholder="请输入节点名称，最多15个字"
        ></a-input>
      </a-form-item>

      <div class="title-block mt24">转人工设置</div>
      <a-form-item class="mb2" label="转接客服">
        <a-radio-group v-model:value="formState.switch_type">
          <a-radio :value="1">自动分配</a-radio>
          <a-radio :value="2">指定客服</a-radio>
          <a-radio :value="3">指定客服组</a-radio>
        </a-radio-group>
      </a-form-item>
      <a-form-item
        v-if="formState.switch_type == 2"
        class="hide-label mb2"
        :colon="false"
        :wrapperCol="{ span: 24 }"
        name="staffs"
        :rules="[{ required: true, message: '请选择指定客服' }]"
      >
        <template #label></template>
        <div class="slect-checkbox-wrap">
          <div class="selct-all-block">
            <a-form-item-rest>
              <a-checkbox v-model:checked="state.staffCheckAll" @change="onCheckAllChange('staff')">
                全选
              </a-checkbox>
            </a-form-item-rest>
          </div>
          <a-checkbox-group @change="onCheckboxChange" v-model:value="formState.staffs">
            <div class="check-list-box">
              <div class="check-item" v-for="item in staffLists">
                <a-checkbox :value="item.user_id">
                  <span class="check-label-item">{{ item.user_name }}</span>
                </a-checkbox>
              </div>
            </div>
          </a-checkbox-group>
        </div>
      </a-form-item>
      <a-form-item
        v-if="formState.switch_type == 3"
        class="hide-label mb2"
        :colon="false"
        :wrapperCol="{ span: 24 }"
        name="staff_group"
        :rules="[{ required: true, message: '请选择指定客服组' }]"
      >
        <template #label></template>
        <div class="slect-checkbox-wrap">
          <div class="selct-all-block">
            <a-form-item-rest>
              <a-checkbox v-model:checked="state.groupCheckAll" @change="onCheckAllChange('group')">
                全选
              </a-checkbox>
            </a-form-item-rest>
          </div>
          <a-checkbox-group @change="onCheckboxChange" v-model:value="formState.staff_group">
            <div class="check-list-box">
              <div class="check-item" v-for="item in KefuGroupList">
                <a-checkbox :value="item._id">
                  <span class="check-label-item">{{ item.group_name }}</span>
                </a-checkbox>
              </div>
            </div>
          </a-checkbox-group>
        </div>
      </a-form-item>
      <div class="title-block mt24">转人工提示语</div>
      <a-form-item class="mb8" label="消息内容" name="switch_content" required>
        <a-textarea
          v-model:value="formState.switch_content"
          placeholder="请输入消息内容"
          :auto-size="{ minRows: 3, maxRows: 5 }"
        />
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup>
import { ref, reactive, inject, watch, onMounted, computed } from 'vue'
import { useUserStore } from '@/stores/modules/user'
import { message, Modal } from 'ant-design-vue'
import {
  CloseCircleFilled,
  CloseCircleOutlined,
  LoadingOutlined,
  PlusOutlined,
  EditOutlined,
} from '@ant-design/icons-vue'
import { getStaffList, getKefuGroupList } from '@/api/robot/robot.js'
import { useRoute } from 'vue-router'
const query = useRoute().query
const userStore = useUserStore()
const { updateNodeItem, updateModifyNum } = inject('nodeInfo')

const props = defineProps({
  properties: {
    type: Object,
    default: () => ({}),
  },
})

const emit = defineEmits(['ok'])
const formRef = ref()

const labelCol = {
  span: 4,
}
const wrapperCol = {
  span: 20,
}

const formState = reactive({
  node_name: '',
  switch_type: 1,
  staffs: [],
  staff_group: [],
  switch_content: '',
})

const state = reactive({
  staffCheckAll: false,
  groupCheckAll: false,
})

const staffLists = ref([])
const KefuGroupList = ref([])

const onCheckAllChange = (key) => {
  if (key == 'staff') {
    if (state.staffCheckAll) {
      formState.staffs = staffLists.value.map((item) => item.user_id)
    } else {
      formState.staffs = []
    }
  } else {
    if (state.groupCheckAll) {
      formState.staff_group = KefuGroupList.value.map((item) => item._id)
    } else {
      formState.staff_group = []
    }
  }
}
let updateNum = 0
watch(
  () => props.properties,
  (val) => {
    try {
      formState.node_name = val.node_name
      formState.switch_content = val.switch_content
      formState.switch_type = +val.switch_type || 1
      let switch_staff = val.switch_staff ? val.switch_staff.split(',') : []
      state.staffCheckAll = false
      state.groupCheckAll = false
      formState.staffs = []
      formState.staff_group = []

      if (switch_staff.length) {
        if (formState.switch_type == 2) {
          formState.staffs = switch_staff
          state.staffCheckAll =
            switch_staff.length > 0 && switch_staff.length == staffLists.value.length
        }
        if (formState.switch_type == 3) {
          formState.staff_group = switch_staff
          state.groupCheckAll =
            switch_staff.length > 0 && switch_staff.length == KefuGroupList.value.length
        }
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

watch(
  () => formState.switch_type,
  () => {
    state.staffCheckAll = formState.staffs.length == staffLists.value.length
    state.groupCheckAll = formState.staff_group.length == KefuGroupList.value.length
  },
)

const onCheckboxChange = () => {
  state.staffCheckAll = formState.staffs.length == staffLists.value.length
  state.groupCheckAll = formState.staff_group.length == KefuGroupList.value.length
}

const checkedDelay = (rule, value) => {
  if (value == null) {
    return Promise.reject('请输入延迟发送时间')
  }
  return Promise.resolve()
}

const saveForm = () => {
  // 保存更新表单数据
  let updateInfo = {
    node_name: formState.node_name,
    switch_type: formState.switch_type,
    switch_content: formState.switch_content,
  }
  let switch_staff = ''
  if (formState.switch_type == 2) {
    // 客服
    switch_staff = formState.staffs.join(',')
  }
  if (formState.switch_type == 3) {
    // 客服组
    switch_staff = formState.staff_group.join(',')
  }
  updateInfo.switch_staff = switch_staff
  updateNodeItem({ ...updateInfo })
}

const onSave = () => {
  formRef.value
    .validate()
    .then(() => {
      saveForm()
    })
    .catch((err) => {
      console.log('error', err)
    })
}

const wechatInfo = computed(() => {
  let result = {
    wechatapp_id: '',
    channel_id: '',
  }
  if (userStore.myAppList.length) {
    result.wechatapp_id = userStore.myAppList[0].wechatapp_id
    if (userStore.myAppList[0].channels && userStore.myAppList[0].channels.length) {
      result.channel_id = userStore.myAppList[0].channels[0]._id
    }
  }
  return result
})

onMounted(() => {
  getStaffList({ wechatapp_id: wechatInfo.value.wechatapp_id, limit: 9999 }).then((res) => {
    staffLists.value = res.data.staffList
    state.staffCheckAll =
      formState.staffs.length > 0 && formState.staffs.length == staffLists.value.length
  })
  getKefuGroupList({ wechatapp_id: wechatInfo.value.wechatapp_id }).then((res) => {
    KefuGroupList.value = res.data || []
    state.groupCheckAll =
      formState.staff_group.length > 0 && formState.staff_group.length == KefuGroupList.value.length
  })
})

defineExpose({
  onSave,
})
</script>

<style lang="less" scoped>
@import './common.less';

.message-cards-box {
  display: flex;
  flex-direction: column;
  gap: 16px;
  .message-card-item {
    background: #f2f4f7;
    border-radius: 6px;
    .message-title {
      display: flex;
      align-items: center;
      justify-content: space-between;
      height: 32px;
      padding: 0 16px;
      color: #000000;
      border-bottom: 1px solid #e4e6eb;
    }
    .message-body {
      padding: 16px;
      padding-bottom: 0;
      ::v-deep(.ant-form-item) {
        margin-bottom: 16px;
      }
      ::v-deep(.mb8.ant-form-item) {
        margin-bottom: 8px;
      }
    }
  }
  .question-guide-box {
    display: flex;
    flex-direction: column;
    gap: 8px;
    .question-guide-item {
      display: flex;
      align-items: center;
      gap: 8px;
      .input-box {
        flex: 1;
      }
    }
  }
}
.gray-block {
  background: #f2f4f7;
  padding: 12px 16px;
  border-radius: 6px;
}

.slect-checkbox-wrap {
  .selct-all-block {
    padding-bottom: 8px;
    border-bottom: 1px solid #d9d9d9;
  }
  .check-list-box {
    display: flex;
    flex-wrap: wrap;
    gap: 4px 0;

    .check-item {
      width: 25%;
      height: 26px;
    }
    .check-label-item {
      display: inline-block;
      width: 85px;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      margin-top: 4px;
    }
  }
}
</style>
