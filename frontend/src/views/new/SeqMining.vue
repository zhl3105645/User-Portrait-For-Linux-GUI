<template>
    <div class="container">
    <div class="box top-left">
      <div class="arrow arrow-bottom-right"></div>
      1
    </div>
    <div class="box top-right">
      <div class="arrow arrow-bottom-left"></div>
      2
    </div>
    <div class="circle">3</div>
    <div class="box bottom-left">
      <div class="arrow arrow-top-right"></div>
      4
    </div>
    <div class="box bottom-right">
      <div class="arrow arrow-top-left"></div>
      5
    </div>
  </div>
</template>

<script>
export default {
    name: "SeqMining",
}


</script>

<style scoped>
.container {
  position: relative;
  width: 100%;
  height: 100vh;
}
.box {
  position: absolute;
  width: 100px;
  height: 200px;
  background-color: lightblue;
}
.top-left {
  top: 0;
  left: 0;
}
.top-right {
  top: 0;
  right: 0;
}
.bottom-left {
  bottom: 0;
  left: 0;
}
.bottom-right {
  bottom: 0;
  right: 0;
}
.circle {
  position: absolute;
  top: calc(50% - 100px);
  left: calc(50% - 100px);
  width: 200px;
  height: 200px;
  border-radius: 50%;
  background-color: lightcoral;
}
.arrow {
  position: absolute;
  width: calc(50% - (50px + (100px / sqrt(2))));
  height: calc(50% - (100px + (100px / sqrt(2))));
}
.arrow-bottom-right:before,
.arrow-bottom-left:before,
.arrow-top-right:before,
.arrow-top-left:before {
    content:'';
    position:absolute;
    border-top:1px solid black;
    border-right:1px solid black;
    width:10px;height:10px;
}
.arrow-bottom-right:before{
    transform-origin:right top;transform:rotate(-45deg);
    right:-6px;top:-6px
}
.arrow-bottom-left:before{
    transform-origin:left top;transform:rotate(45deg);
    left:-6px;top:-6px
}
.arrow-top-right:before{
    transform-origin:right bottom;transform:rotate(45deg);
    right:-6px;bottom:-6px
}
.arrow-top-left:before{
    transform-origin:left bottom;transform:rotate(-45deg);
    left:-6px;bottom:-6px
}
.arrow-bottom-right,
.arrow-bottom-left,
.arrow-top-right,
.arrow-top-left {
    overflow:hidden
}
.arrow-bottom-right{
    border-bottom-right-radius:10px
}
.arrow-bottom-left{
    border-bottom-left-radius:10px
}
.arrow-top-right{
    border-top-right-radius:10px
}
.arrow-top-left{
    border-top-left-radius:10px
}
</style>