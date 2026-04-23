<template>
    <div class="w-screen h-screen flex items-center justify-center" >
        <div v-if="state?.IsVoting" class="text-center">
            <h1 class=" text-[128px]">Scheduled Meeting!</h1>
            <RandomGif class="mx-auto w-150"/>
        </div>
        <div v-else-if="state"
            class="relative">
            <div :class="state.CooldownUntil >= currentTime.valueOf() ? 'blur-lg' : ''">

                <h1 class="mb-16 text-[128px] text-center">STATION {{ stationId }}</h1>
                <p class="text-[72px] text-center p-8 bg-black border-white border-2 rounded">{{state.CurrentCode}}</p>
                <!-- {{ state }} -->
            </div>
            <Countdown class="absolute top-0 left-1/2 -translate-x-1/2" :duration="state.CooldownDuration" :show-until="state.CooldownUntil"/>
        </div>
        <div v-else class="rounded-md p-4 border-2 border-white">
            <form @submit.prevent="login" >
                <div>
                    <label for="pass">Password</label>
                    <input id="pass" type="password" v-model="pass"/>
                </div>
                <input type="submit" class="float-right"/>
            </form>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { useApi, type StationState } from '@/api';
import Countdown from '@/components/Countdown.vue';
import RandomGif from '@/components/RandomGif.vue';
import { getPass, savePass } from '@/pass';
import { ref } from 'vue';
import { useRoute } from 'vue-router';
const api = useApi()
const state = ref<StationState>()
const route = useRoute()

const stationId = parseInt(route.query['station_id'] as string ?? "1")
const pass = ref(getPass())

async function loadState() {
    state.value = await api.getStationState(stationId, pass.value)
}


async function login() {
    await loadState()
    savePass(pass.value)
    setInterval(loadState, 1000)
}

login()


const currentTime = ref(new Date())
setInterval(() => {
    currentTime.value = new Date();
}, 100);
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