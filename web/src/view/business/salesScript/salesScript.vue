<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDrawer('create')"
          >新增话术</el-button
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
        <el-table-column align="left" label="产品ID" prop="ProductID" width="100" />
        <el-table-column align="left" label="场景" prop="Scenario" min-width="200" />
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
              @click="deleteSalesScript(scope.row)"
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
      size="60%"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ drawerType === 'create' ? '新增话术' : '编辑话术' }}</span>
          <div>
            <el-button @click="closeDrawer">取 消</el-button>
            <el-button type="primary" @click="enterDrawer">确 定</el-button>
          </div>
        </div>
      </template>
      <el-form :inline="true" :model="form" label-width="100px">
        <el-form-item label="产品ID">
          <el-input-number v-model="form.ProductID" :min="1" />
        </el-form-item>
        <el-form-item label="场景">
          <el-input v-model="form.Scenario" autocomplete="off" style="width: 300px" />
        </el-form-item>
        <el-form-item label="话术内容" class="w-full">
          <el-input
            v-model="form.Content"
            type="textarea"
            :rows="10"
            placeholder="请输入销售话术内容"
          />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    createSalesScript,
    updateSalesScript,
    deleteSalesScript,
    getSalesScript,
    getSalesScriptsByProductID
  } from '@/api/business/salesScript'
  import { formatDate } from '@/utils/format'
  import { ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'

  defineOptions({
    name: 'SalesScript'
  })

  const form = ref({
    ID: 0,
    ProductID: 1,
    Scenario: '',
    Content: ''
  })

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
    const table = await getSalesScriptsByProductID({
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
      const res = await getSalesScript({ id: row.ID })
      if (res.code === 0) {
        form.value = { ...res.data.salesScript }
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
      ID: 0,
      ProductID: 1,
      Scenario: '',
      Content: ''
    }
  }

  const deleteSalesScriptHandler = async (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const res = await deleteSalesScript({ ID: row.ID })
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
    let res
    if (drawerType.value === 'create') {
      res = await createSalesScript(form.value)
    } else {
      res = await updateSalesScript(form.value)
    }

    if (res.code === 0) {
      closeDrawer()
      getTableData()
    }
  }
</script>

<style scoped>
.w-full {
  width: 100%;
}
</style>
