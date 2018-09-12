<template>
    <div class="panel panel-default" style="padding: 10px; margin-top: 15px;">
        <div class="row" v-if="!atHome">
            <div class="col-xs-12 col-sm-offset-3 col-sm-6">
                <h3 :class="{'hide-results': apropos}" style="text-align: center;">
                    <b>{{ currentTerm }}</b>:
                    <span class="entry-total" v-if="totalResults > 0">
                        {{ totalResults }} {{ pluralize("entrée", totalResults) }} dans {{ totalDicos }} {{ pluralize("dictionnaire", totalDicos) }}
                    </span>
                </h3>
            </div>
            <div class="col-xs-12 col-sm-3" style="margin-top: 15px">
                <div class="pull-right hidden-xs">
                    <div class="btn-group">
                        <button class="btn btn-default dropdown-toggle" type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                            <span class="glyphicon glyphicon-plus add-button"></span>
                            <span style="font-variant: small-caps;">Contribuer au DVLF</span>
                            <span class="caret"></span>
                        </button>
                        <ul class="dropdown-menu">
                            <li>
                                <router-link to="/definition" style="font-variant: small-caps; font-weight: 700">Ajouter une définition</router-link>
                            </li>
                            <li>
                                <router-link to="/exemple" style="font-variant: small-caps; font-weight: 700">Ajouter un exemple</router-link>
                            </li>
                            <li>
                                <router-link to="/synonyme" style="font-variant: small-caps; font-weight: 700">Ajouter un synonyme</router-link>
                            </li>
                            <li>
                                <router-link to="/antonyme" style="font-variant: small-caps; font-weight: 700">Ajouter un antonyme</router-link>
                            </li>
                        </ul>
                    </div>
                </div>
                <div class="hidden-sm hidden-md hidden-lg" style="text-align: center">
                    <div class="btn-group">
                        <button class="btn btn-default dropdown-toggle" type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                            <span class="glyphicon glyphicon-plus add-button"></span>
                            <span style="font-variant: small-caps; font-weight: 700">Contribuer au DVLF</span>
                            <span class="caret"></span>
                        </button>
                        <ul class="dropdown-menu">
                            <li>
                                <router-link to="/definition" style="font-variant: small-caps; font-weight: 700">Ajouter une définition</router-link>
                            </li>
                            <li>
                                <router-link to="/exemple" style="font-variant: small-caps; font-weight: 700">Ajouter un exemple</router-link>
                            </li>
                            <li>
                                <router-link to="/synonyme" style="font-variant: small-caps; font-weight: 700">Ajouter un synonyme</router-link>
                            </li>
                            <li>
                                <router-link to="/antonyme" style="font-variant: small-caps; font-weight: 700">Ajouter un antonyme</router-link>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
        <div id="results" class="row">
            <div class="hidden-xs col-sm-3 col-md-2" style="margin-top: 15px;">
                <word-wheel :headword="currentTerm"></word-wheel>
            </div>
            <div class="col-xs-12 col-sm-6 col-md-7 col-lg-7" v-if="totalResults > 0">
                <dictionary-entries :results="results.dictionaries" :fuzzy-results="results.fuzzyResults"></dictionary-entries>
                <syn-anto-nyms class="hidden-sm hidden-md hidden-lg" :synonyms="results.synonyms" :antonyms="results.antonyms"></syn-anto-nyms>
                <examples :examples="results.examples"></examples>
            </div>
            <div class="hidden-xs col-sm-3 col-md-3 col-lg-3" style="margin-top: 15px;" v-if="totalResults > 0">
                <syn-anto-nyms :synonyms="results.synonyms" :antonyms="results.antonyms"></syn-anto-nyms>
                <collocations :collocates="results.collocates" :headword="currentTerm"></collocations>
                <nearest-neighbors :nearest-neighbors="results.nearestNeighbors" :headword="currentTerm"></nearest-neighbors>
                <time-series :time-series="results.timeSeries" :headword="currentTerm"></time-series>
            </div>
        </div>
    </div>
</template>

<script>
import DictionaryEntries from "./DictionaryEntries.vue"
import SynAntoNyms from "./SynAntoNyms.vue"
import Collocations from "./Collocations.vue"
import NearestNeighbors from "./NearestNeighbors.vue"
import TimeSeries from "./TimeSeries.vue"
import Examples from "./Examples.vue"
import WordWheel from "./WordWheel.vue"
import { EventBus } from "../main.js"

export default {
    name: "Results",
    components: {
        DictionaryEntries,
        SynAntoNyms,
        Collocations,
        NearestNeighbors,
        TimeSeries,
        Examples,
        WordWheel
    },
    data: function() {
        return {
            currentTerm: this.$route.params.queryTerm,
            results: {},
            totalResults: 0,
            totalDicos: 0,
            atHome: true,
            apropos: false
        }
    },
    created() {
        this.fetchData()
    },
    methods: {
        fetchData() {
            let query = `${
                this.$globalConfig.apiServer
            }/api/mot/${this.$route.params.queryTerm.trim()}`
            this.$http
                .get(query, {
                    headers: {
                        "Access-Control-Allow-Origin": "*",
                        "Content-Type": "application/json"
                    }
                })
                .then(response => {
                    this.atHome = false
                    EventBus.$emit("OffHome")
                    this.results = response.data
                    this.totalResults = this.results.dictionaries.totalEntries
                    this.totalDicos = this.results.dictionaries.totalDicos
                })
                .catch(error => {
                    this.error = error.toString()
                    console.log(error)
                })
        },
        pluralize(word, count) {
            if (count > 1) {
                return word + "s"
            }
            return word
        }
    }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.hide-results {
    visibility: hidden;
}
</style>