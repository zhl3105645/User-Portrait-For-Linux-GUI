<template>
  <div style="width: 100%; height: 100vh; overflow: hidden">
    <div style="width: 400px; margin: 100px auto">
      <div style="font-size: 30px; text-align: center; padding: 30px 0">欢迎注册</div>
      <el-form ref="form" :model="form" size="normal" :rules="rules">
        <el-form-item prop="app_name">
          <el-input prefix-icon="el-icon-user-solid" v-model="form.app_name" placeholder="应用名"></el-input>
        </el-form-item>
        <el-form-item prop="account_name">
          <el-input prefix-icon="el-icon-user-solid" v-model="form.account_name" placeholder="账号名"></el-input>
        </el-form-item>
        <el-form-item prop="account_pwd">
          <el-input prefix-icon="el-icon-lock" v-model="form.account_pwd" placeholder="账号密码" show-password></el-input>
        </el-form-item>
        <el-form-item prop="confirm">
          <el-input prefix-icon="el-icon-lock" v-model="form.confirm" placeholder="确认密码" show-password></el-input>
        </el-form-item>
        <el-form-item>
          <el-button style="width: 100%" type="primary" @click="register">注册</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
import request from "@/utils/request";

export default {
  name: "Register",
  data() {
    return {
      form: {},
      rules: {
        app_name: [
          {required: true, message: '请输入应用名', trigger: 'blur'},
        ],
        account_name: [
          {required: true, message: '请输入账号名', trigger: 'blur'},
        ],
        account_pwd: [
          {required: true, message: '请输入密码', trigger: 'blur'},
        ],
        confirm: [
          {required: true, message: '请确认密码', trigger: 'blur'},
        ],
      }
    }
  },
  methods: {
    register() {

      if (this.form.account_pwd !== this.form.confirm) {
        this.$message({
          type: "error",
          message: '2次密码输入不一致！'
        })
        return
      }

      this.$refs['form'].validate((valid) => {
        if (valid) {
          request.post("/register", this.form).then(res => {
            console.log(res)
            if (res.status_code === 0) {
              this.$message({
                type: "success",
                message: "注册成功"
              })
              this.$router.push("/login")  //登录成功之后进行页面的跳转，跳转到主页
            } else {
              this.$message({
                type: "error",
                message: res.status_msg
              })
            }
          })
        }
      })
    }
  }
}
</script>

<style scoped>

</style>
