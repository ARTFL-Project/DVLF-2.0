<template>
    <div id="nearest-neighbors" class="shadow-sm" v-if="nearestNeighbors.length > 0">
        <b-card style="margin-top: 15px;">
            <h4 class="nym-title">Mots associés</h4>
            <span id="nn-title">Mots les plus associés à
                <span style="font-weight: 700">{{ headword }}</span> :</span>
            <p id="nn-content">
                <vue-word-cloud :words="words" :animation-overlap="0.5" :spacing="0.3" :font-size-ratio="0.3">
                    <template slot-scope="{text, weight, word}">
                        <div class="word-cloud" :title="weight" style="cursor: pointer;" @click="onWordClick(word)">
                            {{ text }}
                        </div>
                    </template>
                </vue-word-cloud>
            </p>
            <b-button-group size="sm" style="margin: 10px 0 10px 10px;" class="submit-btn" @click="exploreWord()">
                <b-button variant="primary">
                    <img style="height: 20px" src="../assets/images/baseline-search-24px.svg"/>
                    </b-button>
                    <b-button variant="primary">
                        Associations à travers le temps
                    </b-button>
            </b-button-group>
        </b-card>
    </div>
</template>

<script>
import VueWordCloud from "vuewordcloud"
import { EventBus } from "../main.js"

export default {
    name: "NearestNeighbors",
    components: {
        [VueWordCloud.name]: VueWordCloud
    },
    data() {
        return {
            words: []
        }
    },
    created() {
        for (let wordObject of this.nearestNeighbors) {
            this.words.push([wordObject.word, wordObject.distance])
        }
    },
    props: {
        nearestNeighbors: Array,
        headword: String
    },
    methods: {
        onWordClick(word) {
            this.$router.push(`/mot/${word[0]}`)
        },
        exploreWord() {
            EventBus.$emit("wordExplorer", this.headword)
        }
    }
}
</script>

<style scoped>
#nn-title {
    padding: 0px 10px 0px 10px;
    font-size: 80%;
}
.card-footer {
    padding: 10px;
    text-align: center;
}
#nn-content {
    width: 100%;
    height: 200px;
    padding: 0 5px 10px 5px;
    margin-bottom: 0;
}
.neighbor {
    display: inline-block;
    vertical-align: middle;
}
.neighbor span {
    display: inline-block;
    margin-left: -3px;
}
.word-cloud {
    color: rgb(21, 95, 131) !important;
    font-weight: 400;
    transition: all 200ms;
}
</style>

