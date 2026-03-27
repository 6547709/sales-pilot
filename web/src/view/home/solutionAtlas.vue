<template>
  <div class="solution-atlas">
    <div class="atlas-header">
      <div class="atlas-badge">
        <el-icon><Grid /></el-icon>
        <span>全局视角</span>
      </div>
      <h1 class="atlas-title">
        按架构层级与体系
        <span class="text-primary">选择解决方案</span>
      </h1>
    </div>

    <div v-loading="loading" class="atlas-content">
      <div v-if="!topology" class="loading-state">
        <el-icon class="is-loading"><Loading /></el-icon>
        <span>加载架构拓扑...</span>
      </div>

      <div v-else class="atlas-grid">
        <div class="central-layers">
          <div
            v-for="(block, index) in sortedCentralLayers"
            :key="block.layer.id"
            class="layer-block"
            :style="{ animationDelay: `${index * 0.05}s` }"
          >
            <div class="layer-header">
              <h3 class="layer-title">{{ block.layer.title || block.layer.name }}</h3>
              <p v-if="block.layer.subtitle" class="layer-subtitle">{{ block.layer.subtitle }}</p>
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
                  <span class="category-label">{{ cat.label || cat.name }}</span>
                  <span v-if="cat.hint" class="category-hint">{{ cat.hint }}</span>
                </div>
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
import { Grid, Loading } from '@element-plus/icons-vue'
import { getFullTopology } from '@/api/business/topology'
import { ElMessage } from 'element-plus'

const props = defineProps({
  initialSolutionId: {
    type: String,
    default: ''
  },
  initialKeyword: {
    type: String,
    default: ''
  },
  initialMarket: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['selectSolution'])

const loading = ref(true)
const topology = ref(null)
const activeCategory = ref(null)

// Simple icon mapping
const iconMap = {
  'security': 'Lock',
  'network': 'Connection',
  'cloud': 'Cloud',
  'storage': 'Folder',
  'database': 'Database',
  'server': 'Monitor',
  'default': 'Folder'
}

const getIcon = (iconKey) => {
  return iconMap[iconKey] || iconMap.default
}

const sortedCentralLayers = computed(() => {
  if (!topology.value?.central_layers) return []
  return [...topology.value.central_layers].sort((a, b) => b.layer.level - a.layer.level)
})

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
  if (!topology.value?.central_layers) return null
  for (const block of topology.value.central_layers) {
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
  if (props.initialSolutionId) {
    // Will be set after topology loads
  }
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
  max-width: 1400px;
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
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.central-layers {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.layer-block {
  background: linear-gradient(90deg, var(--el-color-primary-light-9) 0%, #fff 100%);
  border: 1px solid var(--el-color-primary-light-5);
  border-radius: 16px;
  padding: 20px;
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
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.layer-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0;
}

.layer-subtitle {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin: 0;
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 12px;
}

.category-card {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  background: #fff;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s;
}

.category-card:hover {
  border-color: var(--el-color-primary);
  box-shadow: 0 2px 8px var(--el-color-primary-light-8);
}

.category-card.active {
  border-color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
}

.category-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-color-primary-light-9);
  border-radius: 8px;
  color: var(--el-color-primary);
  font-size: 18px;
}

.category-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.category-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--el-text-color-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.category-hint {
  font-size: 11px;
  color: var(--el-text-color-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.text-primary {
  color: var(--el-color-primary);
}
</style>
