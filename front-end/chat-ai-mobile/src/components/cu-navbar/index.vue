<style lang="less" scoped>
.cu-navbar {
    position: relative;
    display: flex;
    align-items: center;
    height: 44px;
    color: #fff;
    background: linear-gradient(270deg, #2475FC 0%, #2497FC 100%);

    .navbar-left,
    .navbar-right {
        position: absolute;
        top: 0;
        bottom: 0;
        display: flex;
        align-items: center;
        padding: 0 16px;
    }

    .navbar-left {
        left: 0;
    }

    .navbar-right {
        right: 0;
    }

    .navbar-arrow {
        font-size: 18px;
    }

    .navbar-body {
        max-width: 60%;
        margin: 0 auto;
    }

    .navbar-title {
        font-size: 16px;
        line-height: 22px;
        font-weight: 400;
    }
}
</style>

<template>
    <div class="cu-navbar" v-if="isMobile" :style="{background: props.backgroundColor}">
        <!-- <div class="navbar-left"><van-icon class="navbar-arrow" name="arrow-left" /></div> -->
        <div class="navbar-body">
            <div class="navbar-title van-ellipsis">{{ title }}</div>
        </div>
        <div class="navbar-right"></div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const isMobile = ref(false);

const checkMobile = () => {
    const userAgent = navigator.userAgent
    // 正则表达式来检测常见的手机设备标识
    isMobile.value = /android|webos|iphone|ipad|ipod|blackberry|iemobile|opera mini/i.test(userAgent);
};

const props = defineProps({
    backgroundColor: {
        type: String,
        default: '#2475FC'
    },
    title: {
        type: String,
        default: ''
    }
})

onMounted(() => {
    checkMobile();
})
</script>
