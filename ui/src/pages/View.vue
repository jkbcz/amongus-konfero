<template>
    <div class="w-screen h-screen flex items-center justify-center " >
        <div v-if="state" class="font-['In_your_face,_Joffrey!'] text-center text-[200px]">
            <h1 class="text-5xl mb-8">AMONG US: Konfero Edition</h1>
            <p>Tasks: {{state.SolvedTasks}}/{{ state.TotalTasks }}</p>
            <p v-if="timeLeft > 0">Time: {{ minutesLeft }}:{{ secondsLeft }}</p>
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
import { useApi, type ResultState } from '@/api';
import { getPass, savePass } from '@/pass';
import { computed, ref } from 'vue';

const api = useApi()
const state = ref<ResultState>()

const pass = ref(getPass())
const timeLeft = ref(0)

async function loadState() {
    state.value = await api.getResultState(pass.value)
    if(!state.value.IsVoting) {
        const secondsSinceStart =  Date.now()/1000 - state.value.GameStart;
        timeLeft.value = Math.round(state.value.GameDuration - secondsSinceStart)
    } else {
        timeLeft.value = Math.round(state.value.GameDuration)
    }
    timeLeft.value = Math.max(timeLeft.value, 0)
}

const minutesLeft = computed(() => {
    return String(Math.floor(timeLeft.value / 60)).padStart(2, "0")
})
const secondsLeft = computed(() => {
    return String(timeLeft.value % 60).padStart(2, "0")
})



async function login() {
    await loadState()
    savePass(pass.value)
    setInterval(loadState, 1000)
}

login()
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