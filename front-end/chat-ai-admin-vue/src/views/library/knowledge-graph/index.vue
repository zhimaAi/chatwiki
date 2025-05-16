<style lang="less" scoped>
.knowledge-graph-page {
  position: relative;
  width: 100%;
  height: 100%;
  background: #fff;
  .loading-wrapper{
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(255, 255, 255, 0.8);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 10;
  }

  .search-box {
    position: absolute;
    top: 24px;
    left: 50%;
    margin-left: -400px;
    width: 800px;
  }

  .left-box {
    position: absolute;
    top: 98px;
    left: 16px;
    bottom: 24px;
  }

  .right-box {
    position: absolute;
    top: 98px;
    right: 16px;
    bottom: 24px;
  }
}
</style>

<template>
  <div class="knowledge-graph-page">
    <div class="loading-wrapper" v-if="loading"> <a-spin :spinning="loading" /></div>
    <ChartBox ref="chartBoxRef" :loading="loading" @nodeClick="onNodeClick" />
    <CollapseBtn @change="toggleFileList" />
    <div class="search-box">
      <SearchInput v-model:value="keyword" @search="handleSearch" />
    </div>
    <div class="left-box">
      <FileList @openFile="onOpenFile" @openFragment="onOpenFragment" :show="fileListShow" />
    </div>
    <div class="right-box">
      <EntityDetails ref="entityDetailsRef" @toDetails="toDetails" />
    </div>
  </div>
</template>

<script setup>
import { getFileGraphInfo } from '@/api/library/index'
import { onMounted, ref, reactive } from 'vue'
import ChartBox from './components/chart-box.vue'
import SearchInput from './components/search-input.vue'
import FileList from './components/file-list.vue'
import EntityDetails from './components/entity-details.vue'
import CollapseBtn from './components/collapse-btn.vue'

const chartBoxRef = ref(null)
const entityDetailsRef = ref(null)
const loading = ref(false)
const keyword = ref('')

const fileListShow = ref(true)
const state = reactive({
  file_id: '',
  data_id: '',
  search_term: ''
})

const chartData = reactive({
  edges: [],
  nodes: []
})

// 关系映射
let edgeMap = {}

const toggleFileList = (val) => {
  fileListShow.value = val
}

const onOpenFile = (item) => {
  state.file_id = item.id
  state.data_id = ''

  getGraphData()
}

const onOpenFragment = (item) => {
  state.data_id = item.id
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
      // item.caption = item.name
      item.id = item.id.toString()

      // 标记节点类型
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

      item.captions = [{
        value: item.name,
        styles: [],
      }]
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

const handleSearch = (val) => {
  state.search_term = val
  getGraphData()
}

const toDetails = (node) => {
  window.open(`/#/library/preview?id=${node.file_id}`)
}


onMounted(() => {
 
})
</script>
