<script setup lang="ts">
import Starfield from './components/Starfield.vue';

</script>

<template>
  <div class="relative">
    <starfield class="-z-10"/>
    <div class="absoulte top-0 w-screen">

      <router-view/>
    </div>
  </div>
</template>

<style scoped></style>
