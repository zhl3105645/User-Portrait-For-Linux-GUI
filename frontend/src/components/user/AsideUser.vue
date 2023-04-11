<template>
  <div>
    <el-menu
        style="width: 200px; min-height: calc(100vh - 50px)"
        :default-active="path"
        router
    >
      <el-menu-item index="/user/home">主页</el-menu-item>
      <el-menu-item index="/user/goods">货物管理</el-menu-item>
      <el-menu-item index="/user/order">订单管理</el-menu-item>
    </el-menu>
  </div>
</template>

<script>
import request from "@/utils/request";

export default {
  name: "AsideUser",
  data() {
    return {
      user: {},
      path: this.$route.path // 设置默认高亮菜单
    }
  },
  created() {
    let userStr = sessionStorage.getItem("user") || {};
    this.user = JSON.parse(userStr);

    //请求服务端吗，确认当前登录用户的合法信息
    request.get("/user/profile").then(res => {
      if (res.code === '0') {
        this.user = res.data;
      }
    })
  }
}
</script>

<style scoped>

</style>