<template>
  <div style="margin: 10px 0">
    <el-input v-model="search" placeholder="请输入关键字" style="width: 20%" clearable></el-input>
    <el-button type="primary" style="margin-left: 5px" @click="load">查询</el-button>
    <el-button type="primary" style="margin-left: 5px" @click="this.dialogAddModelVisible = true">添加模型</el-button>
  
    <el-dialog v-model="dialogAddModelVisible" title="添加模型">
      <el-form :model="addModelForm">
        <el-form-item label="模型名称">
          <el-input v-model="addModelForm.model_name"/>
        </el-form-item>
        <el-form-item label="模型类型">
          <el-radio-group v-model="addModelForm.model_type" class="ml-4">
            <el-radio label="1" size="large">数学统计</el-radio>
            <el-radio label="2" size="large">机器学习</el-radio>
          </el-radio-group>
        </el-form-item>
        <!-- 统计 -->
        <el-form-item v-if="addModelForm.model_type==1" label="模型数据源">
          <el-cascader :options="sta_data_source_options" @change="handleDataSourceChange" />
        </el-form-item>
        <el-form-item v-if="addModelForm.model_type==1" label="统计数据">
          <el-select v-model="addModelForm.calculate_type" placeholder="">
            <el-option v-for="item in data_calculate_options" :key="item.type" :label="item.desc" :value="item.type"></el-option>
          </el-select>
        </el-form-item>
        <!-- 机器学习 -->
        <el-form-item v-if="addModelForm.model_type==2" label="模型功能">
          <el-radio-group v-model="addModelForm.model_feature" class="ml-4">
            <el-radio label="1" size="large">标签</el-radio>
            <el-radio label="2" size="large">预测</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="addModelForm.model_type==2" label="请求参数">
          <el-select v-model="http_type" placeholder="请求类型">
            <el-option v-for="item in http_types" :key="item.type" :label="item.desc" :value="item.type"></el-option>
          </el-select>
          <el-input v-model="http_addr" placeholder="请输入服务地址" style="width: 50%" clearable></el-input>
          <el-button @click="addRow">添加参数</el-button>
          <el-table :data="body_params" >
            <el-table-column prop="name" label="表单参数名">
              <template #default="scope">
                <el-input v-model="scope.row.name" />
              </template>
            </el-table-column>
            <el-table-column prop="dataSource" label="表单参数数据源">
              <template #default="scope">
                <el-select v-model="scope.row.model_id">
                  <el-option
                    v-for="source in learn_data_source_options"
                    :key="source.value"
                    :label="source.label"
                    :value="source.value"
                  />
                </el-select>
              </template>
            </el-table-column>
            <el-table-column label="">
              <template #default="scope">
                <el-button @click="removeRow(scope.$index)">删除参数</el-button>
              </template>
            </el-table-column>
          </el-table>
          
        </el-form-item>
        <el-form-item v-if="addModelForm.model_type==2" label="响应参数">
          <el-select v-model="http_resp_data_type">
            <el-option
              v-for="data in http_resp_data_types"
              :key="data.type"
              :label="data.desc"
              :value="data.type"
            />
          </el-select>
          <el-input v-model="http_resp_name" placeholder="请输入参数名" style="width: 50%" clearable></el-input>
        </el-form-item>

      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogAddModelVisible = false">取消</el-button>
          <el-button type="primary" @click="addModel">
            确定添加
          </el-button>
        </span>
      </template>
    </el-dialog>
    
    <div v-for="(model,index) in model_data" :key="index">
      <el-card>
        <div >
          <span>{{model.model_name}}</span>
          <el-button type="primary" style="margin-left: 5px" @click="handleGeneModelData(model.model_id)">生成数据</el-button>
          <el-button type="primary" style="margin-left: 5px" @click="handleDeleteModel(model.model_id)">删除模型</el-button>
        </div>
        <div :id="model.model_id" style="width:400px;height:300px;margin:auto;"/>
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
  name: "StatisticsModel",
  data() {
    return {
      search: "",
      loading: true,
      currentPage: 1,
      pageSize: 10,
      total: 0,
      model_data: [],
      addModelForm: {
      },
      dialogAddModelVisible: false,
      sta_data_source_options:[],
      learn_data_source_options:[],
      data_calculate_options: [],
      http_types: [
        {
          "type": 1,
          "desc": "POST"
        }
      ],
      http_type: 1,
      http_addr: "",
      http_resp_name: "",
      http_resp_data_type: 1,
      http_resp_data_types: [
        {
          "type": 1,
          "desc": "连续型"
        },
        {
          "type": 2,
          "desc": "枚举型"
        }
      ],
      body_params: [],
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
    addRow() {
      this.body_params.push({ name: "", model_id: this.learn_data_source_options[0] });
    },
    removeRow(index) {
      this.body_params.splice(index, 1);
    },
    load() {
      this.loading = true
      request.get("/api/model", {
        params: {
          "page_num": this.currentPage,
          "page_size": this.pageSize,
          "search": this.search
        }
      }).then(res => {
        console.log(res)
        if (res.status_code === 0) {
          this.total = res.total
          this.model_data = res.models
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
    loadSource() {
      request.get("/api/data_sources").then(res => {
        console.log(res)
        if (res.status_code === 0) {
          let r = res.data_source
          for (let i = 0; i < r.length; i ++ ) {
            if (r[i].value == 6) {
              this.learn_data_source_options = r[i].children
            } else {
              this.sta_data_source_options.push(r[i])
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
    addModel() {
      if (this.addModelForm.model_type == 2) {
        let learning_form = {}
        learning_form.http_type = this.http_type
        learning_form.body_params = this.body_params
        learning_form.http_resp_name = this.http_resp_name
        learning_form.http_resp_data_type = this.http_resp_data_type
        learning_form.http_addr = this.http_addr
        this.addModelForm.learning_param = learning_form
      }
      console.log(this.addModelForm)
      request.post("/api/model", this.addModelForm).then(res => {
        if (res.status_code === 0) {
          this.$message({
            type: "success",
            message: "添加成功"
          })
          this.load()
          this.addModelForm = {}
          this.dialogAddModelVisible = false
        } else {
          this.$message({
            type: "error",
            message: res.status_msg
          })
        }
      })
    },
    handleDataSourceChange(value) {
      console.log(value)
      this.addModelForm.source_type = value[0]
      if (value[1] != null) {
        this.addModelForm.source_value = value[1]
      }

      if (value[0] == 1) {
        if (value[1] == 8) {
          this.data_calculate_options = [
            {
              "type": 2,
              "desc": "众数"
            },
          ]
        } else {
          this.data_calculate_options = [
            {
              "type": 1,
              "desc": "平均数"
            },
          ]
        }
      } else if (value[0] == 2) {
        this.data_calculate_options = [
          {
          "type": 1,
          "desc": "事件次数平均数"
          },
        ]
      } else if (value[0] == 3) {
        this.data_calculate_options = [
          {
          "type": 1,
          "desc": "行为时长平均数"
          },
        ]
      } else if (value[0] == 4) {
        this.data_calculate_options = [
          {
          "type": 3,
          "desc": "各类事件次数平均数"
          },
        ]
      } else if (value[0] == 5) {
        this.data_calculate_options = [
          {
          "type": 4,
          "desc": "各类行为时长平均数"
          },
        ]
      }
    },
    handleDeleteModel(model_id) {
      this.$confirm('确定删除？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(res => {
        request.delete("/api/model/" + model_id).then(res => {
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
    handleGeneModelData(model_id) {
      request.post("/api/model/" + model_id).then(res => {
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
    initCharts() {
      this.model_data.forEach((model, index) => {
        let chart = echarts.init(document.getElementById(model.model_id.toString()));
        if (model.option != null && model.option.series != null) {
          let res = {
            toolbox: model.option.toolbox,
            tooltip: model.option.tooltip,
            xAxis: model.option.xAxis,
            yAxis: model.option.yAxis,
          }
          let series = new Array()
          for (let i = 0; i< model.option.series.length; i++) {
            let oldSer = model.option.series[i]
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