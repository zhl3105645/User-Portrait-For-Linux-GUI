<template>
  <div>
    <!--    头部-->
    <Header :user="user"/>

    <!--    主体-->
    <div style="display: flex">
      <!--      侧边栏-->
      <Aside/>
      <!--      内容区域-->
      <router-view style="flex: 1" @userInfo="refreshUser"/>
    </div>
  </div>
</template>

<script>
import Header from "@/components/user/HeaderUser";
import Aside from "@/components/user/AsideUser";
import request from "@/utils/request";

export default {
  name: "LayoutUser",
  components: {
    Header,
    Aside
  },
  data() {
    return {
      user: {}
    }
  },
  created() {
    this.refreshUser();
  },
  methods: {
    refreshUser() {
      let userJson = sessionStorage.getItem("user");
      if (isEmptyStr(userJson)) {
        return
      }
      request.get("/user/profile").then(res => {
        this.user = res.data;
      })
    }
  }
}

function isEmptyStr(s) {
  return s === undefined || s == null || s === '';
}
</script>

<style scoped>

</style>