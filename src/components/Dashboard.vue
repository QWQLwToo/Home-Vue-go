<template>
  <div class="dashboard" :class="{ loading: isLoading }">
    <!-- 顶部关键指标区 -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(33, 150, 243, 0.1); color: #2196f3;">
          <i class="fas fa-eye"></i>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ animatedStats.totalViews }}</div>
          <div class="stat-label">总访问量</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(76, 175, 80, 0.1); color: #4caf50;">
          <i class="fas fa-users"></i>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ animatedStats.uniqueVisitors }}</div>
          <div class="stat-label">独立访客</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(255, 152, 0, 0.1); color: #ff9800;">
          <i class="fas fa-calendar-day"></i>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ animatedStats.todayViews }}</div>
          <div class="stat-label">今日访问</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(156, 39, 176, 0.1); color: #9c27b0;">
          <i class="fas fa-link"></i>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ animatedStats.totalSites }}</div>
          <div class="stat-label">站点数量</div>
        </div>
      </div>
    </div>

    <!-- 中部图表区域 -->
    <div class="charts-grid">
      <!-- 访问趋势图（主图表，权重更高） -->
      <div class="chart-card chart-card--primary">
        <div class="card-header">
          <h3><i class="fas fa-chart-line"></i> 访问趋势</h3>
          <select v-model="chartPeriod" @change="loadChartData" class="period-select">
            <option value="7">最近7天</option>
            <option value="30">最近30天</option>
            <option value="90">最近90天</option>
          </select>
        </div>
        <div class="card-body chart-body">
          <canvas ref="trendChart"></canvas>
        </div>
      </div>

      <!-- 访问来源（次要图表） -->
      <div class="chart-card chart-card--secondary">
        <div class="card-header">
          <h3><i class="fas fa-chart-pie"></i> 访问来源</h3>
        </div>
        <div class="card-body chart-body">
          <canvas ref="sourceChart"></canvas>
        </div>
      </div>
    </div>

    <!-- 底部列表与日志区域 - 现代化布局 -->
    <div class="bottom-grid-modern">
      <!-- 左侧：最近登录IP（小卡片） -->
      <div class="recent-login-card compact-card">
        <div class="card-header">
          <h3><i class="fas fa-sign-in-alt"></i> 最近登录IP</h3>
          <button @click="loadLoginHistory" class="refresh-logs-btn" title="刷新">
            <i class="fas fa-sync-alt" :class="{ 'fa-spin': loadingLoginHistory }"></i>
          </button>
        </div>
        <div class="card-body compact-body">
          <div class="recent-list login-history-list">
            <div
              v-for="(history, index) in loginHistory"
              :key="index"
              class="recent-item login-item"
            >
              <div
                class="recent-icon"
                :class="{ 'login-success': history.success, 'login-failed': !history.success }"
              >
                <i :class="history.success ? 'fas fa-check-circle' : 'fas fa-times-circle'"></i>
              </div>
              <div class="recent-content">
                <div class="recent-title">
                  <span class="login-ip">{{ history.ip }}</span>
                  <span v-if="history.location" class="login-location">{{ history.location }}</span>
                </div>
                <div class="recent-meta">
                  <span>{{ history.username }}</span>
                  <span>{{ history.loginTime }}</span>
                </div>
              </div>
            </div>
            <div v-if="loginHistory.length === 0" class="empty-state">
              <i class="fas fa-inbox"></i>
              <p>暂无登录记录</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧：后台日志（平铺，占更大空间） -->
      <div class="logs-card-full">
        <div class="card-header">
          <h3><i class="fas fa-terminal"></i> 后台日志</h3>
          <div class="logs-controls">
            <select v-model="logLines" @change="loadLogs" class="log-lines-select">
              <option value="50">最近50条</option>
              <option value="100">最近100条</option>
              <option value="200">最近200条</option>
              <option value="500">最近500条</option>
              <option value="all">全部</option>
            </select>
            <button @click="loadLogs" class="refresh-logs-btn" title="刷新日志">
              <i class="fas fa-sync-alt" :class="{ 'fa-spin': loadingLogs }"></i>
            </button>
            <button
              @click="autoScroll = !autoScroll"
              class="auto-scroll-btn"
              :class="{ active: autoScroll }"
              title="自动滚动"
            >
              <i class="fas fa-arrow-down"></i>
            </button>
          </div>
        </div>
        <div class="card-body logs-body-full" ref="logsContainer">
          <div class="logs-content">
            <div
              v-for="(log, index) in backendLogs"
              :key="index"
              class="log-line"
              :class="getLogLevel(log)"
            >
              <span class="log-content">{{ log }}</span>
            </div>
            <div v-if="backendLogs.length === 0" class="empty-state">
              <i class="fas fa-inbox"></i>
              <p>暂无日志</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { adminAPI } from '../api'

const stats = ref({
  totalViews: 0,
  uniqueVisitors: 0,
  todayViews: 0,
  totalSites: 0,
})

const animatedStats = ref({
  totalViews: 0,
  uniqueVisitors: 0,
  todayViews: 0,
  totalSites: 0,
})

const chartPeriod = ref(7)
const trendChart = ref(null)
const sourceChart = ref(null)
const backendLogs = ref([])
const logLines = ref('100')
const loadingLogs = ref(false)
const autoScroll = ref(true)
const logsContainer = ref(null)
const loginHistory = ref([])
const loadingLoginHistory = ref(false)
const isLoading = ref(true)
let trendChartInstance = null
let sourceChartInstance = null

// 简单的图表实现（使用Canvas绘制，不依赖外部库）
const drawTrendChart = (data) => {
  if (!trendChart.value) return
  
  const canvas = trendChart.value
  const ctx = canvas.getContext('2d')
  const width = canvas.width = canvas.offsetWidth
  const height = canvas.height = canvas.offsetHeight
  
  ctx.clearRect(0, 0, width, height)
  
  if (!data || data.length === 0) {
    ctx.fillStyle = '#999'
    ctx.font = '14px sans-serif'
    ctx.textAlign = 'center'
    ctx.fillText('暂无数据', width / 2, height / 2)
    return
  }
  
  const padding = 50
  const chartWidth = width - padding * 2
  const chartHeight = height - padding * 2
  const maxValue = Math.max(...data.map(d => d.value), 1)
  
  // 绘制背景渐变
  const bgGradient = ctx.createLinearGradient(0, padding, 0, height - padding)
  bgGradient.addColorStop(0, 'rgba(255, 204, 0, 0.05)')
  bgGradient.addColorStop(1, 'rgba(255, 204, 0, 0)')
  ctx.fillStyle = bgGradient
  ctx.fillRect(padding, padding, chartWidth, chartHeight)
  
  // 绘制网格
  ctx.strokeStyle = 'rgba(0, 0, 0, 0.08)'
  ctx.lineWidth = 1
  ctx.setLineDash([2, 4])
  for (let i = 0; i <= 5; i++) {
    const y = padding + (chartHeight / 5) * i
    ctx.beginPath()
    ctx.moveTo(padding, y)
    ctx.lineTo(width - padding, y)
    ctx.stroke()
  }
  ctx.setLineDash([])
  
  // 绘制数据区域（渐变填充）
  const areaGradient = ctx.createLinearGradient(0, padding, 0, height - padding)
  areaGradient.addColorStop(0, 'rgba(255, 204, 0, 0.3)')
  areaGradient.addColorStop(1, 'rgba(255, 204, 0, 0.05)')
  
  ctx.fillStyle = areaGradient
  ctx.beginPath()
  data.forEach((point, index) => {
    const x = padding + (chartWidth / (data.length - 1 || 1)) * index
    const y = padding + chartHeight - (point.value / maxValue) * chartHeight
    if (index === 0) {
      ctx.moveTo(x, height - padding)
      ctx.lineTo(x, y)
    } else {
      ctx.lineTo(x, y)
    }
  })
  ctx.lineTo(width - padding, height - padding)
  ctx.closePath()
  ctx.fill()
  
  // 绘制折线（带渐变）
  const lineGradient = ctx.createLinearGradient(padding, 0, width - padding, 0)
  lineGradient.addColorStop(0, '#ffd700')
  lineGradient.addColorStop(0.5, '#ffcc00')
  lineGradient.addColorStop(1, '#ffaa00')
  
  ctx.strokeStyle = lineGradient
  ctx.lineWidth = 3
  ctx.lineCap = 'round'
  ctx.lineJoin = 'round'
  ctx.shadowColor = 'rgba(255, 204, 0, 0.5)'
  ctx.shadowBlur = 8
  ctx.shadowOffsetX = 0
  ctx.shadowOffsetY = 2
  
  ctx.beginPath()
  data.forEach((point, index) => {
    const x = padding + (chartWidth / (data.length - 1 || 1)) * index
    const y = padding + chartHeight - (point.value / maxValue) * chartHeight
    if (index === 0) {
      ctx.moveTo(x, y)
    } else {
      ctx.lineTo(x, y)
    }
  })
  ctx.stroke()
  
  // 重置阴影
  ctx.shadowColor = 'transparent'
  ctx.shadowBlur = 0
  ctx.shadowOffsetX = 0
  ctx.shadowOffsetY = 0
  
  // 绘制点（带光晕效果）
  data.forEach((point, index) => {
    const x = padding + (chartWidth / (data.length - 1 || 1)) * index
    const y = padding + chartHeight - (point.value / maxValue) * chartHeight
    
    // 外圈光晕
    const glowGradient = ctx.createRadialGradient(x, y, 0, x, y, 8)
    glowGradient.addColorStop(0, 'rgba(255, 204, 0, 0.6)')
    glowGradient.addColorStop(1, 'rgba(255, 204, 0, 0)')
    ctx.fillStyle = glowGradient
    ctx.beginPath()
    ctx.arc(x, y, 8, 0, Math.PI * 2)
    ctx.fill()
    
    // 内圈点
    const pointGradient = ctx.createRadialGradient(x, y, 0, x, y, 5)
    pointGradient.addColorStop(0, '#ffd700')
    pointGradient.addColorStop(1, '#ffcc00')
    ctx.fillStyle = pointGradient
    ctx.beginPath()
    ctx.arc(x, y, 5, 0, Math.PI * 2)
    ctx.fill()
    
    // 白色高光
    ctx.fillStyle = 'rgba(255, 255, 255, 0.6)'
    ctx.beginPath()
    ctx.arc(x - 1.5, y - 1.5, 2, 0, Math.PI * 2)
    ctx.fill()
  })
  
  // 绘制Y轴标签
  ctx.fillStyle = 'rgba(0, 0, 0, 0.5)'
  ctx.font = '11px sans-serif'
  ctx.textAlign = 'right'
  for (let i = 0; i <= 5; i++) {
    const y = padding + (chartHeight / 5) * i
    const value = Math.round(maxValue - (maxValue / 5) * i)
    ctx.fillText(value.toString(), padding - 10, y + 4)
  }
  
  // 绘制X轴标签
  ctx.fillStyle = 'rgba(0, 0, 0, 0.6)'
  ctx.font = '12px sans-serif'
  ctx.textAlign = 'center'
  data.forEach((point, index) => {
    const x = padding + (chartWidth / (data.length - 1 || 1)) * index
    ctx.fillText(point.label, x, height - 15)
  })
}

const drawSourceChart = (data) => {
  if (!sourceChart.value) return
  
  const canvas = sourceChart.value
  const ctx = canvas.getContext('2d')
  const width = canvas.width = canvas.offsetWidth
  const height = canvas.height = canvas.offsetHeight
  
  ctx.clearRect(0, 0, width, height)
  
  if (!data || data.length === 0) {
    ctx.fillStyle = '#999'
    ctx.font = '14px sans-serif'
    ctx.textAlign = 'center'
    ctx.fillText('暂无数据', width / 2, height / 2)
    return
  }
  
  const centerX = width / 2
  const centerY = height / 2 - 10
  const radius = Math.min(width, height) / 2 - 50
  const total = data.reduce((sum, item) => sum + item.value, 0)
  
  let currentAngle = -Math.PI / 2
  const colors = [
    { start: '#ffd700', end: '#ffcc00' },
    { start: '#42a5f5', end: '#2196f3' },
    { start: '#66bb6a', end: '#4caf50' },
    { start: '#ffb74d', end: '#ff9800' },
    { start: '#ba68c8', end: '#9c27b0' }
  ]
  
  data.forEach((item, index) => {
    const sliceAngle = (item.value / total) * Math.PI * 2
    const color = colors[index % colors.length]
    
    // 创建渐变
    const gradient = ctx.createRadialGradient(
      centerX + Math.cos(currentAngle + sliceAngle / 2) * (radius * 0.3),
      centerY + Math.sin(currentAngle + sliceAngle / 2) * (radius * 0.3),
      0,
      centerX,
      centerY,
      radius
    )
    gradient.addColorStop(0, color.start)
    gradient.addColorStop(1, color.end)
    
    // 绘制扇形（带阴影）
    ctx.shadowColor = 'rgba(0, 0, 0, 0.2)'
    ctx.shadowBlur = 8
    ctx.shadowOffsetX = 2
    ctx.shadowOffsetY = 2
    
    ctx.beginPath()
    ctx.moveTo(centerX, centerY)
    ctx.arc(centerX, centerY, radius, currentAngle, currentAngle + sliceAngle)
    ctx.closePath()
    ctx.fillStyle = gradient
    ctx.fill()
    
    // 绘制分隔线
    ctx.strokeStyle = 'rgba(255, 255, 255, 0.8)'
    ctx.lineWidth = 2
    ctx.beginPath()
    ctx.moveTo(centerX, centerY)
    ctx.lineTo(
      centerX + Math.cos(currentAngle) * radius,
      centerY + Math.sin(currentAngle) * radius
    )
    ctx.stroke()
    
    // 重置阴影
    ctx.shadowColor = 'transparent'
    ctx.shadowBlur = 0
    ctx.shadowOffsetX = 0
    ctx.shadowOffsetY = 0
    
    // 绘制标签（带背景）
    const labelAngle = currentAngle + sliceAngle / 2
    const labelX = centerX + Math.cos(labelAngle) * (radius * 0.75)
    const labelY = centerY + Math.sin(labelAngle) * (radius * 0.75)
    
    // 标签背景
    ctx.fillStyle = 'rgba(255, 255, 255, 0.9)'
    ctx.beginPath()
    ctx.arc(labelX, labelY, 20, 0, Math.PI * 2)
    ctx.fill()
    
    // 标签文字
    ctx.fillStyle = '#333'
    ctx.font = 'bold 11px sans-serif'
    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'
    ctx.fillText(item.label, labelX, labelY)
    
    // 百分比
    const percent = ((item.value / total) * 100).toFixed(1)
    ctx.fillStyle = 'rgba(0, 0, 0, 0.6)'
    ctx.font = '10px sans-serif'
    ctx.fillText(percent + '%', labelX, labelY + 15)
    
    currentAngle += sliceAngle
  })
  
  // 绘制中心圆（白色背景）
  ctx.fillStyle = 'rgba(255, 255, 255, 0.95)'
  ctx.beginPath()
  ctx.arc(centerX, centerY, radius * 0.4, 0, Math.PI * 2)
  ctx.fill()
  
  // 绘制图例（改进样式）
  const legendY = height - 25
  const legendStartX = (width - (data.length * 120)) / 2
  
  data.forEach((item, index) => {
    const color = colors[index % colors.length]
    const x = legendStartX + index * 120
    
    // 图例方块（带渐变）
    const legendGradient = ctx.createLinearGradient(x, legendY, x + 16, legendY + 16)
    legendGradient.addColorStop(0, color.start)
    legendGradient.addColorStop(1, color.end)
    
    ctx.fillStyle = legendGradient
    ctx.fillRect(x, legendY, 16, 16)
    
    // 图例文字
    ctx.fillStyle = 'rgba(0, 0, 0, 0.7)'
    ctx.font = '12px sans-serif'
    ctx.textAlign = 'left'
    ctx.textBaseline = 'middle'
    ctx.fillText(item.label, x + 22, legendY + 8)
    
    // 数值
    const percent = ((item.value / total) * 100).toFixed(1)
    ctx.fillStyle = 'rgba(0, 0, 0, 0.5)'
    ctx.font = '11px sans-serif'
    ctx.fillText(`(${item.value}, ${percent}%)`, x + 22, legendY + 22)
  })
}

// 数字动画函数
const animateNumber = (target, current, duration = 1000) => {
  const start = current
  const end = target
  const startTime = performance.now()
  
  const animate = (currentTime) => {
    const elapsed = currentTime - startTime
    const progress = Math.min(elapsed / duration, 1)
    
    // 使用缓动函数
    const easeOutQuart = 1 - Math.pow(1 - progress, 4)
    const value = Math.floor(start + (end - start) * easeOutQuart)
    
    return value
  }
  
  return animate
}

// 更新动画数字
const updateAnimatedStats = (newStats) => {
  const duration = 800
  const steps = 60
  const stepDuration = duration / steps
  
  Object.keys(newStats).forEach((key) => {
    const target = newStats[key]
    const start = animatedStats.value[key]
    let currentStep = 0
    
    const interval = setInterval(() => {
      currentStep++
      const progress = currentStep / steps
      const easeOutQuart = 1 - Math.pow(1 - progress, 4)
      animatedStats.value[key] = Math.floor(start + (target - start) * easeOutQuart)
      
      if (currentStep >= steps) {
        animatedStats.value[key] = target
        clearInterval(interval)
      }
    }, stepDuration)
  })
}

const loadStats = async () => {
  try {
    const res = await adminAPI.getStats()
    stats.value = res.data
    updateAnimatedStats(res.data)
  } catch (error) {
    console.error('加载统计失败:', error)
    // 使用模拟数据
    const mockStats = {
      totalViews: 1234,
      uniqueVisitors: 567,
      todayViews: 89,
      totalSites: 12,
    }
    stats.value = mockStats
    updateAnimatedStats(mockStats)
  }
}

const loadChartData = async () => {
  try {
    const res = await adminAPI.getChartData(chartPeriod.value)
    await nextTick()
    drawTrendChart(res.data.trend)
    drawSourceChart(res.data.sources)
  } catch (error) {
    console.error('加载图表数据失败:', error)
    // 使用模拟数据
    const mockTrend = Array.from({ length: chartPeriod.value }, (_, i) => ({
      label: `${i + 1}日`,
      value: Math.floor(Math.random() * 100) + 50,
    }))
    const mockSources = [
      { label: '直接访问', value: 45 },
      { label: '搜索引擎', value: 30 },
      { label: '社交媒体', value: 15 },
      { label: '其他', value: 10 },
    ]
    await nextTick()
    drawTrendChart(mockTrend)
    drawSourceChart(mockSources)
  }
}


const loadLogs = async () => {
  loadingLogs.value = true
  try {
    const res = await adminAPI.getBackendLogs(logLines.value)
    backendLogs.value = res.data.logs || []
    
    // 自动滚动到底部
    if (autoScroll.value) {
      nextTick(() => {
        scrollToBottom()
      })
    }
  } catch (error) {
    console.error('加载日志失败:', error)
    backendLogs.value = []
  } finally {
    loadingLogs.value = false
  }
}

const scrollToBottom = () => {
  if (logsContainer.value) {
    const container = logsContainer.value
    container.scrollTop = container.scrollHeight
  }
}

const getLogLevel = (log) => {
  if (!log) return ''
  const logUpper = log.toUpperCase()
  if (logUpper.includes('[ERROR]') || logUpper.includes('ERROR')) {
    return 'log-error'
  }
  if (logUpper.includes('[WARN]') || logUpper.includes('WARN')) {
    return 'log-warn'
  }
  if (logUpper.includes('[INFO]') || logUpper.includes('INFO')) {
    return 'log-info'
  }
  return 'log-default'
}

let refreshInterval = null

const loadLoginHistory = async () => {
  loadingLoginHistory.value = true
  try {
    // 最近登录 IP 只展示 5 条
    const res = await adminAPI.getLoginHistory(5)
    loginHistory.value = (res.data.data || res.data || []).slice(0, 5)
  } catch (error) {
    console.error('加载登录历史失败:', error)
    loginHistory.value = []
  } finally {
    loadingLoginHistory.value = false
  }
}

onMounted(async () => {
  isLoading.value = true
  try {
    await Promise.all([
      loadStats(),
      loadChartData(),
      loadLogs(),
      loadLoginHistory()
    ])
  } finally {
    // 延迟一下让动画更自然
    setTimeout(() => {
      isLoading.value = false
    }, 300)
  }
  
  // 每30秒刷新一次数据
  refreshInterval = setInterval(() => {
    loadStats()
    loadChartData()
    loadLogs()
    loadLoginHistory()
  }, 30000)
  
  // 窗口大小改变时重绘图表
  window.addEventListener('resize', () => {
    loadChartData()
  })
  
  // 监听日志容器滚动，如果用户手动滚动，则禁用自动滚动
  if (logsContainer.value) {
    logsContainer.value.addEventListener('scroll', () => {
      const container = logsContainer.value
      const isAtBottom = container.scrollHeight - container.scrollTop <= container.clientHeight + 10
      if (!isAtBottom) {
        autoScroll.value = false
      }
    })
  }
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
  window.removeEventListener('resize', loadChartData)
})
</script>

<style scoped>
.dashboard {
  padding: 4px 0 0;
  display: flex;
  flex-direction: column;
  gap: 28px;
  animation: fadeInUp 0.6s ease-out;
  position: relative;
}

.dashboard.loading {
  pointer-events: none;
  opacity: 0.7;
}

.dashboard.loading * {
  animation-play-state: paused !important;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 24px;
}

.stat-card {
  background: linear-gradient(135deg, rgba(var(--background-color-rgb), 0.8) 0%, rgba(var(--background-color-rgb), 0.6) 100%);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  padding: 28px;
  display: flex;
  align-items: center;
  gap: 20px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08), 0 2px 4px rgba(0, 0, 0, 0.04);
  position: relative;
  overflow: hidden;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, transparent, var(--hover-link-color), transparent);
  opacity: 0;
  transition: opacity 0.4s ease;
}

.stat-card:hover {
  transform: translateY(-4px) scale(1.02);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.12), 0 4px 8px rgba(0, 0, 0, 0.08);
  border-color: var(--hover-link-color);
}

.stat-card:hover::before {
  opacity: 1;
}

.stat-icon {
  width: 64px;
  height: 64px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.75rem;
  flex-shrink: 0;
  position: relative;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-card:hover .stat-icon {
  transform: scale(1.1) rotate(5deg);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.15);
}

.stat-content {
  flex: 1;
  position: relative;
  z-index: 1;
}

.stat-value {
  font-size: 2.25rem;
  font-weight: 700;
  color: var(--text-color);
  line-height: 1.2;
  margin-bottom: 6px;
  background: linear-gradient(135deg, var(--text-color) 0%, rgba(var(--text-color-rgb, 51, 51, 51), 0.8) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  transition: all 0.3s ease;
}

.stat-card:hover .stat-value {
  transform: scale(1.05);
}

.stat-label {
  font-size: 0.95rem;
  color: rgba(var(--text-color-rgb, 51, 51, 51), 0.65);
  font-weight: 500;
  letter-spacing: 0.3px;
}

.charts-grid {
  display: grid;
  grid-template-columns: minmax(0, 2.1fr) minmax(0, 1.7fr);
  gap: 28px;
}

.chart-card--primary .chart-body canvas {
  height: 360px;
  max-height: 360px;
}

.chart-card--secondary .chart-body canvas {
  height: 300px;
  max-height: 300px;
}

.chart-card,
.recent-card {
  background: linear-gradient(135deg, rgba(var(--background-color-rgb), 0.85) 0%, rgba(var(--background-color-rgb), 0.7) 100%);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08), 0 2px 4px rgba(0, 0, 0, 0.04);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

.chart-card:hover,
.recent-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12), 0 4px 8px rgba(0, 0, 0, 0.08);
  border-color: var(--hover-link-color);
}

.card-header {
  padding: 22px 28px;
  background: linear-gradient(135deg, rgba(var(--background-color-rgb), 0.95) 0%, rgba(var(--background-color-rgb), 0.85) 100%);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: relative;
}

.card-header::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--hover-link-color), transparent);
  opacity: 0.3;
}

.card-header h3 {
  margin: 0;
  font-size: 1.15rem;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--text-color);
  letter-spacing: 0.3px;
}

.card-header i {
  color: var(--hover-link-color);
  font-size: 1.2rem;
  transition: transform 0.3s ease;
}

.chart-card:hover .card-header i {
  transform: scale(1.1);
}

.period-select {
  padding: 8px 16px;
  border: 1.5px solid var(--border-color);
  border-radius: 10px;
  background: rgba(var(--background-color-rgb), 0.9);
  color: var(--text-color);
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 500;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.period-select:hover {
  border-color: var(--hover-link-color);
  background: rgba(var(--background-color-rgb), 1);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.period-select:focus {
  outline: none;
  border-color: var(--hover-link-color);
  box-shadow: 0 0 0 3px rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.2);
}

.card-body {
  padding: 28px;
  position: relative;
}

.card-body canvas {
  width: 100%;
  height: 300px;
  max-height: 320px;
  border-radius: 8px;
}

.recent-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.recent-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: linear-gradient(135deg, rgba(var(--background-color-rgb), 0.5) 0%, rgba(var(--background-color-rgb), 0.3) 100%);
  border-radius: 12px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid transparent;
  position: relative;
  overflow: hidden;
}

.recent-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: var(--hover-link-color);
  transform: scaleY(0);
  transition: transform 0.3s ease;
}

.recent-item:hover {
  background: linear-gradient(135deg, rgba(var(--background-color-rgb), 0.7) 0%, rgba(var(--background-color-rgb), 0.5) 100%);
  transform: translateX(4px);
  border-color: var(--hover-link-color);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.recent-item:hover::before {
  transform: scaleY(1);
}

.recent-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: linear-gradient(135deg, rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.15) 0%, rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.1) 100%);
  color: var(--hover-link-color);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 1.1rem;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.recent-item:hover .recent-icon {
  transform: scale(1.1) rotate(5deg);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
}

.recent-content {
  flex: 1;
  position: relative;
  z-index: 1;
}

.recent-title {
  font-weight: 600;
  color: var(--text-color);
  margin-bottom: 6px;
  font-size: 0.95rem;
  letter-spacing: 0.2px;
}

.recent-meta {
  display: flex;
  gap: 16px;
  font-size: 0.85rem;
  color: rgba(var(--text-color-rgb, 51, 51, 51), 0.65);
  font-weight: 500;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: rgba(var(--text-color-rgb, 51, 51, 51), 0.5);
  animation: fadeIn 0.5s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.empty-state i {
  font-size: 3.5rem;
  margin-bottom: 20px;
  opacity: 0.25;
  display: inline-block;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.empty-state p {
  margin: 0;
  font-size: 1rem;
  font-weight: 500;
  letter-spacing: 0.3px;
}

/* 底部布局：现代化两栏布局 */
.bottom-grid-modern {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1.8fr);
  gap: 20px;
  align-items: start;
}

.recent-login-card {
  display: flex;
  flex-direction: column;
  min-height: 0;
  height: 100%;
  background: linear-gradient(135deg, rgba(var(--background-color-rgb), 0.85) 0%, rgba(var(--background-color-rgb), 0.7) 100%);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08), 0 2px 4px rgba(0, 0, 0, 0.04);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

.recent-login-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12), 0 4px 8px rgba(0, 0, 0, 0.08);
  border-color: var(--hover-link-color);
}

.compact-card {
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.compact-body {
  flex: 1;
  padding: 18px;
  overflow-y: auto;
  min-height: 0;
  max-height: 400px;
  /* 自定义滚动条 */
  scrollbar-width: thin;
  scrollbar-color: var(--hover-link-color) rgba(var(--background-color-rgb), 0.3);
}

.compact-body::-webkit-scrollbar {
  width: 6px;
}

.compact-body::-webkit-scrollbar-track {
  background: rgba(var(--background-color-rgb), 0.3);
  border-radius: 3px;
}

.compact-body::-webkit-scrollbar-thumb {
  background: var(--hover-link-color);
  border-radius: 3px;
  transition: background 0.3s ease;
}

.compact-body::-webkit-scrollbar-thumb:hover {
  background: #ffd700;
}

.logs-card-full {
  display: flex;
  flex-direction: column;
  min-height: 0;
  height: 100%;
  background: linear-gradient(135deg, rgba(var(--background-color-rgb), 0.85) 0%, rgba(var(--background-color-rgb), 0.7) 100%);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08), 0 2px 4px rgba(0, 0, 0, 0.04);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  grid-column: span 1;
}

.logs-card-full:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12), 0 4px 8px rgba(0, 0, 0, 0.08);
  border-color: var(--hover-link-color);
}

.logs-body-full {
  flex: 1;
  padding: 0;
  overflow-y: auto;
  min-height: 0;
  max-height: 400px;
  background: #1e1e1e;
  font-family: 'Courier New', 'Consolas', monospace;
  /* 自定义滚动条 */
  scrollbar-width: thin;
  scrollbar-color: var(--hover-link-color) rgba(0, 0, 0, 0.3);
}

.logs-body-full::-webkit-scrollbar {
  width: 8px;
}

.logs-body-full::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.3);
  border-radius: 4px;
}

.logs-body-full::-webkit-scrollbar-thumb {
  background: var(--hover-link-color);
  border-radius: 4px;
  transition: background 0.3s ease;
}

.logs-body-full::-webkit-scrollbar-thumb:hover {
  background: #ffd700;
}

.login-history-list {
  max-height: none;
}

/* 后台日志样式 */
.logs-card {
  background: linear-gradient(135deg, rgba(var(--background-color-rgb), 0.85) 0%, rgba(var(--background-color-rgb), 0.7) 100%);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08), 0 2px 4px rgba(0, 0, 0, 0.04);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

.logs-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12), 0 4px 8px rgba(0, 0, 0, 0.08);
  border-color: var(--hover-link-color);
}

.logs-controls {
  display: flex;
  align-items: center;
  gap: 12px;
}

.log-lines-select {
  padding: 8px 16px;
  border: 1.5px solid var(--border-color);
  border-radius: 10px;
  background: rgba(var(--background-color-rgb), 0.9);
  color: var(--text-color);
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 500;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.log-lines-select:hover {
  border-color: var(--hover-link-color);
  background: rgba(var(--background-color-rgb), 1);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.log-lines-select:focus {
  outline: none;
  border-color: var(--hover-link-color);
  box-shadow: 0 0 0 3px rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.2);
}

.refresh-logs-btn,
.auto-scroll-btn {
  width: 36px;
  height: 36px;
  border: 1.5px solid var(--border-color);
  border-radius: 10px;
  background: rgba(var(--background-color-rgb), 0.9);
  color: var(--text-color);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.refresh-logs-btn:hover,
.auto-scroll-btn:hover {
  background: rgba(var(--background-color-rgb), 1);
  border-color: var(--hover-link-color);
  transform: scale(1.05);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.refresh-logs-btn:active,
.auto-scroll-btn:active {
  transform: scale(0.95);
}

.auto-scroll-btn.active {
  background: linear-gradient(135deg, var(--hover-link-color) 0%, #ffd700 100%);
  color: #333;
  border-color: var(--hover-link-color);
  box-shadow: 0 4px 12px rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.3);
}

.logs-body {
  padding: 0;
  max-height: 400px;
  overflow-y: auto;
  background: #1e1e1e;
  font-family: 'Courier New', 'Consolas', monospace;
}

.logs-content {
  padding: 16px 20px;
}

.log-line {
  display: flex;
  gap: 12px;
  padding: 6px 0;
  font-size: 13px;
  line-height: 1.7;
  word-break: break-all;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.log-time {
  color: #888;
  flex-shrink: 0;
  min-width: 150px;
}

.log-content {
  color: #d4d4d4;
  flex: 1;
  font-family: 'Courier New', 'Consolas', 'Monaco', monospace;
  white-space: pre-wrap;
}

.log-line.log-error .log-content {
  color: #f48771;
}

.log-line.log-warn .log-content {
  color: #dcdcaa;
}

.log-line.log-info .log-content {
  color: #4ec9b0;
}

/* 登录历史样式 - 已移到compact-body中处理滚动 */

.login-item {
  border-left: 4px solid transparent;
  transition: all 0.3s ease;
}

.login-item .login-ip {
  font-weight: 600;
  color: var(--text-color);
  margin-right: 10px;
  font-size: 0.95rem;
}

.login-item .login-location {
  font-size: 0.85rem;
  color: rgba(var(--text-color-rgb, 51, 51, 51), 0.65);
  background: linear-gradient(135deg, rgba(var(--background-color-rgb), 0.6) 0%, rgba(var(--background-color-rgb), 0.4) 100%);
  padding: 4px 10px;
  border-radius: 6px;
  font-weight: 500;
  border: 1px solid var(--border-color);
}

.login-success {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.15) 0%, rgba(76, 175, 80, 0.08) 100%) !important;
  color: #4caf50 !important;
  box-shadow: 0 2px 8px rgba(76, 175, 80, 0.1);
}

.login-failed {
  background: linear-gradient(135deg, rgba(244, 67, 54, 0.15) 0%, rgba(244, 67, 54, 0.08) 100%) !important;
  color: #f44336 !important;
  box-shadow: 0 2px 8px rgba(244, 67, 54, 0.1);
}

.login-item.login-success {
  border-left-color: #4caf50;
}

.login-item.login-success:hover {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.2) 0%, rgba(76, 175, 80, 0.12) 100%) !important;
}

.login-item.login-failed {
  border-left-color: #f44336;
}

.login-item.login-failed:hover {
  background: linear-gradient(135deg, rgba(244, 67, 54, 0.2) 0%, rgba(244, 67, 54, 0.12) 100%) !important;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .charts-grid {
    grid-template-columns: 1fr;
  }
  
  .chart-card--primary .chart-body canvas {
    height: 320px;
  }
  
  .chart-card--secondary .chart-body canvas {
    height: 280px;
  }
}

@media (max-width: 1400px) {
  .bottom-grid-modern {
    grid-template-columns: 1fr;
    grid-template-rows: auto auto;
  }
  
  .logs-card-full {
    grid-column: span 1;
  }
}

@media (max-width: 1200px) {
  .bottom-grid-modern {
    grid-template-columns: 1fr;
    grid-template-rows: auto auto;
  }
  
  .logs-card-full {
    grid-column: span 1;
  }
  
  .compact-body {
    max-height: 280px;
  }
  
  .logs-body-full {
    max-height: 400px;
  }
}

@media (max-width: 768px) {
  .dashboard {
    gap: 20px;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .stat-card {
    padding: 20px;
  }
  
  .stat-icon {
    width: 56px;
    height: 56px;
    font-size: 1.5rem;
  }
  
  .stat-value {
    font-size: 1.9rem;
  }
  
  .charts-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .chart-card--primary .chart-body canvas,
  .chart-card--secondary .chart-body canvas {
    height: 280px;
  }

  .bottom-grid-modern {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .logs-card-full {
    grid-column: span 1;
  }
  
  .card-header {
    padding: 18px 20px;
  }
  
  .card-body {
    padding: 20px;
  }
  
  .compact-body {
    padding: 16px;
    max-height: 250px;
  }
  
  .logs-body-full {
    max-height: 350px;
  }
  
  .recent-item {
    padding: 14px;
  }
}

/* 骨架屏加载效果 */
.dashboard.loading .stat-card,
.dashboard.loading .chart-card,
.dashboard.loading .recent-card,
.dashboard.loading .logs-card {
  position: relative;
  overflow: hidden;
}

.dashboard.loading .stat-card::after,
.dashboard.loading .chart-card::after,
.dashboard.loading .recent-card::after,
.dashboard.loading .logs-card::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.3),
    transparent
  );
  animation: shimmer 1.5s infinite;
}

@keyframes shimmer {
  0% {
    left: -100%;
  }
  100% {
    left: 100%;
  }
}
</style>
