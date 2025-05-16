<template>
  <div ref="mymodal">
    <a-modal wrapClassName="graph-model-wrapper" :zIndex="99999" centered v-model:open="show" title="知识图谱" width="90vw" :footer="null" @ok="handleOk" @cancel="handleCancel">
      <div class="graph-model-body">
        <div class="loading-wrapper" v-if="loading"> <a-spin /></div>
        <chartBox ref="chartBoxRef" :loading="loading" @nodeClick="onNodeClick" />

        <div class="right-box">
          <EntityDetails ref="entityDetailsRef" />
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { getFileGraphInfo } from '@/api/library/index'
import { ref, reactive } from 'vue';
import chartBox from './chart-box.vue';
import EntityDetails from './entity-details.vue'

const show = ref(false);
const mymodal = ref(null)
const chartBoxRef = ref(null)
const entityDetailsRef = ref(null)
const loading = ref(false)

const state = reactive({
  file_id: '',
  data_id: '',
})

const chartData = reactive({
  edges: [],
  nodes: []
})

// 关系映射
let edgeMap = {}

const getData = () => {
  getGraphData()
}

const getGraphData= () => {
  loading.value = true
  getFileGraphInfo({
    file_id: state.file_id,
    data_id: state.data_id,
    search_term: state.search_term
  }).then(res => {
    loading.value = false

    let edges = res.data.edges || []
    let nodes = res.data.nodes || []
    // 去重逻辑（包含反向检测）
    const edgeSet = new Map();
    edges = edges.filter(item => {
      const key = `${item.from_id}_${item.to_id}`;
      const reverseKey = `${item.to_id}_${item.from_id}`;
      
      if(edgeSet.has(key) || edgeSet.has(reverseKey)) {
        return false;
      }
      edgeSet.set(key, true);
      return true;
    });

    edges.forEach(item => {
      item.id = item.id.toString()
      item.from = item.from_id.toString()
      item.to = item.to_id.toString()
      item.caption = item.label
    })

    const fromNodes = new Set(edges.map(e => e.from));
    const toNodes = new Set(edges.map(e => e.to));

    nodes.forEach(item => {
      item.caption = item.name
      item.id = item.id.toString()

      item.nodeType = 
        fromNodes.has(item.id) && toNodes.has(item.id) ? 'both' :
        fromNodes.has(item.id) ? 'from' :
        toNodes.has(item.id) ? 'to' : 'normal';

      // 根据类型设置样式
      if(item.nodeType === 'from') {
        item.size = 35;
        item.color = '#FFB1B2';
      } else if(item.nodeType === 'to') {
        item.size = 45; 
        item.color = '#005AFF';
      } else {
        item.size = 30;
        item.color = '#69E9BE';
      }
    });

    let data = {
      edges,
      nodes
    }
    
    chartData.edges = data.edges
    chartData.nodes = data.nodes

    setTimeout(() => {
      chartBoxRef.value.render(data)
    }, 0);

    edgeMap = {}
    // 提前生成关系映射
    edges.forEach((edge) => {
      if (!edgeMap[edge.to]) {
        edgeMap[edge.to] = [];
      }
      
      edgeMap[edge.to].push(edge.from);
    });

  }).catch(err => {
    console.log(err)
    loading.value = false
  })
}

const onNodeClick = (node) => {
  let edgeIds = edgeMap[node.id] || []
  let fromNodes = []
  
  for (let i = 0; i < chartData.nodes.length; i++) {
    if(edgeIds.indexOf(chartData.nodes[i].id) > -1){
     
      fromNodes.push(chartData.nodes[i])
    }
  }
  
  entityDetailsRef.value.open(node, fromNodes)
}

const handleCancel = () => {
  chartBoxRef.value.destroy()
  show.value = false;
};

const handleOk = () => {
  show.value = false;
};

const open = (data) => {
  state.file_id = data.file_id
  state.data_id = data.id
  show.value = true

  setTimeout(() => {
    getData()
  }, 350);
};

defineExpose({
  open,
})

</script>
<style lang="less" scoped>
.my-modal{
  position: fixed;
  top: 0;
  left: 0;
}
.graph-model-body{
  position: relative;
  height: 90vh;
  width: 100%;
  overflow: hidden;

  .loading-wrapper{
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(242, 244, 247, 0.1);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 10;
  }

  .right-box{
    position: absolute;
    top: 0;
    right: 0;
    bottom: 0;
  }
}
</style>

<style lang="less">
.graph-model-wrapper {
  .ant-modal{
    .ant-modal-body,
    .ant-modal-header,
    .ant-modal-content{
      background: #F2F4F7 !important;
    }


  }
}
</style>
