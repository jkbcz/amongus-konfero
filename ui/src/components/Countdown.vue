<template>
  <div class="timer-container" 
    v-if="timeLeft > 0"
    >
    <div class="relative flex items-center justify-center">
      <svg
        class="-rotate-90"
        width="500"
        height="500"
        viewBox="0 0 100 100"
      >
        
        <circle
          class="timer-circle-progress stroke-gray-600"
          cx="50"
          cy="50"
          :r="radius"
          stroke-width="50"
          :style="circleStyle"
        />
      </svg>
      
      <div class="absolute text-[128px]" >
        {{ formattedTime }}
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue';

const props = defineProps({
  showUntil: {type: Number, required: true},
  duration: {type: Number, required: true}
})

const totalTime = 10; 
const strokeWidth = 50;

const radius = 50 - (strokeWidth / 2);
const circumference = 2 * Math.PI * radius;

// State
const timeLeft = ref(calculateTime());

// Logic
setInterval(() => {
  timeLeft.value = calculateTime()
}, 100);

function calculateTime(): number {
  return (props.showUntil - new Date().valueOf()) / 1000
}

const formattedTime = computed(() => {
  return Math.ceil(timeLeft.value);
});

const circleStyle = computed(() => {
  const progress = timeLeft.value / props.duration;
  const dashOffset = circumference - (progress * circumference);

  return {
    strokeDasharray: circumference,
    strokeDashoffset: dashOffset,
  };
});

</script>

<style scoped>


.timer-circle-progress {
  fill: none;
  /* stroke: #ff4757; */
  transition: stroke-dashoffset 0.1s linear;
}

</style>