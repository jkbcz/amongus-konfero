<template>
    <div v-if="state" class="relative max-w-sm mx-auto">
        <div class="px-8">
            <h1 class="text-5xl mt-8 font-['In_your_face,_Joffrey!']">AMONG US: Konfero Edition</h1>
            <div class="text-sm mb-8 flex justify-between items-center text-gray-400">
                <p>Player ID: {{ state.PlayerId.toString().padStart(4, "0") }}</p>
                <button @click="goToLogin" class="cursor-pointer underline p-0"> logout</button>
            </div>
        </div>


        <div v-if="state.IsDead" class="text-center px-8">
            <div>
                <h1 class="text-5xl mb-16">You Are Dead!</h1>
            </div>
        </div>
        <div v-else-if="state.VotingState && !state.VotingState.Active" class="text-left px-8">
            <h1 class="text-5xl mb-8">Meeting Time!</h1>
            <p class="text-gray-400 mb-4">Go to the meeting area immediately</p>
            <!-- <RandomGif /> -->
        </div>
        <div v-else-if="state.VotingState">
            <div class="px-8">
                <h1 class="text-5xl mb-8">Voting Time!</h1>
                <p class="mb-4">Vote who to eject:</p>
                <div class="grid grid-cols-4 flex-wrap gap-2 mb-8">
                    <template v-for="(isDead, i) in state.VotingState.Players">
                        <div v-if="i > 0" class="cursor-pointer text-center border-2 select-none"
                            :class="isDead ? 'opacity-50' : playerClass(i)" @click="selectPlayer(i)">
                            {{ i.toString().padStart(4, "0") }}
                        </div>
                    </template>
                </div>
            </div>

            <div v-if="state.VotingState.MyVote"
                class="sticky bottom-0 left-0 w-full bg-white text-black border-4 border-b-0 border-white rounded-t-xl p-4">
                <p class="mb-7">Voted to eject: {{ state.VotingState.MyVote.toString().padStart(4, "0") }}</p>
            </div>
            <div v-else-if="playerToEject"
                class="sticky bottom-0 left-0 w-full bg-white text-black border-4 border-b-0 border-white rounded-t-xl p-4">
                <div class="flex justify-between items-center mb-2">
                    <div>
                        <p>Voting to eject: {{ playerToEject.toString().padStart(4, "0") }}</p>
                        <p class="text-gray-400 text-sm">Note: This action cannot be undone</p>
                    </div>
                    <button class="rounded-full border-2 p-2" @click="sendVote">Confirm</button>
                </div>
            </div>
        </div>
        <div v-else class="px-8">
            <form @submit.prevent="sendCode">
                <label for="code">Code</label>
                <input class="w-full text-2xl mx-auto border-white" id="code" :placeholder="state.CodeMask"
                    @input="upper" autocomplete="off" v-model="code" />
                <input class="w-full text-2xl" type="submit" />
            </form>

            <p class="mt-12 text-lg mb-4">Remaining Stations</p>
            <div class="flex flex-wrap gap-8">
                <div v-for="i in state.TotalStations"
                    class="flex items-center justify-center text-lg w-12 h-12 bg-black border-2  rounded"
                    :class="state.FinishedStations.includes(i) ? 'text-gray-800 border-gray-800' : 'border-white'">
                    {{ i }}
                </div>
            </div>
        </div>

    </div>
</template>

<script lang="ts" setup>
import { useApi, type PlayerState } from '@/api';
import { saveUserPass } from '@/pass';
import { onUnmounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useToast } from 'vue-toast-notification';

const api = useApi()
const toast = useToast()
const router = useRouter()
const state = ref<PlayerState>()
const routePass = useRoute().query["pass"]
if (typeof routePass === "string") {
    saveUserPass(routePass)
    router.replace({ query: {} })
}
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
    saveUserPass("")
    router.replace({ path: "/login" })
}

function upper() {
    code.value = code.value.toUpperCase()
}

function playerClass(i: number) {
    if (i == state.value?.PlayerId) {
        return 'opacity-50'
    } else if (i == playerToEject.value || i == state.value?.VotingState?.MyVote) {
        return 'text-black bg-white border-white'
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

const playerToEject = ref(0)
function selectPlayer(playerId: number) {
    if (playerId == state.value?.PlayerId) {
        return;
    }
    if (state.value?.VotingState?.Players[playerId]) {
        return;
    }
    if (state.value?.VotingState?.MyVote) {
        return;
    }
    playerToEject.value = playerId;
}

async function sendVote() {
    try {
        await api.sendVote(playerToEject.value)
        toast.success("Success")
        await loadState()
        playerToEject.value = 0
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