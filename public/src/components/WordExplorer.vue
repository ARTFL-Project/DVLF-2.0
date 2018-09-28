<template>
    <div>
        <b-row>
            <b-col>
                <b-card>
                    <div style="height: 100%; width:100%">
                        <vue-word-cloud :words="seventeenth" :animation-duration="0" :spacing="0.4" :font-size-ratio="0.5">
                            <template slot-scope="{text, weight, word}">
                                <div class="word-cloud" :title="weight" style="cursor: pointer;" @click="onWordClick(word)">
                                    {{ text }}
                                </div>
                            </template>
                        </vue-word-cloud>
                    </div>
                </b-card>
            </b-col>
            <b-col>
                <b-card>
                     <div style="height: 100%; width:100%">
                        <vue-word-cloud :words="eighteenth" :animation-duration="0" :spacing="0.4" :font-size-ratio="0.5">
                            <template slot-scope="{text, weight, word}">
                                <div class="word-cloud" :title="weight" style="cursor: pointer;" @click="onWordClick(word)">
                                    {{ text }}
                                </div>
                            </template>
                        </vue-word-cloud>
                    </div>
                </b-card>
            </b-col>
        </b-row>
        <b-row>
            <b-col>
                <b-card>
                     <div style="height: 100%; width:100%">
                        <vue-word-cloud :words="nineteenth" :animation-duration="0" :spacing="0.4" :font-size-ratio="0.5">
                            <template slot-scope="{text, weight, word}">
                                <div class="word-cloud" :title="weight" style="cursor: pointer;" @click="onWordClick(word)">
                                    {{ text }}
                                </div>
                            </template>
                        </vue-word-cloud>
                    </div>
                </b-card>
            </b-col>
            <b-col>
                <b-card>
                     <div style="height: 100%; width:100%">
                        <vue-word-cloud :words="twenteenth" :animation-duration="0" :spacing="0.4" :font-size-ratio="0.5">
                            <template slot-scope="{text, weight, word}">
                                <div class="word-cloud" :title="weight" style="cursor: pointer;" @click="onWordClick(word)">
                                    {{ text }}
                                </div>
                            </template>
                        </vue-word-cloud>
                    </div>
                </b-card>
            </b-col>
        </b-row>
    </div>
</template>

<script>
import { EventBus } from "../main.js"
import VueWordCloud from "vuewordcloud"

export default {
    name: "WordExplorer",
    components: {
        [VueWordCloud.name]: VueWordCloud
    },
    data() {
        return {
            seventeenth: [],
            eighteenth: [],
            nineteenth: [],
            twenteenth: []
        }
    },
    created() {
        let query = `${
                this.$globalConfig.apiServer
            }/api/explore/${this.$route.params.queryTerm.trim()}`
        this.$http
            .get(query, {
                headers: {
                    "Access-Control-Allow-Origin": "*",
                    "Content-Type": "application/json"
                }
            })
            .then(response => {
                console.log(response.data)
                this.seventeenth = this.convertToArray(response.data["1600"])
                console.log(this.seventeenth)
                this.eighteenth = this.convertToArray(response.data["1700"])
                this.nineteenth = this.convertToArray(response.data["1800"])
                this.twenteenth = this.convertToArray(response.data["1900"])
            })
            .catch(error => {
                this.error = error.toString()
                console.log(error)
            })

    },
    methods: {
        convertToArray(vectors) {
            let words = []
            for (let wordObject of vectors) {
                words.push([wordObject.word, wordObject.distance])
            }
            return words
        },
         onWordClick(word) {
            this.$router.push(`/mot/${word[0]}`)
        },
    }
}
</script>

<style>
</style>
