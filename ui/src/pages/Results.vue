<template>
    <div class="w-screen h-screen flex items-center justify-center" >
        <div v-if="state">

                <h1 class="mb-16 text-[128px] text-center">Remaining Tasks</h1>
                <p class="text-[128px] text-center">{{state.RemainingTasks}}</p>
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
import { useApi, type ResultState, type StationState } from '@/api';
import { getPass, savePass } from '@/pass';
import { ref } from 'vue';

const api = useApi()
const state = ref<ResultState>()

const pass = ref(getPass())

async function loadState() {
    state.value = await api.getResultState(pass.value)
}


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