<template>
  <div class="dashboard">
    <ItemWithIcon
      :icon="firstItem.icon"
      :title="firstItem.title"
      :value="firstItem.value"
    >
    </ItemWithIcon>
    <ItemWithIcon
      :icon="secondItem.icon"
      :title="secondItem.title"
      :value="secondItem.value"
    >
    </ItemWithIcon>
  </div>
</template>
  
<script lang="ts">
import ItemWithIcon from "./ItemWithIcon.vue";
export default {
  components: {
    ItemWithIcon,
  },
  data() {
    return {
      firstItem: {
        icon: "🌡",
        title: "気温",
        value: "25.3℃",
      },
      secondItem: {
        icon: "🌡",
        title: "気温",
        value: "25.3℃",
      },
    };
  },
};
</script>
 
<style scoped>
</style>
  