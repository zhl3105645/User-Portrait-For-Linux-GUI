<template>
  <div style="margin: 10px 0">
    <el-select v-model="user_id" placeholder="请选择用户">
      <el-option v-for="user in users" :key="user.user_id" :label="user.user_name" :value="user.user_id"></el-option>
    </el-select>
    <el-button type="primary" style="margin-left: 5px" @click="load_profile()">生成画像</el-button>
  
 
    <div ref="profile_chart" style="width:950px;height:600px;margin:auto;"/>

  
  
  </div>
  


</template>

<script>
import request from "@/utils/request";
import * as echarts from 'echarts'

export default {
  name: "Profile",
  data() {
    return {
      user_id: null,
      users: [],
      profile: null,
      profile_chart: null,
    }
    
  },
  created() {
    this.load_user()
  },
  mounted() {
    this.profile_chart = echarts.init(this.$refs.profile_chart)
  },
  methods: {
    load_profile(){
      if (this.user_id == null) {
        return 
      }
      request.get("/api/profile/"+this.user_id).then(res => {
        console.log(res)
        if (res.status_code === 0) {
          this.profile = res.label
          this.set_chart()
        } else {
          this.$message({
            type: "error",
            message: res.status_msg
          })
        }
      })
    },
    load_user() {
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
    },
    set_chart() {
      let option = {
        tooltip: {
          trigger: 'item',
          triggerOn: 'mousemove'
        },
        series: [
          {
            type: 'tree',
            data: [this.profile],
            top: '1%',
            left: '7%',
            bottom: '1%',
            right: '20%',
            symbolSize: 7,
            label: {
              position: 'left',
              verticalAlign: 'middle',
              align: 'right',
              fontSize: 14
            },
            leaves: {
              label: {
                position: 'right',
                verticalAlign: 'middle',
                align: 'left'
              }
            },
            emphasis: {
              focus: 'descendant'
            },
            expandAndCollapse: true,
            animationDuration: 550,
            animationDurationUpdate: 750
          }
        ]
      }
      this.profile_chart.setOption(option)
    }
  }
}
</script>

<style scoped>

</style>