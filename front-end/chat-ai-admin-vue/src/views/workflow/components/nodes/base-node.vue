<template>
  <div class="node-common-wrapper" :class="[`node-common-wrapper_${iconName}`]" :data-node-name="title" :data-node-type="node_type">
    <div class="node-common" :class="{ isHovered: isHovered, isSelected: isSelected }">
      <div class="node-header">
        <div class="node-header-left">
          <span class="node-icon" v-if="iconName">
            <img :src="getImggeUrl(iconName)" alt="" />
          </span>
          <span v-if="['1', 1].includes(node_type)" class="node-title">{{ title }}</span>
          <a-input
            v-else
            size="small"
            class="node-title-input"
            :value="title"
            @input="updateTitle($event.target.value)"
          ></a-input>
        </div>
        <div>
          <div class="node-menu-wrapper" v-if="menus.length > 0">
            <div class="node-menu-btn" @click.stop="showMenu">
              <img
                class="btn-icon"
                src="../../../../assets/img/workflow/node-menu-btn.svg"
                alt=""
              />
            </div>

            <div class="node-menus" v-show="isShowMenu">
              <div
                class="node-menu del-btn"
                v-for="item in menus"
                :key="item.key"
                @click.stop="handleMenu(item)"
              >
                {{ item.name }}
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="node-body">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script>
import SvgIcon from '@/components/svg-icon/index.vue'

export default {
  name: 'NodeCommon',
  inject: ['setTitle'],
  components: { SvgIcon },
  props: {
    iconName: {
      type: String,
      default: ''
    },
    node_type: {
      default: ''
    },
    nodeKey: {
      type: String,
      default: ''
    },
    value: {
      type: String,
      default: ''
    },
    title: {
      type: String,
      default: ''
    },
    menus: {
      type: Array,
      default() {
        return []
      }
    },
    isSelected: {
      type: Boolean,
      default: false
    },
    isHovered: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      isShowMenu: false
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
      this.$emit('handleMenu', item)
    },
    showMenu() {
      this.isShowMenu = !this.isShowMenu
    },
    // 点击菜单以外的地方，隐藏菜单
    documentClick(e) {
      if (this.isShowMenu) {
        const menus = this.$el.querySelector('.node-menus')
        if (
          !menus.contains(e.target) &&
          e.target !== this.$el.querySelector('.node-menu-btn') &&
          e.target !== this.$el.querySelector('.btn-icon')
        ) {
          this.isShowMenu = false
        }
      }
    },
    getImggeUrl(name) {
      let url = new URL(`../../../../assets/svg/${name}.svg`, import.meta.url)
      return url.href
    },
    updateTitle(newValue) {
      this.setTitle(newValue)
    }
  }
}
</script>

<style lang="less" scoped>
.node-common {
  position: relative;
  width: 100%;
  padding: 12px 16px 16px 16px;
  border-radius: 8px;
  background: #fff;
  box-sizing: border-box;

  .node-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 24px;
    margin-bottom: 8px;
  }
  .node-title {
    height: 24px;
    line-height: 24px;
    font-size: 16px;
    font-weight: 600;
    color: var(--wf-color-text-1);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .node-header-left {
    flex: 1;
    display: flex;
    align-items: center;
    overflow: hidden;
  }
  .node-icon {
    display: flex;
    margin-right: 8px;
    img {
      width: 20px;
      height: 20px;
    }
  }
  .node-menu-wrapper {
    position: relative;
    .node-menu-btn {
      display: flex;
      justify-content: center;
      align-items: center;
      width: 24px;
      height: 24px;
      border-radius: 6px;
      cursor: pointer;

      &:hover {
        background-color: #e4e6eb;
      }

      .btn-icon {
        width: 16px;
        height: 16px;
      }
    }
    .node-menus {
      position: absolute;
      top: 26px;
      right: 4px;
      padding: 2px;
      z-index: 100;
      border-radius: 6px;
      background: #fff;
      box-shadow:
        0 6px 30px 5px #0000000d,
        0 16px 24px 2px #0000000a,
        0 8px 10px -5px #00000014;
      .node-menu {
        height: 32px;
        line-height: 32px;
        padding: 0 16px;
        font-size: 14px;
        font-weight: 400;
        color: #595959;
        white-space: nowrap;
        cursor: pointer;

        &:hover {
          border-radius: 6px;
          background-color: #e4e6eb;
        }
      }

      .del-btn {
        color: #fb363f;
      }
    }
  }

  .node-title-input {
    border-width: 0;
    color: #262626;
    font-size: 16px;
    font-weight: 600;
    margin-right: 8px;
  }
  .node-title-input:hover {
    border-width: 1px;
  }
}

.node-common.isHovered {
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    border-radius: 8px;
    border: 1px solid #2475fc;
  }
}
.node-common.isSelected {
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    border-radius: 8px;
    border: 2px solid #2475fc;
  }
}


.node-common-wrapper.node-common-wrapper_explain-node{
  .node-icon{
    display: none;
  }
  .node-title-input{
    padding-left: 0;
    font-size: 12px;
    color: #7a8699;
    font-weight: 400;
    background: #FFEFD6;
  }
  .node-common{
    background: #FFEFD6;
    padding: 16px;
  }
  .node-common.isHovered {
    &::before {
      border: 1px solid #BFBFBF;
    }
  }
    .node-common.isSelected {
    &::before {
      border: 2px solid #f90;
    }
  }
}

</style>
