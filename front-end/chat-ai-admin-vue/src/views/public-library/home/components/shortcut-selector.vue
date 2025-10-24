<style lang="less" scoped>
:deep(.ant-tree-treenode-selected) {
  .ant-tree-title {
    background-color: #e6f7ff;
  }
}

:deep(.ant-tree-checkbox) {
  .ant-tree-checkbox-inner {
    border-radius: 50%;
  }
}

.tree-node-title {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>

<template>
  <a-modal v-model:open="visible" :title="title" @ok="handleOk">
    <a-tree
      blockNode
      checkable
      v-model:expandedKeys="expandedKeys"
      :defaultExpandAll="false"
      :checkStrictly="true"
      :checkedKeys="checkedKeys"
      :tree-data="treeData"
      :load-data="loadData"
      v-if="treeData.length > 0"
      @check="onCheck"
    >
      <template #title="{ data, level }">
        <span class="tree-node-title">
          <span class="doc-icon" v-if="data.doc_icon">{{ data.doc_icon }}</span>
          <span class="doc-icon" v-else>{{ data.is_dir == 1 ? props.iconTemplateConfig.levels[level].folder_icon : props.iconTemplateConfig.levels[level].doc_icon }}</span>
          <span>{{ data.title }}</span>
        </span>
      </template>
    </a-tree>
  </a-modal>
</template>

<script setup>
import { ref, computed } from 'vue';
import { getDocList } from '@/api/public-library'

const emit = defineEmits(['ok'])
const props = defineProps({
  iconTemplateConfig: {
    type: Object,
    default: () => ({})
  },
})


const visible = ref(false);
const libraryKey = ref('')
const expandedKeys = ref([]);
const checkedKeys = ref([]);
const treeData = ref([]);
const oldDocKey = ref(null);
const oldDocId = ref(null);

const title = computed(() => {
  return oldDocKey.value ? '修改文档快捷方式' : '添加文档快捷方式'
})

const getMyDocList = async (doc) => {
  let data = {
    library_key: libraryKey.value,
    pid: doc ? doc.id : 0
  }

  const res = await getDocList(data);

  let list = res.data || [];

  list.forEach((item) => {
    item.id = item.id * 1;
    item.pid = item.pid * 1;
    item.key = item.id;
    item.level = doc ? doc.level + 1 : 0;
    item.sort = item.sort * 1;
    item.children_num = item.children_num * 1;
    item.children = [];
    item.hasLoaded = false;
    item.isLeaf = item.children_num == 0;
    item.disabled = item.is_draft == 1;
  });

  if (!doc) {
    treeData.value = list;
  } else {
    doc.hasLoaded = true;
    doc.children = list;

    treeData.value = [...treeData.value];
  }

  return res;
}


const loadData = (treeNode) => {
  return new Promise((resolve) => {
    if (treeNode.dataRef.hasLoaded) {
      resolve()
      return
    }

    getMyDocList(treeNode.dataRef).then(() => {
      resolve()
    })
  })
}

const open = (data = {}) => {
  libraryKey.value = data.library_key;
  expandedKeys.value = [];
  checkedKeys.value = [];
  treeData.value = [];
  oldDocKey.value = data.key || null;
  oldDocId.value = data.doc_id || null;

  if(oldDocKey.value){
    checkedKeys.value = [data.doc_id];
  }

  getMyDocList();

  visible.value = true;
}

const onCheck = (keys, e) => {
  if (e.checked) {
    checkedKeys.value = [e.node.key];
  } else {
    checkedKeys.value = [];
  }
}

const handleOk = () => {
  if (checkedKeys.value.length > 0) {
    emit('ok', {doc_id: checkedKeys.value[0], old_doc_key: oldDocKey.value, old_doc_id: oldDocId.value});
  }
  visible.value = false;
};

defineExpose({
  open
})
</script>
