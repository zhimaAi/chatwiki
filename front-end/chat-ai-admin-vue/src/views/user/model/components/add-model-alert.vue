<style lang="less" scoped>
.add-model-alert {
  .form-wrapper {
    margin-top: 24px;
  }
  .model-logo {
    display: block;
    height: 26px;
    margin-bottom: 16px;
  }
}
</style>

<template>
  <a-modal class="add-model-alert" v-model:open="show" title="添加文心一言" @ok="handleOk">
    <div class="form-wrapper">
      <img class="model-logo" src="@/assets/img/user/model/anthropicc.png" alt="" />

      <div class="form-box">
        <a-form layout="vertical">
          <a-form-item>
            <template #label>
              <div>
                <span style="padding-right: 16px">API Key：</span>
                <a href="#">如何获取?</a>
              </div>
            </template>

            <a-input v-model:value="formState.api_key" placeholder="请输入您的API Key" />
          </a-form-item>

          <a-form-item>
            <template #label>
              <div>
                <span>Secret Key：</span>
              </div>
            </template>

            <a-input v-model:value="formState.api_key" placeholder="请输入您的Secret Key" />
          </a-form-item>
        </a-form>
      </div>
    </div>
  </a-modal>
</template>
<script setup>
import { ref, reactive } from 'vue'

const show = ref(false)
const formState = reactive({
  api_key: '',
  secret_key: ''
})
const handleOk = () => {}
</script>
