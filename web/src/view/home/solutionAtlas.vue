<template>
  <div class="solution-atlas">
    <div class="atlas-header">
      <div class="atlas-badge">
        <el-icon><Grid /></el-icon>
        <span>全局视角</span>
      </div>
      <h1 class="atlas-title">
        解决方案全景图
        <span class="text-primary">选择解决方案</span>
      </h1>
      <p class="atlas-subtitle">点击分类卡片，查看该分类下所有产品</p>
    </div>

    <div v-loading="loading" class="atlas-content">
      <div v-if="!topology" class="loading-state">
        <el-icon class="is-loading"><Loading /></el-icon>
        <span>加载架构拓扑...</span>
      </div>

      <div v-else class="atlas-grid">
        <!-- 左侧：安全体系 -->
        <div class="sidebar security-sidebar">
          <div class="sidebar-title">
            <el-icon><Lock /></el-icon>
            <span>安全体系</span>
          </div>
          <div class="sidebar-categories">
            <div
              v-for="cat in topology.security"
              :key="cat.id"
              class="category-card"
              :class="{ active: activeCategory?.id === cat.id }"
              @click="onCategoryClick(cat)"
            >
              <div class="category-icon">
                <el-icon><component :is="getIcon(cat.icon_key)" /></el-icon>
              </div>
              <div class="category-info">
                <span class="category-label">{{ cat.label }}</span>
                <span v-if="getKeywordText(cat.keywords)" class="category-keywords">{{ getKeywordText(cat.keywords) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 中心：层级结构 -->
        <div class="central-layers">
          <div
            v-for="(block, index) in topology.central_layers"
            :key="block.layer.id"
            class="layer-block"
            :style="{ animationDelay: `${index * 0.05}s` }"
          >
            <div class="layer-header">
              <div class="layer-level">L{{ block.layer.level }}</div>
              <div class="layer-titles">
                <h3 class="layer-title">{{ block.layer.title }}</h3>
                <p v-if="block.layer.subtitle" class="layer-subtitle">{{ block.layer.subtitle }}</p>
              </div>
            </div>
            <div class="category-grid">
              <div
                v-for="cat in block.categories"
                :key="cat.id"
                class="category-card"
                :class="{ active: activeCategory?.id === cat.id }"
                @click="onCategoryClick(cat)"
              >
                <div class="category-icon">
                  <el-icon><component :is="getIcon(cat.icon_key)" /></el-icon>
                </div>
                <div class="category-info">
                  <span class="category-label">{{ cat.label }}</span>
                  <span v-if="getKeywordText(cat.keywords)" class="category-keywords">{{ getKeywordText(cat.keywords) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧：运维体系 -->
        <div class="sidebar ops-sidebar">
          <div class="sidebar-title">
            <el-icon><Setting /></el-icon>
            <span>运维体系</span>
          </div>
          <div class="sidebar-categories">
            <div
              v-for="cat in topology.ops"
              :key="cat.id"
              class="category-card"
              :class="{ active: activeCategory?.id === cat.id }"
              @click="onCategoryClick(cat)"
            >
              <div class="category-icon">
                <el-icon><component :is="getIcon(cat.icon_key)" /></el-icon>
              </div>
              <div class="category-info">
                <span class="category-label">{{ cat.label }}</span>
                <span v-if="getKeywordText(cat.keywords)" class="category-keywords">{{ getKeywordText(cat.keywords) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { Grid, Loading, Lock, Setting } from '@element-plus/icons-vue'
import { getFullTopology } from '@/api/business/topology'
import { ElMessage } from 'element-plus'

const props = defineProps({
  initialSolutionId: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['selectSolution'])

const loading = ref(true)
const topology = ref(null)
const activeCategory = ref(null)

// Element Plus icon mapping for icon_key values
const iconMap = {
  // Security
  'ServerCog': 'Monitor',
  'Shield': 'Shield',
  'Lock': 'Lock',
  'Fingerprint': 'Key',
  'ShieldCheck': 'CircleCheck',
  'Radar': 'Odometer',
  // Ops
  'Brain': 'Cpu',
  'Database': 'Database',
  'Activity': 'DataAnalysis',
  'LineChart': 'TrendCharts',
  'FileStack': 'DocumentCopy',
  'Workflow': 'Operation',
  'CloudCog': 'Cloud',
  // Central categories
  'Server': 'Monitor',
  'Box': 'Box',
  'LayoutGrid': 'Grid',
  'Bot': 'Robot',
  'Cloud': 'Cloud',
  'Network': 'Connection',
  'Router': 'Guide',
  'Scale': 'ScaleToOriginal',
  'Globe2': 'Globe',
  'HardDrive': 'HardDrive',
  'Layers3': 'CopyDocument',
  'Sparkles': 'MagicStick',
  'Building2': 'OfficeBuilding',
  'Zap': 'Lightning',
  'ThermometerSnowflake': 'Cold',
  'Cable': 'Connection',
  'default': 'Folder'
}

const getIcon = (iconKey) => {
  return iconMap[iconKey] || iconMap.default
}

const getKeywordText = (keywords) => {
  if (!keywords) return ''
  if (Array.isArray(keywords)) return keywords.slice(0, 2).join(' / ')
  try {
    const arr = JSON.parse(keywords)
    if (Array.isArray(arr)) return arr.slice(0, 2).join(' / ')
  } catch {}
  return ''
}

const loadTopology = async () => {
  try {
    loading.value = true
    const res = await getFullTopology()
    if (res.code === 0) {
      topology.value = res.data
    } else {
      ElMessage.error('加载拓扑失败')
    }
  } catch (err) {
    ElMessage.error('加载拓扑失败，请确认后端已启动')
    console.error(err)
  } finally {
    loading.value = false
  }
}

const findCategoryById = (id) => {
  if (!topology.value) return null

  // Check security
  for (const cat of (topology.value.security || [])) {
    if (String(cat.id) === String(id)) return cat
  }
  // Check ops
  for (const cat of (topology.value.ops || [])) {
    if (String(cat.id) === String(id)) return cat
  }
  // Check central layers
  for (const block of (topology.value.central_layers || [])) {
    for (const cat of block.categories) {
      if (String(cat.id) === String(id)) return cat
    }
  }
  return null
}

const onCategoryClick = (cat) => {
  if (activeCategory.value?.id === cat.id) {
    activeCategory.value = null
  } else {
    activeCategory.value = cat
  }
  emit('selectSolution', activeCategory.value)
}

onMounted(() => {
  loadTopology()
})

watch(() => props.initialSolutionId, (val) => {
  if (val && topology.value) {
    activeCategory.value = findCategoryById(val)
  }
})
</script>

<style scoped>
.solution-atlas {
  padding: 32px 24px;
  max-width: 1600px;
  margin: 0 auto;
}

.atlas-header {
  text-align: center;
  margin-bottom: 32px;
}

.atlas-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: var(--el-color-primary-light-9);
  color: var(--el-color-primary);
  border-radius: 16px;
  font-size: 13px;
  margin-bottom: 12px;
}

.atlas-title {
  font-size: 28px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 8px 0;
}

.atlas-subtitle {
  font-size: 14px;
  color: var(--el-text-color-secondary);
  margin: 0;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 64px;
  color: var(--el-text-color-secondary);
}

.atlas-grid {
  display: grid;
  grid-template-columns: 240px 1fr 240px;
  gap: 24px;
  align-items: start;
}

/* 侧边栏样式 */
.sidebar {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.sidebar-title {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: linear-gradient(135deg, var(--el-color-primary) 0%, var(--el-color-primary-light-3) 100%);
  color: #fff;
  border-radius: 12px;
  font-weight: 600;
  font-size: 15px;
}

.security-sidebar .sidebar-title {
  background: linear-gradient(135deg, #e74c3c 0%, #c0392b 100%);
}

.ops-sidebar .sidebar-title {
  background: linear-gradient(135deg, #27ae60 0%, #1e8449 100%);
}

.sidebar-categories {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* 中心层级 */
.central-layers {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.layer-block {
  background: linear-gradient(90deg, var(--el-color-primary-light-9) 0%, #fff 100%);
  border: 1px solid var(--el-color-primary-light-5);
  border-radius: 16px;
  padding: 16px 20px;
  animation: fadeInUp 0.4s ease-out forwards;
  opacity: 0;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.layer-header {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 14px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.layer-level {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-color-primary);
  color: #fff;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 700;
  flex-shrink: 0;
}

.layer-titles {
  flex: 1;
}

.layer-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 2px 0;
}

.layer-subtitle {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin: 0;
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 10px;
}

.category-card {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: #fff;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s;
}

.category-card:hover {
  border-color: var(--el-color-primary);
  box-shadow: 0 2px 8px var(--el-color-primary-light-8);
  transform: translateY(-1px);
}

.category-card.active {
  border-color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
}

.category-icon {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-color-primary-light-9);
  border-radius: 6px;
  color: var(--el-color-primary);
  font-size: 16px;
  flex-shrink: 0;
}

.category-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
  flex: 1;
}

.category-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--el-text-color-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.category-keywords {
  font-size: 10px;
  color: var(--el-text-color-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.text-primary {
  color: var(--el-color-primary);
}
</style>
