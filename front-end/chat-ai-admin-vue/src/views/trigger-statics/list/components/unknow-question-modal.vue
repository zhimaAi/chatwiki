<template>
    <a-modal v-model:open="open" title="未知问题机器人分布详情" :width="700" :footer="null">
        <div class="unknow-list-box">
            <a-table :data-source="list" :loading="loading" :pagination="false" :scroll="{ y: 500 }">
                <a-table-column key="index" data-index="index" title="排名" :width="100">
                    <template #default="{ index }">
                        {{ index + 1 }}
                    </template>
                </a-table-column>

                <a-table-column key="robot_name" title="机器人" :width="200">
                    <template #default="{ record }">
                        {{ record.robot_name }}
                    </template>
                </a-table-column>

                <a-table-column key="unknow_count" title="未知问题数量" :width="150">
                    <template #default="{ record }">
                        <a-flex :gap="12">
                            <span>{{ record.unknow_question_total }}</span>
                        </a-flex>
                    </template>
                </a-table-column>

                <a-table-column key="action" title="操作" :width="120">
                    <template #default="{ record }">
                        <a @click="handleViewDetail(record)">查看明细</a>
                    </template>
                </a-table-column>
            </a-table>
        </div>
    </a-modal>

</template>

<script setup>
import { statUnknowQuestionRank } from '@/api/library'
import { getRobotInfo } from '@/api/robot'
import { useRouter } from 'vue-router'
import { ref, reactive } from 'vue'

const router = useRouter()

const open = ref(false)
const loading = ref(false)
const list = ref([])
const queryData = reactive({
    begin_date_ymd: '',
    end_date_ymd: ''
})

// 模拟数据
const mockData = [
    {
        robot_name: '灵魂厨房机器人',
        unknow_count: 75
    },
    {
        robot_name: '智能客服机器人',
        unknow_count: 72
    },
    {
        robot_name: '销售助手机器人',
        unknow_count: 68
    },
    {
        robot_name: '技术支持机器人',
        unknow_count: 55
    },
    {
        robot_name: '产品咨询机器人',
        unknow_count: 48
    },
    {
        robot_name: '订单查询机器人',
        unknow_count: 42
    },
    {
        robot_name: '售后机器人',
        unknow_count: 35
    }
]

const show = (params) => {
    open.value = true
    loading.value = true
    queryData.begin_date_ymd = params.begin_date_ymd
    queryData.end_date_ymd = params.end_date_ymd

    statUnknowQuestionRank(params).then((res) => {
        list.value = res.data || []
    }).finally(() => {
        loading.value = false
    })
}

const handleViewDetail = (record) => {
    getRobotInfo({
        id: record.robot_id
    }).then(res => {
        let str = `id=${record.robot_id}&robot_key=${res.data.robot_key}&start_date=${queryData.begin_date_ymd}&end_date=${queryData.end_date_ymd}`
        window.open(`/#/robot/config/unknown_issue?${str}`)
    })
}

defineExpose({
    show
})
</script>

<style lang="less" scoped>
.unknow-list-box {
    padding-top: 16px;
}
</style>
