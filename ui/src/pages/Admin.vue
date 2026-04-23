<template>
    <div class="w-screen h-screen flex items-center justify-center p-8">
        <div v-if="state">
            <h1 class="text-5xl mb-8">Admin</h1>
            <!-- {{ state }} -->
            <div class="flex justify-between">
                <button @click="updateSetting('is_voting', !state.IsVoting)" class="bg-green-500 mb-4">
                    {{ state.IsVoting ? 'Stop Voting' : 'Start Voting' }}
                </button>
                <button @click="updateSetting('reset', true)" class="bg-red-500 mb-4">
                    Reset game
                </button>
            </div>

            <label>Required Codes</label>
            <div class="flex gap-4 mb-4">
                <input v-model="state.Settings.RequiredCodes" type="number" />
                <input type="submit" @click="updateSetting('required_codes', state.Settings.RequiredCodes)" />
            </div>

            <label>Total Stations</label>
            <div class="flex gap-4 mb-4">
                <input v-model="state.Settings.TotalStations" type="number" />
                <input type="submit" @click="updateSetting('total_stations', state.Settings.TotalStations)" />
            </div>

            <label>Total Players</label>
            <div class="flex gap-4 mb-4">
                <input v-model="state.Settings.TotalPlayers" type="number" />
                <input type="submit" @click="updateSetting('total_players', state.Settings.TotalPlayers)" />
            </div>

            <label>Code Mask</label>
            <div class="flex gap-4 mb-4">
                <input v-model="state.Settings.CodeMask" />
                <input type="submit" @click="updateSetting('code_mask', state.Settings.CodeMask)" />
            </div>

            <label>Cooldown Duration (seconds)</label>
            <div class="flex gap-4 mb-4">
                <input v-model="state.Settings.CooldownDuration" type="number" />
                <input type="submit"
                    @click="updateSetting('cooldown_duration', state.Settings.CooldownDuration + 's')" />
            </div>

            <label>Game Duration (seconds)</label>
            <div class="flex gap-4 mb-4">
                <input v-model="state.Settings.GameDuration" type="number" />
                <input type="submit" @click="updateSetting('game_duration', state.Settings.GameDuration + 's')" />
            </div>

            <h1>Players</h1>

            <div class="flex flex-wrap gap-4">
                <template v-for="(isDead, i) in playersArr">
                    <button v-if="i > 0" class="p-2 border-2 border-white" :class="isDead ? 'opacity-50' : ''"
                        @click="togglePlayer(i)">
                        {{ i.toString().padStart(4, "0") }}
                    </button>
                </template>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { useApi, type AdminState } from '@/api';
import { computed, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useToast } from 'vue-toast-notification';

const api = useApi()
const state = ref<AdminState>()
const toast = useToast()
const route = useRoute()


async function loadState() {
    try {
        state.value = await api.getAdminState()
    } catch {
        useRouter().replace({ path: "/login", query: { returnTo: route.fullPath } })
    }
}

loadState()

async function updateSetting(setting: string, value: any) {
    try {
        await api.updateSetting(setting, value)
        await loadState()
        toast.success("updated settings")
    } catch (err: any) {
        toast.error(err.message)
    }
}

const playersArr = computed(() => {
    const result: boolean[] = []
    if (!state.value) {
        return result
    }
    var bytes = Uint8Array.fromBase64(state.value.Players)

    for (let i = 0; i <= state.value.Settings.TotalPlayers; i++) {
        result.push((bytes[Math.floor(i / 8)] & 1 << i % 8) > 0)
    }
    return result
})

async function togglePlayer(playerId: number) {
    await api.togglePlayer(playerId)
    await loadState()
}

</script>

<style>
@import "tailwindcss";

label {
    @apply text-lg;
}

button {
    @apply p-4 block rounded
}

input {
    @apply block outline-1 rounded p-2 my-2;
}
</style>