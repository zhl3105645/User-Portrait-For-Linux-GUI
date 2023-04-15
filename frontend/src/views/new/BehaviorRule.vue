<template>
  <div style="margin: 10px 0">
    <el-input v-model="search" placeholder="请输入关键字" style="width: 20%" clearable></el-input>
    <el-button type="primary" style="margin-left: 5px" @click="load">查询</el-button>
    <el-button type="primary" style="margin-left: 5px" @click="dialogAddRuleVisible = true">新增规则</el-button>

    <el-dialog v-model="dialogAddRuleVisible" title="添加规则">
      <el-form :model="addRuleform">
        <el-form-item label="规则描述">
          <el-input v-model="addRuleform.rule_desc"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogAddRuleVisible = false">取消</el-button>
          <el-button type="primary" @click="addRule">
            确定添加
          </el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog v-model="dialogUpdateRuleVisible" title="更新规则">
      <el-form :model="updateRuleform">
        <el-form-item label="规则描述">
          <el-input v-model="updateRuleform.rule_desc" :placeholder="curRow.rule_desc"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogAddRuleVisible = false">取消</el-button>
          <el-button type="primary" @click="updateRule">
            确定更新
          </el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog v-model="dialogAddElementVisible" title="添加元素">
      <!-- 标签栏 -->
      <el-tag
        v-for="tag in curEventRules"
        :key="tag.rule_id"
        class="mx-1"
        closable
        :disable-transitions="false"
        @close="handleDeleteTag(tag)"
      >
        {{ tag.rule_desc }}
      </el-tag>

      <el-select v-model="elementTag" filterable placeholder="请选择事件规则">
        <el-option v-for="item in eventRules" :key="item.rule_id" :label="item.rule_desc" :value="item.rule_id"></el-option>
      </el-select>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="addElementTag">
            添加元素
          </el-button>
          <el-button type="primary" @click="addElement">
            完成
          </el-button>
        </span>
      </template>
    </el-dialog>

    <el-table 
      v-loading="loading"
      :span-method="arraySpanMethod"
      :data="tableData" 
      style="width: 100%"
    >
      <el-table-column label="规则">
        <!-- <el-table-column fixed prop="rule_id" label="规则ID" width="80" /> -->
        <el-table-column prop="rule_desc" label="规则描述" width="100" />
        <el-table-column label="规则操作" width="300">
        <template #default="scope">
          <el-button size="small" @click="handleAddElement(scope.row)">新增元素</el-button>
          <el-button size="small" @click="handleUpdateRule(scope.row)">修改规则</el-button>
          <el-button size="small" type="danger" @click="handleDeleteRule(scope.row)">删除规则</el-button>
        </template>
      </el-table-column>
      </el-table-column>
      <el-table-column label="规则元素">
        <!-- <el-table-column prop="element_id" label="规则元素ID" width="100" /> -->
        <el-table-column prop="event_rules" label="元素" width="500">
          <template #default="scope">
            <el-tag
              v-for="tag in scope.row.event_rules"
              :key="tag.rule_id"
              class="mx-1"
              :disable-transitions="false"
            >
              {{ tag.rule_desc }}
            </el-tag>

          </template>
        </el-table-column>
        <el-table-column label="元素操作" width="300">
          <template #default="scope">
            <el-button size="small" @click="handleUpdateElement(scope.row)">
              更新元素
            </el-button>
            <el-button size="small" type="danger" @click="handleDeleteElement(scope.row)">
              删除元素
            </el-button>
          </template>
        </el-table-column>

      </el-table-column>
      
    </el-table>
    

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


export default {
  name: "BehaviorRule",
  data() {
    return {
      search: "",
      tableData: [],
      loading: true,
      currentPage: 1,
      pageSize: 10,
      total: 0,
      dialogAddRuleVisible: false,
      addRuleform: {},
      dialogUpdateRuleVisible: false,
      updateRuleform: {},
      curRow: {},
      dialogAddElementVisible: false,
      status: 0, // 1 添加 2 更新
      curEventRules: [],
      elementTag: null, // 选中标签
      mergeIndexMap: {}
    }
  },
  created() {
    this.loading = true
    this.mergeIndexMap = new Map()
    this.load()
    this.loadEventRule()
  },
  methods: {
    load(){
      this.loading = true
      request.get("/api/elements", {
        params: {
          "rule_type": 2,
          "page_num": this.currentPage,
          "page_size": this.pageSize,
          "search": this.search
        }
      }).then(res => {
        console.log(res)
        if (res.status_code === 0) {
          this.merge(res.behavior_elements)
          this.tableData = res.behavior_elements
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
    loadEventRule() {
      request.get("/api/rules", {
        params: {
          "rule_type": 1,
        }
      }).then(res => {
        console.log(res)
        if (res.status_code === 0) {
          this.eventRules = res.event_rules
        }
      })
    },
    merge(data) {
      this.mergeIndexMap = new Map()
      let curIdx = 0
      for (var idx = 1; idx < data.length; idx++) {
        if (data[idx].rule_id === data[idx-1].rule_id) {
          continue
        } else {
          this.mergeIndexMap.set(curIdx, idx-curIdx)
          curIdx = idx
        }
      }

      this.mergeIndexMap.set(curIdx, data.length - curIdx)
      console.log(this.mergeIndexMap)
    },
    handleSizeChange(pageSize) {  // 改变每页的大小
      this.pageSize = pageSize;
      this.load();
    },
    handleCurrentChange(pageNum) {  //改变当前页码
      this.currentPage = pageNum;
      this.load()
    },
    addRule() {
      this.addRuleform.rule_type = 2
      request.post("/api/rule", this.addRuleform).then(res => {
        console.log(res)
        if (res.status_code === 0) {
           this.$message({
            type: "success",
            message: "添加成功"
          })
          this.addRuleform = {}
          this.dialogAddRuleVisible = false
          this.load()
        } else {
           this.$message({
            type: "error",
            message: res.status_msg
          })
        }
      })
    },
    handleUpdateRule(row) {
      this.curRow = row
      this.dialogUpdateRuleVisible = true
    },
    handleDeleteRule(row) {
      this.curRow = row
      this.$confirm('确定删除？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(res => {
        console.log(res)
        request.delete("/api/rule/" + this.curRow.rule_id).then(res => {
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
    updateRule() {
      request.put("/api/rule/" + this.curRow.rule_id, this.updateRuleform).then(res => {
        console.log(res)
        if (res.status_code === 0) {
           this.$message({
            type: "success",
            message: "更新成功"
          })
          this.dialogUpdateRuleVisible = false
          this.updateRuleform = {}
          this.load()
        } else {
           this.$message({
            type: "error",
            message: res.status_msg
          })
        }
      })
    },
    handleAddElement(row) {
      this.curRow = row
      this.curEventRules = []
      this.dialogAddElementVisible = true
      this.status = 1
    },
    addElementTag() {
      if (this.elementTag == null) {
        return
      }
      let desc = ""
      for (let i = 0; i < this.eventRules.length; i++) {
        if (this.eventRules[i].rule_id === this.elementTag) {
          desc = this.eventRules[i].rule_desc
          break
        }
      }
      this.curEventRules.push({
        "rule_id": this.elementTag,
        "rule_desc": desc
      })
    },
    handleDeleteTag(tag) {
      this.curEventRules.splice(this.curEventRules.indexOf(tag), 1)
    },
    addElement() {
      if (this.status === 1) {
        let eventRuleIds = new Array()
        for (let i = 0; i < this.curEventRules.length; i++) {
          eventRuleIds.push(this.curEventRules[i].rule_id)
        }
        let form = {
          "rule_id": this.curRow.rule_id,
          "event_rule_ids": eventRuleIds
        }
        console.log("form=", form)
        request.post("/api/element", form).then(res => {
          console.log(res)
          if (res.status_code === 0) {
            this.$message({
              type: "success",
              message: "添加成功"
            })
            this.dialogAddElementVisible = false
            this.load()
          } else {
            this.$message({
              type: "error",
              message: res.status_msg
            })
          }
        })
      } else if (this.status ===2 ) {
        let elementId = this.curRow.element_id
        let eventRuleIds = new Array()
        for (let i = 0; i < this.curEventRules.length; i++) {
          eventRuleIds.push(this.curEventRules[i].rule_id)
        }
        let form = {
          "event_rule_ids": eventRuleIds
        }
        console.log("form=", form)
        request.put("/api/element/"+elementId, form).then(res => {
          console.log(res)
          if (res.status_code === 0) {
            this.$message({
              type: "success",
              message: "更新成功"
            })
            this.dialogAddElementVisible = false
            this.load()
          } else {
            this.$message({
              type: "error",
              message: res.status_msg
            })
          }
        })
      }
      
    },
    handleUpdateElement(row) {
      this.curRow = row
      this.curEventRules = row.event_rules
      this.dialogAddElementVisible = true
      this.status = 2
    },
    handleDeleteElement(row) {
      this.$confirm('确定删除？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(res => {
        console.log(res)
        request.delete("/api/element/" + row.element_id).then(res => {
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
    arraySpanMethod({row,column,rowIndex,columnIndex}) {
      if (columnIndex === 0 || columnIndex === 1) {
        let rows = this.mergeIndexMap.get(rowIndex)
        return [rows, 1]
      }
    }
  }
}
</script>

<style scoped>

</style>