<template>
  <div style="padding: 10px">
    <!--    功能区域-->
    <div style="margin: 10px 0">
      <el-button type="primary" @click="add">新增</el-button>
    </div>

    <!--    搜索区域-->
    <div style="margin: 10px 0">
      <el-input v-model="search" placeholder="请输入关键字" style="width: 20%" clearable></el-input>
      <el-button type="primary" style="margin-left: 5px" @click="load">查询</el-button>
    </div>
    <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
        @selection-change="handleSelectionChange"
    >
      <el-table-column
          type="selection"
          width="55">
      </el-table-column>
      <el-table-column
          prop="amount"
          label="数量">
      </el-table-column>
      <el-table-column
          prop="fromAddress"
          label="起点">
      </el-table-column>
      <el-table-column
          prop="toAddress"
          label="终点">
      </el-table-column>
      <el-table-column
          prop="money"
          label="价格">
      </el-table-column>
      <el-table-column label="操作" width="240">
        <template #default="scope">
          <el-button size="mini" @click="handleEdit(scope.row)">编辑</el-button>
          <el-popconfirm title="确定删除吗？" @confirm="handleDelete(scope.row.amount, scope.row.fromAddress, scope.row.toAddress)">
            <template #reference>
              <el-button size="mini" type="danger">删除</el-button>
            </template>
          </el-popconfirm>
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

      <el-dialog title="提示" v-model="dialogVisible" width="30%">
        <el-form :model="form" label-width="120px">
          <el-form-item label="数量" v-if="this.edit === false">
            <el-input v-model="form.amount" style="width: 80%"></el-input>
          </el-form-item>
          <el-form-item label="起点" v-if="this.edit === false">
            <el-input v-model="form.fromAddress" style="width: 80%"></el-input>
          </el-form-item>
          <el-form-item label="终点" v-if="this.edit === false">
            <el-input v-model="form.toAddress" style="width: 80%"></el-input>
          </el-form-item>
          <el-form-item label="价格">
            <el-input v-model="form.money" style="width: 80%"></el-input>
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="dialogVisible = false">取 消</el-button>
            <el-button type="primary" @click="save">确 定</el-button>
          </span>
        </template>
      </el-dialog>

    </div>
  </div>
</template>

<script>
import request from "@/utils/request";

export default {
  name: "Price",
  data() {
    return {
      edit: false,
      loading: true,
      form: {},
      dialogVisible: false,
      search: '',
      currentPage: 1,
      pageSize: 10,
      total: 0,
      tableData: []
    }
  },
  created() {
    this.load();
  },
  methods: {
    handleSelectionChange() {

    },
    updateTotal() {
      request.get("/admin/price/count").then(res => {
        this.total = res.data;
      })
    },
    load() {
      this.loading = true;
      request.get("/admin/price", {
        params: {
          pageNum: this.currentPage,
          pageSize: this.pageSize,
          search: this.search
        }
      }).then(res => {
        this.loading = false;
        this.tableData = res.data;
        this.updateTotal();
      })
    },
    add() {
      this.dialogVisible = true;
      this.form = {};
    },
    save() {
      if (this.edit) {   //update
        request.put("/admin/price/" + this.form.amount + "/" + this.form.fromAddress + "/" + this.form.toAddress, this.form).then(res => {
          console.log(res);
          if (res.code === '0') {
            this.$message({
              type: "success",
              message: "更新成功"
            })
          } else {
            this.$message({
              type: "error",
              message: res.msg
            })
          }
          this.load(); // 刷新表格数据
          this.dialogVisible = false; //关闭弹窗
          this.edit = false;
        })
      } else {        //  add
        request.post("/admin/price", this.form).then(res => {
          console.log(res)
          if (res.code === '0') {
            this.$message({
              type: "success",
              message: "新增成功"
            })
          } else {
            this.$message({
              type: "error",
              message: res.msg
            })
          }
          this.load() // 刷新表格的数据
          this.dialogVisible = false  // 关闭弹窗
        })
      }
    },
    handleEdit(row) {
      this.form = JSON.parse(JSON.stringify(row));
      this.dialogVisible = true;
      this.edit = true;
    },
    handleDelete(amount, fromAddress, toAddress) {
      console.log(id)
      request.delete("/admin/price/" + amount + "/" + fromAddress + "/" + toAddress).then(res => {
        if (res.code === '0') {
          this.$message({
            type: "success",
            message: "删除成功"
          })
        } else {
          this.$message({
            type: "error",
            message: res.msg
          })
        }
        this.load()  // 删除之后重新加载表格的数据
      })
    },
    handleSizeChange(pageSize) {  // 改变每页的大小
      this.pageSize = pageSize;
      this.load();
    },
    handleCurrentChange(pageNum) {  //改变当前页码
      this.currentPage = pageNum;
      this.load()
    }
  }
}
</script>

<style scoped>

</style>