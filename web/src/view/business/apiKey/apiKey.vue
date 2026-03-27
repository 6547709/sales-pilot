<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDrawer('create')"
          >创建API密钥</el-button
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
        <el-table-column align="left" label="名称" prop="Name" width="150" />
        <el-table-column align="left" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.IsActive ? 'success' : 'danger'" size="small">
              {{ scope.row.IsActive ? '活跃' : '已吊销' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="过期时间" width="180">
          <template #default="scope">
            <span>{{ scope.row.ExpiresAt ? formatDate(scope.row.ExpiresAt) : '永不过期' }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="最后使用" width="180">
          <template #default="scope">
            <span>{{ scope.row.LastUsedAt ? formatDate(scope.row.LastUsedAt) : '从未使用' }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" min-width="200">
          <template #default="scope">
            <el-button
              v-if="scope.row.IsActive"
              type="warning"
              link
              icon="close"
              @click="revokeAPIKey(scope.row)"
              >吊销</el-button
            >
            <el-button
              type="danger"
              link
              icon="delete"
              @click="deleteAPIKey(scope.row)"
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
          <span class="text-lg">创建API密钥</span>
          <div>
            <el-button @click="closeDrawer">取 消</el-button>
            <el-button type="primary" @click="enterDrawer">创 建</el-button>
          </div>
        </div>
      </template>
      <el-form :model="form" label-width="100px">
        <el-form-item label="密钥名称">
          <el-input v-model="form.Name" placeholder="给密钥起个名字" />
        </el-form-item>
        <el-form-item label="过期时间">
          <el-date-picker
            v-model="form.ExpiresAt"
            type="datetime"
            placeholder="不设置则永不过期"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- Show API Key Dialog -->
    <el-dialog v-model="showKeyDialog" title="请保存您的API密钥" width="500px" :close-on-click-modal="false">
      <div class="text-center">
        <p class="mb-4">请立即复制并保存您的API密钥，此密钥仅显示一次。</p>
        <el-input v-model="rawKey" readonly class="mb-4" />
        <el-button type="primary" @click="copyKey">复制密钥</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
  import {
    createAPIKey,
    getAPIKeys,
    revokeAPIKey as revokeAPIKeyApi,
    deleteAPIKey as deleteAPIKeyApi
  } from '@/api/business/apiKey'
  import { formatDate } from '@/utils/format'
  import { ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'

  defineOptions({
    name: 'APIKey'
  })

  const form = ref({
    Name: '',
    ExpiresAt: null
  })

  const tableData = ref([])
  const drawerFormVisible = ref(false)
  const showKeyDialog = ref(false)
  const rawKey = ref('')

  // 查询
  const getTableData = async () => {
    const res = await getAPIKeys()
    if (res.code === 0) {
      tableData.value = res.data.apiKeys
    }
  }

  getTableData()

  const closeDrawer = () => {
    drawerFormVisible.value = false
    resetForm()
  }

  const resetForm = () => {
    form.value = {
      Name: '',
      ExpiresAt: null
    }
  }

  const openDrawer = () => {
    resetForm()
    drawerFormVisible.value = true
  }

  const revokeAPIKeyHandler = async (row) => {
    ElMessageBox.confirm('确定要吊销此API密钥吗？吊销后该密钥将立即失效。', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const res = await revokeAPIKeyApi({ ID: row.ID })
      if (res.code === 0) {
        ElMessage({ type: 'success', message: '吊销成功' })
        getTableData()
      }
    })
  }

  const deleteAPIKeyHandler = async (row) => {
    ElMessageBox.confirm('确定要删除此API密钥吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const res = await deleteAPIKeyApi({ ID: row.ID })
      if (res.code === 0) {
        ElMessage({ type: 'success', message: '删除成功' })
        getTableData()
      }
    })
  }

  const enterDrawer = async () => {
    const res = await createAPIKey(form.value)
    if (res.code === 0) {
      closeDrawer()
      rawKey.value = res.data.rawKey
      showKeyDialog.value = true
      getTableData()
    }
  }

  const copyKey = async () => {
    try {
      await navigator.clipboard.writeText(rawKey.value)
      ElMessage({ type: 'success', message: '已复制到剪贴板' })
    } catch (err) {
      ElMessage({ type: 'error', message: '复制失败' })
    }
  }
</script>

<style scoped>
.mb-4 {
  margin-bottom: 16px;
}
</style>
