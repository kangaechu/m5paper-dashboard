<template>
  <div class="container">
    <div class="item">
      <span>{{ $props.icon }}</span
      ><br />
      <span class="title">{{ $props.title }}</span>
      <br />
      <span class="value">{{ $props.value }}</span>
    </div>
  </div>
</template>
  
<script lang="ts">
export default {
  props: {
    icon: {
      type: String,
    },
    title: {
      type: String,
    },
    value: {
      type: String,
    },
  },
};
</script>
 
<style scoped>
div.item {
  padding: 0.5em 1em;
  margin: 1em 0;
  background: white;
}
span.title {
  font-weight: 500;
}
</style>
  