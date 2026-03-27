<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDrawer('create')"
          >新增产品</el-button
        >
      </div>
      <el-table
        ref="multipleTable"
        :data="tableData"
        style="width: 100%"
        tooltip-effect="dark"
        row-key="ID"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="创建日期" width="180">
          <template #default="scope">
            <span>{{ formatDate(scope.row.CreatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="产品名称" prop="Name" width="180" />
        <el-table-column align="left" label="厂商" prop="ManufacturerName" width="120" />
        <el-table-column align="left" label="分类" prop="Category" width="100" />
        <el-table-column align="left" label="市场" prop="VendorMarket" width="80" />
        <el-table-column align="left" label="草稿" width="60">
          <template #default="scope">
            <el-tag :type="scope.row.IsDraft ? 'warning' : 'success'" size="small">
              {{ scope.row.IsDraft ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" min-width="160">
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="edit"
              @click="openDrawer('update', scope.row)"
              >编辑</el-button
            >
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteProduct(scope.row)"
              >删除</el-button
            >
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-drawer
      v-model="drawerFormVisible"
      :before-close="closeDrawer"
      :show-close="false"
      size="80%"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ drawerType === 'create' ? '新增产品' : '编辑产品' }}</span>
          <div>
            <el-button @click="closeDrawer">取 消</el-button>
            <el-button type="primary" @click="enterDrawer">确 定</el-button>
          </div>
        </div>
      </template>
      <el-form :inline="true" :model="form" label-width="120px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="产品名称">
              <el-input v-model="form.Name" autocomplete="off" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="厂商名称">
              <el-input v-model="form.ManufacturerName" autocomplete="off" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="产品分类">
              <el-input v-model="form.Category" autocomplete="off" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="市场">
              <el-select v-model="form.VendorMarket" placeholder="请选择">
                <el-option label="全部" value="all" />
                <el-option label="国内" value="domestic" />
                <el-option label="国外" value="foreign" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="销售联系人">
              <el-input v-model="form.SalesContactName" autocomplete="off" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="销售电话">
              <el-input v-model="form.SalesContactPhone" autocomplete="off" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="销售邮箱">
              <el-input v-model="form.SalesContactEmail" autocomplete="off" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="售前联系人">
              <el-input v-model="form.PresalesContactName" autocomplete="off" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="售前电话">
              <el-input v-model="form.PresalesContactPhone" autocomplete="off" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="售前邮箱">
              <el-input v-model="form.PresalesContactEmail" autocomplete="off" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="产品描述">
          <el-input v-model="form.Description" type="textarea" rows="3" autocomplete="off" />
        </el-form-item>
        <el-form-item label="三大亮点">
          <div v-for="(item, index) in highlights" :key="index" class="flex mb-2">
            <el-input v-model="highlights[index]" placeholder="输入亮点" />
            <el-button type="danger" link icon="delete" @click="removeHighlight(index)">删除</el-button>
          </div>
          <el-button type="primary" link @click="addHighlight">+ 添加亮点</el-button>
        </el-form-item>
        <el-form-item label="目标客户画像">
          <div v-for="(item, index) in targetPersonas" :key="index" class="flex mb-2">
            <el-input v-model="targetPersonas[index]" placeholder="输入目标画像" />
            <el-button type="danger" link icon="delete" @click="removeTargetPersona(index)">删除</el-button>
          </div>
          <el-button type="primary" link @click="addTargetPersona">+ 添加目标画像</el-button>
        </el-form-item>
        <el-form-item label="触发事件">
          <el-input v-model="form.TriggerEvents" type="textarea" rows="2" autocomplete="off" />
        </el-form-item>
        <el-form-item label="黄金三问">
          <div v-for="(item, index) in discoveryQuestions" :key="index" class="flex mb-2">
            <el-input v-model="discoveryQuestions[index]" placeholder="输入问题" />
            <el-button type="danger" link icon="delete" @click="removeDiscoveryQuestion(index)">删除</el-button>
          </div>
          <el-button type="primary" link @click="addDiscoveryQuestion">+ 添加问题</el-button>
        </el-form-item>
        <el-form-item label="竞品分析">
          <el-input v-model="form.CompetitorAnalysis" type="textarea" rows="3" autocomplete="off" />
        </el-form-item>
        <el-form-item label="ROI指标">
          <el-input v-model="form.ROIMetrics" type="textarea" rows="2" autocomplete="off" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="是否草稿">
              <el-switch v-model="form.IsDraft" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    createProduct,
    updateProduct,
    deleteProduct,
    getProduct,
    getProductList
  } from '@/api/business/product'
  import { formatDate } from '@/utils/format'
  import { ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'

  defineOptions({
    name: 'Product'
  })

  const form = ref({
    Name: '',
    Category: '',
    SolutionCategoryID: null,
    ManufacturerName: '',
    SalesContactName: '',
    SalesContactPhone: '',
    SalesContactEmail: '',
    PresalesContactName: '',
    PresalesContactPhone: '',
    PresalesContactEmail: '',
    Description: '',
    Highlights: [],
    TargetPersonas: [],
    TriggerEvents: '',
    DiscoveryQuestions: [],
    CompetitorAnalysis: '',
    ROIMetrics: '',
    VendorMarket: 'all',
    IsDraft: false
  })

  const highlights = ref([])
  const targetPersonas = ref([])
  const discoveryQuestions = ref([])

  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])

  // 分页
  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  // 查询
  const getTableData = async () => {
    const table = await getProductList({
      page: page.value,
      pageSize: pageSize.value
    })
    if (table.code === 0) {
      tableData.value = table.data.list
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  }

  getTableData()

  const drawerFormVisible = ref(false)
  const drawerType = ref('create')

  const openDrawer = async (type, row) => {
    drawerType.value = type
    if (type === 'update' && row) {
      const res = await getProduct({ id: row.ID })
      if (res.code === 0) {
        form.value = { ...res.data.product }
        highlights.value = res.data.product.Highlights || []
        targetPersonas.value = res.data.product.TargetPersonas || []
        discoveryQuestions.value = res.data.product.DiscoveryQuestions || []
      }
    } else {
      resetForm()
    }
    drawerFormVisible.value = true
  }

  const closeDrawer = () => {
    drawerFormVisible.value = false
    resetForm()
  }

  const resetForm = () => {
    form.value = {
      Name: '',
      Category: '',
      SolutionCategoryID: null,
      ManufacturerName: '',
      SalesContactName: '',
      SalesContactPhone: '',
      SalesContactEmail: '',
      PresalesContactName: '',
      PresalesContactPhone: '',
      PresalesContactEmail: '',
      Description: '',
      Highlights: [],
      TargetPersonas: [],
      TriggerEvents: '',
      DiscoveryQuestions: [],
      CompetitorAnalysis: '',
      ROIMetrics: '',
      VendorMarket: 'all',
      IsDraft: false
    }
    highlights.value = []
    targetPersonas.value = []
    discoveryQuestions.value = []
  }

  const deleteProductHandler = async (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const res = await deleteProduct({ ID: row.ID })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === 1 && page.value > 1) {
          page.value--
        }
        getTableData()
      }
    })
  }

  const enterDrawer = async () => {
    form.value.Highlights = highlights.value
    form.value.TargetPersonas = targetPersonas.value
    form.value.DiscoveryQuestions = discoveryQuestions.value

    let res
    if (drawerType.value === 'create') {
      res = await createProduct(form.value)
    } else {
      res = await updateProduct(form.value)
    }

    if (res.code === 0) {
      closeDrawer()
      getTableData()
    }
  }

  const addHighlight = () => {
    highlights.value.push('')
  }

  const removeHighlight = (index) => {
    highlights.value.splice(index, 1)
  }

  const addTargetPersona = () => {
    targetPersonas.value.push('')
  }

  const removeTargetPersona = (index) => {
    targetPersonas.value.splice(index, 1)
  }

  const addDiscoveryQuestion = () => {
    discoveryQuestions.value.push('')
  }

  const removeDiscoveryQuestion = (index) => {
    discoveryQuestions.value.splice(index, 1)
  }
</script>

<style scoped>
.mb-2 {
  margin-bottom: 8px;
}
</style>
