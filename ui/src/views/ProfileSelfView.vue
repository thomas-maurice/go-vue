<script>

import NavBar from '@/components/NavBar.vue'
import AuthGuard from '@/components/AuthGuard.vue'

import { UserApi } from '@/gen/apiclient/src'
import { GetDefaultClient } from '@/apiclient/client'
import moment from 'moment'

export default {
  components: {
    NavBar: NavBar,
    AuthGuard: AuthGuard,
  },
  data() {
    return {
      user: {
        display_name: undefined,
        id: undefined,
        admin: undefined,
        email: undefined,
        last_login: undefined,
      },
      error: undefined,
      moment: moment,
      apiClient: new UserApi(GetDefaultClient()),
    }
  },
  created: () => {
  },
  methods: {
    getUserProfile() {
      this.apiClient.userProfileGet().then( response => {
        this.user = response
        this.error = undefined
      }).catch(error => {
        this.error = error
      })
    },
  },
  created: function() {
    this.getUserProfile()
  }
}
</script>

<template>
  <NavBar />
  <AuthGuard />
  <div class="pt-3 container-md col-8 justify-content-center align-items-center">
    <h1>Your profile</h1>
    <ul>
      <li v-if="!error">Display name: {{ user.display_name }}</li>
      <li>Id: <span class="font-monospace">{{ user.id }}</span></li>
      <li>Kind: <span class="font-monospace">{{ user.kind }}</span></li>
      <li>Username: <span class="font-monospace">{{ user.username}}</span></li>
      <li>Admin: <span class="font-monospace">{{ user.admin }}</span></li>
      <li>Email: <span class="font-monospace">{{ user.email }}</span></li>
      <li>Created: {{ moment(user.created).fromNow() }} ({{ (moment(user.created)) }})</li>
      <li>Last login: {{ moment(user.last_login).fromNow() }} ({{ (moment(user.last_login)) }})</li>
    </ul>
  </div>
</template>
