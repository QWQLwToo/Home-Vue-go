<template>
  <div class="icon-picker">
    <div class="icon-picker-header">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="搜索图标..."
        class="icon-search"
        @input="filterIcons"
      />
    </div>
    <div class="icon-picker-body">
      <div class="icon-categories">
        <button
          v-for="category in categories"
          :key="category.name"
          @click="activeCategory = category.name"
          :class="['category-btn', { active: activeCategory === category.name }]"
        >
          <i :class="category.icon"></i>
          <span>{{ category.label }}</span>
        </button>
      </div>
      <div class="icon-grid" ref="iconGrid">
        <div
          v-for="icon in filteredIcons"
          :key="icon"
          @click="selectIcon(icon)"
          :class="['icon-item', { active: modelValue === icon }]"
          :title="icon"
        >
          <i :class="icon"></i>
          <span class="icon-name">{{ getIconName(icon) }}</span>
        </div>
      </div>
    </div>
    <div class="icon-picker-footer" v-if="modelValue">
      <div class="selected-icon">
        <span>已选择：</span>
        <i :class="modelValue"></i>
        <code>{{ modelValue }}</code>
      </div>
      <div class="icon-picker-actions">
        <button @click="clearIcon" class="btn-clear">
          <i class="fas fa-times"></i>
          清除
        </button>
        <button @click="closePicker" class="btn-close">
          <i class="fas fa-check"></i>
          确定
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['update:modelValue', 'close'])

const searchQuery = ref('')
const activeCategory = ref('all')

// Font Awesome 常用图标分类
const categories = [
  { name: 'all', label: '全部', icon: 'fas fa-th' },
  { name: 'web', label: '网页', icon: 'fas fa-globe' },
  { name: 'social', label: '社交', icon: 'fas fa-share-alt' },
  { name: 'media', label: '媒体', icon: 'fas fa-photo-video' },
  { name: 'business', label: '商业', icon: 'fas fa-briefcase' },
  { name: 'tech', label: '技术', icon: 'fas fa-code' },
  { name: 'other', label: '其他', icon: 'fas fa-ellipsis-h' },
]

// 常用图标列表（按分类）
const iconLibrary = {
  web: [
    'fas fa-home', 'fas fa-globe', 'fas fa-link', 'fas fa-external-link-alt',
    'fas fa-bookmark', 'fas fa-star', 'fas fa-heart', 'fas fa-thumbs-up',
  ],
  social: [
    'fab fa-github', 'fab fa-twitter', 'fab fa-facebook', 'fab fa-instagram',
    'fab fa-linkedin', 'fab fa-youtube', 'fab fa-telegram', 'fab fa-discord',
    'fab fa-weixin', 'fab fa-qq', 'fab fa-weibo', 'fab fa-bilibili',
  ],
  media: [
    'fas fa-image', 'fas fa-video', 'fas fa-music', 'fas fa-film',
    'fas fa-camera', 'fas fa-microphone', 'fas fa-headphones',
  ],
  business: [
    'fas fa-briefcase', 'fas fa-building', 'fas fa-chart-line', 'fas fa-dollar-sign',
    'fas fa-shopping-cart', 'fas fa-credit-card', 'fas fa-handshake',
  ],
  tech: [
    'fas fa-code', 'fas fa-terminal', 'fas fa-server', 'fas fa-database',
    'fas fa-cloud', 'fas fa-mobile-alt', 'fas fa-laptop', 'fas fa-keyboard',
  ],
  other: [
    'fas fa-envelope', 'fas fa-phone', 'fas fa-map-marker-alt', 'fas fa-calendar',
    'fas fa-clock', 'fas fa-bell', 'fas fa-cog', 'fas fa-user', 'fas fa-users',
  ],
}

// 所有图标
const allIcons = computed(() => {
  const icons = []
  Object.values(iconLibrary).forEach(categoryIcons => {
    icons.push(...categoryIcons)
  })
  return icons
})

// 过滤后的图标
const filteredIcons = computed(() => {
  let icons = activeCategory.value === 'all' 
    ? allIcons.value 
    : iconLibrary[activeCategory.value] || []

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    icons = icons.filter(icon => 
      icon.toLowerCase().includes(query) || 
      getIconName(icon).toLowerCase().includes(query)
    )
  }

  return icons
})

const getIconName = (icon) => {
  // 从 "fas fa-home" 提取 "home"
  const parts = icon.split(' ')
  return parts[parts.length - 1] || icon
}

const selectIcon = (icon) => {
  emit('update:modelValue', icon)
}

const clearIcon = () => {
  emit('update:modelValue', '')
}

const closePicker = () => {
  emit('close')
}

const filterIcons = () => {
  // 搜索时自动切换到"全部"分类
  if (searchQuery.value.trim() && activeCategory.value !== 'all') {
    activeCategory.value = 'all'
  }
}

watch(() => props.modelValue, (newVal) => {
  if (newVal) {
    // 如果选择了图标，可以高亮显示
  }
})
</script>

<style scoped>
.icon-picker {
  display: flex;
  flex-direction: column;
  height: 100%;
  max-height: 100%;
  background: rgba(var(--background-color-rgb), 0.98);
  border-radius: var(--border-radius);
  overflow: hidden;
}

.icon-picker-header {
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
}

.icon-search {
  width: 100%;
  padding: 10px 16px;
  border: 2px solid var(--border-color);
  border-radius: 8px;
  background: rgba(var(--background-color-rgb), 0.6);
  color: var(--text-color);
  font-size: 14px;
  transition: all 0.3s ease;
}

.icon-search:focus {
  outline: none;
  border-color: #007aff;
  background: rgba(var(--background-color-rgb), 0.8);
}

.icon-picker-body {
  flex: 1;
  display: flex;
  overflow: hidden;
  min-height: 0;
}

.icon-categories {
  width: 180px;
  padding: 16px;
  border-right: 1px solid var(--border-color);
  overflow-y: auto;
  overflow-x: hidden;
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-height: 0;
  flex-shrink: 0;
  /* 自定义滚动条样式 */
  scrollbar-width: thin;
  scrollbar-color: var(--hover-link-color) rgba(var(--background-color-rgb), 0.3);
}

.icon-categories::-webkit-scrollbar {
  width: 6px;
}

.icon-categories::-webkit-scrollbar-track {
  background: rgba(var(--background-color-rgb), 0.3);
  border-radius: 3px;
}

.icon-categories::-webkit-scrollbar-thumb {
  background: var(--hover-link-color);
  border-radius: 3px;
  transition: background 0.3s ease;
}

.icon-categories::-webkit-scrollbar-thumb:hover {
  background: #ffd700;
}

.category-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: rgba(var(--background-color-rgb), 0.6);
  color: var(--text-color);
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: left;
}

.category-btn:hover {
  background: rgba(var(--background-color-rgb), 0.8);
  border-color: var(--hover-link-color);
}

.category-btn.active {
  background: var(--hover-link-color);
  color: #333;
  border-color: var(--hover-link-color);
  font-weight: 500;
}

.icon-grid {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
  overflow-x: hidden;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 12px;
  align-content: start;
  min-height: 0;
  /* 自定义滚动条样式 */
  scrollbar-width: thin;
  scrollbar-color: var(--hover-link-color) rgba(var(--background-color-rgb), 0.3);
}

.icon-grid::-webkit-scrollbar {
  width: 8px;
}

.icon-grid::-webkit-scrollbar-track {
  background: rgba(var(--background-color-rgb), 0.3);
  border-radius: 4px;
}

.icon-grid::-webkit-scrollbar-thumb {
  background: var(--hover-link-color);
  border-radius: 4px;
  transition: background 0.3s ease;
}

.icon-grid::-webkit-scrollbar-thumb:hover {
  background: #ffd700;
}

.icon-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 16px 8px;
  border: 2px solid var(--border-color);
  border-radius: 8px;
  background: rgba(var(--background-color-rgb), 0.6);
  cursor: pointer;
  transition: all 0.3s ease;
  min-height: 100px;
}

.icon-item:hover {
  background: rgba(var(--background-color-rgb), 0.8);
  border-color: var(--hover-link-color);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px var(--shadow-color);
}

.icon-item.active {
  background: var(--hover-link-color);
  border-color: var(--hover-link-color);
  color: #333;
}

.icon-item i {
  font-size: 24px;
  margin-bottom: 8px;
  color: inherit;
}

.icon-item.active i {
  color: #333;
}

.icon-name {
  font-size: 11px;
  text-align: center;
  word-break: break-all;
  color: inherit;
  opacity: 0.8;
}

.icon-picker-footer {
  padding: 16px;
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  background: rgba(var(--background-color-rgb), 0.98);
  position: sticky;
  bottom: 0;
  z-index: 10;
  flex-shrink: 0;
}

.selected-icon {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  font-size: 14px;
  color: var(--text-color);
}

.selected-icon i {
  font-size: 20px;
  color: var(--hover-link-color);
}

.selected-icon code {
  background: rgba(var(--background-color-rgb), 0.8);
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-family: 'Courier New', monospace;
}

.icon-picker-actions {
  display: flex;
  gap: 8px;
}

.btn-clear,
.btn-close {
  padding: 8px 16px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: rgba(var(--background-color-rgb), 0.8);
  color: var(--text-color);
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
}

.btn-clear:hover {
  background: rgba(244, 67, 54, 0.1);
  border-color: #f44336;
  color: #f44336;
}

.btn-close {
  background: var(--hover-link-color);
  color: #333;
  border-color: var(--hover-link-color);
  font-weight: 500;
}

.btn-close:hover {
  background: #ffd700;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(255, 204, 0, 0.3);
}

@media (max-width: 768px) {
  .icon-picker-body {
    flex-direction: column;
  }
  
  .icon-categories {
    width: 100%;
    flex-direction: row;
    overflow-x: auto;
    border-right: none;
    border-bottom: 1px solid var(--border-color);
  }
  
  .icon-grid {
    grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  }
}
</style>
