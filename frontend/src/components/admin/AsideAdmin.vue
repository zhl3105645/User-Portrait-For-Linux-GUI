<template>
  <div>
    <el-menu
        style="width: 200px; min-height: calc(100vh - 50px)"
        :default-active="path"
        router
    >
      <el-menu-item index="/admin/home">主页</el-menu-item>
      <el-menu-item index="/admin/department" v-if="admin.departmentRole">部门管理</el-menu-item>
      <el-menu-item index="/admin/goods" v-if="admin.staffRole">货物管理</el-menu-item>
      <el-menu-item index="/admin/insurance" v-if="admin.insuranceRole">保险管理</el-menu-item>
      <el-menu-item index="/admin/manager" v-if="admin.managerRole">管理员管理</el-menu-item>
      <el-menu-item index="/admin/order" v-if="admin.orderRole">订单管理</el-menu-item>
      <el-menu-item index="/admin/price" v-if="admin.priceRole">价目表管理</el-menu-item>
      <el-menu-item index="/admin/staff" v-if="admin.staffRole">员工管理</el-menu-item>
      <el-menu-item index="/admin/storehouse" v-if="admin.storehouseRole">仓库管理</el-menu-item>
      <el-menu-item index="/admin/user" v-if="admin.userRole">用户管理</el-menu-item>
      <el-menu-item index="/admin/vehicle" v-if="admin.vehicleRole">车辆管理</el-menu-item>
    </el-menu>
  </div>
</template>

<script>
import request from "@/utils/request";

export default {
  name: "AsideAdmin",
  data() {
    return {
      admin: {},
      path: this.$route.path  //设置默认高亮菜单
    }
  },
  methods:{
    // toHome()  {
    //   this.$router.push({
    //     path: "/admin"
    //   })
    // }
},
  created() {
    let adminStr = sessionStorage.getItem("admin") || {};
    this.admin = JSON.parse(adminStr);

    //请求服务端，确认当前登录用户的合法信息
    request.get("/admin/profile").then(res => {
      if (res.code === '0') {
        this.admin = res.data
      }
    })
  }
}
</script>

<style scoped>

</style>