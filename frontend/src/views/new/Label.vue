<template>
    <div class="container">
      <div class="left-half">
        <div style="margin: 10px 0">
          <el-button type="primary" style="margin-left: 5px" @click="this.dialogAddLabelVisible = true">添加标签</el-button>
        
          <el-dialog v-model="dialogAddLabelVisible" title="添加标签" width="700px">
            <el-form :model="addLabelForm" label-width="100px" label-position="left">
              <el-form-item label="标签名称">
                <el-input v-model="addLabelForm.label_name"/>
              </el-form-item>
              <el-form-item label="父标签" >
                <el-select v-model="addLabelForm.parent_label_id" placeholder="若为一级标签可不选">
                  <el-option v-for="item in label_options" :key="item.label_id" :label="item.label_name" :value="item.label_id"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="叶子标签">
                <el-radio-group class="ml-4" v-model="store_data">
                  <el-radio label="1" size="large">是</el-radio>
                  <el-radio label="2" size="large">否</el-radio>
                </el-radio-group>
              </el-form-item>
              <el-form-item label="数据类型"  v-if="store_data==1">
                <el-radio-group v-model="addLabelForm.data_type" class="ml-4">
                  <el-radio label="1" size="large">枚举</el-radio>
                  <el-radio label="2" size="large">不可枚举</el-radio>
                </el-radio-group>
              </el-form-item>
              <el-form-item label="语义化表述" v-if="addLabelForm.data_type==1">
                <div v-for="(rule, index) in convertRules" :key="index">
                  <el-button>标签数据</el-button>
                  <el-input v-model="rule.data" style="width: 20%"> </el-input>
                  <el-button>标签描述</el-button>
                  <el-input v-model="rule.desc" style="width: 20%"> </el-input>
                </div>
                <el-button @click="addConvertRule">添加规则</el-button>
              </el-form-item>


            </el-form>
            <template #footer>
              <span class="dialog-footer">
                <el-button @click="dialogAddLabelVisible = false">取消</el-button>
                <el-button type="primary" @click="addLabel">
                  确定添加
                </el-button>
              </span>
            </template>
          </el-dialog>

          <div ref="tree_label" style="width:600px;height:450px;margin:auto;">


          </div>
        </div>
      </div>
      <div class="right-half">
        <div style="margin: 10px 0">
          <el-select v-model="cur_label_id" placeholder="选择标签">
            <el-option v-for="item in label_options" :key="item.label_id" :label="item.label_name" :value="item.label_id"></el-option>
          </el-select>  
          <el-button type="primary" style="margin-left: 5px" @click="load_single_label">查询</el-button>
          <el-button type="primary" style="margin-left: 5px" @click="update_label">修改标签</el-button>
          <el-button type="primary" style="margin-left: 5px" @click="handleGeneLabelData">更新标签数据</el-button>
          <el-button type="danger" style="margin-left: 5px" @click="handleDeleteLabel">删除标签</el-button>

          <div ref="single_label_chart" style="width:600px;height:450px;margin:auto;">
        </div>
      </div>
    </div>
    </div>
    

</template>

<script>
import request from "@/utils/request";
import * as echarts from 'echarts'

export default {
  name: "Label",
  data() {
    return {
      addLabelForm: {},
      dialogAddLabelVisible: false,
      label_options: [],
      convertRules: [],
      store_data: false,
      tree_label_chart: null,
      tree_label_data: null,
      cur_label_id: null,
      single_label_chart: null,
    }
  },
  created() {
    this.load()
  },
  mounted() {
    this.tree_label_chart = echarts.init(this.$refs.tree_label)
    this.single_label_chart = echarts.init(this.$refs.single_label_chart)
  },
  methods: {
    load() {
      this.loadLabelOption()
      request.get("/api/tree_label").then(res => {
        console.log(res)
        if (res.status_code === 0) {
          this.tree_label_data = res.data
          this.tree_label_chart.clear()
          this.set_tree_labels()
        } else {
          this.$message({
            type: "error",
            message: res.status_msg
          })
        }
      })
    },
    addLabel() {
      this.addLabelForm.convert_rules = this.convertRules
      request.post("/api/label", this.addLabelForm).then(res => {
        if (res.status_code === 0) {
          this.$message({
            type: "success",
            message: "添加成功"
          })
          this.load()
          this.addLabelForm = {}
          this.convertRules = []
          this.dialogAddLabelVisible = false
        } else {
          this.$message({
            type: "error",
            message: res.status_msg
          })
        }
      })
    },
    handleGeneLabelData() {
      if (this.cur_label_id == null) {
        return 
      }
      request.post("/api/label/" + this.cur_label_id).then(res => {
        if (res.status_code === 0) {
          this.$message({
            type: "success",
            message: "更新中，请5-10分钟后刷新"
          })
        } else {
          this.$message({
            type: "error",
            message: res.status_msg
          })
        }
      })
    },
    handleDeleteLabel() {
      if (this.cur_label_id == null) {
        return 
      }
      this.$confirm('确定删除？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(res => {
        request.delete("/api/label/" + this.cur_label_id).then(res => {
            console.log(res)
            if (res.status_code === 0) {
              this.$message({
                type: "success",
                message: "删除成功"
              })
              this.load()
            } else {
              this.$message({
                type: "error",
                message: res.status_msg
              })
            }
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        });          
      })
    },
    loadLabelOption() {
      request.get("/api/labels").then(res => {
        console.log(res)
        if (res.status_code === 0) {
          this.label_options = res.labels
        } else {
          this.$message({
            type: "error",
            message: res.status_msg
          })
        }
      })
    },
    addConvertRule() {
      this.convertRules.push({})
    },
    set_tree_labels() {      
      let option = {
        tooltip: {
          trigger: 'item',
          triggerOn: 'mousemove'
        },
        series: [
          {
            type: 'tree',
            data: [this.tree_label_data],
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
      //console.log("option",option)
      this.tree_label_chart.setOption(option)
    },
    update_label(){

    },
    load_single_label() {
      request.get("/api/label/" + this.cur_label_id).then(res => {
        console.log(res)
        if (res.status_code === 0) {
          this.single_label_chart.clear()
          if (res.chart_type == 1) {
            this.set_pie_label(res.pie_label)
          } else if (res.chart_type == 2) {
            this.set_bar_label(res.bar_label)
          }
        } else {
          this.$message({
            type: "error",
            message: res.status_msg
          })
        }
      })
    },
    set_pie_label(pie_label) {
      let option1 = {
        tooltip: {
            trigger: 'item'
        },
        series: [
            {
                type: 'pie',
                radius: ['40%', '70%'],
                center: ['50%', '60%'],
                avoidLabelOverlap: false,
                itemStyle: {
                    borderRadius: 10,
                    borderColor: '#fff',
                    borderWidth: 2
                },
                label: {
                  show: true,
                  formatter(param) {
                    // correct the percentage
                    return param.name + ' (' + param.percent + '%)';
                  }
                },
                data: pie_label.data,
            }
        ]
    }
    this.single_label_chart.setOption(option1)
    },
    set_bar_label(bar_label) {
      let option = {
        tooltip: {
            trigger: 'axis',
            axisPointer: {
            type: 'shadow'
            }
        },
        xAxis: [
            {
                name: "小时",
                type: 'category',
                data: bar_label.x_names,
                axisTick: {
                    alignWithLabel: true
                },
            }
        ],
        yAxis: [
            {
                name: '用户数',
                type: 'value',
                nameLocation: 'center',
                nameGap: 30,
            }
        ],
        grid: {
            left: '5%',
            right: '15%',
            bottom: '3%',
            containLabel: true
        },
        series: [
            {
                name: 'Direct',
                type: 'bar',
                barWidth: '60%',
                data: bar_label.data,
                showBackground: true,
                backgroundStyle: {
                    color: 'rgba(180, 180, 180, 0.2)'
                },
                itemStyle: {
                    normal: {
                    //这里是重点
                    color: function(params) {
                        var colorList = ['#c23531','#2f4554', '#61a0a8', '#d48265', '#91c7ae','#749f83', '#ca8622'];
                        // var colorList = ['#c23531','#2f4554', '#61a0a8'];
                        // 自动循环已经有的颜色
                        return colorList[params.dataIndex % colorList.length];
                    }
                    }
                }
            }
        ]
      }
      this.single_label_chart.setOption(option)
    }
  }
}
</script>

<style scoped>

.container {
  position: relative;
  width: 100%;
  height: 100vh;
  /* background-color: #773148; */
}
.left-half {
  /* background-color: #a0868f; */
  position: absolute;
  top: 0px;
  left: 0px;
  width: 50%;
  height: 100vh;
}
.right-half {
  /* background-color: rgb(131, 5, 26); */
  position: absolute;
  top: 0px;
  right: 0px;
  width: 50%;
  height: 100vh;
}
</style>