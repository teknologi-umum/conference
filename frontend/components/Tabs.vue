<script lang="ts" setup>
const slots = useSlots()

const activeTab = ref(0)
console.log(slots.default!())
const tabItems = slots.default!()
</script>
<template>
    <div class="tabs">
        <div class="tab-header flex">
            <a href="#" 
            :class="`tab-title cursor-pointer text-gray-500 hover:text-white  p-3 ${activeTab == i ? 'active' : ''}`" 
            v-for="(item, i) in tabItems"
            @click.prevent="activeTab = i">
            
            {{ item.props?.['name'] }}
            </a>
            <div class="tab-closing flex-grow-1"></div>
        </div>
        <div class="tab-content">
            <component :is="item" v-for="(item, i) in tabItems" v-show="activeTab==i"></component>
        </div>
    </div>
</template>
<style>
.tab-title:not(.active) {
    border-bottom: 1px solid #ccc;
}
.tab-closing {
    border-bottom: 1px solid #ccc;
}
.tab-title:is(.active) {
    border: 1px solid #ccc;
    border-bottom: 0;
    border-top-left-radius: .3rem;
    border-top-right-radius: .3rem;
    color: white;
    background: #ffffff12;
}
.tab-content {
    border: 1px solid #ccc;
    min-height: 200px;
}

</style>