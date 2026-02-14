<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>业绩报表</span>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="日期从">
          <el-date-picker v-model="searchForm.date_from" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item label="到">
          <el-date-picker v-model="searchForm.date_to" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <div class="summary" v-if="reportData">
        <el-row :gutter="20">
          <el-col :span="8">
            <div class="summary-item">
              <div class="label">总金额</div>
              <div class="value">¥{{ reportData.total_amount?.toFixed(2) }}</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="summary-item">
              <div class="label">总业绩</div>
              <div class="value">¥{{ totalPerformance?.toFixed(2) }}</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="summary-item">
              <div class="label">人数</div>
              <div class="value">{{ reportData.reports?.length }}</div>
            </div>
          </el-col>
        </el-row>
      </div>

      <el-table :data="reportData?.reports" v-loading="loading" stripe>
        <el-table-column prop="employee_name" label="姓名" />
        <el-table-column prop="employee_role" label="角色" />
        <el-table-column prop="main_performance" label="主操业绩">
          <template #default="{ row }">¥{{ row.main_performance?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="co_performance" label="协同业绩">
          <template #default="{ row }">¥{{ row.co_performance?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="nurse_performance" label="护士业绩">
          <template #default="{ row }">¥{{ row.nurse_performance?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="total_performance" label="总业绩" sortable>
          <template #default="{ row }">
            <strong>¥{{ row.total_performance?.toFixed(2) }}</strong>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { getPerformanceReport } from '../api/report'

const loading = ref(false)
const reportData = ref(null)

const searchForm = reactive({
  date_from: '',
  date_to: ''
})

const totalPerformance = computed(() => {
  if (!reportData.value?.reports) return 0
  return reportData.value.reports.reduce((sum, r) => sum + (r.total_performance || 0), 0)
})

const loadData = async () => {
  loading.value = true
  try {
    const res = await getPerformanceReport(searchForm)
    reportData.value = res
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  loadData()
}

const handleReset = () => {
  const today = new Date()
  const firstDay = new Date(today.getFullYear(), today.getMonth(), 1)
  searchForm.date_from = firstDay.toISOString().split('T')[0]
  searchForm.date_to = today.toISOString().split('T')[0]
  loadData()
}

onMounted(() => {
  handleReset()
})
</script>

<style scoped>
.page-container { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.summary {
  margin-bottom: 20px;
  padding: 20px;
  background: #f5f7fa;
  border-radius: 4px;
}
.summary-item {
  text-align: center;
}
.summary-item .label {
  color: #666;
  font-size: 14px;
  margin-bottom: 8px;
}
.summary-item .value {
  font-size: 24px;
  font-weight: bold;
  color: #409EFF;
}
</style>
