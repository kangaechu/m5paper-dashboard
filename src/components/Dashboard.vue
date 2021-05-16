<template>
  <div class="dashboard">
    <ItemThermostat :value="firstItem.value"> </ItemThermostat>
    <ItemHumidity :value="secondItem.value"> </ItemHumidity>
  </div>
</template>
  
<script lang="ts">
import ItemThermostat from "./ItemThermostat.vue";
import ItemHumidity from "./ItemHumidity.vue";
export default {
  components: {
    ItemThermostat,
    ItemHumidity,
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
  