<script>

import NavBar from '@/components/NavBar.vue';
import { useUserStore } from '@/stores/user';

import { DefaultApi } from '@/gen/apiclient/src'
import { GetDefaultClient } from '@/apiclient/client'
import AuthGuard from '@/components/AuthGuard.vue';

export default {
  components: {
    NavBar: NavBar,
    AuthGuard: AuthGuard,
  },
  data() {
    return {
      userStore: useUserStore(),
      ping: undefined,
      apiClient: new DefaultApi(GetDefaultClient()),
    }
  },
  created: () => {
  },
  methods: {
    doPing() {
      this.apiClient.pingGet().then( response => {
        this.ping = response.pong
      }).catch(error => {
        this.ping = error
      })
    },
  },
  mounted: function() {
    this.doPing()
  }
}
</script>

<template>
  <NavBar />
  <AuthGuard />
  <div class="pt-3 container-md col-8 justify-content-center align-items-center">
    <h1>Hello {{  userStore.username }}</h1>
    <div>{{ ping }}</div>
  </div>
</template>
