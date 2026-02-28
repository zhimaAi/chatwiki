<template>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    :width="width"
    :height="height"
    :viewBox="`0 0 ${width} ${height}`"
  >
    <!-- 背景圆角矩形 -->
    <rect
      x="0"
      y="0"
      :width="width"
      :height="height"
      :rx="borderRadius"
      :ry="borderRadius"
      :fill="backgroundColor"
    />
    <!-- 文字 -->
    <text
      :x="width / 2"
      :y="height / 2 + fontSize / 3"
      text-anchor="middle"
      :font-size="fontSize"
      :fill="textColor"
    >
      {{ text }}
    </text>
  </svg>
</template>

<script setup>
defineProps({
  text: {
    type: String,
    default: ''
  },
  width: {
    type: Number,
    default: 32
  },
  height: {
    type: Number,
    default: 16
  },
  borderRadius: {
    type: Number,
    default: 4
  },
  backgroundColor: {
    type: String,
    default: '#f0f0f0'
  },
  textColor: {
    type: String,
    default: '#595959'
  },
  fontSize: {
    type: Number,
    default: 12
  }
})
</script>
