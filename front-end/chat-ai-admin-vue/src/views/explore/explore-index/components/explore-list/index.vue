<style lang="less" scoped>
.list-box {
  display: flex;
  flex-flow: row wrap;
  margin: 0 -8px 0 -8px;
  flex: 1;
}
.list-item-wrapper {
  padding: 8px;
  width: 25%;
}
.list-item {
  position: relative;
  width: 100%;
  padding: 24px;
  border: 1px solid #E4E6EB;
  border-radius: 12px;
  background-color: #fff;
  transition: all 0.25s;
  cursor: pointer;

  &:hover {
    box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.12);
  }

  .explore-info {
    position: relative;
    display: flex;
    align-items: center;

    .explore-info-content{
      flex: 1;
      padding-left: 12px;
      overflow: hidden;
    }
  }

  .explore-title {
    height: 24px;
    line-height: 24px;
    margin-bottom: 4px;
    font-size: 16px;
    font-weight: 600;
    color: #262626;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .explore-type {
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 1;
    overflow: hidden;
    color: #8c8c8c;
    text-overflow: ellipsis;
    font-size: 12px;
    font-style: normal;
    font-weight: 400;
    line-height: 20px;
  }
  .item-body{
    margin-top: 12px;
  }
  .explore-desc {
    height: 44px;
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: rgb(89, 89, 89);
    // 超出2行显示省略号
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
  }

  .support-box {
    margin-top: 4px;
    display: flex;
    align-items: center;
    gap: 0px 8px;
    color: #8c8c8c;
    font-size: 12px;
    font-style: normal;
    font-weight: 400;
    line-height: 20px;
    flex-wrap: wrap;

    .support-item {
      white-space: nowrap;
    }
  }
  .item-footer {
    position: relative;
    z-index: 1;
    width: fit-content;
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 12px;
    color: #7a8699;
  }
  .explore-size {
    display: flex;
    line-height: 20px;
    font-size: 12px;
    font-weight: 400;
    color: #7a8699;

    .text-item {
      margin-right: 12px;
      &:last-child{
        margin-right: 0;
      }
    }
  }

  .action-box {
    font-size: 14px;
    height: 24px;
    color: #2475fc;
    display: flex;
    align-items: center;

    .action-item {
      display: flex;
      align-items: center;
      height: 100%;
      padding: 4px;
      border-radius: 6px;
      cursor: pointer;
      color: #595959;
      transition: all 0.2s;
    }
    .action-item:hover {
      background: #E4E6EB;
    }

    .action-icon {
      font-size: 16px;
    }
  }
}

// 大于1920px
@media screen and (min-width: 1920px) {
  .list-box {
    .list-item-wrapper {
      width: 20%;
    }
  }
}

// 小屏保持三列
@media screen and (max-width: 1500px) {
  .list-box {
    .list-item-wrapper {
      width: 33.3333%;
    }
  }
}
</style>

<template>
  <div class="list-box">
    <div class="list-item-wrapper" v-for="item in props.list" :key="item.id">
      <div class="list-item" @click="handleClick($event, item)">
        <div class="explore-info">
          <svg-icon class="explore-icon" :name="item.ability_type" style="font-size: 62px; color: white;" />
          <div class="explore-info-content">
            <div class="explore-title">{{ item.explore_name || item.name }}</div>
            <div class="explore-type">ChatWiki</div>
          </div>
        </div>
        <div class="item-body">
          <div class="explore-desc">{{ item.explore_intro || item.introduction }}</div>

          <span class="support-box">
            <span v-for="ch in item.support_channels_list" :key="ch" class="support-item">
              <svg-icon name="support" class="icon" />
              {{ ch }}
            </span>
          </span>
        </div>

        <div class="item-footer">
            <a-switch
              :checked="(item.user_config?.switch_status == '1') || (item.robot_config?.switch_status == '1')"
              checked-children="开"
              un-checked-children="关"
              class="no-bubble"
              @change="(checked)=>handleSwitchChange(item, checked)"
            />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
const emit = defineEmits(['switchChange', 'clickItem'])

const router = useRouter()

const props = defineProps({
  list: {
    type: Array,
    default: () => []
  },
})

const handleSwitchChange = (item, checked) => emit('switchChange', item, checked)

const handleClick = async (e, item) => {
  try {
    const target = e?.target
    if (target && typeof target.closest === 'function') {
      const blocker = target.closest('.no-bubble')
      if (blocker) return
    }
    // 点击关注后回复跳转
    if (item.ability_type == 'robot_subscribe_reply') {
      router.push({ path: '/explore/index/subscribe-reply' })
      return
    }
  } catch (_) {}
  emit('clickItem', item)
}
</script>
