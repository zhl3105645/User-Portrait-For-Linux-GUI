<template>
    <div style="margin: 10px 0">
        <el-input v-model="search" placeholder="请输入关键字" style="width: 20%" clearable></el-input>
        <el-button type="primary" style="margin-left: 5px" @click="load">查询</el-button>
        <el-button type="primary" style="margin-left: 5px" @click="this.dialog_add_crowd_visible = true">添加人群</el-button>


        <el-dialog v-model="dialog_add_crowd_visible" title="添加人群" width="650px">
            <el-form :model="add_crowd_form" label-width="100px" label-position="left">
                <el-form-item label="人群名">
                    <el-input v-model="add_crowd_form.crowd_name" style="width: 60%"/>
                </el-form-item>
                <el-form-item label="人群描述">
                    <el-input v-model="add_crowd_form.crowd_desc" style="width: 60%"/>
                </el-form-item>
                <el-form-item label="人群划分规则">
                    <div v-for="(rule, index) in divide_rules" :key="index">
                        <!-- 标签选择 -->
                        <el-select v-model="rule.label_id" style="width: 20%;margin-left: 5px" placeholder="标签">
                            <el-option v-for="item in label_options" :key="item.label_id" :label="item.label_name" :value="item.label_id"></el-option>
                        </el-select>
                        <!-- 操作 -->
                        <el-select v-model="rule.divide_operate" style="width: 20%;margin-left: 5px" placeholder="操作符">
                            <el-option v-for="item in divide_operate_types" :key="item.key" :label="item.desc" :value="item.key"></el-option>
                        </el-select>
                        <!-- 标签数据 -->
                        <el-select v-model="rule.label_data" v-if="label_2_enum.get(rule.label_id) == true" style="width: 20%;margin-left: 5px" placeholder="标签数据">
                            <el-option v-for="item in label_2_desc.get(rule.label_id)" :key="item.data" :label="item.desc" :value="item.data"></el-option>
                        </el-select>
                        <el-input v-model="rule.label_data" v-else style="width: 20%;margin-left: 5px" placeholder="标签数据">
                        </el-input>   
                        <el-select v-model="rule.union_operate" style="width: 20%;margin-left: 5px" v-if="index < this.divide_rules.length-1" placeholder="连接方式">
                            <el-option v-for="item in union_operate_types" :key="item.key" :label="item.desc" :value="item.key"></el-option>
                        </el-select>  
                    </div>
                    <el-button @click="add_divide_rule">添加规则</el-button>
                </el-form-item>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                <el-button @click="dialog_add_crowd_visible = false">取消</el-button>
                <el-button type="primary" @click="add_crowd">
                    确定添加
                </el-button>
                </span>
            </template>
        </el-dialog>

        <el-table 
            v-loading="loading"
            :data="crowds" 
            style="width: 100%"
            >
        <el-table-column fixed prop="crowd_id" label="人群ID" width="100" />
        <el-table-column prop="crowd_name" label="人群名" width="150" />
        <el-table-column  prop="crowd_desc" label="人群描述" width="200" />
        <el-table-column  prop="user_num" label="用户数量" width="100" />
        <el-table-column label="操作" width="400">
            <template #default="scope">
                <el-button type="primary" size="small" @click="handel_update_crowd(scope.$index)"
                    >修改人群</el-button
                >
                <el-button type="primary" size="small" @click="gene_crowd(scope.row.crowd_id)"
                    >更新人群数据</el-button
                >
                <el-button type="danger" size="small" @click="handle_delete_crowd(scope.row.crowd_id)"
                    >删除人群</el-button
                >
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
import * as echarts from 'echarts'

export default {
    name: "Crowd",
    data() {
        return {
            crowds: [],
            loading: true,
            search: "",
            currentPage: 1,
            pageSize: 10,
            total: 0,
            dialog_add_crowd_visible: false,
            add_crowd_form: {},
            divide_rules: [], // Or 列表元素为 And
            union_operate_types: [
                {
                    key: 1,
                    desc: "AND",
                },
                {
                    key: 2,
                    desc: "OR",
                }
            ],
            divide_operate_types: [
                {
                    key: 1,
                    desc: ">",
                },
                {
                    key: 2,
                    desc: ">=",
                },
                {
                    key: 3,
                    desc: "=",
                },
                {
                    key: 4,
                    desc: "<",
                },
                {
                    key: 5,
                    desc: "<=",
                }
            ],
            label_options: [],
            label_2_enum: {}, // 是否可枚举
            label_2_desc: {}, // labelId -> [] 枚举数据
        }
    },
    created(){
        this.load_label_option()
        this.load()
    },
    methods: {
        load(){
            this.loading = true
            request.get("/api/crowd", {
                params: {
                "page_num": this.currentPage,
                "page_size": this.pageSize,
                "search": this.search
                }
            }).then(res => {
                console.log(res)
                if (res.status_code === 0) {
                    this.crowds = res.crowds
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
        load_label_option() {
            request.get("/api/labels").then(res => {
            console.log(res)
            if (res.status_code === 0) {
                this.label_options = res.labels
                this.label_2_enum = new Map()
                this.label_2_desc = new Map()
                for (let i = 0; i < this.label_options.length; i++) {
                    let label = this.label_options[i]
                    if (label.label_data_type == 1){
                        this.label_2_enum.set(label.label_id, true)
                        this.label_2_desc.set(label.label_id, label.data)
                    } else {
                        this.label_2_enum.set(label.label_id, false)
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
        add_divide_rule() {
            this.divide_rules.push({})
        },
        add_crowd() {
            if (this.divide_rules.length <= 0) {
                return 
            }
            this.divide_rules[this.divide_rules.length-1].union_operate = 1
            this.add_crowd_form.divide_rules = this.divide_rules
            console.log(this.add_crowd_form)
            request.post("/api/crowd", this.add_crowd_form).then(res => {
                if (res.status_code === 0) {
                this.$message({
                    type: "success",
                    message: "添加成功"
                })
                this.load()
                this.add_crowd_form = {}
                this.divide_rules = []
                this.dialog_add_crowd_visible = false
                } else {
                this.$message({
                    type: "error",
                    message: res.status_msg
                })
                }
            })
        },
        gene_crowd(crowd_id) {
            request.post("/api/crowd/" + crowd_id).then(res => {
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
        handle_delete_crowd(crowd_id) {
            this.$confirm('确定删除？', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(res => {
                request.delete("/api/crowd/" + crowd_id).then(res => {
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
        }
    }
    
}
</script>

