<template>
  <div style="margin: 10px 0">
    <el-input v-model="search" placeholder="请输入关键字" style="width: 20%" clearable></el-input>
    <el-button type="primary" style="margin-left: 5px" @click="load">查询</el-button>
    <el-button type="primary" style="margin-left: 5px" @click="updateBasicBehaviorData">更新基础行为数据</el-button>
    
    
    <el-table 
      v-loading="loading"
      :data="tableData" 
      style="width: 100%"
    >
      <el-table-column fixed prop="component_id" label="ComponentId" width="150" />
      <el-table-column prop="component_name" label="ComponentName" width="320" />
      <el-table-column prop="component_type" label="ComponentType" width="120" />
      <el-table-column prop="component_desc" label="ComponentDesc" width="220" />
    </el-table>

    <div style="margin: 10px 0">
      <el-pagination
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="currentPage"
          :page-sizes="[5, 10, 20]"
          :page-size="pageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total">
      </el-pagination>
    </div>
  
  </div>
</template>

<script>
export default {
  name: "BasicBehaviorData",
  data() {
    return {
      search: '',
      loading: false,
      tableData: [],
      currentPage: 1,
      pageSize: 10,
      total: 0,
    }
  },
  created() {
    this.load()
  },
  methods: {
    load() {
      this.loading = true

      this.loading = false
    },
    handleSizeChange(pageSize) {  // 改变每页的大小
      this.pageSize = pageSize;
      this.load();
    },
    handleCurrentChange(pageNum) {  //改变当前页码
      this.currentPage = pageNum;
      this.load()
    },
    updateBasicBehaviorData() {

    }
  }
}
</script>

<style scoped>

</style>