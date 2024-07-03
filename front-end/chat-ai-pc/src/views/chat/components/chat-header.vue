<style lang="less" scoped>
.chat-header{
    position: relative;
    padding: 16px;
    background-image: linear-gradient(90deg, #2334E6 0%, #00A0FB 100%);
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.12);
    .close-btn{
        position: absolute;
        right: 4px;
        top: 4px;
        color: #fff;
        font-size: 24px;
        cursor: pointer;
    }
    .chat-header-body{
        display: flex;
        align-items: center;

        .avatar{
            width: 40px;
            height: 40px;
            border-radius: 50%;
        }

        .robot-info{
            margin-left: 8px;
        }
        .robot-name{
            height: 22px;
            line-height: 22px;
            font-size: 14px;
            font-weight: 600;
            opacity: 0.95;
            color: #ffffff;
            overflow: hidden;
        }
        .robot-intro{
            max-height: 40px;
            line-height: 20px;
            font-size: 12px;
            font-weight: 400;
            opacity: 0.85;
            color: #ffffff;
            overflow: hidden;
        }
    }
}
</style>

<template>
  <div class="chat-header" :style="{background: backgroundColor}">
    <svg-icon class="close-btn" name="close" @click="handleClose" />
    <div class="chat-header-body">
        <img class="avatar" :src="externalConfigPC.headImage" alt="" v-if="externalConfigPC.headImage" />
        <div class="robot-info">
            <div class="robot-name">{{ externalConfigPC.headTitle }}</div>
            <div class="robot-intro">{{ externalConfigPC.headSubTitle }}</div>
        </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { postCloseChat } from '@/event/postMessage'
import { storeToRefs } from 'pinia'
import { useChatStore } from '@/stores/modules/chat'

const chatStore = useChatStore()
const { externalConfigPC } = storeToRefs(chatStore)
const { headBackgroundColor } = externalConfigPC.value.pageStyle

const backgroundColor = computed(() => {
    const [type, direction, color1, color2] = headBackgroundColor.split(',')

    if(type === 'color'){
        return color1
    }

    return `${type}(${direction}, ${color1}, ${color2})`
})

const handleClose = () => {
  postCloseChat()
}
</script>

