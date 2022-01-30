<template>
  <div class="item">
    <span class="icon"><img :src="icon"></span>
    <br/>
    <span class="title">{{ weather }}</span>
    <br/>
    <span class="value">{{ minTemp }}-{{ maxTemp }}</span>
    <span class="unit">℃</span>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
export default defineComponent({
  data() {
    return {
      weather: '晴れ一時雨',
      minTemp: 15,
      maxTemp: 20,
      rainyPercent: 10,
      icon:'assets/img/sunny.svg',
    }
  },
  computed: {
  }
});
</script>

<style scoped>
div.item {
  text-align: center;
  padding: 0.5em 1em;
  margin: 0 0 1em 0;
  background: white;
}

span.title {
  font-weight: 500;
}

span.value {
  font-size: 2.5rem;
}
</style>
