<template>
  <div>
    <!--    头部-->
    <Header :account="account"/>

    <!--    主体-->
    <div style="display: flex">
      <!--      侧边栏-->
      <Aside/>
      <!--      内容区域-->
      <router-view style="flex: 1" @userInfo="refreshAccount"/>
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
      let accountStr = sessionStorage.getItem("account")
      if (isEmptyStr(accountStr)) {
        // 账号信息
        request.get("/api/account").then(res => {
          console.log(res)
          if (res.status_code === 0) {
            console.log("set session account")
            let accountObj = res.account
            accountStr = JSON.stringify(accountObj)
            sessionStorage.setItem("account",accountStr)
          } else {
            console.log("/api/account code != 0, code=",res.status_code)
          }
        })
      }
      
      console.log("accountStr=", accountStr)
      let accountObj = JSON.parse(accountStr)
      this.account = accountObj
    }
  }
}

function isEmptyStr(s) {
  return s === undefined || s == null || s === '';

}

</script>

<style scoped>

</style>