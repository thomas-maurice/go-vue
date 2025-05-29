import { defineStore } from 'pinia'
import { jwtDecode } from "jwt-decode";
import { faListNumeric } from '@fortawesome/free-solid-svg-icons';

export const useUserStore = defineStore('user', {
  persist: true,
  state: () => {
    return { 
        username: "<none>",
        token: "",
        logged_in: false,
        decoded: {
          admin: false,
          exp: -1,
        }
    }
  },
  actions: {
    logout() {
      this.logged_in = false
      this.username = "<none>"
      this.token = ""
      this.decoded = {
        admin: false,
        exp: -1,
      }
    },
    log_in(username, token) {
        this.logged_in = true
        this.username = username
        this.token = token
        this.decoded = jwtDecode(this.token)
    },
    logged_in() {
      if(this.decoded.exp && this.decoded.exp > Math.floor(Date.now() / 1000)) {
        return this.logged_in
      }

      return false
    }
  },
})