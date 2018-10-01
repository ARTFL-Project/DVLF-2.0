<template>
    <div id="nyms-boxes">
        <b-card class="shadow-sm">
            <b-row class="nym-title">
                <b-col variant="primary" @click="toggle('synonyms')" :class="{active: selected === 'synonyms'}">Synonymes</b-col>
                <b-col variant="primary" @click=" toggle('antonyms')" :class="{active: selected === 'antonyms'}">Antonymes</b-col>
            </b-row>
            <div id="synonyms" v-if="selected === 'synonyms'">
                <p class="synonyms" style="padding: 0 10px 10px 10px;">
                    <router-link class="synonym" :to="`/mot/${synonym.label}`" v-for="(synonym, index) in synonyms" :key="index">
                        {{ synonym.label }}
                        <span v-if="index != synonyms.length -1">,&nbsp;</span>
                    </router-link>
                </p>
                <p style="padding-left: 10px">
                    <b-button-group size="sm" class="submit-btn" @click="addSynonym()">
                        <b-button variant="primary">+</b-button>
                        <b-button variant="primary">
                            Ajoutez un synonyme
                        </b-button>
                    </b-button-group>
                </p>
            </div>
            <div id="antonyms" v-if="selected === 'antonyms'">
                <p class="antomyms" style="padding: 0 10px 10px 10px">
                    <router-link class="antonym" :to="`/mot/${antonym.label}`" v-for="(antonym, index) in antonyms" :key="index">
                        {{ antonym.label }}
                        <span v-if="index != antonyms.length -1">,&nbsp;</span>
                    </router-link>
                </p>
                <p style="padding-left: 10px">
                    <b-button-group size="sm" class="submit-btn" @click="addAntonym()">
                        <b-button variant="primary">+</b-button>
                        <b-button variant="primary">
                            Ajoutez un antonyme
                        </b-button>
                    </b-button-group>
                </p>
            </div>
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
            headword: this.$route.params.queryTerm,
            selected: "synonyms"
        }
    },
    methods: {
        toggle(selection) {
            this.selected = selection
        },
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
#synonyms,
#antonyms {
    margin-top: 10px;
}
.synonym,
.antonym {
    display: inline-block;
}

.synonym span,
.antonym span {
    display: inline-block;
    margin-left: -3px;
}
.nym-title {
    padding: 0;
    margin-left: 0;
    margin-right: 0;
    border-bottom-width: 0;
}
.nym-title .col {
    padding: 7px;
    cursor: pointer;
    transition: all 0.3s;
}
.nym-title .col:first-of-type {
    border-right: 1px #ddd solid;
}
.nym-title .col.active {
    background-color: #fff !important;
}
</style>

