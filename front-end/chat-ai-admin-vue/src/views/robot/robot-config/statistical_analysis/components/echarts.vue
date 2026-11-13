<template>
    <div class="container-echart">
        <div :id="'echart' + timeId" class="interface-echart"></div>
    </div>
</template>
<script setup>
  import * as echarts from 'echarts'
  import { ref, onMounted, watch, onBeforeUnmount } from 'vue'

  defineOptions({ name: 'StatisticsEcharts' })
  
  const props = defineProps({
    xDataArray: {
      // 柱形图x轴数据
      type: Array,
      default: function () {
        return ['12-01','11-29','12-02','12-03','12-04','12-05','12-06']
      }
    },
    yDataArray: {
      //图标信息提示
      type: Array,
      default: function () {
        return [
          {
            name: '日活用户数（人）',
            type: 'line',
            showSymbol: false,
            symbol: 'circle', // 标记的图形
            symbolSize: 10,
            itemStyle: {
              borderColor: '#fff', // 边框颜色
              borderWidth: 2,      // 边框宽度
              borderType: 'solid',  // 边框类型
              borderDashOffset: 0, // 控制边框的虚线样式
              borderDashArray: [10, 5], // 控制边框的虚线样式
              gap: 2,              // 间隔透明度的大小
              borderRadius: 50     // 圆角半径，可以使边框变成圆形
            },
            color: '#2475FC',
            lineStyle: {
              width: 1,
              color: '#2475FC'
            },
            label: {
              show: true, // 在折线拐点上显示数据
              fontSize: 8,
              color: "#fff",
              fontWeight: 10,
            },
            smooth: true, // 开启平滑过渡
            areaStyle: {//区域样式
              origin: "start",//向最小值方向渐变，y轴有负值要写
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                {
                  offset: 0,
                  color: "rgba(36,117,252, 0.3)",
                },
                {
                  offset: 1,
                  color: "rgba(36,117,252, 0)",
                },
              ]),
            },
            data: [542, 1222, 147, 5]
          }
        ]
      }
    }
  })
  let myEchart
  let initTimer
  let updateTimer
  const resizeChart = () => myEchart?.resize()
  const timeId = ref(Math.floor(new Date().getTime() * Math.random())); // 使该图表保持唯id
  const myEchartData = () => {
    if (!myEchart) return
    const option = {
      title: {
        textStyle: {
          color: '#262626',
          textAlign: 'center',
          fontFamily: "PingFang SC",
          fontSize: '14px',
          fontStyle: 'normal',
          fontWeight: '600',
        },
        left: '7%'
      },
      legend: {
        //字体对应折线标识
        show: false,
        right: "1%",
        textStyle: {
          color: "#262626",
          fontSize: "14px",
        },
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
        right: '5%', //图表距离容器右侧距离
        left: '10%',
        top: '10%',
        bottom: '25%'
      },
      tooltip: { 
        trigger: 'axis',
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
        data: props.xDataArray,
        axisLabel: {
            show: true, // 显示坐标轴标签
            fontSize: 10, // 设置标签字体大小
            interval: 0, // 设置标签显示的间隔，为0时强制显示所有标签
            rotate: 60 // 设置倾斜角度
        },
          // 坐标轴轴线
          axisLine: {
            show: false
          },
        //   // 坐标轴刻度
        //   axisTick: {
        //     show: false
        //   },
        //   // 刻度标签
        //   axisLabel: {
        //     show: false
        //   }
      },
      yAxis: {
        type: 'value',
        //   name: '值',
        nameLocation: 'end',
        //   // 坐标轴轴线
        //   axisLine: {
        //     show: false
        //   },
        //   // 坐标轴刻度
        //   axisTick: {
        //     show: false
        //   },
        //   // 刻度标签
        //   axisLabel: {
        //     show: false
        //   },
        min: 0,
        //   splitLine: {
        //     show: false // 不显示网格线
        //   }
      },
      series: props.yDataArray,
    }
    // { notMerge: true } 解决删除数据时,数据不刷新的问题
    myEchart.setOption(option, { notMerge: true })
  }
   
  const scheduleUpdate = () => {
    clearTimeout(updateTimer)
    updateTimer = setTimeout(myEchartData, 500)
  }

  watch(() => props.xDataArray, scheduleUpdate, { deep: true })
  watch(() => props.yDataArray, scheduleUpdate, { deep: true })

  onMounted(() => {
    initTimer = setTimeout(() => {
      const dom = document.getElementById(`echart${timeId.value}`)
      if (!dom) return
      myEchart = echarts.init(dom)
      myEchartData()
    }, 500)
    window.addEventListener('resize', resizeChart)
  })

  onBeforeUnmount(() => {
    clearTimeout(initTimer)
    clearTimeout(updateTimer)
    window.removeEventListener('resize', resizeChart)
    myEchart?.dispose()
    myEchart = undefined
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
  }
  .interface-echart {
    width: 100%;
    height: 310px;
  }
  </style>
  