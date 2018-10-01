<template>
    <div id="word-explorer" class="shadow-lg">
        <div id="explorer-title">
            <div id="title-close" @click="close">
                &times;
            </div>
            Exploration des usages de <span style="font-variant: small-caps">{{ headword }}</span> à travers le temps
        </div>
        <div id="explorer-body">
            <b-row align-h="center">
                <b-col cols="6" v-if="seventeenth" style="margin-bottom: 20px;">
                    <b-card style="height: 400px;" header="Entre 1600 et 1700">
                        <vue-word-cloud :words="seventeenth" :animation-overlap="0.2" :spacing="0.4" :font-size-ratio="0.3">
                            <template slot-scope="{text, weight, word}">
                                <div class="word-cloud" :title="weight" style="cursor: pointer;" @click="onWordClick(word)">
                                    {{ text }}
                                </div>
                            </template>
                        </vue-word-cloud>
                    </b-card>
                </b-col>
                <b-col cols="6" v-if="eighteenth" style="margin-bottom: 20px;">
                    <b-card style="height: 400px;" header="Entre 1700 et 1800">
                        <vue-word-cloud v-if="eighteenth" :words="eighteenth" :animation-overlap="0.2" :spacing="0.4" :font-size-ratio="0.3">
                            <template slot-scope="{text, weight, word}">
                                <div class="word-cloud" :title="weight" style="cursor: pointer;" @click="onWordClick(word)">
                                    {{ text }}
                                </div>
                            </template>
                        </vue-word-cloud>
                    </b-card>
                </b-col>
                <b-col cols="6" v-if="nineteenth" style="margin-bottom: 20px;">
                    <b-card style="height: 400px; min-width:40%" header="Entre 1800 et 1900">
                        <vue-word-cloud :words="nineteenth" :animation-overlap="0.2" :spacing="0.4" :font-size-ratio="0.3">
                            <template slot-scope="{text, weight, word}">
                                <div class="word-cloud" :title="weight" style="cursor: pointer;" @click="onWordClick(word)">
                                    {{ text }}
                                </div>
                            </template>
                        </vue-word-cloud>
                    </b-card>
                </b-col>
                <b-col cols="6" v-if="twenteenth">
                    <b-card style="height: 400px; min-width: 40%" header="Entre 1900 et 2000">
                        <vue-word-cloud :words="twenteenth" :animation-overlap="0.2" :spacing="0.4" :font-size-ratio="0.3">
                            <template slot-scope="{text, weight, word}">
                                <div class="word-cloud" :title="weight" style="cursor: pointer;" @click="onWordClick(word)">
                                    {{ text }}
                                </div>
                            </template>
                        </vue-word-cloud>
                    </b-card>
                </b-col>
            </b-row>
            <b-button style="margin:20px 0" variant="primary" @click="close()">
                Fermer
            </b-button>
        </div>

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
    props: {
        vectors: Object,
        headword: String
    },
    data() {
        return {
            seventeenth: null,
            eighteenth: null,
            nineteenth: null,
            twenteenth: null
        }
    },
    created() {
        window.scrollTo({ top: 0, left: 0, behavior: "smooth" })
        if ("1600" in this.vectors) {
            this.seventeenth = this.convertToArray(this.vectors["1600"])
        }
        if ("1700" in this.vectors) {
            this.eighteenth = this.convertToArray(this.vectors["1700"])
        }
        if ("1800" in this.vectors) {
            this.nineteenth = this.convertToArray(this.vectors["1800"])
        }
        if ("1900" in this.vectors) {
            this.twenteenth = this.convertToArray(this.vectors["1900"])
        }
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
            this.close()
        },
        close() {
            EventBus.$emit("closeWordExplorer")
        }
    }
}
</script>

<style scoped>
#word-explorer {
    position: absolute;
    z-index: 100;
    width: 98%;
    left: 0;
    right: 0;
    margin: auto;
    background-color: #fff;
    border: 1px solid rgba(0, 0, 0, 0.125);
    border-radius: 0.25rem;
}
#explorer-title {
    position: relative;
    margin-bottom: 20px;
    background-color: #f0f0f0;
    border-bottom: 1px solid #eee;
    text-align: center;
    margin-top: 0;
    padding: 7px;
    font-weight: 700;
    font-size: 130%;
}
#title-close {
    position: absolute;
    right: 0;
    top: 0;
    background-color: rgba(21, 95, 131, 0.8) !important;
    color: #fff !important;
    padding: 0px 5px;
    cursor: pointer;
    font-size: 100%;
}
#explorer-body {
    padding: 15px;
}
.card-header {
    text-align: center;
    font-weight: 700;
}
.card-body {
    height: 90%;
    width: 100%;
    padding: 10px !important;
}
.word-cloud {
    color: rgb(21, 95, 131) !important;
    font-weight: 400;
    transition: all 200ms;
}
.btn-primary {
    background-color: rgba(21, 95, 131, 0.8) !important;
    color: #fff !important;
    float: right;
}
</style>
