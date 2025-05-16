<style lang="less" scoped>
.zoom-select {
  position: relative;
  margin: 0 4px;

  .zoom-select-label {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 80px;
    height: 32px;
    border-radius: 6px;
    font-size: 14px;
    color: #595959;
    transition: all 0.2s;
  }

  .zoom-select-options {
    position: absolute;
    bottom: 48px;
    left: 0;
    padding: 2px;
    font-size: 14px;
    border-radius: 6px;
    color: #595959;
    background-color: #fff;
    box-shadow:
      0 6px 30px 5px #0000000d,
      0 16px 24px 2px #0000000a,
      0 8px 10px -5px #00000014;
  }
  .zoom-select-option {
    width: 77px;
    height: 32px;
    line-height: 32px;
    margin-bottom: 2px;
    border-radius: 6px;
    color: #595959;
    text-align: center;
    cursor: pointer;
    transition: all 0.2s;

    &:last-child {
      margin-bottom: 0;
    }
  }

  .option-line {
    border-bottom: 1px solid #f0f0f0;
    margin-bottom: 2px;
  }

  .zoom-select-label:hover,
  .zoom-select-option:hover {
    background-color: #e4e6eb;
    cursor: pointer;
  }
}
</style>

<template>
  <div class="zoom-select">
    <div class="zoom-select-label" @click="showMenu">
      {{ title }}
    </div>
    <div class="zoom-select-options" v-show="isShowMenu">
      <div
        class="zoom-select-option"
        v-for="item in menus"
        :key="item.value"
        @click="handleMenu(item)"
      >
        {{ item.label }}
      </div>
      <div class="option-line"></div>
      <div class="zoom-select-option" @click="handleMenu({ label: '5%', value: 5 })">全览</div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'NodeCommon',
  props: {
    title: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      isShowMenu: false,
      menus: [
        {
          label: '10%',
          value: 0.1,
        },
        {
          label: '50%',
          value: 0.5,
        },
        {
          label: '100%',
          value: 1,
        },
        {
          label: '200%',
          value: 2,
        },
        {
          label: '300%',
          value: 3,
        },
        {
          label: '400%',
          value: 4,
        },
        {
          label: '500%',
          value: 5,
        },
      ],
    }
  },
  mounted() {
    document.addEventListener('click', this.documentClick)
  },
  beforeUnmount() {
    document.removeEventListener('click', this.documentClick)
  },
  methods: {
    handleMenu(item) {
      this.$emit('change', item)
    },
    showMenu() {
      this.isShowMenu = !this.isShowMenu
    },
    // 点击菜单以外的地方，隐藏菜单
    documentClick(e) {
      if (this.isShowMenu) {
        const menus = this.$el.querySelector('.zoom-select-options')
        if (
          !menus.contains(e.target) &&
          e.target !== this.$el.querySelector('.zoom-select-label')
        ) {
          this.isShowMenu = false
        }
      }
    },
  },
}
</script>
