<template>
  <div class="dashboard">
    <CardTemperature :value="firstItem.value"> </CardTemperature>
    <CardHumidity :value="secondItem.value"> </CardHumidity>
  </div>
</template>

<script lang="ts">
import CardTemperature from "./CardTemperature.vue";
import CardHumidity from "./CardHumidity.vue";
export default {
  components: {
    CardTemperature,
    CardHumidity,
  },
  data() {
    return {
      firstItem: {
        value: "25.3",
      },
      secondItem: {
        value: "82",
      },
    };
  },
};
</script>

<style scoped>
.dashboard {
  display: flex;
}
</style>
