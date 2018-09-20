<template>
    <div id="nyms-boxes">
        <b-card class="shadow-sm">
            <h4 class="nym-title">Synonymes</h4>
            <p class="synonyms" style="padding: 0 10px 10px 10px;">
                <router-link class="synonym" :to="`/mot/${synonym.label}`" v-for="(synonym, index) in synonyms" :key="index">
                    {{ synonym.label }}
                    <span v-if="index != synonyms.length -1">,&nbsp;</span>
                </router-link>
            </p>
            <p style="padding-left: 10px">
                <b-button variant="primary" @click="addSynonym()">
                    <span class="hidden-xs hidden-sm">Ajoutez un </span>synonyme
                </b-button>
            </p>
        </b-card>
        <b-card class="shadow-sm" style="margin-top: 15px;">
            <h4 class="nym-title">Antonymes</h4>
            <p class="antomyms" style="padding: 0 10px 10px 10px">
                <router-link class="antonym" :to="`/mot/${antonym.label}`" v-for="(antonym, index) in antonyms" :key="index">
                    {{ antonym.label }}
                    <span v-if="index != antonyms.length -1">,&nbsp;</span>
                </router-link>
            </p>
            <p style="padding-left: 10px">
                <b-button variant="primary" @click="addAntonym()">
                    <span class="hidden-xs hidden-sm">Ajoutez un </span>antonyme
                </b-button>
            </p>
        </b-card>
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

