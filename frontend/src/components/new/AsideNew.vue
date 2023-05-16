<template>
  <div>
    <el-menu
      class="el-menu-vertical-demo"
      style="width: 200px; min-height: calc(100vh - 50px)"
        :default-active="path"
        router
    >
      <el-submenu index="1">
        <template v-slot:title>账号管理</template>
        <el-menu-item-group>
          <el-menu-item index="/front/permission">账号</el-menu-item>
        </el-menu-item-group>
      </el-submenu>
      <el-submenu index="2">
        <template v-slot:title>用户管理</template>
        <el-menu-item-group>
          <el-menu-item index="/front/user_record">用户数据</el-menu-item>
        </el-menu-item-group>
      </el-submenu>
      <el-submenu index="3">
        <template v-slot:title>规则管理</template>
        <el-menu-item-group>
          <template v-slot:title></template>
          <el-menu-item index="/front/seq_mining">序列模式挖掘</el-menu-item>
          <el-menu-item index="/front/event_rule">事件规则</el-menu-item>
          <el-menu-item index="/front/behavior_rule">行为规则</el-menu-item>
          <el-menu-item index="/front/event_rule_data">规则数据</el-menu-item>
        </el-menu-item-group>
      </el-submenu>
      <el-submenu index="4">
        <template v-slot:title>标签管理</template>
        <el-menu-item-group>
          <template v-slot:title></template>
          <el-menu-item index="/front/label">标签</el-menu-item>
          <el-menu-item index="/front/crowd">人群</el-menu-item>
        </el-menu-item-group>
      </el-submenu>
      <el-submenu index="5">
        <template v-slot:title>画像查看</template>
        <el-menu-item-group>
          <template v-slot:title></template>
            <el-menu-item index="/front/profile">单用户画像</el-menu-item>
            <el-menu-item index="/front/group_profile">人群画像</el-menu-item>
        </el-menu-item-group>
      </el-submenu>
      
    </el-menu>
  </div>
</template>

<script>
import request from "@/utils/request";

export default {
  name: "AsideNew",
  data() {
    return {
      path: this.$route.path  //设置默认高亮菜单
    }
  },
  methods:{
    
  },
  created() {
  }
}
</script>

<style scoped>

</style>