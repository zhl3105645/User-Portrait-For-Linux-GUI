<template>
  <div style="height: 50px; line-height: 50px; border-bottom: 1px solid #ccc; display: flex">
    <div style="width: 200px; padding-left: 30px; font-weight: bold; color: dodgerblue">用户行为画像系统</div>
    <div style="flex: 1"></div>
    <div style="width: 100px">
      <el-dropdown>
        <span class="el-dropdown-link">
           {{ account.account_name }}
          <i class="el-icon-arrow-down el-icon--right"></i>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="logout">退出系统</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script>
import request from "@/utils/request";

export default {
  name: "HeaderNew",
  props: ['account'],
  data() {
    return {

    }
  },
  created() {
  },
  methods: {
    logout() {
      sessionStorage.removeItem("account")
      this.$router.push("/login")
    }
  }
}
</script>

<style scoped>

</style>