<template>
  <div>
    <el-menu
      style="width: 200px; min-height: calc(100vh - 50px)"
        :default-active="path"
        router
    >
      <el-menu-item index="/front/user_record">用户</el-menu-item>
      <el-submenu index="2">
        <template v-slot:title>规则</template>
        <el-menu-item-group>
          <template v-slot:title></template>
          <!-- <el-menu-item index="/front/seq_mining">序列模式挖掘</el-menu-item> -->
          <el-menu-item index="/front/event_rule">事件规则</el-menu-item>
          <el-menu-item index="/front/behavior_rule">行为规则</el-menu-item>
          <el-menu-item index="/front/event_rule_data">规则数据</el-menu-item>
        </el-menu-item-group>
      </el-submenu>
      <el-menu-item index="/front/label">标签</el-menu-item>
      <el-menu-item index="/front/crowd">人群</el-menu-item>
      <el-menu-item index="/front/profile">画像</el-menu-item>
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