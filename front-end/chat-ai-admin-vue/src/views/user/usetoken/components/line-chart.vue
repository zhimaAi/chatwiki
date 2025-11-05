<template>
  <div class="container-echart">
    <div ref="chartRef" class="interface-echart"></div>
  </div>
</template>
<script setup>
import * as echarts from 'echarts'
import { ref, onMounted, watch, onBeforeUnmount, nextTick } from 'vue'
// import dayjs from 'dayjs'
const props = defineProps({
  options: {
    // 柱形图x轴数据
    type: Object,
    default: () => {
      return {
        xAxis: [],
        series: []
      }
    }
  }
})

let myEchart
const updateChat = () => {
  const option = {
    title: {
      textStyle: {
        color: '#262626',
        textAlign: 'center',
        fontFamily: 'PingFang SC',
        fontSize: '14px',
        fontStyle: 'normal',
        fontWeight: '600'
      }
    },
    legend: {
      //字体对应折线标识
      show: true,
      bottom: '0%',
      textStyle: {
        color: '#262626',
        fontSize: '14px'
      },
      data: ['输入', '输出', '合计'],
      itemHeight: 6, //圆心
      itemWidth: 6, //折线

      //图标信息提示
      type: 'scroll',
      orient: 'horizontal',
      align: 'auto',
      icon: 'circle'
    },
    grid: {
      //图表距离边框的偏离
      right: '1.5%', //图表距离容器右侧距离
      left: '4%',
      top: '5%',
      bottom: '20%'
    },
    tooltip: {
      trigger: 'axis'
    }, // 设置图案和容器的距离
    xAxis: {
      type: 'category',
      axisTick: {
        show: true,
        lineStyle: {
          color: '#D9D9D9' // 修改为你想要的颜色
        }
      },
      boundaryGap: false,
      //   name: '时间',
      nameLocation: 'end',
      data: props.options.xAxis,
      axisLabel: {
        show: true, // 显示坐标轴标签
        fontSize: 10, // 设置标签字体大小
        interval: 'auto', // 设置标签显示的间隔，为0时强制显示所有标签
        rotate: 0, // 设置倾斜角度
        lineHeight: 26
      },
      // 坐标轴轴线
      axisLine: {
        show: false
      }
    },
    yAxis: {
      type: 'value',
      nameLocation: 'end',
      min: 0
    },
    series: [
      {
        name: '输入',
        type: 'line',
        showSymbol: false,
        data: props.options.series.map((item) => item.prompt_token),
        smooth: true, // 开启平滑过渡
        color: '#2475fc'
      },
      {
        name: '输出',
        type: 'line',
        showSymbol: false,
        data: props.options.series.map((item) => item.completion_token),
        smooth: true, // 开启平滑过渡
        color: '#646161'
      },
      {
        name: '合计',
        type: 'line',
        showSymbol: false,
        data: props.options.series.map((item) => item.total_token),
        smooth: true, // 开启平滑过渡
        color: '#4B7900'
      }
    ]
  }
  // { notMerge: true } 解决删除数据时,数据不刷新的问题
  myEchart.setOption(option, { notMerge: true })
}

watch(
  //监控数据变化
  () => props.options,
  () => {
    setTimeout(() => {
      updateChat()
    }, 200)
  },
  { deep: true }
)

const chatResize = () => {
  myEchart.resize()
}

const chartRef = ref(null)

onMounted(() => {
  nextTick(() => {
    myEchart = echarts.init(chartRef.value)

    updateChat()
  })

  // 当窗口发生变化时
  window.addEventListener('resize', chatResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', chatResize)
})
</script>

<style lang="less" scoped>
.container-echart {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  margin-top: 16px;
}
.interface-echart {
  width: 100%;
  height: 450px;
}
</style>
