<script>

import axios from 'axios'
import NavBar from '@/components/NavBar.vue'
import { useUserStore } from '@/stores/user'
import { useRoute, useRouter } from 'vue-router'
import { API_BASE_URL } from '@/defaults/client'

import {  AuthenticationApi, ApiOIDCCallbackOutput } from '@/gen/apiclient/src'
import { DefaultClient } from '@/apiclient/client'
var authApiClient = new AuthenticationApi(DefaultClient)

export default {
    data() {
        return {
            error: undefined,
            success: undefined,
            authDone: false,
            userStore: useUserStore(),
            route: useRoute(),
            router: useRouter(),
            intervalRedirect: null,
            client: authApiClient,
        }
    },
    components:{
        NavBar: NavBar,
    },
    methods: {
        login_oidc(provider, state, code) {
            this.client.authCallbackNameGet(provider, state, code).then(response => {
                this.error = undefined
                this.router.push({ query: {} })

                this.userStore.log_in(response.username, response.token)

                this.intervalRedirect = setTimeout( () => {
                    this.router.push("/")
                }, 3000)
            }).catch(error => {
                console.log(error)
                if (error.response && error.response.status == 401) {
                    this.error = `Invalid credentials: ${error.body.error}`
                } else if (error.response && error.response.status == 500) {
                    console.log(error)
                    this.error = `Internal server error: ${error.body.error}`
                } else {
                    this.error = `Unknown error: ${error.error}`
                }
            }).finally(() => {
                this.authDone = true
            })
        }
    },  
    created: function() {
        this.login_oidc(this.route.params.name, this.route.query.state, this.route.query.code)
    },
}

</script>

<template>
    <NavBar />
    <div class="pt-3 container-md col-8 justify-content-center align-items-center">
        <h1>OIDC Callback</h1>
        <div v-if="!authDone">
            <div class="spinner-border" role="status">
                <span class="visually-hidden">Completing authentication - hang tight</span>
            </div>
            <h3 class="m-3 mt-3 d-inline-block"> Completing authentication - hang tight</h3>
        </div>
        <div v-if="authDone && !error">
            <div class="spinner-border" role="status">
                <span class="visually-hidden"></span>
            </div>
            <h3 class="m-3 mt-3 d-inline-block"> Authentication successful, you will be redirected</h3>
        </div>
        <div v-if="authDone && error">
            <h3> Authentication failed, contact your admin</h3>
        </div>
        <div class="alert alert-danger font-monosapce" role="alert" v-if="error">
            {{ error }}
        </div>
    </div>
</template>