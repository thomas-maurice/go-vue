<script>

import axios from 'axios'
import { useUserStore } from '@/stores/user'
import router from '@/router'
import NavBar from '@/components/NavBar.vue'
import { API_BASE_URL } from '@/defaults/client'

import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { library } from '@fortawesome/fontawesome-svg-core'
import { faIdCardClip } from '@fortawesome/free-solid-svg-icons'

import {  AuthenticationApi, ApiLoginInput } from '@/gen/apiclient/src'
import { DefaultClient } from '@/apiclient/client'
var authApiClient = new AuthenticationApi(DefaultClient)

library.add(faIdCardClip)

export default {
    data() {
        return {
            input: {
                username: "you@youremail.com",
                password: "password",
            },
            error: undefined,
            loginMessage: "",
            userStore: useUserStore(),
            router: router,
            client: authApiClient,
            oidcProviders: undefined,
        }
    },
    components: {
        NavBar: NavBar,
        FontAwesomeIcon, FontAwesomeIcon,
    },
    methods: {
        async oidc(name) {
            let resp = await axios.get(`${API_BASE_URL}/api/auth/oidc/${name}`)
            location.href = resp.data.url
        },
        getProviders() {
            authApiClient.authOidcProvidersGet().then(response => {
                this.oidcProviders = response
            }).catch(error => {
                console.log(error)
            })
        },
        login() {
            var input = new ApiLoginInput()
            input.username = this.input.username
            input.password = this.input.password
            authApiClient.authLoginPost(input).then(response => {
                this.userStore.log_in(
                    this.input.username,
                    response.token,
                )
                this.error = undefined
                router.push("/")
            })
            .catch(error => {
                this.error = `Login failed: ${error.error}`
            });
        }
    },
    created: function() {
        this.getProviders()
        this.loginMessage = this.$route.query.message
    }
}

</script>

<template>
    <div class="d-flex justify-content-center align-items-center min-vh-100">
      <div class="w-100" style="max-width: 400px;">
        <h1><FontAwesomeIcon :icon="['fas', 'id-card-clip']" /> Log into the app</h1>
        <form>
            <div v-if="loginMessage" class="alert alert-primary">
                {{ loginMessage }}
            </div>
            <div class="mb-3">
                <input type="text" class="form-control" id="login" aria-describedby="loginHelp" v-model="input.username" autofocus>
                <div id="loginHelp" class="form-text">Username or email</div>
            </div>
            <div class="mb-3">
                <input type="password" class="form-control" id="password" v-model="input.password" v-on:keyup.enter="login">
                <div id="passwordHelp" class="form-text">Password</div>
            </div>
            <div v-if="error" class="alert alert-danger" role="alert">
                {{ error }}
            </div>
            <button type="button" v-on:click="login" class="btn btn-primary w-100 mb-2">Submit</button>
            <div v-for="provider in oidcProviders">
                <button type="button" v-on:click="oidc(provider.name)" class="btn btn-primary w-100">Login with {{ provider.display_name  }}</button>
            </div>
        </form>
      </div>
    </div>
</template>