<template>
    <div v-if="state" class="w-full p-12">
        <div class="font-['In_your_face,_Joffrey!'] text-center ">
            <h1 class="text-5xl mb-8">AMONG US: Konfero Edition</h1>
            <div class="grid grid-cols-3 justify-between leading-tight w-full">
                <div class="value-block">
                    <p class="name">Tasks</p>
                    <p class="value">{{ state.SolvedTasks }}/{{ state.TotalTasks }}</p>
                </div>

                <div class="value-block">
                    <p class="name">Players</p>
                    <p class="value">{{ state.AlivePlayers }}/{{ state.TotalPlayers }}</p>
                </div>

                <div v-if="timeLeft > 0" class="value-block">
                    <p class="name">Time</p>
                    <p class="value">{{ minutesLeft }}:{{ secondsLeft }}</p>
                </div>
            </div>
            <div v-if="state.Voting" class="p-8">
                <div v-if="state.Voting.Active">
                    <p class="text-3xl">Votes</p>
                    <div class="text-[250px]/[250px]">
                        {{ voteCount }}/{{ state.AlivePlayers }}
                    </div>
                </div>
                <div v-else class="text-5xl">
                    <p v-if="votesPerPlayer.length" class="text-3xl mb-4">Voting Results</p>
                    <div v-for="player in votesPerPlayer" class="grid grid-cols-[4rem_1fr] items-center w-full mb-6">
                        <div class="text-left">{{ player.playerId.toString().padStart(4, "0") }}</div>
                        <div class="bg-white rounded w-full text-black text-right pr-4 text-5xl/tight"
                            :style="`width: ${player.votes / votesPerPlayer[0]!.votes * 100}%;`">
                            {{ player.votes }}

                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { useApi, type ResultState } from '@/api';
import { computed, onUnmounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const api = useApi()
const state = ref<ResultState>()
const timeLeft = ref(0)
const route = useRoute()

const intervalId = setInterval(loadState, 1000)

onUnmounted(() => {
    clearInterval(intervalId)
})

async function loadState() {
    try {
        state.value = await api.getResultState()
        if (!state.value.Voting) {
            const secondsSinceStart = Date.now() / 1000 - state.value.GameStart;
            timeLeft.value = Math.round(state.value.GameDuration - secondsSinceStart)
        } else {
            timeLeft.value = Math.round(state.value.GameDuration)
        }
        timeLeft.value = Math.max(timeLeft.value, 0)
    } catch {
        useRouter().replace({ path: "/login", query: { returnTo: route.fullPath } })
    }
}

loadState()


const minutesLeft = computed(() => {
    return String(Math.floor(timeLeft.value / 60)).padStart(2, "0")
})
const secondsLeft = computed(() => {
    return String(timeLeft.value % 60).padStart(2, "0")
})

var voteCount = computed(() => {
    if (!state.value?.Voting) return 0
    return state.value.Voting.Votes.filter(v => v != 0).length
})

var votesPerPlayer = computed(() => {
    if (!state.value?.Voting) return []

    var voteMap = new Map<number, number>()
    for (const vote of state.value.Voting.Votes.filter(v => v != 0)) {
        let currentValue = voteMap.get(vote)
        if (!currentValue) {
            currentValue = 0
        }
        voteMap.set(vote, currentValue + 1)
    }
    var results = []

    for (const v of voteMap.entries()) {
        results.push({ playerId: v[0], votes: v[1] })
    }

    results.sort((a, b) => b.votes - a.votes)
    return results

    // return [
    //     {
    //         playerId: 1,
    //         votes: 50
    //     },
    //     {
    //         playerId: 143,
    //         votes: 49
    //     },
    //     {
    //         playerId: 143,
    //         votes: 49
    //     },
    //     {
    //         playerId: 143,
    //         votes: 3
    //     },
    //     {
    //         playerId: 143,
    //         votes: 1
    //     },
    //     {
    //         playerId: 143,
    //         votes: 1
    //     },
    //     {
    //         playerId: 143,
    //         votes: 1
    //     },
    //     {
    //         playerId: 143,
    //         votes: 1
    //     }
    // ]
})

</script>

<style>
@import "tailwindcss";

label {
    @apply text-lg;
}

.value-block {
    @apply text-center
}

.value-block>.value {
    @apply text-[120px]/[100px]
}

.value-block>.name {
    @apply text-3xl
}
</style>