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
  <a-modal class="add-model-alert" v-model:open="show" :title="t('title')" @ok="handleOk">
    <div class="form-wrapper">
      <img class="model-logo" src="@/assets/img/user/model/anthropicc.png" alt="" />

      <div class="form-box">
        <a-form layout="vertical">
          <a-form-item>
            <template #label>
              <div>
                <span style="padding-right: 16px">{{ t('api_key_label') }}</span>
                <a href="#">{{ t('how_to_get') }}</a>
              </div>
            </template>

            <a-input v-model:value="formState.api_key" :placeholder="t('api_key_placeholder')" />
          </a-form-item>

          <a-form-item>
            <template #label>
              <div>
                <span>{{ t('secret_key_label') }}</span>
              </div>
            </template>

            <a-input v-model:value="formState.api_key" :placeholder="t('secret_key_placeholder')" />
          </a-form-item>
        </a-form>
      </div>
    </div>
  </a-modal>
</template>
<script setup>
import { ref, reactive } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.model.components.add-model-alert')

const show = ref(false)
const formState = reactive({
  api_key: '',
  secret_key: ''
})
const handleOk = () => {}
</script>
