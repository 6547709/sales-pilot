<template>
  <div>
    <el-tabs v-model="activeTab" class=" topology-tabs">
      <el-tab-pane label="拓扑层级" name="layer">
        <div class="gva-btn-list mb-4">
          <el-button type="primary" icon="plus" @click="openLayerDrawer('create')"
            >新增层级</el-button
          >
        </div>
        <el-table :data="layers" style="width: 100%" row-key="ID">
          <el-table-column prop="Name" label="层级名称" width="150" />
          <el-table-column prop="SortOrder" label="排序" width="100" />
          <el-table-column label="操作" width="160">
            <template #default="scope">
              <el-button type="primary" link icon="edit" @click="openLayerDrawer('update', scope.row)">编辑</el-button>
              <el-button type="danger" link icon="delete" @click="deleteLayer(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="拓扑分类" name="category">
        <div class="gva-btn-list mb-4">
          <el-button type="primary" icon="plus" @click="openCategoryDrawer('create')"
            >新增分类</el-button
          >
        </div>
        <el-table :data="categories" style="width: 100%" row-key="ID">
          <el-table-column prop="LayerName" label="所属层级" width="150" />
          <el-table-column prop="Name" label="分类名称" width="150" />
          <el-table-column label="操作" width="160">
            <template #default="scope">
              <el-button type="primary" link icon="edit" @click="openCategoryDrawer('update', scope.row)">编辑</el-button>
              <el-button type="danger" link icon="delete" @click="deleteCategory(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="拓扑厂商" name="vendor">
        <div class="gva-btn-list mb-4">
          <el-button type="primary" icon="plus" @click="openVendorDrawer('create')"
            >新增厂商</el-button
          >
        </div>
        <el-table :data="vendors" style="width: 100%" row-key="ID">
          <el-table-column prop="CategoryName" label="所属分类" width="150" />
          <el-table-column prop="Name" label="厂商名称" width="150" />
          <el-table-column prop="Website" label="官网" min-width="200" />
          <el-table-column label="操作" width="160">
            <template #default="scope">
              <el-button type="primary" link icon="edit" @click="openVendorDrawer('update', scope.row)">编辑</el-button>
              <el-button type="danger" link icon="delete" @click="deleteVendor(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <!-- Layer Drawer -->
    <el-drawer v-model="layerDrawerVisible" :before-close="closeLayerDrawer" size="400px">
      <template #header>
        <div class="flex justify-between items-center">
          <span>{{ layerDrawerType === 'create' ? '新增层级' : '编辑层级' }}</span>
          <div>
            <el-button @click="closeLayerDrawer">取 消</el-button>
            <el-button type="primary" @click="enterLayerDrawer">确 定</el-button>
          </div>
        </div>
      </template>
      <el-form :model="layerForm" label-width="100px">
        <el-form-item label="层级名称">
          <el-input v-model="layerForm.Name" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="layerForm.SortOrder" :min="0" />
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- Category Drawer -->
    <el-drawer v-model="categoryDrawerVisible" :before-close="closeCategoryDrawer" size="400px">
      <template #header>
        <div class="flex justify-between items-center">
          <span>{{ categoryDrawerType === 'create' ? '新增分类' : '编辑分类' }}</span>
          <div>
            <el-button @click="closeCategoryDrawer">取 消</el-button>
            <el-button type="primary" @click="enterCategoryDrawer">确 定</el-button>
          </div>
        </div>
      </template>
      <el-form :model="categoryForm" label-width="100px">
        <el-form-item label="所属层级">
          <el-select v-model="categoryForm.LayerID" placeholder="请选择层级">
            <el-option v-for="l in layers" :key="l.ID" :label="l.Name" :value="l.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="分类名称">
          <el-input v-model="categoryForm.Name" />
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- Vendor Drawer -->
    <el-drawer v-model="vendorDrawerVisible" :before-close="closeVendorDrawer" size="400px">
      <template #header>
        <div class="flex justify-between items-center">
          <span>{{ vendorDrawerType === 'create' ? '新增厂商' : '编辑厂商' }}</span>
          <div>
            <el-button @click="closeVendorDrawer">取 消</el-button>
            <el-button type="primary" @click="enterVendorDrawer">确 定</el-button>
          </div>
        </div>
      </template>
      <el-form :model="vendorForm" label-width="100px">
        <el-form-item label="所属分类">
          <el-select v-model="vendorForm.CategoryID" placeholder="请选择分类">
            <el-option v-for="c in categories" :key="c.ID" :label="c.Name" :value="c.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="厂商名称">
          <el-input v-model="vendorForm.Name" />
        </el-form-item>
        <el-form-item label="官网">
          <el-input v-model="vendorForm.Website" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="vendorForm.Description" type="textarea" rows="3" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    createLayer,
    updateLayer,
    deleteLayer,
    getLayers,
    createCategory,
    updateCategory,
    deleteCategory,
    getCategoriesByLayer,
    createVendor,
    updateVendor,
    deleteVendor,
    getVendorsByCategory
  } from '@/api/business/topology'
  import { ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'

  defineOptions({
    name: 'Topology'
  })

  const activeTab = ref('layer')

  // Layers
  const layers = ref([])
  const layerDrawerVisible = ref(false)
  const layerDrawerType = ref('create')
  const layerForm = ref({ ID: 0, Name: '', SortOrder: 0 })

  // Categories
  const categories = ref([])
  const categoryDrawerVisible = ref(false)
  const categoryDrawerType = ref('create')
  const categoryForm = ref({ ID: 0, LayerID: null, Name: '' })

  // Vendors
  const vendors = ref([])
  const vendorDrawerVisible = ref(false)
  const vendorDrawerType = ref('create')
  const vendorForm = ref({ ID: 0, CategoryID: null, Name: '', Website: '', Description: '' })

  // Load data
  const loadLayers = async () => {
    const res = await getLayers()
    if (res.code === 0) {
      layers.value = res.data.layers
    }
  }

  const loadCategories = async () => {
    const res = await getCategoriesByLayer()
    if (res.code === 0) {
      categories.value = res.data.categories
    }
  }

  const loadVendors = async () => {
    const res = await getVendorsByCategory()
    if (res.code === 0) {
      vendors.value = res.data.vendors
    }
  }

  loadLayers()

  // Layer handlers
  const openLayerDrawer = (type, row) => {
    layerDrawerType.value = type
    if (type === 'update' && row) {
      layerForm.value = { ...row }
    } else {
      layerForm.value = { ID: 0, Name: '', SortOrder: 0 }
    }
    layerDrawerVisible.value = true
  }

  const closeLayerDrawer = () => {
    layerDrawerVisible.value = false
  }

  const enterLayerDrawer = async () => {
    let res
    if (layerDrawerType.value === 'create') {
      res = await createLayer(layerForm.value)
    } else {
      res = await updateLayer(layerForm.value)
    }
    if (res.code === 0) {
      closeLayerDrawer()
      loadLayers()
    }
  }

  const deleteLayerHandler = async (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const res = await deleteLayer({ ID: row.ID })
      if (res.code === 0) {
        ElMessage({ type: 'success', message: '删除成功' })
        loadLayers()
      }
    })
  }

  // Category handlers
  const openCategoryDrawer = (type, row) => {
    categoryDrawerType.value = type
    if (type === 'update' && row) {
      categoryForm.value = { ...row }
    } else {
      categoryForm.value = { ID: 0, LayerID: null, Name: '' }
    }
    categoryDrawerVisible.value = true
  }

  const closeCategoryDrawer = () => {
    categoryDrawerVisible.value = false
  }

  const enterCategoryDrawer = async () => {
    let res
    if (categoryDrawerType.value === 'create') {
      res = await createCategory(categoryForm.value)
    } else {
      res = await updateCategory(categoryForm.value)
    }
    if (res.code === 0) {
      closeCategoryDrawer()
      loadCategories()
    }
  }

  const deleteCategoryHandler = async (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const res = await deleteCategory({ ID: row.ID })
      if (res.code === 0) {
        ElMessage({ type: 'success', message: '删除成功' })
        loadCategories()
      }
    })
  }

  // Vendor handlers
  const openVendorDrawer = (type, row) => {
    vendorDrawerType.value = type
    if (type === 'update' && row) {
      vendorForm.value = { ...row }
    } else {
      vendorForm.value = { ID: 0, CategoryID: null, Name: '', Website: '', Description: '' }
    }
    vendorDrawerVisible.value = true
  }

  const closeVendorDrawer = () => {
    vendorDrawerVisible.value = false
  }

  const enterVendorDrawer = async () => {
    let res
    if (vendorDrawerType.value === 'create') {
      res = await createVendor(vendorForm.value)
    } else {
      res = await updateVendor(vendorForm.value)
    }
    if (res.code === 0) {
      closeVendorDrawer()
      loadVendors()
    }
  }

  const deleteVendorHandler = async (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const res = await deleteVendor({ ID: row.ID })
      if (res.code === 0) {
        ElMessage({ type: 'success', message: '删除成功' })
        loadVendors()
      }
    })
  }

  // Tab change
  const handleTabChange = (tab) => {
    if (tab === 'category' && categories.value.length === 0) {
      loadCategories()
    } else if (tab === 'vendor' && vendors.value.length === 0) {
      loadVendors()
    }
  }
</script>

<style scoped>
.mb-4 {
  margin-bottom: 16px;
}
.topology-tabs {
  padding: 0 16px;
}
</style>
