<template>
  <div>
    <!--    头部-->
    <Header :admin="admin"/>

    <!--    主体-->
    <div style="display: flex">
      <!--      侧边栏-->
      <Aside/>
      <!--      内容区域-->
      <router-view style="flex: 1" @userInfo="refreshAdmin"/>
    </div>
  </div>
</template>

<script>
import Header from "@/components/admin/HeaderAdmin";
import Aside from "@/components/admin/AsideAdmin";
import request from "@/utils/request";

export default {
  name: "LayoutAdmin",
  components: {
    Header,
    Aside
  },
  data() {
    return {
      admin: {}
    }
  },
  created() {
    this.refreshAdmin();
  },
  methods: {
    refreshAdmin() {
      let adminJson = sessionStorage.getItem("admin");
      if (isEmptyStr(adminJson)) {
        return
      }
      request.get("/admin/profile").then(res => {
        this.admin = res.data;
        console.log(this.admin)
      }).catch(err => {
        console.log("出现错误")
        console.log(err)
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