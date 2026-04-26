<script setup lang="ts">
import Starfield from './components/Starfield.vue';

</script>

<template>
  <div>
    <starfield class=" fixed top-0 w-screen -z-10"/>

    <router-view/>
    <!-- </div> -->
  </div>
</template>

<style scoped></style>
