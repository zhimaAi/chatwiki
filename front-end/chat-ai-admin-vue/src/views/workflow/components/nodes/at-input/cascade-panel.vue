<template>
  <div class="cascade-panel">
    <div class="cascade-panel-container">
      <!-- 级联面板列表 -->
      <div class="cascade-panel-list" v-for="(level, levelIndex) in activePath" :key="levelIndex">
        <!-- <div class="cascade-panel-header" v-if="level.title">
          {{ level.title }}
        </div> -->
        <div class="cascade-panel-body">
          <div
            class="cascade-panel-item"
            v-for="item in level.options"
            :key="item.value"
            :class="{ 'active': isActive(item, levelIndex) }"
            @click="handleSelect(item, levelIndex)"
          >
            <div class="item-content">
              <span class="item-label">{{ item.label }}</span>
              <span class="item-icon" v-if="hasChildren(item)">
                <svg viewBox="0 0 1024 1024" width="12" height="12">
                  <path d="M765.7 486.8L314.9 134.7c-5.3-4.1-12.9-0.4-12.9 6.3v77.3c0 4.9 2.3 9.6 6.1 12.6l360 281.1-360 281.1c-3.9 3-6.1 7.7-6.1 12.6V883c0 6.7 7.7 10.4 12.9 6.3l450.8-352.1c16.4-12.8 16.4-37.6 0-50.4z" fill="currentColor"></path>
                </svg>
              </span>
            </div>
          </div>
          <div class="cascade-panel-empty" v-if="level.options && level.options.length === 0">
            暂无数据
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "CascadePanel",
  props: {
    // 级联选项数据
    options: {
      type: Array,
      default: () => []
    },
    // 默认选中的值
    value: {
      type: [String, Number, Array],
      default: ""
    },
    // 是否允许选择任意一级
    checkAnyLevel: {
      type: Boolean,
      default: false
    },
    // 子选项的属性名
    childrenKey: {
      type: String,
      default: "children"
    },
    // 值的属性名
    valueKey: {
      type: String,
      default: "value"
    },
    // 标签的属性名
    labelKey: {
      type: String,
      default: "label"
    }
  },
  data() {
    return {
      // 当前激活的路径
      activePath: [],
      // 当前选中的值路径
      selectedValuePath: [],
      // 当前选中的选项路径
      selectedPath: []
    };
  },
  watch: {
    options: {
      handler() {
        this.initPanel();
      },
      immediate: true
    },
    value: {
      handler(val) {
        if (val) {
          this.initSelectedPath();
        }
      },
      immediate: true
    }
  },
  methods: {
    // 初始化面板
    initPanel() {
      if (!this.options || this.options.length === 0) {
        this.activePath = [];
        return;
      }
      
      // 初始化第一级
      this.activePath = [{
        title: "",
        options: this.options,
        level: 0
      }];
      
      // 如果有默认值，初始化选中路径
      if (this.value) {
        this.initSelectedPath();
      }
    },
    
    // 初始化选中路径
    initSelectedPath() {
      // 如果值是数组，表示已经是完整路径
      if (Array.isArray(this.value)) {
        this.selectedValuePath = [...this.value];
        this.findPathByValues(this.selectedValuePath);
      } else {
        // 如果是单个值，需要查找完整路径
        this.findPathByValue(this.value);
      }
    },
    
    // 根据单个值查找路径
    findPathByValue(value) {
      const path = [];
      const valuePath = [];
      
      // 递归查找路径
      const find = (options, level) => {
        if (!options || options.length === 0) return false;
        
        for (let i = 0; i < options.length; i++) {
          const item = options[i];
          const itemValue = item[this.valueKey];
          
          // 找到匹配项
          if (itemValue === value) {
            path[level] = item;
            valuePath[level] = itemValue;
            return true;
          }
          
          // 递归查找子项
          if (item[this.childrenKey] && item[this.childrenKey].length > 0) {
            path[level] = item;
            valuePath[level] = itemValue;
            
            if (find(item[this.childrenKey], level + 1)) {
              return true;
            }
            
            // 未找到，回溯
            path[level] = null;
            valuePath[level] = null;
          }
        }
        
        return false;
      };
      
      find(this.options, 0);
      
      // 过滤掉空值
      this.selectedPath = path.filter(item => item);
      this.selectedValuePath = valuePath.filter(item => item != null);
      
      // 更新激活路径
      this.updateActivePath();
    },
    
    // 根据值数组查找路径
    findPathByValues(values) {
      if (!values || values.length === 0) return;
      
      const path = [];
      let currentOptions = this.options;
      
      // 遍历值数组，查找对应的选项
      for (let i = 0; i < values.length; i++) {
        const value = values[i];
        const found = currentOptions.find(item => item[this.valueKey] === value);
        
        if (found) {
          path.push(found);
          currentOptions = found[this.childrenKey] || [];
        } else {
          break;
        }
      }
      
      this.selectedPath = path;
      
      // 更新激活路径
      this.updateActivePath();
    },
    
    // 更新激活路径
    updateActivePath() {
      if (!this.selectedPath || this.selectedPath.length === 0) {
        this.activePath = [{
          title: "",
          options: this.options,
          level: 0
        }];
        return;
      }
      
      // 重置激活路径
      this.activePath = [];
      
      // 第一级
      this.activePath.push({
        title: "",
        options: this.options,
        level: 0
      });
      
      // 根据选中路径构建激活路径
      for (let i = 0; i < this.selectedPath.length; i++) {
        const item = this.selectedPath[i];
        
        // 如果有子选项，添加到激活路径
        if (item[this.childrenKey] && item[this.childrenKey].length > 0) {
          this.activePath.push({
            title: item[this.labelKey],
            options: item[this.childrenKey],
            level: i + 1
          });
        }
      }
    },
    
    // 处理选择
    handleSelect(item, levelIndex) {
      // 更新选中路径
      this.selectedPath = this.selectedPath.slice(0, levelIndex);
      this.selectedPath.push(item);
      
      // 更新选中值路径
      this.selectedValuePath = this.selectedValuePath.slice(0, levelIndex);
      this.selectedValuePath.push(item[this.valueKey]);
      
      // 更新激活路径
      this.activePath = this.activePath.slice(0, levelIndex + 1);
      
      // 如果有子选项，添加到激活路径
      if (this.hasChildren(item)) {
        this.activePath.push({
          title: item[this.labelKey],
          options: item[this.childrenKey],
          level: levelIndex + 1
        });
      }
      
      // 发送选中事件
      this.$emit('select', item, this.selectedPath, levelIndex);
      
      // 如果允许选择任意一级，或者是叶子节点，发送change事件
      if (this.checkAnyLevel || !this.hasChildren(item)) {
        this.$emit('change', item[this.valueKey], this.selectedValuePath, this.selectedPath);
        this.$emit('input', this.checkAnyLevel ? this.selectedValuePath : item[this.valueKey]);
      }
    },
    
    // 判断是否是激活项
    isActive(item, levelIndex) {
      if (!this.selectedPath || this.selectedPath.length <= levelIndex) {
        return false;
      }
      return item[this.valueKey] === this.selectedPath[levelIndex][this.valueKey];
    },
    
    // 判断是否有子选项
    hasChildren(item) {
      return item[this.childrenKey] && item[this.childrenKey].length > 0;
    },
    
    // 重置组件
    reset() {
      this.selectedPath = [];
      this.selectedValuePath = [];
      this.initPanel();
    }
  }
};
</script>

<style lang="less" scoped>
.cascade-panel {
  width: 100%;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 6px 16px 0 rgba(0, 0, 0, 0.08), 0 3px 6px -4px rgba(0, 0, 0, 0.12), 0 9px 28px 8px rgba(0, 0, 0, 0.05);
  
  &-container {
    display: flex;
    max-height: 300px;
    overflow: hidden;
  }
  
  &-list {
    min-width: 180px;
    max-width: 300px;
    border-right: 1px solid #f0f0f0;
    display: flex;
    flex-direction: column;
    
    &:last-child {
      border-right: none;
    }
  }
  
  &-header {
    padding: 8px 12px;
    font-size: 14px;
    color: rgba(0, 0, 0, 0.45);
    border-bottom: 1px solid #f0f0f0;
    background-color: #fafafa;
  }
  
  &-body {
    flex: 1;
    overflow-y: auto;
    max-height: 300px;
    
    &::-webkit-scrollbar {
      width: 6px;
      height: 6px;
    }
    
    &::-webkit-scrollbar-track {
      background-color: #fafafa;
    }
    
    &::-webkit-scrollbar-thumb {
      border-radius: 3px;
      background: #c1c1c1;
    }
  }
  
  &-item {
    padding: 8px 12px;
    cursor: pointer;
    transition: background 0.3s;
    display: flex;
    align-items: center;
    justify-content: space-between;
    
    &:hover {
      background: rgba(0, 0, 0, 0.04);
    }
    
    &.active {
      background: #e6f7ff;
      color: #1890ff;
      font-weight: 500;
    }
    
    .item-content {
      display: flex;
      align-items: center;
      justify-content: space-between;
      width: 100%;
    }
    
    .item-label {
      flex: 1;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    
    .item-icon {
      margin-left: 8px;
      color: rgba(0, 0, 0, 0.45);
    }
  }
  
  &-empty {
    padding: 16px 12px;
    text-align: center;
    color: rgba(0, 0, 0, 0.45);
    font-size: 14px;
  }
}
</style>