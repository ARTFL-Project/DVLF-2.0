<template>
    <div id="nearest-neighbors" class="shadow-sm" v-if="nearestNeighbors.length > 0">
        <b-card style="margin-top: 15px;">
            <h4 class="nym-title">Mots associés</h4>
            <span id="nn-title">Mots les plus associés à
                <b>{{ headword }}</b> :</span>
            <p id="nn-content">
                <router-link class="neighbor" :to="`/mot/${neighbor}`" v-for="(neighbor, index) in nearestNeighbors" :key="index" v-if="neighbor != ''">
                    {{ neighbor }}
                    <span v-if="index != nearestNeighbors.length -1">,&nbsp;</span>
                </router-link>
            </p>
        </b-card>
    </div>
</template>

<script>
export default {
    name: "NearestNeighbors",
    props: {
        nearestNeighbors: Array,
        headword: String
    }
}
</script>

<style scoped>
#nn-title {
    padding: 0px 10px 15px 10px;
    font-size: 80%;
}
#nn-content {
    padding: 0 10px 10px 10px;
}
.neighbor {
    display: inline-block;
}
.neighbor span {
    display: inline-block;
    margin-left: -3px;
}
</style>

