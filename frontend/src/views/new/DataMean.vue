<template>
  <div style="margin: 10px 0">
    <el-input v-model="search" placeholder="请输入关键字" style="width: 20%" clearable></el-input>
    <el-button type="primary" style="margin-left: 5px" @click="load">查询</el-button>
    <el-button type="primary" style="margin-left: 5px" @click="updateComponent">更新组件数据</el-button>

    <el-table 
      v-loading="loading"
      :data="tableData" 
      style="width: 100%"
    >
      <el-table-column fixed prop="component_id" label="ID" width="150" />
      <el-table-column prop="component_name" label="组件名" width="320" />
      <el-table-column prop="component_type" label="组件类型" width="120" />
      <el-table-column prop="component_desc" label="组件描述" width="220" />
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
import request from "@/utils/request";

export default {
  name: "DataMean",
  data() {
    return {
      search: "",
      tableData: [],
      loading: true,
      currentPage: 1,
      pageSize: 10,
      total: 0,
    }
    
  },
  created() {
    this.loading = true
    this.load()
  },
  methods: {
    load() {
      this.loading = true
      request.get("/api/components", {
        params: {
          "page_num": this.currentPage,
          "page_size": this.pageSize,
          "search": this.search
        }
      }).then(res => {
        console.log(res)
        if (res.status_code === 0) {
          this.tableData = res.components
          this.total = res.total
        } else {
          this.$message({
            type: "error",
            message: res.status_msg
          })
        }
      })

      this.loading = false
    },
    updateComponent() {
      request.post("/api/components").then(res => {
        if (res.status_code === 0) {
          this.$message({
            type: "success",
            message: "更新中，请5-10分钟后刷新"
          })
        } else {
          this.$message({
            type: "error",
            message: res.status_msg
          })
        }
      })
    },
    handleSizeChange(pageSize) {  // 改变每页的大小
      this.pageSize = pageSize;
      this.load();
    },
    handleCurrentChange(pageNum) {  //改变当前页码
      this.currentPage = pageNum;
      this.load()
    },

  }
}
</script>

<style scoped>

</style>