<template>
  <div style="width: 100%; height: 100vh; overflow: hidden"> <!--  :style="bg" 加背景图片-->
    <div style="width: 400px; margin: 100px auto">
      <div style="font-size: 30px; text-align: center; padding: 30px 0">用户行为画像系统登录</div>
      <el-form ref="form" :model="form" :rules="rules">
        <el-form-item prop="app_id">
          <el-select v-model="form.app_id" placeholder="请选择应用">
            <el-option v-for="item in options" :key="item.app_id" :label="item.app_name" :value="item.app_id"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item prop="account_name">
          <el-input prefix-icon="el-icon-user-solid" v-model="form.account_name" placeholder="请输入账号"></el-input>
        </el-form-item>
        <el-form-item prop="account_pwd">
          <el-input prefix-icon="el-icon-lock" v-model="form.account_pwd" show-password placeholder="请输入密码"></el-input>
        </el-form-item>
        <el-form-item>
          <div style="display: flex">
            <el-input prefix-icon="el-icon-key" v-model="form.validCode" style="width: 50%;"
                      placeholder="请输入验证码"></el-input>
            <ValidCode @input="createValidCode"/>
          </div>
        </el-form-item>
        <el-form-item>
          <el-button style="width: 100%" type="primary" @click="login">登 录</el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="text" @click="$router.push('/register')">创建应用 >></el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
import request from "@/utils/request";
import ValidCode from "@/components/ValidCode";

export default {
  name: "Login",
  components: {
    ValidCode
  },
  data() {
    return {
      options: [],
      form: {},
      rules: {
        account_name: [
          {required: true, message: '请输入用户名', trigger: 'blur'},
        ],
        account_pwd: [
          {required: true, message: '请输入密码', trigger: 'blur'},
        ],
        app_id: [
          {required: true, message: '请选择应用', trigger: 'blur'},
        ],
      },
      validCode: ''
      // 加背景图片
      // bg: {
      //   backgroundImage: "url(" + require("@/assets/bg.jpg") + ")",
      //   backgroundRepeat: "no-repeat",
      //   backgroundSize: "100% 100%"
      // }
    }
  },
  created() {
    this.getAppList();
    sessionStorage.removeItem("token")
  },
  methods: {
    // 接收验证码组件提交的 4位验证码
    createValidCode(data) {
      this.validCode = data
    },
    // 获取应用列表 
    getAppList() {
      request.get("/applist").then(res => {
          console.log(res)
          if (res.status_code === 0) {
            this.options = res.apps
          } else {
             this.$message({
                  type: "error",
                  message: res.status_msg
                })
          }
      }
      )
    },
    login() {
      this.$refs['form'].validate((valid) => {
        if (valid) {
          if (!this.form.validCode) {
            this.$message.error("请填写验证码")
            return
          }
          if (this.form.validCode.toLowerCase() !== this.validCode.toLowerCase()) {
            this.$message.error("验证码错误")
            return
          }
          console.log("login...ing")
          request.post("/login", this.form).then(res => {
            console.log(res)
            if (res.status_code === 200) {
              this.$message({
                  type: "success",
                  message: "登录成功"
              })
              sessionStorage.setItem("token", res.token) // 缓存token

              // 账号信息
              request.get("/api/account").then(res => {
                console.log(res)
                if (res.status_code === 0) {
                  console.log("set session account")
                  let accountObj = res.account
                  let accountStr = JSON.stringify(accountObj)
                  sessionStorage.setItem("account",accountStr)
                } else {
                  console.log("/api/account code != 0, code=",res.status_code)
                }
              })
    
              this.$router.push("/front/profile") // 登录界面跳转
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
