<style lang="less" scoped>
.preview-tips {
  position: relative;
}
.preview-tips-content {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px 12px;
  line-height: 22px;
  font-size: 14px;
  color: rgba(0, 0, 0, 0.88);
  background-color: #fffbe6;
  border: 1px solid #ffe58f;
  word-wrap: break-word;
  box-sizing: border-box;
}
.preview-tips .tips-icon {
  margin-right: 8px;
}
.preview-tips .tips-icon svg {
  vertical-align: -2px;
}
</style>

<template>
  <div class="preview-tips">
    <div class="preview-tips-content">
      <span class="tips-icon">
        <svg
          t="1748426707311"
          class="icon"
          viewBox="0 0 1024 1024"
          version="1.1"
          xmlns="http://www.w3.org/2000/svg"
          p-id="5019"
          width="14"
          height="14"
        >
          <path
            d="M512 64C264.8 64 64 264.8 64 512s200.8 448 448 448 448-200.8 448-448S759.2 64 512 64z m32 704h-64v-64h64v64z m-64-128V256h64v384h-64z"
            p-id="5020"
            fill="#faad14"
          ></path>
        </svg>
      </span>
      <span class="tips-text"
        >当前仅支持预览，链接有效期10分钟，如果想要公开，请配置自定义域名后使用分享功能获取公开访问链接</span
      >
    </div>
  </div>
</template>

<script setup></script>
