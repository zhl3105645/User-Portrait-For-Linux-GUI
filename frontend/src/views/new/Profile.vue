<template>
  <div class="container">
    <el-select v-model="user_id" placeholder="请选择用户">
      <el-option v-for="user in users" :key="user.user_id" :label="user.user_name" :value="user.user_id"></el-option>
    </el-select>
    <el-button type="primary" style="margin-left: 5px" @click="load_profile()">生成画像</el-button>

    <!-- <div ref="profile_chart" style="width:950px;height:600px;margin:auto;"/> -->

    <!-- <div ref="behavior_duration_chart" style="width:950px;height:600px;margin:auto;"/> -->
    <div class="box top-left" v-if="this.group_labels.length >= 1">
      <el-card class="box-card">
        <template #header>
          <div class="card-header">
            <span>{{this.group_labels[0].parent_label_name}}</span>
          </div>
        </template>
        <div v-for="item in this.group_labels[0].labels" :key="item.label_id" class="text item">{{item.label_name + ": " + item.label_value }}</div>
      </el-card>
    </div>
    <div class="box_large top-right" v-if="this.group_labels.length >= 2">
      <el-card class="box-card">
        <template #header>
          <div class="card-header">
            <span>{{this.group_labels[1].parent_label_name}}</span>
          </div>
        </template>
        <div v-for="item in this.group_labels[1].labels" :key="item.label_id" class="text item">{{item.label_name + ": " + item.label_value }}</div>
      </el-card>
    </div>
    <div class="circle" ref="behavior_duration_chart"></div>
    <div class="box bottom-left" v-if="this.group_labels.length >= 3">
      <el-card class="box-card">
        <template #header>
          <div class="card-header">
            <span>{{this.group_labels[2].parent_label_name}}</span>
          </div>
        </template>
        <div v-for="item in this.group_labels[2].labels" :key="item.label_id" class="text item">{{item.label_name + ": " + item.label_value }}</div>
      </el-card>
    </div>
    <div class="box bottom-right" v-if="this.group_labels.length >= 4">
      <el-card class="box-card">
        <template #header>
          <div class="card-header">
            <span>{{this.group_labels[3].parent_label_name}}</span>
          </div>
        </template>
        <div v-for="item in this.group_labels[3].labels" :key="item.label_id" class="text item">{{item.label_name + ": " + item.label_value }}</div>
      </el-card>
    </div>
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
      tree_label: null,
      profile_chart: null,
      behavior_duration_chart: null,
      radars: null,
      group_labels: [],
    }
    
  },
  created() {
    this.load_user()
  },
  mounted() {
    // this.profile_chart = echarts.init(this.$refs.profile_chart)
    this.behavior_duration_chart = echarts.init(this.$refs.behavior_duration_chart,null, {renderer: 'svg'})
  },
  methods: {
    load_profile(){
      if (this.user_id == null) {
        return 
      }
      request.get("/api/profile/"+this.user_id).then(res => {
        console.log(res)
        if (res.status_code === 0) {
          this.tree_label = res.label
          this.radars = res.radars
          this.group_labels = res.group_labels
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
      let indicator = new Array()
      let appDuration = new Array()
      let userDuration = new Array()
      for (let i = 0; i < this.radars.length; i++) {
        let radar = this.radars[i]
        indicator.push({
          name: radar.name,
          max: radar.max,
        })
        appDuration.push(radar.ave)
        userDuration.push(radar.cur)
      }

      let option2 = {
        // title: {
        //   text: '行为时长图'
        // },
        // toolbox: {
        //     show: true,
        //     feature: {
        //         mark: {show: true},
        //         saveAsImage: {show: true},
        //     }
        // },
        legend: {
          top: "10%",
          data: ['应用平均时长', '用户平均时长']
        },
        radar: {
          radius: 99,
          indicator: indicator,
          name: {
            textStyle: {
              color: '#333',
              fontSize: 14
            }
          },
          splitLine: {
            lineStyle: {
              color: '#999'
            }
          },
          splitArea: {
            areaStyle: {
              color: ['rgba(250,250,250,0.3)', 'rgba(200,200,200,0.3)']
            }
          },
          axisLine: {
            lineStyle: {
              color: '#999'
            }
          },
        },
        series: [
          {
            name: 'Budget vs spending',
            type: 'radar',
            lineStyle: {
              width: 2
            },
            data: [
              {
                value: appDuration,
                name: "应用平均时长",
              },
              {
                value: userDuration,
                name: "用户平均时长",
              }
            ]
          }
        ]
      }
      this.behavior_duration_chart.setOption(option2)
    }
  }
}
</script>

<style scoped>
.container {
  position: relative;
  width: 100%;
  height: 100vh;
}
.box {
  position: absolute;
  width: 200px;
  height: 200px;
  /* background-color: lightblue; */
}
.box_large {
  position: absolute;
  width: 220px;
  height: 200px;
}
.top-left {
  top: 50px;
  left: 300px;
}
.top-right {
  top: 50px;
  right: 280px;
}
.bottom-left {
  bottom: 120px;
  left: 300px;
}
.bottom-right {
  bottom: 120px;
  right: 300px;
}
.circle {
  position: absolute;
  top: calc(50% - 220px);
  left: calc(50% - 220px);
  width: 440px;
  height: 440px;
  /* background-color: rgb(245, 242, 242); */
}

/*  */
.card-header {
  font-size: 16px;
  font-weight: bold;
  padding: 5px;
  border-bottom: 1px solid #ebeef5;
}

.text.item {
  padding: 5px;
  border-bottom: 1px solid #ebeef5;
  font-size: 14px;
  color: #606266;
}

.box-card {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  border-radius: 4px;
  overflow: hidden;
}

</style>