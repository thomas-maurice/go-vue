<script>

import { DefaultApi } from '@/gen/apiclient/src'
import { GetDefaultClient } from '@/apiclient/client'
import { useRoute, useRouter } from 'vue-router'

export default {
    data() {
        return {
            authCheckInterval: null,
            route: useRoute(),
            router: useRouter(),
            apiClient: new DefaultApi(GetDefaultClient())
        }
    },
    mounted() {
        if(!["login", "oidc-callback"].includes(this.route.name)) {
            this.apiClient.pingGet().then(response => {
            }).catch(error => {
                if (error.response && error.status == 401 ) {
                    this.router.push({
                        name: 'login',
                        query: {
                            message: 'Your session has expired, please log in again'
                        },
                    })
                }
            })
        }

        if (this.authCheckInterval === null) {
            this.authCheckInterval = setInterval(() => {
                if(!["login", "oidc-callback"].includes(this.route.name)) {
                    this.apiClient.pingGet().then(response => {
                    }).catch(error => {
                        if (error.response && error.status == 401 ) {
                            this.router.push({
                                name: 'login',
                                query: {
                                    message: 'Your session has expired, please log in again'
                                },
                            })
                        }
                    })
                }
            }, 5000);
        }
    },
    beforeUnmount() {
        clearInterval(this.authCheckInterval)
        this.authCheckInterval = null
    }
}
</script>

<template>

</template>