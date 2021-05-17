<template>
  <div class="item">
    <span>
      <!--?xml version="1.0" encoding="utf-8"?-->
      <!-- Generator: Adobe Illustrator 15.1.0, SVG Export Plug-In . SVG Version: 6.00 Build 0)  -->

      <svg
        id="_x32_"
        xmlns="http://www.w3.org/2000/svg"
        x="0px"
        y="0px"
        width="40"
        height="40"
        viewBox="0 0 512 512"
        xml:space="preserve"
      >
        <g>
          <path
            class="st0"
            d="M388.823,191.344L255.995,0L123.183,191.344c-53.875,84.438-73.359,192.281,0,265.641
		C159.854,493.656,207.917,512,255.995,512s96.141-18.344,132.828-55.016C462.167,383.625,442.698,275.781,388.823,191.344z
		 M237.339,442.609c-2.359,8.125-10.875,12.797-19,10.438c-21.422-6.234-41.109-17.828-57.25-33.984
		c-11.844-11.844-20.125-25.172-25.25-38.953c-5.156-13.797-7.203-27.969-7.203-41.766c0.016-21.266,4.828-41.672,11.109-59.688
		c2.797-8,11.531-12.203,19.516-9.406c8,2.781,12.203,11.531,9.406,19.516c-5.531,15.734-9.406,33.016-9.391,49.578
		c0,10.781,1.609,21.219,5.266,31.063c3.703,9.844,9.422,19.172,18.234,28c12.531,12.547,27.625,21.406,44.125,26.219
		C235.026,425.984,239.683,434.484,237.339,442.609z"
            style="fill: rgb(75, 75, 75)"
          ></path>
        </g>
      </svg>
    </span>
    <br />
    <span class="title">湿度</span>
    <br />
    <span class="value">{{ $props.value }}</span>
    <span class="unit">%</span>
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
    unit: {
      type: String,
    },
  },
};
</script>

<style scoped>
div.item {
  text-align: center;
  padding: 0.5em 1em;
  margin: 1em 0;
  background: white;
}
span.title {
  font-weight: 500;
}
span.value {
  font-size: 2.5rem;
}
</style>
