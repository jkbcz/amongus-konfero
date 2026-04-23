<template>
    <video :src="randomGif" autoplay muted loop></video>
</template>

<script lang="ts" setup>
import danger from "@/img/gifs/danger.mp4"
import dicaprio from "@/img/gifs/dicaprio.webm"
import emergencyMeeting from "@/img/gifs/emergency-meeting.webm"
import futurama from "@/img/gifs/futurama.webm"
import spider from "@/img/gifs/spider.webm"
import survivor from "@/img/gifs/survivor.webm"
import sweating from "@/img/gifs/sweating.webm"

const allGifs = [
    danger,
    dicaprio,
    emergencyMeeting,
    futurama,
    spider,
    survivor,
    sweating
]

const randomGif = allGifs[Math.floor(Math.random() * allGifs.length)]
</script>