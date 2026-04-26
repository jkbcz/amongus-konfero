<template>
    <div class="grid grid-cols-4 lg:grid-cols-8 flex-wrap gap-2">
        <template v-for="(isDead, i) in playersArr">
            <div v-if="i > 0" class="cursor-pointer text-center border-2 border-white" :class="isDead ? 'opacity-50' : ''"
                @click="emit('playerClick', i)">
                {{ i.toString().padStart(4, "0") }}
            </div>
        </template>
    </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

const props = defineProps({
    rawPlayers: {type: String, required: true},
    totalPlayers: {type: Number, required: true}
})

const playersArr = computed(() => {
    const result: boolean[] = []
    if (!props.rawPlayers) {
        return result
    }
    var bytes = Uint8Array.fromBase64(props.rawPlayers)

    for (let i = 0; i <= props.totalPlayers; i++) {
        result.push((bytes[Math.floor(i / 8)] & 1 << i % 8) > 0)
    }
    return result
})

const emit = defineEmits({
    "playerClick": (id: number) => true
})

</script>
