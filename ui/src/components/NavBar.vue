<script>

import { useUserStore as useUserStore } from '@/stores/user'
import { useRouter, useRoute } from 'vue-router'

import { AuthenticationApi } from '@/gen/apiclient/src'
import { GetDefaultClient } from '@/apiclient/client'

import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { library } from '@fortawesome/fontawesome-svg-core'
import { faRightFromBracket, faRightToBracket, faIdCardClip } from '@fortawesome/free-solid-svg-icons'

import { RouterLink } from 'vue-router';

library.add(
  faRightFromBracket,
  faRightToBracket,
  faIdCardClip,
)

var userStore = useUserStore()

export default {
    props: {
        userContext: Object,
    },
    data() {
        return {
            userStore: useUserStore(),
            router: useRouter(),
            route: useRoute(),
            apiClient: new AuthenticationApi(GetDefaultClient()),
        }
    },
    components: {
        FontAwesomeIcon: FontAwesomeIcon,
        RouterLink: RouterLink,
    },
    methods: {
      logout() {
        this.apiClient.authLogoutPost().then(response => {
          this.userStore.logout()
          this.router.push('/login')
        }).catch(error => {
          console.log("failed to log out", error)
        })
      }
    }
}
</script>

<template>
  <nav class="navbar navbar-expand-lg bg-body-tertiary">
  <div class="container-fluid">
    <a class="navbar-brand">
        <RouterLink class="nav-link" :to="{name: 'home'}">Some App</RouterLink>
    </a>
    <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
      <span class="navbar-toggler-icon"></span>
    </button>
    <div class="collapse navbar-collapse" id="navbarSupportedContent">
      <ul class="navbar-nav me-auto mb-2 mb-lg-0">
        <li class="nav-item">
          <RouterLink :to="{name: 'home'}" style="text-decoration: none; color: inherit;">
            <span>Home</span>
          </RouterLink>
        </li>
      </ul>
      <div v-if="route.name != 'oidc-callback'">
        <RouterLink class="nav-link" :to="{  name: 'login' }">
          <span v-if="!userStore.logged_in" class="btn btn-outline-success">
            <FontAwesomeIcon :icon="['fas', 'right-to-bracket']" /> Login
          </span>
        </RouterLink>
        <RouterLink :to="{name: 'profile-self'}" class="nav-item px-2" style="text-decoration: none; color: inherit;">
          <span v-if="userStore.logged_in" class="btn btn-outline-primary">
            <FontAwesomeIcon :icon="['fas', 'id-card-clip']" /> Profile
          </span>
        </RouterLink>
        <span v-if="userStore.logged_in" class="btn btn-outline-danger" v-on:click="logout()">
          <FontAwesomeIcon :icon="['fas', 'right-from-bracket']" /> Logout ({{ userStore.username }})
        </span>
      </div>
    </div>
  </div>
</nav>
</template>