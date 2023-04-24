<template>
  <div style="margin: 10px 0">
    <el-input v-model="search" placeholder="请输入关键字" style="width: 20%" clearable></el-input>
    <el-button type="primary" style="margin-left: 5px" @click="this.dialogAddLabelVisible = true">添加标签</el-button>
  
    <el-dialog v-model="dialogAddLabelVisible" title="添加标签">
      <el-form :model="addLabelForm">
        <el-form-item label="标签名称">
          <el-input v-model="addLabelForm.label_name"/>
        </el-form-item>
        <el-form-item label="父标签" >
          <el-select v-model="addLabelForm.parent_label_id" placeholder="若为一级标签可不选">
            <el-option v-for="item in label_options" :key="item.label_id" :label="item.label_name" :value="item.label_id"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="是否存储数据">
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
        <el-form-item label="标签语义化表述" v-if="addLabelForm.data_type==1">
          <div v-for="(rule, index) in convertRules" :key="index">
            <el-button>标签数据</el-button>
            <el-input v-model="rule.data" style="width: 10%"> </el-input>
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
    
    <!-- <div v-for="(label,index) in label_data" :key="index">
      <el-card>
        <div >
          <span>{{label.label_name}}</span>
          <el-button type="primary" style="margin-left: 5px" @click="handleGeneLabelData(label.label_id)">生成数据</el-button>
          <el-button type="primary" style="margin-left: 5px" @click="handleDeleteLabel(label.label_id)">删除标签</el-button>
        </div>
        <div :id="label.label_id" style="width:400px;height:300px;margin:auto;"/>
      </el-card>
    </div> -->
  </div>
</template>

<script>
import request from "@/utils/request";
import * as echarts from 'echarts'

export default {
  name: "Label",
  data() {
    return {
      loading: true,
      addLabelForm: {},
      dialogAddLabelVisible: false,
      label_options: [],
      convertRules: [],
      store_data: false,
      tree_label_chart: null,
      tree_label_data: null,
    }
  },
  created() {
    //this.load()
  },
  mounted() {
    this.load()
  },
  methods: {
    load() {
      this.loading = true
      this.loadLabelOption()
      request.get("/api/tree_label").then(res => {
        console.log(res)
        if (res.status_code === 0) {
          this.tree_label_data = res.data
          this.initCharts()
        } else {
          this.$message({
            type: "error",
            message: res.status_msg
          })
        }
      })
      this.loading = false
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
    handleGeneLabelData(label_id) {
      request.post("/api/label/" + label_id).then(res => {
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
    handleDeleteLabel(label_id) {
      this.$confirm('确定删除？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(res => {
        request.delete("/api/label/" + label_id).then(res => {
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
    initCharts() {
      if (this.tree_label_chart == null) {
        console.log("初始化")
        this.tree_label_chart = echarts.init(this.$refs.tree_label)
      }
      
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
      console.log("option",option)
      this.tree_label_chart.setOption(option)
    }

  }
}
</script>

<style scoped>

</style>