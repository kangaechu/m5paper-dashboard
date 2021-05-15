<template>
  <div class="dashboard">{{ currentDate }}</div>
</template>
  
  <script lang="ts">
import { ref, defineComponent } from "vue";
export default defineComponent({
  name: "Dashboard",
  computed: {
    currentDate() {
      const d = new Date();
      return (
        d.getFullYear() +
        "/" +
        (d.getMonth() + 1) +
        "/" +
        d.getDate() +
        " " +
        d.getHours() +
        ":" +
        d.getMinutes() +
        ":" +
        d.getSeconds()
      );
    },
  },
});
</script>
 
<style scoped>
div.dashboard {
  font-size: 10em;
}
</style>
  