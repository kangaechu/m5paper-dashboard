<template>
  <div class="item">
    <span class="title">{{ updatedAt }}</span>
  </div>
</template>

<script lang="ts">
export default {
  computed: {
    updatedAt: function (){
      const today = new Date();
      return today.getFullYear() + "/" + (today.getMonth() + 1) + "/" + today.getDay() + " "
        + today.getHours() + ":" + today.getMinutes();
    }
  }
};
</script>

<style scoped>
.item {
  text-align: center;
  font-size: small;
}
</style>
