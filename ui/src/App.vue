<script setup>

import { RouterView } from 'vue-router'

const html = document.querySelector("html")
html.setAttribute("data-bs-theme", "dark")

</script>

<template>
  <div class="container-fluid p-0">
    <RouterView />
  </div>
</template>

<style scoped>
</style>
