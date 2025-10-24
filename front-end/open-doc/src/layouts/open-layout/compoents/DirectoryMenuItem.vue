<style lang="less" scoped>
.directory-item {
  .directory-item-body {
    display: flex;
    align-items: center;
    height: 32px;
    padding: 0 8px;
    margin-bottom: 4px;
    border-radius: 6px;
    transition: all 0.2s;
  }

  .directory-item-body:hover {
    background-color: #f2f4f7;
    cursor: pointer;
  }

  .directory-expanded-btn {
    width: 24px;
    height: 24px;
    line-height: 24px;
    margin-right: 4px;
    border-radius: 6px;
    text-align: center;
    transition: all 0.2s;
  }
  .directory-expanded-btn:hover {
    background-color: #e4e6eb;
  }
  .doc-body {
    display: flex;
    align-items: center;
    flex: 1;
    height: 100%;
    overflow: hidden;
  }
  .doc-icon {
    margin-right: 4px;
  }
  .doc-icon.active-icon {
    display: none;
  }
  .doc-name {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 14px;
    color: #595959;
  }

  .directory-item-indent {
    align-self: stretch;
    white-space: nowrap;
    user-select: none;
  }

  .directory-item-indent .directory-item-indent-unit {
    display: inline-block;
    width: 24px;
  }
  .directory-expanded-btn .expanded-icon {
    transition: all 0.2s;
  }
  .directory-expanded-btn.directory-expanded-noop:hover {
    background: none;
  }
  .directory-expanded-btn.directory-expanded-noop .expanded-icon {
    display: none;
  }
}
.directory-item.is-expanded {
  & > .directory-item-body {
    .directory-expanded-btn {
      .expanded-icon {
        transform: rotate(90deg);
      }
    }
  }
}

.directory-item.active > .directory-item-body {
  background-color: #e6efff;
}
.directory-item.active > .directory-item-body .doc-name {
  color: #2475fc;
}
.directory-item.active {
  & > .directory-item-body .doc-icon.default-icon {
    display: none;
  }

  & > .directory-item-body .doc-icon.active-icon {
    display: block;
  }
}

@media (max-width: 992px) {
  .directory-item .directory-item-body {
    height: 44px;
  }
}
</style>

<template>
  <div
    :doc-key="props.docKey"
    class="directory-item"
    :class="{ 'is-expanded': isExpanded, active: props.activeKey == props.docKey }"
  >
    <div class="directory-item-body">
      <!-- 子菜单缩进 -->
      <div class="directory-item-indent" :style="{ 'padding-left': paddingLeft + 'px' }"></div>
      <span
        class="directory-expanded-btn"
        :class="{ 'directory-expanded-noop': !hasChildren }"
        @click="toggleExpanded"
      >
        <img class="expanded-icon" src=" @/assets/img/directory_action.svg" alt="" />
      </span>
      <span @click="toggleExpanded" class="doc-body sidebar-link" v-if="isDir == 1">
        <span class="doc-icon">{{ props.icon }}</span>
        <span class="doc-name">{{ props.title }}</span>
      </span>
      <router-link :to="link" class="doc-body sidebar-link" v-else>
        <span class="doc-icon">{{ props.icon }}</span>
        <span class="doc-name">{{ props.title }}</span>
      </router-link>
    </div>
    <div class="directory-item-children" v-if="hasChildren && isExpanded">
      <DirectoryMenuItem
        v-for="item in props.children"
        :icon="item.doc_icon"
        :title="item.title"
        :key="item.doc_key"
        :level="props.level + 1"
        :token="props.token"
        :activeKey="props.activeKey"
        :previewKey="previewKey"
        :docKey="item.doc_key"
        :item="item"
        :isDir="item.is_dir"
        :children="item.children"
        :forceExpanded="forceExpanded"
        :forceCollapsed="forceCollapsed"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
const props = defineProps({
  level: {
    type: Number,
    default: 0,
  },
  icon: {
    type: String,
    default: '',
  },
  isDir: {
    type: String,
    default: '0',
  },
  title: {
    type: String,
    default: '',
  },
  children: {
    type: Array,
    default: () => [],
  },
  docKey: {
    type: String,
    default: '',
  },
  token: {
    type: String,
    default: '',
  },
  activeKey: {
    type: String,
    default: '',
  },
  previewKey: {
    type: String,
    default: '',
  },
  forceExpanded: {
    type: Boolean,
    default: false,
  },
  forceCollapsed: {
    type: Boolean,
    default: false,
  },
  iconTemplateConfig: {
    type: Object,
    default: () => ({}),
  }
})

const isExpanded = ref(true) // 默认展开

const paddingLeft = computed(() => {
  return props.level * 24
})

const link = computed(() => {
  return `/doc/${props.docKey}?${props.previewKey ? 'preview=' + props.previewKey : ''}${props.token.length ? '&token=' + props.token : ''}`
})

const hasChildren = computed(() => {
  return props.children && props.children.length > 0
})

const toggleExpanded = () => {
  if (hasChildren.value) {
    isExpanded.value = !isExpanded.value
  }
}

// 监听强制展开/收起状态
watch(
  () => props.forceExpanded,
  (newVal) => {
    if (newVal && hasChildren.value) {
      isExpanded.value = true
    }
  },
)

watch(
  () => props.forceCollapsed,
  (newVal) => {
    if (newVal && hasChildren.value) {
      isExpanded.value = false
    }
  },
)
</script>
