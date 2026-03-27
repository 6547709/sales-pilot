<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDrawer('create')"
          >新增分类</el-button
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
        <el-table-column align="left" label="分类名称" prop="Name" min-width="200" />
        <el-table-column align="left" label="描述" prop="Description" min-width="300" show-overflow-tooltip />
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
              @click="deleteSolutionCategory(scope.row)"
              >删除</el-button
            >
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-drawer
      v-model="drawerFormVisible"
      :before-close="closeDrawer"
      :show-close="false"
      size="400px"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ drawerType === 'create' ? '新增分类' : '编辑分类' }}</span>
          <div>
            <el-button @click="closeDrawer">取 消</el-button>
            <el-button type="primary" @click="enterDrawer">确 定</el-button>
          </div>
        </div>
      </template>
      <el-form :model="form" label-width="100px">
        <el-form-item label="分类名称">
          <el-input v-model="form.Name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.Description" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    createSolutionCategory,
    updateSolutionCategory,
    deleteSolutionCategory,
    getSolutionCategory,
    getAllSolutionCategories
  } from '@/api/business/solutionCategory'
  import { formatDate } from '@/utils/format'
  import { ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'

  defineOptions({
    name: 'SolutionCategory'
  })

  const form = ref({
    ID: 0,
    Name: '',
    Description: ''
  })

  const tableData = ref([])

  // 查询
  const getTableData = async () => {
    const res = await getAllSolutionCategories()
    if (res.code === 0) {
      tableData.value = res.data.categories
    }
  }

  getTableData()

  const drawerFormVisible = ref(false)
  const drawerType = ref('create')

  const openDrawer = async (type, row) => {
    drawerType.value = type
    if (type === 'update' && row) {
      const res = await getSolutionCategory({ id: row.ID })
      if (res.code === 0) {
        form.value = { ...res.data.solutionCategory }
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
      Name: '',
      Description: ''
    }
  }

  const deleteSolutionCategoryHandler = async (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const res = await deleteSolutionCategory({ ID: row.ID })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        getTableData()
      }
    })
  }

  const enterDrawer = async () => {
    let res
    if (drawerType.value === 'create') {
      res = await createSolutionCategory(form.value)
    } else {
      res = await updateSolutionCategory(form.value)
    }

    if (res.code === 0) {
      closeDrawer()
      getTableData()
    }
  }
</script>
