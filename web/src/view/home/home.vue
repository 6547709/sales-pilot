<template>
  <div class="home-container">
    <solution-atlas
      :initial-solution-id="solutionId"
      :initial-market="market"
      :initial-keyword="keyword"
      @select-solution="onSelectSolution"
    />
    <product-spotlight
      :key="`${solutionId}-${market}-${keyword}`"
      :solution-id="solutionId"
      :vendor-market="market"
      :keyword="keyword"
    />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import SolutionAtlas from './solutionAtlas.vue'
import ProductSpotlight from './productSpotlight.vue'

const route = useRoute()
const router = useRouter()

const solutionId = ref(route.query.solution || '')
const market = ref(route.query.market || '')
const keyword = ref(route.query.kw || '')

const onSelectSolution = (cat) => {
  if (!cat) {
    solutionId.value = ''
    market.value = ''
    router.push('/home')
  } else {
    solutionId.value = String(cat.id)
    router.push(`/home?solution=${cat.id}`)
  }
}

// Watch for route changes
watch(() => route.query, (query) => {
  solutionId.value = query.solution || ''
  market.value = query.market || ''
  keyword.value = query.kw || ''
}, { immediate: true })

defineOptions({
  name: 'Home'
})
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background: linear-gradient(180deg, #f8fafc 0%, #fff 100%);
}
</style>
