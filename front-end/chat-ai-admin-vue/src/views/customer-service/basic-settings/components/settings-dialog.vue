<template>
  <a-modal class="workspace-settings-dialog" centered :mask-style="{ backgroundColor: 'rgba(0, 0, 0, 0.65)' }">
    <template v-if="$slots.title" #title><slot name="title" /></template>
    <template #closeIcon><slot name="closeIcon"><img src="@/assets/img/customer-service/basic-settings/dialogs/close.svg" alt="" /></slot></template>
    <slot />
    <template v-if="$slots.footer" #footer><slot name="footer" /></template>
  </a-modal>
</template>

<style lang="less">
.workspace-settings-dialog.ant-modal {
  max-width: calc(100vw - 32px);
  .ant-modal-content { padding: 24px; border: 1px solid #e8e8e8; border-radius: 16px; box-shadow: 0 8px 24px 8px #0808080a, 0 6px 12px #0808080a, 0 4px 8px -8px #0808080f; }
  .ant-modal-header { margin-bottom: 24px; padding: 0; border: 0; }
  .ant-modal-title { padding-right: 32px; color: #080808; font-size: 20px; font-weight: 500; line-height: 28px; }
  .ant-modal-close { top: 26px; right: 24px; width: 24px; height: 24px; }
  .ant-modal-close-x { display: flex; align-items: center; justify-content: center; line-height: 24px; }
  .ant-modal-body { max-height: calc(100dvh - 188px); padding: 0; overflow-y: auto; }
  .ant-modal-footer { display: flex; justify-content: flex-end; gap: 8px; margin-top: 24px; padding: 0; border: 0; }
  .ant-modal-footer .ant-btn { height: 32px; min-width: 56px; margin: 0; padding: 0 12px; border-radius: 6px; border-color: #d8d8d8; box-shadow: none; }
  .ant-modal-footer .ant-btn-primary { background: #146ef5; border-color: #146ef5; }
}
</style>
