<template>
    <div style="margin: 10px 0">
        <el-input v-model="search" placeholder="请输入关键字" style="width: 20%" clearable></el-input>
        <el-button type="primary" style="margin-left: 5px" @click="load">查询</el-button>
        <el-button type="primary" @click="dialog_add_account = true">
            添加账号
        </el-button>

        <el-dialog v-model="dialog_add_account" title="添加账号" width="500px">
            <el-form :model="add_account_form" label-width="100px" label-position="left">
                <el-form-item label="账号名" style="width: 80%">
                    <el-input v-model="add_account_form.account_name"></el-input>
                </el-form-item>
                <el-form-item label="账号密码" style="width: 80%">
                    <el-input v-model="add_account_form.account_pwd"></el-input>
                </el-form-item>
                <el-form-item label="确认账号密码" style="width: 80%">
                    <el-input v-model="add_account_form.account_pwd2"></el-input>
                </el-form-item>
                <el-form-item label="账号权限" style="width: 80%">
                    <el-select v-model="add_account_form.account_permission" placeholder="">
                        <el-option v-for="item in permission_option" :key="item[0]" :label="item[1]" :value="item[0]"></el-option>
                    </el-select>
                </el-form-item>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                <el-button @click="dialog_add_account = false">取消</el-button>
                <el-button type="primary" @click="add_account">
                    确定添加
                </el-button>
                </span>
            </template>
        </el-dialog>

        <el-table 
            v-loading="loading"
            :data="account_data" 
            style="width: 100%"
        >
            <el-table-column fixed prop="account_id" label="账号ID" width="150" />
            <el-table-column fixed prop="account_name" label="账号名" width="150" />
            <el-table-column fixed label="账号权限" width="150">
                <template #default="scope">
                    <el-tag
                    class="mx-1"
                    :disable-transitions="false"
                    >
                    {{ this.permission_desc.get(scope.row.account_permission) }}
                    </el-tag>
                </template>
            </el-table-column>
            <el-table-column label="操作" width="300">
                <template #default="scope">
                <el-button link type="primary" size="small" @click="handleClickOpenDialog(scope.$index)"
                    >修改</el-button
                >
                <el-button link type="danger" size="small" @click="handleDeleteUser(scope.row.user_id)"
                    >删除</el-button
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

export default {
    name: "Permission",
    data() {
        return {
            search: "",
            loading: true,
            dialog_add_account: false,
            currentPage: 1,
            pageSize: 10,
            total: 0,
            add_account_form: {},
            account_data: new Array(),
            permission_option: new Map([
                [1, "普通"], 
                [2, '管理员'], 
            ]),
            permission_desc: new Map([
                [1, "普通"], 
                [2, '管理员'], 
                [3, '超级管理员'], 
            ]),
        }
    },
    created() {
        this.load()
    },
    methods: {
        load() {
            this.loading = true
            request.get("/api/accounts", {
                params: {
                "page_num": this.currentPage,
                "page_size": this.pageSize,
                "search": this.search
                }
            }).then(res => {
                console.log(res)
                if (res.status_code === 0) {
                    this.account_data = res.accounts
                    this.total = res.total
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
        add_account() {
            console.log(this.add_account_form)
            if (this.add_account_form.account_pwd != this.add_account_form.account_pwd2) {
                this.$message({
                    type: "info",
                    message: "密码不一致"
                })
                return 
            }
            
            request.post("/api/account", this.add_account_form).then(res => {
                console.log(res)
                if (res.status_code === 0) {
                    this.$message({
                        type: "success",
                        message: "添加成功"
                    })
                    this.load()
                    this.dialog_add_account = false;
                    this.add_account_form = {}
                } else {
                    this.$message({
                        type: "error",
                        message: res.status_msg
                    })
                }
            })
        }
    }


}
</script>


