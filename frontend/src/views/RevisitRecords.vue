<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>回访记录</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>新增记录
          </el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="护士">
          <el-select v-model="searchForm.nurse_id" clearable style="width: 150px">
            <el-option v-for="n in nurses" :key="n.id" :label="n.name" :value="n.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期">
          <el-date-picker v-model="searchForm.date" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="nurse.name" label="护士" />
        <el-table-column prop="date" label="日期" />
        <el-table-column prop="reception_count" label="接待人数" />
        <el-table-column prop="add_wechat_count" label="加微人数" />
        <el-table-column prop="revisit_count" label="回访人数" />
        <el-table-column prop="remark" label="备注" show-overflow-tooltip />
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="form" ref="formRef" label-width="100px">
        <el-form-item label="护士" required>
          <el-select v-model="form.nurse_id" style="width: 100%">
            <el-option v-for="n in nurses" :key="n.id" :label="n.name" :value="n.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期" required>
          <el-date-picker v-model="form.date" type="date" style="width: 100%" />
        </el-form-item>
        <el-form-item label="接待人数">
          <el-input-number v-model="form.reception_count" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="加微人数">
          <el-input-number v-model="form.add_wechat_count" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="回访人数">
          <el-input-number v-model="form.revisit_count" :min="0" style="width: 100%" />
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
import { getEmployeeList } from '../api/employee'

const loading = ref(false)
const tableData = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const nurses = ref([])

const searchForm = reactive({
  nurse_id: '',
  date: ''
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref()
const form = reactive({
  id: null,
  nurse_id: null,
  date: new Date(),
  reception_count: 0,
  add_wechat_count: 0,
  revisit_count: 0,
  remark: ''
})

const loadNurses = async () => {
  const res = await getEmployeeList({ role: '护士', page: 1, page_size: 100 })
  nurses.value = res.list
}

const loadData = async () => {
  loading.value = true
  // TODO: 调用回访记录API
  tableData.value = []
  total.value = 0
  loading.value = false
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
  dialogTitle.value = '新增回访记录'
  Object.assign(form, {
    id: null, nurse_id: null, date: new Date(),
    reception_count: 0, add_wechat_count: 0, revisit_count: 0, remark: ''
  })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑回访记录'
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定删除？', '提示', { type: 'warning' })
    ElMessage.success('删除成功')
    loadData()
  } catch (e) {}
}

const handleSubmit = async () => {
  try {
    // TODO: 调用API
    ElMessage.success(form.id ? '更新成功' : '创建成功')
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
  loadNurses()
})
</script>

<style scoped>
.page-container { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
</style>
