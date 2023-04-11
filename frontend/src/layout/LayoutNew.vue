<template>
  <div>
    <!--    头部-->
    <Header :account="account"/>

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
import Header from "@/components/new/HeaderNew";
import Aside from "@/components/new/AsideNew";
import request from "@/utils/request";

export default {
  name: "LayoutNew",
  components: {
    Header,
    Aside
  },
  data() {
    return {
      account: {}
    }
  },
  created() {
    this.refreshAccount();
  },
  methods: {
    refreshAccount() {
      let account = sessionStorage.getItem("account");
      if (account === null) {
        return
      }
      request.get("/api/account").then(res => {
        this.account = res.account;
        console.log(this.account)
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