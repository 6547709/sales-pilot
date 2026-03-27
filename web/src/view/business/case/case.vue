<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDrawer('create')"
          >新增案例</el-button
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
        <el-table-column align="left" label="客户名称" prop="ClientName" width="150" />
        <el-table-column align="left" label="痛点" prop="PainPoints" min-width="200" show-overflow-tooltip />
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
              @click="deleteCase(scope.row)"
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
      size="70%"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ drawerType === 'create' ? '新增案例' : '编辑案例' }}</span>
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
        <el-form-item label="客户名称">
          <el-input v-model="form.ClientName" autocomplete="off" style="width: 200px" />
        </el-form-item>
        <el-form-item label="痛点" class="w-full">
          <el-input
            v-model="form.PainPoints"
            type="textarea"
            :rows="3"
            placeholder="请输入客户痛点"
          />
        </el-form-item>
        <el-form-item label="解决方案" class="w-full">
          <el-input
            v-model="form.Solution"
            type="textarea"
            :rows="3"
            placeholder="请输入解决方案"
          />
        </el-form-item>
        <el-form-item label="价值交付" class="w-full">
          <el-input
            v-model="form.ValueDelivered"
            type="textarea"
            :rows="3"
            placeholder="请输入价值交付"
          />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    createCase,
    updateCase,
    deleteCase,
    getCase,
    getCasesByProductID
  } from '@/api/business/case'
  import { formatDate } from '@/utils/format'
  import { ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'

  defineOptions({
    name: 'Case'
  })

  const form = ref({
    ID: 0,
    ProductID: 1,
    ClientName: '',
    PainPoints: '',
    Solution: '',
    ValueDelivered: ''
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
    const table = await getCasesByProductID({
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
      const res = await getCase({ id: row.ID })
      if (res.code === 0) {
        form.value = { ...res.data.caseInfo }
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
      ClientName: '',
      PainPoints: '',
      Solution: '',
      ValueDelivered: ''
    }
  }

  const deleteCaseHandler = async (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const res = await deleteCase({ ID: row.ID })
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
      res = await createCase(form.value)
    } else {
      res = await updateCase(form.value)
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
