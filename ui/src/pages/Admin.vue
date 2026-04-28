<template>
    <div v-if="state" class="p-8 max-w-5xl m-auto">
        <h1 class="text-5xl mb-8">Admin</h1>
        <!-- {{ state }} -->
        <div class="flex justify-between">
            <div class="flex gap-4">
                <button v-if="!state.Voting" @click="updateSetting('is_voting', true)" class="bg-green-500 mb-4">
                    Start Meeting
                </button>
                <button v-if="state.Voting && !state.Voting.Active" @click="updateSetting('voting_active', true)"
                    class="bg-green-500 mb-4">
                    Start Voting
                </button>
                <button v-if="state.Voting && state.Voting.Active" @click="updateSetting('voting_active', false)"
                    class="bg-green-500 mb-4">
                    Finish Voting
                </button>
                <button v-if="state.Voting" @click="updateSetting('is_voting', false)" class="bg-red-500 mb-4">
                    End Meeting
                </button>
            </div>
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
            <input type="text" v-model="state.Settings.CodeMask" />
            <input type="submit" @click="updateSetting('code_mask', state.Settings.CodeMask)" />
        </div>

        <label>Cooldown Duration (seconds)</label>
        <div class="flex gap-4 mb-4">
            <input v-model="state.Settings.CooldownDuration" type="number" />
            <input type="submit" @click="updateSetting('cooldown_duration', state.Settings.CooldownDuration + 's')" />
        </div>

        <label>Game Duration (seconds)</label>
        <div class="flex gap-4 mb-4">
            <input v-model="state.Settings.GameDuration" type="number" />
            <input type="submit" @click="updateSetting('game_duration', state.Settings.GameDuration + 's')" />
        </div>

        <h1>Players</h1>

        <div class="grid grid-cols-4 lg:grid-cols-8 flex-wrap gap-2">
            <template v-for="(isDead, i) in state.Players">
                <div v-if="i > 0" class="cursor-pointer text-center border-2 border-white"
                    :class="isDead ? 'opacity-50' : ''" @click="togglePlayer(i)">
                    {{ i.toString().padStart(4, "0") }}
                </div>
            </template>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { useApi, type AdminState } from '@/api';
import PlayerList from '@/components/PlayerList.vue';
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

async function togglePlayer(playerId: number) {
    await api.togglePlayer(playerId)
    await loadState()
}

</script>

<style scoped>
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

input[type="number"],
input[type="text"] {
    @apply w-full
}
</style>