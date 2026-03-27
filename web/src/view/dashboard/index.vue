<template>
  <div class="h-full gva-container2 overflow-auto bg-slate-50/60 dark:bg-slate-900">
    <div class="space-y-4 p-4 lg:p-6">
      <section
        class="relative overflow-hidden rounded-xl border border-slate-200/80 bg-white px-5 py-6 shadow-sm dark:border-slate-700 dark:from-slate-900 dark:via-slate-800 dark:to-slate-900"
      >
        <div class="relative flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
          <div>
            <p class="text-xs tracking-[0.2em] text-slate-500 dark:text-slate-400">DASHBOARD</p>
            <h1 class="mt-2 text-xl font-semibold text-slate-900 dark:text-slate-100 lg:text-2xl">
              销售赋能平台
            </h1>
            <p class="mt-2 text-sm text-slate-600 dark:text-slate-300">
              {{ today }} · 业务数据总览
            </p>
          </div>
        </div>
      </section>

      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #409eff20; color: #409eff;">
              <el-icon><Goods /></el-icon>
            </div>
            <div class="stat-info">
              <span class="stat-value">{{ stats.product_count || 0 }}</span>
              <span class="stat-label">产品数量</span>
            </div>
          </div>
        </el-card>
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #67c23a20; color: #67c23a;">
              <el-icon><ChatDotRound /></el-icon>
            </div>
            <div class="stat-info">
              <span class="stat-value">{{ stats.script_count || 0 }}</span>
              <span class="stat-label">销售话术</span>
            </div>
          </div>
        </el-card>
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #e6a23c20; color: #e6a23c;">
              <el-icon><Briefcase /></el-icon>
            </div>
            <div class="stat-info">
              <span class="stat-value">{{ stats.case_count || 0 }}</span>
              <span class="stat-label">客户案例</span>
            </div>
          </div>
        </el-card>
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #f56c6c20; color: #f56c6c;">
              <el-icon><Folder /></el-icon>
            </div>
            <div class="stat-info">
              <span class="stat-value">{{ stats.category_count || 0 }}</span>
              <span class="stat-label">解决方案分类</span>
            </div>
          </div>
        </el-card>
      </div>

      <div class="grid grid-cols-1 items-stretch gap-4 xl:grid-cols-12">
        <div class="grid grid-cols-1 gap-4 content-start xl:col-span-8 xl:h-full">
          <el-card header="快捷入口">
            <div class="quick-links">
              <el-button type="primary" @click="goPage('/admin/product')">产品管理</el-button>
              <el-button type="success" @click="goPage('/admin/script')">销售话术</el-button>
              <el-button type="warning" @click="goPage('/admin/case')">客户案例</el-button>
              <el-button type="danger" @click="goPage('/admin/topology')">拓扑管理</el-button>
              <el-button @click="goPage('/admin/solutionCategory')">解决方案分类</el-button>
              <el-button @click="goPage('/admin/apiKey')">API密钥</el-button>
            </div>
          </el-card>

          <el-card header="访问首页">
            <div class="home-link-section">
              <p class="text-sm text-gray-600 mb-3">查看公开首页（全景图 - 产品赋能）</p>
              <el-button type="primary" plain @click="goHome">打开首页</el-button>
            </div>
          </el-card>
        </div>

        <div class="flex flex-col gap-4 xl:col-span-4 xl:h-full">
          <el-card header="拓扑层级">
            <div class="layer-info">
              <p>当前系统共有 <strong>{{ stats.layer_count || 0 }}</strong> 个拓扑层级</p>
            </div>
          </el-card>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { Goods, ChatDotRound, Briefcase, Folder } from '@element-plus/icons-vue'
  import { getStatisticsOverview } from '@/api/business/statistics'
  import { ElMessage } from 'element-plus'

  const router = useRouter()
  const stats = ref({})

  const today = new Date().toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })

  const loadStats = async () => {
    try {
      const res = await getStatisticsOverview()
      if (res.code === 0) {
        stats.value = res.data
      } else {
        ElMessage.error('加载统计数据失败')
      }
    } catch (err) {
      console.error(err)
    }
  }

  const goPage = (path) => {
    router.push(path)
  }

  const goHome = () => {
    router.push('/home')
  }

  onMounted(() => {
    loadStats()
  })

  defineOptions({
    name: 'Dashboard'
  })
</script>

<style lang="scss" scoped>
.stat-card {
  :deep(.el-card__body) {
    padding: 16px;
  }
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.stat-label {
  font-size: 13px;
  color: var(--el-text-color-secondary);
}

.quick-links {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.home-link-section {
  text-align: left;
}

.layer-info {
  font-size: 14px;
  color: var(--el-text-color-secondary);
}
</style>

