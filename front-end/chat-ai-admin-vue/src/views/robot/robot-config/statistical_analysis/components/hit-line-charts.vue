<template>
  <div class="container-echart">
    <div :id="'echart' + timeId" class="interface-echart"></div>
  </div>
</template>
<script setup>
import * as echarts from 'echarts'
import { ref, onMounted, watch, onBeforeUnmount } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.statistical-analysis.components.hit-line-charts')

const props = defineProps({
  xDataArray: {
    // 柱形图x轴数据
    type: Array,
    default: function () {
      return ['12-01', '11-29', '12-02', '12-03', '12-04', '12-05', '12-06']
    }
  },
  yDataArray: {
    //图标信息提示
    type: Array,
    default: function () {
      return [
        {
          name: 'label_daily_active_users',
          type: 'line',
          showSymbol: false,
          symbol: 'circle', // 标记的图形
          symbolSize: 10,
          itemStyle: {
            borderColor: '#fff', // 边框颜色
            borderWidth: 2, // 边框宽度
            borderType: 'solid', // 边框类型
            borderDashOffset: 0, // 控制边框的虚线样式
            borderDashArray: [10, 5], // 控制边框的虚线样式
            gap: 2, // 间隔透明度的大小
            borderRadius: 50 // 圆角半径，可以使边框变成圆形
          },
          color: '#2475FC',
          lineStyle: {
            width: 1,
            color: '#2475FC'
          },
          label: {
            show: true, // 在折线拐点上显示数据
            fontSize: 8,
            color: '#fff',
            fontWeight: 10
          },
          smooth: true, // 开启平滑过渡
          areaStyle: {
            //区域样式
            origin: 'start', //向最小值方向渐变，y轴有负值要写
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              {
                offset: 0,
                color: 'rgba(36,117,252, 0.3)'
              },
              {
                offset: 1,
                color: 'rgba(36,117,252, 0)'
              }
            ])
          },
          data: [542, 1222, 147, 5]
        }
      ]
    }
  },
  grid: {
    type: Object,
    default: function () {
      return {
        right: '5%', //图表距离容器右侧距离
        left: '0%',
        top: '0%',
        bottom: '35%'
      }
    }
  }
})
let myEchart
const timeId = ref(Math.floor(new Date().getTime() * Math.random())) // 使该图表保持唯id
const myEchartData = () => {
  // 示例配置，需要根据实际组件调整
  const option = {
    tooltip: {
      trigger: 'axis',
      formatter: function (params) {
        let tip = params[0].axisValue + '<br>'
        params.forEach((item) => {
          let value = item.value
          let seriesName = item.seriesName.startsWith('label_') || item.seriesName.startsWith('legend_') ? t(item.seriesName) : item.seriesName
          if (seriesName === t('legend_hit_rate')) {
            value += '%'
          }
          tip += `
              <span style="display:inline-block;margin-right:5px;border-radius:10px;
              width:9px;height:9px;background-color:${item.color};"></span>
              ${seriesName}: ${value}<br>
            `
        })
        return tip
      }
    },
    grid: {
      //图表距离边框的偏离
      right: '7%', //图表距离容器右侧距离
      left: '4%',
      top: '10%',
      bottom: '10%'
    },
    legend: {
      data: [t('legend_total_messages'), t('legend_hit_messages'), t('legend_miss_messages'), t('legend_hit_rate')],
      bottom: 0
    },
    xAxis: {
      type: 'category',
      data: props.xDataArray
    },
    yAxis: [
      {
        type: 'value',
        name: t('yaxis_label_count')
      },
      {
        type: 'value',
        name: t('yaxis_label_percentage'),
        axisLabel: {
          formatter: '{value}%'
        }
      }
    ],
    series: props.yDataArray.map((series, index) => {
      const translatedName = series.name.startsWith('label_') || series.name.startsWith('legend_') ? t(series.name) : series.name
      if (translatedName === t('legend_hit_rate')) {
        return {
          ...series,
          name: translatedName,
          yAxisIndex: 1
        }
      }
      return {
        ...series,
        name: translatedName
      }
    })
  }
  myEchart.setOption(option, { notMerge: true })
}

watch(
  //监控数据变化
  () => props.xDataArray,
  () => {
    setTimeout(() => {
      myEchartData()
    }, 500)
  },
  { deep: true }
)
watch(
  //监控数据变化
  () => props.yDataArray,
  () => {
    setTimeout(() => {
      myEchartData()
    }, 500)
  },
  { deep: true }
)

onMounted(() => {
  setTimeout(() => {
    const dom = document.getElementById(`echart${timeId.value}`)
    myEchart = echarts.init(dom)
    myEchartData()
  }, 500)

  // 当窗口发生变化时
  window.addEventListener('resize', () => {
    myEchart.resize()
  })
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', () => {
    myEchart.resize()
  })
})
</script>
<style lang="less" scoped>
.container-echart {
  width: 100%;
  display: flex;
  // justify-content: center;
  // align-items: center;
  overflow: hidden;
}
.interface-echart {
  width: 100%;
  height: 470px;
}
</style>
