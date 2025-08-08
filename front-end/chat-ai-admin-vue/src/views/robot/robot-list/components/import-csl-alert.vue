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
  <a-modal width="600px" v-model:open="show" :confirmLoading="saveLoading" title="导入csl文件创建机器人" @ok="handleSave"
    @cancel="onCancel">
    <div class="form-box">
      <div class="tip-box">
        <a-alert type="info">
          <template #description>
            您可以将您编排好的机器人导出为csl文件,共享给他人使用。也可以从他人处获取导出的csl文件,通过导入形式创建机器人。
          </template>
        </a-alert>
      </div>
      <a-upload-dragger v-model:fileList="fileList" name="file" :multiple="false" :maxCount="1" accept=".csl"
        :before-upload="beforeUpload" @change="handleChange">
        <p class="ant-upload-drag-icon">
          <inbox-outlined class="upload-icon"></inbox-outlined>
        </p>
        <p class="upload-text">点击或将文件拖拽到这里上传</p>
        <p class="upload-hint">仅支持从chatwiki导出生成的csl文件</p>
      </a-upload-dragger>
    </div>
  </a-modal>
</template>

<script setup>
import { ref } from 'vue';
import { InboxOutlined } from '@ant-design/icons-vue';
import { message, Upload } from 'ant-design-vue';
import { robotImport } from '@/api/robot/index.js';

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
    message.error('请上传csl文件');
    return;
  }

  // 使用robotImport方法处理文件上传
  saveLoading.value = true;

  robotImport({
    file: fileList.value[0].originFileObj
  }).then((res) => {
    saveLoading.value = false;
    show.value = false;
    message.success('导入成功');
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
    message.error('只能上传csl格式的文件');
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