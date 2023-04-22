<template>
  <div style="margin: 10px 0">
    <el-select v-model="user_id" placeholder="请选择用户">
      <el-option v-for="user in users" :key="user.user_id" :label="user.name" :value="user.user_id"></el-option>
    </el-select>
    <el-button type="primary" style="margin-left: 5px" @click="query_profile">生成画像</el-button>
  </div>

</template>

<script>
import request from "@/utils/request";


export default {
  name: "Profile",
  data() {
    return {
      user_id: null,
      users: []
    }
    
  },
  created() {
    this.load()
  },
  methods: {
    load() {
      request.get("/api/all_user").then(res => {
        console.log(res)
        if (res.status_code === 0) {
          this.users = res.users
        } else {
          this.$message({
            type: "error",
            message: res.status_msg
          })
        }
      })
    }
  }
}
</script>

<style scoped>

</style>