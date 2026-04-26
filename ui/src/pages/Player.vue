<template>
    <div v-if="state" class="px-8 max-w-sm mx-auto">
        <h1 class="text-5xl mt-8 font-['In_your_face,_Joffrey!']">AMONG US: Konfero Edition</h1>
        <div class="text-sm mb-8 flex justify-between items-center text-gray-400">
            <p>Player ID: {{ state.PlayerId.toString().padStart(4, "0") }}</p>
            <button @click="goToLogin" class="cursor-pointer underline p-0"> logout</button>
        </div>

        <div v-if="state.IsDead" class="text-center">
            <div>
                <h1 class="text-5xl mb-16">You Are Dead!</h1>
                <RandomGif class="mx-auto w-full" />
            </div>
        </div>
        <div v-else-if="state.VotingState">
            <div>
                <h1 class="text-5xl mb-8">Killing Time!</h1>
                <p class="mb-4">Choose who to eject:</p>
                <div class="grid grid-cols-4 lg:grid-cols-8 flex-wrap gap-2">
                    <template v-for="(isDead, i) in state.VotingState.Players">
                        <div 
                            v-if="i > 0" 
                            class="cursor-pointer text-center border-2" 
                            :class="isDead ? 'opacity-50' : playerClass(i)"
                            @click="playerToEject = i">
                            {{ i.toString().padStart(4, "0") }}
                        </div>
                    </template>
                </div>
            </div>
        </div>
        <div v-else>

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
        </div>
    </div>
</template>

<script lang="ts" setup>
import { useApi, type PlayerState } from '@/api';
import PlayerList from '@/components/PlayerList.vue';
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

const playerToEject = ref(0)

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

function playerClass(i: number) {
    if(i == state.value?.PlayerId) {
        return 'text-yellow-400'
    } else if (i == playerToEject.value) {
        return 'text-green-400'
    }
    return ''
}


const code = ref("")
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
</script>

<style scoped>
@import "tailwindcss";

input {
    @apply block outline-1 rounded p-2;
}
</style>