<template>
  <div style="padding: 10px">
    <!--    功能区域-->
    <div style="margin: 10px 0">
      <el-button type="primary" @click="addOut">新增运出订单</el-button>
      <el-button type="primary" @click="addIn">新增运进订单</el-button>
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
        @selection-change="handleSelectionChange">
      <el-table-column
          type="selection"
          width="40">
      </el-table-column>
      <el-table-column
          prop="id"
          label="ID"
          sortable>
      </el-table-column>
      <el-table-column
          prop="username"
          label="用户名"
          v-if="false">
      </el-table-column>
      <el-table-column
          prop="goodsName"
          label="货物名">
      </el-table-column>
      <el-table-column
          prop="deliverAmount"
          label="运输数量">
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
      <el-table-column
          prop="createTime"
          label="创建时间">
      </el-table-column>
      <el-table-column
          prop="payState"
          label="支付状态">
      </el-table-column>
      <el-table-column
          prop="completeState"
          label="完成状态">
      </el-table-column>
      <el-table-column label="操作" width="240">
        <template #default="scope">
          <el-popconfirm title="确定支付吗？" @confirm="payOrder(scope.row.id)">
            <template #reference>
              <el-button size="mini" type="">支付</el-button>
            </template>
          </el-popconfirm>
          <el-popconfirm title="确定订单已完成吗？" @confirm="completeOrder(scope.row.id)">
            <template #reference>
              <el-button size="mini" type="">完成</el-button>
            </template>
          </el-popconfirm>
          <el-popconfirm title="确定删除吗？" @confirm="handleDelete(scope.row.id)">
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
          <!--运进订单-->
          <el-form-item label="货物名" v-if="this.orderType">
            <el-input v-model="form.name" style="width: 80%"></el-input>
          </el-form-item>
          <el-form-item label="数量" v-if="this.orderType">
            <el-input v-model="form.amount" style="width: 80%"></el-input>
          </el-form-item>
          <el-form-item label="货物类型" v-if="this.orderType">
            <el-input v-model="form.type" style="width: 80%"></el-input>
          </el-form-item>
          <el-form-item label="仓库" v-if="this.orderType">
            <el-select v-model="form.storehouseId" style="width: 80%">
              <el-option v-for="item in storehouseData" :value="item.id" :label="item.name"/>
            </el-select>
          </el-form-item>
          <el-form-item label="起点" v-if="this.orderType">
            <el-input v-model="form.fromAddress" style="width: 80%"></el-input>
          </el-form-item>

          <!--运出订单-->
          <el-form-item label="货物" v-if="!this.orderType">
            <el-select v-model="form.goodsId" style="width: 80%">
              <el-option v-for="item in goodsData" :value="item.id" :label="item.storehouseName + ': ' + item.name"/>
            </el-select>
          </el-form-item>
          <el-form-item label="运输数量" v-if="!this.orderType">
            <el-input v-model="form.deliverAmount" style="width: 80%"></el-input>
          </el-form-item>
          <el-form-item label="终点" v-if="!this.orderType">
            <el-input v-model="form.toAddress" style="width: 80%"></el-input>
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
  name: "Order",
  data() {
    return {
      orderType: false,
      loading: true,
      form: {},
      dialogVisible: false,
      search: '',
      currentPage: 1,
      pageSize: 10,
      total: 0,
      tableData: [],
      storehouseData: [],
      goodsData: []
    }
  },
  created() {
    this.load();
  },
  methods: {
    handleSelectionChange() {

    },
    updateTotal() {
      request.get("/user/order/count").then(res => {
        this.total = res.data;
      })
    },
    load() {
      this.loading = true;
      request.get("/user/order", {
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
      //获取仓库列表信息
      request.get("/user/storehouse").then(res => {
        this.storehouseData = res.data;
      })
      //获取货物信息
      request.get("/user/allgoods").then(res => {
        this.goodsData = res.data;
      })
    },
    addIn() {
      this.dialogVisible = true;
      this.form = {};
      this.orderType = true;
    },
    addOut() {
      this.dialogVisible = true;
      this.form = {};
      this.orderType = false;
    },
    save() {
      if (this.form.id) {   //update
        request.put("/user/order/" + this.form.id, this.form).then(res => {
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
        })
      } else {        //  add
        request.post("/user/order", this.form).then(res => {
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
    payOrder(id) {
      request.put("/user/order/pay", {
        params: {
          id: id
        }
      }).then(res => {
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
        this.load();
      })
    },
    completeOrder(id) {
      request.put("/user/order/complete", {
        params: {
          id: id
        }
      }).then(res => {
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
        this.load();
      })
    },
    handleEdit(row) {
      this.form = JSON.parse(JSON.stringify(row));
      this.dialogVisible = true;
    },
    handleDelete(plateNumber) {
      console.log(plateNumber)
      request.delete("/user/order/" + plateNumber).then(res => {
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