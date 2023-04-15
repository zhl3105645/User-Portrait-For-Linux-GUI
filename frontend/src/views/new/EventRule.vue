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

    <el-dialog v-model="dialogAddElementVisible" title="新增元素">
      <el-form :model="addElementForm">
        <el-form-item label="事件类型">
          <el-select v-model="addElementForm.event_type" placeholder="请选择事件类型">
            <el-option v-for="(name, index) in eventTypes" :key="index" :label="name" :value="index"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="鼠标点击类型">
          <el-select v-model="addElementForm.mouse_click_type" placeholder="请选择鼠标点击类型">
            <el-option v-for="(name, index) in mouseClickTypes" :key="index" :label="name" :value="index"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="鼠标点击键">
          <el-select v-model="addElementForm.mouse_click_button" placeholder="请选择鼠标点击键">
            <el-option v-for="(name, index) in mouseClickBtns" :key="index" :label="name" :value="index"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="键盘点击类型">
          <el-select v-model="addElementForm.key_click_type" placeholder="请选择键盘点击类型">
            <el-option v-for="(name, index) in keyClickTypes" :key="index" :label="name" :value="index"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="键盘点击值">
          <el-input v-model="addElementForm.key_value"/>
        </el-form-item>
        <el-form-item label="组件前缀">
          <el-input v-model="addElementForm.component_name_prefix"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogAddElementVisible = false">取消</el-button>
          <el-button type="primary" @click="addElement">
            确认添加
          </el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog v-model="dialogUpdateElementVisible" title="更新元素">
      <el-form :model="updateElementform">
        <el-form-item label="事件类型">
          <el-select v-model="updateElementform.event_type" :placeholder="eventTypes[curRow.event_type]">
            <el-option v-for="(name, index) in eventTypes" :key="index" :label="name" :value="index"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="鼠标点击类型">
          <el-select v-model="updateElementform.mouse_click_type" :placeholder="mouseClickTypes[curRow.mouse_click_type]">
            <el-option v-for="(name, index) in mouseClickTypes" :key="index" :label="name" :value="index"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="鼠标点击键">
          <el-select v-model="updateElementform.mouse_click_button" :placeholder="mouseClickBtns[curRow.mouse_click_button]">
            <el-option v-for="(name, index) in mouseClickBtns" :key="index" :label="name" :value="index"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="键盘点击类型">
          <el-select v-model="updateElementform.key_click_type" :placeholder="keyClickTypes[curRow.key_click_type]">
            <el-option v-for="(name, index) in keyClickTypes" :key="index" :label="name" :value="index"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="键盘点击值">
          <el-input v-model="updateElementform.key_value" :placeholder="curRow.key_value"/>
        </el-form-item>
        <el-form-item label="组件前缀">
          <el-input v-model="updateElementform.component_name_prefix" :placeholder="curRow.component_name_prefix"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogUpdateElementVisible = false">取消</el-button>
          <el-button type="primary" @click="updateElement">
            确认更新
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
        <!-- <el-table-column prop="element_id" label="规则元素ID" width="90" /> -->
        <el-table-column prop="event_type" label="事件类型" width="80">
          <template #default="scope">
            <div>
              <span style="margin-left: 10px">{{ eventTypes[scope.row.event_type] }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="mouse_click_type" label="鼠标点击类型" width="100" >
          <template #default="scope">
            <div>
              <span style="margin-left: 10px">{{ mouseClickTypes[scope.row.mouse_click_type] }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="mouse_click_button" label="鼠标点击键" width="90" >
          <template #default="scope">
            <div>
              <span style="margin-left: 10px">{{ mouseClickBtns[scope.row.mouse_click_button] }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="key_click_type" label="键盘点击类型" width="100">
          <template #default="scope">
            <div>
              <span style="margin-left: 10px">{{ keyClickTypes[scope.row.key_click_type] }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="key_value" label="键盘点击值" width="90" />
        <el-table-column prop="component_name_prefix" label="组件前缀" width="200" :show-overflow-tooltip="true"/>
        <el-table-column fixed="right" label="元素操作" width="200">
          <template #default="scope">
            <el-button size="small" @click="handleEditElement(scope.row)">编辑元素</el-button>
            <el-button size="small"  type="danger"  @click="handleDeleteElement(scope.row)">删除元素</el-button>
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
  name: "EventRule",
  data() {
    return {
      search: "",
      tableData: [],
      loading: true,
      currentPage: 1,
      pageSize: 10,
      total: 0,
      dialogAddRuleVisible: false,
      addRuleform: {
        "rule_type": 1
      },
      dialogUpdateRuleVisible: false,
      updateRuleform: {},
      dialogAddElementVisible: false,
      addElementForm: {},
      curRow: {},
      dialogUpdateElementVisible: false,
      updateElementform: {},
      eventTypes: {
        "3": "鼠标点击",
        "4": "鼠标移动",
        "5": "键盘点击",
        "6": "鼠标滚动",
        "7": "快捷键",
      },
      mouseClickTypes: {
        "1": "鼠标单击",
        "2": "鼠标双击",
      },
      mouseClickBtns: {
        "1": "鼠标左键",
        "2": "鼠标右键",
      },
      keyClickTypes:{
        "1": "单键",
        "2": "组合键",
      },
      mergeIndexMap: {}
    }
    
  },
  created() {
    this.loading = true
    this.mergeIndexMap = new Map()
    this.load()
  },
  methods: {
    load(){
      this.loading = true
      request.get("/api/elements", {
        params: {
          "rule_type": 1,
          "page_num": this.currentPage,
          "page_size": this.pageSize,
          "search": this.search
        }
      }).then(res => {
        console.log(res)
        if (res.status_code === 0) {
          this.merge(res.event_elements)
          this.tableData = res.event_elements
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
      request.post("/api/rule", this.addRuleform).then(res => {
        console.log(res)
        if (res.status_code === 0) {
           this.$message({
            type: "success",
            message: "添加成功"
          })
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
      this.addElementForm.rule_id = row.rule_id
      this.dialogAddElementVisible = true
    },
    addElement() {
      this.addElementForm.event_type = parseInt(this.addElementForm.event_type)
      this.addElementForm.mouse_click_type = parseInt(this.addElementForm.mouse_click_type)
      this.addElementForm.mouse_click_button = parseInt(this.addElementForm.mouse_click_button)
      this.addElementForm.key_click_type = parseInt(this.addElementForm.key_click_type)

      request.post("/api/element", this.addElementForm).then(res => {
        console.log(res)
        if (res.status_code === 0) {
          this.$message({
            type: "success",
            message: "添加成功"
          })
          this.dialogAddElementVisible = false
          this.addElementForm = {}
          this.load()
        } else {
          this.$message({
            type: "error",
            message: res.status_msg
          })
        }
      })
    },
    handleEditElement(row) {
      this.curRow = row
      this.dialogUpdateElementVisible = true
    },
    updateElement() {
      this.updateElementform.event_type = parseInt(this.updateElementform.event_type)
      this.updateElementform.mouse_click_type = parseInt(this.updateElementform.mouse_click_type)
      this.updateElementform.mouse_click_button = parseInt(this.updateElementform.mouse_click_button)
      this.updateElementform.key_click_type = parseInt(this.updateElementform.key_click_type)

      request.put("/api/element/" + this.curRow.element_id, this.updateElementform).then(res => {
        console.log(res)
        if (res.status_code === 0) {
           this.$message({
            type: "success",
            message: "更新成功"
          })
          this.dialogUpdateElementVisible = false
          this.updateElementform = {}
          this.load()
        } else {
           this.$message({
            type: "error",
            message: res.status_msg
          })
        }
      })
    },
    handleDeleteElement(row) {
      this.curRow = row
      this.$confirm('确定删除？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(res => {
        console.log(res)
        request.delete("/api/element/" + this.curRow.element_id).then(res => {
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
 