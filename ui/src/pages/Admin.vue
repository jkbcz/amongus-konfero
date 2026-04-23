<template>
    <div class="w-screen h-screen flex items-center justify-center p-8" >
        <div v-if="state">
            <h1 class="text-5xl mb-8">Admin</h1>
            <!-- {{ state }} -->
            <div class="flex justify-between">
                <button @click="updateSetting('is_voting', !state.IsVoting)" class="bg-green-500 mb-4">
                    {{state.IsVoting ? 'Stop Voting' : 'Start Voting'}}
                </button>
                <button @click="updateSetting('reset', true)" class="bg-red-500 mb-4">
                    Reset game
                </button>
            </div>

            <label>Required Codes</label>
            <div class="flex gap-4 mb-4">
                <input v-model="state.RequiredCodes" type="number" />
                <input type="submit" @click="updateSetting('required_codes', state.RequiredCodes)"/>
            </div>

            <label>Total Stations</label>
            <div class="flex gap-4 mb-4">
                <input v-model="state.TotalStations" type="number" />
                <input type="submit" @click="updateSetting('total_stations', state.TotalStations)"/>
            </div>

            <label>Code Mask</label>
            <div class="flex gap-4 mb-4">
                <input v-model="state.CodeMask" />
                <input type="submit" @click="updateSetting('code_mask', state.CodeMask)"/>
            </div>

            <label>Cooldown Duration (seconds)</label>
            <div class="flex gap-4 mb-4">
                <input v-model="state.CooldownDuration" type="number" />
                <input type="submit" @click="updateSetting('cooldown_duration', state.CooldownDuration+'s')"/>
            </div>

            <label>Game Duration (seconds)</label>
            <div class="flex gap-4 mb-4">
                <input v-model="state.GameDuration" type="number" />
                <input type="submit" @click="updateSetting('game_duration', state.GameDuration+'s')"/>
            </div>
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
import { useApi, type AdminState } from '@/api';
import { getPass, savePass } from '@/pass';
import { ref } from 'vue';
import { useToast } from 'vue-toast-notification';

const api = useApi()
const state = ref<AdminState>()
const toast = useToast()
const pass = ref(getPass())

async function loadState() {
    state.value = await api.getAdminState(pass.value)
}

async function updateSetting(setting: string, value: any) {
    try {
        await api.updateSetting(setting, value, pass.value)
        await loadState()
        toast.success("updated settings")
    } catch (err: any) {
        toast.error(err.message)
    }
}

async function login() {
    await loadState()
    savePass(pass.value)
}

login()
</script>

<style>
@import "tailwindcss";

label {
    @apply text-lg;
}

button {
    @apply p-4 block  rounded
}

input {
    @apply block outline-1 rounded p-2 my-2;
}
</style>