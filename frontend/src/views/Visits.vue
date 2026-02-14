<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>就诊管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>新增就诊
          </el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="单据号">
          <el-input v-model="searchForm.visit_id" placeholder="请输入" clearable />
        </el-form-item>
        <el-form-item label="顾客">
          <el-input v-model="searchForm.customer_name" placeholder="请输入" clearable />
        </el-form-item>
        <el-form-item label="日期从">
          <el-date-picker v-model="searchForm.date_from" type="date" placeholder="选择日期" />
        </el-form-item>
        <el-form-item label="到">
          <el-date-picker v-model="searchForm.date_to" type="date" placeholder="选择日期" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="visit_id" label="单据号" />
        <el-table-column prop="customer.name" label="顾客" />
        <el-table-column prop="consultant.name" label="咨询师" />
        <el-table-column prop="visit_date" label="就诊时间" />
        <el-table-column prop="total_amount" label="总金额">
          <template #default="{ row }">
            ¥{{ row.total_amount?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link @click="handleViewDetail(row)">明细</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="handleCurrentChange"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>

    <!-- 就诊弹窗 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="form" ref="formRef" label-width="100px">
        <el-form-item label="单据号" required>
          <el-input v-model="form.visit_id" :disabled="!!form.id" />
        </el-form-item>
        <el-form-item label="顾客" required>
          <el-select v-model="form.customer_id" style="width: 100%">
            <el-option v-for="c in customers" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="咨询师">
          <el-select v-model="form.consultant_id" style="width: 100%">
            <el-option label="无" :value="null" />
            <el-option v-for="e in consultants" :key="e.id" :label="e.name" :value="e.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="就诊时间" required>
          <el-date-picker v-model="form.visit_date" type="datetime" style="width: 100%" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getVisitList, createVisit, updateVisit, deleteVisit } from '../api/visit'
import { getCustomerList } from '../api/customer'
import { getEmployeeList } from '../api/employee'

const loading = ref(false)
const tableData = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

const customers = ref([])
const consultants = ref([])

const searchForm = reactive({
  visit_id: '',
  customer_name: '',
  date_from: '',
  date_to: ''
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref()
const form = reactive({
  id: null,
  visit_id: '',
  customer_id: null,
  consultant_id: null,
  visit_date: new Date(),
  remark: ''
})

const loadCustomers = async () => {
  const res = await getCustomerList({ page: 1, page_size: 1000 })
  customers.value = res.list
}

const loadConsultants = async () => {
  const res = await getEmployeeList({ role: '咨询师', page: 1, page_size: 100 })
  consultants.value = res.list
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getVisitList({
      page: page.value,
      page_size: pageSize.value,
      ...searchForm
    })
    tableData.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  page.value = 1
  loadData()
}

const handleReset = () => {
  Object.keys(searchForm).forEach(k => searchForm[k] = '')
  page.value = 1
  loadData()
}

const handleAdd = () => {
  dialogTitle.value = '新增就诊'
  Object.assign(form, {
    id: null, visit_id: '', customer_id: null,
    consultant_id: null, visit_date: new Date(), remark: ''
  })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑就诊'
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleViewDetail = (row) => {
  // TODO: 跳转到明细页或弹窗
  ElMessage.info('明细功能开发中')
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定删除？', '提示', { type: 'warning' })
    await deleteVisit(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (e) {}
}

const handleSubmit = async () => {
  try {
    if (form.id) {
      await updateVisit(form.id, form)
      ElMessage.success('更新成功')
    } else {
      await createVisit(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadData()
  } catch (e) {}
}

const handleCurrentChange = (val) => {
  page.value = val
  loadData()
}

onMounted(() => {
  loadData()
  loadCustomers()
  loadConsultants()
})
</script>

<style scoped>
.page-container { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
</style>
