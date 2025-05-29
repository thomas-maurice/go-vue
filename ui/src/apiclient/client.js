import { API_BASE_URL } from "@/defaults/client"
import { ApiClient } from '@/gen/apiclient/src'
import { useUserStore } from "@/stores/user"

var store = useUserStore()
var client = new ApiClient(API_BASE_URL + "/api")
if (store.token !== "") {
    client.authentications["jwt"].apiKey = store.token
}
client.defaultHeaders = {}

function getDefaultClient() {
    var store = useUserStore()
    var client = new ApiClient(API_BASE_URL + "/api")
    if (store.token !== "") {
        client.authentications["jwt"].apiKey = store.token
    }
    client.defaultHeaders = {}

    return client
}

export const DefaultClient = client
export const GetDefaultClient = getDefaultClient
