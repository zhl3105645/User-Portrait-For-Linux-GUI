<template>
    <div class="container">
        <el-select v-model="crowd_id" placeholder="请选择用户群">
            <el-option v-for="crowd in crowds" :key="crowd.crowd_id" :label="crowd.crowd_name" :value="crowd.crowd_id"></el-option>
        </el-select>
        <el-button type="primary" style="margin-left: 5px" @click="load_crowd_profile()">生成画像</el-button>

        <div class="top-1"  ref="pie_chart_1">

        </div>
        <div class="top-2" ref="pie_chart_2">

        </div>
        <div class="top-3" ref="pie_chart_3">

        </div>
        <div class="top-4" ref="pie_chart_4">

        </div>
        <div class="top-5" ref="bar_chart">

        </div>
        <div class="circle" ref="behavior_duration_chart">

        </div>
        <div class="bottom-1" ref="stack_bar_chart" >

        </div>
    </div>
</template>

<script>
import request from "@/utils/request";
import * as echarts from 'echarts'

export default {
    name: "GroupProfile",
    data() {
        return {
            crowds: [{
                "crowd_id": -1,
                "crowd_name": "全部用户"
            }],
            crowd_id: null,
            radars: null,
            stack_bar_label: null,
            bar_label: null,
            pie_label: null,

            behavior_duration_chart: null,
            pie_chart_1: null,
            pie_chart_2: null,
            pie_chart_3: null,
            pie_chart_4: null,
            stack_bar_chart: null,
            bar_chart: null,
        }
    },
    created() {
        this.load_crowds()
    },
    mounted(){
        this.behavior_duration_chart = echarts.init(this.$refs.behavior_duration_chart,null, {renderer: 'svg'})
        this.pie_chart_1 = echarts.init(this.$refs.pie_chart_1,null, {renderer: 'svg'})
        this.pie_chart_2 = echarts.init(this.$refs.pie_chart_2,null, {renderer: 'svg'})
        this.pie_chart_3 = echarts.init(this.$refs.pie_chart_3,null, {renderer: 'svg'})
        this.pie_chart_4 = echarts.init(this.$refs.pie_chart_4,null, {renderer: 'svg'})
        this.stack_bar_chart = echarts.init(this.$refs.stack_bar_chart,null, {renderer: 'svg'})
        this.bar_chart = echarts.init(this.$refs.bar_chart,null, {renderer: 'svg'})
    },
    methods: {
        load_crowds() {
            request.get("/api/crowds").then(res => {
                console.log(res)
                if (res.status_code === 0) {
                    this.crowds.push.apply(this.crowds, res.crowds)
                } else {
                    this.$message({
                        type: "error",
                        message: res.status_msg
                    })
                }
            })
        },
        load_crowd_profile() {
            request.get("/api/group_profile/" + this.crowd_id).then(res => {
                console.log(res)
                if (res.status_code === 0) {
                    this.radars = res.radars
                    this.stack_bar_label = res.stack_bar_label
                    this.pie_label = res.pie_label
                    this.bar_label = res.bar_label
                    this.behavior_duration_chart.clear()
                    this.pie_chart_1.clear()
                    this.pie_chart_2.clear()
                    this.pie_chart_3.clear()
                    this.pie_chart_4.clear()
                    this.stack_bar_chart.clear()
                    this.bar_chart.clear()
                    this.set_radar()
                    this.set_pie_chart()
                    this.set_stack_bar_chart()
                    this.set_bar_chart()
                } else {
                    this.$message({
                        type: "error",
                        message: res.status_msg
                    })
                }
            })
        },
        set_radar() {
            let indicator = new Array()
            let appDuration = new Array()
            let userDuration = new Array()
            for (let i = 0; i < this.radars.length; i++) {
                let radar = this.radars[i]
                indicator.push({
                name: radar.name,
                max: radar.max,
                })
                appDuration.push(radar.ave)
                userDuration.push(radar.cur)
            }

            let option = {
                // toolbox: {
                //     show: true,
                //     feature: {
                //         mark: {show: true},
                //         saveAsImage: {show: true},
                //     }
                // },
                title: {
                    text: '行为时长',
                },
                legend: {
                    right: "0",
                    data: ['应用平均时长', '人群平均时长']
                },
                radar: {
                    radius: 99,
                    indicator: indicator,
                    name: {
                        textStyle: {
                        color: '#333',
                        fontSize: 14
                        }
                    },
                    splitLine: {
                        lineStyle: {
                        color: '#999'
                        }
                    },
                    splitArea: {
                        areaStyle: {
                        color: ['rgba(250,250,250,0.3)', 'rgba(200,200,200,0.3)']
                        }
                    },
                    axisLine: {
                        lineStyle: {
                        color: '#999'
                        }
                    },
                },
                series: [
                {
                    name: 'Budget vs spending',
                    type: 'radar',
                    lineStyle: {
                    width: 2
                    },
                    data: [
                    {
                        value: appDuration,
                        name: "应用平均时长",
                    },
                    {
                        value: userDuration,
                        name: "人群平均时长",
                    }
                    ]
                }
                ]
            }

            if (this.crowd_id == -1) {
                option.data = ["应用平均时长"]
                option.series[0].data = [option.series[0].data[0]]
            }
            this.behavior_duration_chart.setOption(option)
        },
        set_pie_chart() {
            if (this.pie_label == null || this.pie_label.length < 1) {
                return 
            }
            let option1 = {
                // toolbox: {
                //     show: true,
                //     feature: {
                //         mark: {show: true},
                //         saveAsImage: {show: true},
                //     }
                // },
                tooltip: {
                    trigger: 'item'
                },
                title: {
                    text: this.pie_label[0].label_name,
                    left: 'auto'
                },
                series: [
                    {
                        type: 'pie',
                        radius: ['40%', '70%'],
                        center: ['50%', '60%'],
                        avoidLabelOverlap: false,
                        itemStyle: {
                            borderRadius: 10,
                            borderColor: '#fff',
                            borderWidth: 2
                        },
                        label: {
                            normal: {
                                position: "inside"//此处将展示的文字在内部展示
                            }
                        },
                        data: this.pie_label[0].data,
                    }
                ]
            }
            this.pie_chart_1.setOption(option1)

            if (this.pie_label.length < 2) {
                return 
            }
            let option2 = {
                // toolbox: {
                //     show: true,
                //     feature: {
                //         mark: {show: true},
                //         saveAsImage: {show: true},
                //     }
                // },
                tooltip: {
                    trigger: 'item'
                },
                title: {
                    text: this.pie_label[1].label_name,
                    left: 'auto'
                },
                series: [
                    {
                        type: 'pie',
                        radius: ['40%', '70%'],
                        center: ['50%', '60%'],
                        avoidLabelOverlap: false,
                        itemStyle: {
                            borderRadius: 10,
                            borderColor: '#fff',
                            borderWidth: 2
                        },
                        label: {
                            normal: {
                                position: "inside"//此处将展示的文字在内部展示
                            }
                        },
                        data: this.pie_label[1].data,
                    }
                ]
            }
            this.pie_chart_2.setOption(option2)

            if (this.pie_label.length < 3) {
                return 
            }
            let option3 = {
                // toolbox: {
                //     show: true,
                //     feature: {
                //         mark: {show: true},
                //         saveAsImage: {show: true},
                //     }
                // },
                tooltip: {
                    trigger: 'item'
                },
                title: {
                    text: this.pie_label[2].label_name,
                    left: 'auto'
                },
                series: [
                    {
                        type: 'pie',
                        radius: ['40%', '70%'],
                        center: ['50%', '60%'],
                        avoidLabelOverlap: false,
                        itemStyle: {
                            borderRadius: 10,
                            borderColor: '#fff',
                            borderWidth: 2
                        },
                        label: {
                            normal: {
                                position: "inside"//此处将展示的文字在内部展示
                            }
                        },
                        data: this.pie_label[2].data,
                    }
                ]
            }
            this.pie_chart_3.setOption(option3)

            if (this.pie_label.length < 4) {
                return 
            }
            let option4 = {
                // toolbox: {
                //     show: true,
                //     feature: {
                //         mark: {show: true},
                //         saveAsImage: {show: true},
                //     }
                // },
                tooltip: {
                    trigger: 'item'
                },
                title: {
                    text: this.pie_label[3].label_name,
                    left: 'auto'
                },
                series: [
                    {
                        type: 'pie',
                        radius: ['40%', '70%'],
                        center: ['50%', '60%'],
                        avoidLabelOverlap: false,
                        itemStyle: {
                            borderRadius: 10,
                            borderColor: '#fff',
                            borderWidth: 2
                        },
                        label: {
                            normal: {
                                position: "inside"//此处将展示的文字在内部展示
                            }
                        },
                        data: this.pie_label[3].data,
                    }
                ]
            }
            this.pie_chart_4.setOption(option4)
        },
        set_stack_bar_chart() {
            let series = new Array()
            let x = this.stack_bar_label.label_value_desc;
            for (let i = 0; i < this.stack_bar_label.label_value_desc.length; i++) {
                series.push({
                    name: this.stack_bar_label.label_value_desc[i],
                    type: 'bar',
                    stack: 'common',
                    emphasis: {
                        focus: 'series'
                    },
                    data: this.stack_bar_label.label_cnt[i],
                    label: {
                        show: true,
                        position: 'inside',
                        // formatter: '{a} {b} {c}: {@score}'
                        formatter: function (params) {
                            //console.log(params)
                            //console.log(x.indexOf(params.seriesName))
                            if (x.indexOf(params.seriesName) > -1 && params.data > 0) {
                                return params.seriesName
                            }
                            return "";
                        }
                    }
                })
            }
            //console.log(series)
            let option = {
                // toolbox: {
                //     show: true,
                //     feature: {
                //         mark: {show: true},
                //         saveAsImage: {show: true},
                //     }
                // },
                tooltip: {
                    trigger: 'axis',
                },
                title: {
                    text: '自定义标签',
                },
                xAxis: [
                    {
                        type: 'category',
                        data: this.stack_bar_label.label_names,
                        axisLabel: {
                            interval:0,
                            rotate:15
                        }
                    }
                ],
                yAxis: [
                    {
                        name: '用户数',
                        type: 'value',
                        nameLocation: 'center',
                        nameGap: 30,
                    }
                ],
                series: series,
            }

            this.stack_bar_chart.setOption(option)
        },
        set_bar_chart() {
            let option = {
                // toolbox: {
                //     show: true,
                //     feature: {
                //         mark: {show: true},
                //         saveAsImage: {show: true},
                //     }
                // },
                tooltip: {
                    trigger: 'axis',
                    axisPointer: {
                    type: 'shadow'
                    }
                },
                title: {
                    text: '使用时间段',
                    left: '0',
                },
                xAxis: [
                    {
                        name: "小时",
                        type: 'category',
                        data: this.bar_label.x_names,
                        axisTick: {
                            alignWithLabel: true
                        },
                        
                    }
                ],
                yAxis: [
                    {
                        name: '用户数',
                        type: 'value',
                        nameLocation: 'center',
                        nameGap: 30,
                    }
                ],
                grid: {
                    left: '10%',
                    right: '15%',
                    bottom: '3%',
                    containLabel: true
                },
                series: [
                    {
                        name: 'Direct',
                        type: 'bar',
                        barWidth: '60%',
                        data: this.bar_label.data,
                        showBackground: true,
                        backgroundStyle: {
                            color: 'rgba(180, 180, 180, 0.2)'
                        },
                        itemStyle: {
                            normal: {
                            //这里是重点
                            color: function(params) {
                                var colorList = ['#c23531','#2f4554', '#61a0a8', '#d48265', '#91c7ae','#749f83', '#ca8622'];
                                // var colorList = ['#c23531','#2f4554', '#61a0a8'];
                                // 自动循环已经有的颜色
                                return colorList[params.dataIndex % colorList.length];
                            }
                            }
                        }
                    }
                ]
            }
            this.bar_chart.setOption(option)
        }
    }
}
</script>

<style scoped>
.container {
  position: relative;
  width: 100%;
  height: 100vh;
}



.top-1 {
  top: 50px;
  left: 0px;
  position: absolute;
  width: 200px;
  height: 200px;
  /* background-color: rgb(217, 235, 238); */
}
.top-2 {
  top: 50px;
  left: 240px;
  position: absolute;
  width: 200px;
  height: 200px;
  /* background-color: rgb(217, 235, 238); */
}
.top-3 {
  top: 50px;
  left: 480px;
  position: absolute;
  width: 200px;
  height: 200px;
  /* background-color: rgb(217, 235, 238); */
}
.top-4 {
  top: 50px;
  left: 720px;
  position: absolute;
  width: 200px;
  height: 200px;
  /* background-color: rgb(217, 235, 238); */
}

.top-5 {
  top: 50px;
  left: 960px;
  position: absolute;
  width: 300px;
  height: 250px;
  /* background-color: rgb(217, 235, 238); */
}

.bottom-1 {
  position: absolute;
  bottom: 50px;
  left: 0px;
  width: 1000px;
  height: 400px;
  /* background-color: rgb(192, 167, 167); */
}
.circle {
  position: absolute;
  bottom: 50px;
  left: 960px;
  width: 350px;
  height: 350px;
  /* background-color: rgb(224, 217, 217); */
}

</style>