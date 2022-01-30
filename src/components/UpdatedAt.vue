<template>
  <div class="item">
    <span class="title">{{ updatedAt }}</span>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
export default defineComponent ({
  computed: {
    updatedAt: function (){
      const today = new Date();
      return today.getFullYear() + "/" + (today.getMonth() + 1) + "/" + today.getDay() + " "
        + today.getHours() + ":" + ("0" + today.getMinutes()).slice(-2);
    }
  }
});
</script>

<style scoped>
.item {
  text-align: center;
  font-size: small;
}
</style>
