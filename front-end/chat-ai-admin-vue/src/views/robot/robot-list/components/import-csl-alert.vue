<style lang="less" scoped>
.form-box {
  margin-top: 24px;
}

.tip-box {
  margin-bottom: 16px;
}

::v-deep(.ant-upload-drag) {
  background: #FAFAFA;
  border: 1px dashed #D9D9D9;

  .ant-upload-btn {
    padding: 24px 0;
  }
}

.upload-icon {
  font-size: 48px;
  color: #2475FC;
}

.upload-text {
  margin-top: 12px;
  font-size: 16px;
  color: #262626;
}

.upload-hint {
  margin-top: 4px;
  font-size: 14px;
  color: #8C8C8C;
}
</style>

<template>
  <a-modal width="600px" v-model:open="show" :confirmLoading="saveLoading" :title="t('title_import_csl')" @ok="handleSave"
    @cancel="onCancel">
    <div class="form-box">
      <div class="tip-box">
        <a-alert type="info">
          <template #description>
            {{ t('tip_description') }}
          </template>
        </a-alert>
      </div>
      <a-upload-dragger v-model:fileList="fileList" name="file" :multiple="false" :maxCount="1" accept=".csl"
        :before-upload="beforeUpload" @change="handleChange">
        <p class="ant-upload-drag-icon">
          <inbox-outlined class="upload-icon"></inbox-outlined>
        </p>
        <p class="upload-text">{{ t('upload_text') }}</p>
        <p class="upload-hint">{{ t('upload_hint') }}</p>
      </a-upload-dragger>
    </div>
  </a-modal>
</template>

<script setup>
import { ref } from 'vue';
import { InboxOutlined } from '@ant-design/icons-vue';
import { message, Upload } from 'ant-design-vue';
import { robotImport } from '@/api/robot/index.js';
import { useI18n } from '@/hooks/web/useI18n';

const { t } = useI18n('views.robot.robot-list.components.import-csl-alert');

const emit = defineEmits(['ok']);

const show = ref(false);
const saveLoading = ref(false);
const fileList = ref([]);

const open = () => {
  show.value = true;
};

const onCancel = () => {
  show.value = false;
  fileList.value = [];
};

const handleSave = () => {
  if (fileList.value.length === 0) {
    message.error(t('msg_upload_csl'));
    return;
  }

  // 使用robotImport方法处理文件上传
  saveLoading.value = true;

  robotImport({
    file: fileList.value[0].originFileObj
  }).then((res) => {
    saveLoading.value = false;
    show.value = false;
    message.success(t('msg_import_success'));
    fileList.value = [];
    emit('ok', res.data);
  })
    .catch(() => {
      saveLoading.value = false;
      // message.error('导入失败');
    });
};

const beforeUpload = (file) => {
  const isCsl = file.name.endsWith('.csl');
  if (!isCsl) {
    message.error(t('msg_only_csl_format'));
    return Upload.LIST_IGNORE;
  }
  // 验证通过后返回false，阻止自动上传，由手动触发上传
  return false;
};

const handleChange = (info) => {
  const { status } = info.file;
  if (status !== 'uploading') {
    console.log(info.file, info.fileList);
  }
};

defineExpose({
  open
});
</script>