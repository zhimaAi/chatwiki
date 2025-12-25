<template>
  <component
    :is="currentComp"
    :action="action"
    :actionName="actionName"
    :variableOptions="variableOptions"
    :node="node"
    @updateVar="emit('updateVar')"
  />
</template>
<script setup>
import {ref, reactive, computed} from 'vue';
import OfficialDraftStore from "./official-draft-store.vue";
import OfficialDraftOperate from "./official-draft-operate.vue";

const compMap = {
  add_draft: OfficialDraftStore,
  update_draft: OfficialDraftStore,
  delete_draft: OfficialDraftOperate,
  publish_draft: OfficialDraftOperate,
  preview_message: OfficialDraftOperate,
}
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

const currentComp = computed(() => {
  return compMap[props.actionName]
})
</script>
