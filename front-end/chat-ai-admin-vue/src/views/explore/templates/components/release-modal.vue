<template>
  <a-modal
    v-model:open="visible"
    title="请选择需要上架的应用"
    width="660px"
  >
    <a-alert type="info" class="zm-alert-info mt16" message="上架工作流时请提前清除工作流中引用的apikey等授权、密钥信息，避免个人资产损失已上架应用再次上架时，会覆盖上次上架的内容"/>
    <a-form class="mt24" layout="vertical">
      <a-form-item label="上架应用" required>
        <a-dropdown>
          <div class="app-box">
            <img class="avatar" src="https://xkf-upload.oss-cn-hangzhou.aliyuncs.com/dev/p/chat_wiki_plugin/1/0/202511/45/6b83aca346be5415bdfd349bb49111.jpg"/>
            <div class="info">
              <div class="left">
                <div class="name">同程旅行规划师</div>
                <div class="desc zm-line2">
                  “深度解析八字命理，精准预测人生运势。为您揭开八字中隐藏的能量与未来机遇，更能绘制八字意象图文，具象呈现命运密.“深度解析八字命理，精准预测人生运势。为您揭开八字中隐藏的能量与未来机遇，更能绘制八字意象图文，具象呈现命运密.
                </div>
              </div>
              <div class="right">
                <DownOutlined class="zm-pointer"/>
              </div>
            </div>
          </div>
          <template #overlay>
            <a-menu>
              <a-menu-item v-for="i in 5" :key="i">
                <div class="app-item">
                  <img class="avatar" src="https://xkf-upload.oss-cn-hangzhou.aliyuncs.com/dev/p/chat_wiki_plugin/1/0/202511/45/6b83aca346be5415bdfd349bb49111.jpg"/>
                  <span>小红书账号拆解神器</span>
                </div>
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
        <div class="last-time">上次上架时间：2025-10-11 15:00，上传后会覆盖原有模板</div>
      </a-form-item>
      <a-form-item label="上架分类">
        <a-select v-model:value="formState.appid" placeholder="请选择上架分类">
          <a-select-option value="jack">Jack</a-select-option>
          <a-select-option value="lucy">Lucy</a-select-option>
          <a-select-option value="Yiminghe">yiminghe</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="模板作者" required>
        <a-input placeholder="请输入模板作者名称"/>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import {ref, reactive} from 'vue';
import {DownOutlined} from '@ant-design/icons-vue';

const visible = ref(false)
const formState = reactive({
  appid: undefined
})

</script>

<style scoped lang="less">
.app-box {
  display: flex;
  align-items: center;
  border-radius: 12px;
  border: 1px solid #E4E6EB;
  padding: 24px;

  .avatar {
    flex-shrink: 0;
    width: 62px;
    height: 62px;
    border-radius: 14.59px;
    margin-right: 12px;
  }

  .info {
    display: flex;
    align-items: center;
    justify-content: space-between;
    color: #595959;
    font-size: 14px;
    font-weight: 400;
    line-height: 22px;
    gap: 24px;

    .name {
      color: #262626;
      font-size: 16px;
      font-weight: 600;
    }
  }
}

.app-item {
  display: flex;
  align-items: center;
  gap: 8px;

  .avatar {
    width: 20px;
    height: 20px;
    border-radius: 4px;
  }
}

.last-time {
  padding: 9px 16px;
  border-radius: 6px;
  background: #F2F4F7;
  color: #000000a6;
  font-size: 14px;
  font-weight: 400;
  margin-top: 8px;
}

.mt24 {
  margin-top: 24px;
}

.mt16 {
  margin-top: 16px;
}
</style>
