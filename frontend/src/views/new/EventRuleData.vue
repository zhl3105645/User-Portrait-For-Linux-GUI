<template>
  <div style="margin: 10px 0">
    <el-input v-model="search" placeholder="请输入关键字" style="width: 20%" clearable></el-input>
    <el-button type="primary" style="margin-left: 5px" @click="load">查询</el-button>
    <el-button type="primary" style="margin-left: 5px" @click="updateRuleData">更新规则数据</el-button>
    
    
    <el-table 
      v-loading="loading"
      :span-method="arraySpanMethod"
      :data="tableData" 
      style="width: 100%"
    >
      <el-table-column fixed prop="user_name" label="用户名" width="120" />
      <el-table-column prop="begin_time" label="开始时间" width="180" />
      <el-table-column prop="use_time" label="使用时长" width="180" />

      <el-table-column label="行为规则数据" :show-overflow-tooltip="true" width="500">
        <template #default="scope">
            <el-tag
              v-for="tag in scope.row.behavior_rule_data.rule_elements"
              :key="tag.rule_id"
              class="mx-1"
              :disable-transitions="false"
            >
              {{ tag.rule_desc }}
            </el-tag>
          </template>
      </el-table-column>
      <el-table-column label="操作" width="200" >
        <template #default="scope">
          <el-button link type="primary" size="small" @click="handleViewBehavior(scope.row.behavior_rule_data)">
            查看行为图
          </el-button>
        </template>
      </el-table-column>
    </el-table>


    <el-dialog v-model="dialogBehaviorViewVisible" title="行为图">
      <div ref="view_chart" style="width: 750px; height: 400px;"></div>
    </el-dialog>

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
import * as echarts from 'echarts';

export default {
  name: "EventRuleData",
  data() {
    return {
      search: '',
      loading: false,
      tableData: [],
      currentPage: 1,
      pageSize: 10,
      total: 0,
      mergeIndexMap: {},
      dialogBehaviorViewVisible: false,
      view_chart: null,
    }
  },
  created() {
    this.load()
    this.mergeIndexMap = new Map()
  },
  mounted() {

  },
  methods: {
    load() {
      this.loading = true
      request.get("/api/rule_data", {
        params: {
          "page_num": this.currentPage,
          "page_size": this.pageSize,
          "search": this.search
        }
      }).then(res => {
        console.log(res)
        if (res.status_code === 0) {
          this.merge(res.rule_data)
          this.tableData = res.rule_data
          this.total = res.total
        } else {
          this.$message({
            type: "error",
            message: res.status_msg
          })
        }
      })

      this.loading = false
    },
    merge(data) {
      this.mergeIndexMap = new Map()
      let curIdx = 0
      for (var idx = 1; idx < data.length; idx++) {
        if (data[idx].user_id === data[idx-1].user_id) {
          continue
        } else {
          this.mergeIndexMap.set(curIdx, idx-curIdx)
          curIdx = idx
        }
      }

      this.mergeIndexMap.set(curIdx, data.length - curIdx)
      console.log(this.mergeIndexMap)
    },
    arraySpanMethod({row,column,rowIndex,columnIndex}) {
      if (columnIndex === 0) {
        let rows = this.mergeIndexMap.get(rowIndex)
        return [rows, 1]
      }
    },
    handleSizeChange(pageSize) {  // 改变每页的大小
      this.pageSize = pageSize;
      this.load();
    },
    handleCurrentChange(pageNum) {  //改变当前页码
      this.currentPage = pageNum;
      this.load()
    },
    updateRuleData() {
      request.post("/api/gene_rule").then(res => {
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
    handleViewBehavior(behavior_rule_data) {
      if (behavior_rule_data == null || behavior_rule_data.rule_elements == null) {
        return 
      }
      let ruleMap = new Map()
      behavior_rule_data.rule_elements.forEach((element, index) => {
        if (!ruleMap.has(element.rule_desc)) {
          ruleMap.set(element.rule_desc, true)
        }
      })
      let rules = new Array()
      let rule_desc_2_index = new Map()
      let idx = 0
      for (let key of ruleMap.keys()) {
        rules.push(key)
        rule_desc_2_index.set(key, idx)
        idx = idx + 1
      }

      let data = new Array()
      for (let i = 0; i < behavior_rule_data.rule_elements.length - 1; i++) {
        let element = behavior_rule_data.rule_elements[i]
        let next_element = behavior_rule_data.rule_elements[i+1]
        let index = rule_desc_2_index.get(element.rule_desc)
        let begin_time = element.timestamp 
        let end_time = next_element.timestamp
        let duration = end_time - begin_time
        let d = {
          value: [index, begin_time, end_time, duration]
        }
        data.push(d)
      }


      let option = {
        dataZoom: [
          {
            type: 'slider',
            filterMode: 'weakFilter',
            showDataShadow: false,
            top: 400,
            labelFormatter: ''
          },
          {
            type: 'inside',
            filterMode: 'weakFilter'
          }
        ],
        grid: {
          height: 300
        },
        xAxis: {
          type: 'value',
          //min: startTime,
          scale: true,
          axisLabel: {
            formatter: function (val) {
              return timestampToTime(val)
            },
          }
        },
        yAxis: {
          data: rules
        },
        series: [
          {
            type: 'custom',
            renderItem: renderItem,
            itemStyle: {
              opacity: 0.8
            },
            encode: {
              x: [1, 2],
              y: 0
            },
            data: data
          }
        ]
      };
      this.dialogBehaviorViewVisible = true
      this.view_chart = echarts.init(this.$refs.view_chart);
      console.log("my option=", option)
      this.view_chart.setOption(option)
    }
  }  
}

function renderItem(params, api) {
  var categoryIndex = api.value(0);
  var start = api.coord([api.value(1), categoryIndex]);
  var end = api.coord([api.value(2), categoryIndex]);
  var height = api.size([0, 1])[1] * 0.6;
  var rectShape = echarts.graphic.clipRectByRect(
    {
      x: start[0],
      y: start[1] - height / 2,
      width: end[0] - start[0],
      height: height
    },
    {
      x: params.coordSys.x,
      y: params.coordSys.y,
      width: params.coordSys.width,
      height: params.coordSys.height
    }
  );
  return (
    rectShape && {
      type: 'rect',
      transition: ['shape'],
      shape: rectShape,
      style: api.style()
    }
  );
}

function timestampToTime(timestamp) {
        var date = new Date(timestamp);
        var Y = date.getFullYear() + '-';
        var M = (date.getMonth()+1 < 10 ? '0'+(date.getMonth()+1) : date.getMonth()+1) + '-';
        var D = (date.getDate() < 10 ? '0'+ date.getDate() : date.getDate()) + ' ';
        var h = (date.getHours() < 10 ? '0'+ date.getHours() : date.getHours()) + ':';
        var m = (date.getMinutes() < 10 ? '0'+ date.getMinutes() : date.getMinutes()) + ':';
        var s = (date.getSeconds() < 10 ? '0'+ date.getSeconds() : date.getSeconds());
        return M+D+h+m+s;
}
</script>

<style scoped>

</style>