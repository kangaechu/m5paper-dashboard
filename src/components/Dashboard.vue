<template>
  <UpdatedAt></UpdatedAt>
  <div class="dashboard">
    <CardWeather> </CardWeather>
    <CardTemperature :value="firstItem.value"> </CardTemperature>
    <CardHumidity :value="secondItem.value"> </CardHumidity>
  </div>
</template>

<script lang="ts">
import UpdatedAt from "./UpdatedAt.vue";
import CardWeather from "./CardWeather.vue";
import CardTemperature from "./CardTemperature.vue";
import CardHumidity from "./CardHumidity.vue";
export default {
  components: {
    UpdatedAt,
    CardWeather,
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
