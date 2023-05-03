<template>
    <div style="margin: 10px 0">
      <el-input v-model="search" placeholder="请输入关键字" style="width: 20%" clearable></el-input>
      <el-button type="primary" style="margin-left: 5px" @click="load">查询</el-button>
      <el-button type="primary" style="margin-left: 5px" @click="this.dialog_add_task_visible = true">新增任务</el-button>

      <el-dialog v-model="dialog_add_task_visible" title="添加数据挖掘任务">
        <el-form :model="add_task_form">
          <el-form-item label="任务名">
            <el-input v-model="add_task_form.task_name"></el-input>
          </el-form-item>
          <el-form-item label="最小支持度(百分比)">
            <el-input v-model="add_task_form.percent"></el-input>
          </el-form-item>
        </el-form>
        <template #footer>
        <span class="dialog-footer">
          <el-button @click="this.dialog_add_task_visible = false">取消</el-button>
          <el-button type="primary" @click="add_seq_mining_task">
            确定添加
          </el-button>
        </span>
      </template>
      </el-dialog>

      <el-table 
        v-loading="loading"
        :data="task_data" 
        style="width: 100%"
      >
        <el-table-column prop="task_id" label="任务ID" width="100" />
        <el-table-column prop="task_name" label="任务名" width="150" />
        <el-table-column prop="create_time" label="创建时间" width="150" />
        <el-table-column prop="percent" label="最小支持度(百分比)" width="200" />
        <el-table-column prop="status_desc" label="任务状态" width="150" >
          <template #default="scope">
            <el-tag
              class="mx-1"
              :disable-transitions="false"
            >
              {{ scope.row.status_desc }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column  label="操作" width="200">
          <template #default="scope">
            <el-button type="primary" size="small" :disabled="scope.row.task_status != 3" @click="handle_download_result(scope.row.task_id)">
              下载任务数据
              </el-button>
          </template>
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
    name: "SeqMining",
    data() {
      return {
        loading: true,
        search: "",
        currentPage: 1,
        pageSize: 10,
        total: 0,
        task_data : new Array(),
        dialog_add_task_visible: false,
        add_task_form: {},
        status_desc: new Map([
          [1, "已创建"], 
          [2, '运行中'], 
          [3, '已完成']
        ])
        
      }
    },
    created() {
      this.load()
    },
    methods: {
      load() {
        this.loading = true
        request.get("/api/seq_mining", {
          params: {
            "page_num": this.currentPage,
            "page_size": this.pageSize,
            "search": this.search
          }
        }).then(res => {
          console.log(res)
          if (res.status_code === 0) {
            this.task_data = res.seq_mining_tasks
            for (let i = 0; i < this.task_data.length; i++) {
              this.task_data[i].status_desc = this.status_desc.get(this.task_data[i].task_status)
            }
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
      handleSizeChange(pageSize) {  // 改变每页的大小
        this.pageSize = pageSize;
        this.load();
      },
      handleCurrentChange(pageNum) {  //改变当前页码
        this.currentPage = pageNum;
        this.load()
      },
      add_seq_mining_task() {
        this.add_task_form.percent = parseInt(this.add_task_form.percent)
        console.log(this.add_task_form)
        request.post("/api/seq_mining", this.add_task_form).then(res => {
          console.log(res)
          if (res.status_code === 0) {
            this.$message({
              type: "success",
              message: "添加成功"
            })
            this.dialog_add_task_visible = false
            this.add_seq_mining_task = {}
            this.load()
          } else {
            this.$message({
              type: "error",
              message: res.status_msg
            })
          }
        })
      },
      handle_download_result(task_id) {
        request.get("/api/seq_mining_result/" + task_id).then(res => {

          if (res.status_code != null && res.status_code != 0) {
            this.$message({
              type: "error",
              message: "下载失败"
            })
            return 
          }
          const url = window.URL.createObjectURL(new Blob([res]));
          const link = document.createElement('a');
          link.href = url;
          link.setAttribute('download', 'result.zip');
          document.body.appendChild(link);
          link.click();
        })
      }
    }
}


</script>

<style scoped>

</style>