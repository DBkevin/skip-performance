<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-title">今日就诊</div>
          <div class="stat-value">{{ todayStats.visits }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-title">今日业绩</div>
          <div class="stat-value">¥{{ formatMoney(todayStats.amount) }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-title">本月业绩</div>
          <div class="stat-value">¥{{ formatMoney(monthStats.amount) }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-title">顾客总数</div>
          <div class="stat-value">{{ customerCount }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>近7天业绩趋势</span>
          </template>
          <div class="chart-placeholder">
            <el-empty description="图表功能开发中" />
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>员工业绩排行</span>
          </template>
          <el-table :data="performanceRanking" style="width: 100%">
            <el-table-column prop="name" label="姓名" />
            <el-table-column prop="role" label="角色" />
            <el-table-column prop="performance" label="业绩">
              <template #default="{ row }">
                ¥{{ formatMoney(row.performance) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getCustomerList } from '../api/customer'
import { getPerformanceReport } from '../api/report'

const todayStats = ref({ visits: 0, amount: 0 })
const monthStats = ref({ visits: 0, amount: 0 })
const customerCount = ref(0)
const performanceRanking = ref([])

const formatMoney = (value) => {
  if (!value) return '0.00'
  return value.toFixed(2)
}

const loadData = async () => {
  try {
    // 加载顾客总数
    const customerRes = await getCustomerList({ page: 1, page_size: 1 })
    customerCount.value = customerRes.total || 0

    // 加载业绩报表
    const today = new Date().toISOString().split('T')[0]
    const reportRes = await getPerformanceReport({ 
      date_from: today,
      date_to: today
    })
    if (reportRes && reportRes.reports) {
      todayStats.value.amount = reportRes.total_amount || 0
      performanceRanking.value = reportRes.reports.slice(0, 10)
    }
  } catch (error) {
    console.error('加载数据失败:', error)
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.stat-card {
  text-align: center;
}

.stat-title {
  color: #666;
  font-size: 14px;
  margin-bottom: 10px;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #409EFF;
}

.chart-placeholder {
  height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
