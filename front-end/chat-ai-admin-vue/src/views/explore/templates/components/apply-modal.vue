<template>
  <a-modal
    v-model:open="visible"
    title="模板审核结果"
    width="730px"
  >
    <a-table
        :loading="loading"
        :data-source="list"
        :columns="columns"
        :pagination="false"
    >
      <template #bodyCell="{column, record}">
        <template v-if="column.dataIndex === 'name'">
          <div class="app-item">
            <img class="avatar" src="https://xkf-upload.oss-cn-hangzhou.aliyuncs.com/dev/p/chat_wiki_plugin/1/0/202511/45/6b83aca346be5415bdfd349bb49111.jpg"/>
            <span>小红书账号拆解神器</span>
          </div>
        </template>
        <template v-else-if="column.dataIndex === 'status'">
          <span class="status-tag ok"><CheckCircleFilled/> 审核通过</span>
          <span class="status-tag fail"><CloseCircleFilled /> 审核拒绝</span>
          <span class="status-tag waiting"><ExclamationCircleFilled/> 审核中</span>
          <a-tooltip>
            <div class="tip-info zm-line1">（提示词不合规,提示词不合规提示词不合规提示词不合规</div>
          </a-tooltip>
        </template>
        <template v-else-if="column.dataIndex === 'operate'">
          <a-button type="link" class="pd0">撤回</a-button>
        </template>
      </template>
    </a-table>
  </a-modal>
</template>

<script setup>
import {ref, reactive} from 'vue';
import {CheckCircleFilled, ExclamationCircleFilled, CloseCircleFilled} from '@ant-design/icons-vue';

const visible = ref(false)
const loading = ref(false)
const list = ref([])
const columns = ref([
  {
    title: '应用名称',
    dataIndex: 'name',
    width: 200,
  },
  {
    title: '应用类型',
    dataIndex: 'type',
    width: 100,
  },
  {
    title: '提交时间',
    dataIndex: 'ctrate_time',
    width: 160,
  },
  {
    title: '审核状态',
    dataIndex: 'status',
    width: 100,
  },
  {
    title: '操作',
    dataIndex: 'operate',
    width: 100,
  },
])

for (let i = 0; i < 5; i++) {
  list.value.push({
    id: i+1,
    type: '工作流',
    ctrate_time: '22-11-28 17:07'
  })
}
</script>

<style scoped lang="less">
.app-item {
  display: flex;
  align-items: center;
  gap: 8px;

  .avatar {
    width: 20px;
    height: 20px;
    border-radius: 4px;
  }
}

.status-tag {
  display: inline-block;
  padding: 2px 6px;
  align-items: center;
  gap: 4px;
  border-radius: 6px;
  white-space: nowrap;
  font-size: 14px;
  font-weight: 500;

  &.ok {
    background: #CAFCE4;
    color: #21A665;
  }

  &.waiting {
    background: #D4E3FC;
    color: #2475FC;
  }

  &.fail {
    background: #FBDDDE;
    color: #FB363F;
  }
}

.tip-info {
  font-size: 12px;
  color: #8C8C8C;
}

.pd0 {
  padding: 0;
}
</style>
