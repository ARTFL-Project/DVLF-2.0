<template>
    <div id="nyms">
        <div class="panel panel-default">
            <h4 class="nym-title">Synonymes</h4>
            <p class="synonyms" style="padding: 0 10px 10px 10px;">
                <router-link class="synonym" :to="`/mot/${synonym.label}`" v-for="(synonym, index) in synonyms" :key="index">
                    {{ synonym.label }}
                    <span v-if="index != synonyms.length -1">,&nbsp;</span>
                </router-link>
            </p>
            <p style="padding-left: 10px">
                <button class="btn btn-default" type="button" @click="addSynonym()">
                    <span class="glyphicon glyphicon-plus add-button"></span>
                    <span class="hidden-xs hidden-sm">Ajoutez un </span>synonyme
                </button>
            </p>
        </div>
        <div class="panel panel-default" style="margin-top: 10px;">
            <h4 class="nym-title">Antonymes</h4>
            <p class="antomyms" style="padding: 0 10px 10px 10px">
                <router-link class="antonym" :to="`/mot/${antonym.label}`" v-for="(antonym, index) in antonyms" :key="index">
                    {{ antonym.label }}
                    <span v-if="index != antonyms.length -1">,&nbsp;</span>
                </router-link>
            </p>
            <p style="padding-left: 10px">
                <button class="btn btn-default" type="button" @click="addAntonym()">
                    <span class="glyphicon glyphicon-plus add-button"></span>
                    <span class="hidden-xs hidden-sm">Ajoutez un </span>antonyme
                </button>
            </p>
        </div>
    </div>

</template>

<script>
export default {
    name: "Search",
    props: {
        synonyms: Array,
        antonyms: Array
    },
    data() {
        return {
            headword: this.$route.params.queryTerm
        }
    },
    methods: {
        addSynonym() {
            this.$router.push(`/synonyme/${this.headword}`)
        },
        addAntonym() {
            this.$router.push(`/antonyme/${this.headword}`)
        }
    }
}
</script>

<style scoped>
.synonym,
.antonym {
    display: inline-block;
}

.synonym span,
.antonym span {
    display: inline-block;
    margin-left: -3px;
}
</style>

