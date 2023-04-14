<template>
  <div>
    <el-menu
      style="width: 200px; min-height: calc(100vh - 50px)"
        :default-active="path"
        router
    >
      <el-submenu index="1">
        <template v-slot:title><span>用户</span></template>
        <el-menu-item-group>
          <template v-slot:title></template>
          <el-menu-item index="/front/user_record">用户行为记录</el-menu-item>
        </el-menu-item-group>
      </el-submenu>
      <el-submenu index="2">
        <template v-slot:title>规则</template>
        <el-menu-item-group>
          <template v-slot:title></template>
          <el-menu-item index="/front/data_mean">数据含义</el-menu-item>
          <el-menu-item index="/front/event_rule">事件规则</el-menu-item>
          <el-menu-item index="/front/behavior_rule">行为规则</el-menu-item>
        </el-menu-item-group>
      </el-submenu>
      <el-submenu index="3">
        <template v-slot:title>数据</template>
        <el-menu-item-group>
          <template v-slot:title></template>
          <el-menu-item index="/front/basic_behavior_data">基本行为数据</el-menu-item>
          <el-menu-item index="/front/event_rule_data">事件规则数据</el-menu-item>
          <el-menu-item index="/front/behavior_rule_data">行为规则数据</el-menu-item>
        </el-menu-item-group>
      </el-submenu>
      <el-submenu index="4">
        <template v-slot:title>模型</template>
        <el-menu-item-group>
          <template v-slot:title></template>
          <el-menu-item index="/front/statistics_model">统计模型</el-menu-item>
          <el-menu-item index="/front/learning_model">机器模型</el-menu-item>
        </el-menu-item-group>
      </el-submenu>
      <el-submenu index="5">
        <template v-slot:title>应用</template>
        <el-menu-item-group>
          <template v-slot:title></template>
          <el-menu-item index="/front/label">标签</el-menu-item>
          <el-menu-item index="/front/predict">事件预测</el-menu-item>
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