<template>
  <div style="margin: 10px 0">
    <el-input v-model="search" placeholder="请输入关键字" style="width: 20%" clearable></el-input>
    <el-button type="primary" style="margin-left: 5px" @click="load">查询</el-button>

     <!-- Form -->
    <el-button type="primary" @click="dialogFormVisible = true">
      添加用户
    </el-button>
    <el-dialog v-model="dialogFormVisible" title="添加用户">
      <el-form :model="form">
        <el-form-item label="用户名">
          <el-input v-model="form.username"/>
        </el-form-item>
        <el-form-item label="性别">
          <el-radio-group v-model="form.user_gender" class="ml-4">
            <el-radio label="1" size="large">男</el-radio>
            <el-radio label="2" size="large">女</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="年龄">
          <el-input v-model="form.user_age"/>
        </el-form-item>
        <el-form-item label="职业">
          <el-input v-model="form.user_career"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogFormVisible = false">取消</el-button>
          <el-button type="primary" @click="addUser">
            确定添加
          </el-button>
        </span>
      </template>
    </el-dialog>

    <el-table 
      v-loading="loading"
      :data="tableData" 
      style="width: 100%"
    >
      <el-table-column fixed prop="user_name" label="用户名" width="150" />
      <el-table-column fixed prop="user_gender" label="性别" width="150" />
      <el-table-column fixed prop="user_age" label="年龄" width="150" />
      <el-table-column fixed prop="user_career" label="职业" width="150" />
      <el-table-column prop="record_num" label="使用记录数" width="150" />
      <el-table-column label="操作" width="300">
        <template #default="scope">
          <el-button link type="primary" size="small" @click="handleClickOpenDialog(scope.$index)"
            >导入数据</el-button
          >
          <el-button link type="danger" size="small" @click="handleDeleteUser(scope.row.user_id)"
            >删除用户</el-button
          >
          <el-dialog v-model="dialogUploadVisible" title="导入数据">
            <el-upload
              class="upload-demo"
              drag
              action="https://run.mocky.io/v3/9d059bf9-4660-45f2-925d-ce80ad6c4d15"
              multiple
              :auto-upload="false"
              v-model:file-list="fileList"
              :on-change="handleFileChange"
            >
              <el-icon class="el-icon--upload"><upload-filled /></el-icon>
              <div class="el-upload__text">
                Drop file here or <em>click to upload</em>
              </div>
            </el-upload>
            <template v-slot:footer >
              <el-button @click="dialogUploadVisible=false"> 取消</el-button>
              <el-button type="primary" @click="confirmUpload"> 上传</el-button>
            </template>
          </el-dialog>
          
        </template>
      </el-table-column>
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
  name: "UserRecord",
  data() {
    return {
      search: "",
      loading: true,
      dialogFormVisible: false,
      dialogUploadVisible: false,
      form: {},
      currentPage: 1,
      pageSize: 10,
      total: 0,
      tableData: [],
      fileList: [],
      rowIndex: 0,
    }
  },
  created() {
    this.load()
  },
  methods: {
    load() {
      this.loading = true
      request.get("/api/users", {
        params: {
          "page_num": this.currentPage,
          "page_size": this.pageSize,
          "search": this.search
        }
      }).then(res => {
        console.log(res)
        if (res.status_code === 0) {
          this.tableData = res.users
          this.total = res.total
        } else {
          this.$message({
            type: "error",
            message: res.status_msg
          })``
        }
      })

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
    addUser() {
      request.post("/api/user", this.form).then(res => {
        console.log(res)
        if (res.status_code === 0) {
          this.$message({
            type: "success",
            message: "添加成功"
          })
          this.load()
          this.dialogFormVisible = false;
        } else {
          this.$message({
            type: "error",
            message: res.status_msg
          })
        }
      })
    },
    handleDownload(index) {
      console.log("下载数据",index)
    },
    handleClickOpenDialog(index) {
      this.rowIndex = index
      this.dialogUploadVisible = true
    },
    handleDeleteUser(user_id) {
      this.$confirm('确定删除？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(res => {
        request.delete("/api/user/" + user_id).then(res => {
            console.log(res)
            if (res.status_code === 0) {
              this.$message({
                type: "success",
                message: "删除成功"
              })
              this.load()
            } else {
              this.$message({
                type: "error",
                message: res.status_msg
              })
            }
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        });          
      })
    },
    handleFileChange(file, fileList) { // 文件数量改变
      this.fileList = fileList
    },
    confirmUpload() { // 确认上传
      var param = new FormData()
      this.fileList.forEach((val, index) => {
        param.append("file", val.raw)
      })

      request.post("/api/user/upload/" + this.tableData[this.rowIndex].user_id, param).then(res => {
        if (res.status_code === 0) {
          this.$message({
            type: "success",
            message: "上传成功"
          })
          this.load()
          this.fileList = []
          this.dialogUploadVisible = false
        } else {
          this.$message({
            type: "error",
            message: "上传失败"
          })
        }
      })
    }
  }
}
</script>

<style scoped>

</style>