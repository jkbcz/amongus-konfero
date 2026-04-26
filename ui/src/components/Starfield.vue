<template>
  <div >
    
    <div class="star-layer layer-1" :style="{ '--star-shadow': smallStars }"></div>
    <div class="star-layer layer-2" :style="{ '--star-shadow': mediumStars }"></div>
    <div class="star-layer layer-3" :style="{ '--star-shadow': largeStars }"></div>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';

// Define the boundaries of our starfield (3000px ensures it works on ultrawide monitors)
const SCROLL_WIDTH = 3000; 

// Helper function to generate a massive box-shadow string of random coordinates
const generateStars = (count: number) => {
  let shadows = [];
  for (let i = 0; i < count; i++) {
    const x = Math.floor(Math.random() * SCROLL_WIDTH);
    const y = Math.floor(Math.random() * 2000); // 2000px height is plenty
    // Push a shadow at X Y with white color
    shadows.push(`${x}px ${y}px #FFFFFF`); 
  }
  return shadows.join(', ');
};

const smallStars = ref('');
const mediumStars = ref('');
const largeStars = ref('');

onMounted(() => {
  // Generate the layers. Far away = more stars. Close up = fewer stars.
  smallStars.value = generateStars(700);
  mediumStars.value = generateStars(200);
  largeStars.value = generateStars(50);
});
</script>

<style scoped>
/* Base star layer setup */
.star-layer {
  position: absolute;
  top: 0;
  left: 0;
  background: transparent;
}

/* The ::after pseudo-element creates an exact duplicate of the starfield 
  placed precisely at the end of the first one to create a seamless infinite loop.
*/
.star-layer::after {
  content: " ";
  position: absolute;
  top: 0;
  left: 3000px; /* Must match the SCROLL_WIDTH from the script */
  background: transparent;
}

/* Layer 1: Distant, small, very slow */
.layer-1, .layer-1::after {
  width: 1px;
  height: 1px;
  box-shadow: var(--star-shadow);
}
.layer-1 { animation: moveStars 150s linear infinite; }

/* Layer 2: Mid-distance, medium size, medium speed */
.layer-2, .layer-2::after {
  width: 2px;
  height: 2px;
  box-shadow: var(--star-shadow);
}
.layer-2 { animation: moveStars 100s linear infinite; }

/* Layer 3: Close up, large, fast */
.layer-3, .layer-3::after {
  width: 3px;
  height: 3px;
  box-shadow: var(--star-shadow);
}
.layer-3 { animation: moveStars 50s linear infinite; }

/* The infinite scrolling animation */
@keyframes moveStars {
  from {
    transform: translateX(0);
  }
  to {
    transform: translateX(-3000px); /* Must match the SCROLL_WIDTH */
  }
}
</style>