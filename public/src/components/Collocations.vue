<template>
    <div id="collocations" v-if="collocates.length > 0">
        <div class="panel panel-default">
            <h4 class="nym-title">Collocations</h4>
            <span id="collocation-title">Mots utilisés à proximité de
                <b>{{ headword }}</b> :</span>
            <p id="collocation-content">
                <router-link class="collocate" :to="`/mot/${collocate.key}`" v-for="(collocate, index) in collocates" :key="index">
                    {{ collocate.key }}
                    <span v-if="index != collocates.length -1">,&nbsp;</span>
                </router-link>
            </p>
        </div>
    </div>
</template>

<script>
export default {
    name: "Collocations",
    props: {
        collocates: Array,
        headword: String
    }
}
</script>

<style scoped>
#collocation-title {
    padding: 0px 10px 15px 10px;
    font-size: 80%;
}
#collocation-content {
    padding: 0 10px 10px 10px;
}
.collocate {
    display: inline-block;
}
.collocate span {
    display: inline-block;
    margin-left: -3px;
}
</style>


