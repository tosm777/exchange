
<template>
  <router-view></router-view>
</template>

<script setup>
import { useUserStore } from "@/store/user"
import { useRecipientStore } from "@/store/recipient"
const user = useUserStore()
const recp = useRecipientStore()
</script>


<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  /* text-align: center; */
  color: #2c3e50;
  margin-top: 60px;
}
.main {
  display: block;
  margin: 0 auto;
  width: 50%;
}
</style>
