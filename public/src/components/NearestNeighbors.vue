<template>
    <div id="nearest-neighbors" ng-if="nearestNeighbors.length > 0">
        <div class="panel panel-default">
            <h4 class="nym-title">Mots associés</h4>
            <span id="nn-title">Mots les plus associés à
                <b>{{ headword }}</b> :</span>
            <p id="nn-content">
                <router-link class="neighbor" :to="`/mot/${neighbor}`" v-for="(neighbor, index) in nearestNeighbors" :key="index">
                    {{ neighbor }}
                    <span v-if="index != nearestNeighbors.length -1">,&nbsp;</span>
                </router-link>
            </p>
        </div>
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

