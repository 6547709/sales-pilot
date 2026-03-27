<template>
  <section class="product-spotlight">
    <div class="spotlight-header">
      <div class="spotlight-badge">
        <el-icon><Collection /></el-icon>
        <span>赋能方案</span>
      </div>
      <h2 class="spotlight-title">
        <template v-if="resolvedLabel">
          与「<span class="text-primary">{{ resolvedLabel }}</span>」相关的方案
          <span v-if="marketSuffix" class="market-suffix">({{ marketSuffix }})</span>
        </template>
        <template v-else>
          精选方案
        </template>
      </h2>
      <p class="spotlight-desc">
        已通过「所属解决方案」关联。列表中的「国内 / 国外」来自后台市场定位；全局架构卡仅展示您已录入的真实产品。
      </p>
    </div>

    <div v-loading="loading" class="product-list">
      <div v-if="!loading && products.length === 0" class="empty-state">
        <el-empty description="暂无匹配方案" />
      </div>

      <div v-else class="product-grid">
        <el-card
          v-for="(product, index) in products"
          :key="product.id"
          class="product-card"
          :style="{ animationDelay: `${index * 0.04}s` }"
          @click="goToDetail(product.id)"
        >
          <div class="product-header">
            <h3 class="product-name">{{ product.name }}</h3>
            <el-tag
              v-if="product.vendor_market && product.vendor_market !== 'all'"
              size="small"
              :type="product.vendor_market === 'domestic' ? 'success' : 'warning'"
            >
              {{ product.vendor_market === 'domestic' ? '国内' : '国外' }}
            </el-tag>
          </div>
          <p class="product-category">{{ product.category || '方案' }}</p>
          <p v-if="product.manufacturer_name" class="product-manufacturer">
            {{ product.manufacturer_name }}
          </p>
        </el-card>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Collection } from '@element-plus/icons-vue'
import { getProductsByCategory } from '@/api/business/product'
import { ElMessage } from 'element-plus'

const props = defineProps({
  solutionId: {
    type: String,
    default: ''
  },
  vendorMarket: {
    type: String,
    default: ''
  },
  keyword: {
    type: String,
    default: ''
  }
})

const router = useRouter()
const loading = ref(false)
const products = ref([])
const resolvedLabel = ref('')

const marketSuffix = computed(() => {
  if (props.vendorMarket === 'domestic') return '国内'
  if (props.vendorMarket === 'foreign') return '国外'
  return ''
})

const loadProducts = async () => {
  if (!props.solutionId) {
    products.value = []
    resolvedLabel.value = ''
    return
  }

  try {
    loading.value = true
    const res = await getProductsByCategory({
      solution_category_id: props.solutionId,
      vendor_market: props.vendorMarket
    })
    if (res.code === 0) {
      products.value = res.data || []
    } else {
      ElMessage.error('加载产品失败')
    }
  } catch (err) {
    console.error(err)
    products.value = []
  } finally {
    loading.value = false
  }
}

const goToDetail = (id) => {
  router.push(`/home/product/${id}`)
}

watch(() => [props.solutionId, props.vendorMarket], () => {
  loadProducts()
}, { immediate: true })

onMounted(() => {
  loadProducts()
})
</script>

<style scoped>
.product-spotlight {
  padding: 32px 24px;
  max-width: 1400px;
  margin: 0 auto;
  border-top: 1px solid var(--el-border-color-lighter);
  background: var(--el-fill-color-light);
}

.spotlight-header {
  margin-bottom: 24px;
}

.spotlight-badge {
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

.spotlight-title {
  font-size: 22px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 8px 0;
}

.spotlight-desc {
  font-size: 14px;
  color: var(--el-text-color-secondary);
  margin: 0;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px;
}

.empty-state {
  padding: 48px;
  text-align: center;
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.product-card {
  cursor: pointer;
  transition: all 0.2s;
  animation: fadeInUp 0.4s ease-out forwards;
  opacity: 0;
}

.product-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px var(--el-color-primary-light-7);
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.product-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 8px;
}

.product-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0;
  line-height: 1.4;
}

.product-category {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  margin: 0 0 4px 0;
}

.product-manufacturer {
  font-size: 12px;
  color: var(--el-text-color-placeholder);
  margin: 0;
}

.text-primary {
  color: var(--el-color-primary);
}

.market-suffix {
  color: var(--el-text-color-secondary);
  font-weight: 400;
}
</style>
