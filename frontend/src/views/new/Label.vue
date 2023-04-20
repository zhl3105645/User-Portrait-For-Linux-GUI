<template>
  <div style="margin: 10px 0">
    <el-input v-model="search" placeholder="请输入关键字" style="width: 20%" clearable></el-input>
    <el-button type="primary" style="margin-left: 5px" @click="load">查询</el-button>
    <el-button type="primary" style="margin-left: 5px" @click="this.dialogAddLabelVisible = true">添加标签</el-button>
  
    <el-dialog v-model="dialogAddLabelVisible" title="添加标签">
      <el-form :model="addLabelForm">
        <el-form-item label="标签名称">
          <el-input v-model="addLabelForm.label_name"/>
        </el-form-item>
        <el-form-item label="标签数据源">
          <el-select v-model="addLabelForm.model_id" placeholder="">
            <el-option v-for="item in model_source_options" :key="item.value" :label="item.label" :value="item.value"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="标签转换规则">
          <div v-for="(rule, index) in convertRules" :key="index">
            <el-button>模型数据</el-button>
            <el-select v-model="rule.operator" placeholder="运算符" style="width: 20%">
            <el-option
              v-for="(value, key) in operators"
              :key="key"
              :label="value"
              :value="value"
            />
            </el-select>
            <el-input v-model="rule.x_value" style="width: 10%"> </el-input>
            <el-button>标签数据</el-button>
            <el-input v-model="rule.y_value" style="width: 10%"> </el-input>
            <el-button>标签描述</el-button>
            <el-input v-model="rule.y_desc" style="width: 20%"> </el-input>
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
    
    <div v-for="(label,index) in label_data" :key="index">
      <el-card>
        <div >
          <span>{{label.label_name}}</span>
          <el-button type="primary" style="margin-left: 5px" @click="handleGeneLabelData(label.label_id)">生成数据</el-button>
          <el-button type="primary" style="margin-left: 5px" @click="handleDeleteLabel(label.label_id)">删除标签</el-button>
        </div>
        <div :id="label.label_id" style="width:400px;height:300px;margin:auto;"/>
      </el-card>
    </div>

    <div style="margin: 10px 0">
      <el-pagination
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="currentPage"
          :page-sizes="[5, 10, 20]"
          :page-size="pageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total">
      </el-pagination>
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
      search: "",
      loading: true,
      currentPage: 1,
      pageSize: 10,
      total: 0,
      label_data: [],
      addLabelForm: {},
      dialogAddLabelVisible: false,
      model_source_options: [],
      convertRules: [],
      operators: {
        1: ">=",
        2: ">",
        3: "=",
        4: "<",
        5: "<=",
      }
    }
  },
  created() {
    this.load()
    this.loadSource()
  },
  mounted() {
    this.load()
  },
  methods: {
    load() {
      this.loading = true
      request.get("/api/label", {
        params: {
          "page_num": this.currentPage,
          "page_size": this.pageSize,
          "search": this.search
        }
      }).then(res => {
        console.log(res)
        if (res.status_code === 0) {
          this.total = res.total
          this.label_data = res.labels
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
    handleSizeChange(pageSize) {  // 改变每页的大小
      this.pageSize = pageSize;
      this.load();
    },
    handleCurrentChange(pageNum) {  //改变当前页码
      this.currentPage = pageNum;
      this.load()
    },
    addLabel() {
      this.addLabelForm.convert_rules = this.convertRules
      console.log(this.addLabelForm)
      request.post("/api/label", this.addLabelForm).then(res => {
        if (res.status_code === 0) {
          this.$message({
            type: "success",
            message: "添加成功"
          })
          this.load()
          this.addLabelForm = {}
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
    loadSource() {
      request.get("/api/data_sources").then(res => {
        console.log(res)
        if (res.status_code === 0) {
          let r = res.data_source
          for (let i = 0; i < r.length; i ++ ) {
            if (r[i].value == 6) {
              this.model_source_options = r[i].children
              break
            }
          }

        } else {
          this.$message({
            type: "error",
            message: res.status_msg
          })
        }
      })
    },
    addConvertRule() {
      this.convertRules.push({
        operator: null,
        x_value: 0.0,
        y_value: 0,
        y_desc: "",
      })
    },
    initCharts() {
      this.label_data.forEach((label, index) => {
        let chart = echarts.init(document.getElementById(label.label_id.toString()));
        if (label.option != null && label.option.series != null) {
          let res = {
            toolbox: label.option.toolbox,
            tooltip: label.option.tooltip,
            xAxis: label.option.xAxis,
            yAxis: label.option.yAxis,
          }
          let series = new Array()
          for (let i = 0; i< label.option.series.length; i++) {
            let oldSer = label.option.series[i]
            let ser = {
              type: oldSer.type,
              smooth: oldSer.smooth,
              yAxisIndex: oldSer.yAxisIndex,
              data: new Array(),
            }
            if (oldSer.type == "pie") {
              for (let idx = 0; idx < oldSer.data.length; idx++) {
                ser.data.push(JSON.parse(oldSer.data[idx]))
              }
            } else {
              ser.data = oldSer.data
            }
            series.push(ser)
          }
          res.series = series

          console.log(res)
          chart.setOption(res)
        }
      })
    }

  }
}
</script>

<style scoped>

</style>