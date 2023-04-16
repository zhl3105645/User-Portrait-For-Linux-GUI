<template>
  <div style="margin: 10px 0">
    <el-input v-model="search" placeholder="请输入关键字" style="width: 20%" clearable></el-input>
    <el-button type="primary" style="margin-left: 5px" @click="load">查询</el-button>
    <el-button type="primary" style="margin-left: 5px" @click="updateBasicBehaviorData">更新基础行为数据</el-button>
    
    
    <el-table 
      v-loading="loading"
      :span-method="arraySpanMethod"
      :data="tableData" 
      style="width: 100%"
    >
      <el-table-column prop="user_name" label="user_name" width="120" />
      <el-table-column prop="begin_time" label="begin_time" width="180" />
      <el-table-column prop="use_time" label="use_time" width="120" />
      <el-table-column prop="mouse_click_cnt" label="mouse_click_cnt" width="130" />
      <el-table-column prop="mouse_move_cnt" label="mouse_move_cnt" width="130" />
      <el-table-column prop="mouse_move_dis" label="mouse_move_dis" width="130" />
      <el-table-column prop="mouse_wheel_cnt" label="mouse_wheel_cnt" width="130" />
      <el-table-column prop="key_click_cnt" label="key_click_cnt" width="120" />
      <el-table-column prop="key_click_speed" label="key_click_speed" width="120" />
      <el-table-column prop="shortcut_cnt" label="shortcut_cnt" width="120" />
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
  name: "BasicBehaviorData",
  data() {
    return {
      search: '',
      loading: false,
      tableData: [],
      currentPage: 1,
      pageSize: 10,
      total: 0,
      mergeIndexMap: {}
    }
  },
  created() {
    this.load()
    this.mergeIndexMap = new Map()
  },
  methods: {
    load() {
      this.loading = true
      request.get("/api/basic_behaviors", {
        params: {
          "page_num": this.currentPage,
          "page_size": this.pageSize,
          "search": this.search
        }
      }).then(res => {
        console.log(res)
        if (res.status_code === 0) {
          this.merge(res.basic_behaviors)
          this.tableData = res.basic_behaviors
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
    merge(data) {
      this.mergeIndexMap = new Map()
      let curIdx = 0
      for (var idx = 1; idx < data.length; idx++) {
        if (data[idx].user_id === data[idx-1].user_id) {
          continue
        } else {
          this.mergeIndexMap.set(curIdx, idx-curIdx)
          curIdx = idx
        }
      }

      this.mergeIndexMap.set(curIdx, data.length - curIdx)
      console.log(this.mergeIndexMap)
    },
    arraySpanMethod({row,column,rowIndex,columnIndex}) {
      if (columnIndex === 0) {
        let rows = this.mergeIndexMap.get(rowIndex)
        return [rows, 1]
      }
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
      request.post("/api/gene_basic_behavior").then(res => {
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
    }
  }
}
</script>

<style scoped>

</style>