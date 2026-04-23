<template>
    <div class="p-8 max-w-sm mx-auto">
        <div v-if="state?.IsVoting" class="text-center flex justify-center items-center h-screen">
            <div>
                <h1 class="text-5xl mb-16">Scheduled Meeting!</h1>
                <RandomGif class="mx-auto w-full" />
            </div>
        </div>
        <div v-if="state?.IsDead" class="text-center flex justify-center items-center h-screen">
            <div>
                <h1 class="text-5xl mb-16">You Are Dead!</h1>
                <RandomGif class="mx-auto w-full" />
            </div>
        </div>
        <div v-else-if="state">
            <h1 class="text-5xl mb-8 font-['In_your_face,_Joffrey!']">AMONG US: Konfero Edition</h1>

            <form @submit.prevent="sendCode">
                <label for="code">Code</label>
                <input class="w-full text-2xl mx-auto border-white" id="code" :placeholder="state.CodeMask"
                    @input="upper" autocomplete="off" v-model="code" />
                <input class="w-full text-2xl" type="submit" />
            </form>

            <p class="mt-12 text-lg">Remaining Stations</p>
            <div class="flex flex-wrap">
                <div v-for="i in state.TotalStations"
                    class="flex items-center justify-center text-lg m-4 w-12 h-12 bg-black border-2  rounded"
                    :class="state.FinishedStations.includes(i) ? 'text-gray-800 border-gray-800' : 'border-white'">
                    {{ i }}
                </div>
            </div>

            <div class=" mt-12 flex justify-between items-center text-gray-400">
                <p class="text-lg">Logged in as: {{ state.PlayerId.toString().padStart(4, "0") }}</p>
                <input @click="goToLogin" class="p-1 text-sm" type="submit" value="logout" />
            </div>
        </div>
    </div>

</template>

<script lang="ts" setup>
import { useApi, type PlayerState } from '@/api';
import RandomGif from '@/components/RandomGif.vue';
import { onUnmounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from 'vue-toast-notification';

const api = useApi()
const toast = useToast()
const router = useRouter()
const state = ref<PlayerState>()
loadState()
const intervalId = setInterval(loadState, 1000)

onUnmounted(() => {
    clearInterval(intervalId)
})

async function loadState() {
    try {
        state.value = await api.getPlayerState()
    } catch {
        goToLogin()
    }
}

function goToLogin() {
    router.replace({ path: "/login" })
}

function upper() {
    code.value = code.value.toUpperCase()
}

async function sendCode() {
    try {
        await api.submitCode(code.value)
        toast.success("Success")
        code.value = ""
        await loadState()
    } catch (err: any) {
        toast.error(err.message)
    }
}
const code = ref("")
</script>

<style scoped>
@import "tailwindcss";

input {
    @apply block outline-1 rounded p-2;
}
</style>