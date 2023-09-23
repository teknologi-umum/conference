<script lang="ts" setup>
const slots = useSlots()

const activeTab = ref(0)
console.log(slots.default!())
const tabItems = slots.default!()
</script>
<template>
    <div class="tabs">
        <div class="tab-header flex">
            <a 
            :class="`tab-title cursor-pointer text-gray-500 hover:text-white text-left  p-3 ${activeTab == i ? 'active' : ''}`" 
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
.tab-closing {
    border-bottom: 1px solid #cccccc3b;
}
.tab-title:is(.active) {
    border: 1px solid #cccccc3b;
    border-bottom: 0;
    border-top-left-radius: .3rem;
    border-top-right-radius: .3rem;
    color: white;
    background: #ffffff12;
}
.tab-content {
    border: 1px solid #cccccc7b;
    background-color: rgba(255, 255, 255, 0.035);
    min-height: 200px;
}

</style>