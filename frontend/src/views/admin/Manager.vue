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
          prop="accountName"
          label="账户名"
          sortable
      >
      </el-table-column>
      <el-table-column
          prop="password"
          label="密码">
      </el-table-column>
      <el-table-column
          prop="name"
          label="姓名">
      </el-table-column>
      <el-table-column
          prop="age"
          label="年龄">
      </el-table-column>
      <el-table-column
          prop="sex"
          label="性别">
      </el-table-column>
      <el-table-column
          prop="email"
          label="邮件">
      </el-table-column>
      <el-table-column
          prop="telephone"
          label="电话号码">
      </el-table-column>
      <el-table-column label="操作" width="240">
        <template #default="scope">
          <el-button size="mini" @click="handleEdit(scope.row)">授权</el-button>
          <el-popconfirm title="确定删除吗？" @confirm="handleDelete(scope.row.accountName)">
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
          <el-form-item label="账户名" v-if="granting === false">
            <el-input v-model="form.accountName" style="width: 80%"></el-input>
          </el-form-item>
          <el-form-item label="密码" v-if="granting === false">
            <el-input v-model="form.password" style="width: 80%"></el-input>
          </el-form-item>
          <el-form-item label="姓名" v-if="granting === false">
            <el-input v-model="form.name" style="width: 80%"></el-input>
          </el-form-item>
          <el-form-item label="年龄" v-if="granting === false">
            <el-input v-model="form.age" style="width: 80%"></el-input>
          </el-form-item>
          <el-form-item label="性别" v-if="granting === false">
            <el-input v-model="form.sex" style="width: 80%"></el-input>
          </el-form-item>
          <el-form-item label="邮件" v-if="granting === false">
            <el-input v-model="form.email" style="width: 80%"></el-input>
          </el-form-item>
          <el-form-item label="电话号码" v-if="granting === false">
            <el-input v-model="form.telephone" style="width: 80%"></el-input>
          </el-form-item>
          <el-form-item label="管理员管理权限">
            <el-select v-model="form.managerRole">
              <el-option label="拥有" value="true"></el-option>
              <el-option label="未拥有" value="false"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="员工管理权限">
            <el-select v-model="form.staffRole">
              <el-option label="拥有" value="true"></el-option>
              <el-option label="未拥有" value="false"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="部门管理权限">
            <el-select v-model="form.departmentRole">
              <el-option label="拥有" value="true"></el-option>
              <el-option label="未拥有" value="false"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="用户管理权限">
            <el-select v-model="form.userRole">
              <el-option label="拥有" value="true"></el-option>
              <el-option label="未拥有" value="false"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="车辆管理权限">
            <el-select v-model="form.vehicleRole">
              <el-option label="拥有" value="true"></el-option>
              <el-option label="未拥有" value="false"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="保险管理权限">
            <el-select v-model="form.insuranceRole">
              <el-option label="拥有" value="true"></el-option>
              <el-option label="未拥有" value="false"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="仓库管理权限">
            <el-select v-model="form.storehouseRole">
              <el-option label="拥有" value="true"></el-option>
              <el-option label="未拥有" value="false"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="日志管理权限">
            <el-select v-model="form.logRole">
              <el-option label="拥有" value="true"></el-option>
              <el-option label="未拥有" value="false"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="订单管理权限">
            <el-select v-model="form.orderRole">
              <el-option label="拥有" value="true"></el-option>
              <el-option label="未拥有" value="false"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="价目表管理权限">
            <el-select v-model="form.priceRole">
              <el-option label="拥有" value="true"></el-option>
              <el-option label="未拥有" value="false"></el-option>
            </el-select>
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
  name: "Manager",
  data() {
    return {
      granting: false,
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
      request.get("/admin/manager/count").then(res => {
        this.total = res.data;
      })
    },
    load() {
      this.loading = true;
      request.get("/admin/manager", {
        params: {
          pageNum: this.currentPage,
          pageSize: this.pageSize,
          search: this.search
        }
      }).then(res => {
        this.loading = false;
        this.tableData = res.data;
        console.log("数组", this.tableData)
        this.updateTotal();
      })
    },
    add() {
      this.dialogVisible = true;
      this.form = {};
    },
    save() {
      if (this.granting) {    // 授权
        request.put("/admin/manager/grant", this.form).then(res => {
          console.log(res);
          if (res.code === '0') {
            this.$message({
              type: "success",
              message: "授权成功"
            })
          } else {
            this.$message({
              type: "error",
              message: res.msg
            })
          }
          this.load(); // 刷新表格数据
          this.dialogVisible = false; //关闭弹窗
          this.granting = false;
        })
      } else {
        request.post("/admin/manager", this.form).then(res => {
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
      this.granting = true;
    },
    handleDelete(id) {
      console.log(id)
      request.delete("/admin/manager/" + id).then(res => {
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