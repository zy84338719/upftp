<template>
  <div class="tree-node-wrapper">
    <div
      class="tree-node"
      :class="{ active: node.path === currentPath }"
      :style="{ paddingLeft: 8 + depth * 14 + 'px' }"
      @click="handleClick"
    >
      <span class="tree-arrow" :class="{ open: isOpen }" v-if="hasChildren">
        ▶
      </span>
      <span class="tree-arrow" v-else></span>
      <span class="tree-icon">{{ node.isDir ? '📁' : '📄' }}</span>
      <span class="tree-label">{{ node.name || '/' }}</span>
    </div>
    <div v-if="hasChildren && isOpen" class="tree-children">
      <TreeNode
        v-for="child in node.children"
        :key="child.path"
        :node="child"
        :current-path="currentPath"
        :depth="depth + 1"
        @navigate="$emit('navigate', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { TreeNode } from '@/stores/app'

interface Props {
  node: TreeNode
  currentPath: string
  depth?: number
}

const props = withDefaults(defineProps<Props>(), {
  depth: 0
})

const emit = defineEmits<{
  navigate: [path: string]
}>()

const hasChildren = computed(() => {
  return props.node.children && props.node.children.length > 0
})

const isOpen = ref(
  props.currentPath === props.node.path ||
    props.currentPath.startsWith(props.node.path + '/')
)

function handleClick(e: Event) {
  e.stopPropagation()
  if (hasChildren.value) {
    isOpen.value = !isOpen.value
  }
  emit('navigate', props.node.path)
}
</script>

<style scoped>
.tree-node-wrapper {
  width: 100%;
}

.tree-node {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  cursor: pointer;
  font-size: 12px;
  color: #555;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  transition: background 0.1s;
  border-radius: 4px;
  margin: 1px 6px;
}

.tree-node:hover {
  background: #f0f0f0;
  color: #222;
}

.tree-node.active {
  background: #fef3c7;
  color: #92400e;
  font-weight: 600;
}

.tree-arrow {
  width: 14px;
  text-align: center;
  font-size: 9px;
  color: #bbb;
  flex-shrink: 0;
  transition: transform 0.15s;
}

.tree-arrow.open {
  transform: rotate(90deg);
}

.tree-icon {
  flex-shrink: 0;
  font-size: 14px;
}

.tree-label {
  overflow: hidden;
  text-overflow: ellipsis;
}

.tree-children {
  width: 100%;
}
</style>
