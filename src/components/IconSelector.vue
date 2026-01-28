<template>
  <div class="icon-selector">
    <div class="icon-tabs">
      <button
        v-if="defaultIconPath"
        :class="['tab-btn', { active: iconMode === 'default' }]"
        @click="iconMode = 'default'"
      >
        默认图标
      </button>
      <button
        :class="['tab-btn', { active: iconMode === 'upload' }]"
        @click="iconMode = 'upload'"
      >
        上传图标
      </button>
      <button
        :class="['tab-btn', { active: iconMode === 'url' }]"
        @click="iconMode = 'url'"
      >
        URL图标
      </button>
    </div>

    <!-- 默认图标 -->
    <div v-if="iconMode === 'default' && defaultIconPath" class="icon-content">
      <div class="default-icon-preview">
        <img :src="defaultIconPath" alt="默认图标" class="icon-preview" />
        <p class="icon-hint">使用默认本地图标：{{ defaultIconPath }}</p>
      </div>
      <button @click="selectDefault" class="select-btn">使用默认图标</button>
    </div>

    <!-- 上传图标 -->
    <div v-if="iconMode === 'upload'" class="icon-content">
      <div class="upload-area">
        <input
          type="file"
          ref="fileInput"
          @change="handleFileSelect"
          accept="image/*"
          style="display: none"
        />
        <div v-if="!uploadedIconUrl" class="upload-placeholder">
          <p>点击选择图标文件</p>
          <p class="hint">支持 jpg、png、jpeg、webp、avif、svg、ico 等格式</p>
        </div>
        <div v-else class="uploaded-preview">
          <img :src="uploadedIconUrl" alt="上传的图标" class="icon-preview" />
          <p class="icon-hint">已上传的图标</p>
        </div>
        <button @click="$refs.fileInput.click()" class="upload-btn">
          {{ uploadedIconUrl ? '重新选择' : '选择文件' }}
        </button>
        <div v-if="uploading" class="upload-status">上传中...</div>
      </div>
    </div>

    <!-- URL图标 -->
    <div v-if="iconMode === 'url'" class="icon-content">
      <div class="url-input-group">
        <label>图标URL</label>
        <input
          v-model="iconUrl"
          type="text"
          placeholder="https://example.com/favicon.ico"
          class="url-input"
          @input="handleUrlInput"
        />
        <small class="form-hint">支持重定向图片URL，格式：avif, png, jpg, jpeg, webp, svg, ico等</small>
        <div v-if="urlValidating" class="url-status">
          <i class="fas fa-spinner fa-spin"></i>
          <span>正在验证图片URL...</span>
        </div>
        <div v-if="urlError" class="url-error">
          <i class="fas fa-exclamation-circle"></i>
          <span>{{ urlError }}</span>
        </div>
      </div>
      <div v-if="iconUrl && !urlError && urlValidated" class="url-preview">
        <img :src="validatedUrl" alt="URL图标" class="icon-preview" @error="handleImageError" />
        <p class="icon-hint">URL图标预览（支持重定向）</p>
      </div>
      <button @click="selectUrl" class="select-btn" :disabled="!iconUrl || urlValidating || !!urlError">
        {{ urlValidating ? '验证中...' : '使用URL图标' }}
      </button>
    </div>

    <!-- 当前选择的图标预览 -->
    <div v-if="currentIcon" class="current-icon">
      <p class="current-label">当前图标：</p>
      <img :src="currentIcon" alt="当前图标" class="icon-preview" @error="handleImageError" />
      <p class="icon-hint">{{ currentIcon }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { adminAPI } from '../api'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  defaultIconPath: {
    type: String,
    default: '/favicon.ico',
  },
})

const emit = defineEmits(['update:modelValue'])

const iconMode = ref('default')
const iconUrl = ref('')
const uploadedIconUrl = ref('')
const uploading = ref(false)
const currentIcon = ref(props.modelValue || props.defaultIconPath)
const urlValidating = ref(false)
const urlValidated = ref(false)
const urlError = ref('')
const validatedUrl = ref('')

// 支持的图片格式
const allowedImageFormats = ['avif', 'png', 'jpg', 'jpeg', 'webp', 'svg', 'ico', 'gif', 'bmp']

// 监听外部值变化
watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal) {
      currentIcon.value = newVal
      // 根据值判断模式
      if (props.defaultIconPath && newVal === props.defaultIconPath) {
        iconMode.value = 'default'
      } else if (newVal.startsWith('http://') || newVal.startsWith('https://')) {
        iconMode.value = 'url'
        iconUrl.value = newVal
      } else if (newVal.startsWith('/uploads/')) {
        iconMode.value = 'upload'
        uploadedIconUrl.value = newVal
      } else {
        // 如果值不为空但不是已知格式，默认使用URL模式
        iconMode.value = 'url'
        iconUrl.value = newVal
      }
    } else {
      currentIcon.value = props.defaultIconPath || ''
    }
  },
  { immediate: true }
)

const handleFileSelect = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  // 检查文件类型，支持常见图标格式包括ico
  const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/svg+xml', 'image/x-icon', 'image/vnd.microsoft.icon', 'image/ico', 'image/icon']
  const fileExtension = file.name.split('.').pop()?.toLowerCase()
  const allowedExtensions = ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg', 'ico', 'avif', 'bmp']
  
  if (!allowedTypes.includes(file.type) && !allowedExtensions.includes(fileExtension)) {
    alert('不支持的文件格式，请上传 jpg、png、gif、webp、svg、ico 等格式的图片')
    return
  }

  uploading.value = true
  try {
    const res = await adminAPI.uploadFile(file)
    uploadedIconUrl.value = res.data.url
    currentIcon.value = res.data.url
    emit('update:modelValue', res.data.url)
  } catch (error) {
    alert('上传失败: ' + (error.response?.data?.error || error.message))
  } finally {
    uploading.value = false
  }
}

const selectDefault = () => {
  currentIcon.value = props.defaultIconPath
  emit('update:modelValue', props.defaultIconPath)
}

// 验证图片URL（支持重定向）
const validateImageUrl = async (url) => {
  if (!url) {
    urlError.value = ''
    urlValidated.value = false
    return false
  }

  urlValidating.value = true
  urlError.value = ''
  urlValidated.value = false

  try {
    // 检查URL格式
    if (!url.startsWith('http://') && !url.startsWith('https://')) {
      urlError.value = 'URL必须以http://或https://开头'
      urlValidating.value = false
      return false
    }

    // 使用fetch跟随重定向
    const response = await fetch(url, {
      method: 'HEAD',
      mode: 'cors',
      redirect: 'follow'
    })

    if (!response.ok) {
      urlError.value = '无法访问该URL'
      urlValidating.value = false
      return false
    }

    // 获取最终URL（可能经过重定向）
    const finalUrl = response.url || url
    
    // 检查Content-Type
    const contentType = response.headers.get('content-type') || ''
    const isImage = contentType.startsWith('image/')
    
    // 或者检查URL扩展名
    const urlLower = finalUrl.toLowerCase()
    const hasImageExtension = allowedImageFormats.some(format => 
      urlLower.includes(`.${format}`) || urlLower.includes(`/${format}`)
    )

    if (!isImage && !hasImageExtension) {
      urlError.value = 'URL指向的不是图片文件'
      urlValidating.value = false
      return false
    }

    // 使用Image对象进一步验证（支持重定向）
    return new Promise((resolve) => {
      const img = new Image()
      img.crossOrigin = 'anonymous'
      
      img.onload = () => {
        validatedUrl.value = finalUrl
        urlValidated.value = true
        urlError.value = ''
        urlValidating.value = false
        resolve(true)
      }
      
      img.onerror = () => {
        // 即使Image加载失败，如果Content-Type正确，也允许使用
        if (isImage || hasImageExtension) {
          validatedUrl.value = finalUrl
          urlValidated.value = true
          urlError.value = ''
          urlValidating.value = false
          resolve(true)
        } else {
          urlError.value = '无法加载图片，请检查URL是否正确'
          urlValidated.value = false
          urlValidating.value = false
          resolve(false)
        }
      }
      
      img.src = finalUrl
    })
  } catch (error) {
    urlError.value = '验证失败: ' + (error.message || '网络错误')
    urlValidated.value = false
    urlValidating.value = false
    return false
  }
}

// 处理URL输入（防抖）
let urlValidationTimer = null
const handleUrlInput = () => {
  urlValidated.value = false
  urlError.value = ''
  
  if (urlValidationTimer) {
    clearTimeout(urlValidationTimer)
  }
  
  urlValidationTimer = setTimeout(() => {
    if (iconUrl.value) {
      validateImageUrl(iconUrl.value)
    }
  }, 500)
}

const selectUrl = async () => {
  if (iconUrl.value && !urlError.value) {
    // 如果还未验证，先验证
    if (!urlValidated.value) {
      const isValid = await validateImageUrl(iconUrl.value)
      if (!isValid) {
        return
      }
    }
    
    currentIcon.value = validatedUrl.value || iconUrl.value
    emit('update:modelValue', validatedUrl.value || iconUrl.value)
  }
}

const handleImageError = () => {
  if (iconMode.value === 'url') {
    urlError.value = '无法加载图片，请检查URL是否正确'
    urlValidated.value = false
  }
  console.warn('无法加载图标URL:', iconUrl.value)
}

onMounted(() => {
  // 初始化时根据当前值设置模式
  if (props.modelValue) {
    if (props.defaultIconPath && props.modelValue === props.defaultIconPath) {
      iconMode.value = 'default'
    } else if (
      props.modelValue.startsWith('http://') ||
      props.modelValue.startsWith('https://')
    ) {
      iconMode.value = 'url'
      iconUrl.value = props.modelValue
    } else if (props.modelValue.startsWith('/uploads/')) {
      iconMode.value = 'upload'
      uploadedIconUrl.value = props.modelValue
    } else {
      // 如果值不为空但不是已知格式，默认使用URL模式
      iconMode.value = 'url'
      iconUrl.value = props.modelValue
    }
  } else if (!props.defaultIconPath) {
    // 如果没有默认图标路径，默认使用上传模式
    iconMode.value = 'upload'
  }
})
</script>

<style scoped>
.icon-selector {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 15px;
  background: rgba(var(--background-color-rgb), 0.3);
}

.icon-tabs {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  border-bottom: 2px solid var(--border-color);
}

.tab-btn {
  padding: 8px 16px;
  background: transparent;
  border: none;
  border-bottom: 2px solid transparent;
  color: var(--text-color);
  cursor: pointer;
  transition: all 0.3s;
  font-size: 14px;
}

.tab-btn.active {
  border-bottom-color: #007aff;
  color: #007aff;
  font-weight: bold;
}

.icon-content {
  min-height: 200px;
}

.default-icon-preview,
.uploaded-preview,
.url-preview {
  text-align: center;
  padding: 20px;
  background: rgba(var(--background-color-rgb), 0.5);
  border-radius: 8px;
  margin-bottom: 15px;
}

.icon-preview {
  width: 64px;
  height: 64px;
  object-fit: contain;
  margin: 0 auto 10px;
  display: block;
}

.icon-hint {
  font-size: 12px;
  color: #999;
  margin: 0;
}

.upload-area {
  text-align: center;
}

.upload-placeholder {
  padding: 40px;
  background: rgba(var(--background-color-rgb), 0.5);
  border-radius: 8px;
  border: 2px dashed var(--border-color);
  margin-bottom: 15px;
}

.upload-placeholder p {
  margin: 5px 0;
  color: var(--text-color);
}

.upload-placeholder .hint {
  font-size: 12px;
  color: #999;
}

.upload-btn,
.select-btn {
  padding: 10px 20px;
  background: #007aff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  margin-top: 10px;
}

.upload-btn:hover,
.select-btn:hover {
  background: #0056b3;
}

.select-btn:disabled {
  background: #999;
  cursor: not-allowed;
}

.upload-status {
  margin-top: 10px;
  color: #007aff;
  font-size: 14px;
}

.url-input-group {
  margin-bottom: 15px;
}

.url-input-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
  color: var(--text-color);
}

.url-input {
  width: 100%;
  padding: 8px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: rgba(var(--background-color-rgb), 0.5);
  color: var(--text-color);
  font-size: 14px;
}

.form-hint {
  display: block;
  margin-top: 5px;
  font-size: 12px;
  color: #999;
}

.url-status {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
  color: #007aff;
  font-size: 13px;
}

.url-error {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
  color: #f44336;
  font-size: 13px;
  padding: 8px;
  background: rgba(244, 67, 54, 0.1);
  border-radius: 4px;
  border: 1px solid rgba(244, 67, 54, 0.2);
}

.current-icon {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid var(--border-color);
  text-align: center;
}

.current-label {
  font-size: 14px;
  font-weight: bold;
  color: var(--text-color);
  margin-bottom: 10px;
}
</style>
