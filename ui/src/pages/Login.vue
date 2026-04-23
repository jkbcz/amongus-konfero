<template>
    <div class="w-screen h-screen flex items-center justify-center ">
        <div class="rounded-md p-4 border-2 border-white">
            <form @submit.prevent="login">
                <div>
                    <label for="pass">Password</label>
                    <input id="pass" type="password" v-model="pass" />
                </div>
                <input type="submit" class="float-right" />
            </form>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { useApi } from '@/api';
import { saveAdminPass, saveUserPass } from '@/pass';
import { ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute()
const router = useRouter()

const returnTo: string = route.query["returnTo"] as string ?? "/"

const pass = ref("")
const api = useApi()

async function login() {
    if (returnTo == "/") {
        saveUserPass(pass.value)
        await api.getPlayerState()
    } else {
        saveAdminPass(pass.value)
        await api.getResultState()
    }
    router.replace({ path: returnTo })
}
</script>