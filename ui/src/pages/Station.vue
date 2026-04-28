<template>
    <div class="w-screen h-screen flex items-center justify-center">
        <div v-if="state?.IsVoting" class="text-center">
            <h1 class=" text-[128px]">Meeting Time!</h1>
            <!-- <RandomGif class="mx-auto w-150" /> -->
        </div>
        <div v-else-if="state" class="relative">
            <div :class="state.CooldownUntil >= currentTime.valueOf() ? 'blur-lg' : ''">

                <h1 class="mb-16 text-[128px] text-center">STATION {{ stationId }}</h1>
                <p class="text-[72px] text-center p-8 bg-black border-white border-2 rounded">{{ state.CurrentCode }}
                </p>
            </div>
            <Countdown class="absolute top-0 left-1/2 -translate-x-1/2" :duration="state.CooldownDuration"
                :show-until="state.CooldownUntil" />
        </div>
    </div>
</template>

<script lang="ts" setup>
import { useApi, type StationState } from '@/api';
import Countdown from '@/components/Countdown.vue';
import RandomGif from '@/components/RandomGif.vue';
import { onUnmounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
const api = useApi()
const state = ref<StationState>()
const route = useRoute()

const stationId = parseInt(route.query['station_id'] as string ?? "1")

async function loadState() {
    try {
        state.value = await api.getStationState(stationId)
    } catch {
        useRouter().replace({ path: "/login", query: { returnTo: route.fullPath } })
    }
}

const currentTime = ref(new Date())


loadState()
const stateInterval = setInterval(loadState, 1000)
const timingIntervale = setInterval(() => {
    currentTime.value = new Date();
}, 100);

onUnmounted(() => {
    clearInterval(stateInterval)
    clearInterval(timingIntervale)
})


</script>

<style>
@import "tailwindcss";

label {
    @apply text-lg;
}

input {
    @apply block outline-1 rounded p-2 my-2;
}
</style>