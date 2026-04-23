<template>
    <div class="p-8 max-w-sm mx-auto">
        <div v-if="state?.IsVoting" class="text-center flex justify-center items-center h-screen">
            <div>
                <h1 class="text-5xl mb-16">Scheduled Meeting!</h1>
                <RandomGif class="mx-auto w-full"/>
            </div>
        </div>
        <div v-else-if="state">
            <h1 class="text-5xl mb-8 font-['In_your_face,_Joffrey!']">AMONG US: Konfero Edition</h1>
            <!-- {{ state }} -->

            <form @submit.prevent="sendCode">
                <label for="code">Code</label>
                <input 
                    class="w-full text-xl mx-auto border-white"
                    id="code"
                    :placeholder="state.CodeMask" 
                    @input="upper"
                    autocomplete="off"
                    v-model="code"/>
                <input class="w-full" type="submit"/>
            </form>

            <p class="mt-12 text-lg">Remaining Stations</p>
            <div class="flex flex-wrap">
                <div v-for="i in state.TotalStations"
                    class="flex items-center justify-center text-lg m-4 w-12 h-12 bg-black border-2  rounded"
                    :class="state.FinishedStations.includes(i) ? 'text-gray-800 border-gray-800' : 'border-white'"
                >
                    {{i}}
                </div>
            </div>
        
        </div>
    </div>

</template>

<script lang="ts" setup>
import { useApi, type PlayerState } from '@/api';
import RandomGif from '@/components/RandomGif.vue';
import { ref } from 'vue';
import { useToast } from 'vue-toast-notification';

const api = useApi()
const toast = useToast()
const state = ref<PlayerState>()
loadState()
setInterval(loadState, 1000)

async function loadState() {
    state.value = await api.getPlayerState()
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
    } catch(err: any) {
        toast.error(err.message)
    }
}
const code = ref("")
</script>


<style scoped>
@import "tailwindcss";

input {
    @apply text-2xl block outline-1 rounded p-2;
}
</style>