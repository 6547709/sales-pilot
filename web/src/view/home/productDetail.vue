<template>
  <div class="product-detail" v-loading="loading">
    <div v-if="!loading && !product" class="not-found">
      <el-empty description="产品不存在" />
      <el-button @click="goBack">返回</el-button>
    </div>

    <div v-else-if="product" class="detail-content">
      <div class="detail-header">
        <el-button @click="goBack" text>
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <div class="header-info">
          <h1 class="product-name">{{ product.name }}</h1>
          <div class="product-meta">
            <el-tag v-if="product.vendor_market && product.vendor_market !== 'all'" size="small">
              {{ product.vendor_market === 'domestic' ? '国内' : '国外' }}
            </el-tag>
            <el-tag v-if="product.is_draft" type="warning" size="small">草稿</el-tag>
            <span class="category">{{ product.category }}</span>
          </div>
        </div>
      </div>

      <el-card class="detail-card">
        <template #header>
          <div class="card-header">
            <span>基本信息</span>
          </div>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="产品名称">{{ product.name }}</el-descriptions-item>
          <el-descriptions-item label="分类">{{ product.category }}</el-descriptions-item>
          <el-descriptions-item label="厂商名称">{{ product.manufacturer_name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="市场定位">
            <el-tag v-if="product.vendor_market === 'domestic'" type="success" size="small">国内</el-tag>
            <el-tag v-else-if="product.vendor_market === 'foreign'" type="warning" size="small">国外</el-tag>
            <el-tag v-else size="small">不区分</el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <el-card v-if="product.description" class="detail-card">
        <template #header>
          <div class="card-header">
            <span>产品描述</span>
          </div>
        </template>
        <div class="description-text">{{ product.description }}</div>
      </el-card>

      <el-card v-if="highlights.length > 0" class="detail-card">
        <template #header>
          <div class="card-header">
            <span>产品亮点</span>
          </div>
        </template>
        <ul class="highlights-list">
          <li v-for="(item, index) in highlights" :key="index">{{ item }}</li>
        </ul>
      </el-card>

      <el-card v-if="targetPersonas.length > 0" class="detail-card">
        <template #header>
          <div class="card-header">
            <span>目标客户画像</span>
          </div>
        </template>
        <ul class="highlights-list">
          <li v-for="(item, index) in targetPersonas" :key="index">{{ item }}</li>
        </ul>
      </el-card>

      <el-card v-if="product.trigger_events" class="detail-card">
        <template #header>
          <div class="card-header">
            <span>触发场景</span>
          </div>
        </template>
        <div class="description-text">{{ product.trigger_events }}</div>
      </el-card>

      <el-card v-if="discoveryQuestions.length > 0" class="detail-card">
        <template #header>
          <div class="card-header">
            <span>黄金三问</span>
          </div>
        </template>
        <ul class="highlights-list">
          <li v-for="(item, index) in discoveryQuestions" :key="index">{{ item }}</li>
        </ul>
      </el-card>

      <el-card v-if="product.competitor_analysis" class="detail-card">
        <template #header>
          <div class="card-header">
            <span>竞品分析</span>
          </div>
        </template>
        <div class="description-text">{{ product.competitor_analysis }}</div>
      </el-card>

      <el-card v-if="product.roi_metrics" class="detail-card">
        <template #header>
          <div class="card-header">
            <span>ROI 指标</span>
          </div>
        </template>
        <div class="description-text">{{ product.roi_metrics }}</div>
      </el-card>

      <el-card class="detail-card">
        <template #header>
          <div class="card-header">
            <span>联系方式</span>
          </div>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="销售联系人" :span="2">
            {{ product.sales_contact_name || '-' }}
            <template v-if="product.sales_contact_phone"> ({{ product.sales_contact_phone }})</template>
            <template v-if="product.sales_contact_email"> {{ product.sales_contact_email }}</template>
          </el-descriptions-item>
          <el-descriptions-item label="售前联系人" :span="2">
            {{ product.presales_contact_name || '-' }}
            <template v-if="product.presales_contact_phone"> ({{ product.presales_contact_phone }})</template>
            <template v-if="product.presales_contact_email"> {{ product.presales_contact_email }}</template>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <el-card v-if="scripts.length > 0" class="detail-card">
        <template #header>
          <div class="card-header">
            <span>销售话术 ({{ scripts.length }})</span>
          </div>
        </template>
        <div v-for="script in scripts" :key="script.id" class="script-item">
          <h4>{{ script.scenario }}</h4>
          <p>{{ script.content }}</p>
        </div>
      </el-card>

      <el-card v-if="cases.length > 0" class="detail-card">
        <template #header>
          <div class="card-header">
            <span>客户案例 ({{ cases.length }})</span>
          </div>
        </template>
        <div v-for="c in cases" :key="c.id" class="case-item">
          <h4>{{ c.client_name }}</h4>
          <div class="case-section">
            <strong>痛点：</strong>{{ c.pain_points }}
          </div>
          <div class="case-section">
            <strong>解决方案：</strong>{{ c.solution }}
          </div>
          <div class="case-section">
            <strong>价值交付：</strong>{{ c.value_delivered }}
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { getProduct } from '@/api/business/product'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const product = ref(null)
const scripts = ref([])
const cases = ref([])

const highlights = computed(() => {
  if (!product.value?.highlights) return []
  if (Array.isArray(product.value.highlights)) return product.value.highlights
  try {
    return JSON.parse(product.value.highlights)
  } catch {
    return []
  }
})

const targetPersonas = computed(() => {
  if (!product.value?.target_personas) return []
  if (Array.isArray(product.value.target_personas)) return product.value.target_personas
  try {
    return JSON.parse(product.value.target_personas)
  } catch {
    return []
  }
})

const discoveryQuestions = computed(() => {
  if (!product.value?.discovery_questions) return []
  if (Array.isArray(product.value.discovery_questions)) return product.value.discovery_questions
  try {
    return JSON.parse(product.value.discovery_questions)
  } catch {
    return []
  }
})

const loadProduct = async () => {
  try {
    loading.value = true
    const id = route.params.id
    const res = await getProduct({ id })
    if (res.code === 0) {
      product.value = res.data.product
      scripts.value = res.data.scripts || []
      cases.value = res.data.cases || []
    } else {
      ElMessage.error('加载产品失败')
    }
  } catch (err) {
    console.error(err)
    ElMessage.error('加载产品失败')
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.back()
}

onMounted(() => {
  loadProduct()
})

defineOptions({
  name: 'ProductDetail'
})
</script>

<style scoped>
.product-detail {
  padding: 24px;
  max-width: 1000px;
  margin: 0 auto;
  background: #f5f7fa;
  min-height: 100vh;
}

.not-found {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 64px;
}

.detail-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.detail-header {
  background: #fff;
  padding: 16px 20px;
  border-radius: 8px;
  margin-bottom: 8px;
}

.header-info {
  margin-top: 12px;
}

.product-name {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 8px 0;
  color: var(--el-text-color-primary);
}

.product-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.category {
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

.detail-card {
  border-radius: 8px;
}

.card-header {
  font-weight: 600;
}

.description-text {
  color: var(--el-text-color-primary);
  line-height: 1.6;
  white-space: pre-wrap;
}

.highlights-list {
  margin: 0;
  padding-left: 20px;
}

.highlights-list li {
  margin-bottom: 8px;
  color: var(--el-text-color-primary);
  line-height: 1.6;
}

.highlights-list li:last-child {
  margin-bottom: 0;
}

.script-item,
.case-item {
  padding: 12px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.script-item:last-child,
.case-item:last-child {
  border-bottom: none;
}

.script-item h4,
.case-item h4 {
  margin: 0 0 8px 0;
  color: var(--el-color-primary);
  font-size: 15px;
}

.script-item p {
  margin: 0;
  color: var(--el-text-color-primary);
  line-height: 1.6;
}

.case-section {
  margin-bottom: 6px;
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

.case-section strong {
  color: var(--el-text-color-primary);
}
</style>
