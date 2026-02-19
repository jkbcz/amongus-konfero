<template>
    <div v-if="state" class="p-8 max-w-sm mx-auto">
        <h1 class="text-3xl mb-8">AMONG US: Konfero Edition</h1>
        <!-- {{ state }} -->
        <div class="flex flex-wrap mb-8">
            <div v-for="i in state.TotalStations"
                class="flex items-center justify-center text-lg m-4 w-12 h-12 bg-gray-500 rounded"
                :class="state.FinishedStations.includes(i) ? 'opacity-0' : ''"
            >
                {{i}}
            </div>
        </div>

        <form @submit.prevent="sendCode">
            <input 
                class="w-full text-xl mx-auto border-white border-spacing-4"
                :placeholder="state.CodeMask" 
                @input="upper"
                v-model="code"/>
            <input class="w-full" type="submit"/>
        </form>
    </div>
</template>

<script lang="ts" setup>
import { useApi, type PlayerState } from '@/api';
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
    @apply text-2xl block outline-1 rounded p-2 my-8;
}
</style>