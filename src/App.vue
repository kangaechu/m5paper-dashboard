<template>
  <div>
    <Dashboard />
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Dashboard from "./components/Dashboard.vue";

export default defineComponent({
  name: "App",
  components: {
    Dashboard,
  },
});
</script>

<style>
html {
  font-family: "Helvetica Neue", "Helvetica", "Hiragino Sans",
    "Hiragino Kaku Gothic ProN", "Arial", "Yu Gothic", "Meiryo", sans-serif;
  font-size: 1.5em;
}
</style>
